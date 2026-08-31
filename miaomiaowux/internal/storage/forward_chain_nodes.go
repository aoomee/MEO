package storage

import (
	"context"
	"database/sql"
	"errors"
)

// ForwardChainBinding 记录「某节点绑定到某转发链、用端口 port」。
// 官方 v0.5.2 还带落地地址、四层协议、钉住的出口机，避免删一个入口把整条链出口拆光。
type ForwardChainBinding struct {
	NodeID             int64  `json:"node_id"`
	ChainID            int64  `json:"chain_id"`
	Port               int    `json:"port"`
	TerminusAddr       string `json:"terminus_addr"`
	Protocol           string `json:"protocol"`
	PinnedExitServerID int64  `json:"pinned_exit_server_id"`
	NodeName           string `json:"node_name,omitempty"`
}

// BindForwardChainNode 绑定/改绑(node_id 唯一,存在则覆盖)。
func (r *TrafficRepository) BindForwardChainNode(ctx context.Context, nodeID, chainID int64, port int) error {
	return r.BindForwardChainNodeFull(ctx, ForwardChainBinding{
		NodeID:   nodeID,
		ChainID:  chainID,
		Port:     port,
		Protocol: "tcp",
	})
}

// BindForwardChainNodeFull 写入官方列 terminus_addr / protocol / pinned_exit_server_id。
func (r *TrafficRepository) BindForwardChainNodeFull(ctx context.Context, b ForwardChainBinding) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	if b.Protocol == "" {
		b.Protocol = "tcp"
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM forward_chain_nodes WHERE node_id=?`, b.NodeID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO forward_chain_nodes (node_id, chain_id, port, created_at, terminus_addr, protocol, pinned_exit_server_id)
		 VALUES (?, ?, ?, CURRENT_TIMESTAMP, ?, ?, ?)`,
		b.NodeID, b.ChainID, b.Port, b.TerminusAddr, b.Protocol, b.PinnedExitServerID); err != nil {
		return err
	}
	return tx.Commit()
}

// UnbindForwardChainNode 解绑(节点删除时调用)。
func (r *TrafficRepository) UnbindForwardChainNode(ctx context.Context, nodeID int64) error {
	if err := r.ensureForwardTables(ctx); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM forward_chain_nodes WHERE node_id=?`, nodeID)
	return err
}

const forwardBindingSelect = `SELECT f.node_id, f.chain_id, f.port,
	COALESCE(f.terminus_addr,''), COALESCE(f.protocol,'tcp'), COALESCE(f.pinned_exit_server_id,0),
	COALESCE(n.node_name,'')
	FROM forward_chain_nodes f
	LEFT JOIN nodes n ON n.id=f.node_id`

// GetForwardChainBinding 取某节点的绑定;无绑定返回 (nil, nil)。
func (r *TrafficRepository) GetForwardChainBinding(ctx context.Context, nodeID int64) (*ForwardChainBinding, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return nil, err
	}
	var b ForwardChainBinding
	err := r.db.QueryRowContext(ctx, forwardBindingSelect+` WHERE f.node_id=?`, nodeID).
		Scan(&b.NodeID, &b.ChainID, &b.Port, &b.TerminusAddr, &b.Protocol, &b.PinnedExitServerID, &b.NodeName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// ListForwardChainBindings 列出全部绑定(下发聚合时用)。
func (r *TrafficRepository) ListForwardChainBindings(ctx context.Context) ([]ForwardChainBinding, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return nil, err
	}
	return r.queryForwardBindings(ctx, forwardBindingSelect+` ORDER BY f.node_id`)
}

// ListForwardChainBindingsByChain 列出某链的全部绑定。
func (r *TrafficRepository) ListForwardChainBindingsByChain(ctx context.Context, chainID int64) ([]ForwardChainBinding, error) {
	if err := r.ensureForwardTables(ctx); err != nil {
		return nil, err
	}
	return r.queryForwardBindings(ctx, forwardBindingSelect+` WHERE f.chain_id=? ORDER BY f.node_id`, chainID)
}

func (r *TrafficRepository) queryForwardBindings(ctx context.Context, query string, args ...any) ([]ForwardChainBinding, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ForwardChainBinding
	for rows.Next() {
		var b ForwardChainBinding
		if err := rows.Scan(&b.NodeID, &b.ChainID, &b.Port, &b.TerminusAddr, &b.Protocol, &b.PinnedExitServerID, &b.NodeName); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}
