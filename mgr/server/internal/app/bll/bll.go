package bll

import "github.com/google/wire"

var BllSet = wire.NewSet(
	NewAuthService,
	NewSystemService,
	NewBillingService,
	NewTenantService,
	NewUserService,
)
