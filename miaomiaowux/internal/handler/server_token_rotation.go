package handler

import (
	"context"
	"errors"
	"log"
	"time"

	"miaomiaowux/internal/license"
	"miaomiaowux/internal/storage"
)

// 服务器令牌轮换要同时改三个地方:主控库里的 token、Agent 手上的 token、以及许可证
// 服务上那条 Guard 绑定的授权槽位(槽位按 sha256(server_token) 索引)。
//
// 原实现的顺序是「改库 → 推给 Agent → 最后才迁槽位」,而且迁槽位只在 WS 恰好连着
// 时才会执行、失败也只记一条日志。于是任何一次迁移失败(网络抖动、429、Agent 离线)
// 都会留下「新 token 已生效,但槽位还绑在旧 hash 上」的半完成状态 —— 之后所有需要
// 授权槽位的写操作都会失败,而且没有任何自动修复路径。
//
// 现在的顺序是「改库 → 迁槽位 → 推给 Agent」,迁槽位失败就把库里的 token 回滚。
// 新令牌由 ResetServerToken 生成,所以落库无法排在迁移之后;但本地写回滚廉价又可靠,
// 于是把唯一会真正失败的远端操作(迁槽位)提到推送之前,让它成为轮换的闸门:
// Agent 要么拿到「新令牌 + 已就位的槽位」,要么什么都没收到、这次轮换等于没发生。

// slotRotationSkippable 判断一次迁移失败是否等价于「旧 token 上本来就没有可迁移的
// 授权槽位」。这类失败必须放行 —— 否则从没连上来过的服务器、或者槽位已经释放的
// 服务器就再也重置不了令牌了。
//
// 除此之外的失败(传输错误、429、5xx、配额超限、许可证权威过期)都意味着旧槽位可能
// 仍然活着,必须拦下来回滚,不能放任 token 单方面变更。
func slotRotationSkippable(err error) bool {
	if err == nil {
		return true
	}
	// 免费版 / 许可证暂时不可用:根本没有槽位这回事,轮换不该被挡住。
	if errors.Is(err, license.ErrAgentLeaseEntitlementUnavailable) {
		return true
	}
	var issueErr *license.AgentLeaseIssueError
	if !errors.As(err, &issueErr) {
		// 传输层失败没有结构化 code,旧槽位状态未知 —— 按最坏情况处理。
		return false
	}
	switch issueErr.Code {
	case "server_slot_previous_missing":
		// 新版许可证服务:旧 hash 上没有槽位行。
		return true
	case "server_slot_not_active":
		// 旧版服务只回英文文案,normalizeAgentLeaseIssueCode 归一成这个 code。
		return true
	case "server_slot_state_conflict":
		// 旧槽位存在但不是 active/grace(reserved 未激活 / 已 released),
		// 没有正在生效的授权需要保护。
		return true
	default:
		return false
	}
}

// rotateServerSlot 请求许可证服务把同一条 Guard 绑定的授权槽位从旧 server token
// hash 原地迁到新 hash。只跟许可证服务打交道,不要求 Agent 在线 —— 投递 reservation
// 才需要 WS,那是调用方的事。
//
// needsDelivery 为 true 时,lease 是一张需要送到 Agent 的新 reservation。
//
// 用 context.WithoutCancel:管理员关掉页面不该把一次已经发出去的槽位迁移打断在
// 中途,那正是最难收拾的状态。
func rotateServerSlot(ctx context.Context, manager *license.Manager, oldToken, newToken string) (license.AgentLeaseDelivery, bool, error) {
	if manager == nil {
		return license.AgentLeaseDelivery{}, false, nil
	}
	base := context.WithoutCancel(ctx)
	lease, err := issueSlotRotation(base, manager, oldToken, newToken)

	// 只有传输层失败才重放:这类错误说明响应可能只是丢在路上,迁移也许已经在许可证
	// 服务落库了。服务端把 previous_server_hash 持久化在 reservation 上,原样重放同
	// 一次迁移会返回同一张 reservation;真没落库则等于重新迁一次 —— 两种情况都收敛。
	// 反过来,只要拿到了结构化的服务端响应,就说明请求确实被处理并拒绝了,没有歧义,
	// 重放只会平白多打一次。
	var issueErr *license.AgentLeaseIssueError
	if err != nil && !errors.As(err, &issueErr) && !errors.Is(err, license.ErrAgentLeaseEntitlementUnavailable) {
		log.Printf("[Token Rotation] 槽位迁移传输失败,幂等重放一次: %v", err)
		lease, err = issueSlotRotation(base, manager, oldToken, newToken)
	}

	if err != nil {
		if slotRotationSkippable(err) {
			log.Printf("[Token Rotation] 旧令牌上没有可迁移的授权槽位,继续轮换: %v", err)
			return license.AgentLeaseDelivery{}, false, nil
		}
		return license.AgentLeaseDelivery{}, false, err
	}
	return lease, true, nil
}

func issueSlotRotation(ctx context.Context, manager *license.Manager, oldToken, newToken string) (license.AgentLeaseDelivery, error) {
	leaseCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	return manager.IssueAgentLeaseTokenRotation(leaseCtx, oldToken, newToken)
}

// deliverRotatedSlot 把迁移后的 reservation 送给 Agent。投递失败不算轮换失败:
// 槽位已经在许可证服务侧迁好了,Agent 重连/心跳时会按新 token hash 自己再要一次。
func deliverRotatedSlot(ws *RemoteWSHandler, serverID int64, lease license.AgentLeaseDelivery) bool {
	if ws == nil {
		return false
	}
	conn, ok := ws.GetConnectionByServerID(serverID)
	if !ok {
		return false
	}
	if err := ws.deliverAgentLease(conn, lease, ""); err != nil {
		log.Printf("[Token Rotation] 授权槽位已迁移但投递失败 server=%d: %v", serverID, err)
		return false
	}
	return true
}

// rollbackServerTokenRotation 是 rotateServerSlot 的补偿动作:把库里的 token 改回去。
//
// 这里刻意**不**做反向槽位迁移。走到这一步说明槽位仍然绑在旧 token hash 上(重放已经
// 排除了「其实迁成功了只是响应丢了」这种歧义),旧 token 回到库里之后二者自然一致。
// 而反向迁移在许可证服务侧本来也不成立:刚迁过去的行处于 reserved 状态,而迁移分支
// 只接受 active/grace,反向调用必被 ErrServerSlotState 拒掉。
func rollbackServerTokenRotation(repo *storage.TrafficRepository, serverID int64,
	oldToken string, oldExpiresAt *time.Time) error {
	return repo.RestoreServerToken(context.Background(), serverID, oldToken, oldExpiresAt)
}
