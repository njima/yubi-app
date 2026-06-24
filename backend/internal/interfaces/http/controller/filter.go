package controller

import (
	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

func episodeStatuses(values []openapi.EpisodeCollectionStatus) []repository.EpisodeStatus {
	statuses := make([]repository.EpisodeStatus, 0, len(values))
	for _, value := range values {
		statuses = append(statuses, repository.EpisodeStatus(value))
	}
	return statuses
}

func robotStatus(value *openapi.RobotStatus) *repository.RobotFilterStatus {
	if value == nil {
		return nil
	}
	status := repository.RobotFilterStatus(*value)
	return &status
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

func taskStatuses(values []openapi.TaskStatus) []repository.TaskStatus {
	statuses := make([]repository.TaskStatus, 0, len(values))
	for _, value := range values {
		statuses = append(statuses, repository.TaskStatus(value))
	}
	return statuses
}

func taskPriorities(values []openapi.TaskPriority) []repository.TaskPriority {
	priorities := make([]repository.TaskPriority, 0, len(values))
	for _, value := range values {
		priorities = append(priorities, repository.TaskPriority(value))
	}
	return priorities
}

func taskDifficulties(values []openapi.TaskDifficulty) []repository.TaskDifficulty {
	difficulties := make([]repository.TaskDifficulty, 0, len(values))
	for _, value := range values {
		difficulties = append(difficulties, repository.TaskDifficulty(value))
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
