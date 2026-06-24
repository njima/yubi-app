package authz

import (
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name   string
		role   model.UserRole
		action string
		want   bool
	}{
		// Admin: full access
		{name: "Admin can create location", role: model.UserRoleAdmin, action: "location:create", want: true},
		{name: "Admin can delete location", role: model.UserRoleAdmin, action: "location:delete", want: true},
		{name: "Admin can create user", role: model.UserRoleAdmin, action: "user:create", want: true},
		{name: "Admin can update role", role: model.UserRoleAdmin, action: "user:update_role", want: true},
		{name: "Admin can grant permission", role: model.UserRoleAdmin, action: "user:grant_permission", want: true},
		{name: "Admin can create robot", role: model.UserRoleAdmin, action: "robot:create", want: true},
		{name: "Admin can delete robot", role: model.UserRoleAdmin, action: "robot:delete", want: true},
		{name: "Admin can create episode", role: model.UserRoleAdmin, action: "episode:create", want: true},

		// DataEngineer: no location create/delete, no user create/delete/update_role
		{name: "DataEngineer can list location", role: model.UserRoleDataEngineer, action: "location:list", want: true},
		{name: "DataEngineer cannot create location", role: model.UserRoleDataEngineer, action: "location:create", want: false},
		{name: "DataEngineer cannot delete location", role: model.UserRoleDataEngineer, action: "location:delete", want: false},
		{name: "DataEngineer can create task", role: model.UserRoleDataEngineer, action: "task:create", want: true},
		{name: "DataEngineer cannot create user", role: model.UserRoleDataEngineer, action: "user:create", want: false},
		{name: "DataEngineer cannot update role", role: model.UserRoleDataEngineer, action: "user:update_role", want: false},
		{name: "DataEngineer can create robot", role: model.UserRoleDataEngineer, action: "robot:create", want: true},
		{name: "DataEngineer can create episode", role: model.UserRoleDataEngineer, action: "episode:create", want: true},

		// Manager: same as DataEngineer
		{name: "Manager can list location", role: model.UserRoleManager, action: "location:list", want: true},
		{name: "Manager cannot create location", role: model.UserRoleManager, action: "location:create", want: false},
		{name: "Manager can create task", role: model.UserRoleManager, action: "task:create", want: true},
		{name: "Manager cannot create user", role: model.UserRoleManager, action: "user:create", want: false},
		{name: "Manager can create robot", role: model.UserRoleManager, action: "robot:create", want: true},

		// Operator: read-only for location/task/subtask/robot, can create/update episode
		{name: "Operator can list location", role: model.UserRoleOperator, action: "location:list", want: true},
		{name: "Operator cannot create location", role: model.UserRoleOperator, action: "location:create", want: false},
		{name: "Operator can list task", role: model.UserRoleOperator, action: "task:list", want: true},
		{name: "Operator cannot create task", role: model.UserRoleOperator, action: "task:create", want: false},
		{name: "Operator can list robot", role: model.UserRoleOperator, action: "robot:list", want: true},
		{name: "Operator cannot create robot", role: model.UserRoleOperator, action: "robot:create", want: false},
		{name: "Operator can create episode", role: model.UserRoleOperator, action: "episode:create", want: true},
		{name: "Operator can update episode", role: model.UserRoleOperator, action: "episode:update", want: true},
		{name: "Operator can access robot_device", role: model.UserRoleOperator, action: "robot_device:me", want: true},

		// Viewer: read-only access, no episode create
		{name: "Viewer can list location", role: model.UserRoleViewer, action: "location:list", want: true},
		{name: "Viewer cannot create location", role: model.UserRoleViewer, action: "location:create", want: false},
		{name: "Viewer can list episode", role: model.UserRoleViewer, action: "episode:list", want: true},
		{name: "Viewer cannot create episode", role: model.UserRoleViewer, action: "episode:create", want: false},
		{name: "Viewer can list robot", role: model.UserRoleViewer, action: "robot:list", want: true},
		{name: "Viewer cannot create robot", role: model.UserRoleViewer, action: "robot:create", want: false},
		{name: "Viewer cannot access robot_device me", role: model.UserRoleViewer, action: "robot_device:me", want: false},

		// Self-profile update is allowed for every role
		{name: "Admin can update self", role: model.UserRoleAdmin, action: "user:update_self", want: true},
		{name: "DataEngineer can update self", role: model.UserRoleDataEngineer, action: "user:update_self", want: true},
		{name: "Manager can update self", role: model.UserRoleManager, action: "user:update_self", want: true},
		{name: "Operator can update self", role: model.UserRoleOperator, action: "user:update_self", want: true},
		{name: "Viewer can update self", role: model.UserRoleViewer, action: "user:update_self", want: true},
		// Viewer still cannot update arbitrary users
		{name: "Viewer cannot update other user", role: model.UserRoleViewer, action: "user:update", want: false},

		// Unknown role (use an integer value that doesn't map to any defined role)
		{name: "Unknown role returns false", role: model.UserRole(9999), action: "location:list", want: false},

		// Non-existent action
		{name: "Admin with non-existent action returns false", role: model.UserRoleAdmin, action: "nonexistent:action", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasPermission(tt.role, tt.action)
			if got != tt.want {
				t.Errorf("HasPermission(%v, %v) = %v, want %v", tt.role, tt.action, got, tt.want)
			}
		})
	}
}

