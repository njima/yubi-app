/**
 * Status Constants
 *
 * Defines status values derived from OpenAPI schema
 * Avoids magic numbers and improves code readability and maintainability
 */

/**
 * Episode Collection Status
 * Represents the data collection state of an episode
 */
export const EPISODE_COLLECTION_STATUS = {
  READY: 0,
  RECORDING: 1,
  CANCEL: 2,
  COMPLETED: 3,
} as const;

export type EpisodeCollectionStatusValue =
  (typeof EPISODE_COLLECTION_STATUS)[keyof typeof EPISODE_COLLECTION_STATUS];

/**
 * Robot Status
 * Represents the operational state of a robot
 */
export const ROBOT_STATUS = {
  ONLINE: 0,
  BUSY: 1,
  OFFLINE: 2,
  FAULTED: 3,
  MAINTENANCE: 4,
  READY: 5,
} as const;

export type RobotStatusValue = (typeof ROBOT_STATUS)[keyof typeof ROBOT_STATUS];

/**
 * Robot Type
 * Identifies the robot model.
 */
export const ROBOT_TYPE = {
  YUBI: "yubi",
  YUBI_PORTABLE: "yubi-portable",
} as const;

export type RobotTypeValue = (typeof ROBOT_TYPE)[keyof typeof ROBOT_TYPE];

/**
 * Leader Status
 * Represents the operational state of a robot's leader component
 */
export const LEADER_STATUS = {
  READY: 0,
  FAULTED: 1,
  MAINTENANCE: 2,
} as const;

export type LeaderStatusValue =
  (typeof LEADER_STATUS)[keyof typeof LEADER_STATUS];

/**
 * User Role
 * Represents the permission level of a user
 */
export const USER_ROLE = {
  ADMIN: 0,
  DATA_ENGINEER: 1,
  MANAGER: 2,
  OPERATOR: 3,
  VIEWER: 4,
} as const;

export type UserRoleValue = (typeof USER_ROLE)[keyof typeof USER_ROLE];

/**
 * SubTask Status
 * Represents the completion state of a subtask
 */
export const SUBTASK_STATUS = {
  READY: 0,
  IN_PROGRESS: 1,
  FAILED: 2,
  COMPLETED: 3,
} as const;

export type SubTaskStatusValue =
  (typeof SUBTASK_STATUS)[keyof typeof SUBTASK_STATUS];

/**
 * Approval Status
 * Represents the approval state of a task version
 */
export const APPROVAL_STATUS = {
  DRAFT: 0,
  APPROVED: 1,
} as const;

export type ApprovalStatusValue =
  (typeof APPROVAL_STATUS)[keyof typeof APPROVAL_STATUS];

/**
 * Task Priority
 * Represents the priority level of a task
 */
export const TASK_PRIORITY = {
  LOW: 0,
  NORMAL: 1,
  HIGH: 2,
  URGENT: 3,
} as const;

export type TaskPriorityValue =
  (typeof TASK_PRIORITY)[keyof typeof TASK_PRIORITY];

/**
 * Task Difficulty
 * Represents the difficulty level of a task
 */
export const TASK_DIFFICULTY = {
  S: 0,
  A: 1,
  B: 2,
  C: 3,
} as const;

export type TaskDifficultyValue =
  (typeof TASK_DIFFICULTY)[keyof typeof TASK_DIFFICULTY];

/**
 * Task Status
 * Represents the current state of a task
 */
export const TASK_STATUS = {
  PLANNING: 0,
  DOING: 1,
  COMPLETED: 2,
  CANCELED: 3,
} as const;

export type TaskStatusValue = (typeof TASK_STATUS)[keyof typeof TASK_STATUS];

/**
 * SubTask Collection Status
 * Represents the execution state of an episode's subtask
 */
export const SUBTASK_COLLECTION_STATUS = {
  READY: 0,
  IN_PROGRESS: 1,
  COMPLETED: 2,
  SKIPPED: 3,
  CANCELLED: 4,
} as const;

export type SubTaskCollectionStatusValue =
  (typeof SUBTASK_COLLECTION_STATUS)[keyof typeof SUBTASK_COLLECTION_STATUS];

/**
 * Gate Level
 * Represents the recording gate escalation level
 */
export const GATE_LEVEL = {
  OPEN: 0,
  BLOCK_START: 1,
  HARD_STOP: 2,
} as const;

export type GateLevelValue = (typeof GATE_LEVEL)[keyof typeof GATE_LEVEL];
