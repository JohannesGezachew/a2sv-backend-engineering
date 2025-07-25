# Task Management API Documentation

## Overview

The Task Management API is a RESTful web service that provides comprehensive task management functionality with JWT-based authentication and role-based authorization. The API supports user registration, authentication, and task CRUD operations with different access levels for administrators and regular users.

## Base URL

```
http://localhost:8080
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. After successful login, include the JWT token in the Authorization header for protected endpoints:

```
Authorization: Bearer <your-jwt-token>
```

## User Roles

- **Admin**: Can create, update, delete tasks, view all users, and promote users
- **User**: Can view tasks and their own profile
- **First User**: The first registered user automatically becomes an admin

## API Endpoints

### Public Endpoints (No Authentication Required)

#### 1. Health Check
**GET** `/health`

Check if the API server is running.

**Response:**
```json
{
  "status": "OK",
  "message": "Task Management API is running"
}
```

#### 2. User Registration
**POST** `/api/v1/register`

Register a new user account.

**Request Body:**
```json
{
  "username": "john_doe",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "60f7b3b3b3b3b3b3b3b3b3b3",
    "username": "john_doe",
    "role": "admin",
    "created_at": "2025-07-25T15:30:00Z",
    "updated_at": "2025-07-25T15:30:00Z"
  }
}
```

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "message": "Failed to create user",
  "error": "username already exists"
}
```

#### 3. User Login
**POST** `/api/v1/login`

Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "username": "john_doe",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "60f7b3b3b3b3b3b3b3b3b3b3",
    "username": "john_doe",
    "role": "admin",
    "created_at": "2025-07-25T15:30:00Z",
    "updated_at": "2025-07-25T15:30:00Z"
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "message": "Authentication failed",
  "error": "invalid credentials"
}
```

### Protected Endpoints (Authentication Required)

#### 4. Get User Profile
**GET** `/api/v1/users/profile`

Get the authenticated user's profile information.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "60f7b3b3b3b3b3b3b3b3b3b3",
    "username": "john_doe",
    "role": "admin",
    "created_at": "2025-07-25T15:30:00Z",
    "updated_at": "2025-07-25T15:30:00Z"
  }
}
```

#### 5. Get All Tasks
**GET** `/api/v1/tasks`

Retrieve all tasks (available to all authenticated users).

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Tasks retrieved successfully",
  "data": [
    {
      "id": "60f7b3b3b3b3b3b3b3b3b3b4",
      "title": "Complete Project Documentation",
      "description": "Write comprehensive API documentation",
      "due_date": "2025-08-01T00:00:00Z",
      "status": "pending",
      "created_at": "2025-07-25T15:30:00Z",
      "updated_at": "2025-07-25T15:30:00Z"
    }
  ]
}
```

#### 6. Get Task by ID
**GET** `/api/v1/tasks/:id`

Retrieve a specific task by its ID (available to all authenticated users).

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Task retrieved successfully",
  "data": {
    "id": "60f7b3b3b3b3b3b3b3b3b3b4",
    "title": "Complete Project Documentation",
    "description": "Write comprehensive API documentation",
    "due_date": "2025-08-01T00:00:00Z",
    "status": "pending",
    "created_at": "2025-07-25T15:30:00Z",
    "updated_at": "2025-07-25T15:30:00Z"
  }
}
```

### Admin-Only Endpoints (Admin Role Required)

#### 7. Create Task
**POST** `/api/v1/tasks`

Create a new task (admin only).

**Headers:**
```
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "Complete Project Documentation",
  "description": "Write comprehensive API documentation",
  "status": "pending",
  "due_date": "2025-08-01"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Task created successfully",
  "data": {
    "id": "60f7b3b3b3b3b3b3b3b3b3b4",
    "title": "Complete Project Documentation",
    "description": "Write comprehensive API documentation",
    "due_date": "2025-08-01T00:00:00Z",
    "status": "pending",
    "created_at": "2025-07-25T15:30:00Z",
    "updated_at": "2025-07-25T15:30:00Z"
  }
}
```

#### 8. Update Task
**PUT** `/api/v1/tasks/:id`

Update an existing task (admin only).

**Headers:**
```
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "Updated Task Title",
  "description": "Updated task description",
  "status": "in_progress",
  "due_date": "2025-08-05"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Task updated successfully",
  "data": {
    "id": "60f7b3b3b3b3b3b3b3b3b3b4",
    "title": "Updated Task Title",
    "description": "Updated task description",
    "due_date": "2025-08-05T00:00:00Z",
    "status": "in_progress",
    "created_at": "2025-07-25T15:30:00Z",
    "updated_at": "2025-07-25T15:35:00Z"
  }
}
```

#### 9. Delete Task
**DELETE** `/api/v1/tasks/:id`

Delete a task (admin only).

**Headers:**
```
Authorization: Bearer <admin-jwt-token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Task deleted successfully"
}
```

