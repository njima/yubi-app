import { fetchBackend } from "./core";

import type {
  SubTask,
  SubTaskCreate,
  SubTaskReorder,
  SubTaskUpdate,
} from "./types";

// =============================================================================
// SubTasks API
// =============================================================================

export async function fetchSubTasks(
  taskId?: string,
  taskVersionId?: string
): Promise<SubTask[]> {
  const params = new URLSearchParams();
  if (taskId) params.set("task_id", taskId);
  if (taskVersionId) params.set("task_version_id", taskVersionId);
  const query = params.toString() ? `?${params.toString()}` : "";
  return fetchBackend<SubTask[]>(`/api/sub-tasks${query}`);
}

export async function fetchSubTask(subtaskId: string): Promise<SubTask> {
  return fetchBackend<SubTask>(`/api/sub-tasks/${subtaskId}`);
}

export async function createSubTask(data: SubTaskCreate): Promise<SubTask> {
  return fetchBackend<SubTask>("/api/sub-tasks", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function updateSubTask(
  subtaskId: string,
  data: SubTaskUpdate
): Promise<SubTask> {
  return fetchBackend<SubTask>(`/api/sub-tasks/${subtaskId}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function deleteSubTask(subtaskId: string): Promise<void> {
  await fetchBackend<void>(`/api/sub-tasks/${subtaskId}`, {
    method: "DELETE",
  });
}

export async function reorderSubTasks(
  data: SubTaskReorder
): Promise<SubTask[]> {
  return fetchBackend<SubTask[]>("/api/sub-tasks/reorder", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function completeSubTask(subtaskId: string): Promise<SubTask> {
  return fetchBackend<SubTask>(`/api/sub-tasks/${subtaskId}/complete`, {
    method: "POST",
  });
}
