package storage

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	PackageAssignmentActive    = "active"
	PackageAssignmentSuspended = "suspended"
	PackageAssignmentExpired   = "expired"
)

// UserPackageAssignment is one assigned/bound package instance. A user may
// own the same package template more than once, so ID is its stable identity.
type UserPackageAssignment struct {
	ID                   int64      `json:"id"`
	Username             string     `json:"username"`
	PackageID            int64      `json:"package_id"`
	PackageName          string     `json:"package_name"`
	PackageDescription   string     `json:"package_description"`
	PackageStartDate     *time.Time `json:"package_start_date,omitempty"`
	PackageEndDate       *time.Time `json:"package_end_date,omitempty"`
	IsReset              bool       `json:"is_reset"`
	ResetDay             int        `json:"reset_day"`
	LastResetAt          *time.Time `json:"last_reset_at,omitempty"`
	TrafficLimitOverride *int64     `json:"traffic_limit_override,omitempty"`
	TrafficLimitBytes    int64      `json:"traffic_limit_bytes"`
	Status               string     `json:"status"`
	IsPrimary            bool       `json:"is_primary"`
	ShortCode            string     `json:"short_code"`
	TrafficWarned80      bool       `json:"traffic_warned_80"`
	OverLimitEnforced    bool       `json:"over_limit_enforced"`
	Legacy               bool       `json:"legacy"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

func (a UserPackageAssignment) EffectiveTrafficLimit() int64 {
	if a.TrafficLimitOverride != nil {
		return *a.TrafficLimitOverride
	}
	return a.TrafficLimitBytes
}

// New package instances use dedicated credential tables. This avoids changing
// the legacy tables' uniqueness constraints during a rolling upgrade.
type PackageAssignmentInboundConfig struct {
	ID             int64
	AssignmentID   int64
	Username       string
	ServerID       int64
	InboundTag     string
	Protocol       string
	Email          string
	CredentialJSON string
	CreatedAt      time.Time
}

type PackageAssignmentSubaccount struct {
	ID             int64
	AssignmentID   int64
	Username       string
	RoutedNodeID   int64
	Email          string
	CredentialJSON string
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (r *TrafficRepository) migrateUserPackageAssignments() error {
	const schema = `
CREATE TABLE IF NOT EXISTS user_package_assignments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    package_id INTEGER NOT NULL,
    package_start_date TIMESTAMP,
    package_end_date TIMESTAMP,
    is_reset INTEGER NOT NULL DEFAULT 0,
    reset_day INTEGER NOT NULL DEFAULT 1,
    last_reset_at TIMESTAMP,
    traffic_limit_override INTEGER,
    status TEXT NOT NULL DEFAULT 'active',
    is_primary INTEGER NOT NULL DEFAULT 0,
    short_code TEXT NOT NULL UNIQUE,
    traffic_warned_80 INTEGER NOT NULL DEFAULT 0,
    over_limit_enforced INTEGER NOT NULL DEFAULT 0,
    legacy_source INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE,
    FOREIGN KEY(package_id) REFERENCES packages(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_user_package_assignments_user_status ON user_package_assignments(username,status);
CREATE INDEX IF NOT EXISTS idx_user_package_assignments_package_status ON user_package_assignments(package_id,status);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_package_assignments_legacy ON user_package_assignments(username,legacy_source) WHERE legacy_source IS NOT NULL;
CREATE TABLE IF NOT EXISTS package_assignment_inbound_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    assignment_id INTEGER NOT NULL,
    username TEXT NOT NULL,
    server_id INTEGER NOT NULL,
    inbound_tag TEXT NOT NULL,
    protocol TEXT NOT NULL,
	 email TEXT NOT NULL,
    credential_json TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(assignment_id,server_id,inbound_tag),
    FOREIGN KEY(assignment_id) REFERENCES user_package_assignments(id) ON DELETE CASCADE,
    FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE,
    FOREIGN KEY(server_id) REFERENCES remote_servers(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_pkg_assignment_inbound_user ON package_assignment_inbound_configs(username);
CREATE INDEX IF NOT EXISTS idx_pkg_assignment_inbound_server ON package_assignment_inbound_configs(server_id);
CREATE TABLE IF NOT EXISTS package_assignment_subaccounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    assignment_id INTEGER NOT NULL,
    username TEXT NOT NULL,
    routed_node_id INTEGER NOT NULL,
    email TEXT NOT NULL,
    credential_json TEXT NOT NULL,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(assignment_id,routed_node_id),
    UNIQUE(routed_node_id,email),
    FOREIGN KEY(assignment_id) REFERENCES user_package_assignments(id) ON DELETE CASCADE,
    FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_pkg_assignment_subaccounts_user ON package_assignment_subaccounts(username);
CREATE INDEX IF NOT EXISTS idx_pkg_assignment_subaccounts_email ON package_assignment_subaccounts(email);
`
	if _, err := r.db.Exec(schema); err != nil {
		return fmt.Errorf("migrate user package assignments: %w", err)
	}
	if err := r.ensureTableColumn("package_assignment_inbound_configs", "email", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	return r.backfillLegacyPackageAssignments()
}

func (r *TrafficRepository) backfillLegacyPackageAssignments() error {
	rows, err := r.db.Query(`SELECT username,package_id,package_start_date,package_end_date,COALESCE(is_reset,0),COALESCE(reset_day,1),last_reset_at,traffic_limit_override,COALESCE(traffic_warned_80,0),COALESCE(over_limit_enforced,0)
		FROM users WHERE package_id IS NOT NULL AND package_id > 0`)
	if err != nil {
		return fmt.Errorf("list legacy package assignments: %w", err)
	}
	type legacy struct {
		username              string
		packageID             int64
		start, end, lastReset sql.NullTime
		isReset, resetDay     int
		limit                 sql.NullInt64
		warned, enforced      int
	}
	var values []legacy
	for rows.Next() {
		var v legacy
		if err := rows.Scan(&v.username, &v.packageID, &v.start, &v.end, &v.isReset, &v.resetDay, &v.lastReset, &v.limit, &v.warned, &v.enforced); err != nil {
			rows.Close()
			return err
		}
		values = append(values, v)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, v := range values {
		var exists int
		if err := r.db.QueryRow(`SELECT COUNT(1) FROM user_package_assignments WHERE username=? AND legacy_source=1`, v.username).Scan(&exists); err != nil {
			return err
		}
		if exists > 0 {
			continue
		}
		code, err := generateAssignmentShortCode()
		if err != nil {
			return err
		}
		var start, end, lastReset, limit any
		if v.start.Valid {
			start = v.start.Time
		}
		if v.end.Valid {
			end = v.end.Time
		}
		if v.lastReset.Valid {
			lastReset = v.lastReset.Time
		}
		if v.limit.Valid {
			limit = v.limit.Int64
		}
		if _, err := r.db.Exec(`INSERT INTO user_package_assignments(username,package_id,package_start_date,package_end_date,is_reset,reset_day,last_reset_at,traffic_limit_override,status,is_primary,short_code,traffic_warned_80,over_limit_enforced,legacy_source)
			VALUES(?,?,?,?,?,?,?,?, 'active',1,?,?,?,1)`, v.username, v.packageID, start, end, v.isReset, v.resetDay, lastReset, limit, code, v.warned, v.enforced); err != nil {
			return fmt.Errorf("backfill package assignment for %s: %w", v.username, err)
		}
	}
	return nil
}

func generateAssignmentShortCode() (string, error) {
	buf := make([]byte, 9)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return strings.TrimRight(base64.RawURLEncoding.EncodeToString(buf), "="), nil
}

const assignmentSelect = `a.id,a.username,a.package_id,p.name,p.description,a.package_start_date,a.package_end_date,a.is_reset,a.reset_day,a.last_reset_at,a.traffic_limit_override,p.traffic_limit_bytes,a.status,a.is_primary,a.short_code,a.traffic_warned_80,a.over_limit_enforced,COALESCE(a.legacy_source,0),a.created_at,a.updated_at`

func scanUserPackageAssignment(scanner interface{ Scan(...any) error }) (*UserPackageAssignment, error) {
	var a UserPackageAssignment
	var start, end, lastReset sql.NullTime
	var limit sql.NullInt64
	var isReset, primary, warned, enforced, legacy int
	err := scanner.Scan(&a.ID, &a.Username, &a.PackageID, &a.PackageName, &a.PackageDescription, &start, &end, &isReset, &a.ResetDay, &lastReset, &limit, &a.TrafficLimitBytes, &a.Status, &primary, &a.ShortCode, &warned, &enforced, &legacy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if start.Valid {
		v := start.Time
		a.PackageStartDate = &v
	}
	if end.Valid {
		v := end.Time
		a.PackageEndDate = &v
	}
	if lastReset.Valid {
		v := lastReset.Time
		a.LastResetAt = &v
	}
	if limit.Valid {
		v := limit.Int64
		a.TrafficLimitOverride = &v
	}
	a.IsReset, a.IsPrimary = isReset != 0, primary != 0
	a.TrafficWarned80, a.OverLimitEnforced = warned != 0, enforced != 0
	a.Legacy = legacy != 0
	return &a, nil
}

func (r *TrafficRepository) listPackageAssignments(ctx context.Context, query string, args ...any) ([]UserPackageAssignment, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []UserPackageAssignment
	for rows.Next() {
		a, err := scanUserPackageAssignment(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return out, rows.Err()
}

func (r *TrafficRepository) ListUserPackageAssignments(ctx context.Context, username string, includeInactive bool) ([]UserPackageAssignment, error) {
	q := `SELECT ` + assignmentSelect + ` FROM user_package_assignments a JOIN packages p ON p.id=a.package_id WHERE a.username=?`
	if !includeInactive {
		q += ` AND a.status='active'`
	}
	q += ` ORDER BY a.is_primary DESC,a.created_at,a.id`
	return r.listPackageAssignments(ctx, q, strings.TrimSpace(username))
}

func (r *TrafficRepository) ListActivePackageAssignments(ctx context.Context) ([]UserPackageAssignment, error) {
	return r.listPackageAssignments(ctx, `SELECT `+assignmentSelect+` FROM user_package_assignments a JOIN packages p ON p.id=a.package_id WHERE a.status='active' ORDER BY a.username,a.is_primary DESC,a.id`)
}

func (r *TrafficRepository) ListActivePackageAssignmentsByPackage(ctx context.Context, packageID int64) ([]UserPackageAssignment, error) {
	return r.listPackageAssignments(ctx, `SELECT `+assignmentSelect+` FROM user_package_assignments a JOIN packages p ON p.id=a.package_id WHERE a.package_id=? AND a.status='active' ORDER BY a.id`, packageID)
}

func (r *TrafficRepository) GetUserPackageAssignment(ctx context.Context, id int64) (*UserPackageAssignment, error) {
	a, err := scanUserPackageAssignment(r.db.QueryRowContext(ctx, `SELECT `+assignmentSelect+` FROM user_package_assignments a JOIN packages p ON p.id=a.package_id WHERE a.id=?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPackageNotFound
	}
	return a, err
}

func (r *TrafficRepository) GetUserPackageAssignmentByShortCode(ctx context.Context, code string) (*UserPackageAssignment, error) {
	a, err := scanUserPackageAssignment(r.db.QueryRowContext(ctx, `SELECT `+assignmentSelect+` FROM user_package_assignments a JOIN packages p ON p.id=a.package_id WHERE a.short_code=?`, strings.TrimSpace(code)))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPackageNotFound
	}
	return a, err
}

