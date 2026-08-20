package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"sync"
)

type BatchService struct {
	p  *domain.Ports
	mu sync.Mutex
}

func NewBatchService(p *domain.Ports) *BatchService { return &BatchService{p: p} }

type BatchResult struct {
	Accepted int
	Rejected int
	Errors   []string
}

func (s *BatchService) Submit(ctx context.Context, claims []domain.ExpenseClaim) BatchResult {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := BatchResult{}
	for _, c := range claims {
		if _, e := NewClaimService(s.p).Submit(ctx, c); e != nil {
			out.Rejected++
			out.Errors = append(out.Errors, e.Error())
		} else {
			out.Accepted++
		}
	}
	return out
}
func (s *BatchService) ValidateClaims(claims []domain.ExpenseClaim) []error {
	errs := []error{}
	for _, c := range claims {
		if e := c.Validate(); e != nil {
			errs = append(errs, e)
		}
	}
	return errs
}
