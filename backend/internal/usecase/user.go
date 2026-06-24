package usecase

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/pagination"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type UserUsecase interface {
	Create(ctx context.Context, input CreateInput) (model.User, error)
	Update(ctx context.Context, input UserUpdateInput) (model.User, error)
	UpdateRole(ctx context.Context, input UserRoleUpdateInput) (model.User, error)
	SetLocations(ctx context.Context, userID string, locationIDs []string) (model.User, error)
	SetSites(ctx context.Context, userID string, siteIDs []string) (model.User, error)
	GetByNaturalID(ctx context.Context, idNatural string) (model.User, error)
	List(ctx context.Context, filter repository.UserListFilter, page, limit int) (model.Users, int, error)
	Delete(ctx context.Context, idNatural string) error
}

type CreateInput struct {
	OrganizationID string
	Email          string
	Name           string
	Role           model.UserRole
	LocationIDs    []string
	SiteIDs        []string
}

type UserUpdateInput struct {
	UserID         string
	OrganizationID string
	Name           string
	Email          string
}

type UserRoleUpdateInput struct {
	UserID string
	Role   model.UserRole
}

type user struct {
	userRepo         repository.User
	userLocationRepo repository.UserLocation
	userSiteRepo     repository.UserSite
	db               *bun.DB
	logger           zerolog.Logger
}

func NewUser(
	userRepo repository.User,
	userLocationRepo repository.UserLocation,
	userSiteRepo repository.UserSite,
	db *bun.DB,
	logger zerolog.Logger,
) *user {
	return &user{
		userRepo:         userRepo,
		userLocationRepo: userLocationRepo,
		userSiteRepo:     userSiteRepo,
		db:               db,
		logger:           logger,
	}
}

func (u *user) Create(ctx context.Context, input CreateInput) (model.User, error) {
	exists, err := u.userRepo.ExistsByEmail(ctx, u.db, input.Email)
	if err != nil {
		return model.User{}, err
	}
	if exists {
		return model.User{}, apperror.NewError(apperror.NewMessage(apperror.CodeConflict, "user with this email already exists"))
	}

	nu, err := model.InitUser(input.OrganizationID, input.Name, input.Email, input.Role)
	if err != nil {
		return model.User{}, err
	}

	var cu model.User
	if err := u.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var err error
		cu, err = u.userRepo.Create(ctx, tx, nu)
		if err != nil {
			return err
		}
		if err := u.userLocationRepo.SetUserLocations(ctx, tx, cu.IDNatural, cu.OrganizationID, input.LocationIDs); err != nil {
			return err
		}
		return u.userSiteRepo.SetUserSites(ctx, tx, cu.IDNatural, cu.OrganizationID, input.SiteIDs)
	}); err != nil {
		return model.User{}, err
	}

	return cu, nil
}

func (u *user) Update(ctx context.Context, input UserUpdateInput) (model.User, error) {
	existing, err := u.userRepo.GetByNaturalID(ctx, u.db, input.UserID)
	if err != nil {
		return model.User{}, err
	}

	needsUserUpdate := false
	if input.Name != "" {
		if err := existing.SetName(input.Name); err != nil {
			return model.User{}, err
		}
		needsUserUpdate = true
	}

	if input.Email != "" {
		if err := existing.SetEmail(input.Email); err != nil {
			return model.User{}, err
		}
		needsUserUpdate = true
	}

	if !needsUserUpdate {
		return existing, nil
	}

	updatedUser, err := u.userRepo.Update(ctx, u.db, existing)
	if err != nil {
		return model.User{}, err
	}

	return updatedUser, nil
}

func (u *user) SetLocations(ctx context.Context, userID string, locationIDs []string) (model.User, error) {
	existing, err := u.userRepo.GetByNaturalID(ctx, u.db, userID)
	if err != nil {
		return model.User{}, err
	}

	if err := u.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return u.userLocationRepo.SetUserLocations(ctx, tx, existing.IDNatural, existing.OrganizationID, locationIDs)
	}); err != nil {
		return model.User{}, err
	}

	return u.userRepo.GetByNaturalID(ctx, u.db, userID)
}

func (u *user) SetSites(ctx context.Context, userID string, siteIDs []string) (model.User, error) {
	existing, err := u.userRepo.GetByNaturalID(ctx, u.db, userID)
	if err != nil {
		return model.User{}, err
	}

	if err := u.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return u.userSiteRepo.SetUserSites(ctx, tx, existing.IDNatural, existing.OrganizationID, siteIDs)
	}); err != nil {
		return model.User{}, err
	}

	return u.userRepo.GetByNaturalID(ctx, u.db, userID)
}

func (u *user) UpdateRole(ctx context.Context, input UserRoleUpdateInput) (model.User, error) {
	user, err := u.userRepo.GetByNaturalID(ctx, u.db, input.UserID)
	if err != nil {
		return model.User{}, err
	}

	if err := user.SetRole(input.Role); err != nil {
		return model.User{}, err
	}

	updatedUser, err := u.userRepo.UpdateRole(ctx, u.db, user.IDNatural, user.Role)
	if err != nil {
		return model.User{}, err
	}

	return updatedUser, nil
}

func (u *user) GetByNaturalID(ctx context.Context, idNatural string) (model.User, error) {
	user, err := u.userRepo.GetByNaturalID(ctx, u.db, idNatural)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *user) List(ctx context.Context, filter repository.UserListFilter, page, limit int) (model.Users, int, error) {
	if limit <= 0 {
		limit = pagination.DefaultLimit
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	users, total, err := u.userRepo.List(ctx, u.db, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (u *user) Delete(ctx context.Context, idNatural string) error {
	if err := u.userRepo.Delete(ctx, u.db, idNatural); err != nil {
		return err
	}
	return nil
}
