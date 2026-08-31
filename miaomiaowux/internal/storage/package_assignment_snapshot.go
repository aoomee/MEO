package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// UserPackageAssignmentSnapshot captures every package-assignment column that
// can be changed before remote Xray provisioning finishes. It is used to
// restore the exact previous row when any Agent rejects the coordinated
// configuration transaction.
type UserPackageAssignmentSnapshot struct {
	Username             string
	PackageID            *int64
	PackageStartDate     *time.Time
	PackageEndDate       *time.Time
	IsReset              bool
	ResetDay             int
	TrafficLimitOverride *int64
	LastResetAt          *time.Time
	TrafficWarned80      bool
	EmailTrafficBases    []userEmailTrafficBaseSnapshot
	UserTrafficRows      []userTrafficCycleSnapshot
	TrafficCarry         *trafficCarrySnapshot
	NodeTrafficBaselines []packageNodeTrafficBaselineSnapshot
	NodeSuspensions      []packageNodeTrafficSuspensionSnapshot
}

type userEmailTrafficBaseSnapshot struct {
	ID                                                 int64
	CycleBaseUplink, CycleBaseDownlink                 int64
	CycleBaseWeightedUplink, CycleBaseWeightedDownlink float64
	CycleStart                                         time.Time
}

type userTrafficCycleSnapshot struct {
	ID               int64
	Uplink, Downlink int64
	CycleStart       time.Time
}

type trafficCarrySnapshot struct {
	WeightedUplink, WeightedDownlink float64
}

type packageNodeTrafficBaselineSnapshot struct {
	PackageID, NodeID int64
	Baseline          float64
}

type packageNodeTrafficSuspensionSnapshot struct {
	PackageID, NodeID int64
	Kind, Credential  string
}

