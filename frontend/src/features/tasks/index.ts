// Task Components
export { CreateTaskDialog } from "./components/create-task-dialog";
export { ImportTasksDialog } from "./components/import-tasks-dialog";
export { ExportTasksDialog } from "./components/export-tasks-dialog";
export { CreateTaskForm } from "./components/create-task-form";
export { EditTaskDialog } from "./components/edit-task-dialog";
export { EditTaskForm } from "./components/edit-task-form";
export { TagCategoryBadge } from "./components/tag-category-badge";
export { TaskDifficultyBadge } from "./components/task-difficulty-badge";
export { TaskPriorityBadge } from "./components/task-priority-badge";
export { TaskStatusBadge } from "./components/task-status-badge";
export { TeachMeBizCard } from "./components/teach-me-biz-card";
export { VersionHistoryDialog } from "./components/version-history-dialog";

// Task Detail Components
export { TaskDetailPage } from "./components/detail";

// SubTask Components
export { CreateSubTaskDialog } from "./components/create-subtask-dialog";
export { CreateSubTaskForm } from "./components/create-subtask-form";
export { DeleteSubTaskDialog } from "./components/delete-subtask-dialog";
export { EditSubTaskDialog } from "./components/edit-subtask-dialog";
export { EditSubTaskForm } from "./components/edit-subtask-form";
export { SubTaskList } from "./components/subtask-list";

// Task Table Components
export { TaskDataTable } from "./components/task-data-table";
export { getTaskColumns } from "./components/task-columns";

// Task Summary Components
export { TaskSummaryCard } from "./components/task-summary-card";
export { TaskCompletionTrend } from "./components/task-completion-trend";
export { ParameterizedName } from "./components/parameterized-name";

// Task Hooks
export { useTasksQuery, useTaskQuery } from "./hooks/use-tasks-query";
export { useTaskSearchOptions } from "./hooks/use-task-search-options";
export { useTaskAvailableTagsQuery } from "./hooks/use-task-available-tags-query";
export {
  useTaskVersionsQuery,
  taskVersionsQueryKeys,
} from "./hooks/use-task-versions-query";
export { useCreateTaskMutation } from "./hooks/use-create-task-mutation";
export { useUpdateTaskMutation } from "./hooks/use-update-task-mutation";
export {
  useValidateTaskImportMutation,
  useImportTasksMutation,
} from "./hooks/use-import-tasks-mutation";

export {
  useTaskSummaryQuery,
  taskSummaryQueryKeys,
} from "./hooks/use-task-summary-query";
export {
  useTaskCompletionTrendQuery,
  taskCompletionTrendQueryKeys,
} from "./hooks/use-task-completion-trend-query";

// SubTask Hooks
export {
  useSubTasksQuery,
  subtasksQueryKeys,
} from "./hooks/use-subtasks-query";
export {
  useSubTasksByVersionQuery,
  subtasksByVersionQueryKeys,
} from "./hooks/use-subtasks-by-version-query";
export { useCreateSubTaskMutation } from "./hooks/use-create-subtask-mutation";
export { useUpdateSubTaskMutation } from "./hooks/use-update-subtask-mutation";
export { useDeleteSubTaskMutation } from "./hooks/use-delete-subtask-mutation";

// Schemas - Re-exported from OpenAPI
export {
  taskSchema,
  taskCreateSchema,
  taskUpdateSchema,
  taskVersionSchema,
  type Task,
  type TaskCreate,
  type TaskUpdate,
  type TaskVersion,
} from "./schemas";
