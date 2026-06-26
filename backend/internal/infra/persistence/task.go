package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/database/bunconv"
	"github.com/airoa-org/yubi-app/backend/internal/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/requestctx"
)

type task struct{}

func NewTask() *task { return &task{} }

func (t *task) Create(ctx context.Context, conn repository.DBConn, tk model.Task) (model.Task, error) {
	var inserted entity.Task

	dbTk := entity.Task{
		IDNatural:      tk.IDNatural,
		OrganizationID: tk.OrganizationID,
		Name:           tk.Name,
		Description:    tk.Description,
		ManualURL:      &tk.ManualURL,
		Priority:       derefPriority(tk.Priority),
		Difficulty:     derefDifficulty(tk.Difficulty),
		Status:         derefTaskStatus(tk.Status),
		Deadline:       tk.Deadline,
		RobotType:      tk.RobotType,
	}

	if err := conn.NewInsert().
		Model(&dbTk).
		Returning("*").
		Scan(ctx, &inserted); err != nil {
		return model.Task{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create task: %v", err))
	}

	tvID, err := model.InitID()
	if err != nil {
		return model.Task{}, err
	}

	var insertedVersion entity.TaskVersion
	dbTaskVersion := entity.TaskVersion{
		IDNatural:      tvID,
		TaskID:         inserted.IDNatural,
		OrganizationID: tk.OrganizationID,
		Version:        tk.Version,
		IsActive:       tk.IsActive,
		ApprovalStatus: model.ApprovalStatusDraft,
	}

	if err := conn.NewInsert().
		Model(&dbTaskVersion).
		Returning("*").
		Scan(ctx, &insertedVersion); err != nil {
		return model.Task{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to create task version: %v", err))
	}

	return bunconv.EntityToTaskModel(inserted, insertedVersion), nil
}

func (t *task) Exists(ctx context.Context, conn repository.DBConn, id string) (bool, error) {
	exists, err := conn.NewSelect().
		Model((*entity.Task)(nil)).
		Where("id_natural = ?", id).
		Exists(ctx)
	if err != nil {
		return false, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to check task existence: %v", err))
	}
	return exists, nil
}

func (t *task) GetByID(ctx context.Context, conn repository.DBConn, id string) (model.Task, error) {
	var dbt entity.Task
	if err := conn.NewSelect().
		Model(&dbt).
		Where("id_natural = ?", id).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return model.Task{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskNotFound, "task not found: id_natural=%s", id))
		}
		return model.Task{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task: %v", err))
	}

	tk := model.Task{
		ID:          dbt.ID,
		IDNatural:   dbt.IDNatural,
		Name:        dbt.Name,
		Description: dbt.Description,
		ManualURL:   derefString(dbt.ManualURL),
		Priority:    modelTaskPriorityPtr(dbt.Priority),
		Difficulty:  modelTaskDifficultyPtr(dbt.Difficulty),
		Status:      modelTaskStatusPtr(dbt.Status),
		Deadline:    dbt.Deadline,
		RobotType:   dbt.RobotType,
		CreatedAt:   dbt.CreatedAt,
	}

	var tv entity.TaskVersion
	err := conn.NewSelect().
		Model(&tv).
		Where("task_id = ?", dbt.IDNatural).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return model.Task{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task version: %v", err))
	}

	if err != sql.ErrNoRows {
		tk.Version = tv.Version
		tk.IsActive = tv.IsActive
		tk.TargetDurationSeconds = tv.TargetDurationSeconds
		tk.TargetEpisodeCount = tv.TargetEpisodeCount

		var stats entity.TaskVersionStats
		statsErr := conn.NewSelect().
			Model(&stats).
			Where("task_version_id = ?", tv.IDNatural).
			Scan(ctx)
		if statsErr == nil {
			tk.ActualEpisodeCount = &stats.EpisodeCount
		}
	}

	if !dbt.UpdatedAt.IsZero() {
		t2 := dbt.UpdatedAt
		tk.UpdatedAt = &t2
	}

	return tk, nil
}

