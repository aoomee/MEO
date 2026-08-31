package handler

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"os/exec"
)

// The audited installer scripts are part of the Agent binary. Administrator
// actions never pipe mutable scripts from a third-party repository into bash.
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
