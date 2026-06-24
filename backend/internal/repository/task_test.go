package repository

import (
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

func TestTaskListFilterUsesDomainTaskEnums(t *testing.T) {
	filter := TaskListFilter{
		Statuses:     []model.TaskStatus{model.TaskStatusPlanning},
		Priorities:   []model.TaskPriority{model.TaskPriorityNormal},
		Difficulties: []model.TaskDifficulty{model.TaskDifficultyB},
	}

	var _ []model.TaskStatus = filter.Statuses
	var _ []model.TaskPriority = filter.Priorities
	var _ []model.TaskDifficulty = filter.Difficulties
}

func TestTaskExportRowUsesDomainTaskEnums(t *testing.T) {
	row := TaskExportRow{
		Priority:   model.TaskPriorityNormal,
		Difficulty: model.TaskDifficultyB,
		Status:     model.TaskStatusPlanning,
	}

	var _ model.TaskPriority = row.Priority
	var _ model.TaskDifficulty = row.Difficulty
	var _ model.TaskStatus = row.Status
}
