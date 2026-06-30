package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
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
	FindOrProvisionGoogleUser(ctx context.Context, input GoogleUserInput) (AuthenticatedUserSession, error)
	ResolveActiveMembership(ctx context.Context, userID string, organizationID *string) (model.OrganizationMembership, error)
}

type CreateInput struct {
	OrganizationID string
	GoogleSub      string
	Email          string
	Name           string
	AvatarURL      string
	Role           model.UserRole
	LocationIDs    []string
	SiteIDs        []string
}

type GoogleUserInput struct {
	GoogleSub string
	Email     string
	Name      string
	AvatarURL string
}

type AuthenticatedUserSession struct {
	User               model.User
	ActiveOrganization model.Organization
	ActiveMembership   model.OrganizationMembership
}

type UserUpdateInput struct {
	UserID         string
	OrganizationID string
	Name           *string
	Email          *string
}

type UserRoleUpdateInput struct {
	UserID string
	Role   model.UserRole
}

type user struct {
	userRepo         repository.User
	orgRepo          repository.Organization
	membershipRepo   repository.OrganizationMembership
	userLocationRepo repository.UserLocation
	userSiteRepo     repository.UserSite
	data             repository.DataAccess
	logger           zerolog.Logger
}

func NewUser(
	userRepo repository.User,
	orgRepo repository.Organization,
	membershipRepo repository.OrganizationMembership,
	userLocationRepo repository.UserLocation,
	userSiteRepo repository.UserSite,
	data repository.DataAccess,
	logger zerolog.Logger,
) *user {
	return &user{
		userRepo:         userRepo,
		orgRepo:          orgRepo,
		membershipRepo:   membershipRepo,
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

	googleSub := input.GoogleSub
	if googleSub == "" {
		googleSub = input.Email
	}
	nu, err := model.InitUser(googleSub, input.Name, input.Email, input.AvatarURL)
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
		membership, err := model.InitOrganizationMembership(cu.IDNatural, input.OrganizationID, input.Role)
		if err != nil {
			return err
		}
		if _, err := u.membershipRepo.Create(ctx, conn, membership); err != nil {
			return err
		}
		if err := u.userLocationRepo.SetUserLocations(ctx, conn, cu.IDNatural, input.OrganizationID, input.LocationIDs); err != nil {
			return err
		}
		return u.userSiteRepo.SetUserSites(ctx, conn, cu.IDNatural, input.OrganizationID, input.SiteIDs)
	}); err != nil {
		return model.User{}, err
	}

	return cu, nil
}

func (u *user) FindOrProvisionGoogleUser(ctx context.Context, input GoogleUserInput) (AuthenticatedUserSession, error) {
	var session AuthenticatedUserSession
	if err := u.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()

		existing, err := u.userRepo.GetByGoogleSub(ctx, conn, input.GoogleSub)
		if err != nil {
			if !apperror.SameKind(err, apperror.KindNotFound) {
				return err
			}
			initialized, err := model.InitUser(input.GoogleSub, input.Name, input.Email, input.AvatarURL)
			if err != nil {
				return err
			}
			existing, err = u.userRepo.Create(ctx, conn, initialized)
			if err != nil {
				return err
			}
		}

		memberships, err := u.membershipRepo.ListByUser(ctx, conn, existing.IDNatural)
		if err != nil {
			return err
		}

		if len(memberships) == 0 {
			desc := "Personal workspace"
			org, err := model.InitOrganization(personalWorkspaceName(input.Name, existing.IDNatural, input.GoogleSub), &desc, model.OrganizationKindPersonal)
			if err != nil {
				return err
			}
			createdOrg, err := u.orgRepo.Create(ctx, conn, org)
			if err != nil {
				return err
			}
			membership, err := model.InitOrganizationMembership(existing.IDNatural, createdOrg.IDNatural, model.UserRoleAdmin)
			if err != nil {
				return err
			}
			createdMembership, err := u.membershipRepo.Create(ctx, conn, membership)
			if err != nil {
				return err
			}
			session = AuthenticatedUserSession{
				User:               existing,
				ActiveOrganization: createdOrg,
				ActiveMembership:   createdMembership,
			}
			return nil
		}

		activeMembership := memberships[0]
		activeOrganization, err := u.orgRepo.GetByNaturalID(ctx, conn, activeMembership.OrganizationID)
		if err != nil {
			return err
		}
		session = AuthenticatedUserSession{
			User:               existing,
			ActiveOrganization: activeOrganization,
			ActiveMembership:   activeMembership,
		}
		return nil
	}); err != nil {
		return AuthenticatedUserSession{}, err
	}

	return session, nil
}

