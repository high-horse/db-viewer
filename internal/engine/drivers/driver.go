package drivers

import (
	"context"
	manager "db-lens/internal/engine/connectionManager"
	"db-lens/internal/engine/entities"
	"db-lens/internal/engine/metadata"
	queryParaser "db-lens/internal/engine/parser"
	queryexecutor "db-lens/internal/engine/queryExecutor"
	"db-lens/internal/engine/transports"
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
