package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// 转发管理数据模型:转发组 / 组成员 / 转发链 / 链跳。
// 配合 mmw-agent 的原生四层转发引擎(internal/forward),不走 xray tunnel。
// 建表走 migrate() 幂等调用 + 各 CRUD 开头 lazy 兜底(参考 shared_servers.go)。

// ForwardGroup 一组服务器 + 均衡策略 + 故障策略 +(可选)DNS 域名。
type ForwardGroup struct {
	ID                 int64                `json:"id"`
	Name               string               `json:"name"`
	BalanceStrategy    string               `json:"balance_strategy"`     // round_robin | percentage | cycle
	DNSDomain          string               `json:"dns_domain"`           // 空=不做 DNS LB(仅作链中间/出口)
	DNSProviderID      int64                `json:"dns_provider_id"`      // 关联 dns_providers,0=无
	FailoverEnabled    bool                 `json:"failover_enabled"`     //
	OfflineMsThreshold int                  `json:"offline_ms_threshold"` // RTT 超阈视为不健康卸载,0=只看掉线
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
	Members            []ForwardGroupMember `json:"members,omitempty"`
}

// ForwardGroupMember 组内一台服务器 + 手工权重。
type ForwardGroupMember struct {
	GroupID  int64 `json:"group_id"`
	ServerID int64 `json:"server_id"`
	Weight   int   `json:"weight"`
	Seq      int   `json:"seq"`
}

// ForwardChain 有序转发组列表(第一组=入口,最后一组=出口)。
type ForwardChain struct {
	ID             int64                 `json:"id"`
	Name           string                `json:"name"`
	PortRangeStart int                   `json:"port_range_start"`
	PortRangeEnd   int                   `json:"port_range_end"`
	DNSDomain      string                `json:"dns_domain"`
	DNSDomainV6    string                `json:"dns_domain_v6"`
	DNSProviderID  int64                 `json:"dns_provider_id"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	Hops           []ForwardChainHop     `json:"hops,omitempty"`
	Branches       []ForwardChainBranch  `json:"branches,omitempty"`
	DailyTraffic   []ForwardDailyTraffic `json:"daily_traffic,omitempty"`
	BoundNodes     []ForwardChainBinding `json:"bound_nodes,omitempty"`
}

// ForwardChainHop 链的一跳,seq 0-based。前端 v0.5.2 用 order / group_name 画拓扑。
type ForwardChainHop struct {
	ChainID   int64  `json:"chain_id"`
	Seq       int    `json:"seq"`
	Order     int    `json:"order"`
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name,omitempty"`
}

const forwardTablesSchema = `
CREATE TABLE IF NOT EXISTS forward_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    balance_strategy TEXT NOT NULL DEFAULT 'round_robin',
    dns_domain TEXT NOT NULL DEFAULT '',
    dns_provider_id INTEGER NOT NULL DEFAULT 0,
    failover_enabled INTEGER NOT NULL DEFAULT 1,
    offline_ms_threshold INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS forward_group_members (
    group_id INTEGER NOT NULL,
    server_id INTEGER NOT NULL,
    weight INTEGER NOT NULL DEFAULT 1,
    seq INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY(group_id, server_id)
);
CREATE TABLE IF NOT EXISTS forward_chains (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    port_range_start INTEGER NOT NULL DEFAULT 0,
    port_range_end INTEGER NOT NULL DEFAULT 0,
    dns_domain TEXT NOT NULL DEFAULT '',
    dns_domain_v6 TEXT NOT NULL DEFAULT '',
    dns_provider_id INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS forward_chain_hops (
    chain_id INTEGER NOT NULL,
    seq INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    PRIMARY KEY(chain_id, seq)
);
CREATE TABLE IF NOT EXISTS forward_chain_nodes (
    node_id INTEGER PRIMARY KEY,
    chain_id INTEGER NOT NULL,
    port INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    terminus_addr TEXT NOT NULL DEFAULT '',
    protocol TEXT NOT NULL DEFAULT 'tcp',
    pinned_exit_server_id INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS forward_hop_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    server_id INTEGER NOT NULL,
    rule_id TEXT NOT NULL,
    upstream_addr TEXT NOT NULL,
    healthy INTEGER NOT NULL DEFAULT 0,
    rtt_ms INTEGER NOT NULL DEFAULT 0,
    bytes_up INTEGER NOT NULL DEFAULT 0,
    bytes_down INTEGER NOT NULL DEFAULT 0,
    at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_forward_group_members_group ON forward_group_members(group_id);
CREATE INDEX IF NOT EXISTS idx_forward_chain_hops_chain ON forward_chain_hops(chain_id);
CREATE INDEX IF NOT EXISTS idx_forward_chain_nodes_chain ON forward_chain_nodes(chain_id);
CREATE INDEX IF NOT EXISTS idx_forward_hop_metrics_q ON forward_hop_metrics(server_id, rule_id, at);
CREATE TABLE IF NOT EXISTS forward_chain_branches (
    chain_id INTEGER NOT NULL,
    hop_seq INTEGER NOT NULL,
    server_id INTEGER NOT NULL,
    seq INTEGER NOT NULL,
    via_group_id INTEGER NOT NULL,
    PRIMARY KEY(chain_id, hop_seq, seq)
);
CREATE TABLE IF NOT EXISTS forward_daily_traffic (
    chain_id INTEGER NOT NULL,
    server_id INTEGER NOT NULL,
    date TEXT NOT NULL,
    bytes_up INTEGER NOT NULL DEFAULT 0,
    bytes_down INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(chain_id, server_id, date)
);
CREATE INDEX IF NOT EXISTS idx_fwd_branches_chain ON forward_chain_branches(chain_id);
CREATE INDEX IF NOT EXISTS idx_fwd_daily_chain ON forward_daily_traffic(chain_id, date);
`

