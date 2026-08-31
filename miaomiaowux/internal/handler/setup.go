package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/bcrypt"

	"miaomiaowux/internal/logger"
	"miaomiaowux/internal/storage"
	"miaomiaowux/templates"
)

// setupCompleted 缓存「已完成初始化」的结论。一旦确认库里有用户就永久置位，
// 后续请求不再查库 —— 根路径是最高频入口，不能每次都打一次数据库。
var setupCompleted atomic.Bool

// SetupPending 报告系统是否尚未初始化（库里一个用户都没有）。
//
// 用途：首次安装时让根路径直接进注册页，而不是对外伪装的探针页 —— 全新部署的人
// 打开面板本该看到「创建管理员账号」，而不是一个看不出所以然的探针。
// 注册完成后本函数恒为 false，探针伪装照常生效。
//
// 查询出错时返回 false（当作已初始化）：数据库抖动不应该把注册入口暴露出去。
func SetupPending(ctx context.Context, repo *storage.TrafficRepository) bool {
	if repo == nil || setupCompleted.Load() {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	users, err := repo.ListUsers(ctx, 1)
	if err != nil {
		return false
	}
	if len(users) > 0 {
		setupCompleted.Store(true)
		return false
	}
	return true
}

// InsecureBrowserContext 判断浏览器是否会因「非安全上下文」而禁用 WebCrypto。
//
// 判定规则与浏览器一致：HTTPS 算安全；http://localhost、127.0.0.0/8、[::1] 也算安全；
// 其余纯 HTTP（例如 http://1.2.3.4:12889）一律不安全。反代场景下 r.TLS 为 nil，
// 所以同时认 X-Forwarded-Proto。
//
// 为什么需要它：面板前端的加密通道（securechan）依赖 crypto.subtle，非安全上下文下
// 浏览器直接不提供该 API，前端所有 API 调用都会失败；首次安装时「要不要注册」那次查询
// 也失败，注册页就静默退化成登录框 —— 用户对着一个没反应的登录框完全无从下手。
func InsecureBrowserContext(r *http.Request) bool {
	if r == nil {
		return false
	}
	if r.TLS != nil {
		return false
	}
	proto := strings.ToLower(strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")))
	if i := strings.IndexByte(proto, ','); i >= 0 {
		proto = strings.TrimSpace(proto[:i])
	}
	if proto == "https" {
		return false
	}
	host := r.Host
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	host = strings.Trim(strings.ToLower(strings.TrimSpace(host)), "[]")
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return false
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		return false
	}
	return true
}

// WriteInsecureContextNotice 渲染一张纯 HTML 说明页（不依赖任何 JS）。
// 只在「首次安装 + 非安全上下文」时出现：此时注册根本走不下去，与其让用户对着
// 一个静默失效的登录框，不如直接告诉他两条可行路径。
func WriteInsecureContextNotice(w http.ResponseWriter, r *http.Request) {
	hostPort := r.Host
	port := "12889"
	if _, p, err := net.SplitHostPort(hostPort); err == nil && p != "" {
		port = p
	}
	esc := html.EscapeString
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, insecureContextNoticeHTML, esc(hostPort), esc(port), esc(port), esc(port))
}

