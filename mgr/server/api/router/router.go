package router

import (
	"clearbill/mgr/server/internal/app/action"
	"clearbill/mgr/server/internal/app/config"
	"github.com/gin-gonic/gin"
)

type IRouter interface {
	Register(app *gin.Engine)
	Prefixes() []string
}

type Router struct {
	Config        config.Config
	SystemAction  *action.SystemAction
	BillingAction *action.BillingAction
}

var _ IRouter = (*Router)(nil)

func New(
	cfg config.Config,
	systemAction *action.SystemAction,
	billingAction *action.BillingAction,
) *Router {
	return &Router{
		Config:        cfg,
		SystemAction:  systemAction,
		BillingAction: billingAction,
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
