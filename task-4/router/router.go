package router

import (
	"github.com/gin-gonic/gin"
	"task_manager/controllers"
	"task_manager/data"
)

// SetupRouter initializes and configures the Gin router
func SetupRouter() *gin.Engine {
	// Create Gin router with default middleware (logger and recovery)
	router := gin.Default()

	// Initialize services and controllers
	taskService := data.NewTaskService()
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