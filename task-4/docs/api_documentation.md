# Task Management REST API Documentation

## Overview
This is a simple Task Management REST API built with Go and Gin Framework. The API provides CRUD operations for managing tasks with an in-memory database.

## Base URL
```
http://localhost:8080
```

## Endpoints

### 1. Health Check
**GET** `/health`

Check if the API is running.

**Response:**
```json
{
  "status": "OK",
  "message": "Task Management API is running"
}
```

### 2. Get All Tasks
**GET** `/api/v1/tasks`

Retrieve a list of all tasks.

**Response:**
```json
{
  "success": true,
  "message": "Tasks retrieved successfully",
  "data": [
    {
      "id": 1,
      "title": "Complete project",
      "description": "Finish the task management API",
      "due_date": "2024-12-31T00:00:00Z",
      "status": "pending",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### 3. Get Task by ID
**GET** `/api/v1/tasks/:id`

Retrieve details of a specific task.

**Parameters:**
- `id` (path parameter): Task ID (integer)

**Response (Success):**
```json
{
  "success": true,
  "message": "Task retrieved successfully",
  "data": {
    "id": 1,
    "title": "Complete project",
    "description": "Finish the task management API",
    "due_date": "2024-12-31T00:00:00Z",
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response (Not Found):**
```json
{
  "success": false,
  "message": "Task not found",
  "error": "task not found"
}
```

### 4. Create New Task
**POST** `/api/v1/tasks`

Create a new task.

**Request Body:**
```json
{
  "title": "Complete project",
  "description": "Finish the task management API",
  "due_date": "2024-12-31",
  "status": "pending"
}
```

**Required Fields:**
- `title` (string): Task title
- `status` (string): Task status (pending, in_progress, completed)

**Optional Fields:**
- `description` (string): Task description
- `due_date` (string): Due date in YYYY-MM-DD format

**Response (Success):**
```json
{
  "success": true,
  "message": "Task created successfully",
  "data": {
    "id": 1,
    "title": "Complete project",
    "description": "Finish the task management API",
    "due_date": "2024-12-31T00:00:00Z",
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response (Validation Error):**
```json
{
  "success": false,
  "message": "Invalid request payload",
  "error": "Key: 'TaskRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag"
}
```

### 5. Update Task
**PUT** `/api/v1/tasks/:id`

Update an existing task.

**Parameters:**
- `id` (path parameter): Task ID (integer)

**Request Body:**
```json
{
  "title": "Updated task title",
  "description": "Updated description",
  "due_date": "2024-12-31",
  "status": "in_progress"
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Task updated successfully",
  "data": {
    "id": 1,
    "title": "Updated task title",
    "description": "Updated description",
    "due_date": "2024-12-31T00:00:00Z",
    "status": "in_progress",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:45:00Z"
  }
}
```

### 6. Delete Task
**DELETE** `/api/v1/tasks/:id`

Delete a specific task.

**Parameters:**
- `id` (path parameter): Task ID (integer)

**Response (Success):**
```json
{
  "success": true,
  "message": "Task deleted successfully"
}
```

**Response (Not Found):**
```json
{
  "success": false,
  "message": "Failed to delete task",
  "error": "task not found"
}
```

## Status Codes

- `200 OK`: Successful GET, PUT operations
- `201 Created`: Successful POST operation
- `400 Bad Request`: Invalid request payload or parameters
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Valid Task Statuses

- `pending`: Task is pending
- `in_progress`: Task is in progress
- `completed`: Task is completed

## Error Response Format

All error responses follow this format:
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

## Testing with Postman

### Collection Setup
1. Create a new Postman collection named "Task Management API"
2. Set the base URL as a collection variable: `{{baseUrl}}` = `http://localhost:8080/api/v1`

### Sample Requests

#### 1. Create Task
- Method: POST
- URL: `{{baseUrl}}/tasks`
- Body (JSON):
```json
{
  "title": "Learn Go Programming",
  "description": "Complete Go tutorial and build a REST API",
  "due_date": "2024-02-15",
  "status": "pending"
}
```

#### 2. Get All Tasks
- Method: GET
- URL: `{{baseUrl}}/tasks`

#### 3. Get Task by ID
- Method: GET
- URL: `{{baseUrl}}/tasks/1`

#### 4. Update Task
- Method: PUT
- URL: `{{baseUrl}}/tasks/1`
- Body (JSON):
```json
{
  "title": "Learn Go Programming - Updated",
  "description": "Complete Go tutorial, build REST API, and add tests",
  "due_date": "2024-02-20",
  "status": "in_progress"
}
```

#### 5. Delete Task
- Method: DELETE
- URL: `{{baseUrl}}/tasks/1`

## Running the API

1. Navigate to the project directory
2. Install dependencies: `go mod tidy`
3. Run the application: `go run main.go`
4. The API will be available at `http://localhost:8080`

## Project Structure

```
task_manager/
├── main.go                    # Entry point
├── controllers/
│   └── task_controller.go     # HTTP request handlers
├── models/
│   └── task.go               # Data structures
├── data/
│   └── task_service.go       # Business logic
├── router/
│   └── router.go             # Route configuration
├── docs/
│   └── api_documentation.md  # This documentation
└── go.mod                    # Go module file
```