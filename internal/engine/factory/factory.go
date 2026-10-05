package factory

import (
	"context"
	manager "db-lens/internal/engine/connectionManager"
	"db-lens/internal/engine/drivers"
	"db-lens/internal/engine/entities"
	"db-lens/internal/engine/transports"
	"fmt"
	"log"
)

type Factory struct {
	drivers map[string]drivers.Driver
}

func New() *Factory {
	return &Factory{
		drivers: make(map[string]drivers.Driver),
	}
}

func (f *Factory) Register(driver drivers.Driver) {
	log.Println("registering:", driver.Name())
	f.drivers[driver.Name()] = driver
}

func (f *Factory) Create(
	ctx context.Context,
	config entities.ConnectionConfig,
	transport transports.Transport,
) (manager.Connection, error) {

	log.Println("databse type", config.Type)
	driver, ok := f.drivers[config.Type]

	if !ok {
		return nil, fmt.Errorf("unsupported database type %v", config.Type)
	}

	return driver.Create(ctx, config, transport)
}

func (f *Factory) Driver(name string) (drivers.Driver, error) {
	driver, ok := f.drivers[name]
	if !ok {
		return nil, fmt.Errorf("driver %q not registered", name)
	}

	return driver, nil
}
