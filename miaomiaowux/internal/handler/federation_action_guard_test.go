package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"miaomiaowux/internal/guardclient"
	"miaomiaowux/internal/storage"
)

const federationTestShareToken = "federation-test-share-token"

func newFederationGuardTestRepo(t *testing.T) (*storage.TrafficRepository, storage.RemoteServer, storage.SharedServer) {
	t.Helper()
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "federation.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	server := storage.RemoteServer{Name: "owner-agent", Token: "owner-real-agent-token", Status: storage.RemoteServerStatusConnected}
	if err := repo.CreateRemoteServer(context.Background(), &server); err != nil {
		t.Fatalf("CreateRemoteServer: %v", err)
	}
	shareID, err := repo.CreateSharedServer(context.Background(), server.ID, hashShareToken(federationTestShareToken), "test", false)
	if err != nil {
		t.Fatalf("CreateSharedServer: %v", err)
	}
	share, err := repo.GetSharedServerByTokenHash(context.Background(), hashShareToken(federationTestShareToken))
	if err != nil {
		t.Fatalf("GetSharedServerByTokenHash: %v", err)
	}
	if share.ID != shareID {
		t.Fatalf("share ID = %d, want %d", share.ID, shareID)
	}
	return repo, server, share
}

func federationTestHash(seed byte) string {
	return strings.Repeat(string([]byte{seed}), 64)
}

func TestFederationGuardBindingBindsShareMethodPathAndPayload(t *testing.T) {
	base := federationActionBindingHash(10, 20, actionNodeCreate, http.MethodPost, federationInboundRequestPath, federationTestHash('a'))
	if !validFederationSHA256(base) {
		t.Fatalf("binding is not a canonical SHA-256: %q", base)
	}
	if again := federationActionBindingHash(10, 20, actionNodeCreate, http.MethodPost, federationInboundRequestPath, federationTestHash('a')); again != base {
		t.Fatalf("binding is not deterministic: %q != %q", again, base)
	}
	variants := []string{
		federationActionBindingHash(11, 20, actionNodeCreate, http.MethodPost, federationInboundRequestPath, federationTestHash('a')),
		federationActionBindingHash(10, 21, actionNodeCreate, http.MethodPost, federationInboundRequestPath, federationTestHash('a')),
		federationActionBindingHash(10, 20, actionNodeCreate, http.MethodPut, federationInboundRequestPath, federationTestHash('a')),
		federationActionBindingHash(10, 20, actionNodeCreate, http.MethodPost, federationNodeRequestPath, federationTestHash('a')),
		federationActionBindingHash(10, 20, actionNodeCreate, http.MethodPost, federationInboundRequestPath, federationTestHash('b')),
	}
	for i, variant := range variants {
		if variant == base {
			t.Fatalf("variant %d did not change binding", i)
		}
	}
	if validFederationActionTarget(http.MethodGet, federationInboundRequestPath) ||
		validFederationActionTarget(http.MethodPost, "/api/admin/system") ||
		validFederationSHA256("not-a-hash") {
		t.Fatal("invalid federation binding input was accepted")
	}
}

