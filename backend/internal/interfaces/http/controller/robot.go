package controller

import (
	"context"
	"encoding/json"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/pagination"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func robotResponseFields(r *model.Robot) (openapi.RobotStatus, *openapi.LeaderStatus) {
	return openAPIRobotStatus(r.Status), openAPILeaderStatus(r.LeaderStatus)
}

func (c *controller) ListRobots(ctx context.Context, request openapi.ListRobotsRequestObject) (openapi.ListRobotsResponseObject, error) {
	pg := pagination.Parse(request.Params.Page, request.Params.Limit)

	filter := repository.RobotListFilter{
		SiteID:     request.Params.SiteId,
		LocationID: request.Params.LocationId,
		Status:     robotStatus(request.Params.Status),
		RobotType:  request.Params.RobotType,
		Search:     request.Params.Search,
		SortBy:     robotSortBy(request.Params.SortBy),
		SortOrder:  sortOrder(request.Params.SortOrder),
	}

	robs, total, err := c.robotUsecase.List(ctx, filter, pg.Page, pg.Limit)
	if err != nil {
		return nil, err
	}

	robots := make([]openapi.Robot, 0, len(robs))
	for _, r := range robs {
		status, leaderStatus := robotResponseFields(r)
		robot := openapi.Robot{
			Id:                         r.IDNatural,
			OrganizationId:             &r.OrganizationID,
			OrganizationName:           &r.OrganizationName,
			SiteId:                     &r.SiteID,
			SiteName:                   &r.SiteName,
			LocationId:                 &r.LocationID,
			LocationName:               &r.LocationName,
			Name:                       r.Name,
			RobotType:                  r.RobotType,
			Status:                     &status,
			LeaderStatus:               leaderStatus,
			ConsecutiveFaultDays:       r.ConsecutiveFaultDays(),
			LeaderConsecutiveFaultDays: r.LeaderConsecutiveFaultDays(),
			LeaderFaultStartedAt:       r.LeaderFaultStartedAt,
			LastHeartbeatAt:            r.LastHeartbeatAt,
			OfflineReason:              r.OfflineReason,
			RobotConfig:                mapPtrFromRawMessagePtr(r.RobotConfig),
			ActiveEpisodeId:            r.ActiveEpisodeID,
			ActiveUserId:               r.ActiveUserID,
		}
		if op, err := c.robotOperatorUsecase.Get(ctx, r.IDNatural); err != nil {
			c.logger.Warn().Err(err).Str("robot_id", r.IDNatural).Msg("failed to get operator for robot list")
		} else if op != nil {
			robot.ActiveOperator = &openapi.RobotOperator{
				UserId:           op.UserID,
				DisplayName:      op.DisplayName,
				OrganizationName: op.OrganizationName,
			}
		}
		robots = append(robots, robot)
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
	filter := repository.RobotTypeFilter{
		SiteID:     request.Params.SiteId,
		LocationID: request.Params.LocationId,
		Status:     robotStatus(request.Params.Status),
	}
	types, err := c.robotUsecase.ListTypes(ctx, filter)
	if err != nil {
		return nil, err
	}
	return openapi.ListRobotTypes200JSONResponse(types), nil
}

func (c *controller) CreateRobot(ctx context.Context, request openapi.CreateRobotRequestObject) (openapi.CreateRobotResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	var cam json.RawMessage
	if request.Body.RobotConfig != nil {
		if b, err := json.Marshal(*request.Body.RobotConfig); err == nil {
			cam = b
		}
	}

	input := usecase.RobotCreateInput{
		OrganizationID: request.Body.OrganizationId,
		LocationID:     request.Body.LocationId,
		Name:           request.Body.Name,
		RobotType:      request.Body.RobotType,
		LeaderStatus:   leaderStatus(request.Body.LeaderStatus),
		RobotConfig:    &cam,
	}

	rob, err := c.robotUsecase.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	status, leaderStatus := robotResponseFields(&rob)
	return openapi.CreateRobot201JSONResponse{
		Id:                         rob.IDNatural,
		OrganizationId:             &rob.OrganizationID,
		OrganizationName:           &rob.OrganizationName,
		SiteId:                     &rob.SiteID,
		SiteName:                   &rob.SiteName,
		LocationId:                 &rob.LocationID,
		LocationName:               &rob.LocationName,
		Name:                       rob.Name,
		RobotType:                  rob.RobotType,
		Status:                     &status,
		LeaderStatus:               leaderStatus,
		ConsecutiveFaultDays:       rob.ConsecutiveFaultDays(),
		LeaderConsecutiveFaultDays: rob.LeaderConsecutiveFaultDays(),
		LeaderFaultStartedAt:       rob.LeaderFaultStartedAt,
		LastHeartbeatAt:            rob.LastHeartbeatAt,
		OfflineReason:              rob.OfflineReason,
		RobotConfig:                mapPtrFromRawMessagePtr(rob.RobotConfig),
	}, nil
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

	status, leaderStatus := robotResponseFields(&rob)
	response := openapi.GetRobotById200JSONResponse{
		Id:                         rob.IDNatural,
		OrganizationId:             &rob.OrganizationID,
		OrganizationName:           &rob.OrganizationName,
		SiteId:                     &rob.SiteID,
		SiteName:                   &rob.SiteName,
		LocationId:                 &rob.LocationID,
		LocationName:               &rob.LocationName,
		Name:                       rob.Name,
		RobotType:                  rob.RobotType,
		Status:                     &status,
		LeaderStatus:               leaderStatus,
		ConsecutiveFaultDays:       rob.ConsecutiveFaultDays(),
		LeaderConsecutiveFaultDays: rob.LeaderConsecutiveFaultDays(),
		LeaderFaultStartedAt:       rob.LeaderFaultStartedAt,
		LastHeartbeatAt:            rob.LastHeartbeatAt,
		OfflineReason:              rob.OfflineReason,
		RobotConfig:                mapPtrFromRawMessagePtr(rob.RobotConfig),
		ActiveEpisodeId:            rob.ActiveEpisodeID,
		ActiveUserId:               rob.ActiveUserID,
	}
	if op, err := c.robotOperatorUsecase.Get(ctx, rob.IDNatural); err != nil {
		c.logger.Warn().Err(err).Str("robot_id", rob.IDNatural).Msg("failed to get operator for robot detail")
	} else if op != nil {
		response.ActiveOperator = &openapi.RobotOperator{
			UserId:           op.UserID,
			DisplayName:      op.DisplayName,
			OrganizationName: op.OrganizationName,
		}
	}
	return response, nil
}

func (c *controller) UpdateRobotById(ctx context.Context, request openapi.UpdateRobotByIdRequestObject) (openapi.UpdateRobotByIdResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	input := usecase.RobotUpdateInput{
		ID: request.RobotId,
	}
	if request.Body.Name != nil {
		input.Name = request.Body.Name
	}
	if request.Body.RobotType != nil {
		input.RobotType = request.Body.RobotType
	}
	input.Status = robotStatusModel(request.Body.Status)
	if request.Body.LastHeartbeatAt != nil {
		input.LastHeartbeatAt = request.Body.LastHeartbeatAt
	}
	if request.Body.OfflineReason != nil {
		input.OfflineReason = request.Body.OfflineReason
	}
	if request.Body.RobotConfig != nil {
		cam, err := json.Marshal(request.Body.RobotConfig)
		if err != nil {
			return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeInternal, "failed to marshal robot config"))
		}
		rawMsg := json.RawMessage(cam)
		input.RobotConfig = &rawMsg
	}
	if request.Body.LeaderStatus != nil {
		input.LeaderStatus = leaderStatus(request.Body.LeaderStatus)
		input.HasLeaderStatus = true
	}

	rob, err := c.robotUsecase.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	status, leaderStatus := robotResponseFields(&rob)
	return openapi.UpdateRobotById200JSONResponse{
		Id:                         rob.IDNatural,
		OrganizationId:             &rob.OrganizationID,
		OrganizationName:           &rob.OrganizationName,
		SiteId:                     &rob.SiteID,
		SiteName:                   &rob.SiteName,
		LocationId:                 &rob.LocationID,
		LocationName:               &rob.LocationName,
		Name:                       rob.Name,
		RobotType:                  rob.RobotType,
		Status:                     &status,
		LeaderStatus:               leaderStatus,
		ConsecutiveFaultDays:       rob.ConsecutiveFaultDays(),
		LeaderConsecutiveFaultDays: rob.LeaderConsecutiveFaultDays(),
		LeaderFaultStartedAt:       rob.LeaderFaultStartedAt,
		LastHeartbeatAt:            rob.LastHeartbeatAt,
		OfflineReason:              rob.OfflineReason,
		RobotConfig:                mapPtrFromRawMessagePtr(rob.RobotConfig),
	}, nil
}
