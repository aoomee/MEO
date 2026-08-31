package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type wholeNodeOutboundRequest struct {
	Outbound     map[string]any `json:"outbound,omitempty"`
	BalancerTag  string         `json:"balancer_tag,omitempty"`
	TargetNodeID int64          `json:"target_node_id,omitempty"`
}

type wholeNodeRoutingResponse struct {
	Routing struct {
		Rules []map[string]any `json:"rules"`
	} `json:"routing"`
}

type wholeNodeOutboundsResponse struct {
	Outbounds []map[string]any `json:"outbounds"`
}

// handleWholeNodeOutbound owns the lifecycle of an inbound-wide landing route.
// PUT switches the target; DELETE cancels it. Newly created outbounds use a
// node-specific tag so they can be removed without touching shared outbounds.
func (h *nodesHandler) handleWholeNodeOutbound(w http.ResponseWriter, r *http.Request, idSegment string) {
	nodeID, err := strconv.ParseInt(strings.Trim(idSegment, "/"), 10, 64)
	if err != nil || nodeID <= 0 {
		writeError(w, http.StatusBadRequest, errors.New("invalid node id"))
		return
	}
	node, err := h.repo.GetNodeByID(r.Context(), nodeID)
	if err != nil {
		writeError(w, http.StatusNotFound, errors.New("node not found"))
		return
	}
	if node.NodeType == "routed" || strings.TrimSpace(node.InboundTag) == "" || strings.TrimSpace(node.OriginalServer) == "" {
		writeError(w, http.StatusBadRequest, errors.New("node is not an editable physical inbound"))
		return
	}
	server, err := h.repo.GetRemoteServerByName(r.Context(), node.OriginalServer)
	if err != nil || server == nil {
		writeError(w, http.StatusNotFound, errors.New("source server not found"))
		return
	}
	if h.remoteManage == nil {
		writeError(w, http.StatusServiceUnavailable, errors.New("remote manager unavailable"))
		return
	}

	var req wholeNodeOutboundRequest
	if r.Method == http.MethodPut {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, errors.New("invalid request body"))
			return
		}
		if req.Outbound == nil && strings.TrimSpace(req.BalancerTag) == "" {
			writeError(w, http.StatusBadRequest, errors.New("outbound or balancer_tag is required"))
			return
		}
		if req.Outbound != nil && strings.TrimSpace(req.BalancerTag) != "" {
			writeError(w, http.StatusBadRequest, errors.New("outbound and balancer_tag are mutually exclusive"))
			return
		}
		if req.Outbound != nil && req.TargetNodeID > 0 {
			targetNode, err := h.repo.GetNodeByID(r.Context(), req.TargetNodeID)
			if err != nil {
				writeError(w, http.StatusBadRequest, errors.New("landing target node not found"))
				return
			}
			target, ok := nodeOutboundTarget(r.Context(), h.repo, &targetNode)
			if !ok || !outboundTargetsNode(req.Outbound, target) {
				writeError(w, http.StatusBadRequest, errors.New("landing outbound does not match target node"))
				return
			}
		}
	}

	result, err := h.replaceWholeNodeOutbound(r.Context(), server.ID, nodeID, node.InboundTag, r.Method, req)
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *nodesHandler) replaceWholeNodeOutbound(ctx context.Context, serverID, nodeID int64, inboundTag, method string, req wholeNodeOutboundRequest) (map[string]any, error) {
	routing, err := h.fetchWholeNodeRouting(ctx, serverID)
	if err != nil {
		return nil, err
	}
	outbounds, err := h.fetchWholeNodeOutbounds(ctx, serverID)
	if err != nil {
		return nil, err
	}

	newTag := ""
	newRule := map[string]any{"type": "field", "inboundTag": []string{inboundTag}}
	if method == http.MethodPut && req.Outbound != nil {
		if req.TargetNodeID > 0 {
			newTag = fmt.Sprintf("landing-node-%d-target-%d-%d", nodeID, req.TargetNodeID, time.Now().UnixMilli())
		} else {
			newTag = fmt.Sprintf("landing-node-%d-%d", nodeID, time.Now().UnixMilli())
		}
		req.Outbound["tag"] = newTag // never trust/reuse a client-supplied tag
		body, _ := json.Marshal(map[string]any{"action": "add", "outbound": req.Outbound})
		if _, err := h.remoteManage.ForwardToServer(ctx, serverID, http.MethodPost, "/api/child/outbounds", body); err != nil {
			return nil, fmt.Errorf("add replacement outbound: %w", err)
		}
		newRule["outboundTag"] = newTag
	} else if method == http.MethodPut {
		newRule["balancerTag"] = strings.TrimSpace(req.BalancerTag)
	}

	if method == http.MethodPut {
		body, _ := json.Marshal(map[string]any{"action": "add_rule", "rule": newRule})
		if _, err := h.remoteManage.ForwardToServer(ctx, serverID, http.MethodPost, "/api/child/routing", body); err != nil {
			if newTag != "" {
				_ = h.removeWholeNodeOutboundTag(context.Background(), serverID, newTag)
			}
			return nil, fmt.Errorf("add replacement routing rule: %w", err)
		}
	}

	// Remove old rules from the end. The just-added replacement is not in the
	// snapshot, so it cannot be selected accidentally.
	oldOwnedTags := map[string]bool{}
	// Include orphaned managed outbounds as well. A previous interrupted update
	// may have created the outbound before its routing rule was persisted.
	for _, outbound := range outbounds.Outbounds {
		tag := wholeNodeString(outbound["tag"])
		if isOwnedWholeNodeOutboundTag(tag, nodeID, inboundTag) {
			oldOwnedTags[tag] = true
		}
	}
	for i := len(routing.Routing.Rules) - 1; i >= 0; i-- {
		rule := routing.Routing.Rules[i]
		if !isWholeNodeRuleForInbound(rule, inboundTag) {
			continue
		}
		oldTag := wholeNodeString(rule["outboundTag"])
		body, _ := json.Marshal(map[string]any{"action": "remove_rule", "index": i})
		if _, err := h.remoteManage.ForwardToServer(ctx, serverID, http.MethodPost, "/api/child/routing", body); err != nil {
			return nil, fmt.Errorf("remove previous routing rule: %w", err)
		}
		if isOwnedWholeNodeOutboundTag(oldTag, nodeID, inboundTag) {
			oldOwnedTags[oldTag] = true
		}
	}

	for tag := range oldOwnedTags {
		if tag == newTag || wholeNodeOutboundReferencedElsewhere(routing.Routing.Rules, tag, inboundTag) {
			continue
		}
		if wholeNodeOutboundExists(outbounds.Outbounds, tag) {
			if err := h.removeWholeNodeOutboundTag(ctx, serverID, tag); err != nil {
				return nil, err
			}
		}
	}

	// Routing changes are persisted by the agent but need a restart to become
	// active. Restart only once after the complete replacement.
	restartBody := []byte(`{"service":"xray","action":"restart"}`)
	if _, err := h.remoteManage.ForwardToServer(ctx, serverID, http.MethodPost, "/api/child/services/control", restartBody); err != nil {
		return nil, fmt.Errorf("restart xray after whole-node outbound update: %w", err)
	}
	return map[string]any{"success": true, "outbound_tag": newTag, "balancer_tag": strings.TrimSpace(req.BalancerTag)}, nil
}