func (r *TrafficRepository) GetPrimaryUserPackageAssignment(ctx context.Context, username string) (*UserPackageAssignment, error) {
	a, err := scanUserPackageAssignment(r.db.QueryRowContext(ctx, `SELECT `+assignmentSelect+` FROM user_package_assignments a JOIN packages p ON p.id=a.package_id WHERE a.username=? AND a.status='active' ORDER BY a.is_primary DESC,a.created_at,a.id LIMIT 1`, strings.TrimSpace(username)))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return a, err
}

func (r *TrafficRepository) CreateUserPackageAssignment(ctx context.Context, username string, packageID int64, start, end time.Time, isReset bool, resetDay int) (*UserPackageAssignment, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("username is required")
	}
	if _, err := r.GetUser(ctx, username); err != nil {
		return nil, err
	}
	if _, err := r.GetPackage(ctx, packageID); err != nil {
		return nil, err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var active int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM user_package_assignments WHERE username=? AND status='active'`, username).Scan(&active); err != nil {
		return nil, err
	}
	code, err := generateAssignmentShortCode()
	if err != nil {
		return nil, err
	}
	isPrimary := active == 0
	res, err := tx.ExecContext(ctx, `INSERT INTO user_package_assignments(username,package_id,package_start_date,package_end_date,is_reset,reset_day,last_reset_at,status,is_primary,short_code) VALUES(?,?,?,?,?,?,CURRENT_TIMESTAMP,'active',?,?)`, username, packageID, start, end, isReset, resetDay, isPrimary, code)
	if err != nil {
		return nil, fmt.Errorf("create user package assignment: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	if isPrimary {
		if _, err := tx.ExecContext(ctx, `UPDATE users SET package_id=?,package_start_date=?,package_end_date=?,is_reset=?,reset_day=?,last_reset_at=CURRENT_TIMESTAMP,updated_at=CURRENT_TIMESTAMP WHERE username=?`, packageID, start, end, isReset, resetDay, username); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetUserPackageAssignment(ctx, id)
}

func (r *TrafficRepository) UpdateUserPackageAssignment(ctx context.Context, a UserPackageAssignment) error {
	res, err := r.db.ExecContext(ctx, `UPDATE user_package_assignments SET package_id=?,package_start_date=?,package_end_date=?,is_reset=?,reset_day=?,last_reset_at=?,traffic_limit_override=?,status=?,updated_at=CURRENT_TIMESTAMP WHERE id=? AND username=?`, a.PackageID, a.PackageStartDate, a.PackageEndDate, a.IsReset, a.ResetDay, a.LastResetAt, a.TrafficLimitOverride, a.Status, a.ID, a.Username)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrPackageNotFound
	}
	if a.IsPrimary {
		_, err = r.db.ExecContext(ctx, `UPDATE users SET package_id=?,package_start_date=?,package_end_date=?,is_reset=?,reset_day=?,last_reset_at=?,traffic_limit_override=?,updated_at=CURRENT_TIMESTAMP WHERE username=?`, a.PackageID, a.PackageStartDate, a.PackageEndDate, a.IsReset, a.ResetDay, a.LastResetAt, a.TrafficLimitOverride, a.Username)
	}
	return err
}

func (r *TrafficRepository) DeleteUserPackageAssignment(ctx context.Context, id int64) error {
	a, err := r.GetUserPackageAssignment(ctx, id)
	if err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM user_package_assignments WHERE id=?`, id); err != nil {
		return err
	}
	if a.IsPrimary {
		var nextID, packageID int64
		var start, end, lastReset sql.NullTime
		var isReset, resetDay int
		var limit sql.NullInt64
		err = tx.QueryRowContext(ctx, `SELECT id,package_id,package_start_date,package_end_date,is_reset,reset_day,last_reset_at,traffic_limit_override FROM user_package_assignments WHERE username=? AND status='active' ORDER BY created_at,id LIMIT 1`, a.Username).Scan(&nextID, &packageID, &start, &end, &isReset, &resetDay, &lastReset, &limit)
		if errors.Is(err, sql.ErrNoRows) {
			_, err = tx.ExecContext(ctx, `UPDATE users SET package_id=NULL,package_start_date=NULL,package_end_date=NULL,is_reset=0,reset_day=1,last_reset_at=NULL,traffic_limit_override=NULL,traffic_warned_80=0,is_over_limit=0,over_limit_enforced=0,updated_at=CURRENT_TIMESTAMP WHERE username=?`, a.Username)
		} else if err == nil {
			if _, err = tx.ExecContext(ctx, `UPDATE user_package_assignments SET is_primary=1,updated_at=CURRENT_TIMESTAMP WHERE id=?`, nextID); err == nil {
				var st, en, reset, lim any
				if start.Valid {
					st = start.Time
				}
				if end.Valid {
					en = end.Time
				}
				if lastReset.Valid {
					reset = lastReset.Time
				}
				if limit.Valid {
					lim = limit.Int64
				}
				_, err = tx.ExecContext(ctx, `UPDATE users SET package_id=?,package_start_date=?,package_end_date=?,is_reset=?,reset_day=?,last_reset_at=?,traffic_limit_override=?,traffic_warned_80=0,is_over_limit=0,over_limit_enforced=0,updated_at=CURRENT_TIMESTAMP WHERE username=?`, packageID, st, en, isReset, resetDay, reset, lim, a.Username)
			}
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *TrafficRepository) SavePackageAssignmentInboundConfig(ctx context.Context, c PackageAssignmentInboundConfig) error {
	if c.Email == "" {
		c.Email = CredentialEmail(c.CredentialJSON)
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO package_assignment_inbound_configs(assignment_id,username,server_id,inbound_tag,protocol,email,credential_json) VALUES(?,?,?,?,?,?,?) ON CONFLICT(assignment_id,server_id,inbound_tag) DO NOTHING`, c.AssignmentID, c.Username, c.ServerID, c.InboundTag, c.Protocol, c.Email, c.CredentialJSON)
	return err
}

func (r *TrafficRepository) GetPackageAssignmentInboundConfig(ctx context.Context, assignmentID, serverID int64, tag string) (*PackageAssignmentInboundConfig, error) {
	var c PackageAssignmentInboundConfig
	err := r.db.QueryRowContext(ctx, `SELECT id,assignment_id,username,server_id,inbound_tag,protocol,email,credential_json,created_at FROM package_assignment_inbound_configs WHERE assignment_id=? AND server_id=? AND inbound_tag=?`, assignmentID, serverID, tag).Scan(&c.ID, &c.AssignmentID, &c.Username, &c.ServerID, &c.InboundTag, &c.Protocol, &c.Email, &c.CredentialJSON, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &c, err
}

func (r *TrafficRepository) listPackageAssignmentInboundConfigs(ctx context.Context, query string, args ...any) ([]PackageAssignmentInboundConfig, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PackageAssignmentInboundConfig
	for rows.Next() {
		var c PackageAssignmentInboundConfig
		if err := rows.Scan(&c.ID, &c.AssignmentID, &c.Username, &c.ServerID, &c.InboundTag, &c.Protocol, &c.Email, &c.CredentialJSON, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *TrafficRepository) ListPackageAssignmentInboundConfigs(ctx context.Context, assignmentID int64) ([]PackageAssignmentInboundConfig, error) {
	return r.listPackageAssignmentInboundConfigs(ctx, `SELECT id,assignment_id,username,server_id,inbound_tag,protocol,email,credential_json,created_at FROM package_assignment_inbound_configs WHERE assignment_id=?`, assignmentID)
}
func (r *TrafficRepository) ListPackageAssignmentInboundConfigsByServer(ctx context.Context, serverID int64) ([]PackageAssignmentInboundConfig, error) {
	return r.listPackageAssignmentInboundConfigs(ctx, `SELECT id,assignment_id,username,server_id,inbound_tag,protocol,email,credential_json,created_at FROM package_assignment_inbound_configs WHERE server_id=?`, serverID)
}
func (r *TrafficRepository) ListPackageAssignmentInboundConfigsByUser(ctx context.Context, username string) ([]PackageAssignmentInboundConfig, error) {
	return r.listPackageAssignmentInboundConfigs(ctx, `SELECT id,assignment_id,username,server_id,inbound_tag,protocol,email,credential_json,created_at FROM package_assignment_inbound_configs WHERE username=?`, username)
}
func (r *TrafficRepository) ListAllPackageAssignmentInboundConfigs(ctx context.Context) ([]PackageAssignmentInboundConfig, error) {
	return r.listPackageAssignmentInboundConfigs(ctx, `SELECT id,assignment_id,username,server_id,inbound_tag,protocol,email,credential_json,created_at FROM package_assignment_inbound_configs`)
}
func (r *TrafficRepository) DeletePackageAssignmentInboundConfig(ctx context.Context, assignmentID, serverID int64, tag string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM package_assignment_inbound_configs WHERE assignment_id=? AND server_id=? AND inbound_tag=?`, assignmentID, serverID, tag)
	return err
}

func (r *TrafficRepository) DeletePackageAssignmentInboundConfigsByInbound(ctx context.Context, serverID int64, tag string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `DELETE FROM package_assignment_inbound_configs WHERE server_id=? AND inbound_tag=?`, serverID, tag)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *TrafficRepository) GetPackageAssignmentSubaccount(ctx context.Context, assignmentID, nodeID int64) (*PackageAssignmentSubaccount, error) {
	var s PackageAssignmentSubaccount
	var active int
	err := r.db.QueryRowContext(ctx, `SELECT id,assignment_id,username,routed_node_id,email,credential_json,is_active,created_at,updated_at FROM package_assignment_subaccounts WHERE assignment_id=? AND routed_node_id=?`, assignmentID, nodeID).Scan(&s.ID, &s.AssignmentID, &s.Username, &s.RoutedNodeID, &s.Email, &s.CredentialJSON, &active, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	s.IsActive = active != 0
	return &s, err
}
func (r *TrafficRepository) ClaimPackageAssignmentSubaccount(ctx context.Context, s PackageAssignmentSubaccount) (*PackageAssignmentSubaccount, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO package_assignment_subaccounts(assignment_id,username,routed_node_id,email,credential_json,is_active) VALUES(?,?,?,?,?,?) ON CONFLICT(assignment_id,routed_node_id) DO NOTHING`, s.AssignmentID, s.Username, s.RoutedNodeID, s.Email, s.CredentialJSON, s.IsActive)
	if err != nil {
		return nil, err
	}
	return r.GetPackageAssignmentSubaccount(ctx, s.AssignmentID, s.RoutedNodeID)
}
func (r *TrafficRepository) UpsertPackageAssignmentSubaccount(ctx context.Context, s PackageAssignmentSubaccount) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO package_assignment_subaccounts(assignment_id,username,routed_node_id,email,credential_json,is_active) VALUES(?,?,?,?,?,?) ON CONFLICT(assignment_id,routed_node_id) DO UPDATE SET email=excluded.email,credential_json=excluded.credential_json,is_active=excluded.is_active,updated_at=CURRENT_TIMESTAMP`, s.AssignmentID, s.Username, s.RoutedNodeID, s.Email, s.CredentialJSON, s.IsActive)
	return err
}
func (r *TrafficRepository) SetPackageAssignmentSubaccountActive(ctx context.Context, id int64, active bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE package_assignment_subaccounts SET is_active=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`, active, id)
	return err
}

func (r *TrafficRepository) DeletePackageAssignmentSubaccountByIdentity(ctx context.Context, assignmentID, routedNodeID int64, email string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM package_assignment_subaccounts WHERE assignment_id=? AND routed_node_id=? AND email=?`, assignmentID, routedNodeID, email)
	return err
}

