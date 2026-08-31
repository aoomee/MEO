package handler

import (
	"testing"
	"time"

	"miaomiaowux/internal/storage"
)

func TestAgentLongOfflineDurationRequiresOfflineStatus(t *testing.T) {
	now := time.Date(2026, 8, 12, 13, 20, 0, 0, time.UTC)
	stale := now.Add(-time.Hour)
	threshold := 30 * time.Minute

	for _, tc := range []struct {
		name   string
		server storage.RemoteServer
		want   bool
	}{
		{name: "online with stale heartbeat", server: storage.RemoteServer{Status: storage.RemoteServerStatusConnected, LastHeartbeat: &stale}},
		{name: "offline with stale heartbeat", server: storage.RemoteServer{Status: storage.RemoteServerStatusOffline, LastHeartbeat: &stale}, want: true},
		{name: "offline without heartbeat", server: storage.RemoteServer{Status: storage.RemoteServerStatusOffline}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, got := agentLongOfflineDuration(tc.server, now, threshold)
			if got != tc.want {
				t.Fatalf("should notify = %v, want %v", got, tc.want)
			}
		})
	}
}