func (t *task) List(ctx context.Context, conn repository.DBConn, filter repository.TaskListFilter, limit, offset int) (model.Tasks, int, error) {
	var dbts []entity.Task

	sel := conn.NewSelect().
		Model(&dbts).
		Limit(limit).
		Offset(offset)

	// Dynamic ORDER BY with whitelist to prevent SQL injection
	sel = applyTaskSortOrder(sel, filter.SortBy, filter.SortOrder)

	if filter.HasApprovedVersion != nil && *filter.HasApprovedVersion {
		sel = sel.Where("EXISTS (SELECT 1 FROM task_version tv WHERE tv.task_id = t.id_natural AND tv.approval_status = 1)")
	}
	if len(filter.Statuses) > 0 {
		sel = sel.Where("t.status IN (?)", bun.In(filter.Statuses))
	}
	if len(filter.Priorities) > 0 {
		sel = sel.Where("t.priority IN (?)", bun.In(filter.Priorities))
	}
	if len(filter.Difficulties) > 0 {
		sel = sel.Where("t.difficulty IN (?)", bun.In(filter.Difficulties))
	}
	if filter.RobotType != nil {
		sel = sel.Where("t.robot_type = ?", *filter.RobotType)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		sel = sel.Where("t.name ILIKE ?", "%"+escaped+"%")
	}
	if err := sel.Scan(ctx); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list tasks: %v", err))
	}

	var total int
	countSel := conn.NewSelect().Model((*entity.Task)(nil)).ColumnExpr("COUNT(*)")
	if filter.HasApprovedVersion != nil && *filter.HasApprovedVersion {
		countSel = countSel.Where("EXISTS (SELECT 1 FROM task_version tv WHERE tv.task_id = t.id_natural AND tv.approval_status = 1)")
	}
	if len(filter.Statuses) > 0 {
		countSel = countSel.Where("t.status IN (?)", bun.In(filter.Statuses))
	}
	if len(filter.Priorities) > 0 {
		countSel = countSel.Where("t.priority IN (?)", bun.In(filter.Priorities))
	}
	if len(filter.Difficulties) > 0 {
		countSel = countSel.Where("t.difficulty IN (?)", bun.In(filter.Difficulties))
	}
	if filter.RobotType != nil {
		countSel = countSel.Where("t.robot_type = ?", *filter.RobotType)
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := escapeILIKE(*filter.Search)
		countSel = countSel.Where("t.name ILIKE ?", "%"+escaped+"%")
	}
	if err := countSel.Scan(ctx, &total); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count tasks: %v", err))
	}

	if len(dbts) == 0 {
		return model.Tasks{}, total, nil
	}

	taskIDs := make([]string, 0, len(dbts))
	for _, d := range dbts {
		taskIDs = append(taskIDs, d.IDNatural)
	}

	var taskVersions []entity.TaskVersion
	if err := conn.NewSelect().
		Model(&taskVersions).
		Where("task_id IN (?)", bun.In(taskIDs)).
		Order("created_at DESC").
		Scan(ctx); err != nil && err != sql.ErrNoRows {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get task versions: %v", err))
	}

	// Group by task_id and pick the latest (first in DESC order)
	latestVersionMap := make(map[string]entity.TaskVersion)
	for _, tv := range taskVersions {
		if _, exists := latestVersionMap[tv.TaskID]; !exists {
			latestVersionMap[tv.TaskID] = tv
		}
	}

	// Fetch stats for latest versions
	versionIDs := make([]string, 0, len(latestVersionMap))
	for _, tv := range latestVersionMap {
		versionIDs = append(versionIDs, tv.IDNatural)
	}
	statsMap := make(map[string]*entity.TaskVersionStats)
	if len(versionIDs) > 0 {
		var stats []entity.TaskVersionStats
		if err := conn.NewSelect().
			Model(&stats).
			Where("task_version_id IN (?)", bun.In(versionIDs)).
			Scan(ctx); err != nil && err != sql.ErrNoRows {
			return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to fetch task version stats: %v", err))
		}
		for i := range stats {
			statsMap[stats[i].TaskVersionID] = &stats[i]
		}
	}

	res := make(model.Tasks, 0, len(dbts))
	for _, d := range dbts {
		m := &model.Task{
			ID:             d.ID,
			IDNatural:      d.IDNatural,
			OrganizationID: d.OrganizationID,
			Name:           d.Name,
			Description:    d.Description,
			ManualURL:      derefString(d.ManualURL),
			Priority:       modelTaskPriorityPtr(d.Priority),
			Difficulty:     modelTaskDifficultyPtr(d.Difficulty),
			Status:         modelTaskStatusPtr(d.Status),
			Deadline:       d.Deadline,
			RobotType:      d.RobotType,
			CreatedAt:      d.CreatedAt,
		}
		if tv, ok := latestVersionMap[d.IDNatural]; ok {
			m.Version = tv.Version
			m.IsActive = tv.IsActive
			m.TargetDurationSeconds = tv.TargetDurationSeconds
			m.TargetEpisodeCount = tv.TargetEpisodeCount
			if s, found := statsMap[tv.IDNatural]; found {
				m.ActualEpisodeCount = &s.EpisodeCount
			}
		}
		if !d.UpdatedAt.IsZero() {
			t2 := d.UpdatedAt
			m.UpdatedAt = &t2
		}
		res = append(res, m)
	}

	return res, total, nil
}

