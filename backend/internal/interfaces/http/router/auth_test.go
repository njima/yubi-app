package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/interfaces/http/controller"
	"github.com/gin-gonic/gin"
)

func TestGoogleAuthSessionRouteRequiresInternalSecretWhenConfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := engine.Group("/api")
	ctrl := &stubGoogleAuthController{}
	registerAuthRoutes(api, ctrl, "secret")

	req := httptest.NewRequest(http.MethodPost, "/api/auth/google/session", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if ctrl.called {
		t.Fatal("controller should not be called when internal secret is invalid")
	}
}

type stubGoogleAuthController struct {
	called bool
}

func (s *stubGoogleAuthController) CreateGoogleAuthSession(ctx context.Context, req controller.GoogleAuthSessionRequest) (controller.GoogleAuthSessionResponse, error) {
	s.called = true
	return controller.GoogleAuthSessionResponse{}, nil
}
