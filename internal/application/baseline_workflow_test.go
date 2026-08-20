package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/repository/memory"
	"testing"
)

func TestGrantSettlementWorkflowUsesPublishedRule(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Projects.Create(ctx, &domain.GrantProject{ID: "p", Code: "P", Name: "公益", Year: 2026, AnnualCapCents: 100000})
	_ = p.Beneficiaries.Create(ctx, &domain.Beneficiary{ID: "b", Name: "对象", IdentityHash: "hash", Active: true})
	_ = p.Schemes.Create(ctx, &domain.CoverageScheme{ID: "s", ProjectID: "p", Name: "基础", Active: true})
	_ = p.Rules.Create(ctx, &domain.RuleVersion{ID: "r", ProjectID: "p", SchemeID: "s", Name: "2026", Year: 2026, BaseRateBps: 5000, ExcessRateBps: 2500, ThresholdCents: 1000, CapCents: 100000, Status: domain.RulePublished, VersionNo: 1})
	_ = p.Claims.Create(ctx, &domain.ExpenseClaim{ID: "c", ProjectID: "p", BeneficiaryID: "b", SchemeID: "s", Year: 2026, Summary: "医疗", Category: "medical", AmountCents: 2000, Status: domain.ClaimSubmitted, IdempotencyKey: "k"})
	s, _, e := NewSettlementService(p).Preview(ctx, "c")
	if e != nil {
		t.Fatal(e)
	}
	if s.SubsidyCents <= 0 {
		t.Fatal("expected subsidy")
	}
}
