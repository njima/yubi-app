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
	"github.com/airoa-org/yubi-app/backend/internal/event"
	"github.com/airoa-org/yubi-app/backend/internal/gateway"
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

	userRepo := gateway.NewUser()
	userLocationRepo := gateway.NewUserLocation()
	userSiteRepo := gateway.NewUserSite()
	orgRepo := gateway.NewOrganization()
	siteRepo := gateway.NewSite()
	locRepo := gateway.NewLocation()
	robotRepo := gateway.NewRobot()
	taskRepo := gateway.NewTask()
	taskTagRepo := gateway.NewTaskTag()
	taskVersionRepo := gateway.NewTaskVersion()
	subtaskRepo := gateway.NewSubTask()
	episodeRepo := gateway.NewEpisode()
	episodeGradeRepo := gateway.NewEpisodeGrade()
	episodeSubTaskRepo := gateway.NewEpisodeSubTask()
	episodeSubTaskExecutionRepo := gateway.NewEpisodeSubTaskExecution()
	apiKeyRepo := gateway.NewAPIKey()

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

	robotStatusRepo := gateway.NewRobotStatus(redisClient)
	robotUptimeDeltaRepo := gateway.NewRobotUptimeDelta(redisClient)

	s3Client, err := s3client.NewClient(ctx, conf.S3.Region, conf.S3.BucketName, conf.S3.PresignedURLExpirySec, conf.S3.Endpoint, conf.S3.PresignEndpoint, conf.S3.AccessKeyID, conf.S3.SecretAccessKey)
	if err != nil {
		app.Close()
		return nil, fmt.Errorf("failed to initialize S3 client: %w", err)
	}
	episodeRecordingRepo := gateway.NewEpisodeRecording(s3Client)

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

	app.userUsecase = usecase.NewUser(userRepo, userLocationRepo, userSiteRepo, db, logger)
	app.userImportUsecase = usecase.NewUserImport(userRepo, db, logger)
	app.organizationUsecase = usecase.NewOrganization(orgRepo, db)
	app.siteUsecase = usecase.NewSite(siteRepo, db)
	app.locationUsecase = usecase.NewLocation(locRepo, db)
	app.robotUsecase = usecase.NewRobot(robotRepo, robotStatusRepo, robotUptimeDeltaRepo, db)
	app.robotStatusBus = event.NewBus()
	app.robotDeviceUsecase = usecase.NewRobotDevice(robotRepo, robotStatusRepo, robotUptimeDeltaRepo, db, logger, app.robotStatusBus)
	app.taskTagUsecase = usecase.NewTaskTag(taskTagRepo, db)
	app.taskImportUsecase = usecase.NewTaskImport(taskRepo, taskTagRepo, db)
	app.taskExportUsecase = usecase.NewTaskExport(taskRepo, taskTagRepo, db)
	app.taskUsecase = usecase.NewTask(taskRepo, taskTagRepo, episodeRepo, taskVersionRepo, db)
	app.taskVersionUsecase = usecase.NewTaskVersion(taskVersionRepo, taskRepo, subtaskRepo, episodeRepo, db)
	app.subtaskUsecase = usecase.NewSubTask(subtaskRepo, taskRepo, taskVersionRepo, db)
	app.episodeBus = event.NewBus()
	app.robotEpisodeBus = event.NewBus()
	app.episodeListBus = event.NewBus()
	app.episodeUsecase = usecase.NewEpisode(usecase.EpisodeDependencies{
		Repository:               episodeRepo,
		GradeRepository:          episodeGradeRepo,
		Logger:                   logger,
		TaskVersionRepository:    taskVersionRepo,
		SubTaskRepository:        subtaskRepo,
		EpisodeSubTaskRepository: episodeSubTaskRepo,
		ExecutionRepository:      episodeSubTaskExecutionRepo,
		RobotRepository:          robotRepo,
		RobotStatusRepository:    robotStatusRepo,
		RecordingRepository:      episodeRecordingRepo,
		TaskRepository:           taskRepo,
		LocationRepository:       locRepo,
		SiteRepository:           siteRepo,
		DB:                       db,
		EventBus:                 app.episodeBus,
		RobotEventBus:            app.robotEpisodeBus,
		ListEventBus:             app.episodeListBus,
	})
	app.episodeGradeUsecase = usecase.NewEpisodeGrade(episodeGradeRepo, db)
	app.episodeExportUsecase = usecase.NewEpisodeExport(episodeRepo, db)
	operatorYieldRepo := gateway.NewOperatorYield()
	app.operatorYieldExportUsecase = usecase.NewOperatorYieldExport(operatorYieldRepo, db, logger)
	app.episodeSubTaskUsecase = usecase.NewEpisodeSubTask(episodeRepo, episodeSubTaskRepo, db, app.episodeBus, app.robotEpisodeBus, app.episodeListBus)
	app.episodeExecutionUsecase = usecase.NewEpisodeExecution(episodeRepo, episodeSubTaskRepo, episodeSubTaskExecutionRepo, db, app.episodeBus, app.robotEpisodeBus, app.episodeListBus)
	fleetRepo := gateway.NewFleet()
	app.fleetUsecase = usecase.NewFleet(fleetRepo, db)
	robotOperatorRepo := gateway.NewRobotOperator(redisClient)
	app.robotOperatorUsecase = usecase.NewRobotOperator(robotOperatorRepo)
	app.apiKeyUsecase = usecase.NewAPIKey(apiKeyRepo, userRepo, robotRepo, db, logger)

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
