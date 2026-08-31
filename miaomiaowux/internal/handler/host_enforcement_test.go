package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"miaomiaowux/internal/storage"
)

func TestHTTPSHostEnforcerFollowsLocalOnlySettingImmediately(t *testing.T) {
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "host-enforcement.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	ctx := context.Background()
	if err := repo.SetSystemSettings(ctx, map[string]string{
		"master_url":        "https://panel.example.com",
		"master_local_only": "0",
	}); err != nil {
		t.Fatalf("SetSystemSettings: %v", err)
	}

	enforcer := NewHTTPSHostEnforcer(repo)
	h := enforcer.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	request := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "http://203.0.113.10:12889/settings?q=1", nil)
		req.Host = "203.0.113.10:12889"
		req.RemoteAddr = "198.51.100.20:43210"
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec
	}

	if got := request().Code; got != http.StatusNoContent {
		t.Fatalf("disabled local-only status=%d, want %d", got, http.StatusNoContent)
	}
	if err := repo.SetSystemSetting(ctx, "master_local_only", "1"); err != nil {
		t.Fatalf("enable local-only: %v", err)
	}
	enforcer.Refresh()
	rec := request()
	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("enabled local-only status=%d, want %d", rec.Code, http.StatusTemporaryRedirect)
	}
	if got := rec.Header().Get("Location"); got != "https://panel.example.com/settings?q=1" {
		t.Fatalf("redirect location=%q", got)
	}

	if err := repo.SetSystemSetting(ctx, "master_local_only", "0"); err != nil {
		t.Fatalf("disable local-only: %v", err)
	}
	enforcer.Refresh()
	if got := request().Code; got != http.StatusNoContent {
		t.Fatalf("disabled again status=%d, want %d", got, http.StatusNoContent)
	}
}

func TestHTTPSHostEnforcerForcePublicAccessBypassesRedirect(t *testing.T) {
	t.Setenv("MMWX_FORCE_PUBLIC_ACCESS", "1")
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "host-recovery.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	if err := repo.SetSystemSettings(context.Background(), map[string]string{
		"master_url":        "https://panel.example.com",
		"master_local_only": "1",
	}); err != nil {
		t.Fatalf("SetSystemSettings: %v", err)
	}

	h := NewHTTPSHostEnforcer(repo).Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "http://203.0.113.10:12889/", nil)
	req.RemoteAddr = "198.51.100.20:43210"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("force-public status=%d, want %d", rec.Code, http.StatusNoContent)
	}
}
