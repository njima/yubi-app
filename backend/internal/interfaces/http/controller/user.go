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
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil {
		return nil, err
	}

	var locationIDs []string
	if request.Body.LocationIds != nil {
		locationIDs = *request.Body.LocationIds
	}
	var siteIDs []string
	if request.Body.SiteIds != nil {
		siteIDs = *request.Body.SiteIds
	}
	role, err := userRoleModel(request.Body.Role)
	if err != nil {
		return nil, err
	}

	user, err := c.userUsecase.Create(ctx, usecase.CreateInput{
		OrganizationID: orgID,
		Email:          string(request.Body.Email),
		Name:           request.Body.DisplayName,
		Role:           role,
		LocationIDs:    locationIDs,
		SiteIDs:        siteIDs,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateUser201JSONResponse{
		CreatedAt:      user.CreatedAt,
		DisplayName:    user.Name,
		Email:          user.Email,
		Role:           openAPIUserRolePtr(user.Role),
		OrganizationId: user.OrganizationID,
		UpdatedAt:      user.UpdatedAt,
		UserId:         user.IDNatural,
		Locations:      toLocationSummaries(user.Locations),
		Sites:          toSiteSummaries(user.Sites),
	}, nil
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

	respUsers := make([]openapi.UserResponse, 0, len(users))
	for _, u := range users {
		respUsers = append(respUsers, openapi.UserResponse{
			UserId:           u.IDNatural,
			Email:            u.Email,
			DisplayName:      u.Name,
			Role:             openAPIUserRolePtr(u.Role),
			OrganizationId:   u.OrganizationID,
			OrganizationName: u.OrganizationName,
			CreatedAt:        u.CreatedAt,
			UpdatedAt:        u.UpdatedAt,
			Locations:        toLocationSummaries(u.Locations),
			Sites:            toSiteSummaries(u.Sites),
		})
	}

	userFilter := openapi.UserFilter{
		LocationId: params.LocationId,
		SiteId:     params.SiteId,
	}
	if orgID, err := requestctx.OrganizationID(ctx); err == nil {
		userFilter.OrganizationId = &orgID
	}

	return openapi.ListUsers200JSONResponse{
		Filter: userFilter,
		Pagination: openapi.Pagination{
			Count: total,
			Limit: pg.Limit,
			Page:  pg.Page,
		},
		Users: respUsers,
	}, nil
}

func (c *controller) GetUserById(ctx context.Context, request openapi.GetUserByIdRequestObject) (openapi.GetUserByIdResponseObject, error) {
	user, err := c.userUsecase.GetByNaturalID(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	return openapi.GetUserById200JSONResponse{
		UserId:           user.IDNatural,
		Email:            user.Email,
		DisplayName:      user.Name,
		Role:             openAPIUserRolePtr(user.Role),
		OrganizationId:   user.OrganizationID,
		OrganizationName: user.OrganizationName,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		Locations:        toLocationSummaries(user.Locations),
		Sites:            toSiteSummaries(user.Sites),
	}, nil
}

func (c *controller) GetMe(ctx context.Context, request openapi.GetMeRequestObject) (openapi.GetMeResponseObject, error) {
	userID, err := requestctx.UserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := c.userUsecase.GetByNaturalID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return openapi.GetMe200JSONResponse{
		UserId:           user.IDNatural,
		Email:            user.Email,
		DisplayName:      user.Name,
		Role:             openAPIUserRolePtr(user.Role),
		OrganizationId:   user.OrganizationID,
		OrganizationName: user.OrganizationName,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		Locations:        toLocationSummaries(user.Locations),
		Sites:            toSiteSummaries(user.Sites),
	}, nil
}

func (c *controller) UpdateMe(ctx context.Context, request openapi.UpdateMeRequestObject) (openapi.UpdateMeResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	userID, err := requestctx.UserID(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := c.userUsecase.Update(ctx, usecase.UserUpdateInput{
		UserID: userID,
		Name:   request.Body.DisplayName,
	}); err != nil {
		return nil, err
	}

	user, err := c.userUsecase.GetByNaturalID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateMe200JSONResponse{
		UserId:           user.IDNatural,
		Email:            user.Email,
		DisplayName:      user.Name,
		Role:             openAPIUserRolePtr(user.Role),
		OrganizationId:   user.OrganizationID,
		OrganizationName: user.OrganizationName,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		Locations:        toLocationSummaries(user.Locations),
		Sites:            toSiteSummaries(user.Sites),
	}, nil
}

func (c *controller) UpdateUserById(ctx context.Context, request openapi.UpdateUserByIdRequestObject) (openapi.UpdateUserByIdResponseObject, error) {
	var email, displayName string

	if request.Body.Email != nil {
		email = string(*request.Body.Email)
	}
	if request.Body.DisplayName != nil {
		displayName = *request.Body.DisplayName
	}

	user, err := c.userUsecase.Update(ctx, usecase.UserUpdateInput{
		UserID: request.UserId,
		Email:  email,
		Name:   displayName,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateUserById200JSONResponse{
		UserId:         user.IDNatural,
		Email:          user.Email,
		DisplayName:    user.Name,
		Role:           openAPIUserRolePtr(user.Role),
		OrganizationId: user.OrganizationID,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		Locations:      toLocationSummaries(user.Locations),
		Sites:          toSiteSummaries(user.Sites),
	}, nil
}

func (c *controller) UpdateUserRole(ctx context.Context, request openapi.UpdateUserRoleRequestObject) (openapi.UpdateUserRoleResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}
	role, err := userRoleModel(request.Body.Role)
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

	return openapi.UpdateUserRole200JSONResponse{
		UserId:         user.IDNatural,
		Email:          user.Email,
		DisplayName:    user.Name,
		Role:           openAPIUserRolePtr(user.Role),
		OrganizationId: user.OrganizationID,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		Locations:      toLocationSummaries(user.Locations),
		Sites:          toSiteSummaries(user.Sites),
	}, nil
}

func toLocationSummaries(locs []model.LocationSummary) []openapi.LocationSummary {
	result := make([]openapi.LocationSummary, 0, len(locs))
	for _, l := range locs {
		result = append(result, openapi.LocationSummary{
			LocationId: l.LocationID,
			Name:       l.Name,
		})
	}
	return result
}

func (c *controller) DeleteUserById(ctx context.Context, request openapi.DeleteUserByIdRequestObject) (openapi.DeleteUserByIdResponseObject, error) {
	if err := c.userUsecase.Delete(ctx, request.UserId); err != nil {
		return nil, err
	}

	return openapi.DeleteUserById204Response{}, nil
}

func (c *controller) UpdateUserLocations(ctx context.Context, request openapi.UpdateUserLocationsRequestObject) (openapi.UpdateUserLocationsResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	user, err := c.userUsecase.SetLocations(ctx, request.UserId, request.Body.LocationIds)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateUserLocations200JSONResponse{
		UserId:         user.IDNatural,
		Email:          user.Email,
		DisplayName:    user.Name,
		Role:           openAPIUserRolePtr(user.Role),
		OrganizationId: user.OrganizationID,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		Locations:      toLocationSummaries(user.Locations),
		Sites:          toSiteSummaries(user.Sites),
	}, nil
}

func toSiteSummaries(sites []model.SiteSummary) []openapi.SiteSummary {
	result := make([]openapi.SiteSummary, 0, len(sites))
	for _, s := range sites {
		result = append(result, openapi.SiteSummary{
			SiteId: s.SiteID,
			Name:   s.Name,
		})
	}
	return result
}

func (c *controller) UpdateUserSites(ctx context.Context, request openapi.UpdateUserSitesRequestObject) (openapi.UpdateUserSitesResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	user, err := c.userUsecase.SetSites(ctx, request.UserId, request.Body.SiteIds)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateUserSites200JSONResponse{
		UserId:         user.IDNatural,
		Email:          user.Email,
		DisplayName:    user.Name,
		Role:           openAPIUserRolePtr(user.Role),
		OrganizationId: user.OrganizationID,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		Locations:      toLocationSummaries(user.Locations),
		Sites:          toSiteSummaries(user.Sites),
	}, nil
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
