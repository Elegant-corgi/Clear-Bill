package action

import "github.com/google/wire"

var ActionsSet = wire.NewSet(
	NewAuthAction,
	NewSystemAction,
	NewBillingAction,
	NewTenantAction,
	NewUserAction,
)
