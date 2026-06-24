package model

import (
	"strings"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ApprovalStatus int

const (
	ApprovalStatusDraft    ApprovalStatus = 0
	ApprovalStatusApproved ApprovalStatus = 1
)

type TaskVersionParameter struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type TaskVersion struct {
	ID                              int64
	IDNatural                       string
	OrganizationID                  string
	TaskID                          string
	Version                         string
	DisplayName                     *string
	ApprovalStatus                  ApprovalStatus
	Parameters                      []TaskVersionParameter
	IsCurrent                       bool
	CreatedAt                       time.Time
	TargetDurationSeconds           *int
	TargetEpisodeCount              *int
	TargetDurationPerEpisodeSeconds *int
	ActualDurationSeconds           *int64
	ActualEpisodeCount              *int
}

func (tv *TaskVersion) Approve() {
	tv.ApprovalStatus = ApprovalStatusApproved
}

func (tv TaskVersion) IsDraft() bool {
	return tv.ApprovalStatus == ApprovalStatusDraft
}

func (tv TaskVersion) IsApproved() bool {
	return tv.ApprovalStatus == ApprovalStatusApproved
}

func (tv TaskVersion) DisplayLabel(taskName string) string {
	if tv.DisplayName != nil {
		return *tv.DisplayName
	}
	return taskName + " " + tv.Version
}

type TaskVersions []*TaskVersion

func InitTaskVersion(organizationID, taskID, version string, displayName *string, targetDurationSeconds, targetEpisodeCount, targetDurationPerEpisodeSeconds *int, parameters []TaskVersionParameter) (TaskVersion, error) {
	id, err := InitID()
	if err != nil {
		return TaskVersion{}, err
	}

	tv := TaskVersion{
		IDNatural:                       id,
		OrganizationID:                  organizationID,
		TaskID:                          taskID,
		Version:                         version,
		DisplayName:                     normalizeDisplayName(displayName),
		ApprovalStatus:                  ApprovalStatusDraft,
		Parameters:                      parameters,
		CreatedAt:                       time.Now(),
		TargetDurationSeconds:           targetDurationSeconds,
		TargetEpisodeCount:              targetEpisodeCount,
		TargetDurationPerEpisodeSeconds: targetDurationPerEpisodeSeconds,
	}

	if err := tv.validate(); err != nil {
		return TaskVersion{}, err
	}

	return tv, nil
}

// normalizeDisplayName trims surrounding whitespace and returns nil for empty
// inputs so the persistence layer stores NULL rather than an empty string.
func normalizeDisplayName(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func (tv TaskVersion) validate() error {
	if err := (validation.Errors{
		"id_natural":      validation.Validate(tv.IDNatural, validation.Required.Error("id_natural is required")),
		"organization_id": validation.Validate(tv.OrganizationID, validation.Required.Error("organization_id is required")),
		"task_id":         validation.Validate(tv.TaskID, validation.Required.Error("task_id is required")),
		"version": validation.Validate(
			tv.Version,
			validation.Required.Error("version is required"),
			validation.RuneLength(1, 50).Error("version must be between 1 and 50 characters"),
		),
		"display_name": validation.Validate(
			tv.DisplayName,
			validation.When(tv.DisplayName != nil,
				validation.Length(0, 100).Error("display_name must be at most 100 characters"),
			),
		),
		"approval_status": validation.Validate(
			tv.ApprovalStatus,
			validation.In(ApprovalStatusDraft, ApprovalStatusApproved).Error("approval_status must be 0 (draft) or 1 (approved)"),
		),
		"target_duration_seconds":             validation.Validate(tv.TargetDurationSeconds, validation.Min(1).Error("target_duration_seconds must be at least 1")),
		"target_episode_count":                validation.Validate(tv.TargetEpisodeCount, validation.Min(1).Error("target_episode_count must be at least 1")),
		"target_duration_per_episode_seconds": validation.Validate(tv.TargetDurationPerEpisodeSeconds, validation.Min(1).Error("target_duration_per_episode_seconds must be at least 1")),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "task version validation failed: %v", err))
	}

	if err := ValidateParameterDefinitions(tv.Parameters); err != nil {
		return err
	}

	return nil
}

func ValidateParameterDefinitions(params []TaskVersionParameter) error {
	if len(params) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(params))
	for _, p := range params {
		if p.Key == "" {
			return apperror.NewError(apperror.NewMessage(apperror.CodeValidationError, "parameter key must not be empty"))
		}
		if len(p.Values) == 0 {
			return apperror.NewError(apperror.NewMessage(apperror.CodeValidationError, "parameter %q must have at least one value", p.Key))
		}
		if seen[p.Key] {
			return apperror.NewError(apperror.NewMessage(apperror.CodeValidationError, "duplicate parameter key: %q", p.Key))
		}
		seen[p.Key] = true
	}
	return nil
}

// ValidateParameterValues validates that the given parameter values match the parameter definitions.
func ValidateParameterValues(params []TaskVersionParameter, values map[string]string) error {
	if len(params) == 0 {
		return nil
	}
	allowed := make(map[string]map[string]bool, len(params))
	for _, p := range params {
		valSet := make(map[string]bool, len(p.Values))
		for _, v := range p.Values {
			valSet[v] = true
		}
		allowed[p.Key] = valSet
	}
	for k, v := range values {
		valSet, ok := allowed[k]
		if !ok {
			return apperror.NewError(apperror.NewMessage(apperror.CodeValidationError, "unknown parameter key: %q", k))
		}
		if !valSet[v] {
			return apperror.NewError(apperror.NewMessage(apperror.CodeValidationError, "invalid value %q for parameter %q", v, k))
		}
	}
	return nil
}

// InterpolateSubTaskName replaces {key} placeholders in a subtask name with
// the corresponding parameter values.
func InterpolateSubTaskName(name string, parameterValues map[string]string) string {
	for k, v := range parameterValues {
		name = strings.ReplaceAll(name, "{"+k+"}", v)
	}
	return name
}
