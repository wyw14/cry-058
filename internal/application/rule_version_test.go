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

// 不同版本号或不同业务组合（项目/方案/年度）不应被重复校验拦截。
func TestRulePublishAllowsDistinctVersionAndScope(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Projects.Create(ctx, &domain.GrantProject{ID: "p", Code: "P", Name: "公益", Year: 2026, AnnualCapCents: 1000})
	_ = p.Projects.Create(ctx, &domain.GrantProject{ID: "p2", Code: "Q", Name: "助学", Year: 2026, AnnualCapCents: 1000})
	svc := NewRuleService(p)
	base := domain.RuleVersion{ProjectID: "p", SchemeID: "s", Name: "r", Year: 2026, VersionNo: 1, BaseRateBps: 5000, Status: domain.RulePublished}
	if _, e := svc.Publish(ctx, base); e != nil {
		t.Fatal(e)
	}
	// 同项目/方案/年度但不同版本号 → 允许
	v2 := base
	v2.VersionNo = 2
	if _, e := svc.Publish(ctx, v2); e != nil {
		t.Fatalf("distinct version_no should publish, got %v", e)
	}
	// 不同方案 → 允许
	other := base
	other.SchemeID = "s2"
	if _, e := svc.Publish(ctx, other); e != nil {
		t.Fatalf("distinct scheme should publish, got %v", e)
	}
	// 不同项目 → 允许
	proj := base
	proj.ProjectID = "p2"
	if _, e := svc.Publish(ctx, proj); e != nil {
		t.Fatalf("distinct project should publish, got %v", e)
	}
	// 不同年度 → 允许
	yr := base
	yr.Year = 2027
	if _, e := svc.Publish(ctx, yr); e != nil {
		t.Fatalf("distinct year should publish, got %v", e)
	}
}
