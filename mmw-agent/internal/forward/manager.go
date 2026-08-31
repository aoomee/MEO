package forward

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

// Manager 持有本机全部转发规则。Apply 是整机全量替换。
//
// persistPath 非空时,每次 Apply 成功都把「已下发的规则集」原子落盘;Restore 在 agent
// 启动/升级重启后读回并重建 relay,做到「重启不断转发」(对齐官方 0.6.6)。
type Manager struct {
	mu          sync.Mutex
	rules       map[string]*runningRule
	persistPath string
}

func NewManager() *Manager {
	return &Manager{rules: map[string]*runningRule{}}
}

// SetPersistPath 设置转发规则落盘路径。空则不持久化。
func (m *Manager) SetPersistPath(p string) {
	m.mu.Lock()
	m.persistPath = p
	m.mu.Unlock()
}

// persistLocked 原子写入当前下发的规则集。调用方须持有 m.mu。
func (m *Manager) persistLocked(rules []Rule) {
	if m.persistPath == "" {
		return
	}
	data, err := json.Marshal(rules)
	if err != nil {
		log.Printf("[forward] persist marshal failed: %v", err)
		return
	}
	tmp := m.persistPath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		log.Printf("[forward] persist write failed: %v", err)
		return
	}
	if err := os.Rename(tmp, m.persistPath); err != nil {
		log.Printf("[forward] persist rename failed: %v", err)
		_ = os.Remove(tmp)
	}
}

// Restore 读回落盘的规则并重建 relay。启动时调用一次。规则文件不存在视作无规则。
func (m *Manager) Restore() error {
	m.mu.Lock()
	path := m.persistPath
	m.mu.Unlock()
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var rules []Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return fmt.Errorf("forward: 解析持久化规则失败: %w", err)
	}
	if len(rules) == 0 {
		return nil
	}
	log.Printf("[forward] restoring %d rule(s) from %s", len(rules), path)
	return m.Apply(rules)
}

func (m *Manager) Apply(rules []Rule) error {
	seenListen := map[string]string{}
	cleaned := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if r.ID == "" {
			return fmt.Errorf("forward: 规则缺少 id")
		}
		if r.Listen == "" {
			return fmt.Errorf("forward: 规则 %s 缺少 listen", r.ID)
		}
		if proto := r.Protocol; proto != "" && proto != "tcp" {
			return fmt.Errorf("forward: 规则 %s 不支持协议 %s", r.ID, proto)
		}
		if prev, ok := seenListen[r.Listen]; ok && prev != r.ID {
			return fmt.Errorf("forward: listen %s 被规则 %s 与 %s 共用", r.Listen, prev, r.ID)
		}
		seenListen[r.Listen] = r.ID
		if len(r.Upstreams) == 0 {
			return fmt.Errorf("forward: 规则 %s 没有 upstream", r.ID)
		}
		cleaned = append(cleaned, r)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	want := map[string]Rule{}
	for _, r := range cleaned {
		want[r.ID] = r
	}
	for id, running := range m.rules {
		if _, ok := want[id]; !ok {
			running.stop()
			delete(m.rules, id)
			log.Printf("[forward] removed rule %s", id)
		}
	}
	var firstErr error
	for _, r := range cleaned {
		if cur, ok := m.rules[r.ID]; ok && sameRule(cur.rule, r) {
			continue
		}
		if cur, ok := m.rules[r.ID]; ok {
			cur.stop()
			delete(m.rules, r.ID)
		}
		nr, err := newRunningRule(r)
		if err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("forward: 启动规则 %s 失败: %w", r.ID, err)
			}
			log.Printf("[forward] start %s failed: %v", r.ID, err)
			continue
		}
		m.rules[r.ID] = nr
		log.Printf("[forward] applied rule %s listen=%s upstreams=%d", r.ID, r.Listen, len(r.Upstreams))
	}
	// 校验通过即落盘「意图规则集」——即使个别 relay 因端口占用等瞬时原因启动失败,
	// 重启后 Restore 仍会重试,不至于永久丢转发。
	m.persistLocked(cleaned)
	return firstErr
}

func sameRule(a, b Rule) bool {
	if a.ID != b.ID || a.Listen != b.Listen || a.Protocol != b.Protocol || a.Strategy != b.Strategy {
		return false
	}
	if a.Health != b.Health || len(a.Upstreams) != len(b.Upstreams) {
		return false
	}
	for i := range a.Upstreams {
		if a.Upstreams[i] != b.Upstreams[i] {
			return false
		}
	}
	return true
}

func (m *Manager) Status() []RuleStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]RuleStatus, 0, len(m.rules))
	for _, rr := range m.rules {
		out = append(out, rr.status())
	}
	return out
}
