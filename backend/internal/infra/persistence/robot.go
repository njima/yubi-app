package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/uptrace/bun"
)

type robot struct{}

func NewRobot() *robot { return &robot{} }

func leaderStatusToModel(ls *uint) *model.LeaderStatus {
	if ls == nil {
		return nil
	}
	v := model.LeaderStatus(*ls)
	return &v
}

func leaderStatusToEntity(ls *model.LeaderStatus) *uint {
	if ls == nil {
		return nil
	}
	v := uint(*ls)
	return &v
}

func (r *robot) Create(ctx context.Context, conn repository.DBConn, rob model.Robot) (model.Robot, error) {
	var inserted entity.Robot
	dbRob := entity.Robot{
		IDNatural:            rob.IDNatural,
		OrganizationID:       rob.OrganizationID,
		LocationID:           rob.LocationID,
		Name:                 rob.Name,
		RobotType:            *rob.RobotType,
		Status:               uint(rob.Status),
		LeaderStatus:         leaderStatusToEntity(rob.LeaderStatus),
		LeaderFaultStartedAt: rob.LeaderFaultStartedAt,
		FaultStartedAt:       rob.FaultStartedAt,
		OfflineReason:        rob.OfflineReason,
		RobotConfig:          rob.RobotConfig,
	}

	if err := conn.NewInsert().
		Model(&dbRob).
		Returning("id_natural").
		Scan(ctx, &inserted); err != nil {
		return model.Robot{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create robot: %v", err))
	}

	// Re-fetch with relations to include Site info
	return r.GetByID(ctx, conn, inserted.IDNatural)
}

func (r *robot) GetByID(ctx context.Context, conn repository.DBConn, id string) (model.Robot, error) {
	var dbRob entity.Robot
	if err := conn.NewSelect().
		Model(&dbRob).
		Relation("Organization").
		Relation("Location").
		Relation("Location.Site").
		Where("r.id_natural = ?", id).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.Robot{}, apperror.NewError(apperror.NewMessage(apperror.CodeRobotNotFound, "robot not found: id_natural=%s", id))
		}
		return model.Robot{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get robot: %v", err))
	}

	orgName := ""
	if dbRob.Organization != nil {
		orgName = dbRob.Organization.Name
	}
	locName := ""
	siteID := ""
	siteName := ""
	if dbRob.Location != nil {
		locName = dbRob.Location.Name
		siteID = dbRob.Location.SiteID
		if dbRob.Location.Site != nil {
			siteName = dbRob.Location.Site.Name
		}
	}

	rob := model.Robot{
		ID:                   dbRob.ID,
		IDNatural:            dbRob.IDNatural,
		OrganizationID:       dbRob.OrganizationID,
		OrganizationName:     orgName,
		LocationID:           dbRob.LocationID,
		LocationName:         locName,
		SiteID:               siteID,
		SiteName:             siteName,
		Name:                 dbRob.Name,
		RobotType:            &dbRob.RobotType,
		Status:               model.RobotStatus(dbRob.Status),
		LeaderStatus:         leaderStatusToModel(dbRob.LeaderStatus),
		LeaderFaultStartedAt: dbRob.LeaderFaultStartedAt,
		FaultStartedAt:       dbRob.FaultStartedAt,
		LastHeartbeatAt:      dbRob.LastHeartbeatAt,
		OfflineReason:        dbRob.OfflineReason,
		RobotConfig:          dbRob.RobotConfig,
		ActiveEpisodeID:      dbRob.ActiveEpisodeID,
		ActiveUserID:         dbRob.ActiveUserID,
		CreatedAt:            dbRob.CreatedAt,
		UpdatedAt:            &dbRob.UpdatedAt,
	}

	return rob, nil
}

