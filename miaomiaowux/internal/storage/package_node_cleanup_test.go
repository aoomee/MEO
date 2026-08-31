package storage

import (
	"context"
	"testing"
)

// 删节点必须把它从套餐的节点列表和按 node_id 索引的配置映射里一并摘掉。
//
// 线上真实事故:客户删掉套餐「前洵独享SG网络」引用的两个节点后,packages.nodes 与
// node_traffic_limits 里仍残留这两个 id,之后该套餐每次保存都被
// validatePackageNodeTrafficLimits 拒绝(「节点 74 不在套餐内」),彻底救不回来。
func TestDeleteNodePurgesPackageReferences(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	if err := repo.CreateUser(ctx, "pkgowner", "", "", "hash", RoleUser, ""); err != nil {
		t.Fatal(err)
	}

	mk := func(name string) int64 {
		n, err := repo.CreateNode(ctx, Node{
			Username: "pkgowner", RawURL: "ss://x@1.2.3.4:8388", NodeName: name,
			Protocol: "ss", ParsedConfig: `{"name":"` + name + `"}`,
			ClashConfig: `{"name":"` + name + `","type":"ss"}`, Enabled: true, Tag: "手动输入",
		})
		if err != nil {
			t.Fatalf("create node %s: %v", name, err)
		}
		return n.ID
	}
	doomed := mk("ZZ-Doomed")
	kept := mk("ZZ-Kept")

	pkgID, err := repo.CreatePackage(ctx, Package{
		Name: "purge-pkg", CycleDays: 30, TrafficLimitBytes: 1 << 30,
		Nodes:             []int64{doomed, kept},
		NodeTrafficLimits: map[int64]float64{doomed: 500, kept: 500},
		NodeSpeedLimits:   map[int64]float64{doomed: 100, kept: 100},
		NodeMultipliers:   map[int64]float64{doomed: 2, kept: 1},
		NodeDeviceLimits:  map[int64]int{doomed: 3, kept: 3},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.DeleteNode(ctx, doomed, "pkgowner"); err != nil {
		t.Fatalf("delete node: %v", err)
	}

	pkg, err := repo.GetPackage(ctx, pkgID)
	if err != nil {
		t.Fatal(err)
	}

	for _, id := range pkg.Nodes {
		if id == doomed {
			t.Fatalf("已删节点仍留在套餐节点列表: %v", pkg.Nodes)
		}
	}
	if _, ok := pkg.NodeTrafficLimits[doomed]; ok {
		t.Fatalf("已删节点的流量额度残留: %v", pkg.NodeTrafficLimits)
	}
	if _, ok := pkg.NodeSpeedLimits[doomed]; ok {
		t.Fatalf("已删节点的限速残留: %v", pkg.NodeSpeedLimits)
	}
	if _, ok := pkg.NodeMultipliers[doomed]; ok {
		t.Fatalf("已删节点的倍率残留: %v", pkg.NodeMultipliers)
	}
	if _, ok := pkg.NodeDeviceLimits[doomed]; ok {
		t.Fatalf("已删节点的设备数残留: %v", pkg.NodeDeviceLimits)
	}

	// 保留节点的配置不能被误删
	if len(pkg.Nodes) != 1 || pkg.Nodes[0] != kept {
		t.Fatalf("误删了保留节点: %v", pkg.Nodes)
	}
	if pkg.NodeTrafficLimits[kept] != 500 || pkg.NodeSpeedLimits[kept] != 100 {
		t.Fatalf("保留节点的配置被破坏: traffic=%v speed=%v", pkg.NodeTrafficLimits, pkg.NodeSpeedLimits)
	}
}
