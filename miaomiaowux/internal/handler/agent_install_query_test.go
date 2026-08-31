package handler

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"miaomiaowux/internal/storage"
)

func TestGenerateSecureTokenUsesRawURLEncoding(t *testing.T) {
	oldRead, oldEncode := randRead, base64Encode
	t.Cleanup(func() { randRead, base64Encode = oldRead, oldEncode })
	randRead = func(p []byte) (int, error) {
		for i := range p {
			p[i] = byte(i)
		}
		return len(p), nil
	}
	base64Encode = func(p []byte) string { return base64.RawURLEncoding.EncodeToString(p) }

	token, err := generateSecureToken()
	if err != nil {
		t.Fatalf("generateSecureToken: %v", err)
	}
	if strings.Contains(token, "=") {
		t.Fatalf("new URL-safe token contains padding: %q", token)
	}
	if !validInstallToken(token) {
		t.Fatalf("generated token rejected by installer validator")
	}
}

func TestBuildRemoteInstallQueryKeepsLegacyPaddingReadable(t *testing.T) {
	query := buildRemoteInstallQuery("legacy_token=", url.Values{
		"xray_mode":   {"embedded"},
		"listen_port": {"12345"},
	})
	if !strings.HasPrefix(query, "token=legacy_token=&") {
		t.Fatalf("legacy padding was percent encoded: %q", query)
	}
	values, err := url.ParseQuery(query)
	if err != nil {
		t.Fatalf("ParseQuery: %v", err)
	}
	if got := values.Get("token"); got != "legacy_token=" {
		t.Fatalf("token round trip = %q", got)
	}
}

