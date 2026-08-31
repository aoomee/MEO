// Package guardclient 提供无外部守护进程的本地兼容实现。
//
// 原实现通过 Unix socket 连接独立签名的 mmwx-guardd-agent 守护进程做授权/证明，
// Agent 无 Guard 即拒绝启动。定制版把该客户端改为【本地桩】:所有方法直接在本进程
// 返回“健康 / 已授权 / 全功能”，不再依赖任何外部 Guard 守护进程或 socket。
//
// 这样：主控与 Agent 之间不再需要 Guard 加密会话，Agent 单独运行即可，
// 所有原 PRO/Premium 能力常开。
package guardclient

import (
	"context"
	"time"
)

type AttestationRequest struct {
	Action      string `json:"action"`
	PayloadHash string `json:"payload_hash"`
	ServerHash  string `json:"server_hash"`
}

type Attestation struct {
	Version         int    `json:"version"`
	ID              string `json:"id"`
	Role            string `json:"role"`
	Action          string `json:"action"`
	PayloadHash     string `json:"payload_hash"`
	ServerHash      string `json:"server_hash"`
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
	ServerHash       string `json:"server_hash"`
	LicenseServerURL string `json:"license_server_url"`
}

type SlotDelivery struct {
	Reservation      string `json:"reservation"`
	LicenseServerURL string `json:"license_server_url"`
}

type SlotStatus struct {
	Authorized     bool     `json:"authorized"`
	Renewable      bool     `json:"renewable,omitempty"`
	ServerHash     string   `json:"server_hash,omitempty"`
	LicenseKeyHash string   `json:"license_key_hash,omitempty"`
	SlotID         int64    `json:"slot_id,omitempty"`
	Generation     int64    `json:"generation,omitempty"`
	ExpiresAt      int64    `json:"expires_at,omitempty"`
	Features       []string `json:"features,omitempty"`
}

type Health struct {
	OK             bool   `json:"ok"`
	Role           string `json:"role"`
	Version        string `json:"version"`
	CallerVerified bool   `json:"caller_verified"`
}

// AllFeatures 是自托管模式对外常开的能力集合。
var AllFeatures = []string{
	"limiter", "embedded", "speed_test", "custom_branding",
	"server_share", "reality_pool",
}

// farFuture 为旧协议兼容字段提供稳定的远期时间。
func farFuture() int64 { return time.Now().Add(100 * 365 * 24 * time.Hour).Unix() }

// Client 是去 Guard 后的本地桩,不持有任何 socket / 网络连接。
type Client struct{}

func NewFromEnv() *Client           { return &Client{} }
func NewForSocket(_ string) *Client { return &Client{} }

func (c *Client) Enabled() bool  { return c != nil }
func (c *Client) Required() bool { return c != nil }

// Enforced 对齐官方 v0.6.1：运行时和升级脚本必须用同一个值。
// 自托管构建不安装外部 Guard，恒 false。
func Enforced() bool { return false }

func (c *Client) Enforced() bool { return Enforced() }

// Health 恒返回健康(role=agent, caller_verified=true)。
func (c *Client) Health(_ context.Context) (Health, error) {
	return Health{OK: true, Role: "agent", Version: "offline-noguard", CallerVerified: true}, nil
}

// Attest 本地生成一份结构完整、无签名的证明(主控侧同为去 Guard,不再校验签名)。
func (c *Client) Attest(_ context.Context, request AttestationRequest) (Attestation, error) {
	now := time.Now().Unix()
	return Attestation{
		Version:     1,
		ID:          "offline-noguard",
		Role:        "agent",
		Action:      request.Action,
		PayloadHash: request.PayloadHash,
		ServerHash:  request.ServerHash,
		IssuedAt:    now,
		ExpiresAt:   now + 300,
	}, nil
}

// Consume 恒放行。
func (c *Client) Consume(_ context.Context, _ ConsumeRequest) error { return nil }

func (c *Client) authorizedSlot() SlotStatus {
	return SlotStatus{Authorized: true, Renewable: true, ExpiresAt: farFuture(), Features: AllFeatures}
}

func (c *Client) ActivateSlot(_ context.Context, _ SlotDelivery) (SlotStatus, error) {
	return c.authorizedSlot(), nil
}
func (c *Client) RefreshSlot(_ context.Context) (SlotStatus, error) { return c.authorizedSlot(), nil }
func (c *Client) ReleaseSlot(_ context.Context) error               { return nil }
func (c *Client) SlotStatus(_ context.Context) (SlotStatus, error)  { return c.authorizedSlot(), nil }

// ErrorMetadata 在本地兼容实现中无错误元数据，恒返回空。
func ErrorMetadata(_ error) (code string, status int, retryAfter time.Duration) {
	return "", 0, 0
}
