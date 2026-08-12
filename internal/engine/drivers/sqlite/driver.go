package sqlite

import (
	"context"
	manager "db-viewer/internal/engine/connectionManager"
	"db-viewer/internal/engine/entities"
	"db-viewer/internal/engine/metadata"
	sqliteinspector "db-viewer/internal/engine/metadata/SQLiteInspector"
	queryParaser "db-viewer/internal/engine/parser"
	sqliteQueryParser "db-viewer/internal/engine/parser/sqlite"

	queryexecutor "db-viewer/internal/engine/queryExecutor"
	"db-viewer/internal/engine/queryExecutor/sqlExecutor"
	"db-viewer/internal/engine/transports"
)

type Driver struct{
	executor  queryexecutor.Executor
	inspector metadata.Inspector
	parser queryParaser.Parser
}

func NewDriver() *Driver {
	return &Driver{
		executor: sqlExecutor.New(),
		inspector: sqliteinspector.NewInspector(),
		parser: sqliteQueryParser.NewParser(),
	}
}

func (d *Driver) Name() string {
	return "sqlite"
}

func (d *Driver) Create(ctx context.Context, config entities.ConnectionConfig, transport transports.Transport) (manager.Connection, error) {

	conn := New(config, transport)
	return conn, nil
}

func (d *Driver) Executor() queryexecutor.Executor {
	return d.executor
}

func (d *Driver) Inspector() metadata.Inspector {
	return d.inspector
}

func (d *Driver) Parser() queryParaser.Parser {
	return d.parser
}
