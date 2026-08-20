package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
)

type WorkflowService struct {
	claims      *ClaimService
	reviews     *ReviewService
	settlements *SettlementService
}

func NewWorkflowService(p *domain.Ports) *WorkflowService {
	return &WorkflowService{claims: NewClaimService(p), reviews: NewReviewService(p), settlements: NewSettlementService(p)}
}
func (s *WorkflowService) Approve(ctx context.Context, claim string, actor domain.Actor) (*domain.ExpenseClaim, error) {
	return s.reviews.Decide(ctx, claim, actor, true, "")
}
func (s *WorkflowService) Return(ctx context.Context, claim string, actor domain.Actor, reason string) (*domain.ExpenseClaim, error) {
	if !actor.IsAdmin() {
		return nil, domain.Forbidden("无权退回")
	}
	return s.claims.Return(ctx, claim, reason)
}
func (s *WorkflowService) Resubmit(ctx context.Context, claim string) (*domain.ExpenseClaim, error) {
	return s.claims.Resubmit(ctx, claim)
}
func (s *WorkflowService) Preview(ctx context.Context, claim string) (*domain.Settlement, []string, error) {
	return s.settlements.Preview(ctx, claim)
}
