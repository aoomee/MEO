package handler

import (
	"strings"
	"testing"

	"miaomiaowux/internal/storage"
)

func TestMieruCredentialsAreInboundUnique(t *testing.T) {
	a, _, err := generateCredential("mieru", storage.User{Username: "alice"}, "", "in-1")
	if err != nil {
		t.Fatal(err)
	}
	b, _, err := generateCredential("mieru", storage.User{Username: "alice"}, "", "in-2")
	if err != nil {
		t.Fatal(err)
	}
	if a["username"] == b["username"] {
		t.Fatalf("same user on two mieru inbounds must not share username: %v", a["username"])
	}
	if a["username"] != a["email"] || b["username"] != b["email"] {
		t.Fatalf("mieru username should follow unique email, got %v / %v", a["username"], b["username"])
	}
}

func TestRoutedMieruChildrenDoNotShareAccount(t *testing.T) {
	one, _, err := generateRoutedClientCred("mieru", "", "alice__p1__jp")
	if err != nil {
		t.Fatal(err)
	}
	two, _, err := generateRoutedClientCred("mieru", "", "alice__p1__us")
	if err != nil {
		t.Fatal(err)
	}
	if one["username"] != "alice__p1__jp" || two["username"] != "alice__p1__us" {
		t.Fatalf("routed mieru usernames = %v / %v", one["username"], two["username"])
	}
}

func TestRoutedSS2022UsesInboundKeyLength(t *testing.T) {
	node := storage.Node{
		Protocol:     "shadowsocks",
		ParsedConfig: `{"protocol":"shadowsocks","settings":{"method":"2022-blake3-aes-128-gcm","password":"server"}}`,
	}
	method := protocolMethodFromNode(node)
	if method != "2022-blake3-aes-128-gcm" {
		t.Fatalf("method = %q", method)
	}
	cred, _, err := generateRoutedClientCred("shadowsocks", method, "admin__r1__ss")
	if err != nil {
		t.Fatal(err)
	}
	pass, _ := cred["password"].(string)
	if pass == "" {
		t.Fatal("empty ss2022 password")
	}
	if len(pass) < 20 {
		t.Fatalf("ss2022 password too short (likely uuid): %q", pass)
	}
}

func TestGenerateLegacyShadowsocksCredentialCarriesMethod(t *testing.T) {
	cred, _, err := generateCredential(
		"shadowsocks",
		storage.User{Username: "alice"},
		"aes-256-gcm",
		"ss-8388",
	)
	if err != nil {
		t.Fatal(err)
	}
	if got := cred["method"]; got != "aes-256-gcm" {
		t.Fatalf("legacy method = %v, want aes-256-gcm", got)
	}
	if password, _ := cred["password"].(string); password == "" {
		t.Fatal("legacy user password is empty")
	}
}

func TestApplyShadowsocksCredentialToProxy(t *testing.T) {
	t.Run("legacy uses user password without parent", func(t *testing.T) {
		proxy := map[string]interface{}{
			"cipher":   "aes-128-gcm",
			"password": "ignored-parent",
		}
		applyShadowsocksCredentialToProxy(proxy, map[string]interface{}{
			"method": "aes-256-gcm", "password": "legacy-user",
		})
		if got := proxy["password"]; got != "legacy-user" {
			t.Fatalf("legacy password = %v, want user password only", got)
		}
		if got := proxy["cipher"]; got != "aes-256-gcm" {
			t.Fatalf("legacy cipher = %v, want per-user method", got)
		}
	})

	t.Run("2022 composes server and user passwords", func(t *testing.T) {
		proxy := map[string]interface{}{
			"cipher":   "2022-blake3-aes-128-gcm",
			"password": "server-key:old-user-key",
		}
		applyShadowsocksCredentialToProxy(proxy, map[string]interface{}{
			"password": "new-user-key",
		})
		if got := proxy["password"]; got != "server-key:new-user-key" {
			t.Fatalf("SS2022 password = %v, want server-key:new-user-key", got)
		}
	})
}

