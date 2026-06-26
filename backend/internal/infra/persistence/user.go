package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/uptrace/bun"
)

type user struct{}

func NewUser() *user {
	return &user{}
}

func (u *user) Create(ctx context.Context, conn repository.DBConn, user model.User) (model.User, error) {
	var created entity.User

	dbu := entity.User{
		IDNatural:      user.IDNatural,
		OrganizationID: user.OrganizationID,
		Name:           user.Name,
		Email:          user.Email,
		Role:           uint(user.Role),
	}

	if err := conn.NewInsert().
		Model(&dbu).
		Returning("*").
		Scan(ctx, &created); err != nil {
		return model.User{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create user: %v", err))
	}

	return model.User{
		ID:             created.ID,
		IDNatural:      created.IDNatural,
		OrganizationID: created.OrganizationID,
		Name:           created.Name,
		Email:          created.Email,
		Role:           model.UserRole(created.Role),
		CreatedAt:      created.CreatedAt,
		UpdatedAt:      &created.UpdatedAt,
	}, nil
}

func (u *user) Update(ctx context.Context, conn repository.DBConn, user model.User) (model.User, error) {
	upd := conn.NewUpdate().Model((*entity.User)(nil))
	hasSet := false
	if user.Name != "" {
		upd = upd.Set("name = ?", user.Name)
		hasSet = true
	}
	if user.Email != "" {
		upd = upd.Set("email = ?", user.Email)
		hasSet = true
	}
	if !hasSet {
		return model.User{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "no fields to update"))
	}

	upd = upd.Set("updated_at = ?", time.Now().UTC())

	var updated entity.User
	if err := upd.Where("id_natural = ?", user.IDNatural).Returning("*").Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "user not found: id_natural=%s", user.IDNatural))
		}
		return model.User{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update user: %v", err))
	}

	return model.User{
		ID:             updated.ID,
		IDNatural:      updated.IDNatural,
		OrganizationID: updated.OrganizationID,
		Name:           updated.Name,
		Email:          updated.Email,
		Role:           model.UserRole(updated.Role),
		CreatedAt:      updated.CreatedAt,
		UpdatedAt:      &updated.UpdatedAt,
	}, nil
}

func (u *user) UpdateRole(ctx context.Context, conn repository.DBConn, idNatural string, role model.UserRole) (model.User, error) {
	var updated entity.User
	if err := conn.NewUpdate().
		Model((*entity.User)(nil)).
		Set("role = ?", uint(role)).
		Set("updated_at = ?", time.Now().UTC()).
		Where("id_natural = ?", idNatural).
		Returning("*").
		Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "user not found: id_natural=%s", idNatural))
		}
		return model.User{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update user role: %v", err))
	}

	return model.User{
		ID:             updated.ID,
		IDNatural:      updated.IDNatural,
		OrganizationID: updated.OrganizationID,
		Name:           updated.Name,
		Email:          updated.Email,
		Role:           model.UserRole(updated.Role),
		CreatedAt:      updated.CreatedAt,
		UpdatedAt:      &updated.UpdatedAt,
	}, nil
}

func (u *user) GetByNaturalID(ctx context.Context, conn repository.DBConn, IDNatural string) (model.User, error) {
	var dbUser entity.User

	if err := conn.NewSelect().
		Model(&dbUser).
		Relation("Organization").
		Relation("LocationAssignments.Location").
		Relation("SiteAssignments.Site").
		Where("u.id_natural = ?", IDNatural).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "user not found: id_natural=%s", IDNatural))
		}
		return model.User{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get user: %v", err))
	}

	orgName := ""
	if dbUser.Organization != nil {
		orgName = dbUser.Organization.Name
	}

	locs := make([]model.LocationSummary, 0, len(dbUser.LocationAssignments))
	for _, la := range dbUser.LocationAssignments {
		if la.Location != nil {
			locs = append(locs, model.LocationSummary{
				LocationID: la.LocationID,
				Name:       la.Location.Name,
			})
		}
	}

	sites := make([]model.SiteSummary, 0, len(dbUser.SiteAssignments))
	for _, sa := range dbUser.SiteAssignments {
		if sa.Site != nil {
			sites = append(sites, model.SiteSummary{
				SiteID: sa.SiteID,
				Name:   sa.Site.Name,
			})
		}
	}

	return model.User{
		ID:               dbUser.ID,
		IDNatural:        dbUser.IDNatural,
		OrganizationID:   dbUser.OrganizationID,
		OrganizationName: orgName,
		Name:             dbUser.Name,
		Email:            dbUser.Email,
		Role:             model.UserRole(dbUser.Role),
		CreatedAt:        dbUser.CreatedAt,
		UpdatedAt:        &dbUser.UpdatedAt,
		Locations:        locs,
		Sites:            sites,
	}, nil
}

