package middleware

import (
	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/authz"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/requestctx"

	"github.com/gin-gonic/gin"
)

// NewAuthzMiddleware returns a StrictMiddlewareFunc that enforces role-based permission checks
// using the OpenAPI operationID.
//
// Flow:
//  1. Operations in authzBypassOperations skip the check entirely.
//  2. Operations with no entry in operationPermissions return 403.
//  3. The user role is retrieved from the request context; returns 403 if the role lacks the required permission.
func NewAuthzMiddleware() openapi.StrictMiddlewareFunc {
	return func(f openapi.StrictHandlerFunc, operationID string) openapi.StrictHandlerFunc {
		return func(ctx *gin.Context, request any) (any, error) {
			if authz.IsAuthzBypassOperation(operationID) {
				return f(ctx, request)
			}

			action, ok := authz.ActionForOperation(operationID)
			if !ok {
				return nil, apperror.NewError(
					apperror.NewMessage(apperror.CodeForbidden,
						"no permission mapping defined for this operation"),
				)
			}

			role, err := requestctx.UserRole(ctx.Request.Context())
			if err != nil {
				return nil, apperror.NewError(
					apperror.NewMessage(apperror.CodeForbidden, "user role not found in context"),
				)
			}

			if !authz.HasPermission(role, action) {
				return nil, apperror.NewError(
					apperror.NewMessage(apperror.CodeForbidden, "insufficient permissions"),
				)
			}

			return f(ctx, request)
		}
	}
}
