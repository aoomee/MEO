// Package guardclient 提供无外部守护进程的本地兼容实现。
//
// 原实现通过 Unix socket 连接独立签名的 mmwx-guardd 守护进程,做启动自证明
// (RuntimeProof)、管理写操作授权(Challenge/Attest/Consume)与前端安全通道后端
// 证明(BackendProof);主控无 Guard 即拒绝启动。定制版把该客户端改为【本地桩】:
//
//   - Enabled()/Required() = false  → ActionGuard 旁路(verification_level=0),
//     securechan 握手跳过后端证明(E2E 加密仍照常工作,只是不带官方证明)。
//   - Health() = nil                → 主控启动不再因 Guard 不可用而退出。
//   - 其余方法返回结构完整的良性值(几乎不会被调用,因为上面两个开关已让相关路径旁路)。
package guardclient

import "context"

type Client struct{}

type ChallengeRequest struct {
	Action      string `json:"action"`
	PayloadHash string `json:"payload_hash"`
}

type Challenge struct {
	ChallengeID string `json:"challenge_id"`
	Challenge   string `json:"challenge"`
	PayloadHash string `json:"payload_hash,omitempty"`
	ExpiresAt   int64  `json:"expires_at"`
}

type AttestationRequest struct {
	Action            string `json:"action"`
	PayloadHash       string `json:"payload_hash"`
	ServerHash        string `json:"server_hash,omitempty"`
	ChallengeID       string `json:"challenge_id,omitempty"`
	FrontendPublicKey string `json:"frontend_public_key,omitempty"`
	FrontendSignature string `json:"frontend_signature,omitempty"`
}

type Attestation struct {
	Version         int    `json:"version"`
	ID              string `json:"id"`
	Role            string `json:"role"`
	Action          string `json:"action"`
	PayloadHash     string `json:"payload_hash"`
	ServerHash      string `json:"server_hash,omitempty"`
	ChallengeID     string `json:"challenge_id,omitempty"`
	FrontendKeyHash string `json:"frontend_key_hash,omitempty"`
	IssuedAt        int64  `json:"issued_at"`
	ExpiresAt       int64  `json:"expires_at"`
	PublicKey       string `json:"public_key"`
	Signature       string `json:"signature"`
}

type ConsumeRequest struct {
	Grant            string `json:"grant"`
	Action           string `json:"action"`
	PayloadHash      string `json:"payload_hash"`
	ServerHash       string `json:"server_hash,omitempty"`
	LicenseServerURL string `json:"license_server_url"`
}

type RuntimeProof struct {
	Proof     map[string]any `json:"proof"`
	Signature string         `json:"signature"`
	Manifest  struct {
		Payload   string `json:"payload"`
		Signature string `json:"signature"`
	} `json:"manifest"`
}

type BackendProofRequest struct {
	ClientPublicKey string `json:"client_public_key"`
	ServerPublicKey string `json:"server_public_key"`
	SessionID       string `json:"session_id"`
	Audience        string `json:"audience"`
}

type BackendProofClaims struct {
	Version          int    `json:"version"`
	Issuer           string `json:"issuer"`
	Audience         string `json:"audience"`
	Role             string `json:"role"`
	GuardRelease     string `json:"guard_release"`
	MasterRelease    string `json:"master_release"`
	ExecutableSHA256 string `json:"executable_sha256"`
	ClientPublicKey  string `json:"client_public_key"`
	ServerPublicKey  string `json:"server_public_key"`
	SessionID        string `json:"session_id"`
	Origin           string `json:"origin"`
	KeyID            string `json:"key_id"`
	IssuedAt         int64  `json:"issued_at"`
	ExpiresAt        int64  `json:"expires_at"`
}

type BackendProof struct {
	Proof          BackendProofClaims `json:"proof"`
	Signature      string             `json:"signature"`
	KeyCertificate string             `json:"key_certificate"`
	Manifest       struct {
		Payload   string `json:"payload"`
		Signature string `json:"signature"`
	} `json:"manifest"`
}

func NewFromEnv() *Client { return &Client{} }

// Enabled/Required = false 是去 Guard 的关键开关(见包注释)。
func (c *Client) Enabled() bool  { return false }
func (c *Client) Required() bool { return false }

func (c *Client) Health(_ context.Context) error { return nil }

func (c *Client) CreateChallenge(_ context.Context, request ChallengeRequest) (Challenge, error) {
	return Challenge{ChallengeID: "offline", Challenge: "offline", PayloadHash: request.PayloadHash}, nil
}

func (c *Client) Attest(_ context.Context, request AttestationRequest) (Attestation, error) {
	return Attestation{Version: 1, ID: "offline", Role: "master", Action: request.Action, PayloadHash: request.PayloadHash, ServerHash: request.ServerHash}, nil
}

func (c *Client) Consume(_ context.Context, _ ConsumeRequest) error { return nil }

func (c *Client) RuntimeProof(_ context.Context) (RuntimeProof, error) {
	rp := RuntimeProof{Proof: map[string]any{"role": "master"}, Signature: "offline"}
	rp.Manifest.Payload = "offline"
	rp.Manifest.Signature = "offline"
	return rp, nil
}

func (c *Client) BackendProof(_ context.Context, request BackendProofRequest) (BackendProof, error) {
	bp := BackendProof{Signature: "offline", KeyCertificate: "offline"}
	bp.Proof = BackendProofClaims{Version: 1, Role: "master", ClientPublicKey: request.ClientPublicKey, ServerPublicKey: request.ServerPublicKey, SessionID: request.SessionID}
	bp.Manifest.Payload = "offline"
	bp.Manifest.Signature = "offline"
	return bp, nil
}

func (c *Client) ReloadManifest(_ context.Context) error { return nil }
