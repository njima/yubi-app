package persistence

import (
	"strings"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

// escapeILIKE escapes ILIKE wildcards (% and _) so user-supplied search terms
// don't act as pattern matchers.
func escapeILIKE(s string) string {
	return strings.NewReplacer("%", "\\%", "_", "\\_").Replace(s)
}

func applyLocationListFilters(sel *bun.SelectQuery, filter repository.LocationListFilter) *bun.SelectQuery {
	if filter.SiteID != nil && *filter.SiteID != "" {
		sel = sel.Where("l.site_id = ?", *filter.SiteID)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		sel = sel.Where("l.name ILIKE ?", "%"+escaped+"%")
	}
	return sel
}

func applySiteListFilters(sel *bun.SelectQuery, filter repository.SiteListFilter) *bun.SelectQuery {
	if filter.OrganizationID != nil {
		sel = sel.Where("organization_id = ?", *filter.OrganizationID)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		sel = sel.Where("name ILIKE ?", "%"+escaped+"%")
	}
	return sel
}

func applyRobotListFilters(sel *bun.SelectQuery, filter repository.RobotListFilter) *bun.SelectQuery {
	if filter.SiteID != nil && *filter.SiteID != "" {
		sel = sel.Where(`EXISTS (
			SELECT 1 FROM location l
			WHERE l.id_natural = r.location_id AND l.site_id = ?
		)`, *filter.SiteID)
	}
	if filter.LocationID != nil {
		sel = sel.Where("r.location_id = ?", *filter.LocationID)
	}
	if filter.Status != nil {
		sel = sel.Where("r.status = ?", *filter.Status)
	}
	if filter.RobotType != nil {
		sel = sel.Where("r.robot_type = ?", *filter.RobotType)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		sel = sel.Where("r.name ILIKE ?", "%"+escaped+"%")
	}
	if filter.OnlineRobotIDs != nil {
		ids := *filter.OnlineRobotIDs
		sel = sel.Where("r.status IN (?)", bun.In([]model.RobotStatus{
			model.RobotStatusReady, model.RobotStatusOnline,
		}))
		if filter.ExcludeOnline {
			if len(ids) > 0 {
				sel = sel.Where("r.id_natural NOT IN (?)", bun.In(ids))
			}
		} else {
			sel = sel.Where("r.id_natural IN (?)", bun.In(ids))
		}
	}
	return sel
}

func applyTaskListFilters(sel *bun.SelectQuery, filter repository.TaskListFilter) *bun.SelectQuery {
	if filter.HasApprovedVersion != nil && *filter.HasApprovedVersion {
		sel = sel.Where("EXISTS (SELECT 1 FROM task_version tv WHERE tv.task_id = t.id_natural AND tv.approval_status = 1)")
	}
	if len(filter.Statuses) > 0 {
		sel = sel.Where("t.status IN (?)", bun.In(filter.Statuses))
	}
	if len(filter.Priorities) > 0 {
		sel = sel.Where("t.priority IN (?)", bun.In(filter.Priorities))
	}
	if len(filter.Difficulties) > 0 {
		sel = sel.Where("t.difficulty IN (?)", bun.In(filter.Difficulties))
	}
	if filter.RobotType != nil {
		sel = sel.Where("t.robot_type = ?", *filter.RobotType)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		sel = sel.Where("t.name ILIKE ?", "%"+escaped+"%")
	}
	return sel
}

func applyUserListFilters(sel *bun.SelectQuery, filter repository.UserListFilter) *bun.SelectQuery {
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
	return sel
}

func applySubTaskListFilters(sel *bun.SelectQuery, filter repository.SubTaskListFilter) *bun.SelectQuery {
	// TaskVersionID takes precedence over TaskID.
	if filter.TaskVersionID != nil && *filter.TaskVersionID != "" {
		return sel.Where("task_version_id = ?", *filter.TaskVersionID)
	}
	if filter.TaskID != nil && *filter.TaskID != "" {
		return sel.Where("task_version_id IN (SELECT id_natural FROM task_version WHERE task_id = ?)", *filter.TaskID)
	}
	return sel
}

// namedSQLArgs lets bun's NewRaw resolve ?name placeholders against a flat
// map. bun only consults a NamedArgAppender when args has length 1, so the
// whole map must be passed as the sole argument:
//
//	conn.NewRaw("SELECT ... WHERE id = ?id", namedSQLArgs{"id": x})
//
// Prefer over positional ? when a query has many parameters or reuses the
// same value multiple times.
type namedSQLArgs map[string]any

var _ schema.NamedArgAppender = namedSQLArgs(nil)

func (a namedSQLArgs) AppendNamedArg(gen schema.QueryGen, b []byte, name string) ([]byte, bool) {
	v, ok := a[name]
	if !ok {
		return b, false
	}
	return gen.Append(b, v), true
}
