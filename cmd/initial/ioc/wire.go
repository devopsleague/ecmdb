//go:build wireinject

package ioc

import (
	"github.com/Duke1616/ecmdb/internal/policy"
	"github.com/Duke1616/ecmdb/internal/role"
	"github.com/Duke1616/ecmdb/internal/user"
	"github.com/google/wire"
)

var BaseSet = wire.NewSet(InitMongoDB, InitMySQLDB, InitRedis, InitMQ, InitEtcdClient, InitLdapConfig)

func InitApp() (*App, error) {
	wire.Build(wire.Struct(new(App), "*"),
		BaseSet,
		InitCasbin,
		user.InitModule,
		wire.FieldsOf(new(*user.Module), "Svc"),
		role.InitModule,
		wire.FieldsOf(new(*role.Module), "Svc"),
		policy.InitModule,
	)
	return new(App), nil
}