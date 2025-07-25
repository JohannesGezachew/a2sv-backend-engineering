# MongoDB Integration Guide

## Overview
This Task Management API has been successfully integrated with MongoDB for persistent data storage. The integration replaces the previous in-memory storage with MongoDB using the official MongoDB Go Driver.

## Key Changes Made

### 1. Dependencies Added
- `go.mongodb.org/mongo-driver v1.13.1` - Official MongoDB Go Driver

### 2. Model Updates (`models/task.go`)
- Added MongoDB BSON tags to the Task struct
- Changed ID field from `int` to `primitive.ObjectID`
- Added proper JSON and BSON field mappings

### 3. Database Layer (`data/database.go`)
- Created database configuration utilities
- Added MongoDB connection management
- Environment variable support for configuration
- Connection pooling and error handling

### 4. Service Layer (`data/task_service.go`)
- Complete rewrite to use MongoDB operations
- Replaced in-memory storage with MongoDB collections
- Added proper context handling with timeouts
- Implemented CRUD operations using MongoDB Go Driver
- Added ObjectID validation and error handling

### 5. Controller Updates (`controllers/task_controller.go`)
- Updated to handle string IDs (ObjectIDs) instead of integers
- Enhanced error handling for invalid ObjectID formats
- Updated method signatures to work with new service layer

### 6. Router Updates (`router/router.go`)
- Added MongoDB connection initialization
- Integrated database configuration
- Added graceful shutdown support

### 7. Main Application (`main.go`)
- Added graceful server shutdown
- MongoDB connection lifecycle management
- Signal handling for clean shutdown

## Configuration

The application supports the following environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `MONGODB_URI` | `mongodb://localhost:27017` | MongoDB connection string |
| `MONGODB_DATABASE` | `taskmanager` | Database name |
| `MONGODB_COLLECTION` | `tasks` | Collection name |

## MongoDB Atlas Setup

### Step 1: Create MongoDB Atlas Account
1. Go to [MongoDB Atlas](https://www.mongodb.com/atlas)
2. Create a free account
3. Create a new cluster

### Step 2: Database User Setup
1. Go to Database Access in your Atlas dashboard
2. Add a new database user
3. Choose "Password" authentication
4. Set username and password
5. Grant "Read and write to any database" permissions

### Step 3: Network Access
1. Go to Network Access in your Atlas dashboard
2. Add IP Address
3. For development, you can use "0.0.0.0/0" (allow access from anywhere)
4. For production, specify your server's IP address

### Step 4: Get Connection String
1. Go to your cluster dashboard
2. Click "Connect"
3. Choose "Connect your application"
4. Copy the connection string
5. Replace `<password>` with your actual password

## Connection String Format

For MongoDB Atlas:
```
mongodb+srv://username:password@cluster.mongodb.net/database?retryWrites=true&w=majority
```

For local MongoDB:
```
mongodb://localhost:27017
```

## Testing the Integration

### 1. Set Environment Variables
```bash
# Windows PowerShell
$env:MONGODB_URI="your_connection_string_here"
$env:MONGODB_DATABASE="taskmanager"
$env:MONGODB_COLLECTION="tasks"

# Linux/macOS
export MONGODB_URI="your_connection_string_here"
export MONGODB_DATABASE="taskmanager"
export MONGODB_COLLECTION="tasks"
```

### 2. Run the Application
```bash
go run main.go
```

### 3. Test API Endpoints
```bash
# Test with the provided test script
go run test_api.go
```

## Data Migration

If you have existing data from the in-memory version, you'll need to:
1. Export your existing data
2. Transform integer IDs to ObjectIDs
3. Import into MongoDB

## Performance Considerations

### Recommended Indexes
```javascript
// Status index for filtering
db.tasks.createIndex({ "status": 1 })

// Date indexes for sorting and filtering
db.tasks.createIndex({ "created_at": -1 })
db.tasks.createIndex({ "due_date": 1 })

// Text search index
db.tasks.createIndex({ 
  "title": "text", 
  "description": "text" 
})

// Compound index for common queries
db.tasks.createIndex({ "status": 1, "due_date": 1 })
```

## Error Handling

The integration includes comprehensive error handling for:
- Database connection failures
- Invalid ObjectID formats
- Network timeouts
- Authentication errors
- Document not found scenarios

## Backward Compatibility

The API maintains the same endpoint structure and behavior:
- Same HTTP methods and URLs
- Same request/response formats
- Same validation rules
- Only change: Task IDs are now ObjectIDs instead of integers

## Troubleshooting

### Common Issues

1. **Authentication Failed**
   - Verify username and password
   - Check if user has proper permissions
   - Ensure password is URL-encoded if it contains special characters

2. **Connection Timeout**
   - Check network connectivity
   - Verify MongoDB Atlas IP whitelist
   - Check firewall settings

3. **Invalid ObjectID Format**
   - Ensure you're using 24-character hexadecimal strings
   - Use the ObjectIDs returned by the API

4. **Database Not Found**
   - MongoDB creates databases automatically when first document is inserted
   - Verify database name in connection string

## Security Best Practices

1. **Use Environment Variables**
   - Never hardcode connection strings
   - Use different credentials for different environments

2. **Network Security**
   - Restrict IP access in MongoDB Atlas
   - Use VPC peering for production

3. **Authentication**
   - Use strong passwords
   - Consider using certificate-based authentication

4. **Connection Pooling**
   - The driver handles connection pooling automatically
   - Configure pool size based on your needs

## Monitoring

Consider implementing:
- Connection health checks
- Query performance monitoring
- Error rate tracking
- Database metrics collection

## Future Enhancements

Potential improvements:
- Add database migrations
- Implement soft deletes
- Add audit logging
- Implement caching layer
- Add database backup strategies