package authz

import "github.com/airoa-org/yubi-app/backend/internal/domain/model"

// rolePermissions defines the set of actions each role is allowed to perform.
// Actions follow the "resource:operation" format (e.g. "episode:create", "robot:delete").
// Permission summary by role:
//   - Admin: full access to all resources
//   - DataEngineer: Location read-only; User read+update; full access to Task/SubTask/Robot/Episode/RobotDevice
//   - Manager: same as DataEngineer
//   - Operator: Location/Task/SubTask/Robot read-only; User read+update; Episode create/update; RobotDevice full
//   - Viewer: read-only access to all resources
var rolePermissions = map[model.UserRole]map[string]bool{
	// Admin: full access
	model.UserRoleAdmin: {
		// Fleet
		"fleet:read": true,
		// Site
		"site:create": true, "site:list": true, "site:get_detail": true,
		"site:update": true, "site:delete": true,
		// Location
		"location:create": true, "location:list": true, "location:get_detail": true,
		"location:update": true, "location:delete": true,
		// User
		"user:me": true, "user:update_self": true, "user:create": true, "user:list": true, "user:get_detail": true,
		"user:update": true, "user:update_role": true, "user:update_location": true, "user:update_site": true, "user:delete": true,
		"user:grant_permission": true, "user:revoke_permission": true,
		"user:list_permissions": true,
		// Task
		"task:create": true, "task:list": true, "task:get_detail": true,
		"task:update": true, "task:list_versions": true,
		"task_version:create": true, "task_version:update": true, "task_version:approve": true,
		"task_tag:list": true, "task_tag:create": true,
		// SubTask
		"subtask:create": true, "subtask:list": true, "subtask:get_detail": true,
		"subtask:update": true, "subtask:delete": true,
		// Robot
		"robot:create": true, "robot:list": true, "robot:get_detail": true,
		"robot:update": true, "robot:delete": true, "robot:status_stream": true, "robot:operate": true,
		// Episode
		"episode:create": true, "episode:list": true, "episode:get_detail": true,
		"episode:update": true, "episode:stream": true, "episode:list_stream": true, "robot:teleop_stream": true,
		"episode_grade:list": true, "episode_grade:update": true,
		// Robot Device
		"robot_device:me": true, "robot_device:list_episodes": true,
		"robot_device:update_status": true,
		"robot_device:get_episode":   true,
		"robot_device:start_episode": true, "robot_device:finish_episode": true,
		"robot_device:cancel_episode": true, "robot_device:complete_subtask": true,
		"robot_device:skip_subtask": true, "robot_device:create_execution": true,
		"robot_device:start_execution": true, "robot_device:finish_execution": true,
		"robot_device:cancel_execution":    true,
		"robot_device:repeat_last_episode": true,
		// API Key
		"api_key:list":   true,
		"api_key:create": true,
		"api_key:get":    true,
		"api_key:update": true,
		"api_key:revoke": true,
	},
	// DataEngineer: Location read-only, User read + update, full access to Task/SubTask/Robot/Episode/RobotDevice
	model.UserRoleDataEngineer: {
		// Fleet
		"fleet:read": true,
		// Site
		"site:list": true, "site:get_detail": true,
		// Location
		"location:list": true, "location:get_detail": true,
		// User
		"user:me": true, "user:update_self": true, "user:list": true, "user:get_detail": true,
		"user:update":           true,
		"user:list_permissions": true,
		// Task
		"task:create": true, "task:list": true, "task:get_detail": true,
		"task:update": true, "task:list_versions": true,
		"task_version:create": true, "task_version:update": true, "task_version:approve": true,
		"task_tag:list": true, "task_tag:create": true,
		// SubTask
		"subtask:create": true, "subtask:list": true, "subtask:get_detail": true,
		"subtask:update": true, "subtask:delete": true,
		// Robot
		"robot:create": true, "robot:list": true, "robot:get_detail": true,
		"robot:update": true, "robot:delete": true, "robot:status_stream": true, "robot:operate": true,
		// Episode
		"episode:create": true, "episode:list": true, "episode:get_detail": true,
		"episode:update": true, "episode:stream": true, "episode:list_stream": true, "robot:teleop_stream": true,
		"episode_grade:list": true, "episode_grade:update": true,
		// Robot Device
		"robot_device:me": true, "robot_device:list_episodes": true,
		"robot_device:update_status": true,
		"robot_device:get_episode":   true,
		"robot_device:start_episode": true, "robot_device:finish_episode": true,
		"robot_device:cancel_episode": true, "robot_device:complete_subtask": true,
		"robot_device:skip_subtask": true, "robot_device:create_execution": true,
		"robot_device:start_execution": true, "robot_device:finish_execution": true,
		"robot_device:cancel_execution":    true,
		"robot_device:repeat_last_episode": true,
	},
	// Manager: same as DataEngineer
	model.UserRoleManager: {
		// Fleet
		"fleet:read": true,
		// Site
		"site:list": true, "site:get_detail": true,
		// Location
		"location:list": true, "location:get_detail": true,
		// User
		"user:me": true, "user:update_self": true, "user:list": true, "user:get_detail": true,
		"user:update":           true,
		"user:list_permissions": true,
		// Task
		"task:create": true, "task:list": true, "task:get_detail": true,
		"task:update": true, "task:list_versions": true,
		"task_version:create": true, "task_version:update": true, "task_version:approve": true,
		"task_tag:list": true, "task_tag:create": true,
		// SubTask
		"subtask:create": true, "subtask:list": true, "subtask:get_detail": true,
		"subtask:update": true, "subtask:delete": true,
		// Robot
		"robot:create": true, "robot:list": true, "robot:get_detail": true,
		"robot:update": true, "robot:delete": true, "robot:status_stream": true, "robot:operate": true,
		// Episode
		"episode:create": true, "episode:list": true, "episode:get_detail": true,
		"episode:update": true, "episode:stream": true, "episode:list_stream": true, "robot:teleop_stream": true,
		"episode_grade:list": true, "episode_grade:update": true,
		// Robot Device
		"robot_device:me": true, "robot_device:list_episodes": true,
		"robot_device:update_status": true,
		"robot_device:get_episode":   true,
		"robot_device:start_episode": true, "robot_device:finish_episode": true,
		"robot_device:cancel_episode": true, "robot_device:complete_subtask": true,
		"robot_device:skip_subtask": true, "robot_device:create_execution": true,
		"robot_device:start_execution": true, "robot_device:finish_execution": true,
		"robot_device:cancel_execution":    true,
		"robot_device:repeat_last_episode": true,
	},
	// Operator: Location/Robot/Task/SubTask read-only, User read + update, Episode create/update, RobotDevice full
	model.UserRoleOperator: {
		// Fleet
		"fleet:read": true,
		// Site
		"site:list": true, "site:get_detail": true,
		// Location
		"location:list": true, "location:get_detail": true,
		// User
		"user:me": true, "user:update_self": true, "user:list": true, "user:get_detail": true,
		"user:update":           true,
		"user:list_permissions": true,
		// Task
		"task:list": true, "task:get_detail": true, "task:list_versions": true,
		"task_tag:list": true,
		// SubTask
		"subtask:list": true, "subtask:get_detail": true,
		// Robot
		"robot:list": true, "robot:get_detail": true, "robot:status_stream": true,
		"robot:operate": true,
		// Episode
		"episode:create": true, "episode:list": true, "episode:get_detail": true,
		"episode:update": true, "episode:stream": true, "episode:list_stream": true, "robot:teleop_stream": true,
		"episode_grade:list": true, "episode_grade:update": true,
		// Robot Device
		"robot_device:me": true, "robot_device:list_episodes": true,
		"robot_device:update_status": true,
		"robot_device:get_episode":   true,
		"robot_device:start_episode": true, "robot_device:finish_episode": true,
		"robot_device:cancel_episode": true, "robot_device:complete_subtask": true,
		"robot_device:skip_subtask": true, "robot_device:create_execution": true,
		"robot_device:start_execution": true, "robot_device:finish_execution": true,
		"robot_device:cancel_execution":    true,
		"robot_device:repeat_last_episode": true,
	},
	// Viewer: read-only access
	model.UserRoleViewer: {
		// Fleet
		"fleet:read": true,
		// Site
		"site:list": true, "site:get_detail": true,
		// Location
		"location:list": true, "location:get_detail": true,
		// User
		"user:me": true, "user:update_self": true, "user:list": true, "user:get_detail": true,
		"user:list_permissions": true,
		// Task
		"task:list": true, "task:get_detail": true, "task:list_versions": true,
		"task_tag:list": true,
		// SubTask
		"subtask:list": true, "subtask:get_detail": true,
		// Robot
		"robot:list": true, "robot:get_detail": true, "robot:status_stream": true,
		// Episode
		"episode:list": true, "episode:get_detail": true,
		"episode:stream": true, "episode:list_stream": true, "robot:teleop_stream": true,
		"episode_grade:list": true,
	},
}

