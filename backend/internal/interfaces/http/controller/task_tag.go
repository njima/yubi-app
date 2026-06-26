package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func (c *controller) ListTaskCategoryTypes(ctx context.Context, request openapi.ListTaskCategoryTypesRequestObject) (openapi.ListTaskCategoryTypesResponseObject, error) {
	types, err := c.taskTagUsecase.ListCategoryTypes(ctx)
	if err != nil {
		return nil, err
	}
	resp := make(openapi.ListTaskCategoryTypes200JSONResponse, 0, len(types))
	for _, t := range types {
		resp = append(resp, openapi.TaskCategoryType{
			Id:   t.ID,
			Slug: t.Slug,
			Name: t.Name,
		})
	}
	return resp, nil
}

func (c *controller) ListTaskTags(ctx context.Context, request openapi.ListTaskTagsRequestObject) (openapi.ListTaskTagsResponseObject, error) {
	tags, err := c.taskTagUsecase.ListTags(ctx, request.Params.CategoryTypeId)
	if err != nil {
		return nil, err
	}
	resp := make(openapi.ListTaskTags200JSONResponse, 0, len(tags))
	for _, t := range tags {
		resp = append(resp, toOpenAPITag(*t))
	}
	return resp, nil
}

func (c *controller) CreateTaskTag(ctx context.Context, request openapi.CreateTaskTagRequestObject) (openapi.CreateTaskTagResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	tag, err := c.taskTagUsecase.CreateTag(ctx, usecase.TaskTagCreateInput{
		Name:           request.Body.Name,
		CategoryTypeID: request.Body.CategoryTypeId,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateTaskTag201JSONResponse(toOpenAPITag(tag)), nil
}
