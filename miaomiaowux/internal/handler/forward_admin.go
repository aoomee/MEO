package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"miaomiaowux/internal/storage"
)

// ForwardHandler 转发组/转发链的 admin REST API。
// 路由挂在 /api/admin/forward/,内部按子路径分发。
type ForwardHandler struct {
	repo *storage.TrafficRepository
	// remote 用于「手动 apply 一条链的转发规则」(测试/运维触发);节点生命周期接入后由 P3 自动触发。
	remote *RemoteManageHandler
}

func NewForwardHandler(repo *storage.TrafficRepository, remote *RemoteManageHandler) *ForwardHandler {
	return &ForwardHandler{repo: repo, remote: remote}
}

var validBalanceStrategies = map[string]bool{"round_robin": true, "percentage": true, "cycle": true, "least_conn": true, "sticky": true, "weighted": true}

func (h *ForwardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/forward/")
	seg := strings.Split(strings.Trim(path, "/"), "/")

	switch {
	// /groups
	case len(seg) == 1 && seg[0] == "groups":
		switch r.Method {
		case http.MethodGet:
			h.listGroups(w, r)
		case http.MethodPost:
			h.createGroup(w, r)
		default:
			writeBadRequest(w, "不支持的方法")
		}
	// /groups/{id} 与 /groups/{id}/members
	case len(seg) >= 2 && seg[0] == "groups":
		id, err := strconv.ParseInt(seg[1], 10, 64)
		if err != nil {
			writeBadRequest(w, "非法组 id")
			return
		}
		if len(seg) == 3 && seg[2] == "members" {
			if r.Method == http.MethodPut {
				h.setMembers(w, r, id)
			} else {
				writeBadRequest(w, "不支持的方法")
			}
			return
		}
		switch r.Method {
		case http.MethodGet:
			h.getGroup(w, r, id)
		case http.MethodPut:
			h.updateGroup(w, r, id)
		case http.MethodDelete:
			h.deleteGroup(w, r, id)
		default:
			writeBadRequest(w, "不支持的方法")
		}
	// /status — 拉取全部转发服务器的运行时状态
	case len(seg) == 1 && seg[0] == "status":
		if r.Method == http.MethodGet {
			h.getStatus(w, r)
		} else {
			writeBadRequest(w, "不支持的方法")
		}
	// /metrics — 读取历史采样时序(server_id + rule_id + hours)
	case len(seg) == 1 && seg[0] == "metrics":
		if r.Method == http.MethodGet {
			h.getMetrics(w, r)
		} else {
			writeBadRequest(w, "不支持的方法")
		}
	// /chains
	case len(seg) == 1 && seg[0] == "chains":
		switch r.Method {
		case http.MethodGet:
			h.listChains(w, r)
		case http.MethodPost:
			h.createChain(w, r)
		default:
			writeBadRequest(w, "不支持的方法")
		}
	// /chains/{id}、/chains/{id}/hops、/chains/{id}/apply
	case len(seg) >= 2 && seg[0] == "chains":
		id, err := strconv.ParseInt(seg[1], 10, 64)
		if err != nil {
			writeBadRequest(w, "非法链 id")
			return
		}
		if len(seg) == 3 {
			switch seg[2] {
			case "hops":
				if r.Method == http.MethodPut {
					h.setHops(w, r, id)
				} else {
					writeBadRequest(w, "不支持的方法")
				}
			case "bind":
				switch r.Method {
				case http.MethodPost:
					h.bindNode(w, r, id)
				case http.MethodDelete:
					h.unbindNode(w, r, id)
				default:
					writeBadRequest(w, "不支持的方法")
				}
			case "create-node":
				if r.Method == http.MethodPost {
					h.createNode(w, r, id)
				} else {
					writeBadRequest(w, "不支持的方法")
				}
			case "rebind-node":
				if r.Method == http.MethodPost {
					h.rebindNode(w, r, id)
				} else {
					writeBadRequest(w, "不支持的方法")
				}
			case "daily-traffic":
				if r.Method == http.MethodGet {
					h.getDailyTraffic(w, r, id)
				} else {
					writeBadRequest(w, "不支持的方法")
				}
			default:
				writeBadRequest(w, "未知子路径")
			}
			return
		}
		switch r.Method {
		case http.MethodGet:
			h.getChain(w, r, id)
		case http.MethodPut:
			h.updateChain(w, r, id)
		case http.MethodDelete:
			h.deleteChain(w, r, id)
		default:
			writeBadRequest(w, "不支持的方法")
		}
	default:
		writeBadRequest(w, "未知路径")
	}
}

