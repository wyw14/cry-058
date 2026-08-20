package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/repository/memory"
	"testing"
)

func TestClaimYearMustMatchProject(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Projects.Create(ctx, &domain.GrantProject{ID: "p", Code: "P", Name: "公益", Year: 2026, AnnualCapCents: 1000})
	_ = p.Beneficiaries.Create(ctx, &domain.Beneficiary{ID: "b", Name: "对象", IdentityHash: "h", Active: true})
	_ = p.Schemes.Create(ctx, &domain.CoverageScheme{ID: "s", ProjectID: "p", Name: "方案", Active: true})
	_, e := NewClaimService(p).Submit(ctx, domain.ExpenseClaim{ProjectID: "p", BeneficiaryID: "b", SchemeID: "s", Year: 2025, AmountCents: 10, Summary: "x", IdempotencyKey: "k"})
	if e == nil {
		t.Fatal("cross-year claim accepted")
	}
}
