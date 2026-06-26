package model

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SubTask struct {
	ID                    int64
	IDNatural             string
	OrganizationID        string
	TaskVersionID         string
	Name                  string
	Description           *string
	OrderIndex            int
	CreatedAt             time.Time
	UpdatedAt             *time.Time
	TargetDurationSeconds *int
}

type SubTasks []*SubTask

func InitSubTask(organizationID, taskVersionID, name string, orderIndex int, description *string, targetDurationSeconds *int) (SubTask, error) {
	ID, err := InitID()
	if err != nil {
		return SubTask{}, err
	}

	subtask := SubTask{
		IDNatural:             ID,
		OrganizationID:        organizationID,
		TaskVersionID:         taskVersionID,
		Name:                  name,
		Description:           description,
		OrderIndex:            orderIndex,
		CreatedAt:             time.Now(),
		TargetDurationSeconds: targetDurationSeconds,
	}

	if err := subtask.validate(); err != nil {
		return SubTask{}, err
	}

	return subtask, nil
}

func (s SubTask) validate() error {
	if err := (validation.Errors{
		"id_natural":      validation.Validate(s.IDNatural, validation.Required.Error("id_natural is required")),
		"organization_id": validation.Validate(s.OrganizationID, validation.Required.Error("organization_id is required")),
		"task_version_id": validation.Validate(s.TaskVersionID, validation.Required.Error("task_version_id is required")),
		"name": validation.Validate(
			s.Name,
			validation.Required.Error("name is required"),
			validation.RuneLength(1, 100).Error("name must be between 1 and 100 characters"),
		),
		"order_index":             validation.Validate(s.OrderIndex, validation.Min(0).Error("order_index must be non-negative")),
		"target_duration_seconds": validation.Validate(s.TargetDurationSeconds, validation.Min(1).Error("target_duration_seconds must be at least 1")),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "subtask validation failed: %v", err))
	}

	return nil
}

func (s *SubTask) SetName(name string) error {
	s.Name = name
	return s.validate()
}

func (s *SubTask) SetDescription(description string) error {
	s.Description = &description
	return s.validate()
}

// NextOrderIndex returns the order_index to assign to a new subtask.
// Pass -1 as currentMax when no subtasks exist yet.
func NextOrderIndex(currentMax int) int {
	return currentMax + 1
}

func (s *SubTask) SetOrderIndex(orderIndex int) error {
	s.OrderIndex = orderIndex
	return s.validate()
}