func (r *robot) List(ctx context.Context, conn repository.DBConn, filter repository.RobotListFilter, limit, offset int) (model.Robots, int, error) {
	var dbRobs []entity.Robot

	sel := conn.NewSelect().
		Model(&dbRobs).
		Relation("Organization").
		Relation("Location").
		Relation("Location.Site").
		Join("LEFT JOIN location AS l ON l.id_natural = r.location_id").
		Join("LEFT JOIN \"user\" AS u ON u.id_natural = r.active_user_id").
		Limit(limit).
		Offset(offset)

	if filter.SiteID != nil && *filter.SiteID != "" {
		sel = sel.Where(`EXISTS (
			SELECT 1 FROM location l
			WHERE l.id_natural = r.location_id AND l.site_id = ?
		)`, *filter.SiteID)
	}
	if filter.LocationID != nil {
		sel = sel.Where("r.location_id = ?", *filter.LocationID)
	}
	if filter.Status != nil {
		sel = sel.Where("r.status = ?", *filter.Status)
	}
	if filter.RobotType != nil {
		sel = sel.Where("r.robot_type = ?", *filter.RobotType)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		sel = sel.Where("r.name ILIKE ?", "%"+escaped+"%")
	}
	if filter.OnlineRobotIDs != nil {
		ids := *filter.OnlineRobotIDs
		sel = sel.Where("r.status IN (?)", bun.In([]model.RobotStatus{
			model.RobotStatusReady, model.RobotStatusOnline,
		}))
		if filter.ExcludeOnline {
			if len(ids) > 0 {
				sel = sel.Where("r.id_natural NOT IN (?)", bun.In(ids))
			}
		} else {
			sel = sel.Where("r.id_natural IN (?)", bun.In(ids))
		}
	}

	sel = applyRobotSortOrder(sel, filter.SortBy, filter.SortOrder)

	if err := sel.Scan(ctx); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list robots: %v", err))
	}

	var total int
	countSel := conn.NewSelect().
		Model((*entity.Robot)(nil)).
		ColumnExpr("COUNT(*)")

	if filter.SiteID != nil && *filter.SiteID != "" {
		countSel = countSel.Where(`EXISTS (
			SELECT 1 FROM location l
			WHERE l.id_natural = r.location_id AND l.site_id = ?
		)`, *filter.SiteID)
	}
	if filter.LocationID != nil {
		countSel = countSel.Where("r.location_id = ?", *filter.LocationID)
	}
	if filter.Status != nil {
		countSel = countSel.Where("r.status = ?", *filter.Status)
	}
	if filter.RobotType != nil {
		countSel = countSel.Where("r.robot_type = ?", *filter.RobotType)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		countSel = countSel.Where("r.name ILIKE ?", "%"+escaped+"%")
	}
	if filter.OnlineRobotIDs != nil {
		ids := *filter.OnlineRobotIDs
		countSel = countSel.Where("r.status IN (?)", bun.In([]model.RobotStatus{
			model.RobotStatusReady, model.RobotStatusOnline,
		}))
		if filter.ExcludeOnline {
			if len(ids) > 0 {
				countSel = countSel.Where("r.id_natural NOT IN (?)", bun.In(ids))
			}
		} else {
			countSel = countSel.Where("r.id_natural IN (?)", bun.In(ids))
		}
	}

	if err := countSel.Scan(ctx, &total); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count robots: %v", err))
	}

	res := make(model.Robots, 0, len(dbRobs))
	for _, dr := range dbRobs {
		orgN := ""
		if dr.Organization != nil {
			orgN = dr.Organization.Name
		}
		locName := ""
		siteID := ""
		siteName := ""
		if dr.Location != nil {
			locName = dr.Location.Name
			siteID = dr.Location.SiteID
			if dr.Location.Site != nil {
				siteName = dr.Location.Site.Name
			}
		}
		m := &model.Robot{
			ID:                   dr.ID,
			IDNatural:            dr.IDNatural,
			OrganizationID:       dr.OrganizationID,
			OrganizationName:     orgN,
			LocationID:           dr.LocationID,
			LocationName:         locName,
			SiteID:               siteID,
			SiteName:             siteName,
			Name:                 dr.Name,
			RobotType:            &dr.RobotType,
			Status:               model.RobotStatus(dr.Status),
			LeaderStatus:         leaderStatusToModel(dr.LeaderStatus),
			LeaderFaultStartedAt: dr.LeaderFaultStartedAt,
			FaultStartedAt:       dr.FaultStartedAt,
			LastHeartbeatAt:      dr.LastHeartbeatAt,
			OfflineReason:        dr.OfflineReason,
			RobotConfig:          dr.RobotConfig,
			ActiveEpisodeID:      dr.ActiveEpisodeID,
			ActiveUserID:         dr.ActiveUserID,
			CreatedAt:            dr.CreatedAt,
		}
		if !dr.UpdatedAt.IsZero() {
			t := dr.UpdatedAt
			m.UpdatedAt = &t
		}
		res = append(res, m)
	}

	return res, total, nil
}

