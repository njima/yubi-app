package requestctx

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
)

type robotID struct{}

func SetRobotID(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, robotID{}, val)
}

func IsExistRobotID(ctx context.Context) bool {
	return ctx.Value(robotID{}) != nil
}

func RobotID(ctx context.Context) (string, error) {
	val := ctx.Value(robotID{})
	if val == nil {
		return "", apperror.NewError(
			apperror.NewMessage(apperror.CodeBadRequest, "robot id not found in context"),
		)
	}
	rid, ok := val.(string)
	if !ok {
		return "", apperror.NewError(
			apperror.NewMessage(apperror.CodeInternal, "robot id type assertion failed"),
		)
	}
	return rid, nil
}
