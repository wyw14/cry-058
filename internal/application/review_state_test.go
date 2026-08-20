package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/repository/memory"
	"testing"
)

func TestReturnedClaimRequiresNewSubmissionBeforeApproval(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Claims.Create(ctx, &domain.ExpenseClaim{ID: "c", Status: domain.ClaimReturned, Version: 1})
	_, e := NewReviewService(p).Decide(ctx, "c", domain.Actor{ID: "a", Role: "admin"}, true, "ok")
	if e == nil {
		t.Fatal("returned claim approved without resubmission")
	}
}
