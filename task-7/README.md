# Task Management API - Clean Architecture

This project demonstrates a Task Management API implemented using Clean Architecture principles in Go. The application is organized into distinct layers with clear separation of concerns, making it maintainable, testable, and scalable.

## Architecture Overview

The application follows Clean Architecture principles with the following layers:

### 1. Domain Layer (`Domain/`)
- **Purpose**: Contains core business entities and rules
- **Files**: `domain.go`
- **Responsibilities**:
  - Define core business entities (Task, User)
  - Define request/response models
  - Define business constants and validation rules
  - Remain independent of external frameworks

### 2. Use Cases Layer (`Usecases/`)
- **Purpose**: Contains application-specific business logic
- **Files**: `task_usecases.go`, `user_usecases.go`
- **Responsibilities**:
  - Orchestrate interactions between different layers
  - Implement business rules and validation
  - Define interfaces for external dependencies
  - Coordinate data flow between repositories and controllers

### 3. Repository Layer (`Repositories/`)
- **Purpose**: Abstracts data access logic
- **Files**: `task_repository.go`, `user_repository.go`
- **Responsibilities**:
  - Define interfaces for data access operations
  - Implement MongoDB-specific data access logic
  - Handle data persistence and retrieval
  - Abstract database operations from business logic

### 4. Infrastructure Layer (`Infrastructure/`)
- **Purpose**: Implements external dependencies and services
- **Files**: `auth_middleware.go`, `jwt_service.go`, `password_service.go`
- **Responsibilities**:
  - Handle authentication and authorization
  - Implement JWT token generation and validation
  - Provide password hashing and comparison services
  - Manage external service integrations

### 5. Delivery Layer (`Delivery/`)
- **Purpose**: Handles incoming requests and responses
- **Files**: `main.go`, `controllers/controller.go`, `routers/router.go`
- **Responsibilities**:
  - Handle HTTP requests and responses
  - Route requests to appropriate use cases
  - Manage server lifecycle and configuration
  - Provide API endpoints

## Key Clean Architecture Principles Implemented

### 1. Dependency Inversion
- Higher-level layers do not depend on lower-level layers
- Dependencies point inward toward the domain
- Interfaces are defined in higher layers and implemented in lower layers

### 2. Separation of Concerns
- Each layer has a single, well-defined responsibility
- Business logic is separated from infrastructure concerns
- Data access is abstracted from business rules

### 3. Independence of Frameworks
- Core business logic is independent of external frameworks
- Database and web framework can be easily replaced
- Business rules remain stable regardless of external changes

### 4. Testability
- Each layer can be tested independently
- Dependencies can be easily mocked using interfaces
- Business logic can be tested without external dependencies

## Project Structure

```
task-7/
├── Delivery/
│   ├── main.go                 # Application entry point
│   ├── controllers/
│   │   └── controller.go       # HTTP request handlers
│   └── routers/
│       └── router.go           # Route configuration
├── Domain/
│   └── domain.go               # Core business entities and models
├── Infrastructure/
│   ├── auth_middleware.go      # Authentication middleware
│   ├── jwt_service.go          # JWT token service
│   └── password_service.go     # Password hashing service
├── Repositories/
│   ├── task_repository.go      # Task data access interface and implementation
│   └── user_repository.go      # User data access interface and implementation
├── Usecases/
│   ├── task_usecases.go        # Task business logic
│   └── user_usecases.go        # User business logic
├── .env.example                # Environment variables template
├── .gitignore                  # Git ignore rules
├── go.mod                      # Go module dependencies
├── go.sum                      # Go module checksums
├── README.md                   # Project documentation
└── test_clean_architecture.go  # Clean Architecture demonstration
```

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
The API uses JWT (JSON Web Token) for authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### Response Format
All API responses follow a consistent format:

#### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... }
}
```

#### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

---

## Authentication Endpoints

### 1. Register User
**POST** `/api/v1/register`

Creates a new user account. The first user registered automatically becomes an admin.

**Request Body:**
```json
{
  "username": "string (required, unique)",
  "password": "string (required, min 6 characters)"
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "username": "johndoe",
    "role": "admin",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request payload
- `409 Conflict`: Username already exists

---

### 2. Login User
**POST** `/api/v1/login`

Authenticates a user and returns a JWT token.

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "username": "johndoe",
    "role": "admin",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request payload
- `401 Unauthorized`: Invalid credentials

---

## User Management Endpoints

### 3. Get User Profile
**GET** `/api/v1/users/profile`

Returns the profile of the authenticated user.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "username": "johndoe",
    "role": "admin",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: User not found

---

### 4. Get All Users (Admin Only)
**GET** `/api/v1/users`

Returns a list of all users. Requires admin privileges.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "username": "johndoe",
      "role": "admin",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "507f1f77bcf86cd799439012",
      "username": "janedoe",
      "role": "user",
      "created_at": "2024-01-15T11:00:00Z",
      "updated_at": "2024-01-15T11:00:00Z"
    }
  ]
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin privileges required

