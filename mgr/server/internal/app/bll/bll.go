package bll

import "github.com/google/wire"

var BllSet = wire.NewSet(
	NewSystemService,
	NewBillingService,
	NewTenantService,
)
