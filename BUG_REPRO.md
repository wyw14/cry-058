# Bug reproduction

go test ./internal/service/attachment -run '^TestAttachmentStoreRejectsTraversal$' -count=1
