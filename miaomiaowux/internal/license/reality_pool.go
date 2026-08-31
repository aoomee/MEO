package license

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
)

const localRealityPoolKey = "meo_reality_domain_pool"

type RealityPoolDomain struct {
	Domain       string `json:"domain"`
	TLSVersion   string `json:"tls_version"`
	CipherSuite  string `json:"cipher_suite"`
	CurveID      string `json:"curve_id"`
	CertLen      int    `json:"cert_len"`
	Contributors int    `json:"contributors"`
}

var ErrRealityPoolUnavailable = errors.New("本地域名池不可用")

func (m *Manager) readRealityDomains(ctx context.Context) ([]string, error) {
	if m == nil || m.settings == nil {
		return nil, ErrRealityPoolUnavailable
	}
	raw, err := m.settings.GetSystemSetting(ctx, localRealityPoolKey)
	if err != nil || strings.TrimSpace(raw) == "" {
		return []string{}, nil
	}
	var domains []string
	if err := json.Unmarshal([]byte(raw), &domains); err != nil {
		return nil, err
	}
	return domains, nil
}

func (m *Manager) writeRealityDomains(ctx context.Context, domains []string) error {
	sort.Strings(domains)
	data, err := json.Marshal(domains)
	if err != nil {
		return err
	}
	return m.settings.SetSystemSetting(ctx, localRealityPoolKey, string(data))
}

func normalizeDomains(domains []string) []string {
	seen := make(map[string]struct{}, len(domains))
	out := make([]string, 0, len(domains))
	for _, domain := range domains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if domain == "" {
			continue
		}
		if _, ok := seen[domain]; ok {
			continue
		}
		seen[domain] = struct{}{}
		out = append(out, domain)
	}
	sort.Strings(out)
	return out
}

func (m *Manager) SubmitRealityDomains(ctx context.Context, domains []string) ([]string, map[string]string, error) {
	current, err := m.readRealityDomains(ctx)
	if err != nil {
		return nil, nil, err
	}
	accepted := normalizeDomains(domains)
	if err := m.writeRealityDomains(ctx, normalizeDomains(append(current, accepted...))); err != nil {
		return nil, nil, err
	}
	return accepted, map[string]string{}, nil
}

func (m *Manager) WithdrawRealityDomains(ctx context.Context, domains []string) ([]string, error) {
	current, err := m.readRealityDomains(ctx)
	if err != nil {
		return nil, err
	}
	remove := make(map[string]struct{})
	for _, domain := range normalizeDomains(domains) {
		remove[domain] = struct{}{}
	}
	kept := current[:0]
	for _, domain := range current {
		if _, ok := remove[domain]; !ok {
			kept = append(kept, domain)
		}
	}
	if err := m.writeRealityDomains(ctx, normalizeDomains(kept)); err != nil {
		return nil, err
	}
	return normalizeDomains(domains), nil
}

func (m *Manager) ListRealityDomains(ctx context.Context) ([]RealityPoolDomain, error) {
	domains, err := m.readRealityDomains(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]RealityPoolDomain, 0, len(domains))
	for _, domain := range normalizeDomains(domains) {
		out = append(out, RealityPoolDomain{Domain: domain, Contributors: 1})
	}
	return out, nil
}
