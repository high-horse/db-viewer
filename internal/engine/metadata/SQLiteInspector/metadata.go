package sqliteinspector

import (
	"context"
	"database/sql"
	manager "db-viewer/internal/engine/connectionManager"
	"db-viewer/internal/engine/entities"
	"fmt"
	"log"
	"strings"
)

type SQLiteInspector struct{}

func NewInspector() *SQLiteInspector {
	return &SQLiteInspector{}
}

func (s *SQLiteInspector) ListDatabases(ctx context.Context, conn manager.Connection) ([]entities.DatabaseInfo, error) {
	log.Println("list databases", conn.DatabaseName())
	return []entities.DatabaseInfo{
		{
			Name: conn.DatabaseName(),
		},
	}, nil
}

func (s *SQLiteInspector) ListTables(ctx context.Context, conn manager.Connection) ([]entities.InspectTableInfo, error) {
	sqlConn, ok := conn.(manager.SQLConnection)
	if !ok {
		return nil, fmt.Errorf("connection is not a SQL connection")
	}

	query := `
		SELECT
			name,
			type,
			sql
		FROM sqlite_master
		WHERE type IN ('table','view')
		  AND name NOT LIKE 'sqlite_%'
		ORDER BY name;
	`

	rows, err := sqlConn.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []entities.InspectTableInfo

	for rows.Next() {
		var (
			table     entities.InspectTableInfo
			createSql sql.NullString
		)

		if err := rows.Scan(&table.Name, &table.Type, &createSql); err != nil {
			return nil, err
		}

		table.Database = conn.DatabaseName()
		table.Schema = "main"
		table.Engine = "SQLite"
		if createSql.Valid {
			table.Comment = createSql.String
		}

		var count int64
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", table.Name)
		err := sqlConn.DB().QueryRowContext(ctx, countQuery).Scan(&count)
		if err == nil {
			table.Rows = count
		}

		tables = append(tables, table)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (s *SQLiteInspector) ListColumns(ctx context.Context, conn manager.Connection, table entities.InspectTableInfo) ([]entities.InspectColumnInfo, error) {
	sqlConn, ok := conn.(manager.SQLConnection)
	if !ok {
		return nil, fmt.Errorf("connection is not a SQL connection")
	}

	query := fmt.Sprintf("PRAGMA %s.table_xinfo(%s)", quoteIdentifier(table.Schema), quoteIdentifier(table.Name))
	rows, err := sqlConn.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []entities.InspectColumnInfo

	for rows.Next() {
		var (
			cid      int
			col      entities.InspectColumnInfo
			notNull  int
			pk       int
			hidden   int
			defaultV sql.NullString
		)

		err := rows.Scan(
			&cid,
			&col.Name,
			&col.DatabaseType,
			&notNull,
			&defaultV,
			&pk,
			&hidden,
		)
		if err != nil {
			return nil, err
		}

		col.Nullable = notNull == 0
		col.PrimaryKey = pk > 0
		col.Generated = hidden != 0
		if defaultV.Valid {
			col.DefaultValue = defaultV.String
		}
		// SQLite has no explicit auto_increment flag in PRAGMA.
		// AUTOINCREMENT only applies to INTEGER PRIMARY KEY AUTOINCREMENT.
		col.AutoIncrement = col.PrimaryKey && pk == 1 && strings.EqualFold(col.DatabaseType, "INTEGER") && !strings.Contains(strings.ToUpper(table.Comment), "WITHOUT ROWID")

		columns = append(columns, col)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	keys := 0
	for _, col := range columns {
		if col.PrimaryKey {
			keys++
		}
	}
	if keys > 1 {
		for index := range columns {
			columns[index].AutoIncrement = false
		}
	}
	return columns, nil
}

func quoteIdentifier(name string) string {
	if name == "" {
		name = "main"
	}
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}
