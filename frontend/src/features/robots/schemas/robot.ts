/**
 * Robot Schemas
 * Re-exports OpenAPI schemas and provides type helpers
 */

import { schemas } from "@/lib/api/generated/api";

import type { z } from "zod";

// Re-export Robot schemas from OpenAPI
export const robotSchema = schemas.Robot;
export const robotCreateSchema = schemas.RobotCreate;
export const robotUpdateSchema = schemas.RobotUpdate;

// Infer types
export type Robot = z.infer<typeof robotSchema>;
export type RobotCreate = z.infer<typeof robotCreateSchema>;
export type RobotUpdate = z.infer<typeof robotUpdateSchema>;

// SSE stream types
export type RobotStatusStreamResponse = z.infer<
  typeof schemas.RobotStatusStreamResponse
>;
export type RobotStatusStreamDetail = z.infer<
  typeof schemas.RobotStatusStreamDetail
>;

// Type helpers
export type RobotConfig = Robot["robot_config"];

/**
 * Robot Status Constants
 * Re-exported from lib/status/constants.ts to avoid duplication
 */
export {
  ROBOT_STATUS,
  type RobotStatusValue as RobotStatus,
} from "@/lib/status/constants";

/**
 * Battery level thresholds
 */
export const BATTERY_THRESHOLDS = {
  LOW: 20,
  CRITICAL: 10,
} as const;

/**
 * Heartbeat timeout threshold (in minutes)
 */
export const HEARTBEAT_TIMEOUT_MINUTES = 15;
