package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database"
	"github.com/airoa-org/yubi-app/backend/internal/infra/persistence"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/uptrace/bun"
)

type Config struct {
	Period    string // hourly, daily, monthly
	Backfill  bool   // enable backfill mode
	From      string // backfill start date (YYYY-MM-DD or YYYY-MM-DDTHH:00:00Z)
	To        string // backfill end date
	DBDriver  string
	DBUser    string
	DBPass    string
	DBHost    string
	DBPort    string
	DBName    string
	DBSSLMode string
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := parseFlags()

	if err := run(ctx, cfg); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func parseFlags() Config {
	cfg := Config{}

	flag.StringVar(&cfg.Period, "period", "hourly", "Aggregation period: hourly, daily, or monthly")
	flag.BoolVar(&cfg.Backfill, "backfill", false, "Enable backfill mode to aggregate historical data")
	flag.StringVar(&cfg.From, "from", "", "Backfill start date (YYYY-MM-DD or YYYY-MM-DDTHH:00:00Z for hourly)")
	flag.StringVar(&cfg.To, "to", "", "Backfill end date (YYYY-MM-DD or YYYY-MM-DDTHH:00:00Z for hourly)")

	// Database config from environment with fallbacks
	flag.StringVar(&cfg.DBDriver, "db-driver", getEnv("DB_DRIVER", "postgres"), "Database driver")
	flag.StringVar(&cfg.DBUser, "db-user", getEnv("DB_USER", ""), "Database user")
	flag.StringVar(&cfg.DBPass, "db-password", getEnv("DB_PASSWORD", ""), "Database password")
	flag.StringVar(&cfg.DBHost, "db-host", getEnv("DB_HOST", "localhost"), "Database host")
	flag.StringVar(&cfg.DBPort, "db-port", getEnv("DB_PORT", "5432"), "Database port")
	flag.StringVar(&cfg.DBName, "db-name", getEnv("DB_NAME", ""), "Database name")
	flag.StringVar(&cfg.DBSSLMode, "db-sslmode", getEnv("DB_SSL_MODE", "disable"), "Database SSL mode")

	flag.Parse()

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func run(ctx context.Context, cfg Config) error {
	// Validate period
	period := model.AggregationPeriod(cfg.Period)
	switch period {
	case model.AggregationPeriodHourly, model.AggregationPeriodDaily, model.AggregationPeriodMonthly:
		// valid
	default:
		return fmt.Errorf("invalid period: %s (must be hourly, daily, or monthly)", cfg.Period)
	}

	// Connect to database
	db, err := database.NewDatabase(
		cfg.DBDriver,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	repo := persistence.NewEpisodeStats()

	if cfg.Backfill {
		return runBackfill(ctx, db, repo, period, cfg.From, cfg.To)
	}

	return runRegular(ctx, db, repo, period)
}

// runRegular runs aggregation for the previous period
func runRegular(ctx context.Context, db *bun.DB, repo repository.EpisodeStats, period model.AggregationPeriod) error {
	now := time.Now().UTC()
	var from, to time.Time

	switch period {
	case model.AggregationPeriodHourly:
		// Aggregate the previous hour
		to = now.Truncate(time.Hour)
		from = to.Add(-time.Hour)
	case model.AggregationPeriodDaily:
		// Aggregate the previous day
		to = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		from = to.AddDate(0, 0, -1)
	case model.AggregationPeriodMonthly:
		// Aggregate the previous month
		to = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		from = to.AddDate(0, -1, 0)
	}

	fmt.Printf("Running %s aggregation for period: %s to %s\n", period, from.Format(time.RFC3339), to.Format(time.RFC3339))

	if err := aggregatePeriod(ctx, db, repo, period, from, to); err != nil {
		return err
	}

	// Aggregate task version stats (all-time cumulative, runs on every period)
	if err := aggregateTaskVersionStats(ctx, db, repo); err != nil {
		return fmt.Errorf("failed to aggregate task version stats: %w", err)
	}

	return nil
}

// runBackfill runs aggregation for a historical period
func runBackfill(ctx context.Context, db *bun.DB, repo repository.EpisodeStats, period model.AggregationPeriod, fromStr, toStr string) error {
	if fromStr == "" || toStr == "" {
		return fmt.Errorf("--from and --to are required for backfill mode")
	}

	var from, to time.Time
	var err error

	// Parse dates based on period
	switch period {
	case model.AggregationPeriodHourly:
		from, err = parseDateTime(fromStr)
		if err != nil {
			return fmt.Errorf("invalid --from date: %w", err)
		}
		to, err = parseDateTime(toStr)
		if err != nil {
			return fmt.Errorf("invalid --to date: %w", err)
		}
		from = from.Truncate(time.Hour)
		to = to.Truncate(time.Hour)
	case model.AggregationPeriodDaily:
		from, err = parseDate(fromStr)
		if err != nil {
			return fmt.Errorf("invalid --from date: %w", err)
		}
		to, err = parseDate(toStr)
		if err != nil {
			return fmt.Errorf("invalid --to date: %w", err)
		}
	case model.AggregationPeriodMonthly:
		from, err = parseMonth(fromStr)
		if err != nil {
			return fmt.Errorf("invalid --from date: %w", err)
		}
		to, err = parseMonth(toStr)
		if err != nil {
			return fmt.Errorf("invalid --to date: %w", err)
		}
	}

	if from.After(to) {
		return fmt.Errorf("--from must be before --to")
	}

	fmt.Printf("Running %s backfill from %s to %s\n", period, from.Format(time.RFC3339), to.Format(time.RFC3339))

	// Iterate through each period and aggregate
	current := from
	for current.Before(to) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var periodEnd time.Time
		switch period {
		case model.AggregationPeriodHourly:
			periodEnd = current.Add(time.Hour)
		case model.AggregationPeriodDaily:
			periodEnd = current.AddDate(0, 0, 1)
		case model.AggregationPeriodMonthly:
			periodEnd = current.AddDate(0, 1, 0)
		}

		fmt.Printf("  Aggregating: %s\n", current.Format(time.RFC3339))
		if err := aggregatePeriod(ctx, db, repo, period, current, periodEnd); err != nil {
			return fmt.Errorf("failed to aggregate period %s: %w", current.Format(time.RFC3339), err)
		}

		current = periodEnd
	}

	// Aggregate task version stats once after backfill (all-time cumulative)
	if err := aggregateTaskVersionStats(ctx, db, repo); err != nil {
		return fmt.Errorf("failed to aggregate task version stats: %w", err)
	}

	fmt.Println("Backfill completed successfully")
	return nil
}

// aggregatePeriod aggregates episode data for a specific period
func aggregatePeriod(ctx context.Context, db *bun.DB, repo repository.EpisodeStats, period model.AggregationPeriod, from, to time.Time) error {
	// Get aggregated data from episodes
	aggregatedData, err := repo.AggregateEpisodesForPeriod(ctx, db, from, to)
	if err != nil {
		return fmt.Errorf("failed to aggregate episodes: %w", err)
	}

	// Convert to stats models
	statsList := make([]model.EpisodeStats, len(aggregatedData))
	for i, data := range aggregatedData {
		statsList[i] = model.EpisodeStats{
			OrganizationID:       data.OrganizationID,
			LocationID:           data.LocationID,
			RobotID:              data.RobotID,
			PeriodStart:          from,
			TotalDurationSeconds: data.TotalDurationSeconds,
			EpisodeCount:         data.EpisodeCount,
		}
	}

	// Replace stats for this period within a transaction (delete + insert atomically)
	if err := db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		switch period {
		case model.AggregationPeriodHourly:
			return repo.BulkReplaceHourly(ctx, tx, from, statsList)
		case model.AggregationPeriodDaily:
			return repo.BulkReplaceDaily(ctx, tx, from, statsList)
		case model.AggregationPeriodMonthly:
			return repo.BulkReplaceMonthly(ctx, tx, from, statsList)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to replace %s stats: %w", period, err)
	}

	fmt.Printf("  Aggregated %d records\n", len(statsList))

	return nil
}

// aggregateTaskVersionStats aggregates all-time episode data per task version.
// The delete + upsert is wrapped in a transaction to avoid partial updates.
func aggregateTaskVersionStats(ctx context.Context, db *bun.DB, repo repository.EpisodeStats) error {
	aggregatedData, err := repo.AggregateByTaskVersion(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to aggregate by task version: %w", err)
	}

	statsList := make([]model.TaskVersionStats, len(aggregatedData))
	for i, data := range aggregatedData {
		statsList[i] = model.TaskVersionStats{
			TaskVersionID:        data.TaskVersionID,
			TotalDurationSeconds: data.TotalDurationSeconds,
			EpisodeCount:         data.EpisodeCount,
		}
	}

	if err := db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return repo.BulkUpsertTaskVersionStats(ctx, tx, statsList)
	}); err != nil {
		return fmt.Errorf("failed to upsert task version stats: %w", err)
	}

	fmt.Printf("  Aggregated %d task version stats\n", len(statsList))
	return nil
}

// Date parsing helpers

func parseDateTime(s string) (time.Time, error) {
	// Try RFC3339 first
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC(), nil
	}
	// Try date format
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.UTC(), nil
	}
	// Try datetime without timezone
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("unable to parse datetime: %s", s)
}

func parseDate(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("expected format YYYY-MM-DD: %w", err)
	}
	return t.UTC(), nil
}

func parseMonth(s string) (time.Time, error) {
	// Try YYYY-MM format
	if t, err := time.Parse("2006-01", s); err == nil {
		return t.UTC(), nil
	}
	// Try YYYY-MM-DD format (use first day of month)
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
	}
	return time.Time{}, fmt.Errorf("expected format YYYY-MM or YYYY-MM-DD: %s", s)
}
