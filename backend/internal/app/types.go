package app

import (
	"github.com/airoa-org/yubi-app/backend/internal/config"
	"github.com/airoa-org/yubi-app/backend/internal/event"
	"github.com/airoa-org/yubi-app/backend/internal/redis"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type application struct {
	conf   *config.Config
	logger zerolog.Logger

	db          *bun.DB
	redisClient *redis.Client

	datadogStarted bool
	sentryStarted  bool

	userUsecase                usecase.UserUsecase
	userImportUsecase          usecase.UserImportUsecase
	organizationUsecase        usecase.OrganizationUsecase
	siteUsecase                usecase.SiteUsecase
	locationUsecase            usecase.LocationUsecase
	robotUsecase               usecase.RobotUsecase
	robotDeviceUsecase         usecase.RobotDeviceUsecase
	taskUsecase                usecase.TaskUsecase
	taskVersionUsecase         usecase.TaskVersionUsecase
	taskTagUsecase             usecase.TaskTagUsecase
	taskImportUsecase          usecase.TaskImportUsecase
	taskExportUsecase          usecase.TaskExportUsecase
	subtaskUsecase             usecase.SubTaskUsecase
	episodeUsecase             usecase.EpisodeUsecase
	episodeGradeUsecase        usecase.EpisodeGradeUsecase
	episodeExportUsecase       usecase.EpisodeExportUsecase
	episodeSubTaskUsecase      usecase.EpisodeSubTaskUsecase
	episodeExecutionUsecase    usecase.EpisodeExecutionUsecase
	fleetUsecase               usecase.FleetUsecase
	robotOperatorUsecase       usecase.RobotOperatorUsecase
	operatorYieldExportUsecase usecase.OperatorYieldExportUsecase
	apiKeyUsecase              usecase.APIKeyUsecase

	robotStatusBus  *event.Bus
	episodeBus      *event.Bus
	robotEpisodeBus *event.Bus
	episodeListBus  *event.Bus
}
