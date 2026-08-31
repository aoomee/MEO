package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"miaomiaowux/internal/child"
	"miaomiaowux/internal/license"
	"miaomiaowux/internal/storage"
	"miaomiaowux/internal/version"
)

// ChildAPIHandler 处理来自主服务器的 API 请求（对于pull模式）
type ChildAPIHandler struct {
	client      *child.Client
	configToken string // 用于身份验证的令牌
}

// 创建一个新的子 API 处理程序
func NewChildAPIHandler(client *child.Client, configToken string) *ChildAPIHandler {
	return &ChildAPIHandler{
		client:      client,
		configToken: configToken,
	}
}

// 处理流量数据的 HTTP 请求
func (h *ChildAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 只允许 GET 方法
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 验证请求
	if !h.authenticate(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// 获取流量统计
	stats, err := h.client.GetStats()
	if err != nil {
		log.Printf("[Child API] Failed to get stats: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to collect stats",
		})
		return
	}

	// 返回统计数据
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"stats":   stats,
	})
}

// 处理速度数据的 HTTP 请求
func (h *ChildAPIHandler) ServeSpeedHTTP(w http.ResponseWriter, r *http.Request) {
	// 只允许 GET 方法
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 验证请求
	if !h.authenticate(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// 获取速度数据
	uploadSpeed, downloadSpeed := h.client.GetSpeed()

	// 返回速度数据
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":        true,
		"upload_speed":   uploadSpeed,
		"download_speed": downloadSpeed,
	})
}

// 验证检查请求是否被授权
func (h *ChildAPIHandler) authenticate(r *http.Request) bool {
	if h.configToken == "" {
		// 如果未配置令牌，则允许所有请求（不建议用于生产）
		return true
	}

	// 检查授权标头
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return false
	}

	// 支持“Bearer <token>”格式
	if strings.HasPrefix(auth, "Bearer ") {
		token := strings.TrimPrefix(auth, "Bearer ")
		return token == h.configToken
	}

	// 还支持普通令牌
	return auth == h.configToken
}

// RemoteHeartbeatRequest代表来自远程服务器的心跳请求
type RemoteHeartbeatRequest struct {
	BootTime     int64 `json:"boot_time"`      // MMWX进程启动时间（Unix时间戳）
	XrayBootTime int64 `json:"xray_boot_time"` // Xray 进程开始时间（Unix 时间戳）
	XrayPID      int   `json:"xray_pid"`       // 当前 X 射线进程 ID
	ListenPort   int   `json:"listen_port"`    // 代理HTTP监听端口
	LocalTime    int64 `json:"local_time"`     // agent 本地 Unix 时间戳，用于时钟偏差检测
	// PublicIPv4/v6 由 agent 端 ipProbeLoop 缓存后随心跳上报(WS auth/heartbeat 同款字段)。
	// master 优先用上报值写 db,fallback 才用 r.RemoteAddr 并强校验类型(避免 v6 写 v4 字段)。
	// 老 agent 不发 → 字段为空 → 走 fallback 路径,行为退化为现状。
	PublicIPv4        string `json:"public_ipv4,omitempty"`
	PublicIPv6        string `json:"public_ipv6,omitempty"`
	AgentNeedsLease   bool   `json:"agent_needs_lease,omitempty"`
	AgentLeaseCapable bool   `json:"agent_lease_capable,omitempty"`
	// AgentLeaseIdentityCapable is narrower than AgentLeaseCapable: it is true
	// only when the running Guard explicitly reports signed license identity.
	// Older Guards still support leases, but must use generation transitions
	// rather than an identity field during rolling upgrades.
	AgentLeaseIdentityCapable  bool                         `json:"agent_lease_identity_capable,omitempty"`
	AgentLeaseIdentity         string                       `json:"agent_lease_identity,omitempty"`
	AgentLeaseRequestIDCapable bool                         `json:"agent_lease_request_id_capable,omitempty"`
	AgentLeaseAckID            string                       `json:"agent_lease_ack_id,omitempty"`
	AgentLeaseFailure          *agentLeaseActivationFailure `json:"agent_lease_failure,omitempty"`
}

// RemoteHeartbeatResponse 表示心跳响应
type RemoteHeartbeatResponse struct {
	Success          bool                    `json:"success"`
	Message          string                  `json:"message"`
	MmwxRestarted    bool                    `json:"mmwx_restarted,omitempty"`     // 检测到 MMWX 重启
	XrayRestarted    bool                    `json:"xray_restarted,omitempty"`     // 检测到 X 射线重新启动
	TokenExpiresSoon bool                    `json:"token_expires_soon,omitempty"` // 令牌将在 24 小时内过期
	TokenExpiresAt   int64                   `json:"token_expires_at,omitempty"`   // 令牌过期时间戳
	ServerTime       int64                   `json:"server_time"`                  // 当前服务器时间
	MasterURL        string                  `json:"master_url,omitempty"`         // 恢复连接时让 Agent 持久化主控的新地址
	AgentLease       *httpAgentLeaseDelivery `json:"agent_lease,omitempty"`
}

type httpAgentLeasePending struct {
	token    string
	identity string
	force    bool
	attempt  int
	version  uint64
}

type httpAgentLeaseReady struct {
	token          string
	lease          license.AgentLeaseDelivery
	requestID      string
	deliveries     int
	failureHandled bool
	force          bool
}

type httpAgentLeaseDelivery struct {
	license.AgentLeaseDelivery
	RequestID string `json:"request_id,omitempty"`
}

type httpAgentLeaseIdentityMarker struct {
	token    string
	identity string
}

type httpAgentLeaseDeliveryCooldown struct {
	marker  httpAgentLeaseIdentityMarker
	retryAt time.Time
	force   bool
}

func (h *XrayServerHandler) initHTTPAgentLeaseQueue() {
	h.httpLeaseOnce.Do(func() {
		h.httpLeaseQueue = make(chan int64, 256)
		h.httpLeasePending = make(map[int64]*httpAgentLeasePending)
		h.httpLeaseReady = make(map[int64]httpAgentLeaseReady)
		h.httpLeaseRejected = make(map[int64]httpAgentLeaseIdentityMarker)
		h.httpLeaseCooldown = make(map[int64]httpAgentLeaseDeliveryCooldown)
		h.httpLeaseLegacyIdentity = make(map[int64]httpAgentLeaseIdentityMarker)
		go h.httpAgentLeaseWorker()
	})
}

// markHTTPAgentLeaseIdentity returns true once for each server/token/license
// identity tuple. The queue owns retries after the marker is installed. A
// process restart may cause one extra replacement, which is idempotent.
func (h *XrayServerHandler) markHTTPAgentLeaseIdentity(serverID int64, token, identity string) bool {
	if serverID <= 0 || token == "" || identity == "" {
		return false
	}
	h.initHTTPAgentLeaseQueue()
	h.httpLeaseMu.Lock()
	defer h.httpLeaseMu.Unlock()
	want := httpAgentLeaseIdentityMarker{token: token, identity: identity}
	if h.httpLeaseLegacyIdentity[serverID] == want {
		return false
	}
	h.httpLeaseLegacyIdentity[serverID] = want
	return true
}

func (h *XrayServerHandler) httpAgentLeasePassiveRetryDue(serverID int64, token, identity string) bool {
	if serverID <= 0 || token == "" || identity == "" {
		return false
	}
	h.initHTTPAgentLeaseQueue()
	h.httpLeaseMu.Lock()
	defer h.httpLeaseMu.Unlock()
	cooldown, ok := h.httpLeaseCooldown[serverID]
	return ok && cooldown.marker == (httpAgentLeaseIdentityMarker{token: token, identity: identity}) &&
		!time.Now().Before(cooldown.retryAt)
}

func (h *XrayServerHandler) queueHTTPAgentLease(serverID int64, token, identity string, force bool) {
	if serverID <= 0 || token == "" || identity == "" {
		return
	}
	h.initHTTPAgentLeaseQueue()
	h.httpLeaseMu.Lock()
	marker := httpAgentLeaseIdentityMarker{token: token, identity: identity}
	if h.httpLeaseRejected[serverID] == marker {
		h.httpLeaseMu.Unlock()
		return
	}
	if cooldown, ok := h.httpLeaseCooldown[serverID]; ok {
		if cooldown.marker == marker && time.Now().Before(cooldown.retryAt) {
			h.httpLeaseMu.Unlock()
			return
		}
		if cooldown.marker == marker {
			force = force || cooldown.force
		}
		// The old reservation was delivered repeatedly without an activation
		// ACK. Once its cooldown expires, discard it and issue a fresh one on the
		// next needs-lease heartbeat instead of freezing this tuple forever.
		delete(h.httpLeaseCooldown, serverID)
		delete(h.httpLeaseReady, serverID)
	}
	delete(h.httpLeaseRejected, serverID)
	if pending, ok := h.httpLeasePending[serverID]; ok {
		if pending.token != token || pending.identity != identity || (force && !pending.force) {
			pending.token = token
			pending.identity = identity
			pending.force = pending.force || force
			pending.attempt = 0
			pending.version++
		}
		h.httpLeaseMu.Unlock()
		return
	}
	h.httpLeasePending[serverID] = &httpAgentLeasePending{token: token, identity: identity, force: force, version: 1}
	h.httpLeaseMu.Unlock()
	h.enqueueHTTPAgentLease(serverID)
}

