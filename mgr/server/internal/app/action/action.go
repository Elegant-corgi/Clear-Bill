package action

import (
	"clearbill/mgr/server/internal/app/bll"
	"github.com/google/wire"
)

var ActionsSet = wire.NewSet(
	NewAuthAction,
	NewSystemAction,
	NewBillingAction,
	NewTenantActionProvider,
	NewRoleAction,
	NewUserAction,
	NewCredentialAction,
)

func NewTenantActionProvider(
	tenantService *bll.TenantService,
	_ *bll.RoleService,
) *TenantAction {
	return NewTenantAction(tenantService)
}
