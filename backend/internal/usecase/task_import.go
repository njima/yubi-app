package usecase

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/requestctx"
)

const maxImportFileLen = 5 * 1024 * 1024 // 5MB

// TaskImportUsecase handles CSV-based task import.
type TaskImportUsecase interface {
	Validate(ctx context.Context, csvContent string) (TaskImportValidationResult, error)
	Import(ctx context.Context, csvContent string) (TaskImportResult, error)
}

type TaskImportValidationResult struct {
	ValidRows     []TaskImportRow
	DuplicateRows []TaskImportRowError
	ErrorRows     []TaskImportRowError
}

type TaskImportResult struct {
	ImportedCount int
	SkippedCount  int
	ErrorCount    int
	Errors        []TaskImportRowError
}

type TaskImportRow struct {
	RowNumber                       int
	Name                            string
	Subtasks                        []string // up to 10
	Description                     string
	ManualURL                       string
	Priority                        string
	Difficulty                      string
	Status                          string
	Deadline                        string
	RobotType                       string
	Tags                            string
	TargetDurationSeconds           *int
	TargetEpisodeCount              *int
	CategoryTags                    string
	TargetDurationPerEpisodeSeconds *int
}

type TaskImportRowError struct {
	RowNumber int
	Errors    []string
	Name      string
}

type taskImport struct {
	taskRepo repository.Task
	tagRepo  repository.TaskTag
	data     repository.DataAccess
}

func NewTaskImport(taskRepo repository.Task, tagRepo repository.TaskTag, data repository.DataAccess) *taskImport {
	return &taskImport{taskRepo: taskRepo, tagRepo: tagRepo, data: data}
}

// expectedHeaders defines the required CSV column order.
var expectedHeaders = []string{
	"name", "subtasks",
	"description", "manual_url", "priority", "difficulty",
	"status", "deadline", "robot_type", "tags",
	"target_duration", "target_episode_count", "category_tags", "target_duration_per_episode",
}

// validateInternal classifies already-parsed rows into valid/duplicate/error buckets.
// It fetches tags and existing names once, returning the tagNameToID map so Import
// can reuse it without a second DB round-trip.
func (u *taskImport) validateInternal(
	ctx context.Context,
	rows []TaskImportRow,
	parseErrors []TaskImportRowError,
) (TaskImportValidationResult, map[string]string, error) {
	// Duplicate-name detection against the DB
	names := make([]string, 0, len(rows))
	for _, r := range rows {
		names = append(names, r.Name)
	}
	existingNames, err := u.taskRepo.FindExistingNames(ctx, u.data.Conn(), names)
	if err != nil {
		return TaskImportValidationResult{}, nil, err
	}

	// Fetch only the tag names actually referenced in the CSV (not all tags)
	tagNames := collectTagNames(rows)
	allTags, err := u.tagRepo.GetTagsByNames(ctx, u.data.Conn(), tagNames)
	if err != nil {
		return TaskImportValidationResult{}, nil, err
	}
	tagNameSet := make(map[string]bool, len(allTags))
	tagNameToID := make(map[string]string, len(allTags))
	for _, t := range allTags {
		tagNameSet[t.Name] = true
		tagNameToID[t.Name] = t.ID
	}

	var result TaskImportValidationResult
	seenNames := make(map[string]bool)

	for _, row := range rows {
		rowErrors := validateImportRow(row, tagNameSet)
		if len(rowErrors) > 0 {
			result.ErrorRows = append(result.ErrorRows, TaskImportRowError{
				RowNumber: row.RowNumber,
				Errors:    rowErrors,
				Name:      row.Name,
			})
			continue
		}

		if existingNames[row.Name] || seenNames[row.Name] {
			result.DuplicateRows = append(result.DuplicateRows, TaskImportRowError{
				RowNumber: row.RowNumber,
				Errors:    []string{fmt.Sprintf("task name %q already exists", row.Name)},
				Name:      row.Name,
			})
			continue
		}

		seenNames[row.Name] = true
		result.ValidRows = append(result.ValidRows, row)
	}

	// Append parse-level errors
	result.ErrorRows = append(result.ErrorRows, parseErrors...)

	return result, tagNameToID, nil
}