const insecureContextNoticeHTML = `<!doctype html>
<html lang="zh"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>首次安装 · 需要安全连接</title>
<style>
 :root{color-scheme:light dark}
 body{margin:0;min-height:100vh;display:grid;place-items:center;padding:24px;
      background:#f7f8fa;color:#1f2430;
      font:15px/1.7 -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"PingFang SC","Microsoft YaHei",sans-serif}
 .card{max-width:660px;background:#fff;border:1px solid #e6e8ec;border-radius:12px;padding:32px 30px}
 h1{margin:0 0 6px;font-size:19px}
 .sub{margin:0 0 22px;color:#6b7280;font-size:14px}
 h2{margin:24px 0 8px;font-size:15px}
 p{margin:8px 0}
 code,pre{font-family:ui-monospace,Menlo,Consolas,monospace}
 pre{background:#0f1520;color:#e6eaf1;border-radius:8px;padding:12px 14px;
     font-size:13px;overflow-x:auto;white-space:pre-wrap;word-break:break-all}
 .muted{color:#6b7280;font-size:13px}
 .ok{color:#12924a;font-weight:600}
 @media (prefers-color-scheme:dark){
   body{background:#14161a;color:#e6e8ec}
   .card{background:#1b1e24;border-color:#2a2f37}
   .muted,.sub{color:#9aa1ab}}
</style></head><body><div class="card">
<h1>首次安装：请用 HTTPS 或 localhost 打开</h1>
<p class="sub">当前地址 <code>http://%s</code> 是纯 HTTP，浏览器判定为「不安全上下文」。</p>
<p>面板的加密通道需要浏览器的 <code>WebCrypto</code>（<code>crypto.subtle</code>），
而浏览器只在 <b>HTTPS</b> 或 <b>localhost</b> 下才提供它。现在它不可用，
前端所有接口调用都会失败 —— 包括「是否需要注册」那次查询，
于是注册页会静默退化成一个<b>点了没反应的登录框</b>。</p>
<p class="ok">服务端本身是正常的，数据也没问题，只是浏览器不让在这个地址上做加密。</p>
<h2>办法一：用 SSH 隧道（最快，不用配域名）</h2>
<p>在你自己的电脑上执行（把 <code>服务器IP</code> 换成实际地址）：</p>
<pre>ssh -L %s:127.0.0.1:%s root@服务器IP</pre>
<p>保持这个窗口不关，然后在本机浏览器打开：</p>
<pre>http://localhost:%s</pre>
<p class="muted">localhost 属于安全上下文，注册页就能正常出来。注册完再按办法二配好 HTTPS 即可。</p>
<h2>办法二：配 HTTPS（长期方案）</h2>
<p>给面板配一个域名并启用证书，例如用 Nginx / Caddy 反代到本机端口，
或者套 Cloudflare。之后用 <code>https://你的域名</code> 访问即可。</p>
<p class="muted">反代时记得把 <code>X-Forwarded-Proto: https</code> 传给后端，本页会自动消失。</p>
</div></body></html>`

type setupStatusResponse struct {
	NeedsSetup bool `json:"needs_setup"`
}

type setupRequest struct {
	Username  string                  `json:"username"`
	Password  string                  `json:"password"`
	Nickname  string                  `json:"nickname"`
	Email     string                  `json:"email"`
	AvatarURL string                  `json:"avatar_url"`
	Domain    string                  `json:"domain"`
	Database  *storage.DatabaseConfig `json:"database,omitempty"`
}

