package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/repository/memory"
	"testing"
)

func TestAnnualReportExcludesRevokedSubsidy(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Settlements.Create(ctx, &domain.Settlement{ID: "confirmed", ProjectID: "p", Year: 2026, Status: domain.SettlementConfirmed, SubsidyCents: 100})
	_ = p.Settlements.Create(ctx, &domain.Settlement{ID: "revoked", ProjectID: "p", Year: 2026, Status: domain.SettlementRevoked, SubsidyCents: 40})
	r, e := NewReportService(p).Build(ctx, "p", 2026)
	if e != nil {
		t.Fatal(e)
	}
	if r.SubsidyCents != 100 || r.SettledCount != 1 {
		t.Fatalf("revoked subsidy included: %d", r.SubsidyCents)
	}
}
