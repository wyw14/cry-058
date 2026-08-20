package memory

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"time"
)

func (s *Store) AppendAudit(ctx context.Context, actor, action, entity, entityID, detail string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.audits = append(s.audits, &domain.AuditEvent{ID: entityID + "-" + action, ActorID: actor, Action: action, Entity: entity, EntityID: entityID, Detail: detail, At: time.Now().UTC()})
	return nil
}
func (s *Store) AuditSince(ctx context.Context, since time.Time) []*domain.AuditEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []*domain.AuditEvent{}
	for _, a := range s.audits {
		if a.At.After(since) {
			c := *a
			out = append(out, &c)
		}
	}
	return out
}
func (s *Store) ClearAudit() { s.mu.Lock(); defer s.mu.Unlock(); s.audits = nil }
