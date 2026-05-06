package middleware

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"

	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/ginx"
	"github.com/gin-gonic/gin"
)

const (
	accessKeyHeader = "X-Access-Key"
	secretKeyHeader = "X-Secret-Key"
	timestampHeader = "X-Timestamp"
	signatureHeader = "X-Signature"
)

func ExternalAuthMiddleware(authService *bll.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := extractAPIToken(c); token != "" {
			user, err := authService.ValidateCredentialToken(c.Request.Context(), token)
			if err != nil {
				ginx.ResError(c, errors.New("invalid api token"), 401)
				return
			}
			c.Set(currentUserKey, user)
			c.Set(currentTokenKey, token)
			c.Next()
			return
		}

		accessKey := strings.TrimSpace(c.GetHeader(accessKeyHeader))
		if accessKey != "" {
			if signature := strings.TrimSpace(c.GetHeader(signatureHeader)); signature != "" {
				timestamp, err := strconv.ParseInt(strings.TrimSpace(c.GetHeader(timestampHeader)), 10, 64)
				if err != nil {
					ginx.ResError(c, errors.New("invalid credential signature"), 401)
					return
				}
				body, err := io.ReadAll(c.Request.Body)
				if err != nil {
					ginx.ResError(c, errors.New("invalid credential signature"), 401)
					return
				}
				c.Request.Body = io.NopCloser(bytes.NewReader(body))
				user, err := authService.ValidateAccessKeySignature(
					c.Request.Context(),
					accessKey,
					timestamp,
					signature,
					c.Request.Method,
					c.Request.URL.Path,
					c.Request.URL.RawQuery,
					body,
				)
				if err != nil {
					ginx.ResError(c, errors.New("invalid credential signature"), 401)
					return
				}
				c.Set(currentUserKey, user)
				c.Set(currentTokenKey, accessKey)
				c.Next()
				return
			}

			secretKey := strings.TrimSpace(c.GetHeader(secretKeyHeader))
			if secretKey == "" {
				ginx.ResError(c, errors.New("missing credential secret"), 401)
				return
			}
			user, err := authService.ValidateAccessKeySecret(c.Request.Context(), accessKey, secretKey)
			if err != nil {
				ginx.ResError(c, errors.New("invalid credential"), 401)
				return
			}
			c.Set(currentUserKey, user)
			c.Set(currentTokenKey, accessKey)
			c.Next()
			return
		}

		if token := extractSessionToken(c); token != "" {
			user, err := authService.ValidateSession(c.Request.Context(), token)
			if err != nil {
				ginx.ResError(c, errors.New("invalid session"), 401)
				return
			}
			c.Set(currentUserKey, user)
			c.Set(currentTokenKey, token)
			c.Next()
			return
		}

		ginx.ResError(c, errors.New("missing credentials"), 401)
	}
}

func extractAPIToken(c *gin.Context) string {
	token := strings.TrimSpace(c.GetHeader("X-API-Token"))
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			token = strings.TrimSpace(authHeader[7:])
		}
	}
	return token
}

func extractSessionToken(c *gin.Context) string {
	token, err := c.Cookie(sessionCookieKey)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(token)
}
