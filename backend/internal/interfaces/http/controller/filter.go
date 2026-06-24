package controller

import (
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
