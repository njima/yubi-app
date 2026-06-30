package persistence

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
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
	return captureQueryWithContext(t, context.Background(), query, scanDest...)
}

func captureQueryWithContext(t *testing.T, ctx context.Context, query func(db *bun.DB) *bun.SelectQuery, scanDest ...any) string {
	t.Helper()
	hook := &queryCaptureHook{}
	db := newPersistenceTestDB(t, hook)
	_ = query(db).Scan(ctx, scanDest...)
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

func TestApplyRobotListFilters_UsedByListAndCountQueries(t *testing.T) {
	siteID := "site-1"
	locationID := "loc-1"
	status := repository.RobotFilterStatusReady
	robotType := "arm"
	search := "robot_100%"
	filter := repository.RobotListFilter{
		SiteID:     &siteID,
		LocationID: &locationID,
		Status:     &status,
		RobotType:  &robotType,
		Search:     &search,
	}

	listSQL := captureQuery(t, func(db *bun.DB) *bun.SelectQuery {
		var robots []entity.Robot
		return applyRobotListFilters(db.NewSelect().Model(&robots), filter)
	})
	var total int
	countSQL := captureQuery(t, func(db *bun.DB) *bun.SelectQuery {
		return applyRobotListFilters(db.NewSelect().Model((*entity.Robot)(nil)).ColumnExpr("COUNT(*)"), filter)
	}, &total)

	for _, sql := range []string{listSQL, countSQL} {
		if !strings.Contains(sql, "l.site_id = 'site-1'") {
			t.Fatalf("expected site EXISTS filter in SQL, got:\n%s", sql)
		}
		if !strings.Contains(sql, "r.location_id = 'loc-1'") {
			t.Fatalf("expected location filter in SQL, got:\n%s", sql)
		}
		if !strings.Contains(sql, "r.status = 5") {
			t.Fatalf("expected status filter in SQL, got:\n%s", sql)
		}
		if !strings.Contains(sql, "r.robot_type = 'arm'") {
			t.Fatalf("expected robot type filter in SQL, got:\n%s", sql)
		}
		if !strings.Contains(sql, "r.name ILIKE '%robot\\_100\\%%'") {
			t.Fatalf("expected escaped search filter in SQL, got:\n%s", sql)
		}
	}
}

func TestApplyRobotConnectionStateFilter_UsesReadyResolvableStatuses(t *testing.T) {
	onlineIDs := []string{"robot-1", "robot-2"}
	filter := repository.RobotListFilter{
		OnlineRobotIDs: &onlineIDs,
		ExcludeOnline:  true,
	}

	sql := captureQuery(t, func(db *bun.DB) *bun.SelectQuery {
		var robots []entity.Robot
		return applyRobotConnectionStateFilter(db.NewSelect().Model(&robots), filter.OnlineRobotIDs, filter.ExcludeOnline)
	})

	for _, want := range []string{
		"r.status IN (5, 0)",
		"r.id_natural NOT IN ('robot-1', 'robot-2')",
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("expected %q in SQL, got:\n%s", want, sql)
		}
	}
}

func TestApplyTaskListFilters_UsedByListAndCountQueries(t *testing.T) {
	hasApprovedVersion := true
	robotType := "arm"
	search := "task_100%"
	filter := repository.TaskListFilter{
		HasApprovedVersion: &hasApprovedVersion,
		Statuses:           []repository.TaskStatus{repository.TaskStatusDoing},
		Priorities:         []repository.TaskPriority{repository.TaskPriorityHigh},
		Difficulties:       []repository.TaskDifficulty{repository.TaskDifficultyA},
		RobotType:          &robotType,
		Search:             &search,
	}

	listSQL := captureQuery(t, func(db *bun.DB) *bun.SelectQuery {
		var tasks []entity.Task
		return applyTaskListFilters(db.NewSelect().Model(&tasks), filter)
	})
	var total int
	countSQL := captureQuery(t, func(db *bun.DB) *bun.SelectQuery {
		return applyTaskListFilters(db.NewSelect().Model((*entity.Task)(nil)).ColumnExpr("COUNT(*)"), filter)
	}, &total)

	for _, sql := range []string{listSQL, countSQL} {
		for _, want := range []string{
			"EXISTS (SELECT 1 FROM task_version tv WHERE tv.task_id = t.id_natural AND tv.approval_status = 1)",
			"t.status IN (1)",
			"t.priority IN (2)",
			"t.difficulty IN (1)",
			"t.robot_type = 'arm'",
			"t.name ILIKE '%task\\_100\\%%'",
		} {
			if !strings.Contains(sql, want) {
				t.Fatalf("expected %q in SQL, got:\n%s", want, sql)
			}
		}
	}
}

