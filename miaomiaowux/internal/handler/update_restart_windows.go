//go:build windows

package handler

import (
	"os"
	"os/exec"

	"miaomiaowux/internal/logger"
)

// restartSelf —— Windows 没有 syscall.Exec，只能起一个新进程再退出自己。
// （正式部署是 Linux + systemd；这里主要是为了本机开发时能编出 Windows 版。）
func restartSelf(execPath string) {
	logger.Info("[系统重启] 正在重启服务", "exec_path", execPath)

	cmd := exec.Command(execPath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		logger.Error("[系统重启] 启动新进程失败", "error", err)
		return
	}
	logger.Info("[系统重启] 新进程已启动，退出当前进程", "new_pid", cmd.Process.Pid)
	os.Exit(0)
}
