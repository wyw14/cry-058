package memory

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"sort"
	"sync"
	"time"
)

type Store struct {
	mu            sync.RWMutex
	projects      map[string]*domain.GrantProject
	beneficiaries map[string]*domain.Beneficiary
	schemes       map[string]*domain.CoverageScheme
	claims        map[string]*domain.ExpenseClaim
	rules         map[string]*domain.RuleVersion
	settlements   map[string]*domain.Settlement
	ledgers       map[string]*domain.AnnualLedger
	audits        []*domain.AuditEvent
}

func NewStore() *Store {
	return &Store{projects: map[string]*domain.GrantProject{}, beneficiaries: map[string]*domain.Beneficiary{}, schemes: map[string]*domain.CoverageScheme{}, claims: map[string]*domain.ExpenseClaim{}, rules: map[string]*domain.RuleVersion{}, settlements: map[string]*domain.Settlement{}, ledgers: map[string]*domain.AnnualLedger{}}
}
func cloneProject(v *domain.GrantProject) *domain.GrantProject    { c := *v; return &c }
func cloneBeneficiary(v *domain.Beneficiary) *domain.Beneficiary  { c := *v; return &c }
func cloneScheme(v *domain.CoverageScheme) *domain.CoverageScheme { c := *v; return &c }
func cloneClaim(v *domain.ExpenseClaim) *domain.ExpenseClaim      { c := *v; return &c }
func cloneRule(v *domain.RuleVersion) *domain.RuleVersion         { c := *v; return &c }
func cloneSettlement(v *domain.Settlement) *domain.Settlement {
	c := *v
	c.Lines = append([]domain.SettlementLine(nil), v.Lines...)
	return &c
}
func cloneLedger(v *domain.AnnualLedger) *domain.AnnualLedger { c := *v; return &c }

type ProjectRepo struct{ s *Store }

func (r *ProjectRepo) Create(ctx context.Context, p *domain.GrantProject) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if _, ok := r.s.projects[p.ID]; ok {
		return domain.Conflict("项目已存在")
	}
	if p.Version == 0 {
		p.Version = 1
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	p.UpdatedAt = p.CreatedAt
	r.s.projects[p.ID] = cloneProject(p)
	return nil
}
func (r *ProjectRepo) Get(ctx context.Context, id string) (*domain.GrantProject, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	p, ok := r.s.projects[id]
	if !ok {
		return nil, domain.NotFound("项目", id)
	}
	return cloneProject(p), nil
}
func (r *ProjectRepo) List(ctx context.Context, limit int) ([]*domain.GrantProject, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	out := []*domain.GrantProject{}
	for _, p := range r.s.projects {
		out = append(out, cloneProject(p))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
func (r *ProjectRepo) Update(ctx context.Context, p *domain.GrantProject) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	old, ok := r.s.projects[p.ID]
	if !ok {
		return domain.NotFound("项目", p.ID)
	}
	if old.Version+1 != p.Version {
		return domain.Conflict("项目版本冲突")
	}
	p.UpdatedAt = time.Now().UTC()
	r.s.projects[p.ID] = cloneProject(p)
	return nil
}

type BeneficiaryRepo struct{ s *Store }

func (r *BeneficiaryRepo) Create(ctx context.Context, b *domain.Beneficiary) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if _, ok := r.s.beneficiaries[b.ID]; ok {
		return domain.Conflict("对象已存在")
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now().UTC()
	}
	r.s.beneficiaries[b.ID] = cloneBeneficiary(b)
	return nil
}
func (r *BeneficiaryRepo) Get(ctx context.Context, id string) (*domain.Beneficiary, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	b, ok := r.s.beneficiaries[id]
	if !ok {
		return nil, domain.NotFound("对象", id)
	}
	return cloneBeneficiary(b), nil
}

type SchemeRepo struct{ s *Store }

