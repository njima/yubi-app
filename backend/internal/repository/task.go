package repository

import (
	"context"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

// MaxTaskBatchSize is the maximum number of tasks allowed in a single import or export operation.
const MaxTaskBatchSize = 5000

type TaskListFilter struct {
	HasApprovedVersion *bool
	SortBy             *TaskSortBy
	SortOrder          *SortOrder
	Statuses           []TaskStatus
	Priorities         []TaskPriority
	Difficulties       []TaskDifficulty
	RobotType          *string
	Search             *string
}

type TaskSummaryFilter struct {
	RobotTypes     []string
	CategoryTypeID *string
	TagIDs         []string
	DeadlineFrom   *time.Time
	DeadlineTo     *time.Time
}

// FilteredTask holds minimal task info for summary/trend queries.
type FilteredTask struct {
	ID       string    `bun:"id_natural"`
	Deadline time.Time `bun:"deadline"`
	Status   int       `bun:"status"`
}

// TaskTargets holds aggregated target values per task.
type TaskTargets struct {
	TaskID         string `bun:"task_id"`
	TargetDuration int64  `bun:"target_duration"`
	TargetEpisodes int    `bun:"target_episodes"`
}

// TaskActuals holds aggregated actual values per task.
type TaskActuals struct {
	TaskID         string `bun:"task_id"`
	ActualDuration int64  `bun:"actual_duration"`
	ActualEpisodes int    `bun:"actual_episodes"`
}

type Task interface {
	Create(ctx context.Context, conn DBConn, t model.Task) (model.Task, error)
	Exists(ctx context.Context, conn DBConn, id string) (bool, error)
	GetByID(ctx context.Context, conn DBConn, id string) (model.Task, error)
	List(ctx context.Context, conn DBConn, filter TaskListFilter, limit, offset int) (model.Tasks, int, error)
	Update(ctx context.Context, conn DBConn, t model.Task) (model.Task, error)
	Delete(ctx context.Context, conn DBConn, id string) error
	ListByIDs(ctx context.Context, conn DBConn, ids []string) (model.Tasks, error)
	GetFilteredTasks(ctx context.Context, conn DBConn, filter TaskSummaryFilter) ([]FilteredTask, error)
	GetTargetsByTaskIDs(ctx context.Context, conn DBConn, taskIDs []string) (map[string]TaskTargets, error)
	GetActualsByTaskIDs(ctx context.Context, conn DBConn, taskIDs []string) (map[string]TaskActuals, error)
	FindExistingNames(ctx context.Context, conn DBConn, names []string) (map[string]bool, error)
	BulkCreate(ctx context.Context, conn DBConn, items []BulkTaskItem) ([]model.Task, error)
	Export(ctx context.Context, conn DBConn, filter TaskListFilter) ([]TaskExportRow, error)
}

// TaskExportRow holds the data for a single row in the task export CSV.
// Version-specific fields (subtasks, targets) are taken from the latest approved version.
type TaskExportRow struct {
	IDNatural                       string
	Name                            string
	Description                     *string
	ManualURL                       string
	Priority                        model.TaskPriority
	Difficulty                      model.TaskDifficulty
	Status                          model.TaskStatus
	Deadline                        time.Time
	RobotType                       *string
	SubtaskNames                    []string // from latest approved version, sorted by order_index
	TargetDurationSeconds           *int     // from latest approved version
	TargetEpisodeCount              *int     // from latest approved version
	TargetDurationPerEpisodeSeconds *int     // from latest approved version
}

// BulkTaskItem is a single task to be imported via BulkCreate,
// containing optional subtask names and task_version target fields.
type BulkTaskItem struct {
	Task                            model.Task
	SubtaskNames                    []string
	TargetDurationSeconds           *int
	TargetEpisodeCount              *int
	TargetDurationPerEpisodeSeconds *int
}
