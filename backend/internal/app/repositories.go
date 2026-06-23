package app

import (
	"github.com/airoa-org/yubi-app/backend/internal/infra/cache"
	"github.com/airoa-org/yubi-app/backend/internal/infra/persistence"
	"github.com/airoa-org/yubi-app/backend/internal/infra/storage"
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
		User:                    persistence.NewUser(),
		UserLocation:            persistence.NewUserLocation(),
		UserSite:                persistence.NewUserSite(),
		Organization:            persistence.NewOrganization(),
		Site:                    persistence.NewSite(),
		Location:                persistence.NewLocation(),
		Robot:                   persistence.NewRobot(),
		Task:                    persistence.NewTask(),
		TaskTag:                 persistence.NewTaskTag(),
		TaskVersion:             persistence.NewTaskVersion(),
		SubTask:                 persistence.NewSubTask(),
		Episode:                 persistence.NewEpisode(),
		EpisodeGrade:            persistence.NewEpisodeGrade(),
		EpisodeSubTask:          persistence.NewEpisodeSubTask(),
		EpisodeSubTaskExecution: persistence.NewEpisodeSubTaskExecution(),
		APIKey:                  persistence.NewAPIKey(),
		OperatorYield:           persistence.NewOperatorYield(),
		Fleet:                   persistence.NewFleet(),
	}
}

func newRedisRepositories(redisClient *redis.Client) redisRepositories {
	return redisRepositories{
		RobotStatus:      cache.NewRobotStatus(redisClient),
		RobotUptimeDelta: cache.NewRobotUptimeDelta(redisClient),
		RobotOperator:    cache.NewRobotOperator(redisClient),
	}
}

type storageRepositories struct {
	EpisodeRecording repository.EpisodeRecording
}

func newStorageRepositories(s3Client *s3client.Client) storageRepositories {
	return storageRepositories{
		EpisodeRecording: storage.NewEpisodeRecording(s3Client),
	}
}
