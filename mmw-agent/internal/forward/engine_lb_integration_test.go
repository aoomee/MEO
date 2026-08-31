package forward

import (
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

// 端到端验证负载均衡:真起 Manager + 真 TCP,走完整 accept→picker→dial→relay 链路,
// 断言 least_conn / sticky / round_robin 的实际分发结果。不依赖面板或多机。

// upstreamServer 每来一条连接就先写一个标识字节(区分是哪台),再把连接挂住(读到 EOF 才退),
// 这样 least_conn 的活跃计数能累积。返回监听地址和一个关闭函数。
func startUpstream(t *testing.T, id byte) (addr string, closeFn func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen upstream: %v", err)
	}
	var wg sync.WaitGroup
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			wg.Add(1)
			go func(c net.Conn) {
				defer wg.Done()
				defer c.Close()
				_, _ = c.Write([]byte{id})    // 告诉下游我是谁
				_, _ = io.Copy(io.Discard, c) // 挂住直到对端关闭
			}(c)
		}
	}()
	return ln.Addr().String(), func() { _ = ln.Close(); wg.Wait() }
}

func freePort(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("free port: %v", err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	return addr
}

// dialAndReadID 连到转发监听口,读回 1 字节上游标识(证明 relay 已建立)。
// keepOpen=true 时不关闭连接(留给 least_conn 累积活跃计数),返回的 conn 由调用方收尾。
func dialAndReadID(t *testing.T, listen string, keepOpen bool) (byte, net.Conn) {
	t.Helper()
	c, err := net.DialTimeout("tcp", listen, 2*time.Second)
	if err != nil {
		t.Fatalf("dial listen: %v", err)
	}
	_ = c.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 1)
	if _, err := io.ReadFull(c, buf); err != nil {
		c.Close()
		t.Fatalf("read upstream id: %v", err)
	}
	_ = c.SetReadDeadline(time.Time{})
	if !keepOpen {
		c.Close()
	}
	return buf[0], c
}

func applyRule(t *testing.T, strategy string) (*Manager, string, func()) {
	t.Helper()
	a, ca := startUpstream(t, 'A')
	b, cb := startUpstream(t, 'B')
	d, cd := startUpstream(t, 'C')
	listen := freePort(t)
	m := NewManager()
	rule := Rule{
		ID:       "lb-test",
		Listen:   listen,
		Protocol: "tcp",
		Strategy: strategy,
		Upstreams: []Upstream{
			{Addr: a, Weight: 1}, {Addr: b, Weight: 1}, {Addr: d, Weight: 1},
		},
		Health: Health{Enabled: false}, // 关健康探测 → 三台恒健康,分发确定
	}
	if err := m.Apply([]Rule{rule}); err != nil {
		ca(); cb(); cd()
		t.Fatalf("apply: %v", err)
	}
	time.Sleep(150 * time.Millisecond) // 等监听起来
	cleanup := func() { _ = m.Apply(nil); ca(); cb(); cd() }
	return m, listen, cleanup
}

func TestForwardLeastConnDistributes(t *testing.T) {
	_, listen, cleanup := applyRule(t, "least_conn")
	defer cleanup()

	// 顺序开 9 条并全部挂住 → 每次都落到当前活跃最少的一台 → A/B/C 各 3
	counts := map[byte]int{}
	var held []net.Conn
	for i := 0; i < 9; i++ {
		id, c := dialAndReadID(t, listen, true)
		counts[id]++
		held = append(held, c)
	}
	for _, c := range held {
		c.Close()
	}
	if counts['A'] != 3 || counts['B'] != 3 || counts['C'] != 3 {
		t.Fatalf("least_conn 分发不均: A=%d B=%d C=%d (期望各 3)", counts['A'], counts['B'], counts['C'])
	}
	t.Logf("least_conn OK: A=%d B=%d C=%d", counts['A'], counts['B'], counts['C'])
}

func TestForwardStickySameSource(t *testing.T) {
	_, listen, cleanup := applyRule(t, "sticky")
	defer cleanup()

	// 同一来源(都从 127.0.0.1)反复连 → 必须恒定落到同一台
	first, _ := dialAndReadID(t, listen, false)
	for i := 0; i < 10; i++ {
		id, _ := dialAndReadID(t, listen, false)
		if id != first {
			t.Fatalf("sticky 不稳定: 期望恒为 %c, 第 %d 次得到 %c", first, i, id)
		}
	}
	t.Logf("sticky OK: 同源恒定落到 %c", first)
}

func TestForwardRoundRobinSpreads(t *testing.T) {
	_, listen, cleanup := applyRule(t, "round_robin")
	defer cleanup()

	counts := map[byte]int{}
	for i := 0; i < 9; i++ {
		id, _ := dialAndReadID(t, listen, false)
		counts[id]++
		time.Sleep(10 * time.Millisecond) // 让每条 relay 建立,pick 计数推进
	}
	// round_robin 9 条应大致 3/3/3,至少三台都用到
	if counts['A'] == 0 || counts['B'] == 0 || counts['C'] == 0 {
		t.Fatalf("round_robin 有台没被用到: A=%d B=%d C=%d", counts['A'], counts['B'], counts['C'])
	}
	for _, id := range []byte{'A', 'B', 'C'} {
		if counts[id] < 2 || counts[id] > 4 {
			t.Fatalf("round_robin 分发偏斜: A=%d B=%d C=%d", counts['A'], counts['B'], counts['C'])
		}
	}
	t.Logf("round_robin OK: A=%d B=%d C=%d", counts['A'], counts['B'], counts['C'])
}

var _ = fmt.Sprintf
