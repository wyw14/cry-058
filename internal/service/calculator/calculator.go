package calculator

import "github.com/wyw14/cry058/internal/domain"

type Result struct {
	SubsidyCents int64
	Lines        []domain.SettlementLine
	Warnings     []string
}

func Calculate(amount int64, rule *domain.RuleVersion, used int64) Result {
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
