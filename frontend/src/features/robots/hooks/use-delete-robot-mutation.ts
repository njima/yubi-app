"use client";

/**
 * Robot Deletion Mutation Hook
 * TanStack Query mutation for deleting robots
 */
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { fleetSummaryQueryKeys } from "./use-fleet-summary-query";
import { robotTypesQueryKeys } from "./use-robot-types-query";
import { robotsQueryKeys } from "./use-robots-query";

/**
 * Hook to delete a robot
 *
 * @example
 * ```tsx
 * const { mutate, isPending } = useDeleteRobotMutation();
 *
 * mutate(
 *   { robotId: 123 },
 *   {
 *     onSuccess: () => console.log("Robot deleted!"),
 *     onError: (error) => console.error(error),
 *   }
 * );
 * ```
 */
export function useDeleteRobotMutation() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, { robotId: string }>({
    mutationFn: async ({ robotId }) => {
      const response = await fetch(`/web/api/robots/${robotId}`, {
        method: "DELETE",
      });
      if (!response.ok) {
        throw new Error(`Failed to delete robot: ${response.statusText}`);
      }
    },
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({
        queryKey: robotsQueryKeys.lists(),
      });
      queryClient.invalidateQueries({
        queryKey: robotsQueryKeys.detail(variables.robotId),
      });
      queryClient.invalidateQueries({
        queryKey: fleetSummaryQueryKeys.all,
      });
      queryClient.invalidateQueries({
        queryKey: robotTypesQueryKeys.all,
      });
      toast.success("Robot deleted successfully");
    },
    onError: (error) => {
      toast.error("Failed to delete robot", {
        description: error.message || "An unexpected error occurred",
      });
    },
  });
}
