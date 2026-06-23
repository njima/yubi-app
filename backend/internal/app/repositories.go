package app

import (
	"github.com/airoa-org/yubi-app/backend/internal/gateway"
	"github.com/airoa-org/yubi-app/backend/internal/redis"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	s3client "github.com/airoa-org/yubi-app/backend/internal/s3"
)

type repositories struct {
	databaseRepositories
	redisRepositories
	storageRepositories
}

type databaseRepositories struct {
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
	OperatorYield           repository.OperatorYield
	Fleet                   repository.Fleet
}

type redisRepositories struct {
	RobotStatus      repository.RobotStatusRepository
	RobotUptimeDelta repository.RobotUptimeDeltaRepository
	RobotOperator    repository.RobotOperatorRepository
}

func newRepositories(redisClient *redis.Client, s3Client *s3client.Client) repositories {
	return repositories{
		databaseRepositories: newDatabaseRepositories(),
		redisRepositories:    newRedisRepositories(redisClient),
		storageRepositories:  newStorageRepositories(s3Client),
	}
}

func newDatabaseRepositories() databaseRepositories {
	return databaseRepositories{
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
		OperatorYield:           gateway.NewOperatorYield(),
		Fleet:                   gateway.NewFleet(),
	}
}

func newRedisRepositories(redisClient *redis.Client) redisRepositories {
	return redisRepositories{
		RobotStatus:      gateway.NewRobotStatus(redisClient),
		RobotUptimeDelta: gateway.NewRobotUptimeDelta(redisClient),
		RobotOperator:    gateway.NewRobotOperator(redisClient),
	}
}

type storageRepositories struct {
	EpisodeRecording repository.EpisodeRecording
}

func newStorageRepositories(s3Client *s3client.Client) storageRepositories {
	return storageRepositories{
		EpisodeRecording: gateway.NewEpisodeRecording(s3Client),
	}
}