func (h *XrayServerHandler) enqueueHTTPAgentLease(serverID int64) {
	select {
	case h.httpLeaseQueue <- serverID:
	default:
		time.AfterFunc(time.Second, func() { h.enqueueHTTPAgentLease(serverID) })
	}
}

func (h *XrayServerHandler) httpAgentLeaseWorker() {
	throttle := time.NewTicker(750 * time.Millisecond)
	defer throttle.Stop()
	for serverID := range h.httpLeaseQueue {
		<-throttle.C
		h.httpLeaseMu.Lock()
		pending, ok := h.httpLeasePending[serverID]
		if !ok {
			h.httpLeaseMu.Unlock()
			continue
		}
		snapshot := *pending
		h.httpLeaseMu.Unlock()

		var delivery license.AgentLeaseDelivery
		var err error
		if h.licenseManager == nil {
			err = errors.New("license manager unavailable")
		} else {
			issueCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			if snapshot.force {
				delivery, err = h.licenseManager.IssueAgentLeaseReplacement(issueCtx, snapshot.token)
			} else {
				delivery, err = h.licenseManager.IssueAgentLease(issueCtx, snapshot.token)
			}
			cancel()
		}

		h.httpLeaseMu.Lock()
		current, stillPending := h.httpLeasePending[serverID]
		if !stillPending {
			h.httpLeaseMu.Unlock()
			continue
		}
		if current.version != snapshot.version {
			h.httpLeaseMu.Unlock()
			h.enqueueHTTPAgentLease(serverID)
			continue
		}
		if err == nil {
			delete(h.httpLeasePending, serverID)
			delete(h.httpLeaseRejected, serverID)
			delete(h.httpLeaseCooldown, serverID)
			h.httpLeaseSequence++
			h.httpLeaseReady[serverID] = httpAgentLeaseReady{
				token: snapshot.token, lease: delivery,
				requestID: fmt.Sprintf("http-%d-%d-%d", serverID, time.Now().UnixNano(), h.httpLeaseSequence),
				force:     snapshot.force,
			}
			h.httpLeaseMu.Unlock()
			continue
		}
		var issueErr *license.AgentLeaseIssueError
		if errors.As(err, &issueErr) && issueErr.Code == "server_slot_state_conflict" {
			if !snapshot.force {
				// A lost/reinstalled local Guard can ask for an ordinary lease while the
				// service still owns an active slot. Upgrade exactly once to the guarded
				// replacement path; retrying the ordinary request can never succeed.
				current.force = true
				current.attempt = 0
				current.version++
				h.httpLeaseMu.Unlock()
				log.Printf("[RemoteHeartbeat] lease state conflict upgraded to one forced replacement server=%d", serverID)
				h.enqueueHTTPAgentLease(serverID)
				continue
			}
			// During a rolling deployment an old license service may reject the new
			// replacement flag with this legacy conflict. Do not retain a background
			// job or timer: an offline Agent would otherwise keep hitting the license
			// service forever. The next real heartbeat after cooldown issues fresh.
			retryAt := h.finishHTTPAgentLeaseForcedConflictLocked(serverID, snapshot)
			h.httpLeaseMu.Unlock()
			log.Printf("[RemoteHeartbeat] forced lease state conflict server=%d; waiting for heartbeat after %v: %v",
				serverID, time.Until(retryAt).Round(time.Millisecond), err)
			continue
		}
		if h.licenseManager == nil || !license.AgentLeaseShouldRetry(err) {
			delete(h.httpLeasePending, serverID)
			h.httpLeaseRejected[serverID] = httpAgentLeaseIdentityMarker{token: snapshot.token, identity: snapshot.identity}
			h.httpLeaseMu.Unlock()
			log.Printf("[RemoteHeartbeat] lease issue permanently rejected server=%d: %v", serverID, err)
			continue
		}
		current.attempt++
		attempt := current.attempt
		retryAt := h.finishHTTPAgentLeaseRetryLocked(serverID, snapshot, err, attempt)
		h.httpLeaseMu.Unlock()

		log.Printf("[RemoteHeartbeat] lease issue waiting for heartbeat retry server=%d attempt=%d in=%v: %v",
			serverID, attempt, time.Until(retryAt).Round(time.Millisecond), err)
	}
}

// finishHTTPAgentLeaseForcedConflictLocked converts a forced legacy conflict
// into passive cooldown state. It deliberately creates neither a pending job
// nor a timer; queueHTTPAgentLease is re-armed only by a later Agent heartbeat.
// h.httpLeaseMu must be held by the caller.
func (h *XrayServerHandler) finishHTTPAgentLeaseForcedConflictLocked(serverID int64, snapshot httpAgentLeasePending) time.Time {
	return h.finishHTTPAgentLeasePassiveCooldownLocked(serverID, snapshot, agentLeaseAckTerminalCooldown)
}

// finishHTTPAgentLeaseRetryLocked records Retry-After/backoff without keeping
// a pending job. Network errors, 429 and 5xx therefore stop consuming the
// serialized global queue as soon as the Agent stops heartbeating.
// h.httpLeaseMu must be held by the caller.
func (h *XrayServerHandler) finishHTTPAgentLeaseRetryLocked(serverID int64, snapshot httpAgentLeasePending, err error, attempt int) time.Time {
	delay := license.AgentLeaseRetryAfter(err)
	if delay <= 0 {
		delay = time.Second << min(max(attempt, 1)-1, 6)
		if delay > 2*time.Minute {
			delay = 2 * time.Minute
		}
	}
	delay += time.Duration(serverID%1000) * time.Millisecond
	return h.finishHTTPAgentLeasePassiveCooldownLocked(serverID, snapshot, delay)
}

func (h *XrayServerHandler) finishHTTPAgentLeasePassiveCooldownLocked(serverID int64, snapshot httpAgentLeasePending, delay time.Duration) time.Time {
	retryAt := time.Now().Add(delay)
	marker := httpAgentLeaseIdentityMarker{token: snapshot.token, identity: snapshot.identity}
	delete(h.httpLeasePending, serverID)
	delete(h.httpLeaseReady, serverID)
	delete(h.httpLeaseRejected, serverID)
	h.httpLeaseCooldown[serverID] = httpAgentLeaseDeliveryCooldown{marker: marker, retryAt: retryAt, force: snapshot.force}
	return retryAt
}

// ResetHTTPAgentLeaseFailures lets a verified key or quota transition retry
// requests that were terminal under the previous entitlement. The next Agent
// heartbeat repopulates the serialized queue, so this does not create fan-out.
func (h *XrayServerHandler) ResetHTTPAgentLeaseFailures() {
	h.initHTTPAgentLeaseQueue()
	h.httpLeaseMu.Lock()
	clear(h.httpLeaseRejected)
	clear(h.httpLeaseCooldown)
	clear(h.httpLeaseLegacyIdentity)
	h.httpLeaseMu.Unlock()
}

// handleHTTPAgentLeaseFailure consumes only a failure for the exact ready
// reservation. Structured permanent rejections suppress the same tuple until
// a license/quota reset; temporary failures retain the reservation and honor
// the Guard-provided Retry-After before redelivery.
func (h *XrayServerHandler) handleHTTPAgentLeaseFailure(serverID int64, token, currentIdentity string, failure *agentLeaseActivationFailure) bool {
	if failure == nil || failure.RequestID == "" {
		return false
	}
	h.initHTTPAgentLeaseQueue()
	h.httpLeaseMu.Lock()
	ready, ok := h.httpLeaseReady[serverID]
	if !ok || ready.requestID != failure.RequestID || ready.token != token ||
		ready.lease.LicenseIdentity != currentIdentity {
		h.httpLeaseMu.Unlock()
		return false
	}
	if ready.failureHandled {
		h.httpLeaseMu.Unlock()
		return true
	}
	err := newAgentLeaseActivationError(failure.Error, failure.Code, failure.UpstreamStatus, failure.RetryAfter)
	marker := httpAgentLeaseIdentityMarker{token: token, identity: currentIdentity}
	if !license.AgentLeaseShouldRetry(err) {
		delete(h.httpLeaseReady, serverID)
		delete(h.httpLeasePending, serverID)
		delete(h.httpLeaseCooldown, serverID)
		h.httpLeaseRejected[serverID] = marker
		h.httpLeaseMu.Unlock()
		log.Printf("[RemoteHeartbeat] Agent lease activation permanently rejected server=%d request_id=%q: %v", serverID, failure.RequestID, err)
		return true
	}

	delay := license.AgentLeaseRetryAfter(err)
	if delay <= 0 {
		attempt := max(ready.deliveries, 1)
		delay = time.Second << min(attempt-1, 6)
	}
	if agentLeaseActivationCanSelfRepair(failure.Code) {
		// Retrying the same reservation cannot fix a stale serial/authority.
		// Remove it now and let the first needsLease heartbeat after the short
		// cooldown run the normal issue→force-replacement state machine.
		delete(h.httpLeaseReady, serverID)
		delete(h.httpLeasePending, serverID)
		h.httpLeaseCooldown[serverID] = httpAgentLeaseDeliveryCooldown{marker: marker, retryAt: time.Now().Add(delay), force: true}
		h.httpLeaseMu.Unlock()
		log.Printf("[RemoteHeartbeat] Agent lease activation requires fresh reservation server=%d request_id=%q in=%v: %v", serverID, failure.RequestID, delay, err)
		return true
	}

	ready.failureHandled = true
	h.httpLeaseReady[serverID] = ready
	h.httpLeaseCooldown[serverID] = httpAgentLeaseDeliveryCooldown{marker: marker, retryAt: time.Now().Add(delay), force: ready.force}
	h.httpLeaseMu.Unlock()
	log.Printf("[RemoteHeartbeat] Agent lease activation retry server=%d request_id=%q in=%v: %v", serverID, failure.RequestID, delay, err)
	return true
}

