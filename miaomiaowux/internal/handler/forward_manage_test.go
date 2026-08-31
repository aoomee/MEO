package handler

import (
	"testing"

	"miaomiaowux/internal/storage"
)

func TestBuildForwardRulesChain(t *testing.T) {
	// 链 A(1,2) → B(3,4) → C(5,6),端口 10000。
	// B 组 percentage 策略、权重 2:1;开启故障转移、RTT 阈值 500。
	groups := map[int64]*storage.ForwardGroup{
		10: {ID: 10, Name: "A", BalanceStrategy: "round_robin", Members: []storage.ForwardGroupMember{{ServerID: 1}, {ServerID: 2}}},
		20: {ID: 20, Name: "B", BalanceStrategy: "percentage", FailoverEnabled: true, OfflineMsThreshold: 500,
			Members: []storage.ForwardGroupMember{{ServerID: 3, Weight: 2}, {ServerID: 4, Weight: 1}}},
		30: {ID: 30, Name: "C", BalanceStrategy: "round_robin", Members: []storage.ForwardGroupMember{{ServerID: 5}, {ServerID: 6}}},
	}
	servers := map[int64]*storage.RemoteServer{
		1: {ID: 1, IPAddress: "10.0.0.1"}, 2: {ID: 2, IPAddress: "10.0.0.2"},
		3: {ID: 3, IPAddress: "10.0.0.3"}, 4: {ID: 4, IPAddress: "10.0.0.4"},
		5: {ID: 5, IPAddress: "10.0.0.5"}, 6: {ID: 6, IPAddress: "10.0.0.6"},
	}
	chain := &storage.ForwardChain{ID: 7, Name: "abc", Hops: []storage.ForwardChainHop{
		{Seq: 0, GroupID: 10}, {Seq: 1, GroupID: 20}, {Seq: 2, GroupID: 30},
	}}

	rules, err := buildForwardRules(chain, groups, servers, 10000, "tcp")
	if err != nil {
		t.Fatalf("buildForwardRules: %v", err)
	}

	// 出口组 C(5,6)不应有规则
	if _, ok := rules[5]; ok {
		t.Fatal("出口组服务器 5 不应有转发规则")
	}
	if _, ok := rules[6]; ok {
		t.Fatal("出口组服务器 6 不应有转发规则")
	}

	// A 组每台(1,2)→ B 组成员;策略 round_robin(A→B 取目标组 B 的... 等一下,策略取目标组)
	// 实际:A→B 跳的目标组是 B,B 是 percentage → weighted。
	r1 := rules[1]
	if len(r1) != 1 {
		t.Fatalf("服务器 1 应有 1 条规则,得到 %d", len(r1))
	}
	if r1[0].Listen != ":10000" || r1[0].Protocol != "tcp" {
		t.Fatalf("规则监听/协议不对: %+v", r1[0])
	}
	if r1[0].Strategy != "weighted" {
		t.Fatalf("A→B 跳策略应取目标组 B(percentage→weighted),得到 %q", r1[0].Strategy)
	}
	if len(r1[0].Upstreams) != 2 {
		t.Fatalf("A→B 应有 2 个上游,得到 %d", len(r1[0].Upstreams))
	}
	if r1[0].Upstreams[0].Addr != "10.0.0.3:10000" || r1[0].Upstreams[0].Weight != 2 {
		t.Fatalf("上游[0] 不对: %+v", r1[0].Upstreams[0])
	}
	if r1[0].Upstreams[1].Addr != "10.0.0.4:10000" || r1[0].Upstreams[1].Weight != 1 {
		t.Fatalf("上游[1] 不对: %+v", r1[0].Upstreams[1])
	}
	if !r1[0].Health.Enabled || r1[0].Health.RTTThresholdMs != 500 {
		t.Fatalf("健康配置应取目标组 B: %+v", r1[0].Health)
	}
	// A 组两台规则一致(同一跳)
	if rules[1][0].ID != rules[2][0].ID {
		t.Fatal("A 组同跳两台规则 ID 应一致")
	}

	// B 组每台(3,4)→ C 组成员;目标组 C 是 round_robin
	r3 := rules[3]
	if len(r3) != 1 || r3[0].Strategy != "round_robin" {
		t.Fatalf("B→C 跳策略应取目标组 C(round_robin): %+v", r3)
	}
	if len(r3[0].Upstreams) != 2 || r3[0].Upstreams[0].Addr != "10.0.0.5:10000" {
		t.Fatalf("B→C 上游不对: %+v", r3[0].Upstreams)
	}
	// C 组无故障转移 → Health.Enabled=false
	if r3[0].Health.Enabled {
		t.Fatalf("C 组未开故障转移,Health.Enabled 应为 false: %+v", r3[0].Health)
	}
}

