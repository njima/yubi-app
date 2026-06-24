package bunconv

import (
	"github.com/airoa-org/yubi-app/backend/internal/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

func EntityToTaskModel(et entity.Task, etv entity.TaskVersion) model.Task {
	manualURL := ""
	if et.ManualURL != nil {
		manualURL = *et.ManualURL
	}

	return model.Task{
		ID:                    et.ID,
		IDNatural:             et.IDNatural,
		Name:                  et.Name,
		Description:           et.Description,
		ManualURL:             manualURL,
		Priority:              taskPriorityPtr(et.Priority),
		Difficulty:            taskDifficultyPtr(et.Difficulty),
		Status:                taskStatusPtr(et.Status),
		Deadline:              et.Deadline,
		RobotType:             et.RobotType,
		TargetDurationSeconds: etv.TargetDurationSeconds,
		Version:               etv.Version,
		VersionDisplayName:    etv.DisplayName,
		IsActive:              etv.IsActive,
		CreatedAt:             et.CreatedAt,
		UpdatedAt:             &et.UpdatedAt,
	}
}

func taskPriorityPtr(priority model.TaskPriority) *model.TaskPriority {
	return &priority
}

func taskDifficultyPtr(difficulty model.TaskDifficulty) *model.TaskDifficulty {
	return &difficulty
}

func taskStatusPtr(status model.TaskStatus) *model.TaskStatus {
	return &status
}
