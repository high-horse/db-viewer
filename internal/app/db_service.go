package app

import (
	"context"
	"db-viewer/internal/db"
	manager "db-viewer/internal/engine/connectionManager"
	"db-viewer/internal/engine/drivers/mysql"
	"db-viewer/internal/engine/drivers/postgres"
	"db-viewer/internal/engine/drivers/sqlite"
	"db-viewer/internal/engine/entities"
	"db-viewer/internal/engine/factory"
	"db-viewer/internal/engine/transports"
	"fmt"
	"time"
)

type DbService struct {
	factory     *factory.Factory
	manager     *manager.ConnectionManager
	historyRepo *db.HistoryRepository
}

func NewDbService(historyRepo *db.HistoryRepository) *DbService {
	f := factory.New()
	f.Register(mysql.NewDriver())
	f.Register(postgres.NewDriver())
	f.Register(sqlite.NewDriver())
	return &DbService{
		factory:     f,
		manager:     manager.NewConnectionManager(),
		historyRepo: historyRepo,
	}
}

func (d *DbService) Connect(
	ctx context.Context,
	config entities.ConnectionConfig,
) error {

	transport := transports.NewDirect(config.Host, config.Port)

	conn, err := d.factory.Create(context.Background(), config, transport)
	if err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}

	d.manager.Add(conn)

	return conn.Connect(ctx)
}

func (s *DbService) InspectDatabase(ctx context.Context, connID string) ([]entities.InspectTableInfo, error) {
	conn, ok := s.manager.Get(connID)
	if !ok {
		return nil, fmt.Errorf("connection not found")
	}

	driver, err := s.factory.Driver(conn.Type())
	if err != nil {
		return nil, err
	}

	return driver.Inspector().ListTables(ctx, conn)
}

func (s *DbService) ExecuteQuery(ctx context.Context, connID string, rawQuery string) (*entities.QueryResult, error) {
	conn, ok := s.manager.Get(connID)
	if !ok {
		return nil, fmt.Errorf("connection not found")
	}

	driver, err := s.factory.Driver(conn.Type())
	if err != nil {
		return nil, err
	}

	result, err := driver.Executor().Execute(ctx, conn, rawQuery)
	historyEntry := db.QueryHistoryEntity{
		ConnectionId: connID,
		DatabaseName: conn.DatabaseName(),
		QueryText:    rawQuery,
		Status:       "SUCCESS",
	}
	if err != nil {
		historyEntry.Status = "ERROR"
		return nil, fmt.Errorf("SQL execution evaluation error: %w", err)
	}
	historyEntry.Duration = int(result.Duration)
	go func() {
		s.historyRepo.Log(ctx, historyEntry)
	}()
	return result, nil
}

func (s *DbService) Disconnect(ctx context.Context, connID string) error {
	conn, ok := s.manager.Get(connID)
	if !ok {
		return nil
	}

	_ = conn.Disconnect()
	s.manager.Remove(connID) // Make sure to implement a Remove method inside your engine manager wrapper
	return nil
}

func (s *DbService) GetQueryHistory(
	ctx context.Context, limit int, since time.Time,
) ([]db.QueryHistoryEntity, error) {
	return s.historyRepo.GetHistory(ctx, limit, since)
}
