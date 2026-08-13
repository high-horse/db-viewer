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
	"strconv"
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

func (s *DbService) Connect( ctx context.Context, config entities.ConnectionConfig) (bool, error) {

	transport := transports.NewDirect( config.Host, config.Port)

	conn, err := s.factory.Create(
		ctx,
		config,
		transport,
	)
	if err != nil {
		return false, fmt.Errorf("failed to create connection: %w", err)
	}

	// First test the new connection.
	if err := conn.Connect(ctx); err != nil {
		return false, fmt.Errorf("failed to connect: %w", err)
	}

	if old, ok := s.manager.Active(); ok {
		_ = s.manager.Remove(old.ID())
	}

	if err := s.manager.Add(conn); err != nil {
		_ = conn.Disconnect()
		return false, fmt.Errorf("failed to register connection: %w", err)
	}

	if err := s.manager.SetActive(conn.ID()); err != nil {
		_ = s.manager.Remove(conn.ID())
		return false, fmt.Errorf("failed to set active connection: %w", err)
	}
	
	return true, nil
}

func (s *DbService) Disconnect( ctx context.Context, connID string) error {
	return s.manager.Remove(connID)
}

func (s *DbService) InspectDatabase( ctx context.Context) ([]entities.InspectTableInfo, error) {
	conn, ok := s.manager.Active()
	if !ok {
		return nil, fmt.Errorf("active connection not found")
	}

	driver, err := s.factory.Driver(conn.Type())
	if err != nil {
		return nil, err
	}

	return driver.Inspector().ListTables(ctx, conn)
}

func (s *DbService) ExecuteQuery( ctx context.Context, rawQuery string) (*entities.QueryResult, error) {
	conn, ok := s.manager.Active()
	if !ok {
		return nil, fmt.Errorf("active connection not found")
	}

	driver, err := s.factory.Driver(conn.Type())
	if err != nil {
		return nil, err
	}

	parsed, err := driver.Parser().Parse(rawQuery)
	if err != nil {
		return nil, fmt.Errorf("query parsing error: %w", err)
	}

	result, err := driver.Executor().Execute(
		ctx,
		conn,
		rawQuery,
	)

	historyEntry := db.QueryHistoryEntity{
		ConnectionId: conn.ID(),
		DatabaseName: conn.DatabaseName(),
		QueryText:    rawQuery,
		Status:       "SUCCESS",
	}

	if err != nil {
		historyEntry.Status = "ERROR"
		return nil, fmt.Errorf("SQL execution evaluation error: %w", err)
	}

	_ = parsed
	historyEntry.Duration = int(result.Duration)

	go func() {
		_ = s.historyRepo.Log(ctx, historyEntry)
	}()

	return result, nil
}

func (s *DbService) PingConnection(ctx context.Context, connID string) (bool, error) {
	conn, ok := s.manager.Get(connID)
	if !ok {
		return false, fmt.Errorf("connection not found")
	}

	if err := conn.Ping(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func (s *DbService) PingConfig(ctx context.Context, config entities.ConnectionConfig) (bool, error) {

	transport := transports.NewDirect(config.Host, config.Port)

	conn, err := s.factory.Create(
		ctx,
		config,
		transport,
	)

	if err != nil {
		return false, err
	}

	if err := conn.Connect(ctx); err != nil {
		return false, err
	}
	defer conn.Disconnect()

	return true, conn.Ping(ctx)
}

func (s *DbService) GetQueryHistory(ctx context.Context, limit int, since time.Time) ([]db.QueryHistoryEntity, error) {
	return s.historyRepo.GetHistory(ctx, limit, since)
}


func (s *DbService) SaveAndConnect(ctx context.Context, config entities.ConnectionConfig) (bool, error) {
	newID, err := db.StoreConnection(config)
	if err != nil {
		return false, fmt.Errorf("failed to store connection: %w", err)
	}
	config.ID = strconv.FormatInt(newID, 10)
	return s.Connect(ctx, config)
}

func (s *DbService) GetActiveConnection() (string, error) {
	conn, ok := s.manager.Active()
	if !ok {
		return "", fmt.Errorf("no active connection")
	}
	return conn.ID(), nil
}


