package handler

import "testing"

func TestPickAgentReleasePrefersSemverNotListOrder(t *testing.T) {
	rels := []ghRelease{
		{TagName: "mmwx-v0.4.8", Name: "v0.4.8"},
		{TagName: "mmwx-agent-v0.5.9", Name: "v0.5.9"},
		{TagName: "mmwx-agent-v0.6.0", Name: "mmwx-agent-v0.6.0"},
		{TagName: "mmwx-agent-v0.4.1", Name: "v0.4.1"},
	}
	got := pickAgentRelease(rels)
	if got == nil {
		t.Fatal("expected an agent release")
	}
	if v := agentReleaseVersion(got); v != "0.6.0" {
		t.Fatalf("want 0.6.0, got tag=%s name=%s ver=%s", got.TagName, got.Name, v)
	}
}

func TestAgentReleaseVersionFromTagWhenNameIsNotSemver(t *testing.T) {
	r := ghRelease{TagName: "mmwx-agent-v0.6.0", Name: "mmwx-agent-v0.6.0"}
	if v := agentReleaseVersion(&r); v != "0.6.0" {
		t.Fatalf("got %q", v)
	}
}

func TestIsUpgradeAvailableAgainstGithubLatest(t *testing.T) {
	if !isUpgradeAvailable("0.5.9", "0.6.0") {
		t.Fatal("0.5.9 should be upgradable to 0.6.0")
	}
	if isUpgradeAvailable("0.6.0", "0.6.0") {
		t.Fatal("same version should not be upgradable")
	}
}
