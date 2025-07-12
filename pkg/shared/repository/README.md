# Optimized Repository Pattern

This package provides an optimized MongoDB repository pattern that follows clean architecture principles and provides a generic, reusable base for all data access operations.

## Features

- **Generic Base Repository**: A reusable base repository that handles common MongoDB operations
- **Type Safety**: Support for structs, `bson.M`, and `map[string]interface{}`
- **Automatic Timestamps**: Automatic `created_at` and `updated_at` field management
- **Helper Functions**: Efficient conversion between BSON documents and structs
- **Error Handling**: Comprehensive error handling with context
- **Pagination Support**: Built-in pagination with MongoDB options

## Base Repository Methods

### Core Operations

- `Create(ctx, document)` - Creates a new document and returns the ID
- `GetOne(ctx, filter, opts...)` - Retrieves a single document
- `GetAll(ctx, filter, opts...)` - Retrieves multiple documents
- `UpdateOne(ctx, filter, update)` - Updates a single document
- `UpdateMany(ctx, filter, update)` - Updates multiple documents
- `DeleteOne(ctx, filter)` - Deletes a single document
- `DeleteMany(ctx, filter)` - Deletes multiple documents
- `CountDocuments(ctx, filter)` - Counts documents matching filter
- `Distinct(ctx, field, filter)` - Gets distinct values for a field
- `Aggregate(ctx, pipeline)` - Performs aggregation pipeline

### Helper Functions

- `DocumentToStruct(doc, target)` - Converts BSON document to struct
- `StructToDocument(source)` - Converts struct to BSON document
- `DocumentsToStructs(docs, target)` - Converts slice of documents to slice of structs

## Usage Example

### 1. Create a Base Repository

```go
// Create a base repository for any collection
baseRepo := repository.NewBaseRepository(db, "users")
```

### 2. Implement Domain-Specific Repository

```go
type UserRepository struct {
    *repository.BaseRepository
}

func NewUserRepository(db *mongo.Database) *UserRepository {
    return &UserRepository{
        BaseRepository: repository.NewBaseRepository(db, "users"),
    }
}

// Create a user
func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
    id, err := r.Create(ctx, user)
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

// Get user by ID
func (r *UserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
    filter := bson.M{"_id": id}
    
    doc, err := r.GetOne(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    if doc == nil {
        return nil, errors.New("user not found")
    }
    
    var user User
    err = repository.DocumentToStruct(doc, &user)
    if err != nil {
        return nil, fmt.Errorf("failed to convert document: %w", err)
    }
    
    return &user, nil
}

// Get users with pagination
func (r *UserRepository) List(ctx context.Context, page, limit int) ([]*User, int64, error) {
    skip := (page - 1) * limit
    
    // Get total count
    total, err := r.CountDocuments(ctx, bson.M{})
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count users: %w", err)
    }
    
    // Set up pagination options
    findOptions := options.Find().
        SetSkip(int64(skip)).
        SetLimit(int64(limit)).
        SetSort(bson.D{{"created_at", -1}})
    
    // Get documents
    documents, err := r.GetAll(ctx, bson.M{}, findOptions)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get users: %w", err)
    }
    
    // Convert to structs
    var users []*User
    err = repository.DocumentsToStructs(documents, &users)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to convert documents: %w", err)
    }
    
    return users, total, nil
}
```

### 3. Advanced Usage

```go
// Using with custom filters
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
    filter := bson.M{"email": email}
    doc, err := r.GetOne(ctx, filter)
    // ... handle result
}

// Using aggregation
func (r *UserRepository) GetUserStats(ctx context.Context) ([]bson.M, error) {
    pipeline := []bson.M{
        {"$group": bson.M{
            "_id": "$status",
            "count": bson.M{"$sum": 1},
        }},
    }
    
    return r.Aggregate(ctx, pipeline)
}

// Using distinct values
func (r *UserRepository) GetUniqueEmails(ctx context.Context) ([]interface{}, error) {
    return r.Distinct(ctx, "email", bson.M{})
}
```

## Benefits

1. **DRY Principle**: No need to repeat common MongoDB operations
2. **Consistency**: All repositories follow the same pattern
3. **Maintainability**: Changes to base operations affect all repositories
4. **Type Safety**: Support for both generic and strongly-typed operations
5. **Performance**: Optimized BSON conversion with helper functions
6. **Flexibility**: Support for complex queries, aggregation, and pagination

## Best Practices

1. **Always use the base repository**: Extend `BaseRepository` for domain-specific repositories
2. **Use helper functions**: Use `DocumentToStruct` and `DocumentsToStructs` for efficient conversion
3. **Handle errors properly**: Wrap errors with context using `fmt.Errorf`
4. **Use pagination**: Always implement pagination for list operations
5. **Validate input**: Validate structs before saving to database
6. **Use transactions**: For complex operations involving multiple documents

## Error Handling

The repository pattern provides comprehensive error handling:

```go
// Example error handling
user, err := repo.GetByID(ctx, id)
if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments) {
        return nil, errors.New("user not found")
    }
    return nil, fmt.Errorf("database error: %w", err)
}
```

## Performance Considerations

1. **Indexing**: Ensure proper indexes on frequently queried fields
2. **Projection**: Use projection in `FindOneOptions` to limit returned fields
3. **Batch Operations**: Use `UpdateMany` and `DeleteMany` for bulk operations
4. **Connection Pooling**: Configure MongoDB connection pool appropriately
5. **Caching**: Consider caching frequently accessed data 