package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/service/id"
	"time"
)

type AuditService struct{ p *domain.Ports }

func NewAuditService(p *domain.Ports) *AuditService { return &AuditService{p} }
func (s *AuditService) Record(ctx context.Context, a domain.Actor, action, entity, entityID, detail string) error {
	return s.p.Audits.Append(ctx, &domain.AuditEvent{ID: id.New("aud"), ActorID: a.ID, Action: action, Entity: entity, EntityID: entityID, Detail: detail, At: time.Now().UTC()})
}
func (s *AuditService) List(ctx context.Context, idv string, limit int) ([]*domain.AuditEvent, error) {
	return s.p.Audits.List(ctx, idv, limit)
}
