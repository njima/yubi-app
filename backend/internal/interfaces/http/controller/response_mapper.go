package controller

import (
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
)

func locationResponse(loc model.Location) openapi.Location {
	return openapi.Location{
		Id:       loc.IDNatural,
		Name:     loc.Name,
		SiteId:   loc.SiteID,
		SiteName: loc.SiteName,
	}
}

func siteResponse(site model.Site) openapi.Site {
	return openapi.Site{
		Id:             site.IDNatural,
		Name:           site.Name,
		OrganizationId: site.OrganizationID,
		CreatedAt:      &site.CreatedAt,
		UpdatedAt:      site.UpdatedAt,
	}
}

func organizationResponse(org model.Organization) openapi.OrganizationResponse {
	return openapi.OrganizationResponse{
		OrganizationId: org.IDNatural,
		DisplayName:    org.Name,
		Description:    org.Description,
		CreatedAt:      &org.CreatedAt,
		UpdatedAt:      org.UpdatedAt,
	}
}
