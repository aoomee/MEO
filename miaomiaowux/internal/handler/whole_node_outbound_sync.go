package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"miaomiaowux/internal/storage"
)

type wholeNodeOutboundChange struct {
	serverID int64
	tag      string
	previous map[string]any
}

// updateWholeNodeOutboundsTargetingNode keeps managed whole-node routes alive
// when their external landing node is edited. The routing rule continues to
// reference the same outbound tag; only the outbound payload is replaced.
func (h *nodesHandler) updateWholeNodeOutboundsTargetingNode(ctx context.Context, oldNode *storage.Node, replacement map[string]any) ([]wholeNodeOutboundChange, error) {
	if h.remoteManage == nil || oldNode == nil || replacement == nil {
		return nil, nil
	}
	oldTarget, hasOldTarget := nodeOutboundTarget(ctx, h.repo, oldNode)
	servers, err := h.repo.ListRemoteServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list remote servers: %w", err)
	}

	var changes []wholeNodeOutboundChange
	for _, server := range servers {
		if server.Status != storage.RemoteServerStatusConnected {
			continue
		}
		raw, err := h.remoteManage.forwardToRemoteServer(ctx, server.ID, http.MethodGet, "/api/child/outbounds", nil)
		if err != nil {
			h.rollbackWholeNodeOutboundChanges(context.Background(), changes)
			return nil, fmt.Errorf("load outbounds from %s: %w", server.Name, err)
		}
		var response struct {
			Outbounds []map[string]any `json:"outbounds"`
		}
		if err := json.Unmarshal(raw, &response); err != nil {
			h.rollbackWholeNodeOutboundChanges(context.Background(), changes)
			return nil, fmt.Errorf("decode outbounds from %s: %w", server.Name, err)
		}
		for _, outbound := range response.Outbounds {
			tag := strings.TrimSpace(fmt.Sprint(outbound["tag"]))
			if tag == "" || !wholeNodeOutboundReferencesNode(tag, outbound, oldNode.ID, oldTarget, hasOldTarget) {
				continue
			}
			next := cloneOutboundMap(replacement)
			next["tag"] = tag
			body, _ := json.Marshal(map[string]any{
				"action":   "update",
				"tag":      tag,
				"outbound": next,
			})
			if _, err := h.remoteManage.forwardToRemoteServer(ctx, server.ID, http.MethodPost, "/api/child/outbounds", body); err != nil {
				h.rollbackWholeNodeOutboundChanges(context.Background(), changes)
				return nil, fmt.Errorf("update outbound %s on %s: %w", tag, server.Name, err)
			}
			changes = append(changes, wholeNodeOutboundChange{serverID: server.ID, tag: tag, previous: cloneOutboundMap(outbound)})
		}
	}
	return changes, nil
}

func wholeNodeOutboundReferencesNode(tag string, outbound map[string]any, nodeID int64, oldTarget outboundTarget, hasOldTarget bool) bool {
	if targetID, tagged := wholeNodeOutboundTargetNodeID(tag); tagged {
		return targetID == nodeID
	}
	if !strings.HasPrefix(tag, "landing-node-") && !strings.HasPrefix(tag, "landing-") {
		return false
	}
	return hasOldTarget && outboundTargetsNode(outbound, oldTarget)
}

func cloneOutboundMap(source map[string]any) map[string]any {
	if source == nil {
		return map[string]any{}
	}
	raw, err := json.Marshal(source)
	if err != nil {
		return map[string]any{}
	}
	var clone map[string]any
	if json.Unmarshal(raw, &clone) != nil || clone == nil {
		return map[string]any{}
	}
	return clone
}

func (h *nodesHandler) rollbackWholeNodeOutboundChanges(ctx context.Context, changes []wholeNodeOutboundChange) {
	if h.remoteManage == nil || len(changes) == 0 {
		return
	}
	rollbackCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	for i := len(changes) - 1; i >= 0; i-- {
		change := changes[i]
		body, _ := json.Marshal(map[string]any{
			"action":   "update",
			"tag":      change.tag,
			"outbound": change.previous,
		})
		_, _ = h.remoteManage.forwardToRemoteServer(rollbackCtx, change.serverID, http.MethodPost, "/api/child/outbounds", body)
	}
}
