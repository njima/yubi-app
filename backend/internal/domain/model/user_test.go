package model

import (
	"strings"
	"testing"
	"time"
)

func TestInitUser(t *testing.T) {
	tests := []struct {
		name      string
		googleSub string
		userName  string
		email     string
		avatarURL string
		wantErr   bool
	}{
		{
			name:      "success with valid inputs",
			googleSub: "google-oauth2|123",
			userName:  "Test User",
			email:     "test@example.com",
			avatarURL: "https://example.com/a.png",
			wantErr:   false,
		},
		{
			name:      "success without avatar",
			googleSub: "google-oauth2|456",
			userName:  "No Avatar User",
			email:     "no-avatar@example.com",
			avatarURL: "",
			wantErr:   false,
		},
		{
			name:      "error when google_sub is empty",
			googleSub: "",
			userName:  "Test User",
			email:     "test@example.com",
			avatarURL: "",
			wantErr:   true,
		},
		{
			name:      "error when name is empty",
			googleSub: "google-oauth2|123",
			userName:  "",
			email:     "test@example.com",
			avatarURL: "",
			wantErr:   true,
		},
		{
			name:      "error when name exceeds 60 characters",
			googleSub: "google-oauth2|123",
			userName:  strings.Repeat("a", 61),
			email:     "test@example.com",
			avatarURL: "",
			wantErr:   true,
		},
		{
			name:      "success when name is exactly 60 characters",
			googleSub: "google-oauth2|123",
			userName:  strings.Repeat("a", 60),
			email:     "test@example.com",
			avatarURL: "",
			wantErr:   false,
		},
		{
			name:      "error when email is empty",
			googleSub: "google-oauth2|123",
			userName:  "Test User",
			email:     "",
			avatarURL: "",
			wantErr:   true,
		},
		{
			name:      "error when email format is invalid",
			googleSub: "google-oauth2|123",
			userName:  "Test User",
			email:     "invalid-email",
			avatarURL: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitUser(tt.googleSub, tt.userName, tt.email, tt.avatarURL)

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
			if got.GoogleSub != tt.googleSub {
				t.Errorf("InitUser() GoogleSub = %v, want %v", got.GoogleSub, tt.googleSub)
			}
			if got.Name != tt.userName {
				t.Errorf("InitUser() Name = %v, want %v", got.Name, tt.userName)
			}
			if got.Email != tt.email {
				t.Errorf("InitUser() Email = %v, want %v", got.Email, tt.email)
			}
			if tt.avatarURL == "" {
				if got.AvatarURL != nil {
					t.Errorf("InitUser() AvatarURL = %v, want nil", got.AvatarURL)
				}
			} else if got.AvatarURL == nil || *got.AvatarURL != tt.avatarURL {
				t.Errorf("InitUser() AvatarURL = %v, want %v", got.AvatarURL, tt.avatarURL)
			}
			if got.CreatedAt.IsZero() {
				t.Errorf("InitUser() CreatedAt is zero")
			}
		})
	}
}

func TestInitUserIdentity(t *testing.T) {
	user, err := InitUser("google-oauth2|123", "Ada Lovelace", "ada@example.com", "https://example.com/a.png")
	if err != nil {
		t.Fatalf("InitUser() error = %v", err)
	}
	if user.GoogleSub != "google-oauth2|123" {
		t.Fatalf("GoogleSub = %q", user.GoogleSub)
	}
	if user.Name != "Ada Lovelace" {
		t.Fatalf("Name = %q", user.Name)
	}
	if user.Email != "ada@example.com" {
		t.Fatalf("Email = %q", user.Email)
	}
	if user.AvatarURL == nil || *user.AvatarURL != "https://example.com/a.png" {
		t.Fatalf("AvatarURL = %v", user.AvatarURL)
	}
}

