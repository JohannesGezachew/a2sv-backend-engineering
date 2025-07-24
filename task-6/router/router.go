package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"task_manager/controllers"
	"task_manager/data"
	"task_manager/middleware"
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
	
	userService := data.NewUserService(client, dbConfig.Database)
	userController := controllers.NewUserController(userService)

	// API versioning group
	v1 := router.Group("/api/v1")
	{
		// Public authentication routes (no middleware required)
		v1.POST("/register", userController.Register) // POST /api/v1/register
		v1.POST("/login", userController.Login)       // POST /api/v1/login

		// Protected user routes (authentication required)
		userRoutes := v1.Group("/users")
		userRoutes.Use(middleware.AuthMiddleware())
		{
			userRoutes.GET("/profile", userController.GetProfile)                                    // GET /api/v1/users/profile
			userRoutes.GET("", middleware.AdminMiddleware(), userController.GetAllUsers)             // GET /api/v1/users (admin only)
			userRoutes.POST("/promote", middleware.AdminMiddleware(), userController.PromoteUser)    // POST /api/v1/users/promote (admin only)
		}

		// Protected task routes
		tasks := v1.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware()) // All task routes require authentication
		{
			// Read operations - accessible by all authenticated users (admin and regular users)
			tasks.GET("", middleware.UserMiddleware(), taskController.GetAllTasks)       // GET /api/v1/tasks
			tasks.GET("/:id", middleware.UserMiddleware(), taskController.GetTaskByID)   // GET /api/v1/tasks/:id
			
			// Write operations - accessible only by admins
			tasks.POST("", middleware.AdminMiddleware(), taskController.CreateTask)       // POST /api/v1/tasks (admin only)
			tasks.PUT("/:id", middleware.AdminMiddleware(), taskController.UpdateTask)    // PUT /api/v1/tasks/:id (admin only)
			tasks.DELETE("/:id", middleware.AdminMiddleware(), taskController.DeleteTask) // DELETE /api/v1/tasks/:id (admin only)
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
	
	userService := data.NewUserService(client, dbConfig.Database)
	userController := controllers.NewUserController(userService)

	v1 := router.Group("/api/v1")
	{
		// Public authentication routes (no middleware required)
		v1.POST("/register", userController.Register) // POST /api/v1/register
		v1.POST("/login", userController.Login)       // POST /api/v1/login

		// Protected user routes (authentication required)
		userRoutes := v1.Group("/users")
		userRoutes.Use(middleware.AuthMiddleware())
		{
			userRoutes.GET("/profile", userController.GetProfile)                                    // GET /api/v1/users/profile
			userRoutes.GET("", middleware.AdminMiddleware(), userController.GetAllUsers)             // GET /api/v1/users (admin only)
			userRoutes.POST("/promote", middleware.AdminMiddleware(), userController.PromoteUser)    // POST /api/v1/users/promote (admin only)
		}

		// Protected task routes
		tasks := v1.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware()) // All task routes require authentication
		{
			// Read operations - accessible by all authenticated users (admin and regular users)
			tasks.GET("", middleware.UserMiddleware(), taskController.GetAllTasks)       // GET /api/v1/tasks
			tasks.GET("/:id", middleware.UserMiddleware(), taskController.GetTaskByID)   // GET /api/v1/tasks/:id
			
			// Write operations - accessible only by admins
			tasks.POST("", middleware.AdminMiddleware(), taskController.CreateTask)       // POST /api/v1/tasks (admin only)
			tasks.PUT("/:id", middleware.AdminMiddleware(), taskController.UpdateTask)    // PUT /api/v1/tasks/:id (admin only)
			tasks.DELETE("/:id", middleware.AdminMiddleware(), taskController.DeleteTask) // DELETE /api/v1/tasks/:id (admin only)
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