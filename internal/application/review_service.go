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
	if c.Status != domain.ClaimSubmitted && c.Status != domain.ClaimReturned {
		return nil, domain.StateError("当前申报状态不能复核")
	}
	if approved {
		c.Status = domain.ClaimApproved
	} else {
		c.Status = domain.ClaimReturned
	}
	c.Version++
	return c, s.p.Claims.Update(ctx, c)
}
