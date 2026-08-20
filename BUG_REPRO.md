# Bug reproduction

go test ./internal/application -run '^TestRulePublishRejectsDuplicateVersion$' -count=1
