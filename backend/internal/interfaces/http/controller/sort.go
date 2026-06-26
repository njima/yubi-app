package controller

import (
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func episodeSortBy(v *openapi.ListEpisodesParamsSortBy) *usecase.EpisodeSortBy {
	if v == nil {
		return nil
	}
	sortBy := usecase.EpisodeSortBy(*v)
	return &sortBy
}

func locationSortBy(v *openapi.ListLocationsParamsSortBy) *usecase.LocationSortBy {
	if v == nil {
		return nil
	}
	sortBy := usecase.LocationSortBy(*v)
	return &sortBy
}

func robotSortBy(v *openapi.ListRobotsParamsSortBy) *usecase.RobotSortBy {
	if v == nil {
		return nil
	}
	sortBy := usecase.RobotSortBy(*v)
	return &sortBy
}

func taskSortBy(v *openapi.ListTasksParamsSortBy) *usecase.TaskSortBy {
	if v == nil {
		return nil
	}
	sortBy := usecase.TaskSortBy(*v)
	return &sortBy
}

func userSortBy(v *openapi.ListUsersParamsSortBy) *usecase.UserSortBy {
	if v == nil {
		return nil
	}
	sortBy := usecase.UserSortBy(*v)
	return &sortBy
}

func sortOrder[T ~string](v *T) *usecase.SortOrder {
	if v == nil {
		return nil
	}
	order := usecase.SortOrder(*v)
	return &order
}
