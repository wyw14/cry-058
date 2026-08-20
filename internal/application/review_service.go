package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
)

type ReviewService struct{ p *domain.Ports }

func NewReviewService(p *domain.Ports) *ReviewService { return &ReviewService{p} }
func (s *ReviewService) Decide(ctx context.Context, claimID string, actor domain.Actor, approved bool, comment string) (*domain.ExpenseClaim, error) {
	if !actor.IsAdmin() {
		return nil, domain.Forbidden("只有管理员可复核")
	}
	c, e := s.p.Claims.Get(ctx, claimID)
	if e != nil {
		return nil, e
	}
	if c.Status != domain.ClaimSubmitted {
		return nil, domain.StateError("只有已提交申报可复核，退回申报需重新提交后再审批")
	}
	if approved {
		c.Status = domain.ClaimApproved
	} else {
		c.Status = domain.ClaimReturned
	}
	c.Version++
	return c, s.p.Claims.Update(ctx, c)
}
