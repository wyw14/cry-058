package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/repository/memory"
	"sync"
	"testing"
)

func TestConcurrentLedgerWritesRejectStaleVersion(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	if e := p.Settlements.SaveLedger(ctx, &domain.AnnualLedger{ProjectID: "p", BeneficiaryID: "b", Year: 2026, Version: 1}); e != nil {
		t.Fatal(e)
	}
	ready := make(chan struct{}, 2)
	start := make(chan struct{})
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l, _ := p.Settlements.Ledger(ctx, "p", "b", 2026)
			ready <- struct{}{}
			<-start
			l.PaidCents += 100
			l.Version++
			errs <- p.Settlements.SaveLedger(ctx, l)
		}()
	}
	<-ready
	<-ready
	close(start)
	wg.Wait()
	close(errs)
	ok := 0
	for e := range errs {
		if e == nil {
			ok++
		}
	}
	if ok != 1 {
		t.Fatalf("expected one stale write rejection, successful writes=%d", ok)
	}
}
