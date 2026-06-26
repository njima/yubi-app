package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

func (c *controller) ListOrganizations(ctx context.Context, request openapi.ListOrganizationsRequestObject) (openapi.ListOrganizationsResponseObject, error) {
	params := request.Params
	pg := pagination.Parse(params.Page, params.Limit)

	orgs, total, err := c.organizationUsecase.List(ctx, pg.Page, pg.Limit)
	if err != nil {
		return nil, err
	}

	respOrgs := make([]openapi.OrganizationResponse, 0, len(orgs))
	for _, o := range orgs {
		respOrgs = append(respOrgs, openapi.OrganizationResponse{
			OrganizationId: o.IDNatural,
			DisplayName:    o.Name,
			Description:    o.Description,
			CreatedAt:      &o.CreatedAt,
			UpdatedAt:      o.UpdatedAt,
		})
	}

	return openapi.ListOrganizations200JSONResponse{
		Pagination:    openapi.Pagination{Count: total, Limit: pg.Limit, Page: pg.Page},
		Organizations: respOrgs,
	}, nil
}

func (c *controller) CreateOrganization(ctx context.Context, request openapi.CreateOrganizationRequestObject) (openapi.CreateOrganizationResponseObject, error) {
	body := request.Body

	var desc *string
	if body.Description != nil {
		desc = body.Description
	}

	org, err := c.organizationUsecase.Create(ctx, usecase.OrganizationCreateInput{
		DisplayName: body.DisplayName,
		Description: desc,
	})
	if err != nil {
		if err == usecase.ErrOrganizationInvalidInput {
			return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "invalid input"))
		}
		return nil, err
	}

	return openapi.CreateOrganization201JSONResponse{
		OrganizationId: org.IDNatural,
		DisplayName:    org.Name,
		Description:    org.Description,
		CreatedAt:      &org.CreatedAt,
		UpdatedAt:      org.UpdatedAt,
	}, nil
}

func (c *controller) DeleteOrganizationById(ctx context.Context, request openapi.DeleteOrganizationByIdRequestObject) (openapi.DeleteOrganizationByIdResponseObject, error) {
	if err := c.organizationUsecase.Delete(ctx, request.OrganizationId); err != nil {
		return nil, err
	}
	return openapi.DeleteOrganizationById204Response{}, nil
}

func (c *controller) GetOrganizationById(ctx context.Context, request openapi.GetOrganizationByIdRequestObject) (openapi.GetOrganizationByIdResponseObject, error) {
	org, err := c.organizationUsecase.GetByNaturalID(ctx, request.OrganizationId)
	if err != nil {
		return nil, err
	}

	return openapi.GetOrganizationById200JSONResponse{
		OrganizationId: org.IDNatural,
		DisplayName:    org.Name,
		Description:    org.Description,
		CreatedAt:      &org.CreatedAt,
		UpdatedAt:      org.UpdatedAt,
	}, nil
}

func (c *controller) UpdateOrganizationById(ctx context.Context, request openapi.UpdateOrganizationByIdRequestObject) (openapi.UpdateOrganizationByIdResponseObject, error) {
	var displayName string
	var description *string
	if request.Body.DisplayName != nil {
		displayName = *request.Body.DisplayName
	}
	if request.Body.Description != nil {
		description = request.Body.Description
	}

	org, err := c.organizationUsecase.Update(ctx, usecase.OrganizationUpdateInput{
		OrganizationID: request.OrganizationId,
		DisplayName:    displayName,
		Description:    description,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateOrganizationById200JSONResponse{
		OrganizationId: org.IDNatural,
		DisplayName:    org.Name,
		Description:    org.Description,
		CreatedAt:      &org.CreatedAt,
		UpdatedAt:      org.UpdatedAt,
	}, nil
}
