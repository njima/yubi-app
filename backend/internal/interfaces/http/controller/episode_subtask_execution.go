package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func (c *controller) CreateRobotExecution(ctx context.Context, request openapi.CreateRobotExecutionRequestObject) (openapi.CreateRobotExecutionResponseObject, error) {
	input := usecase.CreateExecutionInput{
		EpisodeID: request.EpisodeId,
		SubTaskID: request.SubtaskId,
	}

	executionID, err := c.episodeExecutionUsecase.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return openapi.CreateRobotExecution201JSONResponse{
		ExecutionId: executionID,
	}, nil
}

func (c *controller) StartRobotExecution(ctx context.Context, request openapi.StartRobotExecutionRequestObject) (openapi.StartRobotExecutionResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	input := usecase.ExecutionActionInput{
		EpisodeID:   request.EpisodeId,
		SubTaskID:   request.SubtaskId,
		ExecutionID: request.ExecutionId,
		OccurredAt:  body.OccurredAt,
	}

	if err := c.episodeExecutionUsecase.Start(ctx, input); err != nil {
		return nil, err
	}

	return openapi.StartRobotExecution200Response{}, nil
}

func (c *controller) FinishRobotExecution(ctx context.Context, request openapi.FinishRobotExecutionRequestObject) (openapi.FinishRobotExecutionResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	input := usecase.ExecutionActionInput{
		EpisodeID:   request.EpisodeId,
		SubTaskID:   request.SubtaskId,
		ExecutionID: request.ExecutionId,
		OccurredAt:  body.OccurredAt,
	}

	if err := c.episodeExecutionUsecase.Finish(ctx, input); err != nil {
		return nil, err
	}

	return openapi.FinishRobotExecution200Response{}, nil
}

func (c *controller) CancelRobotExecution(ctx context.Context, request openapi.CancelRobotExecutionRequestObject) (openapi.CancelRobotExecutionResponseObject, error) {
	input := usecase.CancelExecutionInput{
		EpisodeID:   request.EpisodeId,
		SubTaskID:   request.SubtaskId,
		ExecutionID: request.ExecutionId,
	}

	if err := c.episodeExecutionUsecase.Cancel(ctx, input); err != nil {
		return nil, err
	}

	return openapi.CancelRobotExecution200Response{}, nil
}
