package guardclient

import "testing"

func TestEnforcedIsOptional(t *testing.T) {
	c := NewFromEnv()
	if Enforced() || c.Enforced() {
		t.Fatal("happy-edition stub must report Enforced()=false so upgrade skips Guard")
	}
	if !c.Enabled() {
		t.Fatal("Enabled() still true: local stub is present")
	}
}
