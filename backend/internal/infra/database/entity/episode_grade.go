package entity

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type EpisodeGrade struct {
	bun.BaseModel `bun:"table:episode_grade,alias:eg"`
	Timestamp
	OrgScoped

	EpisodeID string `bun:"episode_id,pk,type:varchar(36),notnull"`
	UserID    string `bun:"user_id,pk,type:varchar(36),notnull"`

	OrganizationID string    `bun:"organization_id,type:varchar(36),notnull"`
	Grade          float64   `bun:"grade,type:double precision,notnull"`
	Comment        *string   `bun:"comment,type:text"`
	GradedAt       time.Time `bun:"graded_at,type:timestamptz,notnull"`

	Episode      *Episode      `bun:"rel:belongs-to,join:episode_id=id_natural"`
	User         *User         `bun:"rel:belongs-to,join:user_id=id_natural"`
	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id_natural"`
}

var _ bun.BeforeAppendModelHook = (*EpisodeGrade)(nil)

func (eg *EpisodeGrade) BeforeAppendModel(_ context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if eg.CreatedAt.IsZero() {
			eg.CreatedAt = now
		}
		if eg.UpdatedAt.IsZero() {
			eg.UpdatedAt = now
		}
	case *bun.UpdateQuery:
		if eg.UpdatedAt.IsZero() {
			eg.UpdatedAt = now
		}
	}

	return nil
}

var EpisodeGradeTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*EpisodeGrade)(nil)).
		IfNotExists().
		ForeignKey(`("episode_id") REFERENCES "episode" ("id_natural") ON DELETE CASCADE`).
		ForeignKey(`("user_id") REFERENCES "user" ("id_natural")`).
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`)
}

var EpisodeGradeIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*EpisodeGrade)(nil)).
			Index("episode_grade_user_id_idx").
			Column("user_id")
	},
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*EpisodeGrade)(nil)).
			Index("episode_grade_organization_id_idx").
			Column("organization_id")
	},
}
