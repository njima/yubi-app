package usecase

import (
	"context"
	"sort"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/uptrace/bun"
)

type FleetUsecase interface {
	GetSummary(ctx context.Context) ([]model.FleetSiteSummary, error)
	GetStats(ctx context.Context, from, to time.Time) ([]model.FleetSiteStats, error)
	GetCollectionTrend(ctx context.Context, granularity model.FleetTrendGranularity, from, to time.Time) (model.CollectionTrend, error)
}

type fleetUsecase struct {
	repo repository.Fleet
	db   *bun.DB
}

func NewFleet(repo repository.Fleet, db *bun.DB) *fleetUsecase {
	return &fleetUsecase{repo: repo, db: db}
}

func (f *fleetUsecase) GetSummary(ctx context.Context) ([]model.FleetSiteSummary, error) {
	rows, err := f.repo.GetSummary(ctx, f.db)
	if err != nil {
		return nil, err
	}

	return buildSummary(rows), nil
}

func (f *fleetUsecase) GetStats(ctx context.Context, from, to time.Time) ([]model.FleetSiteStats, error) {
	filter := repository.FleetStatsFilter{From: from, To: to}

	rows, err := f.repo.GetStats(ctx, f.db, filter)
	if err != nil {
		return nil, err
	}

	uptimeRows, err := f.repo.GetUptimeStats(ctx, f.db, filter)
	if err != nil {
		return nil, err
	}

	return buildStats(rows, uptimeRows, filter), nil
}

func (f *fleetUsecase) GetCollectionTrend(ctx context.Context, granularity model.FleetTrendGranularity, from, to time.Time) (model.CollectionTrend, error) {
	rows, err := f.repo.GetCollectionTrend(ctx, f.db, repository.FleetTrendFilter{
		Granularity: granularity,
		From:        from,
		To:          to,
	})
	if err != nil {
		return model.CollectionTrend{}, err
	}

	return buildCollectionTrend(rows, granularity), nil
}

// isFollowerOperational returns true if the robot follower status is operational.
// Online(0), Busy(1), Offline(2) = operational; Faulted(3), Maintenance(4) = not operational.
func isFollowerOperational(status model.RobotStatus) bool {
	return status != model.RobotStatusFaulted && status != model.RobotStatusMaintenance
}

// isLeaderOperational returns true if the robot leader status is operational.
// Ready(0) = operational; Faulted(1), Maintenance(2) = not operational.
func isLeaderOperational(status model.LeaderStatus) bool {
	return status == model.LeaderStatusReady
}

func buildSummary(rows []repository.FleetSummaryRow) []model.FleetSiteSummary {
	siteMap := make(map[string]*model.FleetSiteSummary)
	siteOrder := make([]string, 0)

	for _, row := range rows {
		site, exists := siteMap[row.SiteID]
		if !exists {
			site = &model.FleetSiteSummary{
				Site:       row.SiteName,
				SiteID:     row.SiteID,
				RobotTypes: make(map[string]model.FleetRobotTypeSummary),
			}
			siteMap[row.SiteID] = site
			siteOrder = append(siteOrder, row.SiteID)
		}

		robotType := row.RobotType
		if robotType == "" {
			robotType = "Unknown"
		}

		ms := site.RobotTypes[robotType]

		ms.Follower.Total += row.Count
		if isFollowerOperational(row.Status) {
			ms.Follower.Operational += row.Count
		}

		if row.LeaderStatus != nil {
			if ms.Leader == nil {
				ms.Leader = &model.FleetStatusCount{}
			}
			ms.Leader.Total += row.Count
			if isLeaderOperational(*row.LeaderStatus) {
				ms.Leader.Operational += row.Count
			}
		}

		site.RobotTypes[robotType] = ms
	}

	result := make([]model.FleetSiteSummary, 0, len(siteOrder))
	for _, siteID := range siteOrder {
		result = append(result, *siteMap[siteID])
	}

	return result
}

