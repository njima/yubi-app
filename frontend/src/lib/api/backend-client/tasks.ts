import { z } from "zod";

import { fetchBackend } from "./core";
import { schemas } from "../generated/api";

import type {
  Task,
  TaskCategoryType,
  TaskCreate,
  TaskListResponse,
  TaskTag,
  TaskTagCreate,
  TaskUpdate,
} from "./types";

// =============================================================================
// Tasks API
// =============================================================================

export async function fetchTasks(params?: {
  has_approved_version?: boolean;
  page?: number;
  limit?: number;
  sort_by?: string;
  sort_order?: string;
  status?: number[];
  priority?: number[];
  difficulty?: number[];
  robot_type?: string;
  search?: string;
}): Promise<TaskListResponse> {
  const searchParams = new URLSearchParams();
  if (params?.has_approved_version) {
    searchParams.append("has_approved_version", "true");
  }
  if (params?.page !== undefined) {
    searchParams.append("page", String(params.page));
  }
  if (params?.limit !== undefined) {
    searchParams.append("limit", String(params.limit));
  }
  if (params?.sort_by) searchParams.append("sort_by", params.sort_by);
  if (params?.sort_order) searchParams.append("sort_order", params.sort_order);
  params?.status?.forEach((s) => searchParams.append("status", String(s)));
  params?.priority?.forEach((p) => searchParams.append("priority", String(p)));
  params?.difficulty?.forEach((d) =>
    searchParams.append("difficulty", String(d))
  );
  if (params?.robot_type) searchParams.append("robot_type", params.robot_type);
  if (params?.search) searchParams.append("search", params.search);
  const query = searchParams.toString();
  return fetchBackend<TaskListResponse>(
    `/api/tasks${query ? `?${query}` : ""}`
  );
}

export async function fetchTask(taskId: string): Promise<Task> {
  return fetchBackend<Task>(`/api/tasks/${taskId}`);
}

export type TaskSummaryResponse = z.infer<typeof schemas.TaskSummaryResponse>;
export type TaskCompletionTrend = z.infer<typeof schemas.TaskCompletionTrend>;

export async function fetchTaskSummary(params?: {
  robot_type?: string[];
  category_type_id?: string;
  tag_id?: string[];
  from?: string;
  to?: string;
}): Promise<TaskSummaryResponse> {
  const searchParams = new URLSearchParams();
  params?.robot_type?.forEach((m) => searchParams.append("robot_type", m));
  if (params?.category_type_id)
    searchParams.append("category_type_id", params.category_type_id);
  params?.tag_id?.forEach((id) => searchParams.append("tag_id", id));
  if (params?.from) searchParams.append("from", params.from);
  if (params?.to) searchParams.append("to", params.to);
  const query = searchParams.toString();
  return fetchBackend<TaskSummaryResponse>(
    `/api/tasks/summary${query ? `?${query}` : ""}`
  );
}

export async function fetchTaskCompletionTrend(params: {
  group_by: string;
  robot_type?: string[];
  category_type_id?: string;
  tag_id?: string[];
  from?: string;
  to?: string;
  interval?: string;
}): Promise<TaskCompletionTrend> {
  const searchParams = new URLSearchParams();
  searchParams.append("group_by", params.group_by);
  params.robot_type?.forEach((m) => searchParams.append("robot_type", m));
  if (params.category_type_id)
    searchParams.append("category_type_id", params.category_type_id);
  params.tag_id?.forEach((id) => searchParams.append("tag_id", id));
  if (params.from) searchParams.append("from", params.from);
  if (params.to) searchParams.append("to", params.to);
  if (params.interval) searchParams.append("interval", params.interval);
  const query = searchParams.toString();
  return fetchBackend<TaskCompletionTrend>(
    `/api/tasks/completion-trend?${query}`
  );
}

export async function fetchTaskAvailableTags(params?: {
  robot_type?: string[];
  category_type_id?: string;
}): Promise<TaskTag[]> {
  const searchParams = new URLSearchParams();
  params?.robot_type?.forEach((m) => searchParams.append("robot_type", m));
  if (params?.category_type_id)
    searchParams.append("category_type_id", params.category_type_id);
  const query = searchParams.toString();
  return fetchBackend<TaskTag[]>(
    `/api/tasks/available-tags${query ? `?${query}` : ""}`
  );
}

