package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/bunconv"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
	"github.com/uptrace/bun"
)

type episode struct{}

func NewEpisode() *episode { return &episode{} }

func (e *episode) Create(ctx context.Context, conn repository.DBConn, ep model.Episode) (model.Episode, error) {
	var inserted entity.Episode

	dbEp := entity.Episode{
		IDNatural:        ep.IDNatural,
		OrganizationID:   ep.OrganizationID,
		TaskVersionID:    ep.TaskVersionID,
		LocationID:       ep.LocationID,
		RobotID:          ep.RobotID,
		UserID:           ep.UserID,
		RecordedByID:     ep.RecordedByID,
		StartedAt:        ep.StartedAt,
		FinishedAt:       ep.FinishedAt,
		CollectionStatus: ep.Status,
		ErrorDetails:     ep.ErrorDetails,
		ParameterValues:  bunconv.ParameterValuesToJSON(ep.ParameterValues),
	}

	if err := bunConn(conn).NewInsert().
		Model(&dbEp).
		Returning("*").
		Scan(ctx, &inserted); err != nil {
		return model.Episode{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create episode: %v", err))
	}

	var taskVersion entity.TaskVersion
	if err := bunConn(conn).NewSelect().
		Model(&taskVersion).
		Where("id_natural = ?", inserted.TaskVersionID).
		Scan(ctx); err != nil {
		return model.Episode{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task version: %v", err))
	}

	return bunconv.EntityToEpisodeModel(inserted, taskVersion), nil
}

func (e *episode) GetByID(ctx context.Context, conn repository.DBConn, id string) (model.Episode, error) {
	var d entity.Episode
	if err := bunConn(conn).NewSelect().
		Model(&d).
		Where("id_natural = ?", id).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.Episode{}, apperror.NewError(apperror.NewMessage(apperror.CodeEpisodeNotFound, "episode not found: id_natural=%s", id))
		}
		return model.Episode{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get episode: %v", err))
	}

	var taskVersion entity.TaskVersion
	if err := bunConn(conn).NewSelect().
		Model(&taskVersion).
		Where("id_natural = ?", d.TaskVersionID).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.Episode{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskNotFound, "task version not found: id_natural=%s", d.TaskVersionID))
		}
		return model.Episode{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task version: %v", err))
	}

	return bunconv.EntityToEpisodeModel(d, taskVersion), nil
}

// GetCurrentRobotEpisode returns the episode the teleop UI should display
// for the robot, in priority order:
//  1. Recording  — oldest by created_at  (in-flight work)
//  2. Ready      — oldest by created_at  (next queued)
//  3. Completed  — most recent by finished_at  (last successful run)
//  4. Cancelled  — most recent by finished_at  (last attempt)
//
// Returns (nil, nil) if the robot has no episodes at all. The hot path (the
// robot has a Recording episode) is one indexed lookup; the cold path runs
// up to four indexed lookups and only fires on idle robots.
func (e *episode) GetCurrentRobotEpisode(ctx context.Context, conn repository.DBConn, robotID string) (*model.Episode, error) {
	type candidate struct {
		status    model.EpisodeStatus
		orderExpr string
		terminal  bool // require finished_at IS NOT NULL
	}
	candidates := []candidate{
		{model.EpisodeStatusRecording, "e.created_at ASC", false},
		{model.EpisodeStatusReady, "e.created_at ASC", false},
		{model.EpisodeStatusCompleted, "e.finished_at DESC", true},
		{model.EpisodeStatusCancel, "e.finished_at DESC", true},
	}
	for _, c := range candidates {
		var d entity.Episode
		q := bunConn(conn).NewSelect().
			Model(&d).
			Where("e.robot_id = ?", robotID).
			Where("e.collection_status = ?", int(c.status))
		if c.terminal {
			q = q.Where("e.finished_at IS NOT NULL")
		}
		err := q.OrderExpr(c.orderExpr).Limit(1).Scan(ctx)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get current robot episode: %v", err))
		}

		var taskVersion entity.TaskVersion
		if err := bunConn(conn).NewSelect().
			Model(&taskVersion).
			Where("id_natural = ?", d.TaskVersionID).
			Scan(ctx); err != nil {
			return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task version: %v", err))
		}

		ep := bunconv.EntityToEpisodeModel(d, taskVersion)
		return &ep, nil
	}
	return nil, nil
}

