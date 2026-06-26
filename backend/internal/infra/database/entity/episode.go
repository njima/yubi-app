package entity

import (
	"context"
	"encoding/json"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/uptrace/bun"
)

type Episode struct {
	bun.BaseModel `bun:"table:episode,alias:e"`
	Timestamp
	OrgScoped

	// columns
	ID               int64               `bun:"id,pk,autoincrement"`
	IDNatural        string              `bun:"id_natural,unique,type:varchar(36),notnull"` // episode_id (UUID)
	OrganizationID   string              `bun:"organization_id,type:varchar(36),notnull"`   // Organization ID (stores organization's id_natural)
	TaskVersionID    string              `bun:"task_version_id,type:varchar(36),notnull"`   // FK to task_version
	LocationID       string              `bun:"location_id,type:varchar(36),notnull"`       // FK to location
	RobotID          string              `bun:"robot_id,type:varchar(36),notnull"`
	UserID           string              `bun:"user_id,type:varchar(36),notnull"` // user or automation
	RecordedByID     *string             `bun:"recorded_by,type:varchar(36)"`
	StartedAt        *time.Time          `bun:"started_at,type:timestamptz"`
	FinishedAt       *time.Time          `bun:"finished_at,type:timestamptz"`
	CollectionStatus model.EpisodeStatus `bun:"collection_status,type:smallint,notnull,default:0"` // 0: ready, 1: recording, 2: cancel, 3: completed
	ErrorDetails     *string             `bun:"error_details,type:jsonb"`
	ParameterValues  json.RawMessage     `bun:"parameter_values,type:jsonb"` // JSON map of resolved parameter values

	// relations
	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id_natural"`
	TaskVersion  *TaskVersion  `bun:"rel:belongs-to,join:task_version_id=id_natural"`
	Location     *Location     `bun:"rel:belongs-to,join:location_id=id_natural"`
	Robot        *Robot        `bun:"rel:belongs-to,join:robot_id=id_natural"`
	User         *User         `bun:"rel:belongs-to,join:user_id=id_natural"`
}

var EpisodeTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*Episode)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`)
}

var EpisodeIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Episode)(nil)).
			Index("episode_organization_id_idx").
			Column("organization_id")
	},
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Episode)(nil)).
			Index("episode_location_id_idx").
			Column("location_id")
	},
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Episode)(nil)).
			Index("episode_robot_id_idx").
			Column("robot_id")
	},
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Episode)(nil)).
			Index("episode_user_id_idx").
			Column("user_id")
	},
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Episode)(nil)).
			Index("episode_started_at_idx").
			Column("started_at")
	},
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Episode)(nil)).
			Index("episode_task_version_stats_idx").
			Column("task_version_id", "collection_status").
			Include("started_at", "finished_at").
			Where("collection_status = 3")
	},
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Episode)(nil)).
			Index("episode_org_created_at_idx").
			ColumnExpr("organization_id, created_at DESC")
	},
	// At most one Recording episode per robot. Belt-and-braces against
	// stuck/duplicate Recording rows that bypass robot.CanStartTeleoperation
	// (e.g. data drift or concurrent Start() calls).
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Episode)(nil)).
			Unique().
			Index("episode_one_recording_per_robot").
			Column("robot_id").
			Where("collection_status = 1")
	},
}

var _ bun.BeforeAppendModelHook = (*Episode)(nil)

func (e *Episode) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
		if e.UpdatedAt.IsZero() {
			e.UpdatedAt = now
		}
	}

	return nil
}
