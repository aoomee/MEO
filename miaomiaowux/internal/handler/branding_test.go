package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"miaomiaowux/internal/storage"
)

func TestBrandingAdminPersistsBrowserIcon(t *testing.T) {
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "branding.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	h := NewBrandingHandler(repo)

	body, _ := json.Marshal(map[string]string{
		"site_title":  "Custom",
		"brand_title": "Brand",
		"logo_url":    "https://example.com/logo.png",
		"icon_url":    "https://example.com/favicon.ico",
	})
	setRec := httptest.NewRecorder()
	h.AdminSet(setRec, httptest.NewRequest(http.MethodPost, "/api/admin/system-settings/branding", bytes.NewReader(body)))
	if setRec.Code != http.StatusOK {
		t.Fatalf("set status=%d body=%s", setRec.Code, setRec.Body.String())
	}

	getRec := httptest.NewRecorder()
	h.AdminGet(getRec, httptest.NewRequest(http.MethodGet, "/api/admin/system-settings/branding", nil))
	var response struct {
		Branding brandingConfig `json:"branding"`
	}
	if err := json.Unmarshal(getRec.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Branding.IconURL != "https://example.com/favicon.ico" {
		t.Fatalf("icon_url=%q", response.Branding.IconURL)
	}
}

func TestProbeIconSettingValidatesAndPersists(t *testing.T) {
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "probe-icon.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	h := NewSystemSettingsHandler(repo, nil)

	for _, tc := range []struct {
		name       string
		icon       string
		wantStatus int
	}{
		{"https icon", "https://example.com/favicon.png", http.StatusOK},
		{"data icon", "data:image/png;base64,AA==", http.StatusOK},
		{"unsafe scheme", "javascript:alert(1)", http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(map[string]any{"title": "Probe", "icon": tc.icon})
			req := httptest.NewRequest(http.MethodPut, "/api/admin/system-settings/probe-disguise", bytes.NewReader(body))
			rec := httptest.NewRecorder()
			h.SetProbeDisguise(rec, req)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
			}
		})
	}

	got, _ := repo.GetSystemSetting(context.Background(), probeDisguiseIconKey)
	if got != "data:image/png;base64,AA==" {
		t.Fatalf("persisted icon=%q", got)
	}
}

func TestProbePayloadIncludesBrowserIcon(t *testing.T) {
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "probe-payload-icon.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	ctx := context.Background()
	for key, value := range map[string]string{
		probeDisguiseEnabledKey: "1",
		probeDisguiseTitleKey:   "Custom Probe",
		probeDisguiseIconKey:    "https://example.com/probe-icon.png",
	} {
		if err := repo.SetSystemSetting(ctx, key, value); err != nil {
			t.Fatal(err)
		}
	}
	payload, err := NewProbePublicHandler(repo, nil, nil).buildPayload(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if payload["title"] != "Custom Probe" || payload["icon"] != "https://example.com/probe-icon.png" {
		t.Fatalf("unexpected probe branding: title=%v icon=%v", payload["title"], payload["icon"])
	}
}

func TestProbeDisguisePersistSettingsSQLite(t *testing.T) {
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "probe-persist.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	h := NewSystemSettingsHandler(repo, nil)

	getRec := httptest.NewRecorder()
	h.GetProbeDisguise(getRec, httptest.NewRequest(http.MethodGet, "/api/admin/system-settings/probe-disguise", nil))
	var got map[string]any
	if err := json.Unmarshal(getRec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["metrics_persist_supported"] != false {
		t.Fatalf("sqlite should not support persist: %v", got["metrics_persist_supported"])
	}
	if got["metrics_persist_enabled"] != false {
		t.Fatalf("sqlite persist enabled=%v", got["metrics_persist_enabled"])
	}
	if int(got["forward_metrics_retention_days"].(float64)) != 7 || int(got["forward_daily_retention_days"].(float64)) != 30 {
		t.Fatalf("default retention = %v %v", got["forward_metrics_retention_days"], got["forward_daily_retention_days"])
	}

	body, _ := json.Marshal(map[string]any{
		"title":                          "Probe",
		"metrics_persist_enabled":        true,
		"forward_metrics_retention_days": 3,
		"forward_daily_retention_days":   10,
	})
	setRec := httptest.NewRecorder()
	h.SetProbeDisguise(setRec, httptest.NewRequest(http.MethodPut, "/api/admin/system-settings/probe-disguise", bytes.NewReader(body)))
	if setRec.Code != http.StatusOK {
		t.Fatalf("set status=%d body=%s", setRec.Code, setRec.Body.String())
	}

	getRec = httptest.NewRecorder()
	h.GetProbeDisguise(getRec, httptest.NewRequest(http.MethodGet, "/api/admin/system-settings/probe-disguise", nil))
	if err := json.Unmarshal(getRec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["metrics_persist_enabled"] != false {
		t.Fatalf("sqlite must ignore persist enable: %v", got)
	}
	if int(got["forward_metrics_retention_days"].(float64)) != 3 || int(got["forward_daily_retention_days"].(float64)) != 10 {
		t.Fatalf("retention not saved: %v", got)
	}
}
