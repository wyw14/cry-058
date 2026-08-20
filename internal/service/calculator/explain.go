package calculator

import (
	"fmt"
	"github.com/wyw14/cry058/internal/domain"
)

type Explanation struct {
	Title    string
	Steps    []string
	Warnings []string
}

func Explain(amount int64, rule *domain.RuleVersion, result Result) Explanation {
	steps := []string{
		fmt.Sprintf("申报金额 %d 分", amount),
		fmt.Sprintf("首段阈值 %d 分，比例 %d 基点", rule.ThresholdCents, rule.BaseRateBps),
		fmt.Sprintf("超限比例 %d 基点", rule.ExcessRateBps),
		fmt.Sprintf("计算补助 %d 分", result.SubsidyCents),
	}
	return Explanation{Title: "年度公益补助试算", Steps: steps, Warnings: append([]string(nil), result.Warnings...)}
}

func Clamp(value, low, high int64) int64 {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}