func (r *TrafficRepository) ensureForwardTables(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, forwardTablesSchema); err != nil {
		return fmt.Errorf("migrate forward tables: %w", err)
	}
	for _, col := range []struct{ name, def string }{
		{"port_range_start", "INTEGER NOT NULL DEFAULT 0"},
		{"port_range_end", "INTEGER NOT NULL DEFAULT 0"},
		{"dns_domain", "TEXT NOT NULL DEFAULT ''"},
		{"dns_domain_v6", "TEXT NOT NULL DEFAULT ''"},
		{"dns_provider_id", "INTEGER NOT NULL DEFAULT 0"},
	} {
		if err := r.ensureForwardTableColumn("forward_chains", col.name, col.def); err != nil {
			return err
		}
	}
	for _, col := range []struct{ name, def string }{
		{"terminus_addr", "TEXT NOT NULL DEFAULT ''"},
		{"protocol", "TEXT NOT NULL DEFAULT 'tcp'"},
		{"pinned_exit_server_id", "INTEGER NOT NULL DEFAULT 0"},
	} {
		if err := r.ensureForwardTableColumn("forward_chain_nodes", col.name, col.def); err != nil {
			return err
		}
	}
	if err := r.ensureForwardTableColumn("forward_group_members", "seq", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	return nil
}

func (r *TrafficRepository) ensureForwardChainColumn(name, definition string) error {
	return r.ensureForwardTableColumn("forward_chains", name, definition)
}

func (r *TrafficRepository) ensureForwardTableColumn(table, name, definition string) error {
	rows, err := r.db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return fmt.Errorf("%s table info: %w", table, err)
	}
	defer rows.Close()
	for rows.Next() {
		var cid, notNull, pk int
		var colName, colType string
		var defaultVal sql.NullString
		if err := rows.Scan(&cid, &colName, &colType, &notNull, &defaultVal, &pk); err != nil {
			return fmt.Errorf("scan %s info: %w", table, err)
		}
		if strings.EqualFold(colName, name) {
			return nil
		}
	}
	alter := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, name, definition)
	if _, err := r.db.Exec(alter); err != nil {
		return fmt.Errorf("add %s.%s: %w", table, name, err)
	}
	return nil
}

// ---------- 转发组 ----------

