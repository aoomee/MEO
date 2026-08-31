package handler

import (
	"archive/zip"
	"bytes"
	"testing"
)

func backupZipForTest(t *testing.T, names ...string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, name := range names {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.Write([]byte("test"))
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestValidateBackupRejectsTraversalAndDuplicatePaths(t *testing.T) {
	if err := validateBackupArchive(backupZipForTest(t, "data/../../escape")); err == nil {
		t.Fatal("path traversal archive was accepted")
	}
	if err := validateBackupArchive(backupZipForTest(t, "data/a", "data/a")); err == nil {
		t.Fatal("duplicate archive path was accepted")
	}
}
