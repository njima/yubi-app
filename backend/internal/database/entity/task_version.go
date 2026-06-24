package entity

import (
	"context"
	"encoding/json"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/uptrace/bun"
)

type TaskVersion struct {
	bun.BaseModel `bun:"table:task_version,alias:tv"`
	Timestamp
	OrgScoped

	// columns
	ID                              int64                `bun:"id,pk,autoincrement"`
	IDNatural                       string               `bun:"id_natural,unique,type:varchar(36),notnull"` // task_version_id (UUID)
	OrganizationID                  string               `bun:"organization_id,type:varchar(36),notnull"`   // Organization ID (stores organization's id_natural)
	TaskID                          string               `bun:"task_id,type:varchar(36),notnull"`           // FK to task
	Version                         string               `bun:"version,type:varchar(50),notnull"`           // "v1", "v2", "2025-01"
	SchemaHash                      string               `bun:"schema_hash,type:varchar(255)"`
	IsActive                        bool                 `bun:"is_active,type:boolean,notnull,default:true"`
	ApprovalStatus                  model.ApprovalStatus `bun:"approval_status,type:integer,notnull,default:0"` // 0=draft, 1=approved
	TargetDurationSeconds           *int                 `bun:"target_duration_seconds,type:integer"`
	TargetEpisodeCount              *int                 `bun:"target_episode_count,type:integer"`
	TargetDurationPerEpisodeSeconds *int                 `bun:"target_duration_per_episode_seconds,type:integer"`
	DisplayName                     *string              `bun:"display_name,type:varchar(100)"` // Optional human-readable label
	Parameters                      json.RawMessage      `bun:"parameters,type:jsonb"`          // JSON array of {key, values[]}

	// relations
	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id_natural"`
	Task         *Task         `bun:"rel:belongs-to,join:task_id=id_natural"`
	SubTasks     []*SubTask    `bun:"rel:has-many,join:id_natural=task_version_id"`
	Episodes     []*Episode    `bun:"rel:has-many,join:id_natural=task_version_id"`
}

var TaskVersionIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*TaskVersion)(nil)).
			Index("task_version_organization_id_idx").
			Column("organization_id")
	},
}

var TaskVersionTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*TaskVersion)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`)
}

var _ bun.BeforeAppendModelHook = (*TaskVersion)(nil)

func (tv *TaskVersion) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if tv.CreatedAt.IsZero() {
			tv.CreatedAt = now
		}

		if tv.UpdatedAt.IsZero() {
			tv.UpdatedAt = now
		}

	case *bun.UpdateQuery:
		if tv.UpdatedAt.IsZero() {
			tv.UpdatedAt = now
		}
	}

	return nil
}
