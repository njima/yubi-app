package repository

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

type SubTaskListFilter struct {
	TaskID        *string
	TaskVersionID *string
}

type SubTask interface {
	Create(ctx context.Context, conn Conn, s model.SubTask) (model.SubTask, error)
	GetByID(ctx context.Context, conn Conn, id string) (model.SubTask, error)
	GetByTaskVersionID(ctx context.Context, conn Conn, taskVersionID string) (model.SubTasks, error)
	GetMaxOrderIndex(ctx context.Context, conn Conn, taskVersionID string) (int, error)
	List(ctx context.Context, conn Conn, filter SubTaskListFilter, limit, offset int) (model.SubTasks, int, error)
	Update(ctx context.Context, conn Conn, s model.SubTask) (model.SubTask, error)
	UpdateOrderIndices(ctx context.Context, conn Conn, ids []string) error
	Delete(ctx context.Context, conn Conn, id string) error
}
