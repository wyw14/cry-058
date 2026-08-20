# Bug reproduction

go test ./internal/service/calculator -run '^TestCalculateNilRuleReturnsWarning$' -count=1
