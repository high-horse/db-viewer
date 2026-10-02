package app

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	manager "db-viewer/internal/engine/connectionManager"
	"db-viewer/internal/engine/entities"
)

func editService(t *testing.T) (*DbService, entities.TableRef, manager.SQLConnection) {
	t.Helper()
	s := NewDbService(nil)
	database := filepath.Join(t.TempDir(), "edit.db")
	_, err := s.Connect(context.Background(), entities.ConnectionConfig{ID: "test", Name: "test", Type: "sqlite", Database: database})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.ServiceShutdown() })
	conn, _ := s.manager.Active()
	return s, entities.TableRef{ConnectionID: "test", Name: "items", Schema: "main", Database: database}, conn.(manager.SQLConnection)
}
func raw(value any) json.RawMessage { encoded, _ := json.Marshal(value); return encoded }

func TestTableEditorSaveAndRollback(t *testing.T) {
	s, table, conn := editService(t)
	ctx := context.Background()
	_, err := conn.DB().Exec(`CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT NOT NULL UNIQUE, note TEXT DEFAULT 'default', doubled INTEGER GENERATED ALWAYS AS (id * 2) STORED); INSERT INTO items(id,name,note) VALUES(1,'first',NULL),(2,'second','keep')`)
	if err != nil {
		t.Fatal(err)
	}
	info, err := s.DescribeTableEdit(ctx, table)
	if err != nil {
		t.Fatal(err)
	}
	if !info.CanInsert || !info.CanModify || len(info.Keys) != 1 || info.Keys[0] != "id" || !info.Columns[3].Generated || !info.Columns[0].AutoIncrement {
		t.Fatalf("metadata: %+v", info)
	}
	page, err := s.ExecuteQuery(ctx, entities.QueryInput{Query: `SELECT * FROM main.items`, PageSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	changes := []entities.RowChange{
		{Operation: "update", Keys: map[string]json.RawMessage{"id": raw(1)}, Values: map[string]json.RawMessage{"name": raw("renamed"), "note": raw("")}},
		{Operation: "delete", Keys: map[string]json.RawMessage{"id": raw(2)}},
		{Operation: "insert", Values: map[string]json.RawMessage{"name": raw("new")}},
	}
	result, err := s.SaveTableChanges(ctx, entities.TableChanges{Table: table, Cursor: page.Cursor, Changes: changes})
	if err != nil || result.Applied != 3 {
		t.Fatalf("save: %+v %v", result, err)
	}
	var name, note string
	if err := conn.DB().QueryRow(`SELECT name,note FROM items WHERE id=1`).Scan(&name, &note); err != nil || name != "renamed" || note != "" {
		t.Fatalf("update failed: %s %s %v", name, note, err)
	}
	if err := conn.DB().QueryRow(`SELECT note FROM items WHERE name='new'`).Scan(&note); err != nil || note != "default" {
		t.Fatalf("default omitted: %s %v", note, err)
	}
	if _, err := s.pager.Fetch(page.Cursor, "test", 2); err == nil {
		t.Fatal("old result cursor was not closed")
	}
	rollback := []entities.RowChange{
		{Operation: "update", Keys: map[string]json.RawMessage{"id": raw(1)}, Values: map[string]json.RawMessage{"name": raw("should rollback")}},
		{Operation: "insert", Values: map[string]json.RawMessage{"name": raw("new")}},
	}
	if _, err := s.SaveTableChanges(ctx, entities.TableChanges{Table: table, Changes: rollback}); err == nil {
		t.Fatal("expected unique constraint failure")
	}
	if err := conn.DB().QueryRow(`SELECT name FROM items WHERE id=1`).Scan(&name); err != nil || name != "renamed" {
		t.Fatalf("batch did not roll back: %s %v", name, err)
	}
}

func TestTableEditorCompositeKeysAndRestrictions(t *testing.T) {
	s, table, conn := editService(t)
	ctx := context.Background()
	if _, err := conn.DB().Exec(`CREATE TABLE items (tenant INTEGER, id INTEGER, name TEXT, PRIMARY KEY(tenant,id)); INSERT INTO items VALUES(1,1,'first'),(2,1,'second'); CREATE TABLE no_key (name TEXT); CREATE VIEW item_view AS SELECT * FROM items`); err != nil {
		t.Fatal(err)
	}
	info, err := s.DescribeTableEdit(ctx, table)
	if err != nil || len(info.Keys) != 2 || info.Columns[0].AutoIncrement {
		t.Fatalf("composite keys: %+v %v", info, err)
	}
	change := entities.RowChange{Operation: "update", Keys: map[string]json.RawMessage{"tenant": raw(1), "id": raw(1)}, Values: map[string]json.RawMessage{"name": raw("changed")}}
	if _, err := s.SaveTableChanges(ctx, entities.TableChanges{Table: table, Changes: []entities.RowChange{change}}); err != nil {
		t.Fatal(err)
	}
	var name string
	if err := conn.DB().QueryRow(`SELECT name FROM items WHERE tenant=2`).Scan(&name); err != nil || name != "second" {
		t.Fatal("updated a different tenant")
	}
	delete(change.Keys, "tenant")
	if _, err := s.SaveTableChanges(ctx, entities.TableChanges{Table: table, Changes: []entities.RowChange{change}}); err == nil {
		t.Fatal("accepted incomplete primary key")
	}
	table.Name = "no_key"
	info, err = s.DescribeTableEdit(ctx, table)
	if err != nil || !info.CanInsert || info.CanModify {
		t.Fatalf("keyless table: %+v %v", info, err)
	}
	table.Name = "item_view"
	info, err = s.DescribeTableEdit(ctx, table)
	if err != nil || info.CanInsert || !strings.Contains(info.Reason, "Views") {
		t.Fatalf("view writable: %+v %v", info, err)
	}
	table.ConnectionID = "other"
	if _, err := s.DescribeTableEdit(ctx, table); err == nil {
		t.Fatal("cross-connection edit permitted")
	}
	table.ConnectionID = "test"
	table.Name = "items"
	config := conn.Config()
	config.ReadOnly = true
	if _, err := s.Connect(ctx, config); err != nil {
		t.Fatal(err)
	}
	info, err = s.DescribeTableEdit(ctx, table)
	if err != nil || info.CanInsert || info.CanModify {
		t.Fatalf("read-only connection editable: %+v %v", info, err)
	}
	if _, err := s.SaveTableChanges(ctx, entities.TableChanges{Table: table, Changes: []entities.RowChange{{Operation: "insert"}}}); err == nil {
		t.Fatal("saved on read-only connection")
	}
}

func TestTableEditorLargeIntegerKeys(t *testing.T) {
	s, table, conn := editService(t)
	ctx := context.Background()
	if _, err := conn.DB().Exec(`CREATE TABLE items(id INTEGER PRIMARY KEY, name TEXT); INSERT INTO items VALUES(9007199254740993,'original')`); err != nil {
		t.Fatal(err)
	}
	page, err := s.ExecuteQuery(ctx, entities.QueryInput{Query: `SELECT * FROM main.items`, PageSize: 100})
	if err != nil || page.Rows[0][0] != "9007199254740993" {
		t.Fatalf("primary key lost precision: %+v %v", page, err)
	}
	change := entities.RowChange{Operation: "update", Keys: map[string]json.RawMessage{"id": raw(page.Rows[0][0])}, Values: map[string]json.RawMessage{"name": raw("edited")}}
	if _, err := s.SaveTableChanges(ctx, entities.TableChanges{Table: table, Changes: []entities.RowChange{change}}); err != nil {
		t.Fatal(err)
	}
	var name string
	if err := conn.DB().QueryRow(`SELECT name FROM items WHERE id=9007199254740993`).Scan(&name); err != nil || name != "edited" {
		t.Fatalf("wrong row updated: %s %v", name, err)
	}
}
