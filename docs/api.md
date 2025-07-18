# API Documentation

This document describes the main API endpoints, their request/response formats, and usage notes.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
- Some endpoints require authentication via a bearer token.
- Obtain a token by logging in (`POST /auth/login`).
- Include the token in the `Authorization` header:
  ```
  Authorization: Bearer <token>
  ```

## User Endpoints

### Register a New User
- **POST** `/users`
- **Request Body:**
  ```json
  {
    "name": "Alice Smith",
    "email": "alice@example.com",
    "password": "securepassword"
  }
  ```
- **Response:**
  ```json
  {
    "id": "user_id",
    "name": "Alice Smith",
    "email": "alice@example.com"
  }
  ```

### Login
- **POST** `/auth/login`
- **Request Body:**
  ```json
  {
    "email": "alice@example.com",
    "password": "securepassword"
  }
  ```
- **Response:**
  ```json
  {
    "token": "jwt_token"
  }
  ```

### Get User by ID
- **GET** `/users/{id}`
- **Response:**
  ```json
  {
    "id": "user_id",
    "name": "Alice Smith",
    "email": "alice@example.com"
  }
  ```

## Error Handling
- Errors are returned in the following format:
  ```json
  {
    "error": "Error message here."
  }
  ```

## Notes
- All request and response bodies are in JSON format.
- For more details on domain logic, see the respective domain documentation in `docs/`. 