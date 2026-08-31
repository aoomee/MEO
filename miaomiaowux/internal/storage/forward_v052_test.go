package storage

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func TestParseForwardRuleChainID(t *testing.T) {
	id, ok := ParseForwardRuleChainID("fwd-c12-p1000-h0")
	if !ok || id != 12 {
		t.Fatalf("got %d ok=%v", id, ok)
	}
	if _, ok := ParseForwardRuleChainID("other"); ok {
		t.Fatal("expected miss")
	}
}

func TestForwardChainBranchesAndDailyTraffic(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	g1, err := repo.CreateForwardGroup(ctx, &ForwardGroup{Name: "entry", BalanceStrategy: "round_robin"})
	if err != nil {
		t.Fatal(err)
	}
	g2, err := repo.CreateForwardGroup(ctx, &ForwardGroup{Name: "exit", BalanceStrategy: "round_robin"})
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.SetForwardGroupMembers(ctx, g1, []ForwardGroupMember{{ServerID: 11}, {ServerID: 12}}); err != nil {
		t.Fatal(err)
	}
	if err := repo.SetForwardGroupMembers(ctx, g2, []ForwardGroupMember{{ServerID: 21}}); err != nil {
		t.Fatal(err)
	}

	chainID, err := repo.CreateForwardChain(ctx, "c1")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.SetForwardChainHops(ctx, chainID, []int64{g1, g2}); err != nil {
		t.Fatal(err)
	}

	chain, err := repo.GetForwardChain(ctx, chainID)
	if err != nil {
		t.Fatal(err)
	}
	if len(chain.Branches) != 3 {
		t.Fatalf("期望 3 条 branch(入口2+出口1), 得到 %d: %+v", len(chain.Branches), chain.Branches)
	}
	if chain.Branches[0].HopSeq != 0 || chain.Branches[0].ViaGroupID != g1 || chain.Branches[0].ServerID != 11 {
		t.Fatalf("入口首条 branch 不对: %+v", chain.Branches[0])
	}
	if chain.Branches[2].HopSeq != 1 || chain.Branches[2].ServerID != 21 {
		t.Fatalf("出口 branch 不对: %+v", chain.Branches[2])
	}

	ruleID := "fwd-c" + strconv.FormatInt(chainID, 10) + "-p1000-h0"
	if err := repo.InsertForwardHopMetrics(ctx, []ForwardHopMetric{
		{ServerID: 11, RuleID: ruleID, BytesUp: 10, BytesDown: 20},
		{ServerID: 11, RuleID: ruleID, BytesUp: 30, BytesDown: 5},
	}); err != nil {
		t.Fatal(err)
	}
	daily, err := repo.ListForwardDailyTraffic(ctx, chainID, 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(daily) != 1 || daily[0].BytesUp != 30 || daily[0].BytesDown != 20 {
		t.Fatalf("日流量应取当日 max, 得到 %+v", daily)
	}
	if daily[0].Date != time.Now().Format("2006-01-02") {
		t.Fatalf("日期 %s", daily[0].Date)
	}

	// 改组成员后 branches 跟着变
	if err := repo.SetForwardGroupMembers(ctx, g1, []ForwardGroupMember{{ServerID: 19}}); err != nil {
		t.Fatal(err)
	}
	chain, _ = repo.GetForwardChain(ctx, chainID)
	if len(chain.Branches) != 2 || chain.Branches[0].ServerID != 19 {
		t.Fatalf("组成员变更后 branches 未重建: %+v", chain.Branches)
	}
}