func (r *TrafficRepository) GetUserPackageAssignmentSnapshot(ctx context.Context, username string) (*UserPackageAssignmentSnapshot, error) {
	var packageID sql.NullInt64
	var start, end sql.NullTime
	var isReset int
	var resetDay int
	var override sql.NullInt64
	var lastReset sql.NullTime
	var warned int
	err := r.db.QueryRowContext(ctx, `
		SELECT package_id, package_start_date, package_end_date,
		       COALESCE(is_reset,0), COALESCE(reset_day,1), traffic_limit_override,
		       last_reset_at, COALESCE(traffic_warned_80,0)
		FROM users WHERE username=?`, username).
		Scan(&packageID, &start, &end, &isReset, &resetDay, &override, &lastReset, &warned)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("snapshot user package assignment: %w", err)
	}
	out := &UserPackageAssignmentSnapshot{Username: username, IsReset: isReset != 0, ResetDay: resetDay, TrafficWarned80: warned != 0}
	if packageID.Valid {
		v := packageID.Int64
		out.PackageID = &v
	}
	if start.Valid {
		v := start.Time
		out.PackageStartDate = &v
	}
	if end.Valid {
		v := end.Time
		out.PackageEndDate = &v
	}
	if override.Valid {
		v := override.Int64
		out.TrafficLimitOverride = &v
	}
	if lastReset.Valid {
		v := lastReset.Time
		out.LastResetAt = &v
	}
	rows, err := r.db.QueryContext(ctx, `SELECT id,cycle_base_uplink,cycle_base_downlink,cycle_base_weighted_uplink,cycle_base_weighted_downlink,cycle_start FROM user_email_traffic WHERE attributed_username=?`, username)
	if err != nil {
		return nil, fmt.Errorf("snapshot email traffic bases: %w", err)
	}
	for rows.Next() {
		var v userEmailTrafficBaseSnapshot
		if err := rows.Scan(&v.ID, &v.CycleBaseUplink, &v.CycleBaseDownlink, &v.CycleBaseWeightedUplink, &v.CycleBaseWeightedDownlink, &v.CycleStart); err != nil {
			rows.Close()
			return nil, err
		}
		out.EmailTrafficBases = append(out.EmailTrafficBases, v)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	rows, err = r.db.QueryContext(ctx, `SELECT id,uplink,downlink,cycle_start FROM user_traffic WHERE username=?`, username)
	if err != nil {
		return nil, fmt.Errorf("snapshot user traffic cycle: %w", err)
	}
	for rows.Next() {
		var v userTrafficCycleSnapshot
		if err := rows.Scan(&v.ID, &v.Uplink, &v.Downlink, &v.CycleStart); err != nil {
			rows.Close()
			return nil, err
		}
		out.UserTrafficRows = append(out.UserTrafficRows, v)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	var carry trafficCarrySnapshot
	if err := r.db.QueryRowContext(ctx, `SELECT weighted_uplink,weighted_downlink FROM user_traffic_cycle_carry WHERE username=?`, username).Scan(&carry.WeightedUplink, &carry.WeightedDownlink); err == nil {
		out.TrafficCarry = &carry
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("snapshot traffic carry: %w", err)
	}
	rows, err = r.db.QueryContext(ctx, `SELECT package_id,node_id,baseline FROM package_user_node_traffic_baselines WHERE username=?`, username)
	if err != nil {
		return nil, fmt.Errorf("snapshot package node baselines: %w", err)
	}
	for rows.Next() {
		var v packageNodeTrafficBaselineSnapshot
		if err := rows.Scan(&v.PackageID, &v.NodeID, &v.Baseline); err != nil {
			rows.Close()
			return nil, err
		}
		out.NodeTrafficBaselines = append(out.NodeTrafficBaselines, v)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	rows, err = r.db.QueryContext(ctx, `SELECT package_id,node_id,kind,credential_json FROM package_node_traffic_suspensions WHERE username=?`, username)
	if err != nil {
		return nil, fmt.Errorf("snapshot package node suspensions: %w", err)
	}
	for rows.Next() {
		var v packageNodeTrafficSuspensionSnapshot
		if err := rows.Scan(&v.PackageID, &v.NodeID, &v.Kind, &v.Credential); err != nil {
			rows.Close()
			return nil, err
		}
		out.NodeSuspensions = append(out.NodeSuspensions, v)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *TrafficRepository) RestoreUserPackageAssignment(ctx context.Context, snapshot *UserPackageAssignmentSnapshot) error {
	if snapshot == nil {
		return fmt.Errorf("package assignment snapshot is nil")
	}
	var packageID interface{}
	if snapshot.PackageID != nil {
		packageID = *snapshot.PackageID
	}
	var start, end interface{}
	if snapshot.PackageStartDate != nil {
		start = *snapshot.PackageStartDate
	}
	if snapshot.PackageEndDate != nil {
		end = *snapshot.PackageEndDate
	}
	var override interface{}
	if snapshot.TrafficLimitOverride != nil {
		override = *snapshot.TrafficLimitOverride
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin restore user package assignment: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var lastReset interface{}
	if snapshot.LastResetAt != nil {
		lastReset = *snapshot.LastResetAt
	}
	result, err := tx.ExecContext(ctx, `
		UPDATE users SET package_id=?, package_start_date=?, package_end_date=?,
		       is_reset=?, reset_day=?, traffic_limit_override=?, last_reset_at=?, traffic_warned_80=?, updated_at=CURRENT_TIMESTAMP
		WHERE username=?`, packageID, start, end, snapshot.IsReset, snapshot.ResetDay, override, lastReset, snapshot.TrafficWarned80, snapshot.Username)
	if err != nil {
		return fmt.Errorf("restore user package assignment: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return ErrUserNotFound
	}
	for _, v := range snapshot.EmailTrafficBases {
		if _, err := tx.ExecContext(ctx, `UPDATE user_email_traffic SET cycle_base_uplink=?,cycle_base_downlink=?,cycle_base_weighted_uplink=?,cycle_base_weighted_downlink=?,cycle_start=? WHERE id=?`, v.CycleBaseUplink, v.CycleBaseDownlink, v.CycleBaseWeightedUplink, v.CycleBaseWeightedDownlink, v.CycleStart, v.ID); err != nil {
			return err
		}
	}
	for _, v := range snapshot.UserTrafficRows {
		if _, err := tx.ExecContext(ctx, `UPDATE user_traffic SET uplink=?,downlink=?,cycle_start=? WHERE id=?`, v.Uplink, v.Downlink, v.CycleStart, v.ID); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM user_traffic_cycle_carry WHERE username=?`, snapshot.Username); err != nil {
		return err
	}
	if snapshot.TrafficCarry != nil {
		if _, err := tx.ExecContext(ctx, `INSERT INTO user_traffic_cycle_carry(username,weighted_uplink,weighted_downlink) VALUES(?,?,?)`, snapshot.Username, snapshot.TrafficCarry.WeightedUplink, snapshot.TrafficCarry.WeightedDownlink); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM package_user_node_traffic_baselines WHERE username=?`, snapshot.Username); err != nil {
		return err
	}
	for _, v := range snapshot.NodeTrafficBaselines {
		if _, err := tx.ExecContext(ctx, `INSERT INTO package_user_node_traffic_baselines(username,package_id,node_id,baseline) VALUES(?,?,?,?)`, snapshot.Username, v.PackageID, v.NodeID, v.Baseline); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM package_node_traffic_suspensions WHERE username=?`, snapshot.Username); err != nil {
		return err
	}
	for _, v := range snapshot.NodeSuspensions {
		if _, err := tx.ExecContext(ctx, `INSERT INTO package_node_traffic_suspensions(username,package_id,node_id,kind,credential_json) VALUES(?,?,?,?,?)`, snapshot.Username, v.PackageID, v.NodeID, v.Kind, v.Credential); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit restore user package assignment: %w", err)
	}
	return nil
}
