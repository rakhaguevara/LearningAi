package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	apperr "github.com/adaptive-ai-learn/backend/internal/common/errors"
	"github.com/adaptive-ai-learn/backend/internal/common/response"
	jwtpkg "github.com/adaptive-ai-learn/backend/pkg/jwt"
)

const (
	AuthHeaderKey = "Authorization"
	BearerPrefix  = "Bearer "
	ContextUserID = "userID"
	ContextEmail  = "email"
	ContextRole   = "role"
)

func Auth(jwtSvc *jwtpkg.Service, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(AuthHeaderKey)
		if header == "" {
			response.Err(c, apperr.NewUnauthorized("missing authorization header"))
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, BearerPrefix) {
			response.Err(c, apperr.NewUnauthorized("invalid authorization format"))
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, BearerPrefix)
		claims, err := jwtSvc.ValidateToken(tokenStr)
		if err != nil {
			log.Debug("token validation failed", zap.Error(err))
			response.Err(c, apperr.NewUnauthorized("invalid or expired token"))
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID.String())
		c.Set(ContextEmail, claims.Email)
		c.Set(ContextRole, claims.Role)
		c.Next()
	}
}

func CORS(frontendURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := frontendURL
		if origin == "" {
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