func (r *TrafficRepository) DeletePackageAssignmentSubaccountsByRoutedNode(ctx context.Context, routedNodeID int64) (int64, error) {
	result, err := r.db.ExecContext(ctx, `DELETE FROM package_assignment_subaccounts WHERE routed_node_id=?`, routedNodeID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
func (r *TrafficRepository) ListPackageAssignmentSubaccounts(ctx context.Context, assignmentID int64) ([]PackageAssignmentSubaccount, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,assignment_id,username,routed_node_id,email,credential_json,is_active,created_at,updated_at FROM package_assignment_subaccounts WHERE assignment_id=?`, assignmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PackageAssignmentSubaccount
	for rows.Next() {
		var s PackageAssignmentSubaccount
		var active int
		if err := rows.Scan(&s.ID, &s.AssignmentID, &s.Username, &s.RoutedNodeID, &s.Email, &s.CredentialJSON, &active, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		s.IsActive = active != 0
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *TrafficRepository) ListPackageAssignmentSubaccountsByUser(ctx context.Context, username string) ([]PackageAssignmentSubaccount, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,assignment_id,username,routed_node_id,email,credential_json,is_active,created_at,updated_at FROM package_assignment_subaccounts WHERE username=?`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PackageAssignmentSubaccount
	for rows.Next() {
		var s PackageAssignmentSubaccount
		var active int
		if err := rows.Scan(&s.ID, &s.AssignmentID, &s.Username, &s.RoutedNodeID, &s.Email, &s.CredentialJSON, &active, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		s.IsActive = active != 0
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *TrafficRepository) ListPackageAssignmentSubaccountsByRoutedNode(ctx context.Context, routedNodeID int64) ([]PackageAssignmentSubaccount, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,assignment_id,username,routed_node_id,email,credential_json,is_active,created_at,updated_at FROM package_assignment_subaccounts WHERE routed_node_id=?`, routedNodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PackageAssignmentSubaccount
	for rows.Next() {
		var s PackageAssignmentSubaccount
		var active int
		if err := rows.Scan(&s.ID, &s.AssignmentID, &s.Username, &s.RoutedNodeID, &s.Email, &s.CredentialJSON, &active, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		s.IsActive = active != 0
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *TrafficRepository) ListServerIDsForPackageAssignmentSubaccounts(ctx context.Context, username string) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT DISTINCT rs.id FROM package_assignment_subaccounts sa INNER JOIN nodes n ON n.id=sa.routed_node_id INNER JOIN remote_servers rs ON rs.name=n.original_server WHERE sa.username=? AND n.node_type='routed'`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}
func (r *TrafficRepository) ListAllPackageAssignmentSubaccounts(ctx context.Context) ([]PackageAssignmentSubaccount, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,assignment_id,username,routed_node_id,email,credential_json,is_active,created_at,updated_at FROM package_assignment_subaccounts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PackageAssignmentSubaccount
	for rows.Next() {
		var s PackageAssignmentSubaccount
		var active int
		if err := rows.Scan(&s.ID, &s.AssignmentID, &s.Username, &s.RoutedNodeID, &s.Email, &s.CredentialJSON, &active, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		s.IsActive = active != 0
		out = append(out, s)
	}
	return out, rows.Err()
}

func PackageAssignmentCredentialEmail(assignmentID int64, username, inboundTag string) string {
	return fmt.Sprintf("%s__pkg%d__%s", username, assignmentID, inboundTag)
}

func CredentialEmail(raw string) string {
	var v map[string]any
	if json.Unmarshal([]byte(raw), &v) != nil {
		return ""
	}
	email := strings.TrimSpace(fmt.Sprint(v["email"]))
	if email == "<nil>" {
		return ""
	}
	return email
}

func (r *TrafficRepository) GetPackageAssignmentBillableTrafficByDirection(ctx context.Context, assignmentID int64) (int64, int64, error) {
	a, err := r.GetUserPackageAssignment(ctx, assignmentID)
	if err != nil {
		return 0, 0, err
	}
	var up, down float64
	query := `SELECT
		COALESCE(SUM(CASE WHEN weighted_uplink>cycle_base_weighted_uplink THEN weighted_uplink-cycle_base_weighted_uplink ELSE 0 END),0),
		COALESCE(SUM(CASE WHEN weighted_downlink>cycle_base_weighted_downlink THEN weighted_downlink-cycle_base_weighted_downlink ELSE 0 END),0)
		FROM user_email_traffic WHERE email IN (
			SELECT email FROM package_assignment_inbound_configs WHERE assignment_id=?
			UNION SELECT email FROM package_assignment_subaccounts WHERE assignment_id=?
		)`
	args := []any{assignmentID, assignmentID}
	if a.Legacy {
		// The legacy assignment owns every historical row for the user except
		// emails explicitly claimed by an independent package instance. This
		// preserves pre-migration traffic while keeping new package quotas isolated.
		query = `SELECT
			COALESCE(SUM(CASE WHEN weighted_uplink>cycle_base_weighted_uplink THEN weighted_uplink-cycle_base_weighted_uplink ELSE 0 END),0),
			COALESCE(SUM(CASE WHEN weighted_downlink>cycle_base_weighted_downlink THEN weighted_downlink-cycle_base_weighted_downlink ELSE 0 END),0)
			FROM user_email_traffic WHERE attributed_username=? AND email NOT IN (
				SELECT c.email FROM package_assignment_inbound_configs c JOIN user_package_assignments a ON a.id=c.assignment_id WHERE a.username=? AND COALESCE(a.legacy_source,0)=0
				UNION SELECT s.email FROM package_assignment_subaccounts s JOIN user_package_assignments a ON a.id=s.assignment_id WHERE a.username=? AND COALESCE(a.legacy_source,0)=0
			)`
		args = []any{a.Username, a.Username, a.Username}
	}
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&up, &down)
	if err != nil {
		return 0, 0, fmt.Errorf("query package assignment traffic: %w", err)
	}
	return int64(up), int64(down), nil
}

func (r *TrafficRepository) GetPackageAssignmentBillableTraffic(ctx context.Context, assignmentID int64) (int64, error) {
	up, down, err := r.GetPackageAssignmentBillableTrafficByDirection(ctx, assignmentID)
	return up + down, err
}

// ResetPackageAssignmentTrafficCycleAt advances only this package instance's
// credential baselines. Other packages owned by the same user are untouched.
func (r *TrafficRepository) ResetPackageAssignmentTrafficCycleAt(ctx context.Context, assignmentID int64, resetAt time.Time) error {
	a, err := r.GetUserPackageAssignment(ctx, assignmentID)
	if err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `UPDATE user_email_traffic
		SET cycle_base_uplink=uplink,cycle_base_downlink=downlink,
		    cycle_base_weighted_uplink=weighted_uplink,cycle_base_weighted_downlink=weighted_downlink,
		    cycle_start=CURRENT_TIMESTAMP,updated_at=CURRENT_TIMESTAMP
		WHERE email IN (
			SELECT email FROM package_assignment_inbound_configs WHERE assignment_id=?
			UNION SELECT email FROM package_assignment_subaccounts WHERE assignment_id=?
		)`
	args := []any{assignmentID, assignmentID}
	if a.Legacy {
		query = `UPDATE user_email_traffic
			SET cycle_base_uplink=uplink,cycle_base_downlink=downlink,
			    cycle_base_weighted_uplink=weighted_uplink,cycle_base_weighted_downlink=weighted_downlink,
			    cycle_start=CURRENT_TIMESTAMP,updated_at=CURRENT_TIMESTAMP
			WHERE attributed_username=? AND email NOT IN (
				SELECT c.email FROM package_assignment_inbound_configs c JOIN user_package_assignments a ON a.id=c.assignment_id WHERE a.username=? AND COALESCE(a.legacy_source,0)=0
				UNION SELECT s.email FROM package_assignment_subaccounts s JOIN user_package_assignments a ON a.id=s.assignment_id WHERE a.username=? AND COALESCE(a.legacy_source,0)=0
			)`
		args = []any{a.Username, a.Username, a.Username}
	}
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("reset package assignment traffic: %w", err)
	}
	if a.Legacy {
		if _, err := tx.ExecContext(ctx, `UPDATE users SET last_reset_at=?,traffic_warned_80=0,updated_at=CURRENT_TIMESTAMP WHERE username=?`, resetAt, a.Username); err != nil {
			return fmt.Errorf("mark legacy package assignment reset: %w", err)
		}
	}
	res, err := tx.ExecContext(ctx, `UPDATE user_package_assignments SET last_reset_at=?,traffic_warned_80=0,updated_at=CURRENT_TIMESTAMP WHERE id=?`, resetAt, assignmentID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrPackageNotFound
	}
	return tx.Commit()
}

func (r *TrafficRepository) UpdatePackageAssignmentLimitState(ctx context.Context, assignmentID int64, warned80, enforced *bool) error {
	sets := make([]string, 0, 2)
	args := make([]any, 0, 3)
	if warned80 != nil {
		sets = append(sets, "traffic_warned_80=?")
		args = append(args, *warned80)
	}
	if enforced != nil {
		sets = append(sets, "over_limit_enforced=?")
		args = append(args, *enforced)
	}
	if len(sets) == 0 {
		return nil
	}
	sets = append(sets, "updated_at=CURRENT_TIMESTAMP")
	args = append(args, assignmentID)
	_, err := r.db.ExecContext(ctx, `UPDATE user_package_assignments SET `+strings.Join(sets, ",")+` WHERE id=?`, args...)
	return err
}

func (r *TrafficRepository) UpdatePackageAssignmentStatus(ctx context.Context, assignmentID int64, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE user_package_assignments SET status=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`, status, assignmentID)
	return err
}

type PackageAssignmentDailyTraffic struct {
	Date     string `json:"date"`
	Uplink   int64  `json:"uplink"`
	Downlink int64  `json:"downlink"`
}

func (r *TrafficRepository) ListPackageAssignmentDailyTraffic(ctx context.Context, assignmentID int64, startDate, endDate string) ([]PackageAssignmentDailyTraffic, error) {
	a, err := r.GetUserPackageAssignment(ctx, assignmentID)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if a.Legacy {
		rows, err = r.db.QueryContext(ctx, `SELECT date,COALESCE(SUM(weighted_uplink),0),COALESCE(SUM(weighted_downlink),0)
			FROM traffic_daily_user_emails WHERE attributed_username=? AND email NOT IN (
				SELECT c.email FROM package_assignment_inbound_configs c JOIN user_package_assignments a ON a.id=c.assignment_id WHERE a.username=? AND COALESCE(a.legacy_source,0)=0
				UNION SELECT s.email FROM package_assignment_subaccounts s JOIN user_package_assignments a ON a.id=s.assignment_id WHERE a.username=? AND COALESCE(a.legacy_source,0)=0
			) AND date>=? AND date<=? GROUP BY date ORDER BY date`, a.Username, a.Username, a.Username, startDate, endDate)
	} else {
		rows, err = r.db.QueryContext(ctx, `SELECT date,COALESCE(SUM(weighted_uplink),0),COALESCE(SUM(weighted_downlink),0) FROM traffic_daily_user_emails WHERE email IN (
			SELECT email FROM package_assignment_inbound_configs WHERE assignment_id=?
			UNION SELECT email FROM package_assignment_subaccounts WHERE assignment_id=?
		) AND date>=? AND date<=? GROUP BY date ORDER BY date`, assignmentID, assignmentID, startDate, endDate)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PackageAssignmentDailyTraffic
	for rows.Next() {
		var v PackageAssignmentDailyTraffic
		var up, down float64
		if err := rows.Scan(&v.Date, &up, &down); err != nil {
			return nil, err
		}
		v.Uplink, v.Downlink = int64(up), int64(down)
		out = append(out, v)
	}
	return out, rows.Err()
}
