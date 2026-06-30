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
	Create(ctx context.Context, conn Conn, user model.User) (model.User, error)
	Update(ctx context.Context, conn Conn, user model.User) (model.User, error)
	GetByNaturalID(ctx context.Context, conn Conn, IDNatural string) (model.User, error)
	GetByGoogleSub(ctx context.Context, conn Conn, googleSub string) (model.User, error)
	ExistsByEmail(ctx context.Context, conn Conn, email string) (bool, error)
	ExistsByEmails(ctx context.Context, conn Conn, emails []string) (map[string]bool, error)
	List(ctx context.Context, conn Conn, filter UserListFilter, limit, offset int) (model.Users, int, error)
	Delete(ctx context.Context, conn Conn, idNatural string) error
}