// operationPermissions maps OpenAPI operationIDs to permission actions.
// Operations registered in authzBypassOperations are not listed here.
var operationPermissions = map[string]string{
	// Fleet
	"GetFleetSummary":         "fleet:read",
	"GetFleetStats":           "fleet:read",
	"GetFleetCollectionTrend": "fleet:read",
	// Episode
	"ListEpisodes":         "episode:list",
	"ExportEpisodes":       "episode:list",
	"ExportOperatorYield":  "episode:list",
	"CreateEpisode":        "episode:create",
	"CreateEpisodesBulk":   "episode:create",
	"DeleteEpisodeById":    "episode:delete", // NOTE: No role is assigned intentionally — this operation is deprecated.
	"GetEpisodeById":       "episode:get_detail",
	"GetEpisodeRecordings": "episode:get_detail",
	"GetEpisodeStats":      "episode:get_detail",
	"UpdateEpisodeById":    "episode:update",
	"GetMyEpisodeGrade":    "episode_grade:list",
	"ListEpisodeGrades":    "episode_grade:list",
	"UpdateMyEpisodeGrade": "episode_grade:update",
	// Site
	"ListSites":      "site:list",
	"CreateSite":     "site:create",
	"DeleteSiteById": "site:delete",
	"GetSiteById":    "site:get_detail",
	"UpdateSiteById": "site:update",
	// Location
	"ListLocations":      "location:list",
	"CreateLocation":     "location:create",
	"DeleteLocationById": "location:delete",
	"GetLocationById":    "location:get_detail",
	"UpdateLocationById": "location:update",
	// User
	"GetMe":               "user:me",
	"UpdateMe":            "user:update_self",
	"ListUsers":           "user:list",
	"CreateUser":          "user:create",
	"DeleteUserById":      "user:delete",
	"GetUserById":         "user:get_detail",
	"UpdateUserById":      "user:update",
	"UpdateUserRole":      "user:update_role",
	"UpdateUserLocations": "user:update_location",
	"UpdateUserSites":     "user:update_site",
	// User Permission
	"GrantUserPermission":  "user:grant_permission",
	"RevokeUserPermission": "user:revoke_permission",
	"ListUserPermissions":  "user:list_permissions",
	// Robot
	"ListRobots":         "robot:list",
	"ListRobotTypes":     "robot:list",
	"CreateRobot":        "robot:create",
	"DeleteRobotById":    "robot:delete",
	"GetRobotById":       "robot:get_detail",
	"UpdateRobotById":    "robot:update",
	"GetRobotOperator":   "robot:operate",
	"SetRobotOperator":   "robot:operate",
	"ClearRobotOperator": "robot:operate",
	// Task
	"ListTasks":                   "task:list",
	"CreateTask":                  "task:create",
	"DeleteTaskById":              "task:delete", // NOTE: No role is assigned intentionally — this operation is deprecated.
	"GetTaskById":                 "task:get_detail",
	"UpdateTaskById":              "task:update",
	"ListTaskVersions":            "task:list_versions",
	"CreateTaskVersion":           "task_version:create",
	"UpdateTaskVersion":           "task_version:update",
	"ApproveTaskVersion":          "task_version:approve",
	"UpdateTaskVersionParameters": "task_version:create",
	// Task Summary
	"GetTaskSummary":         "task:list",
	"GetTaskCompletionTrend": "task:list",
	"GetTaskAvailableTags":   "task_tag:list",
	// Task Export
	"ExportTasks": "task:list",
	// Task Import
	"ValidateTaskImport": "task:create",
	"ImportTasks":        "task:create",
	// User Import
	"ValidateUserImport": "user:create",
	"ImportUsers":        "user:create",
	// Task Tag
	"ListTaskCategoryTypes": "task_tag:list",
	"ListTaskTags":          "task_tag:list",
	"CreateTaskTag":         "task_tag:create",
	// SubTask
	"ListSubTasks":      "subtask:list",
	"CreateSubTask":     "subtask:create",
	"DeleteSubTaskById": "subtask:delete",
	"GetSubTaskById":    "subtask:get_detail",
	"UpdateSubTaskById": "subtask:update",
	"ReorderSubTasks":   "subtask:update",
	"CompleteSubTask":   "subtask:complete", // NOTE: No role is assigned intentionally — this operation is deprecated.
	// Robot Device
	"GetRobotMe":             "robot_device:me",
	"UpdateRobotStatus":      "robot_device:update_status",
	"ListRobotEpisodes":      "robot_device:list_episodes",
	"GetRobotEpisodeById":    "robot_device:get_episode",
	"StartRobotEpisode":      "robot_device:start_episode",
	"FinishRobotEpisode":     "robot_device:finish_episode",
	"CancelRobotEpisode":     "robot_device:cancel_episode",
	"CompleteRobotSubTask":   "robot_device:complete_subtask",
	"SkipRobotSubTask":       "robot_device:skip_subtask",
	"CreateRobotExecution":   "robot_device:create_execution",
	"StartRobotExecution":    "robot_device:start_execution",
	"FinishRobotExecution":   "robot_device:finish_execution",
	"CancelRobotExecution":   "robot_device:cancel_execution",
	"RepeatLastRobotEpisode": "robot_device:repeat_last_episode",
	// Organization
	// NOTE: ListOrganizations and GetOrganizationById are in authzBypassOperations (all roles allowed).
	"CreateOrganization":     "org:create",
	"DeleteOrganizationById": "org:delete",
	"UpdateOrganizationById": "org:update",
	// API Key
	"ListApiKeys":  "api_key:list",
	"CreateApiKey": "api_key:create",
	"GetApiKey":    "api_key:get",
	"UpdateApiKey": "api_key:update",
	"RevokeApiKey": "api_key:revoke",
}

// authzBypassOperations is the set of operationIDs that skip authorization checks.
// Any authenticated user can access these operations regardless of their role.
var authzBypassOperations = map[string]bool{
	// Organization (read-only) — all roles can view organization info.
	"ListOrganizations":   true,
	"GetOrganizationById": true,
}

// HasPermission reports whether the given role is allowed to perform action.
// action must follow the "resource:operation" format (e.g. "episode:create", "robot:delete").
// Returns false for unknown roles or unregistered actions.
func HasPermission(role model.UserRole, action string) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	return perms[action]
}

// ActionForOperation returns the permission action mapped to the given OpenAPI operationID.
// Returns an empty string and false if no mapping exists.
// Operations registered in authzBypassOperations are not included here; they always return ("", false).
func ActionForOperation(operationID string) (string, bool) {
	action, ok := operationPermissions[operationID]
	return action, ok
}

// IsAuthzBypassOperation reports whether the given operationID bypasses authorization checks.
// When true, the operation is executed without any role permission verification.
func IsAuthzBypassOperation(operationID string) bool {
	return authzBypassOperations[operationID]
}
