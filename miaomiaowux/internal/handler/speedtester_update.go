package handler

import (
	"errors"
	"net/http"
	"strings"

	"miaomiaowux/internal/storage"
)

// ============================================================================
// 测速端（speedtester）不提供在线更新。
//
// 原实现会去 GitHub Releases API 查最新版、下载并下发给各测速端自我替换。
// 这里全部移除：update-info 只列出各测速端「当前上报的版本」，恒不提示可更新；
// update-all 直接拒绝。需要升级测速端时，请自行编译并在测速端机器上替换二进制。
// ============================================================================

const speedTesterOfflineNotice = "本版本已移除测速端在线更新：请自行编译新版测速端并在其所在机器上替换二进制。"

// testerPlatform 从测速端上报的能力里解析出 os/arch，以及它是否声明支持自我更新。
func testerPlatform(t storage.SpeedTester) (string, string, bool) {
	var goos, arch string
	for _, c := range t.Caps {
		if strings.HasPrefix(c, "os:") {
			goos = strings.TrimPrefix(c, "os:")
		}
		if strings.HasPrefix(c, "arch:") {
			arch = strings.TrimPrefix(c, "arch:")
		}
	}
	return goos, arch, t.HasCap("update") && goos != "" && arch != ""
}

func normalizeTesterVersion(v string) string { return strings.TrimPrefix(strings.TrimSpace(v), "v") }

// handleTesterUpdateInfo 只回显各测速端当前版本，不发起任何外部请求。
func (h *SpeedTestHandler) handleTesterUpdateInfo(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.ListSpeedTesters(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	items := make([]map[string]any, 0, len(list))
	for _, t := range list {
		_, _, supported := testerPlatform(t)
		items = append(items, map[string]any{
			"id":               t.ID,
			"name":             t.Name,
			"version":          t.Version,
			"online":           h.testerWS != nil && h.testerWS.Online(t.ID),
			"update_supported": supported,
			"update_available": false,
		})
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"success":        true,
		"latest_version": "",
		"has_update":     false,
		"outdated_count": 0,
		"testers":        items,
		"message":        speedTesterOfflineNotice,
	})
}

// handleTesterUpdateAll 在自托管构建中一律拒绝。
func (h *SpeedTestHandler) handleTesterUpdateAll(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errors.New(speedTesterOfflineNotice))
}
