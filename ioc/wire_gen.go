// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"github.com/Duke1616/ecmdb/internal/attribute"
	"github.com/Duke1616/ecmdb/internal/model"
	"github.com/Duke1616/ecmdb/internal/resource"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitApp() (*App, error) {
	v := InitGinMiddlewares()
	client := InitMongoDB()
	handler := model.InitHandler(client)
	webHandler := attribute.InitHandler(client)
	handler2 := resource.InitHandler(client)
	engine := InitWebServer(v, handler, webHandler, handler2)
	app := &App{
		Web: engine,
	}
	return app, nil
}

// wire.go:

var BaseSet = wire.NewSet(InitMongoDB)
