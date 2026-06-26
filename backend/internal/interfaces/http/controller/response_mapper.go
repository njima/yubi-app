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

func userResponse(user model.User) openapi.UserResponse {
	return openapi.UserResponse{
		UserId:           user.IDNatural,
		Email:            user.Email,
		DisplayName:      user.Name,
		Role:             openAPIUserRolePtr(user.Role),
		OrganizationId:   user.OrganizationID,
		OrganizationName: user.OrganizationName,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		Locations:        locationSummaries(user.Locations),
		Sites:            siteSummaries(user.Sites),
	}
}

func userResponses(users model.Users) []openapi.UserResponse {
	result := make([]openapi.UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, userResponse(*user))
	}
	return result
}

func locationSummaries(locs []model.LocationSummary) []openapi.LocationSummary {
	result := make([]openapi.LocationSummary, 0, len(locs))
	for _, loc := range locs {
		result = append(result, openapi.LocationSummary{
			LocationId: loc.LocationID,
			Name:       loc.Name,
		})
	}
	return result
}

func siteSummaries(sites []model.SiteSummary) []openapi.SiteSummary {
	result := make([]openapi.SiteSummary, 0, len(sites))
	for _, site := range sites {
		result = append(result, openapi.SiteSummary{
			SiteId: site.SiteID,
			Name:   site.Name,
		})
	}
	return result
}
