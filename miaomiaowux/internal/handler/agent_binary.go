package handler

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"miaomiaowux/internal/logger"
)

// ============================================================================
// Agent 二进制的分发。
//
// 主控可以自己发 Agent：把二进制放到数据目录下，
//
//	$MMWX_DATA_DIR/agent-bin/mmwx-agent-linux-amd64   （mmw-agent-linux-amd64 这个旧名也认）
//	$MMWX_DATA_DIR/agent-bin/mmwx-agent-linux-arm64
//
// 「添加服务器」生成的安装脚本就会带着安装 token 从主控自身拉
// （URL: {master_url}/api/remote/agent-binary?token=...&arch=amd64）。
//
// 安装脚本按这个顺序找二进制（任一成功即停）：
//
//	1. 主控本机 agent-bin/   —— 只有该目录里确实放了文件才会写进脚本
//	2. GitHub Release        —— 仓库见 defaultAgentGitHubRepo，tag 为 mmwx-agent 的那条
//	3. 自定义镜像            —— 设了 MMWX_AGENT_DOWNLOAD_BASE 才追加
//
// 相关环境变量（都可不设）：
//
//	MMWX_AGENT_GITHUB_REPO=owner/repo        换成你自己的仓库；设成 "off" 可完全关掉 GitHub 源
//	MMWX_GH_PROXY=https://ghproxy.example/   GitHub 加速前缀（国内网络）
//	MMWX_AGENT_DOWNLOAD_BASE=https://host/p  自建镜像，取 {base}/mmwx-agent-linux-{arch}
// ============================================================================

// defaultAgentGitHubRepo 是发布 Agent 的仓库（面板和 Agent 共用一个仓库，靠 tag 前缀区分：
// 面板 mmwx-v*，Agent mmwx-agent-v*）。
const defaultAgentGitHubRepo = ""

// agentBinaryDirName 是数据目录下存放 Agent 二进制的子目录名。
const agentBinaryDirName = "agent-bin"

// dataDirPath 返回数据目录（与 main.go 的解析口径一致：环境变量优先，默认 ./data）。
func dataDirPath() string {
	if v := strings.TrimSpace(os.Getenv("MMWX_DATA_DIR")); v != "" {
		return filepath.Clean(v)
	}
	return "data"
}

// AgentBinaryDir 返回存放 Agent 二进制的目录绝对/相对路径。
func AgentBinaryDir() string {
	return filepath.Join(dataDirPath(), agentBinaryDirName)
}

// AgentDownloadBase 返回可选的「自建镜像」地址（环境变量 MMWX_AGENT_DOWNLOAD_BASE）。默认为空。
func AgentDownloadBase() string {
	return strings.TrimRight(strings.TrimSpace(os.Getenv("MMWX_AGENT_DOWNLOAD_BASE")), "/")
}

// AgentGitHubRepo 返回发布 Agent 的 GitHub 仓库（owner/repo）。
// 设成 "off" / "none" / "-" 可以彻底关掉 GitHub 下载源，回到纯本地分发。
func AgentGitHubRepo() string {
	v := strings.TrimSpace(os.Getenv("MMWX_AGENT_GITHUB_REPO"))
	if v == "" {
		return defaultAgentGitHubRepo
	}
	switch strings.ToLower(v) {
	case "off", "none", "-", "false":
		return ""
	}
	return strings.Trim(v, "/")
}

// GitHubProxyBase 返回可选的 GitHub 加速前缀（环境变量 MMWX_GH_PROXY）。
func GitHubProxyBase() string {
	return strings.TrimRight(strings.TrimSpace(os.Getenv("MMWX_GH_PROXY")), "/")
}

// HasLocalAgentBinary 判断 agent-bin/ 里是否真的放了二进制。
// 没放就不把「从主控下载」写进安装脚本，免得每次装 Agent 都先撞一次 404。
func HasLocalAgentBinary() bool {
	for _, arch := range []string{"amd64", "arm64"} {
		if info, err := os.Stat(AgentBinaryPath(arch)); err == nil && !info.IsDir() && info.Size() > 0 {
			return true
		}
	}
	return false
}

// normalizeAgentArch 把请求里的 arch 参数收敛到白名单，避免路径穿越。
func normalizeAgentArch(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "amd64", "x86_64", "":
		return "amd64", nil
	case "arm64", "aarch64":
		return "arm64", nil
	default:
		return "", errors.New("unsupported arch")
	}
}

// agentBinaryNames 是同一个架构下可接受的文件名，按优先级排列。
// 优先 mmwx-agent-linux-*（GitHub Release 资产就叫这个，下载完直接丢进来即可），
// 同时兼容早先的 mmw-agent-linux-*。
func agentBinaryNames(arch string) []string {
	return []string{"mmwx-agent-linux-" + arch, "mmw-agent-linux-" + arch}
}

// AgentBinaryPath 返回指定架构的 Agent 二进制路径：存在哪个用哪个，
// 都不存在时返回首选名（用于报错信息里提示应该放成什么名字）。
func AgentBinaryPath(arch string) string {
	dir := AgentBinaryDir()
	names := agentBinaryNames(arch)
	for _, name := range names {
		p := filepath.Join(dir, name)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p
		}
	}
	return filepath.Join(dir, names[0])
}

// ServeAgentBinary 把管理员放在数据目录里的 Agent 二进制发给正在安装的子服务器。
//
// 鉴权：必须带一个真实存在的服务器安装 token（与安装脚本用的是同一个 token），
// 因此不会变成任意人可下载的公开端点。
func (h *XrayServerHandler) ServeAgentBinary(w http.ResponseWriter, r *http.Request) {
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
	arch, err := normalizeAgentArch(r.URL.Query().Get("arch"))
	if err != nil {
		http.Error(w, "unsupported arch", http.StatusBadRequest)
		return
	}

	path := AgentBinaryPath(arch)
	f, err := os.Open(path)
	if err != nil {
		logger.Warn("[Agent 安装] 本地未找到 Agent 二进制", "path", path, "error", err)
		http.Error(w, "Agent 二进制未就绪：请把 mmwx-agent-linux-"+arch+
			" 放到主控的 "+AgentBinaryDir()+"/ 目录下。", http.StatusNotFound)
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || info.IsDir() {
		http.Error(w, "Agent 二进制不可读", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\"mmwx-agent-linux-"+arch+"\"")
	// http.ServeContent 负责 Range / If-Modified-Since，断点续传直接可用。
	http.ServeContent(w, r, filepath.Base(path), info.ModTime(), f)
}
