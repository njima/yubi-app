package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func (c *controller) ListSubTasks(ctx context.Context, request openapi.ListSubTasksRequestObject) (openapi.ListSubTasksResponseObject, error) {
	page := 1
	limit := 50
	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}

	sts, _, err := c.subtaskUsecase.List(ctx, request.Params.TaskId, request.Params.TaskVersionId, page, limit)
	if err != nil {
		return nil, err
	}

	return openapi.ListSubTasks200JSONResponse(subTaskResponses(sts)), nil
}

func (c *controller) CreateSubTask(ctx context.Context, request openapi.CreateSubTaskRequestObject) (openapi.CreateSubTaskResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	st, err := c.subtaskUsecase.Create(ctx, usecase.SubTaskCreateInput{
		OrganizationID:        body.OrganizationId,
		TaskID:                body.TaskId,
		TaskVersionID:         body.TaskVersionId,
		Name:                  body.Name,
		Description:           body.Description,
		TargetDurationSeconds: body.TargetDurationSeconds,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateSubTask201JSONResponse(subTaskResponse(st)), nil
}

func (c *controller) DeleteSubTaskById(ctx context.Context, request openapi.DeleteSubTaskByIdRequestObject) (openapi.DeleteSubTaskByIdResponseObject, error) {
	if err := c.subtaskUsecase.Delete(ctx, request.SubtaskId); err != nil {
		return nil, err
	}
	return openapi.DeleteSubTaskById204Response{}, nil
}

func (c *controller) GetSubTaskById(ctx context.Context, request openapi.GetSubTaskByIdRequestObject) (openapi.GetSubTaskByIdResponseObject, error) {
	st, err := c.subtaskUsecase.GetByID(ctx, request.SubtaskId)
	if err != nil {
		return nil, err
	}
	return openapi.GetSubTaskById200JSONResponse(subTaskResponse(st)), nil
}

func (c *controller) UpdateSubTaskById(ctx context.Context, request openapi.UpdateSubTaskByIdRequestObject) (openapi.UpdateSubTaskByIdResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	input := usecase.SubTaskUpdateInput{ID: request.SubtaskId}
	if body.Name != nil {
		input.Name = body.Name
	}
	if body.Description != nil {
		input.Description = body.Description
	}
	if body.TargetDurationSeconds != nil {
		input.TargetDurationSeconds = body.TargetDurationSeconds
	}

	st, err := c.subtaskUsecase.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateSubTaskById200JSONResponse(subTaskResponse(st)), nil
}

func (c *controller) ReorderSubTasks(ctx context.Context, request openapi.ReorderSubTasksRequestObject) (openapi.ReorderSubTasksResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	sts, err := c.subtaskUsecase.Reorder(ctx, usecase.SubTaskReorderInput{
		TaskVersionID: body.TaskVersionId,
		SubTaskIDs:    body.SubtaskIds,
	})
	if err != nil {
		return nil, err
	}

	return openapi.ReorderSubTasks200JSONResponse(subTaskResponses(sts)), nil
}

func (c *controller) CompleteSubTask(ctx context.Context, request openapi.CompleteSubTaskRequestObject) (openapi.CompleteSubTaskResponseObject, error) {
	return openapi.CompleteSubTask200JSONResponse{}, nil
}
