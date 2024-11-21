// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package rota

import (
	"fmt"
	"github.com/Duke1616/ecmdb/internal/rota/internal/repository"
	"github.com/Duke1616/ecmdb/internal/rota/internal/repository/dao"
	"github.com/Duke1616/ecmdb/internal/rota/internal/service"
	"github.com/Duke1616/ecmdb/internal/rota/internal/service/schedule"
	"github.com/Duke1616/ecmdb/internal/rota/internal/web"
	"github.com/Duke1616/ecmdb/pkg/mongox"
	"github.com/google/wire"
	"time"
)

// Injectors from wire.go:

func InitModule(db *mongox.Mongo) (*Module, error) {
	rotaDao := dao.NewRotaDao(db)
	rotaRepository := repository.NewRotaRepository(rotaDao)
	scheduler := InitScheduleRule()
	serviceService := service.NewService(rotaRepository, scheduler)
	handler := web.NewHandler(serviceService)
	module := &Module{
		Hdl: handler,
		Svc: serviceService,
	}
	return module, nil
}

// wire.go:

var ProviderSet = wire.NewSet(web.NewHandler, service.NewService, repository.NewRotaRepository, dao.NewRotaDao)

func InitScheduleRule() schedule.Scheduler {

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Print()
	}

	return schedule.NewRruleSchedule(location)
}
