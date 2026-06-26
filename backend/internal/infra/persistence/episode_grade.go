package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/uptrace/bun"
)

type episodeGrade struct{}

func NewEpisodeGrade() *episodeGrade { return &episodeGrade{} }

func episodeGradeModelToEntity(g model.EpisodeGrade) entity.EpisodeGrade {
	return entity.EpisodeGrade{
		EpisodeID:      g.EpisodeID,
		UserID:         g.UserID,
		OrganizationID: g.OrganizationID,
		Grade:          g.Grade,
		Comment:        g.Comment,
		GradedAt:       g.GradedAt,
	}
}

func entityToEpisodeGradeModel(e entity.EpisodeGrade) model.EpisodeGrade {
	var updatedAt *time.Time
	if !e.UpdatedAt.IsZero() {
		t := e.UpdatedAt
		updatedAt = &t
	}
	return model.NewEpisodeGrade(
		e.OrganizationID,
		e.EpisodeID,
		e.UserID,
		e.Grade,
		e.Comment,
		e.GradedAt,
		e.CreatedAt,
		updatedAt,
	)
}

func (g *episodeGrade) GetAverageMap(ctx context.Context, conn repository.DBConn, episodeIDs []string) (map[string]repository.GradeAggregate, error) {
	result := make(map[string]repository.GradeAggregate)

	if len(episodeIDs) == 0 {
		return result, nil
	}

	type aggRow struct {
		EpisodeID string  `bun:"episode_id"`
		Average   float64 `bun:"average"`
		Count     int     `bun:"count"`
	}

	var rows []aggRow
	if err := conn.NewSelect().
		Model((*entity.EpisodeGrade)(nil)).
		ColumnExpr("episode_id").
		ColumnExpr("AVG(grade) AS average").
		ColumnExpr("COUNT(*) AS count").
		Where("episode_id IN (?)", bun.In(episodeIDs)).
		Group("episode_id").
		Scan(ctx, &rows); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to aggregate episode grades: %v", err))
	}

	for _, r := range rows {
		result[r.EpisodeID] = repository.GradeAggregate{
			Average: r.Average,
			Count:   r.Count,
		}
	}

	return result, nil
}

func (g *episodeGrade) Upsert(ctx context.Context, conn repository.DBConn, grade model.EpisodeGrade) (model.EpisodeGrade, error) {
	dbGrade := episodeGradeModelToEntity(grade)

	var inserted entity.EpisodeGrade
	if err := conn.NewInsert().
		Model(&dbGrade).
		On("CONFLICT (episode_id, user_id) DO UPDATE").
		Set("grade = EXCLUDED.grade").
		Set("comment = EXCLUDED.comment").
		Set("graded_at = EXCLUDED.graded_at").
		Set("updated_at = EXCLUDED.updated_at").
		Returning("*").
		Scan(ctx, &inserted); err != nil {
		return model.EpisodeGrade{}, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to upsert episode grade: %v", err))
	}

	return entityToEpisodeGradeModel(inserted), nil
}

func (g *episodeGrade) GetMyGrade(ctx context.Context, conn repository.DBConn, episodeID, userID string) (*model.EpisodeGrade, error) {
	var row entity.EpisodeGrade
	if err := conn.NewSelect().
		Model(&row).
		Where("episode_id = ?", episodeID).
		Where("user_id = ?", userID).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to get my episode grade: %v", err))
	}

	out := entityToEpisodeGradeModel(row)
	return &out, nil
}

func (g *episodeGrade) ListByEpisodeID(ctx context.Context, conn repository.DBConn, episodeID string, limit, offset int) ([]repository.EpisodeGradeListItem, int, error) {
	type joinedRow struct {
		entity.EpisodeGrade `bun:",extend"`
		UserName            string `bun:"user_name"`
	}

	var rows []joinedRow
	if err := conn.NewSelect().
		Model((*entity.EpisodeGrade)(nil)).
		ColumnExpr("eg.*").
		ColumnExpr(`"user".name AS user_name`).
		Join(`JOIN "user" ON "user".id_natural = eg.user_id`).
		Where("eg.episode_id = ?", episodeID).
		Order("eg.graded_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx, &rows); err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to list episode grades: %v", err))
	}

	var total int
	count, err := conn.NewSelect().
		Model((*entity.EpisodeGrade)(nil)).
		Where("eg.episode_id = ?", episodeID).
		Count(ctx)
	if err != nil {
		return nil, 0, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeDatabaseError, "failed to count episode grades: %v", err))
	}
	total = count

	items := make([]repository.EpisodeGradeListItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, repository.EpisodeGradeListItem{
			Grade:    entityToEpisodeGradeModel(r.EpisodeGrade),
			UserName: r.UserName,
		})
	}
	return items, total, nil
}
