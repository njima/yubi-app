"use client";

/**
 * Task Update Mutation Hook
 * TanStack Query mutation for updating tasks
 */

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { type Task } from "../schemas";
import { taskCompletionTrendQueryKeys } from "./use-task-completion-trend-query";
import { taskSummaryQueryKeys } from "./use-task-summary-query";
import { tasksQueryKeys } from "./use-tasks-query";

type TaskUpdateInput = z.infer<typeof schemas.TaskUpdate>;

/**
 * Hook to update an existing task
 *
 * @example
 * ```tsx
 * const { mutate, isPending } = useUpdateTaskMutation();
 *
 * mutate(
 *   { taskId: "task-123", data: { name: "Updated task" } },
 *   {
 *     onSuccess: () => console.log("Task updated!"),
 *     onError: (error) => console.error(error),
 *   }
 * );
 * ```
 */
export function useUpdateTaskMutation() {
  const queryClient = useQueryClient();

  return useMutation<Task, Error, { taskId: string; data: TaskUpdateInput }>({
    mutationFn: async ({ taskId, data }) => {
      const payload = schemas.TaskUpdate.parse(data);
      const response = await fetch(`/web/api/tasks/${taskId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!response.ok) {
        throw new Error(`Failed to update task: ${response.statusText}`);
      }
      const updatedTask = await response.json();
      return schemas.Task.parse(updatedTask);
    },
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({
        queryKey: tasksQueryKeys.lists(),
      });
      queryClient.invalidateQueries({
        queryKey: tasksQueryKeys.detail(variables.taskId),
      });
      queryClient.invalidateQueries({
        queryKey: taskSummaryQueryKeys.all,
      });
      queryClient.invalidateQueries({
        queryKey: taskCompletionTrendQueryKeys.all,
      });
      toast.success("Task updated successfully");
    },
    onError: (error) => {
      toast.error("Failed to update task", {
        description: error.message || "An unexpected error occurred",
      });
    },
  });
}
