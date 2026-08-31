package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"miaomiaowux/internal/storage"
)

// WireGuard 协议(v0.5.2 新增)——「地址池自动分配 / 订阅下发 / 解绑回收 / 热增删 peer / 健康探测」。
//
// 数据模型见 storage/wireguard.go(wg_devices 地址池 + wg_leases peer),与官方逐字段对齐。
// xray 入站配置格式取自 fork infra/conf/wireguard.go:
//   {protocol:"wireguard", settings:{secretKey, address:[server_ip/32...], peers:[{email,publicKey,allowedIPs:[peer_ip/32]}], mtu}}
// 一个 wg_device = 一台服务器上的一个 WG 入站(= 一个地址池);每个有效 wg_lease = 一个 peer。

// WireGuardHandler 管 WG 入站地址池 + peer 分配。
type WireGuardHandler struct {
	repo   *storage.TrafficRepository
	remote *RemoteManageHandler
}

func NewWireGuardHandler(repo *storage.TrafficRepository, remote *RemoteManageHandler) *WireGuardHandler {
	return &WireGuardHandler{repo: repo, remote: remote}
}

// ---------------- 密钥 ----------------

// wgGenKeypair 生成 WireGuard 密钥对(X25519,base64 标准编码)——复用 xray 的 genX25519Pair。
func wgGenKeypair() (priv, pub string, err error) {
	pk, pubk, err := genX25519Pair()
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(pk), base64.StdEncoding.EncodeToString(pubk), nil
}

// ---------------- xray 入站配置 ----------------

// buildWGInbound 把一个地址池 + 它当前有效的 peers 组装成 xray wireguard 入站配置。
// listenPort 是 WG 监听的 UDP 端口(从 inbound_tag 或调用方给)。
func buildWGInbound(dev *storage.WGDevice, leases []storage.WGLease, listenPort int, mtu int) map[string]any {
	if mtu == 0 {
		mtu = 1420
	}
	// 服务端地址:池里第 first_index 号(约定服务端占池首址)
	addrs := []string{}
	if ip, err := wgHostIP(dev.IPv4CIDR, dev.FirstIndex); err == nil {
		addrs = append(addrs, ip+"/32")
	}
	if dev.IPv6CIDR != "" {
		if ip, err := wgHostIP(dev.IPv6CIDR, dev.FirstIndex); err == nil {
			addrs = append(addrs, ip+"/128")
		}
	}
	peers := make([]map[string]any, 0, len(leases))
	for _, l := range leases {
		allowed := []string{}
		if l.IPv4 != "" {
			allowed = append(allowed, l.IPv4+"/32")
		}
		if l.IPv6 != "" {
			allowed = append(allowed, l.IPv6+"/128")
		}
		peers = append(peers, map[string]any{
			"email":      l.Email,
			"publicKey":  l.PublicKey,
			"allowedIPs": allowed,
		})
	}
	// 探测专用 peer:不占用户地址,只为面板发合法握手拿 RTT。
	if strings.TrimSpace(dev.ProbePublicKey) != "" {
		peers = append(peers, map[string]any{
			"email":      "probe@" + dev.InboundTag,
			"publicKey":  dev.ProbePublicKey,
			"allowedIPs": []string{"172.31.255.1/32"},
		})
	}
	return map[string]any{
		"tag":      dev.InboundTag,
		"protocol": "wireguard",
		"listen":   "0.0.0.0",
		"port":     listenPort,
		"settings": map[string]any{
			"secretKey":   dev.ServerPrivateKey,
			"address":     addrs,
			"peers":       peers,
			"mtu":         mtu,
			"noKernelTun": true,
			"kernelMode":  false,
		},
	}
}

// ---------------- 客户端订阅(标准 .conf) ----------------

// buildWGClientConf 给一个 peer 生成标准 WireGuard 客户端配置(通用 .conf,各家客户端都认)。
// endpoint 是「服务器公网地址:UDP端口」;serverPub 是地址池服务端公钥。
func buildWGClientConf(dev *storage.WGDevice, l *storage.WGLease, endpoint string) string {
	var b strings.Builder
	b.WriteString("[Interface]\n")
	b.WriteString("PrivateKey = " + l.PrivateKey + "\n")
	addr := []string{}
	if l.IPv4 != "" {
		addr = append(addr, l.IPv4+"/32")
	}
	if l.IPv6 != "" {
		addr = append(addr, l.IPv6+"/128")
	}
	b.WriteString("Address = " + strings.Join(addr, ", ") + "\n")
	b.WriteString("DNS = 1.1.1.1\n\n")
	b.WriteString("[Peer]\n")
	b.WriteString("PublicKey = " + dev.ServerPublicKey + "\n")
	b.WriteString("AllowedIPs = 0.0.0.0/0, ::/0\n")
	b.WriteString("Endpoint = " + endpoint + "\n")
	b.WriteString("PersistentKeepalive = 25\n")
	return b.String()
}

