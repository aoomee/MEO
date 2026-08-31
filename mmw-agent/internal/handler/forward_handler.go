package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"mmw-agent/internal/forward"
)

type forwardApplyRequest struct {
	Rules []forward.Rule `json:"rules"`
}

// HandleForwardApply 处理 POST /api/child/forward/apply。整机全量替换转发规则。
func (h *ManageHandler) HandleForwardApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	if !h.authenticate(r) {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, "read body failed")
		return
	}
	var req forwardApplyRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if h.forwardMgr == nil {
		h.forwardMgr = forward.NewManager()
	}
	if err := h.forwardMgr.Apply(req.Rules); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

// SetForwardPersistPath 设置转发规则落盘路径并立即从盘上恢复(重启不断转发)。
// 由 main.go 在拿到配置文件路径后调用一次。
func (h *ManageHandler) SetForwardPersistPath(path string) {
	if h.forwardMgr == nil {
		h.forwardMgr = forward.NewManager()
	}
	h.forwardMgr.SetPersistPath(path)
	if err := h.forwardMgr.Restore(); err != nil {
		log.Printf("[forward] restore on startup failed: %v", err)
	}
}

// HandleForwardStatus 处理 GET /api/child/forward/status。
func (h *ManageHandler) HandleForwardStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	if !h.authenticate(r) {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var rules []forward.RuleStatus
	if h.forwardMgr != nil {
		rules = h.forwardMgr.Status()
	}
	if rules == nil {
		rules = []forward.RuleStatus{}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"rules":   rules,
	})
}
