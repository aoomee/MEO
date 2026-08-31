package handler

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

const localRealityNginxDest = "127.0.0.1:8001"

func normalizeRealityHost(value string) string {
	return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(value)), ".")
}

func realityPublicFallbackHost(inboundReq map[string]interface{}) (string, map[string]interface{}) {
	action := strings.ToLower(strings.TrimSpace(stringValue(inboundReq["action"])))
	if action != "" && action != "add" && action != "replace" && action != "update" {
		return "", nil
	}
	inbound, _ := inboundReq["inbound"].(map[string]interface{})
	if inbound == nil || !isRealityInbound(inbound) {
		return "", nil
	}
	stream, _ := inbound["streamSettings"].(map[string]interface{})
	reality, _ := stream["realitySettings"].(map[string]interface{})
	if reality == nil {
		return "", nil
	}
	dest := strings.TrimSpace(stringValue(reality["dest"]))
	host, port, err := net.SplitHostPort(dest)
	if err != nil || port != "443" {
		return "", nil
	}
	host = normalizeRealityHost(host)
	if host == "" || net.ParseIP(host) != nil {
		return "", nil
	}

	matchesServerName := false
	switch names := reality["serverNames"].(type) {
	case []interface{}:
		for _, raw := range names {
			if normalizeRealityHost(stringValue(raw)) == host {
				matchesServerName = true
				break
			}
		}
	case []string:
		for _, name := range names {
			if normalizeRealityHost(name) == host {
				matchesServerName = true
				break
			}
		}
	case string:
		for _, name := range strings.Split(names, ",") {
			if normalizeRealityHost(name) == host {
				matchesServerName = true
				break
			}
		}
	}
	if !matchesServerName {
		return "", nil
	}
	return host, reality
}

func rewriteRealityFallbackForLocalWebsites(inboundReq map[string]interface{}, localDomains map[string]struct{}) bool {
	host, reality := realityPublicFallbackHost(inboundReq)
	if host == "" {
		return false
	}
	if _, ok := localDomains[host]; !ok {
		return false
	}
	reality["dest"] = localRealityNginxDest
	reality["xver"] = 1
	return true
}

// rewriteLocalRealityFallback prevents a local camouflage website from being
// configured as its own public :443 fallback. Once tunnel-in routes that SNI
// into Reality, such a destination re-enters tunnel-in indefinitely. The
// website inventory is read from the physical Agent; on a federated consumer
// this is deliberately deferred until the owner forwards to its real Agent.
func (h *RemoteManageHandler) rewriteLocalRealityFallback(ctx context.Context, serverID int64, body []byte) ([]byte, bool) {
	var inboundReq map[string]interface{}
	if len(body) == 0 || json.Unmarshal(body, &inboundReq) != nil {
		return body, false
	}
	if host, _ := realityPublicFallbackHost(inboundReq); host == "" {
		return body, false
	}
	if h.repo == nil {
		return body, false
	}
	if _, err := h.repo.GetFederatedServer(ctx, serverID); err == nil {
		return body, false
	}

	raw, err := h.forwardToRemoteServer(ctx, serverID, http.MethodGet, "/api/child/nginx/websites", nil)
	if err != nil {
		return body, false
	}
	var inventory struct {
		Websites []struct {
			Domain string `json:"domain"`
		} `json:"websites"`
	}
	if json.Unmarshal(raw, &inventory) != nil {
		return body, false
	}
	localDomains := make(map[string]struct{}, len(inventory.Websites))
	for _, website := range inventory.Websites {
		if domain := normalizeRealityHost(website.Domain); domain != "" {
			localDomains[domain] = struct{}{}
		}
	}
	if !rewriteRealityFallbackForLocalWebsites(inboundReq, localDomains) {
		return body, false
	}
	rewritten, err := json.Marshal(inboundReq)
	if err != nil {
		return body, false
	}
	return rewritten, true
}
