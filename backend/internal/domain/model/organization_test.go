package model

import (
	"strings"
	"testing"
	"time"
)

func TestInitOrganization(t *testing.T) {
	description := "Test organization description"

	tests := []struct {
		name        string
		orgName     string
		description *string
		kind        OrganizationKind
		wantErr     bool
	}{
		{
			name:        "success with valid inputs",
			orgName:     "Test Organization",
			description: &description,
			kind:        OrganizationKindTeam,
			wantErr:     false,
		},
		{
			name:        "success with nil description",
			orgName:     "Test Organization",
			description: nil,
			kind:        OrganizationKindTeam,
			wantErr:     false,
		},
		{
			name:        "error when name is empty",
			orgName:     "",
			description: &description,
			kind:        OrganizationKindTeam,
			wantErr:     true,
		},
		{
			name:        "error when name exceeds 100 characters",
			orgName:     strings.Repeat("a", 101),
			description: &description,
			kind:        OrganizationKindTeam,
			wantErr:     true,
		},
		{
			name:        "success when name is exactly 100 characters",
			orgName:     strings.Repeat("a", 100),
			description: &description,
			kind:        OrganizationKindTeam,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitOrganization(tt.orgName, tt.description, tt.kind)

			if tt.wantErr {
				if err == nil {
					t.Errorf("InitOrganization() error = nil, wantErr %v", tt.wantErr)
				}
				if got.IDNatural != "" {
					t.Errorf("InitOrganization() IDNatural = %v, want empty", got.IDNatural)
				}
				return
			}

			if err != nil {
				t.Errorf("InitOrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.IDNatural == "" {
				t.Errorf("InitOrganization() IDNatural is empty")
			}
			if got.Name != tt.orgName {
				t.Errorf("InitOrganization() Name = %v, want %v", got.Name, tt.orgName)
			}
			if got.Kind != tt.kind {
				t.Errorf("InitOrganization() Kind = %v, want %v", got.Kind, tt.kind)
			}
			if got.Description == nil {
				t.Errorf("InitOrganization() Description is nil")
			}
			if got.CreatedAt.IsZero() {
				t.Errorf("InitOrganization() CreatedAt is zero")
			}
		})
	}
}

func TestInitOrganizationWithKind(t *testing.T) {
	org, err := InitOrganization("Ada's Workspace", nil, OrganizationKindPersonal)
	if err != nil {
		t.Fatalf("InitOrganization() error = %v", err)
	}
	if org.Kind != OrganizationKindPersonal {
		t.Fatalf("Kind = %v", org.Kind)
	}
}

func TestInitOrganizationRejectsInvalidKind(t *testing.T) {
	_, err := InitOrganization("Bad Workspace", nil, OrganizationKind("invalid"))
	if err == nil {
		t.Fatal("InitOrganization() expected error for invalid kind")
	}
}

func TestNewOrganization(t *testing.T) {
	now := time.Now()
	updatedAt := now
	description := "Test description"

	tests := []struct {
		name        string
		id          int64
		idNatural   string
		orgName     string
		kind        OrganizationKind
		description *string
		createdAt   time.Time
		updatedAt   *time.Time
	}{
		{
			name:        "create with all fields",
			id:          1,
			idNatural:   "550e8400-e29b-41d4-a716-446655440000",
			orgName:     "Test Organization",
			kind:        OrganizationKindTeam,
			description: &description,
			createdAt:   now,
			updatedAt:   &updatedAt,
		},
		{
			name:        "create with nil optional fields",
			id:          2,
			idNatural:   "550e8400-e29b-41d4-a716-446655440001",
			orgName:     "Another Organization",
			kind:        OrganizationKindPersonal,
			description: nil,
			createdAt:   now,
			updatedAt:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOrganization(
				tt.id,
				tt.idNatural,
				tt.orgName,
				tt.kind,
				tt.description,
				tt.createdAt,
				tt.updatedAt,
			)

			if got.ID != tt.id {
				t.Errorf("NewOrganization() ID = %v, want %v", got.ID, tt.id)
			}
			if got.IDNatural != tt.idNatural {
				t.Errorf("NewOrganization() IDNatural = %v, want %v", got.IDNatural, tt.idNatural)
			}
			if got.Name != tt.orgName {
				t.Errorf("NewOrganization() Name = %v, want %v", got.Name, tt.orgName)
			}
			if got.Kind != tt.kind {
				t.Errorf("NewOrganization() Kind = %v, want %v", got.Kind, tt.kind)
			}
			if got.CreatedAt != tt.createdAt {
				t.Errorf("NewOrganization() CreatedAt = %v, want %v", got.CreatedAt, tt.createdAt)
			}
		})
	}
}

