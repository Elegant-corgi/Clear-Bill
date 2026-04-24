package router

import (
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
	r.registerBillingRoutes(v1)
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