var allowedEpisodeSortColumns = map[string]string{
	"task":        "COALESCE(t_sort.name, e.task_version_id)",
	"robot":       "COALESCE(r_sort.name, e.robot_id)",
	"recorded_by": "COALESCE(u_sort.name, e.recorded_by)",
	"started_at":  "e.started_at",
	"ended_at":    "e.finished_at",
	"error":       "CASE WHEN e.error_details IS NOT NULL THEN 1 ELSE 0 END",
}

var nullableEpisodeSortColumns = map[string]bool{
	"recorded_by": true,
	"started_at":  true,
	"ended_at":    true,
}

func applyEpisodeSortOrder(sel *bun.SelectQuery, sortBy *repository.EpisodeSortBy, sortOrder *repository.SortOrder) *bun.SelectQuery {
	if sortBy == nil {
		return sel.OrderExpr("e.created_at DESC")
	}

	col, ok := allowedEpisodeSortColumns[string(*sortBy)]
	if !ok {
		return sel.OrderExpr("e.created_at DESC")
	}

	order := "ASC"
	if sortOrder != nil && *sortOrder == repository.SortOrderDesc {
		order = "DESC"
	}

	nullsClause := ""
	if nullableEpisodeSortColumns[string(*sortBy)] {
		nullsClause = " NULLS LAST"
	}

	return sel.OrderExpr(fmt.Sprintf("%s %s%s, e.created_at DESC, e.id DESC", col, order, nullsClause))
}

// applyStartedAtFilter applies half-open [from, to) bounds on e.started_at to a
// bun SelectQuery. Both bounds are optional. Used by List (body + count) and Export
// to keep the SQL consistent across all episode query paths.
func applyStartedAtFilter(sel *bun.SelectQuery, filter repository.EpisodeListFilter) *bun.SelectQuery {
	if filter.StartedAtFrom != nil {
		sel = sel.Where("e.started_at >= ?", *filter.StartedAtFrom)
	}
	if filter.StartedAtTo != nil {
		sel = sel.Where("e.started_at < ?", *filter.StartedAtTo)
	}
	return sel
}

