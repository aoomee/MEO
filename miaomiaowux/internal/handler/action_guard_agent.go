package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"miaomiaowux/internal/guardclient"
	"miaomiaowux/internal/license"
)

func (h *RemoteManageHandler) authorizeManagedNodeAction(ctx context.Context, browserRequest *http.Request, serverID int64, payloadHash string) (license.ActionGrantDelivery, string, error) {
	if h.actionGuard == nil || !h.actionGuard.RequiresMaster() {
		return license.ActionGrantDelivery{}, "", nil
	}
	if !h.actionGuard.Enabled() {
		return license.ActionGrantDelivery{}, "", errors.New("Action Guard 不可用，已拒绝创建节点")
	}
	server, err := h.repo.GetRemoteServer(ctx, serverID)
	if err != nil || server == nil {
		return license.ActionGrantDelivery{}, "", errors.New("找不到节点所属服务器")
	}
	serverHash := serverTokenHash(server.Token)
	var agent guardclient.Attestation
	if h.actionGuard.RequiresAgent() {
		// Repair a license-identity mismatch before consuming the browser's
		// one-time frontend challenge. Retrying AuthorizeMaster after a failed
		// grant would reuse that challenge and surface the misleading
		// "frontend challenge invalid or expired" error.
		if repairErr := h.prepareManagedNodeAuthorityForServer(ctx, serverID, serverHash); repairErr != nil {
			return license.ActionGrantDelivery{}, "", fmt.Errorf("服务器许可证租约未就绪(server_id=%d, token_hash=%s): %w",
				serverID, shortServerHash(serverHash), repairErr)
		}
		request := guardclient.AttestationRequest{Action: actionNodeCreate, PayloadHash: payloadHash, ServerHash: serverHash}
		body, _ := json.Marshal(request)
		// 严格等级只允许主控经 Agent 已建立的 WS/RPC 取证明。Agent 不监听公网，
		// 这里也绝不降级为反向 HTTP。
		response, rpcErr := h.forwardGuardOverWS(ctx, serverID, http.MethodPost, "/api/child/action-guard/attest", body)
		if rpcErr != nil {
			return license.ActionGrantDelivery{}, "", errors.New("Agent Guard WS 联合验证失败: " + rpcErr.Error())
		}
		if json.Unmarshal(response, &agent) != nil || agent.ID == "" {
			return license.ActionGrantDelivery{}, "", errors.New("Agent Guard 返回无效证明")
		}
	}
	delivery, err := authorizeMasterOnceWithSlotRepair(ctx, serverID, serverHash,
		func() (license.ActionGrantDelivery, error) {
			return h.actionGuard.AuthorizeMaster(ctx, browserRequest, actionNodeCreate, payloadHash, serverHash, agent)
		}, h.ensureAuthoritativeSlot)
	if err != nil {
		return license.ActionGrantDelivery{}, "", err
	}
	if h.actionGuard.RequiresAgent() {
		consumeBody, _ := json.Marshal(guardclient.ConsumeRequest{
			Grant: delivery.Grant, Action: actionNodeCreate, PayloadHash: payloadHash,
			ServerHash: serverHash, LicenseServerURL: delivery.LicenseServerURL,
		})
		if _, err := h.forwardGuardOverWS(ctx, serverID, http.MethodPost, "/api/child/action-guard/consume", consumeBody); err != nil {
			return license.ActionGrantDelivery{}, "", errors.New("Agent Guard WS 消费授权失败: " + err.Error())
		}
	}
	if err := h.actionGuard.ConsumeMaster(ctx, delivery, actionNodeCreate, payloadHash, serverHash); err != nil {
		return license.ActionGrantDelivery{}, "", errors.New("主控 Guard 消费授权失败: " + err.Error())
	}
	return delivery, serverHash, nil
}

// prepareManagedNodeAuthority is also used by the federation owner before it
// issues a one-time browser challenge. Existing shares deliberately survive a
// license change; only their real Agent slot must move to the new authority.
func (h *RemoteManageHandler) prepareManagedNodeAuthority(ctx context.Context, serverID int64) error {
	if h.actionGuard == nil || !h.actionGuard.RequiresAgent() {
		return nil
	}
	server, err := h.repo.GetRemoteServer(ctx, serverID)
	if err != nil || server == nil {
		return errors.New("找不到共享服务器对应的真实服务器")
	}
	return h.prepareManagedNodeAuthorityForServer(ctx, serverID, serverTokenHash(server.Token))
}

