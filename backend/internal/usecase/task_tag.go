package usecase

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type TaskTagUsecase interface {
	ListCategoryTypes(ctx context.Context) (model.TaskCategoryTypes, error)
	ListTags(ctx context.Context, categoryTypeID *string) (model.TaskTags, error)
	CreateTag(ctx context.Context, input TaskTagCreateInput) (model.TaskTag, error)
	GetAvailableTags(ctx context.Context, robotTypes []string, categoryTypeID *string) (model.TaskTags, error)
}

type TaskTagCreateInput struct {
	Name           string
	CategoryTypeID string
}

type taskTag struct {
	repo repository.TaskTag
	db   repository.DBConn
}

func NewTaskTag(repo repository.TaskTag, db repository.DBConn) *taskTag {
	return &taskTag{repo: repo, db: db}
}

func (t *taskTag) ListCategoryTypes(ctx context.Context) (model.TaskCategoryTypes, error) {
	return t.repo.ListCategoryTypes(ctx, t.db)
}

func (t *taskTag) ListTags(ctx context.Context, categoryTypeID *string) (model.TaskTags, error) {
	return t.repo.ListTags(ctx, t.db, categoryTypeID)
}

func (t *taskTag) CreateTag(ctx context.Context, input TaskTagCreateInput) (model.TaskTag, error) {
	id, err := model.InitID()
	if err != nil {
		return model.TaskTag{}, err
	}
	tag := model.TaskTag{
		ID:             id,
		Name:           input.Name,
		CategoryTypeID: input.CategoryTypeID,
	}
	return t.repo.CreateTag(ctx, t.db, tag)
}

func (t *taskTag) GetAvailableTags(ctx context.Context, robotTypes []string, categoryTypeID *string) (model.TaskTags, error) {
	return t.repo.GetAvailableTags(ctx, t.db, robotTypes, categoryTypeID)
}
