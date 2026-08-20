package application

import (
	"context"
	"testing"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/repository/memory"
)

func TestGrantSettlementWorkflowUsesPublishedRule(t *testing.T) {
	p:=NewPortsForTest();ctx:=context.Background()
	prj:=domain.GrantProject{ID:"p",Code:"P",Name:"公益",Year:2026,AnnualCapCents:100000};if e:=p.Projects.Create(ctx,&prj);e!=nil{t.Fatal(e)}
	ben:=domain.Beneficiary{ID:"b",Name:"对象",IdentityHash:"hash",Active:true};if e:=p.Beneficiaries.Create(ctx,&ben);e!=nil{t.Fatal(e)}
	sch:=domain.CoverageScheme{ID:"s",ProjectID:"p",Name:"基础",BaseRateBps:5000,ExcessRateBps:2500,ThresholdCents:1000,CapCents:100000,Active:true};if e:=p.Schemes.Create(ctx,&sch);e!=nil{t.Fatal(e)}
	rule:=domain.RuleVersion{ID:"r",ProjectID:"p",SchemeID:"s",Name:"2026",Year:2026,BaseRateBps:5000,ExcessRateBps:2500,ThresholdCents:1000,CapCents:100000,Status:domain.RulePublished,VersionNo:1};if e:=p.Rules.Create(ctx,&rule);e!=nil{t.Fatal(e)}
	claim:=domain.ExpenseClaim{ID:"c",ProjectID:"p",BeneficiaryID:"b",SchemeID:"s",Year:2026,Summary:"医疗",Category:"medical",AmountCents:2000,Status:domain.ClaimSubmitted,IdempotencyKey:"k"};if e:=p.Claims.Create(ctx,&claim);e!=nil{t.Fatal(e)}
	s,_,e:=NewSettlementService(p).Preview(ctx,"c");if e!=nil{t.Fatal(e)};if s.SubsidyCents<=0{t.Fatal("expected subsidy")}
}

func NewPortsForTest()*domain.Ports{return memory.NewPorts()}
