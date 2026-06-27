import { USER_ROLE, type UserRoleValue } from "./status/constants";

export type PermissionAction =
  | "robot:create"
  | "robot:update"
  | "robot:delete"
  | "task:create"
  | "task:update"
  | "subtask:create"
  | "subtask:update"
  | "subtask:delete"
  | "episode:create"
  | "episode:update"
  | "user:create"
  | "user:delete"
  | "user:update_role"
  | "location:create"
  | "location:update"
  | "location:delete"
  | "api_key:list"
  | "api_key:create"
  | "api_key:revoke";

const PERMISSIONS: Record<PermissionAction, UserRoleValue[]> = {
  "robot:create": [USER_ROLE.ADMIN, USER_ROLE.DATA_ENGINEER, USER_ROLE.MANAGER],
  "robot:update": [USER_ROLE.ADMIN, USER_ROLE.DATA_ENGINEER, USER_ROLE.MANAGER],
  "robot:delete": [USER_ROLE.ADMIN, USER_ROLE.DATA_ENGINEER, USER_ROLE.MANAGER],
  "task:create": [USER_ROLE.ADMIN, USER_ROLE.DATA_ENGINEER, USER_ROLE.MANAGER],
  "task:update": [USER_ROLE.ADMIN, USER_ROLE.DATA_ENGINEER, USER_ROLE.MANAGER],
  "subtask:create": [
    USER_ROLE.ADMIN,
    USER_ROLE.DATA_ENGINEER,
    USER_ROLE.MANAGER,
  ],
  "subtask:update": [
    USER_ROLE.ADMIN,
    USER_ROLE.DATA_ENGINEER,
    USER_ROLE.MANAGER,
  ],
  "subtask:delete": [
    USER_ROLE.ADMIN,
    USER_ROLE.DATA_ENGINEER,
    USER_ROLE.MANAGER,
  ],
  "episode:create": [
    USER_ROLE.ADMIN,
    USER_ROLE.DATA_ENGINEER,
    USER_ROLE.MANAGER,
    USER_ROLE.OPERATOR,
  ],
  "episode:update": [
    USER_ROLE.ADMIN,
    USER_ROLE.DATA_ENGINEER,
    USER_ROLE.MANAGER,
    USER_ROLE.OPERATOR,
  ],
  "user:create": [USER_ROLE.ADMIN],
  "user:delete": [USER_ROLE.ADMIN],
  "user:update_role": [USER_ROLE.ADMIN],
  "location:create": [USER_ROLE.ADMIN],
  "location:update": [USER_ROLE.ADMIN],
  "location:delete": [USER_ROLE.ADMIN],
  "api_key:list": [USER_ROLE.ADMIN],
  "api_key:create": [USER_ROLE.ADMIN],
  "api_key:revoke": [USER_ROLE.ADMIN],
};

export function hasPermission(
  role: UserRoleValue | undefined,
  action: PermissionAction
): boolean {
  if (role === undefined) return false;
  return PERMISSIONS[action].includes(role);
}
