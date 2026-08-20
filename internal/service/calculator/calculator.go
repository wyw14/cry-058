package calculator

import "github.com/wyw14/cry058/internal/domain"

type Result struct {
	SubsidyCents int64
	Lines        []domain.SettlementLine
	Warnings     []string
}

// nilRuleWarning 标识未匹配到规则的场景：补助必须安全归零并给出可解释提示，
// 而不能触发空指针崩溃使调用方拿不到任何结果。
const nilRuleWarning = "未匹配到适用规则，本次试算补助为零"

func Calculate(amount int64, rule *domain.RuleVersion, used int64) Result {
	if rule == nil {
		return Result{SubsidyCents: 0, Lines: nil, Warnings: []string{nilRuleWarning}}
	}
	if amount < 0 {
		amount = 0
	}
	first := amount
	if rule.ThresholdCents > 0 && first > rule.ThresholdCents {
		first = rule.ThresholdCents
	}
	second := amount - first
	subsidy := first*int64(rule.BaseRateBps)/10000 + second*int64(rule.ExcessRateBps)/10000
	remaining := rule.CapCents - used
	if remaining < 0 {
		remaining = 0
	}
	if subsidy > remaining {
		subsidy = remaining
	}
	lines := []domain.SettlementLine{{Category: "eligible", EligibleCents: amount, RateBps: int64(rule.BaseRateBps), SubsidyCents: subsidy, RemainingCents: remaining}}
	w := []string{}
	if remaining == 0 {
		w = append(w, "年度额度已用尽")
	}
	return Result{SubsidyCents: subsidy, Lines: lines, Warnings: w}
}