---

### 5. Promote User to Admin (Admin Only)
**POST** `/api/v1/users/promote`

Promotes a regular user to admin role. Requires admin privileges.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "username": "string (required)"
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/users/promote \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "username": "janedoe"
  }'
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User promoted to admin successfully",
  "data": {
    "id": "507f1f77bcf86cd799439012",
    "username": "janedoe",
    "role": "admin",
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T12:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request payload or user already admin
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin privileges required
- `404 Not Found`: User not found

---

## Task Management Endpoints

### 6. Get All Tasks
**GET** `/api/v1/tasks`

Returns a list of all tasks. Accessible by authenticated users.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Tasks retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439013",
      "title": "Complete project documentation",
      "description": "Write comprehensive API documentation",
      "due_date": "2024-12-31T23:59:59Z",
      "status": "pending",
      "created_at": "2024-01-15T12:00:00Z",
      "updated_at": "2024-01-15T12:00:00Z"
    },
    {
      "id": "507f1f77bcf86cd799439014",
      "title": "Review code changes",
      "description": "Review pull requests from team members",
      "due_date": "2024-01-20T17:00:00Z",
      "status": "in_progress",
      "created_at": "2024-01-15T13:00:00Z",
      "updated_at": "2024-01-15T14:00:00Z"
    }
  ]
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Database error

---

### 7. Get Task by ID
**GET** `/api/v1/tasks/{id}`

Returns a specific task by its ID. Accessible by authenticated users.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Path Parameters:**
- `id`: MongoDB ObjectId of the task

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/v1/tasks/507f1f77bcf86cd799439013 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Task retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439013",
    "title": "Complete project documentation",
    "description": "Write comprehensive API documentation",
    "due_date": "2024-12-31T23:59:59Z",
    "status": "pending",
    "created_at": "2024-01-15T12:00:00Z",
    "updated_at": "2024-01-15T12:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid task ID format
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: Task not found

---

### 8. Create Task (Admin Only)
**POST** `/api/v1/tasks`

Creates a new task. Requires admin privileges.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "title": "string (required)",
  "description": "string (optional)",
  "due_date": "string (optional, format: YYYY-MM-DD)",
  "status": "string (required, one of: pending, in_progress, completed)"
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "title": "Implement user authentication",
    "description": "Add JWT-based authentication to the API",
    "due_date": "2024-02-15",
    "status": "pending"
  }'
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Task created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439015",
    "title": "Implement user authentication",
    "description": "Add JWT-based authentication to the API",
    "due_date": "2024-02-15T00:00:00Z",
    "status": "pending",
    "created_at": "2024-01-15T15:00:00Z",
    "updated_at": "2024-01-15T15:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request payload or invalid status/date format
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin privileges required

---

### 9. Update Task (Admin Only)
**PUT** `/api/v1/tasks/{id}`

Updates an existing task. Requires admin privileges.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Path Parameters:**
- `id`: MongoDB ObjectId of the task

**Request Body:**
```json
{
  "title": "string (required)",
  "description": "string (optional)",
  "due_date": "string (optional, format: YYYY-MM-DD)",
  "status": "string (required, one of: pending, in_progress, completed)"
}
```

**Example Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/507f1f77bcf86cd799439015 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "title": "Implement user authentication",
    "description": "Add JWT-based authentication to the API - Updated",
    "due_date": "2024-02-20",
    "status": "in_progress"
  }'
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Task updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439015",
    "title": "Implement user authentication",
    "description": "Add JWT-based authentication to the API - Updated",
    "due_date": "2024-02-20T00:00:00Z",
    "status": "in_progress",
    "created_at": "2024-01-15T15:00:00Z",
    "updated_at": "2024-01-15T16:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request payload, task ID format, or invalid status/date format
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin privileges required
- `404 Not Found`: Task not found

---

### 10. Delete Task (Admin Only)
**DELETE** `/api/v1/tasks/{id}`

Deletes a task by its ID. Requires admin privileges.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Path Parameters:**
- `id`: MongoDB ObjectId of the task

**Example Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/507f1f77bcf86cd799439015 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Task deleted successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid task ID format
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin privileges required
- `404 Not Found`: Task not found

---

## Health Check Endpoint

### 11. Health Check
**GET** `/health`

Returns the API health status. No authentication required.

