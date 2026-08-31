package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"miaomiaowux/internal/storage"
)

// 转发链编排:把「链 + 组 + 成员」算成每台服务器要承载的转发规则,
// 经现有 RemoteManage → WS-RPC 下发给 agent 的原生转发引擎(不走 xray tunnel)。
//
// 端口模型:一条链上的节点用同一个端口 P。出口组(最后一跳)每台跑真实 xray 节点(:P),
// 中间/入口组每台跑 agent relay(:P → 下一组成员 :P)。因此:
//   - 出口组服务器【没有】转发规则(由 xray 节点承担);
//   - 第 i 跳(i=0..k-2)的每台服务器,转发到第 i+1 组的全部成员;
//   - 该跳的策略/权重/健康阈值都取自【目标组】(i+1),因为它决定如何在目标成员间分发。

// 下发给 agent 的 wire 类型,字段与 mmw-agent internal/forward.Rule 的 JSON 严格对齐。
// agent 与 master 是不同模块,按既有惯例镜像结构(如 WSRPCCallPayload)。
type forwardRule struct {
	ID        string            `json:"id"`
	Listen    string            `json:"listen"`
	Protocol  string            `json:"protocol"`
	Upstreams []forwardUpstream `json:"upstreams"`
	Strategy  string            `json:"strategy"`
	Health    forwardHealth     `json:"health"`
}

type forwardUpstream struct {
	Addr   string `json:"addr"`
	Weight int    `json:"weight"`
}

type forwardHealth struct {
	Enabled        bool `json:"enabled"`
	IntervalMs     int  `json:"interval_ms"`
	TimeoutMs      int  `json:"timeout_ms"`
	FailoverMs     int  `json:"failover_ms"`
	RecoverMs      int  `json:"recover_ms"`
	RTTThresholdMs int  `json:"rtt_threshold_ms"`
}

// agent 上报的运行时状态,字段与 mmw-agent internal/forward.RuleStatus 的 JSON 对齐。
type forwardRuleStatus struct {
	RuleID    string                  `json:"rule_id"`
	Listen    string                  `json:"listen"`
	Upstreams []forwardUpstreamStatus `json:"upstreams"`
}

type forwardUpstreamStatus struct {
	Addr      string `json:"addr"`
	Healthy   bool   `json:"healthy"`
	RTTMs     int64  `json:"rtt_ms"`
	BytesUp   uint64 `json:"bytes_up"`
	BytesDown uint64 `json:"bytes_down"`
}

// forwardServerHost 取一台服务器用于被上游转发连接的地址(优先 IPv4,回退域名)。
// IPv6 择优留待后续(节点若走 v6 再扩展)。
func forwardServerHost(s *storage.RemoteServer) string {
	if s == nil {
		return ""
	}
	if s.IPAddress != "" {
		return s.IPAddress
	}
	return s.Domain
}

// mapForwardStrategy 把组的均衡策略映射到 agent picker 支持的策略。
// agent 认 round_robin / weighted / least_conn / sticky。
//   percentage、cycle 是「按权重」的两种,权重由 master 按流量/周期填入 Upstream.Weight,故都走 weighted;
//   weighted 是手工权重,同样走 weighted;
//   least_conn、sticky 直接透传(agent v0.6.8 起 picker 原生支持,旧 agent 会安全退回 round_robin)。
func mapForwardStrategy(groupStrategy string) string {
	switch groupStrategy {
	case "round_robin", "":
		return "round_robin"
	case "least_conn":
		return "least_conn"
	case "sticky":
		return "sticky"
	default: // percentage / cycle / weighted
		return "weighted"
	}
}

