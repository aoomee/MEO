package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"miaomiaowux/internal/storage"
)

// ProbeForwardHandler 提供 GET /api/public/probe-forward。
// 返回前端探针转发视图需要的链路数据：{chains:[{name,end_to_end_ms,loss_pct,groups,trend,traffic}]}。
type ProbeForwardHandler struct {
	repo *storage.TrafficRepository
}

func NewProbeForwardHandler(repo *storage.TrafficRepository) *ProbeForwardHandler {
	return &ProbeForwardHandler{repo: repo}
}

func (h *ProbeForwardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	payload, err := h.buildForwardPayload(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"chains": []any{}})
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *ProbeForwardHandler) buildForwardPayload(ctx context.Context) (map[string]any, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	chains, err := h.repo.ListForwardChains(ctx)
	if err != nil {
		return nil, err
	}
	since := time.Now().Add(-2 * time.Hour)
	metrics, _ := h.repo.ListRecentForwardHopMetrics(ctx, since)
	latest := latestMetricsByServer(metrics)
	trendByChain := trendByChain(metrics)

	out := make([]map[string]any, 0, len(chains))
	for _, chain := range chains {
		if chain == nil {
			continue
		}
		groups := make([]map[string]any, 0, len(chain.Hops))
		var e2e int64
		var lossSum float64
		var lossN int
		for i, hop := range chain.Hops {
			role := "mid"
			if i == 0 {
				role = "entry"
			} else if i == len(chain.Hops)-1 {
				role = "exit"
			}
			g, err := h.repo.GetForwardGroup(ctx, hop.GroupID)
			if err != nil || g == nil {
				continue
			}
			servers := make([]map[string]any, 0, len(g.Members))
			var hopRTT int64
			var hopN int
			for _, m := range g.Members {
				name := "server"
				if s, serr := h.repo.GetRemoteServer(ctx, m.ServerID); serr == nil && s != nil && s.Name != "" {
					name = s.Name
				}
				met := latest[m.ServerID]
				toNext := int64(0)
				healthy := false
				if met != nil {
					toNext = met.RTTMs
					healthy = met.Healthy
					if toNext > 0 {
						hopRTT += toNext
						hopN++
					}
					if !healthy {
						lossSum += 100
					}
					lossN++
				}
				servers = append(servers, map[string]any{
					"name":       name,
					"healthy":    healthy,
					"to_next_ms": toNext,
				})
			}
			avg := int64(0)
			if hopN > 0 {
				avg = hopRTT / int64(hopN)
			}
			if role != "exit" {
				e2e += avg
			}
			groups = append(groups, map[string]any{
				"name":       firstNonEmpty(hop.GroupName, g.Name),
				"role":       role,
				"to_next_ms": avg,
				"servers":    servers,
			})
		}
		loss := 0.0
		if lossN > 0 {
			loss = lossSum / float64(lossN)
		}
		daily, _ := h.repo.ListForwardDailyTraffic(ctx, chain.ID, 7)
		var totalBytes uint64
		trafficServers := map[int64]uint64{}
		for _, d := range daily {
			totalBytes += d.BytesUp + d.BytesDown
			trafficServers[d.ServerID] += d.BytesUp + d.BytesDown
		}
		ts := make([]map[string]any, 0, len(trafficServers))
		for sid, bytes := range trafficServers {
			name := "server"
			if s, serr := h.repo.GetRemoteServer(ctx, sid); serr == nil && s != nil && s.Name != "" {
				name = s.Name
			}
			ts = append(ts, map[string]any{"name": name, "bytes": bytes})
		}
		trend := trendByChain[chain.ID]
		if trend == nil {
			trend = []map[string]any{}
		}
		out = append(out, map[string]any{
			"name":          chain.Name,
			"end_to_end_ms": e2e,
			"loss_pct":      loss,
			"groups":        groups,
			"trend":         trend,
			"traffic":       map[string]any{"total_gb": float64(totalBytes) / (1024 * 1024 * 1024), "servers": ts},
		})
	}
	return map[string]any{"chains": out}, nil
}

func latestMetricsByServer(rows []storage.ForwardHopMetric) map[int64]*storage.ForwardHopMetric {
	out := map[int64]*storage.ForwardHopMetric{}
	for i := range rows {
		m := &rows[i]
		prev := out[m.ServerID]
		if prev == nil || m.At.After(prev.At) {
			out[m.ServerID] = m
		}
	}
	return out
}

func trendByChain(rows []storage.ForwardHopMetric) map[int64][]map[string]any {
	type bucket struct {
		sum int64
		n   int
	}
	byChain := map[int64]map[int64]*bucket{}
	for _, m := range rows {
		chainID, ok := storage.ParseForwardRuleChainID(m.RuleID)
		if !ok || m.RTTMs <= 0 {
			continue
		}
		key := m.At.Unix() / 300
		if byChain[chainID] == nil {
			byChain[chainID] = map[int64]*bucket{}
		}
		b := byChain[chainID][key]
		if b == nil {
			b = &bucket{}
			byChain[chainID][key] = b
		}
		b.sum += m.RTTMs
		b.n++
	}
	out := map[int64][]map[string]any{}
	for chainID, buckets := range byChain {
		keys := make([]int64, 0, len(buckets))
		for k := range buckets {
			keys = append(keys, k)
		}
		// insertion order is enough for small maps; frontend sorts visually
		for _, k := range keys {
			b := buckets[k]
			avg := int64(0)
			if b.n > 0 {
				avg = b.sum / int64(b.n)
			}
			out[chainID] = append(out[chainID], map[string]any{"e2e_ms": avg})
		}
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
