package memory

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"sort"
)

type Query struct{ store *Store }

func NewQuery(s *Store) *Query { return &Query{store: s} }
func (q *Query) ClaimsForBeneficiary(ctx context.Context, id string, year int) ([]*domain.ExpenseClaim, error) {
	q.store.mu.RLock()
	defer q.store.mu.RUnlock()
	out := []*domain.ExpenseClaim{}
	for _, c := range q.store.claims {
		if c.BeneficiaryID == id && c.Year == year {
			out = append(out, cloneClaim(c))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OccurredOn.Before(out[j].OccurredOn) })
	return out, nil
}
func (q *Query) RulesForScheme(ctx context.Context, id string, year int) ([]*domain.RuleVersion, error) {
	q.store.mu.RLock()
	defer q.store.mu.RUnlock()
	out := []*domain.RuleVersion{}
	for _, r := range q.store.rules {
		if r.SchemeID == id && r.Year == year {
			out = append(out, cloneRule(r))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].VersionNo < out[j].VersionNo })
	return out, nil
}
func (q *Query) OpenSettlements(ctx context.Context, project string, year int) ([]*domain.Settlement, error) {
	q.store.mu.RLock()
	defer q.store.mu.RUnlock()
	out := []*domain.Settlement{}
	for _, v := range q.store.settlements {
		if v.ProjectID == project && v.Year == year && v.Status != domain.SettlementRevoked {
			out = append(out, cloneSettlement(v))
		}
	}
	return out, nil
}
func (q *Query) AuditCount(ctx context.Context, entityID string) int {
	q.store.mu.RLock()
	defer q.store.mu.RUnlock()
	n := 0
	for _, a := range q.store.audits {
		if a.EntityID == entityID {
			n++
		}
	}
	return n
}