// ---------------- HTTP 路由 ----------------

func (h *WireGuardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// /api/admin/wireguard/devices            GET(list?server_id=) / POST(create)
	// /api/admin/wireguard/devices/{id}        DELETE
	// /api/admin/wireguard/devices/{id}/leases GET(list) / POST(assign)
	// /api/admin/wireguard/devices/{id}/leases/{email} DELETE(release)
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/wireguard/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	ctx := r.Context()

	if len(parts) >= 1 && parts[0] == "devices" {
		switch {
		case len(parts) == 1 && r.Method == http.MethodGet:
			h.listDevices(w, r)
			return
		case len(parts) == 1 && r.Method == http.MethodPost:
			h.createDevice(w, r)
			return
		case len(parts) == 2 && r.Method == http.MethodDelete:
			id, _ := strconv.ParseInt(parts[1], 10, 64)
			h.deleteDevice(w, r, id)
			return
		case len(parts) == 3 && parts[2] == "leases" && r.Method == http.MethodGet:
			id, _ := strconv.ParseInt(parts[1], 10, 64)
			h.listLeases(w, ctx, id)
			return
		case len(parts) == 3 && parts[2] == "leases" && r.Method == http.MethodPost:
			id, _ := strconv.ParseInt(parts[1], 10, 64)
			h.assignLease(w, r, id)
			return
		case len(parts) == 4 && parts[2] == "leases" && r.Method == http.MethodDelete:
			id, _ := strconv.ParseInt(parts[1], 10, 64)
			h.releaseLease(w, r, id, parts[3])
			return
		case len(parts) == 3 && parts[2] == "probe" && r.Method == http.MethodPost:
			id, _ := strconv.ParseInt(parts[1], 10, 64)
			h.probeDevice(w, r, id)
			return
		}
	}
	writeJSONError(w, http.StatusNotFound, "wireguard: 未知路由")
}

type createWGDeviceReq struct {
	ServerID   int64  `json:"server_id"`
	InboundTag string `json:"inbound_tag"`
	Port       int    `json:"port"`      // WG 监听 UDP 端口
	IPv4CIDR   string `json:"ipv4_cidr"` // 如 10.6.0.0/16
	IPv6CIDR   string `json:"ipv6_cidr"` // 可空
}

// pushWGInbound 把地址池当前 peers 组装成 xray wireguard 入站,热下发到 agent(add 幂等替换)。
// peer 增删都走它 → 官方说的「热增删 peer,绑定/解绑不中断在用隧道」。
func (h *WireGuardHandler) pushWGInbound(ctx context.Context, dev *storage.WGDevice) error {
	if h.remote == nil {
		return nil // 无下发通道(单测)时跳过
	}
	leases, err := h.repo.ListActiveLeasesByDevice(ctx, dev.ID)
	if err != nil {
		return err
	}
	inbound := buildWGInbound(dev, leases, dev.ListenPort, 0)
	body, _ := json.Marshal(map[string]any{"action": "add", "inbound": inbound})
	_, err = h.remote.ForwardToServer(ctx, dev.ServerID, http.MethodPost, "/api/child/inbounds", body)
	return err
}

func (h *WireGuardHandler) createDevice(w http.ResponseWriter, r *http.Request) {
	var req createWGDeviceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "请求体非法")
		return
	}
	if req.ServerID <= 0 || req.InboundTag == "" || req.IPv4CIDR == "" || req.Port <= 0 || req.Port > 65535 {
		writeJSONError(w, http.StatusBadRequest, "server_id / inbound_tag / ipv4_cidr / port 必填且合法")
		return
	}
	ip, ipnet, err := net.ParseCIDR(req.IPv4CIDR)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "ipv4_cidr 非法")
		return
	}
	_ = ip
	// 池区间:1..(2^hostbits - 2),0 号服务端占首址(见 first_index=1)。/16 → last≈65534。
	ones, bits := ipnet.Mask.Size()
	hostBits := bits - ones
	last := 2
	if hostBits > 0 && hostBits < 24 { // 上限保护,别把 /8 撑爆
		last = (1 << hostBits) - 2
	} else if hostBits >= 24 {
		last = (1 << 24) - 2
	}
	priv, pub, err := wgGenKeypair()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "生成服务端密钥失败")
		return
	}
	probePriv, probePub, err := wgGenKeypair()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "生成探测密钥失败")
		return
	}
	dev := &storage.WGDevice{
		ServerID: req.ServerID, InboundTag: req.InboundTag, ListenPort: req.Port,
		IPv4CIDR: req.IPv4CIDR, IPv6CIDR: req.IPv6CIDR,
		FirstIndex: 1, LastIndex: last,
		ServerPrivateKey: priv, ServerPublicKey: pub,
		ProbePrivateKey: probePriv, ProbePublicKey: probePub,
	}
	id, err := h.repo.CreateWGDevice(r.Context(), dev)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("创建 WG 入站失败: %v", err))
		return
	}
	dev.ID = id
	// 热下发空 peers 的 WG 入站到 agent(建池即起监听)
	pushErr := ""
	if err := h.pushWGInbound(r.Context(), dev); err != nil {
		pushErr = err.Error() // 池已建,下发失败不回滚,报给前端提示
	}
	respondJSON(w, http.StatusOK, map[string]any{"id": id, "server_public_key": pub, "device": dev, "push_error": pushErr})
}

