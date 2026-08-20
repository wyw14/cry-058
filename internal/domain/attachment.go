package domain

import "time"

type Attachment struct {
	ID        string
	OwnerID   string
	FileName  string
	MediaType string
	Size      int64
	Digest    string
	Path      string
	CreatedAt time.Time
}

func (a Attachment) Validate() error {
	if a.ID == "" || a.OwnerID == "" {
		return Invalid("附件缺少归属", "owner_id")
	}
	if a.FileName == "" || a.Size <= 0 {
		return Invalid("附件名称和大小不能为空", "file")
	}
	if a.Size > 10*1024*1024 {
		return Invalid("附件超过 10MB 限制", "size")
	}
	switch a.MediaType {
	case "application/pdf", "image/jpeg", "image/png":
	default:
		return Invalid("附件类型不在白名单", "media_type")
	}
	return nil
}
