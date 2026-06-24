package controller

import (
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

func episodeSortBy(v *openapi.ListEpisodesParamsSortBy) *repository.EpisodeSortBy {
	if v == nil {
		return nil
	}
	sortBy := repository.EpisodeSortBy(*v)
	return &sortBy
}

func locationSortBy(v *openapi.ListLocationsParamsSortBy) *repository.LocationSortBy {
	if v == nil {
		return nil
	}
	sortBy := repository.LocationSortBy(*v)
	return &sortBy
}

func robotSortBy(v *openapi.ListRobotsParamsSortBy) *repository.RobotSortBy {
	if v == nil {
		return nil
	}
	sortBy := repository.RobotSortBy(*v)
	return &sortBy
}

func taskSortBy(v *openapi.ListTasksParamsSortBy) *repository.TaskSortBy {
	if v == nil {
		return nil
	}
	sortBy := repository.TaskSortBy(*v)
	return &sortBy
}

func userSortBy(v *openapi.ListUsersParamsSortBy) *repository.UserSortBy {
	if v == nil {
		return nil
	}
	sortBy := repository.UserSortBy(*v)
	return &sortBy
}

func sortOrder[T ~string](v *T) *repository.SortOrder {
	if v == nil {
		return nil
	}
	order := repository.SortOrder(*v)
	return &order
}
