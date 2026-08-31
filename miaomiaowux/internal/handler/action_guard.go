package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"miaomiaowux/internal/guardclient"
	"miaomiaowux/internal/license"
)

var masterExecutableReplaced = runningMasterExecutableReplaced

func runningMasterExecutableReplaced() bool {
	path, err := os.Executable()
	if err != nil {
		return false
	}
	running, err := os.Stat("/proc/self/exe")
	if err != nil {
		return false
	}
	onDisk, err := os.Stat(path)
	return err == nil && !os.SameFile(running, onDisk)
}

func diagnoseActionGuardError(err error) error {
	if err == nil || !masterExecutableReplaced() {
		return err
	}
	return fmt.Errorf("master/guard version mismatch: 主控磁盘二进制已更新，但当前 mmwx 仍是旧进程；请重启 mmwx 服务（原始错误: %v）", err)
}

const (
	actionServerCreate = "server.create"
	actionNodeCreate   = "node.create"
)

type ActionGuard struct {
	client          *guardclient.Client
	license         *license.Manager
	challengeMu     sync.RWMutex
	federationProxy func(context.Context, int64, FederationActionChallengeRequest) (guardclient.Challenge, bool, error)
}

func NewActionGuard(client *guardclient.Client, licenseManager *license.Manager) *ActionGuard {
	return &ActionGuard{client: client, license: licenseManager}
}

func (g *ActionGuard) Required() bool { return g != nil && g.client != nil && g.client.Required() }
func (g *ActionGuard) Enabled() bool  { return g != nil && g.client != nil && g.client.Enabled() }

func (g *ActionGuard) VerificationLevel() int {
	level := 0
	if g != nil && g.license != nil {
		level = g.license.VerificationLevel()
	}
	// Current releases always require the local Guard, so signed license policy
	// may raise this to L2 but can never lower it below L1.
	// 自托管构建不依赖外部 Guard（client.Required()=false）。
	if g != nil && g.client != nil && g.client.Required() && level < 1 {
		level = 1
	}
	return level
}

func (g *ActionGuard) RequiresMaster() bool { return g != nil && g.VerificationLevel() >= 1 }
func (g *ActionGuard) RequiresAgent() bool  { return g != nil && g.VerificationLevel() >= 2 }

func (g *ActionGuard) SetFederationChallengeProxy(proxy func(context.Context, int64, FederationActionChallengeRequest) (guardclient.Challenge, bool, error)) {
	if g == nil {
		return
	}
	g.challengeMu.Lock()
	g.federationProxy = proxy
	g.challengeMu.Unlock()
}

func (g *ActionGuard) CreateChallenge(ctx context.Context, request guardclient.ChallengeRequest) (guardclient.Challenge, error) {
	if g == nil || g.client == nil || !g.client.Enabled() {
		return guardclient.Challenge{}, errors.New("Action Guard 未启用")
	}
	return g.client.CreateChallenge(ctx, request)
}

func (g *ActionGuard) HandleChallenge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if g == nil || !g.RequiresMaster() {
		respondJSON(w, http.StatusOK, map[string]any{"bypass": true, "verification_level": 0})
		return
	}
	if g.client == nil || !g.client.Enabled() {
		writeJSONError(w, http.StatusServiceUnavailable, "Action Guard 未启用")
		return
	}
	var request struct {
		Action        string `json:"action"`
		PayloadHash   string `json:"payload_hash"`
		ServerID      int64  `json:"server_id,omitempty"`
		RequestMethod string `json:"request_method,omitempty"`
		RequestPath   string `json:"request_path,omitempty"`
	}
	if json.NewDecoder(r.Body).Decode(&request) != nil {
		writeJSONError(w, http.StatusBadRequest, "请求格式不正确")
		return
	}
	baseRequest := guardclient.ChallengeRequest{Action: request.Action, PayloadHash: request.PayloadHash}
	if request.ServerID > 0 {
		g.challengeMu.RLock()
		proxy := g.federationProxy
		g.challengeMu.RUnlock()
		if proxy != nil {
			challenge, handled, err := proxy(r.Context(), request.ServerID, FederationActionChallengeRequest{
				Action: request.Action, PayloadHash: request.PayloadHash,
				RequestMethod: request.RequestMethod, RequestPath: request.RequestPath,
			})
			if handled {
				if err != nil {
					status, payload := federationErrorResponse(err, http.StatusBadGateway, codeFederationGuardFailed)
					writeJSON(w, status, payload)
					return
				}
				respondJSON(w, http.StatusOK, challenge)
				return
			}
		}
	}
	challenge, err := g.client.CreateChallenge(r.Context(), baseRequest)
	if err != nil {
		writeJSONError(w, http.StatusServiceUnavailable, diagnoseActionGuardError(err).Error())
		return
	}
	respondJSON(w, http.StatusOK, challenge)
}

func (g *ActionGuard) AuthorizeMaster(ctx context.Context, r *http.Request, action, payloadHash, serverHash string, agent guardclient.Attestation) (license.ActionGrantDelivery, error) {
	if g == nil || !g.RequiresMaster() {
		return license.ActionGrantDelivery{}, nil
	}
	if g.client == nil || !g.client.Enabled() {
		return license.ActionGrantDelivery{}, errors.New("Action Guard 不可用，已拒绝写操作")
	}
	attestation, err := g.client.Attest(ctx, guardclient.AttestationRequest{
		Action: action, PayloadHash: payloadHash, ServerHash: serverHash,
		ChallengeID:       r.Header.Get("X-MMWX-Challenge-ID"),
		FrontendPublicKey: r.Header.Get("X-MMWX-Frontend-Key"),
		FrontendSignature: r.Header.Get("X-MMWX-Frontend-Signature"),
	})
	if err != nil {
		return license.ActionGrantDelivery{}, diagnoseActionGuardError(err)
	}
	delivery, err := g.license.IssueActionGrant(ctx, action, payloadHash, serverHash, toLicenseAttestation(attestation), toLicenseAttestation(agent))
	if err != nil {
		return license.ActionGrantDelivery{}, err
	}
	return delivery, nil
}

func (g *ActionGuard) ConsumeMaster(ctx context.Context, delivery license.ActionGrantDelivery, action, payloadHash, serverHash string) error {
	if delivery.Grant == "" {
		if g == nil || !g.RequiresMaster() {
			return nil
		}
		return errors.New("ActionGrant 为空")
	}
	return diagnoseActionGuardError(g.client.Consume(ctx, guardclient.ConsumeRequest{
		Grant: delivery.Grant, Action: action, PayloadHash: payloadHash,
		ServerHash: serverHash, LicenseServerURL: delivery.LicenseServerURL,
	}))
}

func toLicenseAttestation(value guardclient.Attestation) license.ActionAttestation {
	return license.ActionAttestation{
		Version: value.Version, ID: value.ID, Role: value.Role, Action: value.Action,
		PayloadHash: value.PayloadHash, ServerHash: value.ServerHash,
		ChallengeID: value.ChallengeID, FrontendKeyHash: value.FrontendKeyHash,
		IssuedAt: value.IssuedAt, ExpiresAt: value.ExpiresAt,
		PublicKey: value.PublicKey, Signature: value.Signature,
	}
}

func hashActionPayload(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func serverTokenHash(token string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(token)))
	return hex.EncodeToString(sum[:])
}
