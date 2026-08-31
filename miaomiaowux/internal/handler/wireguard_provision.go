package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"miaomiaowux/internal/storage"
)

// ensureWGDeviceFromInbound 在创建/更新 protocol=wireguard 入站后落/更新 wg_devices。
// 密钥优先复用 inbound.settings.secretKey,地址池从 settings.address 推断。
func ensureWGDeviceFromInbound(ctx context.Context, repo *storage.TrafficRepository, serverID int64, inbound map[string]any) (*storage.WGDevice, error) {
	if repo == nil || inbound == nil {
		return nil, nil
	}
	tag, _ := inbound["tag"].(string)
	if strings.TrimSpace(tag) == "" {
		return nil, fmt.Errorf("wg: inbound tag 为空")
	}
	port := inboundPort(inbound)
	settings, _ := inbound["settings"].(map[string]any)
	if settings == nil {
		settings = map[string]any{}
	}
	cidr := strings.TrimSpace(fmt.Sprint(settings["ipv4_cidr"]))
	if cidr == "" || cidr == "<nil>" {
		cidr = inferWGPoolCIDR(firstWGAddress(settings))
	}
	if cidr == "" {
		octet := byte(serverID%200 + 20)
		cidr = fmt.Sprintf("10.66.%d.0/24", octet)
	}
	ipv6 := ""
	if v, _ := settings["ipv6_cidr"].(string); strings.TrimSpace(v) != "" {
		ipv6 = strings.TrimSpace(v)
	}

	existing, err := repo.GetWGDevice(ctx, serverID, tag)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if port > 0 && existing.ListenPort != port {
			_ = repo.UpdateWGDeviceListenPort(ctx, existing.ID, port)
			existing.ListenPort = port
		}
		return existing, nil
	}

	priv, _ := settings["secretKey"].(string)
	pub := ""
	if strings.TrimSpace(priv) != "" {
		if p, err := wgPubFromPrivB64(priv); err == nil {
			pub = p
		}
	}
	if priv == "" || pub == "" {
		var err error
		priv, pub, err = wgGenKeypair()
		if err != nil {
			return nil, err
		}
	}
	probePriv, probePub, err := wgGenKeypair()
	if err != nil {
		return nil, err
	}
	_, ipnet, err := net.ParseCIDR(cidr)
	last := 254
	if err == nil {
		ones, bits := ipnet.Mask.Size()
		hostBits := bits - ones
		if hostBits > 1 && hostBits < 24 {
			last = (1 << hostBits) - 2
		}
	}
	if port <= 0 {
		port = 51820
	}
	dev := &storage.WGDevice{
		ServerID: serverID, InboundTag: tag, ListenPort: port,
		IPv4CIDR: cidr, IPv6CIDR: ipv6,
		FirstIndex: 1, LastIndex: last,
		ServerPrivateKey: priv, ServerPublicKey: pub,
		ProbePrivateKey: probePriv, ProbePublicKey: probePub,
	}
	id, err := repo.CreateWGDevice(ctx, dev)
	if err != nil {
		return nil, err
	}
	dev.ID = id
	return dev, nil
}

func firstWGAddress(settings map[string]any) string {
	switch v := settings["address"].(type) {
	case []any:
		if len(v) > 0 {
			return fmt.Sprint(v[0])
		}
	case []string:
		if len(v) > 0 {
			return v[0]
		}
	case string:
		return v
	}
	return ""
}

func inboundPort(inbound map[string]any) int {
	switch v := inbound["port"].(type) {
	case float64:
		return int(v)
	case int:
		return v
	case json.Number:
		n, _ := v.Int64()
		return int(n)
	}
	return 0
}

