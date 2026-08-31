package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/netip"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"miaomiaowux/internal/storage"
)

var globalSubscriptionIPWhitelist atomic.Pointer[SubscriptionIPWhitelist]

const (
	subscriptionIPWhitelistKey      = "subscription_ip_whitelist"
	subscriptionAutoWhitelistKey    = "subscription_auto_whitelist_managed_servers"
	managedWhitelistRefreshInterval = time.Minute
)

// ManagedServerWhitelistEntry describes an automatically trusted public IP.
// Shared/federated servers are deliberately excluded: only Agents managed by
// this master may bypass subscription abuse protection by default.
type ManagedServerWhitelistEntry struct {
	ServerID int64  `json:"server_id"`
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Family   string `json:"family"`
}

type SubscriptionIPWhitelistSnapshot struct {
	ManualEntries             []string                      `json:"manual_entries"`
	AutoIncludeManagedServers bool                          `json:"auto_include_managed_servers"`
	ManagedServers            []ManagedServerWhitelistEntry `json:"managed_servers"`
}

// SubscriptionIPWhitelist is the shared allow-list used by both subscription
// brute-force blocking and subscription request rate limiting.
type SubscriptionIPWhitelist struct {
	mu sync.RWMutex

	repo        *storage.TrafficRepository
	manualText  []string
	manual      []netip.Prefix
	autoManaged bool
	managed     map[netip.Addr]ManagedServerWhitelistEntry
}

func NewSubscriptionIPWhitelist(repo *storage.TrafficRepository) *SubscriptionIPWhitelist {
	w := &SubscriptionIPWhitelist{
		repo:        repo,
		autoManaged: true,
		managed:     make(map[netip.Addr]ManagedServerWhitelistEntry),
	}
	if repo == nil {
		globalSubscriptionIPWhitelist.Store(w)
		return w
	}

	ctx := context.Background()
	if raw, err := repo.GetSystemSetting(ctx, subscriptionAutoWhitelistKey); err == nil && strings.TrimSpace(raw) != "" {
		if enabled, parseErr := strconv.ParseBool(raw); parseErr == nil {
			w.autoManaged = enabled
		}
	}
	if raw, err := repo.GetSystemSetting(ctx, subscriptionIPWhitelistKey); err == nil && strings.TrimSpace(raw) != "" {
		var entries []string
		if json.Unmarshal([]byte(raw), &entries) == nil {
			if normalized, prefixes, normalizeErr := normalizeSubscriptionWhitelist(entries); normalizeErr == nil {
				w.manualText = normalized
				w.manual = prefixes
			}
		}
	}
	_, _ = w.RefreshManaged(ctx)
	globalSubscriptionIPWhitelist.Store(w)
	return w
}

func GetSubscriptionIPWhitelist() *SubscriptionIPWhitelist {
	return globalSubscriptionIPWhitelist.Load()
}

func normalizeSubscriptionWhitelist(entries []string) ([]string, []netip.Prefix, error) {
	seen := make(map[string]struct{}, len(entries))
	normalized := make([]string, 0, len(entries))
	prefixes := make([]netip.Prefix, 0, len(entries))
	for _, raw := range entries {
		entry := strings.TrimSpace(raw)
		if entry == "" {
			continue
		}

		var prefix netip.Prefix
		if strings.Contains(entry, "/") {
			parsed, err := netip.ParsePrefix(entry)
			if err != nil {
				return nil, nil, fmt.Errorf("无效的 IP/CIDR: %s", entry)
			}
			prefix = parsed.Masked()
		} else {
			addr, err := netip.ParseAddr(entry)
			if err != nil {
				return nil, nil, fmt.Errorf("无效的 IP: %s", entry)
			}
			addr = addr.Unmap()
			bits := 128
			if addr.Is4() {
				bits = 32
			}
			prefix = netip.PrefixFrom(addr, bits)
		}

		addr := prefix.Addr().Unmap()
		if addr.Is4() && prefix.Addr().Is4In6() {
			prefix = netip.PrefixFrom(addr, prefix.Bits()-96).Masked()
		}
		text := prefix.String()
		if (addr.Is4() && prefix.Bits() == 32) || (addr.Is6() && prefix.Bits() == 128) {
			text = addr.String()
		}
		if _, ok := seen[text]; ok {
			continue
		}
		seen[text] = struct{}{}
		normalized = append(normalized, text)
		prefixes = append(prefixes, prefix)
	}
	sort.Strings(normalized)
	// Keep prefixes in the same deterministic order as the returned text.
	prefixByText := make(map[string]netip.Prefix, len(prefixes))
	for _, prefix := range prefixes {
		addr := prefix.Addr().Unmap()
		text := prefix.String()
		if (addr.Is4() && prefix.Bits() == 32) || (addr.Is6() && prefix.Bits() == 128) {
			text = addr.String()
		}
		prefixByText[text] = prefix
	}
	prefixes = prefixes[:0]
	for _, text := range normalized {
		prefixes = append(prefixes, prefixByText[text])
	}
	return normalized, prefixes, nil
}

