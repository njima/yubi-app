package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/infra/cache"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/ddtrace"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/infra/persistence"
	"github.com/airoa-org/yubi-app/backend/internal/infra/storage"
	"github.com/airoa-org/yubi-app/backend/internal/platform/config"
	"github.com/airoa-org/yubi-app/backend/internal/platform/log"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/requestctx"

	"github.com/getsentry/sentry-go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/redis/go-redis.v9"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func newApplication(ctx context.Context) (*application, error) {
	conf, err := config.NewConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logger := log.NewZerologLogger(conf.AppName, conf.Environment)
	app := &application{
		conf:   conf,
		logger: logger,
	}

	entity.OrgIDFromContext = func(ctx context.Context) (string, bool) {
		id, err := requestctx.OrganizationID(ctx)
		return id, err == nil && id != ""
	}

	if conf.Datadog.Enabled {
		ddtracer.Start(
			ddtracer.WithService(conf.AppName),
			ddtracer.WithEnv(conf.Environment),
		)
		app.datadogStarted = true
		logger.Info().Msg("datadog tracer started")
	}

	if conf.Sentry.DSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              conf.Sentry.DSN,
			Environment:      conf.Sentry.Environment,
			TracesSampleRate: conf.Sentry.TracesSampleRate,
			EnableTracing:    conf.Sentry.TracesSampleRate > 0,
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to initialize sentry")
		} else {
			app.sentryStarted = true
			logger.Info().Msg("sentry initialized")
		}
	}

	logger.Info().Msg("starting application")

	redisClient, err := cache.NewClient(
		conf.Redis.Host,
		conf.Redis.Port,
		conf.Redis.Password,
		conf.Redis.DB,
	)
	if err != nil {
		app.Close()
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	app.redisClient = redisClient

	if conf.Datadog.Enabled {
		redisClient.EnableDDTrace(redistrace.WithServiceName(conf.AppName + "-redis"))
	}

	s3Client, err := storage.NewClient(ctx, conf.S3.Region, conf.S3.BucketName, conf.S3.PresignedURLExpirySec, conf.S3.Endpoint, conf.S3.PresignEndpoint, conf.S3.AccessKeyID, conf.S3.SecretAccessKey)
	if err != nil {
		app.Close()
		return nil, fmt.Errorf("failed to initialize S3 client: %w", err)
	}

	repos := newRepositories(redisClient, s3Client)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
		conf.Database.SSLMode,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	app.db = db
	if conf.Datadog.Enabled {
		db.AddQueryHook(ddtrace.NewBunHook(conf.AppName + "-db"))
	}
	txRunner := persistence.NewTxRunner(db)
	dataAccess := repository.NewDataAccess(db, txRunner)

	eventBuses := newEventBuses()
	app.robotStatusBus = eventBuses.RobotStatus
	app.episodeBus = eventBuses.Episode
	app.robotEpisodeBus = eventBuses.RobotEpisode
	app.episodeListBus = eventBuses.EpisodeList

	app.usecases = newUsecases(repos, dataAccess, eventBuses, logger)

	return app, nil
}

func (a *application) Close() {
	if a.redisClient != nil {
		if err := a.redisClient.Close(); err != nil {
			a.logger.Error().Err(err).Msg("failed to close redis")
		}
	}
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.logger.Error().Err(err).Msg("failed to close db")
		}
	}
	if a.sentryStarted {
		sentry.Flush(2 * time.Second)
	}
	if a.datadogStarted {
		ddtracer.Stop()
	}
}
