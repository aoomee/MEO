// Package licenselease 提供本地租约兼容层。
//
// 原实现向 Agent Guard 申请/续期“授权槽位(slot)”，据此开关 Premium 能力与防克隆审计，
// Guard 不可用即视为未授权。定制版把它改为【常开授权】:不依赖 Guard，恒返回已授权、
// 全功能、永不过期。保留原有公开方法签名，调用方无需改动。
package licenselease

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
	"time"

	"mmw-agent/internal/guardclient"
)

type Delivery struct {
	Reservation      string `json:"reservation"`
	LicenseServerURL string `json:"license_server_url"`
	ExpiresAt        int64  `json:"reservation_expires_at"`
	RequestID        string `json:"request_id,omitempty"`
}

// guardAuthority 保留接口以兼容 New 的签名(参数被忽略)。
type guardAuthority interface {
	Enabled() bool
	SlotStatus(context.Context) (guardclient.SlotStatus, error)
	ActivateSlot(context.Context, guardclient.SlotDelivery) (guardclient.SlotStatus, error)
	RefreshSlot(context.Context) (guardclient.SlotStatus, error)
	ReleaseSlot(context.Context) error
}

type Manager struct {
	mu         sync.RWMutex
	statePath  string
	serverHash string
	onChange   func(bool)
}

// New 直接返回一个常开授权的管理器,不再要求 Guard。
func New(_ string, statePath, serverToken string, _ guardAuthority) (*Manager, error) {
	return &Manager{statePath: statePath, serverHash: hashServerToken(serverToken)}, nil
}

func (m *Manager) Required() bool { return false }

func (m *Manager) UpdateServerToken(token string) error {
	m.mu.Lock()
	m.serverHash = hashServerToken(token)
	m.mu.Unlock()
	return nil
}

// Status 返回一份合成的“已授权、全功能、永不过期”槽位状态。
func (m *Manager) Status() guardclient.SlotStatus {
	m.mu.RLock()
	hash := m.serverHash
	m.mu.RUnlock()
	return guardclient.SlotStatus{
		Authorized: true,
		Renewable:  true,
		ServerHash: hash,
		ExpiresAt:  time.Now().Add(100 * 365 * 24 * time.Hour).Unix(),
		Features:   guardclient.AllFeatures,
	}
}

func (m *Manager) ServerHash() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.serverHash
}

func hashServerToken(token string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(token)))
	return hex.EncodeToString(sum[:])
}

func (m *Manager) SetAuthorizationHandler(fn func(bool)) {
	m.mu.Lock()
	m.onChange = fn
	m.mu.Unlock()
	// 常开授权:立即通知一次“已授权”。
	if fn != nil {
		fn(true)
	}
}

func (m *Manager) Authorized() bool { return true }

func (m *Manager) NeedsLease() bool { return false }

func (m *Manager) HasFeature(_ string) bool { return true }

// HandleDelivery 定制版无需真正激活槽位,直接成功。
func (m *Manager) HandleDelivery(_ Delivery) error { return nil }

// Start 无后台续期逻辑(常开授权不需要)。
func (m *Manager) Start(_ context.Context) {}

// Release 无槽位可释放。
func (m *Manager) Release(_ context.Context) error { return nil }

// LeaseIdentityCapable 在自托管构建中始终具备租约身份能力。
func (m *Manager) LeaseIdentityCapable() bool { return true }