func (t *task) Update(ctx context.Context, conn repository.DBConn, tk model.Task) (model.Task, error) {
	upd := conn.NewUpdate().Model((*entity.Task)(nil))
	hasSet := false
	if tk.Name != "" {
		upd = upd.Set("name = ?", tk.Name)
		hasSet = true
	}
	if tk.Description != nil {
		upd = upd.Set("description = ?", *tk.Description)
		hasSet = true
	}
	if tk.ManualURL != "" {
		upd = upd.Set("manual_url = ?", tk.ManualURL)
		hasSet = true
	}
	if tk.Priority != nil {
		upd = upd.Set("priority = ?", *tk.Priority)
		hasSet = true
	}
	if tk.Difficulty != nil {
		upd = upd.Set("difficulty = ?", *tk.Difficulty)
		hasSet = true
	}
	if tk.Status != nil {
		upd = upd.Set("status = ?", *tk.Status)
		hasSet = true
	}
	if !tk.Deadline.IsZero() {
		upd = upd.Set("deadline = ?", tk.Deadline)
		hasSet = true
	}
	if tk.RobotType != nil {
		upd = upd.Set("robot_type = ?", *tk.RobotType)
		hasSet = true
	}
	if !hasSet {
		return model.Task{}, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "no fields to update"))
	}

	upd = upd.Set("updated_at = ?", time.Now().UTC())

	var updated entity.Task
	if err := upd.Where("id_natural = ?", tk.IDNatural).Returning("*").Scan(ctx, &updated); err != nil {
		if err == sql.ErrNoRows {
			return model.Task{}, apperror.NewError(apperror.NewMessage(apperror.CodeTaskNotFound, "task not found: id_natural=%s", tk.IDNatural))
		}
		return model.Task{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to update task: %v", err))
	}

	result := model.Task{
		ID:          updated.ID,
		IDNatural:   updated.IDNatural,
		Name:        updated.Name,
		Description: updated.Description,
		ManualURL:   derefString(updated.ManualURL),
		Priority:    modelTaskPriorityPtr(updated.Priority),
		Difficulty:  modelTaskDifficultyPtr(updated.Difficulty),
		Status:      modelTaskStatusPtr(updated.Status),
		Deadline:    updated.Deadline,
		RobotType:   updated.RobotType,
		CreatedAt:   updated.CreatedAt,
	}
	if !updated.UpdatedAt.IsZero() {
		t2 := updated.UpdatedAt
		result.UpdatedAt = &t2
	}

	return result, nil
}

