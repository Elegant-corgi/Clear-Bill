package middleware

import (
	"errors"
	"strings"

	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"clearbill/mgr/server/internal/app/ginx"
	"github.com/gin-gonic/gin"
)

const (
	currentUserKey   = "currentUser"
	currentTokenKey  = "currentToken"
	sessionCookieKey = "clear_bill_session"
)

func AuthMiddleware(authService *bll.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(sessionCookieKey)
		if err != nil || strings.TrimSpace(token) == "" {
			ginx.ResError(c, errors.New("missing session"), 401)
			return
		}

		user, err := authService.ValidateSession(c.Request.Context(), token)
		if err != nil {
			ginx.ResError(c, errors.New("invalid session"), 401)
			return
		}

		c.Set(currentUserKey, user)
		c.Set(currentTokenKey, token)
		c.Next()
	}
}

func APITokenMiddleware(authService *bll.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(c.GetHeader("X-API-Token"))
		if token == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				token = strings.TrimSpace(authHeader[7:])
			}
		}
		if token == "" {
			ginx.ResError(c, errors.New("missing api token"), 401)
			return
		}

		user, err := authService.ValidateAPIToken(c.Request.Context(), token)
		if err != nil {
			ginx.ResError(c, errors.New("invalid api token"), 401)
			return
		}

		c.Set(currentUserKey, user)
		c.Set(currentTokenKey, token)
		c.Next()
	}
}

func RBACMiddleware(roleService *bll.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := CurrentUser(c)
		if user == nil {
			ginx.ResError(c, errors.New("unauthorized"), 401)
			return
		}

		permission, ok := roleService.ResolvePermissionByRoute(c.Request.Method, c.FullPath())
		if !ok {
			ginx.ResError(c, errors.New("permission not registered"), 403)
			return
		}

		allowed, err := roleService.HasPermission(c.Request.Context(), user, permission.ID)
		if err != nil {
			ginx.ResError(c, err, 500)
			return
		}
		if !allowed {
			ginx.ResError(c, errors.New("permission denied"), 403)
			return
		}

		c.Next()
	}
}

func CurrentUser(c *gin.Context) *dbmodel.User {
	value, ok := c.Get(currentUserKey)
	if !ok {
		return nil
	}
	user, _ := value.(*dbmodel.User)
	return user
}

func CurrentToken(c *gin.Context) string {
	value, ok := c.Get(currentTokenKey)
	if !ok {
		return ""
	}
	token, _ := value.(string)
	return token
}
