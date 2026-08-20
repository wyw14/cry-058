package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/service/calculator"
	"github.com/wyw14/cry058/internal/service/id"
	"time"
)

type SettlementService struct{ p *domain.Ports }

func NewSettlementService(p *domain.Ports) *SettlementService { return &SettlementService{p} }
func (s *SettlementService) Preview(ctx context.Context, claimID string) (*domain.Settlement, []string, error) {
	c, e := s.p.Claims.Get(ctx, claimID)
	if e != nil {
		return nil, nil, e
	}
	r, e := s.p.Rules.Latest(ctx, c.ProjectID, c.Year)
	if e != nil {
		return nil, nil, e
	}
	ledger, _ := s.p.Settlements.Ledger(ctx, c.ProjectID, c.BeneficiaryID, c.Year)
	res := calculator.Calculate(c.AmountCents, r, ledger.PaidCents+ledger.ReservedCents)
	v := &domain.Settlement{ID: id.New("set"), ClaimID: c.ID, ProjectID: c.ProjectID, BeneficiaryID: c.BeneficiaryID, RuleID: r.ID, Year: c.Year, ClaimedCents: c.AmountCents, SubsidyCents: res.SubsidyCents, Lines: res.Lines, Status: domain.SettlementPreview, CreatedAt: time.Now().UTC()}
	return v, res.Warnings, nil
}
func (s *SettlementService) Confirm(ctx context.Context, v *domain.Settlement, actor domain.Actor) (*domain.Settlement, error) {
	if !actor.IsAdmin() {
		return nil, domain.Forbidden("只有管理员可确认结算")
	}
	if v.Status != domain.SettlementPreview && v.Status != domain.SettlementPending {
		return nil, domain.StateError("当前状态不能确认")
	}
	ledger, e := s.p.Settlements.Ledger(ctx, v.ProjectID, v.BeneficiaryID, v.Year)
	if e != nil {
		return nil, e
	}
	if ledger.PaidCents+ledger.ReservedCents+v.SubsidyCents > 0 && ledger.PaidCents+ledger.ReservedCents+v.SubsidyCents > v.Lines[0].RemainingCents+ledger.PaidCents+ledger.ReservedCents {
		return nil, domain.Conflict("年度累计额度不足")
	}
	v.Status = domain.SettlementConfirmed
	v.ConfirmedAt = time.Now().UTC()
	v.Version++
	ledger.PaidCents += v.SubsidyCents
	ledger.Version++
	ledger.UpdatedAt = v.ConfirmedAt
	if e = s.p.Settlements.Create(ctx, v); e != nil {
		return nil, e
	}
	if e = s.p.Settlements.SaveLedger(ctx, ledger); e != nil {
		return nil, e
	}
	return v, nil
}
func (s *SettlementService) Revoke(ctx context.Context, idv string, actor domain.Actor) (*domain.Settlement, error) {
	if !actor.IsAdmin() {
		return nil, domain.Forbidden("只有管理员可撤销")
	}
	v, e := s.p.Settlements.Get(ctx, idv)
	if e != nil {
		return nil, e
	}
	if v.Status != domain.SettlementConfirmed {
		return nil, domain.StateError("只有已确认结算可撤销")
	}
	ledger, e := s.p.Settlements.Ledger(ctx, v.ProjectID, v.BeneficiaryID, v.Year)
	if e != nil {
		return nil, e
	}
	if ledger.PaidCents < v.SubsidyCents {
		return nil, domain.Conflict("累计账异常")
	}
	ledger.PaidCents -= v.SubsidyCents
	ledger.ReversedCents += v.SubsidyCents
	ledger.Version++
	_ = s.p.Settlements.SaveLedger(ctx, ledger)
	v.Status = domain.SettlementRevoked
	v.Version++
	return v, s.p.Settlements.Update(ctx, v)
}
