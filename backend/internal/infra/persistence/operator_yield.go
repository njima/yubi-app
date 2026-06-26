package persistence

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/requestctx"
)

type operatorYield struct{}

func NewOperatorYield() *operatorYield { return &operatorYield{} }

func (g *operatorYield) Export(
	ctx context.Context,
	conn repository.DBConn,
	filter repository.OperatorYieldExportFilter,
) ([]repository.OperatorYieldExportRow, error) {
	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "organization context required"))
	}

	// JST calendar dates → UTC half-open instant range [from, toExclusive).
	// Filter dates are inclusive on both ends, so toExclusive is
	// the day AFTER DateTo at 00:00 JST.
	fromUTC := time.Date(
		filter.DateFrom.Year(), filter.DateFrom.Month(), filter.DateFrom.Day(),
		0, 0, 0, 0, repository.JSTLocation,
	).UTC()
	toExclusive := time.Date(
		filter.DateTo.Year(), filter.DateTo.Month(), filter.DateTo.Day(),
		0, 0, 0, 0, repository.JSTLocation,
	).AddDate(0, 0, 1).UTC()

	type row struct {
		WorkDate         time.Time `bun:"work_date"`
		OperatorUserID   string    `bun:"operator_user_id"`
		OperatorName     string    `bun:"operator_name"`
		TaskID           string    `bun:"task_id"`
		TaskName         string    `bun:"task_name"`
		FirstStart       time.Time `bun:"first_start"`
		LastEnd          time.Time `bun:"last_end"`
		WorkingSeconds   int64     `bun:"working_seconds"`
		CollectedSeconds int64     `bun:"collected_seconds"`
		DiscardedSeconds int64     `bun:"discarded_seconds"`
		EpisodeCount     int64     `bun:"episode_count"`
	}

	// end_ts = COALESCE(finished_at, updated_at): Cancel episodes don't stamp
	// finished_at, so we fall back to updated_at (NOT NULL on episode). Caveat:
	// a Cancel row that's later edited (admin update) inflates discarded
	// duration — accepted per spec until a dedicated cancel timestamp lands.
	//
	// Optional filter idiom: (CAST(?id AS text) IS NULL OR col = ?id) lets a
	// nil *string passthrough as "no filter" without rebuilding SQL.
	const sqlText = `
WITH filtered AS (
    SELECT
        e.started_at,
        COALESCE(e.finished_at, e.updated_at) AS end_ts,
        e.collection_status,
        COALESCE(e.recorded_by, e.user_id) AS operator_user_id,
        tv.task_id AS task_id,
        e.location_id
    FROM episode AS e
    JOIN task_version AS tv ON tv.id_natural = e.task_version_id
    WHERE e.organization_id = ?org_id
      AND e.collection_status IN (?status_cancel, ?status_completed)
      AND e.started_at IS NOT NULL
      AND e.started_at >= ?from_utc
      AND e.started_at <  ?to_utc_exclusive
      AND (CAST(?location_id AS text) IS NULL OR e.location_id = ?location_id)
      AND (CAST(?task_id     AS text) IS NULL OR tv.task_id    = ?task_id)
      AND (CAST(?user_id     AS text) IS NULL OR COALESCE(e.recorded_by, e.user_id) = ?user_id)
),
grouped AS (
    SELECT
        date_trunc('day', f.started_at AT TIME ZONE 'Asia/Tokyo')::date AS work_date,
        f.operator_user_id,
        f.task_id,
        MIN(f.started_at) AS first_start,
        MAX(f.end_ts)     AS last_end,
        EXTRACT(EPOCH FROM (MAX(f.end_ts) - MIN(f.started_at)))::bigint AS working_seconds,
        COALESCE(SUM(CASE
            WHEN f.collection_status = ?status_completed
                THEN EXTRACT(EPOCH FROM (f.end_ts - f.started_at))
            ELSE 0
        END), 0)::bigint AS collected_seconds,
        COALESCE(SUM(CASE
            WHEN f.collection_status = ?status_cancel
                THEN EXTRACT(EPOCH FROM (f.end_ts - f.started_at))
            ELSE 0
        END), 0)::bigint AS discarded_seconds,
        COUNT(*) FILTER (WHERE f.collection_status = ?status_completed) AS episode_count
    FROM filtered f
    GROUP BY work_date, f.operator_user_id, f.task_id
)
SELECT
    g.work_date,
    g.operator_user_id,
    COALESCE(u.name, g.operator_user_id) AS operator_name,
    g.task_id,
    COALESCE(t.name, g.task_id)          AS task_name,
    g.first_start,
    g.last_end,
    g.working_seconds,
    g.collected_seconds,
    g.discarded_seconds,
    g.episode_count
FROM grouped g
LEFT JOIN "user" u ON u.id_natural = g.operator_user_id
LEFT JOIN task   t ON t.id_natural = g.task_id
ORDER BY g.work_date ASC, operator_name ASC, task_name ASC
LIMIT ?row_limit
`

	args := namedSQLArgs{
		"org_id":           orgID,
		"status_completed": int(model.EpisodeStatusCompleted),
		"status_cancel":    int(model.EpisodeStatusCancel),
		"from_utc":         fromUTC,
		"to_utc_exclusive": toExclusive,
		"location_id":      filter.LocationID,
		"task_id":          filter.TaskID,
		"user_id":          filter.UserID,
		"row_limit":        repository.MaxOperatorYieldExportRows + 1,
	}

	var rows []row
	if err := conn.NewRaw(sqlText, args).Scan(ctx, &rows); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to export operator yield: %v", err))
	}

	result := make([]repository.OperatorYieldExportRow, 0, len(rows))
	for _, r := range rows {
		result = append(result, repository.OperatorYieldExportRow{
			WorkDate:         r.WorkDate,
			OperatorUserID:   r.OperatorUserID,
			OperatorName:     r.OperatorName,
			TaskID:           r.TaskID,
			TaskName:         r.TaskName,
			FirstStart:       r.FirstStart,
			LastEnd:          r.LastEnd,
			WorkingSeconds:   r.WorkingSeconds,
			CollectedSeconds: r.CollectedSeconds,
			DiscardedSeconds: r.DiscardedSeconds,
			EpisodeCount:     r.EpisodeCount,
		})
	}
	return result, nil
}
