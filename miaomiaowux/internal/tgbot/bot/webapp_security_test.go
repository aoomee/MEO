package bot

import (
	"net/http/httptest"
	"testing"
)

func TestClientIPDoesNotTrustForwardingHeaders(t *testing.T) {
	r := httptest.NewRequest("GET", "http://panel/api/tg-webapp/me", nil)
	r.RemoteAddr = "203.0.113.8:1234"
	r.Header.Set("X-Real-IP", "198.51.100.1")
	r.Header.Set("X-Forwarded-For", "198.51.100.2")
	if got := clientIP(r); got != "203.0.113.8" {
		t.Fatalf("clientIP() = %q, want direct peer address", got)
	}
}
