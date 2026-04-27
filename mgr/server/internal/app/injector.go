package app

import (
	"clearbill/mgr/server/api/router"
	"clearbill/mgr/server/internal/app/action"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/config"
	"clearbill/mgr/server/internal/app/dal"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var InjectorSet = wire.NewSet(
	InitGinEngine,
	InitRouter,
	wire.Struct(new(Injector), "*"),
)

type Injector struct {
	Config         config.Config
	DBClient       *gorm.DB
	Engine         *gin.Engine
	Router         *router.Router
	SystemAction   *action.SystemAction
	BillingAction  *action.BillingAction
	TenantAction   *action.TenantAction
	SystemService  *bll.SystemService
	BillingService *bll.BillingService
	TenantService  *bll.TenantService
	SystemDAL      *dal.SystemDAL
	BillingDAL     *dal.BillingDAL
	TenantDAL      *dal.TenantDAL
}

func InitGinEngine(cfg config.Config) *gin.Engine {
	switch cfg.RunMode {
	case gin.ReleaseMode, gin.TestMode, gin.DebugMode:
		gin.SetMode(cfg.RunMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	return engine
}

func InitRouter(
	cfg config.Config,
	systemAction *action.SystemAction,
	billingAction *action.BillingAction,
	tenantAction *action.TenantAction,
	engine *gin.Engine,
) *router.Router {
	appRouter := router.New(cfg, systemAction, billingAction, tenantAction)
	appRouter.Register(engine)
	return appRouter
}
