package middleware

import (
	"strings"

	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/dal/dbmodel"
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
			c.AbortWithStatusJSON(401, gin.H{"success": false, "error": "missing session"})
			return
		}

		user, err := authService.ValidateSession(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"success": false, "error": "invalid session"})
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
			c.AbortWithStatusJSON(401, gin.H{"success": false, "error": "missing api token"})
			return
		}

		user, err := authService.ValidateAPIToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"success": false, "error": "invalid api token"})
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
			c.AbortWithStatusJSON(401, gin.H{"success": false, "error": "unauthorized"})
			return
		}

		permission, ok := roleService.ResolvePermissionByRoute(c.Request.Method, c.FullPath())
		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"success": false, "error": "permission not registered"})
			return
		}

		allowed, err := roleService.HasPermission(c.Request.Context(), user, permission.ID)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"success": false, "error": err.Error()})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(403, gin.H{"success": false, "error": "permission denied"})
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
