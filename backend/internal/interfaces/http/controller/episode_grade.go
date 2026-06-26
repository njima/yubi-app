package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

func (c *controller) GetMyEpisodeGrade(ctx context.Context, request openapi.GetMyEpisodeGradeRequestObject) (openapi.GetMyEpisodeGradeResponseObject, error) {
	userID, err := requestctx.UserID(ctx)
	if err != nil {
		return nil, err
	}

	grade, err := c.episodeGradeUsecase.GetMyGrade(ctx, request.EpisodeId, userID)
	if err != nil {
		return nil, err
	}
	if grade == nil {
		return openapi.GetMyEpisodeGrade404Response{}, nil
	}

	userName, err := c.lookupUserName(ctx, userID)
	if err != nil {
		return nil, err
	}

	return openapi.GetMyEpisodeGrade200JSONResponse(episodeGradeResponse(*grade, userName)), nil
}

func (c *controller) UpdateMyEpisodeGrade(ctx context.Context, request openapi.UpdateMyEpisodeGradeRequestObject) (openapi.UpdateMyEpisodeGradeResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	userID, err := requestctx.UserID(ctx)
	if err != nil {
		return nil, err
	}
	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := c.episodeUsecase.GetByID(ctx, request.EpisodeId); err != nil {
		return nil, err
	}

	saved, err := c.episodeGradeUsecase.Upsert(ctx, usecase.EpisodeGradeUpsertInput{
		EpisodeID:      request.EpisodeId,
		UserID:         userID,
		OrganizationID: orgID,
		Grade:          body.Grade,
		Comment:        body.Comment,
	})
	if err != nil {
		return nil, err
	}

	userName, err := c.lookupUserName(ctx, userID)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateMyEpisodeGrade200JSONResponse(episodeGradeResponse(saved, userName)), nil
}

func (c *controller) lookupUserName(ctx context.Context, userID string) (string, error) {
	user, err := c.userUsecase.GetByNaturalID(ctx, userID)
	if err != nil {
		return "", err
	}
	return user.Name, nil
}

func (c *controller) ListEpisodeGrades(ctx context.Context, request openapi.ListEpisodeGradesRequestObject) (openapi.ListEpisodeGradesResponseObject, error) {
	if _, err := c.episodeUsecase.GetByID(ctx, request.EpisodeId); err != nil {
		return nil, err
	}

	pg := pagination.Parse(request.Params.Page, request.Params.Limit)

	items, total, err := c.episodeGradeUsecase.List(ctx, request.EpisodeId, pg.Page, pg.Limit)
	if err != nil {
		return nil, err
	}

	grades := make([]openapi.EpisodeGrade, 0, len(items))
	for _, it := range items {
		grades = append(grades, episodeGradeResponse(it.Grade, it.UserName))
	}

	return openapi.ListEpisodeGrades200JSONResponse{
		Grades: grades,
		Pagination: openapi.Pagination{
			Count: total,
			Page:  pg.Page,
			Limit: pg.Limit,
		},
	}, nil
}
