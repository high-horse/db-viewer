package drivers

import (
	"context"
	manager "db-viewer/internal/engine/connectionManager"
	"db-viewer/internal/engine/entities"
	"db-viewer/internal/engine/metadata"
	queryParaser "db-viewer/internal/engine/parser"
	queryexecutor "db-viewer/internal/engine/queryExecutor"
	"db-viewer/internal/engine/transports"
)


type Driver interface {
	Name() string

	Create(
		ctx context.Context,
		config entities.ConnectionConfig,
		transport transports.Transport,
	) (manager.Connection, error)

	Parser() queryParaser.Parser

	Executor() queryexecutor.Executor

	Inspector() metadata.Inspector
}