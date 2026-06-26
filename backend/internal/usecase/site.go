package usecase

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

type SiteUsecase interface {
	Create(ctx context.Context, input SiteCreateInput) (model.Site, error)
	GetByID(ctx context.Context, id string) (model.Site, error)
	List(ctx context.Context, filter SiteListFilter, page, limit int) (model.Sites, int, error)
	Update(ctx context.Context, input SiteUpdateInput) (model.Site, error)
	Delete(ctx context.Context, id string) error
}

type SiteCreateInput struct {
	OrganizationID string
	Name           string
}

type SiteUpdateInput struct {
	ID   string
	Name string
}

type siteUsecase struct {
	siteRepo repository.Site
	data     repository.DataAccess
}

func NewSite(siteRepo repository.Site, data repository.DataAccess) *siteUsecase {
	return &siteUsecase{
		siteRepo: siteRepo,
		data:     data,
	}
}

func (s *siteUsecase) Create(ctx context.Context, input SiteCreateInput) (model.Site, error) {
	si, err := model.InitSite(input.OrganizationID, input.Name)
	if err != nil {
		return model.Site{}, err
	}

	created, err := s.siteRepo.Create(ctx, s.data.Conn(), si)
	if err != nil {
		return model.Site{}, err
	}

	return created, nil
}

func (s *siteUsecase) GetByID(ctx context.Context, id string) (model.Site, error) {
	return s.siteRepo.GetByID(ctx, s.data.Conn(), id)
}

func (s *siteUsecase) List(ctx context.Context, filter SiteListFilter, page, limit int) (model.Sites, int, error) {
	pg := pagination.Normalize(page, limit)
	return s.siteRepo.List(ctx, s.data.Conn(), filter.repositoryFilter(), pg.Limit, pg.Offset)
}

func (s *siteUsecase) Update(ctx context.Context, input SiteUpdateInput) (model.Site, error) {
	si, err := s.siteRepo.GetByID(ctx, s.data.Conn(), input.ID)
	if err != nil {
		return model.Site{}, err
	}

	if err := si.SetName(input.Name); err != nil {
		return model.Site{}, err
	}

	updated, err := s.siteRepo.Update(ctx, s.data.Conn(), si)
	if err != nil {
		return model.Site{}, err
	}

	return updated, nil
}

func (s *siteUsecase) Delete(ctx context.Context, id string) error {
	return s.siteRepo.Delete(ctx, s.data.Conn(), id)
}
