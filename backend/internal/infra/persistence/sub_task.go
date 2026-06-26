package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
)

type subtask struct{}

func NewSubTask() *subtask { return &subtask{} }

func (s *subtask) Create(ctx context.Context, conn repository.DBConn, st model.SubTask) (model.SubTask, error) {
	var inserted entity.SubTask
	dbSt := entity.SubTask{
		IDNatural:             st.IDNatural,
		OrganizationID:        st.OrganizationID,
		TaskVersionID:         st.TaskVersionID,
		Name:                  st.Name,
		OrderIndex:            st.OrderIndex,
		Description:           st.Description,
		TargetDurationSeconds: st.TargetDurationSeconds,
	}

	if err := conn.NewInsert().
		Model(&dbSt).
		Returning("*").
		Scan(ctx, &inserted); err != nil {
		return model.SubTask{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create subtask: %v", err))
	}

	return model.SubTask{
		ID:                    inserted.ID,
		IDNatural:             inserted.IDNatural,
		TaskVersionID:         inserted.TaskVersionID,
		Name:                  inserted.Name,
		OrderIndex:            inserted.OrderIndex,
		Description:           inserted.Description,
		TargetDurationSeconds: inserted.TargetDurationSeconds,
		CreatedAt:             inserted.CreatedAt,
		UpdatedAt:             &inserted.UpdatedAt,
	}, nil
}

func (s *subtask) GetByID(ctx context.Context, conn repository.DBConn, id string) (model.SubTask, error) {
	var dbs entity.SubTask
	if err := conn.NewSelect().
		Model(&dbs).
		Where("id_natural = ?", id).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.SubTask{}, apperror.NewError(apperror.NewMessage(apperror.CodeSubTaskNotFound, "subtask not found: id_natural=%s", id))
		}
		return model.SubTask{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get subtask: %v", err))
	}

	res := model.SubTask{
		ID:                    dbs.ID,
		IDNatural:             dbs.IDNatural,
		TaskVersionID:         dbs.TaskVersionID,
		Name:                  dbs.Name,
		OrderIndex:            dbs.OrderIndex,
		Description:           dbs.Description,
		TargetDurationSeconds: dbs.TargetDurationSeconds,
		CreatedAt:             dbs.CreatedAt,
	}
	if !dbs.UpdatedAt.IsZero() {
		t2 := dbs.UpdatedAt
		res.UpdatedAt = &t2
	}

	return res, nil
}

func (s *subtask) GetMaxOrderIndex(ctx context.Context, conn repository.DBConn, taskVersionID string) (int, error) {
	var maxIndex int
	err := conn.NewSelect().
		Model((*entity.SubTask)(nil)).
		ColumnExpr("COALESCE(MAX(order_index), -1)").
		Where("task_version_id = ?", taskVersionID).
		Scan(ctx, &maxIndex)
	if err != nil {
		return 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get max order index: %v", err))
	}
	return maxIndex, nil
}

func (s *subtask) GetByTaskVersionID(ctx context.Context, conn repository.DBConn, taskVersionID string) (model.SubTasks, error) {
	var dbs []entity.SubTask
	if err := conn.NewSelect().
		Model(&dbs).
		Where("task_version_id = ?", taskVersionID).
		Order("order_index ASC").
		Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get subtasks by task version: %v", err))
	}

	res := make(model.SubTasks, 0, len(dbs))
	for _, d := range dbs {
		m := &model.SubTask{
			ID:                    d.ID,
			IDNatural:             d.IDNatural,
			TaskVersionID:         d.TaskVersionID,
			Name:                  d.Name,
			OrderIndex:            d.OrderIndex,
			Description:           d.Description,
			TargetDurationSeconds: d.TargetDurationSeconds,
			CreatedAt:             d.CreatedAt,
		}
		if !d.UpdatedAt.IsZero() {
			t2 := d.UpdatedAt
			m.UpdatedAt = &t2
		}
		res = append(res, m)
	}

	return res, nil
}

