package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"miaomiaowux/internal/guardclient"
	"miaomiaowux/internal/storage"
	"miaomiaowux/internal/version"
)

const (
	federationGuardProtocolVersion = 1

	actionGuardServerIDHeader = "X-MMWX-Guard-Server-ID"

	codeFederatedChallengeRequired = "action_guard_federated_challenge_required"
	codeFederationGuardUnsupported = "federation_action_guard_unsupported"
	codeFederationGuardInvalid     = "federation_action_guard_invalid"
	codeFederationGuardFailed      = "federation_action_guard_failed"

	federationInboundRequestPath = "/api/admin/remote/inbounds"
	federationNodeRequestPath    = "/api/admin/nodes"
)

// FederationActionChallengeRequest is sent by the consumer master to the
// owner. RequestMethod/RequestPath are the browser mutation, not the child RPC
// path. Binding both prevents a proof obtained for a local node import from
// being replayed as a remote inbound mutation.
type FederationActionChallengeRequest struct {
	Action        string `json:"action"`
	PayloadHash   string `json:"payload_hash"`
	RequestMethod string `json:"request_method"`
	RequestPath   string `json:"request_path"`
}

type federationActionProof struct {
	Action              string `json:"action"`
	OriginalPayloadHash string `json:"original_payload_hash"`
	RequestMethod       string `json:"request_method"`
	RequestPath         string `json:"request_path"`
	ChallengeID         string `json:"challenge_id"`
	FrontendPublicKey   string `json:"frontend_public_key"`
	FrontendSignature   string `json:"frontend_signature"`
}

type federationManageRequest struct {
	Method        string                 `json:"method"`
	Path          string                 `json:"path"`
	BodyB64       string                 `json:"body"`
	AuthorizeOnly bool                   `json:"authorize_only,omitempty"`
	GuardProof    *federationActionProof `json:"guard_proof,omitempty"`
}

type federationChallengeBinding struct {
	ShareID     int64
	ServerID    int64
	Action      string
	PayloadHash string
	Method      string
	Path        string
	BindingHash string
	ExpiresAt   int64
}

// FederationRequestError preserves an owner's status/code across the
// consumer federation hop. Callers must not flatten this to a generic 502:
// action_guard_retry is a protocol signal that requires a fresh challenge.
type FederationRequestError struct {
	Status   int
	Code     string
	Message  string
	ServerID int64
}

func (e *FederationRequestError) Error() string {
	if e == nil {
		return "federation request failed"
	}
	return e.Message
}

func validFederationActionTarget(method, requestPath string) bool {
	return strings.EqualFold(strings.TrimSpace(method), http.MethodPost) &&
		(requestPath == federationInboundRequestPath || requestPath == federationNodeRequestPath)
}

func validFederationSHA256(value string) bool {
	raw, err := hex.DecodeString(value)
	return err == nil && len(raw) == sha256.Size && value == strings.ToLower(value)
}

func federationActionBindingHash(shareID, serverID int64, action, method, requestPath, payloadHash string) string {
	message := "mmwx-federation-action-v1\n" + strconv.FormatInt(shareID, 10) + "\n" +
		strconv.FormatInt(serverID, 10) + "\n" + action + "\n" + strings.ToUpper(method) + "\n" +
		requestPath + "\n" + payloadHash
	sum := sha256.Sum256([]byte(message))
	return hex.EncodeToString(sum[:])
}

func federationProofFromRequest(r *http.Request, serverID int64, payloadHash, requestPath string) (*federationActionProof, error) {
	if r == nil || strings.TrimSpace(r.Header.Get(actionGuardServerIDHeader)) != strconv.FormatInt(serverID, 10) {
		return nil, &FederationRequestError{
			Status: http.StatusPreconditionRequired, Code: codeFederatedChallengeRequired, ServerID: serverID,
			Message: "共享服务器需要由拥有方签发新的安全 challenge",
		}
	}
	proof := &federationActionProof{
		Action: actionNodeCreate, OriginalPayloadHash: payloadHash,
		RequestMethod: http.MethodPost, RequestPath: requestPath,
		ChallengeID:       strings.TrimSpace(r.Header.Get("X-MMWX-Challenge-ID")),
		FrontendPublicKey: strings.TrimSpace(r.Header.Get("X-MMWX-Frontend-Key")),
		FrontendSignature: strings.TrimSpace(r.Header.Get("X-MMWX-Frontend-Signature")),
	}
	if !validFederationSHA256(proof.OriginalPayloadHash) || !validFederationActionTarget(proof.RequestMethod, proof.RequestPath) ||
		proof.ChallengeID == "" || proof.FrontendPublicKey == "" || proof.FrontendSignature == "" {
		return nil, &FederationRequestError{
			Status: http.StatusPreconditionRequired, Code: codeFederatedChallengeRequired, ServerID: serverID,
			Message: "共享服务器安全 challenge 缺失或与当前操作不匹配",
		}
	}
	return proof, nil
}

