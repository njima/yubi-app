package model

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Organization struct {
	ID          int64
	IDNatural   string
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type Organizations []*Organization

func InitOrganization(name string, description *string) (Organization, error) {
	ID, err := InitID()
	if err != nil {
		return Organization{}, err
	}

	if description == nil {
		emptyStr := ""
		description = &emptyStr
	}

	org := Organization{
		IDNatural:   ID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if err := org.validate(); err != nil {
		return Organization{}, err
	}

	return org, nil
}

func NewOrganization(
	id int64,
	idNatural, name string,
	description *string,
	createdAt time.Time,
	updatedAt *time.Time,
) Organization {
	return Organization{
		ID:          id,
		IDNatural:   idNatural,
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (o Organization) validate() error {
	if err := (validation.Errors{
		"id_natural": validation.Validate(o.IDNatural, validation.Required.Error("id_natural is required")),
		"name": validation.Validate(
			o.Name,
			validation.Required.Error("name is required"),
			validation.RuneLength(1, 100).Error("name must be between 1 and 100 characters"),
		),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "organization validation failed: %v", err))
	}

	return nil
}

func (o *Organization) SetName(name string) error {
	o.Name = name
	return o.validate()
}

func (o *Organization) SetDescription(description string) error {
	o.Description = &description
	return o.validate()
}