func (u *user) List(ctx context.Context, conn repository.DBConn, filter repository.UserListFilter, limit, offset int) (model.Users, int, error) {
	var dbUsers []entity.User
	sel := conn.NewSelect().
		Model(&dbUsers).
		Relation("Organization").
		Relation("LocationAssignments.Location").
		Relation("SiteAssignments.Site").
		Limit(limit).
		Offset(offset)

	// Dynamic ORDER BY with whitelist to prevent SQL injection
	sel = applyUserSortOrder(sel, filter.SortBy, filter.SortOrder)
	if filter.LocationID != nil && *filter.LocationID != "" {
		sel = sel.Where(`EXISTS (
			SELECT 1 FROM user_location_assignment ula
			WHERE ula.user_id = u.id_natural AND ula.location_id = ?
		)`, *filter.LocationID)
	}
	if filter.SiteID != nil && *filter.SiteID != "" {
		sel = sel.Where(`EXISTS (
			SELECT 1 FROM user_site_assignment usa
			WHERE usa.user_id = u.id_natural AND usa.site_id = ?
		)`, *filter.SiteID)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		sel = sel.Where("u.name ILIKE ?", "%"+escaped+"%")
	}

	if err := sel.Scan(ctx); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list users: %v", err))
	}

	countQ := conn.NewSelect().Model((*entity.User)(nil)).ColumnExpr("COUNT(*)")
	if filter.LocationID != nil && *filter.LocationID != "" {
		countQ = countQ.Where(`EXISTS (
			SELECT 1 FROM user_location_assignment ula
			WHERE ula.user_id = u.id_natural AND ula.location_id = ?
		)`, *filter.LocationID)
	}
	if filter.SiteID != nil && *filter.SiteID != "" {
		countQ = countQ.Where(`EXISTS (
			SELECT 1 FROM user_site_assignment usa
			WHERE usa.user_id = u.id_natural AND usa.site_id = ?
		)`, *filter.SiteID)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		countQ = countQ.Where("u.name ILIKE ?", "%"+escaped+"%")
	}
	var total int
	if err := countQ.Scan(ctx, &total); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count users: %v", err))
	}

	users := make(model.Users, 0, len(dbUsers))
	for _, du := range dbUsers {
		orgName := ""
		if du.Organization != nil {
			orgName = du.Organization.Name
		}
		locs := make([]model.LocationSummary, 0, len(du.LocationAssignments))
		for _, la := range du.LocationAssignments {
			if la.Location != nil {
				locs = append(locs, model.LocationSummary{
					LocationID: la.LocationID,
					Name:       la.Location.Name,
				})
			}
		}
		sites := make([]model.SiteSummary, 0, len(du.SiteAssignments))
		for _, sa := range du.SiteAssignments {
			if sa.Site != nil {
				sites = append(sites, model.SiteSummary{
					SiteID: sa.SiteID,
					Name:   sa.Site.Name,
				})
			}
		}
		usr := &model.User{
			ID:               du.ID,
			IDNatural:        du.IDNatural,
			OrganizationID:   du.OrganizationID,
			OrganizationName: orgName,
			Name:             du.Name,
			Email:            du.Email,
			Role:             model.UserRole(du.Role),
			CreatedAt:        du.CreatedAt,
			Locations:        locs,
			Sites:            sites,
		}
		if !du.UpdatedAt.IsZero() {
			t := du.UpdatedAt
			usr.UpdatedAt = &t
		}
		users = append(users, usr)
	}

	return users, total, nil
}

var allowedUserSortColumns = map[string]string{
	"name":       "u.name",
	"email":      "u.email",
	"role":       "u.role",
	"location":   "(SELECT l.name FROM user_location_assignment ula JOIN location l ON l.id_natural = ula.location_id WHERE ula.user_id = u.id_natural ORDER BY l.name ASC LIMIT 1)",
	"created_at": "u.created_at",
}

var nullableUserSortColumns = map[string]bool{
	"location": true,
}

func applyUserSortOrder(sel *bun.SelectQuery, sortBy *repository.UserSortBy, sortOrder *repository.SortOrder) *bun.SelectQuery {
	if sortBy == nil {
		return sel.OrderExpr("u.created_at DESC")
	}

	col, ok := allowedUserSortColumns[string(*sortBy)]
	if !ok {
		return sel.OrderExpr("u.created_at DESC")
	}

	order := "ASC"
	if sortOrder != nil && *sortOrder == repository.SortOrderDesc {
		order = "DESC"
	}

	nullsClause := ""
	if nullableUserSortColumns[string(*sortBy)] {
		nullsClause = " NULLS LAST"
	}

	return sel.OrderExpr(fmt.Sprintf("%s %s%s", col, order, nullsClause))
}

func (u *user) ExistsByEmail(ctx context.Context, conn repository.DBConn, email string) (bool, error) {
	exists, err := conn.NewSelect().
		Model((*entity.User)(nil)).
		Where("email = ?", email).
		Exists(ctx)
	if err != nil {
		return false, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to check email existence: %v", err))
	}
	return exists, nil
}

func (u *user) ExistsByEmails(ctx context.Context, conn repository.DBConn, emails []string) (map[string]bool, error) {
	if len(emails) == 0 {
		return map[string]bool{}, nil
	}
	var rows []struct {
		Email string `bun:"email"`
	}
	if err := conn.NewSelect().
		TableExpr(`"user" AS u`).
		ColumnExpr("lower(u.email) AS email").
		Where("lower(u.email) IN (?)", bun.In(emails)).
		Scan(ctx, &rows); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to check emails existence: %v", err))
	}
	result := make(map[string]bool, len(rows))
	for _, r := range rows {
		result[r.Email] = true
	}
	return result, nil
}

func (u *user) Delete(ctx context.Context, conn repository.DBConn, idNatural string) error {
	var deletedID int64
	if err := conn.NewDelete().
		Model((*entity.User)(nil)).
		Where("id_natural = ?", idNatural).
		Returning("id").
		Scan(ctx, &deletedID); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "user not found: id_natural=%s", idNatural))
		}
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete user: %v", err))
	}

	return nil
}
