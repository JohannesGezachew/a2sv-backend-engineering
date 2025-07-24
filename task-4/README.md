# Task Management REST API

A simple Task Management REST API built with Go and Gin Framework, featuring CRUD operations with an in-memory database.

## Features

- ✅ Create, Read, Update, Delete (CRUD) operations for tasks
- ✅ In-memory data storage with thread-safe operations
- ✅ Proper error handling and HTTP status codes
- ✅ JSON request/response format
- ✅ Input validation
- ✅ Clean architecture with separation of concerns
- ✅ Comprehensive API documentation

## Project Structure

```
task_manager/
├── main.go                    # Application entry point
├── controllers/
│   └── task_controller.go     # HTTP request handlers
├── models/
│   └── task.go               # Data structures and models
├── data/
│   └── task_service.go       # Business logic and data operations
├── router/
│   └── router.go             # Route configuration
├── docs/
│   └── api_documentation.md  # Detailed API documentation
├── test_api.go               # API testing script
├── go.mod                    # Go module dependencies
└── README.md                 # This file
```

## Prerequisites

- Go 1.21 or higher
- Git (optional, for cloning)

## Installation & Setup

1. **Navigate to the project directory:**
   ```bash
   cd task-4
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the application:**
   ```bash
   go run main.go
   ```

4. **The API will be available at:**
   ```
   http://localhost:8080
   ```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/tasks` | Get all tasks |
| GET | `/api/v1/tasks/:id` | Get task by ID |
| POST | `/api/v1/tasks` | Create new task |
| PUT | `/api/v1/tasks/:id` | Update existing task |
| DELETE | `/api/v1/tasks/:id` | Delete task |

## Task Model

```json
{
  "id": 1,
  "title": "Task title",
  "description": "Task description",
  "due_date": "2024-12-31T00:00:00Z",
  "status": "pending",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Valid Status Values
- `pending`
- `in_progress` 
- `completed`

## Testing the API

### Option 1: Using the Test Script

Run the included test script to verify all endpoints:

```bash
# Make sure the API server is running first
go run main.go

# In another terminal, run the test script
go run test_api.go
```

### Option 2: Using cURL

**Create a task:**
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go",
    "description": "Complete Go tutorial",
    "due_date": "2024-12-31",
    "status": "pending"
  }'
```

**Get all tasks:**
```bash
curl http://localhost:8080/api/v1/tasks
```

**Get task by ID:**
```bash
curl http://localhost:8080/api/v1/tasks/1
```

**Update a task:**
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go - Updated",
    "description": "Complete Go tutorial and build API",
    "due_date": "2024-12-31",
    "status": "in_progress"
  }'
```

**Delete a task:**
```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/1
```

### Option 3: Using Postman

1. Import the API endpoints into Postman
2. Set base URL: `http://localhost:8080/api/v1`
3. Test each endpoint with appropriate request bodies
4. Refer to `docs/api_documentation.md` for detailed examples

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { /* response data */ }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

## Error Handling

The API handles various error scenarios:

- **400 Bad Request**: Invalid request payload or parameters
- **404 Not Found**: Task not found
- **500 Internal Server Error**: Server errors

## Development Notes

- **Thread Safety**: The in-memory storage uses mutex locks for concurrent access
- **Validation**: Input validation for required fields and valid status values
- **Date Format**: Due dates should be in `YYYY-MM-DD` format
- **Auto-increment IDs**: Task IDs are automatically generated and incremented

## Future Enhancements

- Database integration (PostgreSQL, MySQL, etc.)
- Authentication and authorization
- Pagination for task listing
- Task filtering and sorting
- Task categories and priorities
- User management
- API rate limiting

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request
