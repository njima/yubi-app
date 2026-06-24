package model

import (
	"strings"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type TaskPriority int

const (
	TaskPriorityLow    TaskPriority = 0
	TaskPriorityNormal TaskPriority = 1
	TaskPriorityHigh   TaskPriority = 2
	TaskPriorityUrgent TaskPriority = 3
)

type TaskDifficulty int

const (
	TaskDifficultyS TaskDifficulty = 0
	TaskDifficultyA TaskDifficulty = 1
	TaskDifficultyB TaskDifficulty = 2
	TaskDifficultyC TaskDifficulty = 3
)

type TaskStatus int

const (
	TaskStatusPlanning  TaskStatus = 0
	TaskStatusDoing     TaskStatus = 1
	TaskStatusCompleted TaskStatus = 2
	TaskStatusCanceled  TaskStatus = 3
)

type Task struct {
	ID                    int64
	IDNatural             string
	OrganizationID        string
	Name                  string
	Description           *string
	ManualURL             string
	Priority              *TaskPriority
	Difficulty            *TaskDifficulty
	Status                *TaskStatus
	Deadline              time.Time
	RobotType             *string
	TargetDurationSeconds *int
	TargetEpisodeCount    *int
	ActualEpisodeCount    *int
	Version               string
	VersionDisplayName    *string
	IsActive              bool
	Tags                  TaskTags
	CreatedAt             time.Time
	UpdatedAt             *time.Time
}

type Tasks []*Task

const (
	InitialVersion = "v1.0.0"
)

func InitTask(organizationID string, name string, description *string, manualURL string, priority *TaskPriority, difficulty *TaskDifficulty, status *TaskStatus, deadline time.Time, robotType *string) (Task, error) {
	nID, err := InitID()
	if err != nil {
		return Task{}, err
	}

	if description == nil {
		emptyStr := ""
		description = &emptyStr
	}

	task := Task{
		IDNatural:      nID,
		OrganizationID: organizationID,
		Name:           name,
		Description:    description,
		ManualURL:      manualURL,
		Priority:       priority,
		Difficulty:     difficulty,
		Status:         status,
		Deadline:       deadline,
		RobotType:      robotType,
		Version:        InitialVersion,
		IsActive:       true,
		CreatedAt:      time.Now(),
	}

	if err := task.validate(); err != nil {
		return Task{}, err
	}

	return task, nil
}

func NewTask(
	id int64,
	idNatural, name string,
	description *string,
	deadline time.Time,
	createdAt time.Time,
	updatedAt *time.Time,
) Task {
	return Task{
		ID:          id,
		IDNatural:   idNatural,
		Name:        name,
		Description: description,
		Deadline:    deadline,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (t Task) validate() error {
	if err := (validation.Errors{
		"id_natural":      validation.Validate(t.IDNatural, validation.Required.Error("id_natural is required")),
		"organization_id": validation.Validate(t.OrganizationID, validation.Required.Error("organization_id is required")),
		"name": validation.Validate(
			t.Name,
			validation.Required.Error("name is required"),
			validation.RuneLength(1, 100).Error("name must be between 1 and 100 characters"),
		),
		"manual_url": validation.Validate(t.ManualURL,
			validation.Required.Error("manual_url is required"),
			validation.By(validateHTTPSURL),
		),
		"deadline": validation.Validate(t.Deadline,
			validation.By(validateDeadline),
		),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "task validation failed: %v", err))
	}

	if t.Priority == nil {
		return apperror.NewError(apperror.NewMessage(apperror.CodeValidationError, "priority is required"))
	}
	if err := validatePriority(*t.Priority); err != nil {
		return err
	}

	if t.Difficulty == nil {
		return apperror.NewError(apperror.NewMessage(apperror.CodeValidationError, "difficulty is required"))
	}
	if err := validateDifficulty(*t.Difficulty); err != nil {
		return err
	}

	if t.Status == nil {
		return apperror.NewError(apperror.NewMessage(apperror.CodeValidationError, "status is required"))
	}
	if err := validateTaskStatus(*t.Status); err != nil {
		return err
	}

	return nil
}

func validateDifficulty(difficulty TaskDifficulty) error {
	switch difficulty {
	case TaskDifficultyS, TaskDifficultyA, TaskDifficultyB, TaskDifficultyC:
		return nil
	default:
		return apperror.NewError(apperror.NewMessage(apperror.CodeValidationError, "difficulty must be 0 (S), 1 (A), 2 (B), or 3 (C)"))
	}
}

// DetermineTaskStatus determines the task status based on actual vs target duration.
// actualDuration: sum of completed episode durations for all approved versions (seconds)
// targetDuration: sum of target_duration_seconds for all approved versions (null treated as 0)
func DetermineTaskStatus(actualDuration, targetDuration int64) TaskStatus {
	if actualDuration == 0 {
		return TaskStatusPlanning
	}
	if actualDuration < targetDuration {
		return TaskStatusDoing
	}
	return TaskStatusCompleted
}

func validateTaskStatus(status TaskStatus) error {
	switch status {
	case TaskStatusPlanning, TaskStatusDoing, TaskStatusCompleted, TaskStatusCanceled:
		return nil
	default:
		return apperror.NewError(apperror.NewMessage(apperror.CodeValidationError, "status must be 0 (Planning), 1 (Doing), 2 (Completed), or 3 (Canceled)"))
	}
}

func validateDeadline(value any) error {
	t, ok := value.(time.Time)
	if !ok || t.IsZero() {
		return validation.NewError("validation_deadline_required", "deadline is required")
	}
	return nil
}

// validateHTTPSURL validates that a URL uses HTTPS.
// Empty string returns nil because all callers apply validation.Required first,
// which rejects empty values before this function is invoked.
func validateHTTPSURL(value any) error {
	s, _ := value.(string)
	if s == "" {
		return nil
	}
	if err := is.URL.Validate(s); err != nil {
		return err
	}
	if !strings.HasPrefix(s, "https://") {
		return validation.NewError("validation_https_required", "manual_url must use https")
	}
	return nil
}

func (t *Task) SetName(name string) error {
	if err := validateName(name); err != nil {
		return err
	}
	t.Name = name
	return nil
}

func (t *Task) SetStatus(status *TaskStatus) error {
	if status != nil {
		if err := validateTaskStatus(*status); err != nil {
			return err
		}
	}
	t.Status = status
	return nil
}

func (t *Task) SetDescription(description string) error {
	t.Description = &description
	return nil
}

func validateName(name string) error {
	if err := (validation.Errors{
		"name": validation.Validate(
			name,
			validation.Required.Error("name is required"),
			validation.RuneLength(1, 100).Error("name must be between 1 and 100 characters"),
		),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "name validation failed: %v", err))
	}
	return nil
}

func (t *Task) SetPriority(priority *TaskPriority) error {
	if priority != nil {
		if err := validatePriority(*priority); err != nil {
			return err
		}
	}
	t.Priority = priority
	return nil
}

func validatePriority(priority TaskPriority) error {
	switch priority {
	case TaskPriorityLow, TaskPriorityNormal, TaskPriorityHigh, TaskPriorityUrgent:
		return nil
	default:
		return apperror.NewError(apperror.NewMessage(apperror.CodeValidationError, "priority must be 0 (Low), 1 (Normal), 2 (High), or 3 (Urgent)"))
	}
}

func (t *Task) SetDifficulty(difficulty *TaskDifficulty) error {
	if difficulty != nil {
		if err := validateDifficulty(*difficulty); err != nil {
			return err
		}
	}
	t.Difficulty = difficulty
	return nil
}

func (t *Task) SetManualURL(url string) error {
	if err := validateManualURL(url); err != nil {
		return err
	}
	t.ManualURL = url
	return nil
}

func (t *Task) SetDeadline(deadline time.Time) error {
	if err := validateDeadline(deadline); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "deadline validation failed: %v", err))
	}
	t.Deadline = deadline
	return nil
}

func (t *Task) SetRobotType(robotType *string) {
	t.RobotType = robotType
}

func validateManualURL(url string) error {
	if err := (validation.Errors{
		"manual_url": validation.Validate(url,
			validation.Required.Error("manual_url is required"),
			validation.By(validateHTTPSURL),
		),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "manual_url validation failed: %v", err))
	}
	return nil
}

// TaskSummary holds aggregated task metrics.
type TaskSummary struct {
	TotalTasks            int
	TargetDurationSeconds int64
	TargetEpisodeCount    int
}

// TaskCompletionTrend holds trend data grouped by 2-week deadline periods.
type TaskCompletionTrend struct {
	Periods []TrendPeriod
}

// TrendPeriod holds trend data for a single 2-week period.
type TrendPeriod struct {
	Start  time.Time
	End    time.Time
	Groups []TrendGroup
}

// TrendGroup holds metrics for a single group (tag name or status) within a period.
type TrendGroup struct {
	Label          string
	TargetTasks    int
	ActualTasks    int
	TargetDuration int64
	ActualDuration int64
	TargetEpisodes int
	ActualEpisodes int
}
