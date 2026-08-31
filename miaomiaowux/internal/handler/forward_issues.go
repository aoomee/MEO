package handler

import (
	"context"
	"fmt"
	"strconv"

	"miaomiaowux/internal/storage"
)

// evaluateForwardChainIssues 对照一条链的拓扑,返回给面板展示的中文原因。
// 前端 GET /api/admin/forward/chains 读的是顶层 issues[chain.id] 字符串数组。
func evaluateForwardChainIssues(chain *storage.ForwardChain, groups map[int64]*storage.ForwardGroup, servers map[int64]*storage.RemoteServer) []string {
	if chain == nil {
		return []string{"转发链不存在"}
	}
	if len(chain.Hops) < 2 {
		return []string{"链至少需要入口和出口两组"}
	}
	var issues []string
	for i, hop := range chain.Hops {
		role := "中间"
		if i == 0 {
			role = "入口"
		} else if i == len(chain.Hops)-1 {
			role = "出口"
		}
		g := groups[hop.GroupID]
		if g == nil {
			issues = append(issues, fmt.Sprintf("第 %d 跳（%s）引用了不存在的组", i+1, role))
			continue
		}
		name := g.Name
		if name == "" {
			name = hop.GroupName
		}
		if name == "" {
			name = strconv.FormatInt(hop.GroupID, 10)
		}
		if len(g.Members) == 0 {
			issues = append(issues, fmt.Sprintf("%s组「%s」没有服务器", role, name))
			continue
		}
		usable, missing := 0, 0
		for _, m := range g.Members {
			if forwardServerHost(servers[m.ServerID]) == "" {
				missing++
				continue
			}
			usable++
		}
		if usable == 0 {
			issues = append(issues, fmt.Sprintf("%s组「%s」的服务器没有可用地址", role, name))
		} else if missing > 0 {
			issues = append(issues, fmt.Sprintf("%s组「%s」有 %d 台服务器缺少地址", role, name, missing))
		}
	}
	return issues
}

func (h *ForwardHandler) collectForwardChainIssues(ctx context.Context, chains []*storage.ForwardChain) map[int64][]string {
	out := map[int64][]string{}
	if h == nil || h.repo == nil || len(chains) == 0 {
		return out
	}
	groups := map[int64]*storage.ForwardGroup{}
	if all, err := h.repo.ListForwardGroups(ctx); err == nil {
		for _, g := range all {
			if g == nil {
				continue
			}
			members, _ := h.repo.ListForwardGroupMembers(ctx, g.ID)
			g.Members = members
			groups[g.ID] = g
		}
	}
	servers := map[int64]*storage.RemoteServer{}
	if list, err := h.repo.ListRemoteServers(ctx); err == nil {
		for i := range list {
			s := list[i]
			servers[s.ID] = &s
		}
	}
	for _, chain := range chains {
		if chain == nil {
			continue
		}
		if issues := evaluateForwardChainIssues(chain, groups, servers); len(issues) > 0 {
			out[chain.ID] = issues
		}
	}
	return out
}

func (h *ForwardHandler) chainWarnings(ctx context.Context, chainID int64) []string {
	if h == nil || h.repo == nil || chainID <= 0 {
		return nil
	}
	chain, err := h.repo.GetForwardChain(ctx, chainID)
	if err != nil || chain == nil {
		return nil
	}
	return h.collectForwardChainIssues(ctx, []*storage.ForwardChain{chain})[chainID]
}
