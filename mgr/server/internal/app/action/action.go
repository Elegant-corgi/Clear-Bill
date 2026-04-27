package action

import "github.com/google/wire"

var ActionsSet = wire.NewSet(
	NewSystemAction,
	NewBillingAction,
	NewTenantAction,
)
