package repository

import (
	"context"
	"time"
)

// Export hard cap. Aligned with /episodes/export.
const MaxOperatorYieldExportRows = 30_000

// Report is JST-fixed regardless of server clock. Shared between
// gateway (date bucketing) and usecase (CSV formatting) so the two cannot
// drift.
var JSTLocation = time.FixedZone("Asia/Tokyo", 9*60*60)

// DateFrom/DateTo are calendar dates: only Y/M/D is read, and they're always
// interpreted as JST regardless of the time.Time's Location. The repository
// reconstructs the JST instants for SQL.
type OperatorYieldExportFilter struct {
	DateFrom   time.Time
	DateTo     time.Time
	LocationID *string
	TaskID     *string
	UserID     *string
}

// Raw aggregated row. Display formatting (working window string, percent
// rounding, etc.) belongs to the usecase layer.
type OperatorYieldExportRow struct {
	WorkDate         time.Time
	OperatorUserID   string
	OperatorName     string
	TaskID           string
	TaskName         string
	FirstStart       time.Time
	LastEnd          time.Time
	WorkingSeconds   int64
	CollectedSeconds int64
	DiscardedSeconds int64
	EpisodeCount     int64 // bigint from PG COUNT(*)
}

type OperatorYield interface {
	Export(ctx context.Context, conn Conn, filter OperatorYieldExportFilter) ([]OperatorYieldExportRow, error)
}
