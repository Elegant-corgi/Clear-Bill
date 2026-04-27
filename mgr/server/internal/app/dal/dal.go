package dal

import "github.com/google/wire"

var DalSet = wire.NewSet(
	NewSystemDAL,
	NewBillingDAL,
	NewTenantDAL,
	NewUserDAL,
	NewSessionDAL,
	NewAPITokenDAL,
)
