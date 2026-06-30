package main

import (
	"github.com/airoa-org/yubi-app/backend/internal/infra/cache"
	"github.com/airoa-org/yubi-app/backend/internal/infra/persistence"
	"github.com/airoa-org/yubi-app/backend/internal/infra/storage"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
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
	OrganizationMembership  repository.OrganizationMembership
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

func newRepositories(redisClient *cache.Client, s3Client *storage.Client) repositories {
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
		OrganizationMembership:  persistence.NewOrganizationMembership(),
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

func newRedisRepositories(redisClient *cache.Client) redisRepositories {
	return redisRepositories{
		RobotStatus:      cache.NewRobotStatus(redisClient),
		RobotUptimeDelta: cache.NewRobotUptimeDelta(redisClient),
		RobotOperator:    cache.NewRobotOperator(redisClient),
	}
}

type storageRepositories struct {
	EpisodeRecording repository.EpisodeRecording
}

func newStorageRepositories(s3Client *storage.Client) storageRepositories {
	return storageRepositories{
		EpisodeRecording: storage.NewEpisodeRecording(s3Client),
	}
}
