package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

func (c *controller) CreateUser(ctx context.Context, request openapi.CreateUserRequestObject) (openapi.CreateUserResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil {
		return nil, err
	}

	var locationIDs []string
	if body.LocationIds != nil {
		locationIDs = *body.LocationIds
	}
	var siteIDs []string
	if body.SiteIds != nil {
		siteIDs = *body.SiteIds
	}
	role, err := userRoleModel(body.Role)
	if err != nil {
		return nil, err
	}

	user, err := c.userUsecase.Create(ctx, usecase.CreateInput{
		OrganizationID: orgID,
		Email:          string(body.Email),
		Name:           body.DisplayName,
		Role:           role,
		LocationIDs:    locationIDs,
		SiteIDs:        siteIDs,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.userResponseInActiveWorkspace(ctx, user)
	if err != nil {
		return nil, err
	}

	return openapi.CreateUser201JSONResponse(resp), nil
}

func (c *controller) ListUsers(ctx context.Context, request openapi.ListUsersRequestObject) (openapi.ListUsersResponseObject, error) {
	params := request.Params
	pg := pagination.Parse(params.Page, params.Limit)

	filter := usecase.UserListFilter{
		Search:    params.Search,
		SortBy:    userSortBy(params.SortBy),
		SortOrder: sortOrder(params.SortOrder),
	}
	if params.LocationId != nil && *params.LocationId != "" {
		filter.LocationID = params.LocationId
	}
	if params.SiteId != nil && *params.SiteId != "" {
		filter.SiteID = params.SiteId
	}

	users, total, err := c.userUsecase.List(ctx, filter, pg.Page, pg.Limit)
	if err != nil {
		return nil, err
	}

	userFilter := openapi.UserFilter{
		LocationId: params.LocationId,
		SiteId:     params.SiteId,
	}
	if orgID, err := requestctx.OrganizationID(ctx); err == nil {
		userFilter.OrganizationId = &orgID
	}

	userResponses, err := c.userResponsesInActiveWorkspace(ctx, users)
	if err != nil {
		return nil, err
	}

	return openapi.ListUsers200JSONResponse{
		Filter: userFilter,
		Pagination: openapi.Pagination{
			Count: total,
			Limit: pg.Limit,
			Page:  pg.Page,
		},
		Users: userResponses,
	}, nil
}

func (c *controller) GetUserById(ctx context.Context, request openapi.GetUserByIdRequestObject) (openapi.GetUserByIdResponseObject, error) {
	user, err := c.userUsecase.GetByNaturalID(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	resp, err := c.userResponseInActiveWorkspace(ctx, user)
	if err != nil {
		return nil, err
	}

	return openapi.GetUserById200JSONResponse(resp), nil
}

func (c *controller) GetMe(ctx context.Context, request openapi.GetMeRequestObject) (openapi.GetMeResponseObject, error) {
	userID, err := requestctx.UserID(ctx)
	if err != nil {
		return nil, err
	}

	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil {
		return nil, err
	}

	session, err := c.userUsecase.GetAuthenticatedSession(ctx, userID, &orgID)
	if err != nil {
		return nil, err
	}

	return openapi.GetMe200JSONResponse(meResponse(session)), nil
}

func (c *controller) UpdateMe(ctx context.Context, request openapi.UpdateMeRequestObject) (openapi.UpdateMeResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	userID, err := requestctx.UserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := c.userUsecase.Update(ctx, usecase.UserUpdateInput{
		UserID: userID,
		Name:   &body.DisplayName,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.userResponseInActiveWorkspace(ctx, user)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateMe200JSONResponse(resp), nil
}

func (c *controller) UpdateUserById(ctx context.Context, request openapi.UpdateUserByIdRequestObject) (openapi.UpdateUserByIdResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	var email *string
	var displayName *string

	if body.Email != nil {
		emailValue := string(*body.Email)
		email = &emailValue
	}
	if body.DisplayName != nil {
		displayName = body.DisplayName
	}

	user, err := c.userUsecase.Update(ctx, usecase.UserUpdateInput{
		UserID: request.UserId,
		Email:  email,
		Name:   displayName,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.userResponseInActiveWorkspace(ctx, user)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateUserById200JSONResponse(resp), nil
}

func (c *controller) UpdateUserRole(ctx context.Context, request openapi.UpdateUserRoleRequestObject) (openapi.UpdateUserRoleResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}
	role, err := userRoleModel(body.Role)
	if err != nil {
		return nil, err
	}

	user, err := c.userUsecase.UpdateRole(ctx, usecase.UserRoleUpdateInput{
		UserID: request.UserId,
		Role:   role,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.userResponseInActiveWorkspace(ctx, user)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateUserRole200JSONResponse(resp), nil
}

func (c *controller) DeleteUserById(ctx context.Context, request openapi.DeleteUserByIdRequestObject) (openapi.DeleteUserByIdResponseObject, error) {
	if err := c.userUsecase.Delete(ctx, request.UserId); err != nil {
		return nil, err
	}

	return openapi.DeleteUserById204Response{}, nil
}

func (c *controller) UpdateUserLocations(ctx context.Context, request openapi.UpdateUserLocationsRequestObject) (openapi.UpdateUserLocationsResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	user, err := c.userUsecase.SetLocations(ctx, request.UserId, body.LocationIds)
	if err != nil {
		return nil, err
	}

	resp, err := c.userResponseInActiveWorkspace(ctx, user)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateUserLocations200JSONResponse(resp), nil
}

func (c *controller) UpdateUserSites(ctx context.Context, request openapi.UpdateUserSitesRequestObject) (openapi.UpdateUserSitesResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	user, err := c.userUsecase.SetSites(ctx, request.UserId, body.SiteIds)
	if err != nil {
		return nil, err
	}

	resp, err := c.userResponseInActiveWorkspace(ctx, user)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateUserSites200JSONResponse(resp), nil
}

func (c *controller) userResponseInActiveWorkspace(ctx context.Context, user model.User) (openapi.UserResponse, error) {
	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil {
		return openapi.UserResponse{}, err
	}
	org, err := c.organizationUsecase.GetByNaturalID(ctx, orgID)
	if err != nil {
		return openapi.UserResponse{}, err
	}
	membership, err := c.userUsecase.ResolveActiveMembership(ctx, user.IDNatural, &orgID)
	if err != nil {
		return openapi.UserResponse{}, err
	}
	return userResponseWithWorkspace(user, org, membership), nil
}

func (c *controller) userResponsesInActiveWorkspace(ctx context.Context, users model.Users) ([]openapi.UserResponse, error) {
	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil {
		return nil, err
	}
	org, err := c.organizationUsecase.GetByNaturalID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	responses := make([]openapi.UserResponse, 0, len(users))
	for _, user := range users {
		membership, err := c.userUsecase.ResolveActiveMembership(ctx, user.IDNatural, &orgID)
		if err != nil {
			return nil, err
		}
		responses = append(responses, userResponseWithWorkspace(*user, org, membership))
	}
	return responses, nil
}

// Permissions
func (c *controller) RevokeUserPermission(ctx context.Context, request openapi.RevokeUserPermissionRequestObject) (openapi.RevokeUserPermissionResponseObject, error) {
	return nil, apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "not implemented"))
}

func (c *controller) ListUserPermissions(ctx context.Context, request openapi.ListUserPermissionsRequestObject) (openapi.ListUserPermissionsResponseObject, error) {
	return openapi.ListUserPermissions200JSONResponse{}, nil
}

func (c *controller) GrantUserPermission(ctx context.Context, request openapi.GrantUserPermissionRequestObject) (openapi.GrantUserPermissionResponseObject, error) {
	return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "not implemented"))
}
