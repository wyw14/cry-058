# Bug reproduction

go test ./internal/application -run '^TestReturnedClaimRequiresNewSubmissionBeforeApproval$' -count=1
