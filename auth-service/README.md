
# Authentication Service Documentation

## Overview
The authentication service handles user management, including registration, login, JWT token generation, and admin features. Redis is used for session management and caching user data.

## API Endpoints

### Authentication Routes

| Method | Endpoint                 | Description            |
|--------|--------------------------|------------------------|
| POST   | `/auth/register`         | Register a new user    |
| POST   | `/auth/login`            | Login user and get tokens |
| POST   | `/auth/logout`           | Logout the user        |
| POST   | `/auth/refresh-token`    | Refresh access token   |

#### Example Request/Response for Register
**Request:**
```json
{
  "username": "testuser",
  "password": "testpassword",
  "email": "test@example.com"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

### User Routes (Protected)

| Method | Endpoint                   | Description              |
|--------|----------------------------|--------------------------|
| GET    | `/user/profile`            | Get user profile         |
| PATCH  | `/user/profile`            | Update user profile      |
| PATCH  | `/user/change-password`    | Change user password     |

### Admin Routes (Protected)

| Method | Endpoint                                  | Description                      |
|--------|-------------------------------------------|----------------------------------|
| GET    | `/admin/users`                            | Get list of all users            |
| GET    | `/admin/users/:user_id`                   | Get a user by ID                 |
| PATCH  | `/admin/users/:user_id`                   | Update user information by ID    |
| PATCH  | `/admin/users/:user_id/promote`           | Promote user to admin            |
| PATCH  | `/admin/users/:user_id/deactivate`        | Deactivate user account          |
| PATCH  | `/admin/users/:user_id/activate`          | Activate user account            |
| DELETE | `/admin/users/:user_id`                   | Delete user by ID                |

## Redis Integration
Redis is used in this service for:
- **Session Management**: Storing JWT tokens and session data for logged-in users.
- **Caching**: Caching frequently accessed data like user profiles to improve performance.

### Redis Configuration
Ensure that Redis is running and accessible using the following environment variables in your `.env` file:

```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_password
```

## Setup and Run

1. **Navigate to the Service Directory**:
    ```bash
    cd Q/auth-service
    ```

2. **Environment Configuration**:
    The repository includes a file named `.env-sample`. You need to rename this file to `.env` and fill in the required environment variables.

    ```bash
    mv .env-sample .env
    ```

    After renaming, update the variables according to your setup.

    Example:
    ```bash
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=username
    DB_PASS=password
    DB_NAME=authdb
    JWT_SECRET=your_secret_key
    REDIS_HOST=localhost
    REDIS_PORT=6379
    REDIS_PASSWORD=your_password
    ```

3. **Run the Service**:
    - **Using Docker**:
        ```bash
        docker-compose up --build
        ```
    - **Using Go**:
        ```bash
        go run cmd/main.go
        ```

## Middleware
- **JWT Authentication**: Protects user and admin routes.
- **Validation**: Custom validation for registration and login.

## Error Handling
Errors are returned in a consistent format with appropriate HTTP status codes.
