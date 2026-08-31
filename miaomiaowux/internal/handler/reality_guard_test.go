package handler

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func testRealityGuardConfig() map[string]any {
	return map[string]any{
		"inbounds": []any{
			map[string]any{
				"tag": "vless-reality-443", "port": float64(443), "protocol": "vless",
				"streamSettings": map[string]any{
					"security": "reality",
					"realitySettings": map[string]any{
						"dest": "speed.cloudflare.com:443", "serverNames": []any{"speed.cloudflare.com"},
					},
				},
			},
		},
		"routing": map[string]any{"rules": []any{map[string]any{"type": "field", "inboundTag": []any{"api"}, "outboundTag": "api"}}},
	}
}

func TestMutateRealityGuardConfigEnableAndDisable(t *testing.T) {
	config := testRealityGuardConfig()
	tag := "vless-reality-443"
	if err := mutateRealityGuardConfig(config, tag, true); err != nil {
		t.Fatalf("enable guard: %v", err)
	}
	inbounds := config["inbounds"].([]any)
	if len(inbounds) != 2 {
		t.Fatalf("inbounds=%d, want 2", len(inbounds))
	}
	guard := findInboundByTag(inbounds, realityGuardTag(tag))
	if guard == nil || guard["protocol"] != "tunnel" || guard["listen"] != "127.0.0.1" {
		t.Fatalf("invalid guard inbound: %#v", guard)
	}
	if got := guard["settings"].(map[string]any)["network"]; got != "tcp,udp" {
		t.Fatalf("guard network=%v, want tcp,udp to match Xray normalization", got)
	}
	main := findInboundByTag(inbounds, tag)
	reality := realitySettingsOf(main)
	if got := reality["dest"]; got != "127.0.0.1:39000" {
		t.Fatalf("guarded dest=%v", got)
	}
	rules := config["routing"].(map[string]any)["rules"].([]any)
	if len(rules) != 3 {
		t.Fatalf("rules=%d, want 3", len(rules))
	}
	if got := rules[0].(map[string]any)["outboundTag"]; got != "direct" {
		t.Fatalf("first rule outbound=%v", got)
	}
	if got := rules[1].(map[string]any)["outboundTag"]; got != "block" {
		t.Fatalf("second rule outbound=%v", got)
	}

	if err := mutateRealityGuardConfig(config, tag, false); err != nil {
		t.Fatalf("disable guard: %v", err)
	}
	inbounds = config["inbounds"].([]any)
	if len(inbounds) != 1 || findInboundByTag(inbounds, realityGuardTag(tag)) != nil {
		t.Fatalf("guard inbound was not removed: %#v", inbounds)
	}
	if got := realitySettingsOf(findInboundByTag(inbounds, tag))["dest"]; got != "speed.cloudflare.com:443" {
		t.Fatalf("restored dest=%v", got)
	}
	if got := len(config["routing"].(map[string]any)["rules"].([]any)); got != 1 {
		t.Fatalf("rules after disable=%d", got)
	}
}

func TestMutateRealityGuardConfigUpdatesTargetWithoutDuplicating(t *testing.T) {
	config := testRealityGuardConfig()
	tag := "vless-reality-443"
	if err := mutateRealityGuardConfig(config, tag, true); err != nil {
		t.Fatal(err)
	}
	main := findInboundByTag(config["inbounds"].([]any), tag)
	realitySettingsOf(main)["dest"] = "www.microsoft.com:443"
	realitySettingsOf(main)["serverNames"] = []any{"www.microsoft.com"}
	if err := mutateRealityGuardConfig(config, tag, true); err != nil {
		t.Fatal(err)
	}
	inbounds := config["inbounds"].([]any)
	if len(inbounds) != 2 {
		t.Fatalf("duplicate helper created: %d inbounds", len(inbounds))
	}
	guard := findInboundByTag(inbounds, realityGuardTag(tag))
	settings := guard["settings"].(map[string]any)
	if settings["address"] != "www.microsoft.com" {
		t.Fatalf("helper target=%v", settings["address"])
	}
	rules := config["routing"].(map[string]any)["rules"].([]any)
	if len(rules) != 3 {
		t.Fatalf("duplicate rules created: %d", len(rules))
	}
}

func TestRealityGuardUsesReservedTag(t *testing.T) {
	tag := realityGuardTag("custom-inbound")
	if !isRealityGuardTag(tag) || tag == "custom-inbound" {
		t.Fatalf("unexpected guard tag %q", tag)
	}
}

