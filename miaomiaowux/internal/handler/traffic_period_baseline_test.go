package handler

import (
	"context"
	"testing"
	"time"

	"miaomiaowux/internal/storage"
)

func TestNodeTunnelTargetUsesOriginalEndpointAfterInPlaceRelay(t *testing.T) {
	node := storage.Node{
		ClashConfig:     `{"name":"n","type":"vless","server":"relay.example","port":32000}`,
		RelayOrigServer: "origin.example",
		RelayOrigPort:   443,
	}
	target, ok := nodeTunnelTarget(context.Background(), nil, &node)
	if !ok {
		t.Fatal("expected tunnel target")
	}
	if target.port != 443 || !target.addrSet["origin.example"] {
		t.Fatalf("target=%+v, want origin.example:443", target)
	}
}

func TestPublicProbeTrafficAccountingReconcilesConfiguredPeriod(t *testing.T) {
	now := time.Date(2026, time.August, 10, 12, 0, 0, 0, time.UTC)
	daily := []probeDailyTraffic{
		{Date: "2026-07-31", Uplink: 999, Downlink: 999},
		{Date: "2026-08-01", Uplink: 60, Downlink: 20},
		{Date: "2026-08-02", Uplink: 40, Downlink: 30},
	}
	for _, tc := range []struct {
		name       string
		mode       string
		used       int64
		adjustment int64
	}{
		{name: "both plus adjustment", mode: "both", used: 175, adjustment: 25},
		{name: "both with negative adjustment", mode: "both", used: 140, adjustment: -10},
		{name: "upload plus adjustment", mode: "upload", used: 120, adjustment: 20},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ps := probeServer{DailyTraffic: daily}
			server := storage.RemoteServer{TrafficResetDay: 1, TrafficLimit: 500, TrafficStatsMode: tc.mode, TrafficSource: "system"}
			fillProbeTrafficAccounting(&ps, server, tc.used, now)
			if ps.TrafficUsed == nil || *ps.TrafficUsed != tc.used || ps.TrafficAdjustment == nil || *ps.TrafficAdjustment != tc.adjustment {
				t.Fatalf("unexpected used/adjustment: used=%v adjustment=%v", ps.TrafficUsed, ps.TrafficAdjustment)
			}
			if ps.TrafficUsedUp == nil || *ps.TrafficUsedUp != 100 || ps.TrafficUsedDown == nil || *ps.TrafficUsedDown != 50 {
				t.Fatalf("unexpected raw period traffic: up=%v down=%v", ps.TrafficUsedUp, ps.TrafficUsedDown)
			}
			if got := applyProbeTrafficMode(ps.TrafficStatsMode, *ps.TrafficUsedUp, *ps.TrafficUsedDown) + *ps.TrafficAdjustment; got != *ps.TrafficUsed {
				t.Fatalf("accounting identity broken: got %d want %d", got, *ps.TrafficUsed)
			}
			if ps.TrafficUsedScope != "configured_period" || ps.PeriodStart != "2026-08-01" || ps.PeriodEnd != "2026-09-01" {
				t.Fatalf("unexpected period metadata: %+v", ps)
			}
		})
	}
}

func TestPublicProbeTrafficWithoutResetOmitsPeriod(t *testing.T) {
	now := time.Date(2026, time.August, 10, 12, 0, 0, 0, time.UTC)
	ps := probeServer{DailyTraffic: []probeDailyTraffic{{Date: "2026-08-10", Uplink: 10, Downlink: 20}}}
	server := storage.RemoteServer{TrafficResetDay: 0, TrafficUsedOffset: 7, TrafficStatsMode: "both", TrafficSource: "xray"}
	fillProbeTrafficAccounting(&ps, server, 37, now)
	fillDailyTrafficScope(&ps, server, now)
	if ps.TrafficUsedScope != "counter_since_reset" || ps.TrafficAdjustment == nil || *ps.TrafficAdjustment != 7 {
		t.Fatalf("unexpected counter metadata: %+v", ps)
	}
	if ps.PeriodStart != "" || ps.PeriodEnd != "" || ps.TrafficUsedUp != nil || ps.TrafficUsedDown != nil || ps.TrafficUsedTotal != nil {
		t.Fatalf("reset_day=0 must not publish fake period totals: %+v", ps)
	}
	if ps.DailyTrafficScope != "recent_7d" {
		t.Fatalf("daily scope=%q, want recent_7d", ps.DailyTrafficScope)
	}
}

func TestPublicProbeDailyWindowKeepsSevenDaysWhenPeriodIsShorter(t *testing.T) {
	now := time.Date(2026, time.August, 16, 12, 0, 0, 0, time.UTC)
	start, scope := probeDailyTrafficWindow(now, 15)
	if got := start.Format("2006-01-02"); got != "2026-08-10" {
		t.Fatalf("start=%s, want recent seven-day boundary 2026-08-10", got)
	}
	if scope != "configured_period_and_recent_7d" {
		t.Fatalf("scope=%q", scope)
	}
}

