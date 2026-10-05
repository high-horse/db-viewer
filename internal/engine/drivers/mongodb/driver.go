package mongodb

import (
	"context"
	manager "db-lens/internal/engine/connectionManager"
	"db-lens/internal/engine/entities"
	"db-lens/internal/engine/metadata"
	parser "db-lens/internal/engine/parser"
	executor "db-lens/internal/engine/queryExecutor"
	"db-lens/internal/engine/transports"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type Driver struct{ pager *Pager }

func NewDriver(pager *Pager) *Driver { return &Driver{pager: pager} }
func (d *Driver) Name() string       { return "mongodb" }
func (d *Driver) Create(_ context.Context, config entities.ConnectionConfig, transport transports.Transport) (manager.Connection, error) {
	return New(config, transport), nil
}
func (d *Driver) Parser() parser.Parser         { return commandParser{} }
func (d *Driver) Executor() executor.Executor   { return d }
func (d *Driver) Inspector() metadata.Inspector { return inspector{} }
func (d *Driver) Execute(ctx context.Context, conn manager.Connection, query string, _ ...any) (*entities.QueryResult, error) {
	return d.pager.Open(ctx, conn, query, 100)
}

type commandParser struct{}

func (commandParser) Parse(query string) (*entities.SQLQueryEntity, error) {
	if _, err := parseCommand(query); err != nil {
		return nil, err
	}
	return &entities.SQLQueryEntity{RawSQL: query, StatementType: entities.StatementUnknown}, nil
}

// Extended JSON keeps BSON types and preserves the first command key.
func parseCommand(query string) (bson.D, error) {
	var command bson.D
	if err := bson.UnmarshalExtJSON([]byte(query), false, &command); err != nil {
		return nil, fmt.Errorf("enter a MongoDB command as a JSON object: %w", err)
	}
	if len(command) == 0 {
		return nil, fmt.Errorf("MongoDB command must not be empty")
	}
	switch command[0].Key {
	case "find", "aggregate", "count", "distinct", "listCollections", "listIndexes", "collStats", "dbStats", "ping", "hello", "insert", "update", "delete", "findAndModify", "create", "createIndexes", "drop", "dropIndexes":
	default:
		return nil, fmt.Errorf("unsupported MongoDB command %q", command[0].Key)
	}
	// Live streams do not have a finite result for the grid.
	for _, field := range command {
		if (field.Key == "tailable" || field.Key == "awaitData") && field.Value == true {
			return nil, fmt.Errorf("tailable MongoDB cursors are not supported")
		}
		if field.Key == "pipeline" {
			if stages, ok := field.Value.(bson.A); ok {
				for _, stage := range stages {
					if document, ok := stage.(bson.D); ok {
						for _, entry := range document {
							if entry.Key == "$changeStream" {
								return nil, fmt.Errorf("change streams are not supported")
							}
						}
					}
				}
			}
		}
	}
	return command, nil
}

func readCommand(command bson.D) bool {
	switch command[0].Key {
	case "find", "count", "distinct", "listCollections", "listIndexes", "collStats", "dbStats", "ping", "hello":
		return true
	case "aggregate":
		for _, field := range command {
			if field.Key == "pipeline" {
				stages, ok := field.Value.(bson.A)
				if !ok {
					return false
				}
				for _, stage := range stages {
					document, ok := stage.(bson.D)
					if !ok {
						return false
					}
					for _, entry := range document {
						if entry.Key == "$out" || entry.Key == "$merge" {
							return false
						}
					}
				}
			}
		}
		return true
	}
	return false
}
