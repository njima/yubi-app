package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/airoa-org/yubi-app/backend/internal/config"
	"github.com/airoa-org/yubi-app/backend/internal/database"
	"github.com/airoa-org/yubi-app/backend/internal/infra/cache"
	"github.com/airoa-org/yubi-app/backend/internal/infra/persistence"
	"github.com/airoa-org/yubi-app/backend/internal/platform/log"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "write-robot-status-metrics: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	conf, err := config.NewConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger := log.NewZerologLogger(conf.AppName, conf.Environment)

	if conf.Sentry.DSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              conf.Sentry.DSN,
			Environment:      conf.Sentry.Environment,
			TracesSampleRate: conf.Sentry.TracesSampleRate,
			EnableTracing:    conf.Sentry.TracesSampleRate > 0,
		}); err != nil {
			logger.Error().Err(err).Msg("failed to initialize sentry")
		} else {
			defer sentry.Flush(2 * time.Second)
		}
	}

	redisClient, err := cache.NewClient(
		conf.Redis.Host,
		conf.Redis.Port,
		conf.Redis.Password,
		conf.Redis.DB,
	)
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	defer redisClient.Close()

	db, err := database.NewDatabase(
		"postgres",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
		conf.Database.SSLMode,
	)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	robotRepo := persistence.NewRobot()
	robotUptimeDeltaRepo := cache.NewRobotUptimeDelta(redisClient)
	robotUptimeMetricsRepo := persistence.NewRobotUptimeMetrics()
	dataAccess := repository.NewDataAccess(db, persistence.NewTxRunner(db))

	writer := usecase.NewRobotUptimeMetricsWriter(robotRepo, robotUptimeDeltaRepo, robotUptimeMetricsRepo, dataAccess, logger)

	logger.Info().Msg("robot uptime metrics writer started")
	writer.Run(ctx)
	logger.Info().Msg("robot uptime metrics writer stopped")

	return nil
}
