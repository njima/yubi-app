package entity

import "github.com/uptrace/bun"

type UserLocationAssignment struct {
	bun.BaseModel `bun:"table:user_location_assignment,alias:ula"`
	OrgScoped

	UserID         string `bun:"user_id,pk,type:varchar(36),notnull"`
	LocationID     string `bun:"location_id,pk,type:varchar(36),notnull"`
	OrganizationID string `bun:"organization_id,type:varchar(36),notnull"`

	// relations
	Location *Location `bun:"rel:belongs-to,join:location_id=id_natural"`
}

var UserLocationAssignmentTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*UserLocationAssignment)(nil)).
		IfNotExists().
		ForeignKey(`("user_id") REFERENCES "user" ("id_natural") ON DELETE CASCADE`).
		ForeignKey(`("location_id") REFERENCES "location" ("id_natural") ON DELETE CASCADE`).
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural") ON DELETE CASCADE`)
}

var UserLocationAssignmentIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*UserLocationAssignment)(nil)).
			Index("idx_user_location_assignment_location_id").
			Column("location_id")
	},
}
