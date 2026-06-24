package model

type FleetTrendGranularity string

const (
	FleetTrendGranularityHourly  FleetTrendGranularity = "hourly"
	FleetTrendGranularityDaily   FleetTrendGranularity = "daily"
	FleetTrendGranularityMonthly FleetTrendGranularity = "monthly"
)

type FleetSiteSummary struct {
	Site       string
	SiteID     string
	RobotTypes map[string]FleetRobotTypeSummary
}

type FleetRobotTypeSummary struct {
	Leader   *FleetStatusCount
	Follower FleetStatusCount
}

type FleetStatusCount struct {
	Operational int
	Total       int
}

type FleetSiteStats struct {
	Site       string
	SiteID     string
	RobotTypes []FleetRobotTypeStats
}

type FleetRobotTypeStats struct {
	RobotType          string
	RobotUptime        *float32
	UptimeRate         *float32
	RobotCount         *int
	DataCollectionTime float32
}

type CollectionTrend struct {
	Labels      []string
	BySite      []TrendSeries
	ByRobotType []TrendSeries
}

type TrendSeries struct {
	Label string
	Data  []float32
}
