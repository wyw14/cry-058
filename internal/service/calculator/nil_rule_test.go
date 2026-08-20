package calculator

import "testing"

func TestCalculateNilRuleReturnsWarning(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("nil rule must not panic: %v", r)
		}
	}()
	res := Calculate(100, nil, 0)
	if res.SubsidyCents != 0 {
		t.Fatalf("unexpected subsidy")
	}
}
