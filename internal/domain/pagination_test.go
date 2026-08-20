package domain

import "testing"

func TestPageNormalizesNegativeValues(t *testing.T) {
	p := Page{Number: -1, Size: -2}.Normalize()
	if p.Number != 1 || p.Size != 20 {
		t.Fatalf("bad normalization: %#v", p)
	}
}