func (u *user) ResolveActiveMembership(ctx context.Context, userID string, organizationID *string) (model.OrganizationMembership, error) {
	return u.resolveActiveMembership(ctx, u.data, userID, organizationID)
}

func (u *user) resolveActiveMembership(ctx context.Context, data repository.DataAccess, userID string, organizationID *string) (model.OrganizationMembership, error) {
	conn := data.Conn()
	if organizationID != nil && *organizationID != "" {
		return u.membershipRepo.GetByUserAndOrganization(ctx, conn, userID, *organizationID)
	}

	memberships, err := u.membershipRepo.ListByUser(ctx, conn, userID)
	if err != nil {
		return model.OrganizationMembership{}, err
	}
	if len(memberships) == 0 {
		return model.OrganizationMembership{}, apperror.NewError(apperror.NewMessage(apperror.CodeForbidden, "user has no organization membership"))
	}

	return memberships[0], nil
}

func activeOrganizationID(ctx context.Context) *string {
	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return nil
	}
	return &orgID
}

func personalWorkspaceName(displayName, userID, googleSub string) string {
	token := stableShortToken(userID, googleSub)
	return fmt.Sprintf("%s's Workspace %s", displayName, token)
}

func stableShortToken(primary, fallback string) string {
	source := primary
	if source == "" {
		source = fallback
	}
	sum := sha256.Sum256([]byte(source))
	return hex.EncodeToString(sum[:])[:8]
}

func (u *user) Update(ctx context.Context, input UserUpdateInput) (model.User, error) {
	existing, err := u.userRepo.GetByNaturalID(ctx, u.data.Conn(), input.UserID)
	if err != nil {
		return model.User{}, err
	}

	needsUserUpdate := false
	if input.Name != nil {
		if err := existing.SetName(*input.Name); err != nil {
			return model.User{}, err
		}
		needsUserUpdate = true
	}

	if input.Email != nil {
		if err := existing.SetEmail(*input.Email); err != nil {
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
		activeMembership, err := u.resolveActiveMembership(ctx, txData, existing.IDNatural, activeOrganizationID(ctx))
		if err != nil {
			return err
		}
		return u.userLocationRepo.SetUserLocations(ctx, conn, existing.IDNatural, activeMembership.OrganizationID, locationIDs)
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
		activeMembership, err := u.resolveActiveMembership(ctx, txData, existing.IDNatural, activeOrganizationID(ctx))
		if err != nil {
			return err
		}
		return u.userSiteRepo.SetUserSites(ctx, conn, existing.IDNatural, activeMembership.OrganizationID, siteIDs)
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

	if err := u.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		activeMembership, err := u.resolveActiveMembership(ctx, txData, user.IDNatural, activeOrganizationID(ctx))
		if err != nil {
			return err
		}
		if _, err := model.InitOrganizationMembership(user.IDNatural, activeMembership.OrganizationID, input.Role); err != nil {
			return err
		}
		_, err = u.membershipRepo.UpdateRole(ctx, conn, user.IDNatural, activeMembership.OrganizationID, input.Role)
		return err
	}); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *user) GetByNaturalID(ctx context.Context, idNatural string) (model.User, error) {
	user, err := u.userRepo.GetByNaturalID(ctx, u.data.Conn(), idNatural)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *user) List(ctx context.Context, filter UserListFilter, page, limit int) (model.Users, int, error) {
	pg := pagination.Normalize(page, limit)
	users, total, err := u.userRepo.List(ctx, u.data.Conn(), filter.repositoryFilter(), pg.Limit, pg.Offset)
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
