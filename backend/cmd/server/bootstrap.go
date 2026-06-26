package main

import (
	"context"
	"fmt"

	"github.com/airoa-org/yubi-app/backend/internal/platform/config"
	"github.com/airoa-org/yubi-app/backend/internal/platform/log"
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

	configureOrganizationScope()
	app.datadogStarted = startDatadog(conf, logger)
	app.sentryStarted = initSentry(conf, logger)

	logger.Info().Msg("starting application")

	redisClient, err := newRedisClient(conf)
	if err != nil {
		app.Close()
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	app.redisClient = redisClient
	enableRedisTracing(redisClient, conf)

	s3Client, err := newStorageClient(ctx, conf)
	if err != nil {
		app.Close()
		return nil, fmt.Errorf("failed to initialize S3 client: %w", err)
	}

	repos := newRepositories(redisClient, s3Client)

	db := newDatabase(conf)
	app.db = db
	enableDatabaseTracing(db, conf)
	dataAccess := newDataAccess(db)

	eventBuses := newEventBuses()
	app.robotStatusBus = eventBuses.RobotStatus
	app.episodeBus = eventBuses.Episode
	app.robotEpisodeBus = eventBuses.RobotEpisode
	app.episodeListBus = eventBuses.EpisodeList

	app.usecases = newUsecases(repos, dataAccess, eventBuses, logger)

	return app, nil
}
