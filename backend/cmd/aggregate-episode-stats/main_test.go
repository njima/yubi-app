package main

import (
	"context"
	"testing"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/database"
	"github.com/airoa-org/yubi-app/backend/internal/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/persistence"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type seededData struct {
	OrganizationID string
	SiteID         string
	LocationID1    string
	RobotID1       string
	LocationID2    string
	RobotID2       string
	TaskVersionID  string
	UserID         string
	TaskID         string

	Expected map[string]expectedStat
}

type expectedStat struct {
	Duration int64
	Count    int
}

func TestAggregatePeriod_CreatesRecordsAndIsIdempotent(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	defer db.Close()

	repo := persistence.NewEpisodeStats()

	tests := []struct {
		name   string
		period model.AggregationPeriod
		from   time.Time
		to     time.Time
	}{
		{
			name:   "hourly",
			period: model.AggregationPeriodHourly,
			from:   time.Date(2026, 3, 10, 15, 0, 0, 0, time.UTC),
			to:     time.Date(2026, 3, 10, 16, 0, 0, 0, time.UTC),
		},
		{
			name:   "daily",
			period: model.AggregationPeriodDaily,
			from:   time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC),
			to:     time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			name:   "monthly",
			period: model.AggregationPeriodMonthly,
			from:   time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
			to:     time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seeded := seedTestData(t, ctx, db, tt.from, tt.to)
			defer cleanupTestData(t, ctx, db, seeded)

			if err := aggregatePeriod(ctx, db, repo, tt.period, tt.from, tt.to); err != nil {
				t.Fatalf("first aggregatePeriod failed: %v", err)
			}

			first, err := listStatsByPeriod(ctx, repo, db, tt.period, seeded.OrganizationID, tt.from, tt.to)
			if err != nil {
				t.Fatalf("failed to list stats after first run: %v", err)
			}
			if len(first) < 2 {
				t.Fatalf("record count mismatch after first run: got=%d want>=2", len(first))
			}
			firstByKey := map[string]*model.EpisodeStats{}
			for _, s := range first {
				firstByKey[s.LocationID+"|"+s.RobotID] = s
			}
			for key, want := range seeded.Expected {
				got, ok := firstByKey[key]
				if !ok {
					t.Fatalf("expected stat not found after first run: %s", key)
				}
				if got.TotalDurationSeconds != want.Duration {
					t.Fatalf("duration mismatch after first run for %s: got=%d want=%d", key, got.TotalDurationSeconds, want.Duration)
				}
				if got.EpisodeCount != want.Count {
					t.Fatalf("episode_count mismatch after first run for %s: got=%d want=%d", key, got.EpisodeCount, want.Count)
				}
			}
			if err := aggregatePeriod(ctx, db, repo, tt.period, tt.from, tt.to); err != nil {
				t.Fatalf("second aggregatePeriod failed: %v", err)
			}

			second, err := listStatsByPeriod(ctx, repo, db, tt.period, seeded.OrganizationID, tt.from, tt.to)
			if err != nil {
				t.Fatalf("failed to list stats after second run: %v", err)
			}
			if len(second) != len(first) {
				t.Fatalf("idempotency broken: record count changed after second run: got=%d want=%d", len(second), len(first))
			}
			for _, row := range second {
				key := row.LocationID + "|" + row.RobotID
				want, ok := seeded.Expected[key]
				if !ok {
					t.Fatalf("unexpected stat row found after second run: %s", key)
				}
				if row.TotalDurationSeconds != want.Duration {
					t.Fatalf("idempotency broken: duration changed after second run for %s: got=%d want=%d", key, row.TotalDurationSeconds, want.Duration)
				}
				if row.EpisodeCount != want.Count {
					t.Fatalf("idempotency broken: episode_count changed after second run for %s: got=%d want=%d", key, row.EpisodeCount, want.Count)
				}
			}
		})
	}
}

