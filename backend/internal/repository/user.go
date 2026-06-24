package repository

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

type UserListFilter struct {
	LocationID *string
	SiteID     *string
	Search     *string
	SortBy     *UserSortBy
	SortOrder  *SortOrder
}

type User interface {
	Create(ctx context.Context, conn DBConn, user model.User) (model.User, error)
	Update(ctx context.Context, conn DBConn, user model.User) (model.User, error)
	UpdateRole(ctx context.Context, conn DBConn, idNatural string, role model.UserRole) (model.User, error)
	GetByNaturalID(ctx context.Context, conn DBConn, IDNatural string) (model.User, error)
	ExistsByEmail(ctx context.Context, conn DBConn, email string) (bool, error)
	ExistsByEmails(ctx context.Context, conn DBConn, emails []string) (map[string]bool, error)
	List(ctx context.Context, conn DBConn, filter UserListFilter, limit, offset int) (model.Users, int, error)
	Delete(ctx context.Context, conn DBConn, idNatural string) error
}
