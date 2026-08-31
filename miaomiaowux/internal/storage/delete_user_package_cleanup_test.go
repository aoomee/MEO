package storage

import (
	"context"
	"testing"
	"time"
)

// 删除用户必须一并清掉 user_package_assignments。
//
// 该表按 username 记录、不随 users 级联,留下孤儿行会让这个套餐之后
// 编辑/删除时按 assignment.Username 反查用户失败("用户不存在"),
// 整个套餐操作被彻底卡死(许可证站 #269 / #273)。
func TestDeleteUserRemovesPackageAssignments(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	if err := repo.CreateUser(ctx, "doomed", "", "", "hash", RoleUser, ""); err != nil {
		t.Fatal(err)
	}
	if err := repo.CreateUser(ctx, "keeper", "", "", "hash", RoleUser, ""); err != nil {
		t.Fatal(err)
	}
	packageID, err := repo.CreatePackage(ctx, Package{Name: "cleanup-pkg", CycleDays: 30, TrafficLimitBytes: 1 << 30})
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	for _, u := range []string{"doomed", "keeper"} {
		if err := repo.AssignPackageToUser(ctx, u, packageID, now, now.AddDate(0, 1, 0), true, 1); err != nil {
			t.Fatalf("assign %s: %v", u, err)
		}
	}

	before, err := repo.ListActivePackageAssignmentsByPackage(ctx, packageID)
	if err != nil {
		t.Fatal(err)
	}
	if len(before) != 2 {
		t.Fatalf("前置条件不成立: 期望 2 条关联, 实际 %d", len(before))
	}

	if err := repo.DeleteUser(ctx, "doomed"); err != nil {
		t.Fatalf("delete user: %v", err)
	}

	after, err := repo.ListActivePackageAssignmentsByPackage(ctx, packageID)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range after {
		if a.Username == "doomed" {
			t.Fatalf("已删用户仍残留套餐关联(孤儿行): %+v", a)
		}
	}
	if len(after) != 1 || after[0].Username != "keeper" {
		t.Fatalf("误删了其他用户的关联: %+v", after)
	}
}
