package entity

import (
	"context"
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

type Robot struct {
	bun.BaseModel `bun:"table:robot,alias:r"`
	Timestamp
	OrgScoped

	// columns
	ID                   int64            `bun:"id,pk,autoincrement"`                        // Auto-increment ID
	IDNatural            string           `bun:"id_natural,unique,type:varchar(36),notnull"` // UUID
	OrganizationID       string           `bun:"organization_id,type:varchar(36),notnull"`   // Organization ID (stores organization's id_natural)
	LocationID           string           `bun:"location_id,type:varchar(36),notnull"`       // Location ID (stores location's id_natural)
	Name                 string           `bun:"name,type:varchar(255),notnull"`             // Robot name
	RobotType            string           `bun:"robot_type,type:varchar(255)"`               // Robot type name
	Status               uint             `bun:"status,type:smallint,notnull,default:0"`     // Status (0: online, 1: busy, 2: offline)
	LeaderStatus         *uint            `bun:"leader_status,type:smallint"`                // Leader status (0: ready, 1: faulted, 2: maintenance)
	LeaderFaultStartedAt *time.Time       `bun:"leader_fault_started_at,type:timestamptz"`   // Leader status fault start timestamp
	FaultStartedAt       *time.Time       `bun:"fault_started_at,type:timestamptz"`          // Fault start timestamp (manual Faulted state only)
	LastHeartbeatAt      *time.Time       `bun:"last_heartbeat_at,type:timestamptz"`         // Last heartbeat timestamp
	OfflineReason        *string          `bun:"offline_reason,type:varchar(255)"`           // Offline reason
	RobotConfig          *json.RawMessage `bun:"robot_config,type:jsonb"`                    // Robot configuration including camera settings
	ActiveEpisodeID      *string          `bun:"active_episode_id,type:varchar(36)"`         // Active episode ID (currently recording episode's id_natural)
	ActiveUserID         *string          `bun:"active_user_id,type:varchar(36)"`            // Active user ID (currently operating user's id_natural)

	// relations
	Organization  *Organization `bun:"rel:belongs-to,join:organization_id=id_natural"`
	Location      *Location     `bun:"rel:belongs-to,join:location_id=id_natural"`
	ActiveEpisode *Episode      `bun:"rel:belongs-to,join:active_episode_id=id_natural"`
	ActiveUser    *User         `bun:"rel:belongs-to,join:active_user_id=id_natural"`
}

var RobotTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*Robot)(nil)).
		IfNotExists().
		ForeignKey(`("organization_id") REFERENCES "organization" ("id_natural")`)
}

var RobotIdxCreators = []IndexQueryCreator{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Robot)(nil)).
			Index("robot_organization_id_idx").
			Column("organization_id")
	},
}

var _ bun.BeforeAppendModelHook = (*Robot)(nil)

func (r *Robot) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()

	switch query.(type) {
	case *bun.InsertQuery:
		if r.CreatedAt.IsZero() {
			r.CreatedAt = now
		}

		if r.UpdatedAt.IsZero() {
			r.UpdatedAt = now
		}

	case *bun.UpdateQuery:
		if r.UpdatedAt.IsZero() {
			r.UpdatedAt = now
		}
	}

	return nil
}