func provisionWGLeasesForNodes(ctx context.Context, repo *storage.TrafficRepository, remote *RemoteManageHandler, user storage.User, nodeIDs []int64) {
	if repo == nil || user.Username == "" {
		return
	}
	wg := NewWireGuardHandler(repo, remote)
	for _, nid := range nodeIDs {
		node, err := repo.GetNodeByID(ctx, nid)
		if err != nil || !isWireGuardProtocol(node.Protocol) || strings.TrimSpace(node.InboundTag) == "" {
			continue
		}
		server, _ := resolveNodeServer(ctx, repo, remote, node)
		if server == nil {
			continue
		}
		dev, err := repo.GetWGDevice(ctx, server.ID, node.InboundTag)
		if err != nil || dev == nil {
			continue
		}
		email := credentialEmailForUser(user, node.InboundTag)
		if existing, _ := repo.GetActiveWGLease(ctx, dev.ID, email); existing != nil {
			continue
		}
		priv, pub, err := wgGenKeypair()
		if err != nil {
			log.Printf("[wg] 生成 peer 密钥失败 user=%s node=%s: %v", user.Username, node.NodeName, err)
			continue
		}
		lease := &storage.WGLease{
			DeviceID: dev.ID, Username: user.Username, Email: email,
			AssignmentID: user.PackageAssignmentID, PrivateKey: priv, PublicKey: pub,
		}
		if _, err := repo.AllocateWGLease(ctx, lease); err != nil {
			log.Printf("[wg] 分配 lease 失败 user=%s node=%s: %v", user.Username, node.NodeName, err)
			continue
		}
		if err := wg.pushWGInbound(ctx, dev); err != nil {
			log.Printf("[wg] 下发入站失败 server=%d tag=%s: %v", dev.ServerID, dev.InboundTag, err)
		}
	}
}

func releaseWGLeasesForNodes(ctx context.Context, repo *storage.TrafficRepository, remote *RemoteManageHandler, user storage.User, nodeIDs []int64) {
	if repo == nil {
		return
	}
	wg := NewWireGuardHandler(repo, remote)
	touched := map[int64]*storage.WGDevice{}
	for _, nid := range nodeIDs {
		node, err := repo.GetNodeByID(ctx, nid)
		if err != nil || !isWireGuardProtocol(node.Protocol) || strings.TrimSpace(node.InboundTag) == "" {
			continue
		}
		server, _ := resolveNodeServer(ctx, repo, remote, node)
		if server == nil {
			continue
		}
		dev, err := repo.GetWGDevice(ctx, server.ID, node.InboundTag)
		if err != nil || dev == nil {
			continue
		}
		email := credentialEmailForUser(user, node.InboundTag)
		if err := repo.ReleaseWGLease(ctx, dev.ID, email); err != nil {
			log.Printf("[wg] 回收 lease 失败 email=%s: %v", email, err)
			continue
		}
		touched[dev.ID] = dev
	}
	for _, dev := range touched {
		if err := wg.pushWGInbound(ctx, dev); err != nil {
			log.Printf("[wg] 回收后下发失败 server=%d tag=%s: %v", dev.ServerID, dev.InboundTag, err)
		}
	}
}

func releaseWGLeasesByAssignment(ctx context.Context, repo *storage.TrafficRepository, remote *RemoteManageHandler, assignmentID int64) {
	if repo == nil || assignmentID <= 0 {
		return
	}
	leases, err := repo.ListActiveLeasesByAssignment(ctx, assignmentID)
	if err != nil {
		log.Printf("[wg] 列 assignment=%d leases 失败: %v", assignmentID, err)
		return
	}
	wg := NewWireGuardHandler(repo, remote)
	touched := map[int64]*storage.WGDevice{}
	for _, l := range leases {
		if err := repo.ReleaseWGLease(ctx, l.DeviceID, l.Email); err != nil {
			continue
		}
		if _, ok := touched[l.DeviceID]; !ok {
			if dev, e := repo.GetWGDeviceByID(ctx, l.DeviceID); e == nil && dev != nil {
				touched[l.DeviceID] = dev
			}
		}
	}
	for _, dev := range touched {
		_ = wg.pushWGInbound(ctx, dev)
	}
}

func syncWGLeasesForPackageNodes(ctx context.Context, repo *storage.TrafficRepository, remote *RemoteManageHandler, users []storage.User, added, removed []int64) {
	for _, user := range users {
		if len(removed) > 0 {
			releaseWGLeasesForNodes(ctx, repo, remote, user, removed)
		}
		if len(added) > 0 {
			provisionWGLeasesForNodes(ctx, repo, remote, user, added)
		}
	}
}

func removeWGDeviceForInbound(ctx context.Context, repo *storage.TrafficRepository, serverID int64, tag string) {
	if repo == nil || tag == "" {
		return
	}
	dev, err := repo.GetWGDevice(ctx, serverID, tag)
	if err != nil || dev == nil {
		return
	}
	if err := repo.DeleteWGDevice(ctx, dev.ID); err != nil {
		log.Printf("[wg] 删除 device server=%d tag=%s 失败: %v", serverID, tag, err)
	}
}
