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
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
	"github.com/uptrace/bun"
)

type user struct{}

func NewUser() *user {
	return &user{}
}

func toModelUser(dbu entity.User) model.User {
	return model.User{
		ID:        dbu.ID,
		IDNatural: dbu.IDNatural,
		GoogleSub: dbu.GoogleSub,
		Name:      dbu.Name,
		Email:     dbu.Email,
		AvatarURL: dbu.AvatarURL,
		CreatedAt: dbu.CreatedAt,
		UpdatedAt: updatedAtPtr(dbu.UpdatedAt),
	}
}

func (u *user) Create(ctx context.Context, conn repository.Conn, user model.User) (model.User, error) {
	var created entity.User

	dbu := entity.User{
		IDNatural: user.IDNatural,
		GoogleSub: user.GoogleSub,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}

	if err := bunConn(conn).NewInsert().
		Model(&dbu).
		Returning("*").
		Scan(ctx, &created); err != nil {
		return model.User{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create user: %v", err))
	}

	return toModelUser(created), nil
}

func (u *user) Update(ctx context.Context, conn repository.Conn, user model.User) (model.User, error) {
	upd := bunConn(conn).NewUpdate().Model((*entity.User)(nil))
	hasSet := false
	if user.GoogleSub != "" {
		upd = upd.Set("google_sub = ?", user.GoogleSub)
		hasSet = true
	}
	if user.Name != "" {
		upd = upd.Set("name = ?", user.Name)
		hasSet = true
	}
	if user.Email != "" {
		upd = upd.Set("email = ?", user.Email)
		hasSet = true
	}
	if user.AvatarURL != nil {
		upd = upd.Set("avatar_url = ?", *user.AvatarURL)
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

	return toModelUser(updated), nil
}

func (u *user) GetByNaturalID(ctx context.Context, conn repository.Conn, IDNatural string) (model.User, error) {
	var dbUser entity.User

	if err := bunConn(conn).NewSelect().
		Model(&dbUser).
		Where("u.id_natural = ?", IDNatural).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "user not found: id_natural=%s", IDNatural))
		}
		return model.User{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get user: %v", err))
	}

	return toModelUser(dbUser), nil
}

func (u *user) GetByGoogleSub(ctx context.Context, conn repository.Conn, googleSub string) (model.User, error) {
	var dbUser entity.User

	if err := bunConn(conn).NewSelect().
		Model(&dbUser).
		Where("u.google_sub = ?", googleSub).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, apperror.NewError(apperror.NewMessage(apperror.CodeUserNotFound, "user not found: google_sub=%s", googleSub))
		}
		return model.User{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get user by google_sub: %v", err))
	}

	return toModelUser(dbUser), nil
}

func (u *user) List(ctx context.Context, conn repository.Conn, filter repository.UserListFilter, limit, offset int) (model.Users, int, error) {
	var dbUsers []entity.User
	sel := bunConn(conn).NewSelect().
		Model(&dbUsers).
		Limit(limit).
		Offset(offset)

	// Dynamic ORDER BY with whitelist to prevent SQL injection
	sel = applyUserOrganizationMembershipScope(ctx, sel)
	sel = applyUserSortOrder(sel, filter.SortBy, filter.SortOrder)
	sel = applyUserListFilters(sel, filter)

	if err := sel.Scan(ctx); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list users: %v", err))
	}

	countQ := bunConn(conn).NewSelect().Model((*entity.User)(nil)).ColumnExpr("COUNT(*)")
	countQ = applyUserOrganizationMembershipScope(ctx, countQ)
	countQ = applyUserListFilters(countQ, filter)
	var total int
	if err := countQ.Scan(ctx, &total); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count users: %v", err))
	}

	users := make(model.Users, 0, len(dbUsers))
	for _, du := range dbUsers {
		usr := toModelUser(du)
		users = append(users, &usr)
	}

	return users, total, nil
}

var allowedUserSortColumns = map[string]string{
	"name":       "u.name",
	"email":      "u.email",
	"created_at": "u.created_at",
}

var nullableUserSortColumns = map[string]bool{}

func applyUserOrganizationMembershipScope(ctx context.Context, sel *bun.SelectQuery) *bun.SelectQuery {
	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return sel
	}

	return sel.
		Join("JOIN organization_membership AS om ON om.user_id = u.id_natural").
		Where("om.organization_id = ?", orgID)
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

func (u *user) ExistsByEmail(ctx context.Context, conn repository.Conn, email string) (bool, error) {
	exists, err := bunConn(conn).NewSelect().
		Model((*entity.User)(nil)).
		Where("email = ?", email).
		Exists(ctx)
	if err != nil {
		return false, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to check email existence: %v", err))
	}
	return exists, nil
}

func (u *user) ExistsByEmails(ctx context.Context, conn repository.Conn, emails []string) (map[string]bool, error) {
	if len(emails) == 0 {
		return map[string]bool{}, nil
	}
	var rows []struct {
		Email string `bun:"email"`
	}
	if err := bunConn(conn).NewSelect().
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

func (u *user) Delete(ctx context.Context, conn repository.Conn, idNatural string) error {
	var deletedID int64
	if err := bunConn(conn).NewDelete().
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
