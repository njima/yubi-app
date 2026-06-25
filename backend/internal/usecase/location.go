package usecase

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/pagination"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type LocationUsecase interface {
	Create(ctx context.Context, input LocationCreateInput) (model.Location, error)
	GetByID(ctx context.Context, id string) (model.Location, error)
	List(ctx context.Context, filter repository.LocationListFilter, page, limit int) (model.Locations, int, error)
	Update(ctx context.Context, input LocationUpdateInput) (model.Location, error)
	Delete(ctx context.Context, id string) error
}

type LocationCreateInput struct {
	OrganizationID string
	SiteID         string
	Name           string
}

type LocationUpdateInput struct {
	ID   string
	Name string
}

type location struct {
	locRepo repository.Location
	db      repository.DBConn
}

func NewLocation(locRepo repository.Location, db repository.DBConn) *location {
	return &location{locRepo: locRepo, db: db}
}

func (l *location) Create(ctx context.Context, input LocationCreateInput) (model.Location, error) {
	lo, err := model.InitLocation(input.OrganizationID, input.SiteID, input.Name)
	if err != nil {
		return model.Location{}, err
	}

	ulo, err := l.locRepo.Create(ctx, l.db, lo)
	if err != nil {
		return model.Location{}, err
	}

	return ulo, nil
}

func (l *location) GetByID(ctx context.Context, id string) (model.Location, error) {
	return l.locRepo.GetByID(ctx, l.db, id)
}

func (l *location) List(ctx context.Context, filter repository.LocationListFilter, page, limit int) (model.Locations, int, error) {
	if limit <= 0 {
		limit = pagination.DefaultLimit
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return l.locRepo.List(ctx, l.db, filter, limit, offset)
}

func (l *location) Update(ctx context.Context, input LocationUpdateInput) (model.Location, error) {
	lo, err := l.locRepo.GetByID(ctx, l.db, input.ID)
	if err != nil {
		return model.Location{}, err
	}

	if err := lo.SetName(input.Name); err != nil {
		return model.Location{}, err
	}

	ulo, err := l.locRepo.Update(ctx, l.db, lo)
	if err != nil {
		return model.Location{}, err
	}

	return ulo, nil
}

func (l *location) Delete(ctx context.Context, id string) error {
	return l.locRepo.Delete(ctx, l.db, id)
}
