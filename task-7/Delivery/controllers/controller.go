package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"task_manager/Domain"
	"task_manager/Usecases"
)

// Controller handles HTTP requests for both task and user operations
type Controller struct {
	taskUsecase Usecases.TaskUsecaseInterface
	userUsecase Usecases.UserUsecaseInterface
}

// NewController creates a new instance of Controller
func NewController(taskUsecase Usecases.TaskUsecaseInterface, userUsecase Usecases.UserUsecaseInterface) *Controller {
	return &Controller{
		taskUsecase: taskUsecase,
		userUsecase: userUsecase,
	}
}

// User-related handlers

// Register handles POST /register
func (ctrl *Controller) Register(c *gin.Context) {
	var userReq Domain.UserRequest
	
	if err := c.ShouldBindJSON(&userReq); err != nil {
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user, err := ctrl.userUsecase.RegisterUser(userReq)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "username already exists" {
			statusCode = http.StatusConflict
		}
		
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		}
		c.JSON(statusCode, errorResponse)
		return
	}

	response := Domain.UserResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	}
	
	c.JSON(http.StatusCreated, response)
}

// Login handles POST /login
func (ctrl *Controller) Login(c *gin.Context) {
	var loginReq Domain.LoginRequest
	
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user, token, err := ctrl.userUsecase.LoginUser(loginReq)
	if err != nil {
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Authentication failed",
			Error:   err.Error(),
		}
		c.JSON(http.StatusUnauthorized, errorResponse)
		return
	}

	response := Domain.LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User:    user,
	}
	
	c.JSON(http.StatusOK, response)
}

// PromoteUser handles POST /promote (admin only)
func (ctrl *Controller) PromoteUser(c *gin.Context) {
	var promoteReq Domain.PromoteRequest
	
	if err := c.ShouldBindJSON(&promoteReq); err != nil {
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user, err := ctrl.userUsecase.PromoteUserToAdmin(promoteReq.Username)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Failed to promote user",
			Error:   err.Error(),
		}
		c.JSON(statusCode, errorResponse)
		return
	}

	response := Domain.UserResponse{
		Success: true,
		Message: "User promoted to admin successfully",
		Data:    user,
	}
	
	c.JSON(http.StatusOK, response)
}

// GetAllUsers handles GET /users (admin only)
func (ctrl *Controller) GetAllUsers(c *gin.Context) {
	users, err := ctrl.userUsecase.GetAllUsers()
	if err != nil {
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve users",
			Error:   err.Error(),
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	
	response := Domain.UserResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
	}
	
	c.JSON(http.StatusOK, response)
}

// GetProfile handles GET /profile (authenticated users)
func (ctrl *Controller) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "User ID not found in token",
			Error:   "Authentication required",
		}
		c.JSON(http.StatusUnauthorized, errorResponse)
		return
	}

	user, err := ctrl.userUsecase.GetUserProfile(userID.(string))
	if err != nil {
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve user profile",
			Error:   err.Error(),
		}
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}

	response := Domain.UserResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    user,
	}
	
	c.JSON(http.StatusOK, response)
}

// Task-related handlers

// GetAllTasks handles GET /tasks
func (ctrl *Controller) GetAllTasks(c *gin.Context) {
	tasks, err := ctrl.taskUsecase.GetAllTasks()
	if err != nil {
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve tasks",
			Error:   err.Error(),
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	
	response := Domain.TaskResponse{
		Success: true,
		Message: "Tasks retrieved successfully",
		Data:    tasks,
	}
	
	c.JSON(http.StatusOK, response)
}

// GetTaskByID handles GET /tasks/:id
func (ctrl *Controller) GetTaskByID(c *gin.Context) {
	id := c.Param("id")

	task, err := ctrl.taskUsecase.GetTaskByID(id)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() == "invalid task ID format" {
			statusCode = http.StatusBadRequest
		}
		
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Task not found",
			Error:   err.Error(),
		}
		c.JSON(statusCode, errorResponse)
		return
	}

	response := Domain.TaskResponse{
		Success: true,
		Message: "Task retrieved successfully",
		Data:    task,
	}
	
	c.JSON(http.StatusOK, response)
}

// CreateTask handles POST /tasks (admin only)
func (ctrl *Controller) CreateTask(c *gin.Context) {
	var taskReq Domain.TaskRequest
	
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	task, err := ctrl.taskUsecase.CreateTask(taskReq)
	if err != nil {
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Failed to create task",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := Domain.TaskResponse{
		Success: true,
		Message: "Task created successfully",
		Data:    task,
	}
	
	c.JSON(http.StatusCreated, response)
}

// UpdateTask handles PUT /tasks/:id (admin only)
func (ctrl *Controller) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var taskReq Domain.TaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	task, err := ctrl.taskUsecase.UpdateTask(id, taskReq)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "task not found" {
			statusCode = http.StatusNotFound
		}
		if err.Error() == "invalid task ID format" {
			statusCode = http.StatusBadRequest
		}
		
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Failed to update task",
			Error:   err.Error(),
		}
		c.JSON(statusCode, errorResponse)
		return
	}

	response := Domain.TaskResponse{
		Success: true,
		Message: "Task updated successfully",
		Data:    task,
	}
	
	c.JSON(http.StatusOK, response)
}

// DeleteTask handles DELETE /tasks/:id (admin only)
func (ctrl *Controller) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.taskUsecase.DeleteTask(id)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() == "invalid task ID format" {
			statusCode = http.StatusBadRequest
		}
		
		errorResponse := Domain.ErrorResponse{
			Success: false,
			Message: "Failed to delete task",
			Error:   err.Error(),
		}
		c.JSON(statusCode, errorResponse)
		return
	}

	response := Domain.TaskResponse{
		Success: true,
		Message: "Task deleted successfully",
	}
	
	c.JSON(http.StatusOK, response)
}