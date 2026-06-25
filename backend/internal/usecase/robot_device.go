package usecase

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/event"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/rs/zerolog"
)

// robotSessionThreshold is the maximum gap between heartbeats considered the same session.
// Gaps larger than this indicate the robot restarted — the first heartbeat of a new session
// is excluded from uptime counting since we cannot know how long it was actually online.
//
// Timing relationships:
//   - robotSessionThreshold (2 min) < robotStatusTTL (5 min, gateway/robot_status.go):
//     ensures that when a gap exceeds the threshold, the previous status has already
//     expired from Redis and prev == nil, so no stale delta is produced.
//   - robotSessionThreshold (2 min) < metricsWriteInterval (5 min, usecase/robot_status_metrics.go):
//     guarantees at least one valid heartbeat delta accumulates per flush cycle.
//     Changing any of these three values requires re-evaluating the others.
const robotSessionThreshold = 2 * time.Minute

type RobotDeviceUsecase interface {
	UpdateRobotStatus(ctx context.Context, status repository.RobotStatus) error
	GetRobotStatus(ctx context.Context, robotID string) (*repository.RobotStatus, error)
	RobotExists(ctx context.Context, robotID string) (bool, error)
}

type robotDevice struct {
	robotRepo       repository.Robot
	robotStatusRepo repository.RobotStatusRepository
	uptimeDeltaRepo repository.RobotUptimeDeltaRepository
	data            repository.DataAccess
	logger          zerolog.Logger
	statusBus       *event.Bus
}

func NewRobotDevice(
	robotRepo repository.Robot,
	robotStatusRepo repository.RobotStatusRepository,
	uptimeDeltaRepo repository.RobotUptimeDeltaRepository,
	data repository.DataAccess,
	logger zerolog.Logger,
	statusBus *event.Bus,
) *robotDevice {
	return &robotDevice{
		robotRepo:       robotRepo,
		robotStatusRepo: robotStatusRepo,
		uptimeDeltaRepo: uptimeDeltaRepo,
		data:            data,
		logger:          logger,
		statusBus:       statusBus,
	}
}

func (r *robotDevice) UpdateRobotStatus(ctx context.Context, status repository.RobotStatus) error {
	r.accumulateUptimeDelta(ctx, status)

	if err := r.robotStatusRepo.Save(ctx, status); err != nil {
		return err
	}
	r.statusBus.Notify(status.RobotID)
	return nil
}

// accumulateUptimeDelta computes the time elapsed since the previous heartbeat and
// adds it to the robot's Redis uptime buffer. The period_start is derived from
// status.ReportedAt so that flushes attribute uptime to the correct hour bucket
// regardless of when the flush actually runs.
// The buffer is flushed to PostgreSQL periodically by the write-robot-status-metrics service.
func (r *robotDevice) accumulateUptimeDelta(ctx context.Context, status repository.RobotStatus) {
	prev, err := r.robotStatusRepo.GetByRobotID(ctx, status.RobotID)
	if err != nil || prev == nil {
		// First heartbeat or Redis miss — skip; we cannot compute a delta.
		return
	}

	delta := status.ReportedAt.Sub(prev.ReportedAt)
	if delta <= 0 || delta >= robotSessionThreshold {
		// Negative delta (clock skew) or gap too large (robot restarted) — skip.
		return
	}

	periodStart := status.ReportedAt.UTC().Truncate(time.Hour)
	if err := r.uptimeDeltaRepo.IncrBy(ctx, status.RobotID, int64(delta.Seconds()), periodStart); err != nil {
		r.logger.Warn().Err(err).Str("robot_id", status.RobotID).Msg("failed to accumulate uptime delta")
	}
}

func (r *robotDevice) GetRobotStatus(ctx context.Context, robotID string) (*repository.RobotStatus, error) {
	return r.robotStatusRepo.GetByRobotID(ctx, robotID)
}

func (r *robotDevice) RobotExists(ctx context.Context, robotID string) (bool, error) {
	robot, err := r.robotRepo.GetByID(ctx, r.data.Conn(), robotID)
	if err != nil {
		if apperror.SameKind(err, apperror.KindNotFound) {
			return false, nil
		}
		return false, err
	}
	return robot.ID != 0, nil
}
