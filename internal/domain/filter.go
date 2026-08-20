package domain

import "strings"

type Filter struct {
	Field    string
	Value    string
	Operator string
}

func (f Filter) Valid() bool {
	switch f.Operator {
	case "eq", "contains", "gte", "lte":
	default:
		return false
	}
	return f.Field != ""
}
func MatchText(value string, f Filter) bool {
	if !f.Valid() {
		return false
	}
	switch f.Operator {
	case "eq":
		return value == f.Value
	case "contains":
		return strings.Contains(strings.ToLower(value), strings.ToLower(f.Value))
	default:
		return false
	}
}
func NormalizeFilters(in []Filter) []Filter {
	out := []Filter{}
	for _, f := range in {
		if f.Valid() {
			f.Field = strings.TrimSpace(f.Field)
			f.Value = strings.TrimSpace(f.Value)
			out = append(out, f)
		}
	}
	return out
}
