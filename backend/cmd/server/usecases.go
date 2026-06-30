package main

import (
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
	"github.com/rs/zerolog"
)

type usecases struct {
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
}

func newUsecases(repos repositories, dataAccess repository.DataAccess, eventBuses eventBuses, logger zerolog.Logger) usecases {
	return usecases{
		userUsecase:         usecase.NewUser(repos.User, repos.Organization, repos.OrganizationMembership, repos.UserLocation, repos.UserSite, dataAccess, logger),
		userImportUsecase:   usecase.NewUserImport(repos.User, repos.OrganizationMembership, dataAccess, logger),
		organizationUsecase: usecase.NewOrganization(repos.Organization, dataAccess),
		siteUsecase:         usecase.NewSite(repos.Site, dataAccess),
		locationUsecase:     usecase.NewLocation(repos.Location, dataAccess),
		robotUsecase:        usecase.NewRobot(repos.Robot, repos.RobotStatus, repos.RobotUptimeDelta, dataAccess),
		robotDeviceUsecase:  usecase.NewRobotDevice(repos.Robot, repos.RobotStatus, repos.RobotUptimeDelta, dataAccess, logger, eventBuses.RobotStatus),
		taskTagUsecase:      usecase.NewTaskTag(repos.TaskTag, dataAccess),
		taskImportUsecase:   usecase.NewTaskImport(repos.Task, repos.TaskTag, dataAccess),
		taskExportUsecase:   usecase.NewTaskExport(repos.Task, repos.TaskTag, dataAccess),
		taskUsecase:         usecase.NewTask(repos.Task, repos.TaskTag, repos.Episode, repos.TaskVersion, dataAccess),
		taskVersionUsecase:  usecase.NewTaskVersion(repos.TaskVersion, repos.Task, repos.SubTask, repos.Episode, dataAccess),
		subtaskUsecase:      usecase.NewSubTask(repos.SubTask, repos.Task, repos.TaskVersion, dataAccess),
		episodeUsecase: usecase.NewEpisode(usecase.EpisodeDependencies{
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
			DataAccess:               dataAccess,
			EventBus:                 eventBuses.Episode,
			RobotEventBus:            eventBuses.RobotEpisode,
			ListEventBus:             eventBuses.EpisodeList,
		}),
		episodeGradeUsecase:        usecase.NewEpisodeGrade(repos.EpisodeGrade, dataAccess),
		episodeExportUsecase:       usecase.NewEpisodeExport(repos.Episode, dataAccess),
		operatorYieldExportUsecase: usecase.NewOperatorYieldExport(repos.OperatorYield, dataAccess, logger),
		episodeSubTaskUsecase:      usecase.NewEpisodeSubTask(repos.Episode, repos.EpisodeSubTask, dataAccess, eventBuses.Episode, eventBuses.RobotEpisode, eventBuses.EpisodeList),
		episodeExecutionUsecase:    usecase.NewEpisodeExecution(repos.Episode, repos.EpisodeSubTask, repos.EpisodeSubTaskExecution, dataAccess, eventBuses.Episode, eventBuses.RobotEpisode, eventBuses.EpisodeList),
		fleetUsecase:               usecase.NewFleet(repos.Fleet, dataAccess),
		robotOperatorUsecase:       usecase.NewRobotOperator(repos.RobotOperator),
		apiKeyUsecase:              usecase.NewAPIKey(repos.APIKey, repos.User, repos.Robot, dataAccess, logger),
	}
}
