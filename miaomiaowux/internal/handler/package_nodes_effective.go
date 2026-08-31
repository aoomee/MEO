package handler

import (
	"context"

	"miaomiaowux/internal/storage"
)

// effectivePackageNodeIDs implements the package UI contract: an empty node
// selection means every administrator-owned node, not zero nodes. User-imported
// private nodes are deliberately excluded from this global fallback.
//
// wasConfigured 区分两种空列表(所有读路径都会过滤掉已删节点,二者到这里都是空):
//   - false:管理员从未勾选过节点 → 保持"放行全部节点"的产品语义
//   - true :勾选过,但节点后来被删光 → 必须返回零节点
//
// 后者若也走 fallback,套餐会从"限定 N 个节点"静默变成"不限节点" —— 线上真实事故:
// 客户删掉某套餐仅有的两个节点后,该套餐的用户订阅里出现了全部节点。
func effectivePackageNodeIDs(ctx context.Context, repo *storage.TrafficRepository, configured []int64, wasConfigured bool) ([]int64, error) {
	if len(configured) > 0 {
		return append([]int64(nil), configured...), nil
	}
	if wasConfigured {
		return []int64{}, nil
	}
	users, err := repo.ListUsers(ctx, 1_000_000)
	if err != nil {
		return nil, err
	}
	admins := make(map[string]bool)
	for _, user := range users {
		if user.Role == storage.RoleAdmin {
			admins[user.Username] = true
		}
	}
	nodes, err := repo.ListAllNodes(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(nodes))
	// ListAllNodes is newest-first. The node management fallback order is the
	// authoritative order when a package has no explicit order.
	nodes = orderNodesByUserOrder(ctx, repo, firstAdminUsername(users), nodes)
	for _, node := range nodes {
		if admins[node.Username] && node.RoutedOwner != "user" {
			ids = append(ids, node.ID)
		}
	}
	return ids, nil
}

func firstAdminUsername(users []storage.User) string {
	for _, user := range users {
		if user.Role == storage.RoleAdmin {
			return user.Username
		}
	}
	return ""
}
