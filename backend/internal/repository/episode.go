package repository

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
)

// MaxEpisodeExportRows is the maximum number of episodes returned in a single export operation.
const MaxEpisodeExportRows = 30_000

type EpisodeListFilter struct {
	TaskID        *string
	TaskVersionID *string
	RobotID       *string
	UserID        *string
	Statuses      []EpisodeStatus
	StartedAtFrom *time.Time
	StartedAtTo   *time.Time
	SortBy        *EpisodeSortBy
	SortOrder     *SortOrder
}

// EpisodeExportFilter is a thin wrapper around EpisodeListFilter, reserved for
// fields that may apply only to the export operation in the future.
type EpisodeExportFilter struct {
	EpisodeListFilter
}

// EpisodeExportRow holds the data for a single row in the episode export CSV.
type EpisodeExportRow struct {
	IDNatural     string
	TaskID        string
	TaskVersionID string
	RobotID       string
	LocationID    string
	UserID        string
	RecordedByID  *string
	Status        openapi.EpisodeCollectionStatus
	StartedAt     *time.Time
	FinishedAt    *time.Time
	CreatedAt     time.Time
}

type Episode interface {
	Create(ctx context.Context, conn DBConn, e model.Episode) (model.Episode, error)
	GetByID(ctx context.Context, conn DBConn, id string) (model.Episode, error)
	// GetCurrentRobotEpisode returns the episode the teleop UI should display
	// for the robot, in priority order:
	//   1. Recording  — oldest by created_at  (in-flight work)
	//   2. Ready      — oldest by created_at  (next queued)
	//   3. Completed  — most recent by finished_at  (last successful run)
	//   4. Cancelled  — most recent by finished_at  (last attempt)
	// Returns nil if the robot has no episodes at all.
	GetCurrentRobotEpisode(ctx context.Context, conn DBConn, robotID string) (*model.Episode, error)
	List(ctx context.Context, conn DBConn, filter EpisodeListFilter, limit, offset int) (model.Episodes, int, error)
	Update(ctx context.Context, conn DBConn, e model.Episode) (model.Episode, error)
	Delete(ctx context.Context, conn DBConn, id string) error
	// SumDurationByTaskID returns the total duration (seconds) of completed episodes
	// across all approved task versions for the given task.
	SumDurationByTaskID(ctx context.Context, conn DBConn, taskID string) (int64, error)
	Export(ctx context.Context, conn DBConn, filter EpisodeExportFilter) ([]EpisodeExportRow, error)
}
