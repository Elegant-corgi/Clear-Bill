package router

import (
	"clearbill/mgr/server/internal/app/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"strings"
)

func (r *Router) RegisterValidator() {
}

func (r *Router) RegisterAPI(app *gin.Engine) {
	app.ForwardedByClientIP = true

	app.GET("/health", r.SystemAction.Health)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	r.registerSystemRoutes(v1)
	r.registerAuthRoutes(v1)
	r.registerBillingRoutes(v1)
	r.registerTenantRoutes(v1)
	r.registerUserRoutes(v1)
	r.registerWebsite(app)
}

func (r *Router) registerSystemRoutes(v1 *gin.RouterGroup) {
	v1.GET("/health", r.SystemAction.Health)
}

func (r *Router) registerBillingRoutes(v1 *gin.RouterGroup) {
	v1.GET("/dashboard/summary", r.BillingAction.DashboardSummary)
	v1.GET("/bills", r.BillingAction.ListBills)
	v1.GET("/customers", r.BillingAction.ListCustomers)
	v1.GET("/reconciliations", r.BillingAction.ListReconciliationTasks)
}

func (r *Router) registerAuthRoutes(v1 *gin.RouterGroup) {
	auth := v1.Group("/auth")
	auth.POST("/login", r.AuthAction.Login)

	protected := auth.Group("")
	protected.Use(middleware.AuthMiddleware(r.AuthService))
	protected.POST("/logout", r.AuthAction.Logout)
	protected.GET("/me", r.AuthAction.CurrentUser)
	protected.PUT("/password", r.AuthAction.ChangeOwnPassword)
	protected.POST("/tokens", r.AuthAction.CreateAPIToken)
}

func (r *Router) registerTenantRoutes(v1 *gin.RouterGroup) {
	tenants := v1.Group("/tenants")
	tenants.Use(middleware.AuthMiddleware(r.AuthService))
	tenants.POST("", r.TenantAction.CreateTenant)
	tenants.GET("", r.TenantAction.ListTenants)
	tenants.GET("/:id", r.TenantAction.GetTenant)
	tenants.PUT("/:id", r.TenantAction.UpdateTenant)
	tenants.DELETE("/:id", r.TenantAction.DeleteTenant)
}

func (r *Router) registerUserRoutes(v1 *gin.RouterGroup) {
	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware(r.AuthService))
	users.POST("", r.UserAction.CreateUser)
	users.GET("", r.UserAction.ListUsers)
	users.GET("/:id", r.UserAction.GetUser)
	users.PUT("/:id", r.UserAction.UpdateUser)
	users.DELETE("/:id", r.UserAction.DeleteUser)
	users.PUT("/:id/password", r.UserAction.ResetPassword)
}

func (r *Router) registerWebsite(app *gin.Engine) {
	indexPath := path.Join(r.Config.WebRoot, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		return
	}

	fileServer := http.FileServer(http.Dir(r.Config.WebRoot))

	app.NoRoute(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/api/") {
			ctx.Status(http.StatusNotFound)
			return
		}

		cleanPath := strings.TrimPrefix(path.Clean(ctx.Request.URL.Path), "/")
		if cleanPath != "" && cleanPath != "." {
			targetPath := path.Join(r.Config.WebRoot, cleanPath)
			if info, err := os.Stat(targetPath); err == nil && !info.IsDir() {
				fileServer.ServeHTTP(ctx.Writer, ctx.Request)
				return
			}
		}

		indexRequest := ctx.Request.Clone(ctx.Request.Context())
		indexRequest.URL.Path = "/index.html"
		fileServer.ServeHTTP(ctx.Writer, indexRequest)
	})
}
