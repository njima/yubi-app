package model

import (
	"strings"
	"testing"
	"time"
)

func TestInitTask(t *testing.T) {
	description := "Test description"
	manualURL := "https://example.com/manual"
	deadline := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name           string
		organizationID string
		taskName       string
		description    *string
		manualURL      string
		deadline       time.Time
		wantErr        bool
	}{
		{
			name:           "success with valid inputs",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			taskName:       "Test Task",
			description:    &description,
			manualURL:      manualURL,
			deadline:       deadline,
			wantErr:        false,
		},
		{
			name:           "success with nil description",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			taskName:       "Test Task",
			description:    nil,
			manualURL:      manualURL,
			deadline:       deadline,
			wantErr:        false,
		},
		{
			name:           "error when organization_id is empty",
			organizationID: "",
			taskName:       "Test Task",
			description:    &description,
			manualURL:      manualURL,
			deadline:       deadline,
			wantErr:        true,
		},
		{
			name:           "error when name is empty",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			taskName:       "",
			description:    &description,
			manualURL:      manualURL,
			deadline:       deadline,
			wantErr:        true,
		},
		{
			name:           "error when name exceeds 100 characters",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			taskName:       strings.Repeat("a", 101),
			description:    &description,
			manualURL:      manualURL,
			deadline:       deadline,
			wantErr:        true,
		},
		{
			name:           "success when name is exactly 100 characters",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			taskName:       strings.Repeat("a", 100),
			description:    &description,
			manualURL:      manualURL,
			deadline:       deadline,
			wantErr:        false,
		},
		{
			name:           "error when manual_url is empty",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			taskName:       "Test Task",
			description:    &description,
			manualURL:      "",
			deadline:       deadline,
			wantErr:        true,
		},
		{
			name:           "error when deadline is zero",
			organizationID: "550e8400-e29b-41d4-a716-446655440001",
			taskName:       "Test Task",
			description:    &description,
			manualURL:      manualURL,
			deadline:       time.Time{},
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priority := TaskPriorityNormal
			difficulty := TaskDifficultyB
			status := TaskStatusPlanning
			got, err := InitTask(tt.organizationID, tt.taskName, tt.description, tt.manualURL, &priority, &difficulty, &status, tt.deadline, nil)

			if tt.wantErr {
				if err == nil {
					t.Errorf("InitTask() error = nil, wantErr %v", tt.wantErr)
				}
				if got.IDNatural != "" {
					t.Errorf("InitTask() IDNatural = %v, want empty", got.IDNatural)
				}
				return
			}

			if err != nil {
				t.Errorf("InitTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.IDNatural == "" {
				t.Errorf("InitTask() IDNatural is empty")
			}
			if got.OrganizationID != tt.organizationID {
				t.Errorf("InitTask() OrganizationID = %v, want %v", got.OrganizationID, tt.organizationID)
			}
			if got.Name != tt.taskName {
				t.Errorf("InitTask() Name = %v, want %v", got.Name, tt.taskName)
			}
			if got.ManualURL != tt.manualURL {
				t.Errorf("InitTask() ManualURL = %v, want %v", got.ManualURL, tt.manualURL)
			}
			if !got.Deadline.Equal(tt.deadline) {
				t.Errorf("InitTask() Deadline = %v, want %v", got.Deadline, tt.deadline)
			}
			if got.Version != InitialVersion {
				t.Errorf("InitTask() Version = %v, want %v", got.Version, InitialVersion)
			}
			if !got.IsActive {
				t.Errorf("InitTask() IsActive = %v, want true", got.IsActive)
			}
			if got.CreatedAt.IsZero() {
				t.Errorf("InitTask() CreatedAt is zero")
			}
		})
	}
}

func TestNewTask(t *testing.T) {
	description := "Test description"
	now := time.Now()
	updatedAt := now.Add(time.Hour)

	tests := []struct {
		name        string
		id          int64
		idNatural   string
		taskName    string
		description *string
		deadline    time.Time
		createdAt   time.Time
		updatedAt   *time.Time
	}{
		{
			name:        "success with all fields",
			id:          1,
			idNatural:   "550e8400-e29b-41d4-a716-446655440000",
			taskName:    "Test Task",
			description: &description,
			deadline:    now.Add(2 * time.Hour),
			createdAt:   now,
			updatedAt:   &updatedAt,
		},
		{
			name:        "success with nil description and updatedAt",
			id:          2,
			idNatural:   "550e8400-e29b-41d4-a716-446655440002",
			taskName:    "Another Task",
			description: nil,
			deadline:    now.Add(3 * time.Hour),
			createdAt:   now,
			updatedAt:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTask(tt.id, tt.idNatural, tt.taskName, tt.description, tt.deadline, tt.createdAt, tt.updatedAt)

			if got.ID != tt.id {
				t.Errorf("NewTask() ID = %v, want %v", got.ID, tt.id)
			}
			if got.IDNatural != tt.idNatural {
				t.Errorf("NewTask() IDNatural = %v, want %v", got.IDNatural, tt.idNatural)
			}
			if got.Name != tt.taskName {
				t.Errorf("NewTask() Name = %v, want %v", got.Name, tt.taskName)
			}
			if !got.Deadline.Equal(tt.deadline) {
				t.Errorf("NewTask() Deadline = %v, want %v", got.Deadline, tt.deadline)
			}
			if got.CreatedAt != tt.createdAt {
				t.Errorf("NewTask() CreatedAt = %v, want %v", got.CreatedAt, tt.createdAt)
			}
		})
	}
}

