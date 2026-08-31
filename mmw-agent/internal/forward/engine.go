package forward

import (
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type runningRule struct {
	rule   Rule
	health Health
	picker Picker

	ln     net.Listener
	stopCh chan struct{}
	wg     sync.WaitGroup

	mu     sync.Mutex
	probes map[string]*probeState
}

type probeState struct {
	healthy   bool
	rawOK     bool
	since     time.Time
	rttMs     int64
	bytesUp   atomic.Uint64
	bytesDown atomic.Uint64
}

func newRunningRule(rule Rule) (*runningRule, error) {
	h := rule.Health.normalized()
	rr := &runningRule{
		rule:   rule,
		health: h,
		picker: newPicker(rule.Strategy, rule.Upstreams),
		stopCh: make(chan struct{}),
		probes: make(map[string]*probeState, len(rule.Upstreams)),
	}
	now := time.Now()
	for _, u := range rule.Upstreams {
		if u.Addr == "" {
			continue
		}
		rr.probes[u.Addr] = &probeState{healthy: true, rawOK: true, since: now}
	}
	ln, err := net.Listen("tcp", rule.Listen)
	if err != nil {
		return nil, err
	}
	rr.ln = ln
	rr.wg.Add(1)
	go rr.acceptLoop()
	if h.Enabled {
		rr.wg.Add(1)
		go rr.healthLoop()
	}
	return rr, nil
}

func (rr *runningRule) stop() {
	close(rr.stopCh)
	if rr.ln != nil {
		_ = rr.ln.Close()
	}
	rr.wg.Wait()
}

func (rr *runningRule) acceptLoop() {
	defer rr.wg.Done()
	for {
		conn, err := rr.ln.Accept()
		if err != nil {
			select {
			case <-rr.stopCh:
				return
			default:
				log.Printf("[forward] %s accept: %v", rr.rule.ID, err)
				time.Sleep(50 * time.Millisecond)
				continue
			}
		}
		rr.wg.Add(1)
		go func(c net.Conn) {
			defer rr.wg.Done()
			rr.handleConn(c)
		}(conn)
	}
}

func (rr *runningRule) handleConn(client net.Conn) {
	defer client.Close()
	addr, release, ok := rr.picker.Pick(rr.isHealthy, clientHost(client))
	if !ok {
		return
	}
	defer release()
	timeout := time.Duration(rr.health.TimeoutMs) * time.Millisecond
	if timeout <= 0 {
		timeout = time.Second
	}
	dialer := net.Dialer{Timeout: timeout}
	up, err := dialer.Dial("tcp", addr)
	if err != nil {
		rr.markDialFail(addr)
		return
	}
	defer up.Close()

	st := rr.probeOf(addr)
	var upN, downN atomic.Uint64
	done := make(chan struct{}, 2)
	go func() {
		n, _ := io.Copy(countWriter{w: up, n: &upN}, client)
		_ = n
		_ = closeWrite(up)
		done <- struct{}{}
	}()
	go func() {
		n, _ := io.Copy(countWriter{w: client, n: &downN}, up)
		_ = n
		_ = closeWrite(client)
		done <- struct{}{}
	}()
	<-done
	<-done
	if st != nil {
		st.bytesUp.Add(upN.Load())
		st.bytesDown.Add(downN.Load())
	}
}

func closeWrite(c net.Conn) error {
	if tc, ok := c.(*net.TCPConn); ok {
		return tc.CloseWrite()
	}
	return nil
}

// clientHost 取客户端来源 IP(去掉端口),给 sticky 源哈希用。取不到就返回空串。
func clientHost(c net.Conn) string {
	if c == nil || c.RemoteAddr() == nil {
		return ""
	}
	host, _, err := net.SplitHostPort(c.RemoteAddr().String())
	if err != nil {
		return c.RemoteAddr().String()
	}
	return host
}

type countWriter struct {
	w io.Writer
	n *atomic.Uint64
}

func (c countWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	if n > 0 {
		c.n.Add(uint64(n))
	}
	return n, err
}

func (rr *runningRule) healthLoop() {
	defer rr.wg.Done()
	iv := time.Duration(rr.health.IntervalMs) * time.Millisecond
	t := time.NewTicker(iv)
	defer t.Stop()
	rr.probeAll()
	for {
		select {
		case <-rr.stopCh:
			return
		case <-t.C:
			rr.probeAll()
		}
	}
}

func (rr *runningRule) probeAll() {
	timeout := time.Duration(rr.health.TimeoutMs) * time.Millisecond
	for _, u := range rr.rule.Upstreams {
		if u.Addr == "" {
			continue
		}
		ok, rtt := tcpProbe(u.Addr, timeout)
		if ok && rr.health.RTTThresholdMs > 0 && rtt > int64(rr.health.RTTThresholdMs) {
			ok = false
		}
		rr.applyProbe(u.Addr, ok, rtt)
	}
}

func tcpProbe(addr string, timeout time.Duration) (bool, int64) {
	start := time.Now()
	c, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false, 0
	}
	_ = c.Close()
	return true, time.Since(start).Milliseconds()
}

func (rr *runningRule) applyProbe(addr string, ok bool, rtt int64) {
	now := time.Now()
	rr.mu.Lock()
	defer rr.mu.Unlock()
	st := rr.probes[addr]
	if st == nil {
		return
	}
	st.rttMs = rtt
	if st.rawOK != ok {
		st.rawOK = ok
		st.since = now
	}
	if !rr.health.Enabled {
		st.healthy = true
		return
	}
	elapsed := now.Sub(st.since)
	if ok {
		if !st.healthy && elapsed >= time.Duration(rr.health.RecoverMs)*time.Millisecond {
			st.healthy = true
		}
	} else if st.healthy && elapsed >= time.Duration(rr.health.FailoverMs)*time.Millisecond {
		st.healthy = false
	}
}

func (rr *runningRule) markDialFail(addr string) {
	rr.applyProbe(addr, false, 0)
}

func (rr *runningRule) isHealthy(addr string) bool {
	if !rr.health.Enabled {
		return true
	}
	rr.mu.Lock()
	defer rr.mu.Unlock()
	st := rr.probes[addr]
	return st != nil && st.healthy
}

func (rr *runningRule) probeOf(addr string) *probeState {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	return rr.probes[addr]
}

func (rr *runningRule) status() RuleStatus {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	out := RuleStatus{RuleID: rr.rule.ID, Listen: rr.rule.Listen}
	for _, u := range rr.rule.Upstreams {
		st := rr.probes[u.Addr]
		us := UpstreamStatus{Addr: u.Addr, Healthy: true}
		if st != nil {
			us.Healthy = st.healthy
			us.RTTMs = st.rttMs
			us.BytesUp = st.bytesUp.Load()
			us.BytesDown = st.bytesDown.Load()
		}
		out.Upstreams = append(out.Upstreams, us)
	}
	return out
}