type setupResponse struct {
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Email       string `json:"email"`
	NginxSetup  bool   `json:"nginx_setup,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
	Restarting  bool   `json:"restarting,omitempty"`
}

// 返回一个处理程序，用于检查是否需要初始设置
func NewSetupStatusHandler(repo *storage.TrafficRepository) http.Handler {
	if repo == nil {
		panic("setup status handler requires repository")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("[初始化检查] 收到初始化状态检查请求",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
		)

		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, errors.New("only GET is supported"))
			return
		}

		users, err := repo.ListUsers(r.Context(), 10)
		if err != nil {
			logger.Error("[初始化检查] 查询用户列表失败", "error", err)
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		needsSetup := len(users) == 0

		if len(users) > 0 {
			usernames := make([]string, len(users))
			for i, u := range users {
				usernames[i] = u.Username
			}
			logger.Info("[初始化检查] 数据库中已存在用户",
				"user_count", len(users),
				"usernames", usernames,
				"needs_setup", needsSetup,
			)
		} else {
			logger.Info("[初始化检查] 数据库中没有用户，需要初始化",
				"user_count", 0,
				"needs_setup", needsSetup,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(setupStatusResponse{NeedsSetup: needsSetup})
	})
}

// 处理第一个管理员用户的创建
func NewInitialSetupHandler(repo *storage.TrafficRepository, dataDir ...string) http.Handler {
	if repo == nil {
		panic("initial setup handler requires repository")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("[初始化] 收到初始化请求",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
		)

		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, errors.New("only POST is supported"))
			return
		}

		// 检查是否还需要设置
		users, err := repo.ListUsers(r.Context(), 1)
		if err != nil {
			logger.Error("[初始化] 查询用户列表失败", "error", err)
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		if len(users) > 0 {
			logger.Warn("[初始化] 系统已初始化，拒绝重复初始化",
				"existing_user_count", len(users),
				"first_user", users[0].Username,
			)
			writeError(w, http.StatusConflict, errors.New("系统已初始化，无法再次注册"))
			return
		}

		var payload setupRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			logger.Error("[初始化] 解析请求体失败", "error", err)
			writeError(w, http.StatusBadRequest, err)
			return
		}

		username := strings.TrimSpace(payload.Username)
		password := strings.TrimSpace(payload.Password)
		nickname := strings.TrimSpace(payload.Nickname)
		email := strings.TrimSpace(payload.Email)
		avatarURL := strings.TrimSpace(payload.AvatarURL)

		logger.Info("[初始化] 准备创建管理员用户",
			"username", username,
			"nickname", nickname,
			"email", email,
		)

		if username == "" {
			writeError(w, http.StatusBadRequest, errors.New("用户名不能为空"))
			return
		}
		if err := validateUsername(username); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		if password == "" {
			writeError(w, http.StatusBadRequest, errors.New("密码不能为空"))
			return
		}

		if nickname == "" {
			nickname = username
		}

		// 对密码进行哈希处理
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("[初始化] 密码哈希失败", "error", err)
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		setupRepo := repo
		var selectedDatabase *storage.DatabaseConfig
		if payload.Database != nil {
			if storage.DatabaseConfigUsesEnvironment() {
				writeError(w, http.StatusConflict, errors.New("数据库连接正由环境变量管理，不能在初始化页面覆盖"))
				return
			}
			cfg := *payload.Database
			if err := cfg.Validate(); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			if cfg.Driver == "postgres" {
				candidate, err := storage.NewTrafficRepositoryFromConfig(cfg)
				if err != nil {
					writeError(w, http.StatusBadRequest, fmt.Errorf("PostgreSQL 初始化失败: %w", err))
					return
				}
				defer candidate.Close()
				existing, err := candidate.ListUsers(r.Context(), 1)
				if err != nil || len(existing) > 0 {
					if err == nil {
						err = errors.New("目标 PostgreSQL 已存在用户")
					}
					writeError(w, http.StatusConflict, err)
					return
				}
				setupRepo = candidate
				selectedDatabase = &cfg
			}
		}

		// 创建管理员用户
		if err := setupRepo.CreateUser(r.Context(), username, email, nickname, string(hash), storage.RoleAdmin, ""); err != nil {
			if errors.Is(err, storage.ErrUserExists) {
				logger.Warn("[初始化] 用户已存在", "username", username)
				writeError(w, http.StatusConflict, errors.New("用户已存在"))
				return
			}
			logger.Error("[初始化] 创建用户失败", "username", username, "error", err)
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		// 确保用户设置为管理员且处于活动状态
		_ = setupRepo.UpdateUserRole(r.Context(), username, storage.RoleAdmin)
		_ = setupRepo.UpdateUserStatus(r.Context(), username, true)
		if err := setupRepo.SetPrimaryAdminUsername(r.Context(), username); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("保存主控所有者失败: %w", err))
			return
		}

		if avatarURL != "" || email != "" || nickname != "" {
			_ = setupRepo.UpdateUserProfile(r.Context(), username, storage.UserProfileUpdate{
				Email:     email,
				Nickname:  nickname,
				AvatarURL: avatarURL,
			})
		}

		logger.Info("[初始化] 管理员用户创建成功",
			"username", username,
			"nickname", nickname,
			"role", storage.RoleAdmin,
		)

		resp := setupResponse{
			Username: username,
			Nickname: nickname,
			Email:    email,
		}

		domain := strings.TrimSpace(payload.Domain)
		if domain != "" {
			domain = strings.ToLower(domain)

			port := "12889"
			if _, p, err := net.SplitHostPort(r.Host); err == nil && p != "" {
				port = p
			}
			masterURL := fmt.Sprintf("http://%s:%s", domain, port)
			_ = setupRepo.SetSystemSetting(r.Context(), "master_url", masterURL)
			resp.RedirectURL = masterURL
			logger.Info("[初始化] 已保存 master_url", "master_url", masterURL)
		}
		if selectedDatabase != nil {
			dir := "data"
			if len(dataDir) > 0 && strings.TrimSpace(dataDir[0]) != "" {
				dir = dataDir[0]
			}
			if err := storage.SaveDatabaseConfig(dir, *selectedDatabase); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Errorf("保存数据库配置失败: %w", err))
				return
			}
			resp.Restarting = true
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
		if resp.Restarting {
			go func() { time.Sleep(500 * time.Millisecond); _ = SignalGracefulRestart() }()
		}
	})
}

// NewVerifyDomainHandler 验证域名 DNS 解析是否指向本机出口 IP
func NewVerifyDomainHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, errors.New("only POST is supported"))
			return
		}

		var req struct {
			Domain string `json:"domain"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		domain := strings.TrimSpace(strings.ToLower(req.Domain))
		if domain == "" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"success": false, "message": "域名不能为空"})
			return
		}

		domainIPs, err := net.LookupHost(domain)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"success": true, "match": false,
				"domain_ips": []string{}, "server_ip": "",
				"message": fmt.Sprintf("域名解析失败: %s", err.Error()),
			})
			return
		}

		serverIP := getOutboundIP()

		match := false
		for _, ip := range domainIPs {
			if ip == serverIP {
				match = true
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success":    true,
			"match":      match,
			"domain_ips": domainIPs,
			"server_ip":  serverIP,
		})
	})
}