// buildForwardRules 纯函数:把一条链算成 serverID → []forwardRule。
// 入参 groups/servers 为按 id 索引的查找表,port 为节点端口。出口组不产生规则。
func buildForwardRules(
	chain *storage.ForwardChain,
	groups map[int64]*storage.ForwardGroup,
	servers map[int64]*storage.RemoteServer,
	port int,
	relayProtocol string,
) (map[int64][]forwardRule, error) {
	if chain == nil {
		return nil, fmt.Errorf("forward: 链为空")
	}
	if port <= 0 || port > 65535 {
		return nil, fmt.Errorf("forward: 非法端口 %d", port)
	}
	if len(chain.Hops) < 2 {
		return nil, fmt.Errorf("forward: 链至少需要入口和出口两组")
	}

	out := make(map[int64][]forwardRule)
	// 遍历相邻跳:src=chain.Hops[i] 组的每台 → dst=chain.Hops[i+1] 组的全部成员
	for i := range len(chain.Hops) - 1 {
		src := groups[chain.Hops[i].GroupID]
		dst := groups[chain.Hops[i+1].GroupID]
		if src == nil || dst == nil {
			return nil, fmt.Errorf("forward: 链 %d 第 %d 跳引用了不存在的组", chain.ID, i)
		}
		if len(dst.Members) == 0 {
			// 官方 v0.5.2：空目标组只跳过这一跳，不让整条链一条规则都不下发。
			continue
		}

		upstreams := make([]forwardUpstream, 0, len(dst.Members))
		for _, m := range dst.Members {
			host := forwardServerHost(servers[m.ServerID])
			if host == "" {
				continue
			}
			w := m.Weight
			if w <= 0 {
				w = 1
			}
			upstreams = append(upstreams, forwardUpstream{
				Addr:   fmt.Sprintf("%s:%d", host, port),
				Weight: w,
			})
		}
		if len(upstreams) == 0 {
			continue
		}

		health := forwardHealth{
			Enabled:        dst.FailoverEnabled,
			RTTThresholdMs: dst.OfflineMsThreshold,
			// 窗口留 0 → agent 侧 normalized() 兜默认,默认值单一来源在 agent。
		}
		ruleID := fmt.Sprintf("fwd-c%d-p%d-h%d", chain.ID, port, i)
		proto := strings.ToLower(strings.TrimSpace(relayProtocol))
		if proto == "" {
			proto = "tcp"
		}

		for _, sm := range src.Members {
			out[sm.ServerID] = append(out[sm.ServerID], forwardRule{
				ID:        ruleID,
				Listen:    fmt.Sprintf(":%d", port),
				Protocol:  proto,
				Upstreams: upstreams,
				Strategy:  mapForwardStrategy(dst.BalanceStrategy),
				Health:    health,
			})
		}
	}
	return out, nil
}

// pushForwardRules 把一台服务器的全量转发规则下发到其 agent(幂等)。
func (h *RemoteManageHandler) pushForwardRules(ctx context.Context, serverID int64, rules []forwardRule) error {
	body, err := json.Marshal(map[string]any{"rules": rules})
	if err != nil {
		return err
	}
	resp, err := h.ForwardToServer(ctx, serverID, http.MethodPost, "/api/child/forward/apply", body)
	if err != nil {
		return fmt.Errorf("下发转发规则到服务器 %d 失败: %w", serverID, err)
	}
	var r struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}
	_ = json.Unmarshal(resp, &r) // 解析失败则 Success 保持 false,按拒绝处理
	if !r.Success {
		return fmt.Errorf("服务器 %d 拒绝转发规则: %s", serverID, r.Error)
	}
	return nil
}

// forwardTopology 是下发聚合需要的查找表:全部绑定 + 相关链/组/服务器。
type forwardTopology struct {
	bindings []storage.ForwardChainBinding
	chains   map[int64]*storage.ForwardChain
	groups   map[int64]*storage.ForwardGroup
	servers  map[int64]*storage.RemoteServer
}

// aggregateServerForwardRules 纯函数:算出某台 server 在【所有绑定链】下应承载的全部转发规则。
// 关键:agent 的 Apply 是整机全量状态,所以必须跨链聚合,不能按单链下发(否则多链互相覆盖)。
// 某条链结构不完整(如成员缺地址)时跳过该链、继续其余,返回首个跳过原因供日志。
func aggregateServerForwardRules(serverID int64, topo forwardTopology) ([]forwardRule, error) {
	var out []forwardRule
	var firstErr error
	for _, b := range topo.bindings {
		chain := topo.chains[b.ChainID]
		if chain == nil {
			continue
		}
		rulesByServer, err := buildForwardRules(chain, topo.groups, topo.servers, b.Port, b.Protocol)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		out = append(out, rulesByServer[serverID]...)
	}
	return out, firstErr
}

