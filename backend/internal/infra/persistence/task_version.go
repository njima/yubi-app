package persistence

import (
	"context"
	"database/sql"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/ccontext"
	"github.com/airoa-org/yubi-app/backend/internal/database/bunconv"
	"github.com/airoa-org/yubi-app/backend/internal/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/uptrace/bun"
)

type taskVersion struct{}

func NewTaskVersion() repository.TaskVersion {
	return &taskVersion{}
}

func (tv *taskVersion) Create(ctx context.Context, conn repository.DBConn, m model.TaskVersion) (model.TaskVersion, error) {
	var inserted entity.TaskVersion
	e := entity.TaskVersion{
		IDNatural:                       m.IDNatural,
		OrganizationID:                  m.OrganizationID,
		TaskID:                          m.TaskID,
		Version:                         m.Version,
		DisplayName:                     m.DisplayName,
		IsActive:                        true,
		ApprovalStatus:                  m.ApprovalStatus,
		TargetDurationSeconds:           m.TargetDurationSeconds,
		TargetEpisodeCount:              m.TargetEpisodeCount,
		TargetDurationPerEpisodeSeconds: m.TargetDurationPerEpisodeSeconds,
		Parameters:                      bunconv.TaskVersionParametersToJSON(m.Parameters),
	}

	if err := conn.NewInsert().
		Model(&e).
		Returning("*").
		Scan(ctx, &inserted); err != nil {
		return model.TaskVersion{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create task version: %v", err))
	}

	return *bunconv.EntityToTaskVersionModel(inserted), nil
}

func (tv *taskVersion) Update(ctx context.Context, conn repository.DBConn, m model.TaskVersion) (model.TaskVersion, error) {
	upd := conn.NewUpdate().Model((*entity.TaskVersion)(nil))
	hasSet := false
	if m.TargetDurationSeconds != nil {
		if *m.TargetDurationSeconds == 0 {
			upd = upd.Set("target_duration_seconds = NULL")
		} else {
			upd = upd.Set("target_duration_seconds = ?", *m.TargetDurationSeconds)
		}
		hasSet = true
	}
	if m.TargetEpisodeCount != nil {
		if *m.TargetEpisodeCount == 0 {
			upd = upd.Set("target_episode_count = NULL")
		} else {
			upd = upd.Set("target_episode_count = ?", *m.TargetEpisodeCount)
		}
		hasSet = true
	}
	if m.TargetDurationPerEpisodeSeconds != nil {
		if *m.TargetDurationPerEpisodeSeconds == 0 {
			upd = upd.Set("target_duration_per_episode_seconds = NULL")
		} else {
			upd = upd.Set("target_duration_per_episode_seconds = ?", *m.TargetDurationPerEpisodeSeconds)
		}
		hasSet = true
	}
	if m.DisplayName != nil {
		if *m.DisplayName == "" {
			upd = upd.Set("display_name = NULL")
		} else {
			upd = upd.Set("display_name = ?", *m.DisplayName)
		}
		hasSet = true
	}
	if !hasSet {
		return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "no fields to update"))
	}

	var updated entity.TaskVersion
	if err := upd.Where("id_natural = ?", m.IDNatural).Returning("*").Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskVersionNotFound, "task version not found: id_natural=%s", m.IDNatural))
		}
		return model.TaskVersion{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update task version: %v", err))
	}
	return *bunconv.EntityToTaskVersionModel(updated), nil
}

func (tv *taskVersion) GetByID(ctx context.Context, conn repository.DBConn, id string) (model.TaskVersion, error) {
	var e entity.TaskVersion
	if err := conn.NewSelect().
		Model(&e).
		Where("id_natural = ?", id).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskVersionNotFound, "task version not found: id_natural=%s", id))
		}
		return model.TaskVersion{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task version: %v", err))
	}
	m := bunconv.EntityToTaskVersionModel(e)
	if err := tv.attachStats(ctx, conn, []*model.TaskVersion{m}); err != nil {
		return model.TaskVersion{}, err
	}
	return *m, nil
}

func (tv *taskVersion) GetByIDForUpdate(ctx context.Context, conn repository.DBConn, id string) (model.TaskVersion, error) {
	var e entity.TaskVersion
	if err := conn.NewSelect().
		Model(&e).
		Where("id_natural = ?", id).
		For("UPDATE").
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskVersionNotFound, "task version not found: id_natural=%s", id))
		}
		return model.TaskVersion{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task version for update: %v", err))
	}
	return *bunconv.EntityToTaskVersionModel(e), nil
}

