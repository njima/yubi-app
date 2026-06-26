package entity

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Site struct {
	bun.BaseModel `bun:"table:site,alias:si"`
	Timestamp
	OrgScoped

	// columns
	ID             int64  `bun:"id,pk,autoincrement"`                        // Auto-increment ID
	IDNatural      string `bun:"id_natural,unique,type:varchar(36),notnull"` // UUID
	OrganizationID string `bun:"organization_id,type:varchar(36),notnull"`   // Organization ID
	Name           string `bun:"name,type:varchar(255),notnull"`             // Site name

	// relations
	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id_natural"`
}

var SiteTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*Site)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`)
}

var SiteIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Site)(nil)).
			Index("site_organization_id_idx").
			Column("organization_id")
	},
}

var _ bun.BeforeAppendModelHook = (*Site)(nil)

func (s *Site) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if s.CreatedAt.IsZero() {
			s.CreatedAt = now
		}

		if s.UpdatedAt.IsZero() {
			s.UpdatedAt = now
		}

	case *bun.UpdateQuery:
		if s.UpdatedAt.IsZero() {
			s.UpdatedAt = now
		}
	}

	return nil
}
