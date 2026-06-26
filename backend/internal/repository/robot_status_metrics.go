package repository

import (
	"context"
	"time"
)

type RobotUptimeMetric struct {
	RobotID        string
	OrganizationID string
	LocationID     string
	PeriodStart    time.Time
	UptimeSeconds  int64
}

type RobotUptimeMetricsRepository interface {
	WriteBatch(ctx context.Context, conn Conn, metrics []RobotUptimeMetric) error
}
