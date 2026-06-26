package controller

import (
	"testing"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

func TestLocationResponse(t *testing.T) {
	loc := model.Location{
		IDNatural: "loc-1",
		SiteID:    "site-1",
		SiteName:  "Main Site",
		Name:      "Dock",
	}

	got := locationResponse(loc)

	if got.Id != loc.IDNatural || got.SiteId != loc.SiteID || got.SiteName != loc.SiteName || got.Name != loc.Name {
		t.Fatalf("locationResponse() = %+v, want values from %+v", got, loc)
	}
}

func TestSiteResponse(t *testing.T) {
	createdAt := time.Date(2026, 6, 26, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	site := model.Site{
		IDNatural:      "site-1",
		OrganizationID: "org-1",
		Name:           "Main Site",
		CreatedAt:      createdAt,
		UpdatedAt:      &updatedAt,
	}

	got := siteResponse(site)

	if got.Id != site.IDNatural || got.OrganizationId != site.OrganizationID || got.Name != site.Name {
		t.Fatalf("siteResponse() = %+v, want values from %+v", got, site)
	}
	if got.CreatedAt == nil || !got.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, createdAt)
	}
	if got.UpdatedAt == nil || !got.UpdatedAt.Equal(updatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, updatedAt)
	}
}

func TestOrganizationResponse(t *testing.T) {
	createdAt := time.Date(2026, 6, 26, 10, 0, 0, 0, time.UTC)
	description := "robotics"
	org := model.Organization{
		IDNatural:   "org-1",
		Name:        "Airoa",
		Description: &description,
		CreatedAt:   createdAt,
	}

	got := organizationResponse(org)

	if got.OrganizationId != org.IDNatural || got.DisplayName != org.Name {
		t.Fatalf("organizationResponse() = %+v, want values from %+v", got, org)
	}
	if got.Description == nil || *got.Description != description {
		t.Errorf("Description = %v, want %q", got.Description, description)
	}
	if got.CreatedAt == nil || !got.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, createdAt)
	}
}
