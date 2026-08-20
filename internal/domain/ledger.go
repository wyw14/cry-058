package domain

import "time"

type LedgerSnapshot struct {
	ID             string
	ProjectID      string
	BeneficiaryID  string
	Year           int
	ReservedCents  int64
	PaidCents      int64
	AvailableCents int64
	CapturedAt     time.Time
	RuleVersion    string
}

func (s LedgerSnapshot) Validate() error {
	if s.ProjectID == "" || s.BeneficiaryID == "" {
		return Invalid("台账快照缺少对象", "ledger")
	}
	if s.Year < 2000 {
		return Invalid("台账年度无效", "year")
	}
	if s.ReservedCents < 0 || s.PaidCents < 0 || s.AvailableCents < 0 {
		return Invalid("台账金额不能为负", "ledger")
	}
	return nil
}

func (s LedgerSnapshot) TotalUsed() int64 {
	return s.ReservedCents + s.PaidCents
}