**Example Request:**
```bash
curl -X GET http://localhost:8080/health
```

**Response (200 OK):**
```json
{
  "status": "OK",
  "message": "Task Management API is running"
}
```

---

## Data Models

### User Model
```json
{
  "id": "MongoDB ObjectId",
  "username": "string (unique)",
  "password": "string (hashed, not returned in responses)",
  "role": "admin|user",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Task Model
```json
{
  "id": "MongoDB ObjectId",
  "title": "string",
  "description": "string",
  "due_date": "timestamp",
  "status": "pending|in_progress|completed",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

---

## Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request payload or parameters
- `401 Unauthorized`: Authentication required or invalid token
- `403 Forbidden`: Insufficient privileges
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists
- `500 Internal Server Error`: Server error

---

## Authentication Flow

1. **Register** a new user account using `/api/v1/register`
2. **Login** with credentials using `/api/v1/login` to receive a JWT token
3. **Include** the JWT token in the Authorization header for protected endpoints
4. **Token expires** after 24 hours, requiring re-authentication

---

## Role-Based Access Control

### Admin Users
- Can perform all operations
- Can create, update, and delete tasks
- Can view all users and promote users to admin
- First registered user automatically becomes admin

### Regular Users
- Can view their own profile
- Can view all tasks
- Cannot create, update, or delete tasks
- Cannot access user management endpoints

---

## Complete API Workflow Example

Here's a complete example of using the API:

### Step 1: Register a User
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

### Step 2: Login to Get Token
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

### Step 3: Create a Task (using token from login)
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Complete project",
    "description": "Finish the task management API",
    "due_date": "2024-12-31",
    "status": "pending"
  }'
```

### Step 4: Get All Tasks
```bash
curl -X GET http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Step 5: Update Task Status
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/TASK_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Complete project",
    "description": "Finish the task management API",
    "due_date": "2024-12-31",
    "status": "completed"
  }'
```

## Setup and Installation

### Prerequisites
- Go 1.21 or higher
- MongoDB (local or MongoDB Atlas)

### Installation Steps

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd task-7
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your MongoDB connection details
   ```

4. **Run the application**
   ```bash
   cd Delivery
   go run main.go
   ```

5. **Test the Clean Architecture implementation**
   ```bash
   go run test_clean_architecture.go
   ```

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=taskmanager
MONGODB_COLLECTION=tasks

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here-make-it-long-and-random

# Server Configuration
PORT=8080
```

## Clean Architecture Benefits Demonstrated

### 1. **Maintainability**
- Clear separation of concerns makes code easier to understand and modify
- Changes in one layer don't affect other layers
- Business logic is centralized and reusable

### 2. **Testability**
- Each layer can be tested independently
- Dependencies can be easily mocked
- Business logic can be tested without external dependencies

### 3. **Scalability**
- New features can be added without affecting existing code
- Different delivery mechanisms can be easily added
- Database can be changed without affecting business logic

### 4. **Flexibility**
- External dependencies can be easily replaced
- Different authentication mechanisms can be plugged in
- API versioning is straightforward

## Usage Examples

### 1. Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'
```

### 3. Create a task (admin only)
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{"title": "New Task", "description": "Task description", "due_date": "2024-12-31", "status": "pending"}'
```

### 4. Get all tasks
```bash
curl -X GET http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Testing

The project includes a comprehensive test file (`test_clean_architecture.go`) that demonstrates:
- Layer interaction and dependency flow
- Authentication and authorization
- CRUD operations through all layers
- Error handling and validation

Run the test with:
```bash
go run test_clean_architecture.go
```

## Design Decisions

### 1. **Interface-Based Design**
- All external dependencies are abstracted behind interfaces
- Enables easy testing and dependency injection
- Supports the Dependency Inversion Principle

### 2. **Layer Communication**
- Each layer only communicates with adjacent layers
- Dependencies point inward toward the domain
- No circular dependencies between layers

### 3. **Error Handling**
- Consistent error handling across all layers
- Business errors are handled in use cases
- Infrastructure errors are handled in repositories

### 4. **Security**
- JWT-based authentication
- Role-based authorization (admin/user)
- Password hashing using bcrypt

## Future Enhancements

1. **Add comprehensive unit tests** for each layer
2. **Implement caching layer** for improved performance
3. **Add logging and monitoring** capabilities
4. **Implement API rate limiting**
5. **Add database migrations** for schema management
6. **Implement event-driven architecture** for complex workflows

## Contributing

1. Follow Clean Architecture principles
2. Maintain clear separation of concerns
3. Write comprehensive tests
4. Document any architectural decisions
5. Ensure backward compatibility

## License

This project is licensed under the MIT License.