func TestInboundToClashProxyShadowsocksPasswordSemantics(t *testing.T) {
	h := newTestRMH()
	tests := []struct {
		name       string
		settings   map[string]interface{}
		wantCipher string
		wantPass   string
	}{
		{
			name: "legacy multi-user",
			settings: map[string]interface{}{
				"method": "aes-256-gcm",
				"clients": []interface{}{map[string]interface{}{
					"method": "aes-256-gcm", "password": "legacy-user", "email": "alice",
				}},
			},
			wantCipher: "aes-256-gcm",
			wantPass:   "legacy-user",
		},
		{
			name: "2022 multi-user",
			settings: map[string]interface{}{
				"method":   "2022-blake3-aes-128-gcm",
				"password": "server-key",
				"clients": []interface{}{map[string]interface{}{
					"password": "user-key", "email": "alice",
				}},
			},
			wantCipher: "2022-blake3-aes-128-gcm",
			wantPass:   "server-key:user-key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxy, err := h.inboundToClashProxy(map[string]interface{}{
				"protocol": "shadowsocks", "tag": "ss-8388", "port": float64(8388), "settings": tt.settings,
			}, "1.2.3.4", "server", 0)
			if err != nil {
				t.Fatal(err)
			}
			if got := proxy["cipher"]; got != tt.wantCipher {
				t.Fatalf("cipher = %v, want %s", got, tt.wantCipher)
			}
			if got := proxy["password"]; got != tt.wantPass {
				t.Fatalf("password = %v, want %s", got, tt.wantPass)
			}
		})
	}
}

func TestBuildInboundShadowsocksFamilies(t *testing.T) {
	t.Run("legacy has no parent password", func(t *testing.T) {
		inbound, creds, err := buildInbound(&buildInboundRequest{
			Protocol: "shadowsocks", Port: 8388, Method: "chacha20-ietf-poly1305", Password: "user-secret",
		})
		if err != nil {
			t.Fatal(err)
		}
		settings := inbound["settings"].(map[string]interface{})
		if _, exists := settings["password"]; exists {
			t.Fatalf("legacy multi-user config must not contain parent password: %#v", settings)
		}
		client := settings["clients"].([]interface{})[0].(map[string]interface{})
		if client["method"] != "chacha20-ietf-poly1305" || client["password"] != "user-secret" {
			t.Fatalf("unexpected legacy client: %#v", client)
		}
		if creds["password"] != "user-secret" {
			t.Fatalf("legacy client credential = %#v", creds)
		}
	})

	t.Run("2022 returns composed client password", func(t *testing.T) {
		inbound, creds, err := buildInbound(&buildInboundRequest{
			Protocol: "shadowsocks", Port: 8388, Method: "2022-blake3-aes-256-gcm", Password: "user-key", Network: "udp",
		})
		if err != nil {
			t.Fatal(err)
		}
		settings := inbound["settings"].(map[string]interface{})
		serverKey, _ := settings["password"].(string)
		if serverKey == "" {
			t.Fatal("SS2022 server key is empty")
		}
		if got := creds["password"]; got != serverKey+":user-key" {
			t.Fatalf("SS2022 client password = %v, want composed password", got)
		}
		if strings.Contains(serverKey, ":") {
			t.Fatalf("server key must be a single segment: %q", serverKey)
		}
		if got := settings["network"]; got != "udp" {
			t.Fatalf("SS2022 network = %v, want udp", got)
		}
	})

	t.Run("rejects invalid network", func(t *testing.T) {
		_, _, err := buildInbound(&buildInboundRequest{
			Protocol: "shadowsocks", Port: 8388, Network: "quic",
		})
		if err == nil || !strings.Contains(err.Error(), "network") {
			t.Fatalf("expected network validation error, got %v", err)
		}
	})
}
