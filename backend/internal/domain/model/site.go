package model

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Site struct {
	ID             int64
	IDNatural      string
	OrganizationID string
	Name           string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type Sites []*Site

func InitSite(organizationID, name string) (Site, error) {
	ID, err := InitID()
	if err != nil {
		return Site{}, err
	}

	site := Site{
		IDNatural:      ID,
		OrganizationID: organizationID,
		Name:           name,
		CreatedAt:      time.Now(),
	}

	if err := site.validate(); err != nil {
		return Site{}, err
	}

	return site, nil
}

func NewSite(
	id int64,
	idNatural,
	organizationID,
	name string,
	createdAt time.Time,
	updatedAt *time.Time,
) Site {
	return Site{
		ID:             id,
		IDNatural:      idNatural,
		OrganizationID: organizationID,
		Name:           name,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

func (s Site) validate() error {
	if err := (validation.Errors{
		"id_natural":      validation.Validate(s.IDNatural, validation.Required.Error("id_natural is required")),
		"organization_id": validation.Validate(s.OrganizationID, validation.Required.Error("organization_id is required")),
		"name": validation.Validate(
			s.Name,
			validation.Required.Error("name is required"),
			validation.RuneLength(1, 100).Error("name must be between 1 and 100 characters"),
		),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "site validation failed: %v", err))
	}

	return nil
}

func (s *Site) SetName(name string) error {
	s.Name = name
	return s.validate()
}