func TestRemoteInstallScriptSupportsNoInitGuardSupervisor(t *testing.T) {
	// 本 fork 已去 Guard(免费版 AgentGuardRequired()==false),安装脚本不再生成 Guard supervisor /
	// mmwx-guard-agent 相关步骤。这是上游「装 Guard」的测试,对去 Guard 的我们不适用。
	t.Skip("de-Guard fork:安装脚本不含 Guard supervisor,此上游测试不适用")
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "install-script.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })

	h := NewXrayServerHandler(repo, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/remote/install?token=test-token", nil)
	req.Host = "panel.example.com"
	rec := httptest.NewRecorder()
	h.GetRemoteInstallScript(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}

	script := rec.Body.String()
	guardSupervisor := "nohup /usr/local/bin/mmwx-guard-agent-supervisor.sh"
	agentSupervisor := "nohup /usr/local/bin/mmw-agent-supervisor.sh"
	if !strings.Contains(script, guardSupervisor) || !strings.Contains(script, agentSupervisor) {
		t.Fatal("no-init Guard/Agent supervisors are missing")
	}
	if !strings.Contains(script, "sed -i '/mmw-agent-supervisor.sh/i nohup /usr/local/bin/mmwx-guard-agent-supervisor.sh") {
		t.Fatal("rc.local does not insert Guard before the Agent supervisor")
	}
	if !strings.Contains(script, "while [ ! -S \"$MMWX_GUARD_SOCKET\" ]; do sleep 1; done") {
		t.Fatal("Agent supervisor does not wait for the Guard socket")
	}

	path := filepath.Join(t.TempDir(), "install.sh")
	if err := os.WriteFile(path, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command("bash", "-n", path).CombinedOutput(); err != nil {
		t.Fatalf("generated installer has invalid shell syntax: %v: %s", err, output)
	}
}

func TestRemoteInstallScriptVerifiesReleaseTupleBeforeAtomicInstall(t *testing.T) {
	// 同上:去 Guard 后安装脚本不含 .sig/.manifest 下载、GUARD_VERIFY 签名校验、mmwx-guard-agent 步骤,
	// 此上游测试断言的正是这批 Guard 内容,对去 Guard 的我们不适用。
	t.Skip("de-Guard fork:安装脚本不含 Guard 签名校验步骤,此上游测试不适用")
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "install-script.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = repo.Close() })

	h := NewXrayServerHandler(repo, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/remote/install?token=test-token", nil)
	req.Host = "panel.example.com"
	rec := httptest.NewRecorder()
	h.GetRemoteInstallScript(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	script := rec.Body.String()

	checks := []string{
		`download_file "$agent_url" "$AGENT_NEW" 180`,
		`download_file "${agent_url}.sig" "$AGENT_SIG" 60`,
		`download_file "${agent_url}.manifest" "$MANIFEST_NEW" 60`,
		`openssl pkeyutl -verify -pubin`,
		`"$GUARD_VERIFY" --role agent --manifest "$MANIFEST_NEW" --verify-manifest-for "$AGENT_NEW"`,
		`install -m 0755 "$AGENT_NEW" "$AGENT_BIN.new"`,
		`rollback_install_transaction`,
		`systemctl restart mmwx-guard-agent`,
		`restore_managed_services`,
		`restore_previous_config`,
		`trap early_install_cleanup EXIT`,
		`TRANSACTION_ROLLED_BACK=1`,
		`INSTALL_SUCCEEDED=1`,
		`backup_managed_file /etc/rc.local rc-local`,
		`restore_managed_file /etc/rc.local rc-local`,
		`systemctl disable mmw-agent mmwx-guard-agent`,
		`rc-update del mmw-agent default`,
		`GUARD_VERIFY=$(mktemp /usr/local/bin/.mmwx-guard-verify.XXXXXX)`,
		`install -m 0755 "$GUARD_NEW" "$GUARD_VERIFY"`,
		`verify_update_file "$GUARD_VERIFY" "$GUARD_SIG"`,
		`break 2`,
		`[ -S /run/mmwx-guard-agent/guard.sock ] && \
		   /usr/local/bin/mmw-agent __guard-health /run/mmwx-guard-agent/guard.sock`,
	}
	for _, want := range checks {
		if !strings.Contains(script, want) {
			t.Errorf("generated installer is missing %q", want)
		}
	}
	if strings.Contains(script, `systemctl enable --now mmwx-guard-agent`) {
		t.Error("generated installer does not explicitly restart an existing Guard")
	}
	step2AtForDisable := strings.Index(script, "# Step 2: Create config directory first")
	if step2AtForDisable < 0 {
		t.Fatal("generated installer is missing Step2")
	}
	preTransaction := script[:step2AtForDisable]
	if strings.Contains(preTransaction, `systemctl disable mmw-agent 2>/dev/null`) ||
		strings.Contains(preTransaction, `rc-update del mmw-agent 2>/dev/null`) {
		t.Error("generated installer disables the previous Agent while initially stopping it")
	}
	if strings.Contains(script, `-o /usr/local/share/mmwx-guard/agent.manifest "${agent_url}.manifest"`) {
		t.Error("generated installer downloads a manifest directly to its final path")
	}
	verifyAt := strings.Index(script, `"$GUARD_VERIFY" --role agent --manifest "$MANIFEST_NEW" --verify-manifest-for "$AGENT_NEW"`)
	installAt := strings.Index(script, `install -m 0755 "$AGENT_NEW" "$AGENT_BIN.new"`)
	if verifyAt < 0 || installAt < 0 || verifyAt >= installAt {
		t.Fatalf("release tuple is installed before exact-caller verification: verify=%d install=%d", verifyAt, installAt)
	}
	agentLoopAt := strings.Index(script, `for agent_url in "${MIRRORS[@]}"; do`)
	agentLoopEnd := -1
	if agentLoopAt >= 0 {
		agentLoopEnd = strings.Index(script[agentLoopAt:], `if [ "$download_ok" != "1" ]`)
	}
	if agentLoopAt < 0 || agentLoopEnd < 0 {
		t.Fatal("generated installer is missing the Agent mirror loop")
	}
	agentLoop := script[agentLoopAt : agentLoopAt+agentLoopEnd]
	if !strings.Contains(agentLoop, `--verify-manifest-for "$AGENT_NEW"`) ||
		!strings.Contains(agentLoop, `manifest 与该 Guard 不兼容`) ||
		!strings.Contains(agentLoop, `尝试下一 Guard 镜像`) ||
		!strings.Contains(agentLoop, "break 2") {
		t.Fatal("signature-valid but manifest-mismatched Agent tuple does not fall back to the next mirror")
	}
	restoreAt := strings.Index(script, "restore_previous_agent")
	successAt := strings.LastIndex(script, "INSTALL_SUCCEEDED=1")
	startAt := strings.Index(script, `systemctl start mmw-agent`)
	if restoreAt < 0 || startAt < 0 || successAt <= startAt {
		t.Fatalf("failed install cannot restore the previously active Agent, or commits before startup: restore=%d start=%d success=%d", restoreAt, startAt, successAt)
	}
	earlyTrapAt := strings.Index(script, "trap early_install_cleanup EXIT")
	step2At := strings.Index(script, "# Step 2: Create config directory first")
	if earlyTrapAt < 0 || step2At < 0 || earlyTrapAt >= step2At {
		t.Fatalf("early failure restore trap is not installed before Step2: trap=%d step2=%d", earlyTrapAt, step2At)
	}
	serviceBackupAt := strings.Index(script, `backup_managed_file /etc/systemd/system/mmw-agent.service systemd-agent`)
	serviceWriteAt := strings.Index(script, `cat > /etc/systemd/system/mmw-agent.service`)
	if serviceBackupAt < 0 || serviceWriteAt < 0 || serviceBackupAt >= serviceWriteAt {
		t.Fatalf("service state is not captured before it is overwritten: backup=%d write=%d", serviceBackupAt, serviceWriteAt)
	}
	rollbackAt := strings.Index(script, "rollback_install_transaction() {")
	rollbackEnd := -1
	if rollbackAt >= 0 {
		rollbackEnd = strings.Index(script[rollbackAt:], "\n\t}\n\t# Prepare all .new files")
	}
	if rollbackAt < 0 || rollbackEnd < 0 {
		t.Fatal("generated installer is missing the transaction rollback body")
	}
	rollbackBody := script[rollbackAt : rollbackAt+rollbackEnd]
	stopAgentAt := strings.Index(rollbackBody, "systemctl stop mmw-agent")
	stopGuardAt := strings.Index(rollbackBody, "systemctl stop mmwx-guard-agent")
	if stopAgentAt < 0 || stopGuardAt < 0 || stopAgentAt >= stopGuardAt {
		t.Fatalf("rollback does not stop Agent before restoring its executable: agent=%d guard=%d", stopAgentAt, stopGuardAt)
	}
}
