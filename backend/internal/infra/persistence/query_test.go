package persistence

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type queryCaptureHook struct {
	lastQuery string
}

func (h *queryCaptureHook) BeforeQuery(ctx context.Context, _ *bun.QueryEvent) context.Context {
	return ctx
}

func (h *queryCaptureHook) AfterQuery(_ context.Context, event *bun.QueryEvent) {
	h.lastQuery = event.Query
}

func newPersistenceTestDB(t *testing.T, hook *queryCaptureHook) *bun.DB {
	t.Helper()
	dsn := "postgres://user:pass@localhost:5432/testdb?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(hook)
	return db
}

func captureQuery(t *testing.T, query func(db *bun.DB) *bun.SelectQuery, scanDest ...any) string {
	t.Helper()
	hook := &queryCaptureHook{}
	db := newPersistenceTestDB(t, hook)
	_ = query(db).Scan(context.Background(), scanDest...)
	return hook.lastQuery
}

func TestApplyLocationListFilters_UsedByListAndCountQueries(t *testing.T) {
	siteID := "site-1"
	search := "100%_ready"
	filter := repository.LocationListFilter{
		SiteID: &siteID,
		Search: &search,
	}

	listSQL := captureQuery(t, func(db *bun.DB) *bun.SelectQuery {
		var locations []entity.Location
		return applyLocationListFilters(db.NewSelect().Model(&locations).Relation("Site"), filter)
	})
	var total int
	countSQL := captureQuery(t, func(db *bun.DB) *bun.SelectQuery {
		return applyLocationListFilters(db.NewSelect().Model((*entity.Location)(nil)).ColumnExpr("COUNT(*)"), filter)
	}, &total)

	for _, sql := range []string{listSQL, countSQL} {
		if !strings.Contains(sql, "l.site_id = 'site-1'") {
			t.Fatalf("expected site filter in SQL, got:\n%s", sql)
		}
		if !strings.Contains(sql, "l.name ILIKE '%100\\%\\_ready%'") {
			t.Fatalf("expected escaped search filter in SQL, got:\n%s", sql)
		}
	}
}

func TestApplySiteListFilters_UsedByListAndCountQueries(t *testing.T) {
	organizationID := "org-1"
	search := "main_site%"
	filter := repository.SiteListFilter{
		OrganizationID: &organizationID,
		Search:         &search,
	}

	listSQL := captureQuery(t, func(db *bun.DB) *bun.SelectQuery {
		var sites []entity.Site
		return applySiteListFilters(db.NewSelect().Model(&sites), filter)
	})
	var total int
	countSQL := captureQuery(t, func(db *bun.DB) *bun.SelectQuery {
		return applySiteListFilters(db.NewSelect().Model((*entity.Site)(nil)).ColumnExpr("COUNT(*)"), filter)
	}, &total)

	for _, sql := range []string{listSQL, countSQL} {
		if !strings.Contains(sql, "organization_id = 'org-1'") {
			t.Fatalf("expected organization filter in SQL, got:\n%s", sql)
		}
		if !strings.Contains(sql, "name ILIKE '%main\\_site\\%%'") {
			t.Fatalf("expected escaped search filter in SQL, got:\n%s", sql)
		}
	}
}
