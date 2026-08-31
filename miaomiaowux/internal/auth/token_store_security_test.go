package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueryTokenAllowedOnlyForBrowserStreamingEndpoints(t *testing.T) {
	for _, path := range []string{"/api/ws/dashboard", "/api/admin/update/apply-sse"} {
		r := httptest.NewRequest(http.MethodGet, path+"?token=secret", nil)
		if !queryTokenAllowed(r) {
			t.Fatalf("query token unexpectedly rejected for %s", path)
		}
	}
	r := httptest.NewRequest(http.MethodGet, "/api/user/profile?token=secret", nil)
	if queryTokenAllowed(r) {
		t.Fatal("query token accepted for normal API endpoint")
	}
}
