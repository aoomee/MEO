package storage

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestPackageAssignmentsIsolateSamePackageTraffic(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()
	if err := repo.CreateUser(ctx, "multi-user", "", "", "hash", RoleUser, ""); err != nil {
		t.Fatal(err)
	}
	server := &RemoteServer{Name: "multi-server", Token: "multi-token", IPAddress: "127.0.0.1"}
	if err := repo.CreateRemoteServer(ctx, server); err != nil {
		t.Fatal(err)
	}
	packageID, err := repo.CreatePackage(ctx, Package{Name: "multi-package", CycleDays: 30, TrafficLimitBytes: 1 << 30})
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	if err := repo.AssignPackageToUser(ctx, "multi-user", packageID, now, now.AddDate(0, 1, 0), true, 1); err != nil {
		t.Fatal(err)
	}
	legacy, err := repo.GetPrimaryUserPackageAssignment(ctx, "multi-user")
	if err != nil || legacy == nil || !legacy.Legacy {
		t.Fatalf("legacy assignment not mirrored: assignment=%+v err=%v", legacy, err)
	}
	legacyEmail := "multi-user__vless-443"
	if err := repo.SaveUserInboundConfig(ctx, UserInboundConfig{Username: "multi-user", ServerID: server.ID, InboundTag: "vless-443", Protocol: "vless", CredentialJSON: `{"email":"` + legacyEmail + `"}`}); err != nil {
		t.Fatal(err)
	}
	first, err := repo.CreateUserPackageAssignment(ctx, "multi-user", packageID, now, now.AddDate(0, 2, 0), true, 2)
	if err != nil {
		t.Fatal(err)
	}
	second, err := repo.CreateUserPackageAssignment(ctx, "multi-user", packageID, now, now.AddDate(0, 3, 0), true, 3)
	if err != nil {
		t.Fatal(err)
	}
	firstEmail := PackageAssignmentCredentialEmail(first.ID, "multi-user", "vless-443")
	secondEmail := PackageAssignmentCredentialEmail(second.ID, "multi-user", "vless-443")
	if firstEmail == secondEmail {
		t.Fatal("same node in two assignments reused an xray identity")
	}
	for _, config := range []PackageAssignmentInboundConfig{
		{AssignmentID: first.ID, Username: "multi-user", ServerID: server.ID, InboundTag: "vless-443", Protocol: "vless", Email: firstEmail, CredentialJSON: `{"email":"` + firstEmail + `"}`},
		{AssignmentID: second.ID, Username: "multi-user", ServerID: server.ID, InboundTag: "vless-443", Protocol: "vless", Email: secondEmail, CredentialJSON: `{"email":"` + secondEmail + `"}`},
	} {
		if err := repo.SavePackageAssignmentInboundConfig(ctx, config); err != nil {
			t.Fatal(err)
		}
	}
	for _, sample := range []struct {
		email string
		up    int64
		down  int64
	}{{legacyEmail, 300, 100}, {firstEmail, 120, 80}, {secondEmail, 40, 10}} {
		if err := repo.UpsertUserEmailTraffic(ctx, server.ID, sample.email, 0, 0, false, 1, "multi-user"); err != nil {
			t.Fatal(err)
		}
		if err := repo.UpsertUserEmailTraffic(ctx, server.ID, sample.email, sample.up, sample.down, false, 1, "multi-user"); err != nil {
			t.Fatal(err)
		}
	}
	if used, _ := repo.GetPackageAssignmentBillableTraffic(ctx, first.ID); used != 200 {
		t.Fatalf("first assignment used=%d, want 200", used)
	}
	if used, _ := repo.GetPackageAssignmentBillableTraffic(ctx, legacy.ID); used != 400 {
		t.Fatalf("legacy assignment included independent package traffic: %d", used)
	}
	if used, _ := repo.GetPackageAssignmentBillableTraffic(ctx, second.ID); used != 50 {
		t.Fatalf("second assignment used=%d, want 50", used)
	}
	if err := repo.ResetPackageAssignmentTrafficCycleAt(ctx, first.ID, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if used, _ := repo.GetPackageAssignmentBillableTraffic(ctx, first.ID); used != 0 {
		t.Fatalf("first assignment was not reset: %d", used)
	}
	if used, _ := repo.GetPackageAssignmentBillableTraffic(ctx, second.ID); used != 50 {
		t.Fatalf("reset leaked into second assignment: %d", used)
	}
	if err := repo.ResetPackageAssignmentTrafficCycleAt(ctx, legacy.ID, now.Add(2*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if used, _ := repo.GetPackageAssignmentBillableTraffic(ctx, legacy.ID); used != 0 {
		t.Fatalf("legacy assignment was not reset: %d", used)
	}
	if used, _ := repo.GetPackageAssignmentBillableTraffic(ctx, second.ID); used != 50 {
		t.Fatalf("legacy reset leaked into independent assignment: %d", used)
	}
}

func TestDeletingPrimaryAssignmentPromotesNextPackage(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()
	if err := repo.CreateUser(ctx, "promote-user", "", "", "hash", RoleUser, ""); err != nil {
		t.Fatal(err)
	}
	firstPackage, _ := repo.CreatePackage(ctx, Package{Name: "first-package", CycleDays: 30})
	secondPackage, _ := repo.CreatePackage(ctx, Package{Name: "second-package", CycleDays: 30})
	now := time.Now()
	if err := repo.AssignPackageToUser(ctx, "promote-user", firstPackage, now, now.AddDate(0, 1, 0), false, 1); err != nil {
		t.Fatal(err)
	}
	primary, _ := repo.GetPrimaryUserPackageAssignment(ctx, "promote-user")
	secondary, err := repo.CreateUserPackageAssignment(ctx, "promote-user", secondPackage, now, now.AddDate(0, 2, 0), false, 1)
	if err != nil {
		t.Fatal(err)
	}
	secondaryReset := now.Add(-72 * time.Hour).Truncate(time.Second)
	secondary.LastResetAt = &secondaryReset
	if err := repo.UpdateUserPackageAssignment(ctx, *secondary); err != nil {
		t.Fatal(err)
	}
	if err := repo.DeleteUserPackageAssignment(ctx, primary.ID); err != nil {
		t.Fatal(err)
	}
	promoted, err := repo.GetPrimaryUserPackageAssignment(ctx, "promote-user")
	if err != nil || promoted == nil || promoted.ID != secondary.ID || !promoted.IsPrimary {
		t.Fatalf("secondary assignment was not promoted: %+v err=%v", promoted, err)
	}
	user, err := repo.GetUser(ctx, "promote-user")
	if err != nil {
		t.Fatal(err)
	}
	if user.PackageID != secondPackage {
		t.Fatalf("legacy package mirror=%d, want %d", user.PackageID, secondPackage)
	}
	if user.LastResetAt == nil || !user.LastResetAt.Truncate(time.Second).Equal(secondaryReset) {
		t.Fatalf("promoted reset boundary=%v, want %v", user.LastResetAt, secondaryReset)
	}
}

func TestPackageAssignmentsPostgresCompatibility(t *testing.T) {
	if os.Getenv("MMWX_TEST_POSTGRES") == "" {
		t.Skip("set MMWX_TEST_POSTGRES=1 to run")
	}
	host := os.Getenv("MMWX_TEST_POSTGRES_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := 55432
	if value, err := strconv.Atoi(os.Getenv("MMWX_TEST_POSTGRES_PORT")); err == nil && value > 0 {
		port = value
	}
	repo, err := NewTrafficRepositoryFromConfig(DatabaseConfig{Driver: "postgres", Host: host, Port: port, Database: "mmwx", Username: "mmwx", Password: "mmwx-test", SSLMode: "disable"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	ctx := context.Background()
	suffix := strconv.FormatInt(time.Now().UnixNano(), 10)
	username := "multi-pg-" + suffix
	if err := repo.CreateUser(ctx, username, "", "", "hash", RoleUser, ""); err != nil {
		t.Fatal(err)
	}
	packageID, err := repo.CreatePackage(ctx, Package{Name: "multi-pg-package-" + suffix, CycleDays: 30, TrafficLimitBytes: 1 << 30})
	if err != nil {
		t.Fatal(err)
	}
	server := &RemoteServer{Name: "multi-pg-server-" + suffix, Token: "token-" + suffix, IPAddress: "127.0.0.1"}
	if err := repo.CreateRemoteServer(ctx, server); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = repo.db.ExecContext(context.Background(), `DELETE FROM users WHERE username=?`, username)
		_, _ = repo.db.ExecContext(context.Background(), `DELETE FROM packages WHERE id=?`, packageID)
		_, _ = repo.db.ExecContext(context.Background(), `DELETE FROM remote_servers WHERE id=?`, server.ID)
	})
	now := time.Now()
	if err := repo.AssignPackageToUser(ctx, username, packageID, now, now.AddDate(0, 1, 0), true, 1); err != nil {
		t.Fatal(err)
	}
	assignment, err := repo.CreateUserPackageAssignment(ctx, username, packageID, now, now.AddDate(0, 2, 0), true, 2)
	if err != nil {
		t.Fatal(err)
	}
	email := PackageAssignmentCredentialEmail(assignment.ID, username, "vless-443")
	if err := repo.SavePackageAssignmentInboundConfig(ctx, PackageAssignmentInboundConfig{AssignmentID: assignment.ID, Username: username, ServerID: server.ID, InboundTag: "vless-443", Protocol: "vless", Email: email, CredentialJSON: `{"email":"` + email + `"}`}); err != nil {
		t.Fatal(err)
	}
	if err := repo.UpsertUserEmailTraffic(ctx, server.ID, email, 0, 0, false, 1, username); err != nil {
		t.Fatal(err)
	}
	if err := repo.UpsertUserEmailTraffic(ctx, server.ID, email, 25, 75, false, 1, username); err != nil {
		t.Fatal(err)
	}
	if used, err := repo.GetPackageAssignmentBillableTraffic(ctx, assignment.ID); err != nil || used != 100 {
		t.Fatalf("PostgreSQL assignment traffic=%d err=%v, want 100", used, err)
	}
}
