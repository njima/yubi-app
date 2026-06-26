package usecase

import (
	"context"
	"sort"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/usecase/pagination"
)

type TaskUsecase interface {
	Create(ctx context.Context, input TaskCreateInput) (model.Task, error)
	GetByID(ctx context.Context, id string) (model.Task, error)
	ListByIDs(ctx context.Context, ids []string) (model.Tasks, error)
	List(ctx context.Context, filter TaskListFilter, page, limit int) (model.Tasks, int, error)
	Update(ctx context.Context, input TaskUpdateInput) (model.Task, error)
	Delete(ctx context.Context, id string) error
	GetSummary(ctx context.Context, filter TaskSummaryFilter) (model.TaskSummary, error)
	GetCompletionTrend(ctx context.Context, filter TaskSummaryFilter, groupBy string, from, to time.Time, interval string) (model.TaskCompletionTrend, error)
}

type TaskCreateInput struct {
	OrganizationID string
	Name           string
	LocationID     string
	Description    *string
	ManualURL      string
	Priority       model.TaskPriority
	Difficulty     model.TaskDifficulty
	Status         model.TaskStatus
	Deadline       time.Time
	RobotType      *string
	TagIDs         []string
}

type TaskUpdateInput struct {
	ID          string
	Name        *string
	Description *string
	ManualURL   *string
	Priority    *model.TaskPriority
	Difficulty  *model.TaskDifficulty
	Status      *model.TaskStatus
	Deadline    *time.Time
	RobotType   *string
	TagIDs      *[]string
}

type task struct {
	repo        repository.Task
	tagRepo     repository.TaskTag
	episodeRepo repository.Episode
	tvRepo      repository.TaskVersion
	data        repository.DataAccess
}

func NewTask(repo repository.Task, tagRepo repository.TaskTag, episodeRepo repository.Episode, tvRepo repository.TaskVersion, data repository.DataAccess) *task {
	return &task{repo: repo, tagRepo: tagRepo, episodeRepo: episodeRepo, tvRepo: tvRepo, data: data}
}

func deduplicateTagIDs(tagIDs []string) []string {
	seen := make(map[string]struct{}, len(tagIDs))
	result := make([]string, 0, len(tagIDs))
	for _, id := range tagIDs {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			result = append(result, id)
		}
	}
	return result
}

func (t *task) Create(ctx context.Context, input TaskCreateInput) (model.Task, error) {
	tk, err := model.InitTask(input.OrganizationID, input.Name, input.Description, input.ManualURL, &input.Priority, &input.Difficulty, &input.Status, input.Deadline, input.RobotType)
	if err != nil {
		return model.Task{}, err
	}

	input.TagIDs = deduplicateTagIDs(input.TagIDs)

	var result model.Task
	err = t.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		var txErr error
		result, txErr = t.repo.Create(ctx, conn, tk)
		if txErr != nil {
			return txErr
		}
		if err := t.tagRepo.SetTaskTags(ctx, conn, result.IDNatural, input.TagIDs); err != nil {
			return err
		}
		tags, err := t.tagRepo.GetTagsByTaskID(ctx, conn, result.IDNatural)
		if err != nil {
			return err
		}
		result.Tags = tags
		return nil
	})

	if err != nil {
		return model.Task{}, err
	}

	return result, nil
}

func (t *task) GetByID(ctx context.Context, id string) (model.Task, error) {
	tk, err := t.repo.GetByID(ctx, t.data.Conn(), id)
	if err != nil {
		return model.Task{}, err
	}
	tags, err := t.tagRepo.GetTagsByTaskID(ctx, t.data.Conn(), tk.IDNatural)
	if err != nil {
		return model.Task{}, err
	}
	tk.Tags = tags
	return tk, nil
}

func (t *task) ListByIDs(ctx context.Context, ids []string) (model.Tasks, error) {
	return t.repo.ListByIDs(ctx, t.data.Conn(), ids)
}