#### 10. Get All Users
**GET** `/api/v1/users`

Retrieve all users (admin only).

**Headers:**
```
Authorization: Bearer <admin-jwt-token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "60f7b3b3b3b3b3b3b3b3b3b3",
      "username": "john_doe",
      "role": "admin",
      "created_at": "2025-07-25T15:30:00Z",
      "updated_at": "2025-07-25T15:30:00Z"
    },
    {
      "id": "60f7b3b3b3b3b3b3b3b3b3b5",
      "username": "jane_smith",
      "role": "user",
      "created_at": "2025-07-25T15:35:00Z",
      "updated_at": "2025-07-25T15:35:00Z"
    }
  ]
}
```

#### 11. Promote User
**POST** `/api/v1/users/promote`

Promote a user to admin role (admin only).

**Headers:**
```
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "username": "jane_smith"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User promoted to admin successfully",
  "data": {
    "id": "60f7b3b3b3b3b3b3b3b3b3b5",
    "username": "jane_smith",
    "role": "admin",
    "created_at": "2025-07-25T15:35:00Z",
    "updated_at": "2025-07-25T15:40:00Z"
  }
}
```

## Task Status Values

Tasks can have one of the following status values:
- `pending` - Task is pending
- `in_progress` - Task is currently being worked on
- `completed` - Task has been completed

## Error Responses

### Authentication Errors

**401 Unauthorized - Missing Token:**
```json
{
  "success": false,
  "message": "Authorization header required",
  "error": "Missing Authorization header"
}
```

**401 Unauthorized - Invalid Token:**
```json
{
  "success": false,
  "message": "Invalid or expired token",
  "error": "token is malformed: token contains an invalid number of segments"
}
```

**401 Unauthorized - Malformed Header:**
```json
{
  "success": false,
  "message": "Invalid authorization header format",
  "error": "Authorization header must be in format: Bearer <token>"
}
```

### Authorization Errors

**403 Forbidden - Admin Required:**
```json
{
  "success": false,
  "message": "Access denied",
  "error": "Admin privileges required"
}
```

### Validation Errors

**400 Bad Request - Invalid Input:**
```json
{
  "success": false,
  "message": "Invalid request payload",
  "error": "Key: 'TaskRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag"
}
```

**404 Not Found - Resource Not Found:**
```json
{
  "success": false,
  "message": "Task not found",
  "error": "task not found"
}
```

## Usage Examples

### 1. Complete User Registration and Login Flow

```bash
# 1. Register a new user
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username": "john_doe", "password": "password123"}'

# 2. Login to get JWT token
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "john_doe", "password": "password123"}'

# 3. Use the token for protected endpoints
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

### 2. Task Management Flow (Admin)

```bash
# 1. Create a task (admin only)
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete Project Documentation",
    "description": "Write comprehensive API documentation",
    "status": "pending",
    "due_date": "2025-08-01"
  }'

# 2. Get all tasks
curl -X GET http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 3. Update a task (admin only)
curl -X PUT http://localhost:8080/api/v1/tasks/TASK_ID \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Task Title",
    "description": "Updated description",
    "status": "in_progress",
    "due_date": "2025-08-05"
  }'

# 4. Delete a task (admin only)
curl -X DELETE http://localhost:8080/api/v1/tasks/TASK_ID \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

### 3. User Management Flow (Admin)

```bash
# 1. Get all users (admin only)
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"

# 2. Promote a user to admin (admin only)
curl -X POST http://localhost:8080/api/v1/users/promote \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username": "jane_smith"}'
```

## Security Features

1. **Password Hashing**: User passwords are hashed using bcrypt before storage
2. **JWT Authentication**: Secure token-based authentication with expiration
3. **Role-Based Authorization**: Different access levels for admin and regular users
4. **Protected Routes**: All sensitive operations require valid authentication
5. **Input Validation**: Request payloads are validated for required fields and formats
6. **Error Handling**: Comprehensive error responses with appropriate HTTP status codes

## Environment Variables

The API uses the following environment variables:

```env
# MongoDB Configuration
MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority
MONGODB_DATABASE=taskmanager
MONGODB_COLLECTION=tasks

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here-make-it-long-and-random

# Server Configuration
PORT=8080
```

## Testing

To test the API, ensure the server is running:

```bash
go run main.go
```

Then run the comprehensive test suite:

```bash
go run test_api.go
```

The test suite covers:
- Health check functionality
- User registration and login
- JWT token validation
- Protected route access
- Task CRUD operations
- Role-based authorization
- Error handling scenarios

## Notes

- The first registered user automatically becomes an admin
- JWT tokens expire after 24 hours
- All timestamps are in ISO 8601 format
- Task due dates should be in YYYY-MM-DD format
- The API uses MongoDB for data persistence
- All responses include success/failure indicators and descriptive messages