// ---------- 组 ----------

func (h *ForwardHandler) listGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.repo.ListForwardGroups(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true, "groups": groups})
}

func (h *ForwardHandler) getGroup(w http.ResponseWriter, r *http.Request, id int64) {
	g, err := h.repo.GetForwardGroup(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true, "group": g})
}

func (h *ForwardHandler) createGroup(w http.ResponseWriter, r *http.Request) {
	g, err := decodeGroup(r)
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}
	id, err := h.repo.CreateForwardGroup(r.Context(), g)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	// 允许创建时一并带成员
	if len(g.Members) > 0 {
		if err := h.repo.SetForwardGroupMembers(r.Context(), id, g.Members); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true, "id": id})
}

func (h *ForwardHandler) updateGroup(w http.ResponseWriter, r *http.Request, id int64) {
	g, err := decodeGroup(r)
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}
	g.ID = id
	if err := h.repo.UpdateForwardGroup(r.Context(), g); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *ForwardHandler) deleteGroup(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.repo.DeleteForwardGroup(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *ForwardHandler) setMembers(w http.ResponseWriter, r *http.Request, id int64) {
	var req struct {
		Members []storage.ForwardGroupMember `json:"members"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBadRequest(w, "请求格式错误")
		return
	}
	if err := h.repo.SetForwardGroupMembers(r.Context(), id, req.Members); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true})
}

func decodeGroup(r *http.Request) (*storage.ForwardGroup, error) {
	var g storage.ForwardGroup
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		return nil, errors.New("请求格式错误")
	}
	if strings.TrimSpace(g.Name) == "" {
		return nil, errors.New("组名不能为空")
	}
	if g.BalanceStrategy == "" {
		g.BalanceStrategy = "round_robin"
	}
	if !validBalanceStrategies[g.BalanceStrategy] {
		return nil, errors.New("非法均衡策略,应为 round_robin/percentage/cycle/least_conn/sticky/weighted")
	}
	return &g, nil
}

// ---------- 链 ----------

func (h *ForwardHandler) listChains(w http.ResponseWriter, r *http.Request) {
	chains, err := h.repo.ListForwardChains(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"chains":  chains,
		"issues":  h.collectForwardChainIssues(r.Context(), chains),
	})
}

func (h *ForwardHandler) getChain(w http.ResponseWriter, r *http.Request, id int64) {
	c, err := h.repo.GetForwardChain(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true, "chain": c})
}

func (h *ForwardHandler) getDailyTraffic(w http.ResponseWriter, r *http.Request, id int64) {
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days <= 0 {
		days = 14
	}
	rows, err := h.repo.ListForwardDailyTraffic(r.Context(), id, days)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true, "daily_traffic": rows})
}

func (h *ForwardHandler) createChain(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name           string  `json:"name"`
		GroupIDs       []int64 `json:"group_ids"`
		PortRangeStart int     `json:"port_range_start"`
		PortRangeEnd   int     `json:"port_range_end"`
		DNSDomain      string  `json:"dns_domain"`
		DNSDomainV6    string  `json:"dns_domain_v6"`
		DNSProviderID  int64   `json:"dns_provider_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBadRequest(w, "请求格式错误")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeBadRequest(w, "链名不能为空")
		return
	}
	id, err := h.repo.CreateForwardChainMeta(r.Context(), &storage.ForwardChain{
		Name:           req.Name,
		PortRangeStart: req.PortRangeStart,
		PortRangeEnd:   req.PortRangeEnd,
		DNSDomain:      strings.TrimSpace(req.DNSDomain),
		DNSDomainV6:    strings.TrimSpace(req.DNSDomainV6),
		DNSProviderID:  req.DNSProviderID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if len(req.GroupIDs) > 0 {
		if err := h.repo.SetForwardChainHops(r.Context(), id, req.GroupIDs); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		h.applyChainDNSToEntryGroup(r.Context(), id, strings.TrimSpace(req.DNSDomain), req.DNSProviderID)
	}
	resp := map[string]any{"success": true, "id": id}
	if warnings := h.chainWarnings(r.Context(), id); len(warnings) > 0 {
		resp["warnings"] = warnings
	}
	respondJSON(w, http.StatusOK, resp)
}

func (h *ForwardHandler) updateChain(w http.ResponseWriter, r *http.Request, id int64) {
	var req struct {
		Name           string `json:"name"`
		PortRangeStart int    `json:"port_range_start"`
		PortRangeEnd   int    `json:"port_range_end"`
		DNSDomain      string `json:"dns_domain"`
		DNSDomainV6    string `json:"dns_domain_v6"`
		DNSProviderID  int64  `json:"dns_provider_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBadRequest(w, "请求格式错误")
		return
	}
	if err := h.repo.UpdateForwardChain(r.Context(), &storage.ForwardChain{
		ID:             id,
		Name:           req.Name,
		PortRangeStart: req.PortRangeStart,
		PortRangeEnd:   req.PortRangeEnd,
		DNSDomain:      strings.TrimSpace(req.DNSDomain),
		DNSDomainV6:    strings.TrimSpace(req.DNSDomainV6),
		DNSProviderID:  req.DNSProviderID,
	}); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	h.applyChainDNSToEntryGroup(r.Context(), id, strings.TrimSpace(req.DNSDomain), req.DNSProviderID)
	respondJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *ForwardHandler) applyChainDNSToEntryGroup(ctx context.Context, chainID int64, domain string, providerID int64) {
	if domain == "" && providerID == 0 {
		return
	}
	chain, err := h.repo.GetForwardChain(ctx, chainID)
	if err != nil || chain == nil || len(chain.Hops) == 0 {
		return
	}
	g, err := h.repo.GetForwardGroup(ctx, chain.Hops[0].GroupID)
	if err != nil || g == nil {
		return
	}
	if domain != "" {
		g.DNSDomain = domain
	}
	if providerID != 0 {
		g.DNSProviderID = providerID
	}
	_ = h.repo.UpdateForwardGroup(ctx, g)
}

func (h *ForwardHandler) setHops(w http.ResponseWriter, r *http.Request, id int64) {
	var req struct {
		GroupIDs []int64 `json:"group_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBadRequest(w, "请求格式错误")
		return
	}
	if err := h.repo.SetForwardChainHops(r.Context(), id, req.GroupIDs); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"success":  true,
		"warnings": h.chainWarnings(r.Context(), id),
	})
}

func (h *ForwardHandler) deleteChain(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()
	bindings, _ := h.repo.ListForwardChainBindingsByChain(ctx, id)
	if h.remote != nil {
		_ = h.remote.TeardownExitNodes(ctx, id)
	}
	for _, b := range bindings {
		_ = h.repo.DeleteNode(ctx, b.NodeID, "admin")
	}
	if err := h.repo.DeleteForwardChain(ctx, id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true})
}

