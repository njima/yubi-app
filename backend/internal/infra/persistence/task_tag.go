package persistence

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/requestctx"
)

type taskTag struct{}

func NewTaskTag() *taskTag { return &taskTag{} }

func (t *taskTag) ListCategoryTypes(ctx context.Context, conn repository.DBConn) (model.TaskCategoryTypes, error) {
	var rows []entity.TaskCategoryType
	if err := conn.NewSelect().Model(&rows).Order("name ASC").Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list task category types: %v", err))
	}
	result := make(model.TaskCategoryTypes, 0, len(rows))
	for _, r := range rows {
		result = append(result, &model.TaskCategoryType{ID: r.ID, Slug: r.Slug, Name: r.Name})
	}
	return result, nil
}

func (t *taskTag) GetCategoryTypeByID(ctx context.Context, conn repository.DBConn, id string) (model.TaskCategoryType, error) {
	var row entity.TaskCategoryType
	if err := conn.NewSelect().Model(&row).Where("id = ?", id).Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.TaskCategoryType{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "task category type not found: id=%s", id))
		}
		return model.TaskCategoryType{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task category type: %v", err))
	}
	return model.TaskCategoryType{ID: row.ID, Slug: row.Slug, Name: row.Name}, nil
}

func (t *taskTag) ListTags(ctx context.Context, conn repository.DBConn, categoryTypeID *string) (model.TaskTags, error) {
	var rows []entity.TaskTag
	q := conn.NewSelect().
		Model(&rows).
		Relation("CategoryType").
		Order("tt.name ASC")
	if categoryTypeID != nil && *categoryTypeID != "" {
		q = q.Where("tt.category_type_id = ?", *categoryTypeID)
	}
	if err := q.Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list task tags: %v", err))
	}
	return toTagModels(rows), nil
}

func (t *taskTag) CreateTag(ctx context.Context, conn repository.DBConn, tag model.TaskTag) (model.TaskTag, error) {
	row := entity.TaskTag{
		ID:             tag.ID,
		Name:           tag.Name,
		CategoryTypeID: tag.CategoryTypeID,
	}
	if _, err := conn.NewInsert().Model(&row).Exec(ctx); err != nil {
		return model.TaskTag{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create task tag: %v", err))
	}
	return t.GetTagByID(ctx, conn, row.ID)
}

func (t *taskTag) GetTagByID(ctx context.Context, conn repository.DBConn, id string) (model.TaskTag, error) {
	var row entity.TaskTag
	if err := conn.NewSelect().
		Model(&row).
		Relation("CategoryType").
		Where("tt.id = ?", id).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.TaskTag{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "task tag not found: id=%s", id))
		}
		return model.TaskTag{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task tag: %v", err))
	}
	return toTagModel(row), nil
}

func (t *taskTag) SetTaskTags(ctx context.Context, conn repository.DBConn, taskID string, tagIDs []string) error {
	if _, err := conn.NewDelete().
		Model((*entity.TaskTagAssignment)(nil)).
		Where("task_id = ?", taskID).
		Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete task tag assignments: %v", err))
	}
	if len(tagIDs) == 0 {
		return nil
	}
	rows := make([]entity.TaskTagAssignment, 0, len(tagIDs))
	for _, tid := range tagIDs {
		rows = append(rows, entity.TaskTagAssignment{TaskID: taskID, TagID: tid})
	}
	if _, err := conn.NewInsert().Model(&rows).Exec(ctx); err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to insert task tag assignments: %v", err))
	}
	return nil
}

func (t *taskTag) GetTagsByTaskID(ctx context.Context, conn repository.DBConn, taskID string) (model.TaskTags, error) {
	m, err := t.GetTagsByTaskIDs(ctx, conn, []string{taskID})
	if err != nil {
		return nil, err
	}
	return m[taskID], nil
}

