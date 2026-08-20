package calculator

import (
	"testing"

	"github.com/wyw14/cry058/internal/domain"
)

// TestCalculateNilRuleReturnsZeroAndWarning 覆盖缺失规则场景：试算不得崩溃，
// 必须安全返回零补助，并给出明确、可解释的提示。
func TestCalculateNilRuleReturnsZeroAndWarning(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("nil rule must not panic: %v", r)
		}
	}()
	res := Calculate(2000, nil, 0)
	if res.SubsidyCents != 0 {
		t.Fatalf("nil rule subsidy = %d, want 0", res.SubsidyCents)
	}
	if len(res.Lines) != 0 {
		t.Fatalf("nil rule lines = %v, want none", res.Lines)
	}
	if len(res.Warnings) == 0 {
		t.Fatal("nil rule must produce an explanatory warning")
	}
}

// TestCalculateNilRuleExplainCoversNil 确保 Explain 在缺失规则下也能生成可读明细，且不崩溃。
func TestCalculateNilRuleExplainCoversNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Explain(nil rule) must not panic: %v", r)
		}
	}()
	res := Calculate(2000, nil, 0)
	ex := Explain(2000, nil, res)
	if ex.Title == "" {
		t.Fatal("explain title should be populated")
	}
	if len(ex.Steps) == 0 {
		t.Fatal("explain steps should be populated even without a rule")
	}
	if len(ex.Warnings) == 0 {
		t.Fatal("explain warnings should reflect the missing-rule situation")
	}
}

// TestCalculateNormalRuleSegmentMath 验证已有规则下的分段金额：阈值以内按基础比例，
// 超过阈值部分按超限比例，且年度额度封顶、明细行字段完整。
func TestCalculateNormalRuleSegmentMath(t *testing.T) {
	rule := &domain.RuleVersion{
		BaseRateBps:    5000,
		ExcessRateBps:  2500,
		ThresholdCents: 1000,
		CapCents:       100000,
	}
	// 2000 分：首段 1000*5000/10000=500，超限段 1000*2500/10000=250，合计 750。
	res := Calculate(2000, rule, 0)
	if res.SubsidyCents != 750 {
		t.Fatalf("subsidy = %d, want 750", res.SubsidyCents)
	}
	if len(res.Lines) != 1 {
		t.Fatalf("lines = %d, want 1", len(res.Lines))
	}
	ln := res.Lines[0]
	if ln.Category != "eligible" {
		t.Fatalf("category = %q, want eligible", ln.Category)
	}
	if ln.EligibleCents != 2000 {
		t.Fatalf("eligible = %d, want 2000", ln.EligibleCents)
	}
	if ln.RateBps != 5000 {
		t.Fatalf("rate = %d, want 5000", ln.RateBps)
	}
	if ln.SubsidyCents != 750 {
		t.Fatalf("line subsidy = %d, want 750", ln.SubsidyCents)
	}
	if ln.RemainingCents != 100000 {
		t.Fatalf("remaining = %d, want 100000", ln.RemainingCents)
	}
	if len(res.Warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", res.Warnings)
	}
}

// TestCalculateNormalRuleBelowThresholdOnly 验证未超阈值时只走基础比例。
func TestCalculateNormalRuleBelowThresholdOnly(t *testing.T) {
	rule := &domain.RuleVersion{
		BaseRateBps:    5000,
		ExcessRateBps:  2500,
		ThresholdCents: 1000,
		CapCents:       100000,
	}
	// 800 分：全部落在首段，800*5000/10000=400。
	res := Calculate(800, rule, 0)
	if res.SubsidyCents != 400 {
		t.Fatalf("subsidy = %d, want 400", res.SubsidyCents)
	}
}

// TestCalculateNormalRuleAnnualCapClamps 验证年度额度封顶：已用额度足够时补助被封顶归零并给出提示。
func TestCalculateNormalRuleAnnualCapClamps(t *testing.T) {
	rule := &domain.RuleVersion{
		BaseRateBps:    5000,
		ExcessRateBps:  2500,
		ThresholdCents: 1000,
		CapCents:       100000,
	}
	// 用尽年度额度：remaining = 100000 - 100000 = 0，补助应被封顶为 0。
	res := Calculate(2000, rule, 100000)
	if res.SubsidyCents != 0 {
		t.Fatalf("clamped subsidy = %d, want 0", res.SubsidyCents)
	}
	if len(res.Warnings) == 0 || res.Warnings[0] != "年度额度已用尽" {
		t.Fatalf("expected cap-exhausted warning, got %v", res.Warnings)
	}
	if len(res.Lines) != 1 || res.Lines[0].RemainingCents != 0 {
		t.Fatalf("line remaining should be 0 when cap exhausted, got %v", res.Lines)
	}
}

// TestCalculateNormalRulePartialCap 验证补助超过剩余额度时被夹到剩余额度。
func TestCalculateNormalRulePartialCap(t *testing.T) {
	rule := &domain.RuleVersion{
		BaseRateBps:    5000,
		ExcessRateBps:  2500,
		ThresholdCents: 1000,
		CapCents:       100000,
	}
	// 理论补助 750，但年度仅剩 300：应被夹到 300。
	res := Calculate(2000, rule, 99700)
	if res.SubsidyCents != 300 {
		t.Fatalf("partial-cap subsidy = %d, want 300", res.SubsidyCents)
	}
	if res.Lines[0].RemainingCents != 300 {
		t.Fatalf("remaining = %d, want 300", res.Lines[0].RemainingCents)
	}
}
