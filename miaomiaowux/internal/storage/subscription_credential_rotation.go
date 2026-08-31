package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// CredentialJSONRotation describes one credential row that must be updated as
// part of a subscription credential rotation. The handler prepares the remote
// Xray transaction first; this structure lets the database side commit every
// credential reference and every subscription link in one transaction.
type CredentialJSONRotation struct {
	ID             int64
	CredentialJSON string
}

type NodeCredentialRotation struct {
	ID             int64
	CredentialJSON string
	ClashConfig    string
	RoutedAdmin    bool
}

type CredentialJSONReplacement struct {
	Old string
	New string
}

type SubscriptionCredentialRotation struct {
	LegacyInbound        []CredentialJSONRotation
	AssignmentInbound    []CredentialJSONRotation
	LegacySubaccount     []CredentialJSONRotation
	AssignmentSubaccount []CredentialJSONRotation
	Suspensions          []CredentialJSONReplacement
	Nodes                []NodeCredentialRotation
}

// CommitSubscriptionCredentialRotation atomically updates all stored Xray
// credentials and invalidates every known subscription credential belonging to
// username. Custom short codes are cleared because retaining one would keep an
// already leaked /x/ URL valid. Package assignment short codes are rotated too.
func (r *TrafficRepository) CommitSubscriptionCredentialRotation(ctx context.Context, username string, changes SubscriptionCredentialRotation) (string, error) {
	if r == nil || r.db == nil {
		return "", errors.New("traffic repository not initialized")
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return "", errors.New("username is required")
	}

	for attempt := 0; attempt < 10; attempt++ {
		token := uuid.NewString()
		shortCode, err := generateUserShortCode()
		if err != nil {
			return "", err
		}
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return "", err
		}
		if err = applySubscriptionCredentialChanges(ctx, tx, username, token, shortCode, changes); err == nil {
			err = tx.Commit()
		}
		if err == nil {
			return token, nil
		}
		_ = tx.Rollback()
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "unique") && strings.Contains(lower, "short") {
			continue
		}
		return "", err
	}
	return "", errors.New("failed to generate unique subscription short codes after retries")
}

func applySubscriptionCredentialChanges(ctx context.Context, tx *dialectTx, username, token, shortCode string, changes SubscriptionCredentialRotation) error {
	updateRows := func(query, label string, rows []CredentialJSONRotation) error {
		for _, row := range rows {
			result, err := tx.ExecContext(ctx, query, row.CredentialJSON, row.ID, username)
			if err != nil {
				return fmt.Errorf("update %s %d: %w", label, row.ID, err)
			}
			if count, countErr := result.RowsAffected(); countErr != nil || count != 1 {
				return fmt.Errorf("update %s %d: credential row no longer exists", label, row.ID)
			}
		}
		return nil
	}
	if err := updateRows(`UPDATE user_inbound_configs SET credential_json=? WHERE id=? AND username=?`, "legacy inbound credential", changes.LegacyInbound); err != nil {
		return err
	}
	if err := updateRows(`UPDATE package_assignment_inbound_configs SET credential_json=? WHERE id=? AND username=?`, "package inbound credential", changes.AssignmentInbound); err != nil {
		return err
	}
	if err := updateRows(`UPDATE user_subaccounts SET credential_json=?,updated_at=CURRENT_TIMESTAMP WHERE id=? AND username=?`, "legacy routed credential", changes.LegacySubaccount); err != nil {
		return err
	}
	if err := updateRows(`UPDATE package_assignment_subaccounts SET credential_json=?,updated_at=CURRENT_TIMESTAMP WHERE id=? AND username=?`, "package routed credential", changes.AssignmentSubaccount); err != nil {
		return err
	}
	for _, replacement := range changes.Suspensions {
		if _, err := tx.ExecContext(ctx, `UPDATE package_node_traffic_suspensions SET credential_json=? WHERE username=? AND credential_json=?`, replacement.New, username, replacement.Old); err != nil {
			return fmt.Errorf("update suspended credential snapshot: %w", err)
		}
	}
	for _, node := range changes.Nodes {
		var result sql.Result
		var err error
		if node.RoutedAdmin {
			result, err = tx.ExecContext(ctx, `UPDATE nodes SET routed_admin_credential=?,clash_config=?,updated_at=CURRENT_TIMESTAMP WHERE id=? AND node_type='routed'`, node.CredentialJSON, node.ClashConfig, node.ID)
		} else {
			result, err = tx.ExecContext(ctx, `UPDATE nodes SET clash_config=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`, node.ClashConfig, node.ID)
		}
		if err != nil {
			return fmt.Errorf("update node credential %d: %w", node.ID, err)
		}
		if count, countErr := result.RowsAffected(); countErr != nil || count != 1 {
			return fmt.Errorf("update node credential %d: node no longer exists", node.ID)
		}
	}

	if _, err := tx.ExecContext(ctx, `INSERT INTO user_tokens(username,token,user_short_code,custom_user_short_code,updated_at)
		VALUES(?,?,?,'',CURRENT_TIMESTAMP)
		ON CONFLICT(username) DO UPDATE SET token=excluded.token,user_short_code=excluded.user_short_code,custom_user_short_code='',updated_at=CURRENT_TIMESTAMP`, username, token, shortCode); err != nil {
		return fmt.Errorf("rotate user subscription token: %w", err)
	}

	rows, err := tx.QueryContext(ctx, `SELECT id FROM user_package_assignments WHERE username=? AND status='active' ORDER BY id`, username)
	if err != nil {
		return fmt.Errorf("list package subscription links: %w", err)
	}
	var assignmentIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		assignmentIDs = append(assignmentIDs, id)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, id := range assignmentIDs {
		code, err := generateAssignmentShortCode()
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `UPDATE user_package_assignments SET short_code=?,updated_at=CURRENT_TIMESTAMP WHERE id=? AND username=?`, code, id, username); err != nil {
			return fmt.Errorf("rotate package subscription link %d: %w", id, err)
		}
	}
	return nil
}