func (r *TrafficRepository) CreateForwardGroup(ctx context.Context, g *ForwardGroup) (int64, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return 0, err
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO forward_groups (name, balance_strategy, dns_domain, dns_provider_id, failover_enabled, offline_ms_threshold, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		g.Name, g.BalanceStrategy, g.DNSDomain, g.DNSProviderID, boolToInt(g.FailoverEnabled), g.OfflineMsThreshold)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *TrafficRepository) UpdateForwardGroup(ctx context.Context, g *ForwardGroup) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE forward_groups SET name=?, balance_strategy=?, dns_domain=?, dns_provider_id=?, failover_enabled=?, offline_ms_threshold=?, updated_at=CURRENT_TIMESTAMP
		 WHERE id=?`,
		g.Name, g.BalanceStrategy, g.DNSDomain, g.DNSProviderID, boolToInt(g.FailoverEnabled), g.OfflineMsThreshold, g.ID)
	return err
}

func (r *TrafficRepository) DeleteForwardGroup(ctx context.Context, id int64) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM forward_group_members WHERE group_id=?`, id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM forward_chain_branches WHERE via_group_id=?`, id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM forward_groups WHERE id=?`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *TrafficRepository) GetForwardGroup(ctx context.Context, id int64) (*ForwardGroup, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return nil, err
	}
	g, err := scanForwardGroup(r.db.QueryRowContext(ctx,
		`SELECT id, name, balance_strategy, dns_domain, dns_provider_id, failover_enabled, offline_ms_threshold, created_at, updated_at
		 FROM forward_groups WHERE id=?`, id))
	if err != nil {
		return nil, err
	}
	members, err := r.ListForwardGroupMembers(ctx, id)
	if err != nil {
		return nil, err
	}
	g.Members = members
	return g, nil
}

func (r *TrafficRepository) ListForwardGroups(ctx context.Context) ([]*ForwardGroup, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, balance_strategy, dns_domain, dns_provider_id, failover_enabled, offline_ms_threshold, created_at, updated_at
		 FROM forward_groups ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*ForwardGroup
	for rows.Next() {
		g, err := scanForwardGroup(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

func scanForwardGroup(s rowScanner) (*ForwardGroup, error) {
	var g ForwardGroup
	var failover int
	if err := s.Scan(&g.ID, &g.Name, &g.BalanceStrategy, &g.DNSDomain, &g.DNSProviderID, &failover, &g.OfflineMsThreshold, &g.CreatedAt, &g.UpdatedAt); err != nil {
		return nil, err
	}
	g.FailoverEnabled = failover != 0
	return &g, nil
}

// ---------- 组成员(全量替换)----------

func (r *TrafficRepository) SetForwardGroupMembers(ctx context.Context, groupID int64, members []ForwardGroupMember) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM forward_group_members WHERE group_id=?`, groupID); err != nil {
		return err
	}
	for i, m := range members {
		w := m.Weight
		if w <= 0 {
			w = 1
		}
		seq := m.Seq
		if seq == 0 {
			seq = i
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO forward_group_members (group_id, server_id, weight, seq) VALUES (?, ?, ?, ?)`,
			groupID, m.ServerID, w, seq); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return r.RebuildForwardBranchesForGroup(ctx, groupID)
}

func (r *TrafficRepository) ListForwardGroupMembers(ctx context.Context, groupID int64) ([]ForwardGroupMember, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT group_id, server_id, weight, COALESCE(seq,0) FROM forward_group_members WHERE group_id=? ORDER BY seq, server_id`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ForwardGroupMember
	for rows.Next() {
		var m ForwardGroupMember
		if err := rows.Scan(&m.GroupID, &m.ServerID, &m.Weight, &m.Seq); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// ---------- 转发链 ----------

func (r *TrafficRepository) CreateForwardChain(ctx context.Context, name string) (int64, error) {
	return r.CreateForwardChainMeta(ctx, &ForwardChain{Name: name})
}

func (r *TrafficRepository) CreateForwardChainMeta(ctx context.Context, c *ForwardChain) (int64, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return 0, err
	}
	if c == nil {
		return 0, fmt.Errorf("chain is nil")
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO forward_chains (name, port_range_start, port_range_end, dns_domain, dns_domain_v6, dns_provider_id, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		c.Name, c.PortRangeStart, c.PortRangeEnd, c.DNSDomain, c.DNSDomainV6, c.DNSProviderID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *TrafficRepository) UpdateForwardChain(ctx context.Context, c *ForwardChain) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	if c == nil || c.ID == 0 {
		return fmt.Errorf("非法链 id")
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE forward_chains SET name=COALESCE(NULLIF(?, ''), name), port_range_start=?, port_range_end=?, dns_domain=?, dns_domain_v6=?, dns_provider_id=?, updated_at=CURRENT_TIMESTAMP
		 WHERE id=?`,
		c.Name, c.PortRangeStart, c.PortRangeEnd, c.DNSDomain, c.DNSDomainV6, c.DNSProviderID, c.ID)
	return err
}