func openTestDB(t *testing.T) *bun.DB {
	t.Helper()

	db, err := database.NewDatabase(
		getEnv("DB_DRIVER", "postgres"),
		getEnv("DB_USER", ""),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", ""),
		getEnv("DB_SSL_MODE", "disable"),
	)
	if err != nil {
		t.Skipf("skipping integration test because database is unavailable: %v", err)
	}

	return db
}

func seedTestData(t *testing.T, ctx context.Context, db *bun.DB, from, to time.Time) seededData {
	t.Helper()

	suffix := uuid.NewString()
	orgID := uuid.NewString()
	siteID := uuid.NewString()
	locID1 := uuid.NewString()
	robotID1 := uuid.NewString()
	locID2 := uuid.NewString()
	robotID2 := uuid.NewString()
	userID := uuid.NewString()
	taskID := uuid.NewString()
	taskVersionID := uuid.NewString()

	org := &entity.Organization{IDNatural: orgID, Name: "test-org-" + suffix}
	if _, err := db.NewInsert().Model(org).Exec(ctx); err != nil {
		t.Fatalf("failed to insert organization: %v", err)
	}

	site := &entity.Site{IDNatural: siteID, OrganizationID: orgID, Name: "test-site-" + suffix}
	if _, err := db.NewInsert().Model(site).Exec(ctx); err != nil {
		t.Fatalf("failed to insert site: %v", err)
	}

	loc1 := &entity.Location{IDNatural: locID1, OrganizationID: orgID, SiteID: siteID, Name: "test-loc-1-" + suffix}
	if _, err := db.NewInsert().Model(loc1).Exec(ctx); err != nil {
		t.Fatalf("failed to insert location 1: %v", err)
	}
	loc2 := &entity.Location{IDNatural: locID2, OrganizationID: orgID, SiteID: siteID, Name: "test-loc-2-" + suffix}
	if _, err := db.NewInsert().Model(loc2).Exec(ctx); err != nil {
		t.Fatalf("failed to insert location 2: %v", err)
	}

	robot1 := &entity.Robot{IDNatural: robotID1, OrganizationID: orgID, LocationID: locID1, Name: "test-robot-1-" + suffix}
	if _, err := db.NewInsert().Model(robot1).Exec(ctx); err != nil {
		t.Fatalf("failed to insert robot 1: %v", err)
	}
	robot2 := &entity.Robot{IDNatural: robotID2, OrganizationID: orgID, LocationID: locID2, Name: "test-robot-2-" + suffix}
	if _, err := db.NewInsert().Model(robot2).Exec(ctx); err != nil {
		t.Fatalf("failed to insert robot 2: %v", err)
	}

	user := &entity.User{
		IDNatural:      userID,
		OrganizationID: orgID,
		Name:           "test-user-" + suffix,
		Email:          "test-" + suffix + "@example.com",
	}
	if _, err := db.NewInsert().Model(user).Exec(ctx); err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	task := &entity.Task{IDNatural: taskID, OrganizationID: orgID, Name: "test-task-" + suffix}
	if _, err := db.NewInsert().Model(task).Exec(ctx); err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}

	taskVersion := &entity.TaskVersion{
		IDNatural:      taskVersionID,
		OrganizationID: orgID,
		TaskID:         taskID,
		Version:        "v1",
	}
	if _, err := db.NewInsert().Model(taskVersion).Exec(ctx); err != nil {
		t.Fatalf("failed to insert task_version: %v", err)
	}

	d1Start := from.Add(10 * time.Minute)
	d1End := d1Start.Add(30 * time.Minute)
	d2Start := from.Add(20 * time.Minute)
	d2End := d2Start.Add(20 * time.Minute)
	dCrossFromStart := from.Add(-10 * time.Minute)
	dCrossFromEnd := from.Add(5 * time.Minute)
	dCrossToStart := to.Add(-5 * time.Minute)
	dCrossToEnd := to.Add(10 * time.Minute)
	d3Start := to.Add(10 * time.Minute)
	d3End := d3Start.Add(15 * time.Minute)

	if _, err := db.NewInsert().Model(&entity.Episode{
		IDNatural:        uuid.NewString(),
		OrganizationID:   orgID,
		TaskVersionID:    taskVersionID,
		LocationID:       locID1,
		RobotID:          robotID1,
		UserID:           userID,
		CollectionStatus: model.EpisodeStatusCompleted,
		StartedAt:        &d1Start,
		FinishedAt:       &d1End,
	}).Exec(ctx); err != nil {
		t.Fatalf("failed to insert episode 1-1: %v", err)
	}

	if _, err := db.NewInsert().Model(&entity.Episode{
		IDNatural:        uuid.NewString(),
		OrganizationID:   orgID,
		TaskVersionID:    taskVersionID,
		LocationID:       locID2,
		RobotID:          robotID2,
		UserID:           userID,
		CollectionStatus: model.EpisodeStatusCompleted,
		StartedAt:        &d2Start,
		FinishedAt:       &d2End,
	}).Exec(ctx); err != nil {
		t.Fatalf("failed to insert episode 1-2: %v", err)
	}

	if _, err := db.NewInsert().Model(&entity.Episode{
		IDNatural:        uuid.NewString(),
		OrganizationID:   orgID,
		TaskVersionID:    taskVersionID,
		LocationID:       locID1,
		RobotID:          robotID1,
		UserID:           userID,
		CollectionStatus: model.EpisodeStatusCompleted,
		StartedAt:        &dCrossFromStart,
		FinishedAt:       &dCrossFromEnd,
	}).Exec(ctx); err != nil {
		t.Fatalf("failed to insert period-start crossing episode: %v", err)
	}

	d4Start := from.Add(35 * time.Minute)
	d4End := d4Start.Add(10 * time.Minute)
	if _, err := db.NewInsert().Model(&entity.Episode{
		IDNatural:        uuid.NewString(),
		OrganizationID:   orgID,
		TaskVersionID:    taskVersionID,
		LocationID:       locID2,
		RobotID:          robotID2,
		UserID:           userID,
		CollectionStatus: model.EpisodeStatusCompleted,
		StartedAt:        &d4Start,
		FinishedAt:       &d4End,
	}).Exec(ctx); err != nil {
		t.Fatalf("failed to insert episode 2-2: %v", err)
	}

	if _, err := db.NewInsert().Model(&entity.Episode{
		IDNatural:        uuid.NewString(),
		OrganizationID:   orgID,
		TaskVersionID:    taskVersionID,
		LocationID:       locID1,
		RobotID:          robotID1,
		UserID:           userID,
		CollectionStatus: model.EpisodeStatusCompleted,
		StartedAt:        &d3Start,
		FinishedAt:       &d3End,
	}).Exec(ctx); err != nil {
		t.Fatalf("failed to insert out-of-range episode: %v", err)
	}

	if _, err := db.NewInsert().Model(&entity.Episode{
		IDNatural:        uuid.NewString(),
		OrganizationID:   orgID,
		TaskVersionID:    taskVersionID,
		LocationID:       locID2,
		RobotID:          robotID2,
		UserID:           userID,
		CollectionStatus: model.EpisodeStatusCompleted,
		StartedAt:        &dCrossToStart,
		FinishedAt:       &dCrossToEnd,
	}).Exec(ctx); err != nil {
		t.Fatalf("failed to insert period-end crossing episode: %v", err)
	}

	return seededData{
		OrganizationID: orgID,
		SiteID:         siteID,
		LocationID1:    locID1,
		RobotID1:       robotID1,
		LocationID2:    locID2,
		RobotID2:       robotID2,
		TaskVersionID:  taskVersionID,
		UserID:         userID,
		TaskID:         taskID,
		Expected: map[string]expectedStat{
			locID1 + "|" + robotID1: {
				Duration: (d1End.Unix() - d1Start.Unix()) + (dCrossFromEnd.Unix() - from.Unix()),
				Count:    2,
			},
			locID2 + "|" + robotID2: {
				Duration: (d2End.Unix() - d2Start.Unix()) + (d4End.Unix() - d4Start.Unix()) + (to.Unix() - dCrossToStart.Unix()),
				Count:    3,
			},
		},
	}
}

