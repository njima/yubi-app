package model

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Location struct {
	ID             int64
	IDNatural      string
	OrganizationID string
	SiteID         string
	SiteName       string
	Name           string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type Locations []*Location

func InitLocation(organizationID, siteID, name string) (Location, error) {
	ID, err := InitID()
	if err != nil {
		return Location{}, err
	}

	loc := Location{
		IDNatural:      ID,
		OrganizationID: organizationID,
		SiteID:         siteID,
		Name:           name,
		CreatedAt:      time.Now(),
	}

	if err := loc.validate(); err != nil {
		return Location{}, err
	}

	return loc, nil
}

func NewLocation(
	ID int64,
	IDNatural,
	organizationID,
	siteID,
	name string,
	createdAt time.Time,
	updatedAt *time.Time,
) Location {
	return Location{
		ID:             ID,
		IDNatural:      IDNatural,
		OrganizationID: organizationID,
		SiteID:         siteID,
		Name:           name,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

func (l Location) validate() error {
	if err := (validation.Errors{
		"id_natural":      validation.Validate(l.IDNatural, validation.Required.Error("id_natural is required")),
		"organization_id": validation.Validate(l.OrganizationID, validation.Required.Error("organization_id is required")),
		"site_id":         validation.Validate(l.SiteID, validation.Required.Error("site_id is required")),
		"name": validation.Validate(
			l.Name,
			validation.Required.Error("name is required"),
			validation.RuneLength(1, 100).Error("name must be between 1 and 100 characters"),
		),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "location validation failed: %v", err))
	}

	return nil
}

func (l *Location) SetName(name string) error {
	l.Name = name
	return l.validate()
}
