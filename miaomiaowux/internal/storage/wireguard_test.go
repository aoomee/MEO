package storage

import (
	"context"
	"testing"
)

func TestAllocateWGLeaseSkipsServerIndex(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()
	id, err := repo.CreateWGDevice(ctx, &WGDevice{
		ServerID: 1, InboundTag: "wg-1", IPv4CIDR: "10.6.0.0/24",
		FirstIndex: 1, LastIndex: 5, ListenPort: 51820,
		ServerPrivateKey: "s", ServerPublicKey: "p",
	})
	if err != nil {
		t.Fatal(err)
	}
	lease := &WGLease{DeviceID: id, Username: "u", Email: "u__wg-1", PrivateKey: "k", PublicKey: "K"}
	if _, err := repo.AllocateWGLease(ctx, lease); err != nil {
		t.Fatal(err)
	}
	if lease.HostIndex != 2 {
		t.Fatalf("first client should be index 2 (server occupies 1), got %d ip=%s", lease.HostIndex, lease.IPv4)
	}
	if lease.IPv4 != "10.6.0.2" {
		t.Fatalf("ipv4=%s", lease.IPv4)
	}
}

func TestHostIP(t *testing.T) {
	ip, err := hostIP("10.6.0.0/16", 2)
	if err != nil || ip != "10.6.0.2" {
		t.Fatalf("got %s %v", ip, err)
	}
}
