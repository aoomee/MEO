package handler

import (
	"errors"
	"strings"
	"time"

	"miaomiaowux/internal/license"
)

type agentLeaseActivationFailure struct {
	RequestID      string `json:"request_id"`
	Error          string `json:"error"`
	Code           string `json:"code,omitempty"`
	UpstreamStatus int    `json:"upstream_status,omitempty"`
	RetryAfter     int64  `json:"retry_after,omitempty"`
}

func newAgentLeaseActivationError(message, code string, status int, retryAfterSeconds int64) error {
	message = strings.TrimSpace(message)
	if message == "" {
		message = "Agent Guard rejected lease activation"
	}
	if code == "" && retryAfterSeconds <= 0 {
		// Old Guards and local validation failures may provide an HTTP status but
		// no stable service code. Keep their historical bounded retry/cooldown
		// behavior instead of turning (for example) an expired reservation's 401
		// into a permanent tuple terminal.
		return errors.New(message)
	}
	if agentLeaseActivationCanSelfRepair(code) && retryAfterSeconds <= 0 {
		// A fresh reservation/replacement can repair these authoritative-state
		// conflicts. A short typed delay routes them through the existing bounded
		// activation retry path rather than the permanent issue-error classifier.
		retryAfterSeconds = 1
	}
	retryAfter := time.Duration(0)
	if retryAfterSeconds > 0 {
		retryAfter = time.Duration(retryAfterSeconds) * time.Second
	}
	return &license.AgentLeaseIssueError{
		Message: message, Code: code, StatusCode: status, RetryAfter: retryAfter,
	}
}

func agentLeaseActivationCanSelfRepair(code string) bool {
	switch code {
	case "server_slot_stale_generation", "server_slot_state_conflict", "server_slot_not_active", "stale_license_authority":
		return true
	default:
		return false
	}
}