// TestAggregateServerForwardRulesMultiChain 验证同一台 server 在两条链里的规则被聚合、互不覆盖。
func TestAggregateServerForwardRulesMultiChain(t *testing.T) {
	servers := map[int64]*storage.RemoteServer{
		1: {ID: 1, IPAddress: "10.0.0.1"},
		2: {ID: 2, IPAddress: "10.0.0.2"},
		3: {ID: 3, IPAddress: "10.0.0.3"},
	}
	// 链1: A(s1) → B(s2),端口 10001;链2: C(s1) → D(s3),端口 10002。s1 同时是两条链的入口。
	groups := map[int64]*storage.ForwardGroup{
		1: {ID: 1, Name: "A", BalanceStrategy: "round_robin", Members: []storage.ForwardGroupMember{{ServerID: 1}}},
		2: {ID: 2, Name: "B", BalanceStrategy: "round_robin", Members: []storage.ForwardGroupMember{{ServerID: 2}}},
		3: {ID: 3, Name: "C", BalanceStrategy: "round_robin", Members: []storage.ForwardGroupMember{{ServerID: 1}}},
		4: {ID: 4, Name: "D", BalanceStrategy: "round_robin", Members: []storage.ForwardGroupMember{{ServerID: 3}}},
	}
	chains := map[int64]*storage.ForwardChain{
		1: {ID: 1, Hops: []storage.ForwardChainHop{{Seq: 0, GroupID: 1}, {Seq: 1, GroupID: 2}}},
		2: {ID: 2, Hops: []storage.ForwardChainHop{{Seq: 0, GroupID: 3}, {Seq: 1, GroupID: 4}}},
	}
	topo := forwardTopology{
		bindings: []storage.ForwardChainBinding{
			{NodeID: 10, ChainID: 1, Port: 10001},
			{NodeID: 11, ChainID: 2, Port: 10002},
		},
		chains:  chains,
		groups:  groups,
		servers: servers,
	}

	rules, err := aggregateServerForwardRules(1, topo)
	if err != nil {
		t.Fatalf("aggregate: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("s1 应从两条链各得 1 条规则(共 2 条),得到 %d: %+v", len(rules), rules)
	}
	// 两条规则监听不同端口、指向不同上游,ID 不冲突
	byListen := map[string]string{}
	for _, r := range rules {
		if len(r.Upstreams) != 1 {
			t.Fatalf("每条规则应 1 个上游: %+v", r)
		}
		byListen[r.Listen] = r.Upstreams[0].Addr
	}
	if byListen[":10001"] != "10.0.0.2:10001" {
		t.Fatalf("链1 规则不对: %v", byListen)
	}
	if byListen[":10002"] != "10.0.0.3:10002" {
		t.Fatalf("链2 规则不对: %v", byListen)
	}

	// s2 只在链1出口 → 出口不产生规则 → 空
	if r2, _ := aggregateServerForwardRules(2, topo); len(r2) != 0 {
		t.Fatalf("s2 是出口,不应有规则,得到 %+v", r2)
	}

	// 一条链结构不完整(引用不存在的组)时跳过,不影响另一条
	topo.chains[2].Hops[1].GroupID = 999 // D 组不存在
	rules, aggErr := aggregateServerForwardRules(1, topo)
	if aggErr == nil {
		t.Fatal("不完整的链应返回被跳过的原因")
	}
	if len(rules) != 1 || rules[0].Listen != ":10001" {
		t.Fatalf("坏链应被跳过、好链保留,得到 %+v", rules)
	}
}

// TestForwardingServerIDs 只应挑出真正承载转发规则的服务器(排除纯出口)。
func TestForwardingServerIDs(t *testing.T) {
	servers := map[int64]*storage.RemoteServer{
		1: {ID: 1, IPAddress: "10.0.0.1"}, 2: {ID: 2, IPAddress: "10.0.0.2"}, 3: {ID: 3, IPAddress: "10.0.0.3"},
	}
	groups := map[int64]*storage.ForwardGroup{
		1: {ID: 1, Members: []storage.ForwardGroupMember{{ServerID: 1}}}, // 入口
		2: {ID: 2, Members: []storage.ForwardGroupMember{{ServerID: 2}}}, // 中间
		3: {ID: 3, Members: []storage.ForwardGroupMember{{ServerID: 3}}}, // 出口
	}
	chains := map[int64]*storage.ForwardChain{
		1: {ID: 1, Hops: []storage.ForwardChainHop{{Seq: 0, GroupID: 1}, {Seq: 1, GroupID: 2}, {Seq: 2, GroupID: 3}}},
	}
	topo := forwardTopology{
		bindings: []storage.ForwardChainBinding{{NodeID: 1, ChainID: 1, Port: 10000}},
		chains:   chains, groups: groups, servers: servers,
	}
	ids := forwardingServerIDs(topo)
	set := map[int64]bool{}
	for _, id := range ids {
		set[id] = true
	}
	// s1(入口)、s2(中间)承载转发;s3(出口)不承载
	if !set[1] || !set[2] || set[3] || len(ids) != 2 {
		t.Fatalf("期望 {1,2},得到 %v", ids)
	}
}

func TestBuildForwardRulesErrors(t *testing.T) {
	servers := map[int64]*storage.RemoteServer{1: {ID: 1, IPAddress: "10.0.0.1"}}
	groups := map[int64]*storage.ForwardGroup{
		10: {ID: 10, Members: []storage.ForwardGroupMember{{ServerID: 1}}},
	}

	// 少于两跳
	single := &storage.ForwardChain{ID: 1, Hops: []storage.ForwardChainHop{{Seq: 0, GroupID: 10}}}
	if _, err := buildForwardRules(single, groups, servers, 10000, "tcp"); err == nil {
		t.Fatal("单跳链应报错")
	}

	// 非法端口
	valid := &storage.ForwardChain{ID: 1, Hops: []storage.ForwardChainHop{{Seq: 0, GroupID: 10}, {Seq: 1, GroupID: 10}}}
	if _, err := buildForwardRules(valid, groups, servers, 0, "tcp"); err == nil {
		t.Fatal("端口 0 应报错")
	}

	// 目标组成员无地址：官方仍下发能下发的跳，这一跳被跳过而不是整链失败
	groups[20] = &storage.ForwardGroup{ID: 20, Members: []storage.ForwardGroupMember{{ServerID: 99}}}
	badChain := &storage.ForwardChain{ID: 1, Hops: []storage.ForwardChainHop{{Seq: 0, GroupID: 10}, {Seq: 1, GroupID: 20}}}
	rules, err := buildForwardRules(badChain, groups, servers, 10000, "tcp")
	if err != nil {
		t.Fatalf("空出口组不应让整条链失败: %v", err)
	}
	if len(rules) != 0 {
		t.Fatalf("无可用上游时应不下发规则,得到 %+v", rules)
	}
}

func TestBuildForwardRulesSkipsEmptyDestHop(t *testing.T) {
	servers := map[int64]*storage.RemoteServer{
		1: {ID: 1, IPAddress: "10.0.0.1"},
		2: {ID: 2, IPAddress: "10.0.0.2"},
	}
	groups := map[int64]*storage.ForwardGroup{
		10: {ID: 10, Name: "A", BalanceStrategy: "round_robin", Members: []storage.ForwardGroupMember{{ServerID: 1}}},
		20: {ID: 20, Name: "B", BalanceStrategy: "round_robin", Members: []storage.ForwardGroupMember{{ServerID: 2}}},
		30: {ID: 30, Name: "C", BalanceStrategy: "round_robin"},
	}
	chain := &storage.ForwardChain{ID: 1, Hops: []storage.ForwardChainHop{
		{Seq: 0, GroupID: 10}, {Seq: 1, GroupID: 20}, {Seq: 2, GroupID: 30},
	}}
	rules, err := buildForwardRules(chain, groups, servers, 10000, "tcp")
	if err != nil {
		t.Fatalf("空出口组不应失败: %v", err)
	}
	if len(rules[1]) != 1 || rules[1][0].Upstreams[0].Addr != "10.0.0.2:10000" {
		t.Fatalf("入口→中转规则应仍下发: %+v", rules)
	}
}
