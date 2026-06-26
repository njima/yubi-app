package usecase

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/eventbus"
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
	UpdateRobotStatus(ctx context.Context, status RobotDeviceStatus) error
	GetRobotStatus(ctx context.Context, robotID string) (*RobotDeviceStatus, error)
	RobotExists(ctx context.Context, robotID string) (bool, error)
}

type RobotDeviceStatus struct {
	RobotID    string
	RobotType  string
	ReportedAt time.Time
	Status     RobotDeviceStatusDetail
}

type RobotDeviceStatusDetail struct {
	Battery        RobotBatteryStatus
	Connection     RobotConnectionStatus
	UptimeSec      float64
	Metrics        []RobotMetric
	GateConditions *model.GateConditionStatus
}

type RobotBatteryStatus struct {
	Pct      int
	Charging bool
}

type RobotConnectionStatus struct {
	QualityPct int
}

type RobotMetric struct {
	Name   string
	Type   string
	Unit   string
	Value  any
	Labels map[string]string
}

type robotDevice struct {
	robotRepo       repository.Robot
	robotStatusRepo repository.RobotStatusRepository
	uptimeDeltaRepo repository.RobotUptimeDeltaRepository
	data            repository.DataAccess
	logger          zerolog.Logger
	statusBus       *eventbus.Bus
}

func NewRobotDevice(
	robotRepo repository.Robot,
	robotStatusRepo repository.RobotStatusRepository,
	uptimeDeltaRepo repository.RobotUptimeDeltaRepository,
	data repository.DataAccess,
	logger zerolog.Logger,
	statusBus *eventbus.Bus,
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

func (r *robotDevice) UpdateRobotStatus(ctx context.Context, status RobotDeviceStatus) error {
	repoStatus := status.repositoryStatus()
	r.accumulateUptimeDelta(ctx, repoStatus)

	if err := r.robotStatusRepo.Save(ctx, repoStatus); err != nil {
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

func (r *robotDevice) GetRobotStatus(ctx context.Context, robotID string) (*RobotDeviceStatus, error) {
	status, err := r.robotStatusRepo.GetByRobotID(ctx, robotID)
	if err != nil || status == nil {
		return nil, err
	}
	return robotDeviceStatus(status), nil
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

func (s RobotDeviceStatus) repositoryStatus() repository.RobotStatus {
	metrics := make([]repository.RobotMetric, len(s.Status.Metrics))
	for i, metric := range s.Status.Metrics {
		metrics[i] = repository.RobotMetric{
			Name:   metric.Name,
			Type:   metric.Type,
			Unit:   metric.Unit,
			Value:  metric.Value,
			Labels: metric.Labels,
		}
	}

	return repository.RobotStatus{
		RobotID:    s.RobotID,
		RobotType:  s.RobotType,
		ReportedAt: s.ReportedAt,
		Status: repository.RobotStatusDetail{
			Battery: repository.BatteryStatus{
				Pct:      s.Status.Battery.Pct,
				Charging: s.Status.Battery.Charging,
			},
			Connection: repository.ConnectionStatus{
				QualityPct: s.Status.Connection.QualityPct,
			},
			UptimeSec:      s.Status.UptimeSec,
			Metrics:        metrics,
			GateConditions: s.Status.GateConditions,
		},
	}
}

func robotDeviceStatus(status *repository.RobotStatus) *RobotDeviceStatus {
	metrics := make([]RobotMetric, len(status.Status.Metrics))
	for i, metric := range status.Status.Metrics {
		metrics[i] = RobotMetric{
			Name:   metric.Name,
			Type:   metric.Type,
			Unit:   metric.Unit,
			Value:  metric.Value,
			Labels: metric.Labels,
		}
	}

	return &RobotDeviceStatus{
		RobotID:    status.RobotID,
		RobotType:  status.RobotType,
		ReportedAt: status.ReportedAt,
		Status: RobotDeviceStatusDetail{
			Battery: RobotBatteryStatus{
				Pct:      status.Status.Battery.Pct,
				Charging: status.Status.Battery.Charging,
			},
			Connection: RobotConnectionStatus{
				QualityPct: status.Status.Connection.QualityPct,
			},
			UptimeSec:      status.Status.UptimeSec,
			Metrics:        metrics,
			GateConditions: status.Status.GateConditions,
		},
	}
}
