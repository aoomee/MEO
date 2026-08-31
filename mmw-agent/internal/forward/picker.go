package forward

import (
	"hash/fnv"
	"sort"
	"sync"
	"sync/atomic"
)

// Picker 在健康上游之间做负载均衡。
//
// Pick 返回选中的 addr、一个连接结束时调用的 release(least_conn 用来减活跃计数,其它策略为
// 空操作),以及是否选到。srcIP 是客户端来源 IP(sticky 源哈希用,其它策略忽略)。
//
// 支持的 strategy(与面板 mapForwardStrategy 对齐):
//   round_robin —— 顺序轮流(默认)
//   weighted    —— 平滑加权轮询(权重由面板填入 Upstream.Weight;percentage/cycle 也走它)
//   least_conn  —— 最少活跃连接
//   sticky      —— 源 IP 哈希,同一来源恒定落到同一台
type Picker interface {
	Pick(healthy func(addr string) bool, srcIP string) (addr string, release func(), ok bool)
}

func newPicker(strategy string, ups []Upstream) Picker {
	switch strategy {
	case "weighted":
		return newWeightedPicker(ups)
	case "least_conn":
		return newLeastConnPicker(ups)
	case "sticky":
		return newStickyPicker(ups)
	default: // round_robin 及未知策略都退回顺序轮流(对旧 agent 也是安全默认)
		return newRoundRobinPicker(ups)
	}
}

func noopRelease() {}

// addrsOf 抽出非空 upstream 地址(保持下发顺序)。
func addrsOf(ups []Upstream) []string {
	addrs := make([]string, 0, len(ups))
	for _, u := range ups {
		if u.Addr != "" {
			addrs = append(addrs, u.Addr)
		}
	}
	return addrs
}

// ---------------- round_robin ----------------

type rrPicker struct {
	addrs []string
	n     atomic.Uint64
}

func newRoundRobinPicker(ups []Upstream) *rrPicker {
	return &rrPicker{addrs: addrsOf(ups)}
}

func (p *rrPicker) Pick(healthy func(string) bool, _ string) (string, func(), bool) {
	n := len(p.addrs)
	if n == 0 {
		return "", noopRelease, false
	}
	start := int(p.n.Add(1) % uint64(n))
	for i := 0; i < n; i++ {
		addr := p.addrs[(start+i)%n]
		if healthy == nil || healthy(addr) {
			return addr, noopRelease, true
		}
	}
	return p.addrs[start], noopRelease, true
}

// ---------------- weighted (平滑加权轮询) ----------------

type wrrPicker struct {
	mu    sync.Mutex
	items []wrrItem
	total int
}

type wrrItem struct {
	addr    string
	weight  int
	current int
}

func newWeightedPicker(ups []Upstream) *wrrPicker {
	p := &wrrPicker{}
	for _, u := range ups {
		if u.Addr == "" {
			continue
		}
		w := u.Weight
		if w <= 0 {
			w = 1
		}
		p.items = append(p.items, wrrItem{addr: u.Addr, weight: w})
		p.total += w
	}
	return p
}

func (p *wrrPicker) Pick(healthy func(string) bool, _ string) (string, func(), bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.items) == 0 || p.total <= 0 {
		return "", noopRelease, false
	}
	var best *wrrItem
	for i := range p.items {
		it := &p.items[i]
		if healthy != nil && !healthy(it.addr) {
			continue
		}
		it.current += it.weight
		if best == nil || it.current > best.current {
			best = it
		}
	}
	if best == nil {
		best = &p.items[0]
		best.current += best.weight
	}
	best.current -= p.total
	return best.addr, noopRelease, true
}

// ---------------- least_conn (最少活跃连接) ----------------

type lcPicker struct {
	addrs  []string
	active []atomic.Int64 // 与 addrs 一一对应,当前活跃连接数
}

func newLeastConnPicker(ups []Upstream) *lcPicker {
	addrs := addrsOf(ups)
	return &lcPicker{addrs: addrs, active: make([]atomic.Int64, len(addrs))}
}

func (p *lcPicker) Pick(healthy func(string) bool, _ string) (string, func(), bool) {
	n := len(p.addrs)
	if n == 0 {
		return "", noopRelease, false
	}
	best := -1
	var bestCnt int64
	for i := 0; i < n; i++ {
		if healthy != nil && !healthy(p.addrs[i]) {
			continue
		}
		c := p.active[i].Load()
		if best < 0 || c < bestCnt {
			best, bestCnt = i, c
		}
	}
	if best < 0 { // 全不健康,退回第一个,避免整链断流
		best = 0
	}
	p.active[best].Add(1)
	idx := best
	return p.addrs[idx], func() { p.active[idx].Add(-1) }, true
}

// ---------------- sticky (源 IP 哈希) ----------------

type stickyPicker struct {
	addrs []string
}

func newStickyPicker(ups []Upstream) *stickyPicker {
	addrs := addrsOf(ups)
	// 排序保证同一组上游在不同 agent / 重启后哈希落点一致。
	sort.Strings(addrs)
	return &stickyPicker{addrs: addrs}
}

func (p *stickyPicker) Pick(healthy func(string) bool, srcIP string) (string, func(), bool) {
	n := len(p.addrs)
	if n == 0 {
		return "", noopRelease, false
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(srcIP))
	start := int(h.Sum32() % uint32(n))
	for i := 0; i < n; i++ { // 从哈希点起找第一个健康的,保持源黏性又能故障转移
		addr := p.addrs[(start+i)%n]
		if healthy == nil || healthy(addr) {
			return addr, noopRelease, true
		}
	}
	return p.addrs[start], noopRelease, true
}
