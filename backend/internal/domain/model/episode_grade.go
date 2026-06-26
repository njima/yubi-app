package model

import (
	"math"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type EpisodeGrade struct {
	EpisodeID      string
	UserID         string
	OrganizationID string
	Grade          float64
	Comment        *string
	GradedAt       time.Time
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type EpisodeGrades []*EpisodeGrade

func InitEpisodeGrade(organizationID, episodeID, userID string, grade float64, comment *string) (EpisodeGrade, error) {
	now := time.Now()
	eg := EpisodeGrade{
		EpisodeID:      episodeID,
		UserID:         userID,
		OrganizationID: organizationID,
		Grade:          grade,
		Comment:        comment,
		GradedAt:       now,
		CreatedAt:      now,
	}

	if err := eg.validate(); err != nil {
		return EpisodeGrade{}, err
	}

	return eg, nil
}

func NewEpisodeGrade(
	organizationID,
	episodeID,
	userID string,
	grade float64,
	comment *string,
	gradedAt time.Time,
	createdAt time.Time,
	updatedAt *time.Time,
) EpisodeGrade {
	return EpisodeGrade{
		OrganizationID: organizationID,
		EpisodeID:      episodeID,
		UserID:         userID,
		Grade:          grade,
		Comment:        comment,
		GradedAt:       gradedAt,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

// MaxCommentLength caps the grade comment length to prevent storage/memory
// abuse via the unauthenticated comment field.
const MaxCommentLength = 10000

func (eg EpisodeGrade) validate() error {
	if math.IsNaN(eg.Grade) || math.IsInf(eg.Grade, 0) {
		return apperror.NewError(
			apperror.NewMessage(apperror.CodeValidationError, "grade must be a finite number"),
		)
	}

	commentRules := []validation.Rule{validation.Length(0, MaxCommentLength).Error("comment must be at most 10000 chars")}

	if err := (validation.Errors{
		"organization_id": validation.Validate(eg.OrganizationID, validation.Required.Error("organization_id is required")),
		"episode_id":      validation.Validate(eg.EpisodeID, validation.Required.Error("episode_id is required")),
		"user_id":         validation.Validate(eg.UserID, validation.Required.Error("user_id is required")),
		"grade":           validation.Validate(eg.Grade, validation.Min(0.0).Error("grade must be >= 0"), validation.Max(1.0).Error("grade must be <= 1")),
		"comment":         validation.Validate(eg.Comment, commentRules...),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "episode_grade validation failed: %v", err))
	}

	return nil
}

func (eg *EpisodeGrade) UpdateGrade(grade float64) error {
	eg.Grade = grade
	eg.GradedAt = time.Now()
	return eg.validate()
}

func (eg *EpisodeGrade) UpdateComment(comment *string) error {
	eg.Comment = comment
	eg.GradedAt = time.Now()
	return eg.validate()
}
