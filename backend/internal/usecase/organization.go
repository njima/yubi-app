package usecase

import (
	"context"
	"errors"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/pagination"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

var (
	ErrOrganizationNotFound      = errors.New("organization not found")
	ErrOrganizationInvalidInput  = errors.New("invalid input")
	ErrOrganizationAlreadyExists = errors.New("organization already exists")
)

type OrganizationUsecase interface {
	Create(ctx context.Context, input OrganizationCreateInput) (model.Organization, error)
	GetByNaturalID(ctx context.Context, idNatural string) (model.Organization, error)
	List(ctx context.Context, page, limit int) (model.Organizations, int, error)
	Update(ctx context.Context, input OrganizationUpdateInput) (model.Organization, error)
	Delete(ctx context.Context, idNatural string) error
}

type OrganizationCreateInput struct {
	DisplayName string
	Description *string
}

type OrganizationUpdateInput struct {
	OrganizationID string
	DisplayName    string
	Description    *string
}

type organization struct {
	orgRepo repository.Organization
	db      repository.DBConn
}

func NewOrganization(orgRepo repository.Organization, db repository.DBConn) *organization {
	return &organization{
		orgRepo: orgRepo,
		db:      db,
	}
}

func (o *organization) Create(ctx context.Context, input OrganizationCreateInput) (model.Organization, error) {
	org, err := model.InitOrganization(input.DisplayName, input.Description)
	if err != nil {
		return model.Organization{}, err
	}

	uorg, err := o.orgRepo.Create(ctx, o.db, org)
	if err != nil {
		return model.Organization{}, err
	}

	return uorg, nil
}

func (o *organization) GetByNaturalID(ctx context.Context, idNatural string) (model.Organization, error) {
	return o.orgRepo.GetByNaturalID(ctx, o.db, idNatural)
}

func (o *organization) List(ctx context.Context, page, limit int) (model.Organizations, int, error) {
	if limit <= 0 {
		limit = pagination.DefaultLimit
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return o.orgRepo.List(ctx, o.db, limit, offset)
}

func (o *organization) Update(ctx context.Context, input OrganizationUpdateInput) (model.Organization, error) {
	org, err := o.orgRepo.GetByNaturalID(ctx, o.db, input.OrganizationID)
	if err != nil {
		return model.Organization{}, err
	}

	uorg, err := o.update(ctx, org, input)
	if err != nil {
		return model.Organization{}, err
	}

	updatedOrg, err := o.orgRepo.Update(ctx, o.db, uorg)
	if err != nil {
		return model.Organization{}, err
	}

	return updatedOrg, nil
}

func (o *organization) update(ctx context.Context, org model.Organization, input OrganizationUpdateInput) (model.Organization, error) {
	if err := org.SetName(input.DisplayName); err != nil {
		return model.Organization{}, err
	}

	if input.Description != nil {
		if err := org.SetDescription(*input.Description); err != nil {
			return model.Organization{}, err
		}
	}

	return org, nil
}

func (o *organization) Delete(ctx context.Context, idNatural string) error {
	return o.orgRepo.Delete(ctx, o.db, idNatural)
}
