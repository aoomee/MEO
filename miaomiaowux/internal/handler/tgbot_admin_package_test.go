package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"miaomiaowux/internal/storage"
)

func TestAdminPackageAssignmentAppearsInTGBotSubscriptions(t *testing.T) {
	repo := newResetTestRepo(t)
	ctx := context.Background()
	if err := repo.CreateUser(ctx, "admin", "admin@example.com", "admin", "hash", storage.RoleAdmin, ""); err != nil {
		t.Fatal(err)
	}
	pkgID, err := repo.CreatePackage(ctx, storage.Package{Name: "管理员套餐", CycleDays: 30})
	if err != nil {
		t.Fatal(err)
	}

	assign := NewPackageAssignHandler(repo, nil, nil)
	rec := httptest.NewRecorder()
	assign.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/admin/packages/assign",
		strings.NewReader(`{"username":"admin","package_id":`+strconv.FormatInt(pkgID, 10)+`}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("assign admin package returned %d: %s", rec.Code, rec.Body.String())
	}
	admin, err := repo.GetUser(ctx, "admin")
	if err != nil || admin.PackageID != pkgID {
		t.Fatalf("admin package was not persisted: user=%+v err=%v", admin, err)
	}

	tg := NewTGBotAPIHandler(repo)
	rec = httptest.NewRecorder()
	tg.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/admin/tgbot/user-subscriptions?username=admin", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("subscriptions returned %d: %s", rec.Code, rec.Body.String())
	}
	var response struct {
		DefaultSubscriptions []struct {
			Name         string `json:"name"`
			CombinedCode string `json:"combined_code"`
		} `json:"default_subscriptions"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if len(response.DefaultSubscriptions) != 1 || response.DefaultSubscriptions[0].Name != "管理员套餐" || response.DefaultSubscriptions[0].CombinedCode == "" {
		t.Fatalf("unexpected package subscriptions: %+v", response.DefaultSubscriptions)
	}
}

func TestUnassignedAdminTGBotSubviewUsesFirstManagedSubscription(t *testing.T) {
	repo := newResetTestRepo(t)
	ctx := context.Background()
	if err := repo.CreateUser(ctx, "admin", "admin@example.com", "admin", "hash", storage.RoleAdmin, ""); err != nil {
		t.Fatal(err)
	}
	for _, file := range []storage.SubscribeFile{
		{Name: "后显示", Type: storage.SubscribeTypeCreate, Filename: "later.yaml", SortOrder: 20, CreatedBy: "admin"},
		{Name: "第一条", Type: storage.SubscribeTypeCreate, Filename: "first.yaml", SortOrder: 10, CreatedBy: "admin"},
	} {
		if _, err := repo.CreateSubscribeFile(ctx, file); err != nil {
			t.Fatal(err)
		}
	}

	tg := NewTGBotAPIHandler(repo)
	rec := httptest.NewRecorder()
	tg.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/admin/tgbot/admin-subview?username=admin", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("admin subview returned %d: %s", rec.Code, rec.Body.String())
	}
	var response struct {
		Subscription *struct {
			Name         string `json:"name"`
			CombinedCode string `json:"combined_code"`
		} `json:"subscription"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Subscription == nil || response.Subscription.Name != "第一条" || response.Subscription.CombinedCode == "" {
		t.Fatalf("unexpected admin fallback subscription: %+v", response.Subscription)
	}
}
