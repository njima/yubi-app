package requestctx

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
)

type userID struct{}

func SetUserID(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, userID{}, val)
}

func IsExistUserID(ctx context.Context) bool {
	return ctx.Value(userID{}) != nil
}

func UserID(ctx context.Context) (string, error) {
	val := ctx.Value(userID{})
	if val == nil {
		return "", apperror.NewError(
			apperror.NewMessage(apperror.CodeBadRequest, "user id not found in context"),
		)
	}
	uid, ok := val.(string)
	if !ok {
		return "", apperror.NewError(
			apperror.NewMessage(apperror.CodeInternal, "user id type assertion failed"),
		)
	}
	return uid, nil
}
