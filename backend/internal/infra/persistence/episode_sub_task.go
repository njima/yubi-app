package persistence

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type episodeSubTask struct{}

func NewEpisodeSubTask() *episodeSubTask { return &episodeSubTask{} }

func (e *episodeSubTask) BulkCreate(ctx context.Context, conn repository.DBConn, subtasks []model.EpisodeSubTask) error {
	if len(subtasks) == 0 {
		return nil
	}

	dbSubtasks := make([]entity.EpisodeSubTask, 0, len(subtasks))
	for _, st := range subtasks {
		dbSubtasks = append(dbSubtasks, entity.EpisodeSubTask{
			IDNatural:        st.IDNatural,
			OrganizationID:   st.OrganizationID,
			EpisodeID:        st.EpisodeID,
			SubTaskID:        st.SubTaskID,
			CollectionStatus: st.CollectionStatus,
		})
	}

	if _, err := conn.NewInsert().
		Model(&dbSubtasks).
		Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to bulk create episode subtasks: %v", err))
	}

	return nil
}

func (e *episodeSubTask) GetByID(ctx context.Context, conn repository.DBConn, id string) (model.EpisodeSubTask, error) {
	var dbSt entity.EpisodeSubTask

	if err := conn.NewSelect().
		Model(&dbSt).
		Where("id_natural = ?", id).
		Scan(ctx); err != nil {
		return model.EpisodeSubTask{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeEpisodeSubTaskNotFound, "episode subtask not found: %v", err))
	}

	return model.NewEpisodeSubTask(
		dbSt.ID,
		dbSt.IDNatural,
		dbSt.OrganizationID,
		dbSt.EpisodeID,
		dbSt.SubTaskID,
		dbSt.CollectionStatus,
		dbSt.CreatedAt,
		&dbSt.UpdatedAt,
	), nil
}

func (e *episodeSubTask) GetByEpisodeID(ctx context.Context, conn repository.DBConn, episodeID string) (model.EpisodeSubTasks, error) {
	var dbSubtasks []entity.EpisodeSubTask

	if err := conn.NewSelect().
		Model(&dbSubtasks).
		Where("episode_id = ?", episodeID).
		Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get episode subtasks: %v", err))
	}

	result := make(model.EpisodeSubTasks, 0, len(dbSubtasks))
	for _, dbSt := range dbSubtasks {
		st := model.NewEpisodeSubTask(
			dbSt.ID,
			dbSt.IDNatural,
			dbSt.OrganizationID,
			dbSt.EpisodeID,
			dbSt.SubTaskID,
			dbSt.CollectionStatus,
			dbSt.CreatedAt,
			&dbSt.UpdatedAt,
		)
		result = append(result, &st)
	}

	return result, nil
}

func (e *episodeSubTask) Update(ctx context.Context, conn repository.DBConn, subtask model.EpisodeSubTask) error {
	dbSt := entity.EpisodeSubTask{
		ID:               subtask.ID,
		IDNatural:        subtask.IDNatural,
		OrganizationID:   subtask.OrganizationID,
		EpisodeID:        subtask.EpisodeID,
		SubTaskID:        subtask.SubTaskID,
		CollectionStatus: subtask.CollectionStatus,
	}

	if _, err := conn.NewUpdate().
		Model(&dbSt).
		WherePK().
		Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update episode subtask: %v", err))
	}

	return nil
}

func (e *episodeSubTask) BulkCancelByEpisodeID(ctx context.Context, conn repository.DBConn, episodeID string) error {
	if _, err := conn.NewUpdate().
		Model((*entity.EpisodeSubTask)(nil)).
		Set("collection_status = ?", model.SubTaskCollectionStatusCancelled).
		Where("episode_id = ?", episodeID).
		Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to bulk cancel episode subtasks: %v", err))
	}

	return nil
}
