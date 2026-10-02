package mongodb

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestFindPageCommands(t *testing.T) {
	command, err := parseCommand(`{"find":"users","filter":{"active":true},"sort":{"age":-1},"projection":{"name":1},"skip":10,"limit":250,"batchSize":20}`)
	if err != nil {
		t.Fatal(err)
	}
	count, err := findCountCommand(command)
	if err != nil {
		t.Fatal(err)
	}
	expectedCount := bson.D{{Key: "count", Value: "users"}, {Key: "query", Value: bson.D{{Key: "active", Value: true}}}, {Key: "skip", Value: int64(10)}, {Key: "limit", Value: int64(250)}}
	if !reflect.DeepEqual(count, expectedCount) {
		t.Fatalf("count command: %#v", count)
	}
	for _, test := range []struct {
		page                int
		total               int64
		offset, skip, limit int64
	}{
		{1, 250, 0, 10, 250}, {3, 250, 200, 210, 50}, {99, 250, 200, 210, 50}, {0, 250, 0, 10, 250}, {1, 0, 0, 10, 250},
	} {
		page, offset, err := findPageCommand(command, test.total, 100, test.page)
		if err != nil {
			t.Fatal(err)
		}
		if offset != test.offset {
			t.Fatalf("offset: %d", offset)
		}
		values := map[string]any{}
		for _, field := range page {
			values[field.Key] = field.Value
		}
		if values["skip"] != test.skip || values["limit"] != test.limit || values["batchSize"] != 100 {
			t.Fatalf("page: %#v", page)
		}
		if !reflect.DeepEqual(values["sort"], command[2].Value) || !reflect.DeepEqual(values["projection"], command[3].Value) {
			t.Fatalf("query options lost: %#v", page)
		}
	}
	if command[4].Value != int32(10) || command[5].Value != int32(250) {
		t.Fatal("original command mutated")
	}
}

func TestFindNegativeLimitAndInvalidOffsets(t *testing.T) {
	command, _ := parseCommand(`{"find":"users","limit":-250}`)
	count, err := findCountCommand(command)
	if err != nil || count[1].Value != int64(250) {
		t.Fatalf("negative limit: %#v %v", count, err)
	}
	page, offset, err := findPageCommand(command, 250, 100, 3)
	if err != nil || offset != 200 || page[len(page)-1].Value != int64(50) {
		t.Fatalf("last page: %#v %d %v", page, offset, err)
	}
	for _, query := range []string{`{"find":"users","skip":-1}`, `{"find":"users","limit":1.5}`} {
		command, _ := parseCommand(query)
		if _, err := findCountCommand(command); err == nil {
			t.Fatalf("accepted invalid offset: %s", query)
		}
	}
}

func TestCountedMongoCursorPages(t *testing.T) {
	docs := []any{bson.D{{Key: "name", Value: "third"}}, bson.D{{Key: "name", Value: "fourth"}}, bson.D{{Key: "name", Value: "fifth"}}}
	cursor, err := mongo.NewCursorFromDocuments(docs, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cursor.Close(context.Background())
	total := int64(5)
	stream := &documentCursor{cursor: cursor, size: 2, read: 2, page: 1, totalRows: &total}
	second, err := stream.fetch(context.Background(), "result", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if !second.CanNavigate || second.TotalRows == nil || *second.TotalRows != 5 || second.StartRow != 3 || !second.HasMore || second.Rows[0][0] != "third" {
		t.Fatalf("second page: %+v", second)
	}
	last, err := stream.fetch(context.Background(), "result", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if !last.CanNavigate || *last.TotalRows != 5 || last.StartRow != 5 || last.HasMore || len(last.Rows) != 1 || last.Rows[0][0] != "fifth" {
		t.Fatalf("last page: %+v", last)
	}
}
