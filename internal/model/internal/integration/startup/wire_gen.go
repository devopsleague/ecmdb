// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/Duke1616/ecmdb/internal/model"
	"github.com/Duke1616/ecmdb/internal/model/internal/web"
	"github.com/Duke1616/ecmdb/internal/relation"
)

// Injectors from wire.go:

func InitHandler(rmModule *relation.Module) (*web.Handler, error) {
	mongo := InitMongoDB()
	module, err := model.InitModule(mongo, rmModule)
	if err != nil {
		return nil, err
	}
	handler := module.Hdl
	return handler, nil
}