//go:build !windows

package handler

import (
	"os"
	"os/exec"
	"syscall"

	"miaomiaowux/internal/logger"
)

// restartSelf 用新二进制替换当前进程。
// 优先 syscall.Exec —— PID 不变，systemd 不会把它当成"服务退出"，也不会触发 Restart 计数。
func restartSelf(execPath string) {
	logger.Info("[系统重启] 正在重启服务", "exec_path", execPath)

	err := syscall.Exec(execPath, os.Args, os.Environ())
	if err != nil {
		logger.Warn("[系统重启] syscall.Exec 失败，改为启动新进程", "error", err)

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
}
