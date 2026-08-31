package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

// 套餐里所有按 node_id 索引的 JSON 映射列。节点被删后这些列里的残留项
// 会让该套餐再也保存不了,必须一并清掉。
var packageNodeKeyedColumns = []string{
	"node_multipliers",
	"node_name_overrides",
	"node_speed_limits",
	"node_device_limits",
	"node_traffic_limits",
}

// purgeNodeFromPackages 把已删节点从所有套餐中摘干净:节点列表 nodes 以及
// 五张按 node_id 索引的配置映射。
//
// 为什么必须做:这些都是 JSON 文本列、没有外键,删节点不会级联。留下悬空引用会让
// 该套餐**永远保存不了** —— validatePackageNodeTrafficLimits 一旦发现 limits 里存在
// 不在 nodes 列表中的 id,就直接报「节点 N 不在套餐内」并拒绝整次保存。
//
// JSON 文本列在 SQLite 与 PostgreSQL 之间没有可移植的原地改写语法,因此读出来在
// Go 侧修改,并且只回写真正发生变化的行。
func purgeNodeFromPackages(ctx context.Context, db *dialectDB, nodeID int64) error {
	if db == nil || nodeID <= 0 {
		return nil
	}

	query := `SELECT id, COALESCE(nodes, '[]')`
	for _, col := range packageNodeKeyedColumns {
		query += fmt.Sprintf(`, COALESCE(%s, '{}')`, col)
	}
	query += ` FROM packages`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("list packages for node purge: %w", err)
	}

	type pending struct {
		id      int64
		columns map[string]string
	}
	var updates []pending

	for rows.Next() {
		var id int64
		var nodesJSON string
		maps := make([]string, len(packageNodeKeyedColumns))
		dest := make([]any, 0, len(maps)+2)
		dest = append(dest, &id, &nodesJSON)
		for i := range maps {
			dest = append(dest, &maps[i])
		}
		if err := rows.Scan(dest...); err != nil {
			rows.Close()
			return fmt.Errorf("scan package for node purge: %w", err)
		}

		changed := make(map[string]string)
		if next, ok := removeIDFromJSONArray(nodesJSON, nodeID); ok {
			changed["nodes"] = next
		}
		key := strconv.FormatInt(nodeID, 10)
		for i, col := range packageNodeKeyedColumns {
			if next, ok := removeKeyFromJSONMap(maps[i], key); ok {
				changed[col] = next
			}
		}
		if len(changed) > 0 {
			updates = append(updates, pending{id: id, columns: changed})
		}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate packages for node purge: %w", err)
	}
	rows.Close()

	for _, u := range updates {
		setSQL := ""
		args := make([]any, 0, len(u.columns)+1)
		for col, value := range u.columns {
			if setSQL != "" {
				setSQL += ", "
			}
			setSQL += col + " = ?"
			args = append(args, value)
		}
		args = append(args, u.id)
		if _, err := db.ExecContext(ctx, `UPDATE packages SET `+setSQL+` WHERE id = ?`, args...); err != nil {
			return fmt.Errorf("purge node %d from package %d: %w", nodeID, u.id, err)
		}
	}
	return nil
}

// purgeNodeFromPackagesBestEffort 供删除路径调用:清理失败不应该让「删节点」这个
// 主操作失败(节点行本身已经删掉了),但必须留下日志,不能像原来那样静默吞掉。
func purgeNodeFromPackagesBestEffort(ctx context.Context, db *dialectDB, nodeID int64) {
	if err := purgeNodeFromPackages(ctx, db, nodeID); err != nil {
		log.Printf("[storage] 清理套餐中的已删节点 %d 失败(套餐可能因悬空引用而无法保存): %v", nodeID, err)
	}
}

// removeIDFromJSONArray 从 JSON 数字数组里移除指定 id。返回是否发生变化。
// 解析失败时保持原样 —— 坏数据不该阻塞删除动作。
func removeIDFromJSONArray(raw string, id int64) (string, bool) {
	if raw == "" || raw == "[]" {
		return raw, false
	}
	var ids []int64
	if err := json.Unmarshal([]byte(raw), &ids); err != nil {
		return raw, false
	}
	filtered := make([]int64, 0, len(ids))
	for _, v := range ids {
		if v != id {
			filtered = append(filtered, v)
		}
	}
	if len(filtered) == len(ids) {
		return raw, false
	}
	encoded, err := json.Marshal(filtered)
	if err != nil {
		return raw, false
	}
	return string(encoded), true
}

// removeKeyFromJSONMap 从 JSON 对象里移除指定键。用 json.RawMessage 承载值,
// 这样同一个函数能处理 float / int / string 各种映射,不必关心值类型。
func removeKeyFromJSONMap(raw, key string) (string, bool) {
	if raw == "" || raw == "{}" {
		return raw, false
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return raw, false
	}
	if _, ok := m[key]; !ok {
		return raw, false
	}
	delete(m, key)
	encoded, err := json.Marshal(m)
	if err != nil {
		return raw, false
	}
	return string(encoded), true
}
