package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/service/id"
)

type BeneficiaryService struct{ p *domain.Ports }

func NewBeneficiaryService(p *domain.Ports) *BeneficiaryService { return &BeneficiaryService{p} }
func (s *BeneficiaryService) Create(ctx context.Context, v domain.Beneficiary) (*domain.Beneficiary, error) {
	if e := v.Validate(); e != nil {
		return nil, e
	}
	v.ID = id.New("ben")
	v.Active = true
	return &v, s.p.Beneficiaries.Create(ctx, &v)
}
func (s *BeneficiaryService) Get(ctx context.Context, idv string) (*domain.Beneficiary, error) {
	return s.p.Beneficiaries.Get(ctx, idv)
}
func (s *BeneficiaryService) Deactivate(ctx context.Context, idv string) error {
	v, e := s.p.Beneficiaries.Get(ctx, idv)
	if e != nil {
		return e
	}
	v.Active = false
	return nil
}
