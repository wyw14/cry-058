package domain

import "time"

type RuleStatus string

const (
	RuleDraft     RuleStatus = "draft"
	RulePublished RuleStatus = "published"
	RuleRetired   RuleStatus = "retired"
)

type RuleVersion struct {
	ID, ProjectID, SchemeID, Name, Condition string
	Year                                     int
	BaseRateBps, ExcessRateBps               int
	ThresholdCents, CapCents                 int64
	Status                                   RuleStatus
	VersionNo                                int
	PublishedAt                              time.Time
	CreatedAt                                time.Time
}

func (r RuleVersion) Validate() error {
	if r.ProjectID == "" || r.SchemeID == "" || r.Name == "" {
		return Invalid("规则关联字段不能为空", "rule")
	}
	if r.VersionNo < 1 {
		return Invalid("版本号必须从 1 开始", "version_no")
	}
	if r.BaseRateBps < 0 || r.BaseRateBps > 10000 || r.ExcessRateBps < 0 || r.ExcessRateBps > 10000 {
		return Invalid("规则比例无效", "rate")
	}
	return nil
}
