package persistence

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
)

type userLocation struct{}

func NewUserLocation() *userLocation { return &userLocation{} }

func (ul *userLocation) SetUserLocations(ctx context.Context, conn repository.DBConn, userID string, organizationID string, locationIDs []string) error {
	if _, err := bunConn(conn).NewDelete().
		Model((*entity.UserLocationAssignment)(nil)).
		Where("user_id = ?", userID).
		Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete user location assignments: %v", err))
	}

	if len(locationIDs) == 0 {
		return nil
	}

	// Fetch organization_id for each location from the location table.
	type locationRow struct {
		IDNatural      string `bun:"id_natural"`
		OrganizationID string `bun:"organization_id"`
	}
	var locations []locationRow
	if err := bunConn(conn).NewSelect().
		TableExpr("location").
		Column("id_natural", "organization_id").
		Where("id_natural IN (?)", bun.In(locationIDs)).
		Where("organization_id = ?", organizationID).
		Scan(ctx, &locations); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to fetch locations for assignment: %v", err))
	}

	orgByLocation := make(map[string]string, len(locations))
	for _, l := range locations {
		orgByLocation[l.IDNatural] = l.OrganizationID
	}

	rows := make([]entity.UserLocationAssignment, 0, len(locationIDs))
	for _, lid := range locationIDs {
		locOrgID, ok := orgByLocation[lid]
		if !ok {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "location not found: id=%s", lid))
		}
		if locOrgID != organizationID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "location %s does not belong to the user's organization", lid))
		}
		rows = append(rows, entity.UserLocationAssignment{
			UserID:         userID,
			LocationID:     lid,
			OrganizationID: locOrgID,
		})
	}
	if _, err := bunConn(conn).NewInsert().Model(&rows).Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to insert user location assignments: %v", err))
	}
	return nil
}
