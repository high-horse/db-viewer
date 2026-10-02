package db

import (
	"database/sql"
	_ "embed"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed migration.sql
var migrationSQL string

var Conn *sql.DB

func appDbPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "db-viewer")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "app.db"), nil
}

func InitDb() (*sql.DB, error) {
	path, err := appDbPath()
	if err != nil {
		return nil, err
	}
	log.Println("using app database:", path)

	Conn, err = sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := runMigration(Conn); err != nil {
		return nil, err
	}

	return Conn, nil
}

func runMigration(db *sql.DB) error {
	_, err := db.Exec(migrationSQL)
	if err != nil {
		return err
	}
	// Existing installations need the new persisted setting too.
	rows, err := db.Query("PRAGMA table_info(connections)")
	if err != nil {
		return err
	}
	found := false
	for rows.Next() {
		var id, notNull, primaryKey int
		var name, kind string
		var defaultValue any
		if err := rows.Scan(&id, &name, &kind, &notNull, &defaultValue, &primaryKey); err != nil {
			rows.Close()
			return err
		}
		if name == "read_only" {
			found = true
		}
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	if !found {
		_, err = db.Exec("ALTER TABLE connections ADD COLUMN read_only BOOLEAN NOT NULL DEFAULT false")
	}
	return err
}
