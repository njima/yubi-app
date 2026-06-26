package entity

import "github.com/uptrace/bun"

type UserSiteAssignment struct {
	bun.BaseModel `bun:"table:user_site_assignment,alias:usa"`
	OrgScoped

	UserID         string `bun:"user_id,pk,type:varchar(36),notnull"`
	SiteID         string `bun:"site_id,pk,type:varchar(36),notnull"`
	OrganizationID string `bun:"organization_id,type:varchar(36),notnull"`

	// relations
	Site *Site `bun:"rel:belongs-to,join:site_id=id_natural"`
}

var UserSiteAssignmentTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*UserSiteAssignment)(nil)).
		IfNotExists().
		ForeignKey(`("user_id") REFERENCES "user" ("id_natural") ON DELETE CASCADE`).
		ForeignKey(`("site_id") REFERENCES "site" ("id_natural") ON DELETE CASCADE`).
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural") ON DELETE CASCADE`)
}

var UserSiteAssignmentIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*UserSiteAssignment)(nil)).
			Index("idx_user_site_assignment_site_id").
			Column("site_id")
	},
}
