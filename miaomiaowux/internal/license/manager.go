// Package license is a legacy import path retained for source compatibility.
// MEO runs entirely in self-hosted mode: all local capabilities are available,
// no quotas are enforced, and this package performs no network requests.
package license

import (
	"context"
	"time"
)

const (
	FreeEdition           = true
	FeatureCustomBranding = "custom_branding"
	FeatureSpeedTest      = "speed_test"
	FeatureRealityPool    = "reality_pool"
)

var localFeatures = []string{
	"limiter",
	"embedded",
	FeatureSpeedTest,
	FeatureCustomBranding,
	"server_share",
	FeatureRealityPool,
}

type PlanInfo struct {
	Name          string            `json:"name"`
	DisplayName   string            `json:"display_name"`
	Description   string            `json:"description,omitempty"`
	MaxServers    int               `json:"max_servers"`
	MaxNodes      int               `json:"max_nodes"`
	MaxUsers      int               `json:"max_users"`
	Features      []string          `json:"features"`
	FeatureTokens map[string]string `json:"feature_tokens,omitempty"`
}

type Status struct {
	Valid                 bool      `json:"valid"`
	Error                 string    `json:"error,omitempty"`
	MaxServers            int       `json:"max_servers"`
	ExpiresAt             string    `json:"expires_at,omitempty"`
	Plan                  *PlanInfo `json:"plan,omitempty"`
	LastCheck             time.Time `json:"last_check"`
	HardRevoked           bool      `json:"hard_revoked,omitempty"`
	Entitlement           string    `json:"entitlement,omitempty"`
	SigningKeyCertificate string    `json:"signing_key_certificate,omitempty"`
}

func (s *Status) HasFeature(name string) bool {
	if s == nil || s.Plan == nil {
		return false
	}
	for _, feature := range s.Plan.Features {
		if feature == name {
			return true
		}
	}
	return false
}

type SettingsGetter interface {
	GetSystemSetting(ctx context.Context, key string) (string, error)
}

type SettingsStore interface {
	SettingsGetter
	SetSystemSetting(ctx context.Context, key, value string) error
}

type Manager struct {
	settings SettingsStore
}

func NewManager(settings SettingsStore, _ string) *Manager { return &Manager{settings: settings} }
func (m *Manager) Start(context.Context)                   {}
func (m *Manager) Stop()                                   {}
func (m *Manager) Refresh(context.Context)                 {}

func localPlan() *PlanInfo {
	return &PlanInfo{
		DisplayName: "MEO",
		Description: "自托管",
		MaxServers:  999999,
		MaxNodes:    999999,
		MaxUsers:    999999,
		Features:    append([]string(nil), localFeatures...),
	}
}

func (m *Manager) GetStatus() Status {
	return Status{Valid: true, MaxServers: 999999, Plan: localPlan(), LastCheck: time.Now()}
}

func (m *Manager) IsValid() bool                      { return true }
func (m *Manager) QuotaEnforced() bool                { return false }
func (m *Manager) EffectiveServerQuota() int          { return 999999 }
func (m *Manager) ServerQuotaDecision() (int, bool)   { return 999999, false }
func (m *Manager) EffectiveNodeQuota() int            { return 999999 }
func (m *Manager) EffectiveUserQuota() int            { return 999999 }
func (m *Manager) VerificationLevel() int             { return 0 }
func (m *Manager) AgentGuardRequired() bool           { return false }
func (m *Manager) AgentLeaseIdentity() string         { return "" }
func (m *Manager) HasFeature(string) bool             { return true }
func (m *Manager) HasFeatureForDataPlane(string) bool { return true }

func (m *Manager) StatusForAgent() map[string]any {
	plan := localPlan()
	return map[string]any{
		"valid":       true,
		"max_servers": plan.MaxServers,
		"plan": map[string]any{
			"name":         plan.Name,
			"display_name": plan.DisplayName,
			"max_servers":  plan.MaxServers,
			"max_nodes":    plan.MaxNodes,
			"max_users":    plan.MaxUsers,
			"features":     plan.Features,
		},
	}
}

func GetMachineID() string { return persistentMachineID() }
