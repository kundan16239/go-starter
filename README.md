# Go Project Structure

A modular Go project following clean architecture principles, with MongoDB integration and the Gin web framework.

## Project Structure

```
project-structure-new/
├── cmd/                  # Application entry point (main.go)
├── di/                   # Dependency injection setup
├── docs/                 # Project documentation (architecture, onboarding, API, domains)
├── internal/             # Application-specific business logic (e.g., user domain)
│   └── user/             # User domain: models, handlers, services, repositories, routes, DTOs
├── pkg/                  # Reusable packages and shared utilities
│   ├── db/               # Database connection (e.g., mongo.go)
│   └── shared/           # Shared code: errors, helpers, logger, messaging, middleware, repository
├── go.mod                # Go module file
├── go.sum                # Go module checksums
├── env.example           # Example environment variables
└── README.md             # This file
```

## Key Features
- **Clean Architecture**: Clear separation of concerns between handlers, services, repositories, and models.
- **MongoDB Integration**: Database access via `pkg/db/mongo.go`.
- **Gin Web Framework**: Fast HTTP routing and middleware support.
- **Dependency Injection**: Managed in `di/container.go`.
- **Reusable Utilities**: Logging, error handling, helpers, and messaging in `pkg/shared/`.
- **Domain-Driven**: Business logic organized by domain in `internal/`.
- **Comprehensive Documentation**: See the `docs/` folder for architecture, onboarding, API, and domain guides.

## Getting Started

1. **Clone the repository:**
   ```sh
   git clone <repository-url>
   cd project-structure-new
   ```
2. **Copy environment variables:**
   ```sh
   cp env.example .env
   # Edit .env as needed
   ```
3. **Install dependencies:**
   ```sh
   go mod tidy
   ```
4. **Run the application:**
   ```sh
   go run cmd/main.go
   ```

## Documentation
- **Architecture:** [`docs/architecture.md`](docs/architecture.md)
- **Onboarding:** [`docs/onboarding.md`](docs/onboarding.md)
- **API Reference:** [`docs/api.md`](docs/api.md)
- **Domain Docs:** [`docs/user.md`](docs/user.md) (add more as needed)

## Contributing
- Follow Go best practices and project conventions.
- Write clear code comments and update documentation as needed.
- Add or update tests for your changes.
- Open a pull request with a clear description of your changes.

---
For more details, see the documentation in the `docs/` folder. 