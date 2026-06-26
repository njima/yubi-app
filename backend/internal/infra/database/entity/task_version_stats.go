package entity

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type TaskVersionStats struct {
	bun.BaseModel `bun:"table:task_version_stats,alias:tvs"`
	Timestamp

	// columns
	ID                   int64  `bun:"id,pk,autoincrement"`
	IDNatural            string `bun:"id_natural,unique,type:varchar(36),notnull"`
	TaskVersionID        string `bun:"task_version_id,unique,type:varchar(36),notnull"`
	TotalDurationSeconds int64  `bun:"total_duration_seconds,type:bigint,notnull,default:0"`
	EpisodeCount         int    `bun:"episode_count,type:int,notnull,default:0"`

	// relations
	TaskVersion *TaskVersion `bun:"rel:belongs-to,join:task_version_id=id_natural"`
}

var TaskVersionStatsTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*TaskVersionStats)(nil)).
		IfNotExists().
		ForeignKey(`("task_version_id") REFERENCES "task_version" ("id_natural") ON DELETE CASCADE`)
}

// No additional index creators — unique constraint on task_version_id
// is already enforced by the column tag and migration index.
var TaskVersionStatsIdxCreators = []IndexQueryCreator{}

var _ bun.BeforeAppendModelHook = (*TaskVersionStats)(nil)

func (e *TaskVersionStats) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if e.CreatedAt.IsZero() {
			e.CreatedAt = now
		}
		if e.UpdatedAt.IsZero() {
			e.UpdatedAt = now
		}
	case *bun.UpdateQuery:
		e.UpdatedAt = now
	}

	return nil
}