func (h *XrayServerHandler) takeHTTPAgentLease(serverID int64, token, currentIdentity, ackID string, requestIDCapable bool) (*httpAgentLeaseDelivery, bool) {
	h.initHTTPAgentLeaseQueue()
	h.httpLeaseMu.Lock()
	defer h.httpLeaseMu.Unlock()
	ready, ok := h.httpLeaseReady[serverID]
	if !ok {
		return nil, false
	}
	if currentIdentity == "" || ready.token != token || ready.lease.ExpiresAt <= time.Now().Unix() ||
		ready.lease.LicenseIdentity != currentIdentity {
		delete(h.httpLeaseReady, serverID)
		delete(h.httpLeaseCooldown, serverID)
		return nil, false
	}
	if requestIDCapable && ackID != "" && ackID == ready.requestID {
		delete(h.httpLeaseReady, serverID)
		delete(h.httpLeaseRejected, serverID)
		delete(h.httpLeaseCooldown, serverID)
		return nil, false
	}
	if cooldown, exists := h.httpLeaseCooldown[serverID]; exists &&
		cooldown.marker == (httpAgentLeaseIdentityMarker{token: token, identity: currentIdentity}) {
		if time.Now().Before(cooldown.retryAt) {
			return nil, false
		}
		delete(h.httpLeaseCooldown, serverID)
		if requestIDCapable && ready.deliveries >= 3 {
			// The bounded delivery cycle is over. Discard the old request here so
			// queueHTTPAgentLease later in this heartbeat can issue a fresh one.
			delete(h.httpLeaseReady, serverID)
			return nil, false
		}
		ready.failureHandled = false
		h.httpLeaseReady[serverID] = ready
	}
	if requestIDCapable && ready.deliveries >= 3 {
		// A capable Agent only records ackID after Guard activation succeeds.
		// Stop redelivering a deterministically rejected reservation on every
		// heartbeat, but retain its ID so a delayed successful ACK can clear the
		// terminal marker without a license/quota transition.
		h.httpLeaseCooldown[serverID] = httpAgentLeaseDeliveryCooldown{
			marker:  httpAgentLeaseIdentityMarker{token: token, identity: currentIdentity},
			retryAt: time.Now().Add(agentLeaseAckTerminalCooldown),
			force:   ready.force,
		}
		return nil, false
	}
	ready.deliveries++
	// New Agents ACK on the following heartbeat, so retain the reservation until
	// that exact ID arrives. Legacy Agents cannot ACK; repeat a bounded number of
	// times to survive a lost response without sending a consumed reservation
	// forever after successful Guard activation.
	if !requestIDCapable && ready.deliveries >= 6 {
		delete(h.httpLeaseReady, serverID)
	} else {
		h.httpLeaseReady[serverID] = ready
	}
	delivery := httpAgentLeaseDelivery{AgentLeaseDelivery: ready.lease, RequestID: ready.requestID}
	return &delivery, true
}

// RemoteHeartbeat 处理来自远程服务器的心跳请求
// 该端点不需要管理员身份验证，只需要远程令牌验证
func (h *XrayServerHandler) RemoteHeartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("User-Agent") != version.AgentUserAgent {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(RemoteHeartbeatResponse{
			Success:    false,
			Message:    "Forbidden",
			ServerTime: time.Now().Unix(),
		})
		return
	}

	// 加密中间件处理
	crypto, cryptoErr := handleHTTPCrypto(r, w, h.crypto)
	if crypto == nil {
		return
	}
	_ = cryptoErr

	token := crypto.Token
	if token == "" {
		token = r.Header.Get("MM-Remote-Token")
	}
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(RemoteHeartbeatResponse{
			Success:    false,
			Message:    "缺少认证Token",
			ServerTime: time.Now().Unix(),
		})
		return
	}

	// 解析请求体
	var req RemoteHeartbeatRequest
	json.Unmarshal(crypto.Body, &req)

	// 获取客户端IP — X-Forwarded-For > X-Real-IP > r.RemoteAddr
	rawIP := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// 从逗号分隔列表中获取第一个 IP
		rawIP = strings.Split(forwarded, ",")[0]
		rawIP = strings.TrimSpace(rawIP)
	} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		rawIP = realIP
	}
	// 用 stripPort 正确处理 v4 / [v6]:port / 裸 v6 三种形式。
	// 老代码用 strings.LastIndex(":") 截断,对裸 v6 会把最后一段误删,留下半截 v6 字符串塞进 db.ip_address。
	clientIP := stripPort(rawIP)
	clientParsed := net.ParseIP(clientIP)
	clientIsV4 := clientParsed != nil && clientParsed.To4() != nil

	// 严格选 v4 / v6 字段(同 WS handleHeartbeat 模式):
	//   1) 优先用 agent 上报的 public_ipv4/public_ipv6(经 ipProbeLoop 校验过的本机出口 IP)
	//   2) fallback 用 clientIP,但**只在类型匹配时**才写对应字段 — 避免 agent v6 拨号 master →
	//      master 把 clientIP(v6) 当 v4 塞进 ip_address → IPv4-only master 反向请求全部失败
	v4 := ""
	if reported := strings.TrimSpace(req.PublicIPv4); reported != "" {
		if p := net.ParseIP(reported); p != nil && p.To4() != nil {
			v4 = reported
		}
	}
	if v4 == "" && clientIsV4 {
		v4 = clientIP
	}

	v6 := ""
	if reported := strings.TrimSpace(req.PublicIPv6); reported != "" {
		if p := net.ParseIP(reported); p != nil && p.To4() == nil {
			v6 = reported
		}
	}
	if v6 == "" && clientParsed != nil && !clientIsV4 {
		v6 = clientIP
	}

	ctx := r.Context()

	// 构建心跳更新 — v4/v6 字段空字符串走 storage 层 COALESCE/NULLIF 保留 db 旧值
	update := storage.HeartbeatUpdate{
		Token:       token,
		IPAddress:   v4,
		IPAddressV6: v6,
		ListenPort:  req.ListenPort,
	}

	// 从 Unix 时间戳转换启动时间
	if req.BootTime > 0 {
		bootTime := time.Unix(req.BootTime, 0)
		update.BootTime = &bootTime
	}
	if req.XrayBootTime > 0 {
		xrayBootTime := time.Unix(req.XrayBootTime, 0)
		update.XrayBootTime = &xrayBootTime
	}
	if req.LocalTime > 0 {
		offset := req.LocalTime - time.Now().Unix()
		update.TimeOffsetSeconds = &offset
	}

	// 通过重启检测更新心跳
	result, err := h.repo.UpdateRemoteServerHeartbeatWithRestart(ctx, update)
	if err != nil {
		if err == storage.ErrRemoteServerNotFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(RemoteHeartbeatResponse{
				Success:    false,
				Message:    "无效的Token",
				ServerTime: time.Now().Unix(),
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RemoteHeartbeatResponse{
			Success:    false,
			Message:    fmt.Sprintf("更新心跳失败: %s", err.Error()),
			ServerTime: time.Now().Unix(),
		})
		return
	}

	// 记录重启事件
	if result.MmwxRestarted {
		log.Printf("[RemoteHeartbeat] Detected MMWX restart for token %s... (boot count: %d)", token[:8], result.BootCount)
	}
	if result.XrayRestarted {
		log.Printf("[RemoteHeartbeat] Detected Xray restart for token %s... (xray boot count: %d)", token[:8], result.XrayBootCount)
	}

	// 与 WS 和流量上报的恢复策略保持一致：只有本轮离线已经达到容忍阈值、
	// 且确实发过离线通知，恢复时才发送配对的上线通知。短暂 WS 抖动虽然会
	// 把状态置为 offline，但 offline_notified 仍为 false，必须保持静默。
	if result.PreviousStatus == storage.RemoteServerStatusOffline && result.PreviousOfflineNotified {
		if IsServerUpgrading(result.ServerID) {
			log.Printf("[RemoteHeartbeat] %s came back online during upgrade window — suppressing online notification", result.ServerName)
		} else {
			go SendServerOnlineNotification(context.Background(), result.ServerName, clientIP)
		}
	}

	// agent IP 漂移 → 同步刷新已存在节点的 clash_config.server,避免节点继续指向旧 IP
	if result.IPChanged && result.Server != nil {
		if h.remoteManager != nil && result.PreviousServer != nil {
			before, after := *result.PreviousServer, *result.Server
			go h.remoteManager.SyncServerAddressChange(context.Background(), &before, &after)
		}
		// DDNS:把新 IP 同步到 pull_address 域名的 A/AAAA 记录
		if h.ddnsManager != nil && result.Server.DDNSEnabled {
			go h.ddnsManager.Trigger(context.Background(), result.Server)
		}
	}

	// 首次连接或 Xray 重启时推送限速配置（非 WebSocket 模式的补偿）
	if result.ServerID > 0 && h.limiterPusher != nil {
		if result.PreviousStatus != "connected" || result.XrayRestarted {
			go h.limiterPusher.PushToServer(context.Background(), result.ServerID)
		}
	}

	// 重置成功心跳时的推送失败计数（连接正常）
	if result.ServerID > 0 {
		if err := h.repo.ResetRemoteServerPushFailCount(ctx, result.ServerID); err != nil {
			log.Printf("[RemoteHeartbeat] Failed to reset push fail count for server %d: %v", result.ServerID, err)
		}
	}

	resp := RemoteHeartbeatResponse{
		Success:          true,
		Message:          "心跳成功",
		MmwxRestarted:    result.MmwxRestarted,
		XrayRestarted:    result.XrayRestarted,
		TokenExpiresSoon: result.TokenExpiresSoon,
		ServerTime:       time.Now().Unix(),
	}
	resp.MasterURL, _ = h.repo.GetSystemSetting(ctx, "master_url")
	if h.licenseManager != nil {
		currentIdentity := h.licenseManager.AgentLeaseIdentity()
		identityMismatch := agentLeaseIdentityIsStale(req.AgentLeaseIdentityCapable, currentIdentity, req.AgentLeaseIdentity)
		legacyIdentityRefresh := h.licenseManager.AgentGuardRequired() && req.AgentLeaseCapable &&
			!req.AgentLeaseIdentityCapable && h.markHTTPAgentLeaseIdentity(result.ServerID, token, currentIdentity)
		passiveRetry := h.httpAgentLeasePassiveRetryDue(result.ServerID, token, currentIdentity)
		if h.licenseManager.AgentGuardRequired() && !req.AgentLeaseCapable && !req.AgentNeedsLease {
			writeJSONError(w, http.StatusUpgradeRequired, "当前许可证已启用严格联合验证，请升级 Agent")
			return
		}
		h.handleHTTPAgentLeaseFailure(result.ServerID, token, currentIdentity, req.AgentLeaseFailure)
		if lease, ready := h.takeHTTPAgentLease(result.ServerID, token, currentIdentity,
			req.AgentLeaseAckID, req.AgentLeaseRequestIDCapable); ready {
			resp.AgentLease = lease
		} else if req.AgentNeedsLease || identityMismatch || legacyIdentityRefresh || passiveRetry {
			// HTTP-only Agents receive the result on a later heartbeat. Issuance is
			// serialized, deduplicated and retried with Retry-After/backoff so a
			// license change cannot fan out into a request storm.
			h.queueHTTPAgentLease(result.ServerID, token, currentIdentity, identityMismatch || legacyIdentityRefresh)
		}
	}

	if result.TokenExpiresAt != nil {
		resp.TokenExpiresAt = result.TokenExpiresAt.Unix()
	}

	respData, _ := json.Marshal(resp)
	writeHTTPCryptoResponse(w, crypto.Session, respData)
}

