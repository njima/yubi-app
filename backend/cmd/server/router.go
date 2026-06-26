package main

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/interfaces/http/controller"
	httprouter "github.com/airoa-org/yubi-app/backend/internal/interfaces/http/router"
	"github.com/gin-gonic/gin"
)

func (a *application) newRouter(ctx context.Context) *gin.Engine {
	return httprouter.New(ctx, httprouter.Dependencies{
		Config: httprouter.Config{
			AppName:        a.conf.AppName,
			DatadogEnabled: a.conf.Datadog.Enabled,
			SentryEnabled:  a.conf.Sentry.DSN != "",
		},
		Logger: a.logger,
		Controller: controller.Dependencies{
			Logger:                     a.logger,
			UserUsecase:                a.userUsecase,
			UserImportUsecase:          a.userImportUsecase,
			OrganizationUsecase:        a.organizationUsecase,
			SiteUsecase:                a.siteUsecase,
			LocationUsecase:            a.locationUsecase,
			RobotUsecase:               a.robotUsecase,
			RobotDeviceUsecase:         a.robotDeviceUsecase,
			TaskUsecase:                a.taskUsecase,
			TaskVersionUsecase:         a.taskVersionUsecase,
			TaskTagUsecase:             a.taskTagUsecase,
			TaskImportUsecase:          a.taskImportUsecase,
			TaskExportUsecase:          a.taskExportUsecase,
			SubTaskUsecase:             a.subtaskUsecase,
			EpisodeUsecase:             a.episodeUsecase,
			EpisodeGradeUsecase:        a.episodeGradeUsecase,
			EpisodeExportUsecase:       a.episodeExportUsecase,
			EpisodeSubTaskUsecase:      a.episodeSubTaskUsecase,
			EpisodeExecutionUsecase:    a.episodeExecutionUsecase,
			FleetUsecase:               a.fleetUsecase,
			RobotOperatorUsecase:       a.robotOperatorUsecase,
			OperatorYieldExportUsecase: a.operatorYieldExportUsecase,
			APIKeyUsecase:              a.apiKeyUsecase,
		},
		Auth: httprouter.AuthDependencies{
			UserUsecase:   a.userUsecase,
			RobotUsecase:  a.robotUsecase,
			APIKeyUsecase: a.apiKeyUsecase,
		},
		SSE: httprouter.SSEDependencies{
			RobotDeviceUsecase:  a.robotDeviceUsecase,
			EpisodeUsecase:      a.episodeUsecase,
			TaskUsecase:         a.taskUsecase,
			TaskVersionUsecase:  a.taskVersionUsecase,
			EpisodeBus:          a.episodeBus,
			RobotEpisodeBus:     a.robotEpisodeBus,
			EpisodeListBus:      a.episodeListBus,
			RobotStatusEventBus: a.robotStatusBus,
		},
	})
}
