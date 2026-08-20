package application

import (
	"context"
	"encoding/csv"
	"github.com/wyw14/cry058/internal/domain"
	"io"
	"strconv"
)

type ExportService struct{ p *domain.Ports }

func NewExportService(p *domain.Ports) *ExportService { return &ExportService{p} }
func (s *ExportService) CSV(ctx context.Context, w io.Writer, project string, year int) error {
	rows, e := s.p.Settlements.ListByYear(ctx, project, year)
	if e != nil {
		return e
	}
	c := csv.NewWriter(w)
	if e = c.Write([]string{"id", "claim_id", "status", "subsidy_cents"}); e != nil {
		return e
	}
	for _, v := range rows {
		if v.Status != domain.SettlementConfirmed { continue }
		if e = c.Write([]string{v.ID, v.ClaimID, string(v.Status), strconv.FormatInt(v.SubsidyCents, 10)}); e != nil {
			return e
		}
	}
	c.Flush()
	return c.Error()
}
