package model

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type OrganizationKind string

const (
	OrganizationKindPersonal OrganizationKind = "personal"
	OrganizationKindTeam     OrganizationKind = "team"
)

type Organization struct {
	ID          int64
	IDNatural   string
	Name        string
	Kind        OrganizationKind
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type Organizations []*Organization

func InitOrganization(name string, description *string, kind OrganizationKind) (Organization, error) {
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
		Kind:        kind,
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
	kind OrganizationKind,
	description *string,
	createdAt time.Time,
	updatedAt *time.Time,
) Organization {
	return Organization{
		ID:          id,
		IDNatural:   idNatural,
		Name:        name,
		Kind:        kind,
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
		"kind": validation.Validate(
			string(o.Kind),
			validation.Required.Error("kind is required"),
			validation.In(string(OrganizationKindPersonal), string(OrganizationKindTeam)).Error("invalid organization kind"),
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