func (tv *taskVersion) GetLatestApprovedByTaskID(ctx context.Context, conn repository.DBConn, taskID string) (model.TaskVersion, error) {
	var e entity.TaskVersion
	if err := conn.NewSelect().
		Model(&e).
		Where("task_id = ?", taskID).
		Where("approval_status = ?", model.ApprovalStatusApproved).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskVersionNotApproved, "no approved task version found for task: id=%s", taskID))
		}
		return model.TaskVersion{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get latest approved task version: %v", err))
	}
	return *bunconv.EntityToTaskVersionModel(e), nil
}

func (tv *taskVersion) Approve(ctx context.Context, conn repository.DBConn, id string) (model.TaskVersion, error) {
	var updated entity.TaskVersion
	if err := conn.NewUpdate().
		Model((*entity.TaskVersion)(nil)).
		Set("approval_status = ?", model.ApprovalStatusApproved).
		Where("id_natural = ?", id).
		Returning("*").
		Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskVersionNotFound, "task version not found: id_natural=%s", id))
		}
		return model.TaskVersion{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to approve task version: %v", err))
	}
	return *bunconv.EntityToTaskVersionModel(updated), nil
}

func (tv *taskVersion) UpdateParameters(ctx context.Context, conn repository.DBConn, id string, parameters []model.TaskVersionParameter) (model.TaskVersion, error) {
	var updated entity.TaskVersion
	if err := conn.NewUpdate().
		Model((*entity.TaskVersion)(nil)).
		Set("parameters = ?", bunconv.TaskVersionParametersToJSON(parameters)).
		Where("id_natural = ?", id).
		Returning("*").
		Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.TaskVersion{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskVersionNotFound, "task version not found: id_natural=%s", id))
		}
		return model.TaskVersion{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update task version parameters: %v", err))
	}
	return *bunconv.EntityToTaskVersionModel(updated), nil
}

func (tv *taskVersion) ListByIDs(ctx context.Context, conn repository.DBConn, ids []string) (model.TaskVersions, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var versions []entity.TaskVersion
	if err := conn.NewSelect().
		Model(&versions).
		Where("id_natural IN (?)", bun.In(ids)).
		Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list task versions by IDs: %v", err))
	}
	return bunconv.EntitiesToTaskVersionModels(versions), nil
}

func (tv *taskVersion) ListByTaskID(ctx context.Context, conn repository.DBConn, taskID string) (model.TaskVersions, error) {
	var versions []entity.TaskVersion
	err := conn.NewSelect().
		Model(&versions).
		Where("task_id = ?", taskID).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list task versions: %v", err))
	}
	models := bunconv.EntitiesToTaskVersionModels(versions)
	if err := tv.attachStats(ctx, conn, models); err != nil {
		return nil, err
	}
	return models, nil
}

// attachStats fetches task_version_stats and attaches actual values to models.
func (tv *taskVersion) attachStats(ctx context.Context, conn repository.DBConn, models []*model.TaskVersion) error {
	if len(models) == 0 {
		return nil
	}

	ids := make([]string, len(models))
	for i, m := range models {
		ids[i] = m.IDNatural
	}

	var stats []entity.TaskVersionStats
	err := conn.NewSelect().
		Model(&stats).
		Where("task_version_id IN (?)", bun.In(ids)).
		Scan(ctx)
	if err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to fetch task version stats: %v", err))
	}

	statsMap := make(map[string]*entity.TaskVersionStats, len(stats))
	for i := range stats {
		statsMap[stats[i].TaskVersionID] = &stats[i]
	}

	for _, m := range models {
		if s, ok := statsMap[m.IDNatural]; ok {
			m.ActualDurationSeconds = &s.TotalDurationSeconds
			m.ActualEpisodeCount = &s.EpisodeCount
		}
	}

	return nil
}

func (tv *taskVersion) SumTargetByTaskID(ctx context.Context, conn repository.DBConn, taskID string) (int64, error) {
	orgID, err := ccontext.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return 0, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "organization context required"))
	}

	var total int64

	err = conn.NewSelect().
		TableExpr("task_version").
		ColumnExpr("COALESCE(SUM(COALESCE(target_duration_seconds, 0)), 0)").
		Where("organization_id = ?", orgID).
		Where("task_id = ?", taskID).
		Where("approval_status = ?", model.ApprovalStatusApproved).
		Scan(ctx, &total)

	if err != nil {
		return 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to sum target by task: %v", err))
	}

	return total, nil
}
