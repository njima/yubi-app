package entity

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:user,alias:u"`
	Timestamp

	// columns
	ID        int64   `bun:"id,pk,autoincrement"`                        // Auto-increment ID
	IDNatural string  `bun:"id_natural,unique,type:varchar(36),notnull"` // UUID
	GoogleSub string  `bun:"google_sub,unique,type:varchar(255),notnull"`
	Name      string  `bun:"name,type:varchar(255),notnull"`         // User name
	Email     string  `bun:"email,unique,type:varchar(255),notnull"` // Email address
	AvatarURL *string `bun:"avatar_url,type:text"`

	// relations
	Memberships         []OrganizationMembership `bun:"rel:has-many,join:id_natural=user_id"`
	LocationAssignments []UserLocationAssignment `bun:"rel:has-many,join:id_natural=user_id"`
	SiteAssignments     []UserSiteAssignment     `bun:"rel:has-many,join:id_natural=user_id"`
}

var UserTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*User)(nil)).
		IfNotExists()
}

var UserIdxCreators = []IndexQueryCreator{}

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
