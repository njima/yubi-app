package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/requestctx"
)

func (c *controller) GetRobotOperator(ctx context.Context, request openapi.GetRobotOperatorRequestObject) (openapi.GetRobotOperatorResponseObject, error) {
	operator, err := c.robotOperatorUsecase.Get(ctx, request.RobotId)
	if err != nil {
		return nil, err
	}
	if operator == nil {
		return openapi.GetRobotOperator204Response{}, nil
	}
	return openapi.GetRobotOperator200JSONResponse{
		UserId:           operator.UserID,
		DisplayName:      operator.DisplayName,
		OrganizationName: operator.OrganizationName,
	}, nil
}

func (c *controller) SetRobotOperator(ctx context.Context, request openapi.SetRobotOperatorRequestObject) (openapi.SetRobotOperatorResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	userID, err := requestctx.UserID(ctx)
	if err != nil {
		return nil, err
	}

	operator := model.RobotOperator{
		UserID:           userID,
		DisplayName:      request.Body.DisplayName,
		OrganizationName: request.Body.OrganizationName,
	}

	existing, err := c.robotOperatorUsecase.Set(ctx, request.RobotId, operator)
	if err != nil {
		if apperror.SameKind(err, apperror.KindConflict) && existing != nil {
			return openapi.SetRobotOperator409JSONResponse{
				UserId:           existing.UserID,
				DisplayName:      existing.DisplayName,
				OrganizationName: existing.OrganizationName,
			}, nil
		}
		return nil, err
	}

	return openapi.SetRobotOperator200Response{}, nil
}

func (c *controller) ClearRobotOperator(ctx context.Context, request openapi.ClearRobotOperatorRequestObject) (openapi.ClearRobotOperatorResponseObject, error) {
	userID, err := requestctx.UserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := c.robotOperatorUsecase.Clear(ctx, request.RobotId, userID); err != nil {
		return nil, err
	}

	return openapi.ClearRobotOperator204Response{}, nil
}
