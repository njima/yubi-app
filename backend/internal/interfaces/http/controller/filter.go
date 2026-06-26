package controller

import (
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
)

func episodeStatuses(values []openapi.EpisodeCollectionStatus) ([]model.EpisodeStatus, error) {
	statuses := make([]model.EpisodeStatus, 0, len(values))
	for _, value := range values {
		status, err := episodeStatusModel(&value)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, *status)
	}
	return statuses, nil
}

func openAPIEpisodeStatus(value model.EpisodeStatus) openapi.EpisodeCollectionStatus {
	return openapi.EpisodeCollectionStatus(value)
}

func episodeStatusModel(value *openapi.EpisodeCollectionStatus) (*model.EpisodeStatus, error) {
	if value == nil {
		return nil, nil
	}
	var status model.EpisodeStatus
	switch *value {
	case openapi.EpisodeCollectionStatusReady:
		status = model.EpisodeStatusReady
	case openapi.EpisodeCollectionStatusRecording:
		status = model.EpisodeStatusRecording
	case openapi.EpisodeCollectionStatusCancel:
		status = model.EpisodeStatusCancel
	case openapi.EpisodeCollectionStatusCompleted:
		status = model.EpisodeStatusCompleted
	default:
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "unknown episode status: %d", *value))
	}
	return &status, nil
}

func openAPISubTaskCollectionStatus(value model.SubTaskCollectionStatus) openapi.SubTaskCollectionStatus {
	return openapi.SubTaskCollectionStatus(value)
}

func subTaskCollectionStatusModel(value openapi.SubTaskCollectionStatus) model.SubTaskCollectionStatus {
	return model.SubTaskCollectionStatus(value)
}

func openAPIExecutionStatus(value model.ExecutionStatus) openapi.ExecutionStatus {
	return openapi.ExecutionStatus(value)
}

func executionStatusModel(value openapi.ExecutionStatus) model.ExecutionStatus {
	return model.ExecutionStatus(value)
}

func openAPIApprovalStatus(value model.ApprovalStatus) openapi.ApprovalStatus {
	return openapi.ApprovalStatus(value)
}

func approvalStatusModel(value openapi.ApprovalStatus) model.ApprovalStatus {
	return model.ApprovalStatus(value)
}

func fleetTrendGranularityModel(value openapi.GetFleetCollectionTrendParamsGranularity) model.FleetTrendGranularity {
	return model.FleetTrendGranularity(value)
}

func openAPIUserRole(value model.UserRole) openapi.UserRole {
	return openapi.UserRole(value)
}

func openAPIUserRolePtr(value model.UserRole) *openapi.UserRole {
	role := openAPIUserRole(value)
	return &role
}

func userRoleModel(value openapi.UserRole) (model.UserRole, error) {
	switch value {
	case openapi.Admin:
		return model.UserRoleAdmin, nil
	case openapi.DataEngineer:
		return model.UserRoleDataEngineer, nil
	case openapi.Manager:
		return model.UserRoleManager, nil
	case openapi.Operator:
		return model.UserRoleOperator, nil
	case openapi.Viewer:
		return model.UserRoleViewer, nil
	default:
		return 0, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "unknown user role: %d", value))
	}
}

func robotStatus(value *openapi.RobotStatus) (*model.RobotStatus, error) {
	return robotStatusModel(value)
}

func openAPIRobotStatus(value model.RobotStatus) openapi.RobotStatus {
	return openapi.RobotStatus(value)
}

func openAPILeaderStatus(value *model.LeaderStatus) *openapi.LeaderStatus {
	if value == nil {
		return nil
	}
	status := openapi.LeaderStatus(*value)
	return &status
}

func leaderStatus(value *openapi.LeaderStatus) (*model.LeaderStatus, error) {
	if value == nil {
		return nil, nil
	}
	var status model.LeaderStatus
	switch *value {
	case openapi.LeaderReady:
		status = model.LeaderStatusReady
	case openapi.LeaderFaulted:
		status = model.LeaderStatusFaulted
	case openapi.LeaderMaintenance:
		status = model.LeaderStatusMaintenance
	default:
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "unknown leader status: %d", *value))
	}
	return &status, nil
}

func robotStatusModel(value *openapi.RobotStatus) (*model.RobotStatus, error) {
	if value == nil {
		return nil, nil
	}
	var status model.RobotStatus
	switch *value {
	case openapi.RobotStatusOnline:
		status = model.RobotStatusOnline
	case openapi.RobotStatusBusy:
		status = model.RobotStatusBusy
	case openapi.RobotStatusOffline:
		status = model.RobotStatusOffline
	case openapi.RobotStatusFaulted:
		status = model.RobotStatusFaulted
	case openapi.RobotStatusMaintenance:
		status = model.RobotStatusMaintenance
	case openapi.RobotStatusReady:
		status = model.RobotStatusReady
	default:
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "unknown robot status: %d", *value))
	}
	return &status, nil
}

func taskStatuses(values []openapi.TaskStatus) []model.TaskStatus {
	statuses := make([]model.TaskStatus, 0, len(values))
	for _, value := range values {
		statuses = append(statuses, model.TaskStatus(value))
	}
	return statuses
}

func taskPriorities(values []openapi.TaskPriority) []model.TaskPriority {
	priorities := make([]model.TaskPriority, 0, len(values))
	for _, value := range values {
		priorities = append(priorities, model.TaskPriority(value))
	}
	return priorities
}

func taskDifficulties(values []openapi.TaskDifficulty) []model.TaskDifficulty {
	difficulties := make([]model.TaskDifficulty, 0, len(values))
	for _, value := range values {
		difficulties = append(difficulties, model.TaskDifficulty(value))
	}
	return difficulties
}

func taskPriorityPtr(value openapi.TaskPriority) *model.TaskPriority {
	priority := model.TaskPriority(value)
	return &priority
}

func taskDifficultyPtr(value openapi.TaskDifficulty) *model.TaskDifficulty {
	difficulty := model.TaskDifficulty(value)
	return &difficulty
}

func taskStatusPtr(value openapi.TaskStatus) *model.TaskStatus {
	status := model.TaskStatus(value)
	return &status
}
