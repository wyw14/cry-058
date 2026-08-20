# Bug reproduction

go test ./internal/domain -run '^TestPageNormalizesNegativeValues$' -count=1
