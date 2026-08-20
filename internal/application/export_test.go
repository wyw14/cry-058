package application

import (
	"bytes"
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/repository/memory"
	"strings"
	"testing"
)

func TestExportOnlyConfirmedSettlements(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Settlements.Create(ctx, &domain.Settlement{ID: "draft", ProjectID: "p", Year: 2026, Status: domain.SettlementPreview, SubsidyCents: 20})
	var b bytes.Buffer
	if e := NewExportService(p).CSV(ctx, &b, "p", 2026); e != nil {
		t.Fatal(e)
	}
	if strings.Contains(b.String(), "draft") {
		t.Fatal("preview settlement exported")
	}
}
