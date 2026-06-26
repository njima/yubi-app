package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
)

type site struct{}

func NewSite() *site { return &site{} }

func (s *site) Create(ctx context.Context, conn repository.Conn, si model.Site) (model.Site, error) {
	var inserted entity.Site

	dbSite := siteModelToEntity(si)

	if err := bunConn(conn).NewInsert().
		Model(&dbSite).
		Returning("*").
		Scan(ctx, &inserted); err != nil {
		return model.Site{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create site: %v", err))
	}

	return siteEntityToModel(inserted), nil
}

func (s *site) GetByID(ctx context.Context, conn repository.Conn, id string) (model.Site, error) {
	var dbSite entity.Site
	if err := bunConn(conn).NewSelect().
		Model(&dbSite).
		Where("id_natural = ?", id).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.Site{}, apperror.NewError(apperror.NewMessage(apperror.CodeSiteNotFound, "site not found: id_natural=%s", id))
		}
		return model.Site{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get site: %v", err))
	}

	return siteEntityToModel(dbSite), nil
}

func (s *site) List(ctx context.Context, conn repository.Conn, filter repository.SiteListFilter, limit, offset int) (model.Sites, int, error) {
	var dbSites []entity.Site

	sel := bunConn(conn).NewSelect().
		Model(&dbSites).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)

	sel = applySiteListFilters(sel, filter)

	if err := sel.Scan(ctx); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list sites: %v", err))
	}

	var total int
	countSel := bunConn(conn).NewSelect().
		Model((*entity.Site)(nil)).
		ColumnExpr("COUNT(*)")
	countSel = applySiteListFilters(countSel, filter)
	if err := countSel.Scan(ctx, &total); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count sites: %v", err))
	}

	res := make(model.Sites, 0, len(dbSites))
	for _, ds := range dbSites {
		m := siteEntityToModel(ds)
		res = append(res, &m)
	}

	return res, total, nil
}

func (s *site) Update(ctx context.Context, conn repository.Conn, si model.Site) (model.Site, error) {
	var updated entity.Site

	upd := bunConn(conn).NewUpdate().Model((*entity.Site)(nil))
	hasSet := false
	if si.Name != "" {
		upd = upd.Set("name = ?", si.Name)
		hasSet = true
	}
	if !hasSet {
		return model.Site{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "no fields to update"))
	}

	upd = upd.Set("updated_at = ?", time.Now().UTC())

	if err := upd.Where("id_natural = ?", si.IDNatural).Returning("*").Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.Site{}, apperror.NewError(apperror.NewMessage(apperror.CodeSiteNotFound, "site not found: id_natural=%s", si.IDNatural))
		}
		return model.Site{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update site: %v", err))
	}

	return siteEntityToModel(updated), nil
}

func (s *site) Delete(ctx context.Context, conn repository.Conn, id string) error {
	var deletedID int64
	if err := bunConn(conn).NewDelete().
		Model((*entity.Site)(nil)).
		Where("id_natural = ?", id).
		Returning("id").
		Scan(ctx, &deletedID); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewError(apperror.NewMessage(apperror.CodeSiteNotFound, "site not found: id_natural=%s", id))
		}
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete site: %v", err))
	}
	return nil
}
