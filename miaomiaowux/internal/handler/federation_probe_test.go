package handler

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"miaomiaowux/internal/storage"
)

func TestFederationProbePingAndReturnRoutesRoundTrip(t *testing.T) {
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "fed-probe.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	ctx := context.Background()

	owner := storage.RemoteServer{Name: "owner", Token: "tok-owner", Status: storage.RemoteServerStatusConnected, IPAddress: "10.0.0.1"}
	if err := repo.CreateRemoteServer(ctx, &owner); err != nil {
		t.Fatal(err)
	}
	consumer := storage.RemoteServer{Name: "shared", Token: "tok-shared", Status: storage.RemoteServerStatusConnected}
	if err := repo.CreateRemoteServer(ctx, &consumer); err != nil {
		t.Fatal(err)
	}

	if err := repo.SetSystemSetting(ctx, probeDisguisePingTargetsKey, `[{"key":"he-cu-v4","label":"河北联通","isp":"unicom"}]`); err != nil {
		t.Fatal(err)
	}
	if err := repo.UpsertServerReturnRoute(ctx, storage.ServerReturnRoute{
		ServerID: owner.ID, Carrier: "telecom", Region: "北京", RouteType: "CMIN2",
		TestedAt: time.Unix(1_700_000_000, 0),
	}); err != nil {
		t.Fatal(err)
	}

	store := NewProbeMetricsStore(16)
	store.IngestLatency(owner.ID, []ProbeLatencySample{{
		Key: "he-cu-v4", Success: true, LatencyMs: 23, At: 1_700_000_100,
	}})
	store.IngestSys(owner.ID, ProbeSysSnapshot{CPUPct: 12, HasCPU: true, At: 1_700_000_100})

	info := map[string]any{
		"name":   owner.Name,
		"status": "connected",
	}
	attachFederationProbeSnapshot(info, store, owner.ID, loadProbePingMeta(ctx, repo, owner.ID))
	routes, err := repo.ListServerReturnRoutes(ctx, []int64{owner.ID})
	if err != nil {
		t.Fatal(err)
	}
	info["return_routes"] = encodeFederationReturnRoutes(routes[owner.ID])

	ping, _ := info["probe_ping"].([]map[string]any)
	if len(ping) != 1 || ping[0]["key"] != "he-cu-v4" || ping[0]["label"] != "河北联通" {
		t.Fatalf("owner probe_ping = %#v", info["probe_ping"])
	}
	if _, ok := info["probe_sys"].(map[string]any); !ok {
		t.Fatalf("owner probe_sys missing: %#v", info)
	}
	if raw, ok := info["return_routes"].([]map[string]any); !ok || len(raw) != 1 || raw[0]["carrier"] != "telecom" {
		t.Fatalf("owner return_routes = %#v", info["return_routes"])
	}

	wireRaw, err := json.Marshal(info)
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]any
	if err := json.Unmarshal(wireRaw, &wire); err != nil {
		t.Fatal(err)
	}

	consumerStore := NewProbeMetricsStore(16)
	ingestFederationProbeSys(consumer.ID, wire, consumerStore)
	ingestFederationProbePing(consumer.ID, wire, consumerStore)
	ingestFederationReturnRoutes(ctx, repo, consumer.ID, wire)

	view, ok := consumerStore.Snapshot(consumer.ID, 0)
	if !ok {
		t.Fatal("consumer snapshot missing")
	}
	series, ok := view.Latency["he-cu-v4"]
	if !ok || series.CurrentMs != 23 {
		t.Fatalf("consumer latency = %#v", view.Latency)
	}
	meta := consumerStore.PingMeta(consumer.ID)
	if meta["he-cu-v4"].Label != "河北联通" {
		t.Fatalf("consumer ping meta = %#v", meta)
	}
	gotRoutes, err := repo.ListServerReturnRoutes(ctx, []int64{consumer.ID})
	if err != nil {
		t.Fatal(err)
	}
	if list := gotRoutes[consumer.ID]; len(list) != 1 || list[0].RouteType != "CMIN2" {
		t.Fatalf("consumer return routes = %#v", gotRoutes[consumer.ID])
	}
}
