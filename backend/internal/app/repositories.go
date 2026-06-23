package app

import (
	"github.com/airoa-org/yubi-app/backend/internal/gateway"
	"github.com/airoa-org/yubi-app/backend/internal/redis"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	s3client "github.com/airoa-org/yubi-app/backend/internal/s3"
)

type repositories struct {
	User                    repository.User
	UserLocation            repository.UserLocation
	UserSite                repository.UserSite
	Organization            repository.Organization
	Site                    repository.Site
	Location                repository.Location
	Robot                   repository.Robot
	Task                    repository.Task
	TaskTag                 repository.TaskTag
	TaskVersion             repository.TaskVersion
	SubTask                 repository.SubTask
	Episode                 repository.Episode
	EpisodeGrade            repository.EpisodeGrade
	EpisodeSubTask          repository.EpisodeSubTask
	EpisodeSubTaskExecution repository.EpisodeSubTaskExecution
	APIKey                  repository.APIKey
	RobotStatus             repository.RobotStatusRepository
	RobotUptimeDelta        repository.RobotUptimeDeltaRepository
	EpisodeRecording        repository.EpisodeRecording
	OperatorYield           repository.OperatorYield
	Fleet                   repository.Fleet
	RobotOperator           repository.RobotOperatorRepository
}

func newRepositories(redisClient *redis.Client, s3Client *s3client.Client) repositories {
	return repositories{
		User:                    gateway.NewUser(),
		UserLocation:            gateway.NewUserLocation(),
		UserSite:                gateway.NewUserSite(),
		Organization:            gateway.NewOrganization(),
		Site:                    gateway.NewSite(),
		Location:                gateway.NewLocation(),
		Robot:                   gateway.NewRobot(),
		Task:                    gateway.NewTask(),
		TaskTag:                 gateway.NewTaskTag(),
		TaskVersion:             gateway.NewTaskVersion(),
		SubTask:                 gateway.NewSubTask(),
		Episode:                 gateway.NewEpisode(),
		EpisodeGrade:            gateway.NewEpisodeGrade(),
		EpisodeSubTask:          gateway.NewEpisodeSubTask(),
		EpisodeSubTaskExecution: gateway.NewEpisodeSubTaskExecution(),
		APIKey:                  gateway.NewAPIKey(),
		RobotStatus:             gateway.NewRobotStatus(redisClient),
		RobotUptimeDelta:        gateway.NewRobotUptimeDelta(redisClient),
		EpisodeRecording:        gateway.NewEpisodeRecording(s3Client),
		OperatorYield:           gateway.NewOperatorYield(),
		Fleet:                   gateway.NewFleet(),
		RobotOperator:           gateway.NewRobotOperator(redisClient),
	}
}
