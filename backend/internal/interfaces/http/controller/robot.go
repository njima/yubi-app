package controller

import (
	"context"
	"encoding/json"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

func (c *controller) ListRobots(ctx context.Context, request openapi.ListRobotsRequestObject) (openapi.ListRobotsResponseObject, error) {
	pg := pagination.Parse(request.Params.Page, request.Params.Limit)

	status, err := robotStatus(request.Params.Status)
	if err != nil {
		return nil, err
	}
	filter := usecase.RobotListFilter{
		SiteID:     request.Params.SiteId,
		LocationID: request.Params.LocationId,
		Status:     status,
		RobotType:  request.Params.RobotType,
		Search:     request.Params.Search,
		SortBy:     robotSortBy(request.Params.SortBy),
		SortOrder:  sortOrder(request.Params.SortOrder),
	}

	robs, total, err := c.robotUsecase.List(ctx, filter, pg.Page, pg.Limit)
	if err != nil {
		return nil, err
	}

	robots := robotResponses(robs)
	for i, r := range robs {
		if op, err := c.robotOperatorUsecase.Get(ctx, r.IDNatural); err != nil {
			c.logger.Warn().Err(err).Str("robot_id", r.IDNatural).Msg("failed to get operator for robot list")
		} else if op != nil {
			operator := robotOperatorResponse(*op)
			robots[i].ActiveOperator = &operator
		}
	}

	return openapi.ListRobots200JSONResponse{
		Robots: robots,
		Pagination: openapi.Pagination{
			Count: total,
			Page:  pg.Page,
			Limit: pg.Limit,
		},
	}, nil
}

func (c *controller) ListRobotTypes(ctx context.Context, request openapi.ListRobotTypesRequestObject) (openapi.ListRobotTypesResponseObject, error) {
	status, err := robotStatus(request.Params.Status)
	if err != nil {
		return nil, err
	}
	filter := usecase.RobotTypeFilter{
		SiteID:     request.Params.SiteId,
		LocationID: request.Params.LocationId,
		Status:     status,
	}
	types, err := c.robotUsecase.ListTypes(ctx, filter)
	if err != nil {
		return nil, err
	}
	return openapi.ListRobotTypes200JSONResponse(types), nil
}

func (c *controller) CreateRobot(ctx context.Context, request openapi.CreateRobotRequestObject) (openapi.CreateRobotResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	var cam json.RawMessage
	if body.RobotConfig != nil {
		if b, err := json.Marshal(*body.RobotConfig); err == nil {
			cam = b
		}
	}

	input := usecase.RobotCreateInput{
		OrganizationID: body.OrganizationId,
		LocationID:     body.LocationId,
		Name:           body.Name,
		RobotType:      body.RobotType,
		RobotConfig:    &cam,
	}
	mappedLeaderStatus, err := leaderStatus(body.LeaderStatus)
	if err != nil {
		return nil, err
	}
	input.LeaderStatus = mappedLeaderStatus

	rob, err := c.robotUsecase.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return openapi.CreateRobot201JSONResponse(robotResponse(rob)), nil
}

func (c *controller) DeleteRobotById(ctx context.Context, request openapi.DeleteRobotByIdRequestObject) (openapi.DeleteRobotByIdResponseObject, error) {
	if err := c.robotUsecase.Delete(ctx, request.RobotId); err != nil {
		return nil, err
	}

	return openapi.DeleteRobotById204Response{}, nil
}

func (c *controller) GetRobotById(ctx context.Context, request openapi.GetRobotByIdRequestObject) (openapi.GetRobotByIdResponseObject, error) {
	rob, err := c.robotUsecase.GetByID(ctx, request.RobotId)
	if err != nil {
		return nil, err
	}

	response := robotResponse(rob)
	if op, err := c.robotOperatorUsecase.Get(ctx, rob.IDNatural); err != nil {
		c.logger.Warn().Err(err).Str("robot_id", rob.IDNatural).Msg("failed to get operator for robot detail")
	} else if op != nil {
		operator := robotOperatorResponse(*op)
		response.ActiveOperator = &operator
	}
	return openapi.GetRobotById200JSONResponse(response), nil
}

func (c *controller) UpdateRobotById(ctx context.Context, request openapi.UpdateRobotByIdRequestObject) (openapi.UpdateRobotByIdResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	input := usecase.RobotUpdateInput{
		ID: request.RobotId,
	}
	if body.Name != nil {
		input.Name = body.Name
	}
	if body.RobotType != nil {
		input.RobotType = body.RobotType
	}
	status, err := robotStatusModel(body.Status)
	if err != nil {
		return nil, err
	}
	input.Status = status
	if body.LastHeartbeatAt != nil {
		input.LastHeartbeatAt = body.LastHeartbeatAt
	}
	if body.OfflineReason != nil {
		input.OfflineReason = body.OfflineReason
	}
	if body.RobotConfig != nil {
		cam, err := json.Marshal(body.RobotConfig)
		if err != nil {
			return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeInternal, "failed to marshal robot config"))
		}
		rawMsg := json.RawMessage(cam)
		input.RobotConfig = &rawMsg
	}
	if body.LeaderStatus != nil {
		input.LeaderStatus, err = leaderStatus(body.LeaderStatus)
		if err != nil {
			return nil, err
		}
		input.HasLeaderStatus = true
	}

	rob, err := c.robotUsecase.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	return openapi.UpdateRobotById200JSONResponse(robotResponse(rob)), nil
}
