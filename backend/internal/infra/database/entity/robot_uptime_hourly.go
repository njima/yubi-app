package entity

import (
	"time"

	"github.com/uptrace/bun"
)

// RobotUptimeHourly stores accumulated robot uptime seconds per hour bucket.
// Rows are upserted (not replaced) by the write-robot-status-metrics service,
// which flushes the Redis uptime delta buffer every 5 minutes.
type RobotUptimeHourly struct {
	bun.BaseModel `bun:"table:robot_uptime_hourly,alias:ruh"`
	OrgScoped

	RobotID        string    `bun:"robot_id,pk,type:varchar(36),notnull"`
	OrganizationID string    `bun:"organization_id,type:varchar(36),notnull"`
	LocationID     string    `bun:"location_id,type:varchar(36),notnull"`
	PeriodStart    time.Time `bun:"period_start,pk,type:timestamptz,notnull"`
	UptimeSeconds  int64     `bun:"uptime_seconds,type:bigint,notnull,default:0"`

	// relations
	Robot    *Robot    `bun:"rel:belongs-to,join:robot_id=id_natural"`
	Location *Location `bun:"rel:belongs-to,join:location_id=id_natural"`
}

var RobotUptimeHourlyTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*RobotUptimeHourly)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION`).
		ForeignKey(`("location_id") REFERENCES "location" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION`).
		ForeignKey(`("robot_id") REFERENCES "robot" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION`)
}

var RobotUptimeHourlyIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*RobotUptimeHourly)(nil)).
			Index("robot_uptime_hourly_period_start_idx").
			Column("period_start")
	},
}
