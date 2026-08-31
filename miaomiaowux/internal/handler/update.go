package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"miaomiaowux/internal/logger"
	"miaomiaowux/internal/version"
)

// ============================================================================
// 面板在线更新 —— 直接对接 GitHub Releases。
//
// 源码不开源，仓库里只放安装/更新脚本，编译好的二进制发到 Release。面板这边：
//
//	tag    mmwx                    固定 tag，每次发版更新同一个（Agent 用 mmwx-agent）
//	名字   v<版本>                 版本号写在 release 名字里，例 v0.4.8-beta.18
//	资产   mmwx-linux-<arch>       例 mmwx-linux-amd64
//
// 也兼容「每版一个 tag」的写法（mmwx-v0.4.8-beta.18），那时版本号从 tag 里取。
//
// 「检查更新」查 Releases 列表挑出 tag 前缀最新的一条；「立即更新」下载对应架构的资产、
// 原子替换正在运行的二进制、然后 exec 自己重启。失败会回滚。
//
// 相关环境变量（都可不设）：
//
//	MMWX_UPDATE_REPO=owner/repo   换仓库；设成 off 可彻底关掉在线更新（回到"恒是最新版"）
//	MMWX_GH_PROXY=https://p/      GitHub 加速前缀，下载直连失败时自动套上重试
//	MMWX_GH_TOKEN=ghp_xxx         GitHub token，只为提高 API 限额（匿名 60 次/小时）
//	MMWX_GH_API_BASE=https://…    换 API 根地址（自建镜像 / GitHub Enterprise），默认 api.github.com
// ============================================================================

const (
	// defaultUpdateRepo 是发布面板二进制的仓库。
	defaultUpdateRepo = ""
	// panelTag 是面板的固定 tag（滚动 release：每次发新版就更新这同一个 tag，
	// 版本号写在 release 名字里，例 name="v0.4.8-beta.18"）。
	panelTag = "mmwx"
	// panelTagPrefix 兼容「每版一个 tag」的写法，例 mmwx-v0.4.8-beta.18。
	// 注意 mmwx-agent / mmwx-agent-v* 两种写法都不会命中这里。
	panelTagPrefix = "mmwx-v"
	// agentTag 是 Agent 的固定滚动 tag；agentTagPrefix 是每版一个 tag 的写法。
	agentTag       = "mmwx-agent"
	agentTagPrefix = "mmwx-agent-v"
	// panelAssetPrefix 是面板二进制在 Release 里的资产名前缀。
	panelAssetPrefix = "mmwx-linux-"
	// releaseCacheTTL 是 Releases 查询的缓存时间。匿名 API 只有 60 次/小时，
	// 前端可能反复调「检查更新」，这里挡一层。
	releaseCacheTTL = 10 * time.Minute
)

// githubAPIBase 返回 GitHub API 根地址。默认官方；设 MMWX_GH_API_BASE 可指向自建镜像
// / GitHub Enterprise（也方便本地起个假 API 做联调）。
func githubAPIBase() string {
	if v := strings.TrimRight(strings.TrimSpace(os.Getenv("MMWX_GH_API_BASE")), "/"); v != "" {
		return v
	}
	return "https://api.github.com"
}

// UpdateRepo 返回面板更新用的 GitHub 仓库（owner/repo）。
// 设成 off / none / - / false 可以彻底关掉在线更新。
func UpdateRepo() string {
	v := strings.TrimSpace(os.Getenv("MMWX_UPDATE_REPO"))
	if v == "" {
		return defaultUpdateRepo
	}
	switch strings.ToLower(v) {
	case "off", "none", "-", "false":
		return ""
	}
	return strings.Trim(v, "/")
}

// updateDisabledNotice 是关掉在线更新后各端点的统一说明。
const updateDisabledNotice = "本实例已关闭在线更新（MMWX_UPDATE_REPO=off）：请手动替换二进制后重启服务。"

type UpdateChannel string

const (
	UpdateChannelStable     UpdateChannel = "stable"
	UpdateChannelPrerelease UpdateChannel = "prerelease"
)

// UpdateInfo 是 /api/admin/update/check 的返回结构。
type UpdateInfo struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	HasUpdate      bool   `json:"has_update"`
	ReleaseURL     string `json:"release_url"`
	DownloadURL    string `json:"download_url"`
	ReleaseNotes   string `json:"release_notes"`
	Channel        string `json:"channel"`
	Prerelease     bool   `json:"prerelease"`
}

