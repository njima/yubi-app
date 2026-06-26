package repository

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

type GradeAggregate struct {
	Average float64
	Count   int
}

// EpisodeGradeListItem couples the grade row with the grader's display name
// pulled from the user table in a single JOIN, so the detail view does not
// have to issue per-row user lookups.
type EpisodeGradeListItem struct {
	Grade    model.EpisodeGrade
	UserName string
}

type EpisodeGrade interface {
	// GetAverageMap omits episodes that have no grades from the returned map.
	GetAverageMap(ctx context.Context, conn Conn, episodeIDs []string) (map[string]GradeAggregate, error)

	Upsert(ctx context.Context, conn Conn, grade model.EpisodeGrade) (model.EpisodeGrade, error)

	// GetMyGrade returns (nil, nil) when the user has not graded the episode.
	GetMyGrade(ctx context.Context, conn Conn, episodeID, userID string) (*model.EpisodeGrade, error)

	// ListByEpisodeID returns paginated grades for the episode together with
	// each grader's display name. Total is the unpaginated row count.
	ListByEpisodeID(ctx context.Context, conn Conn, episodeID string, limit, offset int) ([]EpisodeGradeListItem, int, error)
}
