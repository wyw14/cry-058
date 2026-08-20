package domain

import "time"

type AuditEvent struct {
	ID, ActorID, Action, Entity, EntityID, RequestID, Detail string
	At                                                       time.Time
}
type Actor struct{ ID, Role string }

func (a Actor) IsAdmin() bool { return a.Role == "admin" }
