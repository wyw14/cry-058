# Bug reproduction

go test ./internal/application -run '^TestClosedProjectCannotBeActivated$' -count=1
