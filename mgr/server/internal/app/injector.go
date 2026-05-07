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
	Config            config.Config
	DBClient          *gorm.DB
	Engine            *gin.Engine
	Router            *router.Router
	APITokenDAL       *dal.APITokenDAL
	AuthAction        *action.AuthAction
	CredentialAction  *action.CredentialAction
	AuditAction       *action.AuditAction
	SystemAction      *action.SystemAction
	BillingAction     *action.BillingAction
	TenantAction      *action.TenantAction
	RoleAction        *action.RoleAction
	UserAction        *action.UserAction
	AuthService       *bll.AuthService
	RoleService       *bll.RoleService
	SystemService     *bll.SystemService
	BillingService    *bll.BillingService
	TenantService     *bll.TenantService
	UserService       *bll.UserService
	SystemDAL         *dal.SystemDAL
	BillingDAL        *dal.BillingDAL
	TenantDAL         *dal.TenantDAL
	RoleDAL           *dal.RoleDAL
	RolePermissionDAL *dal.RolePermissionDAL
	UserDAL           *dal.UserDAL
	SessionDAL        *dal.SessionDAL
	CredentialDAL     *dal.CredentialDAL
	CredentialService *bll.CredentialService
	AuditDAL          *dal.AuditDAL
	AuditService      *bll.AuditService
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
	authService *bll.AuthService,
	roleService *bll.RoleService,
	authAction *action.AuthAction,
	credentialAction *action.CredentialAction,
	auditAction *action.AuditAction,
	systemAction *action.SystemAction,
	billingAction *action.BillingAction,
	tenantAction *action.TenantAction,
	roleAction *action.RoleAction,
	userAction *action.UserAction,
	engine *gin.Engine,
) *router.Router {
	appRouter := router.New(cfg, authService, roleService, authAction, credentialAction, auditAction, systemAction, billingAction, tenantAction, roleAction, userAction)
	appRouter.Register(engine)
	return appRouter
}
