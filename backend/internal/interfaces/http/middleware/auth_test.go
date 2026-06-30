package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

type stubAuthUserUsecase struct {
	user       model.User
	membership model.OrganizationMembership
	resolveErr error

	resolvedUserID         string
	resolvedOrganizationID *string
}

func (s *stubAuthUserUsecase) Create(ctx context.Context, input usecase.CreateInput) (model.User, error) {
	return model.User{}, nil
}

func (s *stubAuthUserUsecase) Update(ctx context.Context, input usecase.UserUpdateInput) (model.User, error) {
	return model.User{}, nil
}

func (s *stubAuthUserUsecase) UpdateRole(ctx context.Context, input usecase.UserRoleUpdateInput) (model.User, error) {
	return model.User{}, nil
}

func (s *stubAuthUserUsecase) SetLocations(ctx context.Context, userID string, locationIDs []string) (model.User, error) {
	return model.User{}, nil
}

func (s *stubAuthUserUsecase) SetSites(ctx context.Context, userID string, siteIDs []string) (model.User, error) {
	return model.User{}, nil
}

func (s *stubAuthUserUsecase) GetByNaturalID(ctx context.Context, idNatural string) (model.User, error) {
	return s.user, nil
}

func (s *stubAuthUserUsecase) List(ctx context.Context, filter usecase.UserListFilter, page, limit int) (model.Users, int, error) {
	return nil, 0, nil
}

func (s *stubAuthUserUsecase) Delete(ctx context.Context, idNatural string) error {
	return nil
}

func (s *stubAuthUserUsecase) FindOrProvisionGoogleUser(ctx context.Context, input usecase.GoogleUserInput) (usecase.AuthenticatedUserSession, error) {
	return usecase.AuthenticatedUserSession{}, nil
}

func (s *stubAuthUserUsecase) ResolveActiveMembership(ctx context.Context, userID string, organizationID *string) (model.OrganizationMembership, error) {
	s.resolvedUserID = userID
	if organizationID != nil {
		orgID := *organizationID
		s.resolvedOrganizationID = &orgID
	}
	if s.resolveErr != nil {
		return model.OrganizationMembership{}, s.resolveErr
	}
	return s.membership, nil
}

type stubAuthRobotUsecase struct {
	robot model.Robot
}

func (s *stubAuthRobotUsecase) Create(ctx context.Context, input usecase.RobotCreateInput) (model.Robot, error) {
	return model.Robot{}, nil
}

func (s *stubAuthRobotUsecase) GetByID(ctx context.Context, id string) (model.Robot, error) {
	return s.robot, nil
}

func (s *stubAuthRobotUsecase) List(ctx context.Context, filter usecase.RobotListFilter, page, limit int) (model.Robots, int, error) {
	return nil, 0, nil
}

func (s *stubAuthRobotUsecase) ListTypes(ctx context.Context, filter usecase.RobotTypeFilter) ([]string, error) {
	return nil, nil
}

func (s *stubAuthRobotUsecase) Update(ctx context.Context, input usecase.RobotUpdateInput) (model.Robot, error) {
	return model.Robot{}, nil
}

func (s *stubAuthRobotUsecase) Delete(ctx context.Context, id string) error {
	return nil
}

