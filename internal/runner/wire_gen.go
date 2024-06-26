// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package runner

import (
	"context"
	"github.com/Duke1616/ecmdb/internal/codebook"
	"github.com/Duke1616/ecmdb/internal/runner/internal/event"
	"github.com/Duke1616/ecmdb/internal/runner/internal/repository"
	"github.com/Duke1616/ecmdb/internal/runner/internal/repository/dao"
	"github.com/Duke1616/ecmdb/internal/runner/internal/service"
	"github.com/Duke1616/ecmdb/internal/runner/internal/web"
	"github.com/Duke1616/ecmdb/internal/worker"
	"github.com/Duke1616/ecmdb/pkg/mongox"
	"github.com/ecodeclub/mq-api"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitModule(db *mongox.Mongo, q mq.MQ, workerModule *worker.Module, codebookModule *codebook.Module) (*Module, error) {
	runnerDAO := dao.NewRunnerDAO(db)
	runnerRepository := repository.NewRunnerRepository(runnerDAO)
	serviceService := service.NewService(runnerRepository)
	service2 := workerModule.Svc
	service3 := codebookModule.Svc
	handler := web.NewHandler(serviceService, service2, service3)
	taskRunnerConsumer := initTaskRunnerConsumer(serviceService, q, service2, service3)
	module := &Module{
		Svc: serviceService,
		Hdl: handler,
		c:   taskRunnerConsumer,
	}
	return module, nil
}

// wire.go:

var ProviderSet = wire.NewSet(web.NewHandler, service.NewService, repository.NewRunnerRepository, dao.NewRunnerDAO)

func initTaskRunnerConsumer(svc service.Service, mq2 mq.MQ, workerSvc worker.Service, codebookSvc codebook.Service) *event.TaskRunnerConsumer {
	consumer, err := event.NewTaskRunnerConsumer(svc, mq2, workerSvc, codebookSvc)
	if err != nil {
		panic(err)
	}

	consumer.Start(context.Background())
	return consumer
}