// UpdateProgress 表示更新操作的进度（SSE 事件体）。
type UpdateProgress struct {
	Step     string `json:"step"`     // checking / downloading / backing_up / replacing / restarting / done / error
	Progress int    `json:"progress"` // 下载进度 0-100
	Message  string `json:"message"`
}

// ghRelease 只取我们需要的字段。
type ghRelease struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"` // release 标题；滚动 tag 的写法把版本号写在这里
	HTMLURL    string `json:"html_url"`
	Body       string `json:"body"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
	Assets     []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
	} `json:"assets"`
}

var releaseCache struct {
	sync.Mutex
	value   []ghRelease
	fetched time.Time
	repo    string
}

// SyncInstalledVersionMarker 在启动时把「已安装版本」标记刷成当前二进制的版本。
// 纯本地文件操作，不联网；手动替换二进制后诊断信息不会停留在旧版本。
func SyncInstalledVersionMarker() error {
	if isDocker() {
		return nil
	}
	return persistInstalledVersion(version.Version)
}

// ---------------------------------------------------------------------------
// 查 Release
// ---------------------------------------------------------------------------

func parseUpdateChannel(raw string) UpdateChannel {
	if strings.EqualFold(strings.TrimSpace(raw), string(UpdateChannelPrerelease)) {
		return UpdateChannelPrerelease
	}
	return UpdateChannelStable
}

// fetchReleases 拉仓库的 Release 列表（GitHub 返回按时间倒序），带 10 分钟缓存。
func fetchReleases(repo string) ([]ghRelease, error) {
	releaseCache.Lock()
	if releaseCache.repo == repo && releaseCache.value != nil && time.Since(releaseCache.fetched) < releaseCacheTTL {
		cached := releaseCache.value
		releaseCache.Unlock()
		return cached, nil
	}
	releaseCache.Unlock()

	url := githubAPIBase() + "/repos/" + repo + "/releases?per_page=100"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "miaomiaowux-updater")
	if token := strings.TrimSpace(os.Getenv("MMWX_GH_TOKEN")); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("连接 GitHub 失败: %w", err)
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests:
		return nil, errors.New("GitHub API 限流（匿名 60 次/小时），稍后再试，或给面板设置 MMWX_GH_TOKEN")
	case resp.StatusCode == http.StatusNotFound:
		return nil, fmt.Errorf("仓库 %s 不存在或不可访问", repo)
	case resp.StatusCode != http.StatusOK:
		return nil, fmt.Errorf("GitHub API 状态码 %d", resp.StatusCode)
	}

	var releases []ghRelease
	if err := json.NewDecoder(io.LimitReader(resp.Body, 8<<20)).Decode(&releases); err != nil {
		return nil, fmt.Errorf("解析 Release 列表: %w", err)
	}

	releaseCache.Lock()
	releaseCache.value = releases
	releaseCache.fetched = time.Now()
	releaseCache.repo = repo
	releaseCache.Unlock()
	return releases, nil
}

// isPanelReleaseTag 判断这条 release 是不是面板的。
// 固定 tag（mmwx）和每版一个 tag（mmwx-v*）两种都认；mmwx-agent* 一律不认。
func isPanelReleaseTag(tag string) bool {
	tag = strings.TrimSpace(tag)
	return tag == panelTag || strings.HasPrefix(tag, panelTagPrefix)
}

func isAgentReleaseTag(tag string) bool {
	tag = strings.TrimSpace(tag)
	return tag == agentTag || strings.HasPrefix(tag, agentTagPrefix)
}

// agentReleaseVersion 从 Agent Release 读出版本号。
// 名字写成 v0.6.0 时直接用；名字写成 mmwx-agent-v0.6.0 这种不像版本号时，从 tag 后缀取。
func agentReleaseVersion(r *ghRelease) string {
	if r == nil {
		return ""
	}
	if v := normalizeVersionString(r.Name); v != "" {
		return v
	}
	if strings.HasPrefix(r.TagName, agentTagPrefix) {
		return normalizeVersionString(strings.TrimPrefix(r.TagName, agentTagPrefix))
	}
	return ""
}

// pickAgentRelease 在仓库里挑出版本号最高的那条 Agent Release（不是列表里最先出现的）。
func pickAgentRelease(releases []ghRelease) *ghRelease {
	var best *ghRelease
	bestVer := ""
	for i := range releases {
		r := &releases[i]
		if r.Draft || !isAgentReleaseTag(r.TagName) {
			continue
		}
		ver := agentReleaseVersion(r)
		if ver == "" {
			continue
		}
		if best == nil || compareSemver(ver, bestVer) > 0 {
			best = r
			bestVer = ver
		}
	}
	return best
}

// normalizeVersionString 把 "v0.4.8-beta.18" / "0.4.8-beta.18" 收敛成 "0.4.8-beta.18"。
// 不像版本号（空、或首字符不是数字）时返回空。
func normalizeVersionString(raw string) string {
	v := strings.TrimPrefix(strings.TrimSpace(raw), "v")
	if v == "" || v[0] < '0' || v[0] > '9' {
		return ""
	}
	return v
}

// releaseVersion 取这条 release 的版本号：优先 release 名字（滚动 tag 的写法把版本写在这），
// 名字不像版本号时再退回从 tag 后缀取。都取不到返回空。
func releaseVersion(r *ghRelease) string {
	if v := normalizeVersionString(r.Name); v != "" {
		return v
	}
	if strings.HasPrefix(r.TagName, panelTagPrefix) {
		return normalizeVersionString(strings.TrimPrefix(r.TagName, panelTagPrefix))
	}
	return ""
}

// pickPanelRelease 从列表里挑出面板要用的那条。
//   - stable：优先没标 prerelease 的；一条都没有时退回最新的一条（省得仓库里全标了 pre 就查不到）
//   - prerelease：直接取最新的一条
func pickPanelRelease(releases []ghRelease, channel UpdateChannel) *ghRelease {
	var newest *ghRelease
	for i := range releases {
		r := &releases[i]
		if r.Draft || !isPanelReleaseTag(r.TagName) {
			continue
		}
		if newest == nil {
			newest = r
		}
		if channel == UpdateChannelStable && r.Prerelease {
			continue
		}
		return r
	}
	return newest
}

// assetForCurrentArch 找当前架构对应的资产下载地址。
func assetForCurrentArch(r *ghRelease) (string, int64) {
	want := panelAssetPrefix + runtime.GOARCH
	for _, a := range r.Assets {
		if a.Name == want {
			return a.BrowserDownloadURL, a.Size
		}
	}
	return "", 0
}

// checkLatestVersion 查最新版本并和当前版本比对。
func checkLatestVersion(channel UpdateChannel) (*UpdateInfo, error) {
	current := strings.TrimPrefix(strings.TrimSpace(version.Version), "v")

	repo := UpdateRepo()
	if repo == "" {
		return &UpdateInfo{
			CurrentVersion: current,
			LatestVersion:  current,
			HasUpdate:      false,
			ReleaseNotes:   updateDisabledNotice,
			Channel:        "off",
		}, nil
	}

	releases, err := fetchReleases(repo)
	if err != nil {
		return nil, err
	}
	rel := pickPanelRelease(releases, channel)
	if rel == nil {
		// 仓库能访问、只是还没发过面板的版本。这不是错误，按「已是最新」返回，
		// 免得前端弹一个看不懂的 500（顺带把说明写在 release_notes 里）。
		logger.Info("[系统更新] 仓库里暂无面板 Release", "repo", repo, "tag", panelTag)
		return &UpdateInfo{
			CurrentVersion: current,
			LatestVersion:  current,
			HasUpdate:      false,
			ReleaseNotes: fmt.Sprintf("仓库 %s 里还没有面板的 Release（tag 应为 %s，或 %s 开头）。",
				repo, panelTag, panelTagPrefix),
			Channel: string(channel),
		}, nil
	}

	latest := releaseVersion(rel)
	if latest == "" {
		// 认出了 release 但读不出版本号（名字既不是 v1.2.3 也没写在 tag 里）。
		// 不敢乱比大小，就当没有更新，把原因摆出来让管理员能看懂。
		logger.Warn("[系统更新] Release 里读不出版本号", "tag", rel.TagName, "name", rel.Name)
		return &UpdateInfo{
			CurrentVersion: current,
			LatestVersion:  current,
			HasUpdate:      false,
			ReleaseURL:     rel.HTMLURL,
			ReleaseNotes: fmt.Sprintf("找到了 Release（tag=%s），但它的名字 %q 不是版本号，无法判断新旧。"+
				"请把 release 名字改成 v0.0.0 这种形式。", rel.TagName, rel.Name),
			Channel: string(channel),
		}, nil
	}
	downloadURL, _ := assetForCurrentArch(rel)

	return &UpdateInfo{
		CurrentVersion: current,
		LatestVersion:  latest,
		HasUpdate:      compareSemVersion(latest, current) > 0,
		ReleaseURL:     rel.HTMLURL,
		DownloadURL:    downloadURL,
		ReleaseNotes:   rel.Body,
		Channel:        string(channel),
		Prerelease:     rel.Prerelease,
	}, nil
}

// ---------------------------------------------------------------------------
// HTTP 端点
// ---------------------------------------------------------------------------

// NewUpdateCheckHandler 返回检查更新的处理程序。
func NewUpdateCheckHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeUpdateError(w, http.StatusMethodNotAllowed, errors.New("only GET is supported"))
			return
		}
		info, err := checkLatestVersion(parseUpdateChannel(r.URL.Query().Get("channel")))
		if err != nil {
			writeUpdateError(w, http.StatusInternalServerError, fmt.Errorf("检查更新失败: %w", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(info)
	})
}

// preflightUpdate 是 apply / apply-sse 共用的前置检查：拿到可下载的新版本信息。
func preflightUpdate(channel UpdateChannel, force bool, wantTarget string) (*UpdateInfo, error) {
	if UpdateRepo() == "" {
		return nil, errors.New(updateDisabledNotice)
	}
	if isDocker() {
		return nil, errors.New("容器内不能替换主控二进制，请更新镜像")
	}
	info, err := checkLatestVersion(channel)
	if err != nil {
		return nil, fmt.Errorf("检查更新失败: %w", err)
	}
	if !info.HasUpdate && !force {
		return nil, errors.New("已是最新版本")
	}
	if target := strings.TrimPrefix(strings.TrimSpace(wantTarget), "v"); target != "" && target != info.LatestVersion {
		return nil, errors.New("目标版本已变化，请重新检查更新")
	}
	if info.DownloadURL == "" {
		return nil, fmt.Errorf("该 Release 里没有 %s%s，请确认发版时上传了对应架构的二进制",
			panelAssetPrefix, runtime.GOARCH)
	}
	return info, nil
}

// applyUpdate 下载并替换二进制，成功后返回替换到的路径（调用方负责重启）。
func applyUpdate(info *UpdateInfo, onProgress func(downloaded, total int64), onRetry func(string)) (string, error) {
	tempFile, err := downloadBinaryWithProgressAndRetry(info.DownloadURL, onProgress, onRetry)
	if err != nil {
		return "", fmt.Errorf("下载失败: %w", err)
	}
	defer os.Remove(tempFile)

	targetPath, err := getUpdateTargetPath()
	if err != nil {
		return "", fmt.Errorf("获取程序路径失败: %w", err)
	}

	backupPath := targetPath + ".bak"
	if err := copyFile(targetPath, backupPath); err != nil {
		logger.Warn("[系统更新] 备份当前版本失败（非致命）", "error", err)
	}

	logger.Info("[系统更新] 替换二进制", "from", tempFile, "to", targetPath)
	if err := replaceBinary(tempFile, targetPath); err != nil {
		// 换失败就把备份放回去，别把面板搞成起不来。
		// 注意 copyFile 走 os.Create，恢复出来的文件默认没有可执行位，必须补 chmod。
		if _, statErr := os.Stat(backupPath); statErr == nil {
			restoreErr := copyFile(backupPath, targetPath)
			if restoreErr == nil {
				restoreErr = os.Chmod(targetPath, 0o755)
			}
			if restoreErr != nil {
				logger.Error("[系统更新] 替换失败，自动回滚也失败，请手动恢复",
					"backup", backupPath, "target", targetPath, "error", restoreErr)
				return "", fmt.Errorf("替换失败: %w（自动回滚未成功，旧版本备份在 %s，可手动复制回 %s）",
					err, backupPath, targetPath)
			}
			logger.Warn("[系统更新] 替换失败，已回滚到旧版本", "error", err)
		}
		return "", fmt.Errorf("替换失败: %w", err)
	}
	if err := os.Chmod(targetPath, 0o755); err != nil {
		return "", fmt.Errorf("设置权限失败: %w", err)
	}
	if err := persistInstalledVersion(info.LatestVersion); err != nil {
		logger.Warn("[系统更新] 写入已安装版本标记失败", "error", err)
	}
	return targetPath, nil
}

// NewUpdateApplyHandler 一次性执行更新（无进度）。
func NewUpdateApplyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeUpdateError(w, http.StatusMethodNotAllowed, errors.New("only POST is supported"))
			return
		}
		info, err := preflightUpdate(parseUpdateChannel(r.URL.Query().Get("channel")), false, "")
		if err != nil {
			writeUpdateError(w, http.StatusBadRequest, err)
			return
		}
		logger.Info("[系统更新] 开始更新", "from", info.CurrentVersion, "to", info.LatestVersion, "url", info.DownloadURL)
		targetPath, err := applyUpdate(info, nil, nil)
		if err != nil {
			writeUpdateError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "更新完成，正在重启...",
			"version": info.LatestVersion,
		})

		go func() {
			time.Sleep(500 * time.Millisecond)
			restartSelf(targetPath)
		}()
	})
}

// NewUpdateApplySSEHandler 带进度的更新（前端进度条走这个）。
func NewUpdateApplySSEHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no") // 禁用 nginx 缓冲

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}
		send := func(step string, progress int, message string) {
			data, _ := json.Marshal(UpdateProgress{Step: step, Progress: progress, Message: message})
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}

		send("checking", 0, "正在检查版本信息...")
		info, err := preflightUpdate(
			parseUpdateChannel(r.URL.Query().Get("channel")),
			r.URL.Query().Get("force") == "true",
			r.URL.Query().Get("target"),
		)
		if err != nil {
			send("error", 0, err.Error())
			return
		}

		send("downloading", 0, "正在下载 v"+info.LatestVersion+" ...")
		logger.Info("[系统更新] 开始下载", "url", info.DownloadURL)
		last := 0
		targetPath, err := applyUpdate(info,
			func(downloaded, total int64) {
				if total <= 0 {
					return
				}
				p := int(downloaded * 100 / total)
				// 每 5% 推一次，别把 SSE 刷爆
				if p >= last+5 || p == 100 {
					last = p
					send("downloading", p, fmt.Sprintf("正在下载... %d%%", p))
				}
			},
			func(proxyURL string) {
				last = 0
				send("downloading", 0, "直连失败，正在通过加速地址重试...")
			},
		)
		if err != nil {
			send("error", 0, err.Error())
			return
		}

		send("restarting", 0, "更新完成，正在重启服务...")
		send("done", 100, "更新完成，已升级到 v"+info.LatestVersion)
		logger.Info("[系统更新] 更新成功，准备重启", "version", info.LatestVersion)

		go func() {
			time.Sleep(500 * time.Millisecond)
			restartSelf(targetPath)
		}()
	})
}

// ---------------------------------------------------------------------------
// 下载
// ---------------------------------------------------------------------------

// downloadBinaryWithProgressAndRetry 先直连下载；失败且配了 MMWX_GH_PROXY 就套加速前缀重试。
func downloadBinaryWithProgressAndRetry(url string, onProgress func(downloaded, total int64), onRetry func(proxyURL string)) (string, error) {
	tempFile, err := downloadBinaryDirect(url, onProgress, 5*time.Minute)
	if err == nil {
		return tempFile, nil
	}

	proxy := GitHubProxyBase()
	if proxy == "" {
		return "", err
	}
	logger.Warn("[系统更新] 直连下载失败，改用加速地址", "error", err)
	proxyURL := proxy + "/" + url
	if onRetry != nil {
		onRetry(proxyURL)
	}
	tempFile, err = downloadBinaryDirect(proxyURL, onProgress, 10*time.Minute)
	if err != nil {
		return "", fmt.Errorf("加速地址也失败: %w", err)
	}
	return tempFile, nil
}

// downloadBinaryDirect 下载到临时文件。先试 4 线程分段（HTTP Range），不行回退单线程。
func downloadBinaryDirect(url string, onProgress func(downloaded, total int64), timeout time.Duration) (string, error) {
	if path, perr := downloadRangedParallel(url, 4, onProgress, timeout); perr == nil {
		return path, nil
	} else {
		logger.Info("[系统更新] 分段下载不可用，回退单线程", "reason", perr.Error())
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载返回状态码: %d", resp.StatusCode)
	}

	tempFile, err := os.CreateTemp("", "mmwx-update-*")
	if err != nil {
		return "", err
	}
	total := resp.ContentLength
	var downloaded int64

	if onProgress == nil || total <= 0 {
		if _, err := io.Copy(tempFile, resp.Body); err != nil {
			tempFile.Close()
			os.Remove(tempFile.Name())
			return "", err
		}
	} else {
		buf := make([]byte, 32*1024)
		for {
			n, readErr := resp.Body.Read(buf)
			if n > 0 {
				if _, writeErr := tempFile.Write(buf[:n]); writeErr != nil {
					tempFile.Close()
					os.Remove(tempFile.Name())
					return "", writeErr
				}
				downloaded += int64(n)
				onProgress(downloaded, total)
			}
			if readErr != nil {
				if readErr == io.EOF {
					break
				}
				tempFile.Close()
				os.Remove(tempFile.Name())
				return "", readErr
			}
		}
	}
	tempFile.Close()
	return tempFile.Name(), nil
}

// downloadRangedParallel 用 threads 个连接分段并发下载。
// 服务器不支持 Range / 大小未知 / 文件过小 / 任一分段失败 → 返回错误，由调用方回退单线程。
func downloadRangedParallel(url string, threads int, onProgress func(downloaded, total int64), timeout time.Duration) (string, error) {
	probe, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	probe.Header.Set("Range", "bytes=0-0")
	pc := &http.Client{Timeout: timeout}
	presp, err := pc.Do(probe)
	if err != nil {
		return "", err
	}
	presp.Body.Close()
	if presp.StatusCode != http.StatusPartialContent {
		return "", fmt.Errorf("range 不支持(状态 %d)", presp.StatusCode)
	}
	var total int64 = -1
	if cr := presp.Header.Get("Content-Range"); cr != "" {
		if i := strings.LastIndex(cr, "/"); i >= 0 {
			total, _ = strconv.ParseInt(strings.TrimSpace(cr[i+1:]), 10, 64)
		}
	}
	if total < 4<<20 || threads < 2 {
		return "", errors.New("文件过小或线程数不足,不走并发")
	}

	tempFile, err := os.CreateTemp("", "mmwx-update-*")
	if err != nil {
		return "", err
	}
	name := tempFile.Name()
	if err := tempFile.Truncate(total); err != nil {
		tempFile.Close()
		os.Remove(name)
		return "", err
	}

	part := (total + int64(threads) - 1) / int64(threads)
	var downloaded int64
	var wg sync.WaitGroup
	errs := make([]error, threads)
	for i := 0; i < threads; i++ {
		start := int64(i) * part
		end := start + part - 1
		if end >= total {
			end = total - 1
		}
		if start > end {
			break
		}
		wg.Add(1)
		go func(idx int, start, end int64) {
			defer wg.Done()
			req, e := http.NewRequest(http.MethodGet, url, nil)
			if e != nil {
				errs[idx] = e
				return
			}
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
			rc := &http.Client{Timeout: timeout}
			rp, e := rc.Do(req)
			if e != nil {
				errs[idx] = e
				return
			}
			defer rp.Body.Close()
			if rp.StatusCode != http.StatusPartialContent {
				errs[idx] = fmt.Errorf("分段 %d 状态 %d", idx, rp.StatusCode)
				return
			}
			buf := make([]byte, 64*1024)
			off := start
			for {
				n, re := rp.Body.Read(buf)
				if n > 0 {
					if _, we := tempFile.WriteAt(buf[:n], off); we != nil {
						errs[idx] = we
						return
					}
					off += int64(n)
					if onProgress != nil {
						onProgress(atomic.AddInt64(&downloaded, int64(n)), total)
					}
				}
				if re != nil {
					if re == io.EOF {
						break
					}
					errs[idx] = re
					return
				}
			}
		}(i, start, end)
	}
	wg.Wait()
	for _, e := range errs {
		if e != nil {
			tempFile.Close()
			os.Remove(name)
			return "", e
		}
	}
	tempFile.Close()
	return name, nil
}

// ---------------------------------------------------------------------------
// 本地文件工具
// ---------------------------------------------------------------------------

// getUpdateTargetPath 返回当前正在运行的可执行文件路径（解开 symlink）。
func getUpdateTargetPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return "", err
	}
	return execPath, nil
}

// replaceBinary 把 dst 换成 src。Linux 上可以删掉正在运行的文件（inode 还在内存里）。
func replaceBinary(src, dst string) error {
	if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
		if err := os.Rename(src, dst); err == nil {
			return nil
		}
		return copyFile(src, dst)
	}
	if err := os.Rename(src, dst); err != nil {
		// 跨设备 rename 会失败，退回复制
		return copyFile(src, dst)
	}
	return nil
}

// persistInstalledVersion 把版本标记写到数据目录同级的 .version 文件。
func persistInstalledVersion(latest string) error {
	dataDir := strings.TrimSpace(os.Getenv("MMWX_DATA_DIR"))
	if dataDir == "" {
		dataDir = "data"
	}
	marker := filepath.Join(filepath.Dir(filepath.Clean(dataDir)), ".version")
	if err := os.MkdirAll(filepath.Dir(marker), 0o755); err != nil {
		return err
	}
	value := "v" + strings.TrimPrefix(strings.TrimSpace(latest), "v") + "\n"
	tmp, err := os.CreateTemp(filepath.Dir(marker), ".version-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0o644); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.WriteString(value); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, marker)
}

// copyFile 把文件从 src 复制到 dst（本地文件工具，供数据迁移等模块复用）。
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Sync()
}

// isDocker 检查是否在 Docker 容器内运行。
func isDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	if os.Getenv("DOCKER") == "1" {
		return true
	}
	data, err := os.ReadFile("/proc/1/cgroup")
	if err == nil && strings.Contains(string(data), "docker") {
		return true
	}
	return false
}

// IsDockerEnvironment exposes the shared container detection to startup wiring.
func IsDockerEnvironment() bool { return isDocker() }

// ---------------------------------------------------------------------------
// 语义版本比较 —— 更新判断 + agent 握手时校验最低主控版本都用它。
// ---------------------------------------------------------------------------

type semVersion struct {
	major, minor, patch int
	pre                 []string
	valid               bool
}

func parseSemVersion(v string) semVersion {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	v = strings.SplitN(v, "+", 2)[0]
	parts := strings.SplitN(v, "-", 2)
	var out semVersion
	if _, err := fmt.Sscanf(parts[0], "%d.%d.%d", &out.major, &out.minor, &out.patch); err != nil {
		return out
	}
	if parts[0] != fmt.Sprintf("%d.%d.%d", out.major, out.minor, out.patch) {
		return out
	}
	if len(parts) == 2 {
		if parts[1] == "" {
			return out
		}
		out.pre = strings.Split(parts[1], ".")
	}
	out.valid = true
	return out
}

func compareSemVersion(a, b string) int {
	av, bv := parseSemVersion(a), parseSemVersion(b)
	if !av.valid || !bv.valid {
		return strings.Compare(a, b)
	}
	for _, pair := range [][2]int{{av.major, bv.major}, {av.minor, bv.minor}, {av.patch, bv.patch}} {
		if pair[0] < pair[1] {
			return -1
		}
		if pair[0] > pair[1] {
			return 1
		}
	}
	if len(av.pre) == 0 && len(bv.pre) > 0 {
		return 1
	}
	if len(av.pre) > 0 && len(bv.pre) == 0 {
		return -1
	}
	for i := 0; i < len(av.pre) || i < len(bv.pre); i++ {
		if i >= len(av.pre) {
			return -1
		}
		if i >= len(bv.pre) {
			return 1
		}
		var ai, bi int
		aNum, aErr := fmt.Sscanf(av.pre[i], "%d", &ai)
		bNum, bErr := fmt.Sscanf(bv.pre[i], "%d", &bi)
		if aErr == nil && bErr == nil && aNum == 1 && bNum == 1 {
			if ai < bi {
				return -1
			}
			if ai > bi {
				return 1
			}
			continue
		}
		if aErr == nil {
			return -1
		}
		if bErr == nil {
			return 1
		}
		if cmp := strings.Compare(av.pre[i], bv.pre[i]); cmp != 0 {
			return cmp
		}
	}
	return 0
}

func writeUpdateError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
