package storage

import (
	"context"
	"path/filepath"
	"testing"
)

func TestLicenseUsageExcludesExternalNodes(t *testing.T) {
	repo, err := NewTrafficRepository(filepath.Join(t.TempDir(), "license-usage.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	ctx := context.Background()

	if _, err := repo.CreateNode(ctx, Node{
		Username: "admin", NodeName: "external", Protocol: "vless",
		ClashConfig: `{}`, Enabled: true,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.CreateNode(ctx, Node{
		Username: "admin", NodeName: "managed", Protocol: "vless",
		ClashConfig: `{}`, Enabled: true, OriginalServer: "managed-server", InboundTag: "inbound-1",
	}); err != nil {
		t.Fatal(err)
	}

	_, nodes, _, err := repo.LicenseUsage(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if nodes != 1 {
		t.Fatalf("LicenseUsage nodes = %d, want 1 managed node; external node must be excluded", nodes)
	}
}
