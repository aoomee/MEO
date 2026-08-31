package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"strings"

	"crypto/ecdh"

	"miaomiaowux/internal/storage"
)

func isWireGuardProtocol(protocol string) bool {
	return strings.EqualFold(strings.TrimSpace(protocol), "wireguard")
}

func wgPubFromPrivB64(privB64 string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(privB64))
	if err != nil || len(raw) != 32 {
		return "", fmt.Errorf("invalid wg private key")
	}
	key, err := ecdh.X25519().NewPrivateKey(raw)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key.PublicKey().Bytes()), nil
}

func wgDeviceEndpoint(ctx context.Context, repo *storage.TrafficRepository, dev *storage.WGDevice) string {
	if repo == nil || dev == nil {
		return ""
	}
	host := ""
	if srv, err := repo.GetRemoteServer(ctx, dev.ServerID); err == nil && srv != nil {
		host = remoteServerPublicHost(srv)
	}
	if host == "" {
		return ""
	}
	return net.JoinHostPort(host, fmt.Sprintf("%d", dev.ListenPort))
}

func remoteServerPublicHost(srv *storage.RemoteServer) string {
	if srv == nil {
		return ""
	}
	if h := strings.TrimSpace(srv.Domain); h != "" {
		return h
	}
	if h := strings.TrimSpace(srv.IPAddress); h != "" {
		return h
	}
	return strings.TrimSpace(srv.PullAddress)
}

func inferWGPoolCIDR(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return ""
	}
	if !strings.Contains(addr, "/") {
		if ip := net.ParseIP(addr); ip != nil && ip.To4() != nil {
			v := ip.To4()
			return fmt.Sprintf("%d.%d.%d.0/24", v[0], v[1], v[2])
		}
		return ""
	}
	ip, ipnet, err := net.ParseCIDR(addr)
	if err != nil {
		return ""
	}
	ones, bits := ipnet.Mask.Size()
	if ip.To4() != nil && (ones >= 32 || bits-ones <= 1) {
		v := ip.To4()
		return fmt.Sprintf("%d.%d.%d.0/24", v[0], v[1], v[2])
	}
	return ipnet.String()
}

func applyWGLeaseToProxy(ctx context.Context, repo *storage.TrafficRepository, proxy map[string]any, node storage.Node, username string) {
	if repo == nil || proxy == nil || !isWireGuardProtocol(node.Protocol) || strings.TrimSpace(username) == "" {
		return
	}
	leases, err := repo.ListActiveLeasesByUser(ctx, username)
	if err != nil {
		return
	}
	var match *storage.WGLease
	var matchDev *storage.WGDevice
	for i := range leases {
		dev, err := repo.GetWGDeviceByID(ctx, leases[i].DeviceID)
		if err != nil || dev == nil {
			continue
		}
		if node.InboundTag != "" && dev.InboundTag != node.InboundTag {
			continue
		}
		if name := strings.TrimSpace(node.OriginalServer); name != "" {
			if srv, serr := repo.GetRemoteServerByName(ctx, name); serr == nil && srv != nil && srv.ID != dev.ServerID {
				continue
			}
		}
		l := leases[i]
		match, matchDev = &l, dev
		break
	}
	if match == nil || matchDev == nil {
		return
	}
	proxy["type"] = "wireguard"
	proxy["udp"] = true
	proxy["private-key"] = match.PrivateKey
	proxy["public-key"] = matchDev.ServerPublicKey
	if match.IPv4 != "" {
		proxy["ip"] = match.IPv4 + "/32"
	}
	if match.IPv6 != "" {
		proxy["ipv6"] = match.IPv6 + "/128"
	}
	if host := ""; matchDev.ListenPort > 0 {
		if srv, err := repo.GetRemoteServer(ctx, matchDev.ServerID); err == nil && srv != nil {
			host = remoteServerPublicHost(srv)
		}
		if host != "" {
			proxy["server"] = host
			proxy["port"] = matchDev.ListenPort
		}
	}
	if _, ok := proxy["mtu"]; !ok {
		proxy["mtu"] = 1420
	}
	proxy["allowed-ips"] = []string{"0.0.0.0/0", "::/0"}
}

func clashProxyFromWGInbound(inbound map[string]interface{}, serverHost string, nodePort int, tag, serverName string) map[string]interface{} {
	settings, _ := inbound["settings"].(map[string]interface{})
	if settings == nil {
		settings = map[string]interface{}{}
	}
	pub := ""
	if sk, _ := settings["secretKey"].(string); sk != "" {
		if p, err := wgPubFromPrivB64(sk); err == nil {
			pub = p
		}
	}
	if p, _ := settings["publicKey"].(string); p != "" {
		pub = p
	}
	mtu := 1420
	switch v := settings["mtu"].(type) {
	case float64:
		if int(v) > 0 {
			mtu = int(v)
		}
	case int:
		if v > 0 {
			mtu = v
		}
	}
	name := fmt.Sprintf("[%s] %s", serverName, tag)
	return map[string]interface{}{
		"name":        name,
		"type":        "wireguard",
		"server":      serverHost,
		"port":        nodePort,
		"udp":         true,
		"public-key":  pub,
		"private-key": "",
		"mtu":         mtu,
	}
}
