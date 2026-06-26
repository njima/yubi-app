package repository

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

type RobotListFilter struct {
	SiteID         *string
	LocationID     *string
	Status         *RobotFilterStatus
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
	Status         *RobotFilterStatus
	OnlineRobotIDs *[]string
	ExcludeOnline  bool
}

type Robot interface {
	Create(ctx context.Context, conn Conn, r model.Robot) (model.Robot, error)
	GetByID(ctx context.Context, conn Conn, id string) (model.Robot, error)
	List(ctx context.Context, conn Conn, filter RobotListFilter, limit, offset int) (model.Robots, int, error)
	ListTypes(ctx context.Context, conn Conn, filter RobotTypeFilter) ([]string, error)
	Update(ctx context.Context, conn Conn, r model.Robot) (model.Robot, error)
	Delete(ctx context.Context, conn Conn, id string) error
}
