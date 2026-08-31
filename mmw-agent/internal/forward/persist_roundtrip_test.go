package forward

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPersistRestoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	pf := filepath.Join(dir, "forward-rules.json")

	m1 := NewManager()
	m1.SetPersistPath(pf)
	rule := Rule{
		ID:        "r1",
		Listen:    "127.0.0.1:0", // 0 = 内核分配空闲端口,重建时再分配一个新的,不会冲突
		Protocol:  "tcp",
		Upstreams: []Upstream{{Addr: "127.0.0.1:65000", Weight: 1}},
		Health:    Health{Enabled: false},
	}
	if err := m1.Apply([]Rule{rule}); err != nil {
		t.Fatalf("apply: %v", err)
	}
	data, err := os.ReadFile(pf)
	if err != nil || len(data) == 0 || string(data) == "null" {
		t.Fatalf("persist file bad: err=%v data=%q", err, data)
	}
	t.Cleanup(func() { m1.Apply(nil) })

	// 新 Manager 从盘恢复(模拟 agent 重启)
	m2 := NewManager()
	m2.SetPersistPath(pf)
	if err := m2.Restore(); err != nil {
		t.Fatalf("restore: %v", err)
	}
	st := m2.Status()
	if len(st) != 1 || st[0].RuleID != "r1" {
		t.Fatalf("restored status wrong: %+v", st)
	}
	m2.Apply(nil)
}

func TestClearPersistsEmpty(t *testing.T) {
	dir := t.TempDir()
	pf := filepath.Join(dir, "forward-rules.json")
	m := NewManager()
	m.SetPersistPath(pf)
	if err := m.Apply([]Rule{{
		ID: "r1", Listen: "127.0.0.1:0", Protocol: "tcp",
		Upstreams: []Upstream{{Addr: "127.0.0.1:65000", Weight: 1}},
	}}); err != nil {
		t.Fatalf("apply: %v", err)
	}
	// 主控清空规则 → 落盘 [] → 重启后不应重建
	m.Apply(nil)
	m2 := NewManager()
	m2.SetPersistPath(pf)
	if err := m2.Restore(); err != nil {
		t.Fatalf("restore: %v", err)
	}
	if n := len(m2.Status()); n != 0 {
		t.Fatalf("expected 0 rules after clear-restore, got %d", n)
	}
}
