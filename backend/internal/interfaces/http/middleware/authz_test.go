package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
	"github.com/gin-gonic/gin"
)

func TestNewAuthzMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	middleware := NewAuthzMiddleware()

	tests := []struct {
		name        string
		operationID string
		setupCtx    func(c *gin.Context)
		wantCalled  bool
		wantErr     bool
	}{
		{
			name:        "bypass operation passes without permission check",
			operationID: "ListOrganizations",
			setupCtx:    nil,
			wantCalled:  true,
			wantErr:     false,
		},
		{
			name:        "unmapped operation returns forbidden",
			operationID: "UnknownOperation",
			setupCtx: func(c *gin.Context) {
				ctx := requestctx.SetUserRole(c.Request.Context(), model.UserRoleAdmin)
				c.Request = c.Request.WithContext(ctx)
			},
			wantCalled: false,
			wantErr:    true,
		},
		{
			name:        "missing role in context returns forbidden",
			operationID: "ListEpisodes",
			setupCtx:    nil,
			wantCalled:  false,
			wantErr:     true,
		},
		{
			name:        "insufficient permission returns forbidden",
			operationID: "CreateEpisode",
			setupCtx: func(c *gin.Context) {
				ctx := requestctx.SetUserRole(c.Request.Context(), model.UserRoleViewer)
				c.Request = c.Request.WithContext(ctx)
			},
			wantCalled: false,
			wantErr:    true,
		},
		{
			name:        "sufficient permission passes to handler",
			operationID: "CreateEpisode",
			setupCtx: func(c *gin.Context) {
				ctx := requestctx.SetUserRole(c.Request.Context(), model.UserRoleAdmin)
				c.Request = c.Request.WithContext(ctx)
			},
			wantCalled: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/", nil)

			if tt.setupCtx != nil {
				tt.setupCtx(c)
			}

			called := false
			handler := func(ctx *gin.Context, request any) (any, error) {
				called = true
				return "ok", nil
			}

			wrapped := middleware(handler, tt.operationID)
			_, err := wrapped(c, nil)

			if called != tt.wantCalled {
				t.Errorf("handler called = %v, want %v", called, tt.wantCalled)
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if !apperror.SameKind(err, apperror.KindForbidden) {
					t.Errorf("expected forbidden error, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
