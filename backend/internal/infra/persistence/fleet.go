package persistence

import (
	"context"
	"fmt"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/ccontext"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type fleet struct{}

func NewFleet() *fleet { return &fleet{} }

// GetSummary aggregates robot counts by site and robot type.
// episode_stats/robot_uptime tables don't have site_id,
// so we JOIN through location to reach site.
func (f *fleet) GetSummary(ctx context.Context, conn repository.DBConn) ([]repository.FleetSummaryRow, error) {
	orgID, err := ccontext.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "organization context required"))
	}

	var rows []repository.FleetSummaryRow

	err = conn.NewSelect().
		TableExpr("robot AS r").
		Join("JOIN location AS l ON l.id_natural = r.location_id").
		Join("JOIN site AS si ON si.id_natural = l.site_id").
		ColumnExpr("l.site_id AS site_id").
		ColumnExpr("si.name AS site_name").
		ColumnExpr("r.robot_type").
		ColumnExpr("r.status").
		ColumnExpr("r.leader_status").
		ColumnExpr("COUNT(*) AS count").
		Where("r.organization_id = ?", orgID).
		GroupExpr("l.site_id, si.name, r.robot_type, r.status, r.leader_status").
		Scan(ctx, &rows)

	if err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get fleet summary: %v", err))
	}

	return rows, nil
}

func (f *fleet) GetStats(ctx context.Context, conn repository.DBConn, filter repository.FleetStatsFilter) ([]repository.FleetStatsRow, error) {
	orgID, err := ccontext.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "organization context required"))
	}

	var rows []repository.FleetStatsRow

	err = conn.NewSelect().
		TableExpr("episode_stats_hourly AS esh").
		Join("JOIN robot AS r ON r.id_natural = esh.robot_id").
		Join("JOIN location AS l ON l.id_natural = esh.location_id").
		Join("JOIN site AS si ON si.id_natural = l.site_id").
		ColumnExpr("l.site_id AS site_id").
		ColumnExpr("si.name AS site_name").
		ColumnExpr("r.robot_type").
		ColumnExpr("SUM(esh.total_duration_seconds) AS total_duration_seconds").
		Where("esh.organization_id = ?", orgID).
		Where("esh.period_start >= ?", filter.From).
		Where("esh.period_start < ?", filter.To).
		GroupExpr("l.site_id, si.name, r.robot_type").
		Scan(ctx, &rows)

	if err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get fleet stats: %v", err))
	}

	return rows, nil
}

func (f *fleet) GetUptimeStats(ctx context.Context, conn repository.DBConn, filter repository.FleetStatsFilter) ([]repository.FleetUptimeStatsRow, error) {
	orgID, err := ccontext.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "organization context required"))
	}

	var rows []repository.FleetUptimeStatsRow

	err = conn.NewSelect().
		TableExpr("robot_uptime_hourly AS ruh").
		Join("JOIN robot AS r ON r.id_natural = ruh.robot_id").
		Join("JOIN location AS l ON l.id_natural = ruh.location_id").
		Join("JOIN site AS si ON si.id_natural = l.site_id").
		ColumnExpr("l.site_id AS site_id").
		ColumnExpr("si.name AS site_name").
		ColumnExpr("r.robot_type").
		ColumnExpr("SUM(ruh.uptime_seconds) AS uptime_seconds").
		ColumnExpr("COUNT(DISTINCT ruh.robot_id) AS robot_count").
		Where("ruh.organization_id = ?", orgID).
		Where("ruh.period_start >= ?", filter.From).
		Where("ruh.period_start < ?", filter.To).
		GroupExpr("l.site_id, si.name, r.robot_type").
		Scan(ctx, &rows)

	if err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get fleet uptime stats: %v", err))
	}

	return rows, nil
}

func (f *fleet) GetCollectionTrend(ctx context.Context, conn repository.DBConn, filter repository.FleetTrendFilter) ([]repository.FleetTrendRow, error) {
	tableName, err := statsTableForGranularity(filter.Granularity)
	if err != nil {
		return nil, err
	}

	orgID, orgErr := ccontext.OrganizationID(ctx)
	if orgErr != nil || orgID == "" {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "organization context required"))
	}

	var rows []repository.FleetTrendRow

	err = conn.NewSelect().
		TableExpr(fmt.Sprintf("%s AS es", tableName)).
		Join("JOIN robot AS r ON r.id_natural = es.robot_id").
		Join("JOIN location AS l ON l.id_natural = es.location_id").
		Join("JOIN site AS si ON si.id_natural = l.site_id").
		ColumnExpr("l.site_id AS site_id").
		ColumnExpr("si.name AS site_name").
		ColumnExpr("r.robot_type").
		ColumnExpr("es.period_start").
		ColumnExpr("es.total_duration_seconds").
		Where("es.organization_id = ?", orgID).
		Where("es.period_start >= ?", filter.From).
		Where("es.period_start < ?", filter.To).
		OrderExpr("es.period_start ASC").
		Scan(ctx, &rows)

	if err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get collection trend: %v", err))
	}

	return rows, nil
}

func statsTableForGranularity(granularity openapi.GetFleetCollectionTrendParamsGranularity) (string, error) {
	switch granularity {
	case openapi.Hourly:
		return "episode_stats_hourly", nil
	case openapi.Daily:
		return "episode_stats_daily", nil
	case openapi.Monthly:
		return "episode_stats_monthly", nil
	default:
		return "", apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "invalid granularity: %s", granularity))
	}
}