// bindNode 把一个节点绑定到本链、用指定端口,并同步下发该链涉及服务器的转发规则。
// 同时在出口组所有服务器上自动创建落地节点。
func (h *ForwardHandler) bindNode(w http.ResponseWriter, r *http.Request, chainID int64) {
	var req struct {
		NodeID int64 `json:"node_id"`
		Port   int   `json:"port"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBadRequest(w, "请求格式错误")
		return
	}
	if req.NodeID == 0 {
		writeBadRequest(w, "缺少 node_id")
		return
	}
	if req.Port <= 0 || req.Port > 65535 {
		writeBadRequest(w, "非法 port")
		return
	}
	chain, _ := h.repo.GetForwardChain(r.Context(), chainID)
	terminus, pin := h.computeTerminus(r.Context(), chain, req.Port, 0)
	if err := h.repo.BindForwardChainNodeFull(r.Context(), storage.ForwardChainBinding{
		NodeID:             req.NodeID,
		ChainID:            chainID,
		Port:               req.Port,
		TerminusAddr:       terminus,
		Protocol:           "tcp",
		PinnedExitServerID: pin,
	}); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.syncChain(r.Context(), chainID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	// 自动在出口组创建落地节点
	if h.remote != nil {
		if err := h.remote.ProvisionExitNodes(r.Context(), chainID, req.Port); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true, "message": "已绑定、下发转发规则并创建出口节点"})
}

// unbindNode 解绑节点并重新同步(撤掉该节点带来的转发规则),同时清理出口组的自动建站。
func (h *ForwardHandler) unbindNode(w http.ResponseWriter, r *http.Request, chainID int64) {
	nodeID, err := strconv.ParseInt(r.URL.Query().Get("node_id"), 10, 64)
	if err != nil || nodeID == 0 {
		writeBadRequest(w, "缺少或非法 node_id")
		return
	}
	ctx := r.Context()
	old, _ := h.repo.GetForwardChainBinding(ctx, nodeID)
	if err := h.repo.UnbindForwardChainNode(ctx, nodeID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.syncChain(ctx, chainID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	// 官方 v0.5.2：只拆这条绑定钉住的出口，剩下来的入口节点保住自己的落地。
	if h.remote != nil {
		remain, _ := h.repo.ListForwardChainBindingsByChain(ctx, chainID)
		if len(remain) == 0 {
			if err := h.remote.TeardownExitNodes(ctx, chainID); err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
		} else if old != nil && old.PinnedExitServerID > 0 {
			stillPinned := false
			for _, b := range remain {
				if b.PinnedExitServerID == old.PinnedExitServerID {
					stillPinned = true
					break
				}
			}
			if !stillPinned {
				_ = h.remote.TeardownExitNodeOnServer(ctx, chainID, old.PinnedExitServerID)
			}
		}
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true, "message": "已解绑、同步规则并清理出口节点"})
}

// createNode 对应官方 v0.4.8 POST /chains/{id}/create-node：
// 新建入口订阅节点(可按入口机拆开)并绑定到链,或把已有节点绑上来。
func (h *ForwardHandler) createNode(w http.ResponseWriter, r *http.Request, chainID int64) {
	var req struct {
		NodeName       string         `json:"node_name"`
		RelayProtocol  string         `json:"relay_protocol"`
		Port           int            `json:"port"`
		ExistingNodeID int64          `json:"existing_node_id"`
		EntrySeparate  bool           `json:"entry_separate"`
		ExitSeparate   bool           `json:"exit_separate"`
		Protocol       string         `json:"protocol"`
		Settings       map[string]any `json:"settings"`
		StreamSettings map[string]any `json:"stream_settings"`
		CertID         int64          `json:"cert_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBadRequest(w, "请求格式错误")
		return
	}
	chain, err := h.repo.GetForwardChain(r.Context(), chainID)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	port, err := pickChainPort(chain, req.Port)
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	ctx := r.Context()
	count := 0
	if req.ExistingNodeID > 0 {
		node, err := h.repo.GetNodeByID(ctx, req.ExistingNodeID)
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		if host := h.firstEntryHost(ctx, chain); host != "" {
			if rewritten, rerr := rewriteClashServerPort(node.ClashConfig, host, port); rerr == nil {
				node.ClashConfig = rewritten
				if _, uerr := h.repo.UpdateNode(ctx, node); uerr != nil {
					writeError(w, http.StatusInternalServerError, uerr)
					return
				}
			}
		}
		relayProto := strings.TrimSpace(req.RelayProtocol)
		if relayProto == "" {
			relayProto = "tcp"
		}
		terminus, pin := h.computeTerminus(ctx, chain, port, 0)
		if err := h.repo.BindForwardChainNodeFull(ctx, storage.ForwardChainBinding{
			NodeID:             req.ExistingNodeID,
			ChainID:            chainID,
			Port:               port,
			TerminusAddr:       terminus,
			Protocol:           relayProto,
			PinnedExitServerID: pin,
		}); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		count = 1
	} else {
		name := strings.TrimSpace(req.NodeName)
		if name == "" {
			name = chain.Name
		}
		hosts := h.entryHosts(ctx, chain, req.EntrySeparate)
		uuidStr := uuid.New().String()
		for i, host := range hosts {
			nodeName := name
			if len(hosts) > 1 {
				nodeName = fmt.Sprintf("%s-%d", name, i+1)
			}
			clash := syntheticVLESSClash(nodeName, host, port, uuidStr)
			created, err := h.repo.CreateNode(ctx, storage.Node{
				Username:       "admin",
				NodeName:       nodeName,
				Protocol:       "vless",
				ClashConfig:    clash,
				ParsedConfig:   clash,
				Enabled:        true,
				Tag:            "转发链入口",
				OriginalServer: host,
				NodeType:       "forward-entry",
			})
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			relayProto := strings.TrimSpace(req.RelayProtocol)
			if relayProto == "" {
				relayProto = "tcp"
			}
			pin := int64(0)
			if req.ExitSeparate {
				pin = h.nthExitServerID(ctx, chain, i)
			}
			terminus, pinned := h.computeTerminus(ctx, chain, port, pin)
			if err := h.repo.BindForwardChainNodeFull(ctx, storage.ForwardChainBinding{
				NodeID:             created.ID,
				ChainID:            chainID,
				Port:               port,
				TerminusAddr:       terminus,
				Protocol:           relayProto,
				PinnedExitServerID: pinned,
			}); err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			count++
		}
	}

	if err := h.syncChain(ctx, chainID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if h.remote != nil {
		if err := h.remote.ProvisionExitNodesSpec(ctx, chainID, port, exitProvisionSpec{
			Protocol:       req.Protocol,
			Settings:       req.Settings,
			StreamSettings: req.StreamSettings,
			CertID:         req.CertID,
		}); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true, "count": count, "port": port})
}

