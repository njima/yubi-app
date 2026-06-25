package usecase

import (
	"bytes"
	"context"
	"encoding/csv"
	"math"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type OperatorYieldExportUsecase interface {
	Export(ctx context.Context, filter repository.OperatorYieldExportFilter) ([]byte, error)
}

type operatorYieldExport struct {
	repo   repository.OperatorYield
	data   repository.DataAccess
	logger zerolog.Logger
}

func NewOperatorYieldExport(repo repository.OperatorYield, data repository.DataAccess, logger zerolog.Logger) *operatorYieldExport {
	return &operatorYieldExport{repo: repo, data: data, logger: logger}
}

var operatorYieldExportHeaders = []string{
	"オペレーター名", "稼働日時", "タスク",
	"稼働時間(分)", "非稼働時間(分)", "破棄データ(分)", "収集データ(分)",
	"EP数", "収率",
}

// Prepended so Excel opens the CSV without mojibake.
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// sanitizeUserText neutralises CSV/Formula Injection by prefixing a single
// quote to user-supplied values that start with a formula trigger. Do NOT
// apply to system-generated cells (e.g. yield "-" would become '-).
func sanitizeUserText(s string) string {
	if s == "" {
		return s
	}
	switch s[0] {
	case '=', '+', '-', '@', '\t', '\r':
		return "'" + s
	}
	return s
}

func (u *operatorYieldExport) Export(ctx context.Context, filter repository.OperatorYieldExportFilter) ([]byte, error) {
	if filter.DateFrom.After(filter.DateTo) {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest,
			"date_from must be on or before date_to"))
	}

	start := time.Now()
	rows, err := u.repo.Export(ctx, u.data.Conn(), filter)
	if err != nil {
		return nil, err
	}
	queryElapsed := time.Since(start)

	if len(rows) > repository.MaxOperatorYieldExportRows {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest,
			"export limit of %d rows exceeded; apply filters to reduce the result set",
			repository.MaxOperatorYieldExportRows))
	}

	var buf bytes.Buffer
	buf.Write(utf8BOM)
	w := csv.NewWriter(&buf)

	if err := w.Write(operatorYieldExportHeaders); err != nil {
		return nil, err
	}

	for _, r := range rows {
		// Subtract in seconds first: floor(a)-floor(b)-floor(c) != floor(a-b-c).
		nonWorkingSec := r.WorkingSeconds - r.CollectedSeconds - r.DiscardedSeconds
		if nonWorkingSec < 0 {
			// Overlapping episodes can push collected+discarded above the
			// [first_start, last_end] span. Clamp rather than emit a negative.
			nonWorkingSec = 0
		}

		workingMin := r.WorkingSeconds / 60
		collectedMin := r.CollectedSeconds / 60
		discardedMin := r.DiscardedSeconds / 60
		nonWorkingMin := nonWorkingSec / 60

		yieldStr := "-"
		if r.WorkingSeconds > 0 {
			pct := int(math.Round(float64(r.CollectedSeconds) / float64(r.WorkingSeconds) * 100))
			yieldStr = strconv.Itoa(pct) + "%"
		}

		startJST := r.FirstStart.In(repository.JSTLocation)
		endJST := r.LastEnd.In(repository.JSTLocation)
		workingPeriod := startJST.Format("20060102") + "_" +
			startJST.Format("15:04") + "-" + endJST.Format("15:04")

		record := []string{
			sanitizeUserText(r.OperatorName),
			workingPeriod,
			sanitizeUserText(r.TaskName),
			strconv.FormatInt(workingMin, 10),
			strconv.FormatInt(nonWorkingMin, 10),
			strconv.FormatInt(discardedMin, 10),
			strconv.FormatInt(collectedMin, 10),
			strconv.FormatInt(r.EpisodeCount, 10),
			yieldStr,
		}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	u.logger.Info().
		Int("rows", len(rows)).
		Int("bytes", buf.Len()).
		Dur("query_elapsed", queryElapsed).
		Dur("total_elapsed", time.Since(start)).
		Time("date_from", filter.DateFrom).
		Time("date_to", filter.DateTo).
		Msg("operator-yield export completed")

	return buf.Bytes(), nil
}
