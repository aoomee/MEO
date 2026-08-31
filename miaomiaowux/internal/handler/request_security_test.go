package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"miaomiaowux/internal/storage"
)

func resetTrustedProxiesForTest(t *testing.T, value string) {
	t.Helper()
	old, had := os.LookupEnv("MMWX_TRUSTED_PROXIES")
	if err := os.Setenv("MMWX_TRUSTED_PROXIES", value); err != nil {
		t.Fatal(err)
	}
	trustedProxyOnce = sync.Once{}
	trustedProxyNets = nil
	t.Cleanup(func() {
		if had {
			_ = os.Setenv("MMWX_TRUSTED_PROXIES", old)
		} else {
			_ = os.Unsetenv("MMWX_TRUSTED_PROXIES")
		}
		trustedProxyOnce = sync.Once{}
		trustedProxyNets = nil
	})
}

func TestGetClientIPIgnoresUntrustedForwardingHeaders(t *testing.T) {
	resetTrustedProxiesForTest(t, "")
	r := httptest.NewRequest(http.MethodGet, "http://panel/", nil)
	r.RemoteAddr = "203.0.113.8:1234"
	r.Header.Set("CF-Connecting-IP", "198.51.100.9")
	r.Header.Set("X-Forwarded-For", "198.51.100.10")
	if got := GetClientIP(r); got != "203.0.113.8" {
		t.Fatalf("GetClientIP() = %q, want direct peer", got)
	}
}

func TestGetClientIPUsesTrustedProxyChainFromRight(t *testing.T) {
	resetTrustedProxiesForTest(t, "10.0.0.0/8")
	r := httptest.NewRequest(http.MethodGet, "http://panel/", nil)
	r.RemoteAddr = "10.0.0.2:1234"
	r.Header.Set("X-Forwarded-For", "192.0.2.99, 198.51.100.7, 10.0.0.3")
	if got := GetClientIP(r); got != "198.51.100.7" {
		t.Fatalf("GetClientIP() = %q, want first untrusted hop from right", got)
	}
}

func TestSameOriginOrNoOriginChecksSchemeAndHost(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "https://panel.example/", nil)
	r.RemoteAddr = "127.0.0.1:1234"
	r.Header.Set("Origin", "https://panel.example")
	if !sameOriginOrNoOrigin(r) {
		t.Fatal("same HTTPS origin should be accepted")
	}
	r.Header.Set("Origin", "http://panel.example")
	if sameOriginOrNoOrigin(r) {
		t.Fatal("downgraded HTTP origin must be rejected")
	}
	r.Header.Set("Origin", "https://attacker.example")
	if sameOriginOrNoOrigin(r) {
		t.Fatal("cross-site origin must be rejected")
	}
	r.Header.Set("Origin", "https://user@panel.example")
	if sameOriginOrNoOrigin(r) {
		t.Fatal("origin userinfo must be rejected")
	}
}

func TestRequestIsHTTPSIgnoresUntrustedForwardedProto(t *testing.T) {
	resetTrustedProxiesForTest(t, "")
	r := httptest.NewRequest(http.MethodGet, "http://panel.example/", nil)
	r.RemoteAddr = "203.0.113.8:1234"
	r.Header.Set("X-Forwarded-Proto", "https")
	if requestIsHTTPS(r) {
		t.Fatal("untrusted X-Forwarded-Proto must not mark request HTTPS")
	}

	r.RemoteAddr = "127.0.0.1:1234"
	if !requestIsHTTPS(r) {
		t.Fatal("loopback trusted X-Forwarded-Proto should mark request HTTPS")
	}
}

func TestProbeSameOriginDoesNotTrustUntrustedForwardedHost(t *testing.T) {
	resetTrustedProxiesForTest(t, "")
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "probe-origin.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	if err := repo.SetSystemSetting(context.Background(), "master_url", "https://panel.example"); err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest(http.MethodGet, "https://panel.example/api/public/probe-ws", nil)
	r.RemoteAddr = "203.0.113.8:1234"
	r.Header.Set("Origin", "https://attacker.example")
	r.Header.Set("X-Forwarded-Host", "attacker.example")
	if probeSameOriginRequest(repo, r) {
		t.Fatal("untrusted X-Forwarded-Host widened the same-origin allow-list")
	}
}

