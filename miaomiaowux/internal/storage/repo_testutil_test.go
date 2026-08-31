package storage

import (
	"path/filepath"
	"testing"
)

func newTestRepo(t *testing.T) *TrafficRepository {
	t.Helper()
	repo, err := NewTrafficRepository(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	return repo
}