func (s *subtask) List(ctx context.Context, conn repository.DBConn, filter repository.SubTaskListFilter, limit, offset int) (model.SubTasks, int, error) {
	var dbs []entity.SubTask
	sel := conn.NewSelect().
		Model(&dbs).
		Order("order_index ASC").
		Limit(limit).
		Offset(offset)

	sel = applySubTaskListFilters(sel, filter)

	if err := sel.Scan(ctx); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list subtasks: %v", err))
	}

	var total int
	countSel := conn.NewSelect().Model((*entity.SubTask)(nil)).ColumnExpr("COUNT(*)")
	countSel = applySubTaskListFilters(countSel, filter)
	if err := countSel.Scan(ctx, &total); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count subtasks: %v", err))
	}

	res := make(model.SubTasks, 0, len(dbs))
	for _, d := range dbs {
		m := &model.SubTask{
			ID:                    d.ID,
			IDNatural:             d.IDNatural,
			OrganizationID:        d.OrganizationID,
			TaskVersionID:         d.TaskVersionID,
			Name:                  d.Name,
			OrderIndex:            d.OrderIndex,
			Description:           d.Description,
			TargetDurationSeconds: d.TargetDurationSeconds,
			CreatedAt:             d.CreatedAt,
		}
		if !d.UpdatedAt.IsZero() {
			t2 := d.UpdatedAt
			m.UpdatedAt = &t2
		}
		res = append(res, m)
	}

	return res, total, nil
}

func (s *subtask) Update(ctx context.Context, conn repository.DBConn, st model.SubTask) (model.SubTask, error) {
	upd := conn.NewUpdate().Model((*entity.SubTask)(nil))
	hasSet := false
	if st.Name != "" {
		upd = upd.Set("name = ?", st.Name)
		hasSet = true
	}
	if st.OrderIndex != 0 {
		upd = upd.Set("order_index = ?", st.OrderIndex)
		hasSet = true
	}
	if st.Description != nil {
		upd = upd.Set("description = ?", *st.Description)
		hasSet = true
	}
	if st.TargetDurationSeconds != nil {
		upd = upd.Set("target_duration_seconds = ?", *st.TargetDurationSeconds)
		hasSet = true
	}
	if !hasSet {
		return model.SubTask{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "no fields to update"))
	}

	upd = upd.Set("updated_at = ?", time.Now().UTC())

	var updated entity.SubTask
	if err := upd.Where("id_natural = ?", st.IDNatural).Returning("*").Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.SubTask{}, apperror.NewError(apperror.NewMessage(apperror.CodeSubTaskNotFound, "subtask not found: id_natural=%s", st.IDNatural))
		}
		return model.SubTask{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update subtask: %v", err))
	}

	result := model.SubTask{
		ID:                    updated.ID,
		IDNatural:             updated.IDNatural,
		TaskVersionID:         updated.TaskVersionID,
		Name:                  updated.Name,
		OrderIndex:            updated.OrderIndex,
		Description:           updated.Description,
		TargetDurationSeconds: updated.TargetDurationSeconds,
		CreatedAt:             updated.CreatedAt,
	}
	if !updated.UpdatedAt.IsZero() {
		t2 := updated.UpdatedAt
		result.UpdatedAt = &t2
	}

	return result, nil
}

func (s *subtask) UpdateOrderIndices(ctx context.Context, conn repository.DBConn, ids []string) error {
	now := time.Now().UTC()
	for i, id := range ids {
		_, err := conn.NewUpdate().
			Model((*entity.SubTask)(nil)).
			Set("order_index = ?", i).
			Set("updated_at = ?", now).
			Where("id_natural = ?", id).
			Exec(ctx)
		if err != nil {
			return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update order_index for subtask: id_natural=%s", id))
		}
	}
	return nil
}

func (s *subtask) Delete(ctx context.Context, conn repository.DBConn, id string) error {
	var deletedID int64
	if err := conn.NewDelete().
		Model((*entity.SubTask)(nil)).
		Where("id_natural = ?", id).
		Returning("id").
		Scan(ctx, &deletedID); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewError(apperror.NewMessage(apperror.CodeSubTaskNotFound, "subtask not found: id_natural=%s", id))
		}
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete subtask: %v", err))
	}
	return nil
}
