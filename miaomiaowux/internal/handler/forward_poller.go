package handler

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"miaomiaowux/internal/storage"
)

// 转发状态统一 poller:每 tick 拉一次全量转发运行时状态,同时干两件事——
//  1. #8 状态持久化:把每 upstream 的健康/RTT/字节写进 forward_hop_metrics 时序表(供前端画曲线);
//  2. #7 入口组 DNS 编排:据健康态把设了 DNSDomain 的入口组的 A 记录集调成健康成员 IP。
//
// 健康态直接采用 agent 侧已带滞回门(failover/recover 双窗)的 Healthy,DNS 层不再叠加防抖。

const forwardPollInterval = 17 * time.Second // 避开整点抖动

func forwardRetentionWindows(ctx context.Context, repo *storage.TrafficRepository) (metrics, daily time.Duration) {
	metricsDays := readClampedSetting(ctx, repo, forwardMetricsRetentionDaysKey, defaultForwardMetricsRetention, 1, 7)
	dailyDays := readClampedSetting(ctx, repo, forwardDailyRetentionDaysKey, defaultForwardDailyRetention, 1, 30)
	return time.Duration(metricsDays) * 24 * time.Hour, time.Duration(dailyDays) * 24 * time.Hour
}

// StartForwardStatusPoller 常驻轮询,ctx 取消时退出。由 main 用 collectorCtx 启动。
func (h *RemoteManageHandler) StartForwardStatusPoller(ctx context.Context) {
	ticker := time.NewTicker(forwardPollInterval)
	defer ticker.Stop()
	var lastClean time.Time
	log.Printf("[forward] status poller started, interval=%s", forwardPollInterval)
	for {
		select {
		case <-ctx.Done():
			log.Printf("[forward] status poller stopped")
			return
		case <-ticker.C:
			h.pollForwardStatusOnce(ctx)
			// 每天清一次过期采样点(首 tick 也会清)。
			if now := time.Now(); now.Sub(lastClean) >= 24*time.Hour {
				metricsKeep, dailyKeep := forwardRetentionWindows(ctx, h.repo)
				if err := h.repo.CleanOldForwardHopMetrics(ctx, now.Add(-metricsKeep)); err != nil {
					log.Printf("[forward] 清理过期采样点失败: %v", err)
				}
				if err := h.repo.CleanOldForwardDailyTraffic(ctx, now.Add(-dailyKeep)); err != nil {
					log.Printf("[forward] 清理过期每日流量失败: %v", err)
				}
				lastClean = now
			}
		}
	}
}

// pollForwardStatusOnce 一轮:拉全量状态 → 落时序表 → 编排入口组 DNS。逐台错误 continue,不中断整轮。
func (h *RemoteManageHandler) pollForwardStatusOnce(ctx context.Context) {
	topo, err := h.loadForwardTopology(ctx)
	if err != nil {
		log.Printf("[forward] poller 加载拓扑失败: %v", err)
		return
	}
	sids := forwardingServerIDs(topo)
	if len(sids) == 0 {
		return
	}

	statusByServer := make(map[int64][]forwardRuleStatus, len(sids))
	anyOK := false
	var rows []storage.ForwardHopMetric
	for _, sid := range sids {
		rules, err := h.FetchForwardStatus(ctx, sid)
		if err != nil {
			log.Printf("[forward] 拉服务器 %d 状态失败: %v", sid, err)
			continue
		}
		anyOK = true
		statusByServer[sid] = rules
		for _, ru := range rules {
			for _, up := range ru.Upstreams {
				rows = append(rows, storage.ForwardHopMetric{
					ServerID:     sid,
					RuleID:       ru.RuleID,
					UpstreamAddr: up.Addr,
					Healthy:      up.Healthy,
					RTTMs:        up.RTTMs,
					BytesUp:      up.BytesUp,
					BytesDown:    up.BytesDown,
				})
			}
		}
	}
	if err := h.repo.InsertForwardHopMetrics(ctx, rows); err != nil {
		log.Printf("[forward] 写入采样点失败: %v", err)
	}

	// 全量拉取失败(所有 agent 都不可达)时不做 DNS 编排,避免把入口域名整体清空(误摘)。
	if anyOK {
		h.reconcileEntryGroupDNS(ctx, topo, statusByServer)
	}
}

// reconcileEntryGroupDNS 对每条绑定链的入口组(设了 DNSDomain 的)按健康成员 IP 编排 A 记录集。
func (h *RemoteManageHandler) reconcileEntryGroupDNS(ctx context.Context, topo forwardTopology, statusByServer map[int64][]forwardRuleStatus) {
	if h.ddnsManager == nil {
		return
	}
	for _, b := range topo.bindings {
		chain := topo.chains[b.ChainID]
		if chain == nil || len(chain.Hops) < 2 {
			continue
		}
		entry := topo.groups[chain.Hops[0].GroupID]
		if entry == nil || strings.TrimSpace(entry.DNSDomain) == "" {
			continue
		}
		// 入口组第 0 跳规则 id 与 buildForwardRules 一致:成员在此规则上有健康 upstream 即视为可用入口。
		hop0RuleID := fmt.Sprintf("fwd-c%d-p%d-h%d", chain.ID, b.Port, 0)
		var desiredIPs []string
		for _, m := range entry.Members {
			if !entryMemberHealthy(statusByServer[m.ServerID], hop0RuleID) {
				continue
			}
			host := forwardServerHost(topo.servers[m.ServerID])
			if ip := net.ParseIP(host); ip != nil && ip.To4() != nil {
				desiredIPs = append(desiredIPs, host)
			}
		}
		if err := h.ddnsManager.ReconcileGroupDNS(ctx, entry.DNSProviderID, entry.DNSDomain, desiredIPs); err != nil {
			log.Printf("[forward] 入口组 %q(%s)DNS 编排失败: %v", entry.Name, entry.DNSDomain, err)
		}
	}
}

// entryMemberHealthy 入口成员在其第 0 跳规则上是否有至少一个健康 upstream(能打通下一跳)。
// 状态缺失(agent 不可达)→ false → 从入口域名摘除。
func entryMemberHealthy(rules []forwardRuleStatus, ruleID string) bool {
	for _, ru := range rules {
		if ru.RuleID != ruleID {
			continue
		}
		for _, up := range ru.Upstreams {
			if up.Healthy {
				return true
			}
		}
	}
	return false
}
