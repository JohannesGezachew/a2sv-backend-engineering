package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"task_manager/controllers"
	"task_manager/data"
)

// SetupRouter initializes and configures the Gin router
func SetupRouter() *gin.Engine {
	// Create Gin router with default middleware (logger and recovery)
	router := gin.Default()

	// Get database configuration
	dbConfig := data.GetDatabaseConfig()

	// Connect to MongoDB
	client, err := data.ConnectToMongoDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Initialize services and controllers
	taskService := data.NewTaskService(client, dbConfig.Database, dbConfig.Collection)
	taskController := controllers.NewTaskController(taskService)

	// API versioning group
	v1 := router.Group("/api/v1")
	{
		// Task routes
		tasks := v1.Group("/tasks")
		{
			tasks.GET("", taskController.GetAllTasks)       // GET /api/v1/tasks
			tasks.GET("/:id", taskController.GetTaskByID)   // GET /api/v1/tasks/:id
			tasks.POST("", taskController.CreateTask)       // POST /api/v1/tasks
			tasks.PUT("/:id", taskController.UpdateTask)    // PUT /api/v1/tasks/:id
			tasks.DELETE("/:id", taskController.DeleteTask) // DELETE /api/v1/tasks/:id
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Task Management API is running",
		})
	})

	return router
}

// SetupRouterWithClient initializes router with existing MongoDB client (useful for testing)
func SetupRouterWithClient(client *mongo.Client, dbConfig *data.DatabaseConfig) *gin.Engine {
	router := gin.Default()

	taskService := data.NewTaskService(client, dbConfig.Database, dbConfig.Collection)
	taskController := controllers.NewTaskController(taskService)

	v1 := router.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.GET("", taskController.GetAllTasks)
			tasks.GET("/:id", taskController.GetTaskByID)
			tasks.POST("", taskController.CreateTask)
			tasks.PUT("/:id", taskController.UpdateTask)
			tasks.DELETE("/:id", taskController.DeleteTask)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Task Management API is running",
		})
	})

	return router
}