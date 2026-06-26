package handler

import (
	"math"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func robotStatusStreamDetail(status usecase.RobotDeviceStatusDetail) openapi.RobotStatusStreamDetail {
	detail := openapi.RobotStatusStreamDetail{
		BatteryPct:    status.Battery.Pct,
		ConnectionPct: status.Connection.QualityPct,
		UptimeSec:     int(math.Round(status.UptimeSec)),
	}
	if g := status.GateConditions; g != nil {
		oGate := convertGateToOpenAPI(g)
		detail.GateConditions = &oGate
	}
	return detail
}

func emptyRobotStatusStreamResponse(robotID string) openapi.RobotStatusStreamResponse {
	return openapi.RobotStatusStreamResponse{
		RobotId:   robotID,
		RobotType: "",
		Status: openapi.RobotStatusStreamDetail{
			BatteryPct:    0,
			ConnectionPct: 0,
			UptimeSec:     0,
		},
	}
}

func buildEpisodeStreamResponse(
	ep *model.Episode,
	subtasks []openapi.EpisodeSubTask,
	task *model.Task,
	taskVersionDisplayName *string,
) openapi.Episode {
	response := openapi.Episode{
		Id:            ep.IDNatural,
		LocationId:    ep.LocationID,
		UserId:        ep.UserID,
		RobotId:       ep.RobotID,
		Status:        openapi.EpisodeCollectionStatus(ep.Status),
		TaskId:        ep.TaskID,
		TaskVersionId: ep.TaskVersionID,
		StartedAt:     ep.StartedAt,
		EndedAt:       ep.FinishedAt,
		ErrorDetails:  ep.ErrorDetails,
		Subtasks:      &subtasks,
		CreatedAt:     ep.CreatedAt,
		RecordedBy:    ep.RecordedByID,
		AverageGrade:  ep.AverageGrade,
		GradeCount:    &ep.GradeCount,
	}
	if task != nil {
		response.TaskName = &task.Name
		response.TaskDescription = task.Description
	}
	if len(ep.ParameterValues) > 0 {
		response.ParameterValues = &ep.ParameterValues
	}
	if taskVersionDisplayName != nil {
		response.TaskVersionDisplayName = taskVersionDisplayName
	}
	return response
}

func buildTeleopTaskMeta(task *model.Task, version *model.TaskVersion) teleopTaskMeta {
	return teleopTaskMeta{
		Id:          task.IDNatural,
		Name:        task.Name,
		Description: task.Description,
		ManualURL:   task.ManualURL,
		Version:     version.Version,
	}
}

func convertGateToOpenAPI(g *model.GateConditionStatus) openapi.GateConditionStatus {
	groups := make(map[string]openapi.GateGroupStatus, len(g.Groups))
	for name, grp := range g.Groups {
		conditions := make([]openapi.GateCondition, len(grp.Conditions))
		for i, c := range grp.Conditions {
			conditions[i] = openapi.GateCondition{
				Name:       c.Name,
				Passed:     c.Passed,
				Reason:     c.Reason,
				Escalation: c.Escalation,
			}
		}
		groups[name] = openapi.GateGroupStatus{
			Level:      grp.Level,
			Settled:    grp.Settled,
			Conditions: conditions,
		}
	}
	return openapi.GateConditionStatus{
		GateLevel: g.GateLevel,
		Groups:    groups,
	}
}
