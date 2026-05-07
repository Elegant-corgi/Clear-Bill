package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"clearbill/mgr/server/internal/app/ginx"
	"github.com/gin-gonic/gin"
)

const (
	currentUserKey   = "currentUser"
	currentTokenKey  = "currentToken"
	auditUserKey     = "auditUser"
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

func SetAuditUser(c *gin.Context, user string) {
	c.Set(auditUserKey, strings.TrimSpace(user))
}

func AuditMiddleware(auditService *bll.AuditService, roleService *bll.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		start := time.Now().UTC()
		c.Next()

		user := auditUser(c)
		operation := auditOperation(c, roleService)
		resource := auditResource(c)
		result := auditResult(c)
		_ = auditService.Record(c.Request.Context(), user, operation, resource, result, start)
	}
}

func auditUser(c *gin.Context) string {
	if user := CurrentUser(c); user != nil && strings.TrimSpace(user.Username) != "" {
		return user.Username
	}
	if value, ok := c.Get(auditUserKey); ok {
		if user, _ := value.(string); strings.TrimSpace(user) != "" {
			return user
		}
	}
	return "anonymous"
}

func auditOperation(c *gin.Context, roleService *bll.RoleService) string {
	if roleService != nil {
		if permission, ok := roleService.ResolvePermissionByRoute(c.Request.Method, c.FullPath()); ok {
			return permission.ID
		}
	}
	fullPath := strings.TrimSpace(c.FullPath())
	if fullPath == "" {
		fullPath = c.Request.URL.Path
	}
	return strings.ToLower(strings.TrimSpace(c.Request.Method)) + ":" + fullPath
}

func auditResource(c *gin.Context) string {
	if fullPath := strings.TrimSpace(c.FullPath()); fullPath != "" {
		return fullPath
	}
	return c.Request.URL.Path
}

func auditResult(c *gin.Context) string {
	status := c.Writer.Status()
	if status >= http.StatusBadRequest {
		return "failed:" + strconv.Itoa(status)
	}
	return "success"
}
