package model

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type EpisodeStatus int

const (
	EpisodeStatusReady     EpisodeStatus = 0
	EpisodeStatusRecording EpisodeStatus = 1
	EpisodeStatusCancel    EpisodeStatus = 2
	EpisodeStatusCompleted EpisodeStatus = 3
)

type Episode struct {
	ID              int64
	IDNatural       string
	OrganizationID  string
	TaskID          string
	TaskVersionID   string
	LocationID      string
	RobotID         string
	UserID          string
	RecordedByID    *string
	StartedAt       *time.Time
	FinishedAt      *time.Time
	Status          EpisodeStatus
	ErrorDetails    *string
	ParameterValues map[string]string
	CreatedAt       time.Time
	UpdatedAt       *time.Time

	// Populated by usecase when listing/fetching. nil + 0 means no grades.
	AverageGrade *float64
	GradeCount   int
}

type Episodes []*Episode

func InitEpisode(organizationID, taskID, locationID, robotID, userID string, recordedByID *string) (Episode, error) {
	idNatural, err := InitID()
	if err != nil {
		return Episode{}, err
	}

	episode := Episode{
		IDNatural:      idNatural,
		OrganizationID: organizationID,
		TaskID:         taskID,
		LocationID:     locationID,
		RobotID:        robotID,
		UserID:         userID,
		RecordedByID:   recordedByID,
		Status:         EpisodeStatusReady,
		CreatedAt:      time.Now(),
	}

	if err := episode.validate(); err != nil {
		return Episode{}, err
	}

	return episode, nil
}

func NewEpisode(
	id int64,
	idNatural,
	taskID,
	taskVersionID,
	locationID,
	robotID,
	userID string,
	errorDetails *string,
	startedAt,
	finishedAt *time.Time,
	status EpisodeStatus,
	createdAt time.Time,
	updatedAt *time.Time,
	recordedByID *string,
) Episode {
	return Episode{
		ID:            id,
		IDNatural:     idNatural,
		TaskID:        taskID,
		TaskVersionID: taskVersionID,
		LocationID:    locationID,
		RobotID:       robotID,
		UserID:        userID,
		RecordedByID:  recordedByID,
		StartedAt:     startedAt,
		FinishedAt:    finishedAt,
		ErrorDetails:  errorDetails,
		Status:        status,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (e Episode) validate() error {
	if err := (validation.Errors{
		"id_natural":      validation.Validate(e.IDNatural, validation.Required.Error("id_natural is required")),
		"organization_id": validation.Validate(e.OrganizationID, validation.Required.Error("organization_id is required")),
		"task_id":         validation.Validate(e.TaskID, validation.Required.Error("task_id is required")),
		"location_id":     validation.Validate(e.LocationID, validation.Required.Error("location_id is required")),
		"robot_id":        validation.Validate(e.RobotID, validation.Required.Error("robot_id is required")),
		"user_id":         validation.Validate(e.UserID, validation.Required.Error("user_id is required")),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "episode validation failed: %v", err))
	}

	return nil
}

func (e *Episode) SetTaskVersionID(taskVersionID string) error {
	e.TaskVersionID = taskVersionID
	return e.validate()
}

func (e *Episode) SetStartedAt(startedAt time.Time) error {
	e.StartedAt = &startedAt
	return e.validate()
}

func (e *Episode) SetFinishedAt(finishedAt time.Time) error {
	e.FinishedAt = &finishedAt
	return e.validate()
}

func (e *Episode) SetStatus(status EpisodeStatus) error {
	e.Status = status
	return e.validate()
}

func (e *Episode) SetErrorDetails(errorDetails string) error {
	e.ErrorDetails = &errorDetails
	return e.validate()
}

func (e *Episode) SetRecordedByID(userID string) error {
	e.RecordedByID = &userID
	return e.validate()
}

// CanStart checks if the episode can be started
func (e *Episode) CanStart() error {
	if e.Status != EpisodeStatusReady {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "episode status must be Ready to start, current: %d", e.Status),
		)
	}
	return nil
}

// Start transitions the episode from Ready to Recording
func (e *Episode) Start(occurredAt time.Time) error {
	if err := e.CanStart(); err != nil {
		return err
	}
	e.Status = EpisodeStatusRecording
	e.StartedAt = &occurredAt
	return nil
}

// CanFinish checks if the episode can be finished
func (e *Episode) CanFinish() error {
	if e.Status != EpisodeStatusRecording {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "episode status must be Recording to finish, current: %d", e.Status),
		)
	}
	return nil
}

// Finish transitions the episode from Recording to Completed
func (e *Episode) Finish(occurredAt time.Time) error {
	if err := e.CanFinish(); err != nil {
		return err
	}
	e.Status = EpisodeStatusCompleted
	e.FinishedAt = &occurredAt
	return nil
}

// CanCancel checks if the episode can be cancelled
func (e *Episode) CanCancel() error {
	if e.Status != EpisodeStatusRecording {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeConflict, "episode status must be Recording to cancel, current: %d", e.Status),
		)
	}
	return nil
}

// Cancel transitions the episode from Recording to Cancelled
func (e *Episode) Cancel() error {
	if err := e.CanCancel(); err != nil {
		return err
	}
	e.Status = EpisodeStatusCancel
	return nil
}
