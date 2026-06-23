package persistence

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type episodeStats struct{}

func NewEpisodeStats() *episodeStats { return &episodeStats{} }

// UpsertHourly inserts or updates hourly stats (idempotent)
func (e *episodeStats) UpsertHourly(ctx context.Context, conn repository.DBConn, stats model.EpisodeStats) error {
	dbEntity := entity.EpisodeStatsHourly{
		IDNatural:            stats.IDNatural,
		OrganizationID:       stats.OrganizationID,
		LocationID:           stats.LocationID,
		RobotID:              stats.RobotID,
		PeriodStart:          stats.PeriodStart,
		TotalDurationSeconds: stats.TotalDurationSeconds,
		EpisodeCount:         stats.EpisodeCount,
	}

	_, err := conn.NewInsert().
		Model(&dbEntity).
		On("CONFLICT (organization_id, location_id, robot_id, period_start) DO UPDATE").
		Set("total_duration_seconds = EXCLUDED.total_duration_seconds").
		Set("episode_count = EXCLUDED.episode_count").
		Set("updated_at = NOW()").
		Exec(ctx)

	if err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to upsert hourly stats: %v", err))
	}

	return nil
}

// UpsertDaily inserts or updates daily stats (idempotent)
func (e *episodeStats) UpsertDaily(ctx context.Context, conn repository.DBConn, stats model.EpisodeStats) error {
	dbEntity := entity.EpisodeStatsDaily{
		IDNatural:            stats.IDNatural,
		OrganizationID:       stats.OrganizationID,
		LocationID:           stats.LocationID,
		RobotID:              stats.RobotID,
		PeriodStart:          stats.PeriodStart,
		TotalDurationSeconds: stats.TotalDurationSeconds,
		EpisodeCount:         stats.EpisodeCount,
	}

	_, err := conn.NewInsert().
		Model(&dbEntity).
		On("CONFLICT (organization_id, location_id, robot_id, period_start) DO UPDATE").
		Set("total_duration_seconds = EXCLUDED.total_duration_seconds").
		Set("episode_count = EXCLUDED.episode_count").
		Set("updated_at = NOW()").
		Exec(ctx)

	if err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to upsert daily stats: %v", err))
	}

	return nil
}

// UpsertMonthly inserts or updates monthly stats (idempotent)
func (e *episodeStats) UpsertMonthly(ctx context.Context, conn repository.DBConn, stats model.EpisodeStats) error {
	dbEntity := entity.EpisodeStatsMonthly{
		IDNatural:            stats.IDNatural,
		OrganizationID:       stats.OrganizationID,
		LocationID:           stats.LocationID,
		RobotID:              stats.RobotID,
		PeriodStart:          stats.PeriodStart,
		TotalDurationSeconds: stats.TotalDurationSeconds,
		EpisodeCount:         stats.EpisodeCount,
	}

	_, err := conn.NewInsert().
		Model(&dbEntity).
		On("CONFLICT (organization_id, location_id, robot_id, period_start) DO UPDATE").
		Set("total_duration_seconds = EXCLUDED.total_duration_seconds").
		Set("episode_count = EXCLUDED.episode_count").
		Set("updated_at = NOW()").
		Exec(ctx)

	if err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to upsert monthly stats: %v", err))
	}

	return nil
}

// BulkReplaceHourly replaces hourly stats for a given period.
func (e *episodeStats) BulkReplaceHourly(ctx context.Context, conn repository.DBConn, periodStart time.Time, statsList []model.EpisodeStats) error {
	// Delete all existing rows for this period (handles removed/excluded episodes)
	if _, err := conn.NewDelete().
		Model((*entity.EpisodeStatsHourly)(nil)).
		Where("period_start = ?", periodStart).
		Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete existing hourly stats for period: %v", err))
	}

	if len(statsList) == 0 {
		return nil
	}

	entities := make([]entity.EpisodeStatsHourly, len(statsList))
	for i, stats := range statsList {
		entities[i] = entity.EpisodeStatsHourly{
			IDNatural:            uuid.New().String(),
			OrganizationID:       stats.OrganizationID,
			LocationID:           stats.LocationID,
			RobotID:              stats.RobotID,
			PeriodStart:          stats.PeriodStart,
			TotalDurationSeconds: stats.TotalDurationSeconds,
			EpisodeCount:         stats.EpisodeCount,
		}
	}

	_, err := conn.NewInsert().
		Model(&entities).
		Exec(ctx)

	if err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to replace hourly stats: %v", err))
	}

	return nil
}

