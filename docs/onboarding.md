# Onboarding Guide

Welcome to the project! This guide will help you get started as a contributor.

## 1. Project Overview
- This project is organized using a modular Go structure with clear separation between business logic, shared utilities, and entry points.
- See `docs/architecture.md` for a high-level system overview.

## 2. Getting Started
### Prerequisites
- Go (see `go.mod` for version)
- MongoDB (or your configured database)
- (Optional) Docker for running services locally

### Setup Steps
1. **Clone the repository:**
   ```sh
   git clone <repo-url>
   cd <project-directory>
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

## 3. Project Structure
- `cmd/` — Application entry point
- `internal/` — Domain logic (e.g., user management)
- `pkg/` — Shared utilities and packages
- `di/` — Dependency injection setup
- `docs/` — Documentation

## 4. Where to Find Help
- **Architecture:** `docs/architecture.md`
- **Domain Docs:** See `docs/<domain>.md` (e.g., `docs/user.md`)
- **API Reference:** `docs/api.md`
- **Error Handling:** `pkg/shared/errors/`
- **Logging:** `pkg/shared/logger/`

## 5. Contributing
- Follow Go best practices and project conventions.
- Write clear code comments and update documentation as needed.
- Add or update tests for your changes.
- Open a pull request with a clear description of your changes.

## 6. Additional Resources
- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)

---
If you have questions, ask in the project chat or open an issue. Welcome aboard! 