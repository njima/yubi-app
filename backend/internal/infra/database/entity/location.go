package entity

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Location struct {
	bun.BaseModel `bun:"table:location,alias:l"`
	Timestamp
	OrgScoped

	// columns
	ID             int64  `bun:"id,pk,autoincrement"`                        // Auto-increment ID
	IDNatural      string `bun:"id_natural,unique,type:varchar(36),notnull"` // UUID
	OrganizationID string `bun:"organization_id,type:varchar(36),notnull"`   // Organization ID
	SiteID         string `bun:"site_id,type:varchar(36),notnull"`           // Site ID
	Name           string `bun:"name,type:varchar(255),notnull"`             // Location name

	// relations
	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id_natural"`
	Site         *Site         `bun:"rel:belongs-to,join:site_id=id_natural"`
}

var LocationTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*Location)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`).
		ForeignKey(`("site_id") REFERENCES "site" ("id_natural")`)
}

var LocationIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Location)(nil)).
			Index("location_organization_id_idx").
			Column("organization_id")
	},
}

var _ bun.BeforeAppendModelHook = (*Location)(nil)

func (l *Location) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if l.CreatedAt.IsZero() {
			l.CreatedAt = now
		}

		if l.UpdatedAt.IsZero() {
			l.UpdatedAt = now
		}

	case *bun.UpdateQuery:
		if l.UpdatedAt.IsZero() {
			l.UpdatedAt = now
		}
	}

	return nil
}
