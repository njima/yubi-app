package repository

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

type FleetSummaryRow struct {
	SiteID       string              `bun:"site_id"`
	SiteName     string              `bun:"site_name"`
	RobotType    string              `bun:"robot_type"`
	Status       model.RobotStatus   `bun:"status"`
	LeaderStatus *model.LeaderStatus `bun:"leader_status"`
	Count        int                 `bun:"count"`
}

type FleetStatsFilter struct {
	From time.Time
	To   time.Time
}

type FleetStatsRow struct {
	SiteID               string `bun:"site_id"`
	SiteName             string `bun:"site_name"`
	RobotType            string `bun:"robot_type"`
	TotalDurationSeconds int64  `bun:"total_duration_seconds"`
}

type FleetUptimeStatsRow struct {
	SiteID        string `bun:"site_id"`
	SiteName      string `bun:"site_name"`
	RobotType     string `bun:"robot_type"`
	UptimeSeconds int64  `bun:"uptime_seconds"`
	RobotCount    int64  `bun:"robot_count"`
}

type FleetTrendFilter struct {
	Granularity model.FleetTrendGranularity
	From        time.Time
	To          time.Time
}

type FleetTrendRow struct {
	SiteID               string    `bun:"site_id"`
	SiteName             string    `bun:"site_name"`
	RobotType            string    `bun:"robot_type"`
	PeriodStart          time.Time `bun:"period_start"`
	TotalDurationSeconds int64     `bun:"total_duration_seconds"`
}

type Fleet interface {
	GetSummary(ctx context.Context, conn DBConn) ([]FleetSummaryRow, error)
	GetStats(ctx context.Context, conn DBConn, filter FleetStatsFilter) ([]FleetStatsRow, error)
	GetUptimeStats(ctx context.Context, conn DBConn, filter FleetStatsFilter) ([]FleetUptimeStatsRow, error)
	GetCollectionTrend(ctx context.Context, conn DBConn, filter FleetTrendFilter) ([]FleetTrendRow, error)
}
