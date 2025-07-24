package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"task_manager/data"
	"task_manager/models"
)

// TaskController handles HTTP requests for task operations
type TaskController struct {
	taskService *data.TaskService
}

// NewTaskController creates a new instance of TaskController
func NewTaskController(taskService *data.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

// GetAllTasks handles GET /tasks
func (tc *TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := tc.taskService.GetAllTasks()
	if err != nil {
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve tasks",
			Error:   err.Error(),
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	
	response := models.TaskResponse{
		Success: true,
		Message: "Tasks retrieved successfully",
		Data:    tasks,
	}
	
	c.JSON(http.StatusOK, response)
}

// GetTaskByID handles GET /tasks/:id
func (tc *TaskController) GetTaskByID(c *gin.Context) {
	id := c.Param("id")

	task, err := tc.taskService.GetTaskByID(id)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() == "invalid task ID format" {
			statusCode = http.StatusBadRequest
		}
		
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Task not found",
			Error:   err.Error(),
		}
		c.JSON(statusCode, errorResponse)
		return
	}

	response := models.TaskResponse{
		Success: true,
		Message: "Task retrieved successfully",
		Data:    task,
	}
	
	c.JSON(http.StatusOK, response)
}

// CreateTask handles POST /tasks
func (tc *TaskController) CreateTask(c *gin.Context) {
	var taskReq models.TaskRequest
	
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	task, err := tc.taskService.CreateTask(taskReq)
	if err != nil {
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Failed to create task",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := models.TaskResponse{
		Success: true,
		Message: "Task created successfully",
		Data:    task,
	}
	
	c.JSON(http.StatusCreated, response)
}

// UpdateTask handles PUT /tasks/:id
func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var taskReq models.TaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	task, err := tc.taskService.UpdateTask(id, taskReq)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "task not found" {
			statusCode = http.StatusNotFound
		}
		if err.Error() == "invalid task ID format" {
			statusCode = http.StatusBadRequest
		}
		
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Failed to update task",
			Error:   err.Error(),
		}
		c.JSON(statusCode, errorResponse)
		return
	}

	response := models.TaskResponse{
		Success: true,
		Message: "Task updated successfully",
		Data:    task,
	}
	
	c.JSON(http.StatusOK, response)
}

// DeleteTask handles DELETE /tasks/:id
func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := tc.taskService.DeleteTask(id)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() == "invalid task ID format" {
			statusCode = http.StatusBadRequest
		}
		
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Failed to delete task",
			Error:   err.Error(),
		}
		c.JSON(statusCode, errorResponse)
		return
	}

	response := models.TaskResponse{
		Success: true,
		Message: "Task deleted successfully",
	}
	
	c.JSON(http.StatusOK, response)
}