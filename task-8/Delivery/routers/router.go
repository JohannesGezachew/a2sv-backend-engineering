package routers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	
	"task_manager/Delivery/controllers"
	"task_manager/Infrastructure"
	"task_manager/Repositories"
	"task_manager/Usecases"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URI        string
	Database   string
	Collection string
}

// SetupRouter initializes and configures the Gin router with Clean Architecture
func SetupRouter(client *mongo.Client, dbConfig *DatabaseConfig) *gin.Engine {
	router := gin.Default()

	// Initialize Infrastructure layer
	passwordService := Infrastructure.NewPasswordService()
	jwtService := Infrastructure.NewJWTService()
	authMiddleware := Infrastructure.NewAuthMiddleware(jwtService)

	// Initialize Repository layer
	taskRepo := Repositories.NewTaskRepository(client, dbConfig.Database, dbConfig.Collection)
	userRepo := Repositories.NewUserRepository(client, dbConfig.Database)

	// Initialize Usecase layer
	taskUsecase := Usecases.NewTaskUsecase(taskRepo)
	userUsecase := Usecases.NewUserUsecase(userRepo, passwordService, jwtService)

	// Initialize Controller layer
	controller := controllers.NewController(taskUsecase, userUsecase)

	// API versioning group
	v1 := router.Group("/api/v1")
	{
		// Public authentication routes (no middleware required)
		v1.POST("/register", controller.Register) // POST /api/v1/register
		v1.POST("/login", controller.Login)       // POST /api/v1/login

		// Protected user routes (authentication required)
		userRoutes := v1.Group("/users")
		userRoutes.Use(authMiddleware.AuthenticateToken())
		{
			userRoutes.GET("/profile", controller.GetProfile)                                          // GET /api/v1/users/profile
			userRoutes.GET("", authMiddleware.RequireAdmin(), controller.GetAllUsers)                  // GET /api/v1/users (admin only)
			userRoutes.POST("/promote", authMiddleware.RequireAdmin(), controller.PromoteUser)         // POST /api/v1/users/promote (admin only)
		}

		// Protected task routes
		tasks := v1.Group("/tasks")
		tasks.Use(authMiddleware.AuthenticateToken()) // All task routes require authentication
		{
			// Read operations - accessible by all authenticated users (admin and regular users)
			tasks.GET("", authMiddleware.RequireUser(), controller.GetAllTasks)       // GET /api/v1/tasks
			tasks.GET("/:id", authMiddleware.RequireUser(), controller.GetTaskByID)   // GET /api/v1/tasks/:id
			
			// Write operations - accessible only by admins
			tasks.POST("", authMiddleware.RequireAdmin(), controller.CreateTask)       // POST /api/v1/tasks (admin only)
			tasks.PUT("/:id", authMiddleware.RequireAdmin(), controller.UpdateTask)    // PUT /api/v1/tasks/:id (admin only)
			tasks.DELETE("/:id", authMiddleware.RequireAdmin(), controller.DeleteTask) // DELETE /api/v1/tasks/:id (admin only)
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