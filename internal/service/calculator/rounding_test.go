package calculator

import "testing"

func TestRoundBpsUsesNearestCent(t *testing.T) {
	if got := RoundBps(1, 5000); got != 1 {
		t.Fatalf("got %d", got)
	}
	if got := RoundBps(3, 5000); got != 2 {
		t.Fatalf("got %d", got)
	}
	// 余数低于半数（5000）时向下舍入。
	if got := RoundBps(9, 4999); got != 4 {
		t.Fatalf("below half, got %d", got)
	}
	// 余数达到半数（5000）时进位，修复此前截断导致少一分钱的缺陷。
	if got := RoundBps(1, 5000); got != 1 {
		t.Fatalf("at half, got %d", got)
	}
	if got := RoundBps(7, 5000); got != 4 {
		t.Fatalf("at half 7, got %d", got)
	}
	// 余数超过半数时正常进位。
	if got := RoundBps(7, 6000); got != 4 {
		t.Fatalf("above half, got %d", got)
	}
	// 边界：输入或比率为零返回零。
	if got := RoundBps(0, 5000); got != 0 {
		t.Fatalf("zero cents, got %d", got)
	}
	if got := RoundBps(100, 0); got != 0 {
		t.Fatalf("zero bps, got %d", got)
	}
}
