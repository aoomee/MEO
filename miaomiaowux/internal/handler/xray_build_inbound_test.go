package handler

import "testing"

func TestBuildInboundSniffingDoesNotRewriteDestination(t *testing.T) {
	inbound, _, err := buildInbound(&buildInboundRequest{
		Protocol:  "vmess",
		Port:      10086,
		Transport: "tcp",
		Security:  "none",
		UUID:      "00000000-0000-4000-8000-000000000001",
	})
	if err != nil {
		t.Fatal(err)
	}
	sniffing, ok := inbound["sniffing"].(map[string]any)
	if !ok {
		t.Fatalf("sniffing missing or has unexpected type: %#v", inbound["sniffing"])
	}
	if sniffing["routeOnly"] != true {
		t.Fatalf("sniffing.routeOnly=%#v, want true", sniffing["routeOnly"])
	}
}
