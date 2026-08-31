package handler

import (
	"archive/zip"
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"miaomiaowux/internal/storage"
)

func TestRestoreSQLiteBackupRequiresRestart(t *testing.T) {
	if runtime.GOOS == "windows" {
		// 还原时要 rename 覆盖打开中的 mmwx.db;Windows 不允许覆盖打开的文件(Access denied),
		// 而生产是 Linux(rename-over-open 正常)。跳过,避免本机跑测试假失败。
		t.Skip("Windows 不能 rename 覆盖打开中的 db;此路径在 Linux 生产上正常")
	}
	dataDir := t.TempDir()
	databasePath := filepath.Join(dataDir, "mmwx.db")
	repo, err := storage.NewTrafficRepository(databasePath)
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	defer repo.Close()

	var archive bytes.Buffer
	zw := zip.NewWriter(&archive)
	restoredDatabase, err := os.ReadFile(databasePath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if err := addBytesToZip(zw, "data/mmwx.db", restoredDatabase); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	part, err := mw.CreateFormFile("backup", "backup.zip")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write(archive.Bytes())
	_ = mw.Close()
	req := httptest.NewRequest("POST", "/api/admin/backup/restore", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	recorder := httptest.NewRecorder()

	result, err := restoreFromRequest(recorder, req, repo, dataDir)
	if err != nil {
		t.Fatalf("restore failed (HTTP %d): %v; body=%s", recorder.Code, err, recorder.Body.String())
	}
	if !result.databaseSwitched {
		t.Fatal("SQLite restore did not request a process restart")
	}
}
