package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/repository/memory"
	"testing"
)

// TestReturnedClaimCannotBeApprovedDirectly ensures a returned claim does NOT
// enter the approve branch without being resubmitted first. The earlier guard
// `c.Status != ClaimSubmitted && c.Status != ClaimReturned` wrongly admitted
// ClaimReturned straight to ClaimApproved.
func TestReturnedClaimCannotBeApprovedDirectly(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Claims.Create(ctx, &domain.ExpenseClaim{ID: "c", Status: domain.ClaimReturned, Version: 1})
	_, e := NewReviewService(p).Decide(ctx, "c", domain.Actor{ID: "a", Role: "admin"}, true, "ok")
	if e == nil {
		t.Fatal("returned claim approved without resubmission")
	}
}

// TestReturnedClaimCannotBeRejectedEither ensures a returned claim cannot be
// acted on by review at all — it must be resubmitted before any review decision.
func TestReturnedClaimCannotBeRejectedEither(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Claims.Create(ctx, &domain.ExpenseClaim{ID: "c", Status: domain.ClaimReturned, Version: 1})
	_, e := NewReviewService(p).Decide(ctx, "c", domain.Actor{ID: "a", Role: "admin"}, false, "退回")
	if e == nil {
		t.Fatal("returned claim rejected again without resubmission")
	}
}

// TestResubmitReturnsClaimToSubmittedThenApprove covers the happy path:
// submitted -> returned -> resubmit -> submitted -> approved.
func TestResubmitReturnsClaimToSubmittedThenApprove(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	admin := domain.Actor{ID: "a", Role: "admin"}

	// start submitted
	c := &domain.ExpenseClaim{ID: "c", ProjectID: "p", BeneficiaryID: "b", SchemeID: "s", Year: 2026, Summary: "x", AmountCents: 100, Status: domain.ClaimSubmitted, Version: 1}
	if e := p.Claims.Create(ctx, c); e != nil {
		t.Fatal(e)
	}
	// admin returns it (normal return flow intact)
	rev := NewReviewService(p)
	if _, e := rev.Decide(ctx, "c", admin, false, "需补充材料"); e != nil {
		t.Fatalf("return: %v", e)
	}
	if got, _ := p.Claims.Get(ctx, "c"); got.Status != domain.ClaimReturned {
		t.Fatalf("after return = %s, want returned", got.Status)
	}
	// resubmit must move returned -> submitted
	cls := NewClaimService(p)
	if _, e := cls.Resubmit(ctx, "c"); e != nil {
		t.Fatalf("resubmit: %v", e)
	}
	if got, _ := p.Claims.Get(ctx, "c"); got.Status != domain.ClaimSubmitted {
		t.Fatalf("after resubmit = %s, want submitted", got.Status)
	}
	// now approval is allowed
	if _, e := rev.Decide(ctx, "c", admin, true, "ok"); e != nil {
		t.Fatalf("approve after resubmit: %v", e)
	}
	if got, _ := p.Claims.Get(ctx, "c"); got.Status != domain.ClaimApproved {
		t.Fatalf("after approve = %s, want approved", got.Status)
	}
}

// TestResubmitOnlyFromReturned ensures resubmit is gated to returned claims,
// so a fresh submitted claim cannot be "resubmitted" to bump its version.
func TestResubmitOnlyFromReturned(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Claims.Create(ctx, &domain.ExpenseClaim{ID: "c", Status: domain.ClaimSubmitted, Version: 1})
	if _, e := NewClaimService(p).Resubmit(ctx, "c"); e == nil {
		t.Fatal("resubmit allowed on non-returned claim")
	}
}

// TestSubmittedClaimApprovesNormally is the regression guard for the normal
// approve branch — only submitted claims may be approved, but they still can.
func TestSubmittedClaimApprovesNormally(t *testing.T) {
	p := memory.NewPorts()
	ctx := context.Background()
	_ = p.Claims.Create(ctx, &domain.ExpenseClaim{ID: "c", Status: domain.ClaimSubmitted, Version: 1})
	c, e := NewReviewService(p).Decide(ctx, "c", domain.Actor{ID: "a", Role: "admin"}, true, "")
	if e != nil {
		t.Fatalf("submitted approve: %v", e)
	}
	if c.Status != domain.ClaimApproved {
		t.Fatalf("status = %s, want approved", c.Status)
	}
}
