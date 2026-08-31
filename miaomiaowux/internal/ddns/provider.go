// Package ddns 把 agent 上报的 IPv4/IPv6 同步到 DNS provider 的 A/AAAA 记录。
// 跟 internal/acme 解耦:acme 只用 TXT 做 DNS-01 验证,DDNS 是常驻 + 并发 + 不同 record type,
// 直接复用 acme/dns_providers.go 里的凭据 JSON key 命名(CF_DNS_API_TOKEN 等)即可,
// 不要复用 SetDNSCredentialEnv(env var 模式不是线程安全的)。
package ddns

import (
	"context"
	"errors"
	"fmt"
)

// Provider 定义一个 DNS 服务商的 A/AAAA 记录 upsert 能力。
// "upsert" 语义:同 fqdn+recordType 已存在 → 更新 content;不存在 → 创建。
// 实现里要内部 lookup zone 并切分 sub(参考 zone.go 的 splitZone)。
type Provider interface {
	// UpsertRecord 把 fqdn 的 recordType 记录指到 content。
	// recordType ∈ {"A", "AAAA"};content 是 IP 字符串。
	// ttl=0 表示让 provider 用默认值(通常是 auto / 120-300s)。
	UpsertRecord(ctx context.Context, fqdn string, recordType string, content string, ttl int) error

	// CanManage 只读探测:该 provider 的账号是否托管了 fqdn 所属的 zone(能否改这个域名的记录)。
	// 供「自动模式无证书兜底」时遍历 dns_providers 找能管辖该域名的那个。
	// true=能管;false+err=不能管/探测失败(遍历时直接跳过)。不写任何记录。
	CanManage(ctx context.Context, fqdn string) (bool, error)

	// ReconcileRecordSet 把 fqdn 的 recordType 记录集调成【恰好】desiredContents(缺则增、多则删,幂等)。
	// 用于入口组 DNS 负载均衡:多台健康成员 → 多条 A 记录;成员掉线则摘除对应记录。
	// desiredContents 为空 → 删光该 fqdn+recordType 的全部记录。ttl=0 用 provider 默认值。
	ReconcileRecordSet(ctx context.Context, fqdn string, recordType string, desiredContents []string, ttl int) error
}

// diffRecordSet 算出「要新增的 content」和「要删除的现存记录」。
// desired/existing 都按 content 去重比较;existing 里 content 命中 desired 的保留,未命中的进删除集。
func diffRecordSet(desired []string, existing map[string]bool) (toAdd []string, keep map[string]bool) {
	want := map[string]bool{}
	for _, c := range desired {
		want[c] = true
	}
	keep = map[string]bool{}
	for c := range want {
		if existing[c] {
			keep[c] = true
		} else {
			toAdd = append(toAdd, c)
		}
	}
	return toAdd, keep
}

// ProviderType 跟 dns_providers.provider_type 列、acme/dns_providers.go DNSProviderEnvKeys 的 key 一致。
const (
	ProviderTypeCloudflare   = "cloudflare"
	ProviderTypeAlidns       = "alidns"
	ProviderTypeTencentCloud = "tencentcloud"
	ProviderTypeDNSPod       = "dnspod"
	ProviderTypeNamesilo     = "namesilo"
	ProviderTypeGoDaddy      = "godaddy"
)

// ErrUnsupportedProvider — providerType 不在已实现列表里
var ErrUnsupportedProvider = errors.New("unsupported DNS provider type for DDNS")

// NewProvider 工厂:按 providerType 派发到具体实现。
// credentials 的 key 沿用 acme/dns_providers.go 里的环境变量名(如 CF_DNS_API_TOKEN / ALICLOUD_ACCESS_KEY),
// UI 复用现有 DNS provider 表单,不需要为 DDNS 单独存一份凭据。
func NewProvider(providerType string, credentials map[string]string) (Provider, error) {
	switch providerType {
	case ProviderTypeCloudflare:
		return newCloudflareProvider(credentials)
	case ProviderTypeAlidns:
		return newAlidnsProvider(credentials)
	case ProviderTypeTencentCloud:
		return newTencentProvider(credentials)
	case ProviderTypeDNSPod:
		return newDNSPodProvider(credentials)
	case ProviderTypeNamesilo:
		return newNamesiloProvider(credentials)
	case ProviderTypeGoDaddy:
		return newGoDaddyProvider(credentials)
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedProvider, providerType)
	}
}
