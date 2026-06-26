package requestctx

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
)

type organizationID struct{}

func SetOrganizationID(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, organizationID{}, val)
}

func IsExistOrganizationID(ctx context.Context) bool {
	return ctx.Value(organizationID{}) != nil
}

func OrganizationID(ctx context.Context) (string, error) {
	val := ctx.Value(organizationID{})
	if val == nil {
		return "", apperror.NewError(
			apperror.NewMessage(apperror.CodeBadRequest, "organization id not found in context"),
		)
	}
	orgID, ok := val.(string)
	if !ok {
		return "", apperror.NewError(
			apperror.NewMessage(apperror.CodeInternal, "organization id type assertion failed"),
		)
	}
	return orgID, nil
}
