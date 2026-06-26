package model

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ExecutionStatus int

const (
	ExecutionStatusReady     ExecutionStatus = 0
	ExecutionStatusStarted   ExecutionStatus = 1
	ExecutionStatusCancelled ExecutionStatus = 2
	ExecutionStatusFinished  ExecutionStatus = 3
)

func (s ExecutionStatus) IsTerminal() bool {
	return s == ExecutionStatusFinished || s == ExecutionStatusCancelled
}

func (s ExecutionStatus) IsWorkflowResolved() bool {
	return s.IsTerminal()
}

func (s ExecutionStatus) IsSuccessfulCompletion() bool {
	return s == ExecutionStatusFinished
}

// EpisodeSubTaskExecution represents the database record for episode_sub_task_execution table
type EpisodeSubTaskExecution struct {
	ID               int64
	IDNatural        string
	OrganizationID   string
	EpisodeSubTaskID string
	ExecutionStatus  ExecutionStatus
	StartedAt        *time.Time
	FinishedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        *time.Time
}

type EpisodeSubTaskExecutions []*EpisodeSubTaskExecution

func InitEpisodeSubTaskExecution(organizationID, episodeSubTaskID string) (EpisodeSubTaskExecution, error) {
	idNatural, err := InitID()
	if err != nil {
		return EpisodeSubTaskExecution{}, err
	}

	exe := EpisodeSubTaskExecution{
		IDNatural:        idNatural,
		OrganizationID:   organizationID,
		EpisodeSubTaskID: episodeSubTaskID,
		ExecutionStatus:  ExecutionStatusReady,
		CreatedAt:        time.Now(),
	}

	if err := exe.validate(); err != nil {
		return EpisodeSubTaskExecution{}, err
	}

	return exe, nil
}

func NewEpisodeSubTaskExecution(
	id int64,
	idNatural,
	organizationID,
	episodeSubTaskID string,
	executionStatus ExecutionStatus,
	startedAt,
	finishedAt *time.Time,
	createdAt time.Time,
	updatedAt *time.Time,
) EpisodeSubTaskExecution {
	return EpisodeSubTaskExecution{
		ID:               id,
		IDNatural:        idNatural,
		OrganizationID:   organizationID,
		EpisodeSubTaskID: episodeSubTaskID,
		ExecutionStatus:  executionStatus,
		StartedAt:        startedAt,
		FinishedAt:       finishedAt,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

func (exe EpisodeSubTaskExecution) validate() error {
	if err := (validation.Errors{
		"id_natural":          validation.Validate(exe.IDNatural, validation.Required.Error("id_natural is required")),
		"organization_id":     validation.Validate(exe.OrganizationID, validation.Required.Error("organization_id is required")),
		"episode_sub_task_id": validation.Validate(exe.EpisodeSubTaskID, validation.Required.Error("episode_sub_task_id is required")),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "episode_sub_task_execution validation failed: %v", err))
	}
	return nil
}

func (exe EpisodeSubTaskExecution) IsTerminal() bool {
	return exe.ExecutionStatus.IsTerminal()
}

func (exe EpisodeSubTaskExecution) IsWorkflowResolved() bool {
	return exe.ExecutionStatus.IsWorkflowResolved()
}

func (exe EpisodeSubTaskExecution) IsSuccessfulCompletion() bool {
	return exe.ExecutionStatus.IsSuccessfulCompletion()
}

// CanStart checks if the execution can be started
func (exe *EpisodeSubTaskExecution) CanStart() error {
	if exe.ExecutionStatus != ExecutionStatusReady {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "execution status must be Ready to start, current: %v", exe.ExecutionStatus),
		)
	}
	return nil
}

// Start transitions the execution from Ready to Started
func (exe *EpisodeSubTaskExecution) Start(occurredAt time.Time) error {
	if err := exe.CanStart(); err != nil {
		return err
	}
	exe.ExecutionStatus = ExecutionStatusStarted
	exe.StartedAt = &occurredAt
	return nil
}

// CanFinish checks if the execution can be finished
func (exe *EpisodeSubTaskExecution) CanFinish() error {
	if exe.ExecutionStatus != ExecutionStatusStarted {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "execution status must be Started to finish, current: %v", exe.ExecutionStatus),
		)
	}
	return nil
}

// Finish transitions the execution from Started to Finished
func (exe *EpisodeSubTaskExecution) Finish(occurredAt time.Time) error {
	if err := exe.CanFinish(); err != nil {
		return err
	}
	exe.ExecutionStatus = ExecutionStatusFinished
	exe.FinishedAt = &occurredAt
	return nil
}

// Cancel transitions the execution to Cancelled
func (exe *EpisodeSubTaskExecution) Cancel() error {
	if exe.ExecutionStatus == ExecutionStatusCancelled {
		return nil // Already cancelled
	}
	if exe.ExecutionStatus == ExecutionStatusFinished {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "cannot cancel finished execution"),
		)
	}
	exe.ExecutionStatus = ExecutionStatusCancelled
	return nil
}
