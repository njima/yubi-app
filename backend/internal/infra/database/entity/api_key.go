package entity

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type APIKey struct {
	bun.BaseModel `bun:"table:api_key,alias:ak"`
	Timestamp
	OrgScoped

	ID             int64      `bun:"id,pk,autoincrement"`
	IDNatural      string     `bun:"id_natural,unique,type:varchar(36),notnull"`
	OrganizationID string     `bun:"organization_id,type:varchar(36),notnull"`
	UserID         string     `bun:"user_id,type:varchar(36),notnull"`
	RobotID        *string    `bun:"robot_id,type:varchar(36)"`
	Name           string     `bun:"name,type:varchar(255),notnull"`
	KeyHash        string     `bun:"key_hash,unique,type:char(64),notnull"`
	KeyHint        string     `bun:"key_hint,type:varchar(16),notnull"`
	ExpiresAt      *time.Time `bun:"expires_at,nullzero"`
	LastUsedAt     *time.Time `bun:"last_used_at,nullzero"`
	RevokedAt      *time.Time `bun:"revoked_at,nullzero"`

	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id_natural"`
	User         *User         `bun:"rel:belongs-to,join:user_id=id_natural"`
	Robot        *Robot        `bun:"rel:belongs-to,join:robot_id=id_natural"`
}

var APIKeyTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*APIKey)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`).
		ForeignKey(`("user_id") REFERENCES "user" ("id_natural") ON DELETE CASCADE`).
		ForeignKey(`("robot_id") REFERENCES "robot" ("id_natural") ON DELETE CASCADE`)
}

var APIKeyIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*APIKey)(nil)).
			Index("api_key_organization_id_idx").
			Column("organization_id")
	},
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*APIKey)(nil)).
			Index("api_key_robot_id_idx").
			Column("robot_id")
	},
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*APIKey)(nil)).
			Index("api_key_user_id_idx").
			Column("user_id")
	},
}

var _ bun.BeforeAppendModelHook = (*APIKey)(nil)

func (a *APIKey) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if a.CreatedAt.IsZero() {
			a.CreatedAt = now
		}
		if a.UpdatedAt.IsZero() {
			a.UpdatedAt = now
		}
	case *bun.UpdateQuery:
		if a.UpdatedAt.IsZero() {
			a.UpdatedAt = now
		}
	}

	return nil
}