func (t *task) ListByIDs(ctx context.Context, conn repository.DBConn, ids []string) (model.Tasks, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var dbts []entity.Task
	if err := conn.NewSelect().
		Model(&dbts).
		Where("id_natural IN (?)", bun.In(ids)).
		Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list tasks by IDs: %v", err))
	}

	res := make(model.Tasks, 0, len(dbts))
	for _, d := range dbts {
		res = append(res, &model.Task{
			ID:          d.ID,
			IDNatural:   d.IDNatural,
			Name:        d.Name,
			Description: d.Description,
			ManualURL:   derefString(d.ManualURL),
			Priority:    modelTaskPriorityPtr(d.Priority),
			Difficulty:  modelTaskDifficultyPtr(d.Difficulty),
			Status:      modelTaskStatusPtr(d.Status),
			Deadline:    d.Deadline,
			RobotType:   d.RobotType,
		})
	}

	return res, nil
}

func derefDifficulty(d *model.TaskDifficulty) model.TaskDifficulty {
	if d == nil {
		return model.TaskDifficultyB
	}
	return *d
}

func derefTaskStatus(s *model.TaskStatus) model.TaskStatus {
	if s == nil {
		return model.TaskStatusPlanning
	}
	return *s
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefPriority(p *model.TaskPriority) model.TaskPriority {
	if p == nil {
		return model.TaskPriorityNormal
	}
	return *p
}

func modelTaskPriorityPtr(priority model.TaskPriority) *model.TaskPriority {
	return &priority
}

func modelTaskDifficultyPtr(difficulty model.TaskDifficulty) *model.TaskDifficulty {
	return &difficulty
}

func modelTaskStatusPtr(status model.TaskStatus) *model.TaskStatus {
	return &status
}

func (t *task) Delete(ctx context.Context, conn repository.DBConn, id string) error {
	var deletedID int64
	if err := conn.NewDelete().
		Model((*entity.Task)(nil)).
		Where("id_natural = ?", id).
		Returning("id").
		Scan(ctx, &deletedID); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewError(apperror.NewMessage(apperror.CodeTaskNotFound, "task not found: id_natural=%s", id))
		}
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to delete task: %v", err))
	}
	return nil
}

var allowedTaskSortColumns = map[string]string{
	"name":                    "t.name",
	"robot_type":              "t.robot_type",
	"priority":                "t.priority",
	"difficulty":              "t.difficulty",
	"status":                  "t.status",
	"target_duration_seconds": "(SELECT tv2.target_duration_seconds FROM task_version tv2 WHERE tv2.task_id = t.id_natural ORDER BY tv2.created_at DESC LIMIT 1)",
}

