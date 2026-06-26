package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
)

const (
	robotUptimeDeltaKeyFormat  = "robot:uptime_delta:%s"
	robotUptimePeriodKeyFormat = "robot:uptime_period:%s"
)

type robotUptimeDelta struct {
	redisClient *Client
}

func NewRobotUptimeDelta(redisClient *Client) *robotUptimeDelta {
	return &robotUptimeDelta{redisClient: redisClient}
}

func buildUptimeDeltaKey(robotID string) string {
	return fmt.Sprintf(robotUptimeDeltaKeyFormat, robotID)
}

func buildUptimePeriodKey(robotID string) string {
	return fmt.Sprintf(robotUptimePeriodKeyFormat, robotID)
}

// IncrBy adds seconds to the robot's uptime buffer and stores the period_start derived
// from the heartbeat's ReportedAt. The period key is overwritten on each heartbeat so
// the stored value always reflects the latest heartbeat's hour bucket.
func (r *robotUptimeDelta) IncrBy(ctx context.Context, robotID string, seconds int64, periodStart time.Time) error {
	if err := r.redisClient.IncrBy(ctx, buildUptimeDeltaKey(robotID), seconds); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to increment uptime delta for robot %s", robotID))
	}
	if err := r.redisClient.Set(ctx, buildUptimePeriodKey(robotID), periodStart.UTC().Format(time.RFC3339), 0); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to set uptime period for robot %s", robotID))
	}
	return nil
}

// Get reads the accumulated seconds and period_start without modifying Redis.
// Returns seconds=0 and zero Time if the robot has no buffered uptime.
func (r *robotUptimeDelta) Get(ctx context.Context, robotID string) (int64, time.Time, error) {
	seconds, err := r.redisClient.GetInt64(ctx, buildUptimeDeltaKey(robotID))
	if err != nil {
		return 0, time.Time{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to get uptime delta for robot %s", robotID))
	}
	if seconds <= 0 {
		return 0, time.Time{}, nil
	}

	periodStr, err := r.redisClient.Get(ctx, buildUptimePeriodKey(robotID))
	if err != nil {
		if IsNotFound(err) {
			return 0, time.Time{}, nil
		}
		return 0, time.Time{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to get uptime period for robot %s", robotID))
	}

	periodStart, err := time.Parse(time.RFC3339, periodStr)
	if err != nil {
		return 0, time.Time{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to parse uptime period for robot %s: %s", robotID, periodStr))
	}

	return seconds, periodStart, nil
}

// Delete removes both the delta counter and period key for a robot.
// Call only after a successful WriteBatch to avoid losing data on DB failure.
func (r *robotUptimeDelta) Delete(ctx context.Context, robotID string) error {
	if err := r.redisClient.Delete(ctx, buildUptimeDeltaKey(robotID), buildUptimePeriodKey(robotID)); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to delete uptime delta for robot %s", robotID))
	}
	return nil
}
