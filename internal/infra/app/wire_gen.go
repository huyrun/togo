// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package app

import (
	"context"
	"github.com/manabie-com/togo/internal/infra"
)

// Injectors from wire.go:

func InitApplication(ctx context.Context) (*ApplicationContext, func(), error) {
	appConfig, err := infra.ProvideConfig()
	if err != nil {
		return nil, nil, err
	}
	db, cleanup, err := infra.ProvidePostgres(appConfig)
	if err != nil {
		return nil, nil, err
	}
	userRepo := infra.ProvideUserRepo(db)
	service := infra.ProvideAuthService(appConfig, userRepo)
	eventStore, err := infra.ProvideEventStore(db)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userConfigRepo := infra.ProvideUserConfigRepo(db)
	userTaskRepo := infra.ProvideUserTaskRepo(db)
	userTaskProjector := infra.ProvideUserTaskProjector(userConfigRepo, userTaskRepo)
	taskRepo := infra.ProvideTaskRepo(db)
	taskHandler := infra.ProvideTaskHandler(taskRepo)
	eventBus := infra.ProvideEventBus(userTaskProjector, taskHandler)
	aggregateStore, err := infra.ProvideAggregateStore(eventStore, eventBus)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userTaskCommandHandler, err := infra.ProvideUserTaskCommandHandler(aggregateStore)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	commandHandler := infra.ProvideCommandbus(userTaskCommandHandler)
	user_tasksService := infra.ProvideUserTaskService(commandHandler, userConfigRepo)
	dbSlave, cleanup2, err := infra.ProvidePostgresSlave(appConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	readRepo := infra.ProvideReadRepo(dbSlave)
	restAPIHandler := infra.ProvideRestAPIHandler(service, user_tasksService, readRepo)
	restService, cleanup3, err := infra.ProvideRestService(appConfig, restAPIHandler)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	applicationContext := &ApplicationContext{
		ctx:     ctx,
		cfg:     appConfig,
		restSrv: restService,
		authSrv: service,
	}
	return applicationContext, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}