package usecase

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
	"github.com/rs/zerolog"
)

type UserUsecase interface {
	Create(ctx context.Context, input CreateInput) (model.User, error)
	Update(ctx context.Context, input UserUpdateInput) (model.User, error)
	UpdateRole(ctx context.Context, input UserRoleUpdateInput) (model.User, error)
	SetLocations(ctx context.Context, userID string, locationIDs []string) (model.User, error)
	SetSites(ctx context.Context, userID string, siteIDs []string) (model.User, error)
	GetByNaturalID(ctx context.Context, idNatural string) (model.User, error)
	List(ctx context.Context, filter UserListFilter, page, limit int) (model.Users, int, error)
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
	data             repository.DataAccess
	logger           zerolog.Logger
}

func NewUser(
	userRepo repository.User,
	userLocationRepo repository.UserLocation,
	userSiteRepo repository.UserSite,
	data repository.DataAccess,
	logger zerolog.Logger,
) *user {
	return &user{
		userRepo:         userRepo,
		userLocationRepo: userLocationRepo,
		userSiteRepo:     userSiteRepo,
		data:             data,
		logger:           logger,
	}
}

func (u *user) Create(ctx context.Context, input CreateInput) (model.User, error) {
	exists, err := u.userRepo.ExistsByEmail(ctx, u.data.Conn(), input.Email)
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
	if err := u.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		var err error
		cu, err = u.userRepo.Create(ctx, conn, nu)
		if err != nil {
			return err
		}
		if err := u.userLocationRepo.SetUserLocations(ctx, conn, cu.IDNatural, cu.OrganizationID, input.LocationIDs); err != nil {
			return err
		}
		return u.userSiteRepo.SetUserSites(ctx, conn, cu.IDNatural, cu.OrganizationID, input.SiteIDs)
	}); err != nil {
		return model.User{}, err
	}

	return cu, nil
}

func (u *user) Update(ctx context.Context, input UserUpdateInput) (model.User, error) {
	existing, err := u.userRepo.GetByNaturalID(ctx, u.data.Conn(), input.UserID)
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

	updatedUser, err := u.userRepo.Update(ctx, u.data.Conn(), existing)
	if err != nil {
		return model.User{}, err
	}

	return updatedUser, nil
}

func (u *user) SetLocations(ctx context.Context, userID string, locationIDs []string) (model.User, error) {
	existing, err := u.userRepo.GetByNaturalID(ctx, u.data.Conn(), userID)
	if err != nil {
		return model.User{}, err
	}

	if err := u.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		return u.userLocationRepo.SetUserLocations(ctx, conn, existing.IDNatural, existing.OrganizationID, locationIDs)
	}); err != nil {
		return model.User{}, err
	}

	return u.userRepo.GetByNaturalID(ctx, u.data.Conn(), userID)
}

func (u *user) SetSites(ctx context.Context, userID string, siteIDs []string) (model.User, error) {
	existing, err := u.userRepo.GetByNaturalID(ctx, u.data.Conn(), userID)
	if err != nil {
		return model.User{}, err
	}

	if err := u.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		return u.userSiteRepo.SetUserSites(ctx, conn, existing.IDNatural, existing.OrganizationID, siteIDs)
	}); err != nil {
		return model.User{}, err
	}

	return u.userRepo.GetByNaturalID(ctx, u.data.Conn(), userID)
}

func (u *user) UpdateRole(ctx context.Context, input UserRoleUpdateInput) (model.User, error) {
	user, err := u.userRepo.GetByNaturalID(ctx, u.data.Conn(), input.UserID)
	if err != nil {
		return model.User{}, err
	}

	if err := user.SetRole(input.Role); err != nil {
		return model.User{}, err
	}

	updatedUser, err := u.userRepo.UpdateRole(ctx, u.data.Conn(), user.IDNatural, user.Role)
	if err != nil {
		return model.User{}, err
	}

	return updatedUser, nil
}

func (u *user) GetByNaturalID(ctx context.Context, idNatural string) (model.User, error) {
	user, err := u.userRepo.GetByNaturalID(ctx, u.data.Conn(), idNatural)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *user) List(ctx context.Context, filter UserListFilter, page, limit int) (model.Users, int, error) {
	if limit <= 0 {
		limit = pagination.DefaultLimit
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	users, total, err := u.userRepo.List(ctx, u.data.Conn(), filter.repositoryFilter(), limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (u *user) Delete(ctx context.Context, idNatural string) error {
	if err := u.userRepo.Delete(ctx, u.data.Conn(), idNatural); err != nil {
		return err
	}
	return nil
}
