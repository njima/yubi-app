package model

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserRole int

const (
	UserRoleAdmin        UserRole = 0
	UserRoleDataEngineer UserRole = 1
	UserRoleManager      UserRole = 2
	UserRoleOperator     UserRole = 3
	UserRoleViewer       UserRole = 4
)

type LocationSummary struct {
	LocationID string
	Name       string
}

type SiteSummary struct {
	SiteID string
	Name   string
}

type User struct {
	ID               int64
	IDNatural        string
	OrganizationID   string
	OrganizationName string
	Name             string
	Email            string
	Role             UserRole
	CreatedAt        time.Time
	UpdatedAt        *time.Time
	Locations        []LocationSummary
	Sites            []SiteSummary
}

type Users []*User

func InitUser(
	organizationID,
	name,
	email string,
	role UserRole,
) (User, error) {
	ID, err := InitID()
	if err != nil {
		return User{}, err
	}

	user := User{
		IDNatural:      ID,
		OrganizationID: organizationID,
		Name:           name,
		Email:          email,
		Role:           role,
		CreatedAt:      time.Now(),
	}

	if err := user.validate(); err != nil {
		return User{}, err
	}

	return user, nil
}

func NewUser(
	ID int64,
	IDNatural,
	organizationID,
	name,
	email string,
	role UserRole,
	createdAt time.Time,
	updatedAt *time.Time,
) User {
	return User{
		ID:             ID,
		IDNatural:      IDNatural,
		OrganizationID: organizationID,
		Name:           name,
		Email:          email,
		Role:           role,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

func (u User) validate() error {
	if err := (validation.Errors{
		"id_natural":      validation.Validate(u.IDNatural, validation.Required.Error("id_natural is required")),
		"organization_id": validation.Validate(u.OrganizationID, validation.Required.Error("organization_id is required")),
		"name": validation.Validate(
			u.Name,
			validation.Required.Error("name is required"),
			validation.RuneLength(1, 60).Error("name must be between 1 and 60 characters"),
		),
		"email": validation.Validate(
			u.Email,
			validation.Required.Error("email is required"),
			is.EmailFormat.Error("invalid email format"),
		),
	}).Filter(); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeValidationError, "user validation failed: %v", err))
	}

	return nil
}

func (u *User) SetName(name string) error {
	u.Name = name
	return u.validate()
}

func (u *User) SetEmail(email string) error {
	u.Email = email
	return u.validate()
}

func (u *User) SetRole(role UserRole) error {
	u.Role = role
	return u.validate()
}
