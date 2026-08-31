package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestUpdateRemoteServerHeartbeatReturnsPreviousOfflineNotificationState(t *testing.T) {
	repo, err := NewTrafficRepository(filepath.Join(t.TempDir(), "heartbeat-notify.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()
	_, err = repo.db.ExecContext(ctx, `
		INSERT INTO remote_servers (name, token, status, offline_notified)
		VALUES ('notified-server', 'notified-token', 'offline', 1),
		       ('debounced-server', 'debounced-token', 'offline', 0)`)
	if err != nil {
		t.Fatalf("insert remote servers: %v", err)
	}

	for _, tc := range []struct {
		name     string
		token    string
		notified bool
	}{
		{name: "notified offline cycle", token: "notified-token", notified: true},
		{name: "recovered within tolerance", token: "debounced-token", notified: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := repo.UpdateRemoteServerHeartbeatWithRestart(ctx, HeartbeatUpdate{Token: tc.token})
			if err != nil {
				t.Fatalf("UpdateRemoteServerHeartbeatWithRestart: %v", err)
			}
			if result.PreviousStatus != RemoteServerStatusOffline {
				t.Fatalf("PreviousStatus = %q, want %q", result.PreviousStatus, RemoteServerStatusOffline)
			}
			if result.PreviousOfflineNotified != tc.notified {
				t.Fatalf("PreviousOfflineNotified = %v, want %v", result.PreviousOfflineNotified, tc.notified)
			}

			var status string
			var notified bool
			if err := repo.db.QueryRowContext(ctx,
				`SELECT status, offline_notified FROM remote_servers WHERE token = ?`, tc.token,
			).Scan(&status, &notified); err != nil {
				t.Fatalf("query updated server: %v", err)
			}
			if status != RemoteServerStatusConnected || notified {
				t.Fatalf("updated state = (%q, %v), want (%q, false)", status, notified, RemoteServerStatusConnected)
			}
		})
	}
}

func TestUpdateRemoteServerHeartbeatIgnoresSubsecondBootTimeDifferences(t *testing.T) {
	repo, err := NewTrafficRepository(filepath.Join(t.TempDir(), "heartbeat-boot.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()
	boot := time.Date(2026, 8, 17, 14, 0, 0, 987654321, time.UTC)
	xrayBoot := time.Date(2026, 8, 17, 13, 0, 0, 456789123, time.UTC)
	_, err = repo.db.ExecContext(ctx, `
		INSERT INTO remote_servers
			(name, token, status, boot_time, xray_boot_time, boot_count, xray_boot_count)
		VALUES (?, ?, 'connected', ?, ?, 4, 7)`,
		"same-second-server", "same-second-token", boot, xrayBoot)
	if err != nil {
		t.Fatalf("insert remote server: %v", err)
	}

	// HTTP heartbeat serializes Unix seconds, while WS can retain nanoseconds.
	// They represent the same process start and must not increment restart counts.
	bootSeconds := time.Unix(boot.Unix(), 0).UTC()
	xrayBootSeconds := time.Unix(xrayBoot.Unix(), 0).UTC()
	result, err := repo.UpdateRemoteServerHeartbeatWithRestart(ctx, HeartbeatUpdate{
		Token:        "same-second-token",
		BootTime:     &bootSeconds,
		XrayBootTime: &xrayBootSeconds,
	})
	if err != nil {
		t.Fatalf("UpdateRemoteServerHeartbeatWithRestart: %v", err)
	}
	if result.MmwxRestarted || result.XrayRestarted {
		t.Fatalf("same-second boot times reported as restart: %+v", result)
	}
	if result.BootCount != 4 || result.XrayBootCount != 7 {
		t.Fatalf("restart counts changed: mmwx=%d xray=%d", result.BootCount, result.XrayBootCount)
	}

	// Simulate the next WS heartbeat restoring nanosecond precision. Alternating
	// HTTP/WS heartbeats for the same process must remain stable across cycles.
	result, err = repo.UpdateRemoteServerHeartbeatWithRestart(ctx, HeartbeatUpdate{
		Token:        "same-second-token",
		BootTime:     &boot,
		XrayBootTime: &xrayBoot,
	})
	if err != nil {
		t.Fatalf("second UpdateRemoteServerHeartbeatWithRestart: %v", err)
	}
	if result.MmwxRestarted || result.XrayRestarted || result.BootCount != 4 || result.XrayBootCount != 7 {
		t.Fatalf("alternating heartbeat precision changed restart state: %+v", result)
	}

	newBoot := boot.Add(time.Hour)
	newXrayBoot := xrayBoot.Add(time.Hour)
	result, err = repo.UpdateRemoteServerHeartbeatWithRestart(ctx, HeartbeatUpdate{
		Token:        "same-second-token",
		BootTime:     &newBoot,
		XrayBootTime: &newXrayBoot,
	})
	if err != nil {
		t.Fatalf("restart UpdateRemoteServerHeartbeatWithRestart: %v", err)
	}
	if !result.MmwxRestarted || !result.XrayRestarted || result.BootCount != 5 || result.XrayBootCount != 8 {
		t.Fatalf("real restart was not counted exactly once: %+v", result)
	}
}

func TestCheckpointPassiveDoesNotTruncateWAL(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "passive.db")
	repo, err := NewTrafficRepository(dbPath)
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	defer repo.Close()

	conn, err := repo.db.Conn(context.Background())
	if err != nil {
		t.Fatalf("reserve SQLite connection: %v", err)
	}
	defer conn.Close()
	var autoCheckpointPages int
	if err := conn.QueryRowContext(context.Background(), "PRAGMA wal_autocheckpoint").Scan(&autoCheckpointPages); err != nil {
		t.Fatalf("read wal_autocheckpoint: %v", err)
	}
	if autoCheckpointPages != 0 {
		t.Fatalf("wal_autocheckpoint = %d, want 0", autoCheckpointPages)
	}
	if _, err := conn.ExecContext(context.Background(), `CREATE TABLE IF NOT EXISTS wal_probe (id INTEGER PRIMARY KEY, payload BLOB)`); err != nil {
		t.Fatalf("create probe table: %v", err)
	}
	payload := make([]byte, 256*1024)
	for i := 0; i < 8; i++ {
		if _, err := conn.ExecContext(context.Background(), `INSERT INTO wal_probe(payload) VALUES (?)`, payload); err != nil {
			t.Fatalf("insert probe row %d: %v", i, err)
		}
	}

	before, err := os.Stat(dbPath + "-wal")
	if err != nil {
		t.Fatalf("stat WAL before checkpoint: %v", err)
	}
	if before.Size() == 0 {
		t.Fatal("expected a non-empty WAL")
	}
	if _, err := repo.CheckpointPassive(); err != nil {
		t.Fatalf("CheckpointPassive: %v", err)
	}
	after, err := os.Stat(dbPath + "-wal")
	if err != nil {
		t.Fatalf("stat WAL after checkpoint: %v", err)
	}
	if after.Size() != before.Size() {
		t.Fatalf("PASSIVE checkpoint resized WAL: before=%d after=%d", before.Size(), after.Size())
	}
}

func TestNewTrafficRepositoryKeepsSQLiteConnectionsIdle(t *testing.T) {
	repo, err := NewTrafficRepository(filepath.Join(t.TempDir(), "pool.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	defer repo.Close()

	if got := repo.DatabaseConfig().MaxIdleConns; got != sqliteMaxIdleConns {
		t.Fatalf("MaxIdleConns = %d, want %d", got, sqliteMaxIdleConns)
	}
	for i := 0; i < 10; i++ {
		if _, err := repo.db.Exec(`CREATE TABLE IF NOT EXISTS pool_probe (id INTEGER PRIMARY KEY)`); err != nil {
			t.Fatalf("exec %d: %v", i, err)
		}
	}
	stats := repo.db.Stats()
	if stats.Idle == 0 {
		t.Fatalf("SQLite pool retained no idle connection: %+v", stats)
	}
	if stats.MaxIdleClosed != 0 {
		t.Fatalf("SQLite pool closed connections as non-idle: %+v", stats)
	}
}