// collectTagNames extracts all unique tag names referenced in the CSV rows
// (from both "tags" and "category_tags" columns).
func collectTagNames(rows []TaskImportRow) []string {
	seen := make(map[string]struct{})
	var names []string
	for _, row := range rows {
		for _, raw := range []string{row.Tags, row.CategoryTags} {
			if raw == "" {
				continue
			}
			for _, tn := range strings.Split(raw, ";") {
				tn = strings.TrimSpace(tn)
				if tn != "" {
					if _, ok := seen[tn]; !ok {
						seen[tn] = struct{}{}
						names = append(names, tn)
					}
				}
			}
		}
	}
	return names
}

func (u *taskImport) Validate(ctx context.Context, csvContent string) (TaskImportValidationResult, error) {
	rows, parseErrors, err := u.parseCSV(csvContent)
	if err != nil {
		return TaskImportValidationResult{}, err
	}

	result, _, err := u.validateInternal(ctx, rows, parseErrors)
	return result, err
}

func (u *taskImport) Import(ctx context.Context, csvContent string) (TaskImportResult, error) {
	orgID, err := requestctx.OrganizationID(ctx)
	if err != nil {
		return TaskImportResult{}, err
	}

	rows, parseErrors, err := u.parseCSV(csvContent)
	if err != nil {
		return TaskImportResult{}, err
	}

	// Single pass: validate + build tagNameToID (no second DB round-trip)
	validation, tagNameToID, err := u.validateInternal(ctx, rows, parseErrors)
	if err != nil {
		return TaskImportResult{}, err
	}

	// No valid rows: return duplicates and errors only
	if len(validation.ValidRows) == 0 {
		return TaskImportResult{
			ImportedCount: 0,
			SkippedCount:  len(validation.DuplicateRows),
			ErrorCount:    len(validation.ErrorRows),
			Errors:        validation.ErrorRows,
		}, nil
	}

	// Convert valid rows to BulkTaskItems
	items := make([]repository.BulkTaskItem, 0, len(validation.ValidRows))
	tagIDsByTask := make(map[string][]string)

	for _, row := range validation.ValidRows {
		priority := parsePriority(row.Priority)
		difficulty := parseDifficulty(row.Difficulty)
		status := parseStatus(row.Status)
		deadline := parseDeadline(row.Deadline)

		tk, err := model.InitTask(orgID, row.Name, strPtrOrNil(row.Description), row.ManualURL, &priority, &difficulty, &status, deadline, strPtrOrNil(row.RobotType))
		if err != nil {
			return TaskImportResult{}, err
		}

		// Collect tag IDs using the already-fetched tagNameToID map
		var tagIDs []string
		for _, rawTags := range []string{row.Tags, row.CategoryTags} {
			if rawTags == "" {
				continue
			}
			for _, tn := range strings.Split(rawTags, ";") {
				tn = strings.TrimSpace(tn)
				if id, ok := tagNameToID[tn]; ok {
					tagIDs = append(tagIDs, id)
				}
			}
		}
		if len(tagIDs) > 0 {
			tagIDsByTask[tk.IDNatural] = tagIDs
		}

		items = append(items, repository.BulkTaskItem{
			Task:                            tk,
			SubtaskNames:                    row.Subtasks,
			TargetDurationSeconds:           row.TargetDurationSeconds,
			TargetEpisodeCount:              row.TargetEpisodeCount,
			TargetDurationPerEpisodeSeconds: row.TargetDurationPerEpisodeSeconds,
		})
	}

	// Bulk create in transaction
	var createdTasks []model.Task
	err = u.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		var err error
		createdTasks, err = u.taskRepo.BulkCreate(ctx, conn, items)
		if err != nil {
			return err
		}
		for _, tk := range createdTasks {
			if ids, ok := tagIDsByTask[tk.IDNatural]; ok {
				if err := u.tagRepo.SetTaskTags(ctx, conn, tk.IDNatural, ids); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return TaskImportResult{}, err
	}

	return TaskImportResult{
		ImportedCount: len(createdTasks),
		SkippedCount:  len(validation.DuplicateRows),
		ErrorCount:    len(validation.ErrorRows),
		Errors:        validation.ErrorRows,
	}, nil
}

// parseCSV streams the CSV row-by-row to avoid loading all records into memory at once.
// It stops as soon as the row limit is reached, preventing DoS via large payloads.
func (u *taskImport) parseCSV(csvContent string) ([]TaskImportRow, []TaskImportRowError, error) {
	if len(csvContent) > maxImportFileLen {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "CSV content exceeds maximum size of 5MB"))
	}

	// Strip UTF-8 BOM if present
	csvContent = strings.TrimPrefix(csvContent, "\xef\xbb\xbf")

	if !utf8.ValidString(csvContent) {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "CSV content is not valid UTF-8"))
	}

	reader := csv.NewReader(strings.NewReader(csvContent))
	reader.TrimLeadingSpace = true

	// Read and validate header row
	header, err := reader.Read()
	if err == io.EOF {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "CSV must have a header row and at least one data row"))
	}
	if err != nil {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "failed to parse CSV: %v", err))
	}
	if len(header) != len(expectedHeaders) {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest,
			"CSV header must have exactly %d columns: %s", len(expectedHeaders), strings.Join(expectedHeaders, ", ")))
	}
	for i, h := range header {
		if strings.TrimSpace(strings.ToLower(h)) != expectedHeaders[i] {
			return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest,
				"CSV column %d must be %q, got %q", i+1, expectedHeaders[i], h))
		}
	}

	var rows []TaskImportRow
	var parseErrors []TaskImportRowError
	rowNum := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "failed to parse CSV: %v", err))
		}
		rowNum++
		if rowNum > repository.MaxTaskBatchSize {
			return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest,
				"CSV exceeds maximum of %d rows", repository.MaxTaskBatchSize))
		}

		if len(record) != len(expectedHeaders) {
			parseErrors = append(parseErrors, TaskImportRowError{
				RowNumber: rowNum,
				Errors:    []string{fmt.Sprintf("expected %d columns, got %d", len(expectedHeaders), len(record))},
			})
			continue
		}

		// columns: 0=name, 1=subtasks, 2=description, 3=manual_url,
		// 4=priority, 5=difficulty, 6=status, 7=deadline, 8=robot_type,
		// 9=tags, 10=target_duration, 11=target_episode_count,
		// 12=category_tags, 13=target_duration_per_episode

		// Parse semicolon-separated subtask names
		var subtasks []string
		if raw := strings.TrimSpace(record[1]); raw != "" {
			for _, s := range strings.Split(raw, ";") {
				if v := strings.TrimSpace(s); v != "" {
					subtasks = append(subtasks, v)
				}
			}
		}

		row := TaskImportRow{
			RowNumber:    rowNum,
			Name:         strings.TrimSpace(record[0]),
			Subtasks:     subtasks,
			Description:  strings.TrimSpace(record[2]),
			ManualURL:    strings.TrimSpace(record[3]),
			Priority:     strings.TrimSpace(record[4]),
			Difficulty:   strings.TrimSpace(record[5]),
			Status:       strings.TrimSpace(record[6]),
			Deadline:     strings.TrimSpace(record[7]),
			RobotType:    strings.TrimSpace(record[8]),
			Tags:         strings.TrimSpace(record[9]),
			CategoryTags: strings.TrimSpace(record[12]),
		}

		hasError := false
		type intField struct {
			col  int
			name string
			dest **int
		}
		for _, f := range []intField{
			{10, "target_duration", &row.TargetDurationSeconds},
			{11, "target_episode_count", &row.TargetEpisodeCount},
			{13, "target_duration_per_episode", &row.TargetDurationPerEpisodeSeconds},
		} {
			if v := strings.TrimSpace(record[f.col]); v != "" {
				n, err := strconv.Atoi(v)
				if err != nil || n <= 0 {
					parseErrors = append(parseErrors, TaskImportRowError{
						RowNumber: rowNum,
						Errors:    []string{fmt.Sprintf("%s %q is invalid (must be a positive integer)", f.name, v)},
						Name:      strings.TrimSpace(record[0]),
					})
					hasError = true
					break
				}
				*f.dest = &n
			}
		}
		if hasError {
			continue
		}

		rows = append(rows, row)
	}

	if rowNum == 0 {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "CSV must have a header row and at least one data row"))
	}

	return rows, parseErrors, nil
}

