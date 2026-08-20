package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
)

type ClaimQueryService struct{ p *domain.Ports }

func NewClaimQueryService(p *domain.Ports) *ClaimQueryService { return &ClaimQueryService{p} }
func (s *ClaimQueryService) ByYear(ctx context.Context, project string, year int) ([]*domain.ExpenseClaim, error) {
	return s.p.Claims.ListByYear(ctx, project, year)
}
func (s *ClaimQueryService) Get(ctx context.Context, id string) (*domain.ExpenseClaim, error) {
	return s.p.Claims.Get(ctx, id)
}
