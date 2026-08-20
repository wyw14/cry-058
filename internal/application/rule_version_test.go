package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/repository/memory"
	"testing"
)

func TestRulePublishRejectsDuplicateVersion(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Projects.Create(ctx, &domain.GrantProject{ID: "p", Code: "P", Name: "公益", Year: 2026, AnnualCapCents: 1000})
	svc := NewRuleService(p)
	r := domain.RuleVersion{ProjectID: "p", SchemeID: "s", Name: "r", Year: 2026, VersionNo: 1, BaseRateBps: 5000, Status: domain.RulePublished}
	if _, e := svc.Publish(ctx, r); e != nil {
		t.Fatal(e)
	}
	if _, e := svc.Publish(ctx, r); e == nil {
		t.Fatal("duplicate rule version accepted")
	}
}
