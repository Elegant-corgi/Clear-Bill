package router

import (
	"clearbill/mgr/server/internal/app/action"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/config"
	"github.com/gin-gonic/gin"
)

type IRouter interface {
	Register(app *gin.Engine)
	Prefixes() []string
}

type Router struct {
	AuthService   *bll.AuthService
	AuthAction    *action.AuthAction
	Config        config.Config
	SystemAction  *action.SystemAction
	BillingAction *action.BillingAction
	TenantAction  *action.TenantAction
	UserAction    *action.UserAction
}

var _ IRouter = (*Router)(nil)

func New(
	cfg config.Config,
	authService *bll.AuthService,
	authAction *action.AuthAction,
	systemAction *action.SystemAction,
	billingAction *action.BillingAction,
	tenantAction *action.TenantAction,
	userAction *action.UserAction,
) *Router {
	return &Router{
		AuthService:   authService,
		AuthAction:    authAction,
		Config:        cfg,
		SystemAction:  systemAction,
		BillingAction: billingAction,
		TenantAction:  tenantAction,
		UserAction:    userAction,
	}
}

func (r *Router) Register(app *gin.Engine) {
	r.RegisterValidator()
	r.RegisterAPI(app)
}

func (r *Router) Prefixes() []string {
	return []string{
		"/api",
	}
}
