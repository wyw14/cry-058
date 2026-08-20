package domain

import "context"

type ProjectRepo interface {
	Create(context.Context, *GrantProject) error
	Get(context.Context, string) (*GrantProject, error)
	List(context.Context, int) ([]*GrantProject, error)
	Update(context.Context, *GrantProject) error
}
type BeneficiaryRepo interface {
	Create(context.Context, *Beneficiary) error
	Get(context.Context, string) (*Beneficiary, error)
}
type SchemeRepo interface {
	Create(context.Context, *CoverageScheme) error
	Get(context.Context, string) (*CoverageScheme, error)
	Update(context.Context, *CoverageScheme) error
}
type ClaimRepo interface {
	Create(context.Context, *ExpenseClaim) error
	Get(context.Context, string) (*ExpenseClaim, error)
	GetByIdempotency(context.Context, string) (*ExpenseClaim, error)
	Update(context.Context, *ExpenseClaim) error
	ListByYear(context.Context, string, int) ([]*ExpenseClaim, error)
}
type RuleRepo interface {
	Create(context.Context, *RuleVersion) error
	Get(context.Context, string) (*RuleVersion, error)
	Latest(context.Context, string, int) (*RuleVersion, error)
	List(context.Context, string, int) ([]*RuleVersion, error)
}
type SettlementRepo interface {
	Create(context.Context, *Settlement) error
	Get(context.Context, string) (*Settlement, error)
	GetByIdempotency(context.Context, string) (*Settlement, error)
	Update(context.Context, *Settlement) error
	Ledger(context.Context, string, string, int) (*AnnualLedger, error)
	SaveLedger(context.Context, *AnnualLedger) error
	ListByYear(context.Context, string, int) ([]*Settlement, error)
}
type AuditRepo interface {
	Append(context.Context, *AuditEvent) error
	List(context.Context, string, int) ([]*AuditEvent, error)
}
type Ports struct {
	Projects      ProjectRepo
	Beneficiaries BeneficiaryRepo
	Schemes       SchemeRepo
	Claims        ClaimRepo
	Rules         RuleRepo
	Settlements   SettlementRepo
	Audits        AuditRepo
}
