package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"miaomiaowux/internal/storage"
)

// HandleWireGuardTransit POST /api/admin/remote/wireguard-transit?server_id=<源>
// body: {"peer_server_id": <对端>}
// 前端文案:「已创建 WireGuard 中转出站」—— 对端建 WG 入站,源端建 WG 出站。
func (h *RemoteManageHandler) HandleWireGuardTransit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		remoteWriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	srcID, err := strconv.ParseInt(r.URL.Query().Get("server_id"), 10, 64)
	if err != nil || srcID <= 0 {
		remoteWriteError(w, http.StatusBadRequest, "server_id required")
		return
	}
	var req struct {
		PeerServerID int64 `json:"peer_server_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PeerServerID <= 0 {
		remoteWriteError(w, http.StatusBadRequest, "peer_server_id required")
		return
	}
	if req.PeerServerID == srcID {
		remoteWriteError(w, http.StatusBadRequest, "对端不能是自己")
		return
	}
	src, err := h.repo.GetRemoteServer(r.Context(), srcID)
	if err != nil || src == nil {
		remoteWriteError(w, http.StatusNotFound, "源服务器不存在")
		return
	}
	dst, err := h.repo.GetRemoteServer(r.Context(), req.PeerServerID)
	if err != nil || dst == nil {
		remoteWriteError(w, http.StatusNotFound, "对端服务器不存在")
		return
	}
	if strings.EqualFold(strings.TrimSpace(dst.XrayMode), "external") {
		remoteWriteError(w, http.StatusBadRequest, "对端需要内嵌 xray 模式才能建 WireGuard 入站")
		return
	}

	tag := fmt.Sprintf("wg-transit-%d-%d", srcID, dst.ID)
	outTag := fmt.Sprintf("wg-to-%d", dst.ID)
	dev, err := h.ensureTransitDevice(r.Context(), dst, srcID, tag)
	if err != nil {
		remoteWriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	email := fmt.Sprintf("transit-s%d@wg", srcID)
	wg := NewWireGuardHandler(h.repo, h)
	lease, err := h.repo.GetActiveWGLease(r.Context(), dev.ID, email)
	if err != nil {
		remoteWriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if lease == nil {
		priv, pub, kerr := wgGenKeypair()
		if kerr != nil {
			remoteWriteError(w, http.StatusInternalServerError, "生成中转密钥失败")
			return
		}
		lease = &storage.WGLease{
			DeviceID: dev.ID, Username: fmt.Sprintf("transit-%d", srcID),
			Email: email, PrivateKey: priv, PublicKey: pub,
		}
		if _, err := h.repo.AllocateWGLease(r.Context(), lease); err != nil {
			remoteWriteError(w, http.StatusInternalServerError, fmt.Sprintf("分配中转地址失败: %v", err))
			return
		}
	}
	if err := wg.pushWGInbound(r.Context(), dev); err != nil {
		remoteWriteError(w, http.StatusBadGateway, "对端 WireGuard 入站下发失败: "+err.Error())
		return
	}

	endpointHost := remoteServerPublicHost(dst)
	if endpointHost == "" {
		remoteWriteError(w, http.StatusBadRequest, "对端没有可用的公网地址")
		return
	}
	addrs := []string{}
	if lease.IPv4 != "" {
		addrs = append(addrs, lease.IPv4+"/32")
	}
	if lease.IPv6 != "" {
		addrs = append(addrs, lease.IPv6+"/128")
	}
	outbound := map[string]any{
		"tag":      outTag,
		"protocol": "wireguard",
		"settings": map[string]any{
			"secretKey": lease.PrivateKey,
			"address":   addrs,
			"peers": []map[string]any{{
				"publicKey":  dev.ServerPublicKey,
				"endpoint":   net.JoinHostPort(endpointHost, fmt.Sprintf("%d", dev.ListenPort)),
				"keepAlive":  25,
				"allowedIPs": []string{"0.0.0.0/0", "::/0"},
			}},
			"mtu":         1420,
			"noKernelTun": true,
		},
	}
	body, _ := json.Marshal(map[string]any{"action": "add", "outbound": outbound})
	if _, err := h.ForwardToServer(r.Context(), srcID, http.MethodPost, "/api/child/outbounds", body); err != nil {
		remoteWriteError(w, http.StatusBadGateway, "源端 WireGuard 出站下发失败: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"message":      "已创建 WireGuard 中转出站",
		"outbound_tag": outTag,
		"inbound_tag":  tag,
		"device_id":    dev.ID,
	})
}

func (h *RemoteManageHandler) ensureTransitDevice(ctx context.Context, dst *storage.RemoteServer, srcID int64, tag string) (*storage.WGDevice, error) {
	if existing, err := h.repo.GetWGDevice(ctx, dst.ID, tag); err != nil {
		return nil, err
	} else if existing != nil {
		return existing, nil
	}
	priv, pub, err := wgGenKeypair()
	if err != nil {
		return nil, err
	}
	probePriv, probePub, err := wgGenKeypair()
	if err != nil {
		return nil, err
	}
	var lastErr error
	for i := 0; i < 16; i++ {
		octet := int((srcID*31+dst.ID*7+int64(i))%200 + 20)
		port := 51820 + int((srcID*13+dst.ID+int64(i))%800)
		dev := &storage.WGDevice{
			ServerID: dst.ID, InboundTag: tag, ListenPort: port,
			IPv4CIDR:   fmt.Sprintf("10.66.%d.0/24", octet),
			FirstIndex: 1, LastIndex: 254,
			ServerPrivateKey: priv, ServerPublicKey: pub,
			ProbePrivateKey: probePriv, ProbePublicKey: probePub,
		}
		id, err := h.repo.CreateWGDevice(ctx, dev)
		if err != nil {
			lastErr = err
			continue
		}
		dev.ID = id
		return dev, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("无法分配中转地址池")
	}
	return nil, lastErr
}
