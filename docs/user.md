# User Domain Documentation

## Overview
The User domain is responsible for all functionality related to users in the application. This includes creating new user accounts, managing user information, handling authentication (login/logout), and enforcing user-specific business rules. The User domain is central to application security and user experience.

## Key Models

### User
Represents an individual who can log in and interact with the application.

| Field         | Type       | Description                                 |
|-------------- |----------- |---------------------------------------------|
| ID            | string     | Unique identifier for the user              |
| Name          | string     | The user's full name                        |
| Email         | string     | The user's email address (must be unique)   |
| PasswordHash  | string     | Securely hashed password                    |
| CreatedAt     | time.Time  | When the user account was created           |
| UpdatedAt     | time.Time  | When the user account was last updated      |

## Services

### UserService
Handles all user-related operations, such as:
- Registering new users
- Retrieving user profiles
- Updating user information
- Deleting user accounts

### AuthService
Manages authentication processes, including:
- Logging users in and out
- Generating and validating authentication tokens
- Enforcing authentication and authorization rules

## API Endpoints

| Method | Endpoint                | Description                        |
|--------|-------------------------|------------------------------------|
| POST   | /api/v1/users           | Register a new user                |
| GET    | /api/v1/users/{id}      | Get details for a specific user    |
| PUT    | /api/v1/users/{id}      | Update user information            |
| DELETE | /api/v1/users/{id}      | Delete a user account              |
| POST   | /api/v1/auth/login      | Authenticate user and get a token  |

## Business Logic
- **Password Security:** All passwords are hashed before being stored in the database. Plaintext passwords are never saved.
- **Unique Emails:** Each user must have a unique email address. Duplicate emails are not allowed.
- **Authorization:**
  - Users can only update or delete their own accounts.
  - Admin users may have additional permissions (e.g., managing other users).
- **Data Validation:** All user input is validated to ensure data integrity and security.

## Example Usage

### Register a New User
```
POST /api/v1/users
{
  "name": "Alice Smith",
  "email": "alice@example.com",
  "password": "securepassword"
}
```

### Login
```
POST /api/v1/auth/login
{
  "email": "alice@example.com",
  "password": "securepassword"
}
```

## Notes & Recommendations
- Always handle sensitive data (like passwords and tokens) securely.
- Consider implementing features such as email verification and password reset for better security and user experience.
- Keep this documentation up to date as the user domain evolves.
- For more details on authentication, see [Authentication documentation](./auth.md) if available. 