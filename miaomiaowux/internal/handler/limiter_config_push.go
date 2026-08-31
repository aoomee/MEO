package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"miaomiaowux/internal/license"
	"miaomiaowux/internal/storage"
	"miaomiaowux/internal/version"
)

type LimiterConfigPusher struct {
	repo           *storage.TrafficRepository
	wsHandler      *RemoteWSHandler
	httpClient     *http.Client
	licenseManager *license.Manager
	pending        sync.Map // serverID -> *time.Timer
}

func NewLimiterConfigPusher(repo *storage.TrafficRepository, wsHandler *RemoteWSHandler) *LimiterConfigPusher {
	return &LimiterConfigPusher{
		repo:      repo,
		wsHandler: wsHandler,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// SetLicenseManager 注入 license 管理器,启用 limiter feature 运行时再校验。
// 未注入时(开发场景)不做散布校验。
func (p *LimiterConfigPusher) SetLicenseManager(mgr *license.Manager) {
	p.licenseManager = mgr
}

// resolveLimit 按 4 段优先级算用户在指定节点上的限速 + 客户端数:
//
//	user.NodeSpeedLimitOverrides[node_id]  ← 用户级 per-node(map 含 key 即生效)
//	  ?? user.SpeedLimitOverride           ← 用户级全局
//	  ?? pkg.NodeSpeedLimits[node_id]      ← 套餐级 per-node(含 routed→父 一次跳)
//	  ?? pkg.SpeedLimitMbps                ← 套餐通用
//	  ?? 0 (unlimited)
//
// 每一层用 "map 是否含 key" / "指针是否非 nil" 判断;**不能用 value > 0 判断**,
// 因为 0 是显式不限速的有意义值。客户端数同结构。
// nodeID = 0 时跳过 per-node 层,只用全局/通用层(常见于反查未命中)。
// connGroupKey 连接数计数分组键 = "<username>|<物理节点ID>"。同一用户在同一物理节点(含其路由出站
// 子账户)的所有 email 共用此 key → agent 侧共享一份并发连接配额(问题1:20 而非 20×N)。
func connGroupKey(username string, physicalNodeID int64) string {
	return fmt.Sprintf("%s|%d", username, physicalNodeID)
}

func assignmentConnGroupKey(username string, assignmentID, physicalNodeID int64) string {
	if assignmentID <= 0 {
		return connGroupKey(username, physicalNodeID)
	}
	return fmt.Sprintf("%s|pkg:%d|%d", username, assignmentID, physicalNodeID)
}

func routedRefForSubaccount(sa storage.ActiveSubaccountForLimiter, byID map[int64]storage.InboundNodeRef) storage.InboundNodeRef {
	ref := byID[sa.RoutedNodeID]
	if ref.NodeID == 0 {
		// Keep statistics isolated even for temporarily inconsistent historical
		// data. Falling back by inbound tag is unsafe because routed nodes share it.
		ref.NodeID = sa.RoutedNodeID
	}
	return ref
}

func resolveLimit(user *storage.User, pkg *storage.Package, nodeID, parentID int64) (speedMbps float64, deviceLimit int) {
	// 限速
	switch {
	case user != nil && nodeID > 0:
		if v, ok := user.NodeSpeedLimitOverrides[nodeID]; ok {
			speedMbps = v
			break
		}
		if parentID > 0 {
			if v, ok := user.NodeSpeedLimitOverrides[parentID]; ok {
				speedMbps = v
				break
			}
		}
		if user.SpeedLimitOverride != nil {
			speedMbps = *user.SpeedLimitOverride
			break
		}
		if pkg != nil {
			if v, ok := pkg.SpeedLimitMbpsForNode(nodeID, &parentID); ok {
				speedMbps = v
				break
			}
			speedMbps = pkg.SpeedLimitMbps
		}
	case user != nil:
		if user.SpeedLimitOverride != nil {
			speedMbps = *user.SpeedLimitOverride
		} else if pkg != nil {
			speedMbps = pkg.SpeedLimitMbps
		}
	}

	// 客户端数
	switch {
	case user != nil && nodeID > 0:
		if v, ok := user.NodeDeviceLimitOverrides[nodeID]; ok {
			deviceLimit = v
			break
		}
		if parentID > 0 {
			if v, ok := user.NodeDeviceLimitOverrides[parentID]; ok {
				deviceLimit = v
				break
			}
		}
		if user.DeviceLimitOverride != nil {
			deviceLimit = *user.DeviceLimitOverride
			break
		}
		if pkg != nil {
			if v, ok := pkg.DeviceLimitForNode(nodeID, &parentID); ok {
				deviceLimit = v
				break
			}
			deviceLimit = pkg.DeviceLimit
		}
	case user != nil:
		if user.DeviceLimitOverride != nil {
			deviceLimit = *user.DeviceLimitOverride
		} else if pkg != nil {
			deviceLimit = pkg.DeviceLimit
		}
	}

	return
}

func (p *LimiterConfigPusher) BuildLimiterConfigForServer(ctx context.Context, serverID int64) ([]WSLimiterConfigPayload, error) {
	legacyConfigs, err := p.repo.GetUserInboundConfigsByServer(ctx, serverID)
	if err != nil {
		return nil, err
	}
	assignmentConfigs, err := p.repo.ListPackageAssignmentInboundConfigsByServer(ctx, serverID)
	if err != nil {
		return nil, err
	}

	// 查 server name,用于反查子账号(子账号通过 routed_node 的 original_server 关联)
	var serverName string
	if servers, err := p.repo.ListRemoteServers(ctx); err == nil {
		for _, s := range servers {
			if s.ID == serverID {
				serverName = s.Name
				break
			}
		}
	}

	// routed 节点的 active 子账号:也要为它们下发限速规则,key 是子账号 email
	var subaccs []storage.ActiveSubaccountForLimiter
	if serverName != "" {
		subaccs, _ = p.repo.ListActiveSubaccountsByServerName(ctx, serverName)
		assignmentSubaccounts, _ := p.repo.ListPackageAssignmentSubaccountsByServerName(ctx, serverName)
		subaccs = append(subaccs, assignmentSubaccounts...)
	}

	if len(legacyConfigs) == 0 && len(assignmentConfigs) == 0 && len(subaccs) == 0 {
		return nil, nil
	}

	// 主账号可按 inbound_tag 唯一反查物理节点；路由子账号必须按
	// user_subaccounts.routed_node_id 精确反查。同一物理入站下可以有多个 routed
	// 节点共享 inbound_tag，若用 tag 做单值 map，所有连接会错误归到最后一个节点。
	physicalByTag := make(map[string]storage.InboundNodeRef)
	routedByID := make(map[int64]storage.InboundNodeRef)
	if serverName != "" {
		if refs, err := p.repo.ListInboundNodeRefsForServer(ctx, serverName); err == nil {
			for _, r := range refs {
				if r.NodeType == "routed" {
					routedByID[r.NodeID] = r
				} else {
					physicalByTag[r.InboundTag] = r
				}
			}
		}
	}

	// A limiter identity is (username, assignment). Legacy credentials use
	// assignment=0, while each additional package has its own package limits
	// and connection group even when both packages contain the same node.
	key := func(username string, assignmentID int64) string {
		return fmt.Sprintf("%s\x00%d", username, assignmentID)
	}
	userMap := make(map[string]*storage.User)
	pkgCache := make(map[int64]*storage.Package)
	baseUsers := make(map[string]storage.User)
	loadBase := func(username string) (storage.User, bool) {
		if user, ok := baseUsers[username]; ok {
			return user, true
		}
		user, err := p.repo.GetUser(ctx, username)
		if err != nil {
			return storage.User{}, false
		}
		baseUsers[username] = user
		return user, true
	}
	loadLegacy := func(username string) {
		if _, ok := userMap[key(username, 0)]; ok {
			return
		}
		user, ok := loadBase(username)
		if !ok {
			return
		}
		if !user.IsActive {
			return
		}
		if overLimit, _ := p.repo.IsUserOverLimit(ctx, username); overLimit {
			return
		}
		u := user
		userMap[key(username, 0)] = &u
		if user.PackageID > 0 {
			if _, ok := pkgCache[user.PackageID]; !ok {
				if pkg, err := p.repo.GetPackage(ctx, user.PackageID); err == nil {
					pkgCache[user.PackageID] = pkg
				}
			}
		}
	}
	for _, config := range legacyConfigs {
		loadLegacy(config.Username)
	}
	for _, sa := range subaccs {
		if sa.AssignmentID == 0 {
			loadLegacy(sa.Username)
		}
	}
	assignments, err := p.repo.ListActivePackageAssignments(ctx)
	if err != nil {
		return nil, err
	}
	for _, assignment := range assignments {
		if assignment.Legacy || assignment.OverLimitEnforced {
			continue
		}
		base, ok := loadBase(assignment.Username)
		if !ok || !base.IsActive {
			continue
		}
		u := userForPackageAssignment(base, assignment)
		userMap[key(assignment.Username, assignment.ID)] = &u
		if _, ok := pkgCache[assignment.PackageID]; !ok {
			if pkg, err := p.repo.GetPackage(ctx, assignment.PackageID); err == nil {
				pkgCache[assignment.PackageID] = pkg
			}
		}
	}

	tagUsers := make(map[string][]WSUserLimitInfo)
	tagPkgIDs := make(map[string]map[int64]bool)
	// 即使某 tag 当前没有可用用户，也必须发送 Users:[]。否则 Agent 会保留
	// 上一次的 limiter 用户和存量连接，禁用/超限看起来仍有连接数。
	for _, c := range legacyConfigs {
		tagUsers[c.InboundTag] = nil
	}
	for _, c := range assignmentConfigs {
		tagUsers[c.InboundTag] = nil
	}
	for _, sa := range subaccs {
		tagUsers[sa.InboundTag] = nil
	}

	// 主账号:走 c.InboundTag,反查 physical 节点的 (nodeID, parentID)
	appendPhysical := func(username string, assignmentID int64, inboundTag, email string) {
		user, ok := userMap[key(username, assignmentID)]
		if !ok {
			return
		}
		var pkg *storage.Package
		if user.PackageID > 0 {
			pkg = pkgCache[user.PackageID]
		}
		ref := physicalByTag[inboundTag]
		speedMbps, deviceLimit := resolveLimit(user, pkg, ref.NodeID, ref.ParentID)
		var speedBytes uint64
		if speedMbps > 0 {
			speedBytes = uint64(speedMbps * 1000000 / 8)
		}
		tagUsers[inboundTag] = append(tagUsers[inboundTag], WSUserLimitInfo{
			Email:         email,
			SpeedLimit:    speedBytes,
			DeviceLimit:   deviceLimit,
			ConnGroup:     assignmentConnGroupKey(user.Username, assignmentID, ref.NodeID),
			ConnStatGroup: assignmentConnGroupKey(user.Username, assignmentID, ref.NodeID),
		})
		if user.PackageID > 0 {
			if tagPkgIDs[inboundTag] == nil {
				tagPkgIDs[inboundTag] = make(map[int64]bool)
			}
			tagPkgIDs[inboundTag][user.PackageID] = true
		}
	}
	for _, c := range legacyConfigs {
		appendPhysical(c.Username, 0, c.InboundTag, credentialEmailForUser(storage.User{Username: c.Username}, c.InboundTag))
	}
	for _, c := range assignmentConfigs {
		appendPhysical(c.Username, c.AssignmentID, c.InboundTag, c.Email)
	}

	// 子账号:走 sa.InboundTag,反查 routed 节点的 (nodeID, parentID)。
	// routed 节点的 per-node 限速继承 parent 物理节点(在 resolveLimit 内自动处理)。
	for _, sa := range subaccs {
		if !sa.IsActive {
			continue
		}
		user, ok := userMap[key(sa.Username, sa.AssignmentID)]
		if !ok {
			continue
		}
		var pkg *storage.Package
		if user.PackageID > 0 {
			pkg = pkgCache[user.PackageID]
		}
		ref := routedRefForSubaccount(sa, routedByID)
		speedMbps, deviceLimit := resolveLimit(user, pkg, ref.NodeID, ref.ParentID)
		var speedBytes uint64
		if speedMbps > 0 {
			speedBytes = uint64(speedMbps * 1000000 / 8)
		}
		// 路由出站的 group 归到**父物理节点**(ref.ParentID),从而与父节点及其它路由出站共享连接配额。
		physID := ref.ParentID
		if physID == 0 {
			physID = ref.NodeID // 兜底:父未知时按自身,避免 group="user|0" 把不同节点误并
		}
		tagUsers[sa.InboundTag] = append(tagUsers[sa.InboundTag], WSUserLimitInfo{
			Email:         sa.Email,
			SpeedLimit:    speedBytes,
			DeviceLimit:   deviceLimit,
			ConnGroup:     assignmentConnGroupKey(user.Username, sa.AssignmentID, physID),
			ConnStatGroup: assignmentConnGroupKey(user.Username, sa.AssignmentID, ref.NodeID),
		})
		if user.PackageID > 0 {
			if tagPkgIDs[sa.InboundTag] == nil {
				tagPkgIDs[sa.InboundTag] = make(map[int64]bool)
			}
			tagPkgIDs[sa.InboundTag][user.PackageID] = true
		}
	}

	var payloads []WSLimiterConfigPayload
	for tag, users := range tagUsers {
		var rules []storage.AutoSpeedLimitRule
		for pkgID := range tagPkgIDs[tag] {
			if pkg, ok := pkgCache[pkgID]; ok && len(pkg.AutoSpeedRules) > 0 {
				rules = append(rules, pkg.AutoSpeedRules...)
			}
		}
		payloads = append(payloads, WSLimiterConfigPayload{
			InboundTag:     tag,
			Users:          users,
			AutoSpeedRules: rules,
		})
	}

	return payloads, nil
}

func (p *LimiterConfigPusher) PushToServer(ctx context.Context, serverID int64) {
	if p == nil {
		return
	}
	if old, ok := p.pending.Load(serverID); ok {
		if t, ok := old.(*time.Timer); ok {
			t.Stop()
		}
	}
	t := time.AfterFunc(1500*time.Millisecond, func() {
		p.pending.Delete(serverID)
		p.pushToServerNow(context.Background(), serverID)
	})
	p.pending.Store(serverID, t)
}

func (p *LimiterConfigPusher) pushToServerNow(ctx context.Context, serverID int64) {
	server, err := p.repo.GetRemoteServer(ctx, serverID)
	if err != nil {
		return
	}
	if server.XrayMode != "embedded" {
		return
	}

	// Connection identity/group metadata is operational telemetry, not a paid
	// limiter feature. Without entitlement we still push a zero-limit tracking
	// config so live connection counts never collapse to a misleading zero.
	enforce := p.licenseManager == nil || (p.licenseManager.HasFeature("limiter") && p.licenseManager.HasFeature("embedded"))

	configs, err := p.BuildLimiterConfigForServer(ctx, serverID)
	if err != nil {
		log.Printf("[LimiterPush] Failed to build config for server %d: %v", serverID, err)
		return
	}
	if len(configs) == 0 {
		return
	}
	if !enforce {
		for i := range configs {
			configs[i].TrackingOnly = true
			configs[i].NodeLimit = 0
			configs[i].AutoSpeedRules = nil
			for j := range configs[i].Users {
				configs[i].Users[j].SpeedLimit = 0
				configs[i].Users[j].DeviceLimit = 0
			}
		}
	}
	configs = prepareLimiterConfigs(serverID, configs)

	// WS-first:写成功直接 return;写失败 fallback HTTP(socket 撕连瞬间 GetConnection 还
	// 返回 ok 但 Send 已失败,这一次推送不能丢)。跟 deployToRemoteServer 路径对齐语义。
	if _, ok := p.wsHandler.GetConnectionByServerID(serverID); ok {
		if err := p.wsHandler.SendLimiterConfig(serverID, configs); err == nil {
			return
		} else {
			log.Printf("[LimiterPush] WebSocket send failed for server %d (%v), falling back to HTTP", serverID, err)
		}
	}

	p.pushViaHTTP(ctx, server, configs)
}

// RevokeAllEmbeddedServers 给所有 embedded 服务器下发**清零**的限速配置。
//
// 用于 limiter 特性丢失时(降级 / 到期 / 解绑)。这是 PushToAllEmbeddedServers 的反向操作,
// 补的是长期缺失的撤销路径:此前降级只是让 PushToServer 静默 return,而 agent 侧限速是
// 纯内存态(HandleLimiter → l.AddInboundLimiter),**不推新配置 ≠ 限速失效** ——
// PRO 时期推下去的那份会一直生效到 agent 进程重启。
//
// 刻意**不走** PushToServer:那里有 limiter/embedded 的 PRO gate,而本函数恰恰是在
// 「已经没有该特性」时调用的,走它会被自己拦掉,等于什么都没做。
func (p *LimiterConfigPusher) RevokeAllEmbeddedServers(ctx context.Context) {
	servers, err := p.repo.ListRemoteServers(ctx)
	if err != nil {
		log.Printf("[LimiterPush] revoke: ListRemoteServers failed: %v", err)
		return
	}
	for _, s := range servers {
		if s.XrayMode != "embedded" {
			continue
		}
		configs, berr := p.BuildLimiterConfigForServer(ctx, s.ID)
		if berr != nil {
			log.Printf("[LimiterPush] revoke: build config for server %d failed: %v", s.ID, berr)
			continue
		}
		// 保留 inbound_tag 与用户名单,只把限速值清零 —— agent 侧 AddInboundLimiter 是按
		// inbound_tag 整体替换的,必须把每个已下发过的 tag 都覆盖一遍才算真正撤销。
		zeroed := make([]WSLimiterConfigPayload, 0, len(configs))
		for _, c := range configs {
			users := make([]WSUserLimitInfo, len(c.Users))
			for i, u := range c.Users {
				// 保留 Email/ConnGroup(agent 按它们索引),只清限速与连接数上限
				users[i] = WSUserLimitInfo{Email: u.Email, ConnGroup: u.ConnGroup, ConnStatGroup: u.ConnStatGroup, SpeedLimit: 0, DeviceLimit: 0}
			}
			zeroed = append(zeroed, WSLimiterConfigPayload{
				InboundTag: c.InboundTag,
				NodeLimit:  0,
				Users:      users,
				// 自动限速规则一并清空,否则 SpeedMonitor 会继续按旧规则压限速
				AutoSpeedRules: nil,
			})
		}
		if len(zeroed) == 0 {
			continue
		}
		if _, ok := p.wsHandler.GetConnectionByServerID(s.ID); ok {
			if err := p.wsHandler.SendLimiterConfig(s.ID, zeroed); err == nil {
				log.Printf("[LimiterPush] revoked limiter on server %s (license feature lost)", s.Name)
				continue
			}
		}
		srv := s
		p.pushViaHTTP(ctx, &srv, zeroed)
		log.Printf("[LimiterPush] revoked limiter on server %s via HTTP (license feature lost)", s.Name)
	}
}

func (p *LimiterConfigPusher) pushViaHTTP(ctx context.Context, server *storage.RemoteServer, configs []WSLimiterConfigPayload) {
	hdr := http.Header{}
	hdr.Set("Content-Type", "application/json")
	hdr.Set("Authorization", "Bearer "+server.Token)
	hdr.Set("User-Agent", version.AgentUserAgent)

	for _, cfg := range configs {
		body, err := json.Marshal(cfg)
		if err != nil {
			log.Printf("[LimiterPush] Failed to marshal config for server %s: %v", server.Name, err)
			continue
		}
		// tryHTTPWithFallback 内部 v4-first → v6-fallback,消灭旧 strings.LastIndex IPv6 截断 bug
		resp, err := tryHTTPWithFallback(ctx, p.httpClient, server, http.MethodPost, "/api/child/limiter", body, hdr)
		if err != nil {
			log.Printf("[LimiterPush] HTTP push failed for server %s: %v", server.Name, err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("[LimiterPush] HTTP push returned %d for server %s", resp.StatusCode, server.Name)
		}
	}
}

func (p *LimiterConfigPusher) PushToAllServersForPackage(ctx context.Context, packageID int64) {
	assignments, err := p.repo.ListActivePackageAssignmentsByPackage(ctx, packageID)
	if err != nil {
		return
	}

	serverIDs := make(map[int64]bool)
	usernames := map[string]bool{}
	for _, assignment := range assignments {
		usernames[assignment.Username] = true
		configs, _ := p.repo.ListPackageAssignmentInboundConfigs(ctx, assignment.ID)
		for _, c := range configs {
			serverIDs[c.ServerID] = true
		}
	}
	for username := range usernames {
		if legacyConfigs, err := p.repo.GetUserInboundConfigs(ctx, username); err == nil {
			for _, config := range legacyConfigs {
				serverIDs[config.ServerID] = true
			}
		}
		// 同 PushToAllServersForUser:补上该用户 routed 子账号所在 server,避免 routed-only 用户漏推。
		if subIDs, err := p.repo.ListServerIDsForUserSubaccounts(ctx, username); err == nil {
			for _, id := range subIDs {
				serverIDs[id] = true
			}
		}
		if subIDs, err := p.repo.ListServerIDsForPackageAssignmentSubaccounts(ctx, username); err == nil {
			for _, id := range subIDs {
				serverIDs[id] = true
			}
		}
	}

	for sid := range serverIDs {
		p.PushToServer(ctx, sid)
	}
}

func (p *LimiterConfigPusher) PushToAllServersForUser(ctx context.Context, username string) {
	configs, err := p.repo.GetUserInboundConfigs(ctx, username)
	if err != nil {
		return
	}

	serverIDs := make(map[int64]bool)
	for _, c := range configs {
		serverIDs[c.ServerID] = true
	}
	if assignmentConfigs, err := p.repo.ListPackageAssignmentInboundConfigsByUser(ctx, username); err == nil {
		for _, c := range assignmentConfigs {
			serverIDs[c.ServerID] = true
		}
	}
	// 只有 routed 子账号、没有物理 inbound 的用户,上面查不到 server —— 补上子账号所在 server,
	// 否则这些用户在用户管理/套餐里设的限速对该 server 永不下发。
	if subIDs, err := p.repo.ListServerIDsForUserSubaccounts(ctx, username); err == nil {
		for _, id := range subIDs {
			serverIDs[id] = true
		}
	}
	if subIDs, err := p.repo.ListServerIDsForPackageAssignmentSubaccounts(ctx, username); err == nil {
		for _, id := range subIDs {
			serverIDs[id] = true
		}
	}

	for sid := range serverIDs {
		p.PushToServer(ctx, sid)
	}
}

// PushToAllEmbeddedServers 给所有 embedded 模式远程服务器重推限速配置。
// 用于 license 从失效恢复时补下发:失效期间 PushToServer 内部被 license gate 跳过(实际不限速),
// 恢复后主动补一遍,否则要等 agent 下次重连 auth 或用户改配置才恢复。
func (p *LimiterConfigPusher) PushToAllEmbeddedServers(ctx context.Context) {
	servers, err := p.repo.ListRemoteServers(ctx)
	if err != nil {
		return
	}
	for _, s := range servers {
		if s.XrayMode == "embedded" {
			p.PushToServer(ctx, s.ID)
		}
	}
}
