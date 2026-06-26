package entity_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// captureHook records the SQL of every query executed via bun (including failed ones).
// bun calls AfterQuery even when the DB connection fails, so we can capture the
// generated SQL without a real database.
type captureHook struct {
	lastQuery string
}

func (h *captureHook) BeforeQuery(ctx context.Context, _ *bun.QueryEvent) context.Context {
	return ctx
}

func (h *captureHook) AfterQuery(_ context.Context, event *bun.QueryEvent) {
	h.lastQuery = event.Query
}

// newTestDB creates a bun.DB backed by a non-existent Postgres instance.
// Queries will fail at the network level but the SQL is still built and
// captured by the QueryHook before the call goes out.
func newTestDB(t *testing.T) *bun.DB {
	t.Helper()
	dsn := "postgres://user:pass@localhost:5432/testdb?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	return bun.NewDB(sqldb, pgdialect.New())
}

const testOrgID = "org-abc"

func withOrg(t *testing.T) {
	t.Helper()
	entity.OrgIDFromContext = func(_ context.Context) (string, bool) {
		return testOrgID, true
	}
	t.Cleanup(func() { entity.OrgIDFromContext = nil })
}

// --- SELECT ---

func TestOrgScoped_Select_WithOrg(t *testing.T) {
	withOrg(t)
	hook := &captureHook{}
	db := newTestDB(t)
	db.AddQueryHook(hook)
	var robot entity.Robot
	_ = db.NewSelect().Model(&robot).Scan(context.Background())
	t.Logf("SELECT SQL:\n  %s", hook.lastQuery)
	// Must use table alias (r) to avoid ambiguity when JOINed tables also have organization_id.
	want := `r.organization_id = '` + testOrgID + `'`
	if !strings.Contains(hook.lastQuery, want) {
		t.Errorf("expected alias-qualified org filter %q in SELECT, got:\n  %s", want, hook.lastQuery)
	}
}

func TestOrgScoped_Select_AliasQualified(t *testing.T) {
	withOrg(t)
	hook := &captureHook{}
	db := newTestDB(t)
	db.AddQueryHook(hook)
	// Task entity uses alias "t" — verify each entity gets its own alias.
	var task entity.Task
	_ = db.NewSelect().Model(&task).Scan(context.Background())
	t.Logf("SELECT SQL (task):\n  %s", hook.lastQuery)
	want := `t.organization_id = '` + testOrgID + `'`
	if !strings.Contains(hook.lastQuery, want) {
		t.Errorf("expected alias-qualified org filter %q in SELECT, got:\n  %s", want, hook.lastQuery)
	}
}

func TestOrgScoped_Select_NoOrg(t *testing.T) {
	entity.OrgIDFromContext = func(_ context.Context) (string, bool) { return "", false }
	t.Cleanup(func() { entity.OrgIDFromContext = nil })
	hook := &captureHook{}
	db := newTestDB(t)
	db.AddQueryHook(hook)
	var robot entity.Robot
	_ = db.NewSelect().Model(&robot).Scan(context.Background())
	t.Logf("SELECT SQL (no org):\n  %s", hook.lastQuery)
	// "organization_id" appears as a column name in SELECT — check WHERE is absent.
	if strings.Contains(hook.lastQuery, "WHERE") {
		t.Errorf("expected no WHERE in SELECT when org not in context, got:\n  %s", hook.lastQuery)
	}
}

// --- UPDATE ---

func TestOrgScoped_Update_WithOrg(t *testing.T) {
	withOrg(t)
	hook := &captureHook{}
	db := newTestDB(t)
	db.AddQueryHook(hook)
	robot := &entity.Robot{Name: "test-bot"}
	_, _ = db.NewUpdate().Model(robot).Where("id_natural = ?", "robot-1").Exec(context.Background())
	t.Logf("UPDATE SQL:\n  %s", hook.lastQuery)
	if !strings.Contains(hook.lastQuery, testOrgID) {
		t.Errorf("expected org filter in UPDATE, got:\n  %s", hook.lastQuery)
	}
}

func TestOrgScoped_Update_NoOrg(t *testing.T) {
	entity.OrgIDFromContext = func(_ context.Context) (string, bool) { return "", false }
	t.Cleanup(func() { entity.OrgIDFromContext = nil })
	hook := &captureHook{}
	db := newTestDB(t)
	db.AddQueryHook(hook)
	robot := &entity.Robot{Name: "test-bot"}
	_, _ = db.NewUpdate().Model(robot).Where("id_natural = ?", "robot-1").Exec(context.Background())
	t.Logf("UPDATE SQL (no org):\n  %s", hook.lastQuery)
	if strings.Contains(hook.lastQuery, testOrgID) {
		t.Errorf("expected no org filter in UPDATE when org not in context, got:\n  %s", hook.lastQuery)
	}
}

// --- DELETE ---

func TestOrgScoped_Delete_WithOrg(t *testing.T) {
	withOrg(t)
	hook := &captureHook{}
	db := newTestDB(t)
	db.AddQueryHook(hook)
	_, _ = db.NewDelete().Model((*entity.Robot)(nil)).Where("id_natural = ?", "robot-1").Exec(context.Background())
	t.Logf("DELETE SQL:\n  %s", hook.lastQuery)
	if !strings.Contains(hook.lastQuery, testOrgID) {
		t.Errorf("expected org filter in DELETE, got:\n  %s", hook.lastQuery)
	}
}

func TestOrgScoped_Delete_NoOrg(t *testing.T) {
	entity.OrgIDFromContext = func(_ context.Context) (string, bool) { return "", false }
	t.Cleanup(func() { entity.OrgIDFromContext = nil })
	hook := &captureHook{}
	db := newTestDB(t)
	db.AddQueryHook(hook)
	_, _ = db.NewDelete().Model((*entity.Robot)(nil)).Where("id_natural = ?", "robot-1").Exec(context.Background())
	t.Logf("DELETE SQL (no org):\n  %s", hook.lastQuery)
	if strings.Contains(hook.lastQuery, testOrgID) {
		t.Errorf("expected no org filter in DELETE when org not in context, got:\n  %s", hook.lastQuery)
	}
}

// --- nil guard ---

func TestOrgScoped_NilGuard(t *testing.T) {
	entity.OrgIDFromContext = nil
	hook := &captureHook{}
	db := newTestDB(t)
	db.AddQueryHook(hook)
	// All three operations must not panic when OrgIDFromContext is nil.
	var robot entity.Robot
	_ = db.NewSelect().Model(&robot).Scan(context.Background())
	t.Logf("SELECT SQL (nil guard):\n  %s", hook.lastQuery)
	_, _ = db.NewUpdate().Model(&entity.Robot{Name: "x"}).Where("id_natural = ?", "r1").Exec(context.Background())
	t.Logf("UPDATE SQL (nil guard):\n  %s", hook.lastQuery)
	_, _ = db.NewDelete().Model((*entity.Robot)(nil)).Where("id_natural = ?", "r1").Exec(context.Background())
	t.Logf("DELETE SQL (nil guard):\n  %s", hook.lastQuery)
}
