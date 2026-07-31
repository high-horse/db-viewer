package pgxInspector

import (
	"context"
	"database/sql"
	manager "db-viewer/internal/engine/connectionManager"
	"db-viewer/internal/engine/entities"
	"fmt"
)

type PostgresInspector struct{}

func NewInspector() *PostgresInspector {
	return &PostgresInspector{}
}

func (i *PostgresInspector) ListDatabases(ctx context.Context, conn manager.Connection) ([]entities.DatabaseInfo, error) {
	sqlConn, ok := conn.(manager.SQLConnection)
	if !ok {
		return nil, fmt.Errorf("connection is not sql")
	}

	q := `
		SELECT datname
		FROM pg_database
		where datistemplate = false
		order by datname
	`
	rows, err := sqlConn.DB().QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []entities.DatabaseInfo

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		databases = append(databases, entities.DatabaseInfo{Name: name})
	}
	return databases, nil
}

func (i *PostgresInspector) ListTables(ctx context.Context, conn manager.Connection) ([]entities.InspectTableInfo, error) {
	sqlConn, ok := conn.(manager.SQLConnection)
	if !ok {
		return nil, fmt.Errorf("connection is not sql")
	}
	query := `
		SELECT
			t.table_name,
			CASE
				WHEN t.table_type = 'BASE TABLE' THEN 'TABLE'
				ELSE t.table_type
			END,
			current_database(),
			t.table_schema,
			c.reltuples::bigint,
			pg_catalog.obj_description(c.oid),
			NULL,
			NULL
		FROM information_schema.tables t
		JOIN pg_class c
			ON c.relname = t.table_name
		JOIN pg_namespace n
			ON n.oid = c.relnamespace
			AND n.nspname = t.table_schema
		WHERE t.table_schema NOT IN (
			'pg_catalog',
			'information_schema'
		)
		ORDER BY
			t.table_schema,
			t.table_name
	`
	rows, err := sqlConn.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []entities.InspectTableInfo

	for rows.Next() {

		var (
			table entities.InspectTableInfo

			rowsCnt sql.NullInt64
			comment sql.NullString
			created sql.NullTime
			updated sql.NullTime
		)

		err := rows.Scan(
			&table.Name,
			&table.Type,
			&table.Database,
			&table.Schema,
			&rowsCnt,
			&comment,
			&created,
			&updated,
		)
		if err != nil {
			return nil, err
		}

		if rowsCnt.Valid {
			table.Rows = rowsCnt.Int64
		}

		if comment.Valid {
			table.Comment = comment.String
		}

		if created.Valid {
			table.CreatedAt = &created.Time
		}

		if updated.Valid {
			table.UpdatedAt = &updated.Time
		}

		// PostgreSQL doesn't have a storage engine.
		table.Engine = ""

		tables = append(tables, table)
	}

	return tables, nil
}

func (i *PostgresInspector) ListColumns(
	ctx context.Context,
	conn manager.Connection,
	table entities.InspectTableInfo,
) ([]entities.InspectColumnInfo, error) {

	sqlConn, ok := conn.(manager.SQLConnection)
	if !ok {
		return nil, fmt.Errorf("connection is not SQL")
	}

	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default,
			COALESCE(
				(
					SELECT tc.constraint_type
					FROM information_schema.table_constraints tc
					JOIN information_schema.key_column_usage kcu
						ON tc.constraint_name = kcu.constraint_name
						AND tc.table_schema = kcu.table_schema
					WHERE tc.table_name = c.table_name
					AND tc.table_schema = c.table_schema
					AND kcu.column_name = c.column_name
					AND tc.constraint_type = 'PRIMARY KEY'
					LIMIT 1
				),
				''
			)
		FROM information_schema.columns c
		WHERE table_catalog = $1
		AND table_schema = $2
		AND table_name = $3
		ORDER BY ordinal_position
	`

	rows, err := sqlConn.DB().QueryContext(
		ctx,
		query,
		table.Database,
		table.Schema,
		table.Name,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []entities.InspectColumnInfo

	for rows.Next() {

		var (
			col entities.InspectColumnInfo

			isNullable string
			columnKey  string

			defaultValue sql.NullString
		)

		err := rows.Scan(
			&col.Name,
			&col.DatabaseType,
			&isNullable,
			&defaultValue,
			&columnKey,
		)
		if err != nil {
			return nil, err
		}

		col.Nullable = isNullable == "YES"

		if defaultValue.Valid {
			col.DefaultValue = defaultValue.String
		}

		col.PrimaryKey = columnKey == "PRIMARY KEY"

		if defaultValue.Valid {
			// PostgreSQL serial/bigserial/identity columns
			col.AutoIncrement =
				len(defaultValue.String) > 7 &&
					defaultValue.String[:7] == "nextval"
		}

		columns = append(columns, col)
	}

	return columns, rows.Err()
}
