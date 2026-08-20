package application

import (
	"context"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/service/id"
)

type ProjectService struct{ p *domain.Ports }

func NewProjectService(p *domain.Ports) *ProjectService { return &ProjectService{p} }
func (s *ProjectService) Create(ctx context.Context, v domain.GrantProject) (*domain.GrantProject, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	v.ID = id.New("prj")
	v.Status = domain.ProjectDraft
	return &v, s.p.Projects.Create(ctx, &v)
}
func (s *ProjectService) Get(ctx context.Context, idv string) (*domain.GrantProject, error) {
	return s.p.Projects.Get(ctx, idv)
}
func (s *ProjectService) Activate(ctx context.Context, idv string) (*domain.GrantProject, error) {
	v, e := s.p.Projects.Get(ctx, idv)
	if e != nil {
		return nil, e
	}
	if v.Status == domain.ProjectClosed {
		return nil, domain.StateError("已关闭项目不能激活")
	}
	v.Status = domain.ProjectActive
	v.Version++
	return v, s.p.Projects.Update(ctx, v)
}
