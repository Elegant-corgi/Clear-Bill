package bll

import "github.com/google/wire"

var BllSet = wire.NewSet(
	NewPermissionCatalog,
	NewAuthService,
	NewSystemService,
	NewBillingService,
	NewTenantService,
	NewRoleService,
	NewUserService,
	NewCredentialService,
	NewAuditService,
)