func TestPublicProbeBootTrafficKeepsLegacyAliases(t *testing.T) {
	var ps probeServer
	fillProbeBootTraffic(&ps, 123, 456)
	if ps.BootTrafficUp == nil || *ps.BootTrafficUp != 123 || ps.BootTrafficDown == nil || *ps.BootTrafficDown != 456 {
		t.Fatalf("unexpected boot counters: %+v", ps)
	}
	if ps.CumulativeUp == nil || *ps.CumulativeUp != 123 || ps.CumulativeDown == nil || *ps.CumulativeDown != 456 || ps.CumulativeTrafficScope != "current_boot" {
		t.Fatalf("legacy aliases diverged: %+v", ps)
	}
}

func TestSubEmailBaselineFallsBackAsOneCycle(t *testing.T) {
	row := storage.UserEmailTraffic{ServerID: 7, Email: "tom__in", Uplink: 13, Downlink: 75}
	key := "7|tom__in"

	// Baseline is from the previous reset cycle. Both directions must retain
	// the current-cycle values; independently clamping them produced 0 B rows.
	got := subEmailBaseline(row, map[string]int64{key: 100}, map[string]int64{key: 200})
	if got.Uplink != 13 || got.Downlink != 75 {
		t.Fatalf("cross-cycle fallback=%d/%d want 13/75", got.Uplink, got.Downlink)
	}

	// A baseline in the same cycle is still subtracted normally.
	got = subEmailBaseline(row, map[string]int64{key: 3}, map[string]int64{key: 5})
	if got.Uplink != 10 || got.Downlink != 70 {
		t.Fatalf("same-cycle delta=%d/%d want 10/70", got.Uplink, got.Downlink)
	}
}

func TestRoutedConnectionStatsUseSubaccountNodeID(t *testing.T) {
	refs := map[int64]storage.InboundNodeRef{
		101: {InboundTag: "shared-in", NodeID: 101, ParentID: 10, NodeType: "routed"},
		102: {InboundTag: "shared-in", NodeID: 102, ParentID: 10, NodeType: "routed"},
	}
	a := routedRefForSubaccount(storage.ActiveSubaccountForLimiter{RoutedNodeID: 101, InboundTag: "shared-in"}, refs)
	b := routedRefForSubaccount(storage.ActiveSubaccountForLimiter{RoutedNodeID: 102, InboundTag: "shared-in"}, refs)
	if a.NodeID != 101 || b.NodeID != 102 {
		t.Fatalf("routed nodes sharing an inbound were merged: a=%+v b=%+v", a, b)
	}
	if connGroupKey("alice", a.ParentID) != connGroupKey("alice", b.ParentID) {
		t.Fatal("routed nodes on the same physical inbound must still share the quota group")
	}
	if connGroupKey("alice", a.NodeID) == connGroupKey("alice", b.NodeID) {
		t.Fatal("routed nodes must have distinct statistics groups")
	}
}

func TestDailyLedgerHistoryKeepsInstallationFirstDay(t *testing.T) {
	const gb = int64(1 << 30)
	rows := []storage.ServerDailyTraffic{
		{ServerID: 1, Date: "2026-08-06", Uplink: 1 * gb, Downlink: 2 * gb},
		{ServerID: 1, Date: "2026-08-07", Uplink: 3 * gb, Downlink: 4 * gb},
		{ServerID: 2, Date: "2026-08-07", Uplink: 5 * gb, Downlink: 6 * gb},
	}
	got := aggregateDailyLedgerHistory(rows, map[int64]string{1: "both", 2: "both"}, 30)
	if len(got) != 2 || got[0].Date != "2026-08-06" || got[0].UsedGB == nil || *got[0].UsedGB != 3 {
		t.Fatalf("installation first day disappeared: %#v", got)
	}
	if got[1].Date != "2026-08-07" || got[1].UsedGB == nil || *got[1].UsedGB != 18 {
		t.Fatalf("second day aggregate mismatch: %#v", got)
	}
}

func TestDailyLedgerHistoryHonorsServerTrafficMode(t *testing.T) {
	const gb = int64(1 << 30)
	rows := []storage.ServerDailyTraffic{
		{ServerID: 1, Date: "2026-08-07", Uplink: 2 * gb, Downlink: 9 * gb},
		{ServerID: 2, Date: "2026-08-07", Uplink: 7 * gb, Downlink: 3 * gb},
		{ServerID: 3, Date: "2026-08-07", Uplink: 4 * gb, Downlink: 6 * gb},
	}
	got := aggregateDailyLedgerHistory(rows, map[int64]string{1: "upload", 2: "download", 3: "max"}, 30)
	if len(got) != 1 || got[0].UsedGB == nil || *got[0].UsedGB != 11 { // 2 + 3 + 6
		t.Fatalf("traffic modes were ignored: %#v", got)
	}
}

func TestDailyLedgerHistoryIncludesOnlySelectedServers(t *testing.T) {
	const gb = int64(1 << 30)
	rows := []storage.ServerDailyTraffic{
		{ServerID: 1, Date: "2026-08-07", Uplink: 2 * gb, Downlink: 3 * gb},
		{ServerID: 2, Date: "2026-08-07", Uplink: 40 * gb, Downlink: 50 * gb},
	}
	// Presence in the modes map is also the traffic-summary selection. Server 2
	// has ledger data but is deliberately absent and must not affect the chart.
	got := aggregateDailyLedgerHistory(rows, map[int64]string{1: "both"}, 30)
	if len(got) != 1 || got[0].UsedGB == nil || *got[0].UsedGB != 5 {
		t.Fatalf("unselected server leaked into daily history: %#v", got)
	}
}
