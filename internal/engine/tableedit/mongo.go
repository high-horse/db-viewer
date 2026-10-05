package tableedit

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"db-lens/internal/engine/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDocument(values map[string]json.RawMessage) (bson.M, error) {
	encoded, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	var document bson.M
	if err := bson.UnmarshalExtJSON(encoded, false, &document); err != nil {
		return nil, fmt.Errorf("invalid MongoDB Extended JSON: %w", err)
	}
	return document, nil
}
func MongoKey(change entities.RowChange) (bson.M, error) {
	if len(change.Keys) != 1 {
		return nil, fmt.Errorf("an original _id is required")
	}
	if _, ok := change.Keys["_id"]; !ok {
		return nil, fmt.Errorf("an original _id is required")
	}
	key, err := MongoDocument(change.Keys)
	if err != nil {
		return nil, err
	}
	if key["_id"] == nil {
		return nil, fmt.Errorf("_id cannot be null")
	}
	return key, nil
}

// Standalone MongoDB servers do not support multi-document transactions.
// Report each committed prefix so a retry never resubmits successful rows.
type MongoWriter interface {
	InsertOne(context.Context, any, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	DeleteOne(context.Context, any, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	ReplaceOne(context.Context, any, any, ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
}

func SaveMongo(ctx context.Context, collection MongoWriter, changes []entities.RowChange) entities.TableSaveResult {
	result := entities.TableSaveResult{}
	for i, change := range changes {
		var err error
		switch change.Operation {
		case "insert":
			var document bson.M
			document, err = MongoDocument(change.Values)
			if err == nil {
				_, err = collection.InsertOne(ctx, document)
			}
		case "update", "delete":
			var key bson.M
			key, err = MongoKey(change)
			if err == nil && change.Operation == "delete" {
				var deleted *mongo.DeleteResult
				deleted, err = collection.DeleteOne(ctx, key)
				if err == nil && deleted.DeletedCount != 1 {
					err = fmt.Errorf("document no longer exists; reload the collection")
				}
			} else if err == nil {
				var replacement bson.M
				replacement, err = MongoDocument(change.Values)
				if err == nil {
					if id, ok := replacement["_id"]; ok && !reflect.DeepEqual(id, key["_id"]) {
						err = fmt.Errorf("_id cannot be changed")
					} else {
						replacement["_id"] = key["_id"]
						var updated *mongo.UpdateResult
						updated, err = collection.ReplaceOne(ctx, key, replacement)
						if err == nil && updated.MatchedCount != 1 {
							err = fmt.Errorf("document no longer exists; reload the collection")
						}
					}
				}
			}
		default:
			err = fmt.Errorf("unknown row operation")
		}
		if err != nil {
			result.Error = fmt.Sprintf("change %d: %v", i+1, err)
			return result
		}
		result.Applied++
	}
	return result
}
