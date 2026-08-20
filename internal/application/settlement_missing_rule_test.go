package application

import (
	"context"
	"testing"

	"github.com/wyw14/cry058/internal/domain"
)

// TestSettlementPreviewWithoutRuleReturnsZeroAndWarning 覆盖缺失规则回归场景：
// 当项目当年度没有已发布规则时，预览不得崩溃，应安全返回零补助、空明细与明确提示。
func TestSettlementPreviewWithoutRuleReturnsZeroAndWarning(t *testing.T) {
	p := NewPortsForTest()
	ctx := context.Background()
	prj := domain.GrantProject{ID: "p2", Code: "P2", Name: "公益二", Year: 2026, AnnualCapCents: 100000}
	if e := p.Projects.Create(ctx, &prj); e != nil {
		t.Fatal(e)
	}
	ben := domain.Beneficiary{ID: "b2", Name: "对象二", IdentityHash: "hash2", Active: true}
	if e := p.Beneficiaries.Create(ctx, &ben); e != nil {
		t.Fatal(e)
	}
	claim := domain.ExpenseClaim{ID: "c2", ProjectID: "p2", BeneficiaryID: "b2", SchemeID: "s2", Year: 2026, Summary: "医疗", Category: "medical", AmountCents: 2000, Status: domain.ClaimSubmitted, IdempotencyKey: "k2"}
	if e := p.Claims.Create(ctx, &claim); e != nil {
		t.Fatal(e)
	}
	// 注意：未创建任何规则，Latest 将找不到已发布规则。
	v, w, e := NewSettlementService(p).Preview(ctx, "c2")
	if e != nil {
		t.Fatalf("missing-rule preview must not error: %v", e)
	}
	if v == nil {
		t.Fatal("missing-rule preview must return a settlement, got nil")
	}
	if v.SubsidyCents != 0 {
		t.Fatalf("missing-rule subsidy = %d, want 0", v.SubsidyCents)
	}
	if len(v.Lines) != 0 {
		t.Fatalf("missing-rule lines = %v, want none", v.Lines)
	}
	if len(w) == 0 {
		t.Fatal("missing-rule preview must carry an explanatory warning")
	}
}

// TestSettlementPreviewWithRuleComputesCorrectSubsidy 覆盖正常规则回归场景：
// 首段阈值内按基础比例、超阈值按超限比例，年度额度未触顶时补助与明细正确。
func TestSettlementPreviewWithRuleComputesCorrectSubsidy(t *testing.T) {
	p := NewPortsForTest()
	ctx := context.Background()
	prj := domain.GrantProject{ID: "p3", Code: "P3", Name: "公益三", Year: 2026, AnnualCapCents: 100000}
	if e := p.Projects.Create(ctx, &prj); e != nil {
		t.Fatal(e)
	}
	ben := domain.Beneficiary{ID: "b3", Name: "对象三", IdentityHash: "hash3", Active: true}
	if e := p.Beneficiaries.Create(ctx, &ben); e != nil {
		t.Fatal(e)
	}
	rule := domain.RuleVersion{ID: "r3", ProjectID: "p3", SchemeID: "s3", Name: "2026", Year: 2026, BaseRateBps: 5000, ExcessRateBps: 2500, ThresholdCents: 1000, CapCents: 100000, Status: domain.RulePublished, VersionNo: 1}
	if e := p.Rules.Create(ctx, &rule); e != nil {
		t.Fatal(e)
	}
	claim := domain.ExpenseClaim{ID: "c3", ProjectID: "p3", BeneficiaryID: "b3", SchemeID: "s3", Year: 2026, Summary: "医疗", Category: "medical", AmountCents: 2000, Status: domain.ClaimSubmitted, IdempotencyKey: "k3"}
	if e := p.Claims.Create(ctx, &claim); e != nil {
		t.Fatal(e)
	}
	v, w, e := NewSettlementService(p).Preview(ctx, "c3")
	if e != nil {
		t.Fatalf("normal preview error: %v", e)
	}
	// 2000 分：首段 1000*5000/10000=500，超限段 1000*2500/10000=250，合计 750。
	if v.SubsidyCents != 750 {
		t.Fatalf("subsidy = %d, want 750", v.SubsidyCents)
	}
	if v.RuleID != "r3" {
		t.Fatalf("rule id = %q, want r3", v.RuleID)
	}
	if len(v.Lines) != 1 {
		t.Fatalf("lines = %d, want 1", len(v.Lines))
	}
	ln := v.Lines[0]
	if ln.EligibleCents != 2000 || ln.RateBps != 5000 || ln.SubsidyCents != 750 || ln.RemainingCents != 100000 {
		t.Fatalf("unexpected line detail: %+v", ln)
	}
	if len(w) != 0 {
		t.Fatalf("unexpected warnings: %v", w)
	}
}
