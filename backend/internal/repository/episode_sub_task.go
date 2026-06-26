package repository

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

type EpisodeSubTask interface {
	BulkCreate(ctx context.Context, conn Conn, subtasks []model.EpisodeSubTask) error
	GetByID(ctx context.Context, conn Conn, id string) (model.EpisodeSubTask, error)
	GetByEpisodeID(ctx context.Context, conn Conn, episodeID string) (model.EpisodeSubTasks, error)
	Update(ctx context.Context, conn Conn, subtask model.EpisodeSubTask) error
	BulkCancelByEpisodeID(ctx context.Context, conn Conn, episodeID string) error
}
