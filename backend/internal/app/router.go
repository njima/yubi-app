package app

import (
	"context"
	"net/http"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/authz"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/interfaces/http/controller"
	"github.com/airoa-org/yubi-app/backend/internal/interfaces/http/handler"
	"github.com/airoa-org/yubi-app/backend/internal/interfaces/http/middleware"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

func (a *application) newRouter(ctx context.Context) *gin.Engine {
	ctrl := controller.NewController(
		a.logger,
		a.userUsecase,
		a.userImportUsecase,
		a.organizationUsecase,
		a.siteUsecase,
		a.locationUsecase,
		a.robotUsecase,
		a.robotDeviceUsecase,
		a.taskUsecase,
		a.taskVersionUsecase,
		a.taskTagUsecase,
		a.taskImportUsecase,
		a.taskExportUsecase,
		a.subtaskUsecase,
		a.episodeUsecase,
		a.episodeGradeUsecase,
		a.episodeExportUsecase,
		a.episodeSubTaskUsecase,
		a.episodeExecutionUsecase,
		a.fleetUsecase,
		a.robotOperatorUsecase,
		a.operatorYieldExportUsecase,
		a.apiKeyUsecase,
	)

	router := gin.Default()
	router.ContextWithFallback = true

	if a.conf.Datadog.Enabled {
		router.Use(gintrace.Middleware(a.conf.AppName))
	}

	if a.conf.Sentry.DSN != "" {
		router.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
		}))
	}

	router.Use(middleware.ErrorLogger(a.logger))

	allowedOrigins := []string{"http://localhost:3000"}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-ID", "X-Robot-ID", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health-check", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	errorHandler := middleware.NewErrorHandler(a.logger)
	strictMiddlewares := []openapi.StrictMiddlewareFunc{
		authz.NewAuthzMiddleware(),
		errorHandler.ConvertErrorResponseWithLogging(),
	}
	strictHandler := openapi.NewStrictHandler(ctrl, strictMiddlewares)
	api := router.Group("/api")

	api.Use(middleware.Auth(a.userUsecase, a.robotUsecase, a.apiKeyUsecase))

	const apiBodyLimit = 6 * 1024 * 1024 // 6MB
	api.Use(middleware.MaxBodySize(apiBodyLimit))

	openapi.RegisterHandlers(api, strictHandler)

	sseHandler := handler.NewSSEHandler(ctx, a.logger, a.robotDeviceUsecase, a.episodeUsecase, a.taskUsecase, a.taskVersionUsecase, a.episodeBus, a.robotEpisodeBus, a.episodeListBus, a.robotStatusBus)
	api.GET("/robots/:robotId/status/stream", sseHandler.StreamRobotStatus)
	api.GET("/robots/status/stream", sseHandler.StreamRobotStatusByIds)
	api.GET("/episodes/stream", sseHandler.StreamEpisodeListUpdates)
	api.GET("/episodes/:episodeId/stream", sseHandler.StreamEpisodeUpdates)
	api.GET("/robots/:robotId/teleop/stream", sseHandler.StreamRobotTeleop)

	return router
}
