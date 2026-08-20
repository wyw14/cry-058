package calculator

import "github.com/wyw14/cry058/internal/domain"

type Variance struct {
	OldSubsidy  int64
	NewSubsidy  int64
	Delta       int64
	Explanation string
}

func Compare(amount, used int64, oldRule, newRule *domain.RuleVersion) Variance {
	a := Calculate(amount, oldRule, used)
	b := Calculate(amount, newRule, used)
	return Variance{OldSubsidy: a.SubsidyCents, NewSubsidy: b.SubsidyCents, Delta: b.SubsidyCents - a.SubsidyCents, Explanation: "按两个不可变规则版本分别试算"}
}
func Cap(value, cap int64) int64 {
	if cap > 0 && value > cap {
		return cap
	}
	return value
}
func Difference(a, b []domain.SettlementLine) int64 {
	var x, y int64
	for _, v := range a {
		x += v.SubsidyCents
	}
	for _, v := range b {
		y += v.SubsidyCents
	}
	return y - x
}
