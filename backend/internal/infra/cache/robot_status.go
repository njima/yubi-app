package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

const (
	robotStatusKeyFormat = "robot:status:%s"
	robotStatusTTL       = 5 * time.Minute
)

type robotStatus struct {
	redisClient *Client
}

func NewRobotStatus(redisClient *Client) *robotStatus {
	return &robotStatus{
		redisClient: redisClient,
	}
}

func buildRobotStatusKey(robotID string) string {
	return fmt.Sprintf(robotStatusKeyFormat, robotID)
}

func (r *robotStatus) Save(ctx context.Context, status repository.RobotStatus) error {
	key := buildRobotStatusKey(status.RobotID)

	if err := r.redisClient.SetJSON(ctx, key, status, robotStatusTTL); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to save robot status"))
	}

	return nil
}

func (r *robotStatus) GetByRobotID(ctx context.Context, robotID string) (*repository.RobotStatus, error) {
	key := buildRobotStatusKey(robotID)

	var status repository.RobotStatus
	if err := r.redisClient.GetJSON(ctx, key, &status); err != nil {
		if IsNotFound(err) {
			return nil, nil
		}
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to get robot status"))
	}

	return &status, nil
}

func (r *robotStatus) GetByRobotIDs(ctx context.Context, robotIDs []string) (map[string]*repository.RobotStatus, error) {
	if len(robotIDs) == 0 {
		return map[string]*repository.RobotStatus{}, nil
	}

	keys := make([]string, len(robotIDs))
	for i, id := range robotIDs {
		keys[i] = buildRobotStatusKey(id)
	}

	values, err := r.redisClient.MGetBytes(ctx, keys...)
	if err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to batch get robot statuses"))
	}

	result := make(map[string]*repository.RobotStatus, len(robotIDs))
	for i, data := range values {
		if data == nil {
			continue
		}
		var status repository.RobotStatus
		if err := json.Unmarshal(data, &status); err != nil {
			continue
		}
		result[robotIDs[i]] = &status
	}
	return result, nil
}

func (r *robotStatus) Delete(ctx context.Context, robotID string) error {
	key := buildRobotStatusKey(robotID)

	if err := r.redisClient.Delete(ctx, key); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to delete robot status"))
	}

	return nil
}

func (r *robotStatus) GetAllOnlineRobotIDs(ctx context.Context) ([]string, error) {
	keys, err := r.redisClient.Scan(ctx, "robot:status:*")
	if err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to scan robot status keys"))
	}
	ids := make([]string, 0, len(keys))
	for _, key := range keys {
		ids = append(ids, strings.TrimPrefix(key, "robot:status:"))
	}
	return ids, nil
}
