# Task Management REST API

A robust Task Management REST API built with Go and Gin Framework, featuring CRUD operations with MongoDB persistent storage.

## Features

- ✅ Create, Read, Update, Delete (CRUD) operations for tasks
- ✅ MongoDB persistent data storage
- ✅ Proper error handling and HTTP status codes
- ✅ JSON request/response format
- ✅ Input validation
- ✅ Clean architecture with separation of concerns
- ✅ Graceful server shutdown
- ✅ Environment-based configuration
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
- MongoDB 4.4 or higher (local installation or cloud service)
- Git (optional, for cloning)

## MongoDB Setup

### Option 1: Local MongoDB Installation

1. **Install MongoDB Community Edition:**
   - Windows: Download from [MongoDB Download Center](https://www.mongodb.com/try/download/community)
   - macOS: `brew install mongodb-community`
   - Linux: Follow [official installation guide](https://docs.mongodb.com/manual/installation/)

2. **Start MongoDB service:**
   ```bash
   # Windows (as service)
   net start MongoDB
   
   # macOS/Linux
   brew services start mongodb-community
   # or
   sudo systemctl start mongod
   ```

### Option 2: Docker (Recommended for Development)

1. **Install Docker and Docker Compose**
2. **Start MongoDB using Docker Compose:**
   ```bash
   docker-compose up -d
   ```
3. **MongoDB will be available at:**
   ```
   mongodb://localhost:27017
   ```

### Option 3: MongoDB Atlas (Cloud)

1. Create a free account at [MongoDB Atlas](https://www.mongodb.com/atlas)
2. Create a new cluster
3. Get your connection string
4. Set the `MONGODB_URI` environment variable

## Installation & Setup

1. **Navigate to the project directory:**
   ```bash
   cd task-5
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Configure environment variables (optional):**
   ```bash
   # Default values are used if not set
   export MONGODB_URI="mongodb://localhost:27017"
   export MONGODB_DATABASE="taskmanager"
   export MONGODB_COLLECTION="tasks"
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```

5. **The API will be available at:**
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
  "id": "507f1f77bcf86cd799439011",
  "title": "Task title",
  "description": "Task description",
  "due_date": "2024-12-31T00:00:00Z",
  "status": "pending",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Note:** Task IDs are now MongoDB ObjectIDs (24-character hexadecimal strings) instead of integers.

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
curl http://localhost:8080/api/v1/tasks/507f1f77bcf86cd799439011
```

**Update a task:**
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/507f1f77bcf86cd799439011 \
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
curl -X DELETE http://localhost:8080/api/v1/tasks/507f1f77bcf86cd799439011
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

- **MongoDB Integration**: Uses MongoDB Go Driver for persistent data storage
- **ObjectID**: Task IDs are MongoDB ObjectIDs (24-character hex strings)
- **Validation**: Input validation for required fields and valid status values
- **Date Format**: Due dates should be in `YYYY-MM-DD` format
- **Error Handling**: Comprehensive error handling for database operations
- **Graceful Shutdown**: Server handles shutdown signals gracefully
- **Environment Configuration**: Database connection configurable via environment variables

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MONGODB_URI` | `mongodb://localhost:27017` | MongoDB connection string |
| `MONGODB_DATABASE` | `taskmanager` | Database name |
| `MONGODB_COLLECTION` | `tasks` | Collection name |

## Future Enhancements

- Authentication and authorization
- Pagination for task listing
- Task filtering and sorting
- Task categories and priorities
- User management
- API rate limiting
- MongoDB indexes for better performance
- Data validation with MongoDB schema validation
- Aggregation pipelines for advanced queries

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is for educational purposes as part of a Go programming task.