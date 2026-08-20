package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/service/id"
	"time"
)

type RuleService struct{ p *domain.Ports }

func NewRuleService(p *domain.Ports) *RuleService { return &RuleService{p} }
func (s *RuleService) Publish(ctx context.Context, v domain.RuleVersion) (*domain.RuleVersion, error) {
	if e := v.Validate(); e != nil {
		return nil, e
	}
	if _, e := s.p.Projects.Get(ctx, v.ProjectID); e != nil {
		return nil, e
	}
	v.ID = id.New("rule")
	v.Status = domain.RulePublished
	v.PublishedAt = time.Now().UTC()
	v.CreatedAt = v.PublishedAt
	return &v, s.p.Rules.Create(ctx, &v)
}
func (s *RuleService) History(ctx context.Context, p string, y int) ([]*domain.RuleVersion, error) {
	return s.p.Rules.List(ctx, p, y)
}
