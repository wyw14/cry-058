package domain

import "time"

type ReviewQueueItem struct {
	ID         string
	ClaimID    string
	ProjectID  string
	Priority   int
	Reason     string
	AssignedTo string
	CreatedAt  time.Time
	ResolvedAt *time.Time
}

func (i ReviewQueueItem) Open() bool {
	return i.ResolvedAt == nil
}

func (i ReviewQueueItem) Validate() error {
	if i.ID == "" || i.ClaimID == "" || i.ProjectID == "" {
		return Invalid("复核队列缺少申报关联", "review")
	}
	if i.Priority < 1 || i.Priority > 5 {
		return Invalid("优先级必须为 1-5", "priority")
	}
	return nil
}
