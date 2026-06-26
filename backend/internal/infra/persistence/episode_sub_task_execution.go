package persistence

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/uptrace/bun"
)

type episodeSubTaskExecution struct{}

func NewEpisodeSubTaskExecution() *episodeSubTaskExecution { return &episodeSubTaskExecution{} }

func (e *episodeSubTaskExecution) Create(ctx context.Context, conn repository.DBConn, execution model.EpisodeSubTaskExecution) (model.EpisodeSubTaskExecution, error) {
	var inserted entity.EpisodeSubTaskExecution
	dbExe := entity.EpisodeSubTaskExecution{
		IDNatural:        execution.IDNatural,
		OrganizationID:   execution.OrganizationID,
		EpisodeSubTaskID: execution.EpisodeSubTaskID,
		ExecutionStatus:  execution.ExecutionStatus,
		StartedAt:        execution.StartedAt,
		FinishedAt:       execution.FinishedAt,
	}

	if err := conn.NewInsert().
		Model(&dbExe).
		Returning("*").
		Scan(ctx, &inserted); err != nil {
		return model.EpisodeSubTaskExecution{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create execution: %v", err))
	}

	return model.NewEpisodeSubTaskExecution(
		inserted.ID,
		inserted.IDNatural,
		inserted.OrganizationID,
		inserted.EpisodeSubTaskID,
		inserted.ExecutionStatus,
		inserted.StartedAt,
		inserted.FinishedAt,
		inserted.CreatedAt,
		&inserted.UpdatedAt,
	), nil
}

func (e *episodeSubTaskExecution) GetByID(ctx context.Context, conn repository.DBConn, id string) (model.EpisodeSubTaskExecution, error) {
	var dbExe entity.EpisodeSubTaskExecution

	if err := conn.NewSelect().
		Model(&dbExe).
		Where("id_natural = ?", id).
		Scan(ctx); err != nil {
		return model.EpisodeSubTaskExecution{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeExecutionNotFound, "execution not found: %v", err))
	}

	return model.NewEpisodeSubTaskExecution(
		dbExe.ID,
		dbExe.IDNatural,
		dbExe.OrganizationID,
		dbExe.EpisodeSubTaskID,
		dbExe.ExecutionStatus,
		dbExe.StartedAt,
		dbExe.FinishedAt,
		dbExe.CreatedAt,
		&dbExe.UpdatedAt,
	), nil
}

func (e *episodeSubTaskExecution) GetByEpisodeSubTaskIDs(ctx context.Context, conn repository.DBConn, ids []string) (model.EpisodeSubTaskExecutions, error) {
	if len(ids) == 0 {
		return model.EpisodeSubTaskExecutions{}, nil
	}

	var dbExecs []entity.EpisodeSubTaskExecution

	if err := conn.NewSelect().
		Model(&dbExecs).
		Where("episode_sub_task_id IN (?)", bun.In(ids)).
		Order("created_at ASC").
		Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get executions: %v", err))
	}

	result := make(model.EpisodeSubTaskExecutions, 0, len(dbExecs))
	for _, dbExe := range dbExecs {
		exec := model.NewEpisodeSubTaskExecution(
			dbExe.ID,
			dbExe.IDNatural,
			dbExe.OrganizationID,
			dbExe.EpisodeSubTaskID,
			dbExe.ExecutionStatus,
			dbExe.StartedAt,
			dbExe.FinishedAt,
			dbExe.CreatedAt,
			&dbExe.UpdatedAt,
		)
		result = append(result, &exec)
	}

	return result, nil
}

func (e *episodeSubTaskExecution) Update(ctx context.Context, conn repository.DBConn, execution model.EpisodeSubTaskExecution) error {
	dbExe := entity.EpisodeSubTaskExecution{
		ID:               execution.ID,
		IDNatural:        execution.IDNatural,
		OrganizationID:   execution.OrganizationID,
		EpisodeSubTaskID: execution.EpisodeSubTaskID,
		ExecutionStatus:  execution.ExecutionStatus,
		StartedAt:        execution.StartedAt,
		FinishedAt:       execution.FinishedAt,
	}

	if _, err := conn.NewUpdate().
		Model(&dbExe).
		WherePK().
		Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update execution: %v", err))
	}

	return nil
}

func (e *episodeSubTaskExecution) CountStartedBySubTaskID(ctx context.Context, conn repository.DBConn, episodeSubTaskID string) (int, error) {
	count, err := conn.NewSelect().
		Model((*entity.EpisodeSubTaskExecution)(nil)).
		Where("episode_sub_task_id = ?", episodeSubTaskID).
		Where("execution_status = ?", model.ExecutionStatusStarted).
		Count(ctx)

	if err != nil {
		return 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count started executions: %v", err))
	}

	return count, nil
}

func (e *episodeSubTaskExecution) BulkCancelByEpisodeID(ctx context.Context, conn repository.DBConn, episodeID string) error {
	// Cancel all executions that belong to subtasks of the given episode
	if _, err := conn.NewUpdate().
		Model((*entity.EpisodeSubTaskExecution)(nil)).
		Set("execution_status = ?", model.ExecutionStatusCancelled).
		Where("episode_sub_task_id IN (SELECT id_natural FROM episode_sub_task WHERE episode_id = ?)", episodeID).
		Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to bulk cancel executions: %v", err))
	}

	return nil
}
