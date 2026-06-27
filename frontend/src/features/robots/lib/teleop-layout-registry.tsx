"use client";

import type { schemas } from "@/lib/api/generated/api";

import type { LayoutItem } from "./teleop-layout-types";
import type { ReactNode } from "react";
import type { z } from "zod";

type Episode = z.infer<typeof schemas.Episode>;
type RobotStatusStreamDetail = z.infer<typeof schemas.RobotStatusStreamDetail>;

export interface LayoutCamera {
  namespace: string;
  name?: string;
}

/**
 * Context bag passed to every layout component.
 * Each component reads only the fields it needs; all fields are optional.
 */
export interface LayoutContext {
  // Robot
  robot?: {
    id: string;
    name: string;
    robot_type?: string;
    robot_config?: Record<string, unknown> | null;
  };
  realtimeStatus?: RobotStatusStreamDetail | null;
  isRobotConnected?: boolean;
  isLoadingRobot?: boolean;

  // Episode
  episode?: Episode | null;
  isLoadingEpisode?: boolean;
  activeEpisodeId?: string;
  currentSubtask?: { order_index: number; name: string } | null;
  nextSubtask?: { order_index: number; name: string } | null;

  // Task
  taskName?: string;
  taskVersion?: string;
  taskManualUrl?: string;
  taskDescription?: string;

  // Camera
  cameras?: LayoutCamera[];
  host?: string;
  port?: number;
  rosbridgePort?: number;
  streamConfig?: { quality?: number; width?: number; height?: number };

  // Gate
  gateLevel?: number;
}

/**
 * Registry of renderable components keyed by string ID.
 * Lazy-loaded to avoid circular imports — populated by feature modules.
 */
type LayoutComponentRenderer = (
  ctx: LayoutContext,
  item: LayoutItem
) => ReactNode;

const registry = new Map<string, LayoutComponentRenderer>();

export function registerLayoutComponent(
  id: string,
  render: LayoutComponentRenderer
) {
  registry.set(id, render);
}

export function getLayoutComponent(
  id: string
): LayoutComponentRenderer | undefined {
  return registry.get(id);
}
