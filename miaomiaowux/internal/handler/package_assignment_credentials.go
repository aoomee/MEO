package handler

import (
	"context"
	"fmt"

	"miaomiaowux/internal/storage"
)

func credentialEmailForUser(user storage.User, inboundTag string) string {
	if user.PackageAssignmentID > 0 {
		return storage.PackageAssignmentCredentialEmail(user.PackageAssignmentID, user.Username, inboundTag)
	}
	return user.Username + "__" + inboundTag
}

func userForPackageAssignment(base storage.User, assignment storage.UserPackageAssignment) storage.User {
	base.PackageID = assignment.PackageID
	base.PackageEndDate = assignment.PackageEndDate
	base.IsReset = assignment.IsReset
	base.ResetDay = assignment.ResetDay
	base.LastResetAt = assignment.LastResetAt
	base.TrafficLimitOverride = assignment.TrafficLimitOverride
	if assignment.Legacy {
		base.PackageAssignmentID = 0
	} else {
		base.PackageAssignmentID = assignment.ID
	}
	return base
}

func getProvisionedInboundConfig(ctx context.Context, repo *storage.TrafficRepository, user storage.User, serverID int64, inboundTag string) (*storage.UserInboundConfig, error) {
	if user.PackageAssignmentID <= 0 {
		return repo.GetUserInboundConfig(ctx, user.Username, serverID, inboundTag)
	}
	c, err := repo.GetPackageAssignmentInboundConfig(ctx, user.PackageAssignmentID, serverID, inboundTag)
	if err != nil || c == nil {
		return nil, err
	}
	return &storage.UserInboundConfig{ID: c.ID, AssignmentID: c.AssignmentID, Username: c.Username, ServerID: c.ServerID, InboundTag: c.InboundTag, Protocol: c.Protocol, CredentialJSON: c.CredentialJSON, CreatedAt: c.CreatedAt}, nil
}

func saveProvisionedInboundConfig(ctx context.Context, repo *storage.TrafficRepository, cfg storage.UserInboundConfig) error {
	if cfg.AssignmentID <= 0 {
		return repo.SaveUserInboundConfig(ctx, cfg)
	}
	return repo.SavePackageAssignmentInboundConfig(ctx, storage.PackageAssignmentInboundConfig{
		AssignmentID: cfg.AssignmentID, Username: cfg.Username, ServerID: cfg.ServerID,
		InboundTag: cfg.InboundTag, Protocol: cfg.Protocol, CredentialJSON: cfg.CredentialJSON,
	})
}

func deleteProvisionedInboundConfig(ctx context.Context, repo *storage.TrafficRepository, cfg storage.UserInboundConfig) error {
	if cfg.AssignmentID <= 0 {
		return repo.DeleteUserInboundConfig(ctx, cfg.Username, cfg.ServerID, cfg.InboundTag)
	}
	return repo.DeletePackageAssignmentInboundConfig(ctx, cfg.AssignmentID, cfg.ServerID, cfg.InboundTag)
}

func getProvisionedSubaccount(ctx context.Context, repo *storage.TrafficRepository, user storage.User, routedNodeID int64) (*storage.UserSubaccount, error) {
	if user.PackageAssignmentID <= 0 {
		return repo.GetUserSubaccount(ctx, routedNodeID, user.Username)
	}
	s, err := repo.GetPackageAssignmentSubaccount(ctx, user.PackageAssignmentID, routedNodeID)
	if err != nil || s == nil {
		return nil, err
	}
	return &storage.UserSubaccount{ID: s.ID, AssignmentID: s.AssignmentID, Username: s.Username, RoutedNodeID: s.RoutedNodeID, Email: s.Email, CredentialJSON: s.CredentialJSON, IsActive: s.IsActive, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt}, nil
}

func claimProvisionedSubaccount(ctx context.Context, repo *storage.TrafficRepository, user storage.User, sa storage.UserSubaccount) (*storage.UserSubaccount, error) {
	if user.PackageAssignmentID <= 0 {
		return repo.ClaimUserSubaccount(ctx, sa)
	}
	claimed, err := repo.ClaimPackageAssignmentSubaccount(ctx, storage.PackageAssignmentSubaccount{
		AssignmentID: user.PackageAssignmentID, Username: sa.Username, RoutedNodeID: sa.RoutedNodeID,
		Email: sa.Email, CredentialJSON: sa.CredentialJSON, IsActive: sa.IsActive,
	})
	if err != nil || claimed == nil {
		return nil, err
	}
	return &storage.UserSubaccount{ID: claimed.ID, AssignmentID: claimed.AssignmentID, Username: claimed.Username, RoutedNodeID: claimed.RoutedNodeID, Email: claimed.Email, CredentialJSON: claimed.CredentialJSON, IsActive: claimed.IsActive, CreatedAt: claimed.CreatedAt, UpdatedAt: claimed.UpdatedAt}, nil
}

func upsertProvisionedSubaccount(ctx context.Context, repo *storage.TrafficRepository, user storage.User, sa storage.UserSubaccount) error {
	if user.PackageAssignmentID <= 0 {
		_, err := repo.UpsertUserSubaccount(ctx, sa)
		return err
	}
	return repo.UpsertPackageAssignmentSubaccount(ctx, storage.PackageAssignmentSubaccount{
		AssignmentID: user.PackageAssignmentID, Username: sa.Username, RoutedNodeID: sa.RoutedNodeID,
		Email: sa.Email, CredentialJSON: sa.CredentialJSON, IsActive: sa.IsActive,
	})
}

func setProvisionedSubaccountActive(ctx context.Context, repo *storage.TrafficRepository, sa storage.UserSubaccount, active bool) error {
	if sa.AssignmentID <= 0 {
		return repo.SetSubaccountActive(ctx, sa.ID, active)
	}
	return repo.SetPackageAssignmentSubaccountActive(ctx, sa.ID, active)
}

func packageAssignmentRoutedEmail(user storage.User, suffix string, outbound bool) string {
	if user.PackageAssignmentID <= 0 {
		if outbound {
			return fmt.Sprintf("%s-%s", user.Username, suffix)
		}
		return fmt.Sprintf("%s__%s", user.Username, suffix)
	}
	if outbound {
		return fmt.Sprintf("%s__pkg%d-%s", user.Username, user.PackageAssignmentID, suffix)
	}
	return fmt.Sprintf("%s__pkg%d__%s", user.Username, user.PackageAssignmentID, suffix)
}
