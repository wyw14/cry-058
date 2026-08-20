package domain

import "time"

type SettlementStatus string

const (
	SettlementPreview   SettlementStatus = "preview"
	SettlementPending   SettlementStatus = "pending_review"
	SettlementConfirmed SettlementStatus = "confirmed"
	SettlementRevoked   SettlementStatus = "revoked"
	SettlementCorrected SettlementStatus = "corrected"
)

type SettlementLine struct {
	Category                                             string
	EligibleCents, RateBps, SubsidyCents, RemainingCents int64
}
type Settlement struct {
	ID, ClaimID, ProjectID, BeneficiaryID, RuleID, IdempotencyKey string
	Year                                                          int
	ClaimedCents, SubsidyCents                                    int64
	Status                                                        SettlementStatus
	Lines                                                         []SettlementLine
	CreatedAt, ConfirmedAt                                        time.Time
	Version                                                       int
}
type AnnualLedger struct {
	ProjectID, BeneficiaryID                string
	Year                                    int
	ReservedCents, PaidCents, ReversedCents int64
	Version                                 int
	UpdatedAt                               time.Time
}
type ReviewDecision struct {
	ClaimID, ReviewerID, Comment string
	Approved                     bool
	At                           time.Time
}
