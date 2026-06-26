package repository

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

// EpisodeStatsFilter represents filter options for listing episode stats
type EpisodeStatsFilter struct {
	OrganizationID *string
	LocationID     *string
	RobotID        *string
	From           *time.Time
	To             *time.Time
}

// EpisodeStats defines the repository interface for episode statistics
type EpisodeStats interface {
	// UpsertHourly inserts or updates hourly stats (idempotent)
	UpsertHourly(ctx context.Context, conn Conn, stats model.EpisodeStats) error

	// UpsertDaily inserts or updates daily stats (idempotent)
	UpsertDaily(ctx context.Context, conn Conn, stats model.EpisodeStats) error

	// UpsertMonthly inserts or updates monthly stats (idempotent)
	UpsertMonthly(ctx context.Context, conn Conn, stats model.EpisodeStats) error

	BulkReplaceHourly(ctx context.Context, conn Conn, periodStart time.Time, statsList []model.EpisodeStats) error
	BulkReplaceDaily(ctx context.Context, conn Conn, periodStart time.Time, statsList []model.EpisodeStats) error
	BulkReplaceMonthly(ctx context.Context, conn Conn, periodStart time.Time, statsList []model.EpisodeStats) error

	// ListHourly retrieves hourly stats with filters
	ListHourly(ctx context.Context, conn Conn, filter EpisodeStatsFilter) (model.EpisodeStatsList, error)

	// ListDaily retrieves daily stats with filters
	ListDaily(ctx context.Context, conn Conn, filter EpisodeStatsFilter) (model.EpisodeStatsList, error)

	// ListMonthly retrieves monthly stats with filters
	ListMonthly(ctx context.Context, conn Conn, filter EpisodeStatsFilter) (model.EpisodeStatsList, error)

	// AggregateEpisodesForPeriod aggregates episode data for a specific time period
	// Returns aggregated data grouped by organization, location, and robot
	AggregateEpisodesForPeriod(ctx context.Context, conn Conn, from, to time.Time) ([]model.AggregatedEpisodeData, error)

	// AggregateByTaskVersion aggregates all-time episode data grouped by task_version_id
	// Only counts completed episodes (collection_status = 3)
	AggregateByTaskVersion(ctx context.Context, conn Conn) ([]model.AggregatedTaskVersionData, error)

	// BulkUpsertTaskVersionStats inserts or updates task version stats (idempotent)
	BulkUpsertTaskVersionStats(ctx context.Context, conn Conn, statsList []model.TaskVersionStats) error
}