// loadForwardTopology 一次性加载全部绑定及其引用的链/组/服务器。
func (h *RemoteManageHandler) loadForwardTopology(ctx context.Context) (forwardTopology, error) {
	topo := forwardTopology{
		chains:  map[int64]*storage.ForwardChain{},
		groups:  map[int64]*storage.ForwardGroup{},
		servers: map[int64]*storage.RemoteServer{},
	}
	bindings, err := h.repo.ListForwardChainBindings(ctx)
	if err != nil {
		return topo, err
	}
	topo.bindings = bindings
	for _, b := range bindings {
		if _, ok := topo.chains[b.ChainID]; ok {
			continue
		}
		chain, err := h.repo.GetForwardChain(ctx, b.ChainID)
		if err != nil {
			return topo, err
		}
		topo.chains[b.ChainID] = chain
		for _, hop := range chain.Hops {
			if _, ok := topo.groups[hop.GroupID]; ok {
				continue
			}
			g, err := h.repo.GetForwardGroup(ctx, hop.GroupID)
			if err != nil {
				return topo, err
			}
			topo.groups[g.ID] = g
			for _, m := range g.Members {
				if _, ok := topo.servers[m.ServerID]; ok {
					continue
				}
				s, err := h.repo.GetRemoteServer(ctx, m.ServerID)
				if err != nil {
					return topo, err
				}
				topo.servers[m.ServerID] = s
			}
		}
	}
	return topo, nil
}

