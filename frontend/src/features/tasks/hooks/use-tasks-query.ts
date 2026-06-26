"use client";

import { useQuery, type UseQueryOptions } from "@tanstack/react-query";
import { z } from "zod";

import { fetchAndParse } from "@/lib/api/client-fetch";
import { schemas } from "@/lib/api/generated/api";
import { withQueryString } from "@/lib/api/query-string";

import { taskSchema, type Task } from "../schemas";

type TaskListResponse = z.infer<typeof schemas.TaskListResponse>;

export const tasksQueryKeys = {
  all: ["tasks"] as const,
  lists: () => [...tasksQueryKeys.all, "list"] as const,
  list: (params?: {
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
  }) => [...tasksQueryKeys.lists(), params] as const,
  details: () => [...tasksQueryKeys.all, "detail"] as const,
  detail: (id: string) => [...tasksQueryKeys.details(), id] as const,
};

export function useTasksQuery(
  params?: {
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
  },
  options?: Omit<
    UseQueryOptions<TaskListResponse, Error>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery({
    queryKey: tasksQueryKeys.list(params),
    queryFn: async () => {
      return fetchAndParse(
        withQueryString("/web/api/tasks", {
          ...params,
          has_approved_version: params?.has_approved_version || undefined,
        }),
        schemas.TaskListResponse,
        "Failed to fetch tasks"
      );
    },
    ...options,
  });
}

type GetTaskResponse = Task;

export function useTaskQuery(
  taskId: string,
  options?: Omit<
    UseQueryOptions<GetTaskResponse, Error>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery({
    queryKey: tasksQueryKeys.detail(taskId),
    queryFn: async () => {
      return fetchAndParse(
        `/web/api/tasks/${taskId}`,
        taskSchema,
        "Failed to fetch task"
      );
    },
    enabled: !!taskId,
    ...options,
  });
}
