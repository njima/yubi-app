package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func (c *controller) StartRobotEpisode(ctx context.Context, request openapi.StartRobotEpisodeRequestObject) (openapi.StartRobotEpisodeResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	input := usecase.StartEpisodeInput{
		EpisodeID:  request.EpisodeId,
		OccurredAt: body.OccurredAt,
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
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	input := usecase.FinishEpisodeInput{
		EpisodeID:  request.EpisodeId,
		OccurredAt: body.OccurredAt,
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