// RefreshRemoteTokenResponse 是令牌刷新端点的响应
type RefreshRemoteTokenResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	NewToken  string `json:"new_token,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"` // Unix时间戳
}

// 处理远程服务器的令牌刷新
func (h *XrayServerHandler) RefreshRemoteToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("User-Agent") != version.AgentUserAgent {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(RefreshRemoteTokenResponse{
			Success: false,
			Message: "Forbidden",
		})
		return
	}

	// 从标头获取令牌
	token := r.Header.Get("MM-Remote-Token")
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(RefreshRemoteTokenResponse{
			Success: false,
			Message: "Missing MM-Remote-Token header",
		})
		return
	}

	// 刷新前先读一次:回滚需要 server id 和原过期时间。
	ctx := r.Context()
	server, serverErr := h.repo.GetRemoteServerByToken(ctx, token)

	// 尝试刷新令牌
	newToken, expiresAt, err := h.repo.RefreshRemoteServerToken(ctx, token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		// 检查具体错误
		if err.Error() == "token can only be refreshed within 24 hours of expiration" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(RefreshRemoteTokenResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		if errors.Is(err, storage.ErrRemoteServerNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(RefreshRemoteTokenResponse{
				Success: false,
				Message: "Invalid token",
			})
			return
		}

		log.Printf("[Remote] Failed to refresh token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RefreshRemoteTokenResponse{
			Success: false,
			Message: "Failed to refresh token",
		})
		return
	}

	// Agent 每 7 天自助续期一次,是 server token 轮换最高频的入口。这里同样必须把
	// Guard 绑定的授权槽位迁到新 token hash —— 以前这条路径完全不迁,槽位每续期一次
	// 就和实际 token 错位一次。迁不动就回滚,让 Agent 拿着旧令牌稍后重试,
	// 也好过换到一个没有授权槽位的新令牌。
	if serverErr == nil && server != nil {
		lease, needsDelivery, leaseErr := rotateServerSlot(ctx, h.licenseManager, token, newToken)
		if leaseErr != nil {
			if rollbackErr := rollbackServerTokenRotation(h.repo, server.ID, token, server.TokenExpiresAt); rollbackErr != nil {
				log.Printf("[Remote] 授权槽位迁移失败且令牌回滚失败 server=%d: %v (回滚: %v)", server.ID, leaseErr, rollbackErr)
			} else {
				log.Printf("[Remote] 授权槽位迁移失败,已回滚令牌 server=%d: %v", server.ID, leaseErr)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(RefreshRemoteTokenResponse{
				Success: false,
				Message: "Authoritative slot rotation failed, token unchanged",
			})
			return
		}
		if needsDelivery {
			deliverRotatedSlot(h.wsHandler, server.ID, lease)
		}
	} else if serverErr != nil {
		log.Printf("[Remote] 令牌已刷新,但读不到服务器记录、无法迁移授权槽位: %v", serverErr)
	}

	log.Printf("[Remote] Token refreshed successfully, new expiration: %s", expiresAt.Format(time.RFC3339))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RefreshRemoteTokenResponse{
		Success:   true,
		Message:   "Token refreshed successfully",
		NewToken:  newToken,
		ExpiresAt: expiresAt.Unix(),
	})
}

func (h *XrayServerHandler) masterPublicKeyBase64() string {
	if h.crypto != nil && h.crypto.Identity != nil {
		return h.crypto.Identity.PublicKeyBase64()
	}
	return ""
}

// validInstallToken 校验安装 token 字符集。token 会被写进 curl|bash 执行的安装脚本,
// 必须白名单化,否则 $(...)/`...` 会被当命令替换执行(命令注入)。
//
// 字符集 [A-Za-z0-9._-],外加**结尾**最多两个 '=' 的 base64 padding。
//
// padding 规则为存量兼容:2026-08-12 之前 generateSecureToken 使用
// base64.URLEncoding,32 字节会编成带结尾 '=' 的 44 字符。新 token 已改用
// RawURLEncoding 避免安装命令出现 %3D，但旧 token 仍必须可用。
// '=' 只允许出现在结尾:中间出现的 '=' 不是合法 base64,没有放行的理由。
func validInstallToken(s string) bool {
	if s == "" || len(s) > 512 {
		return false
	}
	body := strings.TrimRight(s, "=")
	if len(s)-len(body) > 2 || body == "" {
		return false
	}
	for _, c := range body {
		switch {
		case c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z', c >= '0' && c <= '9',
			c == '.', c == '_', c == '-':
		default:
			return false
		}
	}
	return true
}

// shSingleQuote 把任意字符串安全包成 bash 单引号字面量(单引号内除 ' 外无特殊字符)。
// 规则:整体套单引号,内部每个 ' 替换成 '\” 序列。用于把外部输入写进安装脚本时防注入。
func shSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