func (h *RemoteManageHandler) prepareManagedNodeAuthorityForServer(ctx context.Context, serverID int64, serverHash string) error {
	if h.wsHandler == nil {
		return errors.New("主控 WS 服务不可用，请重启 Agent 后重试")
	}
	conn, ok := h.wsHandler.GetConnectionByServerID(serverID)
	if !ok {
		return errors.New("Agent 未连接，无法重新签发授权槽位")
	}
	conn.mu.Lock()
	identityCapable := conn.Capabilities.LeaseIdentity
	reportedIdentity := conn.AgentLeaseIdentity
	leaseReady := conn.LeaseReady
	conn.mu.Unlock()
	currentIdentity := ""
	if h.actionGuard != nil && h.actionGuard.license != nil {
		currentIdentity = h.actionGuard.license.AgentLeaseIdentity()
	}
	identityMatches := !identityCapable || (currentIdentity != "" && reportedIdentity == currentIdentity)
	if leaseReady && identityMatches {
		return nil
	}
	return h.ensureAuthoritativeSlot(ctx, serverID, serverHash, false)
}

type authorizeMasterFunc func() (license.ActionGrantDelivery, error)
type authoritativeSlotRepairFunc func(context.Context, int64, string, bool) error

var errAuthoritativeSlotRepaired = errors.New("authoritative slot repaired; retry with a new challenge")

type authoritativeSlotRepairedError struct {
	serverID  int64
	tokenHash string
	cause     error
}

func (e *authoritativeSlotRepairedError) Error() string {
	return fmt.Sprintf(
		"服务器许可证租约已自动修复(server_id=%d, token_hash=%s)，请重试本次创建操作以使用新的安全 challenge；首次 ActionGrant 错误: %v",
		e.serverID, shortServerHash(e.tokenHash), e.cause)
}

func (e *authoritativeSlotRepairedError) Unwrap() []error {
	return []error{errAuthoritativeSlotRepaired, e.cause}
}

// IsAuthoritativeSlotRepaired lets HTTP handlers return a stable
// action_guard_retry code without parsing a localized error message.
func IsAuthoritativeSlotRepaired(err error) bool {
	return errors.Is(err, errAuthoritativeSlotRepaired)
}

// authorizeMasterOnceWithSlotRepair deliberately invokes authorize exactly
// once. A frontend challenge is single-use, so an authoritative-slot failure
// may trigger lease repair, but the repaired request must be retried by the
// browser with a newly signed challenge.
func authorizeMasterOnceWithSlotRepair(
	ctx context.Context,
	serverID int64,
	serverHash string,
	authorize authorizeMasterFunc,
	repair authoritativeSlotRepairFunc,
) (license.ActionGrantDelivery, error) {
	delivery, err := authorize()
	if err == nil {
		return delivery, nil
	}
	if !license.IsAuthoritativeSlotError(err) {
		return license.ActionGrantDelivery{}, fmt.Errorf("创建节点授权失败(server_id=%d, token_hash=%s): %w",
			serverID, shortServerHash(serverHash), err)
	}

	// The license service has authoritative knowledge that the slot is not
	// active. Force replacement even when Agent Guard's local slot file still
	// looks healthy; local status cannot disprove a server-side stale slot.
	repairErr := repair(ctx, serverID, serverHash, true)
	if repairErr != nil {
		// Keep both errors in the unwrap chain: the first error explains why
		// repair started, while the second explains why automatic recovery failed.
		return license.ActionGrantDelivery{}, fmt.Errorf(
			"服务器授权槽位自动修复失败(server_id=%d, token_hash=%s): 首次 ActionGrant 错误: %w；强制重签错误: %w",
			serverID, shortServerHash(serverHash), err, repairErr)
	}
	return license.ActionGrantDelivery{}, &authoritativeSlotRepairedError{serverID: serverID, tokenHash: serverHash, cause: err}
}

type authoritativeSlotStatus struct {
	Matches bool `json:"matches"`
	Slot    struct {
		Authorized     bool   `json:"authorized"`
		ServerHash     string `json:"server_hash"`
		LicenseKeyHash string `json:"license_key_hash"`
		SlotID         int64  `json:"slot_id"`
		Generation     int64  `json:"generation"`
	} `json:"slot"`
}

type authoritativeSlotReader func(context.Context) (authoritativeSlotStatus, bool)

