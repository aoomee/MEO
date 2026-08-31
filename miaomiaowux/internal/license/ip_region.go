package license

import (
	"context"
	"errors"
	"strings"
)

type IPRegion struct {
	Country         string `json:"country"`
	Region          string `json:"region"`
	City            string `json:"city"`
	ProviderName    string `json:"provider_name,omitempty"`
	ProviderURL     string `json:"provider_url,omitempty"`
	TelecomPaidPeer bool   `json:"telecom_paid_peer,omitempty"`
}

func (r IPRegion) Label() string {
	parts := make([]string, 0, 3)
	for _, part := range []string{r.Country, r.Region, r.City} {
		if part != "" && (len(parts) == 0 || parts[len(parts)-1] != part) {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, " · ")
}

func (r IPRegion) Flag() string {
	code := strings.ToUpper(strings.TrimSpace(r.Country))
	if len(code) != 2 || code[0] < 'A' || code[0] > 'Z' || code[1] < 'A' || code[1] > 'Z' {
		return ""
	}
	return string([]rune{0x1F1E6 + rune(code[0]-'A'), 0x1F1E6 + rune(code[1]-'A')})
}

func (m *Manager) ResolveIPRegion(context.Context, string) (IPRegion, error) {
	return IPRegion{}, errors.New("automatic IP region lookup is disabled in self-hosted mode")
}
