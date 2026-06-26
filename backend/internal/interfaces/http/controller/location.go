package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

func (c *controller) ListLocations(ctx context.Context, request openapi.ListLocationsRequestObject) (openapi.ListLocationsResponseObject, error) {
	pg := pagination.Parse(request.Params.Page, request.Params.Limit)

	filter := usecase.LocationListFilter{
		SiteID:    request.Params.SiteId,
		Search:    request.Params.Search,
		SortBy:    locationSortBy(request.Params.SortBy),
		SortOrder: sortOrder(request.Params.SortOrder),
	}

	locs, total, err := c.locationUsecase.List(ctx, filter, pg.Page, pg.Limit)
	if err != nil {
		return nil, err
	}

	locations := make([]openapi.Location, 0, len(locs))
	for _, l := range locs {
		locations = append(locations, locationResponse(*l))
	}

	return openapi.ListLocations200JSONResponse{
		Locations: locations,
		Pagination: openapi.Pagination{
			Count: total,
			Page:  pg.Page,
			Limit: pg.Limit,
		},
	}, nil
}

func (c *controller) CreateLocation(ctx context.Context, request openapi.CreateLocationRequestObject) (openapi.CreateLocationResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	loc, err := c.locationUsecase.Create(ctx, usecase.LocationCreateInput{
		OrganizationID: request.Body.OrganizationId,
		SiteID:         request.Body.SiteId,
		Name:           request.Body.Name,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateLocation201JSONResponse(locationResponse(loc)), nil
}

func (c *controller) DeleteLocationById(ctx context.Context, request openapi.DeleteLocationByIdRequestObject) (openapi.DeleteLocationByIdResponseObject, error) {
	if err := c.locationUsecase.Delete(ctx, request.LocationId); err != nil {
		return nil, err
	}

	return openapi.DeleteLocationById204Response{}, nil
}

func (c *controller) GetLocationById(ctx context.Context, request openapi.GetLocationByIdRequestObject) (openapi.GetLocationByIdResponseObject, error) {
	loc, err := c.locationUsecase.GetByID(ctx, request.LocationId)
	if err != nil {
		return nil, err
	}

	return openapi.GetLocationById200JSONResponse(locationResponse(loc)), nil
}

func (c *controller) UpdateLocationById(ctx context.Context, request openapi.UpdateLocationByIdRequestObject) (openapi.UpdateLocationByIdResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	input := usecase.LocationUpdateInput{
		ID:   request.LocationId,
		Name: request.Body.Name,
	}

	loc, err := c.locationUsecase.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateLocationById200JSONResponse(locationResponse(loc)), nil
}
