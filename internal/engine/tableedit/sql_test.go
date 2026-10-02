package tableedit

import (
	"encoding/json"
	"strings"
	"testing"

	"db-viewer/internal/engine/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSQLParametersAndColumnValidation(t *testing.T) {
	info := entities.TableEditInfo{Table: entities.TableRef{Name: `odd"table`, Schema: "public"}, Driver: "pgx", CanInsert: true, CanModify: true, Keys: []string{"id"}, Columns: []entities.InspectColumnInfo{{Name: "id"}, {Name: "name"}, {Name: "generated", Generated: true}}}
	change := entities.RowChange{Operation: "update", Keys: map[string]json.RawMessage{"id": json.RawMessage(`9007199254740993`)}, Values: map[string]json.RawMessage{"name": json.RawMessage(`"'; DROP TABLE items; --"`)}}
	query, args, err := BuildSQL(info, change)
	if err != nil {
		t.Fatal(err)
	}
	if query != `UPDATE "public"."odd""table" SET "name" = $1 WHERE "id" = $2` || len(args) != 2 || args[1] != "9007199254740993" || strings.Contains(query, "DROP") {
		t.Fatalf("unsafe SQL or integer precision loss: %s %+v", query, args)
	}
	for _, name := range []string{"missing", "generated"} {
		change.Values = map[string]json.RawMessage{name: json.RawMessage(`1`)}
		if _, _, err := BuildSQL(info, change); err == nil {
			t.Fatalf("accepted unwritable column: %s", name)
		}
	}
	info.Driver = "mysql"
	query, _, err = BuildSQL(info, entities.RowChange{Operation: "insert", Values: map[string]json.RawMessage{}})
	if err != nil || !strings.HasSuffix(query, "() VALUES ()") {
		t.Fatalf("default MySQL insert: %s %v", query, err)
	}
}

func TestMongoExtendedJSONPreservesValues(t *testing.T) {
	document, err := MongoDocument(map[string]json.RawMessage{
		"_id":    json.RawMessage(`{"$oid":"507f1f77bcf86cd799439011"}`),
		"large":  json.RawMessage(`{"$numberLong":"9007199254740993"}`),
		"nested": json.RawMessage(`{"values":[true,null,{"$numberDecimal":"1.25"}]}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := document["_id"].(primitive.ObjectID); !ok {
		t.Fatal("ObjectID became text")
	}
	if document["large"] != int64(9007199254740993) {
		t.Fatalf("integer precision lost: %+v", document)
	}
	for _, keys := range []map[string]json.RawMessage{nil, {"other": json.RawMessage(`1`)}, {"_id": json.RawMessage(`null`)}} {
		if _, err := MongoKey(entities.RowChange{Keys: keys}); err == nil {
			t.Fatal("accepted missing MongoDB identity")
		}
	}
}

func TestSQLIdentityColumnsCannotBeUpdated(t *testing.T) {
	info := entities.TableEditInfo{Driver: "sqlite", Table: entities.TableRef{Name: "items"}, CanInsert: true, CanModify: true, Keys: []string{"tenant", "key"}, Columns: []entities.InspectColumnInfo{{Name: "tenant", PrimaryKey: true}, {Name: "key", PrimaryKey: true}, {Name: "id"}, {Name: "name"}}}
	for _, name := range []string{"tenant", "key", "id"} {
		change := entities.RowChange{Operation: "update", Keys: map[string]json.RawMessage{"tenant": json.RawMessage(`1`), "key": json.RawMessage(`2`)}, Values: map[string]json.RawMessage{name: json.RawMessage(`3`)}}
		if _, _, err := BuildSQL(info, change); err == nil {
			t.Fatalf("allowed an update to identity column %s", name)
		}
	}
	// A new row may provide a required key that the database does not generate.
	if _, _, err := BuildSQL(info, entities.RowChange{Operation: "insert", Values: map[string]json.RawMessage{"tenant": json.RawMessage(`1`), "key": json.RawMessage(`2`)}}); err != nil {
		t.Fatal(err)
	}
	info.Columns[0].AutoIncrement = true
	if _, _, err := BuildSQL(info, entities.RowChange{Operation: "insert", Values: map[string]json.RawMessage{"tenant": json.RawMessage(`1`)}}); err == nil {
		t.Fatal("allowed overriding an auto-generated identity")
	}
}
