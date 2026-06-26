package entity

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type SubTask struct {
	bun.BaseModel `bun:"table:subtask,alias:s"`
	Timestamp
	OrgScoped

	// columns
	ID                    int64   `bun:"id,pk,autoincrement"`
	IDNatural             string  `bun:"id_natural,unique,type:varchar(36),notnull"` // sub_task_id (UUID)
	OrganizationID        string  `bun:"organization_id,type:varchar(36),notnull"`   // Organization ID (stores organization's id_natural)
	TaskVersionID         string  `bun:"task_version_id,type:varchar(36),notnull"`   // FK to task_version
	OrderIndex            int     `bun:"order_index,type:integer,notnull"`           // 1, 2, 3...
	Name                  string  `bun:"name,type:varchar(255),notnull"`
	Description           *string `bun:"description,type:text"`
	TargetDurationSeconds *int    `bun:"target_duration_seconds,type:integer"`

	// relations
	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id_natural"`
	TaskVersion  *TaskVersion  `bun:"rel:belongs-to,join:task_version_id=id_natural"`
}

var SubTaskIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*SubTask)(nil)).
			Index("subtask_organization_id_idx").
			Column("organization_id")
	},
}

var SubTaskTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*SubTask)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`)
}

var _ bun.BeforeAppendModelHook = (*SubTask)(nil)

func (st *SubTask) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if st.CreatedAt.IsZero() {
			st.CreatedAt = now
		}

		if st.UpdatedAt.IsZero() {
			st.UpdatedAt = now
		}

	case *bun.UpdateQuery:
		if st.UpdatedAt.IsZero() {
			st.UpdatedAt = now
		}
	}

	return nil
}
