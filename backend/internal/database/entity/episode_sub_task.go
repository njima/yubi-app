package entity

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/uptrace/bun"
)

type EpisodeSubTask struct {
	bun.BaseModel `bun:"table:episode_sub_task,alias:est"`
	Timestamp
	OrgScoped

	// columns
	ID               int64                         `bun:"id,pk,autoincrement"`
	IDNatural        string                        `bun:"id_natural,unique,type:varchar(36),notnull"`        // episode_sub_task_id (UUID)
	OrganizationID   string                        `bun:"organization_id,type:varchar(36),notnull"`          // Organization ID (stores organization's id_natural)
	EpisodeID        string                        `bun:"episode_id,type:varchar(36),notnull"`               // FK to episode
	SubTaskID        string                        `bun:"sub_task_id,type:varchar(36),notnull"`              // FK to subtask
	CollectionStatus model.SubTaskCollectionStatus `bun:"collection_status,type:smallint,notnull,default:0"` // 0:ready, 1:in_progress, 2:completed, 3:skipped, 4:cancelled
	TaskResult       openapi.SubTaskTaskResult     `bun:"task_result,type:smallint,notnull,default:0"`       // 0:undetermined, 1:success, 2:failed

	// relations
	Organization *Organization              `bun:"rel:belongs-to,join:organization_id=id_natural"`
	Episode      *Episode                   `bun:"rel:belongs-to,join:episode_id=id_natural"`
	SubTask      *SubTask                   `bun:"rel:belongs-to,join:sub_task_id=id_natural"`
	Executions   []*EpisodeSubTaskExecution `bun:"rel:has-many,join:id_natural=episode_sub_task_id"`
}

var EpisodeSubTaskTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*EpisodeSubTask)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`).
		ForeignKey(`("episode_id") REFERENCES "episode" ("id_natural")`).
		ForeignKey(`("sub_task_id") REFERENCES "subtask" ("id_natural")`).
		ColumnExpr(`CONSTRAINT "episode_sub_task_episode_sub_task_unique" UNIQUE ("episode_id", "sub_task_id")`)
}

var _ bun.BeforeAppendModelHook = (*EpisodeSubTask)(nil)

func (est *EpisodeSubTask) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if est.CreatedAt.IsZero() {
			est.CreatedAt = now
		}

		if est.UpdatedAt.IsZero() {
			est.UpdatedAt = now
		}

	case *bun.UpdateQuery:
		if est.UpdatedAt.IsZero() {
			est.UpdatedAt = now
		}
	}

	return nil
}
