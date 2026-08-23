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
		return  nil, err
	}

	return Conn, nil
}

func runMigration(db *sql.DB) error {
	_, err := db.Exec(migrationSQL)
	if err != nil {
		return  err
	}
	return  nil
}