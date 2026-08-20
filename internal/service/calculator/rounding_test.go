package calculator

import "testing"

func TestRoundBpsUsesNearestCent(t *testing.T) {
	if got := RoundBps(1, 5000); got != 1 {
		t.Fatalf("got %d", got)
	}
	if got := RoundBps(3, 5000); got != 2 {
		t.Fatalf("got %d", got)
	}
}