func TestActionForOperation(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		wantAction  string
		wantFound   bool
	}{
		{name: "ListEpisodes", operationID: "ListEpisodes", wantAction: "episode:list", wantFound: true},
		{name: "CreateEpisode", operationID: "CreateEpisode", wantAction: "episode:create", wantFound: true},
		{name: "GetEpisodeById", operationID: "GetEpisodeById", wantAction: "episode:get_detail", wantFound: true},
		{name: "GetEpisodeRecordings", operationID: "GetEpisodeRecordings", wantAction: "episode:get_detail", wantFound: true},
		{name: "GetEpisodeStats", operationID: "GetEpisodeStats", wantAction: "episode:get_detail", wantFound: true},
		{name: "UpdateEpisodeById", operationID: "UpdateEpisodeById", wantAction: "episode:update", wantFound: true},
		{name: "ExportEpisodes", operationID: "ExportEpisodes", wantAction: "episode:list", wantFound: true},
		{name: "ExportOperatorYield", operationID: "ExportOperatorYield", wantAction: "episode:list", wantFound: true},
		{name: "ListLocations", operationID: "ListLocations", wantAction: "location:list", wantFound: true},
		{name: "CreateLocation", operationID: "CreateLocation", wantAction: "location:create", wantFound: true},
		{name: "GetMe", operationID: "GetMe", wantAction: "user:me", wantFound: true},
		{name: "UpdateMe", operationID: "UpdateMe", wantAction: "user:update_self", wantFound: true},
		{name: "ListUsers", operationID: "ListUsers", wantAction: "user:list", wantFound: true},
		{name: "UpdateUserRole", operationID: "UpdateUserRole", wantAction: "user:update_role", wantFound: true},
		{name: "UpdateUserSites", operationID: "UpdateUserSites", wantAction: "user:update_site", wantFound: true},
		{name: "ListRobots", operationID: "ListRobots", wantAction: "robot:list", wantFound: true},
		{name: "CreateRobot", operationID: "CreateRobot", wantAction: "robot:create", wantFound: true},
		{name: "ListTasks", operationID: "ListTasks", wantAction: "task:list", wantFound: true},
		{name: "GetRobotMe", operationID: "GetRobotMe", wantAction: "robot_device:me", wantFound: true},
		{name: "StartRobotEpisode", operationID: "StartRobotEpisode", wantAction: "robot_device:start_episode", wantFound: true},
		{name: "GrantUserPermission", operationID: "GrantUserPermission", wantAction: "user:grant_permission", wantFound: true},

		// Undefined operation
		{name: "undefined operation returns empty and false", operationID: "UndefinedOperation", wantAction: "", wantFound: false},
		{name: "empty operationID returns empty and false", operationID: "", wantAction: "", wantFound: false},

		// Bypass operations are NOT in operationPermissions
		{name: "ListOrganizations is bypass, not in operationPermissions", operationID: "ListOrganizations", wantAction: "", wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAction, gotFound := ActionForOperation(tt.operationID)
			if gotFound != tt.wantFound {
				t.Errorf("ActionForOperation(%v) found = %v, want %v", tt.operationID, gotFound, tt.wantFound)
			}
			if gotAction != tt.wantAction {
				t.Errorf("ActionForOperation(%v) action = %v, want %v", tt.operationID, gotAction, tt.wantAction)
			}
		})
	}
}

func TestIsAuthzBypassOperation(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		want        bool
	}{
		{name: "ListOrganizations is bypass", operationID: "ListOrganizations", want: true},
		{name: "GetOrganizationById is bypass", operationID: "GetOrganizationById", want: true},

		// Non-bypass operations
		{name: "ListEpisodes is not bypass", operationID: "ListEpisodes", want: false},
		{name: "CreateEpisode is not bypass", operationID: "CreateEpisode", want: false},
		{name: "ListUsers is not bypass", operationID: "ListUsers", want: false},
		{name: "GetMe is not bypass", operationID: "GetMe", want: false},
		{name: "empty string is not bypass", operationID: "", want: false},
		{name: "unknown operation is not bypass", operationID: "UnknownOperation", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAuthzBypassOperation(tt.operationID)
			if got != tt.want {
				t.Errorf("IsAuthzBypassOperation(%v) = %v, want %v", tt.operationID, got, tt.want)
			}
		})
	}
}
