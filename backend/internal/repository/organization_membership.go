package repository

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

type OrganizationMembership interface {
	Create(ctx context.Context, conn Conn, membership model.OrganizationMembership) (model.OrganizationMembership, error)
	GetByUserAndOrganization(ctx context.Context, conn Conn, userID, organizationID string) (model.OrganizationMembership, error)
	ListByUser(ctx context.Context, conn Conn, userID string) ([]model.OrganizationMembership, error)
	CountByUser(ctx context.Context, conn Conn, userID string) (int, error)
}
