"use client";

/**
 * Robot Update Mutation Hook
 * TanStack Query mutation for updating robots
 */

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { fleetSummaryQueryKeys } from "./use-fleet-summary-query";
import { robotTypesQueryKeys } from "./use-robot-types-query";
import { robotsQueryKeys } from "./use-robots-query";

type RobotUpdate = z.infer<typeof schemas.RobotUpdate>;
type Robot = z.infer<typeof schemas.Robot>;

/**
 * Hook to update an existing robot
 *
 * @example
 * ```tsx
 * const { mutate, isPending } = useUpdateRobotMutation();
 *
 * mutate(
 *   { robotId: 1, data: { name: "Robot Alpha v2" } },
 *   {
 *     onSuccess: () => console.log("Robot updated!"),
 *     onError: (error) => console.error(error),
 *   }
 * );
 * ```
 */
export function useUpdateRobotMutation() {
  const queryClient = useQueryClient();

  return useMutation<Robot, Error, { robotId: string; data: RobotUpdate }>({
    mutationFn: async ({ robotId, data }) => {
      const response = await fetch(`/web/api/robots/${robotId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      if (!response.ok) {
        throw new Error(`Failed to update robot: ${response.statusText}`);
      }
      return response.json();
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
      toast.success("Robot updated successfully");
    },
    onError: (error) => {
      toast.error("Failed to update robot", {
        description: error.message || "An unexpected error occurred",
      });
    },
  });
}
