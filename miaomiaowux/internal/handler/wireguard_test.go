package handler

import (
	"testing"
)

func TestInferWGPoolCIDR(t *testing.T) {
	if got := inferWGPoolCIDR("10.6.0.1/32"); got != "10.6.0.0/24" {
		t.Fatalf("got %s", got)
	}
	if got := inferWGPoolCIDR("10.7.1.0/24"); got != "10.7.1.0/24" {
		t.Fatalf("got %s", got)
	}
	if got := inferWGPoolCIDR("10.8.2.9"); got != "10.8.2.0/24" {
		t.Fatalf("got %s", got)
	}
}

func TestWGBuildInitiationSize(t *testing.T) {
	priv, pub, err := wgGenKeypair()
	if err != nil {
		t.Fatal(err)
	}
	pkt, err := wgBuildInitiation(priv, pub)
	if err != nil {
		t.Fatal(err)
	}
	if len(pkt) != 148 {
		t.Fatalf("initiation must be 148 bytes, got %d", len(pkt))
	}
	if pkt[0] != 1 {
		t.Fatalf("type=%d", pkt[0])
	}
}

func TestClashProxyFromWGInbound(t *testing.T) {
	priv, _, err := wgGenKeypair()
	if err != nil {
		t.Fatal(err)
	}
	proxy := clashProxyFromWGInbound(map[string]interface{}{
		"settings": map[string]interface{}{"secretKey": priv, "mtu": 1280.0},
	}, "1.2.3.4", 51820, "wg-in", "s1")
	if proxy["type"] != "wireguard" || proxy["server"] != "1.2.3.4" || proxy["port"] != 51820 {
		t.Fatalf("%v", proxy)
	}
	if proxy["public-key"] == "" {
		t.Fatal("expected derived public-key")
	}
	if proxy["mtu"] != 1280 {
		t.Fatalf("mtu=%v", proxy["mtu"])
	}
}

func TestIsWireGuardProtocol(t *testing.T) {
	if !isWireGuardProtocol("WireGuard") || isWireGuardProtocol("vless") {
		t.Fatal("protocol detect")
	}
}
