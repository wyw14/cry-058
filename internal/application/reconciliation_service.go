package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
)

type ReconciliationService struct{ p *domain.Ports }

func NewReconciliationService(p *domain.Ports) *ReconciliationService {
	return &ReconciliationService{p}
}
func (s *ReconciliationService) Year(ctx context.Context, project string, year int) (map[string]int64, error) {
	rows, e := s.p.Settlements.ListByYear(ctx, project, year)
	if e != nil {
		return nil, e
	}
	out := map[string]int64{}
	for _, v := range rows {
		out[string(v.Status)] += v.SubsidyCents
	}
	return out, nil
}
