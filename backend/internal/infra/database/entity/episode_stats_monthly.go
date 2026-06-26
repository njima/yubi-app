package entity

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type EpisodeStatsMonthly struct {
	bun.BaseModel `bun:"table:episode_stats_monthly,alias:esm"`
	Timestamp
	OrgScoped

	// columns
	ID                   int64     `bun:"id,pk,autoincrement"`
	IDNatural            string    `bun:"id_natural,unique,type:varchar(36),notnull"`
	OrganizationID       string    `bun:"organization_id,type:varchar(36),notnull"`
	LocationID           string    `bun:"location_id,type:varchar(36),notnull"`
	RobotID              string    `bun:"robot_id,type:varchar(36),notnull"`
	PeriodStart          time.Time `bun:"period_start,type:timestamptz,notnull"`
	TotalDurationSeconds int64     `bun:"total_duration_seconds,type:bigint,notnull,default:0"`
	EpisodeCount         int       `bun:"episode_count,type:int,notnull,default:0"`

	// relations
	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id_natural"`
	Location     *Location     `bun:"rel:belongs-to,join:location_id=id_natural"`
	Robot        *Robot        `bun:"rel:belongs-to,join:robot_id=id_natural"`
}

var EpisodeStatsMonthlyTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*EpisodeStatsMonthly)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`).
		ForeignKey(`("location_id") REFERENCES "location" ("id_natural")`).
		ForeignKey(`("robot_id") REFERENCES "robot" ("id_natural")`)
}

var EpisodeStatsMonthlyIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*EpisodeStatsMonthly)(nil)).
			Index("episode_stats_monthly_org_loc_robot_period_idx").
			Unique().
			Column("organization_id", "location_id", "robot_id", "period_start")
	},
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*EpisodeStatsMonthly)(nil)).
			Index("episode_stats_monthly_period_start_idx").
			Column("period_start")
	},
}

var _ bun.BeforeAppendModelHook = (*EpisodeStatsMonthly)(nil)

func (e *EpisodeStatsMonthly) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if e.CreatedAt.IsZero() {
			e.CreatedAt = now
		}
		if e.UpdatedAt.IsZero() {
			e.UpdatedAt = now
		}
	case *bun.UpdateQuery:
		e.UpdatedAt = now
	}

	return nil
}