func (t *taskTag) GetTagsByTaskIDs(ctx context.Context, conn repository.DBConn, taskIDs []string) (map[string]model.TaskTags, error) {
	if len(taskIDs) == 0 {
		return map[string]model.TaskTags{}, nil
	}
	var assignments []entity.TaskTagAssignment
	if err := conn.NewSelect().
		Model(&assignments).
		Relation("Tag.CategoryType").
		Where("tta.task_id IN (?)", bun.In(taskIDs)).
		Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get tags by task IDs: %v", err))
	}
	result := make(map[string]model.TaskTags, len(taskIDs))
	for _, a := range assignments {
		if a.Tag != nil {
			result[a.TaskID] = append(result[a.TaskID], &model.TaskTag{
				ID:               a.Tag.ID,
				Name:             a.Tag.Name,
				CategoryTypeID:   a.Tag.CategoryTypeID,
				CategoryTypeName: categoryTypeName(a.Tag.CategoryType),
			})
		}
	}
	return result, nil
}

func toTagModel(r entity.TaskTag) model.TaskTag {
	return model.TaskTag{
		ID:               r.ID,
		Name:             r.Name,
		CategoryTypeID:   r.CategoryTypeID,
		CategoryTypeName: categoryTypeName(r.CategoryType),
	}
}

func toTagModels(rows []entity.TaskTag) model.TaskTags {
	result := make(model.TaskTags, 0, len(rows))
	for _, r := range rows {
		m := toTagModel(r)
		result = append(result, &m)
	}
	return result
}

func categoryTypeName(ct *entity.TaskCategoryType) string {
	if ct == nil {
		return ""
	}
	return ct.Name
}

func (t *taskTag) GetTagsByNames(ctx context.Context, conn repository.DBConn, names []string) (model.TaskTags, error) {
	if len(names) == 0 {
		return model.TaskTags{}, nil
	}
	var rows []entity.TaskTag
	if err := conn.NewSelect().
		Model(&rows).
		Relation("CategoryType").
		Where("tt.name IN (?)", bun.In(names)).
		Order("tt.name ASC").
		Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task tags by names: %v", err))
	}
	return toTagModels(rows), nil
}

func (tt *taskTag) GetAvailableTags(ctx context.Context, conn repository.DBConn, robotTypes []string, categoryTypeID *string) (model.TaskTags, error) {
	type tagRow struct {
		ID               string `bun:"id"`
		Name             string `bun:"name"`
		CategoryTypeID   string `bun:"category_type_id"`
		CategoryTypeName string `bun:"category_type_name"`
	}

	var rows []tagRow

	sel := conn.NewSelect().
		TableExpr("task_tag AS tt").
		ColumnExpr("DISTINCT tt.id, tt.name, tt.category_type_id").
		ColumnExpr("tct.name AS category_type_name").
		Join("JOIN task_category_type tct ON tct.id = tt.category_type_id").
		Join("JOIN task_tag_assignment tta ON tta.tag_id = tt.id").
		Join("JOIN task t ON t.id_natural = tta.task_id").
		OrderExpr("tt.name ASC")

	if orgID, err := requestctx.OrganizationID(ctx); err == nil && orgID != "" {
		sel = sel.Where("t.organization_id = ?", orgID)
	}
	if len(robotTypes) > 0 {
		sel = sel.Where("t.robot_type IN (?)", bun.In(robotTypes))
	}
	if categoryTypeID != nil {
		sel = sel.Where("tt.category_type_id = ?", *categoryTypeID)
	}

	if err := sel.Scan(ctx, &rows); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get available tags: %v", err))
	}

	result := make(model.TaskTags, 0, len(rows))
	for _, r := range rows {
		result = append(result, &model.TaskTag{
			ID:               r.ID,
			Name:             r.Name,
			CategoryTypeID:   r.CategoryTypeID,
			CategoryTypeName: r.CategoryTypeName,
		})
	}
	return result, nil
}