type federationGuardProofContextKey struct{}

func withFederationGuardProof(ctx context.Context, proof *federationActionProof) context.Context {
	return context.WithValue(ctx, federationGuardProofContextKey{}, proof)
}

func federationGuardProofFromContext(ctx context.Context) *federationActionProof {
	proof, _ := ctx.Value(federationGuardProofContextKey{}).(*federationActionProof)
	return proof
}

func (h *RemoteManageHandler) ProxyFederatedActionChallenge(
	ctx context.Context,
	serverID int64,
	request FederationActionChallengeRequest,
) (guardclient.Challenge, bool, error) {
	fed, err := h.repo.GetFederatedServer(ctx, serverID)
	if err != nil {
		if errors.Is(err, storage.ErrFederatedServerNotFound) {
			return guardclient.Challenge{}, false, nil
		}
		return guardclient.Challenge{}, true, &FederationRequestError{
			Status: http.StatusInternalServerError, Code: codeFederationGuardFailed, ServerID: serverID,
			Message: "无法确认服务器的共享状态",
		}
	}
	if request.Action != actionNodeCreate || !validFederationSHA256(request.PayloadHash) ||
		!validFederationActionTarget(request.RequestMethod, request.RequestPath) {
		return guardclient.Challenge{}, true, &FederationRequestError{
			Status: http.StatusBadRequest, Code: codeFederationGuardInvalid, ServerID: serverID,
			Message: "共享服务器 challenge 的操作绑定无效",
		}
	}
	body, _ := json.Marshal(request)
	endpoint := strings.TrimRight(fed.OwnerURL, "/") + "/api/federation/action-guard/challenge"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return guardclient.Challenge{}, true, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Share-Token", fed.ShareToken)
	req.Header.Set("User-Agent", version.AgentUserAgent)
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return guardclient.Challenge{}, true, &FederationRequestError{
			Status: http.StatusBadGateway, Code: codeFederationGuardFailed, ServerID: serverID,
			Message: "无法连接拥有方主控以签发共享服务器 challenge: " + err.Error(),
		}
	}
	defer resp.Body.Close()
	responseBody, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if readErr != nil {
		return guardclient.Challenge{}, true, readErr
	}
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
		return guardclient.Challenge{}, true, &FederationRequestError{
			Status: http.StatusUpgradeRequired, Code: codeFederationGuardUnsupported, ServerID: serverID,
			Message: "拥有方主控版本不支持共享服务器联合验证，请先升级拥有方主控",
		}
	}
	if resp.StatusCode >= http.StatusBadRequest {
		ferr := federationErrorFromBody(resp.StatusCode, responseBody)
		if typed, ok := ferr.(*FederationRequestError); ok {
			typed.ServerID = serverID
		}
		return guardclient.Challenge{}, true, ferr
	}
	var challenge guardclient.Challenge
	if json.Unmarshal(responseBody, &challenge) != nil || challenge.ChallengeID == "" ||
		challenge.Challenge == "" || challenge.PayloadHash == "" {
		return guardclient.Challenge{}, true, &FederationRequestError{
			Status: http.StatusBadGateway, Code: codeFederationGuardInvalid, ServerID: serverID,
			Message: "拥有方返回了无效的共享服务器 challenge",
		}
	}
	return challenge, true, nil
}

