package controller

import (
	"encoding/json"

	"github.com/rs/zerolog"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

type controller struct {
	logger                     zerolog.Logger
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

type Dependencies struct {
	Logger                     zerolog.Logger
	UserUsecase                usecase.UserUsecase
	UserImportUsecase          usecase.UserImportUsecase
	OrganizationUsecase        usecase.OrganizationUsecase
	SiteUsecase                usecase.SiteUsecase
	LocationUsecase            usecase.LocationUsecase
	RobotUsecase               usecase.RobotUsecase
	RobotDeviceUsecase         usecase.RobotDeviceUsecase
	TaskUsecase                usecase.TaskUsecase
	TaskVersionUsecase         usecase.TaskVersionUsecase
	TaskTagUsecase             usecase.TaskTagUsecase
	TaskImportUsecase          usecase.TaskImportUsecase
	TaskExportUsecase          usecase.TaskExportUsecase
	SubTaskUsecase             usecase.SubTaskUsecase
	EpisodeUsecase             usecase.EpisodeUsecase
	EpisodeGradeUsecase        usecase.EpisodeGradeUsecase
	EpisodeExportUsecase       usecase.EpisodeExportUsecase
	EpisodeSubTaskUsecase      usecase.EpisodeSubTaskUsecase
	EpisodeExecutionUsecase    usecase.EpisodeExecutionUsecase
	FleetUsecase               usecase.FleetUsecase
	RobotOperatorUsecase       usecase.RobotOperatorUsecase
	OperatorYieldExportUsecase usecase.OperatorYieldExportUsecase
	APIKeyUsecase              usecase.APIKeyUsecase
}

// Verify that controller implements StrictServerInterface
var _ openapi.StrictServerInterface = (*controller)(nil)

func NewController(deps Dependencies) *controller {
	return &controller{
		logger:                     deps.Logger,
		userUsecase:                deps.UserUsecase,
		userImportUsecase:          deps.UserImportUsecase,
		organizationUsecase:        deps.OrganizationUsecase,
		siteUsecase:                deps.SiteUsecase,
		locationUsecase:            deps.LocationUsecase,
		robotUsecase:               deps.RobotUsecase,
		robotDeviceUsecase:         deps.RobotDeviceUsecase,
		taskUsecase:                deps.TaskUsecase,
		taskVersionUsecase:         deps.TaskVersionUsecase,
		taskTagUsecase:             deps.TaskTagUsecase,
		taskImportUsecase:          deps.TaskImportUsecase,
		taskExportUsecase:          deps.TaskExportUsecase,
		subtaskUsecase:             deps.SubTaskUsecase,
		episodeUsecase:             deps.EpisodeUsecase,
		episodeGradeUsecase:        deps.EpisodeGradeUsecase,
		episodeExportUsecase:       deps.EpisodeExportUsecase,
		episodeSubTaskUsecase:      deps.EpisodeSubTaskUsecase,
		episodeExecutionUsecase:    deps.EpisodeExecutionUsecase,
		fleetUsecase:               deps.FleetUsecase,
		robotOperatorUsecase:       deps.RobotOperatorUsecase,
		operatorYieldExportUsecase: deps.OperatorYieldExportUsecase,
		apiKeyUsecase:              deps.APIKeyUsecase,
	}
}

func mapPtrFromRawMessagePtr(b *json.RawMessage) *map[string]interface{} {
	if b == nil || len(*b) == 0 {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(*b, &m); err != nil {
		return nil
	}
	return &m
}
