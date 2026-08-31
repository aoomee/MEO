package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"

	"miaomiaowux/internal/storage"
)

func newForwardTestHandler(t *testing.T) *ForwardHandler {
	t.Helper()
	repo, err := storage.NewTrafficRepository(filepath.Join(t.TempDir(), "fwd.db"))
	if err != nil {
		t.Fatalf("NewTrafficRepository: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	return NewForwardHandler(repo, nil) // remote=nil:不测 apply(需真机 agent)
}

// do 发一个请求给 handler,返回状态码 + 解析后的 JSON map。
func do(t *testing.T, h http.Handler, method, path string, body any) (int, map[string]any) {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	var out map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	return rec.Code, out
}

func TestForwardAdminGroupAndChainCRUD(t *testing.T) {
	h := newForwardTestHandler(t)

	// 建组(带成员)
	code, resp := do(t, h, http.MethodPost, "/api/admin/forward/groups", map[string]any{
		"name":             "entry",
		"balance_strategy": "percentage",
		"failover_enabled": true,
		"members":          []map[string]any{{"server_id": 1, "weight": 3}, {"server_id": 2}},
	})
	if code != http.StatusOK || resp["success"] != true {
		t.Fatalf("createGroup: code=%d resp=%v", code, resp)
	}
	gid := int64(resp["id"].(float64))

	// 列组
	code, resp = do(t, h, http.MethodGet, "/api/admin/forward/groups", nil)
	if code != http.StatusOK {
		t.Fatalf("listGroups code=%d", code)
	}
	if groups, _ := resp["groups"].([]any); len(groups) != 1 {
		t.Fatalf("期望 1 个组,得到 %v", resp["groups"])
	}

	// 取单组 + 成员回读
	code, resp = do(t, h, http.MethodGet, "/api/admin/forward/groups/"+i64a(gid), nil)
	if code != http.StatusOK {
		t.Fatalf("getGroup code=%d", code)
	}
	g := resp["group"].(map[string]any)
	if members, _ := g["members"].([]any); len(members) != 2 {
		t.Fatalf("期望 2 个成员,得到 %v", g["members"])
	}

	// 非法策略应 400
	code, _ = do(t, h, http.MethodPost, "/api/admin/forward/groups", map[string]any{
		"name": "bad", "balance_strategy": "nope",
	})
	if code != http.StatusBadRequest {
		t.Fatalf("非法策略应 400,得到 %d", code)
	}

	// 全量替换成员
	code, _ = do(t, h, http.MethodPut, "/api/admin/forward/groups/"+i64a(gid)+"/members", map[string]any{
		"members": []map[string]any{{"server_id": 9, "weight": 1}},
	})
	if code != http.StatusOK {
		t.Fatalf("setMembers code=%d", code)
	}
	_, resp = do(t, h, http.MethodGet, "/api/admin/forward/groups/"+i64a(gid), nil)
	g = resp["group"].(map[string]any)
	if members, _ := g["members"].([]any); len(members) != 1 {
		t.Fatalf("替换后应剩 1 个成员,得到 %v", g["members"])
	}

	// 再建两个组当出口/中间跳
	_, r2 := do(t, h, http.MethodPost, "/api/admin/forward/groups", map[string]any{"name": "mid"})
	gid2 := int64(r2["id"].(float64))
	_, r3 := do(t, h, http.MethodPost, "/api/admin/forward/groups", map[string]any{"name": "exit"})
	gid3 := int64(r3["id"].(float64))

	// 建链(带 hops)
	code, resp = do(t, h, http.MethodPost, "/api/admin/forward/chains", map[string]any{
		"name":      "chain-abc",
		"group_ids": []int64{gid, gid2, gid3},
	})
	if code != http.StatusOK {
		t.Fatalf("createChain code=%d resp=%v", code, resp)
	}
	cid := int64(resp["id"].(float64))

	// 取链回读 hops
	code, resp = do(t, h, http.MethodGet, "/api/admin/forward/chains/"+i64a(cid), nil)
	if code != http.StatusOK {
		t.Fatalf("getChain code=%d", code)
	}
	c := resp["chain"].(map[string]any)
	if hops, _ := c["hops"].([]any); len(hops) != 3 {
		t.Fatalf("期望 3 跳,得到 %v", c["hops"])
	}

	// 更新链端口段
	code, _ = do(t, h, http.MethodPut, "/api/admin/forward/chains/"+i64a(cid), map[string]any{
		"port_range_start": 50000, "port_range_end": 50010,
	})
	if code != http.StatusOK {
		t.Fatalf("updateChain code=%d", code)
	}
	_, resp = do(t, h, http.MethodGet, "/api/admin/forward/chains/"+i64a(cid), nil)
	c = resp["chain"].(map[string]any)
	if int(c["port_range_start"].(float64)) != 50000 || int(c["port_range_end"].(float64)) != 50010 {
		t.Fatalf("端口段未写回: %v", c)
	}

	// 删链、删组
	if code, _ = do(t, h, http.MethodDelete, "/api/admin/forward/chains/"+i64a(cid), nil); code != http.StatusOK {
		t.Fatalf("deleteChain code=%d", code)
	}
	if code, _ = do(t, h, http.MethodDelete, "/api/admin/forward/groups/"+i64a(gid), nil); code != http.StatusOK {
		t.Fatalf("deleteGroup code=%d", code)
	}
}

func TestForwardCreateAndRebindNode(t *testing.T) {
	h := newForwardTestHandler(t)
	_, r1 := do(t, h, http.MethodPost, "/api/admin/forward/groups", map[string]any{"name": "entry"})
	_, r2 := do(t, h, http.MethodPost, "/api/admin/forward/groups", map[string]any{"name": "exit"})
	gid1 := int64(r1["id"].(float64))
	gid2 := int64(r2["id"].(float64))

	code, resp := do(t, h, http.MethodPost, "/api/admin/forward/chains", map[string]any{
		"name": "c1", "group_ids": []int64{gid1, gid2},
		"port_range_start": 51000, "port_range_end": 51000,
	})
	if code != http.StatusOK {
		t.Fatalf("createChain: %d %v", code, resp)
	}
	cid1 := int64(resp["id"].(float64))
	_, resp = do(t, h, http.MethodPost, "/api/admin/forward/chains", map[string]any{
		"name": "c2", "group_ids": []int64{gid1, gid2},
	})
	cid2 := int64(resp["id"].(float64))

	code, resp = do(t, h, http.MethodPost, "/api/admin/forward/chains/"+i64a(cid1)+"/create-node", map[string]any{
		"node_name": "入口A", "relay_protocol": "tcp",
	})
	if code != http.StatusOK {
		t.Fatalf("create-node: %d %v", code, resp)
	}
	if int(resp["count"].(float64)) < 1 {
		t.Fatalf("create-node count=%v", resp["count"])
	}
	if int(resp["port"].(float64)) != 51000 {
		t.Fatalf("create-node 应使用链端口段,得到 %v", resp["port"])
	}

	binds, err := h.repo.ListForwardChainBindingsByChain(t.Context(), cid1)
	if err != nil || len(binds) == 0 {
		t.Fatalf("绑定未写入: %v %v", binds, err)
	}
	nodeID := binds[0].NodeID

	code, resp = do(t, h, http.MethodPost, "/api/admin/forward/chains/"+i64a(cid2)+"/rebind-node", map[string]any{
		"node_id": nodeID,
	})
	if code != http.StatusOK {
		t.Fatalf("rebind-node: %d %v", code, resp)
	}
	b, err := h.repo.GetForwardChainBinding(t.Context(), nodeID)
	if err != nil || b == nil || b.ChainID != cid2 {
		t.Fatalf("改绑失败: %+v %v", b, err)
	}

	// 已有节点绑到链
	n, err := h.repo.CreateNode(t.Context(), storage.Node{
		Username: "admin", NodeName: "已有", Protocol: "ss",
		ClashConfig: `{"name":"已有","type":"ss","server":"1.2.3.4","port":443}`,
		Enabled:     true,
	})
	if err != nil {
		t.Fatalf("CreateNode: %v", err)
	}
	code, resp = do(t, h, http.MethodPost, "/api/admin/forward/chains/"+i64a(cid1)+"/create-node", map[string]any{
		"existing_node_id": n.ID, "port": 52000, "relay_protocol": "tcp",
	})
	if code != http.StatusOK {
		t.Fatalf("create-node existing: %d %v", code, resp)
	}
}

func i64a(v int64) string {
	return strconv.FormatInt(v, 10)
}

func TestForwardChainIssuesReportedOnList(t *testing.T) {
	h := newForwardTestHandler(t)
	_, r1 := do(t, h, http.MethodPost, "/api/admin/forward/groups", map[string]any{"name": "入口空组"})
	_, r2 := do(t, h, http.MethodPost, "/api/admin/forward/groups", map[string]any{"name": "出口空组"})
	gid1 := int64(r1["id"].(float64))
	gid2 := int64(r2["id"].(float64))

	code, resp := do(t, h, http.MethodPost, "/api/admin/forward/chains", map[string]any{
		"name": "残链", "group_ids": []int64{gid1, gid2},
	})
	if code != http.StatusOK {
		t.Fatalf("createChain: %d %v", code, resp)
	}
	warnings, _ := resp["warnings"].([]any)
	if len(warnings) == 0 {
		t.Fatalf("创建不完整链应带回 warnings: %v", resp)
	}

	code, resp = do(t, h, http.MethodGet, "/api/admin/forward/chains", nil)
	if code != http.StatusOK {
		t.Fatalf("listChains: %d", code)
	}
	issues, _ := resp["issues"].(map[string]any)
	if len(issues) == 0 {
		t.Fatalf("列表应带 issues: %v", resp)
	}
	cid := int64(resp["chains"].([]any)[0].(map[string]any)["id"].(float64))
	got, _ := issues[i64a(cid)].([]any)
	if len(got) == 0 {
		t.Fatalf("链 %d 的 issues 为空: %v", cid, issues)
	}
}

func TestEvaluateForwardChainIssues(t *testing.T) {
	entry := &storage.ForwardGroup{ID: 1, Name: "入", Members: []storage.ForwardGroupMember{{ServerID: 11}}}
	exit := &storage.ForwardGroup{ID: 2, Name: "出", Members: []storage.ForwardGroupMember{{ServerID: 22}}}
	servers := map[int64]*storage.RemoteServer{
		11: {ID: 11, IPAddress: "10.0.0.1"},
		22: {ID: 22},
	}
	chain := &storage.ForwardChain{ID: 9, Hops: []storage.ForwardChainHop{{GroupID: 1}, {GroupID: 2}}}
	got := evaluateForwardChainIssues(chain, map[int64]*storage.ForwardGroup{1: entry, 2: exit}, servers)
	if len(got) != 1 || got[0] != "出口组「出」的服务器没有可用地址" {
		t.Fatalf("issues=%v", got)
	}
	short := evaluateForwardChainIssues(&storage.ForwardChain{ID: 1}, nil, nil)
	if len(short) != 1 {
		t.Fatalf("短链 issues=%v", short)
	}
}
