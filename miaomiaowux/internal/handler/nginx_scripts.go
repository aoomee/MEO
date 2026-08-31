package handler

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"net/http"
	"os/exec"
)

// The scripts are shipped inside the MEO binary so privileged installation
// never executes mutable code fetched from another project at runtime.
//
//go:embed assets/install-nginx.sh assets/uninstall-nginx.sh
var bundledNginxScripts embed.FS

func bundledNginxScript(name string) ([]byte, error) {
	switch name {
	case "install-nginx.sh", "uninstall-nginx.sh":
		return bundledNginxScripts.ReadFile("assets/" + name)
	default:
		return nil, errors.New("unsupported nginx script")
	}
}

func bundledNginxCommand(ctx context.Context, name string, args ...string) (*exec.Cmd, error) {
	script, err := bundledNginxScript(name)
	if err != nil {
		return nil, err
	}
	commandArgs := append([]string{"-s", "--"}, args...)
	cmd := exec.CommandContext(ctx, "bash", commandArgs...)
	cmd.Stdin = bytes.NewReader(script)
	return cmd, nil
}

// ServeNginxScript lets a registered Agent fetch the exact script embedded in
// its master. Access requires the same per-server installation token used by
// the Agent binary endpoint.
func (h *XrayServerHandler) ServeNginxScript(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	token := r.URL.Query().Get("token")
	if !validInstallToken(token) {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}
	if _, err := h.repo.GetRemoteServerByToken(r.Context(), token); err != nil {
		http.Error(w, "invalid token", http.StatusForbidden)
		return
	}

	name := "install-nginx.sh"
	if r.URL.Query().Get("action") == "uninstall" {
		name = "uninstall-nginx.sh"
	}
	script, err := bundledNginxScript(name)
	if err != nil {
		http.Error(w, "script unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/x-shellscript; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	w.Header().Set("Cache-Control", "no-store")
	if r.Method == http.MethodGet {
		_, _ = w.Write(script)
	}
}
