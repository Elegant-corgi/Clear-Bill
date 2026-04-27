//go:build wireinject
// +build wireinject

package app

import (
	"clearbill/mgr/server/internal/app/action"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/config"
	"clearbill/mgr/server/internal/app/dal"
	"clearbill/mgr/server/internal/app/db"
	"github.com/google/wire"
)

func BuildInjector(cfg config.Config) (*Injector, error) {
	wire.Build(
		db.DBSet,
		dal.DalSet,
		bll.BllSet,
		action.ActionsSet,
		InjectorSet,
	)
	return nil, nil
}
