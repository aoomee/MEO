package license

import (
	"context"
	"errors"
	"time"
)

type AgentLeaseDelivery struct {
	Reservation      string `json:"reservation"`
	LicenseServerURL string `json:"license_server_url"`
	ExpiresAt        int64  `json:"reservation_expires_at"`
	LicenseIdentity  string `json:"license_identity,omitempty"`
}

type AgentLeaseIssueError struct {
	Message    string
	Code       string
	StatusCode int
	RetryAfter time.Duration
}

var ErrAgentLeaseEntitlementUnavailable = errors.New("remote lease service is disabled")
var ErrLicenseRequestSuperseded = errors.New("request superseded")

func (e *AgentLeaseIssueError) Error() string  { return e.Message }
func AgentLeaseRetryAfter(error) time.Duration { return 0 }
func AgentLeaseShouldRetry(error) bool         { return false }

func (m *Manager) IssueAgentLease(context.Context, string) (AgentLeaseDelivery, error) {
	return AgentLeaseDelivery{}, ErrAgentLeaseEntitlementUnavailable
}

func (m *Manager) IssueAgentLeaseReplacement(context.Context, string) (AgentLeaseDelivery, error) {
	return AgentLeaseDelivery{}, ErrAgentLeaseEntitlementUnavailable
}

func (m *Manager) IssueAgentLeaseTokenRotation(context.Context, string, string) (AgentLeaseDelivery, error) {
	return AgentLeaseDelivery{}, ErrAgentLeaseEntitlementUnavailable
}
