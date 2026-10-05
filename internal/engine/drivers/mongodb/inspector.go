package mongodb

import (
	"context"
	manager "db-lens/internal/engine/connectionManager"
	"db-lens/internal/engine/entities"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type inspector struct{}

func mongoConnection(conn manager.Connection) (manager.NoSQLConnection, error) {
	c, ok := conn.(manager.NoSQLConnection)
	if !ok || c.DB() == nil {
		return nil, fmt.Errorf("active MongoDB connection not available")
	}
	return c, nil
}
func (inspector) ListDatabases(ctx context.Context, conn manager.Connection) ([]entities.DatabaseInfo, error) {
	c, err := mongoConnection(conn)
	if err != nil {
		return nil, err
	}
	names, err := c.Client().ListDatabaseNames(ctx, bson.D{})
	result := make([]entities.DatabaseInfo, 0, len(names))
	for _, name := range names {
		result = append(result, entities.DatabaseInfo{Name: name})
	}
	return result, err
}
func (inspector) ListTables(ctx context.Context, conn manager.Connection) ([]entities.InspectTableInfo, error) {
	c, err := mongoConnection(conn)
	if err != nil {
		return nil, err
	}
	cursor, err := c.DB().ListCollections(ctx, bson.D{}, options.ListCollections().SetNameOnly(true).SetAuthorizedCollections(true))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	result := []entities.InspectTableInfo{}
	for cursor.Next(ctx) {
		var item struct {
			Name string `bson:"name"`
			Type string `bson:"type"`
		}
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		kind := "COLLECTION"
		if item.Type == "view" {
			kind = "VIEW"
		}
		result = append(result, entities.InspectTableInfo{Name: item.Name, Type: kind, Database: c.DatabaseName(), Schema: c.DatabaseName(), Engine: "MongoDB"})
	}
	return result, cursor.Err()
}
func (inspector) ListColumns(context.Context, manager.Connection, entities.InspectTableInfo) ([]entities.InspectColumnInfo, error) {
	return []entities.InspectColumnInfo{{Name: "document", DatabaseType: "Extended JSON"}}, nil
}
func (inspector) GetTableDDL(context.Context, manager.Connection, entities.TableRef) (string, error) {
	return "", fmt.Errorf("MongoDB collections do not have SQL DDL; use listCollections or listIndexes JSON commands")
}
