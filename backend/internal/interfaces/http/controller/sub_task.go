package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
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

	resp := make([]openapi.SubTask, 0, len(sts))
	for _, s := range sts {
		resp = append(resp, openapi.SubTask{
			Id:                    s.IDNatural,
			Name:                  s.Name,
			Description:           s.Description,
			TargetDurationSeconds: s.TargetDurationSeconds,
		})
	}

	return openapi.ListSubTasks200JSONResponse(resp), nil
}

func (c *controller) CreateSubTask(ctx context.Context, request openapi.CreateSubTaskRequestObject) (openapi.CreateSubTaskResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	st, err := c.subtaskUsecase.Create(ctx, usecase.SubTaskCreateInput{
		OrganizationID:        request.Body.OrganizationId,
		TaskID:                request.Body.TaskId,
		TaskVersionID:         request.Body.TaskVersionId,
		Name:                  request.Body.Name,
		Description:           request.Body.Description,
		TargetDurationSeconds: request.Body.TargetDurationSeconds,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateSubTask201JSONResponse{
		Id:                    st.IDNatural,
		Name:                  st.Name,
		Description:           st.Description,
		TargetDurationSeconds: st.TargetDurationSeconds,
	}, nil
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
	return openapi.GetSubTaskById200JSONResponse{
		Id:                    st.IDNatural,
		Name:                  st.Name,
		Description:           st.Description,
		TargetDurationSeconds: st.TargetDurationSeconds,
	}, nil
}

func (c *controller) UpdateSubTaskById(ctx context.Context, request openapi.UpdateSubTaskByIdRequestObject) (openapi.UpdateSubTaskByIdResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	input := usecase.SubTaskUpdateInput{ID: request.SubtaskId}
	if request.Body.Name != nil {
		input.Name = request.Body.Name
	}
	if request.Body.Description != nil {
		input.Description = request.Body.Description
	}
	if request.Body.TargetDurationSeconds != nil {
		input.TargetDurationSeconds = request.Body.TargetDurationSeconds
	}

	st, err := c.subtaskUsecase.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateSubTaskById200JSONResponse{
		Id:                    st.IDNatural,
		Name:                  st.Name,
		Description:           st.Description,
		TargetDurationSeconds: st.TargetDurationSeconds,
	}, nil
}

func (c *controller) ReorderSubTasks(ctx context.Context, request openapi.ReorderSubTasksRequestObject) (openapi.ReorderSubTasksResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	sts, err := c.subtaskUsecase.Reorder(ctx, usecase.SubTaskReorderInput{
		TaskVersionID: request.Body.TaskVersionId,
		SubTaskIDs:    request.Body.SubtaskIds,
	})
	if err != nil {
		return nil, err
	}

	resp := make(openapi.ReorderSubTasks200JSONResponse, 0, len(sts))
	for _, s := range sts {
		resp = append(resp, openapi.SubTask{
			Id:                    s.IDNatural,
			Name:                  s.Name,
			Description:           s.Description,
			TargetDurationSeconds: s.TargetDurationSeconds,
		})
	}

	return resp, nil
}

func (c *controller) CompleteSubTask(ctx context.Context, request openapi.CompleteSubTaskRequestObject) (openapi.CompleteSubTaskResponseObject, error) {
	return openapi.CompleteSubTask200JSONResponse{}, nil
}