func parseManagedServerIP(raw string) (netip.Addr, bool) {
	addr, err := netip.ParseAddr(strings.TrimSpace(strings.Trim(raw, "[]")))
	if err != nil {
		return netip.Addr{}, false
	}
	addr = addr.Unmap()
	if !addr.IsValid() || addr.IsUnspecified() {
		return netip.Addr{}, false
	}
	return addr, true
}

// RefreshManaged refreshes the in-memory list from Agent-reported public IPs.
// It returns true when the effective automatic list changed.
func (w *SubscriptionIPWhitelist) RefreshManaged(ctx context.Context) (bool, error) {
	if w == nil || w.repo == nil {
		return false, nil
	}
	servers, err := w.repo.ListRemoteServers(ctx)
	if err != nil {
		return false, err
	}
	next := make(map[netip.Addr]ManagedServerWhitelistEntry)
	for _, server := range servers {
		if server.IsFederated {
			continue
		}
		for _, candidate := range []struct {
			raw    string
			family string
		}{
			{raw: server.IPAddress, family: "IPv4"},
			{raw: server.IPAddressV6, family: "IPv6"},
		} {
			addr, ok := parseManagedServerIP(candidate.raw)
			if !ok {
				continue
			}
			family := candidate.family
			if addr.Is4() {
				family = "IPv4"
			} else {
				family = "IPv6"
			}
			next[addr] = ManagedServerWhitelistEntry{
				ServerID: server.ID,
				Name:     server.Name,
				IP:       addr.String(),
				Family:   family,
			}
		}
	}

	w.mu.Lock()
	changed := !managedWhitelistEqual(w.managed, next)
	w.managed = next
	w.mu.Unlock()
	return changed, nil
}

func managedWhitelistEqual(a, b map[netip.Addr]ManagedServerWhitelistEntry) bool {
	if len(a) != len(b) {
		return false
	}
	for addr, left := range a {
		right, ok := b[addr]
		if !ok || left != right {
			return false
		}
	}
	return true
}

func (w *SubscriptionIPWhitelist) Contains(ip string) bool {
	if w == nil {
		return false
	}
	addr, err := netip.ParseAddr(strings.TrimSpace(ip))
	if err != nil {
		return false
	}
	addr = addr.Unmap()

	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, prefix := range w.manual {
		if prefix.Contains(addr) {
			return true
		}
	}
	if w.autoManaged {
		_, ok := w.managed[addr]
		return ok
	}
	return false
}

func (w *SubscriptionIPWhitelist) Update(ctx context.Context, manual []string, autoManaged bool) error {
	if w == nil || w.repo == nil {
		return fmt.Errorf("subscription IP whitelist is unavailable")
	}
	normalized, prefixes, err := normalizeSubscriptionWhitelist(manual)
	if err != nil {
		return err
	}
	raw, err := json.Marshal(normalized)
	if err != nil {
		return err
	}
	if err := w.repo.SetSystemSettings(ctx, map[string]string{
		subscriptionIPWhitelistKey:   string(raw),
		subscriptionAutoWhitelistKey: strconv.FormatBool(autoManaged),
	}); err != nil {
		return err
	}

	w.mu.Lock()
	w.manualText = normalized
	w.manual = prefixes
	w.autoManaged = autoManaged
	w.mu.Unlock()
	_, err = w.RefreshManaged(ctx)
	return err
}

func (w *SubscriptionIPWhitelist) Snapshot() SubscriptionIPWhitelistSnapshot {
	if w == nil {
		return SubscriptionIPWhitelistSnapshot{AutoIncludeManagedServers: true}
	}
	w.mu.RLock()
	defer w.mu.RUnlock()

	snapshot := SubscriptionIPWhitelistSnapshot{
		ManualEntries:             append([]string{}, w.manualText...),
		AutoIncludeManagedServers: w.autoManaged,
		ManagedServers:            make([]ManagedServerWhitelistEntry, 0, len(w.managed)),
	}
	for _, entry := range w.managed {
		snapshot.ManagedServers = append(snapshot.ManagedServers, entry)
	}
	sort.Slice(snapshot.ManagedServers, func(i, j int) bool {
		if snapshot.ManagedServers[i].ServerID != snapshot.ManagedServers[j].ServerID {
			return snapshot.ManagedServers[i].ServerID < snapshot.ManagedServers[j].ServerID
		}
		return snapshot.ManagedServers[i].IP < snapshot.ManagedServers[j].IP
	})
	return snapshot
}

func (w *SubscriptionIPWhitelist) Start(ctx context.Context, onChange func()) {
	ticker := time.NewTicker(managedWhitelistRefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if changed, err := w.RefreshManaged(ctx); err == nil && changed && onChange != nil {
				onChange()
			}
		}
	}
}
