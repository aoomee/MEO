package handler

import "testing"

func TestConnectionCountReadinessDoesNotTreatMissingStatusAsZero(t *testing.T) {
	const serverID int64 = 991337
	connCountsByServer.Store(serverID, map[string]int64{})
	defer connCountsByServer.Delete(serverID)
	defer connCountReadyByServer.Delete(serverID)

	ready, notReady := ConnectionCountReadiness()
	if ready {
		t.Fatal("missing readiness was reported ready")
	}
	found := false
	for _, id := range notReady {
		if id == serverID {
			found = true
		}
	}
	if !found {
		t.Fatal("not-ready server id was omitted")
	}

	connCountReadyByServer.Store(serverID, true)
	ready, notReady = ConnectionCountReadiness()
	if !ready || len(notReady) != 0 {
		t.Fatalf("ready status = %v, notReady=%v", ready, notReady)
	}
}
