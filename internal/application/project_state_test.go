package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/repository/memory"
	"testing"
)

func TestClosedProjectCannotBeActivated(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Projects.Create(ctx, &domain.GrantProject{ID: "p", Code: "P", Name: "公益", Year: 2026, AnnualCapCents: 1000, Status: domain.ProjectClosed})
	if _, e := NewProjectService(p).Activate(ctx, "p"); e == nil {
		t.Fatal("closed project activated")
	}
}

func TestDraftProjectCanBeActivated(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Projects.Create(ctx, &domain.GrantProject{ID: "p", Code: "P", Name: "公益", Year: 2026, AnnualCapCents: 1000, Status: domain.ProjectDraft})
	v, e := NewProjectService(p).Activate(ctx, "p")
	if e != nil {
		t.Fatalf("draft project should be activatable: %v", e)
	}
	if v.Status != domain.ProjectActive {
		t.Fatalf("expected active, got %s", v.Status)
	}
}

func TestActiveProjectCanBeReactivated(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Projects.Create(ctx, &domain.GrantProject{ID: "p", Code: "P", Name: "公益", Year: 2026, AnnualCapCents: 1000, Status: domain.ProjectActive})
	if _, e := NewProjectService(p).Activate(ctx, "p"); e != nil {
		t.Fatalf("active project should remain activatable: %v", e)
	}
}
