package usecase

import (
	"context"
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
	"github.com/rs/zerolog"
)

func TestUserImport_ImportCreatesIdentityUserAndActiveOrganizationMembership(t *testing.T) {
	userRepo := &stubUserRepo{}
	membershipRepo := &stubOrganizationMembershipRepo{}
	uc := NewUserImport(
		userRepo,
		membershipRepo,
		repository.NewDataAccess(nil, stubTxRunner{}),
		zerolog.Nop(),
	)
	ctx := requestctx.SetOrganizationID(context.Background(), "org-1")

	got, err := uc.Import(ctx, "email,display_name,role\noperator@example.com,Operator,operator\n")

	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if got.ImportedCount != 1 || got.ErrorCount != 0 {
		t.Fatalf("Import() = %+v, want one imported row without errors", got)
	}
	if userRepo.created.Email != "operator@example.com" || userRepo.created.GoogleSub != "operator@example.com" {
		t.Errorf("created user = %+v, want identity user keyed by email", userRepo.created)
	}
	if len(membershipRepo.memberships) != 1 {
		t.Fatalf("created memberships = %+v, want one membership", membershipRepo.memberships)
	}
	membership := membershipRepo.memberships[0]
	if membership.UserID != userRepo.created.IDNatural || membership.OrganizationID != "org-1" {
		t.Errorf("membership scope = (%q, %q), want (%q, org-1)", membership.UserID, membership.OrganizationID, userRepo.created.IDNatural)
	}
	if membership.Role != model.UserRoleOperator {
		t.Errorf("membership role = %v, want operator", membership.Role)
	}
}