// BulkReplaceDaily replaces daily stats for a given period.
func (e *episodeStats) BulkReplaceDaily(ctx context.Context, conn repository.DBConn, periodStart time.Time, statsList []model.EpisodeStats) error {
	if _, err := conn.NewDelete().
		Model((*entity.EpisodeStatsDaily)(nil)).
		Where("period_start = ?", periodStart).
		Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete existing daily stats for period: %v", err))
	}

	if len(statsList) == 0 {
		return nil
	}

	entities := make([]entity.EpisodeStatsDaily, len(statsList))
	for i, stats := range statsList {
		entities[i] = entity.EpisodeStatsDaily{
			IDNatural:            uuid.New().String(),
			OrganizationID:       stats.OrganizationID,
			LocationID:           stats.LocationID,
			RobotID:              stats.RobotID,
			PeriodStart:          stats.PeriodStart,
			TotalDurationSeconds: stats.TotalDurationSeconds,
			EpisodeCount:         stats.EpisodeCount,
		}
	}

	_, err := conn.NewInsert().
		Model(&entities).
		Exec(ctx)

	if err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to replace daily stats: %v", err))
	}

	return nil
}

// BulkReplaceMonthly replaces monthly stats for a given period.
func (e *episodeStats) BulkReplaceMonthly(ctx context.Context, conn repository.DBConn, periodStart time.Time, statsList []model.EpisodeStats) error {
	if _, err := conn.NewDelete().
		Model((*entity.EpisodeStatsMonthly)(nil)).
		Where("period_start = ?", periodStart).
		Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete existing monthly stats for period: %v", err))
	}

	if len(statsList) == 0 {
		return nil
	}

	entities := make([]entity.EpisodeStatsMonthly, len(statsList))
	for i, stats := range statsList {
		entities[i] = entity.EpisodeStatsMonthly{
			IDNatural:            uuid.New().String(),
			OrganizationID:       stats.OrganizationID,
			LocationID:           stats.LocationID,
			RobotID:              stats.RobotID,
			PeriodStart:          stats.PeriodStart,
			TotalDurationSeconds: stats.TotalDurationSeconds,
			EpisodeCount:         stats.EpisodeCount,
		}
	}

	_, err := conn.NewInsert().
		Model(&entities).
		Exec(ctx)

	if err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to replace monthly stats: %v", err))
	}

	return nil
}

// ListHourly retrieves hourly stats with filters
func (e *episodeStats) ListHourly(ctx context.Context, conn repository.DBConn, filter repository.EpisodeStatsFilter) (model.EpisodeStatsList, error) {
	var entities []entity.EpisodeStatsHourly

	query := conn.NewSelect().Model(&entities).Order("period_start ASC")
	query = applyStatsFilterToSelect(query, filter)

	if err := query.Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list hourly stats: %v", err))
	}

	return hourlyEntitiesToModels(entities), nil
}

// ListDaily retrieves daily stats with filters
func (e *episodeStats) ListDaily(ctx context.Context, conn repository.DBConn, filter repository.EpisodeStatsFilter) (model.EpisodeStatsList, error) {
	var entities []entity.EpisodeStatsDaily

	query := conn.NewSelect().Model(&entities).Order("period_start ASC")
	query = applyStatsFilterToSelect(query, filter)

	if err := query.Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list daily stats: %v", err))
	}

	return dailyEntitiesToModels(entities), nil
}

// ListMonthly retrieves monthly stats with filters
func (e *episodeStats) ListMonthly(ctx context.Context, conn repository.DBConn, filter repository.EpisodeStatsFilter) (model.EpisodeStatsList, error) {
	var entities []entity.EpisodeStatsMonthly

	query := conn.NewSelect().Model(&entities).Order("period_start ASC")
	query = applyStatsFilterToSelect(query, filter)

	if err := query.Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list monthly stats: %v", err))
	}

	return monthlyEntitiesToModels(entities), nil
}

