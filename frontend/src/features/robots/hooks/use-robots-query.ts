"use client";

import { useQuery, type UseQueryOptions } from "@tanstack/react-query";
import { z } from "zod";

import { fetchJson } from "@/lib/api/client-fetch";
import { schemas } from "@/lib/api/generated/api";
import { withQueryString } from "@/lib/api/query-string";

export const robotsQueryKeys = {
  all: ["robots"] as const,
  lists: () => [...robotsQueryKeys.all, "list"] as const,
  list: (
    params: {
      site_id?: string;
      location_id?: string;
      status?: 0 | 1 | 2 | 3 | 4;
      robot_type?: string;
      page?: number;
      limit?: number;
      search?: string;
      sort_by?: string;
      sort_order?: string;
    } = {}
  ) => [...robotsQueryKeys.lists(), params] as const,
  details: () => [...robotsQueryKeys.all, "detail"] as const,
  detail: (id: string) => [...robotsQueryKeys.details(), id] as const,
};

type Robot = z.infer<typeof schemas.Robot>;
type RobotListResponse = z.infer<typeof schemas.RobotListResponse>;

function normalizeRobot(robot: unknown): unknown {
  if (typeof robot !== "object" || robot === null) return robot;
  return {
    ...(robot as Record<string, unknown>),
    robot_config:
      (robot as Record<string, unknown>).robot_config === null
        ? undefined
        : (robot as Record<string, unknown>).robot_config,
  };
}

export function useRobotsQuery(
  params: {
    site_id?: string;
    location_id?: string;
    status?: 0 | 1 | 2 | 3 | 4;
    robot_type?: string;
    page?: number;
    limit?: number;
    search?: string;
    sort_by?: string;
    sort_order?: string;
  } = {},
  options?: Omit<
    UseQueryOptions<RobotListResponse, Error>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery({
    queryKey: robotsQueryKeys.list(params),
    queryFn: async () => {
      const data = await fetchJson<RobotListResponse>(
        withQueryString("/web/api/robots", params),
        "Failed to fetch robots"
      );
      const normalized = {
        ...data,
        robots: (data.robots || []).map(normalizeRobot),
      };
      return schemas.RobotListResponse.parse(normalized);
    },
    ...options,
  });
}

export function useRobotQuery(
  robotId: string,
  options?: Omit<UseQueryOptions<Robot, Error>, "queryKey" | "queryFn">
) {
  return useQuery({
    queryKey: robotsQueryKeys.detail(robotId),
    queryFn: async () => {
      const data = await fetchJson<Robot>(
        `/web/api/robots/${robotId}`,
        "Failed to fetch robot"
      );
      const normalized = normalizeRobot(data);
      return schemas.Robot.parse(normalized);
    },
    enabled: !!robotId,
    staleTime: 0,
    refetchOnMount: "always",
    ...options,
  });
}
