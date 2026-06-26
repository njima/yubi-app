package entity

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/uptrace/bun"
)

type Task struct {
	bun.BaseModel `bun:"table:task,alias:t"`
	Timestamp
	OrgScoped

	// columns
	ID             int64                `bun:"id,pk,autoincrement"`
	IDNatural      string               `bun:"id_natural,unique,type:varchar(36),notnull"` // task_id (UUID)
	OrganizationID string               `bun:"organization_id,type:varchar(36),notnull"`   // Organization ID (stores organization's id_natural)
	Name           string               `bun:"name,type:varchar(255),notnull"`             // task_name
	Description    *string              `bun:"description,type:text"`
	ManualURL      *string              `bun:"manual_url,type:text"`
	Priority       model.TaskPriority   `bun:"priority,type:smallint,notnull"`
	Difficulty     model.TaskDifficulty `bun:"difficulty,type:smallint,notnull"`
	Status         model.TaskStatus     `bun:"status,type:smallint,notnull,default:0"`
	Deadline       time.Time            `bun:"deadline,type:timestamptz,notnull"`
	RobotType      *string              `bun:"robot_type,type:varchar(255)"`

	// relations
	Organization *Organization  `bun:"rel:belongs-to,join:organization_id=id_natural"`
	TaskVersions []*TaskVersion `bun:"rel:has-many,join:id_natural=task_id"`
}

var TaskTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*Task)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`)
}

var TaskIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Task)(nil)).
			Index("task_organization_id_idx").
			Column("organization_id")
	},
}

var _ bun.BeforeAppendModelHook = (*Task)(nil)

func (t *Task) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if t.CreatedAt.IsZero() {
			t.CreatedAt = now
		}

		if t.UpdatedAt.IsZero() {
			t.UpdatedAt = now
		}

	case *bun.UpdateQuery:
		if t.UpdatedAt.IsZero() {
			t.UpdatedAt = now
		}
	}

	return nil
}