func TestApplyUserListFilters_UsedByListAndCountQueries(t *testing.T) {
	locationID := "loc-1"
	siteID := "site-1"
	search := "alice_100%"
	filter := repository.UserListFilter{
		LocationID: &locationID,
		SiteID:     &siteID,
		Search:     &search,
	}

	listSQL := captureQuery(t, func(db *bun.DB) *bun.SelectQuery {
		var users []entity.User
		return applyUserListFilters(db.NewSelect().Model(&users), filter)
	})
	var total int
	countSQL := captureQuery(t, func(db *bun.DB) *bun.SelectQuery {
		return applyUserListFilters(db.NewSelect().Model((*entity.User)(nil)).ColumnExpr("COUNT(*)"), filter)
	}, &total)

	for _, sql := range []string{listSQL, countSQL} {
		for _, want := range []string{
			"ula.user_id = u.id_natural AND ula.location_id = 'loc-1'",
			"usa.user_id = u.id_natural AND usa.site_id = 'site-1'",
			"u.name ILIKE '%alice\\_100\\%%'",
		} {
			if !strings.Contains(sql, want) {
				t.Fatalf("expected %q in SQL, got:\n%s", want, sql)
			}
		}
	}
}

func TestApplyUserOrganizationMembershipScope_UsedByListAndCountQueries(t *testing.T) {
	ctx := requestctx.SetOrganizationID(context.Background(), "org-1")

	listSQL := captureQueryWithContext(t, ctx, func(db *bun.DB) *bun.SelectQuery {
		var users []entity.User
		return applyUserOrganizationMembershipScope(ctx, db.NewSelect().Model(&users))
	})
	var total int
	countSQL := captureQueryWithContext(t, ctx, func(db *bun.DB) *bun.SelectQuery {
		return applyUserOrganizationMembershipScope(ctx, db.NewSelect().Model((*entity.User)(nil)).ColumnExpr("COUNT(*)"))
	}, &total)

	for _, sql := range []string{listSQL, countSQL} {
		for _, want := range []string{
			"JOIN organization_membership AS om ON om.user_id = u.id_natural",
			"om.organization_id = 'org-1'",
		} {
			if !strings.Contains(sql, want) {
				t.Fatalf("expected %q in SQL, got:\n%s", want, sql)
			}
		}
	}
}

func TestApplySubTaskListFilters_UsesTaskVersionBeforeTask(t *testing.T) {
	taskID := "task-1"
	taskVersionID := "version-1"
	filter := repository.SubTaskListFilter{
		TaskID:        &taskID,
		TaskVersionID: &taskVersionID,
	}

	sql := captureQuery(t, func(db *bun.DB) *bun.SelectQuery {
		var subtasks []entity.SubTask
		return applySubTaskListFilters(db.NewSelect().Model(&subtasks), filter)
	})

	if !strings.Contains(sql, "task_version_id = 'version-1'") {
		t.Fatalf("expected task version filter in SQL, got:\n%s", sql)
	}
	if strings.Contains(sql, "task_id = 'task-1'") {
		t.Fatalf("expected task_id filter to be ignored when task_version_id is present, got:\n%s", sql)
	}
}
