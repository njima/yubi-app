package model

import "testing"

func TestInitOrganizationMembership(t *testing.T) {
	m, err := InitOrganizationMembership("user-1", "org-1", UserRoleAdmin)
	if err != nil {
		t.Fatalf("InitOrganizationMembership() error = %v", err)
	}
	if m.UserID != "user-1" || m.OrganizationID != "org-1" || m.Role != UserRoleAdmin {
		t.Fatalf("membership = %+v", m)
	}
}

func TestInitOrganizationMembershipRequiresUserID(t *testing.T) {
	_, err := InitOrganizationMembership("", "org-1", UserRoleAdmin)
	if err == nil {
		t.Fatal("InitOrganizationMembership() expected error")
	}
}

func TestInitOrganizationMembershipRequiresOrganizationID(t *testing.T) {
	_, err := InitOrganizationMembership("user-1", "", UserRoleAdmin)
	if err == nil {
		t.Fatal("InitOrganizationMembership() expected error")
	}
}

func TestInitOrganizationMembershipRejectsInvalidRole(t *testing.T) {
	_, err := InitOrganizationMembership("user-1", "org-1", UserRole(999))
	if err == nil {
		t.Fatal("InitOrganizationMembership() expected error")
	}
}

func TestInitOrganizationMembershipAcceptsValidRoles(t *testing.T) {
	roles := []UserRole{
		UserRoleAdmin,
		UserRoleDataEngineer,
		UserRoleManager,
		UserRoleOperator,
		UserRoleViewer,
	}

	for _, role := range roles {
		_, err := InitOrganizationMembership("user-1", "org-1", role)
		if err != nil {
			t.Fatalf("InitOrganizationMembership() role %v error = %v", role, err)
		}
	}
}
