package middleware

import (
	"net/http"
	"strings"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/shared/requestctx"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

const (
	headerXUserID  = "X-User-ID"
	headerXRobotID = "X-Robot-ID"
	headerXAPIKey  = "X-API-Key"
)

func isRobotAPIPath(path string) bool {
	return strings.HasPrefix(path, "/api/robot/") || path == "/api/robot"
}

// Auth validates requests using either API Key or X-User-ID header.
// For robot API paths, X-API-Key is tried first, falling back to X-User-ID + X-Robot-ID headers.
// For user API paths, X-User-ID is required.
func Auth(userUC usecase.UserUsecase, robotUC usecase.RobotUsecase, apiKeyUC usecase.APIKeyUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isRobotAPIPath(c.Request.URL.Path) {
			robotAuth(c, userUC, robotUC, apiKeyUC)
		} else {
			userAuth(c, userUC)
		}
	}
}

// robotAuth handles authentication for robot API paths (/api/robot/*).
// It tries X-API-Key first, then falls back to X-User-ID + X-Robot-ID headers.
func robotAuth(c *gin.Context, userUC usecase.UserUsecase, robotUC usecase.RobotUsecase, apiKeyUC usecase.APIKeyUsecase) {
	rawKey := c.GetHeader(headerXAPIKey)
	if rawKey != "" {
		robotAuthByAPIKey(c, apiKeyUC, rawKey)
		return
	}

	// Fallback to X-User-ID + X-Robot-ID headers
	robotAuthByHeaders(c, userUC, robotUC)
}

// robotAuthByAPIKey authenticates using an API key from X-API-Key header.
func robotAuthByAPIKey(c *gin.Context, apiKeyUC usecase.APIKeyUsecase, rawKey string) {
	apiKeyAuth, err := apiKeyUC.Authenticate(c.Request.Context(), rawKey)
	if err != nil {
		if apperror.SameKind(err, apperror.KindNotFound) || apperror.SameKind(err, apperror.KindUnauthorized) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to verify API key",
			})
		}
		return
	}

	ctx := c.Request.Context()
	ctx = requestctx.SetUserID(ctx, apiKeyAuth.UserID)
	ctx = requestctx.SetRobotID(ctx, apiKeyAuth.RobotID)
	ctx = requestctx.SetOrganizationID(ctx, apiKeyAuth.OrganizationID)
	ctx = requestctx.SetUserRole(ctx, apiKeyAuth.UserRole)
	c.Request = c.Request.WithContext(ctx)

	// Debounced update of last_used_at (MarkUsed spawns its own goroutine internally)
	apiKeyUC.MarkUsed(apiKeyAuth.APIKeyID)

	c.Next()
}

// robotAuthByHeaders authenticates using X-User-ID and X-Robot-ID headers.
func robotAuthByHeaders(c *gin.Context, userUC usecase.UserUsecase, robotUC usecase.RobotUsecase) {
	userID := c.GetHeader(headerXUserID)
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "X-API-Key or X-User-ID header is required",
		})
		return
	}

	appUser, err := userUC.GetByNaturalID(c.Request.Context(), userID)
	if err != nil {
		if apperror.SameKind(err, apperror.KindNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to verify user",
			})
		}
		return
	}

	robotID := c.GetHeader(headerXRobotID)
	if robotID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "X-Robot-ID header is required",
		})
		return
	}

	rob, err := robotUC.GetByID(c.Request.Context(), robotID)
	if err != nil {
		if apperror.SameKind(err, apperror.KindNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Robot not found",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to verify robot",
			})
		}
		return
	}

	if appUser.OrganizationID != rob.OrganizationID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "User and robot do not belong to the same organization",
		})
		return
	}

	ctx := c.Request.Context()
	ctx = requestctx.SetUserID(ctx, userID)
	ctx = requestctx.SetRobotID(ctx, robotID)
	ctx = requestctx.SetOrganizationID(ctx, rob.OrganizationID)
	ctx = requestctx.SetUserRole(ctx, appUser.Role)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

// userAuth handles authentication for non-robot API paths using X-User-ID header.
func userAuth(c *gin.Context, userUC usecase.UserUsecase) {
	userID := c.GetHeader(headerXUserID)
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "X-User-ID header is required",
		})
		return
	}

	appUser, err := userUC.GetByNaturalID(c.Request.Context(), userID)
	if err != nil {
		if apperror.SameKind(err, apperror.KindNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to verify user",
			})
		}
		return
	}

	ctx := c.Request.Context()
	ctx = requestctx.SetUserID(ctx, appUser.IDNatural)
	ctx = requestctx.SetOrganizationID(ctx, appUser.OrganizationID)
	ctx = requestctx.SetUserRole(ctx, appUser.Role)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
