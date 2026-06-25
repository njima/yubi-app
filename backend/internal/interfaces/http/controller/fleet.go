package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
)

func (c *controller) GetFleetSummary(ctx context.Context, request openapi.GetFleetSummaryRequestObject) (openapi.GetFleetSummaryResponseObject, error) {
	summaries, err := c.fleetUsecase.GetSummary(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]openapi.FleetSiteSummary, 0, len(summaries))
	for _, s := range summaries {
		rtMap := make(map[string]openapi.FleetRobotTypeSummary, len(s.RobotTypes))
		for robotTypeName, ms := range s.RobotTypes {
			fms := openapi.FleetRobotTypeSummary{
				Follower: openapi.FleetStatusCount{
					Operational: ms.Follower.Operational,
					Total:       ms.Follower.Total,
				},
			}
			if ms.Leader != nil {
				fms.Leader = &openapi.FleetStatusCount{
					Operational: ms.Leader.Operational,
					Total:       ms.Leader.Total,
				}
			}
			rtMap[robotTypeName] = fms
		}

		resp = append(resp, openapi.FleetSiteSummary{
			Site:       s.Site,
			SiteId:     s.SiteID,
			RobotTypes: rtMap,
		})
	}

	return openapi.GetFleetSummary200JSONResponse(resp), nil
}

func (c *controller) GetFleetStats(ctx context.Context, request openapi.GetFleetStatsRequestObject) (openapi.GetFleetStatsResponseObject, error) {
	stats, err := c.fleetUsecase.GetStats(ctx, request.Params.From, request.Params.To)
	if err != nil {
		return nil, err
	}

	resp := make([]openapi.FleetSiteStats, 0, len(stats))
	for _, s := range stats {
		rtList := make([]openapi.FleetRobotTypeStats, 0, len(s.RobotTypes))
		for _, ms := range s.RobotTypes {
			rtList = append(rtList, openapi.FleetRobotTypeStats{
				RobotType:          ms.RobotType,
				RobotUptime:        ms.RobotUptime,
				UptimeRate:         ms.UptimeRate,
				RobotCount:         ms.RobotCount,
				DataCollectionTime: ms.DataCollectionTime,
			})
		}
		resp = append(resp, openapi.FleetSiteStats{
			Site:       s.Site,
			SiteId:     s.SiteID,
			RobotTypes: rtList,
		})
	}

	return openapi.GetFleetStats200JSONResponse(resp), nil
}

func (c *controller) GetFleetCollectionTrend(ctx context.Context, request openapi.GetFleetCollectionTrendRequestObject) (openapi.GetFleetCollectionTrendResponseObject, error) {
	trend, err := c.fleetUsecase.GetCollectionTrend(ctx, fleetTrendGranularityModel(request.Params.Granularity), request.Params.From, request.Params.To)
	if err != nil {
		return nil, err
	}

	bySite := make([]openapi.TrendSeries, 0, len(trend.BySite))
	for _, s := range trend.BySite {
		bySite = append(bySite, openapi.TrendSeries{
			Label: s.Label,
			Data:  s.Data,
		})
	}

	byRobotType := make([]openapi.TrendSeries, 0, len(trend.ByRobotType))
	for _, s := range trend.ByRobotType {
		byRobotType = append(byRobotType, openapi.TrendSeries{
			Label: s.Label,
			Data:  s.Data,
		})
	}

	return openapi.GetFleetCollectionTrend200JSONResponse{
		Labels:      trend.Labels,
		BySite:      bySite,
		ByRobotType: byRobotType,
	}, nil
}