func (t *task) Export(ctx context.Context, conn repository.DBConn, filter repository.TaskListFilter) ([]repository.TaskExportRow, error) {
	var dbts []entity.Task

	sel := conn.NewSelect().
		Model(&dbts).
		Limit(repository.MaxTaskBatchSize + 1).
		OrderExpr("t.created_at ASC")

	if len(filter.Statuses) > 0 {
		sel = sel.Where("t.status IN (?)", bun.In(filter.Statuses))
	}
	if len(filter.Priorities) > 0 {
		sel = sel.Where("t.priority IN (?)", bun.In(filter.Priorities))
	}
	if len(filter.Difficulties) > 0 {
		sel = sel.Where("t.difficulty IN (?)", bun.In(filter.Difficulties))
	}
	if filter.RobotType != nil {
		sel = sel.Where("t.robot_type = ?", *filter.RobotType)
	}

	if err := sel.Scan(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to export tasks: %v", err))
	}

	if len(dbts) == 0 {
		return []repository.TaskExportRow{}, nil
	}

	taskIDs := make([]string, 0, len(dbts))
	for _, d := range dbts {
		taskIDs = append(taskIDs, d.IDNatural)
	}

	// Get latest approved version per task
	var approvedVersions []entity.TaskVersion
	if err := conn.NewSelect().
		Model(&approvedVersions).
		Where("task_id IN (?)", bun.In(taskIDs)).
		Where("approval_status = ?", model.ApprovalStatusApproved).
		Order("created_at DESC").
		Scan(ctx); err != nil && err != sql.ErrNoRows {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get approved versions for export: %v", err))
	}

	latestApprovedMap := make(map[string]entity.TaskVersion)
	for _, tv := range approvedVersions {
		if _, exists := latestApprovedMap[tv.TaskID]; !exists {
			latestApprovedMap[tv.TaskID] = tv
		}
	}

	// Collect approved version IDs for subtask query
	approvedVersionIDs := make([]string, 0, len(latestApprovedMap))
	for _, tv := range latestApprovedMap {
		approvedVersionIDs = append(approvedVersionIDs, tv.IDNatural)
	}

	// Get subtasks for approved versions, sorted by order_index
	subtasksByVersionID := make(map[string][]string)
	if len(approvedVersionIDs) > 0 {
		var subtasks []entity.SubTask
		if err := conn.NewSelect().
			Model(&subtasks).
			Where("task_version_id IN (?)", bun.In(approvedVersionIDs)).
			Order("task_version_id ASC", "order_index ASC").
			Scan(ctx); err != nil && err != sql.ErrNoRows {
			return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get subtasks for export: %v", err))
		}
		for _, st := range subtasks {
			subtasksByVersionID[st.TaskVersionID] = append(subtasksByVersionID[st.TaskVersionID], st.Name)
		}
	}

	rows := make([]repository.TaskExportRow, 0, len(dbts))
	for _, d := range dbts {
		row := repository.TaskExportRow{
			IDNatural:   d.IDNatural,
			Name:        d.Name,
			Description: d.Description,
			ManualURL:   derefString(d.ManualURL),
			Priority:    d.Priority,
			Difficulty:  d.Difficulty,
			Status:      d.Status,
			Deadline:    d.Deadline,
			RobotType:   d.RobotType,
		}
		if tv, ok := latestApprovedMap[d.IDNatural]; ok {
			row.SubtaskNames = subtasksByVersionID[tv.IDNatural]
			row.TargetDurationSeconds = tv.TargetDurationSeconds
			row.TargetEpisodeCount = tv.TargetEpisodeCount
			row.TargetDurationPerEpisodeSeconds = tv.TargetDurationPerEpisodeSeconds
		}
		rows = append(rows, row)
	}

	return rows, nil
}

var nullableSortColumns = map[string]bool{
	"robot_type":              true,
	"target_duration_seconds": true,
}

func applyTaskSortOrder(sel *bun.SelectQuery, sortBy *repository.TaskSortBy, sortOrder *repository.SortOrder) *bun.SelectQuery {
	if sortBy == nil {
		return sel.OrderExpr("t.created_at DESC")
	}

	if *sortBy == repository.TaskSortByRecommended {
		return sel.OrderExpr("t.priority DESC, t.deadline ASC NULLS LAST")
	}

	col, ok := allowedTaskSortColumns[string(*sortBy)]
	if !ok {
		return sel.OrderExpr("t.created_at DESC")
	}

	order := "ASC"
	if sortOrder != nil && *sortOrder == repository.SortOrderDesc {
		order = "DESC"
	}

	nullsClause := ""
	if nullableSortColumns[string(*sortBy)] {
		nullsClause = " NULLS LAST"
	}

	return sel.OrderExpr(fmt.Sprintf("%s %s%s", col, order, nullsClause))
}

