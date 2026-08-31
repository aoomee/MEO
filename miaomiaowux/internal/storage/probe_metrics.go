package storage

import (
	"context"
	"fmt"
)

// ProbeMetricSlot 是探针延迟 5 分钟聚合槽的落库行。仅 PostgreSQL 使用;
// SQLite 走内存 + 文件 checkpoint,避免分钟级时序把 WAL 涨到几十 GB。
type ProbeMetricSlot struct {
	ServerID  int64
	TargetKey string
	Slot      int64
	Sum       int64
	Cnt       int64
	Fail      int64
}

func (r *TrafficRepository) ensureProbeMetricsTables(ctx context.Context) error {
	if r == nil || r.db == nil {
		return nil
	}
	_, err := r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS probe_metric_slots (
    server_id INTEGER NOT NULL,
    target_key TEXT NOT NULL,
    slot INTEGER NOT NULL,
    sum_ms INTEGER NOT NULL DEFAULT 0,
    cnt INTEGER NOT NULL DEFAULT 0,
    fail INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (server_id, target_key, slot)
);
CREATE INDEX IF NOT EXISTS idx_probe_metric_slots_slot ON probe_metric_slots(slot);
`)
	if err != nil {
		return fmt.Errorf("ensure probe_metric_slots: %w", err)
	}
	return nil
}

func (r *TrafficRepository) ProbeMetricsPersistReady() bool {
	return r != nil && r.config.Driver == "postgres"
}

func (r *TrafficRepository) UpsertProbeMetricSlots(ctx context.Context, slots []ProbeMetricSlot) error {
	if !r.ProbeMetricsPersistReady() || len(slots) == 0 {
		return nil
	}
	if err := r.ensureProbeMetricsTables(ctx); err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, slot := range slots {
		if slot.ServerID <= 0 || slot.TargetKey == "" || slot.Slot <= 0 {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
INSERT INTO probe_metric_slots (server_id, target_key, slot, sum_ms, cnt, fail)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(server_id, target_key, slot) DO UPDATE SET
  sum_ms=excluded.sum_ms, cnt=excluded.cnt, fail=excluded.fail`,
			slot.ServerID, slot.TargetKey, slot.Slot, slot.Sum, slot.Cnt, slot.Fail); err != nil {
			return fmt.Errorf("upsert probe metric slot: %w", err)
		}
	}
	return tx.Commit()
}

func (r *TrafficRepository) ListProbeMetricSlots(ctx context.Context, since int64) ([]ProbeMetricSlot, error) {
	if !r.ProbeMetricsPersistReady() {
		return nil, nil
	}
	if err := r.ensureProbeMetricsTables(ctx); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT server_id, target_key, slot, sum_ms, cnt, fail
FROM probe_metric_slots WHERE slot>=? ORDER BY server_id, target_key, slot`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ProbeMetricSlot
	for rows.Next() {
		var slot ProbeMetricSlot
		if err := rows.Scan(&slot.ServerID, &slot.TargetKey, &slot.Slot, &slot.Sum, &slot.Cnt, &slot.Fail); err != nil {
			return nil, err
		}
		out = append(out, slot)
	}
	return out, rows.Err()
}

func (r *TrafficRepository) CleanOldProbeMetricSlots(ctx context.Context, before int64) error {
	if !r.ProbeMetricsPersistReady() {
		return nil
	}
	if err := r.ensureProbeMetricsTables(ctx); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM probe_metric_slots WHERE slot < ?`, before)
	if err != nil {
		return fmt.Errorf("clean old probe metric slots: %w", err)
	}
	return nil
}