func (r *SchemeRepo) Create(ctx context.Context, v *domain.CoverageScheme) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if _, ok := r.s.schemes[v.ID]; ok {
		return domain.Conflict("方案已存在")
	}
	if v.Version == 0 {
		v.Version = 1
	}
	r.s.schemes[v.ID] = cloneScheme(v)
	return nil
}
func (r *SchemeRepo) Get(ctx context.Context, id string) (*domain.CoverageScheme, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	v, ok := r.s.schemes[id]
	if !ok {
		return nil, domain.NotFound("方案", id)
	}
	return cloneScheme(v), nil
}
func (r *SchemeRepo) Update(ctx context.Context, v *domain.CoverageScheme) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	old, ok := r.s.schemes[v.ID]
	if !ok {
		return domain.NotFound("方案", v.ID)
	}
	if old.Version+1 != v.Version {
		return domain.Conflict("方案版本冲突")
	}
	r.s.schemes[v.ID] = cloneScheme(v)
	return nil
}

type ClaimRepo struct{ s *Store }

func (r *ClaimRepo) Create(ctx context.Context, c *domain.ExpenseClaim) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if c.IdempotencyKey != "" {
		for _, v := range r.s.claims {
			if v.IdempotencyKey == c.IdempotencyKey {
				return domain.Conflict("重复申报")
			}
		}
	}
	if c.Version == 0 {
		c.Version = 1
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	c.UpdatedAt = c.CreatedAt
	r.s.claims[c.ID] = cloneClaim(c)
	return nil
}
func (r *ClaimRepo) Get(ctx context.Context, id string) (*domain.ExpenseClaim, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	c, ok := r.s.claims[id]
	if !ok {
		return nil, domain.NotFound("申报", id)
	}
	return cloneClaim(c), nil
}
func (r *ClaimRepo) GetByIdempotency(ctx context.Context, key string) (*domain.ExpenseClaim, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, c := range r.s.claims {
		if c.IdempotencyKey == key {
			return cloneClaim(c), nil
		}
	}
	return nil, domain.NotFound("申报幂等键", key)
}
func (r *ClaimRepo) Update(ctx context.Context, c *domain.ExpenseClaim) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	old, ok := r.s.claims[c.ID]
	if !ok {
		return domain.NotFound("申报", c.ID)
	}
	if old.Version+1 != c.Version {
		return domain.Conflict("申报版本冲突")
	}
	c.UpdatedAt = time.Now().UTC()
	r.s.claims[c.ID] = cloneClaim(c)
	return nil
}
func (r *ClaimRepo) ListByYear(ctx context.Context, p string, y int) ([]*domain.ExpenseClaim, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	out := []*domain.ExpenseClaim{}
	for _, c := range r.s.claims {
		if c.ProjectID == p && c.Year == y {
			out = append(out, cloneClaim(c))
		}
	}
	return out, nil
}

type RuleRepo struct{ s *Store }

func (r *RuleRepo) Create(ctx context.Context, v *domain.RuleVersion) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.rules[v.ID] = cloneRule(v)
	return nil
}
func (r *RuleRepo) Get(ctx context.Context, id string) (*domain.RuleVersion, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	v, ok := r.s.rules[id]
	if !ok {
		return nil, domain.NotFound("规则", id)
	}
	return cloneRule(v), nil
}
func (r *RuleRepo) Latest(ctx context.Context, p string, y int) (*domain.RuleVersion, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	var best *domain.RuleVersion
	for _, v := range r.s.rules {
		if v.ProjectID == p && v.Year == y && v.Status == domain.RulePublished && (best == nil || v.VersionNo > best.VersionNo) {
			best = v
		}
	}
	if best == nil {
		return nil, domain.NotFound("已发布规则", p)
	}
	return cloneRule(best), nil
}
func (r *RuleRepo) List(ctx context.Context, p string, y int) ([]*domain.RuleVersion, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	out := []*domain.RuleVersion{}
	for _, v := range r.s.rules {
		if v.ProjectID == p && v.Year == y {
			out = append(out, cloneRule(v))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].VersionNo < out[j].VersionNo })
	return out, nil
}

type SettlementRepo struct{ s *Store }

