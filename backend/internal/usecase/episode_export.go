package usecase

import (
	"bytes"
	"context"
	"encoding/csv"
	"time"

	"github.com/uptrace/bun"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

// EpisodeExportUsecase generates a CSV export of episodes.
type EpisodeExportUsecase interface {
	Export(ctx context.Context, filter repository.EpisodeExportFilter) ([]byte, error)
}

type episodeExport struct {
	episodeRepo repository.Episode
	db          *bun.DB
}

func NewEpisodeExport(episodeRepo repository.Episode, db *bun.DB) *episodeExport {
	return &episodeExport{episodeRepo: episodeRepo, db: db}
}

var episodeExportHeaders = []string{
	"id", "task_id", "task_version_id", "robot_id", "location_id",
	"user_id", "recorded_by", "status", "started_at", "finished_at", "created_at",
}

func (u *episodeExport) Export(ctx context.Context, filter repository.EpisodeExportFilter) ([]byte, error) {
	rows, err := u.episodeRepo.Export(ctx, u.db, filter)
	if err != nil {
		return nil, err
	}

	if len(rows) > repository.MaxEpisodeExportRows {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest,
			"export limit of %d episodes exceeded; apply filters to reduce the result set", repository.MaxEpisodeExportRows))
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write(episodeExportHeaders); err != nil {
		return nil, err
	}

	for _, row := range rows {
		record := []string{
			row.IDNatural,
			row.TaskID,
			row.TaskVersionID,
			row.RobotID,
			row.LocationID,
			row.UserID,
			derefStringExport(row.RecordedByID),
			episodeStatusLabel(row.Status),
			timeToStringExport(row.StartedAt),
			timeToStringExport(row.FinishedAt),
			row.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func episodeStatusLabel(status model.EpisodeStatus) string {
	switch status {
	case model.EpisodeStatusReady:
		return "Ready"
	case model.EpisodeStatusRecording:
		return "Recording"
	case model.EpisodeStatusCancel:
		return "Cancel"
	case model.EpisodeStatusCompleted:
		return "Completed"
	default:
		return "Unknown"
	}
}

func timeToStringExport(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02T15:04:05Z07:00")
}
