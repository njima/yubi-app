package persistence

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type userSite struct{}

func NewUserSite() *userSite { return &userSite{} }

func (us *userSite) SetUserSites(ctx context.Context, conn repository.DBConn, userID string, organizationID string, siteIDs []string) error {
	if _, err := conn.NewDelete().
		Model((*entity.UserSiteAssignment)(nil)).
		Where("user_id = ?", userID).
		Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete user site assignments: %v", err))
	}

	if len(siteIDs) == 0 {
		return nil
	}

	// Fetch organization_id for each site from the site table.
	type siteRow struct {
		IDNatural      string `bun:"id_natural"`
		OrganizationID string `bun:"organization_id"`
	}
	var sites []siteRow
	if err := conn.NewSelect().
		TableExpr("site").
		Column("id_natural", "organization_id").
		Where("id_natural IN (?)", bun.In(siteIDs)).
		Where("organization_id = ?", organizationID).
		Scan(ctx, &sites); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to fetch sites for assignment: %v", err))
	}

	orgBySite := make(map[string]string, len(sites))
	for _, s := range sites {
		orgBySite[s.IDNatural] = s.OrganizationID
	}

	rows := make([]entity.UserSiteAssignment, 0, len(siteIDs))
	for _, sid := range siteIDs {
		siteOrgID, ok := orgBySite[sid]
		if !ok {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "site not found: id=%s", sid))
		}
		if siteOrgID != organizationID {
			return apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "site %s does not belong to the user's organization", sid))
		}
		rows = append(rows, entity.UserSiteAssignment{
			UserID:         userID,
			SiteID:         sid,
			OrganizationID: siteOrgID,
		})
	}
	if _, err := conn.NewInsert().Model(&rows).Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to insert user site assignments: %v", err))
	}
	return nil
}
