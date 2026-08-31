package handler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"miaomiaowux/internal/storage"
)

func newTestRMH() *RemoteManageHandler {
	dir, err := os.MkdirTemp("", "mmwx-rmh-")
	if err != nil {
		panic(err)
	}
	repo, err := storage.NewTrafficRepository(filepath.Join(dir, "t.db"))
	if err != nil {
		panic(err)
	}
	return NewRemoteManageHandler(repo, nil)
}

func newResetTestRepo(t *testing.T) *storage.TrafficRepository {
	t.Helper()
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "reset.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	return repo
}

func mustObMap(t *testing.T, raw string) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatal(err)
	}
	return out
}
