package entity

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:user,alias:u"`
	Timestamp
	OrgScoped

	// columns
	ID             int64  `bun:"id,pk,autoincrement"`                        // Auto-increment ID
	IDNatural      string `bun:"id_natural,unique,type:varchar(36),notnull"` // UUID
	OrganizationID string `bun:"organization_id,type:varchar(36),notnull"`   // Organization ID (stores organization's id_natural)
	Name           string `bun:"name,type:varchar(255),notnull"`             // User name
	Email          string `bun:"email,unique,type:varchar(255),notnull"`     // Email address
	Role           uint   `bun:"role,type:smallint,notnull,default:0"`       // Role (0: Admin, 1: Data Engineer, 2: Operator, 3: Viewer)

	// relations
	Organization        *Organization            `bun:"rel:belongs-to,join:organization_id=id_natural"`
	LocationAssignments []UserLocationAssignment `bun:"rel:has-many,join:id_natural=user_id"`
	SiteAssignments     []UserSiteAssignment     `bun:"rel:has-many,join:id_natural=user_id"`
}

var UserTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*User)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`)
}

var UserIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*User)(nil)).
			Index("user_organization_id_idx").
			Column("organization_id")
	},
}

var _ bun.BeforeAppendModelHook = (*User)(nil)

func (u *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if u.CreatedAt.IsZero() {
			u.CreatedAt = now
		}

		if u.UpdatedAt.IsZero() {
			u.UpdatedAt = now
		}

	case *bun.UpdateQuery:
		if u.UpdatedAt.IsZero() {
			u.UpdatedAt = now
		}
	}

	return nil
}
