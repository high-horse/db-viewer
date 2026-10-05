package sqliteinspector

import (
	"context"
	manager "db-lens/internal/engine/connectionManager"
	"db-lens/internal/engine/entities"
	"fmt"
	"strings"
)

func (p *SQLiteInspector) GetTableDDL(ctx context.Context, conn manager.Connection, table entities.TableRef) (string, error) {
	sqlConn, ok := conn.(manager.SQLConnection)
	if !ok {
		return "", fmt.Errorf("connection is not SQL")
	}
	schema := table.Schema
	if schema == "" {
		schema = "main"
	}
	rows, err := sqlConn.DB().QueryContext(ctx, `SELECT sql FROM `+quoteIdentifier(schema)+`.sqlite_schema WHERE (name = ? OR tbl_name = ?) AND sql IS NOT NULL ORDER BY CASE type WHEN 'table' THEN 0 WHEN 'view' THEN 0 WHEN 'index' THEN 1 ELSE 2 END, name`, table.Name, table.Name)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var statements []string
	for rows.Next() {
		var statement string
		if err := rows.Scan(&statement); err != nil {
			return "", err
		}
		statements = append(statements, strings.TrimRight(statement, "; \n")+";")
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	if len(statements) == 0 {
		return "", fmt.Errorf("no DDL found for %s", table.Name)
	}
	return strings.Join(statements, "\n\n"), nil
}
