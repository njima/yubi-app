package entity

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/uptrace/bun"
)

type EpisodeSubTaskExecution struct {
	bun.BaseModel `bun:"table:episode_sub_task_execution,alias:este"`
	Timestamp
	OrgScoped

	// columns
	ID               int64                 `bun:"id,pk,autoincrement"`
	IDNatural        string                `bun:"id_natural,unique,type:varchar(36),notnull"`       // execution_id (UUID)
	OrganizationID   string                `bun:"organization_id,type:varchar(36),notnull"`         // Organization ID (stores organization's id_natural)
	EpisodeSubTaskID string                `bun:"episode_sub_task_id,type:varchar(36),notnull"`     // FK to episode_sub_task
	ExecutionStatus  model.ExecutionStatus `bun:"execution_status,type:smallint,notnull,default:0"` // 0:ready, 1:started, 2:cancelled, 3:finished
	StartedAt        *time.Time            `bun:"started_at,type:timestamptz"`
	FinishedAt       *time.Time            `bun:"finished_at,type:timestamptz"`

	// relations
	Organization   *Organization   `bun:"rel:belongs-to,join:organization_id=id_natural"`
	EpisodeSubTask *EpisodeSubTask `bun:"rel:belongs-to,join:episode_sub_task_id=id_natural"`
}

var EpisodeSubTaskExecutionTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*EpisodeSubTaskExecution)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`).
		ForeignKey(`("episode_sub_task_id") REFERENCES "episode_sub_task" ("id_natural")`)
}

var EpisodeSubTaskExecutionIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*EpisodeSubTaskExecution)(nil)).
			Index("episode_sub_task_execution_episode_sub_task_id_created_at_idx").
			Column("episode_sub_task_id", "created_at")
	},
}

var _ bun.BeforeAppendModelHook = (*EpisodeSubTaskExecution)(nil)

func (este *EpisodeSubTaskExecution) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if este.CreatedAt.IsZero() {
			este.CreatedAt = now
		}

		if este.UpdatedAt.IsZero() {
			este.UpdatedAt = now
		}

	case *bun.UpdateQuery:
		if este.UpdatedAt.IsZero() {
			este.UpdatedAt = now
		}
	}

	return nil
}
