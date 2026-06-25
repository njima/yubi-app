package usecase

import (
	"bytes"
	"context"
	"encoding/csv"
	"strconv"
	"strings"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

// TaskExportUsecase generates a CSV export of tasks.
type TaskExportUsecase interface {
	Export(ctx context.Context, filter TaskListFilter) ([]byte, error)
}

type taskExport struct {
	taskRepo repository.Task
	tagRepo  repository.TaskTag
	data     repository.DataAccess
}

func NewTaskExport(taskRepo repository.Task, tagRepo repository.TaskTag, data repository.DataAccess) *taskExport {
	return &taskExport{taskRepo: taskRepo, tagRepo: tagRepo, data: data}
}

// exportHeaders must match expectedHeaders in task_import.go exactly so that
// an exported CSV can be re-imported without modification.
// "category_tags" is intentionally kept as an empty column: the import side
// reads both "tags" and "category_tags", but the export puts all tags into
// "tags" only, preserving round-trip compatibility.
var exportHeaders = []string{
	"name", "subtasks",
	"description", "manual_url", "priority", "difficulty",
	"status", "deadline", "robot_type", "tags",
	"target_duration", "target_episode_count", "category_tags", "target_duration_per_episode",
}

func (u *taskExport) Export(ctx context.Context, filter TaskListFilter) ([]byte, error) {
	rows, err := u.taskRepo.Export(ctx, u.data.Conn(), filter.repositoryFilter())
	if err != nil {
		return nil, err
	}

	if len(rows) > repository.MaxTaskBatchSize {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest,
			"export limit of %d tasks exceeded; apply filters to reduce the result set", repository.MaxTaskBatchSize))
	}

	taskIDs := make([]string, 0, len(rows))
	for _, r := range rows {
		taskIDs = append(taskIDs, r.IDNatural)
	}

	tagsByTaskID, err := u.tagRepo.GetTagsByTaskIDs(ctx, u.data.Conn(), taskIDs)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write(exportHeaders); err != nil {
		return nil, err
	}

	for _, row := range rows {
		tags := tagsByTaskID[row.IDNatural]
		tagNames := make([]string, 0, len(tags))
		for _, t := range tags {
			tagNames = append(tagNames, t.Name)
		}

		record := []string{
			row.Name,
			strings.Join(row.SubtaskNames, ";"),
			derefStringExport(row.Description),
			row.ManualURL,
			reverseLookup(priorityMap, row.Priority, "Normal"),
			reverseLookup(difficultyMap, row.Difficulty, "B"),
			statusLabel(int(row.Status)),
			deadlineToString(row.Deadline),
			derefStringExport(row.RobotType),
			strings.Join(tagNames, ";"),
			intPtrToString(row.TargetDurationSeconds),
			intPtrToString(row.TargetEpisodeCount),
			"", // category_tags: intentionally empty (all tags in "tags" column)
			intPtrToString(row.TargetDurationPerEpisodeSeconds),
		}

		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func deadlineToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func derefStringExport(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func intPtrToString(v *int) string {
	if v == nil {
		return ""
	}
	return strconv.Itoa(*v)
}

// reverseLookup returns the map key whose value equals val, with the first
// character uppercased (e.g. "low" → "Low", "s" → "S").
// Falls back to fallback if val is not found.
func reverseLookup[V comparable](m map[string]V, val V, fallback string) string {
	for k, v := range m {
		if v == val {
			return strings.ToUpper(k[:1]) + k[1:]
		}
	}
	return fallback
}
