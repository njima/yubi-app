package repository

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

type TaskTag interface {
	ListCategoryTypes(ctx context.Context, conn Conn) (model.TaskCategoryTypes, error)
	GetCategoryTypeByID(ctx context.Context, conn Conn, id string) (model.TaskCategoryType, error)
	ListTags(ctx context.Context, conn Conn, categoryTypeID *string) (model.TaskTags, error)
	CreateTag(ctx context.Context, conn Conn, tag model.TaskTag) (model.TaskTag, error)
	GetTagByID(ctx context.Context, conn Conn, id string) (model.TaskTag, error)
	SetTaskTags(ctx context.Context, conn Conn, taskID string, tagIDs []string) error
	GetTagsByTaskID(ctx context.Context, conn Conn, taskID string) (model.TaskTags, error)
	GetTagsByTaskIDs(ctx context.Context, conn Conn, taskIDs []string) (map[string]model.TaskTags, error)
	GetAvailableTags(ctx context.Context, conn Conn, robotTypes []string, categoryTypeID *string) (model.TaskTags, error)
	// GetTagsByNames returns only the tags whose names match the given list.
	// Used by the CSV import flow to avoid fetching all tags in the system.
	GetTagsByNames(ctx context.Context, conn Conn, names []string) (model.TaskTags, error)
}
