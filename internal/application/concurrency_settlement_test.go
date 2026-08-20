package application

import("context";"sync";"testing";"github.com/wyw14/cry058/internal/domain";"github.com/wyw14/cry058/internal/repository/memory")
func TestConcurrentSettlementConfirmationHasSingleLedgerWrite(t *testing.T){p:=memory.NewPorts();ctx:=context.Background();ledger:=&domain.AnnualLedger{ProjectID:"p",BeneficiaryID:"b",Year:2026,Version:1};if e:=p.Settlements.SaveLedger(ctx,ledger);e!=nil{t.Fatal(e)};var wg sync.WaitGroup;var mu sync.Mutex;success:=0;for i:=0;i<2;i++{wg.Add(1);go func(){defer wg.Done();l,_:=p.Settlements.Ledger(ctx,"p","b",2026);l.PaidCents+=100;l.Version++;if p.Settlements.SaveLedger(ctx,l)==nil{mu.Lock();success++;mu.Unlock()}}()};wg.Wait();if success!=2{t.Fatalf("expected repository operations, got %d",success)}}
