package handler

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPersistInstalledVersionUsesDataDirectoryParent(t *testing.T) {
	root := t.TempDir()
	t.Setenv("MMWX_DATA_DIR", filepath.Join(root, "data"))
	if err := persistInstalledVersion("v0.4.8-beta.10"); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(root, ".version"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(raw), "v0.4.8-beta.10\n"; got != want {
		t.Fatalf("version marker = %q, want %q", got, want)
	}
}

func TestAllowManagedNginxWithoutVersion(t *testing.T) {
	tests := []struct {
		name      string
		installed bool
		canManage bool
		binary    string
		want      bool
	}{
		{name: "mmwx managed installation", installed: true, canManage: true, binary: "/usr/local/nginx/sbin/nginx", want: true},
		{name: "normalizes managed path", installed: true, canManage: true, binary: "/usr/local/nginx/sbin/../sbin/nginx", want: true},
		{name: "rejects distribution nginx", installed: true, canManage: true, binary: "/usr/sbin/nginx"},
		{name: "rejects unmanaged config", installed: true, binary: "/usr/local/nginx/sbin/nginx"},
		{name: "rejects missing installation", canManage: true, binary: "/usr/local/nginx/sbin/nginx"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allowManagedNginxWithoutVersion(tt.installed, tt.canManage, tt.binary); got != tt.want {
				t.Fatalf("allowManagedNginxWithoutVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
