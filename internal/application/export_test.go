package application

import (
	"bytes"
	"context"
	"encoding/csv"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/repository/memory"
	"strings"
	"testing"
)

func TestExportOnlyConfirmedSettlements(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Settlements.Create(ctx, &domain.Settlement{ID: "draft", ProjectID: "p", Year: 2026, Status: domain.SettlementPreview, SubsidyCents: 20})
	_ = p.Settlements.Create(ctx, &domain.Settlement{ID: "pending", ProjectID: "p", Year: 2026, Status: domain.SettlementPending, SubsidyCents: 30})
	_ = p.Settlements.Create(ctx, &domain.Settlement{ID: "revoked", ProjectID: "p", Year: 2026, Status: domain.SettlementRevoked, SubsidyCents: 40})
	_ = p.Settlements.Create(ctx, &domain.Settlement{ID: "confirmed", ProjectID: "p", Year: 2026, Status: domain.SettlementConfirmed, SubsidyCents: 50, ClaimID: "c1"})
	var b bytes.Buffer
	if e := NewExportService(p).CSV(ctx, &b, "p", 2026); e != nil {
		t.Fatal(e)
	}
	out := b.String()
	for _, banned := range []string{"draft", "pending", "revoked"} {
		if strings.Contains(out, banned) {
			t.Fatalf("non-confirmed settlement exported: %s", banned)
		}
	}
	r := csv.NewReader(strings.NewReader(out))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("expected header + 1 confirmed row, got %d", len(records))
	}
	if len(records[0]) != 4 || records[0][0] != "id" {
		t.Fatalf("unexpected header: %v", records[0])
	}
	want := []string{"confirmed", "c1", string(domain.SettlementConfirmed), "50"}
	for i := range want {
		if records[1][i] != want[i] {
			t.Fatalf("row[%d]=%q want %q", i, records[1][i], want[i])
		}
	}
}
