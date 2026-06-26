package usecase

import (
	"context"
	"errors"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
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
	DisplayName    *string
	Description    *string
}

type organization struct {
	orgRepo repository.Organization
	data    repository.DataAccess
}

func NewOrganization(orgRepo repository.Organization, data repository.DataAccess) *organization {
	return &organization{
		orgRepo: orgRepo,
		data:    data,
	}
}

func (o *organization) Create(ctx context.Context, input OrganizationCreateInput) (model.Organization, error) {
	org, err := model.InitOrganization(input.DisplayName, input.Description)
	if err != nil {
		return model.Organization{}, err
	}

	uorg, err := o.orgRepo.Create(ctx, o.data.Conn(), org)
	if err != nil {
		return model.Organization{}, err
	}

	return uorg, nil
}

func (o *organization) GetByNaturalID(ctx context.Context, idNatural string) (model.Organization, error) {
	return o.orgRepo.GetByNaturalID(ctx, o.data.Conn(), idNatural)
}

func (o *organization) List(ctx context.Context, page, limit int) (model.Organizations, int, error) {
	pg := pagination.Normalize(page, limit)
	return o.orgRepo.List(ctx, o.data.Conn(), pg.Limit, pg.Offset)
}

func (o *organization) Update(ctx context.Context, input OrganizationUpdateInput) (model.Organization, error) {
	org, err := o.orgRepo.GetByNaturalID(ctx, o.data.Conn(), input.OrganizationID)
	if err != nil {
		return model.Organization{}, err
	}

	if input.DisplayName == nil && input.Description == nil {
		return org, nil
	}

	uorg, err := o.update(ctx, org, input)
	if err != nil {
		return model.Organization{}, err
	}

	updatedOrg, err := o.orgRepo.Update(ctx, o.data.Conn(), uorg)
	if err != nil {
		return model.Organization{}, err
	}

	return updatedOrg, nil
}

func (o *organization) update(ctx context.Context, org model.Organization, input OrganizationUpdateInput) (model.Organization, error) {
	if input.DisplayName != nil {
		if err := org.SetName(*input.DisplayName); err != nil {
			return model.Organization{}, err
		}
	}

	if input.Description != nil {
		if err := org.SetDescription(*input.Description); err != nil {
			return model.Organization{}, err
		}
	}

	return org, nil
}

func (o *organization) Delete(ctx context.Context, idNatural string) error {
	return o.orgRepo.Delete(ctx, o.data.Conn(), idNatural)
}
