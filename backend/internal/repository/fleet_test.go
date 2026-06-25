package repository

import (
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

func TestFleetSummaryRowUsesDomainRobotStatuses(t *testing.T) {
	leaderStatus := model.LeaderStatusReady
	row := FleetSummaryRow{
		Status:       model.RobotStatusOnline,
		LeaderStatus: &leaderStatus,
	}

	var _ model.RobotStatus = row.Status
	var _ *model.LeaderStatus = row.LeaderStatus
}

func TestFleetTrendFilterUsesDomainGranularity(t *testing.T) {
	filter := FleetTrendFilter{
		Granularity: model.FleetTrendGranularityDaily,
	}

	var _ model.FleetTrendGranularity = filter.Granularity
}
