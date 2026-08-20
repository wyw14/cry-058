package application

import (
	"bufio"
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"io"
	"strings"
)

type ImportService struct{ claims *ClaimService }

func NewImportService(p *domain.Ports) *ImportService {
	return &ImportService{claims: NewClaimService(p)}
}
func (s *ImportService) CSV(ctx context.Context, r io.Reader, project string, year int) BatchResult {
	scan := bufio.NewScanner(r)
	items := []domain.ExpenseClaim{}
	for scan.Scan() {
		parts := strings.Split(scan.Text(), ",")
		if len(parts) < 4 {
			continue
		}
		items = append(items, domain.ExpenseClaim{ProjectID: project, Year: year, BeneficiaryID: parts[0], SchemeID: parts[1], Summary: parts[2], AmountCents: parseCents(parts[3]), IdempotencyKey: parts[0] + "/" + parts[2]})
	}
	return NewBatchService(s.claims.p).Submit(ctx, items)
}
func parseCents(v string) int64 {
	var n int64
	for _, r := range strings.TrimSpace(v) {
		if r >= '0' && r <= '9' {
			n = n*10 + int64(r-'0')
		}
	}
	return n
}