// applyTaskSummaryFilter applies common filters for task summary queries.
func applyTaskSummaryFilter(sel *bun.SelectQuery, filter repository.TaskSummaryFilter) *bun.SelectQuery {
	if len(filter.RobotTypes) > 0 {
		sel = sel.Where("t.robot_type IN (?)", bun.In(filter.RobotTypes))
	}
	if len(filter.TagIDs) > 0 {
		sel = sel.Where(`EXISTS (
			SELECT 1 FROM task_tag_assignment
			WHERE task_id = t.id_natural AND tag_id IN (?)
		)`, bun.In(filter.TagIDs))
	}
	if filter.CategoryTypeID != nil {
		sel = sel.Where(`EXISTS (
			SELECT 1 FROM task_tag_assignment
			JOIN task_tag ON task_tag.id = task_tag_assignment.tag_id
			WHERE task_tag_assignment.task_id = t.id_natural
			AND task_tag.category_type_id = ?
		)`, *filter.CategoryTypeID)
	}
	if filter.DeadlineFrom != nil {
		sel = sel.Where("t.deadline >= ?", *filter.DeadlineFrom)
	}
	if filter.DeadlineTo != nil {
		sel = sel.Where("t.deadline < ?", *filter.DeadlineTo)
	}
	return sel
}

func (t *task) GetFilteredTasks(ctx context.Context, conn repository.DBConn, filter repository.TaskSummaryFilter) ([]repository.FilteredTask, error) {
	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "organization context required"))
	}

	var rows []repository.FilteredTask
	sel := conn.NewSelect().
		TableExpr("task AS t").
		ColumnExpr("t.id_natural, t.deadline, t.status").
		Where("t.organization_id = ?", orgID)

	sel = applyTaskSummaryFilter(sel, filter)

	if err := sel.Scan(ctx, &rows); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get filtered tasks: %v", err))
	}
	return rows, nil
}

func (t *task) GetTargetsByTaskIDs(ctx context.Context, conn repository.DBConn, taskIDs []string) (map[string]repository.TaskTargets, error) {
	if len(taskIDs) == 0 {
		return map[string]repository.TaskTargets{}, nil
	}

	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "organization context required"))
	}

	var rows []repository.TaskTargets
	err = conn.NewSelect().
		TableExpr("task_version").
		ColumnExpr("task_id").
		ColumnExpr("COALESCE(SUM(COALESCE(target_duration_seconds, 0)), 0) AS target_duration").
		ColumnExpr("COALESCE(SUM(COALESCE(target_episode_count, 0)), 0) AS target_episodes").
		Where("organization_id = ?", orgID).
		Where("task_id IN (?)", bun.In(taskIDs)).
		Where("approval_status = ?", model.ApprovalStatusApproved).
		GroupExpr("task_id").
		Scan(ctx, &rows)

	if err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get targets by task IDs: %v", err))
	}

	result := make(map[string]repository.TaskTargets, len(rows))
	for _, r := range rows {
		result[r.TaskID] = r
	}
	return result, nil
}

func (t *task) FindExistingNames(ctx context.Context, conn repository.DBConn, names []string) (map[string]bool, error) {
	if len(names) == 0 {
		return map[string]bool{}, nil
	}

	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "organization context required"))
	}

	var rows []struct {
		Name string `bun:"name"`
	}
	if err := conn.NewSelect().
		TableExpr("task AS t").
		ColumnExpr("t.name").
		Where("t.organization_id = ?", orgID).
		Where("t.name IN (?)", bun.In(names)).
		Scan(ctx, &rows); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to find existing task names: %v", err))
	}

	result := make(map[string]bool, len(rows))
	for _, r := range rows {
		result[r.Name] = true
	}
	return result, nil
}

