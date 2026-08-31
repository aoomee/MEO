package storage

import (
	"context"
	"errors"
	"testing"
	"time"
)

// 同一用户不能对同一套餐既有主绑定又有独立绑定。
//
// 线上事故:管理套餐弹窗里「新增为独立套餐」和底部「确认保存」共用同一个
// selectedPackageId。管理员挨个点选套餐加成独立绑定后,高亮停在最后点的那个,
// 再点「确认保存」就把主套餐也指到了它 —— 于是「已绑定套餐」「复制订阅链接」
// 里出现两条一模一样的记录,退订时无法分辨该退哪个。
func TestAssignPackageRejectsDuplicateOfIndependentAssignment(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	if err := repo.CreateUser(ctx, "dupuser", "", "", "hash", RoleUser, ""); err != nil {
		t.Fatal(err)
	}
	pkgA, err := repo.CreatePackage(ctx, Package{Name: "pkg-a", CycleDays: 30, TrafficLimitBytes: 1 << 30})
	if err != nil {
		t.Fatal(err)
	}
	pkgB, err := repo.CreatePackage(ctx, Package{Name: "pkg-b", CycleDays: 30, TrafficLimitBytes: 1 << 30})
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	end := now.AddDate(0, 1, 0)

	// 先有一个主套餐 A(走遗留路径,会生成 legacy 镜像行)
	if err := repo.AssignPackageToUser(ctx, "dupuser", pkgA, now, end, true, 1); err != nil {
		t.Fatalf("assign primary: %v", err)
	}
	// 再把 B 加成独立套餐
	if _, err := repo.CreateUserPackageAssignment(ctx, "dupuser", pkgB, now, end, true, 1); err != nil {
		t.Fatalf("create independent: %v", err)
	}

	// 关键:此时把主套餐也改成 B —— 必须被拒绝,否则 B 会出现两次
	err = repo.AssignPackageToUser(ctx, "dupuser", pkgB, now, end, true, 1)
	if !errors.Is(err, ErrPackageAlreadyAssignedIndependently) {
		t.Fatalf("把主套餐设成已独立绑定的套餐本应被拒绝,实际 err=%v", err)
	}

	// 确认没有产生重复
	assignments, err := repo.ListUserPackageAssignments(ctx, "dupuser", false)
	if err != nil {
		t.Fatal(err)
	}
	seen := make(map[int64]int)
	for _, a := range assignments {
		if a.Status == "active" {
			seen[a.PackageID]++
		}
	}
	for pkgID, n := range seen {
		if n > 1 {
			t.Fatalf("套餐 %d 出现 %d 条重复绑定: %+v", pkgID, n, assignments)
		}
	}

	// 换成一个没被独立绑定的套餐,应当正常放行(不能误伤正常改主套餐)
	if err := repo.AssignPackageToUser(ctx, "dupuser", pkgA, now, end, true, 1); err != nil {
		t.Fatalf("改成未被独立绑定的套餐本应成功: %v", err)
	}
}
