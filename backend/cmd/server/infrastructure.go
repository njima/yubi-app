package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/airoa-org/yubi-app/backend/internal/infra/cache"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/ddtrace"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/infra/persistence"
	"github.com/airoa-org/yubi-app/backend/internal/infra/storage"
	"github.com/airoa-org/yubi-app/backend/internal/platform/config"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/redis/go-redis.v9"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func configureOrganizationScope() {
	entity.OrgIDFromContext = func(ctx context.Context) (string, bool) {
		id, err := requestctx.OrganizationID(ctx)
		return id, err == nil && id != ""
	}
}

func startDatadog(conf *config.Config, logger zerolog.Logger) bool {
	if !conf.Datadog.Enabled {
		return false
	}
	ddtracer.Start(
		ddtracer.WithService(conf.AppName),
		ddtracer.WithEnv(conf.Environment),
	)
	logger.Info().Msg("datadog tracer started")
	return true
}

func initSentry(conf *config.Config, logger zerolog.Logger) bool {
	if conf.Sentry.DSN == "" {
		return false
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              conf.Sentry.DSN,
		Environment:      conf.Sentry.Environment,
		TracesSampleRate: conf.Sentry.TracesSampleRate,
		EnableTracing:    conf.Sentry.TracesSampleRate > 0,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to initialize sentry")
		return false
	}
	logger.Info().Msg("sentry initialized")
	return true
}

func newRedisClient(conf *config.Config) (*cache.Client, error) {
	return cache.NewClient(
		conf.Redis.Host,
		conf.Redis.Port,
		conf.Redis.Password,
		conf.Redis.DB,
	)
}

func enableRedisTracing(redisClient *cache.Client, conf *config.Config) {
	if !conf.Datadog.Enabled {
		return
	}
	redisClient.EnableDDTrace(redistrace.WithServiceName(conf.AppName + "-redis"))
}

func newStorageClient(ctx context.Context, conf *config.Config) (*storage.Client, error) {
	return storage.NewClient(
		ctx,
		conf.S3.Region,
		conf.S3.BucketName,
		conf.S3.PresignedURLExpirySec,
		conf.S3.Endpoint,
		conf.S3.PresignEndpoint,
		conf.S3.AccessKeyID,
		conf.S3.SecretAccessKey,
	)
}

func newDatabase(conf *config.Config) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(databaseDSN(conf))))
	return bun.NewDB(sqldb, pgdialect.New())
}

func databaseDSN(conf *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
		conf.Database.SSLMode,
	)
}

func enableDatabaseTracing(db *bun.DB, conf *config.Config) {
	if !conf.Datadog.Enabled {
		return
	}
	db.AddQueryHook(ddtrace.NewBunHook(conf.AppName + "-db"))
}

func newDataAccess(db *bun.DB) repository.DataAccess {
	txRunner := persistence.NewTransactionRunner(db)
	return repository.NewDataAccess(db, txRunner)
}
