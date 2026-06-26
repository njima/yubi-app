package persistence

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
)

type robotUptimeMetrics struct{}

func NewRobotUptimeMetrics() *robotUptimeMetrics {
	return &robotUptimeMetrics{}
}

// WriteBatch upserts accumulated uptime seconds into robot_uptime_hourly.
// Each metric carries its own PeriodStart derived from the heartbeat's ReportedAt,
// so flushes that span an hour boundary write to the correct hour bucket.
// On conflict, uptime_seconds is incremented rather than replaced, so partial
// flushes within the same hour are additive.
func (r *robotUptimeMetrics) WriteBatch(ctx context.Context, conn repository.Conn, metrics []repository.RobotUptimeMetric) error {
	if len(metrics) == 0 {
		return nil
	}

	entities := make([]entity.RobotUptimeHourly, 0, len(metrics))
	for _, m := range metrics {
		entities = append(entities, entity.RobotUptimeHourly{
			RobotID:        m.RobotID,
			OrganizationID: m.OrganizationID,
			LocationID:     m.LocationID,
			PeriodStart:    m.PeriodStart,
			UptimeSeconds:  m.UptimeSeconds,
		})
	}

	_, err := bunConn(conn).NewInsert().
		Model(&entities).
		On("CONFLICT (robot_id, period_start) DO UPDATE").
		Set("uptime_seconds = ruh.uptime_seconds + EXCLUDED.uptime_seconds").
		Exec(ctx)

	if err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to write robot uptime metrics"))
	}

	return nil
}
