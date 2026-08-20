package domain

import "time"

type SubmissionKind string

const (
	SubmissionExpense    SubmissionKind = "expense"
	SubmissionCorrection SubmissionKind = "correction"
	SubmissionReversal   SubmissionKind = "reversal"
)

type SubmissionEvent struct {
	ID          string
	ClaimID     string
	Kind        SubmissionKind
	ActorID     string
	Reason      string
	AmountCents int64
	At          time.Time
	RequestID   string
}

func (e SubmissionEvent) Validate() error {
	if e.ID == "" || e.ClaimID == "" || e.ActorID == "" {
		return Invalid("提交事件缺少关联字段", "event")
	}
	if e.AmountCents < 0 {
		return Invalid("事件金额不能为负", "amount_cents")
	}
	if e.Reason == "" && e.Kind != SubmissionExpense {
		return Invalid("更正或撤销必须填写原因", "reason")
	}
	return nil
}

func (e SubmissionEvent) IsFinancial() bool {
	return e.Kind == SubmissionExpense || e.Kind == SubmissionCorrection || e.Kind == SubmissionReversal
}