func (r *SettlementRepo) Create(ctx context.Context, v *domain.Settlement) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if v.IdempotencyKey != "" {
		for _, x := range r.s.settlements {
			if x.IdempotencyKey == v.IdempotencyKey {
				return domain.Conflict("重复结算")
			}
		}
	}
	if v.Version == 0 {
		v.Version = 1
	}
	if v.CreatedAt.IsZero() {
		v.CreatedAt = time.Now().UTC()
	}
	r.s.settlements[v.ID] = cloneSettlement(v)
	return nil
}
func (r *SettlementRepo) Get(ctx context.Context, id string) (*domain.Settlement, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	v, ok := r.s.settlements[id]
	if !ok {
		return nil, domain.NotFound("结算", id)
	}
	return cloneSettlement(v), nil
}
func (r *SettlementRepo) GetByIdempotency(ctx context.Context, key string) (*domain.Settlement, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, v := range r.s.settlements {
		if v.IdempotencyKey == key {
			return cloneSettlement(v), nil
		}
	}
	return nil, domain.NotFound("结算幂等键", key)
}
func (r *SettlementRepo) Update(ctx context.Context, v *domain.Settlement) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	old, ok := r.s.settlements[v.ID]
	if !ok {
		return domain.NotFound("结算", v.ID)
	}
	if old.Version+1 != v.Version {
		return domain.Conflict("结算版本冲突")
	}
	r.s.settlements[v.ID] = cloneSettlement(v)
	return nil
}
func (r *SettlementRepo) Ledger(ctx context.Context, p, b string, y int) (*domain.AnnualLedger, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	v, ok := r.s.ledgers[p+":"+b+":"+timeKey(y)]
	if !ok {
		return &domain.AnnualLedger{ProjectID: p, BeneficiaryID: b, Year: y, Version: 1}, nil
	}
	return cloneLedger(v), nil
}
func (r *SettlementRepo) SaveLedger(ctx context.Context, v *domain.AnnualLedger) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	key := v.ProjectID + ":" + v.BeneficiaryID + ":" + timeKey(v.Year)
	old, ok := r.s.ledgers[key]
	if !ok {
		// 首次落库接受写入，版本从 1 起算。
		if v.Version < 1 {
			v.Version = 1
		}
		r.s.ledgers[key] = cloneLedger(v)
		return nil
	}
	// 乐观锁：提交版本必须正好是当前版本 +1，否则判定为陈旧写入予以拒绝。
	if old.Version+1 != v.Version {
		return domain.Conflict("账本版本冲突")
	}
	r.s.ledgers[key] = cloneLedger(v)
	return nil
}
func (r *SettlementRepo) ListByYear(ctx context.Context, p string, y int) ([]*domain.Settlement, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	out := []*domain.Settlement{}
	for _, v := range r.s.settlements {
		if v.ProjectID == p && v.Year == y {
			out = append(out, cloneSettlement(v))
		}
	}
	return out, nil
}

type AuditRepo struct{ s *Store }

func (r *AuditRepo) Append(ctx context.Context, v *domain.AuditEvent) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if v.At.IsZero() {
		v.At = time.Now().UTC()
	}
	c := *v
	r.s.audits = append(r.s.audits, &c)
	return nil
}
func (r *AuditRepo) List(ctx context.Context, e string, limit int) ([]*domain.AuditEvent, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	out := []*domain.AuditEvent{}
	for _, v := range r.s.audits {
		if e == "" || v.EntityID == e {
			c := *v
			out = append(out, &c)
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[len(out)-limit:]
	}
	return out, nil
}
func timeKey(y int) string { return string(rune(y)) }
func NewPorts() *domain.Ports {
	s := NewStore()
	return &domain.Ports{Projects: &ProjectRepo{s}, Beneficiaries: &BeneficiaryRepo{s}, Schemes: &SchemeRepo{s}, Claims: &ClaimRepo{s}, Rules: &RuleRepo{s}, Settlements: &SettlementRepo{s}, Audits: &AuditRepo{s}}
}
