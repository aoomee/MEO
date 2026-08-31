package storage

import (
	"context"
	"fmt"
	"time"
)

// 转发跳的运行时时序:agent 上报的每 upstream 健康/RTT/字节,按 tick append。
// 追加写(无 upsert),按 (server_id, rule_id, at) 索引区间读,过期按 at 清理。

// ForwardHopMetric 一条采样点,对应 agent forward.RuleStatus 里的单个 upstream。
type ForwardHopMetric struct {
	ID           int64     `json:"id"`
	ServerID     int64     `json:"server_id"`
	RuleID       string    `json:"rule_id"`
	UpstreamAddr string    `json:"upstream_addr"`
	Healthy      bool      `json:"healthy"`
	RTTMs        int64     `json:"rtt_ms"`
	BytesUp      uint64    `json:"bytes_up"`
	BytesDown    uint64    `json:"bytes_down"`
	At           time.Time `json:"at"`
}

// InsertForwardHopMetrics 批量追加采样点(一个 poller tick 一批)。
func (r *TrafficRepository) InsertForwardHopMetrics(ctx context.Context, rows []ForwardHopMetric) error {
	if len(rows) == 0 {
		return nil
	}
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	const stmt = `INSERT INTO forward_hop_metrics (server_id, rule_id, upstream_addr, healthy, rtt_ms, bytes_up, bytes_down, at)
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`
	for _, m := range rows {
		if _, err := tx.ExecContext(ctx, stmt,
			m.ServerID, m.RuleID, m.UpstreamAddr, boolToInt(m.Healthy), m.RTTMs, m.BytesUp, m.BytesDown); err != nil {
			return fmt.Errorf("insert forward hop metric: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	date := time.Now().Format("2006-01-02")
	for _, m := range rows {
		chainID, ok := ParseForwardRuleChainID(m.RuleID)
		if !ok {
			continue
		}
		if err := r.UpsertForwardDailyTraffic(ctx, chainID, m.ServerID, date, m.BytesUp, m.BytesDown); err != nil {
			return err
		}
	}
	return nil
}

// GetForwardHopMetrics 读某台 server(可选某 rule)自 since 起的采样点,按时间升序。
// ruleID 为空则返回该 server 全部 rule 的点。
func (r *TrafficRepository) GetForwardHopMetrics(ctx context.Context, serverID int64, ruleID string, since time.Time) ([]ForwardHopMetric, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return nil, err
	}
	var query string
	var args []any
	if ruleID != "" {
		query = `SELECT id, server_id, rule_id, upstream_addr, healthy, rtt_ms, bytes_up, bytes_down, at
			FROM forward_hop_metrics WHERE server_id=? AND rule_id=? AND at>=? ORDER BY at ASC`
		args = []any{serverID, ruleID, since}
	} else {
		query = `SELECT id, server_id, rule_id, upstream_addr, healthy, rtt_ms, bytes_up, bytes_down, at
			FROM forward_hop_metrics WHERE server_id=? AND at>=? ORDER BY at ASC`
		args = []any{serverID, since}
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query forward hop metrics: %w", err)
	}
	defer rows.Close()
	var out []ForwardHopMetric
	for rows.Next() {
		var m ForwardHopMetric
		var healthy int
		if err := rows.Scan(&m.ID, &m.ServerID, &m.RuleID, &m.UpstreamAddr, &healthy, &m.RTTMs, &m.BytesUp, &m.BytesDown, &m.At); err != nil {
			return nil, fmt.Errorf("scan forward hop metric: %w", err)
		}
		m.Healthy = healthy != 0
		out = append(out, m)
	}
	return out, rows.Err()
}

// ListRecentForwardHopMetrics 读所有服务器自 since 起的采样点。
func (r *TrafficRepository) ListRecentForwardHopMetrics(ctx context.Context, since time.Time) ([]ForwardHopMetric, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, server_id, rule_id, upstream_addr, healthy, rtt_ms, bytes_up, bytes_down, at
		 FROM forward_hop_metrics WHERE at>=? ORDER BY at ASC`, since)
	if err != nil {
		return nil, fmt.Errorf("query recent forward hop metrics: %w", err)
	}
	defer rows.Close()
	var out []ForwardHopMetric
	for rows.Next() {
		var m ForwardHopMetric
		var healthy int
		if err := rows.Scan(&m.ID, &m.ServerID, &m.RuleID, &m.UpstreamAddr, &healthy, &m.RTTMs, &m.BytesUp, &m.BytesDown, &m.At); err != nil {
			return nil, fmt.Errorf("scan recent forward hop metric: %w", err)
		}
		m.Healthy = healthy != 0
		out = append(out, m)
	}
	return out, rows.Err()
}

// CleanOldForwardHopMetrics 删除早于 before 的采样点。
func (r *TrafficRepository) CleanOldForwardHopMetrics(ctx context.Context, before time.Time) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM forward_hop_metrics WHERE at < ?`, before)
	if err != nil {
		return fmt.Errorf("clean old forward hop metrics: %w", err)
	}
	return nil
}
