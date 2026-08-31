package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
	"miaomiaowux/internal/auth"
	"miaomiaowux/internal/storage"
)

func TestValidatePackageNodeTrafficLimitsRejectsLimitAbovePackage(t *testing.T) {
	err := validatePackageNodeTrafficLimits(map[int64]float64{7: 101}, []int64{7}, 100)
	if err == nil || !strings.Contains(err.Error(), "不能大于套餐流量额度") {
		t.Fatalf("expected package quota validation error, got %v", err)
	}
}

func TestPackageNodeTrafficNameDefaultsOffAndPersists(t *testing.T) {
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "node-traffic-name.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	ctx := context.Background()
	if packageNodeTrafficNameEnabled(ctx, repo) {
		t.Fatal("node traffic suffix must default to disabled")
	}
	if err := repo.SetSystemSetting(ctx, packageNodeTrafficNameEnabledKey, "1"); err != nil {
		t.Fatalf("SetSystemSetting: %v", err)
	}
	if !packageNodeTrafficNameEnabled(ctx, repo) {
		t.Fatal("stored enabled setting was not applied")
	}
}

func TestPackageNodeTrafficNameSettingHandlers(t *testing.T) {
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "node-traffic-name-handler.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	h := NewSystemSettingsHandler(repo, nil)

	put := httptest.NewRequest(http.MethodPut, "/api/admin/system-settings/package-node-traffic-name", strings.NewReader(`{"enabled":true}`))
	putRecorder := httptest.NewRecorder()
	h.SetPackageNodeTrafficName(putRecorder, put)
	if putRecorder.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, body = %s", putRecorder.Code, putRecorder.Body.String())
	}

	getRecorder := httptest.NewRecorder()
	h.GetPackageNodeTrafficName(getRecorder, httptest.NewRequest(http.MethodGet, "/api/admin/system-settings/package-node-traffic-name", nil))
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("GET status = %d, body = %s", getRecorder.Code, getRecorder.Body.String())
	}
	var response struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(getRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode GET response: %v", err)
	}
	if !response.Enabled {
		t.Fatal("GET did not return the persisted enabled setting")
	}
}

func TestApplyPackageNodeTrafficNameDoesNothingWhenDisabled(t *testing.T) {
	proxy := map[string]any{"name": "Tokyo"}
	pkg := &storage.Package{NodeTrafficLimits: map[int64]float64{7: 10}}
	applyPackageNodeTrafficName(context.Background(), nil, "alice", pkg, storage.Node{ID: 7}, proxy, false)
	if got := proxy["name"]; got != "Tokyo" {
		t.Fatalf("disabled suffix changed node name to %v", got)
	}
}

func TestApplyPackageNodeTrafficNameAddsUsageWhenEnabled(t *testing.T) {
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "node-traffic-name-enabled.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	ctx := context.Background()
	pkg := storage.Package{
		Name:              "traffic-name",
		Nodes:             []int64{7},
		NodeTrafficLimits: map[int64]float64{7: 10},
	}
	pkgID, err := repo.CreatePackage(ctx, pkg)
	if err != nil {
		t.Fatalf("CreatePackage: %v", err)
	}
	pkg.ID = pkgID
	proxy := map[string]any{"name": "Tokyo"}
	applyPackageNodeTrafficName(ctx, repo, "alice", &pkg, storage.Node{ID: 7}, proxy, true)
	if got := proxy["name"]; got != "Tokyo [0MB/10.00GB]" {
		t.Fatalf("enabled suffix = %v", got)
	}
}

func TestPackageSubscriptionTrafficNameToggle(t *testing.T) {
	t.Chdir("../..")
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "node-traffic-name-subscription.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	ctx := context.Background()
	const username = "traffic-name-user"
	if err := repo.CreateUser(ctx, username, "", username, "hash", storage.RoleUser, ""); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	node, err := repo.CreateNode(ctx, storage.Node{
		Username:    "admin",
		NodeName:    "Tokyo",
		Protocol:    "ss",
		Enabled:     true,
		ClashConfig: `{"name":"Tokyo","type":"ss","server":"192.0.2.1","port":443,"cipher":"aes-128-gcm","password":"test"}`,
	})
	if err != nil {
		t.Fatalf("CreateNode: %v", err)
	}
	pkg := storage.Package{
		Name:              "traffic-name-package",
		CycleDays:         30,
		Nodes:             []int64{node.ID},
		NodeTrafficLimits: map[int64]float64{node.ID: 10},
	}
	pkgID, err := repo.CreatePackage(ctx, pkg)
	if err != nil {
		t.Fatalf("CreatePackage: %v", err)
	}
	now := time.Now()
	if err := repo.AssignPackageToUser(ctx, username, pkgID, now, now.AddDate(0, 1, 0), false, 1); err != nil {
		t.Fatalf("AssignPackageToUser: %v", err)
	}

	h := NewPackageSubscribeHandler(repo)
	getNames := func() []string {
		t.Helper()
		req := httptest.NewRequest(http.MethodGet, "/api/package-subscription", nil)
		req.Header.Set("User-Agent", "Clash.Meta")
		req = req.WithContext(auth.ContextWithUsername(req.Context(), username))
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Fatalf("subscription status = %d, body = %s", recorder.Code, recorder.Body.String())
		}
		var document struct {
			Proxies []map[string]any `yaml:"proxies"`
		}
		if err := yaml.Unmarshal(recorder.Body.Bytes(), &document); err != nil {
			t.Fatalf("decode subscription: %v", err)
		}
		names := make([]string, 0, len(document.Proxies))
		for _, proxy := range document.Proxies {
			if name, ok := proxy["name"].(string); ok {
				names = append(names, name)
			}
		}
		return names
	}

	if names := getNames(); len(names) != 1 || names[0] != "Tokyo" {
		t.Fatalf("default-off names = %v", names)
	}
	if err := repo.SetSystemSetting(ctx, packageNodeTrafficNameEnabledKey, "1"); err != nil {
		t.Fatalf("enable setting: %v", err)
	}
	if names := getNames(); len(names) != 1 || names[0] != "Tokyo [0MB/10.00GB]" {
		t.Fatalf("enabled names = %v", names)
	}
	if err := repo.SetSystemSetting(ctx, packageNodeTrafficNameEnabledKey, "0"); err != nil {
		t.Fatalf("disable setting: %v", err)
	}
	if names := getNames(); len(names) != 1 || names[0] != "Tokyo" {
		t.Fatalf("disabled-again names = %v", names)
	}
}

func TestValidatePackageNodeTrafficLimitsAcceptsBoundedAndUnlimited(t *testing.T) {
	for name, limits := range map[string]map[int64]float64{
		"equal":     {7: 100},
		"below":     {7: 25},
		"unlimited": {7: 0},
	} {
		t.Run(name, func(t *testing.T) {
			if err := validatePackageNodeTrafficLimits(limits, []int64{7}, 100); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
