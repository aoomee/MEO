package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"miaomiaowux/internal/license"
)

// 令牌轮换要不要被授权槽位迁移失败拦下来，全看这个分类。
// 分错任何一边都会造成真实事故：
//   - 该放行却拦下 → 从没连过的服务器再也重置不了令牌
//   - 该拦下却放行 → 新令牌生效但槽位还绑在旧 hash，写操作全线失败且无法自愈
func TestSlotRotationSkippable(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"迁移成功", nil, true},
		{
			"免费版/许可证不可用：根本没有槽位体系",
			license.ErrAgentLeaseEntitlementUnavailable,
			true,
		},
		{
			"包装过的 entitlement 错误也要认出来",
			fmt.Errorf("issue lease: %w", license.ErrAgentLeaseEntitlementUnavailable),
			true,
		},
		{
			"新版服务：旧 hash 上没有槽位行",
			&license.AgentLeaseIssueError{Code: "server_slot_previous_missing", StatusCode: http.StatusConflict},
			true,
		},
		{
			"旧版服务：文案归一后的 not_active",
			&license.AgentLeaseIssueError{Code: "server_slot_not_active", StatusCode: http.StatusConflict},
			true,
		},
		{
			"旧槽位存在但非 active/grace，没有生效中的授权要保护",
			&license.AgentLeaseIssueError{Code: "server_slot_state_conflict", StatusCode: http.StatusConflict},
			true,
		},
		{
			"传输失败：旧槽位状态未知，必须拦",
			errors.New("dial tcp: connection refused"),
			false,
		},
		{
			"429 限流：旧槽位还活着",
			&license.AgentLeaseIssueError{Code: "rate_limited", StatusCode: http.StatusTooManyRequests},
			false,
		},
		{
			"5xx：旧槽位还活着",
			&license.AgentLeaseIssueError{Code: "server_slot_issue_failed", StatusCode: http.StatusServiceUnavailable},
			false,
		},
		{
			"配额超限",
			&license.AgentLeaseIssueError{Code: "server_slot_quota_exceeded", StatusCode: http.StatusForbidden},
			false,
		},
		{
			"释放冷却中：槽位仍被占用",
			&license.AgentLeaseIssueError{Code: "server_slot_release_cooldown", StatusCode: http.StatusConflict},
			false,
		},
		{
			"许可证权威过期：该先修许可证，不该顺手换掉令牌",
			&license.AgentLeaseIssueError{Code: "stale_license_authority", StatusCode: http.StatusConflict},
			false,
		},
		{
			"无 code 的未知失败",
			&license.AgentLeaseIssueError{Message: "boom", StatusCode: http.StatusInternalServerError},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := slotRotationSkippable(tc.err); got != tc.want {
				t.Fatalf("slotRotationSkippable(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

// 没有 license.Manager（免费版）时，轮换不能被挡住，也不该产生待投递的 reservation。
func TestRotateServerSlotWithoutLicenseManager(t *testing.T) {
	lease, needsDelivery, err := rotateServerSlot(context.Background(), nil, "old-token", "new-token")
	if err != nil {
		t.Fatalf("免费版不该报错，得到: %v", err)
	}
	if needsDelivery {
		t.Fatal("没有许可证体系时不该有 reservation 要投递")
	}
	if lease.Reservation != "" {
		t.Fatalf("不该返回 reservation，得到: %q", lease.Reservation)
	}
}

// 请求 context 已取消时，槽位迁移不能被连带取消 —— 管理员关掉页面正是最容易
// 把轮换卡在「一半」的时机。这里用 nil manager 走通路径，验证 WithoutCancel 生效。
func TestRotateServerSlotSurvivesCanceledRequestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, _, err := rotateServerSlot(ctx, nil, "old-token", "new-token"); err != nil {
		t.Fatalf("已取消的请求 context 不该让迁移直接失败: %v", err)
	}
}

// deliverRotatedSlot 在没有 WS 的情况下必须安全返回 false，而不是 panic ——
// 投递失败只是「稍后补发」，绝不能反过来打断已经成功的轮换。
func TestDeliverRotatedSlotWithoutWSHandler(t *testing.T) {
	if deliverRotatedSlot(nil, 1, license.AgentLeaseDelivery{Reservation: "r"}) {
		t.Fatal("没有 WS handler 时不该报告投递成功")
	}
}
