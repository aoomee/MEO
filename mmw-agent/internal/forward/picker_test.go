package forward

import "testing"

func TestRoundRobinSkipsUnhealthy(t *testing.T) {
	p := newRoundRobinPicker([]Upstream{{Addr: "a:1"}, {Addr: "b:1"}, {Addr: "c:1"}})
	got := map[string]int{}
	for i := 0; i < 6; i++ {
		addr, release, ok := p.Pick(func(a string) bool { return a != "b:1" }, "")
		if !ok {
			t.Fatal("pick failed")
		}
		release()
		got[addr]++
	}
	if got["b:1"] != 0 || got["a:1"] == 0 || got["c:1"] == 0 {
		t.Fatalf("unexpected picks: %v", got)
	}
}

func TestLeastConnPrefersIdle(t *testing.T) {
	p := newLeastConnPicker([]Upstream{{Addr: "a:1"}, {Addr: "b:1"}, {Addr: "c:1"}})
	// 占住 a:1 两条、b:1 一条,不释放 → 下一次应选空闲的 c:1
	a, ra, _ := p.Pick(nil, "")
	a2, ra2, _ := p.Pick(nil, "")
	if a != "a:1" || a2 != "b:1" { // 初始全 0,按顺序取到 a、b
		t.Fatalf("warmup picks: %v %v", a, a2)
	}
	_ = ra
	_ = ra2
	// 再 Pick:a=1,b=1,c=0 → 选 c
	addr, rc, ok := p.Pick(nil, "")
	if !ok || addr != "c:1" {
		t.Fatalf("least_conn should pick idle c:1, got %v", addr)
	}
	rc()
	// 释放 a 的一条后 a=0,b=1,c=0 → 选 a(索引更小)
	ra()
	addr2, r2, _ := p.Pick(nil, "")
	if addr2 != "a:1" {
		t.Fatalf("after release, least_conn should pick a:1, got %v", addr2)
	}
	r2()
	ra2()
}

func TestLeastConnSkipsUnhealthy(t *testing.T) {
	p := newLeastConnPicker([]Upstream{{Addr: "a:1"}, {Addr: "b:1"}})
	addr, rel, ok := p.Pick(func(a string) bool { return a == "b:1" }, "")
	if !ok || addr != "b:1" {
		t.Fatalf("should pick only healthy b:1, got %v", addr)
	}
	rel()
}

func TestStickyConsistentPerSource(t *testing.T) {
	p := newStickyPicker([]Upstream{{Addr: "a:1"}, {Addr: "b:1"}, {Addr: "c:1"}})
	first, r, ok := p.Pick(nil, "203.0.113.7")
	if !ok {
		t.Fatal("pick failed")
	}
	r()
	for i := 0; i < 20; i++ { // 同一来源 IP 必须恒定落到同一台
		addr, rr, _ := p.Pick(nil, "203.0.113.7")
		rr()
		if addr != first {
			t.Fatalf("sticky not consistent: %v vs %v", addr, first)
		}
	}
	// 选中的那台不健康时应故障转移到别的健康台
	down := first
	addr, rr, _ := p.Pick(func(a string) bool { return a != down }, "203.0.113.7")
	rr()
	if addr == down {
		t.Fatalf("sticky should failover off unhealthy %v", down)
	}
}

func TestHealthNormalizedDefaults(t *testing.T) {
	h := Health{}.normalized()
	if h.IntervalMs <= 0 || h.TimeoutMs <= 0 || h.FailoverMs <= 0 || h.RecoverMs <= 0 {
		t.Fatalf("defaults not applied: %+v", h)
	}
}
