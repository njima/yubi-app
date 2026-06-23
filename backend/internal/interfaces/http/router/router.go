package router

import (
	"context"
	"net/http"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"

	"github.com/airoa-org/yubi-app/backend/internal/authz"
	"github.com/airoa-org/yubi-app/backend/internal/event"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/interfaces/http/controller"
	"github.com/airoa-org/yubi-app/backend/internal/interfaces/http/handler"
	"github.com/airoa-org/yubi-app/backend/internal/interfaces/http/middleware"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

const DefaultAPIBodyLimit = 6 * 1024 * 1024 // 6MB

type Config struct {
	AppName        string
	DatadogEnabled bool
	SentryEnabled  bool
	AllowedOrigins []string
	APIBodyLimit   int64
}

type Dependencies struct {
	Config     Config
	Logger     zerolog.Logger
	Controller controller.Dependencies
	Auth       AuthDependencies
	SSE        SSEDependencies
}

type AuthDependencies struct {
	UserUsecase   usecase.UserUsecase
	RobotUsecase  usecase.RobotUsecase
	APIKeyUsecase usecase.APIKeyUsecase
}

type SSEDependencies struct {
	RobotDeviceUsecase  usecase.RobotDeviceUsecase
	EpisodeUsecase      usecase.EpisodeUsecase
	TaskUsecase         usecase.TaskUsecase
	TaskVersionUsecase  usecase.TaskVersionUsecase
	EpisodeBus          *event.Bus
	RobotEpisodeBus     *event.Bus
	EpisodeListBus      *event.Bus
	RobotStatusEventBus *event.Bus
}

func New(ctx context.Context, deps Dependencies) *gin.Engine {
	cfg := deps.Config.withDefaults()
	engine := newEngine(cfg, deps.Logger)
	registerHealthCheck(engine)
	registerAPIRoutes(ctx, engine, cfg, deps)
	return engine
}

func (c Config) withDefaults() Config {
	if len(c.AllowedOrigins) == 0 {
		c.AllowedOrigins = []string{"http://localhost:3000"}
	}
	if c.APIBodyLimit == 0 {
		c.APIBodyLimit = DefaultAPIBodyLimit
	}
	return c
}

func newEngine(cfg Config, logger zerolog.Logger) *gin.Engine {
	engine := gin.Default()
	engine.ContextWithFallback = true

	if cfg.DatadogEnabled {
		engine.Use(gintrace.Middleware(cfg.AppName))
	}

	if cfg.SentryEnabled {
		engine.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
		}))
	}

	engine.Use(middleware.ErrorLogger(logger))
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-ID", "X-Robot-ID", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return engine
}

func registerHealthCheck(engine *gin.Engine) {
	engine.GET("/health-check", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}

func registerAPIRoutes(ctx context.Context, engine *gin.Engine, cfg Config, deps Dependencies) {
	ctrl := controller.NewController(deps.Controller)
	errorHandler := middleware.NewErrorHandler(deps.Logger)
	strictHandler := openapi.NewStrictHandler(ctrl, []openapi.StrictMiddlewareFunc{
		authz.NewAuthzMiddleware(),
		errorHandler.ConvertErrorResponseWithLogging(),
	})

	api := engine.Group("/api")
	api.Use(middleware.Auth(deps.Auth.UserUsecase, deps.Auth.RobotUsecase, deps.Auth.APIKeyUsecase))
	api.Use(middleware.MaxBodySize(cfg.APIBodyLimit))

	openapi.RegisterHandlers(api, strictHandler)
	registerSSERoutes(ctx, api, deps.Logger, deps.SSE)
}

func registerSSERoutes(ctx context.Context, api gin.IRouter, logger zerolog.Logger, deps SSEDependencies) {
	sseHandler := handler.NewSSEHandler(
		ctx,
		logger,
		deps.RobotDeviceUsecase,
		deps.EpisodeUsecase,
		deps.TaskUsecase,
		deps.TaskVersionUsecase,
		deps.EpisodeBus,
		deps.RobotEpisodeBus,
		deps.EpisodeListBus,
		deps.RobotStatusEventBus,
	)

	api.GET("/robots/:robotId/status/stream", sseHandler.StreamRobotStatus)
	api.GET("/robots/status/stream", sseHandler.StreamRobotStatusByIds)
	api.GET("/episodes/stream", sseHandler.StreamEpisodeListUpdates)
	api.GET("/episodes/:episodeId/stream", sseHandler.StreamEpisodeUpdates)
	api.GET("/robots/:robotId/teleop/stream", sseHandler.StreamRobotTeleop)
}
