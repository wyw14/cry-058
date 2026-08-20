package calculator

import "github.com/wyw14/cry058/internal/domain"

type Policy struct {
	MinClaimCents     int64
	MaxClaimCents     int64
	AllowedCategories map[string]bool
	RequireReceipt    bool
}

func DefaultPolicy() Policy {
	return Policy{MinClaimCents: 1, MaxClaimCents: 5000000, AllowedCategories: map[string]bool{"medical": true, "education": true, "housing": true}, RequireReceipt: true}
}
func (p Policy) ValidateClaim(c domain.ExpenseClaim, receipts int) error {
	if c.AmountCents < p.MinClaimCents {
		return domain.Invalid("申报金额低于政策下限", "amount_cents")
	}
	if p.MaxClaimCents > 0 && c.AmountCents > p.MaxClaimCents {
		return domain.Invalid("申报金额超过政策上限", "amount_cents")
	}
	if len(p.AllowedCategories) > 0 && !p.AllowedCategories[c.Category] {
		return domain.Invalid("申报类别不在政策范围", "category")
	}
	if p.RequireReceipt && receipts == 0 {
		return domain.Invalid("申报必须包含凭证", "receipt")
	}
	return nil
}
func (p Policy) Eligible(amount, used int64) bool {
	return amount > 0 && used >= 0 && amount+used <= p.MaxClaimCents
}
func EffectiveRate(rule *domain.RuleVersion, amount int64) int64 {
	if rule == nil {
		return 0
	}
	if rule.ThresholdCents > 0 && amount > rule.ThresholdCents {
		return int64(rule.ExcessRateBps)
	}
	return int64(rule.BaseRateBps)
}
