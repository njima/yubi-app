package controller

import (
	"os"
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"gopkg.in/yaml.v3"
)

var robotStatusMappingCases = []struct {
	name string
	api  openapi.RobotStatus
	want model.RobotStatus
}{
	{name: "online", api: openapi.RobotStatusOnline, want: model.RobotStatusOnline},
	{name: "busy", api: openapi.RobotStatusBusy, want: model.RobotStatusBusy},
	{name: "offline", api: openapi.RobotStatusOffline, want: model.RobotStatusOffline},
	{name: "faulted", api: openapi.RobotStatusFaulted, want: model.RobotStatusFaulted},
	{name: "maintenance", api: openapi.RobotStatusMaintenance, want: model.RobotStatusMaintenance},
	{name: "ready", api: openapi.RobotStatusReady, want: model.RobotStatusReady},
}

var leaderStatusMappingCases = []struct {
	name string
	api  openapi.LeaderStatus
	want model.LeaderStatus
}{
	{name: "ready", api: openapi.LeaderReady, want: model.LeaderStatusReady},
	{name: "faulted", api: openapi.LeaderFaulted, want: model.LeaderStatusFaulted},
	{name: "maintenance", api: openapi.LeaderMaintenance, want: model.LeaderStatusMaintenance},
}

var episodeStatusMappingCases = []struct {
	name string
	api  openapi.EpisodeCollectionStatus
	want model.EpisodeStatus
}{
	{name: "ready", api: openapi.EpisodeCollectionStatusReady, want: model.EpisodeStatusReady},
	{name: "recording", api: openapi.EpisodeCollectionStatusRecording, want: model.EpisodeStatusRecording},
	{name: "cancel", api: openapi.EpisodeCollectionStatusCancel, want: model.EpisodeStatusCancel},
	{name: "completed", api: openapi.EpisodeCollectionStatusCompleted, want: model.EpisodeStatusCompleted},
}

func TestRobotStatusModelMapping(t *testing.T) {
	for _, tt := range robotStatusMappingCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := robotStatusModel(&tt.api)
			if err != nil {
				t.Fatalf("robotStatusModel() error = %v", err)
			}
			if got == nil {
				t.Fatalf("robotStatusModel() = nil")
			}
			if *got != tt.want {
				t.Fatalf("robotStatusModel() = %v, want %v", *got, tt.want)
			}
			if roundTrip := openAPIRobotStatus(*got); roundTrip != tt.api {
				t.Fatalf("openAPIRobotStatus(robotStatusModel()) = %v, want %v", roundTrip, tt.api)
			}
		})
	}
}

func TestRobotStatusMappingCoversOpenAPISchemaEnum(t *testing.T) {
	want := openAPISchemaEnumCount(t, "RobotStatus")

	if got := len(robotStatusMappingCases); got != want {
		t.Fatalf("robot status mapping count = %d, want OpenAPI enum count %d", got, want)
	}
}

func TestRobotStatusModelRejectsUnknownValue(t *testing.T) {
	value := openapi.RobotStatus(999)

	got, err := robotStatusModel(&value)
	if err == nil {
		t.Fatalf("robotStatusModel() error = nil")
	}
	if got != nil {
		t.Fatalf("robotStatusModel() = %v, want nil", *got)
	}
}

func TestLeaderStatusMapping(t *testing.T) {
	for _, tt := range leaderStatusMappingCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := leaderStatus(&tt.api)
			if err != nil {
				t.Fatalf("leaderStatus() error = %v", err)
			}
			if got == nil {
				t.Fatalf("leaderStatus() = nil")
			}
			if *got != tt.want {
				t.Fatalf("leaderStatus() = %v, want %v", *got, tt.want)
			}
			roundTrip := openAPILeaderStatus(got)
			if roundTrip == nil {
				t.Fatalf("openAPILeaderStatus() = nil")
			}
			if *roundTrip != tt.api {
				t.Fatalf("openAPILeaderStatus(leaderStatus()) = %v, want %v", *roundTrip, tt.api)
			}
		})
	}
}

func TestLeaderStatusMappingCoversOpenAPISchemaEnum(t *testing.T) {
	want := openAPISchemaEnumCount(t, "LeaderStatus")

	if got := len(leaderStatusMappingCases); got != want {
		t.Fatalf("leader status mapping count = %d, want OpenAPI enum count %d", got, want)
	}
}

func TestLeaderStatusRejectsUnknownValue(t *testing.T) {
	value := openapi.LeaderStatus(999)

	got, err := leaderStatus(&value)
	if err == nil {
		t.Fatalf("leaderStatus() error = nil")
	}
	if got != nil {
		t.Fatalf("leaderStatus() = %v, want nil", *got)
	}
}

func TestEpisodeStatusModelMapping(t *testing.T) {
	for _, tt := range episodeStatusMappingCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := episodeStatusModel(&tt.api)
			if err != nil {
				t.Fatalf("episodeStatusModel() error = %v", err)
			}
			if got == nil {
				t.Fatalf("episodeStatusModel() = nil")
			}
			if *got != tt.want {
				t.Fatalf("episodeStatusModel() = %v, want %v", *got, tt.want)
			}
			if roundTrip := openAPIEpisodeStatus(*got); roundTrip != tt.api {
				t.Fatalf("openAPIEpisodeStatus(episodeStatusModel()) = %v, want %v", roundTrip, tt.api)
			}
		})
	}
}

func TestEpisodeStatusMappingCoversOpenAPISchemaEnum(t *testing.T) {
	want := openAPISchemaEnumCount(t, "EpisodeCollectionStatus")

	if got := len(episodeStatusMappingCases); got != want {
		t.Fatalf("episode status mapping count = %d, want OpenAPI enum count %d", got, want)
	}
}

func TestEpisodeStatusRejectsUnknownValue(t *testing.T) {
	value := openapi.EpisodeCollectionStatus(999)

	got, err := episodeStatusModel(&value)
	if err == nil {
		t.Fatalf("episodeStatusModel() error = nil")
	}
	if got != nil {
		t.Fatalf("episodeStatusModel() = %v, want nil", *got)
	}
}

func openAPISchemaEnumCount(t *testing.T, schemaName string) int {
	t.Helper()

	content, err := os.ReadFile("../../../../../openapi/openapi.yaml")
	if err != nil {
		t.Fatalf("failed to read OpenAPI schema: %v", err)
	}

	var doc struct {
		Components struct {
			Schemas map[string]struct {
				Enum []int `yaml:"enum"`
			} `yaml:"schemas"`
		} `yaml:"components"`
	}
	if err := yaml.Unmarshal(content, &doc); err != nil {
		t.Fatalf("failed to parse OpenAPI schema: %v", err)
	}

	schema, ok := doc.Components.Schemas[schemaName]
	if !ok {
		t.Fatalf("schema %q not found in OpenAPI schema", schemaName)
	}

	return len(schema.Enum)
}
