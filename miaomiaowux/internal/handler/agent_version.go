package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"miaomiaowux/internal/auth"
	"miaomiaowux/internal/storage"
)

// AgentVersionHandler 给前端用:取目标服务器 agent 上报的当前版本,
// 再对照 GitHub 上 mmwx-agent / mmwx-agent-v* 的最新版本,算出能不能升级。
type AgentVersionHandler struct {
	rm   *RemoteManageHandler
	repo *storage.TrafficRepository
}

func NewAgentVersionHandler(rm *RemoteManageHandler, repo *storage.TrafficRepository) *AgentVersionHandler {
	return &AgentVersionHandler{rm: rm, repo: repo}
}

func (h *AgentVersionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("only GET"))
		return
	}
	if !userIsAdmin(r.Context(), h.repo, auth.UsernameFromContext(r.Context())) {
		writeError(w, http.StatusForbidden, errors.New("admin only"))
		return
	}
	sidStr := r.URL.Query().Get("server_id")
	sid, err := strconv.ParseInt(sidStr, 10, 64)
	if err != nil || sid <= 0 {
		writeBadRequest(w, "server_id required")
		return
	}

	current, currentErr := h.fetchAgentCurrent(r.Context(), sid)
	latest, latestErr := h.fetchLatest()

	resp := map[string]any{
		"server_id":         sid,
		"current":           current,
		"latest":            latest,
		"upgrade_available": isUpgradeAvailable(current, latest),
	}
	if currentErr != "" {
		resp["current_error"] = currentErr
	}
	if latestErr != "" {
		resp["latest_error"] = latestErr
	}
	// agent 不可达 → 返回 502,前端 react-query 视为 error 不写 data 缓存,
	// server.status 翻回 connected 后组件重新挂载时不会消费 stale "?" cache,自动 refetch。
	// 区分两类失败:
	//   - current 空 + currentErr 非空 → agent 不可达(或转发失败)= 502 BadGateway
	//   - current 空 + currentErr 空   → agent 可达但旧版未上报版本 = 200(保留 "?" 显示让升级提示生效)
	if current == "" && currentErr != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// fetchAgentCurrent 取目标 agent 的版本号。
// 旧 agent 不返回 agent_version 字段 → 空字符串,前端按"未知版本/需要升级"处理。
func (h *AgentVersionHandler) fetchAgentCurrent(ctx context.Context, serverID int64) (string, string) {
	// WS-first:新 agent 经 auth 上报了 agent_version 就直接用,不反向拉。
	// 端口隐身(HidePortOnWS)关闭入站后反向 HTTP 不可达;旧 agent 不上报则 fallback 反向 HTTP。
	if h.rm.wsHandler != nil {
		if conn, ok := h.rm.wsHandler.GetConnectionByServerID(serverID); ok && conn.AgentVersion != "" {
			return conn.AgentVersion, ""
		}
	}
	// 5s 超时即可,system-info 本来就是轻量 endpoint
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	body, err := h.rm.forwardToRemoteServer(ctx, serverID, http.MethodGet, "/api/child/system/info", nil)
	if err != nil {
		return "", err.Error()
	}
	var info struct {
		AgentVersion string `json:"agent_version"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return "", "parse system info: " + err.Error()
	}
	return strings.TrimSpace(info.AgentVersion), ""
}

// fetchLatest 查 GitHub 上 Agent 的最新版本号。复用面板更新那套 Release 缓存。
func (h *AgentVersionHandler) fetchLatest() (string, string) {
	repo := AgentGitHubRepo()
	if repo == "" {
		return "", "未配置 Agent 更新源（MMWX_AGENT_GITHUB_REPO=off）"
	}
	releases, err := fetchReleases(repo)
	if err != nil {
		return "", err.Error()
	}
	rel := pickAgentRelease(releases)
	if rel == nil {
		return "", "仓库里没有 Agent Release（tag 应为 mmwx-agent 或 mmwx-agent-v*）"
	}
	v := agentReleaseVersion(rel)
	if v == "" {
		return "", "Agent Release 读不出版本号（请把 tag 写成 mmwx-agent-v0.6.0，或把名字写成 v0.6.0）"
	}
	return v, ""
}

// isUpgradeAvailable 简单语义版本比对(major.minor.patch),非数字段降级到字符串比对。
//
//	latest 空         → 无法判断,不提示
//	current == latest → 不需要更新
//	current < latest  → 需要
func isUpgradeAvailable(current, latest string) bool {
	current = strings.TrimSpace(current)
	latest = strings.TrimSpace(latest)
	if latest == "" {
		// 不知道目标版本 — 没法判断,默认不报"可升级",避免误导
		return false
	}
	if current == "" {
		return true
	}
	return compareSemver(current, latest) < 0
}

// compareSemver:正数=a>b, 负数=a<b, 0=相等;非数字段走字符串比较。
func compareSemver(a, b string) int {
	pa := strings.Split(a, ".")
	pb := strings.Split(b, ".")
	n := len(pa)
	if len(pb) > n {
		n = len(pb)
	}
	for i := 0; i < n; i++ {
		var av, bv string
		if i < len(pa) {
			av = pa[i]
		}
		if i < len(pb) {
			bv = pb[i]
		}
		ai, aerr := strconv.Atoi(av)
		bi, berr := strconv.Atoi(bv)
		if aerr == nil && berr == nil {
			if ai != bi {
				if ai < bi {
					return -1
				}
				return 1
			}
			continue
		}
		if av != bv {
			if av < bv {
				return -1
			}
			return 1
		}
	}
	return 0
}