// 返回远程服务器的安装脚本
func (h *XrayServerHandler) GetRemoteInstallScript(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从查询参数中获取令牌
	token := r.URL.Query().Get("token")
	// 安全:token 会被写进 curl|bash 执行的安装脚本,必须白名单校验,否则命令注入。
	if !validInstallToken(token) {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}
	server, err := h.repo.GetRemoteServerByToken(r.Context(), token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusForbidden)
		return
	}
	stealSelf := r.URL.Query().Get("steal_self") == "1"
	if stealSelf {
		if forbidMasterHTTPSSteal(r.Context(), h.repo, server) {
			http.Error(w, masterHTTPSStealMessage, http.StatusConflict)
			return
		}
	}
	xrayMode := r.URL.Query().Get("xray_mode")
	if xrayMode != "embedded" {
		xrayMode = "external"
	}
	// 自定义 Agent 监听端口(由主控创建服务器时透传过来),非法/缺省值用 agent 内置默认 23889
	listenPortParam := strings.TrimSpace(r.URL.Query().Get("listen_port"))
	if p, perr := strconv.Atoi(listenPortParam); perr != nil || p < 1024 || p > 65535 {
		listenPortParam = ""
	}
	frontService := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("front_service")))
	if frontService != "xray" && frontService != "nginx" {
		frontService = "xray"
	}
	// nginx 前置暂未支持，先固定为 xray
	if frontService == "nginx" {
		frontService = "xray"
	}

	// 计算 install 脚本里写入的 SERVER:
	// 优先用系统设置 master_url 里的 host(用户配置的对外可达域名),
	// 这是 agent 真正访问主控的地址。仅在 master_url 未配置时回退到 r.Host(可能是 nginx upstream 名,如 miaomiaowu_web,不可对外访问)。
	// 若 master_url 已显式配置,EXPLICIT_MASTER=1 在脚本里禁用"同机部署"自动覆盖
	// (避免在主控本机上安装 agent 时把 master_url 改写成 127.0.0.1)。
	scriptServer := strings.TrimSpace(r.Host)
	// nginx 默认 `proxy_set_header Host $host` 不带端口,导致 cf:8443 → nginx → mmwx 时 r.Host 只有域名,
	// agent 安装命令缺端口连不上主控。这里如果检测到 X-Forwarded-Host(带端口最优)或 X-Forwarded-Port
	// 且端口不是 80/443,主动把 :port 拼回去,方便用户不需要必须先去配 master_url 就能拿到正确安装命令。
	if !strings.Contains(scriptServer, ":") {
		if xfh := strings.TrimSpace(r.Header.Get("X-Forwarded-Host")); xfh != "" && strings.Contains(xfh, ":") {
			scriptServer = xfh
		} else if xfp := strings.TrimSpace(r.Header.Get("X-Forwarded-Port")); xfp != "" && xfp != "80" && xfp != "443" {
			scriptServer = scriptServer + ":" + xfp
		}
	}
	scriptProtocol := ""
	// nginx 反代下大概率有 X-Forwarded-Proto,带这个就别走脚本里 "host 有 : 就当 http" 的启发,直接显式 https
	if xfproto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); xfproto == "https" || xfproto == "http" {
		scriptProtocol = xfproto
	}
	if mu, err := h.repo.GetSystemSetting(r.Context(), "master_url"); err == nil {
		mu = strings.TrimSpace(mu)
		if mu != "" {
			s := strings.TrimRight(mu, "/")
			if strings.HasPrefix(s, "https://") {
				scriptProtocol = "https"
				s = strings.TrimPrefix(s, "https://")
			} else if strings.HasPrefix(s, "http://") {
				scriptProtocol = "http"
				s = strings.TrimPrefix(s, "http://")
			}
			if i := strings.Index(s, "/"); i >= 0 {
				s = s[:i]
			}
			if s != "" {
				scriptServer = s
			}
		}
	}
	scriptRecoveryURL, _, recoveryErr := effectiveRecoveryURL(r.Context(), h.repo, masterListenPort())
	if recoveryErr != nil {
		port := masterListenPort()
		host := scriptServer
		if parsed, err := url.Parse("//" + scriptServer); err == nil && parsed.Hostname() != "" {
			host = parsed.Hostname()
		}
		scriptRecoveryURL = "http://" + net.JoinHostPort(host, port)
	}

	// 返回安装脚本内容
	script := `#!/bin/bash
# MMWX Remote Server Installation Script
# This script installs MMWX from GitHub and configures it as a remote server

set -e

TOKEN=` + shSingleQuote(token) + `
SERVER=` + shSingleQuote(scriptServer) + `
SCRIPT_PROTOCOL=` + shSingleQuote(scriptProtocol) + `
AUTO_STEAL_SELF="` + map[bool]string{true: "1", false: "0"}[stealSelf] + `"
FRONT_SERVICE=` + shSingleQuote(frontService) + `
XRAY_MODE=` + shSingleQuote(xrayMode) + `
MASTER_PUBLIC_KEY=` + shSingleQuote(h.masterPublicKeyBase64()) + `
RECOVERY_URL=` + shSingleQuote(scriptRecoveryURL) + `
LISTEN_PORT=` + shSingleQuote(listenPortParam) + `

# 协议:优先用主控注入的 SCRIPT_PROTOCOL(来自系统设置 master_url 的 scheme),
# 否则按 SERVER 是否带端口启发判断(开发场景常见 http)。
if [ -n "$SCRIPT_PROTOCOL" ]; then
    PROTOCOL="$SCRIPT_PROTOCOL"
elif [[ "$SERVER" == *":"* ]]; then
    PROTOCOL="http"
else
    PROTOCOL="https"
fi

# 允许通过环境变量强制覆盖协议
if [ -n "$MMWX_PROTOCOL" ]; then
    PROTOCOL="$MMWX_PROTOCOL"
fi

MASTER_URL="${PROTOCOL}://${SERVER}"

# 最后一层本机保护：主控已使用 HTTPS 时，在主控机器上安装 Agent 绝不能
# 启用偷自己抢占 443。兼容直装主控和常见 Docker 挂载目录。
if [ "$AUTO_STEAL_SELF" = "1" ] && [ "$PROTOCOL" = "https" ]; then
    if pgrep -x mmwx >/dev/null 2>&1 || [ -d /etc/mmwx ] || [ -d /var/lib/mmwx ]; then
        echo "ERROR: ` + masterHTTPSStealMessage + `" >&2
        exit 1
    fi
fi

echo "=========================================="
echo "  MMWX Remote Server Installation"
echo "=========================================="
echo ""
echo "Master Server: $MASTER_URL"
echo ""

# 检测 init 系统:OpenRC(Alpine 首选)/ systemd(主流)/ 兜底用 nohup + rc.local。
# - Alpine 优先用 OpenRC:Alpine 主流就是 OpenRC,即便镜像里塞了 systemd 也不用它
# - Alpine 极简镜像/LXC 可能没装 openrc 包 → 自动 apk add 装上,再走 OpenRC 路径
# - 大部分 LXC 容器没有 systemd,老脚本直接 systemctl 失败"systemctl: command not found"
HAS_SYSTEMD=0
HAS_OPENRC=0
IS_ALPINE=0
if [ -f /etc/alpine-release ]; then
    IS_ALPINE=1
elif [ -f /etc/os-release ] && grep -qE '^ID=alpine' /etc/os-release 2>/dev/null; then
    IS_ALPINE=1
fi
# Alpine 上 openrc 缺失就尝试自动装,失败不致命(下面还有 nohup 兜底)
if [ "$IS_ALPINE" = "1" ] && ! command -v rc-service >/dev/null 2>&1; then
    echo "[Init] Alpine detected without OpenRC, installing openrc..."
    if command -v apk >/dev/null 2>&1; then
        apk add --no-cache openrc 2>/dev/null || echo "[Init] apk add openrc failed, will fall back to nohup"
    fi
fi
# Alpine 优先 OpenRC;非 Alpine 仍然先看 systemd(主流发行版默认)
if [ "$IS_ALPINE" = "1" ]; then
    if command -v rc-service >/dev/null 2>&1; then HAS_OPENRC=1; fi
fi
if [ "$HAS_OPENRC" = "0" ] && command -v systemctl >/dev/null 2>&1 && [ -d /run/systemd/system ]; then
    HAS_SYSTEMD=1
fi
if [ "$HAS_SYSTEMD" = "0" ] && [ "$HAS_OPENRC" = "0" ] && command -v rc-service >/dev/null 2>&1; then
    HAS_OPENRC=1
fi
echo "Init system: $([ "$HAS_OPENRC" = 1 ] && echo openrc || ([ "$HAS_SYSTEMD" = 1 ] && echo systemd || echo none))$([ "$IS_ALPINE" = 1 ] && echo " (Alpine)")"

# Step 1: Stop existing service if running
echo "[1/6] Stopping existing service (if any)..."
AGENT_WAS_ACTIVE=0
AGENT_WAS_ENABLED=0
GUARD_WAS_ACTIVE=0
GUARD_WAS_ENABLED=0
openrc_service_enabled() {
	rc-update show default 2>/dev/null | awk -v service="$1" '$1 == service { found=1 } END { exit !found }'
}
if [ "$HAS_SYSTEMD" = "1" ]; then
	if systemctl is-active --quiet mmw-agent; then AGENT_WAS_ACTIVE=1; fi
	if systemctl is-enabled --quiet mmw-agent; then AGENT_WAS_ENABLED=1; fi
	if systemctl is-active --quiet mmwx-guard-agent; then GUARD_WAS_ACTIVE=1; fi
	if systemctl is-enabled --quiet mmwx-guard-agent; then GUARD_WAS_ENABLED=1; fi
    systemctl stop mmw-agent 2>/dev/null || true
elif [ "$HAS_OPENRC" = "1" ]; then
	if rc-service mmw-agent status >/dev/null 2>&1; then AGENT_WAS_ACTIVE=1; fi
	if openrc_service_enabled mmw-agent; then AGENT_WAS_ENABLED=1; fi
	if rc-service mmwx-guard-agent status >/dev/null 2>&1; then GUARD_WAS_ACTIVE=1; fi
	if openrc_service_enabled mmwx-guard-agent; then GUARD_WAS_ENABLED=1; fi
    rc-service mmw-agent stop 2>/dev/null || true
else
    # nohup 兜底:杀掉现有 mmw-agent 进程
	if pgrep -f '^/usr/local/bin/mmw-agent( |$)' >/dev/null 2>&1 || \
	   pgrep -f '^/bin/sh /usr/local/bin/mmw-agent-supervisor.sh$' >/dev/null 2>&1; then AGENT_WAS_ACTIVE=1; fi
	if pgrep -f '^/usr/local/bin/mmwx-guardd-agent( |$)' >/dev/null 2>&1 || \
	   pgrep -f '^/bin/sh /usr/local/bin/mmwx-guard-agent-supervisor.sh$' >/dev/null 2>&1; then GUARD_WAS_ACTIVE=1; fi
    pkill -f /usr/local/bin/mmw-agent 2>/dev/null || true
    sleep 1
fi

restore_previous_agent() {
	[ "$AGENT_WAS_ACTIVE" = "1" ] || return 0
	if [ "$HAS_SYSTEMD" = "1" ]; then
		systemctl enable mmw-agent >/dev/null 2>&1 || true
		systemctl restart mmw-agent >/dev/null 2>&1 || true
	elif [ "$HAS_OPENRC" = "1" ]; then
		rc-update add mmw-agent default >/dev/null 2>&1 || true
		rc-service mmw-agent restart >/dev/null 2>&1 || true
	elif [ -x /usr/local/bin/mmw-agent-supervisor.sh ]; then
		pkill -f '^/bin/sh /usr/local/bin/mmw-agent-supervisor.sh$' 2>/dev/null || true
		nohup /usr/local/bin/mmw-agent-supervisor.sh >/dev/null 2>&1 </dev/null &
	fi
}

CONFIG_PATH=/etc/mmw-agent/config.yaml
CONFIG_BACKUP=/var/lib/mmw-agent/config.yaml.install-backup
CONFIG_HAD_OLD=0
CONFIG_TOUCHED=0
SERVICE_STATE_CAPTURED=0
restore_previous_config() {
	[ "$CONFIG_TOUCHED" = "1" ] || return 0
	if [ "$CONFIG_HAD_OLD" = "1" ] && [ -f "$CONFIG_BACKUP" ]; then
		cp -p "$CONFIG_BACKUP" "$CONFIG_PATH" || true
	else
		rm -f "$CONFIG_PATH"
	fi
}
early_install_cleanup() {
	early_rc=$?
	trap - EXIT
	set +e
	restore_previous_config
	if declare -F restore_managed_services >/dev/null 2>&1 && [ "$SERVICE_STATE_CAPTURED" = "1" ]; then
		restore_managed_services
	else
		restore_previous_agent
	fi
	if [ -n "${STAGING_DIR:-}" ]; then rm -rf "$STAGING_DIR"; fi
	rm -f "$CONFIG_BACKUP"
	exit "$early_rc"
}
# Install this immediately after stopping the old Agent. Failures while creating
# config/service files or bootstrapping curl/OpenSSL must not leave it offline.
trap early_install_cleanup EXIT

# Service definitions and boot registration are part of the install
# transaction too. A failed reinstall must put custom/legacy service files and
# their exact enabled/active state back; a failed fresh install must not leave
# broken boot entries behind.
mkdir -p /var/lib/mmw-agent-install
chmod 0700 /var/lib/mmw-agent-install
STAGING_DIR=$(mktemp -d /var/lib/mmw-agent-install/in-process.XXXXXX)
chmod 0700 "$STAGING_DIR"
backup_managed_file() { # backup_managed_file <path> <key>
	if [ -e "$1" ] || [ -L "$1" ]; then
		cp -a "$1" "$STAGING_DIR/service-$2.old"
		: > "$STAGING_DIR/service-$2.present"
	fi
}
restore_managed_file() { # restore_managed_file <path> <key>
	if [ -f "$STAGING_DIR/service-$2.present" ]; then
		rm -f "$1"
		cp -a "$STAGING_DIR/service-$2.old" "$1"
	else
		rm -f "$1"
	fi
}
backup_managed_file /etc/systemd/system/mmw-agent.service systemd-agent
backup_managed_file /etc/systemd/system/mmwx-guard-agent.service systemd-guard
backup_managed_file /etc/init.d/mmw-agent openrc-agent
backup_managed_file /etc/init.d/mmwx-guard-agent openrc-guard
backup_managed_file /usr/local/bin/mmw-agent-supervisor.sh direct-agent
backup_managed_file /usr/local/bin/mmwx-guard-agent-supervisor.sh direct-guard
backup_managed_file /etc/rc.local rc-local
SERVICE_STATE_CAPTURED=1

restore_managed_services() {
	set +e
	if [ "$HAS_SYSTEMD" = "1" ]; then
		systemctl stop mmw-agent mmwx-guard-agent >/dev/null 2>&1
		# Remove registrations while the just-written unit files still exist;
		# systemctl cannot reliably disable a dangling unit after fresh rollback.
		systemctl disable mmw-agent mmwx-guard-agent >/dev/null 2>&1
	elif [ "$HAS_OPENRC" = "1" ]; then
		rc-service mmw-agent stop >/dev/null 2>&1
		rc-service mmwx-guard-agent stop >/dev/null 2>&1
		rc-update del mmw-agent default >/dev/null 2>&1
		rc-update del mmwx-guard-agent default >/dev/null 2>&1
	else
		if [ -f /run/mmwx-guard-agent-supervisor.pid ]; then
			old_guard_supervisor=$(cat /run/mmwx-guard-agent-supervisor.pid 2>/dev/null || true)
			case "$old_guard_supervisor" in *[!0-9]*|"") ;; *) kill "$old_guard_supervisor" 2>/dev/null ;; esac
		fi
		pkill -f '^/bin/sh /usr/local/bin/mmw-agent-supervisor.sh$' 2>/dev/null
		pkill -f '^/usr/local/bin/mmw-agent( |$)' 2>/dev/null
		pkill -f '^/bin/sh /usr/local/bin/mmwx-guard-agent-supervisor.sh$' 2>/dev/null
		pkill -f '^/usr/local/bin/mmwx-guardd-agent( |$)' 2>/dev/null
		rm -f /run/mmwx-guard-agent-supervisor.pid /run/mmwx-guard-agent/guard.sock
	fi
	restore_managed_file /etc/systemd/system/mmw-agent.service systemd-agent
	restore_managed_file /etc/systemd/system/mmwx-guard-agent.service systemd-guard
	restore_managed_file /etc/init.d/mmw-agent openrc-agent
	restore_managed_file /etc/init.d/mmwx-guard-agent openrc-guard
	restore_managed_file /usr/local/bin/mmw-agent-supervisor.sh direct-agent
	restore_managed_file /usr/local/bin/mmwx-guard-agent-supervisor.sh direct-guard
	restore_managed_file /etc/rc.local rc-local
	if [ "$HAS_SYSTEMD" = "1" ]; then
		systemctl daemon-reload >/dev/null 2>&1
		if [ "$GUARD_WAS_ENABLED" = "1" ]; then systemctl enable mmwx-guard-agent >/dev/null 2>&1
		fi
		if [ "$AGENT_WAS_ENABLED" = "1" ]; then systemctl enable mmw-agent >/dev/null 2>&1
		fi
		if [ "$GUARD_WAS_ACTIVE" = "1" ]; then systemctl start mmwx-guard-agent >/dev/null 2>&1; fi
		if [ "$AGENT_WAS_ACTIVE" = "1" ]; then systemctl start mmw-agent >/dev/null 2>&1; fi
	elif [ "$HAS_OPENRC" = "1" ]; then
		if [ "$GUARD_WAS_ENABLED" = "1" ]; then rc-update add mmwx-guard-agent default >/dev/null 2>&1
		fi
		if [ "$AGENT_WAS_ENABLED" = "1" ]; then rc-update add mmw-agent default >/dev/null 2>&1
		fi
		if [ "$GUARD_WAS_ACTIVE" = "1" ]; then rc-service mmwx-guard-agent start >/dev/null 2>&1; fi
		if [ "$AGENT_WAS_ACTIVE" = "1" ]; then rc-service mmw-agent start >/dev/null 2>&1; fi
	else
		if [ "$GUARD_WAS_ACTIVE" = "1" ] && [ -x /usr/local/bin/mmwx-guard-agent-supervisor.sh ]; then
			nohup /usr/local/bin/mmwx-guard-agent-supervisor.sh >/var/log/mmwx-guard-agent.log 2>&1 </dev/null &
			echo $! > /run/mmwx-guard-agent-supervisor.pid
		fi
		if [ "$AGENT_WAS_ACTIVE" = "1" ] && [ -x /usr/local/bin/mmw-agent-supervisor.sh ]; then
			nohup /usr/local/bin/mmw-agent-supervisor.sh >/dev/null 2>&1 </dev/null &
		fi
	fi
}

# 被 systemd/OpenRC 中断的旧 bootstrap 可能在小容量 /tmp 中留下多个
# 15MB Guard 副本，最终让全新安装连第一个 Agent 文件都写不进去。此时
# Agent 已经停止，因此这些精确前缀的目录都不可能仍被有效 bootstrap 使用。
find /tmp -mindepth 1 -maxdepth 1 -type d \( -name 'mmwx-guard-bootstrap' -o -name 'mmwx-guard-bootstrap-*' \) -exec rm -rf -- {} + 2>/dev/null || true
rm -f /tmp/mmw-agent /tmp/mmw-agent-new /tmp/mmw-agent-new.sig \
    /tmp/mmw-agent-new.manifest /tmp/mmwx-guardd /tmp/mmwx-guardd.sig \
    /tmp/mmwx-guardd-new /tmp/mmwx-guardd-new.sig

# Step 2: Create config directory first
echo ""
echo "[2/6] Creating configuration..."
mkdir -p /etc/mmw-agent
mkdir -p /var/lib/mmw-agent
mkdir -p /var/lib/mmwx-guard

# 端口探测:从 LISTEN_PORT(或默认 23889)起,被占用就 +1,最多试 20 次。
# 用 ss 看任意接口的 LISTEN socket,避免 agent 启动后 bind 失败造成"WS 活/HTTP 死"的死锁状态。
REQUESTED_PORT="${LISTEN_PORT:-23889}"
ACTUAL_PORT=""
for i in $(seq 0 19); do
    TRY_PORT=$((REQUESTED_PORT + i))
    if ss -tln 2>/dev/null | awk '{print $4}' | grep -qE "[:.]${TRY_PORT}\$"; then
        echo "  端口 ${TRY_PORT} 已被占用,尝试下一个..."
        continue
    fi
    ACTUAL_PORT="$TRY_PORT"
    break
done
if [ -z "$ACTUAL_PORT" ]; then
    echo "ERROR: 从 ${REQUESTED_PORT} 起的 20 个端口全部被占用,安装中止" >&2
    exit 1
fi
if [ "$ACTUAL_PORT" != "$REQUESTED_PORT" ]; then
    echo "⚠ 端口 ${REQUESTED_PORT} 被占用,自动改用 ${ACTUAL_PORT}"
fi
LISTEN_PORT="$ACTUAL_PORT"

rm -f "$CONFIG_BACKUP"
if [ -f "$CONFIG_PATH" ]; then
	cp -p "$CONFIG_PATH" "$CONFIG_BACKUP"
	CONFIG_HAD_OLD=1
fi
CONFIG_TOUCHED=1
cat > "$CONFIG_PATH" << EOF
# MMWX Remote Server Configuration
# Generated by install script

mode: remote
master_url: ${MASTER_URL}
recovery_url: ${RECOVERY_URL}
token: ${TOKEN}
connection_mode: websocket
xray_mode: ${XRAY_MODE}
steal_mode: $([ "$AUTO_STEAL_SELF" = "1" ] && echo "tunnel" || echo "")
master_public_key: ${MASTER_PUBLIC_KEY}
listen_port: "${LISTEN_PORT}"
EOF

echo "Configuration saved to /etc/mmw-agent/config.yaml"

# Step 3: 创建 service 文件 — 按检测到的 init 系统选不同写法
echo ""
echo "[3/6] Creating service..."

if [ "$HAS_SYSTEMD" = "1" ]; then
    cat > /etc/systemd/system/mmw-agent.service << EOF
[Unit]
Description=MMW Agent Remote Server
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/mmw-agent -c /etc/mmw-agent/config.yaml
Restart=always
RestartSec=5
WorkingDirectory=/var/lib/mmw-agent

[Install]
WantedBy=multi-user.target
EOF
    systemctl daemon-reload
elif [ "$HAS_OPENRC" = "1" ]; then
    cat > /etc/init.d/mmw-agent << 'EOF'
#!/sbin/openrc-run
name="mmw-agent"
description="MMW Agent Remote Server"
command="/usr/local/bin/mmw-agent"
command_args="-c /etc/mmw-agent/config.yaml"
# supervise-daemon keeps the agent alive when it intentionally exits after a
# listen_port / xray_mode switch. command_background/start-stop-daemon alone
# only daemonizes once and does not respawn the process.
supervisor="supervise-daemon"
respawn_delay=3
respawn_max=0
# 日志由 agent 自身写文件并轮转(/var/log/mmw-agent/mmw-agent.log),不再用 output_log 重复落地(避免无轮转爆盘)
depend() { need net; }
EOF
    chmod +x /etc/init.d/mmw-agent
else
    # 无 init 系统(典型 LXC 容器):写一个 supervisor 脚本,失败自动重启,放后台跑;同时塞进 rc.local 以便重启
    cat > /usr/local/bin/mmw-agent-supervisor.sh << 'EOF'
#!/bin/sh
# 自托管构建不依赖外部 Guard，不等待 guard socket。
while true; do
    # 日志由 agent 自身写文件并轮转(/var/log/mmw-agent/mmw-agent.log);这里输出走 stdout(由 rc.local 的 nohup 接管)
    /usr/local/bin/mmw-agent -c /etc/mmw-agent/config.yaml
    echo "[supervisor] mmw-agent exited, restarting in 5s..."
    sleep 5
done
EOF
    chmod +x /usr/local/bin/mmw-agent-supervisor.sh

    # 写入 rc.local 实现重启自启动(若文件不存在就建一个)
    if [ ! -f /etc/rc.local ]; then
        echo "#!/bin/sh" > /etc/rc.local
        echo "exit 0" >> /etc/rc.local
        chmod +x /etc/rc.local
    fi
    if ! grep -q "mmw-agent-supervisor.sh" /etc/rc.local; then
        sed -i '/^exit 0/i nohup /usr/local/bin/mmw-agent-supervisor.sh >/dev/null 2>&1 \&' /etc/rc.local
    fi
fi

# Step 4: Download and install binary only (without starting)
echo ""
echo "[4/6] Downloading MMWX binary..."

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH_NAME="amd64"
        ;;
    aarch64|arm64)
        ARCH_NAME="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Agent 二进制的下载来源,按顺序尝试,任一成功即停:
#   1) 主控本机 —— 管理员把二进制放在 <数据目录>/agent-bin/ 时才加入(带安装 token 鉴权)
#   2) GitHub Release —— tag 为 mmwx-agent(或 mmwx-agent-v* 写法)的那条
#   3) 自定义镜像 —— 主控设了 MMWX_AGENT_DOWNLOAD_BASE 才加入
MIRRORS=()

if [ "__MMWX_AGENT_LOCAL__" = "1" ]; then
    MIRRORS+=("${MASTER_URL}/api/remote/agent-binary?token=${TOKEN}&arch=${ARCH_NAME}")
fi

GH_REPO="__MMWX_AGENT_REPO__"
GH_PROXY="__MMWX_GH_PROXY__"
if [ -n "$GH_REPO" ] && command -v curl >/dev/null 2>&1; then
    echo "Resolving latest Agent release from GitHub (${GH_REPO}) ..."
    GH_TAG=$(curl -fsSL --connect-timeout 10 --max-time 30 \
        -H "Accept: application/vnd.github+json" \
        "https://api.github.com/repos/${GH_REPO}/releases?per_page=100" 2>/dev/null \
        | grep -o '"tag_name":[[:space:]]*"[^"]*"' \
        | sed 's/.*"tag_name":[[:space:]]*"//; s/"$//' \
        | grep -E "^(mmwx-agent|mmwx-agent-v.*)$" | head -n 1 || true)
    if [ -n "$GH_TAG" ]; then
        echo "  → $GH_TAG"
        for GH_ASSET in "mmwx-agent-linux-${ARCH_NAME}" "mmw-agent-linux-${ARCH_NAME}"; do
            GH_URL="https://github.com/${GH_REPO}/releases/download/${GH_TAG}/${GH_ASSET}"
            if [ -n "$GH_PROXY" ]; then
                MIRRORS+=("${GH_PROXY%/}/${GH_URL}")
            fi
            MIRRORS+=("$GH_URL")
        done
    else
        echo "  → 没查到 Agent 的 Release(tag 应为 mmwx-agent;可能 API 限流),跳过 GitHub"
    fi
fi

AGENT_MIRROR_BASE="__MMWX_AGENT_MIRROR__"
if [ -n "$AGENT_MIRROR_BASE" ]; then
    MIRRORS+=("${AGENT_MIRROR_BASE}/mmwx-agent-linux-${ARCH_NAME}")
fi

if [ ${#MIRRORS[@]} -eq 0 ]; then
    echo "ERROR: 没有可用的 Agent 下载来源。" >&2
    echo "       要么把 mmwx-agent-linux-${ARCH_NAME} 放到主控的 <数据目录>/agent-bin/," >&2
    echo "       要么在 GitHub 仓库发一个 tag 为 mmwx-agent 的 Release。" >&2
    exit 1
fi

# Download binary — 优先用 curl(更普遍),没有就用 wget;两者都没就按发行版包管理器装一个,
# 杜绝 "wget: command not found" 噪声 / "ERROR: 都没装" 卡死。
ensure_downloader() {
    if command -v curl >/dev/null 2>&1; then return 0; fi
    if command -v wget >/dev/null 2>&1; then return 0; fi
    echo "未检测到 curl/wget,尝试自动安装 curl..."
    if command -v apt-get >/dev/null 2>&1; then
        apt-get update -qq >/dev/null 2>&1 || true
        DEBIAN_FRONTEND=noninteractive apt-get install -y curl
    elif command -v dnf >/dev/null 2>&1; then
        dnf install -y curl
    elif command -v yum >/dev/null 2>&1; then
        yum install -y curl
    elif command -v apk >/dev/null 2>&1; then
        apk add --no-cache curl
    elif command -v pacman >/dev/null 2>&1; then
        pacman -Sy --noconfirm curl
    elif command -v zypper >/dev/null 2>&1; then
        zypper -n install curl
    else
        echo "ERROR: 无法识别系统包管理器,请手动安装 curl 或 wget 后重试" >&2
        return 1
    fi
}
ensure_downloader || exit 1
# 安装过程先把候选二进制放进私有事务目录,全部就绪后再原子替换。
AGENT_NEW="$STAGING_DIR/mmw-agent"
INSTALL_SUCCEEDED=0
TRANSACTION_ROLLED_BACK=0
cleanup_install_assets() {
	cleanup_rc=$?
	trap - EXIT
	set +e
	if [ "$INSTALL_SUCCEEDED" != "1" ]; then
		if declare -F rollback_install_transaction >/dev/null 2>&1 && [ "$TRANSACTION_ROLLED_BACK" != "1" ]; then
			rollback_install_transaction
			set +e
		fi
		restore_previous_config
		if [ "$SERVICE_STATE_CAPTURED" = "1" ]; then restore_managed_services
		else restore_previous_agent; fi
	fi
    rm -rf "$STAGING_DIR"
    rm -f /usr/local/bin/mmw-agent.new
	rm -f "$CONFIG_BACKUP"
	exit "$cleanup_rc"
}
trap cleanup_install_assets EXIT

download_file() { # download_file <url> <output> <max-seconds>
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL --connect-timeout 10 --max-time "$3" -o "$2" "$1"
    else
        wget -q --connect-timeout=10 --read-timeout="$3" -O "$2" "$1"
    fi
}
# 从主控下载 Agent 并原子安装。
download_ok=0
for agent_url in "${MIRRORS[@]}"; do
	echo "Downloading Agent from $agent_url ..."
	if download_file "$agent_url" "$AGENT_NEW" 180; then download_ok=1; break; fi
	echo "  → 该镜像失败,尝试下一个..."
done
if [ "$download_ok" != "1" ]; then
	echo "ERROR: Agent 下载失败,上面所有来源都没成功。可以:" >&2
	echo "       1) 检查本机能不能访问 GitHub(国内网络给主控设 MMWX_GH_PROXY 走加速)" >&2
	echo "       2) 或把 mmwx-agent-linux-${ARCH_NAME} 放到主控的 <数据目录>/agent-bin/ 下,由主控直接分发" >&2
	exit 1
fi
chmod 0755 "$AGENT_NEW"
AGENT_BIN=/usr/local/bin/mmw-agent
if [ -f "$AGENT_BIN" ]; then cp -p "$AGENT_BIN" "$STAGING_DIR/agent.old"; fi
if ! install -m 0755 "$AGENT_NEW" "$AGENT_BIN.new" || ! mv -f "$AGENT_BIN.new" "$AGENT_BIN"; then
	[ -f "$STAGING_DIR/agent.old" ] && cp -p "$STAGING_DIR/agent.old" "$AGENT_BIN"
	echo "ERROR: Agent 安装失败,已恢复原文件" >&2
	exit 1
fi
echo "Agent installed to /usr/local/bin/mmw-agent (offline / no-Guard)."

# Step 5: 启用并启动 service
echo ""
echo "[5/6] Starting service..."
if [ "$HAS_SYSTEMD" = "1" ]; then
    systemctl enable mmw-agent
    systemctl start mmw-agent
elif [ "$HAS_OPENRC" = "1" ]; then
    # rc-update 在 LXC 容器里没初始化 runlevel 时会报错,失败不致命(set -e 兜底)
    rc-update add mmw-agent default 2>/dev/null || echo "  ⚠ rc-update add 失败(常见于 LXC 容器,不影响当前会话启动)"
    rc-service mmw-agent start
else
    nohup /usr/local/bin/mmw-agent-supervisor.sh >/dev/null 2>&1 &
    echo "Started via nohup (PID=$!); 安装重启后通过 /etc/rc.local 自启动"
fi

# Wait a moment for service to start
sleep 3

agent_start_ok=0
if [ "$HAS_SYSTEMD" = "1" ]; then
	if systemctl is-active --quiet mmw-agent; then agent_start_ok=1; fi
elif [ "$HAS_OPENRC" = "1" ]; then
	if rc-service mmw-agent status >/dev/null 2>&1; then agent_start_ok=1; fi
elif pgrep -f '^/usr/local/bin/mmw-agent( |$)' >/dev/null 2>&1; then
	agent_start_ok=1
fi
if [ "$agent_start_ok" != "1" ]; then
	echo "ERROR: Agent 启动验证失败,正在恢复原安装" >&2
	exit 1
fi
# The Agent transaction is committed here. Xray/Nginx setup below is ancillary:
# a later provisioning error is reported to the caller but must not replace a
# healthy Agent with the pre-install version.
INSTALL_SUCCEEDED=1

# Step 6: Verify installation
echo ""
echo "[6/6] Verifying installation..."

echo ""
echo "=========================================="
echo "  Installation Complete!"
echo "=========================================="
echo ""
echo "Service status:"
if [ "$HAS_SYSTEMD" = "1" ]; then
    systemctl status mmw-agent --no-pager -l 2>/dev/null | head -15 || echo "Service started"
elif [ "$HAS_OPENRC" = "1" ]; then
    rc-service mmw-agent status 2>/dev/null || echo "Service started"
else
    pgrep -af /usr/local/bin/mmw-agent | head -5 || echo "Process not found in pgrep, check /var/log/mmw-agent.log"
fi
echo ""
echo "To check status:"
if [ "$HAS_SYSTEMD" = "1" ]; then
    echo "  systemctl status mmw-agent"
elif [ "$HAS_OPENRC" = "1" ]; then
    echo "  rc-service mmw-agent status"
else
    echo "  tail -f /var/log/mmw-agent.log  # 或: pgrep -af mmw-agent"
fi
echo ""
echo "To view logs:"
echo "  journalctl -u mmw-agent -f"
echo ""

# Auto-install Xray (unless embedded mode)
if [ "$XRAY_MODE" != "embedded" ]; then
    XRAY_INSTALLED=0
    if command -v xray >/dev/null 2>&1 || [ -x /usr/local/bin/xray ] || [ -x /usr/bin/xray ] || [ -x /opt/xray/xray ]; then
        XRAY_INSTALLED=1
    fi

    if [ "$XRAY_INSTALLED" = "1" ]; then
        echo "[Auto] Xray already installed, skip."
    else
        echo "[Auto] Installing Xray..."
        # 先落一份占位配置再装:Xray-install 装完会立刻 systemctl enable --now xray,
        # 而这一刻 mmw-agent 往往还没连上主控、没来得及写出真实 config.json,
        # xray 便会以 "failed to load config files ... no such file or directory" 启动失败,
        # 在全新安装的最后留下一条刺眼的红色报错(功能其实不受影响,agent 随后会写配置并重启它)。
        # 占位配置只保证 xray 能起来,agent 一旦下发真实配置就会整份覆盖。
        # 只在文件不存在时创建 —— 重装/已有配置的机器绝不能被覆盖。
        mkdir -p /usr/local/etc/xray
        if [ ! -f /usr/local/etc/xray/config.json ]; then
            cat > /usr/local/etc/xray/config.json <<'XRAYPLACEHOLDERCFG'
{
  "log": { "loglevel": "warning" },
  "inbounds": [],
  "outbounds": [ { "protocol": "freedom", "tag": "direct" } ]
}
XRAYPLACEHOLDERCFG
        fi
        bash -c "$(curl -L https://github.com/XTLS/Xray-install/raw/main/install-release.sh)" @ install
    fi
fi

if [ "$AUTO_STEAL_SELF" = "1" ]; then
    echo "=========================================="
    echo "  Auto Install: Nginx"
    echo "=========================================="
    echo ""

    if [ -x /usr/local/nginx/sbin/nginx ]; then
        NGINX_INFO=$(/usr/local/nginx/sbin/nginx -V 2>&1)
        NGINX_VERSION=$(echo "$NGINX_INFO" | sed -n 's#^nginx version: nginx/\([0-9.]*\).*$#\1#p')
        NGINX_OLDEST=$(printf '%s\n%s\n' "1.25.1" "$NGINX_VERSION" | sort -V | head -n 1)
        if [ -z "$NGINX_VERSION" ] || [ "$NGINX_OLDEST" != "1.25.1" ] \
            || ! echo "$NGINX_INFO" | grep -q -- "--with-http_ssl_module" \
            || ! echo "$NGINX_INFO" | grep -q -- "--with-http_v2_module" \
            || ! echo "$NGINX_INFO" | grep -q -- "--with-http_v3_module" \
            || ! echo "$NGINX_INFO" | grep -q -- "--with-http_realip_module" \
            || ! echo "$NGINX_INFO" | grep -q -- "--with-stream_ssl_module" \
            || ! echo "$NGINX_INFO" | grep -q -- "--with-stream_ssl_preread_module"; then
            echo "ERROR: 现有 Nginx 不兼容MEO（要求 >= 1.25.1，并启用 HTTP/2、HTTP/3、Stream SSL 模块）。" >&2
            echo "请先卸载现有 Nginx，再重新执行本安装脚本，由MEO安装受支持的版本。" >&2
            exit 1
        fi
        echo "[Auto] Compatible Nginx already installed, skip."
    elif command -v nginx >/dev/null 2>&1; then
        echo "ERROR: 检测到非MEO管理的 Nginx: $(command -v nginx)" >&2
        echo "为避免覆盖用户配置，安装已停止。请先卸载系统 Nginx，再重新执行本安装脚本。" >&2
        exit 1
    else
        echo "[Auto] Installing Nginx..."
        curl -fsSL "${MASTER_URL}/api/remote/nginx-script?token=${TOKEN}&action=install" | bash
    fi
    echo ""
    echo "Auto install complete (front service: ${FRONT_SERVICE}, xray mode: ${XRAY_MODE})"
fi
echo ""
`

	// 注入更新 CDN base(开关开→写死域名,开关关→空);空则脚本里 CDN_BASE 为空、MIRRORS 只用 GitHub / gh-proxy
	script = strings.ReplaceAll(script, "__MMWX_AGENT_MIRROR__", AgentDownloadBase())
	script = strings.ReplaceAll(script, "__MMWX_AGENT_REPO__", AgentGitHubRepo())
	script = strings.ReplaceAll(script, "__MMWX_GH_PROXY__", GitHubProxyBase())
	script = strings.ReplaceAll(script, "__MMWX_AGENT_LOCAL__", map[bool]string{true: "1", false: "0"}[HasLocalAgentBinary()])

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=install.sh")
	w.Write([]byte(script))
}