func (e *episode) List(ctx context.Context, conn repository.DBConn, filter repository.EpisodeListFilter, limit, offset int) (model.Episodes, int, error) {
	var ds []entity.Episode
	sel := bunConn(conn).NewSelect().
		Model(&ds).
		Limit(limit).
		Offset(offset)

	if filter.SortBy != nil {
		switch string(*filter.SortBy) {
		case "task":
			sel = sel.Join("LEFT JOIN task_version AS tv_sort ON tv_sort.id_natural = e.task_version_id")
			sel = sel.Join("LEFT JOIN task AS t_sort ON t_sort.id_natural = tv_sort.task_id")
		case "robot":
			sel = sel.Join("LEFT JOIN robot AS r_sort ON r_sort.id_natural = e.robot_id")
		case "recorded_by":
			sel = sel.Join("LEFT JOIN \"user\" AS u_sort ON u_sort.id_natural = e.recorded_by")
		}
	}

	sel = applyEpisodeSortOrder(sel, filter.SortBy, filter.SortOrder)

	if filter.TaskID != nil {
		sel = sel.
			Join("JOIN task_version AS tv ON tv.id_natural = e.task_version_id").
			Where("tv.task_id = ?", *filter.TaskID)
	}
	if filter.TaskVersionID != nil {
		sel = sel.Where("e.task_version_id = ?", *filter.TaskVersionID)
	}
	if filter.RobotID != nil {
		sel = sel.Where("e.robot_id = ?", *filter.RobotID)
	}
	if filter.UserID != nil {
		sel = sel.Where("e.user_id = ?", *filter.UserID)
	}
	if len(filter.Statuses) > 0 {
		sel = sel.Where("e.collection_status IN (?)", bun.In(filter.Statuses))
	}
	sel = applyStartedAtFilter(sel, filter)

	if err := sel.Scan(ctx); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list episodes: %v", err))
	}

	var total int
	countSel := bunConn(conn).NewSelect().
		Model((*entity.Episode)(nil)).
		ColumnExpr("COUNT(*)")

	if filter.TaskID != nil {
		countSel = countSel.
			Join("JOIN task_version AS tv ON tv.id_natural = e.task_version_id").
			Where("tv.task_id = ?", *filter.TaskID)
	}
	if filter.TaskVersionID != nil {
		countSel = countSel.Where("e.task_version_id = ?", *filter.TaskVersionID)
	}
	if filter.RobotID != nil {
		countSel = countSel.Where("e.robot_id = ?", *filter.RobotID)
	}
	if filter.UserID != nil {
		countSel = countSel.Where("e.user_id = ?", *filter.UserID)
	}
	if len(filter.Statuses) > 0 {
		countSel = countSel.Where("e.collection_status IN (?)", bun.In(filter.Statuses))
	}
	countSel = applyStartedAtFilter(countSel, filter)

	if err := countSel.Scan(ctx, &total); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count episodes: %v", err))
	}

	if len(ds) == 0 {
		return model.Episodes{}, total, nil
	}

	taskVersionIDs := make([]string, 0, len(ds))
	for _, d := range ds {
		taskVersionIDs = append(taskVersionIDs, d.TaskVersionID)
	}

	var taskVersions []entity.TaskVersion
	if err := bunConn(conn).NewSelect().
		Model(&taskVersions).
		Where("id_natural IN (?)", bun.In(taskVersionIDs)).
		Scan(ctx); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task versions: %v", err))
	}

	taskVersionMap := make(map[string]entity.TaskVersion)
	for _, tv := range taskVersions {
		taskVersionMap[tv.IDNatural] = tv
	}

	res := make(model.Episodes, 0, len(ds))
	for _, d := range ds {
		tv, ok := taskVersionMap[d.TaskVersionID]
		if !ok {
			continue
		}

		m := bunconv.EntityToEpisodeModel(d, tv)
		res = append(res, &m)
	}

	return res, total, nil
}

func (e *episode) Update(ctx context.Context, conn repository.DBConn, ep model.Episode) (model.Episode, error) {
	dbEp := entity.Episode{
		IDNatural:        ep.IDNatural,
		TaskVersionID:    ep.TaskVersionID,
		LocationID:       ep.LocationID,
		RobotID:          ep.RobotID,
		UserID:           ep.UserID,
		RecordedByID:     ep.RecordedByID,
		StartedAt:        ep.StartedAt,
		FinishedAt:       ep.FinishedAt,
		CollectionStatus: ep.Status,
		ErrorDetails:     ep.ErrorDetails,
	}

	var updated entity.Episode
	if err := bunConn(conn).NewUpdate().
		Model(&dbEp).
		Where("id_natural = ?", ep.IDNatural).
		ExcludeColumn("id", "id_natural", "organization_id", "created_at", "parameter_values").
		Returning("*").
		Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.Episode{}, apperror.NewError(apperror.NewMessage(apperror.CodeEpisodeNotFound, "episode not found: id_natural=%s", ep.IDNatural))
		}
		return model.Episode{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update episode: %v", err))
	}

	var taskVersion entity.TaskVersion
	if err := bunConn(conn).NewSelect().
		Model(&taskVersion).
		Where("id_natural = ?", updated.TaskVersionID).
		Scan(ctx); err != nil {
		return model.Episode{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task version: %v", err))
	}

	return bunconv.EntityToEpisodeModel(updated, taskVersion), nil
}

func (e *episode) Delete(ctx context.Context, conn repository.DBConn, id string) error {
	var deletedID int64
	if err := bunConn(conn).NewDelete().
		Model((*entity.Episode)(nil)).
		Where("id_natural = ?", id).
		Returning("id").
		Scan(ctx, &deletedID); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewError(apperror.NewMessage(apperror.CodeEpisodeNotFound, "episode not found: id_natural=%s", id))
		}
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete episode: %v", err))
	}
	return nil
}

