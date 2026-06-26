package repository

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

type LocationListFilter struct {
	SiteID    *string
	Search    *string
	SortBy    *LocationSortBy
	SortOrder *SortOrder
}

type Location interface {
	Create(ctx context.Context, conn Conn, loc model.Location) (model.Location, error)
	GetByID(ctx context.Context, conn Conn, id string) (model.Location, error)
	List(ctx context.Context, conn Conn, filter LocationListFilter, limit, offset int) (model.Locations, int, error)
	Update(ctx context.Context, conn Conn, loc model.Location) (model.Location, error)
	Delete(ctx context.Context, conn Conn, id string) error
}
