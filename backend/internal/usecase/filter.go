package usecase

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type EpisodeSortBy string
type LocationSortBy string
type RobotSortBy string
type TaskSortBy string
type UserSortBy string

type APIKeyListFilter struct {
	RobotID        *string
	UserID         *string
	IncludeRevoked bool
}

type EpisodeListFilter struct {
	TaskID        *string
	TaskVersionID *string
	RobotID       *string
	UserID        *string
	Statuses      []model.EpisodeStatus
	StartedAtFrom *time.Time
	StartedAtTo   *time.Time
	SortBy        *EpisodeSortBy
	SortOrder     *SortOrder
}

type EpisodeExportFilter struct {
	EpisodeListFilter
}

type LocationListFilter struct {
	SiteID    *string
	Search    *string
	SortBy    *LocationSortBy
	SortOrder *SortOrder
}

type OperatorYieldExportFilter struct {
	DateFrom   time.Time
	DateTo     time.Time
	LocationID *string
	TaskID     *string
	UserID     *string
}

type RobotListFilter struct {
	SiteID         *string
	LocationID     *string
	Status         *model.RobotStatus
	RobotType      *string
	Search         *string
	OnlineRobotIDs *[]string
	ExcludeOnline  bool
	SortBy         *RobotSortBy
	SortOrder      *SortOrder
}

type RobotTypeFilter struct {
	SiteID         *string
	LocationID     *string
	Status         *model.RobotStatus
	OnlineRobotIDs *[]string
	ExcludeOnline  bool
}

type SiteListFilter struct {
	OrganizationID *string
	Search         *string
}

type TaskListFilter struct {
	HasApprovedVersion *bool
	SortBy             *TaskSortBy
	SortOrder          *SortOrder
	Statuses           []model.TaskStatus
	Priorities         []model.TaskPriority
	Difficulties       []model.TaskDifficulty
	RobotType          *string
	Search             *string
}

type TaskSummaryFilter struct {
	RobotTypes     []string
	CategoryTypeID *string
	TagIDs         []string
	DeadlineFrom   *time.Time
	DeadlineTo     *time.Time
}

type UserListFilter struct {
	LocationID *string
	SiteID     *string
	Search     *string
	SortBy     *UserSortBy
	SortOrder  *SortOrder
}

func (f APIKeyListFilter) repositoryFilter() repository.APIKeyListFilter {
	return repository.APIKeyListFilter(f)
}

func (f EpisodeListFilter) repositoryFilter() repository.EpisodeListFilter {
	return repository.EpisodeListFilter{
		TaskID:        f.TaskID,
		TaskVersionID: f.TaskVersionID,
		RobotID:       f.RobotID,
		UserID:        f.UserID,
		Statuses:      f.Statuses,
		StartedAtFrom: f.StartedAtFrom,
		StartedAtTo:   f.StartedAtTo,
		SortBy:        episodeSortByRepository(f.SortBy),
		SortOrder:     sortOrderRepository(f.SortOrder),
	}
}

func (f EpisodeExportFilter) repositoryFilter() repository.EpisodeExportFilter {
	return repository.EpisodeExportFilter{
		EpisodeListFilter: f.EpisodeListFilter.repositoryFilter(),
	}
}

func (f LocationListFilter) repositoryFilter() repository.LocationListFilter {
	return repository.LocationListFilter{
		SiteID:    f.SiteID,
		Search:    f.Search,
		SortBy:    locationSortByRepository(f.SortBy),
		SortOrder: sortOrderRepository(f.SortOrder),
	}
}

func (f OperatorYieldExportFilter) repositoryFilter() repository.OperatorYieldExportFilter {
	return repository.OperatorYieldExportFilter(f)
}

func (f RobotListFilter) repositoryFilter() repository.RobotListFilter {
	return repository.RobotListFilter{
		SiteID:         f.SiteID,
		LocationID:     f.LocationID,
		Status:         robotStatusRepository(f.Status),
		RobotType:      f.RobotType,
		Search:         f.Search,
		OnlineRobotIDs: f.OnlineRobotIDs,
		ExcludeOnline:  f.ExcludeOnline,
		SortBy:         robotSortByRepository(f.SortBy),
		SortOrder:      sortOrderRepository(f.SortOrder),
	}
}

func (f RobotTypeFilter) repositoryFilter() repository.RobotTypeFilter {
	return repository.RobotTypeFilter{
		SiteID:         f.SiteID,
		LocationID:     f.LocationID,
		Status:         robotStatusRepository(f.Status),
		OnlineRobotIDs: f.OnlineRobotIDs,
		ExcludeOnline:  f.ExcludeOnline,
	}
}

func (f SiteListFilter) repositoryFilter() repository.SiteListFilter {
	return repository.SiteListFilter(f)
}

func (f TaskListFilter) repositoryFilter() repository.TaskListFilter {
	return repository.TaskListFilter{
		HasApprovedVersion: f.HasApprovedVersion,
		SortBy:             taskSortByRepository(f.SortBy),
		SortOrder:          sortOrderRepository(f.SortOrder),
		Statuses:           f.Statuses,
		Priorities:         f.Priorities,
		Difficulties:       f.Difficulties,
		RobotType:          f.RobotType,
		Search:             f.Search,
	}
}

func (f TaskSummaryFilter) repositoryFilter() repository.TaskSummaryFilter {
	return repository.TaskSummaryFilter(f)
}

func (f UserListFilter) repositoryFilter() repository.UserListFilter {
	return repository.UserListFilter{
		LocationID: f.LocationID,
		SiteID:     f.SiteID,
		Search:     f.Search,
		SortBy:     userSortByRepository(f.SortBy),
		SortOrder:  sortOrderRepository(f.SortOrder),
	}
}

func sortOrderRepository(v *SortOrder) *repository.SortOrder {
	if v == nil {
		return nil
	}
	order := repository.SortOrder(*v)
	return &order
}

func episodeSortByRepository(v *EpisodeSortBy) *repository.EpisodeSortBy {
	if v == nil {
		return nil
	}
	sortBy := repository.EpisodeSortBy(*v)
	return &sortBy
}

func locationSortByRepository(v *LocationSortBy) *repository.LocationSortBy {
	if v == nil {
		return nil
	}
	sortBy := repository.LocationSortBy(*v)
	return &sortBy
}

func robotSortByRepository(v *RobotSortBy) *repository.RobotSortBy {
	if v == nil {
		return nil
	}
	sortBy := repository.RobotSortBy(*v)
	return &sortBy
}

func taskSortByRepository(v *TaskSortBy) *repository.TaskSortBy {
	if v == nil {
		return nil
	}
	sortBy := repository.TaskSortBy(*v)
	return &sortBy
}

func userSortByRepository(v *UserSortBy) *repository.UserSortBy {
	if v == nil {
		return nil
	}
	sortBy := repository.UserSortBy(*v)
	return &sortBy
}

func robotStatusRepository(v *model.RobotStatus) *repository.RobotFilterStatus {
	if v == nil {
		return nil
	}
	status := repository.RobotFilterStatus(*v)
	return &status
}
