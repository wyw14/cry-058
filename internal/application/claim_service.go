package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/service/id"
	"time"
)

type ClaimService struct{ p *domain.Ports }

func NewClaimService(p *domain.Ports) *ClaimService { return &ClaimService{p} }
func (s *ClaimService) Submit(ctx context.Context, v domain.ExpenseClaim) (*domain.ExpenseClaim, error) {
	if e := v.Validate(); e != nil {
		return nil, e
	}
	if old, e := s.p.Claims.GetByIdempotency(ctx, v.IdempotencyKey); e == nil {
		return old, nil
	}
	if _, e := s.p.Projects.Get(ctx, v.ProjectID); e != nil {
		return nil, e
	}
	if b, e := s.p.Beneficiaries.Get(ctx, v.BeneficiaryID); e != nil || !b.Active {
		return nil, domain.Forbidden("对象不存在或已停用")
	}
	if _, e := s.p.Schemes.Get(ctx, v.SchemeID); e != nil {
		return nil, e
	}
	v.ID = id.New("clm")
	v.Status = domain.ClaimSubmitted
	v.CreatedAt = time.Now().UTC()
	e := s.p.Claims.Create(ctx, &v)
	return &v, e
}
func (s *ClaimService) Return(ctx context.Context, idv string, reason string) (*domain.ExpenseClaim, error) {
	v, e := s.p.Claims.Get(ctx, idv)
	if e != nil {
		return nil, e
	}
	if v.Status != domain.ClaimSubmitted {
		return nil, domain.StateError("只有已提交申报可退回")
	}
	v.Status = domain.ClaimReturned
	v.Version++
	return v, s.p.Claims.Update(ctx, v)
}
