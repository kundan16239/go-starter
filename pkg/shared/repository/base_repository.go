package repository

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseRepository provides a generic MongoDB repository implementation
type BaseRepository struct {
	collection *mongo.Collection
}

// NewBaseRepository creates a new base repository instance
func NewBaseRepository(db *mongo.Database, collectionName string) *BaseRepository {
	return &BaseRepository{
		collection: db.Collection(collectionName),
	}
}

// Create creates a new document in the database
func (r *BaseRepository) Create(ctx context.Context, document any) (string, error) {
	var doc map[string]any

	switch v := document.(type) {
	case bson.M:
		doc = v
	case map[string]any:
		doc = v
	default:
		val := reflect.ValueOf(document)

		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() == reflect.Struct {
			data, err := bson.Marshal(v)
			if err != nil {
				return "", fmt.Errorf("failed to marshal document: %w", err)
			}
			err = bson.Unmarshal(data, &doc)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal document: %w", err)
			}
		} else {
			return "", fmt.Errorf("document must be of type bson.M, map[string]interface{}, or a struct")
		}
	}

	now := time.Now()
	doc["created_at"] = now
	doc["updated_at"] = now

	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("failed to insert document: %w", err)
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// GetOne retrieves a single document from the database
func (r *BaseRepository) GetOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) (bson.M, error) {
	var findOptions *options.FindOneOptions
	if len(opts) > 0 {
		findOptions = opts[0]
	}

	result := r.collection.FindOne(ctx, filter, findOptions)
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}

	var document bson.M
	if err := result.Decode(&document); err != nil {
		return nil, err
	}

	return document, nil
}

// GetAll retrieves multiple documents from the database
func (r *BaseRepository) GetAll(ctx context.Context, filter any, opts ...*options.FindOptions) ([]bson.M, error) {
	var findOptions *options.FindOptions
	if len(opts) > 0 {
		findOptions = opts[0]
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []bson.M
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// UpdateOne updates a single document in the database
func (r *BaseRepository) UpdateOne(ctx context.Context, filter, update any) (int64, error) {
	switch u := update.(type) {
	case bson.M:
		if set, exists := u["$set"]; exists {
			if setMap, ok := set.(bson.M); ok {
				setMap["updated_at"] = time.Now()
			}
		} else {
			u["$set"] = bson.M{"updated_at": time.Now()}
		}
	case map[string]interface{}:
		if set, exists := u["$set"]; exists {
			if setMap, ok := set.(map[string]any); ok {
				setMap["updated_at"] = time.Now()
			}
		} else {
			u["$set"] = map[string]interface{}{"updated_at": time.Now()}
		}
	default:
		return 0, fmt.Errorf("unsupported update type")
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

// UpdateMany updates multiple documents in the database
func (r *BaseRepository) UpdateMany(ctx context.Context, filter, update interface{}) (int64, error) {
	switch u := update.(type) {
	case bson.M:
		if set, exists := u["$set"]; exists {
			if setMap, ok := set.(bson.M); ok {
				setMap["updated_at"] = time.Now()
			}
		} else {
			u["$set"] = bson.M{"updated_at": time.Now()}
		}
	case map[string]interface{}:
		if set, exists := u["$set"]; exists {
			if setMap, ok := set.(map[string]interface{}); ok {
				setMap["updated_at"] = time.Now()
			}
		} else {
			u["$set"] = map[string]interface{}{"updated_at": time.Now()}
		}
	default:
		return 0, fmt.Errorf("unsupported update type")
	}

	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

// CountDocuments counts documents matching the filter
func (r *BaseRepository) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}

// Distinct gets distinct values for a field
func (r *BaseRepository) Distinct(ctx context.Context, field string, filter interface{}) ([]interface{}, error) {
	results, err := r.collection.Distinct(ctx, field, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get distinct values: %w", err)
	}

	return results, nil
}

// Aggregate performs aggregation pipeline
func (r *BaseRepository) Aggregate(ctx context.Context, pipeline interface{}) ([]bson.M, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// DeleteOne deletes a single document
func (r *BaseRepository) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete document: %w", err)
	}

	return result.DeletedCount, nil
}

// DeleteMany deletes multiple documents
func (r *BaseRepository) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete documents: %w", err)
	}

	return result.DeletedCount, nil
}

// GetCollection returns the underlying MongoDB collection
func (r *BaseRepository) GetCollection() *mongo.Collection {
	return r.collection
}
