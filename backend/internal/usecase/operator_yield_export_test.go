package usecase

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/rs/zerolog"
)

func newTestUsecase(repo repository.OperatorYield) *operatorYieldExport {
	return NewOperatorYieldExport(repo, repository.DataAccess{}, zerolog.Nop())
}

type stubOperatorYieldRepo struct {
	rows []repository.OperatorYieldExportRow
	err  error
}

func (s *stubOperatorYieldRepo) Export(_ context.Context, _ repository.Conn, _ repository.OperatorYieldExportFilter) ([]repository.OperatorYieldExportRow, error) {
	return s.rows, s.err
}

func mustParseUTC(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("parse %q: %v", s, err)
	}
	return v.UTC()
}

func TestOperatorYieldExport_FormatsRows(t *testing.T) {
	// Sample row: 120m work / 65m collected / 45m discarded
	// → 10m non-working, 53 EP, 54% yield.
	rows := []repository.OperatorYieldExportRow{
		{
			WorkDate:         time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
			OperatorUserID:   "user-001",
			OperatorName:     "John Doe",
			TaskID:           "task-001",
			TaskName:         "Pick and place bottle",
			FirstStart:       mustParseUTC(t, "2026-04-10T01:00:00Z"), // 10:00 JST
			LastEnd:          mustParseUTC(t, "2026-04-10T03:00:00Z"), // 12:00 JST
			WorkingSeconds:   120 * 60,
			CollectedSeconds: 65 * 60,
			DiscardedSeconds: 45 * 60,
			EpisodeCount:     53,
		},
	}

	uc := newTestUsecase(&stubOperatorYieldRepo{rows: rows})
	got, err := uc.Export(context.Background(), OperatorYieldExportFilter{
		DateFrom: time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := string(got)
	if !strings.HasPrefix(body, "\xef\xbb\xbf") {
		t.Errorf("expected UTF-8 BOM prefix")
	}
	expectedRow := "John Doe,20260410_10:00-12:00,Pick and place bottle,120,10,45,65,53,54%"
	if !strings.Contains(body, expectedRow) {
		t.Errorf("expected row\n  %q\nin output:\n%s", expectedRow, body)
	}
	expectedHeader := "オペレーター名,稼働日時,タスク,稼働時間(分),非稼働時間(分),破棄データ(分),収集データ(分),EP数,収率"
	if !strings.Contains(body, expectedHeader) {
		t.Errorf("expected header\n  %q\nin output:\n%s", expectedHeader, body)
	}
}

func TestOperatorYieldExport_ZeroWorkingTimeYieldsDash(t *testing.T) {
	rows := []repository.OperatorYieldExportRow{
		{
			WorkDate:       time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
			OperatorName:   "Op",
			TaskName:       "T",
			FirstStart:     mustParseUTC(t, "2026-04-10T01:00:00Z"),
			LastEnd:        mustParseUTC(t, "2026-04-10T01:00:00Z"),
			WorkingSeconds: 0,
		},
	}
	uc := newTestUsecase(&stubOperatorYieldRepo{rows: rows})
	got, err := uc.Export(context.Background(), OperatorYieldExportFilter{
		DateFrom: time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(got), ",-\n") && !strings.HasSuffix(strings.TrimRight(string(got), "\r\n"), ",-") {
		t.Errorf("expected yield column to be '-' for zero working time. got:\n%s", string(got))
	}
}

func TestOperatorYieldExport_NonWorkingClampedToZero(t *testing.T) {
	// collected+discarded > working can happen with overlapping/corrupt rows;
	// non-working must clamp to 0 rather than emit a negative value.
	rows := []repository.OperatorYieldExportRow{
		{
			WorkDate:         time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
			OperatorName:     "Op",
			TaskName:         "T",
			FirstStart:       mustParseUTC(t, "2026-04-10T01:00:00Z"),
			LastEnd:          mustParseUTC(t, "2026-04-10T02:00:00Z"),
			WorkingSeconds:   60 * 60,
			CollectedSeconds: 50 * 60,
			DiscardedSeconds: 30 * 60,
			EpisodeCount:     1,
		},
	}
	uc := newTestUsecase(&stubOperatorYieldRepo{rows: rows})
	got, err := uc.Export(context.Background(), OperatorYieldExportFilter{
		DateFrom: time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Columns: name,period,task,working,non_working,discarded,collected,ep,yield
	if !strings.Contains(string(got), ",60,0,30,50,") {
		t.Errorf("expected non-working clamped to 0. got:\n%s", string(got))
	}
}

func TestOperatorYieldExport_RejectsInvertedDateRange(t *testing.T) {
	uc := newTestUsecase(&stubOperatorYieldRepo{})
	_, err := uc.Export(context.Background(), OperatorYieldExportFilter{
		DateFrom: time.Date(2026, 4, 11, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
	})
	if err == nil {
		t.Fatal("expected error for inverted date range, got nil")
	}
}

func TestOperatorYieldExport_PreservesRowOrder(t *testing.T) {
	// Ordering is the repository's job (SQL ORDER BY); usecase must not reorder.
	rows := []repository.OperatorYieldExportRow{
		{
			OperatorName:   "Alice",
			TaskName:       "First",
			FirstStart:     mustParseUTC(t, "2026-04-10T01:00:00Z"),
			LastEnd:        mustParseUTC(t, "2026-04-10T02:00:00Z"),
			WorkingSeconds: 60 * 60,
		},
		{
			OperatorName:   "Bob",
			TaskName:       "Second",
			FirstStart:     mustParseUTC(t, "2026-04-10T01:00:00Z"),
			LastEnd:        mustParseUTC(t, "2026-04-10T02:00:00Z"),
			WorkingSeconds: 60 * 60,
		},
	}
	uc := newTestUsecase(&stubOperatorYieldRepo{rows: rows})
	got, err := uc.Export(context.Background(), OperatorYieldExportFilter{
		DateFrom: time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := string(got)
	idxAlice := strings.Index(body, "Alice")
	idxBob := strings.Index(body, "Bob")
	if idxAlice < 0 || idxBob < 0 || idxAlice >= idxBob {
		t.Errorf("expected Alice row to precede Bob row. got:\n%s", body)
	}
}

func TestOperatorYieldExport_EscapesCommaInTaskName(t *testing.T) {
	// Asserts we delegate quoting to encoding/csv rather than rolling our own.
	rows := []repository.OperatorYieldExportRow{
		{
			OperatorName:   "Alice",
			TaskName:       "task,with,commas",
			FirstStart:     mustParseUTC(t, "2026-04-10T01:00:00Z"),
			LastEnd:        mustParseUTC(t, "2026-04-10T02:00:00Z"),
			WorkingSeconds: 60 * 60,
		},
	}
	uc := newTestUsecase(&stubOperatorYieldRepo{rows: rows})
	got, err := uc.Export(context.Background(), OperatorYieldExportFilter{
		DateFrom: time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(got), `"task,with,commas"`) {
		t.Errorf("expected task name with commas to be quoted. got:\n%s", string(got))
	}
}

func TestOperatorYieldExport_AccurateSecondLevelArithmetic(t *testing.T) {
	// Regression guard for the floor-then-subtract drift:
	//   180s working / 119s collected / 59s discarded
	//   wrong (floor first):  3 - 1 - 0 = 2 minutes non-working
	//   correct (sub first):  floor((180-119-59)/60) = 0
	//   yield 119/180 = 66.1% → 66%, not skewed by intermediate floor()s.
	rows := []repository.OperatorYieldExportRow{
		{
			OperatorName:     "Op",
			TaskName:         "T",
			FirstStart:       mustParseUTC(t, "2026-04-10T01:00:00Z"),
			LastEnd:          mustParseUTC(t, "2026-04-10T01:03:00Z"),
			WorkingSeconds:   180,
			CollectedSeconds: 119,
			DiscardedSeconds: 59,
			EpisodeCount:     1,
		},
	}
	uc := newTestUsecase(&stubOperatorYieldRepo{rows: rows})
	got, err := uc.Export(context.Background(), OperatorYieldExportFilter{
		DateFrom: time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantRow := "Op,20260410_10:00-10:03,T,3,0,0,1,1,66%"
	if !strings.Contains(string(got), wantRow) {
		t.Errorf("expected accurate row\n  %q\nin output:\n%s", wantRow, string(got))
	}
}

func TestOperatorYieldExport_SanitizesFormulaInjection(t *testing.T) {
	// User-supplied names beginning with =/+/-/@ are weaponisable in Excel/
	// Sheets unless prefixed with a single quote.
	rows := []repository.OperatorYieldExportRow{
		{
			OperatorName:   "=cmd|'/c calc'!A1",
			TaskName:       "+1+1",
			FirstStart:     mustParseUTC(t, "2026-04-10T01:00:00Z"),
			LastEnd:        mustParseUTC(t, "2026-04-10T02:00:00Z"),
			WorkingSeconds: 60 * 60,
		},
		{
			OperatorName:   "@SUM(A1)",
			TaskName:       "-2+3",
			FirstStart:     mustParseUTC(t, "2026-04-10T01:00:00Z"),
			LastEnd:        mustParseUTC(t, "2026-04-10T02:00:00Z"),
			WorkingSeconds: 60 * 60,
		},
	}
	uc := newTestUsecase(&stubOperatorYieldRepo{rows: rows})
	got, err := uc.Export(context.Background(), OperatorYieldExportFilter{
		DateFrom: time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := string(got)
	for _, want := range []string{
		"'=cmd|'/c calc'!A1",
		"'+1+1",
		"'@SUM(A1)",
		"'-2+3",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("expected sanitized cell %q in output:\n%s", want, body)
		}
	}
	// Raw injection tokens must NOT appear at the start of a CSV cell.
	for _, banned := range []string{
		"\n=cmd",
		"\n+1+1",
		"\n@SUM",
		"\n-2+3",
	} {
		if strings.Contains(body, banned) {
			t.Errorf("expected %q to be neutralised but found in output:\n%s", banned, body)
		}
	}
}

func TestOperatorYieldExport_EmptyResultStillEmitsHeader(t *testing.T) {
	uc := newTestUsecase(&stubOperatorYieldRepo{rows: nil})
	got, err := uc.Export(context.Background(), OperatorYieldExportFilter{
		DateFrom: time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := string(got)
	if !strings.HasPrefix(body, "\xef\xbb\xbf") {
		t.Errorf("expected BOM even on empty result")
	}
	if !strings.Contains(body, "オペレーター名") {
		t.Errorf("expected header on empty result. got:\n%s", body)
	}
}
