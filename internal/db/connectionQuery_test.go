package db

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"db-viewer/internal/engine/entities"
)

func testConnectionsDB(t *testing.T) {
	t.Helper()
	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	database.SetMaxOpenConns(1)
	previous := Conn
	Conn = database
	t.Cleanup(func() { Conn = previous; database.Close() })
	if err := runMigration(database); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateConnectionPreservesIdentityAndHistory(t *testing.T) {
	testConnectionsDB(t)
	config := entities.ConnectionConfig{Name: "original", Type: "pgx", Host: "localhost", Port: 5433, User: "user", Password: "old", Database: "app", SSHConfig: &entities.SSHConfig{Name: "tunnel", Host: "ssh.example", Port: 22, Username: "sshuser", AuthMethod: "password", Password: "ssh-password"}}
	id, err := StoreConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	config.ID = fmt.Sprint(id)
	if _, err := Conn.Exec(`UPDATE connections SET pinned = 1 WHERE id = ?`, id); err != nil {
		t.Fatal(err)
	}
	if _, err := Conn.Exec(`INSERT INTO query_history (connection_id, database_name, query_text, duration_ms, status) VALUES (?, 'app', 'SELECT 1', 1, 'SUCCESS')`, config.ID); err != nil {
		t.Fatal(err)
	}
	config.Name, config.Host, config.Port, config.Password, config.Database = "edited", "db.example", 5434, "new", "newdb"
	config.ReadOnly = true
	config.SSHConfig.Host = "new-ssh.example"
	if err := UpdateConnection(config); err != nil {
		t.Fatal(err)
	}
	list, err := GetConnectionList()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Id != int(id) || list[0].Name != "edited" || list[0].Host != "db.example" || list[0].Port.Int64 != 5434 || list[0].Password != "new" || list[0].DBName != "newdb" || !list[0].Pinned || !list[0].ReadOnly || list[0].SSHConfig.Host != "new-ssh.example" || list[0].SSHConfig.Password.String != "ssh-password" {
		t.Fatalf("update failed: %+v", list)
	}
	var count int
	if err := Conn.QueryRow(`SELECT count(*) FROM query_history WHERE connection_id = ?`, config.ID).Scan(&count); err != nil || count != 1 {
		t.Fatalf("history lost: %d %v", count, err)
	}
	config.SSHConfig = nil
	if err := UpdateConnection(config); err != nil {
		t.Fatal(err)
	}
	list, err = GetConnectionList()
	if err != nil || list[0].SSHConfigId.Valid {
		t.Fatalf("SSH not disabled: %v", err)
	}
	if err := Conn.QueryRow(`SELECT count(*) FROM ssh_configs`).Scan(&count); err != nil || count != 0 {
		t.Fatalf("orphan SSH config: %d %v", count, err)
	}
}

func TestUpdateConnectionSharedSSHAndRollback(t *testing.T) {
	testConnectionsDB(t)
	config := entities.ConnectionConfig{Name: "one", Type: "mysql", SSHConfig: &entities.SSHConfig{Name: "shared", Host: "original", Username: "user", Port: 22}}
	id, err := StoreConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Conn.Exec(`INSERT INTO connections (name, driver, host, port, user, password, dbname, ssh_config_id) SELECT 'two', driver, host, port, user, password, dbname, ssh_config_id FROM connections WHERE id = ?`, id); err != nil {
		t.Fatal(err)
	}
	config.ID = fmt.Sprint(id)
	config.SSHConfig.Host = "changed"
	if err := UpdateConnection(config); err != nil {
		t.Fatal(err)
	}
	var host string
	if err := Conn.QueryRow(`SELECT s.host FROM ssh_configs s JOIN connections c ON s.id = c.ssh_config_id WHERE c.name = 'two'`).Scan(&host); err != nil || host != "original" {
		t.Fatalf("shared SSH changed: %s %v", host, err)
	}
	if _, err := Conn.Exec(`CREATE TRIGGER reject_update BEFORE UPDATE ON connections BEGIN SELECT RAISE(ABORT, 'test rejection'); END`); err != nil {
		t.Fatal(err)
	}
	if err := UpdateConnection(config); err == nil {
		t.Fatal("expected failed update")
	}
	var count int
	if err := Conn.QueryRow(`SELECT count(*) FROM ssh_configs`).Scan(&count); err != nil || count != 2 {
		t.Fatalf("failed update leaked SSH config: %d %v", count, err)
	}
	config.ID = "99999"
	if err := UpdateConnection(config); err == nil {
		t.Fatal("updated missing connection")
	}
}

func TestConnectionMigrationExistingDatabase(t *testing.T) {
	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	old := strings.ReplaceAll(migrationSQL, "    read_only BOOLEAN NOT NULL DEFAULT false,\n", "")
	if _, err := database.Exec(old); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO connections (name) VALUES ('existing')`); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := runMigration(database); err != nil {
			t.Fatal(err)
		}
	}
	var name string
	var readOnly bool
	if err := database.QueryRow(`SELECT name, read_only FROM connections`).Scan(&name, &readOnly); err != nil || name != "existing" || readOnly {
		t.Fatalf("migration lost saved settings: %s %v %v", name, readOnly, err)
	}
}