func buildStats(rows []repository.FleetStatsRow, uptimeRows []repository.FleetUptimeStatsRow, filter repository.FleetStatsFilter) []model.FleetSiteStats {
	periodSeconds := filter.To.Sub(filter.From).Seconds()

	// Build uptime lookup keyed by "siteID|robotType".
	// Requires an index on robot_uptime_hourly(period_start) for large datasets.
	type uptimeEntry struct {
		uptimeSeconds int64
		robotCount    int64
	}
	uptimeMap := make(map[string]uptimeEntry, len(uptimeRows))
	for _, row := range uptimeRows {
		robotType := row.RobotType
		if robotType == "" {
			robotType = "Unknown"
		}
		key := row.SiteID + "|" + robotType
		entry := uptimeMap[key]
		entry.uptimeSeconds += row.UptimeSeconds
		entry.robotCount += row.RobotCount
		uptimeMap[key] = entry
	}

	siteMap := make(map[string]*model.FleetSiteStats)
	siteOrder := make([]string, 0)

	for _, row := range rows {
		site, exists := siteMap[row.SiteID]
		if !exists {
			site = &model.FleetSiteStats{
				Site:       row.SiteName,
				SiteID:     row.SiteID,
				RobotTypes: make([]model.FleetRobotTypeStats, 0),
			}
			siteMap[row.SiteID] = site
			siteOrder = append(siteOrder, row.SiteID)
		}

		robotType := row.RobotType
		if robotType == "" {
			robotType = "Unknown"
		}

		var robotUptime *float32
		var uptimeRate *float32
		var robotCount *int
		if entry, ok := uptimeMap[row.SiteID+"|"+robotType]; ok {
			uptimeHours := float32(entry.uptimeSeconds) / 3600.0
			robotUptime = &uptimeHours
			if entry.robotCount > 0 && periodSeconds > 0 {
				rate := float32(float64(entry.uptimeSeconds) / (float64(entry.robotCount) * periodSeconds))
				if rate > 1.0 {
					// TODO: emit a metric/log here — raw rate > 1.0 indicates a data
					// anomaly in robot_uptime_hourly (duplicate writes or clock skew).
					rate = 1.0
				}
				uptimeRate = &rate
			}
			rc := int(entry.robotCount)
			robotCount = &rc
		}

		site.RobotTypes = append(site.RobotTypes, model.FleetRobotTypeStats{
			RobotType:          robotType,
			RobotUptime:        robotUptime,
			UptimeRate:         uptimeRate,
			RobotCount:         robotCount,
			DataCollectionTime: float32(row.TotalDurationSeconds) / 3600.0,
		})
	}

	result := make([]model.FleetSiteStats, 0, len(siteOrder))
	for _, siteID := range siteOrder {
		result = append(result, *siteMap[siteID])
	}

	return result
}

func buildCollectionTrend(rows []repository.FleetTrendRow, granularity model.FleetTrendGranularity) model.CollectionTrend {
	if len(rows) == 0 {
		return model.CollectionTrend{
			Labels:      []string{},
			BySite:      []model.TrendSeries{},
			ByRobotType: []model.TrendSeries{},
		}
	}

	// Collect unique labels (period_start) in order
	labelSet := make(map[string]bool)
	labels := make([]string, 0)
	for _, row := range rows {
		label := formatLabel(row.PeriodStart, granularity)
		if !labelSet[label] {
			labelSet[label] = true
			labels = append(labels, label)
		}
	}
	sort.Strings(labels)

	labelIndex := make(map[string]int, len(labels))
	for i, l := range labels {
		labelIndex[l] = i
	}

	// bySite: group by site_id, sum total_duration_seconds per label
	bySiteMap := make(map[string]*model.TrendSeries)
	bySiteOrder := make([]string, 0)
	// byRobotType: group by robot type, sum total_duration_seconds per label
	byRobotTypeMap := make(map[string]*model.TrendSeries)
	byRobotTypeOrder := make([]string, 0)

	for _, row := range rows {
		label := formatLabel(row.PeriodStart, granularity)
		idx := labelIndex[label]
		hours := float32(row.TotalDurationSeconds) / 3600.0

		// bySite
		siteSeries, exists := bySiteMap[row.SiteID]
		if !exists {
			siteSeries = &model.TrendSeries{
				Label: row.SiteName,
				Data:  make([]float32, len(labels)),
			}
			bySiteMap[row.SiteID] = siteSeries
			bySiteOrder = append(bySiteOrder, row.SiteID)
		}
		siteSeries.Data[idx] += hours

		// byModel
		robotType := row.RobotType
		if robotType == "" {
			robotType = "Unknown"
		}
		robotTypeSeries, exists := byRobotTypeMap[robotType]
		if !exists {
			robotTypeSeries = &model.TrendSeries{
				Label: robotType,
				Data:  make([]float32, len(labels)),
			}
			byRobotTypeMap[robotType] = robotTypeSeries
			byRobotTypeOrder = append(byRobotTypeOrder, robotType)
		}
		robotTypeSeries.Data[idx] += hours
	}

	bySite := make([]model.TrendSeries, 0, len(bySiteOrder))
	for _, siteID := range bySiteOrder {
		bySite = append(bySite, *bySiteMap[siteID])
	}

	byRobotType := make([]model.TrendSeries, 0, len(byRobotTypeOrder))
	for _, m := range byRobotTypeOrder {
		byRobotType = append(byRobotType, *byRobotTypeMap[m])
	}

	return model.CollectionTrend{
		Labels:      labels,
		BySite:      bySite,
		ByRobotType: byRobotType,
	}
}

func formatLabel(t time.Time, granularity model.FleetTrendGranularity) string {
	switch granularity {
	case model.FleetTrendGranularityHourly:
		return t.UTC().Format(time.RFC3339)
	case model.FleetTrendGranularityDaily:
		return t.UTC().Format("2006-01-02")
	case model.FleetTrendGranularityMonthly:
		return t.UTC().Format("2006-01")
	default:
		return t.UTC().Format(time.RFC3339)
	}
}