func TestOrganization_validate(t *testing.T) {
	tests := []struct {
		name    string
		org     Organization
		wantErr bool
	}{
		{
			name: "valid with all required fields",
			org: Organization{
				IDNatural: "550e8400-e29b-41d4-a716-446655440000",
				Name:      "Test Organization",
				Kind:      OrganizationKindTeam,
			},
			wantErr: false,
		},
		{
			name: "valid with name at max length",
			org: Organization{
				IDNatural: "550e8400-e29b-41d4-a716-446655440000",
				Name:      strings.Repeat("a", 100),
				Kind:      OrganizationKindTeam,
			},
			wantErr: false,
		},
		{
			name: "error when id_natural is empty",
			org: Organization{
				IDNatural: "",
				Name:      "Test Organization",
				Kind:      OrganizationKindTeam,
			},
			wantErr: true,
		},
		{
			name: "error when name is empty",
			org: Organization{
				IDNatural: "550e8400-e29b-41d4-a716-446655440000",
				Name:      "",
				Kind:      OrganizationKindTeam,
			},
			wantErr: true,
		},
		{
			name: "error when name exceeds 100 characters",
			org: Organization{
				IDNatural: "550e8400-e29b-41d4-a716-446655440000",
				Name:      strings.Repeat("a", 101),
				Kind:      OrganizationKindTeam,
			},
			wantErr: true,
		},
		{
			name: "error when kind is invalid",
			org: Organization{
				IDNatural: "550e8400-e29b-41d4-a716-446655440000",
				Name:      "Test Organization",
				Kind:      OrganizationKind("invalid"),
			},
			wantErr: true,
		},
		{
			name: "error when multiple fields are invalid",
			org: Organization{
				IDNatural: "",
				Name:      "",
				Kind:      "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.org.validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Organization.validate() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Organization.validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func newValidOrganization() Organization {
	return Organization{
		IDNatural: "550e8400-e29b-41d4-a716-446655440000",
		Name:      "Test Organization",
		Kind:      OrganizationKindTeam,
	}
}

func TestOrganization_SetName(t *testing.T) {
	tests := []struct {
		name    string
		orgName string
		wantErr bool
	}{
		{
			name:    "success with valid name",
			orgName: "New Organization",
			wantErr: false,
		},
		{
			name:    "success with name at max length",
			orgName: strings.Repeat("b", 100),
			wantErr: false,
		},
		{
			name:    "error when name is empty",
			orgName: "",
			wantErr: true,
		},
		{
			name:    "error when name exceeds 100 characters",
			orgName: strings.Repeat("c", 101),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := newValidOrganization()
			err := o.SetName(tt.orgName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Organization.SetName() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Organization.SetName() error = %v, wantErr %v", err, tt.wantErr)
				}
				if o.Name != tt.orgName {
					t.Errorf("Organization.Name = %v, want %v", o.Name, tt.orgName)
				}
			}
		})
	}
}

func TestOrganization_SetDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
	}{
		{
			name:        "success with valid description",
			description: "New description",
			wantErr:     false,
		},
		{
			name:        "success with empty description",
			description: "",
			wantErr:     false,
		},
		{
			name:        "success with long description",
			description: strings.Repeat("d", 500),
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := newValidOrganization()
			err := o.SetDescription(tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Organization.SetDescription() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Organization.SetDescription() error = %v, wantErr %v", err, tt.wantErr)
				}
				if o.Description == nil {
					t.Errorf("Organization.Description is nil")
				} else if *o.Description != tt.description {
					t.Errorf("Organization.Description = %v, want %v", *o.Description, tt.description)
				}
			}
		})
	}
}