func issueFederationTestChallenge(
	t *testing.T,
	h *FederationHandler,
	serverID int64,
	payloadHash, requestPath string,
) guardclient.Challenge {
	t.Helper()
	body, _ := json.Marshal(FederationActionChallengeRequest{
		Action: actionNodeCreate, PayloadHash: payloadHash,
		RequestMethod: http.MethodPost, RequestPath: requestPath,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/federation/action-guard/challenge", bytes.NewReader(body))
	req.Header.Set("X-Share-Token", federationTestShareToken)
	recorder := httptest.NewRecorder()
	h.handleActionGuardChallenge(recorder, req, serverID)
	if recorder.Code != http.StatusOK {
		t.Fatalf("challenge status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	var challenge guardclient.Challenge
	if err := json.Unmarshal(recorder.Body.Bytes(), &challenge); err != nil {
		t.Fatalf("decode challenge: %v", err)
	}
	return challenge
}

func TestFederationOwnerChallengeReturnsBoundPayloadHash(t *testing.T) {
	repo, server, share := newFederationGuardTestRepo(t)
	h := NewFederationHandler(repo, nil, nil)
	originalHash := federationTestHash('a')
	var guardRequest guardclient.ChallengeRequest
	h.createChallenge = func(_ context.Context, request guardclient.ChallengeRequest) (guardclient.Challenge, error) {
		guardRequest = request
		return guardclient.Challenge{ChallengeID: "owner-challenge", Challenge: "random-proof", ExpiresAt: time.Now().Add(time.Minute).Unix()}, nil
	}
	challenge := issueFederationTestChallenge(t, h, server.ID, originalHash, federationInboundRequestPath)
	wantBinding := federationActionBindingHash(share.ID, server.ID, actionNodeCreate, http.MethodPost, federationInboundRequestPath, originalHash)
	if challenge.PayloadHash != wantBinding || guardRequest.PayloadHash != wantBinding || guardRequest.Action != actionNodeCreate {
		t.Fatalf("owner binding mismatch: response=%q guard=%+v want=%q", challenge.PayloadHash, guardRequest, wantBinding)
	}
	if challenge.PayloadHash == originalHash {
		t.Fatal("owner returned the unbound browser payload hash")
	}

	invalidBody, _ := json.Marshal(FederationActionChallengeRequest{
		Action: actionNodeCreate, PayloadHash: originalHash,
		RequestMethod: http.MethodPost, RequestPath: "/api/admin/system",
	})
	invalidReq := httptest.NewRequest(http.MethodPost, "/api/federation/action-guard/challenge", bytes.NewReader(invalidBody))
	invalidReq.Header.Set("X-Share-Token", federationTestShareToken)
	invalidRecorder := httptest.NewRecorder()
	h.handleActionGuardChallenge(invalidRecorder, invalidReq, server.ID)
	if invalidRecorder.Code != http.StatusBadRequest {
		t.Fatalf("invalid path status=%d body=%s", invalidRecorder.Code, invalidRecorder.Body.String())
	}
}

func TestFederationOwnerRepairsLicenseSlotBeforeIssuingChallenge(t *testing.T) {
	repo, server, _ := newFederationGuardTestRepo(t)
	h := NewFederationHandler(repo, nil, nil)
	prepareCalls, challengeCalls := 0, 0
	h.prepareAction = func(_ context.Context, gotServerID int64) error {
		prepareCalls++
		if gotServerID != server.ID {
			t.Fatalf("prepared server=%d, want %d", gotServerID, server.ID)
		}
		if prepareCalls == 1 {
			return errors.New("replacement lease is still activating")
		}
		return nil
	}
	h.createChallenge = func(_ context.Context, request guardclient.ChallengeRequest) (guardclient.Challenge, error) {
		challengeCalls++
		return guardclient.Challenge{
			ChallengeID: "after-license-repair", Challenge: "random", PayloadHash: request.PayloadHash,
			ExpiresAt: time.Now().Add(time.Minute).Unix(),
		}, nil
	}

	body, _ := json.Marshal(FederationActionChallengeRequest{
		Action: actionNodeCreate, PayloadHash: federationTestHash('a'),
		RequestMethod: http.MethodPost, RequestPath: federationInboundRequestPath,
	})
	request := func() *http.Request {
		req := httptest.NewRequest(http.MethodPost, "/api/federation/action-guard/challenge", bytes.NewReader(body))
		req.Header.Set("X-Share-Token", federationTestShareToken)
		return req
	}

	first := httptest.NewRecorder()
	h.handleActionGuardChallenge(first, request(), server.ID)
	var failed map[string]any
	_ = json.Unmarshal(first.Body.Bytes(), &failed)
	if first.Code != http.StatusForbidden || failed["code"] != "action_guard_retry" || challengeCalls != 0 {
		t.Fatalf("stale slot allocated challenge: status=%d response=%v challenges=%d", first.Code, failed, challengeCalls)
	}
	if len(h.guardChallenges) != 0 {
		t.Fatal("failed lease preflight left a one-time federation challenge behind")
	}

	second := httptest.NewRecorder()
	h.handleActionGuardChallenge(second, request(), server.ID)
	if second.Code != http.StatusOK || prepareCalls != 2 || challengeCalls != 1 {
		t.Fatalf("repaired slot challenge status=%d prepare=%d challenge=%d body=%s",
			second.Code, prepareCalls, challengeCalls, second.Body.String())
	}
	if len(h.guardChallenges) != 1 {
		t.Fatalf("repaired challenge bindings=%d, want 1", len(h.guardChallenges))
	}
}

func federationTestProof(challenge guardclient.Challenge, originalHash, requestPath string) *federationActionProof {
	return &federationActionProof{
		Action: actionNodeCreate, OriginalPayloadHash: originalHash,
		RequestMethod: http.MethodPost, RequestPath: requestPath,
		ChallengeID: challenge.ChallengeID, FrontendPublicKey: "frontend-key", FrontendSignature: "frontend-signature",
	}
}

func callFederationManage(t *testing.T, h *FederationHandler, request federationManageRequest, token string) *httptest.ResponseRecorder {
	t.Helper()
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/api/federation/manage", bytes.NewReader(body))
	req.Header.Set("X-Share-Token", token)
	recorder := httptest.NewRecorder()
	h.handleManage(recorder, req, requestServerID(t, h.repo, token))
	return recorder
}

func requestServerID(t *testing.T, repo *storage.TrafficRepository, token string) int64 {
	t.Helper()
	share, err := repo.GetSharedServerByTokenHash(context.Background(), hashShareToken(token))
	if err != nil {
		// handleManage still revalidates the token; use a non-existent ID for the
		// deliberately invalid-token case.
		return -1
	}
	return share.ServerID
}

func TestFederationManageProofRejectsMissingTamperAndReplayBeforeForward(t *testing.T) {
	repo, server, _ := newFederationGuardTestRepo(t)
	h := NewFederationHandler(repo, nil, nil)
	h.createChallenge = func(_ context.Context, request guardclient.ChallengeRequest) (guardclient.Challenge, error) {
		return guardclient.Challenge{
			ChallengeID: "challenge-" + request.PayloadHash[:12], Challenge: "random", ExpiresAt: time.Now().Add(time.Minute).Unix(),
		}, nil
	}
	authorizeCalls, forwardCalls := 0, 0
	h.authorizeAction = func(_ context.Context, request *http.Request, gotServerID int64, bindingHash string) error {
		if forwardCalls != 0 {
			t.Fatal("mutation was forwarded before Action Guard authorization")
		}
		if gotServerID != server.ID || request.Header.Get("X-MMWX-Challenge-ID") == "" || !validFederationSHA256(bindingHash) {
			t.Fatalf("invalid owner authorization input server=%d binding=%q headers=%v", gotServerID, bindingHash, request.Header)
		}
		authorizeCalls++
		return nil
	}
	h.forwardAgent = func(_ context.Context, gotServerID int64, method, path string, _ []byte) ([]byte, error) {
		if gotServerID != server.ID || method != http.MethodPost || path != "/api/child/inbounds" {
			t.Fatalf("invalid forward target server=%d %s %s", gotServerID, method, path)
		}
		forwardCalls++
		return []byte(`{"success":true}`), nil
	}
	inboundBody := []byte(`{"action":"add","inbound":{"tag":"shared-test"}}`)
	originalHash := federationTestHash('c')
	challenge := issueFederationTestChallenge(t, h, server.ID, originalHash, federationInboundRequestPath)
	proof := federationTestProof(challenge, originalHash, federationInboundRequestPath)
	manage := federationManageRequest{
		Method: http.MethodPost, Path: "/api/child/inbounds", BodyB64: base64.StdEncoding.EncodeToString(inboundBody), GuardProof: proof,
	}

	missing := manage
	missing.GuardProof = nil
	missingRecorder := callFederationManage(t, h, missing, federationTestShareToken)
	if missingRecorder.Code != http.StatusPreconditionRequired || forwardCalls != 0 {
		t.Fatalf("missing proof status=%d forward=%d body=%s", missingRecorder.Code, forwardCalls, missingRecorder.Body.String())
	}

	success := callFederationManage(t, h, manage, federationTestShareToken)
	if success.Code != http.StatusOK || authorizeCalls != 1 || forwardCalls != 1 {
		t.Fatalf("valid manage status=%d auth=%d forward=%d body=%s", success.Code, authorizeCalls, forwardCalls, success.Body.String())
	}
	replay := callFederationManage(t, h, manage, federationTestShareToken)
	if replay.Code != http.StatusUnauthorized || authorizeCalls != 1 || forwardCalls != 1 {
		t.Fatalf("replay status=%d auth=%d forward=%d body=%s", replay.Code, authorizeCalls, forwardCalls, replay.Body.String())
	}

	tamperChallenge := issueFederationTestChallenge(t, h, server.ID, originalHash, federationInboundRequestPath)
	tampered := manage
	tamperedProof := *federationTestProof(tamperChallenge, originalHash, federationInboundRequestPath)
	tamperedProof.OriginalPayloadHash = federationTestHash('d')
	tampered.GuardProof = &tamperedProof
	tamperRecorder := callFederationManage(t, h, tampered, federationTestShareToken)
	if tamperRecorder.Code != http.StatusUnauthorized || authorizeCalls != 1 || forwardCalls != 1 {
		t.Fatalf("tamper status=%d auth=%d forward=%d body=%s", tamperRecorder.Code, authorizeCalls, forwardCalls, tamperRecorder.Body.String())
	}
}

func TestFederationAuthorizeOnlyRunsFullGuardWithoutMutation(t *testing.T) {
	repo, server, _ := newFederationGuardTestRepo(t)
	h := NewFederationHandler(repo, nil, nil)
	h.createChallenge = func(_ context.Context, _ guardclient.ChallengeRequest) (guardclient.Challenge, error) {
		return guardclient.Challenge{ChallengeID: "authorize-only", Challenge: "random", ExpiresAt: time.Now().Add(time.Minute).Unix()}, nil
	}
	authorizeCalls, forwardCalls := 0, 0
	h.authorizeAction = func(_ context.Context, _ *http.Request, gotServerID int64, _ string) error {
		if gotServerID != server.ID {
			t.Fatalf("authorized synthetic consumer server %d, want physical owner server %d", gotServerID, server.ID)
		}
		authorizeCalls++
		return nil
	}
	h.forwardAgent = func(context.Context, int64, string, string, []byte) ([]byte, error) {
		forwardCalls++
		return nil, errors.New("must not forward")
	}
	originalHash := federationTestHash('e')
	challenge := issueFederationTestChallenge(t, h, server.ID, originalHash, federationNodeRequestPath)
	recorder := callFederationManage(t, h, federationManageRequest{
		Method: http.MethodPost, Path: federationNodeRequestPath, AuthorizeOnly: true,
		GuardProof: federationTestProof(challenge, originalHash, federationNodeRequestPath),
	}, federationTestShareToken)
	if recorder.Code != http.StatusOK || authorizeCalls != 1 || forwardCalls != 0 {
		t.Fatalf("authorize-only status=%d auth=%d forward=%d body=%s", recorder.Code, authorizeCalls, forwardCalls, recorder.Body.String())
	}
}

func TestFederatedChallengeHeaderTargetsConsumerRecordWithoutLocalWS(t *testing.T) {
	repo, server, _ := newFederationGuardTestRepo(t)
	owner := NewFederationHandler(repo, nil, nil)
	owner.createChallenge = func(_ context.Context, request guardclient.ChallengeRequest) (guardclient.Challenge, error) {
		return guardclient.Challenge{
			ChallengeID: "owner-consumer-e2e", Challenge: "random", PayloadHash: request.PayloadHash,
			ExpiresAt: time.Now().Add(time.Minute).Unix(),
		}, nil
	}
	authorizeCalls, forwardCalls := 0, 0
	owner.authorizeAction = func(_ context.Context, _ *http.Request, gotServerID int64, bindingHash string) error {
		if gotServerID != server.ID || !validFederationSHA256(bindingHash) {
			t.Fatalf("owner authorized wrong server/hash: server=%d hash=%q", gotServerID, bindingHash)
		}
		authorizeCalls++
		return nil
	}
	owner.forwardAgent = func(_ context.Context, gotServerID int64, method, path string, _ []byte) ([]byte, error) {
		if gotServerID != server.ID || method != http.MethodPost || path != "/api/child/inbounds" {
			t.Fatalf("owner forwarded wrong target: server=%d %s %s", gotServerID, method, path)
		}
		forwardCalls++
		return []byte(`{"success":true}`), nil
	}
	ownerHTTP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/federation/action-guard/challenge":
			owner.handleActionGuardChallenge(w, r, server.ID)
		case "/api/federation/manage":
			owner.handleManage(w, r, server.ID)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ownerHTTP.Close()

	consumer := storage.RemoteServer{Name: "consumer-placeholder", Token: "synthetic-token", Status: storage.RemoteServerStatusConnected}
	if err := repo.CreateRemoteServer(context.Background(), &consumer); err != nil {
		t.Fatalf("CreateRemoteServer consumer: %v", err)
	}
	if err := repo.SetFederatedServer(context.Background(), consumer.ID, ownerHTTP.URL, federationTestShareToken, ""); err != nil {
		t.Fatalf("SetFederatedServer: %v", err)
	}
	rm := NewRemoteManageHandler(repo, nil) // deliberately no local Agent WS
	rm.httpClient = ownerHTTP.Client()
	inboundBody := []byte(`{"action":"add","inbound":{"tag":"consumer-e2e"}}`)
	payloadHash := hashActionPayload(inboundBody)
	challenge, handled, err := rm.ProxyFederatedActionChallenge(context.Background(), consumer.ID, FederationActionChallengeRequest{
		Action: actionNodeCreate, PayloadHash: payloadHash,
		RequestMethod: http.MethodPost, RequestPath: federationInboundRequestPath,
	})
	if err != nil || !handled || challenge.PayloadHash == payloadHash {
		t.Fatalf("owner challenge proxy failed handled=%v challenge=%+v err=%v", handled, challenge, err)
	}
	request := httptest.NewRequest(http.MethodPost, federationInboundRequestPath, nil)
	request.Header.Set(actionGuardServerIDHeader, strconv.FormatInt(consumer.ID, 10))
	request.Header.Set("X-MMWX-Challenge-ID", challenge.ChallengeID)
	request.Header.Set("X-MMWX-Frontend-Key", "frontend-key")
	request.Header.Set("X-MMWX-Frontend-Signature", "frontend-signature")
	prepared, federated, err := rm.prepareFederatedManagedNodeAction(
		context.Background(), request, consumer.ID, payloadHash, federationInboundRequestPath, false,
	)
	if err != nil || !federated || federationGuardProofFromContext(prepared) == nil {
		t.Fatalf("federated prepare failed federated=%v err=%v", federated, err)
	}
	result, err := rm.forwardToRemoteServer(prepared, consumer.ID, http.MethodPost, "/api/child/inbounds", inboundBody)
	if err != nil || string(result) != `{"success":true}` || authorizeCalls != 1 || forwardCalls != 1 {
		t.Fatalf("consumer federation failed result=%s auth=%d forward=%d err=%v", result, authorizeCalls, forwardCalls, err)
	}

	request.Header.Set(actionGuardServerIDHeader, strconv.FormatInt(consumer.ID+1, 10))
	_, federated, err = rm.prepareFederatedManagedNodeAction(
		context.Background(), request, consumer.ID, payloadHash, federationInboundRequestPath, false,
	)
	var typed *FederationRequestError
	if !federated || !errors.As(err, &typed) || typed.Code != codeFederatedChallengeRequired || typed.ServerID != consumer.ID {
		t.Fatalf("wrong header did not request a targeted challenge: federated=%v err=%#v", federated, err)
	}
}

func TestFederationOwnerPreservesActionGuardRetryWithoutForward(t *testing.T) {
	repo, server, _ := newFederationGuardTestRepo(t)
	h := NewFederationHandler(repo, nil, nil)
	h.createChallenge = func(_ context.Context, _ guardclient.ChallengeRequest) (guardclient.Challenge, error) {
		return guardclient.Challenge{ChallengeID: "retry-owner", Challenge: "random", ExpiresAt: time.Now().Add(time.Minute).Unix()}, nil
	}
	h.authorizeAction = func(context.Context, *http.Request, int64, string) error {
		return &authoritativeSlotRepairedError{serverID: server.ID, tokenHash: federationTestHash('a'), cause: errors.New("slot replaced")}
	}
	forwardCalls := 0
	h.forwardAgent = func(context.Context, int64, string, string, []byte) ([]byte, error) {
		forwardCalls++
		return nil, nil
	}
	originalHash := federationTestHash('a')
	challenge := issueFederationTestChallenge(t, h, server.ID, originalHash, federationInboundRequestPath)
	recorder := callFederationManage(t, h, federationManageRequest{
		Method: http.MethodPost, Path: "/api/child/inbounds",
		BodyB64:    base64.StdEncoding.EncodeToString([]byte(`{"action":"add","inbound":{"tag":"retry"}}`)),
		GuardProof: federationTestProof(challenge, originalHash, federationInboundRequestPath),
	}, federationTestShareToken)
	var response map[string]any
	_ = json.Unmarshal(recorder.Body.Bytes(), &response)
	if recorder.Code != http.StatusForbidden || response["code"] != "action_guard_retry" || forwardCalls != 0 {
		t.Fatalf("retry signal lost status=%d response=%v forward=%d", recorder.Code, response, forwardCalls)
	}
}

func TestFederationGuardRPCNeverFallsBackToAgentHTTP(t *testing.T) {
	repo, server, _ := newFederationGuardTestRepo(t)
	hits := make(chan struct{}, 1)
	legacyHTTP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits <- struct{}{}
		w.WriteHeader(http.StatusOK)
	}))
	defer legacyHTTP.Close()
	rm := NewRemoteManageHandler(repo, nil)
	rm.httpClient = legacyHTTP.Client()
	_, err := rm.forwardGuardOverWS(context.Background(), server.ID, http.MethodPost, "/api/child/action-guard/attest", []byte(`{}`))
	if err == nil || !strings.Contains(err.Error(), "Agent 未连接 WS") {
		t.Fatalf("strict Guard RPC unexpectedly succeeded: %v", err)
	}
	select {
	case <-hits:
		t.Fatal("strict Agent Guard authorization fell back to HTTP")
	default:
	}
}

func TestFederationChallengeProxyFailsClosedAgainstOldOwner(t *testing.T) {
	repo, consumer, _ := newFederationGuardTestRepo(t)
	oldOwner := httptest.NewServer(http.NotFoundHandler())
	defer oldOwner.Close()
	if err := repo.SetFederatedServer(context.Background(), consumer.ID, oldOwner.URL, federationTestShareToken, ""); err != nil {
		t.Fatalf("SetFederatedServer: %v", err)
	}
	rm := NewRemoteManageHandler(repo, nil)
	_, handled, err := rm.ProxyFederatedActionChallenge(context.Background(), consumer.ID, FederationActionChallengeRequest{
		Action: actionNodeCreate, PayloadHash: federationTestHash('a'),
		RequestMethod: http.MethodPost, RequestPath: federationInboundRequestPath,
	})
	var typed *FederationRequestError
	if !handled || !errors.As(err, &typed) || typed.Status != http.StatusUpgradeRequired || typed.Code != codeFederationGuardUnsupported {
		t.Fatalf("old owner did not fail closed: handled=%v err=%#v", handled, err)
	}
}

func TestFederationErrorFromBodyPreservesStructuredFields(t *testing.T) {
	err := federationErrorFromBody(http.StatusForbidden, []byte(`{"error":"retry with fresh challenge","code":"action_guard_retry","server_id":77}`))
	var typed *FederationRequestError
	if !errors.As(err, &typed) || typed.Status != http.StatusForbidden || typed.Code != "action_guard_retry" || typed.ServerID != 77 {
		t.Fatalf("structured federation error was flattened: %#v", err)
	}
	status, payload := federationErrorResponse(err, http.StatusBadGateway, "fallback")
	if status != http.StatusForbidden || payload["code"] != "action_guard_retry" || payload["server_id"] != int64(77) {
		t.Fatalf("structured response was not preserved: status=%d payload=%v", status, payload)
	}
}
