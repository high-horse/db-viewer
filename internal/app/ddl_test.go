package app

import (
	"context"
	"database/sql"
	manager "db-viewer/internal/engine/connectionManager"
	"db-viewer/internal/engine/entities"
	"db-viewer/internal/engine/metadata"
	"db-viewer/internal/engine/metadata/mySQLInspector"
	pgxInspector "db-viewer/internal/engine/metadata/postgres"
	"os"
	"strings"
	"testing"
)

func TestSQLiteDDL(t *testing.T) {
	service, table, conn := editService(t)
	ctx := context.Background()
	_, err := conn.DB().Exec(`CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT NOT NULL DEFAULT 'hello', UNIQUE(name)); CREATE INDEX item_name ON items(name); CREATE TRIGGER item_insert AFTER INSERT ON items BEGIN UPDATE items SET name = upper(name) WHERE id = new.id; END; CREATE VIEW item_view AS SELECT name FROM items`)
	if err != nil {
		t.Fatal(err)
	}
	ddl, err := service.GetDDL(ctx, table)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"CREATE TABLE items", "DEFAULT 'hello'", "UNIQUE(name)", "CREATE INDEX item_name", "CREATE TRIGGER item_insert"} {
		if !strings.Contains(ddl, expected) {
			t.Fatalf("DDL missing %q: %s", expected, ddl)
		}
	}
	// Execute the returned DDL after dropping the original table to verify it is valid SQL.
	if _, err := conn.DB().Exec(`DROP VIEW item_view; DROP TABLE items; ` + ddl); err != nil {
		t.Fatalf("DDL cannot recreate table: %v\n%s", err, ddl)
	}
	if _, err := conn.DB().Exec(`CREATE VIEW item_view AS SELECT name FROM items`); err != nil {
		t.Fatal(err)
	}
	table.Name = "item_view"
	ddl, err = service.GetDDL(ctx, table)
	if err != nil || !strings.Contains(ddl, "CREATE VIEW") {
		t.Fatalf("view DDL: %q %v", ddl, err)
	}
	table.Name = "missing"
	if _, err := service.GetDDL(ctx, table); err == nil {
		t.Fatal("missing table did not return error")
	}
	table.Name = "items"
	table.ConnectionID = "missing"
	if _, err := service.GetDDL(ctx, table); err == nil {
		t.Fatal("missing connection did not return error")
	}
}

// Set the DSNs to disposable databases to run the live SQL inspector checks.
func TestLiveSQLDDL(t *testing.T) {
	for _, dialect := range []string{"postgres", "mysql"} {
		t.Run(dialect, func(t *testing.T) {
			dsn := os.Getenv("DB_VIEWER_DDL_" + strings.ToUpper(dialect) + "_DSN")
			if dsn == "" {
				t.Skip("no disposable test database DSN configured")
			}
			driver := dialect
			if driver == "postgres" {
				driver = "pgx"
			}
			database, err := sql.Open(driver, dsn)
			if err != nil {
				t.Fatal(err)
			}
			defer database.Close()
			ctx := context.Background()
			conn := ddlTestConnection{database: database}
			table := entities.TableRef{Name: "odd.table", Schema: "ddl_check", Database: "ddl_check"}
			var inspector metadata.Inspector
			var setup, cleanup, drop, view string
			if dialect == "postgres" {
				inspector = pgxInspector.NewInspector()
				setup = `CREATE SCHEMA ddl_check; CREATE TABLE ddl_check."odd.table" (id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY, name TEXT NOT NULL DEFAULT 'hello', doubled INTEGER GENERATED ALWAYS AS (id * 2) STORED, CONSTRAINT name_unique UNIQUE(name)); CREATE INDEX item_name ON ddl_check."odd.table"(name)`
				cleanup = `DROP SCHEMA ddl_check CASCADE`
				drop = `DROP TABLE ddl_check."odd.table"`
				view = `CREATE VIEW ddl_check.item_view AS SELECT name FROM ddl_check."odd.table"`
			} else {
				inspector = mySQLInspector.NewInspector()
				setup = "CREATE DATABASE ddl_check"
				cleanup = "DROP DATABASE ddl_check"
				drop = "DROP TABLE ddl_check.`odd.table`"
				view = "CREATE VIEW ddl_check.item_view AS SELECT name FROM ddl_check.`odd.table`"
			}
			if _, err := database.ExecContext(ctx, setup); err != nil {
				t.Fatal(err)
			}
			defer database.ExecContext(ctx, cleanup)
			if dialect == "mysql" {
				if _, err := database.ExecContext(ctx, "CREATE TABLE ddl_check.`odd.table` (id INTEGER AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100) NOT NULL DEFAULT 'hello', doubled INTEGER GENERATED ALWAYS AS (length(name) * 2) STORED, UNIQUE(name), INDEX item_name(name))"); err != nil {
					t.Fatal(err)
				}
			}
			ddl, err := inspector.GetTableDDL(ctx, conn, table)
			if err != nil {
				t.Fatal(err)
			}
			for _, expected := range []string{"CREATE TABLE", "odd.table", "hello", "PRIMARY KEY", "item_name", "GENERATED ALWAYS"} {
				if !strings.Contains(ddl, expected) {
					t.Fatalf("missing %q in DDL:\n%s", expected, ddl)
				}
			}
			if _, err := database.ExecContext(ctx, drop); err != nil {
				t.Fatal(err)
			}
			// MySQL's native statement contains an unqualified table name.
			if dialect == "mysql" {
				ddl = strings.Replace(ddl, "CREATE TABLE `odd.table`", "CREATE TABLE `ddl_check`.`odd.table`", 1)
			}
			for _, statement := range strings.Split(ddl, ";") {
				if strings.TrimSpace(statement) == "" {
					continue
				}
				if _, err := database.ExecContext(ctx, statement); err != nil {
					t.Fatalf("DDL cannot recreate table: %v\n%s", err, ddl)
				}
			}
			if _, err := database.ExecContext(ctx, view); err != nil {
				t.Fatal(err)
			}
			table.Name = "item_view"
			ddl, err = inspector.GetTableDDL(ctx, conn, table)
			if err != nil || !strings.Contains(ddl, "VIEW") {
				t.Fatalf("view DDL: %s %v", ddl, err)
			}
			table.Name = "missing"
			if _, err := inspector.GetTableDDL(ctx, conn, table); err == nil {
				t.Fatal("expected missing table error")
			}
		})
	}
}

type ddlTestConnection struct {
	manager.Connection
	database *sql.DB
}

func (c ddlTestConnection) DB() *sql.DB          { return c.database }
func (c ddlTestConnection) DatabaseName() string { return "ddl_check" }
