package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/service/id"
)

type SchemeService struct{ p *domain.Ports }

func NewSchemeService(p *domain.Ports) *SchemeService { return &SchemeService{p} }
func (s *SchemeService) Create(ctx context.Context, v domain.CoverageScheme) (*domain.CoverageScheme, error) {
	if e := v.Validate(); e != nil {
		return nil, e
	}
	v.ID = id.New("sch")
	v.Active = true
	return &v, s.p.Schemes.Create(ctx, &v)
}
func (s *SchemeService) Disable(ctx context.Context, idv string) (*domain.CoverageScheme, error) {
	v, e := s.p.Schemes.Get(ctx, idv)
	if e != nil {
		return nil, e
	}
	v.Active = false
	v.Version++
	return v, s.p.Schemes.Update(ctx, v)
}
