package mongodb

import (
	"context"
	"db-viewer/internal/engine/entities"
	"db-viewer/internal/engine/transports"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCommands(t *testing.T) {
	command, err := parseCommand(`{"find":"users","filter":{"_id":{"$oid":"507f1f77bcf86cd799439011"}}}`)
	if err != nil {
		t.Fatal(err)
	}
	if command[0].Key != "find" {
		t.Fatal("command order lost")
	}
	filter := command[1].Value.(bson.D)
	if _, ok := filter[0].Value.(primitive.ObjectID); !ok {
		t.Fatal("ObjectID not decoded")
	}
	for _, query := range []string{`{}`, `[]`, `db.users.find({})`, `{"find":"users","tailable":true}`, `{"aggregate":"users","pipeline":[{"$changeStream":{}}]}`, `{"eval":"code"}`} {
		if _, err := parseCommand(query); err == nil {
			t.Errorf("accepted %s", query)
		}
	}
	for _, example := range []struct {
		query string
		read  bool
	}{
		{`{"find":"users"}`, true},
		{`{"aggregate":"users","pipeline":[{"$match":{}}]}`, true},
		{`{"aggregate":"users","pipeline":[{"$out":"copy"}]}`, false},
		{`{"aggregate":"users","pipeline":[{"$merge":"copy"}]}`, false},
		{`{"insert":"users","documents":[{}]}`, false},
		{`{"drop":"users"}`, false},
	} {
		command, err := parseCommand(example.query)
		if err != nil {
			t.Fatal(err)
		}
		if got := readCommand(command); got != example.read {
			t.Errorf("readCommand(%s) = %v", example.query, got)
		}
	}
}

func TestConnectionOptions(t *testing.T) {
	config := entities.ConnectionConfig{Type: "mongodb", Host: "localhost", Port: 27017, Database: "app", User: "user", Password: "p@ss:/word"}
	conn := New(config, transports.NewDirect(config.Host, config.Port))
	opts, err := conn.clientOptions()
	if err != nil {
		t.Fatal(err)
	}
	if opts.Hosts[0] != "localhost:27017" || opts.Auth.Password != config.Password || opts.Auth.AuthSource != "admin" {
		t.Fatal("incorrect host or credentials")
	}
	conn.config.Host = "mongodb://uriuser:uripass@localhost:27017/app?authSource=custom"
	conn.config.User = ""
	opts, err = conn.clientOptions()
	if err != nil {
		t.Fatal(err)
	}
	if opts.Auth.AuthSource != "custom" || opts.Auth.Username != "uriuser" {
		t.Fatal("URI authentication ignored")
	}
	conn.config.Host = "localhost"
	conn.config.SSHConfig = &entities.SSHConfig{}
	opts, err = conn.clientOptions()
	if err != nil {
		t.Fatal(err)
	}
	if opts.Direct == nil || !*opts.Direct {
		t.Fatal("SSH must use a direct connection")
	}
	conn.config.Host = "mongodb://localhost"
	if _, err := conn.clientOptions(); err == nil {
		t.Fatal("SSH accepted URI")
	}
}

func TestDocumentPages(t *testing.T) {
	for _, count := range []int{0, 2, 3, 4} {
		t.Run(string(rune('0'+count)), func(t *testing.T) {
			documents := []any{}
			for i := 0; i < count; i++ {
				documents = append(documents, bson.D{{Key: "_id", Value: primitive.NewObjectID()}, {Key: "nested", Value: bson.D{{Key: "index", Value: i}}}})
			}
			cursor, err := mongo.NewCursorFromDocuments(documents, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			pager := NewPager()
			stream := &documentCursor{connection: "owner", cursor: cursor, size: 2}
			pager.cursors["result"] = stream
			defer pager.CloseConnection("owner")
			first, err := pager.Fetch(context.Background(), "result", "owner", 1)
			if err != nil {
				t.Fatal(err)
			}
			if first.StartRow != 1 || first.HasMore != (count > 2) {
				t.Fatalf("bad first page: %+v", first)
			}
			if len(first.Rows) > 0 && !strings.Contains(first.Rows[0][0].(string), `"$oid"`) {
				t.Fatal("ObjectID lost in result")
			}
			if first.HasMore {
				if _, err := pager.Fetch(context.Background(), "result", "other", 2); err == nil {
					t.Fatal("cursor accessible across connections")
				}
				if _, err := pager.Fetch(context.Background(), "result", "owner", 3); err == nil {
					t.Fatal("accepted skipped page")
				}
				retry, err := pager.Fetch(context.Background(), "result", "owner", 1)
				if err != nil || retry != first {
					t.Fatal("retry lost first page")
				}
				second, err := pager.Fetch(context.Background(), "result", "owner", 2)
				if err != nil {
					t.Fatal(err)
				}
				if second.StartRow != 3 || second.HasMore || len(first.Rows)+len(second.Rows) != count {
					t.Fatalf("bad second page: %+v", second)
				}
			}
			if len(pager.cursors) != 0 {
				t.Fatal("exhausted cursor retained")
			}
		})
	}
}

func TestCursorClose(t *testing.T) {
	cursor, err := mongo.NewCursorFromDocuments([]any{bson.D{{Key: "value", Value: 1}}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	pager := NewPager()
	stream := &documentCursor{connection: "owner", cursor: cursor, size: 1}
	pager.cursors["result"] = stream
	pager.Close("result", "other")
	if len(pager.cursors) != 1 {
		t.Fatal("another connection closed cursor")
	}
	pager.CloseConnection("owner")
	pager.CloseConnection("owner")
	if !stream.closed {
		t.Fatal("cursor not closed")
	}
	if _, err := pager.Fetch(context.Background(), "result", "owner", 1); err == nil {
		t.Fatal("closed cursor accepted")
	}
}

func TestDocumentPrecisionAndWriteErrors(t *testing.T) {
	document, err := bson.Marshal(bson.D{{Key: "large", Value: int64(9007199254740993)}})
	if err != nil {
		t.Fatal(err)
	}
	row, err := documentRow(document)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(row[0].(string), `"$numberLong":"9007199254740993"`) {
		t.Fatal("large integer precision lost")
	}
	for _, document := range []bson.D{
		{{Key: "ok", Value: 1}, {Key: "writeErrors", Value: bson.A{bson.D{{Key: "code", Value: 11000}, {Key: "errmsg", Value: "duplicate key"}}}}},
		{{Key: "ok", Value: 1}, {Key: "writeConcernError", Value: bson.D{{Key: "code", Value: 64}, {Key: "errmsg", Value: "write concern timeout"}}}},
	} {
		raw, err := bson.Marshal(document)
		if err != nil {
			t.Fatal(err)
		}
		if err := commandWriteError(raw); err == nil {
			t.Fatal("write error reported as successful")
		}
	}
	raw, _ := bson.Marshal(bson.D{{Key: "ok", Value: 1}, {Key: "n", Value: 1}})
	if err := commandWriteError(raw); err != nil {
		t.Fatal(err)
	}
}

func TestDocumentTableColumns(t *testing.T) {
	docs := []bson.D{
		{{Key: "name", Value: "Alice"}, {Key: "age", Value: int32(30)}, {Key: "nested", Value: bson.D{{Key: "active", Value: true}}}},
		{{Key: "enabled", Value: false}, {Key: "name", Value: "Bob"}, {Key: "large", Value: int64(9007199254740993)}},
	}
	raw := make([]bson.Raw, len(docs))
	for i, doc := range docs {
		encoded, err := bson.Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		raw[i] = encoded
	}
	columns, rows, err := documentTable(raw)
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"name", "age", "nested", "enabled", "large"}
	if len(columns) != len(expected) {
		t.Fatalf("columns: %+v", columns)
	}
	for i, name := range expected {
		if columns[i].Name != name {
			t.Fatalf("column %d: %+v", i, columns[i])
		}
	}
	if rows[0][0] != "Alice" || rows[1][0] != "Bob" || rows[0][1] != int32(30) || rows[1][3] != false {
		t.Fatalf("misaligned values: %+v", rows)
	}
	if rows[0][3] != nil || rows[1][1] != nil {
		t.Fatalf("missing fields: %+v", rows)
	}
	if columns[2].DatabaseType != "Extended JSON" || rows[0][2] != `{"active":true}` {
		t.Fatalf("nested value: %+v", rows[0][2])
	}
	if rows[1][4] != `{"$numberLong":"9007199254740993"}` {
		t.Fatalf("integer precision: %+v", rows[1][4])
	}
}