func TestObservedPublicBaseURLIgnoresUntrustedForwardingHeaders(t *testing.T) {
	resetTrustedProxiesForTest(t, "")
	r := httptest.NewRequest(http.MethodGet, "http://panel.example/api/admin/tgbot", nil)
	r.RemoteAddr = "203.0.113.8:1234"
	r.Header.Set("X-Forwarded-Host", "attacker.example")
	r.Header.Set("X-Forwarded-Proto", "https")
	if got := observedPublicBaseURL(r); got != "http://panel.example" {
		t.Fatalf("observedPublicBaseURL() = %q, want direct request origin", got)
	}

	r.RemoteAddr = "127.0.0.1:1234"
	if got := observedPublicBaseURL(r); got != "https://attacker.example" {
		t.Fatalf("trusted proxy base URL = %q, want forwarded origin", got)
	}
}

func TestVerificationDomainValidation(t *testing.T) {
	valid := []string{"example.com", "node-1.example.net", "a.b"}
	for _, domain := range valid {
		if !validVerificationDomain(domain) {
			t.Errorf("validVerificationDomain(%q) = false", domain)
		}
	}
	invalid := []string{"", "127.0.0.1", "https://example.com", "a..example.com", "-bad.example", "bad-.example", "bad_/example"}
	for _, domain := range invalid {
		if validVerificationDomain(domain) {
			t.Errorf("validVerificationDomain(%q) = true", domain)
		}
	}
}

func TestSetupProtectorRejectsRemoteWithoutTokenAndAcceptsCookie(t *testing.T) {
	t.Setenv("MMWX_SETUP_TOKEN", strings.Repeat("a", 32))
	p, err := NewSetupProtector()
	if err != nil {
		t.Fatal(err)
	}
	protected := p.Protect(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))

	remote := httptest.NewRequest(http.MethodPost, "http://panel/api/setup/init", nil)
	remote.RemoteAddr = "203.0.113.8:1234"
	w := httptest.NewRecorder()
	protected.ServeHTTP(w, remote)
	if w.Code != http.StatusForbidden {
		t.Fatalf("unauthorized status = %d, want 403", w.Code)
	}

	authReq := httptest.NewRequest(http.MethodGet, "http://panel/api/setup/authorize?token="+p.Token(), nil)
	authReq.RemoteAddr = "203.0.113.8:1234"
	authW := httptest.NewRecorder()
	p.AuthorizeHandler().ServeHTTP(authW, authReq)
	if authW.Code != http.StatusSeeOther || len(authW.Result().Cookies()) != 1 {
		t.Fatalf("authorize response = %d, cookies=%d", authW.Code, len(authW.Result().Cookies()))
	}

	remote.AddCookie(authW.Result().Cookies()[0])
	w = httptest.NewRecorder()
	protected.ServeHTTP(w, remote)
	if w.Code != http.StatusNoContent {
		t.Fatalf("authorized status = %d, want 204", w.Code)
	}
}

func TestSetupProtectorDoesNotTrustLoopbackReverseProxy(t *testing.T) {
	t.Setenv("MMWX_SETUP_TOKEN", strings.Repeat("b", 32))
	p, err := NewSetupProtector()
	if err != nil {
		t.Fatal(err)
	}
	protected := p.Protect(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	r := httptest.NewRequest(http.MethodPost, "https://panel.example/api/setup/init", nil)
	r.RemoteAddr = "127.0.0.1:1234"
	r.Header.Set("X-Forwarded-For", "198.51.100.7")
	w := httptest.NewRecorder()
	protected.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Fatalf("reverse-proxy request status = %d, want 403", w.Code)
	}

	local := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:12889/api/setup/init", nil)
	local.RemoteAddr = "127.0.0.1:1234"
	w = httptest.NewRecorder()
	protected.ServeHTTP(w, local)
	if w.Code != http.StatusNoContent {
		t.Fatalf("direct loopback request status = %d, want 204", w.Code)
	}
}

func TestChildAPIAuthenticationFailsClosed(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://child/api/child/traffic", nil)
	if NewChildAPIHandler(nil, "").authenticate(r) {
		t.Fatal("empty token must not authenticate")
	}
	h := NewChildAPIHandler(nil, "secret")
	r.Header.Set("Authorization", "Bearer secret")
	if !h.authenticate(r) {
		t.Fatal("valid bearer token was rejected")
	}
}
