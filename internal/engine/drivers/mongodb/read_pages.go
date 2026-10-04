package mongodb

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"math"
)

func commandInteger(value any) (int64, error) {
	switch n := value.(type) {
	case int32:
		return int64(n), nil
	case int64:
		return n, nil
	case int:
		return int64(n), nil
	case float64:
		if n >= math.MinInt64 && n < math.MaxInt64 && math.Trunc(n) == n {
			return int64(n), nil
		}
	}
	return 0, fmt.Errorf("MongoDB skip and limit must be integers")
}

func findCountCommand(command bson.D) (bson.D, error) {
	result := bson.D{{Key: "count", Value: command[0].Value}}
	for _, field := range command[1:] {
		switch field.Key {
		case "filter":
			result = append(result, bson.E{Key: "query", Value: field.Value})
		case "skip", "limit":
			n, err := commandInteger(field.Value)
			if err != nil {
				return nil, err
			}
			if field.Key == "skip" && n < 0 {
				return nil, fmt.Errorf("MongoDB skip must be non-negative")
			}
			if n == math.MinInt64 {
				return nil, fmt.Errorf("MongoDB limit is out of range")
			}
			if n < 0 {
				n = -n
			}
			result = append(result, bson.E{Key: field.Key, Value: n})
		case "collation", "hint", "maxTimeMS", "readConcern", "comment":
			result = append(result, field)
		}
	}
	return result, nil
}

func findPageCommand(command bson.D, total int64, size, page int) (bson.D, int64, error) {
	last := int64(1)
	if total > 0 {
		last = (total-1)/int64(size) + 1
	}
	target := int64(page)
	if target < 1 {
		target = 1
	}
	if target > last {
		target = last
	}
	offset := (target - 1) * int64(size)
	var skip, limit int64
	result := bson.D{command[0]}
	for _, field := range command[1:] {
		switch field.Key {
		case "skip", "limit":
			n, err := commandInteger(field.Value)
			if err != nil {
				return nil, 0, err
			}
			if field.Key == "skip" {
				skip = n
			} else {
				limit = n
			}
		case "batchSize", "singleBatch": // Use a bounded page with lookahead.
		default:
			result = append(result, field)
		}
	}
	if skip < 0 || skip > math.MaxInt64-offset || limit == math.MinInt64 {
		return nil, 0, fmt.Errorf("MongoDB page offset is out of range")
	}
	if limit < 0 {
		limit = -limit
	}
	if limit > 0 {
		limit -= offset
		if limit < 1 {
			limit = 1
		}
	}
	result = append(result, bson.E{Key: "skip", Value: skip + offset}, bson.E{Key: "batchSize", Value: size})
	if limit > 0 {
		result = append(result, bson.E{Key: "limit", Value: limit})
	}
	return result, offset, nil
}
