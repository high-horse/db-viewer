package mySQLInspector

import (
	"context"
	manager "db-lens/internal/engine/connectionManager"
	"db-lens/internal/engine/entities"
	"fmt"
	"strings"
)

func (i *MySQLInspector) GetTableDDL(ctx context.Context, conn manager.Connection, table entities.TableRef) (string, error) {
	sqlConn, ok := conn.(manager.SQLConnection)
	if !ok {
		return "", fmt.Errorf("connection is not SQL")
	}
	quote := func(s string) string { return "`" + strings.ReplaceAll(s, "`", "``") + "`" }
	database := table.Database
	if database == "" {
		database = table.Schema
	}
	if database == "" {
		database = conn.DatabaseName()
	}
	rows, err := sqlConn.DB().QueryContext(ctx, "SHOW CREATE TABLE "+quote(database)+"."+quote(table.Name))
	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("no DDL found for %s", table.Name)
	}
	values := make([]any, len(columns))
	destinations := make([]any, len(columns))
	for index := range values {
		destinations[index] = &values[index]
	}
	if err := rows.Scan(destinations...); err != nil {
		return "", err
	}
	for index, column := range columns {
		if column != "Create Table" && column != "Create View" {
			continue
		}
		var ddl string
		switch value := values[index].(type) {
		case string:
			ddl = value
		case []byte:
			ddl = string(value)
		}
		if ddl != "" {
			return strings.TrimRight(ddl, "; \n") + ";", nil
		}
	}
	return "", fmt.Errorf("no CREATE statement returned for %s", table.Name)
}