func (r *TrafficRepository) DeleteForwardChain(ctx context.Context, id int64) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM forward_chain_hops WHERE chain_id=?`, id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM forward_chain_branches WHERE chain_id=?`, id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM forward_daily_traffic WHERE chain_id=?`, id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM forward_chain_nodes WHERE chain_id=?`, id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM forward_chains WHERE id=?`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *TrafficRepository) GetForwardChain(ctx context.Context, id int64) (*ForwardChain, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return nil, err
	}
	var c ForwardChain
	if err := r.db.QueryRowContext(ctx,
		`SELECT id, name, COALESCE(port_range_start,0), COALESCE(port_range_end,0), COALESCE(dns_domain,''), COALESCE(dns_domain_v6,''), COALESCE(dns_provider_id,0), created_at, updated_at
		 FROM forward_chains WHERE id=?`, id).
		Scan(&c.ID, &c.Name, &c.PortRangeStart, &c.PortRangeEnd, &c.DNSDomain, &c.DNSDomainV6, &c.DNSProviderID, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	hops, err := r.ListForwardChainHops(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Hops = hops
	if err := r.decorateForwardHops(ctx, c.Hops); err != nil {
		return nil, err
	}
	if branches, err := r.ListForwardChainBranches(ctx, id); err == nil {
		c.Branches = branches
	}
	if daily, err := r.ListForwardDailyTraffic(ctx, id, 14); err == nil {
		c.DailyTraffic = daily
	}
	if bound, err := r.ListForwardChainBindingsByChain(ctx, id); err == nil {
		c.BoundNodes = bound
	}
	return &c, nil
}

func (r *TrafficRepository) ListForwardChains(ctx context.Context) ([]*ForwardChain, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, COALESCE(port_range_start,0), COALESCE(port_range_end,0), COALESCE(dns_domain,''), COALESCE(dns_domain_v6,''), COALESCE(dns_provider_id,0), created_at, updated_at
		 FROM forward_chains ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*ForwardChain
	for rows.Next() {
		var c ForwardChain
		if err := rows.Scan(&c.ID, &c.Name, &c.PortRangeStart, &c.PortRangeEnd, &c.DNSDomain, &c.DNSDomainV6, &c.DNSProviderID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for _, c := range out {
		hops, err := r.ListForwardChainHops(ctx, c.ID)
		if err != nil {
			return nil, err
		}
		c.Hops = hops
		if err := r.decorateForwardHops(ctx, c.Hops); err != nil {
			return nil, err
		}
		if branches, err := r.ListForwardChainBranches(ctx, c.ID); err == nil {
			c.Branches = branches
		}
		if bound, err := r.ListForwardChainBindingsByChain(ctx, c.ID); err == nil {
			c.BoundNodes = bound
		}
	}
	return out, nil
}

func (r *TrafficRepository) decorateForwardHops(ctx context.Context, hops []ForwardChainHop) error {
	for i := range hops {
		hops[i].Order = hops[i].Seq
		g, err := r.GetForwardGroup(ctx, hops[i].GroupID)
		if err != nil {
			continue
		}
		if g != nil {
			hops[i].GroupName = g.Name
		}
	}
	return nil
}

// SetForwardChainHops 全量替换链的跳序列(seq 由入参顺序决定)。
func (r *TrafficRepository) SetForwardChainHops(ctx context.Context, chainID int64, groupIDs []int64) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM forward_chain_hops WHERE chain_id=?`, chainID); err != nil {
		return err
	}
	for seq, gid := range groupIDs {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO forward_chain_hops (chain_id, seq, group_id) VALUES (?, ?, ?)`,
			chainID, seq, gid); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return r.RebuildForwardChainBranches(ctx, chainID)
}

func (r *TrafficRepository) ListForwardChainHops(ctx context.Context, chainID int64) ([]ForwardChainHop, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT chain_id, seq, group_id FROM forward_chain_hops WHERE chain_id=? ORDER BY seq`, chainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ForwardChainHop
	for rows.Next() {
		var h ForwardChainHop
		if err := rows.Scan(&h.ChainID, &h.Seq, &h.GroupID); err != nil {
			return nil, err
		}
		h.Order = h.Seq
		out = append(out, h)
	}
	return out, rows.Err()
}