export async function createTask(data: TaskCreate): Promise<Task> {
  return fetchBackend<Task>("/api/tasks", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function updateTask(
  taskId: string,
  data: TaskUpdate
): Promise<Task> {
  return fetchBackend<Task>(`/api/tasks/${taskId}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

// =============================================================================
// Task Import API
// =============================================================================

export interface TaskImportValidationResponse {
  valid_rows: TaskImportRow[];
  duplicate_rows: TaskImportRowError[];
  error_rows: TaskImportRowError[];
  summary: {
    valid_count: number;
    duplicate_count: number;
    error_count: number;
  };
}

export interface TaskImportRow {
  row_number: number;
  name: string;
  description?: string;
  manual_url: string;
  priority: string;
  difficulty: string;
  status?: string;
  deadline: string;
  robot_type?: string;
  tags?: string;
}

export interface TaskImportRowError {
  row_number: number;
  errors: string[];
  name?: string;
}

export interface TaskImportResponse {
  imported_count: number;
  skipped_count: number;
  error_count: number;
  errors: TaskImportRowError[];
}

export async function validateTaskImport(
  csvContent: string
): Promise<TaskImportValidationResponse> {
  return fetchBackend<TaskImportValidationResponse>(
    "/api/tasks/import/validate",
    {
      method: "POST",
      body: JSON.stringify({ csv_content: csvContent }),
    }
  );
}

export async function importTasks(
  csvContent: string
): Promise<TaskImportResponse> {
  return fetchBackend<TaskImportResponse>("/api/tasks/import", {
    method: "POST",
    body: JSON.stringify({ csv_content: csvContent }),
  });
}

// =============================================================================
// Task Versions API
// =============================================================================

export type TaskVersion = z.infer<typeof schemas.TaskVersion>;
export type TaskVersionCreate = z.infer<typeof schemas.TaskVersionCreate>;
export type TaskVersionUpdate = z.infer<typeof schemas.TaskVersionUpdate>;

export async function fetchTaskVersions(
  taskId: string
): Promise<TaskVersion[]> {
  return fetchBackend<TaskVersion[]>(`/api/tasks/${taskId}/versions`);
}

export async function createTaskVersion(
  taskId: string,
  data: TaskVersionCreate
): Promise<TaskVersion> {
  return fetchBackend<TaskVersion>(`/api/tasks/${taskId}/versions`, {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function updateTaskVersion(
  taskId: string,
  versionId: string,
  data: TaskVersionUpdate
): Promise<TaskVersion> {
  return fetchBackend<TaskVersion>(
    `/api/tasks/${taskId}/versions/${versionId}`,
    { method: "PATCH", body: JSON.stringify(data) }
  );
}

export async function approveTaskVersion(
  taskId: string,
  versionId: string
): Promise<TaskVersion> {
  return fetchBackend<TaskVersion>(
    `/api/tasks/${taskId}/versions/${versionId}/approve`,
    { method: "POST" }
  );
}

export type TaskVersionParametersUpdate = z.infer<
  typeof schemas.TaskVersionParametersUpdate
>;

export async function updateTaskVersionParameters(
  taskId: string,
  versionId: string,
  data: TaskVersionParametersUpdate
): Promise<TaskVersion> {
  return fetchBackend<TaskVersion>(
    `/api/tasks/${taskId}/versions/${versionId}/parameters`,
    {
      method: "PUT",
      body: JSON.stringify(data),
    }
  );
}

// =============================================================================
// Task Tags API
// =============================================================================

export async function fetchTaskTags(
  categoryTypeId?: string
): Promise<TaskTag[]> {
  const params = new URLSearchParams();
  if (categoryTypeId) params.set("category_type_id", categoryTypeId);
  const query = params.toString() ? `?${params.toString()}` : "";
  return fetchBackend<TaskTag[]>(`/api/task-tags${query}`);
}

export async function createTaskTag(data: TaskTagCreate): Promise<TaskTag> {
  return fetchBackend<TaskTag>("/api/task-tags", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

// =============================================================================
// Task Category Types API
// =============================================================================

export async function fetchTaskCategoryTypes(): Promise<TaskCategoryType[]> {
  return fetchBackend<TaskCategoryType[]>("/api/task-category-types");
}