func TestInitUserIdentityRequiresGoogleSub(t *testing.T) {
	_, err := InitUser("", "Ada Lovelace", "ada@example.com", "")
	if err == nil {
		t.Fatal("InitUser() expected error for empty google_sub")
	}
}

func TestNewUser(t *testing.T) {
	now := time.Now()
	updatedAt := now.Add(time.Hour)
	avatarURL := "https://example.com/a.png"

	tests := []struct {
		name      string
		id        int64
		idNatural string
		googleSub string
		userName  string
		email     string
		avatarURL *string
		createdAt time.Time
		updatedAt *time.Time
	}{
		{
			name:      "success with all fields",
			id:        1,
			idNatural: "550e8400-e29b-41d4-a716-446655440000",
			googleSub: "google-oauth2|123",
			userName:  "Test User",
			email:     "test@example.com",
			avatarURL: &avatarURL,
			createdAt: now,
			updatedAt: &updatedAt,
		},
		{
			name:      "success with nil optional fields",
			id:        2,
			idNatural: "550e8400-e29b-41d4-a716-446655440002",
			googleSub: "google-oauth2|456",
			userName:  "Another User",
			email:     "another@example.com",
			avatarURL: nil,
			createdAt: now,
			updatedAt: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUser(tt.id, tt.idNatural, tt.googleSub, tt.userName, tt.email, tt.avatarURL, tt.createdAt, tt.updatedAt)

			if got.ID != tt.id {
				t.Errorf("NewUser() ID = %v, want %v", got.ID, tt.id)
			}
			if got.IDNatural != tt.idNatural {
				t.Errorf("NewUser() IDNatural = %v, want %v", got.IDNatural, tt.idNatural)
			}
			if got.GoogleSub != tt.googleSub {
				t.Errorf("NewUser() GoogleSub = %v, want %v", got.GoogleSub, tt.googleSub)
			}
			if got.Name != tt.userName {
				t.Errorf("NewUser() Name = %v, want %v", got.Name, tt.userName)
			}
			if got.Email != tt.email {
				t.Errorf("NewUser() Email = %v, want %v", got.Email, tt.email)
			}
			if got.AvatarURL != tt.avatarURL {
				t.Errorf("NewUser() AvatarURL = %v, want %v", got.AvatarURL, tt.avatarURL)
			}
			if got.CreatedAt != tt.createdAt {
				t.Errorf("NewUser() CreatedAt = %v, want %v", got.CreatedAt, tt.createdAt)
			}
		})
	}
}

func newValidUser() User {
	return User{
		IDNatural: "550e8400-e29b-41d4-a716-446655440000",
		GoogleSub: "google-oauth2|123",
		Name:      "Test User",
		Email:     "test@example.com",
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
				IDNatural: "",
				GoogleSub: "google-oauth2|123",
				Name:      "Test User",
				Email:     "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "error when google_sub is empty",
			user: User{
				IDNatural: "550e8400-e29b-41d4-a716-446655440000",
				GoogleSub: "",
				Name:      "Test User",
				Email:     "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "error when name is empty",
			user: User{
				IDNatural: "550e8400-e29b-41d4-a716-446655440000",
				GoogleSub: "google-oauth2|123",
				Name:      "",
				Email:     "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "error when name exceeds 60 characters",
			user: User{
				IDNatural: "550e8400-e29b-41d4-a716-446655440000",
				GoogleSub: "google-oauth2|123",
				Name:      strings.Repeat("a", 61),
				Email:     "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "error when email is empty",
			user: User{
				IDNatural: "550e8400-e29b-41d4-a716-446655440000",
				GoogleSub: "google-oauth2|123",
				Name:      "Test User",
				Email:     "",
			},
			wantErr: true,
		},
		{
			name: "error when email format is invalid",
			user: User{
				IDNatural: "550e8400-e29b-41d4-a716-446655440000",
				GoogleSub: "google-oauth2|123",
				Name:      "Test User",
				Email:     "invalid-email",
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