func validateImportRow(row TaskImportRow, validTags map[string]bool) []string {
	var errs []string

	if row.Name == "" {
		errs = append(errs, "name is required")
	} else if len([]rune(row.Name)) > 100 {
		errs = append(errs, "name must be 100 characters or less")
	}

	if len(row.Subtasks) > 10 {
		errs = append(errs, "subtasks must be 10 or fewer")
	}

	for i, st := range row.Subtasks {
		if len([]rune(st)) > 100 {
			errs = append(errs, fmt.Sprintf("subtask%d name must be 100 characters or less", i+1))
		}
	}

	if row.ManualURL == "" {
		errs = append(errs, "manual_url is required")
	} else if !strings.HasPrefix(row.ManualURL, "https://") {
		errs = append(errs, "manual_url must start with https://")
	}

	if row.Priority == "" {
		errs = append(errs, "priority is required")
	} else if !isValidPriority(row.Priority) {
		errs = append(errs, fmt.Sprintf("priority %q is invalid (must be Low, Normal, High, or Urgent)", row.Priority))
	}

	if row.Difficulty == "" {
		errs = append(errs, "difficulty is required")
	} else if !isValidDifficulty(row.Difficulty) {
		errs = append(errs, fmt.Sprintf("difficulty %q is invalid (must be S, A, B, or C)", row.Difficulty))
	}

	if row.Status != "" && !isValidStatus(row.Status) {
		errs = append(errs, fmt.Sprintf("status %q is invalid (must be Planning, Doing, Completed, or Canceled)", row.Status))
	}

	if row.Deadline != "" && parseDeadline(row.Deadline).IsZero() {
		errs = append(errs, fmt.Sprintf("deadline %q is invalid (use YYYY-MM-DD, YYYY/MM/DD, or ISO 8601 format)", row.Deadline))
	}

	for _, rawTags := range []string{row.Tags, row.CategoryTags} {
		if rawTags == "" {
			continue
		}
		for _, tn := range strings.Split(rawTags, ";") {
			tn = strings.TrimSpace(tn)
			if tn == "" {
				continue
			}
			if !validTags[tn] {
				errs = append(errs, fmt.Sprintf("tag %q does not exist", tn))
			}
		}
	}

	return errs
}