func TestUserAuthUsesActiveOrganizationMembershipFromHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userUC := &stubAuthUserUsecase{
		user: model.User{
			IDNatural: "user-1",
		},
		membership: model.OrganizationMembership{
			UserID:         "user-1",
			OrganizationID: "org-active",
			Role:           model.UserRoleManager,
		},
	}
	w := httptest.NewRecorder()
	called := false
	var gotOrgID string
	var gotRole model.UserRole
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userAuth(c, userUC)
	})
	router.GET("/api/episodes", func(c *gin.Context) {
		called = true
		var err error
		gotOrgID, err = requestctx.OrganizationID(c.Request.Context())
		if err != nil {
			t.Fatalf("OrganizationID() error = %v", err)
		}
		gotRole, err = requestctx.UserRole(c.Request.Context())
		if err != nil {
			t.Fatalf("UserRole() error = %v", err)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/api/episodes", nil)
	req.Header.Set(headerXUserID, "user-1")
	req.Header.Set(headerXOrganizationID, "org-active")

	router.ServeHTTP(w, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}
	if userUC.resolvedUserID != "user-1" {
		t.Errorf("resolved user id = %q, want %q", userUC.resolvedUserID, "user-1")
	}
	if userUC.resolvedOrganizationID == nil || *userUC.resolvedOrganizationID != "org-active" {
		t.Fatalf("resolved organization id = %v, want %q", userUC.resolvedOrganizationID, "org-active")
	}
	if gotOrgID != "org-active" {
		t.Errorf("context organization id = %q, want %q", gotOrgID, "org-active")
	}
	if gotRole != model.UserRoleManager {
		t.Errorf("context role = %v, want %v", gotRole, model.UserRoleManager)
	}
}

func TestRobotAuthByHeadersUsesRobotOrganizationMembership(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userUC := &stubAuthUserUsecase{
		user: model.User{
			IDNatural: "user-1",
		},
		membership: model.OrganizationMembership{
			UserID:         "user-1",
			OrganizationID: "robot-org",
			Role:           model.UserRoleOperator,
		},
	}
	robotUC := &stubAuthRobotUsecase{
		robot: model.Robot{
			IDNatural:      "robot-1",
			OrganizationID: "robot-org",
		},
	}
	w := httptest.NewRecorder()
	called := false
	var gotOrgID string
	var gotRole model.UserRole
	router := gin.New()
	router.Use(func(c *gin.Context) {
		robotAuthByHeaders(c, userUC, robotUC)
	})
	router.GET("/api/robot/status", func(c *gin.Context) {
		called = true
		var err error
		gotOrgID, err = requestctx.OrganizationID(c.Request.Context())
		if err != nil {
			t.Fatalf("OrganizationID() error = %v", err)
		}
		gotRole, err = requestctx.UserRole(c.Request.Context())
		if err != nil {
			t.Fatalf("UserRole() error = %v", err)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/api/robot/status", nil)
	req.Header.Set(headerXUserID, "user-1")
	req.Header.Set(headerXRobotID, "robot-1")

	router.ServeHTTP(w, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}
	if userUC.resolvedOrganizationID == nil || *userUC.resolvedOrganizationID != "robot-org" {
		t.Fatalf("resolved organization id = %v, want %q", userUC.resolvedOrganizationID, "robot-org")
	}
	if gotOrgID != "robot-org" {
		t.Errorf("context organization id = %q, want %q", gotOrgID, "robot-org")
	}
	if gotRole != model.UserRoleOperator {
		t.Errorf("context role = %v, want %v", gotRole, model.UserRoleOperator)
	}
}

func TestRobotAuthByHeadersRejectsUserWithoutRobotOrganizationMembership(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userUC := &stubAuthUserUsecase{
		user: model.User{
			IDNatural: "user-1",
		},
		resolveErr: apperror.NewError(apperror.NewMessage(apperror.CodeForbidden, "organization membership not found")),
	}
	robotUC := &stubAuthRobotUsecase{
		robot: model.Robot{
			IDNatural:      "robot-1",
			OrganizationID: "robot-org",
		},
	}
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		robotAuthByHeaders(c, userUC, robotUC)
	})
	router.GET("/api/robot/status", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/robot/status", nil)
	req.Header.Set(headerXUserID, "user-1")
	req.Header.Set(headerXRobotID, "robot-1")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestUserAuthReturnsInternalServerErrorWhenMembershipLookupFails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userUC := &stubAuthUserUsecase{
		user: model.User{
			IDNatural: "user-1",
		},
		resolveErr: errors.New("database unavailable"),
	}
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userAuth(c, userUC)
	})
	router.GET("/api/episodes", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/episodes", nil)
	req.Header.Set(headerXUserID, "user-1")
	req.Header.Set(headerXOrganizationID, "org-active")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
