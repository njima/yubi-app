package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func (c *controller) StartRobotEpisode(ctx context.Context, request openapi.StartRobotEpisodeRequestObject) (openapi.StartRobotEpisodeResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	input := usecase.StartEpisodeInput{
		EpisodeID:  request.EpisodeId,
		OccurredAt: request.Body.OccurredAt,
	}

	// If a live teleop operator is registered (Redis heartbeat), pass their
	// identity so episode.Start() uses the actual operator — not the API
	// key's static user.
	if robotID, err := requestctx.RobotID(ctx); err == nil {
		if op, err := c.robotOperatorUsecase.Get(ctx, robotID); err == nil && op != nil {
			input.ActiveUserID = &op.UserID
		}
	}

	if err := c.episodeUsecase.Start(ctx, input); err != nil {
		return nil, err
	}

	return openapi.StartRobotEpisode200Response{}, nil
}

func (c *controller) FinishRobotEpisode(ctx context.Context, request openapi.FinishRobotEpisodeRequestObject) (openapi.FinishRobotEpisodeResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	input := usecase.FinishEpisodeInput{
		EpisodeID:  request.EpisodeId,
		OccurredAt: request.Body.OccurredAt,
	}

	if err := c.episodeUsecase.Finish(ctx, input); err != nil {
		return nil, err
	}

	return openapi.FinishRobotEpisode200Response{}, nil
}

func (c *controller) CancelRobotEpisode(ctx context.Context, request openapi.CancelRobotEpisodeRequestObject) (openapi.CancelRobotEpisodeResponseObject, error) {
	input := usecase.CancelEpisodeInput{
		EpisodeID: request.EpisodeId,
	}

	if err := c.episodeUsecase.Cancel(ctx, input); err != nil {
		return nil, err
	}

	return openapi.CancelRobotEpisode200Response{}, nil
}
