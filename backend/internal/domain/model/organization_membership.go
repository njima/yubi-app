package model

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type OrganizationMembership struct {
	ID             int64
	IDNatural      string
	UserID         string
	OrganizationID string
	Role           UserRole
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

func InitOrganizationMembership(userID, organizationID string, role UserRole) (OrganizationMembership, error) {
	ID, err := InitID()
	if err != nil {
		return OrganizationMembership{}, err
	}

	membership := OrganizationMembership{
		IDNatural:      ID,
		UserID:         userID,
		OrganizationID: organizationID,
		Role:           role,
		CreatedAt:      time.Now(),
	}
	if err := membership.validate(); err != nil {
		return OrganizationMembership{}, err
	}

	return membership, nil
}

func (m OrganizationMembership) validate() error {
	if err := (validation.Errors{
		"id_natural":      validation.Validate(m.IDNatural, validation.Required.Error("id_natural is required")),
		"user_id":         validation.Validate(m.UserID, validation.Required.Error("user_id is required")),
		"organization_id": validation.Validate(m.OrganizationID, validation.Required.Error("organization_id is required")),
		"role": validation.Validate(
			m.Role,
			validation.In(
				UserRoleAdmin,
				UserRoleDataEngineer,
				UserRoleManager,
				UserRoleOperator,
				UserRoleViewer,
			).Error("invalid user role"),
		),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "organization membership validation failed: %v", err))
	}

	return nil
}