// AggregateEpisodesForPeriod aggregates episode data for a specific time period
func (e *episodeStats) AggregateEpisodesForPeriod(ctx context.Context, conn repository.DBConn, from, to time.Time) ([]model.AggregatedEpisodeData, error) {
	type aggregateResult struct {
		OrganizationID       string `bun:"organization_id"`
		LocationID           string `bun:"location_id"`
		RobotID              string `bun:"robot_id"`
		TotalDurationSeconds int64  `bun:"total_duration_seconds"`
		EpisodeCount         int    `bun:"episode_count"`
	}

	var results []aggregateResult

	err := conn.NewSelect().
		TableExpr("episode AS e").
		ColumnExpr("e.organization_id").
		ColumnExpr("e.location_id").
		ColumnExpr("e.robot_id").
		ColumnExpr("COALESCE(SUM(EXTRACT(EPOCH FROM (LEAST(e.finished_at, ?) - GREATEST(e.started_at, ?)))::bigint), 0) AS total_duration_seconds", to, from).
		ColumnExpr("COUNT(*) AS episode_count").
		Where("e.finished_at IS NOT NULL").
		Where("e.started_at IS NOT NULL").
		Where("e.collection_status = ?", openapi.EpisodeCollectionStatusCompleted).
		Where("e.started_at < ?", to).
		Where("e.finished_at > ?", from).
		Group("e.organization_id", "e.location_id", "e.robot_id").
		Scan(ctx, &results)

	if err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to aggregate episodes: %v", err))
	}

	data := make([]model.AggregatedEpisodeData, len(results))
	for i, r := range results {
		data[i] = model.AggregatedEpisodeData{
			OrganizationID:       r.OrganizationID,
			LocationID:           r.LocationID,
			RobotID:              r.RobotID,
			TotalDurationSeconds: r.TotalDurationSeconds,
			EpisodeCount:         r.EpisodeCount,
		}
	}

	return data, nil
}

// AggregateByTaskVersion aggregates all-time episode data grouped by task_version_id.
// Only counts completed episodes (collection_status = 3).
func (e *episodeStats) AggregateByTaskVersion(ctx context.Context, conn repository.DBConn) ([]model.AggregatedTaskVersionData, error) {
	type aggregateResult struct {
		TaskVersionID        string `bun:"task_version_id"`
		TotalDurationSeconds int64  `bun:"total_duration_seconds"`
		EpisodeCount         int    `bun:"episode_count"`
	}

	var results []aggregateResult

	err := conn.NewSelect().
		TableExpr("episode AS e").
		ColumnExpr("e.task_version_id").
		ColumnExpr("COALESCE(SUM(EXTRACT(EPOCH FROM (e.finished_at - e.started_at))::bigint), 0) AS total_duration_seconds").
		ColumnExpr("COUNT(*) AS episode_count").
		Where("e.finished_at IS NOT NULL").
		Where("e.started_at IS NOT NULL").
		Where("e.collection_status = ?", openapi.EpisodeCollectionStatusCompleted).
		Group("e.task_version_id").
		Scan(ctx, &results)

	if err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to aggregate episodes by task version: %v", err))
	}

	data := make([]model.AggregatedTaskVersionData, len(results))
	for i, r := range results {
		data[i] = model.AggregatedTaskVersionData{
			TaskVersionID:        r.TaskVersionID,
			TotalDurationSeconds: r.TotalDurationSeconds,
			EpisodeCount:         r.EpisodeCount,
		}
	}

	return data, nil
}

