package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/ccontext"
	"github.com/airoa-org/yubi-app/backend/internal/config"
	"github.com/airoa-org/yubi-app/backend/internal/database/ddtrace"
	"github.com/airoa-org/yubi-app/backend/internal/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/log"
	"github.com/airoa-org/yubi-app/backend/internal/redis"
	s3client "github.com/airoa-org/yubi-app/backend/internal/s3"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"

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
		id, err := ccontext.OrganizationID(ctx)
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

	redisClient, err := redis.NewClient(
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

	s3Client, err := s3client.NewClient(ctx, conf.S3.Region, conf.S3.BucketName, conf.S3.PresignedURLExpirySec, conf.S3.Endpoint, conf.S3.PresignEndpoint, conf.S3.AccessKeyID, conf.S3.SecretAccessKey)
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

	buses := newEventBuses()
	app.robotStatusBus = buses.RobotStatus
	app.episodeBus = buses.Episode
	app.robotEpisodeBus = buses.RobotEpisode
	app.episodeListBus = buses.EpisodeList

	app.userUsecase = usecase.NewUser(repos.User, repos.UserLocation, repos.UserSite, db, logger)
	app.userImportUsecase = usecase.NewUserImport(repos.User, db, logger)
	app.organizationUsecase = usecase.NewOrganization(repos.Organization, db)
	app.siteUsecase = usecase.NewSite(repos.Site, db)
	app.locationUsecase = usecase.NewLocation(repos.Location, db)
	app.robotUsecase = usecase.NewRobot(repos.Robot, repos.RobotStatus, repos.RobotUptimeDelta, db)
	app.robotDeviceUsecase = usecase.NewRobotDevice(repos.Robot, repos.RobotStatus, repos.RobotUptimeDelta, db, logger, app.robotStatusBus)
	app.taskTagUsecase = usecase.NewTaskTag(repos.TaskTag, db)
	app.taskImportUsecase = usecase.NewTaskImport(repos.Task, repos.TaskTag, db)
	app.taskExportUsecase = usecase.NewTaskExport(repos.Task, repos.TaskTag, db)
	app.taskUsecase = usecase.NewTask(repos.Task, repos.TaskTag, repos.Episode, repos.TaskVersion, db)
	app.taskVersionUsecase = usecase.NewTaskVersion(repos.TaskVersion, repos.Task, repos.SubTask, repos.Episode, db)
	app.subtaskUsecase = usecase.NewSubTask(repos.SubTask, repos.Task, repos.TaskVersion, db)
	app.episodeUsecase = usecase.NewEpisode(usecase.EpisodeDependencies{
		Repository:               repos.Episode,
		GradeRepository:          repos.EpisodeGrade,
		Logger:                   logger,
		TaskVersionRepository:    repos.TaskVersion,
		SubTaskRepository:        repos.SubTask,
		EpisodeSubTaskRepository: repos.EpisodeSubTask,
		ExecutionRepository:      repos.EpisodeSubTaskExecution,
		RobotRepository:          repos.Robot,
		RobotStatusRepository:    repos.RobotStatus,
		RecordingRepository:      repos.EpisodeRecording,
		TaskRepository:           repos.Task,
		LocationRepository:       repos.Location,
		SiteRepository:           repos.Site,
		DB:                       db,
		EventBus:                 app.episodeBus,
		RobotEventBus:            app.robotEpisodeBus,
		ListEventBus:             app.episodeListBus,
	})
	app.episodeGradeUsecase = usecase.NewEpisodeGrade(repos.EpisodeGrade, db)
	app.episodeExportUsecase = usecase.NewEpisodeExport(repos.Episode, db)
	app.operatorYieldExportUsecase = usecase.NewOperatorYieldExport(repos.OperatorYield, db, logger)
	app.episodeSubTaskUsecase = usecase.NewEpisodeSubTask(repos.Episode, repos.EpisodeSubTask, db, app.episodeBus, app.robotEpisodeBus, app.episodeListBus)
	app.episodeExecutionUsecase = usecase.NewEpisodeExecution(repos.Episode, repos.EpisodeSubTask, repos.EpisodeSubTaskExecution, db, app.episodeBus, app.robotEpisodeBus, app.episodeListBus)
	app.fleetUsecase = usecase.NewFleet(repos.Fleet, db)
	app.robotOperatorUsecase = usecase.NewRobotOperator(repos.RobotOperator)
	app.apiKeyUsecase = usecase.NewAPIKey(repos.APIKey, repos.User, repos.Robot, db, logger)

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
