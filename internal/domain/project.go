package domain

import "time"

type ProjectStatus string

const (
	ProjectDraft  ProjectStatus = "draft"
	ProjectActive ProjectStatus = "active"
	ProjectClosed ProjectStatus = "closed"
)

type GrantProject struct {
	ID, Code, Name, OwnerID string
	Year                    int
	AnnualCapCents          int64
	Status                  ProjectStatus
	CreatedAt, UpdatedAt    time.Time
	Version                 int
}
type Beneficiary struct {
	ID, Name, IdentityHash, Region string
	Active                         bool
	CreatedAt                      time.Time
}
type CoverageScheme struct {
	ID, ProjectID, Name string
	BaseRateBps         int
	ExcessRateBps       int
	ThresholdCents      int64
	CapCents            int64
	Active              bool
	Version             int
}

func (p GrantProject) Validate() error {
	if p.ID == "" || p.Code == "" || p.Name == "" {
		return Invalid("项目标识和名称不能为空", "project")
	}
	if p.Year < 2000 || p.Year > 2200 {
		return Invalid("年度不在允许范围", "year")
	}
	if p.AnnualCapCents <= 0 {
		return Invalid("年度额度必须为正数", "annual_cap_cents")
	}
	return nil
}
func (b Beneficiary) Validate() error {
	if b.ID == "" || b.Name == "" || b.IdentityHash == "" {
		return Invalid("对象身份字段不能为空", "beneficiary")
	}
	return nil
}
func (s CoverageScheme) Validate() error {
	if s.ID == "" || s.ProjectID == "" || s.Name == "" {
		return Invalid("保障方案字段不能为空", "scheme")
	}
	if s.BaseRateBps < 0 || s.BaseRateBps > 10000 || s.ExcessRateBps < 0 || s.ExcessRateBps > 10000 {
		return Invalid("比例必须在 0-10000 基点", "rate")
	}
	if s.ThresholdCents < 0 || s.CapCents < 0 {
		return Invalid("阈值和封顶线不能为负", "limit")
	}
	return nil
}
