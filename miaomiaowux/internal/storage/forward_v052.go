package storage

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ForwardChainBranch 官方 v0.5.2 forward_chain_branches：把「链跳 → 组」展开成
// 每个组成员一台服务器。hop_seq 是跳序号，seq 是组内序号，via_group_id 是来源组。
type ForwardChainBranch struct {
	ChainID    int64 `json:"chain_id"`
	HopSeq     int   `json:"hop_seq"`
	ServerID   int64 `json:"server_id"`
	Seq        int   `json:"seq"`
	ViaGroupID int64 `json:"via_group_id"`
}

// ForwardDailyTraffic 官方 v0.5.2 forward_daily_traffic：按日累计某链在某服务器上的转发字节。
// 采样点是 gauge（累计值），落库取当日 max。
type ForwardDailyTraffic struct {
	ChainID   int64     `json:"chain_id"`
	ServerID  int64     `json:"server_id"`
	Date      string    `json:"date"`
	BytesUp   uint64    `json:"bytes_up"`
	BytesDown uint64    `json:"bytes_down"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ParseForwardRuleChainID 从 agent 规则 id（fwd-c{chain}-p{port}-h{hop}）抽出 chain_id。
func ParseForwardRuleChainID(ruleID string) (int64, bool) {
	ruleID = strings.TrimSpace(ruleID)
	if !strings.HasPrefix(ruleID, "fwd-c") {
		return 0, false
	}
	rest := strings.TrimPrefix(ruleID, "fwd-c")
	i := strings.IndexByte(rest, '-')
	if i <= 0 {
		return 0, false
	}
	id, err := strconv.ParseInt(rest[:i], 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

// RebuildForwardChainBranches 按当前 hops + 组成员全量重写一条链的 branches。
func (r *TrafficRepository) RebuildForwardChainBranches(ctx context.Context, chainID int64) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	hops, err := r.ListForwardChainHops(ctx, chainID)
	if err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM forward_chain_branches WHERE chain_id=?`, chainID); err != nil {
		return err
	}
	for _, hop := range hops {
		members, err := r.ListForwardGroupMembers(ctx, hop.GroupID)
		if err != nil {
			return err
		}
		for seq, m := range members {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO forward_chain_branches (chain_id, hop_seq, server_id, seq, via_group_id) VALUES (?,?,?,?,?)`,
				chainID, hop.Seq, m.ServerID, seq, hop.GroupID); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

// RebuildForwardBranchesForGroup 组成员变更后，重写所有引用该组的链的 branches。
func (r *TrafficRepository) RebuildForwardBranchesForGroup(ctx context.Context, groupID int64) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT DISTINCT chain_id FROM forward_chain_hops WHERE group_id=?`, groupID)
	if err != nil {
		return err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, id := range ids {
		if err := r.RebuildForwardChainBranches(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *TrafficRepository) ListForwardChainBranches(ctx context.Context, chainID int64) ([]ForwardChainBranch, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT chain_id, hop_seq, server_id, seq, via_group_id FROM forward_chain_branches WHERE chain_id=? ORDER BY hop_seq, seq`,
		chainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ForwardChainBranch
	for rows.Next() {
		var b ForwardChainBranch
		if err := rows.Scan(&b.ChainID, &b.HopSeq, &b.ServerID, &b.Seq, &b.ViaGroupID); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *TrafficRepository) ListForwardDailyTraffic(ctx context.Context, chainID int64, days int) ([]ForwardDailyTraffic, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return nil, err
	}
	if days <= 0 {
		days = 14
	}
	since := time.Now().AddDate(0, 0, -days+1).Format("2006-01-02")
	rows, err := r.db.QueryContext(ctx,
		`SELECT chain_id, server_id, date, bytes_up, bytes_down, updated_at
		 FROM forward_daily_traffic WHERE chain_id=? AND date>=? ORDER BY date, server_id`,
		chainID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ForwardDailyTraffic
	for rows.Next() {
		var d ForwardDailyTraffic
		if err := rows.Scan(&d.ChainID, &d.ServerID, &d.Date, &d.BytesUp, &d.BytesDown, &d.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// UpsertForwardDailyTraffic 按日取 max(累计字节)。agent 上报的是 gauge。
func (r *TrafficRepository) UpsertForwardDailyTraffic(ctx context.Context, chainID, serverID int64, date string, up, down uint64) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO forward_daily_traffic (chain_id, server_id, date, bytes_up, bytes_down, updated_at)
		VALUES (?,?,?,?,?,CURRENT_TIMESTAMP)
		ON CONFLICT(chain_id, server_id, date) DO UPDATE SET
		  bytes_up = MAX(bytes_up, excluded.bytes_up),
		  bytes_down = MAX(bytes_down, excluded.bytes_down),
		  updated_at = CURRENT_TIMESTAMP`,
		chainID, serverID, date, up, down)
	if err != nil {
		return fmt.Errorf("upsert forward daily traffic: %w", err)
	}
	return nil
}

// CleanOldForwardDailyTraffic 按日清理超期的转发链每日流量。before 当天本身保留。
func (r *TrafficRepository) CleanOldForwardDailyTraffic(ctx context.Context, before time.Time) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM forward_daily_traffic WHERE date < ?`, before.Format("2006-01-02"))
	if err != nil {
		return fmt.Errorf("clean old forward daily traffic: %w", err)
	}
	return nil
}