var priorityMap = map[string]model.TaskPriority{
	"low":    model.TaskPriorityLow,
	"normal": model.TaskPriorityNormal,
	"high":   model.TaskPriorityHigh,
	"urgent": model.TaskPriorityUrgent,
}

var difficultyMap = map[string]model.TaskDifficulty{
	"s": model.TaskDifficultyS,
	"a": model.TaskDifficultyA,
	"b": model.TaskDifficultyB,
	"c": model.TaskDifficultyC,
}

var statusMap = map[string]model.TaskStatus{
	"planning":  model.TaskStatusPlanning,
	"doing":     model.TaskStatusDoing,
	"completed": model.TaskStatusCompleted,
	"canceled":  model.TaskStatusCanceled,
}

func isValidPriority(s string) bool {
	_, ok := priorityMap[strings.ToLower(s)]
	return ok
}

func isValidDifficulty(s string) bool {
	_, ok := difficultyMap[strings.ToLower(s)]
	return ok
}

func isValidStatus(s string) bool {
	_, ok := statusMap[strings.ToLower(s)]
	return ok
}

func parsePriority(s string) model.TaskPriority {
	if v, ok := priorityMap[strings.ToLower(s)]; ok {
		return v
	}
	return model.TaskPriorityNormal
}

func parseDifficulty(s string) model.TaskDifficulty {
	if v, ok := difficultyMap[strings.ToLower(s)]; ok {
		return v
	}
	return model.TaskDifficultyB
}

func parseStatus(s string) model.TaskStatus {
	if v, ok := statusMap[strings.ToLower(s)]; ok {
		return v
	}
	return model.TaskStatusPlanning
}

func parseDeadline(s string) time.Time {
	for _, layout := range []string{
		"2006-01-02",
		"2006/01/02",
		time.RFC3339,
		"2006-01-02T15:04:05",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
