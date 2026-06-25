package usecase

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"

	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

// metricsWriteInterval controls how often the Redis uptime buffer is flushed to PostgreSQL.
// Tradeoff: shorter = more DB writes; longer = more data lost if the service crashes.
// At 5 minutes, the worst-case data loss is 5 minutes of uptime — acceptable for
// daily/monthly reporting. Must be greater than robotSessionThreshold (2 min) so that
// at least one valid heartbeat delta can accumulate before each flush.
const metricsWriteInterval = 5 * time.Minute

type RobotUptimeMetricsWriter struct {
	robotRepo       repository.Robot
	uptimeDeltaRepo repository.RobotUptimeDeltaRepository
	metricsRepo     repository.RobotUptimeMetricsRepository
	data            repository.DataAccess
	logger          zerolog.Logger
}

func NewRobotUptimeMetricsWriter(
	robotRepo repository.Robot,
	uptimeDeltaRepo repository.RobotUptimeDeltaRepository,
	metricsRepo repository.RobotUptimeMetricsRepository,
	data repository.DataAccess,
	logger zerolog.Logger,
) *RobotUptimeMetricsWriter {
	return &RobotUptimeMetricsWriter{
		robotRepo:       robotRepo,
		uptimeDeltaRepo: uptimeDeltaRepo,
		metricsRepo:     metricsRepo,
		data:            data,
		logger:          logger,
	}
}

// Run flushes the Redis uptime buffer to PostgreSQL on every metricsWriteInterval.
// Errors are logged but never propagate — uptime writes must not crash the service.
func (w *RobotUptimeMetricsWriter) Run(ctx context.Context) {
	ticker := time.NewTicker(metricsWriteInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info().Msg("robot uptime metrics writer stopped")
			return
		case <-ticker.C:
			if err := w.flush(ctx); err != nil {
				w.logger.Error().Err(err).Msg("failed to flush robot uptime metrics")
				sentry.CaptureException(err)
			}
		}
	}
}

// flush reads accumulated uptime seconds from Redis for each robot and writes them to
// robot_uptime_hourly. Redis keys are deleted only after a successful DB write, so a
// WriteBatch failure leaves the buffer intact for the next flush cycle.
func (w *RobotUptimeMetricsWriter) flush(ctx context.Context) error {
	// limit=0 means no LIMIT clause in bun — fetches all robots.
	robots, _, err := w.robotRepo.List(ctx, w.data.Conn(), repository.RobotListFilter{}, 0, 0)
	if err != nil {
		return err
	}

	var metrics []repository.RobotUptimeMetric
	var processedRobotIDs []string

	for _, robot := range robots {
		seconds, periodStart, err := w.uptimeDeltaRepo.Get(ctx, robot.IDNatural)
		if err != nil {
			w.logger.Warn().Err(err).Str("robot_id", robot.IDNatural).Msg("failed to get uptime delta; skipping")
			continue
		}
		if seconds <= 0 || periodStart.IsZero() {
			continue
		}
		metrics = append(metrics, repository.RobotUptimeMetric{
			RobotID:        robot.IDNatural,
			OrganizationID: robot.OrganizationID,
			LocationID:     robot.LocationID,
			PeriodStart:    periodStart,
			UptimeSeconds:  seconds,
		})
		processedRobotIDs = append(processedRobotIDs, robot.IDNatural)
	}

	if len(metrics) == 0 {
		return nil
	}

	if err := w.metricsRepo.WriteBatch(ctx, w.data.Conn(), metrics); err != nil {
		return err
	}

	// Delete Redis keys only after a successful DB write.
	// A failure here is non-fatal: the delta will be double-counted on the next flush,
	// but that is far preferable to data loss.
	for _, robotID := range processedRobotIDs {
		if err := w.uptimeDeltaRepo.Delete(ctx, robotID); err != nil {
			w.logger.Warn().Err(err).Str("robot_id", robotID).Msg("failed to delete uptime delta after write; may cause double-count on next flush")
		}
	}

	return nil
}
