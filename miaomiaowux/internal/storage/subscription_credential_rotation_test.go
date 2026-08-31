package storage

import (
	"context"
	"testing"
	"time"
)

func TestCommitSubscriptionCredentialRotationInvalidatesLinksAndUpdatesCredentials(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()
	const username = "rotate-user"
	if err := repo.CreateUser(ctx, username, "", username, "hash", RoleUser, ""); err != nil {
		t.Fatal(err)
	}
	_, err := repo.GetOrCreateUserToken(ctx, username)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.UpdateUserCustomShortCode(ctx, username, "leaked-code"); err != nil {
		t.Fatal(err)
	}
	oldToken, err := repo.GetOrCreateUserToken(ctx, username)
	if err != nil {
		t.Fatal(err)
	}
	server := &RemoteServer{Name: "rotate-server", Token: "rotate-token", IPAddress: "127.0.0.1"}
	if err := repo.CreateRemoteServer(ctx, server); err != nil {
		t.Fatal(err)
	}
	pkgID, err := repo.CreatePackage(ctx, Package{Name: "rotate-package", CycleDays: 30, TrafficLimitBytes: 1 << 30})
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	assignment, err := repo.CreateUserPackageAssignment(ctx, username, pkgID, now, now.AddDate(0, 1, 0), false, 1)
	if err != nil {
		t.Fatal(err)
	}
	oldAssignmentCode := assignment.ShortCode
	if err := repo.SaveUserInboundConfig(ctx, UserInboundConfig{Username: username, ServerID: server.ID, InboundTag: "vless-443", Protocol: "vless", CredentialJSON: `{"id":"old-legacy","email":"rotate-user__vless-443"}`}); err != nil {
		t.Fatal(err)
	}
	legacy, err := repo.GetUserInboundConfigs(ctx, username)
	if err != nil || len(legacy) != 1 {
		t.Fatalf("legacy credential=%+v err=%v", legacy, err)
	}
	assignmentEmail := PackageAssignmentCredentialEmail(assignment.ID, username, "vless-443")
	if err := repo.SavePackageAssignmentInboundConfig(ctx, PackageAssignmentInboundConfig{AssignmentID: assignment.ID, Username: username, ServerID: server.ID, InboundTag: "vless-443", Protocol: "vless", Email: assignmentEmail, CredentialJSON: `{"id":"old-assignment","email":"` + assignmentEmail + `"}`}); err != nil {
		t.Fatal(err)
	}
	assignmentConfigs, err := repo.ListPackageAssignmentInboundConfigsByUser(ctx, username)
	if err != nil || len(assignmentConfigs) != 1 {
		t.Fatalf("assignment credential=%+v err=%v", assignmentConfigs, err)
	}
	if err := repo.SavePackageNodeTrafficSuspension(ctx, PackageNodeTrafficSuspension{Username: username, PackageID: pkgID, NodeID: 123, Kind: "physical", CredentialJSON: legacy[0].CredentialJSON}); err != nil {
		t.Fatal(err)
	}

	newToken, err := repo.CommitSubscriptionCredentialRotation(ctx, username, SubscriptionCredentialRotation{
		LegacyInbound:     []CredentialJSONRotation{{ID: legacy[0].ID, CredentialJSON: `{"id":"new-legacy","email":"rotate-user__vless-443"}`}},
		AssignmentInbound: []CredentialJSONRotation{{ID: assignmentConfigs[0].ID, CredentialJSON: `{"id":"new-assignment","email":"` + assignmentEmail + `"}`}},
		Suspensions:       []CredentialJSONReplacement{{Old: legacy[0].CredentialJSON, New: `{"id":"new-legacy","email":"rotate-user__vless-443"}`}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if newToken == "" || newToken == oldToken {
		t.Fatalf("token was not rotated: old=%q new=%q", oldToken, newToken)
	}
	if custom, _ := repo.GetUserCustomShortCode(ctx, username); custom != "" {
		t.Fatalf("leaked custom short code remains valid: %q", custom)
	}
	rotatedAssignment, err := repo.GetUserPackageAssignment(ctx, assignment.ID)
	if err != nil || rotatedAssignment.ShortCode == oldAssignmentCode {
		t.Fatalf("assignment short code was not rotated: old=%q assignment=%+v err=%v", oldAssignmentCode, rotatedAssignment, err)
	}
	legacy, _ = repo.GetUserInboundConfigs(ctx, username)
	if legacy[0].CredentialJSON != `{"id":"new-legacy","email":"rotate-user__vless-443"}` {
		t.Fatalf("legacy credential not updated: %s", legacy[0].CredentialJSON)
	}
	assignmentConfigs, _ = repo.ListPackageAssignmentInboundConfigsByUser(ctx, username)
	if assignmentConfigs[0].CredentialJSON != `{"id":"new-assignment","email":"`+assignmentEmail+`"}` {
		t.Fatalf("assignment credential not updated: %s", assignmentConfigs[0].CredentialJSON)
	}
	suspensions, err := repo.ListPackageNodeTrafficSuspensions(ctx)
	if err != nil || len(suspensions) != 1 || suspensions[0].CredentialJSON != `{"id":"new-legacy","email":"rotate-user__vless-443"}` {
		t.Fatalf("suspended credential snapshot not updated: %+v err=%v", suspensions, err)
	}
}

func TestCommitSubscriptionCredentialRotationRollsBackWhenReferenceDisappears(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()
	const username = "rotate-rollback-user"
	if err := repo.CreateUser(ctx, username, "", username, "hash", RoleUser, ""); err != nil {
		t.Fatal(err)
	}
	oldToken, err := repo.GetOrCreateUserToken(ctx, username)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.CommitSubscriptionCredentialRotation(ctx, username, SubscriptionCredentialRotation{
		LegacyInbound: []CredentialJSONRotation{{ID: 999999, CredentialJSON: `{}`}},
	})
	if err == nil {
		t.Fatal("missing credential row unexpectedly committed")
	}
	current, err := repo.GetOrCreateUserToken(ctx, username)
	if err != nil || current != oldToken {
		t.Fatalf("token changed despite rollback: old=%q current=%q err=%v", oldToken, current, err)
	}
}
