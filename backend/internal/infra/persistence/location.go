package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type location struct{}

func NewLocation() *location { return &location{} }

func (l *location) Create(ctx context.Context, conn repository.DBConn, loc model.Location) (model.Location, error) {
	var inserted entity.Location

	dbLoc := entity.Location{
		IDNatural:      loc.IDNatural,
		OrganizationID: loc.OrganizationID,
		SiteID:         loc.SiteID,
		Name:           loc.Name,
	}

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

	result := model.Location{
		ID:             inserted.ID,
		IDNatural:      inserted.IDNatural,
		OrganizationID: inserted.OrganizationID,
		SiteID:         inserted.SiteID,
		SiteName:       siteName,
		Name:           inserted.Name,
		CreatedAt:      inserted.CreatedAt,
		UpdatedAt:      &inserted.UpdatedAt,
	}

	return result, nil
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

	loc := model.Location{
		ID:             dbLoc.ID,
		IDNatural:      dbLoc.IDNatural,
		OrganizationID: dbLoc.OrganizationID,
		SiteID:         dbLoc.SiteID,
		SiteName:       dbLoc.Site.Name,
		Name:           dbLoc.Name,
		CreatedAt:      dbLoc.CreatedAt,
		UpdatedAt:      &dbLoc.UpdatedAt,
	}

	return loc, nil
}

func (l *location) List(ctx context.Context, conn repository.DBConn, filter repository.LocationListFilter, limit, offset int) (model.Locations, int, error) {
	var dbLocs []entity.Location

	sel := conn.NewSelect().
		Model(&dbLocs).
		Relation("Site").
		Limit(limit).
		Offset(offset)

	if filter.SiteID != nil && *filter.SiteID != "" {
		sel = sel.Where("l.site_id = ?", *filter.SiteID)
	}

	sel = applyLocationSortOrder(sel, filter.SortBy, filter.SortOrder)

	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		sel = sel.Where("l.name ILIKE ?", "%"+escaped+"%")
	}

	if err := sel.Scan(ctx); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list locations: %v", err))
	}

	var total int
	countSel := conn.NewSelect().
		Model((*entity.Location)(nil)).
		ColumnExpr("COUNT(*)")
	if filter.SiteID != nil && *filter.SiteID != "" {
		countSel = countSel.Where("l.site_id = ?", *filter.SiteID)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		countSel = countSel.Where("l.name ILIKE ?", "%"+escaped+"%")
	}
	if err := countSel.Scan(ctx, &total); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count locations: %v", err))
	}

	res := make(model.Locations, 0, len(dbLocs))
	for _, dl := range dbLocs {
		m := &model.Location{
			ID:             dl.ID,
			IDNatural:      dl.IDNatural,
			OrganizationID: dl.OrganizationID,
			SiteID:         dl.SiteID,
			SiteName:       dl.Site.Name,
			Name:           dl.Name,
			CreatedAt:      dl.CreatedAt,
		}
		if !dl.UpdatedAt.IsZero() {
			t := dl.UpdatedAt
			m.UpdatedAt = &t
		}
		res = append(res, m)
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

	return model.Location{
		ID:             updated.ID,
		IDNatural:      updated.IDNatural,
		OrganizationID: updated.OrganizationID,
		SiteID:         updated.SiteID,
		SiteName:       siteName,
		Name:           updated.Name,
		CreatedAt:      updated.CreatedAt,
		UpdatedAt:      &updated.UpdatedAt,
	}, nil
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
