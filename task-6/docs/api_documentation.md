# Task Management REST API Documentation

## Overview
This is a robust Task Management REST API built with Go and Gin Framework. The API provides CRUD operations for managing tasks with MongoDB persistent storage.

## Base URL
```
http://localhost:8080
```

## MongoDB Configuration

The API uses MongoDB for persistent data storage. Configure the connection using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `MONGODB_URI` | `mongodb://localhost:27017` | MongoDB connection string |
| `MONGODB_DATABASE` | `taskmanager` | Database name |
| `MONGODB_COLLECTION` | `tasks` | Collection name |

## Task Model

Tasks are stored as MongoDB documents with the following structure:

```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "Complete project",
  "description": "Finish the task management API",
  "due_date": "2024-12-31T00:00:00Z",
  "status": "pending",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Note:** Task IDs are MongoDB ObjectIDs (24-character hexadecimal strings).

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

Retrieve a list of all tasks from MongoDB.

**Response (Success):**
```json
{
  "success": true,
  "message": "Tasks retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
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

**Response (Database Error):**
```json
{
  "success": false,
  "message": "Failed to retrieve tasks",
  "error": "database connection error"
}
```

### 3. Get Task by ID
**GET** `/api/v1/tasks/:id`

Retrieve details of a specific task by MongoDB ObjectID.

**Parameters:**
- `id` (path parameter): Task ObjectID (24-character hex string)

**Response (Success):**
```json
{
  "success": true,
  "message": "Task retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "title": "Complete project",
    "description": "Finish the task management API",
    "due_date": "2024-12-31T00:00:00Z",
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response (Invalid ID Format):**
```json
{
  "success": false,
  "message": "Task not found",
  "error": "invalid task ID format"
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

Create a new task in MongoDB.

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
    "id": "507f1f77bcf86cd799439011",
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

Update an existing task in MongoDB.

**Parameters:**
- `id` (path parameter): Task ObjectID (24-character hex string)

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
    "id": "507f1f77bcf86cd799439011",
    "title": "Updated task title",
    "description": "Updated description",
    "due_date": "2024-12-31T00:00:00Z",
    "status": "in_progress",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:45:00Z"
  }
}
```

**Response (Invalid ID):**
```json
{
  "success": false,
  "message": "Failed to update task",
  "error": "invalid task ID format"
}
```

### 6. Delete Task
**DELETE** `/api/v1/tasks/:id`

Delete a specific task from MongoDB.

**Parameters:**
- `id` (path parameter): Task ObjectID (24-character hex string)

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

**Response (Invalid ID):**
```json
{
  "success": false,
  "message": "Failed to delete task",
  "error": "invalid task ID format"
}
```

## Status Codes

- `200 OK`: Successful GET, PUT, DELETE operations
- `201 Created`: Successful POST operation
- `400 Bad Request`: Invalid request payload, parameters, or ObjectID format
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Database connection or server error

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

## MongoDB Integration Details

### Connection Management
- The API establishes a connection pool to MongoDB on startup
- Connection parameters are configurable via environment variables
- Graceful shutdown ensures proper connection cleanup

### Data Persistence
- All tasks are stored as BSON documents in MongoDB
- ObjectIDs are automatically generated for new tasks
- Timestamps are managed automatically (created_at, updated_at)

### Error Handling
- Database connection errors are handled gracefully
- Invalid ObjectID formats return appropriate error messages
- Network timeouts are configured with 10-second limits

## Testing with Postman

### Collection Setup
1. Create a new Postman collection named "Task Management API - MongoDB"
2. Set the base URL as a collection variable: `{{baseUrl}}` = `http://localhost:8080/api/v1`

### Sample Requests

#### 1. Create Task
- Method: POST
- URL: `{{baseUrl}}/tasks`
- Body (JSON):
```json
{
  "title": "Learn MongoDB with Go",
  "description": "Complete MongoDB integration tutorial",
  "due_date": "2024-02-15",
  "status": "pending"
}
```

#### 2. Get All Tasks
- Method: GET
- URL: `{{baseUrl}}/tasks`

#### 3. Get Task by ID
- Method: GET
- URL: `{{baseUrl}}/tasks/507f1f77bcf86cd799439011`

#### 4. Update Task
- Method: PUT
- URL: `{{baseUrl}}/tasks/507f1f77bcf86cd799439011`
- Body (JSON):
```json
{
  "title": "Learn MongoDB with Go - Updated",
  "description": "Complete MongoDB integration and add indexes",
  "due_date": "2024-02-20",
  "status": "in_progress"
}
```

#### 5. Delete Task
- Method: DELETE
- URL: `{{baseUrl}}/tasks/507f1f77bcf86cd799439011`

## Running the API

### Prerequisites
1. MongoDB installed and running (local or cloud)
2. Go 1.21 or higher

### Setup Steps
1. Navigate to the project directory
2. Install dependencies: `go mod tidy`
3. Set environment variables (optional):
   ```bash
   export MONGODB_URI="mongodb://localhost:27017"
   export MONGODB_DATABASE="taskmanager"
   export MONGODB_COLLECTION="tasks"
   ```
4. Run the application: `go run main.go`
5. The API will be available at `http://localhost:8080`

## Project Structure

```
task_manager/
├── main.go                    # Entry point with MongoDB connection
├── controllers/
│   └── task_controller.go     # HTTP request handlers
├── models/
│   └── task.go               # Data structures with BSON tags
├── data/
│   ├── database.go           # MongoDB connection utilities
│   └── task_service.go       # Business logic with MongoDB operations
├── router/
│   └── router.go             # Route configuration
├── docs/
│   └── api_documentation.md  # This documentation
└── go.mod                    # Go module with MongoDB driver
```

## MongoDB Collections

### Tasks Collection Schema
```javascript
{
  _id: ObjectId,           // MongoDB ObjectID (auto-generated)
  title: String,           // Required
  description: String,     // Optional
  due_date: Date,         // Optional
  status: String,         // Required (pending|in_progress|completed)
  created_at: Date,       // Auto-generated
  updated_at: Date        // Auto-updated
}
```

### Recommended Indexes
For better performance, consider creating these indexes:

```javascript
// Index on status for filtering
db.tasks.createIndex({ "status": 1 })

// Compound index for status and due_date
db.tasks.createIndex({ "status": 1, "due_date": 1 })

// Text index for searching titles and descriptions
db.tasks.createIndex({ 
  "title": "text", 
  "description": "text" 
})
```