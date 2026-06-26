package persistence

import (
	"testing"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
)

var mapperTestTime = time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)

func TestLocationEntityToModel(t *testing.T) {
	dbLocation := entity.Location{
		ID:             1,
		IDNatural:      "loc-1",
		OrganizationID: "org-1",
		SiteID:         "site-1",
		Name:           "Dock",
		Timestamp: entity.Timestamp{
			CreatedAt: mapperTestTime,
		},
		Site: &entity.Site{Name: "Main Site"},
	}

	got := locationEntityToModel(dbLocation)

	if got.ID != dbLocation.ID || got.IDNatural != dbLocation.IDNatural {
		t.Fatalf("got location identifiers %+v, want %+v", got, dbLocation)
	}
	if got.SiteName != "Main Site" {
		t.Errorf("SiteName = %q, want related site name", got.SiteName)
	}
	if got.UpdatedAt != nil {
		t.Errorf("UpdatedAt = %v, want nil for zero timestamp", got.UpdatedAt)
	}
}

func TestSiteEntityToModel(t *testing.T) {
	dbSite := entity.Site{
		ID:             1,
		IDNatural:      "site-1",
		OrganizationID: "org-1",
		Name:           "Main Site",
		Timestamp: entity.Timestamp{
			CreatedAt: mapperTestTime,
			UpdatedAt: mapperTestTime,
		},
	}

	got := siteEntityToModel(dbSite)

	if got.IDNatural != dbSite.IDNatural || got.OrganizationID != dbSite.OrganizationID {
		t.Fatalf("got site %+v, want values from entity %+v", got, dbSite)
	}
	if got.UpdatedAt == nil || !got.UpdatedAt.Equal(mapperTestTime) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, mapperTestTime)
	}
}

func TestOrganizationEntityToModel(t *testing.T) {
	description := "robotics"
	dbOrganization := entity.Organization{
		ID:          1,
		IDNatural:   "org-1",
		Name:        "Airoa",
		Description: &description,
		Timestamp: entity.Timestamp{
			CreatedAt: mapperTestTime,
		},
	}

	got := organizationEntityToModel(dbOrganization)

	if got.IDNatural != dbOrganization.IDNatural || got.Name != dbOrganization.Name {
		t.Fatalf("got organization %+v, want values from entity %+v", got, dbOrganization)
	}
	if got.Description == nil || *got.Description != description {
		t.Errorf("Description = %v, want %q", got.Description, description)
	}
	if got.UpdatedAt != nil {
		t.Errorf("UpdatedAt = %v, want nil for zero timestamp", got.UpdatedAt)
	}
}
