package queryexecutor

import (
	"context"
	manager "db-lens/internal/engine/connectionManager"
	"db-lens/internal/engine/entities"
)

type Executor interface {
	Execute(ctx context.Context, conn manager.Connection, query string, args ...any) (*entities.QueryResult, error)
}