func (h *RemoteManageHandler) prepareFederatedManagedNodeAction(
	ctx context.Context,
	browserRequest *http.Request,
	serverID int64,
	payloadHash string,
	requestPath string,
	authorizeOnly bool,
) (context.Context, bool, error) {
	fed, err := h.repo.GetFederatedServer(ctx, serverID)
	if err != nil {
		if errors.Is(err, storage.ErrFederatedServerNotFound) {
			return ctx, false, nil
		}
		return ctx, true, &FederationRequestError{
			Status: http.StatusInternalServerError, Code: codeFederationGuardFailed, ServerID: serverID,
			Message: "无法确认服务器的共享状态",
		}
	}
	proof, err := federationProofFromRequest(browserRequest, serverID, payloadHash, requestPath)
	if err != nil {
		return ctx, true, err
	}
	if !authorizeOnly {
		return withFederationGuardProof(ctx, proof), true, nil
	}
	request := federationManageRequest{
		Method: http.MethodPost, Path: requestPath, AuthorizeOnly: true, GuardProof: proof,
	}
	if _, err := h.doFederationWireRequest(ctx, fed, request); err != nil {
		if typed, ok := err.(*FederationRequestError); ok {
			typed.ServerID = serverID
		}
		return ctx, true, err
	}
	return ctx, true, nil
}

func federationErrorResponse(err error, fallbackStatus int, fallbackCode string) (int, map[string]any) {
	status := fallbackStatus
	code := fallbackCode
	message := err.Error()
	serverID := int64(0)
	var typed *FederationRequestError
	if errors.As(err, &typed) {
		if typed.Status >= 400 && typed.Status <= 599 {
			status = typed.Status
		}
		if typed.Code != "" {
			code = typed.Code
		}
		serverID = typed.ServerID
	}
	payload := map[string]any{"error": message, "message": message, "status": status}
	if code != "" {
		payload["code"] = code
	}
	if serverID > 0 {
		payload["server_id"] = serverID
	}
	return status, payload
}

func (h *FederationHandler) pruneFederationChallengesLocked(now int64) {
	for id, binding := range h.guardChallenges {
		if binding.ExpiresAt <= now {
			delete(h.guardChallenges, id)
		}
	}
}

func (h *FederationHandler) consumeFederationChallenge(
	proof *federationActionProof,
	share storage.SharedServer,
	serverID int64,
) (string, error) {
	if proof == nil || proof.ChallengeID == "" {
		return "", &FederationRequestError{Status: http.StatusPreconditionRequired, Code: codeFederatedChallengeRequired, Message: "共享服务器联合验证 proof 缺失"}
	}
	now := time.Now().Unix()
	h.guardChallengeMu.Lock()
	h.pruneFederationChallengesLocked(now)
	binding, ok := h.guardChallenges[proof.ChallengeID]
	if ok {
		delete(h.guardChallenges, proof.ChallengeID)
	}
	h.guardChallengeMu.Unlock()
	if !ok || binding.ExpiresAt <= now {
		return "", &FederationRequestError{Status: http.StatusUnauthorized, Code: codeFederationGuardInvalid, Message: "共享服务器 challenge 无效、过期或已使用"}
	}
	method := strings.ToUpper(strings.TrimSpace(proof.RequestMethod))
	expectedHash := federationActionBindingHash(share.ID, serverID, proof.Action, method, proof.RequestPath, proof.OriginalPayloadHash)
	if binding.ShareID != share.ID || binding.ServerID != serverID || binding.Action != proof.Action ||
		binding.PayloadHash != proof.OriginalPayloadHash || binding.Method != method || binding.Path != proof.RequestPath ||
		binding.BindingHash != expectedHash || !validFederationSHA256(proof.OriginalPayloadHash) ||
		!validFederationActionTarget(method, proof.RequestPath) || proof.FrontendPublicKey == "" || proof.FrontendSignature == "" {
		return "", &FederationRequestError{Status: http.StatusUnauthorized, Code: codeFederationGuardInvalid, Message: "共享服务器 proof 与分享、请求路径或请求内容不匹配"}
	}
	return expectedHash, nil
}

