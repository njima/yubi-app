package usecase

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
)

func (e *episode) buildPreviewPath(ctx context.Context, ep model.Episode) (model.EpisodePreviewPath, error) {
	if ep.StartedAt == nil {
		return model.EpisodePreviewPath{}, apperror.NewError(apperror.NewMessage(apperror.CodeInternal, "episode has no started_at"))
	}

	robot, err := e.rr.GetByID(ctx, e.data.Conn(), ep.RobotID)
	if err != nil {
		return model.EpisodePreviewPath{}, err
	}

	loc, err := e.locRepo.GetByID(ctx, e.data.Conn(), ep.LocationID)
	if err != nil {
		return model.EpisodePreviewPath{}, err
	}

	site, err := e.siteRepo.GetByID(ctx, e.data.Conn(), loc.SiteID)
	if err != nil {
		return model.EpisodePreviewPath{}, err
	}

	robotType := ""
	if robot.RobotType != nil {
		robotType = *robot.RobotType
	}

	return model.EpisodePreviewPath{
		UUID:         ep.IDNatural,
		Organization: robot.OrganizationName,
		Site:         site.Name,
		Location:     loc.Name,
		RobotType:    robotType,
		RobotID:      robot.IDNatural,
		StartedAt:    *ep.StartedAt,
	}, nil
}

func (e *episode) GetRecordings(ctx context.Context, episodeID string) (map[string]string, error) {
	ep, err := e.repo.GetByID(ctx, e.data.Conn(), episodeID)
	if err != nil {
		return nil, err
	}
	path, err := e.buildPreviewPath(ctx, ep)
	if err != nil {
		return nil, err
	}
	return e.recRepo.GetRecordingURLs(ctx, path)
}

func (e *episode) GetStats(ctx context.Context, episodeID string) (model.EpisodeRecordingStats, error) {
	ep, err := e.repo.GetByID(ctx, e.data.Conn(), episodeID)
	if err != nil {
		return nil, err
	}
	path, err := e.buildPreviewPath(ctx, ep)
	if err != nil {
		return nil, err
	}
	return e.recRepo.GetStats(ctx, path)
}
