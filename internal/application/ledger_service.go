package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"time"
)

type LedgerService struct{ p *domain.Ports }

func NewLedgerService(p *domain.Ports) *LedgerService { return &LedgerService{p} }
func (s *LedgerService) Get(ctx context.Context, project, beneficiary string, year int) (*domain.AnnualLedger, error) {
	return s.p.Settlements.Ledger(ctx, project, beneficiary, year)
}
func (s *LedgerService) Snapshot(ctx context.Context, project, beneficiary string, year int) (domain.LedgerSnapshot, error) {
	l, e := s.Get(ctx, project, beneficiary, year)
	if e != nil {
		return domain.LedgerSnapshot{}, e
	}
	return domain.LedgerSnapshot{ID: project + "/" + beneficiary, ProjectID: project, BeneficiaryID: beneficiary, Year: year, ReservedCents: l.ReservedCents, PaidCents: l.PaidCents, AvailableCents: l.ReservedCents + l.PaidCents, CapturedAt: time.Now().UTC()}, nil
}
func (s *LedgerService) CanReserve(ctx context.Context, project, beneficiary string, year int, amount, cap int64) (bool, error) {
	l, e := s.Get(ctx, project, beneficiary, year)
	if e != nil {
		return false, e
	}
	return l.PaidCents+l.ReservedCents+amount <= cap, nil
}
