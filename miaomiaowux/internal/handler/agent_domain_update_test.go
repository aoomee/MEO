package handler

import (
	"testing"

	"miaomiaowux/internal/storage"
)

func TestFollowLinkedServerDomainKeepsIndependentDomain(t *testing.T) {
	old := &storage.RemoteServer{PullAddress: "1.2.3.4", Domain: "keep.example.com"}
	addr := "5.6.7.8"
	req := &RemoteServerUpdateRequest{PullAddress: addr}
	followLinkedServerDomain(old, req)
	if req.Domain != nil {
		t.Fatalf("独立域名不应被部分更新改写, got %v", *req.Domain)
	}
}

func TestFollowLinkedServerDomainMovesLinkedDomain(t *testing.T) {
	old := &storage.RemoteServer{PullAddress: "old.example.com", Domain: "old.example.com"}
	req := &RemoteServerUpdateRequest{PullAddress: "new.example.com"}
	followLinkedServerDomain(old, req)
	if req.Domain == nil || *req.Domain != "new.example.com" {
		t.Fatalf("跟随域名应一起改, got %v", req.Domain)
	}
}
