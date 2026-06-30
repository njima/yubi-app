package persistence

import (
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
)

func updatedAtPtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func locationModelToEntity(loc model.Location) entity.Location {
	return entity.Location{
		IDNatural:      loc.IDNatural,
		OrganizationID: loc.OrganizationID,
		SiteID:         loc.SiteID,
		Name:           loc.Name,
	}
}

func locationEntityToModel(loc entity.Location) model.Location {
	siteName := ""
	if loc.Site != nil {
		siteName = loc.Site.Name
	}

	return model.Location{
		ID:             loc.ID,
		IDNatural:      loc.IDNatural,
		OrganizationID: loc.OrganizationID,
		SiteID:         loc.SiteID,
		SiteName:       siteName,
		Name:           loc.Name,
		CreatedAt:      loc.CreatedAt,
		UpdatedAt:      updatedAtPtr(loc.UpdatedAt),
	}
}

func siteModelToEntity(site model.Site) entity.Site {
	return entity.Site{
		IDNatural:      site.IDNatural,
		OrganizationID: site.OrganizationID,
		Name:           site.Name,
	}
}

func siteEntityToModel(site entity.Site) model.Site {
	return model.Site{
		ID:             site.ID,
		IDNatural:      site.IDNatural,
		OrganizationID: site.OrganizationID,
		Name:           site.Name,
		CreatedAt:      site.CreatedAt,
		UpdatedAt:      updatedAtPtr(site.UpdatedAt),
	}
}

func organizationModelToEntity(org model.Organization) entity.Organization {
	return entity.Organization{
		IDNatural:   org.IDNatural,
		Name:        org.Name,
		Kind:        string(org.Kind),
		Description: org.Description,
	}
}

func organizationEntityToModel(org entity.Organization) model.Organization {
	return model.Organization{
		ID:          org.ID,
		IDNatural:   org.IDNatural,
		Name:        org.Name,
		Kind:        model.OrganizationKind(org.Kind),
		Description: org.Description,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   updatedAtPtr(org.UpdatedAt),
	}
}

func robotEntityToModel(robot entity.Robot) model.Robot {
	orgName := ""
	if robot.Organization != nil {
		orgName = robot.Organization.Name
	}

	locName := ""
	siteID := ""
	siteName := ""
	if robot.Location != nil {
		locName = robot.Location.Name
		siteID = robot.Location.SiteID
		if robot.Location.Site != nil {
			siteName = robot.Location.Site.Name
		}
	}

	return model.NewRobot(
		robot.ID,
		robot.IDNatural,
		robot.OrganizationID,
		orgName,
		siteID,
		siteName,
		robot.LocationID,
		locName,
		robot.Name,
		&robot.RobotType,
		model.RobotStatus(robot.Status),
		leaderStatusToModel(robot.LeaderStatus),
		robot.LeaderFaultStartedAt,
		robot.FaultStartedAt,
		robot.LastHeartbeatAt,
		robot.OfflineReason,
		robot.RobotConfig,
		robot.ActiveEpisodeID,
		robot.ActiveUserID,
		robot.CreatedAt,
		updatedAtPtr(robot.UpdatedAt),
	)
}