func TestValidateRealityGuardStealMode(t *testing.T) {
	on, off := true, false
	for _, mode := range []string{"tunnel", "fallback"} {
		if err := validateRealityGuardStealMode(mode, &on); err == nil {
			t.Fatalf("mode %q allowed Reality guard", mode)
		}
		if err := validateRealityGuardStealMode(mode, &off); err != nil {
			t.Fatalf("mode %q rejected disabled guard: %v", mode, err)
		}
	}
	for _, mode := range []string{"", "default"} {
		if err := validateRealityGuardStealMode(mode, &on); err != nil {
			t.Fatalf("mode %q rejected Reality guard: %v", mode, err)
		}
	}
}

func TestFilterInboundsHidesRealityGuardAndRestoresEditFields(t *testing.T) {
	config := testRealityGuardConfig()
	tag := "vless-reality-443"
	if err := mutateRealityGuardConfig(config, tag, true); err != nil {
		t.Fatal(err)
	}
	payload, _ := json.Marshal(map[string]any{"success": true, "inbounds": config["inbounds"]})
	filtered := (&RemoteManageHandler{}).filterInboundsResponse(payload)
	var response struct {
		Inbounds []map[string]any `json:"inbounds"`
	}
	if err := json.Unmarshal(filtered, &response); err != nil {
		t.Fatal(err)
	}
	if len(response.Inbounds) != 1 {
		t.Fatalf("visible inbounds=%d, want 1", len(response.Inbounds))
	}
	if response.Inbounds[0]["reality_guard"] != true {
		t.Fatalf("reality_guard=%v", response.Inbounds[0]["reality_guard"])
	}
	if got := realitySettingsOf(response.Inbounds[0])["dest"]; got != "speed.cloudflare.com:443" {
		t.Fatalf("editable dest=%v", got)
	}
}

func TestDiagnoseActionGuardErrorExplainsReplacedMasterProcess(t *testing.T) {
	original := masterExecutableReplaced
	t.Cleanup(func() { masterExecutableReplaced = original })
	masterExecutableReplaced = func() bool { return true }

	err := diagnoseActionGuardError(errors.New("guard rejected request: frontend challenge invalid or expired"))
	if err == nil || !strings.Contains(err.Error(), "master/guard version mismatch") || !strings.Contains(err.Error(), "重启 mmwx") {
		t.Fatalf("unexpected diagnostic: %v", err)
	}

	masterExecutableReplaced = func() bool { return false }
	originalErr := errors.New("guard unavailable")
	if got := diagnoseActionGuardError(originalErr); got != originalErr {
		t.Fatalf("ordinary Guard error was rewritten: %v", got)
	}
}

func realityFallbackRequest(dest string, names any) map[string]interface{} {
	return map[string]interface{}{
		"action": "add",
		"inbound": map[string]interface{}{
			"protocol": "vless",
			"streamSettings": map[string]interface{}{
				"security": "reality",
				"realitySettings": map[string]interface{}{
					"dest":        dest,
					"serverNames": names,
					"xver":        float64(0),
				},
			},
		},
	}
}

func TestRewriteRealityFallbackForLocalWebsite(t *testing.T) {
	req := realityFallbackRequest("Static.Example.com:443", []interface{}{"static.example.com"})
	changed := rewriteRealityFallbackForLocalWebsites(req, map[string]struct{}{
		"static.example.com": {},
	})
	if !changed {
		t.Fatal("expected local website fallback to be rewritten")
	}
	stream := req["inbound"].(map[string]interface{})["streamSettings"].(map[string]interface{})
	settings := stream["realitySettings"].(map[string]interface{})
	if got := settings["dest"]; got != localRealityNginxDest {
		t.Fatalf("dest = %v, want %s", got, localRealityNginxDest)
	}
	if got := settings["xver"]; got != 1 {
		t.Fatalf("xver = %v, want 1", got)
	}
}

func TestRewriteRealityFallbackKeepsExternalWebsite(t *testing.T) {
	req := realityFallbackRequest("www.microsoft.com:443", []interface{}{"www.microsoft.com"})
	original, _ := json.Marshal(req)
	if rewriteRealityFallbackForLocalWebsites(req, map[string]struct{}{"static.example.com": {}}) {
		t.Fatal("external camouflage destination must not be rewritten")
	}
	after, _ := json.Marshal(req)
	if string(after) != string(original) {
		t.Fatalf("external request changed: %s", after)
	}
}

func TestRewriteRealityFallbackRequiresDestToMatchServerName(t *testing.T) {
	req := realityFallbackRequest("static.example.com:443", []interface{}{"other.example.com"})
	if rewriteRealityFallbackForLocalWebsites(req, map[string]struct{}{"static.example.com": {}}) {
		t.Fatal("mismatched Reality destination must not be rewritten")
	}
}
