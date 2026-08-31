package handler

// 应用层 host 拦截:仅在用户启用“关闭公网访问”且 HTTPS 启用后,只允许 Host header
// 匹配 master_url 域名的请求
// (或来自 loopback 的本机自调用),其它临时重定向到正确的 https URL。
//
// 历史背景:
//   早期实现是 master 启动时绑 127.0.0.1,物理层禁止 IP+端口直连。但这套实现
//   在 Docker 部署下会让 -p 端口映射失效(容器内的 lo 跟 host 端口隔绝),所以
//   把"防 IP 直连"挪到应用层,跟 bind host 解耦:
//     - master 总是绑 0.0.0.0(网络层尽量通用)
//     - 应用层根据 master_url 配置决定拒不拒
//
// 设计要点:
//   - 仅当 master_local_only=1 且 master_url 以 https:// 开头时才生效
//   - 缓存 5min(读 system_settings 不便每请求都查 db);设置保存后主动刷新
//   - loopback 来源永远放行 — 容器健康检查 / 同主机 nginx 反代 / master 自调用需要
//   - Host 不匹配时 307 重定向到正确 https URL — 设置关闭后不会留下永久缓存

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"miaomiaowux/internal/storage"
)

// hostEnforcer 中间件状态。
type HTTPSHostEnforcer struct {
	repo *storage.TrafficRepository

	mu                     sync.RWMutex
	cachedMasterURL        string
	cachedHost             string // master_url 的 host(不含端口)
	cachedSubscriptionHost string // subscription_url 的 host,仅订阅请求放行
	cachedHTTPS            bool
	cachedLocalOnly        bool
	lastRefresh            time.Time
	// pendingRefresh 用于避免雷暴 — 多请求同时发现缓存过期时只有一个去查 db
	pendingRefresh atomic.Bool
}

const masterURLCacheTTL = 5 * time.Minute

// EnforceHTTPSHost 包一层中间件:HTTPS 启用时拒绝非合法 Host 的请求(loopback 除外)。
// 通常用法:srv.Handler = EnforceHTTPSHost(mux, repo)。
func EnforceHTTPSHost(next http.Handler, repo *storage.TrafficRepository) http.Handler {
	return NewHTTPSHostEnforcer(repo).Middleware(next)
}

func NewHTTPSHostEnforcer(repo *storage.TrafficRepository) *HTTPSHostEnforcer {
	return &HTTPSHostEnforcer{repo: repo}
}

