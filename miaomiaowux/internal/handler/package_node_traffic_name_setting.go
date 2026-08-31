package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"miaomiaowux/internal/storage"
)

const packageNodeTrafficNameEnabledKey = "package_node_traffic_name_enabled"

func packageNodeTrafficNameEnabled(ctx context.Context, repo *storage.TrafficRepository) bool {
	if repo == nil {
		return false
	}
	raw, err := repo.GetSystemSetting(ctx, packageNodeTrafficNameEnabledKey)
	if err != nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func (h *SystemSettingsHandler) GetPackageNodeTrafficName(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"enabled": packageNodeTrafficNameEnabled(r.Context(), h.repo),
	})
}

func (h *SystemSettingsHandler) SetPackageNodeTrafficName(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "请求格式错误")
		return
	}
	value := "0"
	if req.Enabled {
		value = "1"
	}
	if err := h.repo.SetSystemSetting(r.Context(), packageNodeTrafficNameEnabledKey, value); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "保存节点流量名称设置失败")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"enabled": req.Enabled,
	})
}
