package license

import "context"

type ActionAttestation struct {
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

type ActionGrantDelivery struct {
	Grant            string `json:"grant"`
	LicenseServerURL string `json:"license_server_url"`
	ExpiresAt        int64  `json:"expires_at"`
}

type ActionGrantIssueError struct {
	Message string
	Code    string
}

func (e *ActionGrantIssueError) Error() string { return e.Message }
func IsAuthoritativeSlotError(error) bool      { return false }

func (m *Manager) IssueActionGrant(context.Context, string, string, string, ActionAttestation, ActionAttestation) (ActionGrantDelivery, error) {
	return ActionGrantDelivery{}, nil
}
