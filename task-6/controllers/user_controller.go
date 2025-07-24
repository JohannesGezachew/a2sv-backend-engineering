package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"task_manager/data"
	"task_manager/middleware"
	"task_manager/models"
)

// UserController handles HTTP requests for user operations
type UserController struct {
	userService *data.UserService
}

// NewUserController creates a new instance of UserController
func NewUserController(userService *data.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Register handles POST /register
func (uc *UserController) Register(c *gin.Context) {
	var userReq models.UserRequest
	
	if err := c.ShouldBindJSON(&userReq); err != nil {
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user, err := uc.userService.CreateUser(userReq)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "username already exists" {
			statusCode = http.StatusConflict
		}
		
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		}
		c.JSON(statusCode, errorResponse)
		return
	}

	response := models.UserResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	}
	
	c.JSON(http.StatusCreated, response)
}

// Login handles POST /login
func (uc *UserController) Login(c *gin.Context) {
	var loginReq models.LoginRequest
	
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user, err := uc.userService.AuthenticateUser(loginReq)
	if err != nil {
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Authentication failed",
			Error:   err.Error(),
		}
		c.JSON(http.StatusUnauthorized, errorResponse)
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user)
	if err != nil {
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   err.Error(),
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	response := models.LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User:    user,
	}
	
	c.JSON(http.StatusOK, response)
}

// PromoteUser handles POST /promote (admin only)
func (uc *UserController) PromoteUser(c *gin.Context) {
	var promoteReq models.PromoteRequest
	
	if err := c.ShouldBindJSON(&promoteReq); err != nil {
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user, err := uc.userService.PromoteUserToAdmin(promoteReq.Username)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Failed to promote user",
			Error:   err.Error(),
		}
		c.JSON(statusCode, errorResponse)
		return
	}

	response := models.UserResponse{
		Success: true,
		Message: "User promoted to admin successfully",
		Data:    user,
	}
	
	c.JSON(http.StatusOK, response)
}

// GetAllUsers handles GET /users (admin only)
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve users",
			Error:   err.Error(),
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	
	response := models.UserResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
	}
	
	c.JSON(http.StatusOK, response)
}

// GetProfile handles GET /profile (authenticated users)
func (uc *UserController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "User ID not found in token",
			Error:   "Authentication required",
		}
		c.JSON(http.StatusUnauthorized, errorResponse)
		return
	}

	user, err := uc.userService.GetUserByID(userID.(string))
	if err != nil {
		errorResponse := models.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve user profile",
			Error:   err.Error(),
		}
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}

	response := models.UserResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    user,
	}
	
	c.JSON(http.StatusOK, response)
}