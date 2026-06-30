package entity

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type OrganizationMembership struct {
	bun.BaseModel `bun:"table:organization_membership,alias:om"`
	Timestamp
	OrgScoped

	// columns
	ID             int64  `bun:"id,pk,autoincrement"`                        // Auto-increment ID
	IDNatural      string `bun:"id_natural,unique,type:varchar(36),notnull"` // UUID
	UserID         string `bun:"user_id,type:varchar(36),notnull,unique:organization_membership_user_id_organization_id_key"`
	OrganizationID string `bun:"organization_id,type:varchar(36),notnull,unique:organization_membership_user_id_organization_id_key"`
	Role           uint   `bun:"role,type:smallint,notnull,default:0"`

	// relations
	User         *User         `bun:"rel:belongs-to,join:user_id=id_natural"`
	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id_natural"`
}

var OrganizationMembershipTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*OrganizationMembership)(nil)).
		IfNotExists().
		ForeignKey(`("user_id") REFERENCES "user" ("id_natural") ON DELETE CASCADE`).
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural") ON DELETE CASCADE`)
}

var OrganizationMembershipIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*OrganizationMembership)(nil)).
			Index("organization_membership_user_id_idx").
			Column("user_id")
	},
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*OrganizationMembership)(nil)).
			Index("organization_membership_organization_id_idx").
			Column("organization_id")
	},
}

var _ bun.BeforeAppendModelHook = (*OrganizationMembership)(nil)

func (om *OrganizationMembership) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if om.CreatedAt.IsZero() {
			om.CreatedAt = now
		}

		if om.UpdatedAt.IsZero() {
			om.UpdatedAt = now
		}

	case *bun.UpdateQuery:
		if om.UpdatedAt.IsZero() {
			om.UpdatedAt = now
		}
	}

	return nil
}
