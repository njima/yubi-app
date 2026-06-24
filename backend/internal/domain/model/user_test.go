package model

import (
	"strings"
	"testing"
	"time"
)

func TestInitUser(t *testing.T) {
	tests := []struct {
		name           string
		organizationID string
		userName       string
		email          string
		role           UserRole
		wantErr        bool
	}{
		{
			name:           "success with valid inputs",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			userName:       "Test User",
			email:          "test@example.com",
			role:           UserRoleAdmin,
			wantErr:        false,
		},
		{
			name:           "success with operator role",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			userName:       "Operator User",
			email:          "operator@example.com",
			role:           UserRoleOperator,
			wantErr:        false,
		},
		{
			name:           "error when organization_id is empty",
			organizationID: "",
			userName:       "Test User",
			email:          "test@example.com",
			role:           UserRoleAdmin,
			wantErr:        true,
		},
		{
			name:           "error when name is empty",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			userName:       "",
			email:          "test@example.com",
			role:           UserRoleAdmin,
			wantErr:        true,
		},
		{
			name:           "error when name exceeds 60 characters",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			userName:       strings.Repeat("a", 61),
			email:          "test@example.com",
			role:           UserRoleAdmin,
			wantErr:        true,
		},
		{
			name:           "success when name is exactly 60 characters",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			userName:       strings.Repeat("a", 60),
			email:          "test@example.com",
			role:           UserRoleAdmin,
			wantErr:        false,
		},
		{
			name:           "error when email is empty",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			userName:       "Test User",
			email:          "",
			role:           UserRoleAdmin,
			wantErr:        true,
		},
		{
			name:           "error when email format is invalid",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			userName:       "Test User",
			email:          "invalid-email",
			role:           UserRoleAdmin,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitUser(tt.organizationID, tt.userName, tt.email, tt.role)

			if tt.wantErr {
				if err == nil {
					t.Errorf("InitUser() error = nil, wantErr %v", tt.wantErr)
				}
				if got.IDNatural != "" {
					t.Errorf("InitUser() IDNatural = %v, want empty", got.IDNatural)
				}
				return
			}

			if err != nil {
				t.Errorf("InitUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.IDNatural == "" {
				t.Errorf("InitUser() IDNatural is empty")
			}
			if got.OrganizationID != tt.organizationID {
				t.Errorf("InitUser() OrganizationID = %v, want %v", got.OrganizationID, tt.organizationID)
			}
			if got.Name != tt.userName {
				t.Errorf("InitUser() Name = %v, want %v", got.Name, tt.userName)
			}
			if got.Email != tt.email {
				t.Errorf("InitUser() Email = %v, want %v", got.Email, tt.email)
			}
			if got.Role != tt.role {
				t.Errorf("InitUser() Role = %v, want %v", got.Role, tt.role)
			}
			if got.CreatedAt.IsZero() {
				t.Errorf("InitUser() CreatedAt is zero")
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	now := time.Now()
	updatedAt := now.Add(time.Hour)

	tests := []struct {
		name           string
		id             int64
		idNatural      string
		organizationID string
		userName       string
		email          string
		role           UserRole
		createdAt      time.Time
		updatedAt      *time.Time
	}{
		{
			name:           "success with all fields",
			id:             1,
			idNatural:      "550e8400-e29b-41d4-a716-446655440000",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			userName:       "Test User",
			email:          "test@example.com",
			role:           UserRoleAdmin,
			createdAt:      now,
			updatedAt:      &updatedAt,
		},
		{
			name:           "success with nil updatedAt",
			id:             2,
			idNatural:      "550e8400-e29b-41d4-a716-446655440002",
			organizationID: "550e8400-e29b-41d4-a716-446655440003",
			userName:       "Another User",
			email:          "another@example.com",
			role:           UserRoleOperator,
			createdAt:      now,
			updatedAt:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUser(tt.id, tt.idNatural, tt.organizationID, tt.userName, tt.email, tt.role, tt.createdAt, tt.updatedAt)

			if got.ID != tt.id {
				t.Errorf("NewUser() ID = %v, want %v", got.ID, tt.id)
			}
			if got.IDNatural != tt.idNatural {
				t.Errorf("NewUser() IDNatural = %v, want %v", got.IDNatural, tt.idNatural)
			}
			if got.OrganizationID != tt.organizationID {
				t.Errorf("NewUser() OrganizationID = %v, want %v", got.OrganizationID, tt.organizationID)
			}
			if got.Name != tt.userName {
				t.Errorf("NewUser() Name = %v, want %v", got.Name, tt.userName)
			}
			if got.Email != tt.email {
				t.Errorf("NewUser() Email = %v, want %v", got.Email, tt.email)
			}
			if got.Role != tt.role {
				t.Errorf("NewUser() Role = %v, want %v", got.Role, tt.role)
			}
			if got.CreatedAt != tt.createdAt {
				t.Errorf("NewUser() CreatedAt = %v, want %v", got.CreatedAt, tt.createdAt)
			}
		})
	}
}

func newValidUser() User {
	return User{
		IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
		OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
		Name:           "Test User",
		Email:          "test@example.com",
		Role:           UserRoleAdmin,
	}
}

func TestUser_validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name:    "valid with all required fields",
			user:    newValidUser(),
			wantErr: false,
		},
		{
			name: "error when id_natural is empty",
			user: User{
				IDNatural:      "",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				Name:           "Test User",
				Email:          "test@example.com",
				Role:           UserRoleAdmin,
			},
			wantErr: true,
		},
		{
			name: "error when organization_id is empty",
			user: User{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "",
				Name:           "Test User",
				Email:          "test@example.com",
				Role:           UserRoleAdmin,
			},
			wantErr: true,
		},
		{
			name: "error when name is empty",
			user: User{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				Name:           "",
				Email:          "test@example.com",
				Role:           UserRoleAdmin,
			},
			wantErr: true,
		},
		{
			name: "error when name exceeds 60 characters",
			user: User{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				Name:           strings.Repeat("a", 61),
				Email:          "test@example.com",
				Role:           UserRoleAdmin,
			},
			wantErr: true,
		},
		{
			name: "error when email is empty",
			user: User{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				Name:           "Test User",
				Email:          "",
				Role:           UserRoleAdmin,
			},
			wantErr: true,
		},
		{
			name: "error when email format is invalid",
			user: User{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				Name:           "Test User",
				Email:          "invalid-email",
				Role:           UserRoleAdmin,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("User.validate() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("User.validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestUser_SetName(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		wantErr  bool
	}{
		{
			name:     "success with valid name",
			userName: "New User",
			wantErr:  false,
		},
		{
			name:     "success with name at max length",
			userName: strings.Repeat("b", 60),
			wantErr:  false,
		},
		{
			name:     "error when name is empty",
			userName: "",
			wantErr:  true,
		},
		{
			name:     "error when name exceeds 60 characters",
			userName: strings.Repeat("c", 61),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newValidUser()
			err := u.SetName(tt.userName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("User.SetName() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("User.SetName() error = %v, wantErr %v", err, tt.wantErr)
				}
				if u.Name != tt.userName {
					t.Errorf("User.Name = %v, want %v", u.Name, tt.userName)
				}
			}
		})
	}
}

func TestUser_SetEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "success with valid email",
			email:   "newemail@example.com",
			wantErr: false,
		},
		{
			name:    "success with email containing subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "error when email is empty",
			email:   "",
			wantErr: true,
		},
		{
			name:    "error when email format is invalid - no at sign",
			email:   "invalidemail.com",
			wantErr: true,
		},
		{
			name:    "error when email format is invalid - no domain",
			email:   "user@",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newValidUser()
			err := u.SetEmail(tt.email)

			if tt.wantErr {
				if err == nil {
					t.Errorf("User.SetEmail() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("User.SetEmail() error = %v, wantErr %v", err, tt.wantErr)
				}
				if u.Email != tt.email {
					t.Errorf("User.Email = %v, want %v", u.Email, tt.email)
				}
			}
		})
	}
}

func TestUser_SetRole(t *testing.T) {
	tests := []struct {
		name    string
		role    UserRole
		wantErr bool
	}{
		{
			name:    "success with Admin role",
			role:    UserRoleAdmin,
			wantErr: false,
		},
		{
			name:    "success with Operator role",
			role:    UserRoleOperator,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newValidUser()
			err := u.SetRole(tt.role)

			if tt.wantErr {
				if err == nil {
					t.Errorf("User.SetRole() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("User.SetRole() error = %v, wantErr %v", err, tt.wantErr)
				}
				if u.Role != tt.role {
					t.Errorf("User.Role = %v, want %v", u.Role, tt.role)
				}
			}
		})
	}
}
