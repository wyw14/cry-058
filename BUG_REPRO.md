# Bug reproduction

go test ./internal/service/calculator -run '^TestRoundBpsUsesNearestCent$' -count=1
