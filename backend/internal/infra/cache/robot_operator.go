package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/redis"
)

const (
	robotOperatorKeyFormat = "robot:operator:%s"
	robotOperatorTTL       = 60 * time.Second
)

type robotOperator struct {
	redisClient *redis.Client
}

func NewRobotOperator(redisClient *redis.Client) *robotOperator {
	return &robotOperator{
		redisClient: redisClient,
	}
}

func buildRobotOperatorKey(robotID string) string {
	return fmt.Sprintf(robotOperatorKeyFormat, robotID)
}

func (r *robotOperator) Save(ctx context.Context, robotID string, operator model.RobotOperator) error {
	key := buildRobotOperatorKey(robotID)

	if err := r.redisClient.SetJSON(ctx, key, operator, robotOperatorTTL); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to save robot operator"))
	}

	return nil
}

func (r *robotOperator) SaveNX(ctx context.Context, robotID string, operator model.RobotOperator) (bool, error) {
	key := buildRobotOperatorKey(robotID)

	acquired, err := r.redisClient.SetNXJSON(ctx, key, operator, robotOperatorTTL)
	if err != nil {
		return false, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to acquire robot operator lock"))
	}

	return acquired, nil
}

func (r *robotOperator) GetByRobotID(ctx context.Context, robotID string) (*model.RobotOperator, error) {
	key := buildRobotOperatorKey(robotID)

	var operator model.RobotOperator
	if err := r.redisClient.GetJSON(ctx, key, &operator); err != nil {
		if redis.IsNotFound(err) {
			return nil, nil
		}
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to get robot operator"))
	}

	return &operator, nil
}

func (r *robotOperator) Delete(ctx context.Context, robotID string) error {
	key := buildRobotOperatorKey(robotID)

	if err := r.redisClient.Delete(ctx, key); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to delete robot operator"))
	}

	return nil
}
