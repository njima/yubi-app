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
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type location struct{}

func NewLocation() *location { return &location{} }

func (l *location) Create(ctx context.Context, conn repository.DBConn, loc model.Location) (model.Location, error) {
	var inserted entity.Location

	dbLoc := locationModelToEntity(loc)

	if err := conn.NewInsert().
		Model(&dbLoc).
		Returning("*").
		Scan(ctx, &inserted); err != nil {
		return model.Location{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create location: %v", err))
	}

	siteName, err := fetchSiteName(ctx, conn, inserted.SiteID)
	if err != nil {
		log.Warn().Err(err).Str("location_id", inserted.IDNatural).Str("site_id", inserted.SiteID).Msg("failed to fetch site name after location create; returning empty site_name")
		siteName = ""
	}

	inserted.Site = &entity.Site{Name: siteName}

	return locationEntityToModel(inserted), nil
}

func (l *location) GetByID(ctx context.Context, conn repository.DBConn, id string) (model.Location, error) {
	var dbLoc entity.Location
	if err := conn.NewSelect().
		Model(&dbLoc).
		Relation("Site").
		Where("l.id_natural = ?", id).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.Location{}, apperror.NewError(apperror.NewMessage(apperror.CodeLocationNotFound, "location not found: id_natural=%s", id))
		}
		return model.Location{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get location: %v", err))
	}

	return locationEntityToModel(dbLoc), nil
}

func (l *location) List(ctx context.Context, conn repository.DBConn, filter repository.LocationListFilter, limit, offset int) (model.Locations, int, error) {
	var dbLocs []entity.Location

	sel := conn.NewSelect().
		Model(&dbLocs).
		Relation("Site").
		Limit(limit).
		Offset(offset)

	sel = applyLocationListFilters(sel, filter)

	sel = applyLocationSortOrder(sel, filter.SortBy, filter.SortOrder)

	if err := sel.Scan(ctx); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list locations: %v", err))
	}

	var total int
	countSel := conn.NewSelect().
		Model((*entity.Location)(nil)).
		ColumnExpr("COUNT(*)")
	countSel = applyLocationListFilters(countSel, filter)
	if err := countSel.Scan(ctx, &total); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count locations: %v", err))
	}

	res := make(model.Locations, 0, len(dbLocs))
	for _, dl := range dbLocs {
		m := locationEntityToModel(dl)
		res = append(res, &m)
	}

	return res, total, nil
}

func (l *location) Update(ctx context.Context, conn repository.DBConn, loc model.Location) (model.Location, error) {
	var updated entity.Location

	upd := conn.NewUpdate().Model((*entity.Location)(nil))
	hasSet := false
	if loc.Name != "" {
		upd = upd.Set("name = ?", loc.Name)
		hasSet = true
	}
	if !hasSet {
		return model.Location{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "no fields to update"))
	}

	upd = upd.Set("updated_at = ?", time.Now().UTC())

	if err := upd.Where("id_natural = ?", loc.IDNatural).Returning("*").Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.Location{}, apperror.NewError(apperror.NewMessage(apperror.CodeLocationNotFound, "location not found: id_natural=%s", loc.IDNatural))
		}
		return model.Location{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update location: %v", err))
	}

	// Fetch site_name as best-effort: the update has already been committed.
	// See Create() for the same rationale.
	siteName, err := fetchSiteName(ctx, conn, updated.SiteID)
	if err != nil {
		log.Warn().Err(err).Str("location_id", updated.IDNatural).Str("site_id", updated.SiteID).Msg("failed to fetch site name after location update; returning empty site_name")
		siteName = ""
	}

	updated.Site = &entity.Site{Name: siteName}

	return locationEntityToModel(updated), nil
}

func (l *location) Delete(ctx context.Context, conn repository.DBConn, id string) error {
	var deletedID int64
	if err := conn.NewDelete().
		Model((*entity.Location)(nil)).
		Where("id_natural = ?", id).
		Returning("id").
		Scan(ctx, &deletedID); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewError(apperror.NewMessage(apperror.CodeLocationNotFound, "location not found: id_natural=%s", id))
		}
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete location: %v", err))
	}
	return nil
}

func fetchSiteName(ctx context.Context, conn repository.DBConn, siteID string) (string, error) {
	var site entity.Site
	if err := conn.NewSelect().
		Model(&site).
		Column("name").
		Where("id_natural = ?", siteID).
		Scan(ctx); err != nil {
		return "", apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to fetch site name: %v", err))
	}
	return site.Name, nil
}

var allowedLocationSortColumns = map[string]string{
	"name": "l.name",
}

func applyLocationSortOrder(sel *bun.SelectQuery, sortBy *repository.LocationSortBy, sortOrder *repository.SortOrder) *bun.SelectQuery {
	if sortBy == nil {
		return sel.OrderExpr("l.created_at DESC, l.id DESC")
	}

	col, ok := allowedLocationSortColumns[string(*sortBy)]
	if !ok {
		return sel.OrderExpr("l.created_at DESC, l.id DESC")
	}

	order := "ASC"
	if sortOrder != nil && *sortOrder == repository.SortOrderDesc {
		order = "DESC"
	}

	return sel.OrderExpr(fmt.Sprintf("%s %s, l.created_at DESC, l.id DESC", col, order))
}
