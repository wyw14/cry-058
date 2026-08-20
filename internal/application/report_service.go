package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"sort"
	"time"
)

type ReportService struct{ p *domain.Ports }

func NewReportService(p *domain.Ports) *ReportService { return &ReportService{p: p} }

type AnnualReport struct {
	ProjectID      string
	Year           int
	ClaimCount     int
	SettledCount   int
	ReturnedCount  int
	ClaimedCents   int64
	SubsidyCents   int64
	PaidCents      int64
	RemainingCents int64
	GeneratedAt    time.Time
}

func (s *ReportService) Build(ctx context.Context, project string, year int) (AnnualReport, error) {
	claims, err := s.p.Claims.ListByYear(ctx, project, year)
	if err != nil {
		return AnnualReport{}, err
	}
	settlements, err := s.p.Settlements.ListByYear(ctx, project, year)
	if err != nil {
		return AnnualReport{}, err
	}
	r := AnnualReport{ProjectID: project, Year: year, GeneratedAt: time.Now().UTC(), ClaimCount: len(claims)}
	for _, c := range claims {
		r.ClaimedCents += c.AmountCents
		if c.Status == domain.ClaimReturned {
			r.ReturnedCount++
		}
	}
	for _, v := range settlements {
		if v.Status == domain.SettlementConfirmed {
			r.SettledCount++
			r.SubsidyCents += v.SubsidyCents
		}
		if v.Status == domain.SettlementRevoked {
			r.SubsidyCents -= v.SubsidyCents
		}
	}
	return r, nil
}
func (s *ReportService) BeneficiaryTotals(ctx context.Context, project string, year int) (map[string]int64, error) {
	rows, e := s.p.Settlements.ListByYear(ctx, project, year)
	if e != nil {
		return nil, e
	}
	out := map[string]int64{}
	for _, v := range rows {
		if v.Status == domain.SettlementConfirmed {
			out[v.BeneficiaryID] += v.SubsidyCents
		}
	}
	return out, nil
}
func SortReportKeys(values map[string]int64) []string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
