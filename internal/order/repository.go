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

type MongoRepository struct {
	*repository.BaseRepository
}

func NewMongoRepository(db *mongo.Database) Repository {
	return &MongoRepository{
		BaseRepository: repository.NewBaseRepository(db, "orders"),
	}
}

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

func (r *MongoRepository) List(ctx context.Context, skip, limit int64, findOptions *options.FindOptions) ([]*Order, int64, error) {
	total, err := r.BaseRepository.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}
	if findOptions == nil {
		findOptions = options.Find()
	}
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	documents, err := r.BaseRepository.GetAll(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get orders: %w", err)
	}
	var orders []*Order
	err = repository.DocumentsToStructs(documents, &orders)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to convert documents to orders: %w", err)
	}
	return orders, total, nil
}

type Repository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, skip, limit int64, findOptions *options.FindOptions) ([]*Order, int64, error)
}
