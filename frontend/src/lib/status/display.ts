import {
  EPISODE_COLLECTION_STATUS,
  type EpisodeCollectionStatusValue,
  ROBOT_STATUS,
  type RobotStatusValue,
  TASK_STATUS,
  type TaskStatusValue,
} from "./constants";

export interface StatusDisplayConfig {
  labelKey: string;
  className: string;
  isTerminal?: boolean;
  isSuccessful?: boolean;
}

export const STATUS_BADGE_STYLES = {
  blue: "bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300",
  gray: "bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300",
  green: "bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300",
  orange:
    "bg-orange-100 text-orange-700 dark:bg-orange-900 dark:text-orange-300",
  red: "bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300",
  yellow:
    "bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300",
} as const;

export const ROBOT_STATUS_DISPLAY: Record<
  RobotStatusValue,
  StatusDisplayConfig
> = {
  [ROBOT_STATUS.ONLINE]: {
    labelKey: "status.online",
    className: STATUS_BADGE_STYLES.green,
  },
  [ROBOT_STATUS.BUSY]: {
    labelKey: "status.busy",
    className: STATUS_BADGE_STYLES.yellow,
  },
  [ROBOT_STATUS.OFFLINE]: {
    labelKey: "status.offline",
    className: STATUS_BADGE_STYLES.orange,
  },
  [ROBOT_STATUS.FAULTED]: {
    labelKey: "status.faulted",
    className: STATUS_BADGE_STYLES.red,
  },
  [ROBOT_STATUS.MAINTENANCE]: {
    labelKey: "status.maintenance",
    className: STATUS_BADGE_STYLES.blue,
  },
  [ROBOT_STATUS.READY]: {
    labelKey: "status.ready",
    className: STATUS_BADGE_STYLES.gray,
  },
};

export const ROBOT_CONNECTION_STATUS_DISPLAY = {
  connected: {
    labelKey: "status.online",
    className: STATUS_BADGE_STYLES.green,
  },
  disconnected: {
    labelKey: "status.offline",
    className: STATUS_BADGE_STYLES.orange,
  },
} as const;

export const TASK_STATUS_DISPLAY: Record<TaskStatusValue, StatusDisplayConfig> =
  {
    [TASK_STATUS.PLANNING]: {
      labelKey: "status.planning",
      className: STATUS_BADGE_STYLES.gray,
    },
    [TASK_STATUS.DOING]: {
      labelKey: "status.doing",
      className: STATUS_BADGE_STYLES.blue,
    },
    [TASK_STATUS.COMPLETED]: {
      labelKey: "status.completed",
      className: STATUS_BADGE_STYLES.green,
      isTerminal: true,
      isSuccessful: true,
    },
    [TASK_STATUS.CANCELED]: {
      labelKey: "status.canceled",
      className: STATUS_BADGE_STYLES.red,
      isTerminal: true,
    },
  };

export const EPISODE_COLLECTION_STATUS_DISPLAY: Record<
  EpisodeCollectionStatusValue,
  StatusDisplayConfig
> = {
  [EPISODE_COLLECTION_STATUS.READY]: {
    labelKey: "status.ready",
    className: STATUS_BADGE_STYLES.gray,
  },
  [EPISODE_COLLECTION_STATUS.RECORDING]: {
    labelKey: "status.recording",
    className: STATUS_BADGE_STYLES.blue,
  },
  [EPISODE_COLLECTION_STATUS.CANCEL]: {
    labelKey: "status.cancel",
    className: STATUS_BADGE_STYLES.red,
    isTerminal: true,
  },
  [EPISODE_COLLECTION_STATUS.COMPLETED]: {
    labelKey: "status.completed",
    className: STATUS_BADGE_STYLES.green,
    isTerminal: true,
    isSuccessful: true,
  },
};

export function getStatusDisplay<TStatus extends string | number>(
  configs: Partial<Record<TStatus, StatusDisplayConfig>>,
  status: TStatus,
  fallback: StatusDisplayConfig
): StatusDisplayConfig {
  return configs[status] ?? fallback;
}
