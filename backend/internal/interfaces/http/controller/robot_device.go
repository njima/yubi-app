package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/requestctx"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func (c *controller) GetRobotMe(ctx context.Context, request openapi.GetRobotMeRequestObject) (openapi.GetRobotMeResponseObject, error) {
	robotID, err := requestctx.RobotID(ctx)
	if err != nil {
		return nil, err
	}

	rob, err := c.robotUsecase.GetByID(ctx, robotID)
	if err != nil {
		return nil, err
	}

	status := openAPIRobotStatus(rob.Status)
	resp := openapi.GetRobotMe200JSONResponse{
		Id:               rob.IDNatural,
		OrganizationId:   &rob.OrganizationID,
		OrganizationName: &rob.OrganizationName,
		SiteId:           &rob.SiteID,
		SiteName:         &rob.SiteName,
		LocationId:       &rob.LocationID,
		LocationName:     &rob.LocationName,
		Name:             rob.Name,
		RobotType:        rob.RobotType,
		Status:           &status,
		LastHeartbeatAt:  rob.LastHeartbeatAt,
		OfflineReason:    rob.OfflineReason,
		RobotConfig:      mapPtrFromRawMessagePtr(rob.RobotConfig),
		ActiveEpisodeId:  rob.ActiveEpisodeID,
		ActiveUserId:     rob.ActiveUserID,
	}

	operator, err := c.robotOperatorUsecase.Get(ctx, robotID)
	if err != nil {
		c.logger.Warn().Err(err).Str("robot_id", robotID).Msg("failed to get robot operator")
	} else if operator != nil {
		resp.ActiveOperator = &openapi.RobotOperator{
			UserId:           operator.UserID,
			DisplayName:      operator.DisplayName,
			OrganizationName: operator.OrganizationName,
		}
	}

	return resp, nil
}

func (c *controller) UpdateRobotStatus(ctx context.Context, request openapi.UpdateRobotStatusRequestObject) (openapi.UpdateRobotStatusResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	robotID, err := requestctx.RobotID(ctx)
	if err != nil {
		return nil, err
	}

	exists, err := c.robotDeviceUsecase.RobotExists(ctx, robotID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeRobotNotFound, "robot not found"))
	}

	status := convertToRobotStatus(robotID, request.Body)

	if err := c.robotDeviceUsecase.UpdateRobotStatus(ctx, status); err != nil {
		return nil, err
	}

	return openapi.UpdateRobotStatus200Response{}, nil
}

func convertToRobotStatus(robotID string, req *openapi.RobotStatusUpdateRequest) usecase.RobotDeviceStatus {
	var metrics []usecase.RobotMetric
	if req.Status.Metrics != nil && len(*req.Status.Metrics) > 0 {
		for _, m := range *req.Status.Metrics {
			var labels map[string]string
			if m.Labels != nil {
				labels = *m.Labels
			}
			metrics = append(metrics, usecase.RobotMetric{
				Name:   m.Name,
				Type:   string(m.Type),
				Unit:   m.Unit,
				Value:  m.Value,
				Labels: labels,
			})
		}
	}

	var gate *model.GateConditionStatus
	if req.Status.GateConditions != nil {
		g := req.Status.GateConditions
		groups := make(map[string]model.GateGroupStatus, len(g.Groups))
		for name, grp := range g.Groups {
			conditions := make([]model.GateCondition, len(grp.Conditions))
			for i, c := range grp.Conditions {
				conditions[i] = model.GateCondition{
					Name:       c.Name,
					Passed:     c.Passed,
					Reason:     c.Reason,
					Escalation: c.Escalation,
				}
			}
			groups[name] = model.GateGroupStatus{
				Level:      grp.Level,
				Settled:    grp.Settled,
				Conditions: conditions,
			}
		}
		gate = &model.GateConditionStatus{
			GateLevel: g.GateLevel,
			Groups:    groups,
		}
	}

	return usecase.RobotDeviceStatus{
		RobotID:    robotID,
		RobotType:  req.RobotType,
		ReportedAt: req.ReportedAt,
		Status: usecase.RobotDeviceStatusDetail{
			Battery: usecase.RobotBatteryStatus{
				Pct:      req.Status.Battery.Pct,
				Charging: req.Status.Battery.Charging,
			},
			Connection: usecase.RobotConnectionStatus{
				QualityPct: req.Status.Connection.QualityPct,
			},
			UptimeSec:      req.Status.UptimeSec,
			Metrics:        metrics,
			GateConditions: gate,
		},
	}
}
