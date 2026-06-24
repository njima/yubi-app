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
	Create(ctx context.Context, conn DBConn, r model.Robot) (model.Robot, error)
	GetByID(ctx context.Context, conn DBConn, id string) (model.Robot, error)
	List(ctx context.Context, conn DBConn, filter RobotListFilter, limit, offset int) (model.Robots, int, error)
	ListTypes(ctx context.Context, conn DBConn, filter RobotTypeFilter) ([]string, error)
	Update(ctx context.Context, conn DBConn, r model.Robot) (model.Robot, error)
	Delete(ctx context.Context, conn DBConn, id string) error
}
