package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

func (c *controller) ListSites(ctx context.Context, request openapi.ListSitesRequestObject) (openapi.ListSitesResponseObject, error) {
	pg := pagination.Parse(request.Params.Page, request.Params.Limit)

	filter := usecase.SiteListFilter{
		OrganizationID: request.Params.OrganizationId,
		Search:         request.Params.Search,
	}

	sites, total, err := c.siteUsecase.List(ctx, filter, pg.Page, pg.Limit)
	if err != nil {
		return nil, err
	}

	result := make([]openapi.Site, 0, len(sites))
	for _, s := range sites {
		site := openapi.Site{
			Id:             s.IDNatural,
			Name:           s.Name,
			OrganizationId: s.OrganizationID,
			CreatedAt:      &s.CreatedAt,
		}
		if s.UpdatedAt != nil {
			site.UpdatedAt = s.UpdatedAt
		}
		result = append(result, site)
	}

	return openapi.ListSites200JSONResponse{
		Sites: result,
		Pagination: openapi.Pagination{
			Count: total,
			Page:  pg.Page,
			Limit: pg.Limit,
		},
	}, nil
}

func (c *controller) CreateSite(ctx context.Context, request openapi.CreateSiteRequestObject) (openapi.CreateSiteResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	site, err := c.siteUsecase.Create(ctx, usecase.SiteCreateInput{
		OrganizationID: request.Body.OrganizationId,
		Name:           request.Body.Name,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateSite201JSONResponse{
		Id:             site.IDNatural,
		Name:           site.Name,
		OrganizationId: site.OrganizationID,
		CreatedAt:      &site.CreatedAt,
		UpdatedAt:      site.UpdatedAt,
	}, nil
}

func (c *controller) GetSiteById(ctx context.Context, request openapi.GetSiteByIdRequestObject) (openapi.GetSiteByIdResponseObject, error) {
	site, err := c.siteUsecase.GetByID(ctx, request.SiteId)
	if err != nil {
		return nil, err
	}

	return openapi.GetSiteById200JSONResponse{
		Id:             site.IDNatural,
		Name:           site.Name,
		OrganizationId: site.OrganizationID,
		CreatedAt:      &site.CreatedAt,
		UpdatedAt:      site.UpdatedAt,
	}, nil
}

func (c *controller) UpdateSiteById(ctx context.Context, request openapi.UpdateSiteByIdRequestObject) (openapi.UpdateSiteByIdResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	site, err := c.siteUsecase.Update(ctx, usecase.SiteUpdateInput{
		ID:   request.SiteId,
		Name: request.Body.Name,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateSiteById200JSONResponse{
		Id:             site.IDNatural,
		Name:           site.Name,
		OrganizationId: site.OrganizationID,
		CreatedAt:      &site.CreatedAt,
		UpdatedAt:      site.UpdatedAt,
	}, nil
}

func (c *controller) DeleteSiteById(ctx context.Context, request openapi.DeleteSiteByIdRequestObject) (openapi.DeleteSiteByIdResponseObject, error) {
	if err := c.siteUsecase.Delete(ctx, request.SiteId); err != nil {
		return nil, err
	}

	return openapi.DeleteSiteById204Response{}, nil
}