func (h *FederationHandler) handleActionGuardChallenge(w http.ResponseWriter, r *http.Request, serverID int64) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "POST only", "code": codeFederationGuardInvalid})
		return
	}
	share, err := h.authenticatedShare(r)
	if err != nil || share.ServerID != serverID {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid or revoked share token", "code": codeFederationGuardInvalid})
		return
	}
	var request FederationActionChallengeRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&request) != nil || request.Action != actionNodeCreate ||
		!validFederationSHA256(request.PayloadHash) || !validFederationActionTarget(request.RequestMethod, request.RequestPath) {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid federation action challenge request", "code": codeFederationGuardInvalid})
		return
	}
	if h.createChallenge == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": "owner Action Guard is unavailable", "code": codeFederationGuardFailed})
		return
	}
	// A browser challenge is single-use. Repair an Agent slot that still belongs
	// to the owner's previous license before allocating that challenge, otherwise
	// a slow replacement turns an otherwise recoverable license switch into a
	// failed federated mutation. The consumer recognizes action_guard_retry and
	// repeats the whole operation, obtaining a challenge only after the exact
	// replacement lease has been acknowledged.
	if h.prepareAction != nil {
		if err := h.prepareAction(r.Context(), serverID); err != nil {
			writeJSON(w, http.StatusForbidden, map[string]any{
				"error":     "拥有方服务器许可证租约尚未就绪: " + err.Error(),
				"message":   "拥有方服务器许可证租约尚未就绪: " + err.Error(),
				"code":      "action_guard_retry",
				"server_id": serverID,
			})
			return
		}
	}
	method := strings.ToUpper(strings.TrimSpace(request.RequestMethod))
	bindingHash := federationActionBindingHash(share.ID, serverID, request.Action, method, request.RequestPath, request.PayloadHash)
	challenge, err := h.createChallenge(r.Context(), guardclient.ChallengeRequest{Action: request.Action, PayloadHash: bindingHash})
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": diagnoseActionGuardError(err).Error(), "code": codeFederationGuardFailed})
		return
	}
	if challenge.ChallengeID == "" || challenge.Challenge == "" || challenge.ExpiresAt <= time.Now().Unix() {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": "owner Action Guard returned an invalid challenge", "code": codeFederationGuardInvalid})
		return
	}
	challenge.PayloadHash = bindingHash
	h.guardChallengeMu.Lock()
	if h.guardChallenges == nil {
		h.guardChallenges = make(map[string]federationChallengeBinding)
	}
	h.pruneFederationChallengesLocked(time.Now().Unix())
	h.guardChallenges[challenge.ChallengeID] = federationChallengeBinding{
		ShareID: share.ID, ServerID: serverID, Action: request.Action, PayloadHash: request.PayloadHash,
		Method: method, Path: request.RequestPath, BindingHash: bindingHash, ExpiresAt: challenge.ExpiresAt,
	}
	h.guardChallengeMu.Unlock()
	writeJSON(w, http.StatusOK, challenge)
}

func (h *FederationHandler) authorizeFederationAction(
	ctx context.Context,
	share storage.SharedServer,
	serverID int64,
	proof *federationActionProof,
) error {
	bindingHash, err := h.consumeFederationChallenge(proof, share, serverID)
	if err != nil {
		return err
	}
	if h.authorizeAction == nil {
		return &FederationRequestError{Status: http.StatusServiceUnavailable, Code: codeFederationGuardFailed, Message: "owner Action Guard is unavailable"}
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://federation.invalid"+proof.RequestPath, nil)
	if err != nil {
		return &FederationRequestError{Status: http.StatusInternalServerError, Code: codeFederationGuardFailed, Message: err.Error()}
	}
	request.Header.Set("X-MMWX-Challenge-ID", proof.ChallengeID)
	request.Header.Set("X-MMWX-Frontend-Key", proof.FrontendPublicKey)
	request.Header.Set("X-MMWX-Frontend-Signature", proof.FrontendSignature)
	if err := h.authorizeAction(ctx, request, serverID, bindingHash); err != nil {
		code := codeFederationGuardFailed
		status := http.StatusForbidden
		if IsAuthoritativeSlotRepaired(err) {
			code = "action_guard_retry"
		} else if strings.Contains(err.Error(), "Agent 未连接 WS") || strings.Contains(err.Error(), "版本不支持 RPC") {
			code = "federation_owner_agent_ws_unavailable"
			status = http.StatusConflict
		}
		return &FederationRequestError{Status: status, Code: code, Message: err.Error()}
	}
	return nil
}

func isFederationInboundCreate(method, requestPath string, body []byte) bool {
	if !strings.EqualFold(method, http.MethodPost) || stripQuery(requestPath) != "/api/child/inbounds" {
		return false
	}
	var request struct {
		Action string `json:"action"`
	}
	if json.Unmarshal(body, &request) != nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(request.Action)) {
	case "", "add", "replace":
		return true
	default:
		return false
	}
}

func (h *FederationHandler) writeActionGuardError(w http.ResponseWriter, err error) {
	status, payload := federationErrorResponse(err, http.StatusForbidden, codeFederationGuardFailed)
	writeJSON(w, status, payload)
}