func (t *task) List(ctx context.Context, filter TaskListFilter, page, limit int) (model.Tasks, int, error) {
	if limit <= 0 {
		limit = pagination.DefaultLimit
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	tasks, total, err := t.repo.List(ctx, t.data.Conn(), filter.repositoryFilter(), limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if len(tasks) == 0 {
		return tasks, total, nil
	}
	ids := make([]string, 0, len(tasks))
	for _, tk := range tasks {
		ids = append(ids, tk.IDNatural)
	}
	tagsByTask, err := t.tagRepo.GetTagsByTaskIDs(ctx, t.data.Conn(), ids)
	if err != nil {
		return nil, 0, err
	}
	for _, tk := range tasks {
		tk.Tags = tagsByTask[tk.IDNatural]
	}
	return tasks, total, nil
}

func (t *task) Update(ctx context.Context, input TaskUpdateInput) (model.Task, error) {
	tk := model.Task{IDNatural: input.ID}
	if input.Name != nil {
		tk.Name = *input.Name
	}
	if input.Description != nil {
		tk.Description = input.Description
	}
	if input.ManualURL != nil {
		if err := tk.SetManualURL(*input.ManualURL); err != nil {
			return model.Task{}, err
		}
	}
	if input.Priority != nil {
		if err := tk.SetPriority(input.Priority); err != nil {
			return model.Task{}, err
		}
	}
	if input.Difficulty != nil {
		if err := tk.SetDifficulty(input.Difficulty); err != nil {
			return model.Task{}, err
		}
	}
	// Status is handled inside RunInTx (Canceled release needs auto-determination)
	if input.Deadline != nil {
		if err := tk.SetDeadline(*input.Deadline); err != nil {
			return model.Task{}, err
		}
	}
	if input.RobotType != nil {
		tk.SetRobotType(input.RobotType)
	}
	if input.TagIDs != nil {
		deduped := deduplicateTagIDs(*input.TagIDs)
		input.TagIDs = &deduped
	}

	var result model.Task
	err := t.data.RunInTx(ctx, func(ctx context.Context, txData repository.DataAccess) error {
		conn := txData.Conn()
		// Handle status change inside transaction
		if input.Status != nil {
			if *input.Status == model.TaskStatusCanceled {
				// Any → Canceled: set directly
				if err := tk.SetStatus(input.Status); err != nil {
					return err
				}
			} else {
				// Check if uncanceling
				currentTask, err := t.repo.GetByID(ctx, conn, input.ID)
				if err != nil {
					return err
				}
				if currentTask.Status != nil && *currentTask.Status == model.TaskStatusCanceled {
					// Canceled → non-Canceled: auto-determine correct status
					actual, err := t.episodeRepo.SumDurationByTaskID(ctx, conn, currentTask.IDNatural)
					if err != nil {
						return err
					}
					target, err := t.tvRepo.SumTargetByTaskID(ctx, conn, currentTask.IDNatural)
					if err != nil {
						return err
					}
					correctStatus := model.DetermineTaskStatus(actual, target)
					if err := tk.SetStatus(&correctStatus); err != nil {
						return err
					}
				} else {
					// Non-Canceled → non-Canceled: set as requested (backward compatible)
					if err := tk.SetStatus(input.Status); err != nil {
						return err
					}
				}
			}
		}

		var txErr error
		result, txErr = t.repo.Update(ctx, conn, tk)
		if txErr != nil {
			return txErr
		}
		if input.TagIDs != nil {
			if err := t.tagRepo.SetTaskTags(ctx, conn, result.IDNatural, *input.TagIDs); err != nil {
				return err
			}
		}
		tags, err := t.tagRepo.GetTagsByTaskID(ctx, conn, result.IDNatural)
		if err != nil {
			return err
		}
		result.Tags = tags
		return nil
	})
	if err != nil {
		return model.Task{}, err
	}
	return result, nil
}

func (t *task) Delete(ctx context.Context, id string) error {
	return t.repo.Delete(ctx, t.data.Conn(), id)
}

func (t *task) GetSummary(ctx context.Context, filter TaskSummaryFilter) (model.TaskSummary, error) {
	// Step 1: Get filtered tasks
	tasks, err := t.repo.GetFilteredTasks(ctx, t.data.Conn(), filter.repositoryFilter())
	if err != nil {
		return model.TaskSummary{}, err
	}
	if len(tasks) == 0 {
		return model.TaskSummary{}, nil
	}

	// Step 2: Get targets
	taskIDs := extractTaskIDs(tasks)
	targets, err := t.repo.GetTargetsByTaskIDs(ctx, t.data.Conn(), taskIDs)
	if err != nil {
		return model.TaskSummary{}, err
	}

	// Step 3: Aggregate
	summary := model.TaskSummary{TotalTasks: len(tasks)}
	for _, tgt := range targets {
		summary.TargetDurationSeconds += tgt.TargetDuration
		summary.TargetEpisodeCount += tgt.TargetEpisodes
	}
	return summary, nil
}

func (t *task) GetCompletionTrend(ctx context.Context, filter TaskSummaryFilter, groupBy string, from, to time.Time, interval string) (model.TaskCompletionTrend, error) {
	// Step 1: Get filtered tasks
	tasks, err := t.repo.GetFilteredTasks(ctx, t.data.Conn(), filter.repositoryFilter())
	if err != nil {
		return model.TaskCompletionTrend{}, err
	}
	if len(tasks) == 0 {
		return model.TaskCompletionTrend{}, nil
	}

	taskIDs := extractTaskIDs(tasks)

	// Step 2: Get targets per task
	targets, err := t.repo.GetTargetsByTaskIDs(ctx, t.data.Conn(), taskIDs)
	if err != nil {
		return model.TaskCompletionTrend{}, err
	}

	// Step 3: Get actuals per task
	actuals, err := t.repo.GetActualsByTaskIDs(ctx, t.data.Conn(), taskIDs)
	if err != nil {
		return model.TaskCompletionTrend{}, err
	}

	// Step 4: Get tags per task (for category grouping, using existing tagRepo)
	tagsByTask, err := t.tagRepo.GetTagsByTaskIDs(ctx, t.data.Conn(), taskIDs)
	if err != nil {
		return model.TaskCompletionTrend{}, err
	}

	// Step 5: Build trend data
	return buildCompletionTrend(tasks, targets, actuals, tagsByTask, groupBy, from, to, interval), nil
}

func extractTaskIDs(tasks []repository.FilteredTask) []string {
	ids := make([]string, len(tasks))
	for i, tk := range tasks {
		ids[i] = tk.ID
	}
	return ids
}

// overduePeriodKey is a sentinel for tasks with deadline before the range start.
var overduePeriodKey = periodKey{start: time.Time{}, end: time.Time{}}

type periodKey struct {
	start time.Time
	end   time.Time
}

// truncateToDay returns the date with time set to 00:00:00 UTC.
func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// buildTrendPeriods generates period buckets based on from/to/interval.
// interval: "1week" (7 days), "2week" (14 days), "month" (calendar month).
// Returns the period list and a function to assign a deadline to a period.
func buildTrendPeriods(from, to time.Time, interval string) ([]periodKey, func(deadline time.Time) periodKey) {
	fromDay := truncateToDay(from)
	// to is inclusive, so add 1 day for exclusive end
	toExclusive := truncateToDay(to).AddDate(0, 0, 1)

	var periods []periodKey

	switch interval {
	case "month":
		// Calendar month boundaries
		cursor := time.Date(fromDay.Year(), fromDay.Month(), 1, 0, 0, 0, 0, time.UTC)
		for cursor.Before(toExclusive) {
			next := cursor.AddDate(0, 1, 0)
			periods = append(periods, periodKey{start: cursor, end: next})
			cursor = next
		}
	case "1week":
		cursor := fromDay
		for cursor.Before(toExclusive) {
			next := cursor.AddDate(0, 0, 7)
			periods = append(periods, periodKey{start: cursor, end: next})
			cursor = next
		}
	default: // "2week"
		cursor := fromDay
		for cursor.Before(toExclusive) {
			next := cursor.AddDate(0, 0, 14)
			periods = append(periods, periodKey{start: cursor, end: next})
			cursor = next
		}
	}

	assign := func(deadline time.Time) periodKey {
		dl := truncateToDay(deadline)
		if dl.Before(fromDay) {
			return overduePeriodKey
		}
		for _, pk := range periods {
			if !dl.Before(pk.start) && dl.Before(pk.end) {
				return pk
			}
		}
		// Beyond range → last period
		if len(periods) > 0 {
			return periods[len(periods)-1]
		}
		return overduePeriodKey
	}

	return periods, assign
}

// labelsForTask returns the group labels for a task based on groupBy mode.
func labelsForTask(tk repository.FilteredTask, tagsByTask map[string]model.TaskTags, groupBy string) []string {
	if groupBy == "status" {
		return []string{statusLabel(tk.Status)}
	}
	tags := tagsByTask[tk.ID]
	if len(tags) == 0 {
		return []string{"Untagged"}
	}
	labels := make([]string, len(tags))
	for i, tag := range tags {
		labels[i] = tag.Name
	}
	return labels
}

// addToGroup accumulates task metrics into a TrendGroup.
func addToGroup(g *model.TrendGroup, tk repository.FilteredTask, tgt repository.TaskTargets, act repository.TaskActuals) {
	g.TargetTasks++
	if tk.Status == int(model.TaskStatusCompleted) {
		g.ActualTasks++
	}
	g.TargetDuration += tgt.TargetDuration
	g.ActualDuration += act.ActualDuration
	g.TargetEpisodes += tgt.TargetEpisodes
	g.ActualEpisodes += act.ActualEpisodes
}

// groupsMapToSlice converts a label→group map to a sorted slice.
func groupsMapToSlice(m map[string]*model.TrendGroup) []model.TrendGroup {
	groups := make([]model.TrendGroup, 0, len(m))
	for _, g := range m {
		groups = append(groups, *g)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Label < groups[j].Label
	})
	return groups
}

func buildCompletionTrend(
	tasks []repository.FilteredTask,
	targets map[string]repository.TaskTargets,
	actuals map[string]repository.TaskActuals,
	tagsByTask map[string]model.TaskTags,
	groupBy string,
	from, to time.Time,
	interval string,
) model.TaskCompletionTrend {
	periods, assignPeriod := buildTrendPeriods(from, to, interval)

	// Initialize groups for all periods (including overdue)
	periodGroups := make(map[periodKey]map[string]*model.TrendGroup)
	for _, pk := range periods {
		periodGroups[pk] = make(map[string]*model.TrendGroup)
	}
	periodGroups[overduePeriodKey] = make(map[string]*model.TrendGroup)

	for _, tk := range tasks {
		tgt := targets[tk.ID]
		act := actuals[tk.ID]
		labels := labelsForTask(tk, tagsByTask, groupBy)

		for _, label := range labels {
			pk := assignPeriod(tk.Deadline)
			if periodGroups[pk][label] == nil {
				periodGroups[pk][label] = &model.TrendGroup{Label: label}
			}
			addToGroup(periodGroups[pk][label], tk, tgt, act)
		}
	}

	// Build result: overdue → future periods
	result := model.TaskCompletionTrend{
		Periods: make([]model.TrendPeriod, 0, len(periods)+1),
	}

	// Overdue (deadline < today)
	if len(periodGroups[overduePeriodKey]) > 0 {
		result.Periods = append(result.Periods, model.TrendPeriod{
			Groups: groupsMapToSlice(periodGroups[overduePeriodKey]),
		})
	}

	// Future periods (today → 2 months ahead)
	for _, pk := range periods {
		result.Periods = append(result.Periods, model.TrendPeriod{
			Start:  pk.start,
			End:    pk.end,
			Groups: groupsMapToSlice(periodGroups[pk]),
		})
	}

	return result
}

func statusLabel(status int) string {
	switch model.TaskStatus(status) {
	case model.TaskStatusPlanning:
		return "Planning"
	case model.TaskStatusDoing:
		return "Doing"
	case model.TaskStatusCompleted:
		return "Completed"
	case model.TaskStatusCanceled:
		return "Canceled"
	default:
		return "Unknown"
	}
}