func (h *RemoteManageHandler) ensureAuthoritativeSlot(ctx context.Context, serverID int64, expectedHash string, forceReplacement bool) error {
	if h.wsHandler == nil {
		return errors.New("主控 WS 服务不可用，请重启 Agent 后重试")
	}
	conn, ok := h.wsHandler.GetConnectionByServerID(serverID)
	if !ok {
		return errors.New("Agent 未连接，无法重新签发授权槽位")
	}
	conn.mu.Lock()
	identityCapable := conn.Capabilities.LeaseIdentity
	conn.mu.Unlock()
	readSlot := func(queryCtx context.Context) (authoritativeSlotStatus, bool) {
		body, err := h.forwardGuardOverWS(queryCtx, serverID, http.MethodGet, "/api/child/action-guard/status", nil)
		if err != nil {
			return authoritativeSlotStatus{}, false
		}
		var status authoritativeSlotStatus
		if json.Unmarshal(body, &status) != nil {
			return authoritativeSlotStatus{}, false
		}
		return status, true
	}
	currentIdentity := ""
	if h.actionGuard != nil && h.actionGuard.license != nil {
		currentIdentity = h.actionGuard.license.AgentLeaseIdentity()
	}
	// Guard activation can legitimately take up to 25 seconds. Use the shared
	// serialized queue (including Retry-After/backoff and exact delivery ACK)
	// instead of issuing directly from the request handler, and bound the whole
	// repair + status verification to one user-visible operation.
	repairCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()
	return ensureAuthoritativeSlotWithOperations(repairCtx, expectedHash, currentIdentity, identityCapable, forceReplacement,
		readSlot, func() error { return h.wsHandler.QueueAgentLeaseAndWait(repairCtx, conn, true) },
		40*time.Second, 400*time.Millisecond)
}

func ensureAuthoritativeSlotWithOperations(
	ctx context.Context,
	expectedHash string,
	currentIdentity string,
	identityCapable bool,
	forceReplacement bool,
	readSlot authoritativeSlotReader,
	replace func() error,
	timeout time.Duration,
	pollInterval time.Duration,
) error {
	baseline, haveBaseline := readSlot(ctx)
	slotMatches := func(queryCtx context.Context, requireTransition bool) bool {
		status, readable := readSlot(queryCtx)
		if !readable || !status.Matches || !status.Slot.Authorized || status.Slot.ServerHash != expectedHash {
			return false
		}
		if identityCapable && (currentIdentity == "" || status.Slot.LicenseKeyHash != currentIdentity) {
			return false
		}
		if requireTransition {
			return authoritativeSlotTransitioned(haveBaseline, baseline.Slot.SlotID, baseline.Slot.Generation,
				status.Slot.SlotID, status.Slot.Generation)
		}
		return true
	}
	if !forceReplacement && identityCapable && slotMatches(ctx, false) {
		return nil
	}
	if err := replace(); err != nil {
		return fmt.Errorf("重新签发授权槽位失败: %w", err)
	}
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	// A forced repair was caused by a license-service rejection. In that case a
	// locally healthy but unchanged slot is not evidence of recovery: wait for a
	// new slot ID or generation. Legacy Guard versions require the same proof
	// because they cannot report license_key_hash.
	requireTransition := forceReplacement || !identityCapable
	for {
		if slotMatches(waitCtx, requireTransition) {
			return nil
		}
		select {
		case <-waitCtx.Done():
			return fmt.Errorf("新许可证租约未在修复期限内激活: %w；请检查 Agent Guard 日志中的 slot/identity 错误", waitCtx.Err())
		case <-ticker.C:
		}
	}
}

func authoritativeSlotTransitioned(haveBaseline bool, beforeID, beforeGeneration, afterID, afterGeneration int64) bool {
	if !haveBaseline {
		// replace() only returns after the exact forced lease delivery has been
		// acknowledged. If the pre-repair status RPC was unavailable, a later
		// fully-authorized non-zero slot is therefore sufficient evidence; making
		// a missing baseline impossible to satisfy would turn a transient RPC
		// failure into a guaranteed repair timeout.
		return afterID != 0 && afterGeneration > 0
	}
	return (afterID != 0 && afterID != beforeID) || afterGeneration > beforeGeneration
}

func shortServerHash(hash string) string {
	if len(hash) > 12 {
		return hash[:12]
	}
	return hash
}

func (h *RemoteManageHandler) forwardGuardOverWS(ctx context.Context, serverID int64, method, path string, body []byte) ([]byte, error) {
	response, ok, err := h.tryWSRPC(ctx, serverID, method, path, body)
	if !ok {
		return nil, errors.New("Agent 未连接 WS 或版本不支持 RPC，请先升级 Agent")
	}
	return response, err
}