// SyncServersForwarding 把给定服务器各自的【全量】转发规则重算并下发(幂等)。
// 传入被删除/被移出组的服务器也没问题:它们算出空规则集 → agent 清空对应 listener。
func (h *RemoteManageHandler) SyncServersForwarding(ctx context.Context, serverIDs []int64) error {
	if len(serverIDs) == 0 {
		return nil
	}
	topo, err := h.loadForwardTopology(ctx)
	if err != nil {
		return err
	}
	var firstErr error
	for _, sid := range serverIDs {
		rules, aggErr := aggregateServerForwardRules(sid, topo)
		if aggErr != nil {
			log.Printf("[forward] 聚合服务器 %d 规则时跳过了不完整的链: %v", sid, aggErr)
		}
		if err := h.pushForwardRules(ctx, sid, rules); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// SyncChainForwarding 一条链发生变化(建/删节点、改成员/跳/策略)后,
// 重新同步该链涉及的所有服务器。extraServers 用于把「被移出的服务器」也纳入清理。
func (h *RemoteManageHandler) SyncChainForwarding(ctx context.Context, chainID int64, extraServers ...int64) error {
	chain, err := h.repo.GetForwardChain(ctx, chainID)
	if err != nil {
		return err
	}
	idset := map[int64]struct{}{}
	for _, sid := range extraServers {
		idset[sid] = struct{}{}
	}
	for _, hop := range chain.Hops {
		g, err := h.repo.GetForwardGroup(ctx, hop.GroupID)
		if err != nil {
			return err
		}
		for _, m := range g.Members {
			idset[m.ServerID] = struct{}{}
		}
	}
	ids := make([]int64, 0, len(idset))
	for sid := range idset {
		ids = append(ids, sid)
	}
	return h.SyncServersForwarding(ctx, ids)
}

// FetchForwardStatus 拉取一台服务器 agent 的转发运行时状态。
func (h *RemoteManageHandler) FetchForwardStatus(ctx context.Context, serverID int64) ([]forwardRuleStatus, error) {
	resp, err := h.ForwardToServer(ctx, serverID, http.MethodGet, "/api/child/forward/status", nil)
	if err != nil {
		return nil, err
	}
	var r struct {
		Success bool                `json:"success"`
		Rules   []forwardRuleStatus `json:"rules"`
	}
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, fmt.Errorf("解析服务器 %d 转发状态失败: %w", serverID, err)
	}
	return r.Rules, nil
}

// forwardingServerIDs 纯函数:算出拓扑中【真正承载转发规则】的服务器(排除只做出口的)。
// 用它决定向哪些 agent 拉状态,避免给纯出口 server 发无意义请求。
func forwardingServerIDs(topo forwardTopology) []int64 {
	var ids []int64
	for sid := range topo.servers {
		if rules, _ := aggregateServerForwardRules(sid, topo); len(rules) > 0 {
			ids = append(ids, sid)
		}
	}
	return ids
}

type exitProvisionSpec struct {
	Protocol       string
	Settings       map[string]any
	StreamSettings map[string]any
	CertID         int64
}

// ProvisionExitNodes 在转发链绑定时自动在出口组每台服务器上创建落地 xray 入站节点。
func (h *RemoteManageHandler) ProvisionExitNodes(ctx context.Context, chainID int64, port int) error {
	return h.ProvisionExitNodesSpec(ctx, chainID, port, exitProvisionSpec{})
}

// ProvisionExitNodesSpec 官方 v0.5.2：出口可走任意协议 + 托管证书。
func (h *RemoteManageHandler) ProvisionExitNodesSpec(ctx context.Context, chainID int64, port int, spec exitProvisionSpec) error {
	chain, err := h.repo.GetForwardChain(ctx, chainID)
	if err != nil {
		return fmt.Errorf("读取链失败: %w", err)
	}
	if len(chain.Hops) == 0 {
		return fmt.Errorf("链 %d 没有跳", chainID)
	}

	// 出口组 = 最后一跳
	exitGroupID := chain.Hops[len(chain.Hops)-1].GroupID
	exitGroup, err := h.repo.GetForwardGroup(ctx, exitGroupID)
	if err != nil {
		return fmt.Errorf("读取出口组失败: %w", err)
	}

	// 回滚栈
	type provision struct {
		serverID   int64
		inboundTag string
		nodeID     int64
	}
	var done []provision
	rollback := func() {
		for _, p := range done {
			// 删除入站
			body, _ := json.Marshal(map[string]any{"action": "remove", "tag": p.inboundTag})
			rctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			_, _ = h.ForwardToServer(rctx, p.serverID, http.MethodPost, "/api/child/inbounds", body)
			cancel()
			// 删除节点行
			if p.nodeID > 0 {
				_ = h.repo.DeleteNode(context.Background(), p.nodeID, "admin")
			}
		}
	}

	// 逐台出口服务器 provision
	for _, member := range exitGroup.Members {
		serverID := member.ServerID
		server, err := h.repo.GetRemoteServer(ctx, serverID)
		if err != nil {
			rollback()
			return fmt.Errorf("读取服务器 %d 失败: %w", serverID, err)
		}

		inboundTag := fmt.Sprintf("fwd-exit-c%d-s%d", chainID, serverID)
		nodeName := fmt.Sprintf("[%s] 转发链出口-%d", server.Name, chainID)
		if existing := h.findForwardExitNode(ctx, inboundTag); existing != nil {
			continue
		}

		uuidStr := uuid.New().String()
		proto := strings.ToLower(strings.TrimSpace(spec.Protocol))
		if proto == "" {
			proto = "vless"
		}
		settings := spec.Settings
		if settings == nil {
			settings = map[string]any{
				"clients": []map[string]any{
					{"id": uuidStr, "email": fmt.Sprintf("fwd-exit-%d@system", chainID)},
				},
				"decryption": "none",
			}
		}
		stream := spec.StreamSettings
		if stream == nil {
			stream = map[string]any{"network": "tcp"}
		}
		if spec.CertID > 0 {
			stream["security"] = "tls"
			stream["tlsSettings"] = map[string]any{"cert_id": spec.CertID}
		}
		inbound := map[string]any{
			"tag":            inboundTag,
			"protocol":       proto,
			"port":           port,
			"settings":       settings,
			"streamSettings": stream,
		}

		// 下发入站
		addBody, _ := json.Marshal(map[string]any{"action": "add", "inbound": inbound})
		if _, err := h.ForwardToServer(ctx, serverID, http.MethodPost, "/api/child/inbounds", addBody); err != nil {
			rollback()
			return fmt.Errorf("服务器 %s 添加入站失败: %w", server.Name, err)
		}

		// 生成 clash_config
		serverHost := forwardServerHost(server)
		clashProxy, err := h.inboundToClashProxy(inbound, serverHost, server.Name, 0)
		if err != nil {
			rollback()
			return fmt.Errorf("服务器 %s 生成 clash_config 失败: %w", server.Name, err)
		}
		clashBytes, _ := json.Marshal(clashProxy)
		parsedBytes, _ := json.Marshal(inbound)

		// 落节点行
		node := storage.Node{
			Username:       "_system_",
			NodeName:       nodeName,
			Protocol:       proto,
			ParsedConfig:   string(parsedBytes),
			ClashConfig:    string(clashBytes),
			Enabled:        true,
			Tag:            "转发链出口",
			OriginalServer: server.Name,
			InboundTag:     inboundTag,
			NodeType:       "forward-exit",
		}
		createdNode, err := h.repo.CreateNode(ctx, node)
		if err != nil {
			rollback()
			return fmt.Errorf("服务器 %s 创建节点行失败: %w", server.Name, err)
		}

		done = append(done, provision{serverID: serverID, inboundTag: inboundTag, nodeID: createdNode.ID})
		log.Printf("[ForwardProvision] 链 %d 出口服务器 %s 节点已创建 (node_id=%d tag=%s port=%d)",
			chainID, server.Name, createdNode.ID, inboundTag, port)
	}

	return nil
}

// TeardownExitNodes 在转发链解绑或删除节点时清理出口服务器上的自动建站。
// 按 node_type=forward-exit 标记找到自动建的节点,删除 agent 入站和节点行。
func (h *RemoteManageHandler) TeardownExitNodes(ctx context.Context, chainID int64) error {
	// 查找所有系统节点（挂在 admin 下）
	nodes, err := h.repo.ListNodes(ctx, "admin")
	if err != nil {
		return fmt.Errorf("查找节点失败: %w", err)
	}

	tagPrefix := fmt.Sprintf("fwd-exit-c%d-", chainID)
	for _, node := range nodes {
		if !strings.HasPrefix(node.InboundTag, tagPrefix) {
			continue
		}

		// 从 OriginalServer 查找 server_id
		server, err := h.repo.GetRemoteServerByName(ctx, node.OriginalServer)
		if err != nil {
			log.Printf("[ForwardTeardown] 查找服务器 %s 失败: %v", node.OriginalServer, err)
			continue
		}

		// 删除 agent 入站
		body, _ := json.Marshal(map[string]any{"action": "remove", "tag": node.InboundTag})
		if _, err := h.ForwardToServer(ctx, server.ID, http.MethodPost, "/api/child/inbounds", body); err != nil {
			log.Printf("[ForwardTeardown] 删除服务器 %s 入站 %s 失败: %v", node.OriginalServer, node.InboundTag, err)
		}

		if err := h.repo.DeleteNodeByID(ctx, node.ID); err != nil {
			log.Printf("[ForwardTeardown] 删除节点 %d 失败: %v", node.ID, err)
		} else {
			log.Printf("[ForwardTeardown] 链 %d 出口节点 %d (%s) 已清理", chainID, node.ID, node.NodeName)
		}
	}

	return nil
}

func (h *RemoteManageHandler) TeardownExitNodeOnServer(ctx context.Context, chainID, serverID int64) error {
	tag := fmt.Sprintf("fwd-exit-c%d-s%d", chainID, serverID)
	body, _ := json.Marshal(map[string]any{"action": "remove", "tag": tag})
	if _, err := h.ForwardToServer(ctx, serverID, http.MethodPost, "/api/child/inbounds", body); err != nil {
		log.Printf("[ForwardTeardown] 删除服务器 %d 入站 %s 失败: %v", serverID, tag, err)
	}
	if node := h.findForwardExitNode(ctx, tag); node != nil {
		_ = h.repo.DeleteNodeByID(ctx, node.ID)
	}
	return nil
}

func (h *RemoteManageHandler) findForwardExitNode(ctx context.Context, inboundTag string) *storage.Node {
	for _, owner := range []string{"_system_", "admin"} {
		nodes, err := h.repo.ListNodes(ctx, owner)
		if err != nil {
			continue
		}
		for i := range nodes {
			if nodes[i].InboundTag == inboundTag {
				return &nodes[i]
			}
		}
	}
	return nil
}
