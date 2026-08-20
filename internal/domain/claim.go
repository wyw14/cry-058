package domain

import "time"

type ClaimStatus string

const (
	ClaimDraft     ClaimStatus = "draft"
	ClaimSubmitted ClaimStatus = "submitted"
	ClaimReturned  ClaimStatus = "returned"
	ClaimApproved  ClaimStatus = "approved"
	ClaimSettled   ClaimStatus = "settled"
	ClaimRevoked   ClaimStatus = "revoked"
)

type ExpenseClaim struct {
	ID, ProjectID, BeneficiaryID, SchemeID, Category, Summary, IdempotencyKey string
	Year                                                                      int
	OccurredOn                                                                time.Time
	AmountCents                                                               int64
	Status                                                                    ClaimStatus
	Version                                                                   int
	CreatedAt, UpdatedAt                                                      time.Time
}
type Receipt struct {
	ID, ClaimID, Digest, FileName, MediaType string
	Size                                     int64
	UploadedAt                               time.Time
}

func (c ExpenseClaim) Validate() error {
	if c.ProjectID == "" || c.BeneficiaryID == "" || c.SchemeID == "" {
		return Invalid("申报关联字段不能为空", "claim")
	}
	if c.Year < 2000 {
		return Invalid("年度无效", "year")
	}
	if c.AmountCents <= 0 {
		return Invalid("申报金额必须为正数", "amount_cents")
	}
	if c.Summary == "" {
		return Invalid("凭证摘要不能为空", "summary")
	}
	return nil
}
