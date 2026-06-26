/**
 * Status Label Hooks
 *
 * React hooks for getting translatable status labels.
 * These hooks use useTranslation to ensure components re-render on language change.
 */

import { useTranslation } from "react-i18next";

import {
  type EpisodeCollectionStatusValue,
  type RobotStatusValue,
  ROBOT_TYPE,
  type RobotTypeValue,
  SUBTASK_STATUS,
  type SubTaskStatusValue,
  TASK_PRIORITY,
  type TaskPriorityValue,
  type TaskStatusValue,
  USER_ROLE,
  type UserRoleValue,
} from "@/shared/lib/status-constants";
import {
  EPISODE_COLLECTION_STATUS_DISPLAY,
  ROBOT_STATUS_DISPLAY,
  TASK_STATUS_DISPLAY,
} from "@/shared/lib/status-display";

/**
 * Hook to get episode collection status label
 */
export function useEpisodeCollectionStatusLabel() {
  const { t } = useTranslation();

  return (status: EpisodeCollectionStatusValue): string => {
    return t(
      EPISODE_COLLECTION_STATUS_DISPLAY[status]?.labelKey ?? "status.unknown"
    );
  };
}

/**
 * Hook to get robot status label
 */
export function useRobotStatusLabel() {
  const { t } = useTranslation();

  return (status: RobotStatusValue): string => {
    return t(ROBOT_STATUS_DISPLAY[status]?.labelKey ?? "status.unknown");
  };
}

/**
 * Hook to get robot type label
 */
export function useRobotTypeLabel() {
  const { t } = useTranslation();

  return (type: RobotTypeValue): string => {
    switch (type) {
      case ROBOT_TYPE.YUBI_STATIONARY:
        return t("robotForm.robotTypeYubiStationary");
      case ROBOT_TYPE.YUBI_PORTABLE:
        return t("robotForm.robotTypeYubiPortable");
      default:
        return type;
    }
  };
}

/**
 * Hook to get user role label
 */
export function useUserRoleLabel() {
  const { t } = useTranslation();

  return (role: UserRoleValue): string => {
    switch (role) {
      case USER_ROLE.ADMIN:
        return t("status.admin");
      case USER_ROLE.DATA_ENGINEER:
        return t("status.dataEngineer");
      case USER_ROLE.MANAGER:
        return t("status.manager");
      case USER_ROLE.OPERATOR:
        return t("status.operator");
      case USER_ROLE.VIEWER:
        return t("status.viewer");
      default:
        return t("status.unknown");
    }
  };
}

/**
 * Hook to get subtask status label
 */
export function useSubTaskStatusLabel() {
  const { t } = useTranslation();

  return (status: SubTaskStatusValue): string => {
    switch (status) {
      case SUBTASK_STATUS.READY:
        return t("status.ready");
      case SUBTASK_STATUS.IN_PROGRESS:
        return t("status.inProgress");
      case SUBTASK_STATUS.FAILED:
        return t("status.failed");
      case SUBTASK_STATUS.COMPLETED:
        return t("status.completed");
      default:
        return t("status.unknown");
    }
  };
}

/**
 * Hook to get task priority label
 */
export function useTaskPriorityLabel() {
  const { t } = useTranslation();

  return (priority: TaskPriorityValue): string => {
    switch (priority) {
      case TASK_PRIORITY.LOW:
        return t("status.low");
      case TASK_PRIORITY.NORMAL:
        return t("status.normal");
      case TASK_PRIORITY.HIGH:
        return t("status.high");
      case TASK_PRIORITY.URGENT:
        return t("status.urgent");
      default:
        return t("status.unknown");
    }
  };
}

/**
 * Hook to get task status label
 */
export function useTaskStatusLabel() {
  const { t } = useTranslation();

  return (status: TaskStatusValue): string => {
    return t(TASK_STATUS_DISPLAY[status]?.labelKey ?? "status.unknown");
  };
}
