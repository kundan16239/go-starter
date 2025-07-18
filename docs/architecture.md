# Project Architecture

## Overview
This document provides a high-level overview of the system architecture, its main components, and how they interact. It is intended to help new contributors quickly understand the structure and design principles of the project.

## Main Components

- **cmd/**: Application entry points (e.g., main.go for starting the app).
- **internal/**: Application-specific business logic, organized by domain (e.g., user management).
- **pkg/**: Reusable packages and shared utilities (e.g., database, logging, messaging).
- **di/**: Dependency injection setup for wiring components together.
- **docs/**: Project documentation and onboarding guides.

## Component Interaction
- The application starts from `cmd/main.go`, which initializes dependencies using the DI container (`di/container.go`).
- Business logic is organized into domains under `internal/`, each with its own models, services, handlers, and repositories.
- Shared functionality (e.g., database access, logging, messaging) is provided by packages in `pkg/`.
- Middleware and error handling are implemented in `pkg/shared/`.

## Data Flow Example
1. An HTTP request is received by a handler in `internal/user/handler.go`.
2. The handler calls a service (e.g., `UserService`) to process the request.
3. The service interacts with repositories for data access and may use shared utilities from `pkg/`.
4. The response is returned to the client.

## Design Principles
- **Separation of Concerns**: Each package/folder has a clear responsibility.
- **Reusability**: Shared code is placed in `pkg/` for use across domains.
- **Testability**: Code is organized to facilitate unit and integration testing.

## Diagram
Consider adding a diagram here to visualize the architecture (e.g., using Mermaid or another tool).

---
For more details on each component, see the respective documentation files in the `docs/` folder. 