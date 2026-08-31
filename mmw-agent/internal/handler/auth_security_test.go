package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"mmw-agent/internal/constants"
)

func agentRequest(token string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "http://agent/", nil)
	r.Header.Set(constants.HeaderUserAgent, constants.AgentUserAgent)
	if token != "" {
		r.Header.Set(constants.HeaderAuthorization, constants.BearerPrefix+token)
	}
	return r
}

func TestAgentHTTPAuthenticationFailsClosed(t *testing.T) {
	if NewManageHandler("", "", "").authenticate(agentRequest("")) {
		t.Fatal("manage handler accepted empty configured token")
	}
	if (&APIHandler{}).authenticate(agentRequest("")) {
		t.Fatal("pull handler accepted empty configured token")
	}
	if silentAuthenticate(agentRequest(""), "") {
		t.Fatal("silent middleware accepted empty configured token")
	}

	if !NewManageHandler("secret", "", "").authenticate(agentRequest("secret")) {
		t.Fatal("manage handler rejected valid token")
	}
	h := &APIHandler{configToken: "secret"}
	if !h.authenticate(agentRequest("secret")) {
		t.Fatal("pull handler rejected valid token")
	}
	if !silentAuthenticate(agentRequest("secret"), "secret") {
		t.Fatal("silent middleware rejected valid token")
	}
}

func TestWarpRPCHeaderCannotBeSpoofedOverHTTP(t *testing.T) {
	h := &WarpHandler{configToken: "secret"}
	r := httptest.NewRequest(http.MethodGet, "http://agent/api/child/warp/status", nil)
	r.Header.Set("X-WS-RPC", "1")
	r.RemoteAddr = "127.0.0.1:1234"
	if h.auth(r) {
		t.Fatal("external request spoofed the internal WS RPC marker")
	}
	r.RemoteAddr = "ws-rpc"
	if !h.auth(r) {
		t.Fatal("internal WS RPC request was rejected")
	}
}