// BulkUpsertTaskVersionStats replaces all task version stats with the latest snapshot.
// Rows not present in statsList are deleted to avoid stale data (e.g. after episode deletion).
func (e *episodeStats) BulkUpsertTaskVersionStats(ctx context.Context, conn repository.DBConn, statsList []model.TaskVersionStats) error {
	// If no stats, delete all rows (all completed episodes were removed)
	if len(statsList) == 0 {
		_, err := conn.NewDelete().
			Model((*entity.TaskVersionStats)(nil)).
			Where("TRUE").
			Exec(ctx)
		if err != nil {
			return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete task version stats: %v", err))
		}
		return nil
	}

	// Delete rows not in the new snapshot
	activeIDs := make([]string, len(statsList))
	for i, stats := range statsList {
		activeIDs[i] = stats.TaskVersionID
	}
	if _, err := conn.NewDelete().
		Model((*entity.TaskVersionStats)(nil)).
		Where("task_version_id NOT IN (?)", bun.In(activeIDs)).
		Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to clean up stale task version stats: %v", err))
	}

	// Upsert current stats
	entities := make([]entity.TaskVersionStats, len(statsList))
	for i, stats := range statsList {
		entities[i] = entity.TaskVersionStats{
			IDNatural:            uuid.New().String(),
			TaskVersionID:        stats.TaskVersionID,
			TotalDurationSeconds: stats.TotalDurationSeconds,
			EpisodeCount:         stats.EpisodeCount,
		}
	}

	_, err := conn.NewInsert().
		Model(&entities).
		On("CONFLICT (task_version_id) DO UPDATE").
		Set("total_duration_seconds = EXCLUDED.total_duration_seconds").
		Set("episode_count = EXCLUDED.episode_count").
		Set("updated_at = NOW()").
		Exec(ctx)

	if err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to bulk upsert task version stats: %v", err))
	}

	return nil
}

// Helper functions

func applyStatsFilterToSelect(query *bun.SelectQuery, filter repository.EpisodeStatsFilter) *bun.SelectQuery {
	if filter.OrganizationID != nil {
		query = query.Where("organization_id = ?", *filter.OrganizationID)
	}
	if filter.LocationID != nil {
		query = query.Where("location_id = ?", *filter.LocationID)
	}
	if filter.RobotID != nil {
		query = query.Where("robot_id = ?", *filter.RobotID)
	}
	if filter.From != nil {
		query = query.Where("period_start >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("period_start < ?", *filter.To)
	}
	return query
}

// newEpisodeStatsModel converts common entity fields to an EpisodeStats model.
// All episode stats entity types (Hourly, Daily, Monthly) share identical field structures.
func newEpisodeStatsModel(
	id int64, idNatural, orgID, locID, robotID string,
	periodStart time.Time, duration int64, count int,
	createdAt, updatedAt time.Time,
) *model.EpisodeStats {
	return &model.EpisodeStats{
		ID:                   id,
		IDNatural:            idNatural,
		OrganizationID:       orgID,
		LocationID:           locID,
		RobotID:              robotID,
		PeriodStart:          periodStart,
		TotalDurationSeconds: duration,
		EpisodeCount:         count,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}
}

func hourlyEntitiesToModels(entities []entity.EpisodeStatsHourly) model.EpisodeStatsList {
	models := make(model.EpisodeStatsList, len(entities))
	for i, e := range entities {
		models[i] = newEpisodeStatsModel(e.ID, e.IDNatural, e.OrganizationID, e.LocationID, e.RobotID, e.PeriodStart, e.TotalDurationSeconds, e.EpisodeCount, e.CreatedAt, e.UpdatedAt)
	}
	return models
}

func dailyEntitiesToModels(entities []entity.EpisodeStatsDaily) model.EpisodeStatsList {
	models := make(model.EpisodeStatsList, len(entities))
	for i, e := range entities {
		models[i] = newEpisodeStatsModel(e.ID, e.IDNatural, e.OrganizationID, e.LocationID, e.RobotID, e.PeriodStart, e.TotalDurationSeconds, e.EpisodeCount, e.CreatedAt, e.UpdatedAt)
	}
	return models
}

func monthlyEntitiesToModels(entities []entity.EpisodeStatsMonthly) model.EpisodeStatsList {
	models := make(model.EpisodeStatsList, len(entities))
	for i, e := range entities {
		models[i] = newEpisodeStatsModel(e.ID, e.IDNatural, e.OrganizationID, e.LocationID, e.RobotID, e.PeriodStart, e.TotalDurationSeconds, e.EpisodeCount, e.CreatedAt, e.UpdatedAt)
	}
	return models
}