func (e *episode) Export(ctx context.Context, conn repository.DBConn, filter repository.EpisodeExportFilter) ([]repository.EpisodeExportRow, error) {
	type exportRow struct {
		IDNatural     string              `bun:"id_natural"`
		TaskID        string              `bun:"task_id"`
		TaskVersionID string              `bun:"task_version_id"`
		RobotID       string              `bun:"robot_id"`
		LocationID    string              `bun:"location_id"`
		UserID        string              `bun:"user_id"`
		RecordedByID  *string             `bun:"recorded_by"`
		Status        model.EpisodeStatus `bun:"collection_status"`
		StartedAt     *time.Time          `bun:"started_at"`
		FinishedAt    *time.Time          `bun:"finished_at"`
		CreatedAt     time.Time           `bun:"created_at"`
	}

	var rows []exportRow

	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "organization context required"))
	}

	sel := bunConn(conn).NewSelect().
		TableExpr("episode AS e").
		Join("JOIN task_version AS tv ON tv.id_natural = e.task_version_id").
		ColumnExpr("e.id_natural, tv.task_id, e.task_version_id, e.robot_id, e.location_id, e.user_id, e.recorded_by, e.collection_status, e.started_at, e.finished_at, e.created_at").
		OrderExpr("e.created_at DESC").
		Limit(repository.MaxEpisodeExportRows+1).
		Where("e.organization_id = ?", orgID)

	if filter.TaskID != nil {
		sel = sel.Where("tv.task_id = ?", *filter.TaskID)
	}
	if filter.TaskVersionID != nil {
		sel = sel.Where("e.task_version_id = ?", *filter.TaskVersionID)
	}
	if filter.RobotID != nil {
		sel = sel.Where("e.robot_id = ?", *filter.RobotID)
	}
	if filter.UserID != nil {
		sel = sel.Where("e.user_id = ?", *filter.UserID)
	}
	if len(filter.Statuses) > 0 {
		sel = sel.Where("e.collection_status IN (?)", bun.In(filter.Statuses))
	}
	sel = applyStartedAtFilter(sel, filter.EpisodeListFilter)

	if err := sel.Scan(ctx, &rows); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to export episodes: %v", err))
	}

	result := make([]repository.EpisodeExportRow, 0, len(rows))
	for _, r := range rows {
		result = append(result, repository.EpisodeExportRow{
			IDNatural:     r.IDNatural,
			TaskID:        r.TaskID,
			TaskVersionID: r.TaskVersionID,
			RobotID:       r.RobotID,
			LocationID:    r.LocationID,
			UserID:        r.UserID,
			RecordedByID:  r.RecordedByID,
			Status:        r.Status,
			StartedAt:     r.StartedAt,
			FinishedAt:    r.FinishedAt,
			CreatedAt:     r.CreatedAt,
		})
	}

	return result, nil
}

func (e *episode) SumDurationByTaskID(ctx context.Context, conn repository.DBConn, taskID string) (int64, error) {
	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return 0, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "organization context required"))
	}

	var total int64

	err = bunConn(conn).NewSelect().
		TableExpr("episode AS e").
		Join("JOIN task_version AS tv ON tv.id_natural = e.task_version_id").
		ColumnExpr("COALESCE(SUM(EXTRACT(EPOCH FROM (e.finished_at - e.started_at))::bigint), 0)").
		Where("e.organization_id = ?", orgID).
		Where("tv.task_id = ?", taskID).
		Where("tv.approval_status = ?", model.ApprovalStatusApproved).
		Where("e.collection_status = ?", model.EpisodeStatusCompleted).
		Where("e.finished_at IS NOT NULL").
		Where("e.started_at IS NOT NULL").
		Scan(ctx, &total)

	if err != nil {
		return 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to sum episode duration by task: %v", err))
	}

	return total, nil
}
