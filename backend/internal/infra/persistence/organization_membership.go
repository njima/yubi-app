package persistence

import (
	"context"
	"database/sql"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
)

type organizationMembership struct{}

func NewOrganizationMembership() *organizationMembership {
	return &organizationMembership{}
}

func toModelOrganizationMembership(dbm entity.OrganizationMembership) model.OrganizationMembership {
	return model.OrganizationMembership{
		ID:             dbm.ID,
		IDNatural:      dbm.IDNatural,
		UserID:         dbm.UserID,
		OrganizationID: dbm.OrganizationID,
		Role:           model.UserRole(dbm.Role),
		CreatedAt:      dbm.CreatedAt,
		UpdatedAt:      updatedAtPtr(dbm.UpdatedAt),
	}
}

func (o *organizationMembership) Create(ctx context.Context, conn repository.Conn, membership model.OrganizationMembership) (model.OrganizationMembership, error) {
	var created entity.OrganizationMembership

	dbm := entity.OrganizationMembership{
		IDNatural:      membership.IDNatural,
		UserID:         membership.UserID,
		OrganizationID: membership.OrganizationID,
		Role:           uint(membership.Role),
	}

	if err := bunConn(conn).NewInsert().
		Model(&dbm).
		Returning("*").
		Scan(ctx, &created); err != nil {
		return model.OrganizationMembership{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create organization membership: %v", err))
	}

	return toModelOrganizationMembership(created), nil
}

func (o *organizationMembership) GetByUserAndOrganization(ctx context.Context, conn repository.Conn, userID, organizationID string) (model.OrganizationMembership, error) {
	var dbm entity.OrganizationMembership

	if err := bunConn(conn).NewSelect().
		Model(&dbm).
		Where("om.user_id = ?", userID).
		Where("om.organization_id = ?", organizationID).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.OrganizationMembership{}, apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "organization membership not found: user_id=%s organization_id=%s", userID, organizationID))
		}
		return model.OrganizationMembership{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get organization membership: %v", err))
	}

	return toModelOrganizationMembership(dbm), nil
}

func (o *organizationMembership) ListByUser(ctx context.Context, conn repository.Conn, userID string) ([]model.OrganizationMembership, error) {
	var dbMemberships []entity.OrganizationMembership

	if err := bunConn(conn).NewSelect().
		Model(&dbMemberships).
		Where("om.user_id = ?", userID).
		Order("om.created_at DESC").
		Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list organization memberships: %v", err))
	}

	memberships := make([]model.OrganizationMembership, 0, len(dbMemberships))
	for _, dbm := range dbMemberships {
		memberships = append(memberships, toModelOrganizationMembership(dbm))
	}

	return memberships, nil
}

func (o *organizationMembership) CountByUser(ctx context.Context, conn repository.Conn, userID string) (int, error) {
	count, err := bunConn(conn).NewSelect().
		Model((*entity.OrganizationMembership)(nil)).
		Where("user_id = ?", userID).
		Count(ctx)
	if err != nil {
		return 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count organization memberships: %v", err))
	}

	return count, nil
}