func getOutboundIP() string {
	client := &http.Client{Timeout: 5 * time.Second}
	for _, url := range []string{"https://api.ipify.org", "https://ifconfig.me/ip", "https://icanhazip.com"} {
		resp, err := client.Get(url)
		if err != nil {
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		ip := strings.TrimSpace(string(body))
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

func deployLocalNginx(domain string, repo *storage.TrafficRepository) error {
	nginxBin := findNginxBinary()
	if nginxBin == "" {
		return fmt.Errorf("未找到 nginx 可执行文件")
	}
	if err := validateManagedNginx(nginxBin); err != nil {
		return err
	}

	nginxConf, err := templates.ReadFile("single_nginx.conf")
	if err != nil {
		return fmt.Errorf("读取 single_nginx.conf 模板失败: %w", err)
	}

	domainTpl, err := templates.ReadFile("mmwx_domain.conf")
	if err != nil {
		return fmt.Errorf("读取 mmwx_domain.conf 模板失败: %w", err)
	}
	domainConf := strings.ReplaceAll(string(domainTpl), "{domain}", domain)

	dirs := []string{
		"/usr/local/nginx/conf",
		"/usr/local/nginx/servers",
		"/usr/local/nginx/cert",
		"/usr/local/nginx/html",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
	}

	if err := os.WriteFile("/usr/local/nginx/nginx.conf", nginxConf, 0644); err != nil {
		return fmt.Errorf("写入 nginx.conf 失败: %w", err)
	}

	serverConfPath := filepath.Join("/usr/local/nginx/servers", domain+".conf")
	if err := os.WriteFile(serverConfPath, []byte(domainConf), 0644); err != nil {
		return fmt.Errorf("写入 domain.conf 失败: %w", err)
	}

	if repo != nil {
		deployCertToLocal(domain, repo)
	}

	if err := ensureNginxRunning(nginxBin); err != nil {
		return fmt.Errorf("nginx 启动失败: %w", err)
	}
	if err := markDockerNginxEnabled(); err != nil {
		return fmt.Errorf("保存 Docker HTTPS 启用状态失败: %w", err)
	}

	// 确保 nginx 开机自启 — 裸机部署一键脚本可能跳过 enable,服务器重启后 nginx 不起会让主控反代失效。
	// 失败只 warn 不阻塞:enable 是补防御,主流程已经把 nginx 跑起来了。Docker 容器内没 systemd,跳过。
	if !isDocker() {
		if err := exec.Command("systemctl", "enable", "nginx").Run(); err != nil {
			logger.Warn("[本机Nginx] systemctl enable nginx 失败 (开机自启未设置)", "error", err)
		}
	}
	return nil
}

func deployLocalNginxWithCert(domain string, cert *storage.Certificate) error {
	nginxBin := findNginxBinary()
	if nginxBin == "" {
		return fmt.Errorf("未找到 nginx 可执行文件")
	}
	if err := validateManagedNginx(nginxBin); err != nil {
		return err
	}

	nginxConf, err := templates.ReadFile("single_nginx.conf")
	if err != nil {
		return fmt.Errorf("读取 single_nginx.conf 模板失败: %w", err)
	}
	domainTpl, err := templates.ReadFile("mmwx_domain.conf")
	if err != nil {
		return fmt.Errorf("读取 mmwx_domain.conf 模板失败: %w", err)
	}
	certName := domain
	if cert != nil {
		certName = certDeployFilename(cert.Domain)
	}
	domainConf := strings.ReplaceAll(string(domainTpl), "{domain}", domain)
	domainConf = strings.ReplaceAll(domainConf, "{cert_name}", certName)

	dirs := []string{"/usr/local/nginx/conf", "/usr/local/nginx/servers", "/usr/local/nginx/cert", "/usr/local/nginx/html"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
	}
	if err := os.WriteFile("/usr/local/nginx/nginx.conf", nginxConf, 0644); err != nil {
		return fmt.Errorf("写入 nginx.conf 失败: %w", err)
	}
	if err := os.WriteFile(filepath.Join("/usr/local/nginx/servers", domain+".conf"), []byte(domainConf), 0644); err != nil {
		return fmt.Errorf("写入 domain.conf 失败: %w", err)
	}
	if cert != nil && cert.CertPEM != "" && cert.KeyPEM != "" {
		if err := os.WriteFile(filepath.Join("/usr/local/nginx/cert", certName+".pem"), []byte(cert.CertPEM), 0644); err != nil {
			return fmt.Errorf("写入证书失败: %w", err)
		}
		if err := os.WriteFile(filepath.Join("/usr/local/nginx/cert", certName+".key"), []byte(cert.KeyPEM), 0600); err != nil {
			return fmt.Errorf("写入密钥失败: %w", err)
		}
	}
	// 统一走 ensureNginxRunning:先 reload,失败时 Docker 用裸 nginx 拉起、裸机才 systemctl。
	// 之前这里直接 fallback systemctl,Docker 容器无 systemd → 报「systemctl 不在 $PATH」。
	if err := ensureNginxRunning(nginxBin); err != nil {
		return err
	}
	return markDockerNginxEnabled()
}

func deployCertToLocal(domain string, repo *storage.TrafficRepository) {
	ctx := context.Background()
	cert, err := repo.FindDeployableCertByDomain(ctx, domain, 0)
	if err != nil || cert == nil || cert.CertPEM == "" || cert.KeyPEM == "" {
		logger.Warn("[本机Nginx] 未找到域名证书，跳过证书部署", "domain", domain)
		return
	}
	certFilename := certDeployFilename(cert.Domain)
	certPath := filepath.Join("/usr/local/nginx/cert", certFilename+".pem")
	keyPath := filepath.Join("/usr/local/nginx/cert", certFilename+".key")
	if err := os.WriteFile(certPath, []byte(cert.CertPEM), 0644); err != nil {
		logger.Error("[本机Nginx] 写入证书失败", "error", err)
		return
	}
	if err := os.WriteFile(keyPath, []byte(cert.KeyPEM), 0600); err != nil {
		logger.Error("[本机Nginx] 写入密钥失败", "error", err)
		return
	}
	logger.Info("[本机Nginx] 证书部署成功", "domain", domain)
}
