package tableedit

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"db-lens/internal/engine/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type fakeMongoWriter struct {
	inserts     int
	failInsert  int
	key         bson.M
	replacement bson.M
	deleted     int64
	matched     int64
}

func (f *fakeMongoWriter) InsertOne(context.Context, any, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	f.inserts++
	if f.inserts == f.failInsert {
		return nil, fmt.Errorf("duplicate key")
	}
	return &mongo.InsertOneResult{}, nil
}
func (f *fakeMongoWriter) DeleteOne(_ context.Context, key any, _ ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	f.key = key.(bson.M)
	return &mongo.DeleteResult{DeletedCount: f.deleted}, nil
}
func (f *fakeMongoWriter) ReplaceOne(_ context.Context, key, replacement any, _ ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	f.key = key.(bson.M)
	f.replacement = replacement.(bson.M)
	return &mongo.UpdateResult{MatchedCount: f.matched}, nil
}
func TestMongoSaveReportsCommittedPrefix(t *testing.T) {
	writer := &fakeMongoWriter{failInsert: 2}
	changes := []entities.RowChange{{Operation: "insert", Values: map[string]json.RawMessage{"name": json.RawMessage(`"one"`)}}, {Operation: "insert", Values: map[string]json.RawMessage{"name": json.RawMessage(`"two"`)}}, {Operation: "insert", Values: map[string]json.RawMessage{"name": json.RawMessage(`"three"`)}}}
	result := SaveMongo(context.Background(), writer, changes)
	if result.Applied != 1 || result.Error == "" || writer.inserts != 2 {
		t.Fatalf("partial save: %+v calls %d", result, writer.inserts)
	}
}
func TestMongoSaveIdentityAndMissingDocuments(t *testing.T) {
	writer := &fakeMongoWriter{matched: 1, deleted: 1}
	keys := map[string]json.RawMessage{"_id": json.RawMessage(`{"$oid":"507f1f77bcf86cd799439011"}`)}
	change := entities.RowChange{Operation: "update", Keys: keys, Values: map[string]json.RawMessage{"large": json.RawMessage(`{"$numberLong":"9007199254740993"}`)}}
	result := SaveMongo(context.Background(), writer, []entities.RowChange{change})
	if result.Error != "" || result.Applied != 1 {
		t.Fatalf("replace failed: %+v", result)
	}
	if _, ok := writer.key["_id"].(primitive.ObjectID); !ok {
		t.Fatalf("identity lost BSON type: %+v", writer.key)
	}
	if writer.replacement["_id"] != writer.key["_id"] || writer.replacement["large"] != int64(9007199254740993) {
		t.Fatalf("replacement lost values: %+v", writer.replacement)
	}
	change.Values["_id"] = json.RawMessage(`"different"`)
	if result := SaveMongo(context.Background(), writer, []entities.RowChange{change}); result.Error == "" || result.Applied != 0 {
		t.Fatal("allowed _id to change")
	}
	writer.deleted = 0
	if result := SaveMongo(context.Background(), writer, []entities.RowChange{{Operation: "delete", Keys: keys}}); result.Error == "" || result.Applied != 0 {
		t.Fatal("missing document counted as saved")
	}
}
