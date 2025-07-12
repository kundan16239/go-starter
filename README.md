# Go Project Structure with MongoDB and Gin

A well-structured Go project following clean architecture principles with MongoDB integration and the Gin web framework.

## Project Structure

```
project-structure/
├── cmd/
│   └── appname/
│       └── main.go            # Application entry point
├── internal/
│   ├── user/
│   │   ├── dto.go             # Data Transfer Objects
│   │   ├── model.go           # Domain models
│   │   ├── repository.go      # Repository interface
│   │   ├── mongo_repository.go # MongoDB implementation
│   │   ├── service.go         # Business logic
│   │   └── handler.go         # HTTP handlers (Gin)
│   └── order/
│       └── dto.go             # Order DTOs (example)
├── pkg/
│   └── shared/
│       ├── logger/
│       │   └── logger.go      # Logging utilities
│       ├── metrics/
│       │   └── metrics.go     # Application metrics
│       ├── errors/
│       │   └── errors.go      # Custom error types
│       ├── helpers/
│       │   └── helpers.go     # Common utility functions
│       └── middleware/
│           └── middleware.go  # Custom middleware
├── go.mod                     # Go module file
└── README.md                  # This file
```

## Architecture Overview

This project follows **Clean Architecture** principles with clear separation of concerns:

### Layers

1. **Handlers (Adapters)**: HTTP request/response handling using Gin
2. **Services (Use Cases)**: Business logic and orchestration
3. **Repositories (Data Access)**: Database operations
4. **Models**: Domain entities and business rules

### Key Features

- **Gin Web Framework**: Fast HTTP web framework with middleware support
- **Dependency Injection**: Clean dependency management
- **Interface Segregation**: Repository interfaces for testability
- **Error Handling**: Custom error types with proper HTTP status codes
- **Logging**: Structured logging throughout the application
- **Metrics**: Basic application metrics collection
- **MongoDB Integration**: Full CRUD operations with MongoDB
- **Middleware Support**: Built-in logging and recovery middleware
- **Helper Functions**: Comprehensive utility functions for common operations
- **Input Validation**: Built-in validation using helper functions

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MongoDB (local or remote)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd project-structure
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up MongoDB:

**Option 1: Local MongoDB**
```bash
# Install MongoDB locally (macOS with Homebrew)
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb/brew/mongodb-community

# Or install MongoDB manually from https://www.mongodb.com/try/download/community
```

**Option 2: MongoDB Atlas (Cloud)**
- Go to [MongoDB Atlas](https://www.mongodb.com/atlas)
- Create a free account and cluster
- Get your connection string
- Update your `.env` file with the connection string

4. Set up environment variables:

Create a `.env` file in the root directory:
```bash
cp env.example .env
```

Then edit the `.env` file with your configuration:
```bash
# For local MongoDB:
MONGO_URI=mongodb://localhost:27017

# For MongoDB Atlas (replace with your connection string):
# MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/database?retryWrites=true&w=majority

# Database name
DB_NAME=appname

# Server port
PORT=8080
```

5. Run the application:
```bash
go run cmd/appname/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Users

- `POST /users` - Create a new user
- `GET /users` - List users (with pagination)
- `GET /users/:id` - Get user by ID
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user
- `POST /users/login` - User login

### Example Requests

#### Create User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

#### List Users
```bash
curl "http://localhost:8080/users?page=1&limit=10"
```

#### Login
```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

## Gin Framework Features

### Middleware
The application uses Gin's built-in middleware:
- **Logger**: Logs all HTTP requests
- **Recovery**: Recovers from panics and returns 500 error

### Route Groups
Routes are organized using Gin's route groups for better structure:
```go
users := router.Group("/users")
{
    users.POST("", h.CreateUser)
    users.GET("", h.ListUsers)
    users.GET("/:id", h.GetUser)
    // ...
}
```

### Request Binding
Gin provides automatic JSON binding and validation:
```go
var req CreateUserRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
    return
}
```

### Response Handling
Simplified JSON responses with proper status codes:
```go
c.JSON(http.StatusCreated, user)
c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
```

## Project Structure Explanation

### `/cmd`
Contains the main applications of the project. Each application has its own directory with a `main.go` file.

### `/internal`
Private application and library code. This is the code you don't want others importing in their applications or libraries.

- **`/user`**: Complete user management module with Gin handlers
- **`/order`**: Example order management module (can be expanded)

### `/pkg`
Library code that's ok to use by external applications. This is where you put code that you want to be used by other people.

- **`/shared/logger`**: Logging utilities
- **`/shared/metrics`**: Application metrics
- **`/shared/errors`**: Custom error types
- **`/shared/helpers`**: Common utility functions (validation, security, time, etc.)
- **`/shared/middleware`**: Custom middleware functions

## Development Guidelines

### Adding New Modules

1. Create a new directory in `/internal` (e.g., `/internal/product`)
2. Follow the same structure as `/internal/user`:
   - `dto.go` - Request/response structures
   - `model.go` - Domain models
   - `repository.go` - Repository interface
   - `mongo_repository.go` - MongoDB implementation
   - `service.go` - Business logic
   - `handler.go` - Gin HTTP handlers

### Error Handling

Use the custom error types from `/pkg/shared/errors`:

```go
import "project-structure/pkg/shared/errors"

// In your service
if user == nil {
    return nil, errors.NewNotFound("User not found")
}
```

### Logging

Use the logger from `/pkg/shared/logger`:

```go
import "project-structure/pkg/shared/logger"

logger.Info("User created successfully", "user_id", user.ID)
logger.Error("Failed to create user", "error", err)
```

### Helper Functions

Use helper functions from `/pkg/shared/helpers`:

```go
import "project-structure/pkg/shared/helpers"

// Validation
if !helpers.IsValidEmail(email) {
    return errors.New("invalid email")
}

// Security
hashedPassword, err := helpers.HashPassword(password)
token, err := helpers.GenerateToken()

// Time formatting
formattedDate := helpers.FormatDate(time.Now())

// String manipulation
slug := helpers.GenerateSlug("My Article Title")
```

### Adding Middleware

To add custom middleware to your Gin routes:

```go
// In main.go
router.Use(customMiddleware())

// Or for specific route groups
users := router.Group("/users")
users.Use(authMiddleware())
```

## Testing

To run tests:

```bash
go test ./...
```

## Contributing

1. Follow the existing code structure
2. Add tests for new functionality
3. Update documentation as needed
4. Use meaningful commit messages

## License

This project is licensed under the MIT License. 