func (r *robot) ListTypes(ctx context.Context, conn repository.DBConn, filter repository.RobotTypeFilter) ([]string, error) {
	var types []string
	sel := conn.NewSelect().
		Model((*entity.Robot)(nil)).
		ColumnExpr("DISTINCT robot_type").
		Where("robot_type != ''")

	if filter.SiteID != nil && *filter.SiteID != "" {
		sel = sel.Where(`EXISTS (
			SELECT 1 FROM location l
			WHERE l.id_natural = r.location_id AND l.site_id = ?
		)`, *filter.SiteID)
	}
	if filter.LocationID != nil {
		sel = sel.Where("r.location_id = ?", *filter.LocationID)
	}
	if filter.Status != nil {
		sel = sel.Where("r.status = ?", *filter.Status)
	}
	if filter.OnlineRobotIDs != nil {
		ids := *filter.OnlineRobotIDs
		sel = sel.Where("r.status IN (?)", bun.In([]model.RobotStatus{
			model.RobotStatusReady, model.RobotStatusOnline,
		}))
		if filter.ExcludeOnline {
			if len(ids) > 0 {
				sel = sel.Where("r.id_natural NOT IN (?)", bun.In(ids))
			}
		} else {
			sel = sel.Where("r.id_natural IN (?)", bun.In(ids))
		}
	}

	sel = sel.OrderExpr("robot_type ASC")

	if err := sel.Scan(ctx, &types); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list robot types: %v", err))
	}
	return types, nil
}

func (r *robot) Update(ctx context.Context, conn repository.DBConn, rob model.Robot) (model.Robot, error) {
	robotTypeStr := ""
	if rob.RobotType != nil {
		robotTypeStr = *rob.RobotType
	}

	dbRob := entity.Robot{
		IDNatural:            rob.IDNatural,
		OrganizationID:       rob.OrganizationID,
		LocationID:           rob.LocationID,
		Name:                 rob.Name,
		RobotType:            robotTypeStr,
		Status:               uint(rob.Status),
		LeaderStatus:         leaderStatusToEntity(rob.LeaderStatus),
		LeaderFaultStartedAt: rob.LeaderFaultStartedAt,
		FaultStartedAt:       rob.FaultStartedAt,
		LastHeartbeatAt:      rob.LastHeartbeatAt,
		OfflineReason:        rob.OfflineReason,
		RobotConfig:          rob.RobotConfig,
		ActiveEpisodeID:      rob.ActiveEpisodeID,
		ActiveUserID:         rob.ActiveUserID,
	}

	var updated entity.Robot
	if err := conn.NewUpdate().
		Model(&dbRob).
		Where("id_natural = ?", rob.IDNatural).
		ExcludeColumn("id", "id_natural", "organization_id", "created_at").
		Returning("id_natural").
		Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.Robot{}, apperror.NewError(apperror.NewMessage(apperror.CodeRobotNotFound, "robot not found: id_natural=%s", rob.IDNatural))
		}
		return model.Robot{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update robot: %v", err))
	}

	// Re-fetch with relations to include Site info
	return r.GetByID(ctx, conn, updated.IDNatural)
}

func (r *robot) Delete(ctx context.Context, conn repository.DBConn, id string) error {
	var deletedID int64
	if err := conn.NewDelete().
		Model((*entity.Robot)(nil)).
		Where("id_natural = ?", id).
		Returning("id").
		Scan(ctx, &deletedID); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewError(apperror.NewMessage(apperror.CodeRobotNotFound, "robot not found: id_natural=%s", id))
		}
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete robot: %v", err))
	}
	return nil
}

var allowedRobotSortColumns = map[string]string{
	"name":              "r.name",
	"location_id":       "COALESCE(l.name, r.location_id)",
	"robot_type":        "r.robot_type",
	"status":            "r.status",
	"leader_status":     "r.leader_status",
	"last_heartbeat_at": "r.last_heartbeat_at",
	"active_episode_id": "r.active_episode_id",
	"active_user_id":    "COALESCE(u.name, r.active_user_id)",
}

var nullableRobotSortColumns = map[string]bool{
	"robot_type":        true,
	"leader_status":     true,
	"last_heartbeat_at": true,
	"active_episode_id": true,
	"active_user_id":    true,
}

func applyRobotSortOrder(sel *bun.SelectQuery, sortBy *repository.RobotSortBy, sortOrder *repository.SortOrder) *bun.SelectQuery {
	if sortBy == nil {
		return sel.OrderExpr("r.created_at DESC, r.id DESC")
	}

	col, ok := allowedRobotSortColumns[string(*sortBy)]
	if !ok {
		return sel.OrderExpr("r.created_at DESC, r.id DESC")
	}

	order := "ASC"
	if sortOrder != nil && *sortOrder == repository.SortOrderDesc {
		order = "DESC"
	}

	nullsClause := ""
	if nullableRobotSortColumns[string(*sortBy)] {
		nullsClause = " NULLS LAST"
	}

	return sel.OrderExpr(fmt.Sprintf("%s %s%s, r.created_at DESC, r.id DESC", col, order, nullsClause))
}
