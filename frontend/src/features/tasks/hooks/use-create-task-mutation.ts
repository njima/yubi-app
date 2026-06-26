"use client";

/**
 * Task Creation Mutation Hook
 * TanStack Query mutation for creating tasks
 */

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { type Task } from "../schemas";
import { taskCompletionTrendQueryKeys } from "./use-task-completion-trend-query";
import { taskSummaryQueryKeys } from "./use-task-summary-query";
import { tasksQueryKeys } from "./use-tasks-query";

type TaskCreateInput = z.infer<typeof schemas.TaskCreate>;

/**
 * Hook to create a new task
 *
 * @example
 * ```tsx
 * const { mutate, isPending } = useCreateTaskMutation();
 *
 * mutate(
 *   { name: "New task", description: "Optional description" },
 *   {
 *     onSuccess: () => console.log("Task created!"),
 *     onError: (error) => console.error(error),
 *   }
 * );
 * ```
 */
export function useCreateTaskMutation() {
  const queryClient = useQueryClient();

  return useMutation<Task, Error, TaskCreateInput>({
    mutationFn: async (data: TaskCreateInput) => {
      const payload = schemas.TaskCreate.parse(data);
      const response = await fetch("/web/api/tasks", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!response.ok) {
        throw new Error(`Failed to create task: ${response.statusText}`);
      }
      const createdTask = await response.json();
      return schemas.Task.parse(createdTask);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: tasksQueryKeys.lists(),
      });
      queryClient.invalidateQueries({
        queryKey: taskSummaryQueryKeys.all,
      });
      queryClient.invalidateQueries({
        queryKey: taskCompletionTrendQueryKeys.all,
      });
      toast.success("Task created successfully");
    },
    onError: (error) => {
      toast.error("Failed to create task", {
        description: error.message || "An unexpected error occurred",
      });
    },
  });
}
