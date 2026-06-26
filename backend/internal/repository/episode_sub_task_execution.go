package repository

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

type EpisodeSubTaskExecution interface {
	Create(ctx context.Context, conn Conn, execution model.EpisodeSubTaskExecution) (model.EpisodeSubTaskExecution, error)
	GetByID(ctx context.Context, conn Conn, id string) (model.EpisodeSubTaskExecution, error)
	GetByEpisodeSubTaskIDs(ctx context.Context, conn Conn, ids []string) (model.EpisodeSubTaskExecutions, error)
	Update(ctx context.Context, conn Conn, execution model.EpisodeSubTaskExecution) error
	CountStartedBySubTaskID(ctx context.Context, conn Conn, episodeSubTaskID string) (int, error)
	BulkCancelByEpisodeID(ctx context.Context, conn Conn, episodeID string) error
}
