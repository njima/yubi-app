package middleware

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func ErrorLogger(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		logEvent := logger.Info()
		if statusCode >= 500 {
			logEvent = logger.Error()
		} else if statusCode >= 400 {
			logEvent = logger.Warn()
		}

		logEvent.
			Int("status", statusCode).
			Str("method", method).
			Str("path", path).
			Str("client_ip", c.ClientIP()).
			Dur("latency", latency).
			Msg("request completed")
	}
}

type ErrorHandler struct {
	logger zerolog.Logger
}

func NewErrorHandler(logger zerolog.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (e *ErrorHandler) ConvertErrorResponseWithLogging() openapi.StrictMiddlewareFunc {
	return func(f openapi.StrictHandlerFunc, operationID string) openapi.StrictHandlerFunc {
		return func(ctx *gin.Context, request interface{}) (response interface{}, err error) {
			start := time.Now()

			response, err = f(ctx, request)

			latency := time.Since(start)
			method := ctx.Request.Method
			path := ctx.Request.URL.Path

			if err == nil {
				return response, nil
			}

			e.logError(operationID, method, path, ctx.ClientIP(), latency, err)

			// Report error to Sentry (only 5xx errors)
			codes := apperror.GetCodes(err)
			errorResponse := apperror.NewErrorResponse(codes)
			if errorResponse.Code >= 500 {
				if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
					hub.WithScope(func(scope *sentry.Scope) {
						scope.SetTag("operation", operationID)
						scope.SetExtra("request_path", path)
						scope.SetExtra("request_method", method)
						hub.CaptureException(err)
					})
				}
			}

			ctx.JSON(errorResponse.Code, openapi.ErrorResponse{
				Code:    errorResponse.Code,
				Message: errorResponse.Message,
			})

			return nil, nil
		}
	}
}

func (e *ErrorHandler) logError(operationID, method, path, clientIP string, latency time.Duration, err error) {
	appErr := apperror.ConvertAppError(err)
	if appErr != nil {
		codes := apperror.GetCodes(err)
		if len(codes) > 0 {
			e.logger.Error().
				Str("operation", operationID).
				Str("method", method).
				Str("path", path).
				Str("client_ip", clientIP).
				Dur("latency", latency).
				Str("error_code", codes[0].Code).
				Str("error_domain", string(codes[0].Domain)).
				Str("error_message", codes[0].MessageForEndUser).
				Str("developer_message", err.Error()).
				Str("stack_trace", apperror.StackTrace(err)).
				Msg("request failed")
		} else {
			e.logger.Error().
				Str("operation", operationID).
				Str("method", method).
				Str("path", path).
				Str("client_ip", clientIP).
				Dur("latency", latency).
				Str("error", err.Error()).
				Str("stack_trace", apperror.StackTrace(err)).
				Msg("request failed")
		}
	} else {
		e.logger.Error().
			Str("operation", operationID).
			Str("method", method).
			Str("path", path).
			Str("client_ip", clientIP).
			Dur("latency", latency).
			Str("error", err.Error()).
			Msg("request failed")
	}
}