// rebindNode 把已绑定节点改挂到另一条链,沿用原端口。
func (h *ForwardHandler) rebindNode(w http.ResponseWriter, r *http.Request, chainID int64) {
	var req struct {
		NodeID int64 `json:"node_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBadRequest(w, "请求格式错误")
		return
	}
	if req.NodeID == 0 {
		writeBadRequest(w, "缺少 node_id")
		return
	}
	ctx := r.Context()
	old, err := h.repo.GetForwardChainBinding(ctx, req.NodeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if old == nil {
		writeBadRequest(w, "该节点未绑定转发链")
		return
	}
	if old.ChainID == chainID {
		respondJSON(w, http.StatusOK, map[string]any{"success": true, "message": "已在本链"})
		return
	}
	port := old.Port
	if err := h.repo.UnbindForwardChainNode(ctx, req.NodeID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	remain, _ := h.repo.ListForwardChainBindingsByChain(ctx, old.ChainID)
	if len(remain) == 0 {
		if err := h.syncChain(ctx, old.ChainID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		if h.remote != nil {
			if err := h.remote.TeardownExitNodes(ctx, old.ChainID); err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
		}
	}
	if err := h.repo.BindForwardChainNode(ctx, req.NodeID, chainID, port); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.syncChain(ctx, chainID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if h.remote != nil {
		if err := h.remote.ProvisionExitNodes(ctx, chainID, port); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true, "message": "已更换转发链"})
}

func pickChainPort(chain *storage.ForwardChain, requested int) (int, error) {
	if requested > 0 && requested <= 65535 {
		return requested, nil
	}
	start, end := 50520, 51314
	if chain != nil && chain.PortRangeStart > 0 && chain.PortRangeEnd >= chain.PortRangeStart {
		start, end = chain.PortRangeStart, chain.PortRangeEnd
	}
	if start == end {
		return start, nil
	}
	return start + rand.Intn(end-start+1), nil
}

func (h *ForwardHandler) firstEntryHost(ctx context.Context, chain *storage.ForwardChain) string {
	hosts := h.entryHosts(ctx, chain, false)
	if len(hosts) == 0 {
		return ""
	}
	return hosts[0]
}

func (h *ForwardHandler) entryHosts(ctx context.Context, chain *storage.ForwardChain, separate bool) []string {
	if chain != nil && !separate && strings.TrimSpace(chain.DNSDomain) != "" {
		return []string{strings.TrimSpace(chain.DNSDomain)}
	}
	if chain == nil || len(chain.Hops) == 0 {
		return []string{"127.0.0.1"}
	}
	g, err := h.repo.GetForwardGroup(ctx, chain.Hops[0].GroupID)
	if err != nil || g == nil {
		return []string{"127.0.0.1"}
	}
	if !separate && strings.TrimSpace(g.DNSDomain) != "" {
		return []string{strings.TrimSpace(g.DNSDomain)}
	}
	var hosts []string
	for _, m := range g.Members {
		host := fmt.Sprintf("server-%d", m.ServerID)
		if s, err := h.repo.GetRemoteServer(ctx, m.ServerID); err == nil && s != nil {
			if hst := forwardEntryHost(s); hst != "" {
				host = hst
			} else if s.Name != "" {
				host = s.Name
			}
		}
		hosts = append(hosts, host)
	}
	if len(hosts) == 0 {
		return []string{"127.0.0.1"}
	}
	if !separate {
		return hosts[:1]
	}
	return hosts
}

func forwardEntryHost(s *storage.RemoteServer) string {
	if s == nil {
		return ""
	}
	if d := strings.TrimSpace(s.Domain); d != "" {
		return d
	}
	if d := strings.TrimSpace(s.PullAddress); d != "" {
		return d
	}
	return forwardServerHost(s)
}

func (h *ForwardHandler) computeTerminus(ctx context.Context, chain *storage.ForwardChain, port int, pin int64) (string, int64) {
	if chain == nil || len(chain.Hops) == 0 {
		return "", 0
	}
	last := chain.Hops[len(chain.Hops)-1]
	g, err := h.repo.GetForwardGroup(ctx, last.GroupID)
	if err != nil || g == nil || len(g.Members) == 0 {
		if d := strings.TrimSpace(chain.DNSDomain); d != "" && port > 0 {
			return fmt.Sprintf("%s:%d", d, port), 0
		}
		return "", 0
	}
	member := g.Members[0]
	if pin > 0 {
		for _, m := range g.Members {
			if m.ServerID == pin {
				member = m
				break
			}
		}
	}
	host := ""
	if s, err := h.repo.GetRemoteServer(ctx, member.ServerID); err == nil && s != nil {
		host = forwardEntryHost(s)
		if host == "" {
			host = s.Name
		}
	}
	if host == "" {
		return "", member.ServerID
	}
	return fmt.Sprintf("%s:%d", host, port), member.ServerID
}

func (h *ForwardHandler) nthExitServerID(ctx context.Context, chain *storage.ForwardChain, i int) int64 {
	if chain == nil || len(chain.Hops) == 0 {
		return 0
	}
	g, err := h.repo.GetForwardGroup(ctx, chain.Hops[len(chain.Hops)-1].GroupID)
	if err != nil || g == nil || len(g.Members) == 0 {
		return 0
	}
	return g.Members[i%len(g.Members)].ServerID
}

func rewriteClashServerPort(clashJSON, host string, port int) (string, error) {
	var m map[string]any
	if strings.TrimSpace(clashJSON) != "" {
		if err := json.Unmarshal([]byte(clashJSON), &m); err != nil {
			m = map[string]any{}
		}
	} else {
		m = map[string]any{}
	}
	if host != "" {
		m["server"] = host
	}
	m["port"] = port
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func syntheticVLESSClash(name, host string, port int, uuidStr string) string {
	m := map[string]any{
		"name":    name,
		"type":    "vless",
		"server":  host,
		"port":    port,
		"uuid":    uuidStr,
		"network": "tcp",
		"tls":     false,
		"udp":     true,
	}
	b, _ := json.Marshal(m)
	return string(b)
}

// syncChain 触发一条链涉及服务器的转发规则同步(remote 未就绪则跳过,便于纯 CRUD 单测)。
func (h *ForwardHandler) syncChain(ctx context.Context, chainID int64) error {
	if h.remote == nil {
		return nil
	}
	return h.remote.SyncChainForwarding(ctx, chainID)
}

// getStatus 拉取所有承载转发规则的服务器 agent 的运行时状态(健康/RTT/字节),逐台捕获错误。
func (h *ForwardHandler) getStatus(w http.ResponseWriter, r *http.Request) {
	if h.remote == nil {
		respondJSON(w, http.StatusOK, map[string]any{"success": true, "servers": []any{}})
		return
	}
	ctx := r.Context()
	topo, err := h.remote.loadForwardTopology(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	servers := []map[string]any{}
	for _, sid := range forwardingServerIDs(topo) {
		entry := map[string]any{"server_id": sid}
		if s := topo.servers[sid]; s != nil {
			entry["name"] = s.Name
		}
		if rules, err := h.remote.FetchForwardStatus(ctx, sid); err != nil {
			entry["error"] = err.Error()
		} else {
			entry["rules"] = rules
		}
		servers = append(servers, entry)
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true, "servers": servers})
}

// getMetrics 读取某条规则在某台服务器上的历史采样点,供前端画延迟/健康曲线。
// query: server_id(必填)、rule_id(必填)、hours(可选,默认 24,钳到 [1,720])。
func (h *ForwardHandler) getMetrics(w http.ResponseWriter, r *http.Request) {
	serverID, err := strconv.ParseInt(r.URL.Query().Get("server_id"), 10, 64)
	if err != nil || serverID <= 0 {
		writeBadRequest(w, "非法 server_id")
		return
	}
	ruleID := strings.TrimSpace(r.URL.Query().Get("rule_id"))
	if ruleID == "" {
		writeBadRequest(w, "缺少 rule_id")
		return
	}
	hours := 24
	if v := r.URL.Query().Get("hours"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			hours = n
		}
	}
	if hours < 1 {
		hours = 1
	}
	if hours > 720 {
		hours = 720
	}
	since := time.Now().Add(-time.Duration(hours) * time.Hour)
	points, err := h.repo.GetForwardHopMetrics(r.Context(), serverID, ruleID, since)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"success": true, "points": points})
}
