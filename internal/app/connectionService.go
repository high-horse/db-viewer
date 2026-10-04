package app

import (
	"db-viewer/internal/db"
	"db-viewer/internal/engine/entities"
	"db-viewer/internal/types"
)

type DatabaseService struct{}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

func (s *DatabaseService) GetConnections() ([]types.Connection, error) {
	return db.GetConnectionList()
}

// func (s *DatabaseService) CreateConnection(connection types.Connection) error {
// 	return db.CreateConnection(connection)
// }

func (s *DatabaseService) DeleteConnection(id int) error {
	return db.DeleteConnection(id)
}

func (s *DatabaseService) UpdateConnection(config entities.ConnectionConfig) error {
	return db.UpdateConnection(config)
}
