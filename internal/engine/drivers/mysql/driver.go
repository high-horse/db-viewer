package mysql

import (
	"context"
	manager "db-lens/internal/engine/connectionManager"
	"db-lens/internal/engine/entities"
	"db-lens/internal/engine/metadata"
	"db-lens/internal/engine/metadata/mySQLInspector"
	queryParaser "db-lens/internal/engine/parser"
	mysqlQueryParser "db-lens/internal/engine/parser/mysql"

	queryexecutor "db-lens/internal/engine/queryExecutor"
	"db-lens/internal/engine/queryExecutor/sqlExecutor"
	"db-lens/internal/engine/transports"
)

type Driver struct {
	executor  queryexecutor.Executor
	inspector metadata.Inspector
	parser    queryParaser.Parser
}

func NewDriver() *Driver {
	return &Driver{
		// Instantiate once during driver creation
		executor:  sqlExecutor.New(),
		inspector: mySQLInspector.NewInspector(),
		parser:    mysqlQueryParser.NewParser(),
	}
}

func (d *Driver) Name() string {
	return "mysql"
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
