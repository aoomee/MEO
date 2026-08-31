package handler

import "testing"

func TestIsWholeNodeRuleForInbound(t *testing.T) {
	tests := []struct {
		name string
		rule map[string]any
		want bool
	}{
		{"plain outbound", map[string]any{"inboundTag": []any{"in-a"}, "outboundTag": "landing-in-a-1"}, true},
		{"plain balancer", map[string]any{"inboundTag": []any{"in-a"}, "balancerTag": "bal-a"}, true},
		{"another inbound", map[string]any{"inboundTag": []any{"in-b"}, "outboundTag": "landing-in-b-1"}, false},
		{"domain split", map[string]any{"inboundTag": []any{"in-a"}, "domain": []any{"example.com"}, "outboundTag": "proxy"}, false},
		{"user split", map[string]any{"inboundTag": []any{"in-a"}, "user": []any{"u@example"}, "outboundTag": "proxy"}, false},
		{"no target", map[string]any{"inboundTag": []any{"in-a"}}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isWholeNodeRuleForInbound(test.rule, "in-a"); got != test.want {
				t.Fatalf("isWholeNodeRuleForInbound() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestIsOwnedWholeNodeOutboundTag(t *testing.T) {
	if !isOwnedWholeNodeOutboundTag("landing-node-42-123", 42, "in-a") {
		t.Fatal("new node-owned tag should be owned")
	}
	if !isOwnedWholeNodeOutboundTag("landing-node-42-target-99-123", 42, "in-a") {
		t.Fatal("target-aware node-owned tag should be owned")
	}
	if !isOwnedWholeNodeOutboundTag("landing-in-a-123", 42, "in-a") {
		t.Fatal("legacy inbound-owned tag should be owned")
	}
	if isOwnedWholeNodeOutboundTag("shared-outbound", 42, "in-a") {
		t.Fatal("shared outbound must not be owned")
	}
}

func TestWholeNodeOutboundTargetNodeID(t *testing.T) {
	id, ok := wholeNodeOutboundTargetNodeID("landing-node-42-target-99-123456")
	if !ok || id != 99 {
		t.Fatalf("target id = %d, %v; want 99, true", id, ok)
	}
	for _, tag := range []string{
		"landing-node-42-123456",
		"landing-node-42-target-x-123456",
		"shared-outbound",
	} {
		if _, ok := wholeNodeOutboundTargetNodeID(tag); ok {
			t.Fatalf("tag %q must not expose a target node id", tag)
		}
	}
}

func TestWholeNodeOutboundReferencesNode(t *testing.T) {
	target := outboundTarget{addrSet: map[string]bool{"old.example.com": true}, port: 443}
	oldOutbound := mustObMap(t, `{"tag":"landing-node-7-123","protocol":"vless","settings":{"vnext":[{"address":"old.example.com","port":443}]}}`)
	if !wholeNodeOutboundReferencesNode("landing-node-7-123", oldOutbound, 9, target, true) {
		t.Fatal("legacy managed outbound should fall back to endpoint matching")
	}
	if !wholeNodeOutboundReferencesNode("landing-node-7-target-9-123", map[string]any{}, 9, target, false) {
		t.Fatal("target-aware outbound should match by node id")
	}
	if wholeNodeOutboundReferencesNode("landing-node-7-target-10-123", oldOutbound, 9, target, true) {
		t.Fatal("a target-aware outbound for another node must not fall back to endpoint matching")
	}
	if wholeNodeOutboundReferencesNode("manual-outbound", oldOutbound, 9, target, true) {
		t.Fatal("an unmanaged outbound must not be updated")
	}
}
