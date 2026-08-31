package proxygroups

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// 默认代理组配置内置在二进制里，启动时不再向任何远程地址发请求。
// 想改成远程源时,设环境变量 PROXY_GROUPS_SOURCE_URL,或在系统设置/同步接口里显式传入地址。
//
//go:embed proxy-groups-default.json
var embeddedDefaultConfig []byte

// BuiltinSource 是内置配置在日志/接口里显示的来源标识。
const BuiltinSource = "builtin"

var (
	ErrInvalidConfig  = errors.New("proxy groups config is invalid")
	ErrDownloadFailed = errors.New("proxy groups config download failed")
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// ResolveSourceURL 解析配置源地址。
// 优先级: 传入参数 > 环境变量 PROXY_GROUPS_SOURCE_URL > 空(表示用内置配置)。
func ResolveSourceURL(overrideURL string) string {
	if overrideURL != "" {
		return overrideURL
	}

	if env := os.Getenv("PROXY_GROUPS_SOURCE_URL"); env != "" {
		return env
	}

	return ""
}

// FetchConfig 取代理组配置：未配置来源时用内置副本，配置了来源才远程下载。两种情况都会校验。
// 数据仅保存在内存中,不写入磁盘
// 返回值:
//   - []byte: 配置数据
//   - string: 解析后的实际 URL
//   - error: 错误信息
func FetchConfig(overrideURL string) ([]byte, string, error) {
	resolvedURL := ResolveSourceURL(overrideURL)

	// 没有显式配置来源 → 用内置配置,不发任何网络请求。
	if resolvedURL == "" {
		normalized, err := NormalizeConfig(embeddedDefaultConfig)
		if err != nil {
			return nil, BuiltinSource, err
		}
		return normalized, BuiltinSource, nil
	}

	data, err := downloadConfig(resolvedURL)
	if err != nil {
		return nil, resolvedURL, err
	}

	// 规范化并验证配置有效性
	normalized, err := NormalizeConfig(data)
	if err != nil {
		return nil, resolvedURL, err
	}

	return normalized, resolvedURL, nil
}

// 从远程地址下载配置
func downloadConfig(sourceURL string) ([]byte, error) {
	resp, err := httpClient.Get(sourceURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: unexpected status %d from %s", ErrDownloadFailed, resp.StatusCode, sourceURL)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}

	return data, nil
}