func cleanupTestData(t *testing.T, ctx context.Context, db *bun.DB, data seededData) {
	t.Helper()

	if _, err := db.NewDelete().Model((*entity.EpisodeStatsHourly)(nil)).Where("organization_id = ?", data.OrganizationID).Exec(ctx); err != nil {
		t.Fatalf("failed to delete hourly stats: %v", err)
	}
	if _, err := db.NewDelete().Model((*entity.EpisodeStatsDaily)(nil)).Where("organization_id = ?", data.OrganizationID).Exec(ctx); err != nil {
		t.Fatalf("failed to delete daily stats: %v", err)
	}
	if _, err := db.NewDelete().Model((*entity.EpisodeStatsMonthly)(nil)).Where("organization_id = ?", data.OrganizationID).Exec(ctx); err != nil {
		t.Fatalf("failed to delete monthly stats: %v", err)
	}
	if _, err := db.NewDelete().Model((*entity.Episode)(nil)).Where("organization_id = ?", data.OrganizationID).Exec(ctx); err != nil {
		t.Fatalf("failed to delete episodes: %v", err)
	}
	if _, err := db.NewDelete().Model((*entity.Robot)(nil)).Where("id_natural IN (?, ?)", data.RobotID1, data.RobotID2).Exec(ctx); err != nil {
		t.Fatalf("failed to delete robots: %v", err)
	}
	if _, err := db.NewDelete().Model((*entity.Location)(nil)).Where("id_natural IN (?, ?)", data.LocationID1, data.LocationID2).Exec(ctx); err != nil {
		t.Fatalf("failed to delete locations: %v", err)
	}
	if _, err := db.NewDelete().Model((*entity.Site)(nil)).Where("id_natural = ?", data.SiteID).Exec(ctx); err != nil {
		t.Fatalf("failed to delete site: %v", err)
	}
	if _, err := db.NewDelete().Model((*entity.User)(nil)).Where("id_natural = ?", data.UserID).Exec(ctx); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}
	if _, err := db.NewDelete().Model((*entity.TaskVersion)(nil)).Where("id_natural = ?", data.TaskVersionID).Exec(ctx); err != nil {
		t.Fatalf("failed to delete task_version: %v", err)
	}
	if _, err := db.NewDelete().Model((*entity.Task)(nil)).Where("id_natural = ?", data.TaskID).Exec(ctx); err != nil {
		t.Fatalf("failed to delete task: %v", err)
	}
	if _, err := db.NewDelete().Model((*entity.Organization)(nil)).Where("id_natural = ?", data.OrganizationID).Exec(ctx); err != nil {
		t.Fatalf("failed to delete organization: %v", err)
	}
}

func listStatsByPeriod(
	ctx context.Context,
	repo repository.EpisodeStats,
	db *bun.DB,
	period model.AggregationPeriod,
	organizationID string,
	from, to time.Time,
) (model.EpisodeStatsList, error) {
	filter := repository.EpisodeStatsFilter{
		OrganizationID: &organizationID,
		From:           &from,
		To:             &to,
	}

	switch period {
	case model.AggregationPeriodHourly:
		return repo.ListHourly(ctx, db, filter)
	case model.AggregationPeriodDaily:
		return repo.ListDaily(ctx, db, filter)
	case model.AggregationPeriodMonthly:
		return repo.ListMonthly(ctx, db, filter)
	default:
		return nil, nil
	}
}
