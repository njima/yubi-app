package usecase

import (
	"context"
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

type stubLocationRepo struct {
	existing     model.Location
	updateCalled bool
	updateArg    model.Location
}

func (s *stubLocationRepo) Create(ctx context.Context, conn repository.Conn, loc model.Location) (model.Location, error) {
	return loc, nil
}

func (s *stubLocationRepo) GetByID(ctx context.Context, conn repository.Conn, id string) (model.Location, error) {
	return s.existing, nil
}

func (s *stubLocationRepo) List(ctx context.Context, conn repository.Conn, filter repository.LocationListFilter, limit, offset int) (model.Locations, int, error) {
	return nil, 0, nil
}

func (s *stubLocationRepo) Update(ctx context.Context, conn repository.Conn, loc model.Location) (model.Location, error) {
	s.updateCalled = true
	s.updateArg = loc
	return loc, nil
}

func (s *stubLocationRepo) Delete(ctx context.Context, conn repository.Conn, id string) error {
	return nil
}

type stubSiteRepo struct {
	existing     model.Site
	updateCalled bool
	updateArg    model.Site
}

func (s *stubSiteRepo) Create(ctx context.Context, conn repository.Conn, site model.Site) (model.Site, error) {
	return site, nil
}

func (s *stubSiteRepo) GetByID(ctx context.Context, conn repository.Conn, id string) (model.Site, error) {
	return s.existing, nil
}

func (s *stubSiteRepo) List(ctx context.Context, conn repository.Conn, filter repository.SiteListFilter, limit, offset int) (model.Sites, int, error) {
	return nil, 0, nil
}

func (s *stubSiteRepo) Update(ctx context.Context, conn repository.Conn, site model.Site) (model.Site, error) {
	s.updateCalled = true
	s.updateArg = site
	return site, nil
}

func (s *stubSiteRepo) Delete(ctx context.Context, conn repository.Conn, id string) error {
	return nil
}

type stubOrganizationRepo struct {
	existing     model.Organization
	updateCalled bool
	updateArg    model.Organization
}

func (s *stubOrganizationRepo) Create(ctx context.Context, conn repository.Conn, org model.Organization) (model.Organization, error) {
	return org, nil
}

func (s *stubOrganizationRepo) GetByNaturalID(ctx context.Context, conn repository.Conn, idNatural string) (model.Organization, error) {
	return s.existing, nil
}

func (s *stubOrganizationRepo) List(ctx context.Context, conn repository.Conn, limit, offset int) (model.Organizations, int, error) {
	return nil, 0, nil
}

func (s *stubOrganizationRepo) Update(ctx context.Context, conn repository.Conn, org model.Organization) (model.Organization, error) {
	s.updateCalled = true
	s.updateArg = org
	return org, nil
}

func (s *stubOrganizationRepo) Delete(ctx context.Context, conn repository.Conn, idNatural string) error {
	return nil
}

func TestLocationUsecase_Update_UsesPointerFields(t *testing.T) {
	repo := &stubLocationRepo{
		existing: model.NewLocation(1, "loc-1", "org-1", "site-1", "Dock", testTime, nil),
	}
	uc := NewLocation(repo, repository.NewDataAccess(nil, nil))
	name := "Charging Dock"

	got, err := uc.Update(context.Background(), LocationUpdateInput{ID: "loc-1", Name: &name})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.updateCalled {
		t.Fatalf("expected repo.Update to be called")
	}
	if repo.updateArg.Name != name {
		t.Errorf("repo update Name = %q, want %q", repo.updateArg.Name, name)
	}
	if got.Name != name {
		t.Errorf("returned Name = %q, want %q", got.Name, name)
	}
}

func TestLocationUsecase_Update_NoFieldsReturnsExistingLocation(t *testing.T) {
	existing := model.NewLocation(1, "loc-1", "org-1", "site-1", "Dock", testTime, nil)
	repo := &stubLocationRepo{existing: existing}
	uc := NewLocation(repo, repository.NewDataAccess(nil, nil))

	got, err := uc.Update(context.Background(), LocationUpdateInput{ID: "loc-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.updateCalled {
		t.Fatalf("expected no repo.Update call when no fields are provided")
	}
	if got.Name != existing.Name {
		t.Errorf("got Name = %q, want existing name %q", got.Name, existing.Name)
	}
}

func TestLocationUsecase_Update_RejectsExplicitEmptyName(t *testing.T) {
	repo := &stubLocationRepo{
		existing: model.NewLocation(1, "loc-1", "org-1", "site-1", "Dock", testTime, nil),
	}
	uc := NewLocation(repo, repository.NewDataAccess(nil, nil))
	emptyName := ""

	_, err := uc.Update(context.Background(), LocationUpdateInput{ID: "loc-1", Name: &emptyName})
	if err == nil {
		t.Fatalf("expected validation error for explicit empty name")
	}
	if repo.updateCalled {
		t.Fatalf("invalid update should not reach repo.Update")
	}
}

func TestSiteUsecase_Update_NoFieldsReturnsExistingSite(t *testing.T) {
	existing := model.NewSite(1, "site-1", "org-1", "Main Site", testTime, nil)
	repo := &stubSiteRepo{existing: existing}
	uc := NewSite(repo, repository.NewDataAccess(nil, nil))

	got, err := uc.Update(context.Background(), SiteUpdateInput{ID: "site-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.updateCalled {
		t.Fatalf("expected no repo.Update call when no fields are provided")
	}
	if got.Name != existing.Name {
		t.Errorf("got Name = %q, want existing name %q", got.Name, existing.Name)
	}
}

func TestSiteUsecase_Update_RejectsExplicitEmptyName(t *testing.T) {
	repo := &stubSiteRepo{
		existing: model.NewSite(1, "site-1", "org-1", "Main Site", testTime, nil),
	}
	uc := NewSite(repo, repository.NewDataAccess(nil, nil))
	emptyName := ""

	_, err := uc.Update(context.Background(), SiteUpdateInput{ID: "site-1", Name: &emptyName})
	if err == nil {
		t.Fatalf("expected validation error for explicit empty name")
	}
	if repo.updateCalled {
		t.Fatalf("invalid update should not reach repo.Update")
	}
}

func TestOrganizationUsecase_Update_AllowsDescriptionOnlyUpdate(t *testing.T) {
	oldDescription := "old"
	repo := &stubOrganizationRepo{
		existing: model.NewOrganization(1, "org-1", "Airoa", model.OrganizationKindTeam, &oldDescription, testTime, nil),
	}
	uc := NewOrganization(repo, repository.NewDataAccess(nil, nil))
	description := "new"

	got, err := uc.Update(context.Background(), OrganizationUpdateInput{
		OrganizationID: "org-1",
		Description:    &description,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.updateCalled {
		t.Fatalf("expected repo.Update to be called")
	}
	if repo.updateArg.Name != "Airoa" {
		t.Errorf("repo update Name = %q, want existing name", repo.updateArg.Name)
	}
	if got.Description == nil || *got.Description != description {
		t.Errorf("returned Description = %v, want %q", got.Description, description)
	}
}

func TestOrganizationUsecase_Update_NoFieldsReturnsExistingOrganization(t *testing.T) {
	description := "old"
	existing := model.NewOrganization(1, "org-1", "Airoa", model.OrganizationKindTeam, &description, testTime, nil)
	repo := &stubOrganizationRepo{existing: existing}
	uc := NewOrganization(repo, repository.NewDataAccess(nil, nil))

	got, err := uc.Update(context.Background(), OrganizationUpdateInput{OrganizationID: "org-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.updateCalled {
		t.Fatalf("expected no repo.Update call when no fields are provided")
	}
	if got.Name != existing.Name {
		t.Errorf("got Name = %q, want existing name %q", got.Name, existing.Name)
	}
}

func TestOrganizationUsecase_Update_RejectsExplicitEmptyDisplayName(t *testing.T) {
	description := "old"
	repo := &stubOrganizationRepo{
		existing: model.NewOrganization(1, "org-1", "Airoa", model.OrganizationKindTeam, &description, testTime, nil),
	}
	uc := NewOrganization(repo, repository.NewDataAccess(nil, nil))
	emptyName := ""

	_, err := uc.Update(context.Background(), OrganizationUpdateInput{
		OrganizationID: "org-1",
		DisplayName:    &emptyName,
	})
	if err == nil {
		t.Fatalf("expected validation error for explicit empty display name")
	}
	if repo.updateCalled {
		t.Fatalf("invalid update should not reach repo.Update")
	}
}
