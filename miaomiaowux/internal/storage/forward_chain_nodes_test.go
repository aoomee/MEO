package storage

import (
	"context"
	"testing"
)

func TestForwardChainNodeBinding(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	// 未绑定 → nil
	if b, err := repo.GetForwardChainBinding(ctx, 100); err != nil || b != nil {
		t.Fatalf("未绑定应返回 (nil,nil),得到 b=%v err=%v", b, err)
	}

	// 绑定
	if err := repo.BindForwardChainNode(ctx, 100, 7, 10000); err != nil {
		t.Fatalf("BindForwardChainNode: %v", err)
	}
	b, err := repo.GetForwardChainBinding(ctx, 100)
	if err != nil || b == nil {
		t.Fatalf("GetForwardChainBinding: b=%v err=%v", b, err)
	}
	if b.ChainID != 7 || b.Port != 10000 {
		t.Fatalf("绑定回读不一致: %+v", b)
	}

	// 改绑(node_id 唯一 → 覆盖)
	if err := repo.BindForwardChainNode(ctx, 100, 8, 20000); err != nil {
		t.Fatalf("改绑: %v", err)
	}
	b, _ = repo.GetForwardChainBinding(ctx, 100)
	if b.ChainID != 8 || b.Port != 20000 {
		t.Fatalf("改绑未覆盖: %+v", b)
	}
	all, _ := repo.ListForwardChainBindings(ctx)
	if len(all) != 1 {
		t.Fatalf("改绑后应仍只有 1 条绑定,得到 %d", len(all))
	}

	// 另一个节点绑同一链
	if err := repo.BindForwardChainNode(ctx, 101, 8, 20001); err != nil {
		t.Fatalf("bind 101: %v", err)
	}
	byChain, _ := repo.ListForwardChainBindingsByChain(ctx, 8)
	if len(byChain) != 2 {
		t.Fatalf("链 8 应有 2 个绑定,得到 %d", len(byChain))
	}

	// 解绑
	if err := repo.UnbindForwardChainNode(ctx, 100); err != nil {
		t.Fatalf("Unbind: %v", err)
	}
	if b, _ := repo.GetForwardChainBinding(ctx, 100); b != nil {
		t.Fatalf("解绑后应为 nil,得到 %+v", b)
	}
	if all, _ := repo.ListForwardChainBindings(ctx); len(all) != 1 {
		t.Fatalf("解绑后应剩 1 条,得到 %d", len(all))
	}
}

func TestForwardChainBindingOfficialColumns(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()
	if err := repo.BindForwardChainNodeFull(ctx, ForwardChainBinding{
		NodeID: 200, ChainID: 9, Port: 50662,
		TerminusAddr: "entry.example.com:50662", Protocol: "tcp", PinnedExitServerID: 33,
	}); err != nil {
		t.Fatal(err)
	}
	b, err := repo.GetForwardChainBinding(ctx, 200)
	if err != nil || b == nil {
		t.Fatalf("b=%v err=%v", b, err)
	}
	if b.TerminusAddr != "entry.example.com:50662" || b.PinnedExitServerID != 33 {
		t.Fatalf("官方列回读不对: %+v", b)
	}
}

func TestDeleteForwardChainClearsBindings(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()
	chainID, err := repo.CreateForwardChain(ctx, "wipe")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.BindForwardChainNode(ctx, 301, chainID, 1000); err != nil {
		t.Fatal(err)
	}
	if err := repo.DeleteForwardChain(ctx, chainID); err != nil {
		t.Fatal(err)
	}
	if b, _ := repo.GetForwardChainBinding(ctx, 301); b != nil {
		t.Fatalf("删链后绑定应清空, 得到 %+v", b)
	}
}
