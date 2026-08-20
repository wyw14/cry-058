package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
)

type RuleCompareService struct{ p *domain.Ports }

func NewRuleCompareService(p *domain.Ports) *RuleCompareService { return &RuleCompareService{p} }

type RuleDiff struct {
	RateDelta      int
	ThresholdDelta int64
	CapDelta       int64
}

func (s *RuleCompareService) Compare(ctx context.Context, oldID, newID string) (RuleDiff, error) {
	a, e := s.p.Rules.Get(ctx, oldID)
	if e != nil {
		return RuleDiff{}, e
	}
	b, e := s.p.Rules.Get(ctx, newID)
	if e != nil {
		return RuleDiff{}, e
	}
	return RuleDiff{RateDelta: b.BaseRateBps - a.BaseRateBps, ThresholdDelta: b.ThresholdCents - a.ThresholdCents, CapDelta: b.CapCents - a.CapCents}, nil
}
