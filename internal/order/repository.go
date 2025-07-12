package order

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"project-structure/pkg/shared/repository"
)

// Repository defines the interface for order data access operations
type Repository interface {
	// Create creates a new order in the database
	Create(ctx context.Context, order *Order) error

	// GetByID retrieves an order by its ID
	GetByID(ctx context.Context, id primitive.ObjectID) (*Order, error)

	// GetByUserID retrieves orders by user ID
	GetByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*Order, int64, error)

	// Update updates an existing order in the database
	Update(ctx context.Context, order *Order) error

	// Delete removes an order from the database
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List retrieves a paginated list of orders
	List(ctx context.Context, page, limit int) ([]*Order, int64, error)

	// GetByStatus retrieves orders by status
	GetByStatus(ctx context.Context, status string, page, limit int) ([]*Order, int64, error)

	// CountDocuments counts documents matching the filter
	CountDocuments(ctx context.Context, filter interface{}) (int64, error)

	// Aggregate performs aggregation pipeline
	Aggregate(ctx context.Context, pipeline interface{}) ([]bson.M, error)
}

// MongoRepository implements the Repository interface using the optimized base repository
type MongoRepository struct {
	*repository.BaseRepository
}

// NewMongoRepository creates a new MongoDB repository instance
func NewMongoRepository(db *mongo.Database) Repository {
	return &MongoRepository{
		BaseRepository: repository.NewBaseRepository(db, "orders"),
	}
}

// Create creates a new order in the database
func (r *MongoRepository) Create(ctx context.Context, order *Order) error {
	id, err := r.BaseRepository.Create(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to parse object ID: %w", err)
	}

	order.ID = objectID
	return nil
}

// GetByID retrieves an order by its ID
func (r *MongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*Order, error) {
	filter := bson.M{"_id": id}

	doc, err := r.BaseRepository.GetOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}

	if doc == nil {
		return nil, errors.New("order not found")
	}

	var order Order
	err = repository.DocumentToStruct(doc, &order)
	if err != nil {
		return nil, fmt.Errorf("failed to convert document to order: %w", err)
	}

	return &order, nil
}

// GetByUserID retrieves orders by user ID
func (r *MongoRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*Order, int64, error) {
	filter := bson.M{"user_id": userID}

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Get total count
	total, err := r.BaseRepository.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Set up options for pagination and sorting
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{"created_at", -1}})

	// Execute query
	documents, err := r.BaseRepository.GetAll(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get orders: %w", err)
	}

	// Convert documents to Order structs
	var orders []*Order
	err = repository.DocumentsToStructs(documents, &orders)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to convert documents to orders: %w", err)
	}

	return orders, total, nil
}

// Update updates an existing order in the database
func (r *MongoRepository) Update(ctx context.Context, order *Order) error {
	filter := bson.M{"_id": order.ID}
	update := bson.M{"$set": order}

	modifiedCount, err := r.BaseRepository.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	if modifiedCount == 0 {
		return errors.New("order not found")
	}

	return nil
}

// Delete removes an order from the database
func (r *MongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	deletedCount, err := r.BaseRepository.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if deletedCount == 0 {
		return errors.New("order not found")
	}

	return nil
}

// List retrieves a paginated list of orders
func (r *MongoRepository) List(ctx context.Context, page, limit int) ([]*Order, int64, error) {
	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Get total count
	total, err := r.BaseRepository.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Set up options for pagination and sorting
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{"created_at", -1}})

	// Execute query
	documents, err := r.BaseRepository.GetAll(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get orders: %w", err)
	}

	// Convert documents to Order structs
	var orders []*Order
	err = repository.DocumentsToStructs(documents, &orders)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to convert documents to orders: %w", err)
	}

	return orders, total, nil
}

// GetByStatus retrieves orders by status
func (r *MongoRepository) GetByStatus(ctx context.Context, status string, page, limit int) ([]*Order, int64, error) {
	filter := bson.M{"status": status}

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Get total count
	total, err := r.BaseRepository.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders by status: %w", err)
	}

	// Set up options for pagination and sorting
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{"created_at", -1}})

	// Execute query
	documents, err := r.BaseRepository.GetAll(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get orders by status: %w", err)
	}

	// Convert documents to Order structs
	var orders []*Order
	err = repository.DocumentsToStructs(documents, &orders)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to convert documents to orders: %w", err)
	}

	return orders, total, nil
}

// CountDocuments counts documents matching the filter
func (r *MongoRepository) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	return r.BaseRepository.CountDocuments(ctx, filter)
}

// Aggregate performs aggregation pipeline
func (r *MongoRepository) Aggregate(ctx context.Context, pipeline interface{}) ([]bson.M, error) {
	return r.BaseRepository.Aggregate(ctx, pipeline)
}
