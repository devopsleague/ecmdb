// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package relation

import (
	"github.com/Duke1616/ecmdb/internal/attribute"
	"github.com/Duke1616/ecmdb/internal/relation/internal/repository"
	"github.com/Duke1616/ecmdb/internal/relation/internal/repository/dao"
	"github.com/Duke1616/ecmdb/internal/relation/internal/service"
	"github.com/Duke1616/ecmdb/internal/relation/internal/web"
	"github.com/Duke1616/ecmdb/internal/resource"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

// Injectors from wire.go:

func InitModule(db *mongo.Client, attributeModel *attribute.Module, resourceModel *resource.Module) (*Module, error) {
	relationDAO := dao.NewRelationDAO(db)
	relationRepository := repository.NewRelationRepository(relationDAO)
	serviceService := service.NewService(relationRepository)
	service2 := attributeModel.Svc
	service3 := resourceModel.Svc
	handler := web.NewHandler(serviceService, service2, service3)
	module := &Module{
		Svc: serviceService,
		Hdl: handler,
	}
	return module, nil
}

// wire.go:

var ProviderSet = wire.NewSet(web.NewHandler, service.NewService, repository.NewRelationRepository, dao.NewRelationDAO)