func (h *WireGuardHandler) listDevices(w http.ResponseWriter, r *http.Request) {
	serverID, _ := strconv.ParseInt(r.URL.Query().Get("server_id"), 10, 64)
	if serverID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "server_id 必填")
		return
	}
	devs, err := h.repo.ListWGDevicesByServer(r.Context(), serverID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"devices": devs})
}

func (h *WireGuardHandler) deleteDevice(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.repo.DeleteWGDevice(r.Context(), id); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *WireGuardHandler) listLeases(w http.ResponseWriter, ctx context.Context, deviceID int64) {
	leases, err := h.repo.ListActiveLeasesByDevice(ctx, deviceID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"leases": leases})
}

type assignWGLeaseReq struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	AssignmentID int64  `json:"assignment_id"`
}

// assignLease 在池里给一个用户分配 peer(自动挑地址 + 生成 peer 密钥),返回客户端配置。
func (h *WireGuardHandler) assignLease(w http.ResponseWriter, r *http.Request, deviceID int64) {
	var req assignWGLeaseReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "请求体非法")
		return
	}
	if req.Username == "" || req.Email == "" {
		writeJSONError(w, http.StatusBadRequest, "username / email 必填")
		return
	}
	priv, pub, err := wgGenKeypair()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "生成 peer 密钥失败")
		return
	}
	lease := &storage.WGLease{
		DeviceID: deviceID, Username: req.Username, Email: req.Email,
		AssignmentID: req.AssignmentID, PrivateKey: priv, PublicKey: pub,
	}
	if existing, e := h.repo.GetActiveWGLease(r.Context(), deviceID, req.Email); e == nil && existing != nil {
		lease = existing
	} else if _, err := h.repo.AllocateWGLease(r.Context(), lease); err != nil {
		if err == storage.ErrWGPoolExhausted {
			writeJSONError(w, http.StatusConflict, "地址池已耗尽")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("分配 peer 失败: %v", err))
		return
	}
	// 热增 peer:重建入站下发(不中断在用隧道)
	pushErr := ""
	var clientConf string
	if dev, e := h.repo.GetWGDeviceByID(r.Context(), deviceID); e == nil && dev != nil {
		if e := h.pushWGInbound(r.Context(), dev); e != nil {
			pushErr = e.Error()
		}
		clientConf = buildWGClientConf(dev, lease, wgDeviceEndpoint(r.Context(), h.repo, dev))
	}
	respondJSON(w, http.StatusOK, map[string]any{"lease": lease, "client_conf": clientConf, "push_error": pushErr})
}

func (h *WireGuardHandler) releaseLease(w http.ResponseWriter, r *http.Request, deviceID int64, email string) {
	if err := h.repo.ReleaseWGLease(r.Context(), deviceID, email); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 热删 peer:重建入站下发
	pushErr := ""
	if dev, e := h.repo.GetWGDeviceByID(r.Context(), deviceID); e == nil && dev != nil {
		if e := h.pushWGInbound(r.Context(), dev); e != nil {
			pushErr = e.Error()
		}
	}
	respondJSON(w, http.StatusOK, map[string]any{"ok": true, "push_error": pushErr})
}

// wgHostIP 复用 storage 的 CIDR+序号→IP(这里再实现一份避免跨包导出;逻辑一致)。
func wgHostIP(cidr string, idx int) (string, error) {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	base := ip.To4()
	if base == nil {
		base = ip.To16()
	}
	b := make(net.IP, len(base))
	copy(b, base)
	carry := idx
	for i := len(b) - 1; i >= 0 && carry > 0; i-- {
		v := int(b[i]) + (carry & 0xff)
		b[i] = byte(v & 0xff)
		carry = (carry >> 8) + (v >> 8)
	}
	return b.String(), nil
}
