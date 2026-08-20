package memory

import (
	"context"
	"errors"
	"github.com/wyw14/cry058/internal/domain"
)

type Tx struct {
	store     *Store
	committed bool
}

func (s *Store) Begin(ctx context.Context) (*Tx, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return &Tx{store: s}, nil
}
func (t *Tx) Commit() error {
	if t.committed {
		return errors.New("transaction already committed")
	}
	t.committed = true
	return nil
}
func (t *Tx) Rollback() error {
	if t.committed {
		return errors.New("cannot rollback committed transaction")
	}
	t.committed = true
	return nil
}
func (t *Tx) Confirm(v *domain.Settlement, l *domain.AnnualLedger) error {
	if v == nil || l == nil {
		return errors.New("missing settlement or ledger")
	}
	if v.Status != domain.SettlementConfirmed {
		return domain.StateError("settlement not confirmed")
	}
	if e := t.store.SaveLedger(context.Background(), l); e != nil {
		return e
	}
	return t.store.SaveSettlement(context.Background(), v)
}
func (s *Store) SaveLedger(ctx context.Context, v *domain.AnnualLedger) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := v.ProjectID + ":" + v.BeneficiaryID + ":" + timeKey(v.Year)
	if old, ok := s.ledgers[key]; ok && old.Version != v.Version-1 {
		return domain.Conflict("年度台账版本冲突")
	}
	s.ledgers[key] = cloneLedger(v)
	return nil
}
func (s *Store) SaveSettlement(ctx context.Context, v *domain.Settlement) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.settlements[v.ID] = cloneSettlement(v)
	return nil
}
