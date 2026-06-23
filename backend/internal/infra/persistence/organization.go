package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type organization struct{}

func NewOrganization() *organization { return &organization{} }

func (o *organization) Create(ctx context.Context, conn repository.DBConn, org model.Organization) (model.Organization, error) {
	var inserted entity.Organization

	dbOrg := entity.Organization{
		IDNatural:   org.IDNatural,
		Name:        org.Name,
		Description: org.Description,
	}

	if err := conn.NewInsert().
		Model(&dbOrg).
		Returning("*").
		Scan(ctx, &inserted); err != nil {
		return model.Organization{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create organization: %v", err))
	}

	result := model.Organization{
		ID:          inserted.ID,
		IDNatural:   inserted.IDNatural,
		Name:        inserted.Name,
		Description: inserted.Description,
		CreatedAt:   inserted.CreatedAt,
		UpdatedAt:   &inserted.UpdatedAt,
	}

	return result, nil
}

func (o *organization) GetByNaturalID(ctx context.Context, conn repository.DBConn, idNatural string) (model.Organization, error) {
	var dbOrg entity.Organization

	if err := conn.NewSelect().
		Model(&dbOrg).
		Where("id_natural = ?", idNatural).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.Organization{}, apperror.NewError(apperror.NewMessage(apperror.CodeOrganizationNotFound, "organization not found: id_natural=%s", idNatural))
		}
		return model.Organization{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get organization: %v", err))
	}

	domainOrg := model.Organization{
		ID:          dbOrg.ID,
		IDNatural:   dbOrg.IDNatural,
		Name:        dbOrg.Name,
		Description: dbOrg.Description,
		CreatedAt:   dbOrg.CreatedAt,
		UpdatedAt:   &dbOrg.UpdatedAt,
	}

	return domainOrg, nil
}

func (o *organization) List(ctx context.Context, conn repository.DBConn, limit, offset int) (model.Organizations, int, error) {
	var dbOrgs []entity.Organization

	sel := conn.NewSelect().
		Model(&dbOrgs).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)

	if err := sel.Scan(ctx); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list organizations: %v", err))
	}

	var total int
	if err := conn.NewSelect().
		Model((*entity.Organization)(nil)).
		ColumnExpr("COUNT(*)").
		Scan(ctx, &total); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count organizations: %v", err))
	}

	orgs := make(model.Organizations, 0, len(dbOrgs))
	for _, d := range dbOrgs {
		o := &model.Organization{
			ID:          d.ID,
			IDNatural:   d.IDNatural,
			Name:        d.Name,
			Description: d.Description,
			CreatedAt:   d.CreatedAt,
		}
		if !d.UpdatedAt.IsZero() {
			t := d.UpdatedAt
			o.UpdatedAt = &t
		}
		orgs = append(orgs, o)
	}

	return orgs, total, nil
}

func (o *organization) Update(ctx context.Context, conn repository.DBConn, org model.Organization) (model.Organization, error) {
	var updated entity.Organization

	upd := conn.NewUpdate().Model((*entity.Organization)(nil))
	hasSet := false
	if org.Name != "" {
		upd = upd.Set("name = ?", org.Name)
		hasSet = true
	}
	if org.Description != nil {
		upd = upd.Set("description = ?", *org.Description)
		hasSet = true
	}
	if !hasSet {
		return model.Organization{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "no fields to update"))
	}

	upd = upd.Set("updated_at = ?", time.Now().UTC())

	if err := upd.Where("id_natural = ?", org.IDNatural).Returning("*").Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.Organization{}, apperror.NewError(apperror.NewMessage(apperror.CodeOrganizationNotFound, "organization not found: id_natural=%s", org.IDNatural))
		}
		return model.Organization{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update organization: %v", err))
	}

	return model.Organization{
		ID:          updated.ID,
		IDNatural:   updated.IDNatural,
		Name:        updated.Name,
		Description: updated.Description,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   &updated.UpdatedAt,
	}, nil
}

func (o *organization) Delete(ctx context.Context, conn repository.DBConn, idNatural string) error {
	var deletedID int64
	if err := conn.NewDelete().
		Model((*entity.Organization)(nil)).
		Where("id_natural = ?", idNatural).
		Returning("id").
		Scan(ctx, &deletedID); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewError(apperror.NewMessage(apperror.CodeOrganizationNotFound, "organization not found: id_natural=%s", idNatural))
		}
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete organization: %v", err))
	}

	return nil
}
