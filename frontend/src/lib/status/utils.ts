/**
 * Status Utility Functions
 *
 * Provides helper functions for parsing and formatting status values
 */

import {
  EPISODE_COLLECTION_STATUS,
  type EpisodeCollectionStatusValue,
  ROBOT_STATUS,
  type RobotStatusValue,
  SUBTASK_STATUS,
  type SubTaskStatusValue,
  TASK_DIFFICULTY,
  type TaskDifficultyValue,
  TASK_PRIORITY,
  type TaskPriorityValue,
  TASK_STATUS,
  type TaskStatusValue,
  USER_ROLE,
  type UserRoleValue,
} from "./constants";

/**
 * Parse string to EpisodeCollectionStatusValue
 * Returns undefined if the value is invalid
 */
export function parseEpisodeCollectionStatus(
  value: string | null | undefined
): EpisodeCollectionStatusValue | undefined {
  if (!value) return undefined;
  const parsed = parseInt(value, 10);
  if (Number.isNaN(parsed)) return undefined;

  return parsed === EPISODE_COLLECTION_STATUS.READY ||
    parsed === EPISODE_COLLECTION_STATUS.RECORDING ||
    parsed === EPISODE_COLLECTION_STATUS.CANCEL ||
    parsed === EPISODE_COLLECTION_STATUS.COMPLETED
    ? (parsed as EpisodeCollectionStatusValue)
    : undefined;
}

/**
 * Parse string to RobotStatusValue
 * Returns undefined if the value is invalid
 */
export function parseRobotStatus(
  value: string | null | undefined
): RobotStatusValue | undefined {
  if (!value) return undefined;
  const parsed = parseInt(value, 10);
  if (Number.isNaN(parsed)) return undefined;

  return parsed === ROBOT_STATUS.ONLINE ||
    parsed === ROBOT_STATUS.BUSY ||
    parsed === ROBOT_STATUS.OFFLINE ||
    parsed === ROBOT_STATUS.FAULTED ||
    parsed === ROBOT_STATUS.MAINTENANCE ||
    parsed === ROBOT_STATUS.READY
    ? (parsed as RobotStatusValue)
    : undefined;
}

/**
 * Parse string to UserRoleValue
 * Returns undefined if the value is invalid
 */
export function parseUserRole(
  value: string | null | undefined
): UserRoleValue | undefined {
  if (!value) return undefined;
  const parsed = parseInt(value, 10);
  if (Number.isNaN(parsed)) return undefined;

  return parsed === USER_ROLE.VIEWER ||
    parsed === USER_ROLE.DATA_ENGINEER ||
    parsed === USER_ROLE.OPERATOR ||
    parsed === USER_ROLE.ADMIN
    ? (parsed as UserRoleValue)
    : undefined;
}

/**
 * Parse string to TaskStatusValue
 * Returns undefined if the value is invalid
 */
export function parseTaskStatus(
  value: string | null | undefined
): TaskStatusValue | undefined {
  if (!value) return undefined;
  const parsed = parseInt(value, 10);
  if (Number.isNaN(parsed)) return undefined;

  return parsed === TASK_STATUS.PLANNING ||
    parsed === TASK_STATUS.DOING ||
    parsed === TASK_STATUS.COMPLETED ||
    parsed === TASK_STATUS.CANCELED
    ? (parsed as TaskStatusValue)
    : undefined;
}

/**
 * Parse string to TaskPriorityValue
 * Returns undefined if the value is invalid
 */
export function parseTaskPriority(
  value: string | null | undefined
): TaskPriorityValue | undefined {
  if (!value) return undefined;
  const parsed = parseInt(value, 10);
  if (Number.isNaN(parsed)) return undefined;

  return parsed === TASK_PRIORITY.LOW ||
    parsed === TASK_PRIORITY.NORMAL ||
    parsed === TASK_PRIORITY.HIGH ||
    parsed === TASK_PRIORITY.URGENT
    ? (parsed as TaskPriorityValue)
    : undefined;
}

/**
 * Parse string to TaskDifficultyValue
 * Returns undefined if the value is invalid
 */
export function parseTaskDifficulty(
  value: string | null | undefined
): TaskDifficultyValue | undefined {
  if (!value) return undefined;
  const parsed = parseInt(value, 10);
  if (Number.isNaN(parsed)) return undefined;

  return parsed === TASK_DIFFICULTY.S ||
    parsed === TASK_DIFFICULTY.A ||
    parsed === TASK_DIFFICULTY.B ||
    parsed === TASK_DIFFICULTY.C
    ? (parsed as TaskDifficultyValue)
    : undefined;
}

/**
 * Parse string to SubTaskStatusValue
 * Returns undefined if the value is invalid
 */
export function parseSubTaskStatus(
  value: string | null | undefined
): SubTaskStatusValue | undefined {
  if (!value) return undefined;
  const parsed = parseInt(value, 10);
  if (Number.isNaN(parsed)) return undefined;

  return parsed === SUBTASK_STATUS.READY ||
    parsed === SUBTASK_STATUS.IN_PROGRESS ||
    parsed === SUBTASK_STATUS.FAILED ||
    parsed === SUBTASK_STATUS.COMPLETED
    ? (parsed as SubTaskStatusValue)
    : undefined;
}

/**
 * Get display label for task difficulty
 */
export function getTaskDifficultyLabel(
  difficulty: TaskDifficultyValue
): string {
  switch (difficulty) {
    case TASK_DIFFICULTY.S:
      return "S";
    case TASK_DIFFICULTY.A:
      return "A";
    case TASK_DIFFICULTY.B:
      return "B";
    case TASK_DIFFICULTY.C:
      return "C";
    default:
      return "Unknown";
  }
}