func (h *nodesHandler) fetchWholeNodeRouting(ctx context.Context, serverID int64) (wholeNodeRoutingResponse, error) {
	var result wholeNodeRoutingResponse
	raw, err := h.remoteManage.ForwardToServer(ctx, serverID, http.MethodGet, "/api/child/routing", nil)
	if err != nil {
		return result, fmt.Errorf("load routing: %w", err)
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return result, fmt.Errorf("decode routing: %w", err)
	}
	return result, nil
}

func (h *nodesHandler) fetchWholeNodeOutbounds(ctx context.Context, serverID int64) (wholeNodeOutboundsResponse, error) {
	var result wholeNodeOutboundsResponse
	raw, err := h.remoteManage.ForwardToServer(ctx, serverID, http.MethodGet, "/api/child/outbounds", nil)
	if err != nil {
		return result, fmt.Errorf("load outbounds: %w", err)
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return result, fmt.Errorf("decode outbounds: %w", err)
	}
	return result, nil
}

func (h *nodesHandler) removeWholeNodeOutboundTag(ctx context.Context, serverID int64, tag string) error {
	body, _ := json.Marshal(map[string]any{"action": "remove", "tag": tag})
	if _, err := h.remoteManage.ForwardToServer(ctx, serverID, http.MethodPost, "/api/child/outbounds", body); err != nil {
		return fmt.Errorf("remove previous outbound %s: %w", tag, err)
	}
	return nil
}

func isWholeNodeRuleForInbound(rule map[string]any, inboundTag string) bool {
	inbounds, ok := rule["inboundTag"].([]any)
	if !ok {
		if values, typed := rule["inboundTag"].([]string); typed {
			for _, value := range values {
				if value == inboundTag {
					inbounds = []any{value}
					break
				}
			}
		}
	}
	found := false
	for _, value := range inbounds {
		if fmt.Sprint(value) == inboundTag {
			found = true
			break
		}
	}
	if !found || (wholeNodeString(rule["outboundTag"]) == "" && wholeNodeString(rule["balancerTag"]) == "") {
		return false
	}
	for _, key := range []string{"domain", "ip", "port", "sourcePort", "network", "source", "user", "protocol", "attrs"} {
		if wholeNodeRuleValuePresent(rule[key]) {
			return false
		}
	}
	return true
}

func wholeNodeRuleValuePresent(value any) bool {
	if value == nil {
		return false
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed) != ""
	case []any:
		return len(typed) > 0
	case []string:
		return len(typed) > 0
	default:
		return true
	}
}

func wholeNodeString(value any) string {
	if value == nil {
		return ""
	}
	if text, ok := value.(string); ok {
		return strings.TrimSpace(text)
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func isOwnedWholeNodeOutboundTag(tag string, nodeID int64, inboundTag string) bool {
	return strings.HasPrefix(tag, fmt.Sprintf("landing-node-%d-", nodeID)) || strings.HasPrefix(tag, "landing-"+inboundTag+"-")
}

func wholeNodeOutboundTargetNodeID(tag string) (int64, bool) {
	const marker = "-target-"
	if !strings.HasPrefix(tag, "landing-node-") {
		return 0, false
	}
	start := strings.Index(tag, marker)
	if start < 0 {
		return 0, false
	}
	rest := tag[start+len(marker):]
	end := strings.IndexByte(rest, '-')
	if end <= 0 {
		return 0, false
	}
	id, err := strconv.ParseInt(rest[:end], 10, 64)
	return id, err == nil && id > 0
}

func wholeNodeOutboundReferencedElsewhere(rules []map[string]any, tag, removedInbound string) bool {
	for _, rule := range rules {
		if wholeNodeString(rule["outboundTag"]) != tag {
			continue
		}
		if !isWholeNodeRuleForInbound(rule, removedInbound) {
			return true
		}
	}
	return false
}

func wholeNodeOutboundExists(outbounds []map[string]any, tag string) bool {
	for _, outbound := range outbounds {
		if strings.TrimSpace(fmt.Sprint(outbound["tag"])) == tag {
			return true
		}
	}
	return false
}