func (e *HTTPSHostEnforcer) Middleware(next http.Handler) http.Handler {
	if e == nil || e.repo == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e.shouldBlock(r) {
			// 使用 307 而非可被浏览器长期缓存的 308。访问限制是可切换设置，
			// 永久重定向会导致关闭限制后 IP:端口仍被客户端继续跳转。
			target := e.canonicalURL() + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Refresh 让系统设置保存后立即生效，避免旧状态在五分钟缓存期内继续产生重定向。
func (e *HTTPSHostEnforcer) Refresh() {
	if e != nil && e.repo != nil {
		e.refresh()
	}
}

// ForcePublicAccessEnabled 是因错误访问设置而锁在面板外时的启动逃生开关。
// 设置 MMWX_FORCE_PUBLIC_ACCESS=1 并重启，可恢复公网监听并绕过 Host 重定向。
func ForcePublicAccessEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("MMWX_FORCE_PUBLIC_ACCESS"))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// shouldBlock 返回 true 表示该请求被拦截 + 重定向。
func (e *HTTPSHostEnforcer) shouldBlock(r *http.Request) bool {
	host, subscriptionHost, https, localOnly := e.snapshot()
	if ForcePublicAccessEnabled() || !localOnly || !https || host == "" {
		return false
	}
	// loopback 来源放行
	if remoteIPIsLoopback(r.RemoteAddr) {
		return false
	}
	// Host 比对(忽略端口、忽略大小写)
	got := stripPort(r.Host)
	if strings.EqualFold(got, host) {
		return false
	}
	if subscriptionHost != "" && strings.EqualFold(got, subscriptionHost) && isSubscriptionPath(r.URL.Path) {
		return false
	}
	return true
}

func isSubscriptionPath(path string) bool {
	return strings.HasPrefix(path, "/x/") ||
		path == "/api/clash/subscribe" ||
		path == "/api/user/package-subscribe" ||
		path == "/api/subscribe"
}

// snapshot 返回缓存的 (期望 host, 是否启用 https)。过期就异步刷新 + 同步用旧值
// (启动后第一次请求会同步拉一次以保证立刻可用)。
func (e *HTTPSHostEnforcer) snapshot() (string, string, bool, bool) {
	e.mu.RLock()
	expired := time.Since(e.lastRefresh) > masterURLCacheTTL
	host, subscriptionHost, https, localOnly := e.cachedHost, e.cachedSubscriptionHost, e.cachedHTTPS, e.cachedLocalOnly
	firstTime := e.lastRefresh.IsZero()
	e.mu.RUnlock()

	if firstTime {
		// 同步刷一次,避免启动后短暂"全放行"
		e.refresh()
		e.mu.RLock()
		host, subscriptionHost, https, localOnly = e.cachedHost, e.cachedSubscriptionHost, e.cachedHTTPS, e.cachedLocalOnly
		e.mu.RUnlock()
		return host, subscriptionHost, https, localOnly
	}
	if expired && e.pendingRefresh.CompareAndSwap(false, true) {
		go func() {
			defer e.pendingRefresh.Store(false)
			e.refresh()
		}()
	}
	return host, subscriptionHost, https, localOnly
}

func (e *HTTPSHostEnforcer) refresh() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	raw, err := e.repo.GetSystemSetting(ctx, "master_url")
	if err != nil {
		// 读不到先按"未启用"处理,5 分钟后下次再试
		e.mu.Lock()
		e.cachedMasterURL = ""
		e.cachedHost = ""
		e.cachedSubscriptionHost = ""
		e.cachedHTTPS = false
		e.lastRefresh = time.Now()
		e.mu.Unlock()
		return
	}
	masterURL := strings.TrimSpace(raw)
	subscriptionURL, _ := e.repo.GetSystemSetting(ctx, "subscription_url")
	localOnlyRaw, _ := e.repo.GetSystemSetting(ctx, "master_local_only")
	localOnly := localOnlyRaw == "1"
	https := strings.HasPrefix(masterURL, "https://")
	var host string
	if https {
		if u, perr := url.Parse(masterURL); perr == nil {
			host = stripPort(u.Host)
		}
	}
	var subscriptionHost string
	if u, perr := url.Parse(strings.TrimSpace(subscriptionURL)); perr == nil {
		subscriptionHost = stripPort(u.Host)
	}
	e.mu.Lock()
	if e.cachedMasterURL != masterURL {
		log.Printf("[HostEnforcement] master_url=%q https=%v host=%q", masterURL, https, host)
	}
	e.cachedMasterURL = masterURL
	e.cachedHost = host
	e.cachedSubscriptionHost = subscriptionHost
	e.cachedHTTPS = https
	e.cachedLocalOnly = localOnly
	e.lastRefresh = time.Now()
	e.mu.Unlock()
}

// canonicalURL 返回正确的 master_url 前缀(无尾斜杠),供 redirect 用。
func (e *HTTPSHostEnforcer) canonicalURL() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return strings.TrimRight(e.cachedMasterURL, "/")
}

// stripPort 从 "host:port" 或 "[::1]:port" 中取 host 部分,失败时原样返回。
func stripPort(s string) string {
	if s == "" {
		return s
	}
	if h, _, err := net.SplitHostPort(s); err == nil {
		return h
	}
	return s
}

// remoteIPIsLoopback 解析 r.RemoteAddr 的 IP 部分,判断是否 loopback (127.0.0.0/8 / ::1)。
// nginx 反代同主机时 RemoteAddr 是 127.0.0.1;Docker 容器内健康检查也是 loopback。
func remoteIPIsLoopback(remoteAddr string) bool {
	if remoteAddr == "" {
		return false
	}
	ipStr := stripPort(remoteAddr)
	ipStr = strings.Trim(ipStr, "[]")
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}
