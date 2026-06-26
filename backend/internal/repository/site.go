package repository

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

type SiteListFilter struct {
	OrganizationID *string
	Search         *string
}

type Site interface {
	Create(ctx context.Context, conn Conn, site model.Site) (model.Site, error)
	GetByID(ctx context.Context, conn Conn, id string) (model.Site, error)
	List(ctx context.Context, conn Conn, filter SiteListFilter, limit, offset int) (model.Sites, int, error)
	Update(ctx context.Context, conn Conn, site model.Site) (model.Site, error)
	Delete(ctx context.Context, conn Conn, id string) error
}