func newValidTask() Task {
	priority := TaskPriorityNormal
	difficulty := TaskDifficultyB
	status := TaskStatusPlanning
	return Task{
		IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
		OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
		Name:           "Test Task",
		ManualURL:      "https://example.com/manual",
		Priority:       &priority,
		Difficulty:     &difficulty,
		Status:         &status,
		Deadline:       time.Now().Add(time.Hour),
		Version:        InitialVersion,
		IsActive:       true,
	}
}

func TestTask_validate(t *testing.T) {
	tests := []struct {
		name    string
		task    Task
		wantErr bool
	}{
		{
			name:    "valid with all required fields",
			task:    newValidTask(),
			wantErr: false,
		},
		{
			name: "error when id_natural is empty",
			task: Task{
				IDNatural:      "",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				Name:           "Test Task",
			},
			wantErr: true,
		},
		{
			name: "error when organization_id is empty",
			task: Task{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "",
				Name:           "Test Task",
			},
			wantErr: true,
		},
		{
			name: "error when name is empty",
			task: Task{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				Name:           "",
			},
			wantErr: true,
		},
		{
			name: "error when name exceeds 100 characters",
			task: Task{
				IDNatural:      "550e8400-e29b-41d4-a716-446655440000",
				OrganizationID: "550e8400-e29b-41d4-a716-446655440001",
				Name:           strings.Repeat("a", 101),
			},
			wantErr: true,
		},
		{
			name: "valid with https manual_url",
			task: func() Task {
				t := newValidTask()
				t.ManualURL = "https://example.com/manual"
				return t
			}(),
			wantErr: false,
		},
		{
			name: "error with empty manual_url",
			task: func() Task {
				t := newValidTask()
				t.ManualURL = ""
				return t
			}(),
			wantErr: true,
		},
		{
			name: "error when manual_url uses http",
			task: func() Task {
				t := newValidTask()
				t.ManualURL = "http://example.com/manual"
				return t
			}(),
			wantErr: true,
		},
		{
			name: "error when manual_url is not a valid URL",
			task: func() Task {
				t := newValidTask()
				t.ManualURL = "not-a-url"
				return t
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Task.validate() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Task.validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestTask_SetName(t *testing.T) {
	tests := []struct {
		name     string
		taskName string
		wantErr  bool
	}{
		{
			name:     "success with valid name",
			taskName: "New Task",
			wantErr:  false,
		},
		{
			name:     "success with name at max length",
			taskName: strings.Repeat("b", 100),
			wantErr:  false,
		},
		{
			name:     "error when name is empty",
			taskName: "",
			wantErr:  true,
		},
		{
			name:     "error when name exceeds 100 characters",
			taskName: strings.Repeat("c", 101),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := newValidTask()
			err := task.SetName(tt.taskName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Task.SetName() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Task.SetName() error = %v, wantErr %v", err, tt.wantErr)
				}
				if task.Name != tt.taskName {
					t.Errorf("Task.Name = %v, want %v", task.Name, tt.taskName)
				}
			}
		})
	}
}

func TestTask_SetManualURL(t *testing.T) {
	tests := []struct {
		name      string
		manualURL string
		wantErr   bool
	}{
		{
			name:      "success with valid https URL",
			manualURL: "https://example.com/manual",
			wantErr:   false,
		},
		{
			name:      "error with empty string",
			manualURL: "",
			wantErr:   true,
		},
		{
			name:      "error with http URL",
			manualURL: "http://example.com/manual",
			wantErr:   true,
		},
		{
			name:      "error with invalid URL",
			manualURL: "not-a-url",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := newValidTask()
			err := task.SetManualURL(tt.manualURL)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Task.SetManualURL() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Task.SetManualURL() error = %v, wantErr %v", err, tt.wantErr)
				}
				if task.ManualURL != tt.manualURL {
					t.Errorf("Task.ManualURL = %v, want %v", task.ManualURL, tt.manualURL)
				}
			}
		})
	}
}

func TestTask_SetDescription(t *testing.T) {
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
			task := newValidTask()
			err := task.SetDescription(tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Task.SetDescription() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Task.SetDescription() error = %v, wantErr %v", err, tt.wantErr)
				}
				if task.Description == nil {
					t.Errorf("Task.Description is nil")
				} else if *task.Description != tt.description {
					t.Errorf("Task.Description = %v, want %v", *task.Description, tt.description)
				}
			}
		})
	}
}