func (t *task) BulkCreate(ctx context.Context, conn repository.DBConn, items []repository.BulkTaskItem) ([]model.Task, error) {
	if len(items) == 0 {
		return nil, nil
	}

	tasks := make([]model.Task, 0, len(items))
	dbTasks := make([]entity.Task, 0, len(items))
	for _, item := range items {
		tasks = append(tasks, item.Task)
		dbTasks = append(dbTasks, entity.Task{
			IDNatural:      item.Task.IDNatural,
			OrganizationID: item.Task.OrganizationID,
			Name:           item.Task.Name,
			Description:    item.Task.Description,
			ManualURL:      &item.Task.ManualURL,
			Priority:       derefPriority(item.Task.Priority),
			Difficulty:     derefDifficulty(item.Task.Difficulty),
			Status:         derefTaskStatus(item.Task.Status),
			Deadline:       item.Task.Deadline,
			RobotType:      item.Task.RobotType,
		})
	}

	if _, err := conn.NewInsert().
		Model(&dbTasks).
		Exec(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to bulk create tasks: %v", err))
	}

	// Create task_versions with optional target fields
	dbVersions := make([]entity.TaskVersion, 0, len(items))
	versionIDByTaskID := make(map[string]string, len(items))
	for _, item := range items {
		tvID, err := model.InitID()
		if err != nil {
			return nil, err
		}
		versionIDByTaskID[item.Task.IDNatural] = tvID
		dbVersions = append(dbVersions, entity.TaskVersion{
			IDNatural:                       tvID,
			TaskID:                          item.Task.IDNatural,
			OrganizationID:                  item.Task.OrganizationID,
			Version:                         model.InitialVersion,
			IsActive:                        true,
			ApprovalStatus:                  model.ApprovalStatusDraft,
			TargetDurationSeconds:           item.TargetDurationSeconds,
			TargetEpisodeCount:              item.TargetEpisodeCount,
			TargetDurationPerEpisodeSeconds: item.TargetDurationPerEpisodeSeconds,
		})
	}

	if _, err := conn.NewInsert().
		Model(&dbVersions).
		Exec(ctx); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to bulk create task versions: %v", err))
	}

	// Create subtasks
	var dbSubtasks []entity.SubTask
	for _, item := range items {
		if len(item.SubtaskNames) == 0 {
			continue
		}
		tvID := versionIDByTaskID[item.Task.IDNatural]
		for i, name := range item.SubtaskNames {
			stID, err := model.InitID()
			if err != nil {
				return nil, err
			}
			dbSubtasks = append(dbSubtasks, entity.SubTask{
				IDNatural:      stID,
				OrganizationID: item.Task.OrganizationID,
				TaskVersionID:  tvID,
				Name:           name,
				OrderIndex:     i,
			})
		}
	}

	if len(dbSubtasks) > 0 {
		if _, err := conn.NewInsert().
			Model(&dbSubtasks).
			Exec(ctx); err != nil {
			return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to bulk create subtasks: %v", err))
		}
	}

	return tasks, nil
}

func (t *task) GetActualsByTaskIDs(ctx context.Context, conn repository.DBConn, taskIDs []string) (map[string]repository.TaskActuals, error) {
	if len(taskIDs) == 0 {
		return map[string]repository.TaskActuals{}, nil
	}

	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil || orgID == "" {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "organization context required"))
	}

	var rows []repository.TaskActuals
	err = conn.NewSelect().
		TableExpr("task_version AS tv").
		Join("JOIN task_version_stats tvs ON tvs.task_version_id = tv.id_natural").
		ColumnExpr("tv.task_id").
		ColumnExpr("COALESCE(SUM(tvs.total_duration_seconds), 0) AS actual_duration").
		ColumnExpr("COALESCE(SUM(tvs.episode_count), 0) AS actual_episodes").
		Where("tv.organization_id = ?", orgID).
		Where("tv.task_id IN (?)", bun.In(taskIDs)).
		Where("tv.approval_status = ?", model.ApprovalStatusApproved).
		GroupExpr("tv.task_id").
		Scan(ctx, &rows)

	if err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get actuals by task IDs: %v", err))
	}

	result := make(map[string]repository.TaskActuals, len(rows))
	for _, r := range rows {
		result[r.TaskID] = r
	}
	return result, nil
}
