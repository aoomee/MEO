package storage

import (
	"context"
	"testing"
)

func TestForwardGroupCRUD(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	// 建组
	id, err := repo.CreateForwardGroup(ctx, &ForwardGroup{
		Name:               "entry",
		BalanceStrategy:    "percentage",
		DNSDomain:          "fwd.example.com",
		DNSProviderID:      3,
		FailoverEnabled:    true,
		OfflineMsThreshold: 800,
	})
	if err != nil {
		t.Fatalf("CreateForwardGroup: %v", err)
	}
	if id == 0 {
		t.Fatal("期望非零自增 id")
	}

	// 设成员(全量替换)
	if err := repo.SetForwardGroupMembers(ctx, id, []ForwardGroupMember{
		{ServerID: 11, Weight: 3},
		{ServerID: 12}, // weight<=0 → 归一为 1
	}); err != nil {
		t.Fatalf("SetForwardGroupMembers: %v", err)
	}

	got, err := repo.GetForwardGroup(ctx, id)
	if err != nil {
		t.Fatalf("GetForwardGroup: %v", err)
	}
	if got.BalanceStrategy != "percentage" || !got.FailoverEnabled || got.OfflineMsThreshold != 800 {
		t.Fatalf("字段回读不一致: %+v", got)
	}
	if len(got.Members) != 2 {
		t.Fatalf("期望 2 个成员,得到 %d", len(got.Members))
	}
	if got.Members[0].Weight != 3 || got.Members[1].Weight != 1 {
		t.Fatalf("权重不一致: %+v", got.Members)
	}

	// 全量替换成员(应清掉旧的)
	if err := repo.SetForwardGroupMembers(ctx, id, []ForwardGroupMember{{ServerID: 20, Weight: 5}}); err != nil {
		t.Fatalf("SetForwardGroupMembers 二次: %v", err)
	}
	members, _ := repo.ListForwardGroupMembers(ctx, id)
	if len(members) != 1 || members[0].ServerID != 20 {
		t.Fatalf("全量替换未生效: %+v", members)
	}

	// 更新组
	got.Name = "entry-renamed"
	got.BalanceStrategy = "round_robin"
	if err := repo.UpdateForwardGroup(ctx, got); err != nil {
		t.Fatalf("UpdateForwardGroup: %v", err)
	}
	reread, _ := repo.GetForwardGroup(ctx, id)
	if reread.Name != "entry-renamed" || reread.BalanceStrategy != "round_robin" {
		t.Fatalf("更新未生效: %+v", reread)
	}

	// 删组(级联删成员)
	if err := repo.DeleteForwardGroup(ctx, id); err != nil {
		t.Fatalf("DeleteForwardGroup: %v", err)
	}
	if members, _ := repo.ListForwardGroupMembers(ctx, id); len(members) != 0 {
		t.Fatalf("删组后成员应清空,残留 %d", len(members))
	}
}

func TestForwardChainCRUD(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	// 先建 3 个组当跳
	var gids []int64
	for _, name := range []string{"A", "B", "C"} {
		gid, err := repo.CreateForwardGroup(ctx, &ForwardGroup{Name: name, BalanceStrategy: "round_robin"})
		if err != nil {
			t.Fatalf("CreateForwardGroup %s: %v", name, err)
		}
		gids = append(gids, gid)
	}

	chainID, err := repo.CreateForwardChain(ctx, "chain-abc")
	if err != nil {
		t.Fatalf("CreateForwardChain: %v", err)
	}
	if err := repo.SetForwardChainHops(ctx, chainID, gids); err != nil {
		t.Fatalf("SetForwardChainHops: %v", err)
	}

	chain, err := repo.GetForwardChain(ctx, chainID)
	if err != nil {
		t.Fatalf("GetForwardChain: %v", err)
	}
	if len(chain.Hops) != 3 {
		t.Fatalf("期望 3 跳,得到 %d", len(chain.Hops))
	}
	// 顺序 = 入参顺序,seq 0-based
	for i, h := range chain.Hops {
		if h.Seq != i || h.GroupID != gids[i] {
			t.Fatalf("跳序不一致: hop[%d]=%+v", i, h)
		}
	}

	// 全量替换跳(倒序)
	rev := []int64{gids[2], gids[1], gids[0]}
	if err := repo.SetForwardChainHops(ctx, chainID, rev); err != nil {
		t.Fatalf("SetForwardChainHops 二次: %v", err)
	}
	chain, _ = repo.GetForwardChain(ctx, chainID)
	if len(chain.Hops) != 3 || chain.Hops[0].GroupID != gids[2] {
		t.Fatalf("倒序替换未生效: %+v", chain.Hops)
	}

	// 删链(级联删跳)
	if err := repo.DeleteForwardChain(ctx, chainID); err != nil {
		t.Fatalf("DeleteForwardChain: %v", err)
	}
	if hops, _ := repo.ListForwardChainHops(ctx, chainID); len(hops) != 0 {
		t.Fatalf("删链后跳应清空,残留 %d", len(hops))
	}
}
