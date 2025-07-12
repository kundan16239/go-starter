package user

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

// Repository defines the interface for user data access operations
type Repository interface {
	// Create creates a new user in the database
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)

	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Update updates an existing user in the database
	Update(ctx context.Context, user *User) error

	// Delete removes a user from the database
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List retrieves a paginated list of users
	List(ctx context.Context, page, limit int) ([]*User, int64, error)

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)

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
		BaseRepository: repository.NewBaseRepository(db, "users"),
	}
}

// Create creates a new user in the database
func (r *MongoRepository) Create(ctx context.Context, user *User) error {
	id, err := r.BaseRepository.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to parse object ID: %w", err)
	}

	user.ID = objectID
	return nil
}

// GetByID retrieves a user by their ID
func (r *MongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	filter := bson.M{"_id": id}

	doc, err := r.BaseRepository.GetOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	if doc == nil {
		return nil, errors.New("user not found")
	}

	var user User
	err = repository.DocumentToStruct(doc, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to convert document to user: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by their email address
func (r *MongoRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	filter := bson.M{"email": email}

	doc, err := r.BaseRepository.GetOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if doc == nil {
		return nil, errors.New("user not found")
	}

	var user User
	err = repository.DocumentToStruct(doc, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to convert document to user: %w", err)
	}

	return &user, nil
}

// Update updates an existing user in the database
func (r *MongoRepository) Update(ctx context.Context, user *User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	modifiedCount, err := r.BaseRepository.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if modifiedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Delete removes a user from the database
func (r *MongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	deletedCount, err := r.BaseRepository.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if deletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// List retrieves a paginated list of users
func (r *MongoRepository) List(ctx context.Context, page, limit int) ([]*User, int64, error) {
	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Get total count
	total, err := r.BaseRepository.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Set up options for pagination and sorting
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{"created_at", -1}})

	// Execute query
	documents, err := r.BaseRepository.GetAll(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	// Convert documents to User structs
	var users []*User
	err = repository.DocumentsToStructs(documents, &users)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to convert documents to users: %w", err)
	}

	return users, total, nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *MongoRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}

	count, err := r.BaseRepository.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

// CountDocuments counts documents matching the filter
func (r *MongoRepository) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	return r.BaseRepository.CountDocuments(ctx, filter)
}

// Aggregate performs aggregation pipeline
func (r *MongoRepository) Aggregate(ctx context.Context, pipeline interface{}) ([]bson.M, error) {
	return r.BaseRepository.Aggregate(ctx, pipeline)
}
