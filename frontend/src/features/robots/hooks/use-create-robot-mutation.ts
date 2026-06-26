"use client";

/**
 * Robot Creation Mutation Hook
 * TanStack Query mutation for creating robots
 */

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { z } from "zod";

import { schemas } from "@/lib/api/generated/api";

import { fleetSummaryQueryKeys } from "./use-fleet-summary-query";
import { robotTypesQueryKeys } from "./use-robot-types-query";
import { robotsQueryKeys } from "./use-robots-query";

type RobotCreate = z.infer<typeof schemas.RobotCreate>;
type Robot = z.infer<typeof schemas.Robot>;

export function useCreateRobotMutation() {
  const queryClient = useQueryClient();

  return useMutation<Robot, Error, RobotCreate>({
    mutationFn: async (data: RobotCreate) => {
      const payload = schemas.RobotCreate.parse(data);
      const response = await fetch("/web/api/robots", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!response.ok) {
        throw new Error(`Failed to create robot: ${response.statusText}`);
      }
      const responseData = await response.json();
      const normalized = {
        ...responseData,
        robot_config:
          responseData?.robot_config === null
            ? undefined
            : responseData?.robot_config,
      };
      return schemas.Robot.parse(normalized);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: robotsQueryKeys.lists(),
      });
      queryClient.invalidateQueries({
        queryKey: fleetSummaryQueryKeys.all,
      });
      queryClient.invalidateQueries({
        queryKey: robotTypesQueryKeys.all,
      });
      toast.success("Robot created successfully");
    },
    onError: (error) => {
      toast.error("Failed to create robot", {
        description: error.message || "An unexpected error occurred",
      });
    },
  });
}
