package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"task_manager/Domain"
)

// Mock implementations for testing
type MockTaskUsecase struct {
	mock.Mock
}

func (m *MockTaskUsecase) GetAllTasks() ([]*Domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]*Domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) GetTaskByID(id string) (*Domain.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) CreateTask(taskReq Domain.TaskRequest) (*Domain.Task, error) {
	args := m.Called(taskReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) UpdateTask(id string, taskReq Domain.TaskRequest) (*Domain.Task, error) {
	args := m.Called(id, taskReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) DeleteTask(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) RegisterUser(userReq Domain.UserRequest) (*Domain.User, error) {
	args := m.Called(userReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

func (m *MockUserUsecase) LoginUser(loginReq Domain.LoginRequest) (*Domain.User, string, error) {
	args := m.Called(loginReq)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}
	return args.Get(0).(*Domain.User), args.String(1), args.Error(2)
}

func (m *MockUserUsecase) GetUserProfile(userID string) (*Domain.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

func (m *MockUserUsecase) GetAllUsers() ([]*Domain.User, error) {
	args := m.Called()
	return args.Get(0).([]*Domain.User), args.Error(1)
}

func (m *MockUserUsecase) PromoteUserToAdmin(username string) (*Domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

// Test setup helper
func setupTestController() (*Controller, *MockTaskUsecase, *MockUserUsecase) {
	mockTaskUsecase := new(MockTaskUsecase)
	mockUserUsecase := new(MockUserUsecase)
	controller := NewController(mockTaskUsecase, mockUserUsecase)
	return controller, mockTaskUsecase, mockUserUsecase
}

func setupGinContext() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// User Controller Tests

func TestController_Register(t *testing.T) {
	t.Run("Success - register user", func(t *testing.T) {
		// Arrange
		controller, _, mockUserUsecase := setupTestController()
		router := setupGinContext()
		router.POST("/register", controller.Register)

		userReq := Domain.UserRequest{
			Username: "testuser",
			Password: "password123",
		}
		expectedUser := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: "testuser",
			Role:     Domain.RoleUser,
		}

		mockUserUsecase.On("RegisterUser", userReq).Return(expectedUser, nil)

		reqBody, _ := json.Marshal(userReq)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response Domain.UserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "User registered successfully", response.Message)
		
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("Error - invalid JSON", func(t *testing.T) {
		// Arrange
		controller, _, _ := setupTestController()
		router := setupGinContext()
		router.POST("/register", controller.Register)

		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request payload", response.Message)
	})

	t.Run("Error - username already exists", func(t *testing.T) {
		// Arrange
		controller, _, mockUserUsecase := setupTestController()
		router := setupGinContext()
		router.POST("/register", controller.Register)

		userReq := Domain.UserRequest{
			Username: "existinguser",
			Password: "password123",
		}

		mockUserUsecase.On("RegisterUser", userReq).Return(nil, errors.New("username already exists"))

		reqBody, _ := json.Marshal(userReq)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusConflict, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Failed to create user", response.Message)
		
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("Error - other registration error", func(t *testing.T) {
		// Arrange
		controller, _, mockUserUsecase := setupTestController()
		router := setupGinContext()
		router.POST("/register", controller.Register)

		userReq := Domain.UserRequest{
			Username: "testuser",
			Password: "password123",
		}

		mockUserUsecase.On("RegisterUser", userReq).Return(nil, errors.New("database error"))

		reqBody, _ := json.Marshal(userReq)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Failed to create user", response.Message)
		
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestController_Login(t *testing.T) {
	t.Run("Success - login user", func(t *testing.T) {
		// Arrange
		controller, _, mockUserUsecase := setupTestController()
		router := setupGinContext()
		router.POST("/login", controller.Login)

		loginReq := Domain.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		expectedUser := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: "testuser",
			Role:     Domain.RoleUser,
		}
		expectedToken := "jwt.token.here"

		mockUserUsecase.On("LoginUser", loginReq).Return(expectedUser, expectedToken, nil)

		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response Domain.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Login successful", response.Message)
		assert.Equal(t, expectedToken, response.Token)
		
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("Error - invalid JSON", func(t *testing.T) {
		// Arrange
		controller, _, _ := setupTestController()
		router := setupGinContext()
		router.POST("/login", controller.Login)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request payload", response.Message)
	})

	t.Run("Error - authentication failed", func(t *testing.T) {
		// Arrange
		controller, _, mockUserUsecase := setupTestController()
		router := setupGinContext()
		router.POST("/login", controller.Login)

		loginReq := Domain.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}

		mockUserUsecase.On("LoginUser", loginReq).Return(nil, "", errors.New("invalid credentials"))

		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Authentication failed", response.Message)
		
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestController_PromoteUser(t *testing.T) {
	t.Run("Success - promote user", func(t *testing.T) {
		// Arrange
		controller, _, mockUserUsecase := setupTestController()
		router := setupGinContext()
		router.POST("/promote", controller.PromoteUser)

		promoteReq := Domain.PromoteRequest{
			Username: "usertoPromote",
		}
		expectedUser := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: "usertoPromote",
			Role:     Domain.RoleAdmin,
		}

		mockUserUsecase.On("PromoteUserToAdmin", promoteReq.Username).Return(expectedUser, nil)

		reqBody, _ := json.Marshal(promoteReq)
		req := httptest.NewRequest("POST", "/promote", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response Domain.UserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "User promoted to admin successfully", response.Message)
		
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("Error - user not found", func(t *testing.T) {
		// Arrange
		controller, _, mockUserUsecase := setupTestController()
		router := setupGinContext()
		router.POST("/promote", controller.PromoteUser)

		promoteReq := Domain.PromoteRequest{
			Username: "nonexistentuser",
		}

		mockUserUsecase.On("PromoteUserToAdmin", promoteReq.Username).Return(nil, errors.New("user not found"))

		reqBody, _ := json.Marshal(promoteReq)
		req := httptest.NewRequest("POST", "/promote", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Failed to promote user", response.Message)
		
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("Error - invalid JSON", func(t *testing.T) {
		// Arrange
		controller, _, _ := setupTestController()
		router := setupGinContext()
		router.POST("/promote", controller.PromoteUser)

		req := httptest.NewRequest("POST", "/promote", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request payload", response.Message)
	})
}

func TestController_GetAllUsers(t *testing.T) {
	t.Run("Success - get all users", func(t *testing.T) {
		// Arrange
		controller, _, mockUserUsecase := setupTestController()
		router := setupGinContext()
		router.GET("/users", controller.GetAllUsers)

		expectedUsers := []*Domain.User{
			{
				ID:       primitive.NewObjectID(),
				Username: "user1",
				Role:     Domain.RoleUser,
			},
			{
				ID:       primitive.NewObjectID(),
				Username: "admin1",
				Role:     Domain.RoleAdmin,
			},
		}

		mockUserUsecase.On("GetAllUsers").Return(expectedUsers, nil)

		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response Domain.UserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Users retrieved successfully", response.Message)
		
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("Error - repository error", func(t *testing.T) {
		// Arrange
		controller, _, mockUserUsecase := setupTestController()
		router := setupGinContext()
		router.GET("/users", controller.GetAllUsers)

		mockUserUsecase.On("GetAllUsers").Return([]*Domain.User(nil), errors.New("database error"))

		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Failed to retrieve users", response.Message)
		
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestController_GetProfile(t *testing.T) {
	t.Run("Success - get user profile", func(t *testing.T) {
		// Arrange
		controller, _, mockUserUsecase := setupTestController()
		router := setupGinContext()
		
		// Middleware to set user_id in context
		router.Use(func(c *gin.Context) {
			c.Set("user_id", "507f1f77bcf86cd799439011")
			c.Next()
		})
		router.GET("/profile", controller.GetProfile)

		expectedUser := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: "testuser",
			Role:     Domain.RoleUser,
		}

		mockUserUsecase.On("GetUserProfile", "507f1f77bcf86cd799439011").Return(expectedUser, nil)

		req := httptest.NewRequest("GET", "/profile", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response Domain.UserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Profile retrieved successfully", response.Message)
		
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("Error - user ID not found in context", func(t *testing.T) {
		// Arrange
		controller, _, _ := setupTestController()
		router := setupGinContext()
		router.GET("/profile", controller.GetProfile)

		req := httptest.NewRequest("GET", "/profile", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "User ID not found in token", response.Message)
	})

	t.Run("Error - user not found", func(t *testing.T) {
		// Arrange
		controller, _, mockUserUsecase := setupTestController()
		router := setupGinContext()
		
		// Middleware to set user_id in context
		router.Use(func(c *gin.Context) {
			c.Set("user_id", "507f1f77bcf86cd799439011")
			c.Next()
		})
		router.GET("/profile", controller.GetProfile)

		mockUserUsecase.On("GetUserProfile", "507f1f77bcf86cd799439011").Return(nil, errors.New("user not found"))

		req := httptest.NewRequest("GET", "/profile", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Failed to retrieve user profile", response.Message)
		
		mockUserUsecase.AssertExpectations(t)
	})
}

// Task Controller Tests

func TestController_GetAllTasks(t *testing.T) {
	t.Run("Success - get all tasks", func(t *testing.T) {
		// Arrange
		controller, mockTaskUsecase, _ := setupTestController()
		router := setupGinContext()
		router.GET("/tasks", controller.GetAllTasks)

		expectedTasks := []*Domain.Task{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Task 1",
				Description: "Description 1",
				Status:      Domain.StatusPending,
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Task 2",
				Description: "Description 2",
				Status:      Domain.StatusCompleted,
			},
		}

		mockTaskUsecase.On("GetAllTasks").Return(expectedTasks, nil)

		req := httptest.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response Domain.TaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Tasks retrieved successfully", response.Message)
		
		mockTaskUsecase.AssertExpectations(t)
	})

	t.Run("Error - repository error", func(t *testing.T) {
		// Arrange
		controller, mockTaskUsecase, _ := setupTestController()
		router := setupGinContext()
		router.GET("/tasks", controller.GetAllTasks)

		mockTaskUsecase.On("GetAllTasks").Return([]*Domain.Task(nil), errors.New("database error"))

		req := httptest.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Failed to retrieve tasks", response.Message)
		
		mockTaskUsecase.AssertExpectations(t)
	})
}

func TestController_GetTaskByID(t *testing.T) {
	t.Run("Success - get task by ID", func(t *testing.T) {
		// Arrange
		controller, mockTaskUsecase, _ := setupTestController()
		router := setupGinContext()
		router.GET("/tasks/:id", controller.GetTaskByID)

		taskID := primitive.NewObjectID().Hex()
		expectedTask := &Domain.Task{
			ID:          primitive.NewObjectID(),
			Title:       "Test Task",
			Description: "Test Description",
			Status:      Domain.StatusInProgress,
		}

		mockTaskUsecase.On("GetTaskByID", taskID).Return(expectedTask, nil)

		req := httptest.NewRequest("GET", "/tasks/"+taskID, nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response Domain.TaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Task retrieved successfully", response.Message)
		
		mockTaskUsecase.AssertExpectations(t)
	})

	t.Run("Error - task not found", func(t *testing.T) {
		// Arrange
		controller, mockTaskUsecase, _ := setupTestController()
		router := setupGinContext()
		router.GET("/tasks/:id", controller.GetTaskByID)

		taskID := primitive.NewObjectID().Hex()

		mockTaskUsecase.On("GetTaskByID", taskID).Return(nil, errors.New("task not found"))

		req := httptest.NewRequest("GET", "/tasks/"+taskID, nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Task not found", response.Message)
		
		mockTaskUsecase.AssertExpectations(t)
	})

	t.Run("Error - invalid task ID format", func(t *testing.T) {
		// Arrange
		controller, mockTaskUsecase, _ := setupTestController()
		router := setupGinContext()
		router.GET("/tasks/:id", controller.GetTaskByID)

		invalidID := "invalid-id"

		mockTaskUsecase.On("GetTaskByID", invalidID).Return(nil, errors.New("invalid task ID format"))

		req := httptest.NewRequest("GET", "/tasks/"+invalidID, nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Task not found", response.Message)
		
		mockTaskUsecase.AssertExpectations(t)
	})
}

func TestController_CreateTask(t *testing.T) {
	t.Run("Success - create task", func(t *testing.T) {
		// Arrange
		controller, mockTaskUsecase, _ := setupTestController()
		router := setupGinContext()
		router.POST("/tasks", controller.CreateTask)

		taskReq := Domain.TaskRequest{
			Title:       "New Task",
			Description: "New Description",
			DueDate:     "2024-12-31",
			Status:      Domain.StatusPending,
		}
		expectedTask := &Domain.Task{
			ID:          primitive.NewObjectID(),
			Title:       "New Task",
			Description: "New Description",
			Status:      Domain.StatusPending,
		}

		mockTaskUsecase.On("CreateTask", taskReq).Return(expectedTask, nil)

		reqBody, _ := json.Marshal(taskReq)
		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response Domain.TaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Task created successfully", response.Message)
		
		mockTaskUsecase.AssertExpectations(t)
	})

	t.Run("Error - invalid JSON", func(t *testing.T) {
		// Arrange
		controller, _, _ := setupTestController()
		router := setupGinContext()
		router.POST("/tasks", controller.CreateTask)

		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request payload", response.Message)
	})

	t.Run("Error - create task failed", func(t *testing.T) {
		// Arrange
		controller, mockTaskUsecase, _ := setupTestController()
		router := setupGinContext()
		router.POST("/tasks", controller.CreateTask)

		taskReq := Domain.TaskRequest{
			Title:  "New Task",
			Status: Domain.StatusPending,
		}

		mockTaskUsecase.On("CreateTask", taskReq).Return(nil, errors.New("validation error"))

		reqBody, _ := json.Marshal(taskReq)
		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Failed to create task", response.Message)
		
		mockTaskUsecase.AssertExpectations(t)
	})
}

func TestController_UpdateTask(t *testing.T) {
	t.Run("Success - update task", func(t *testing.T) {
		// Arrange
		controller, mockTaskUsecase, _ := setupTestController()
		router := setupGinContext()
		router.PUT("/tasks/:id", controller.UpdateTask)

		taskID := primitive.NewObjectID().Hex()
		taskReq := Domain.TaskRequest{
			Title:       "Updated Task",
			Description: "Updated Description",
			Status:      Domain.StatusCompleted,
		}
		expectedTask := &Domain.Task{
			ID:          primitive.NewObjectID(),
			Title:       "Updated Task",
			Description: "Updated Description",
			Status:      Domain.StatusCompleted,
		}

		mockTaskUsecase.On("UpdateTask", taskID, taskReq).Return(expectedTask, nil)

		reqBody, _ := json.Marshal(taskReq)
		req := httptest.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response Domain.TaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Task updated successfully", response.Message)
		
		mockTaskUsecase.AssertExpectations(t)
	})

	t.Run("Error - task not found", func(t *testing.T) {
		// Arrange
		controller, mockTaskUsecase, _ := setupTestController()
		router := setupGinContext()
		router.PUT("/tasks/:id", controller.UpdateTask)

		taskID := primitive.NewObjectID().Hex()
		taskReq := Domain.TaskRequest{
			Title:  "Updated Task",
			Status: Domain.StatusCompleted,
		}

		mockTaskUsecase.On("UpdateTask", taskID, taskReq).Return(nil, errors.New("task not found"))

		reqBody, _ := json.Marshal(taskReq)
		req := httptest.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Failed to update task", response.Message)
		
		mockTaskUsecase.AssertExpectations(t)
	})

	t.Run("Error - invalid JSON", func(t *testing.T) {
		// Arrange
		controller, _, _ := setupTestController()
		router := setupGinContext()
		router.PUT("/tasks/:id", controller.UpdateTask)

		taskID := primitive.NewObjectID().Hex()

		req := httptest.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request payload", response.Message)
	})
}

func TestController_DeleteTask(t *testing.T) {
	t.Run("Success - delete task", func(t *testing.T) {
		// Arrange
		controller, mockTaskUsecase, _ := setupTestController()
		router := setupGinContext()
		router.DELETE("/tasks/:id", controller.DeleteTask)

		taskID := primitive.NewObjectID().Hex()

		mockTaskUsecase.On("DeleteTask", taskID).Return(nil)

		req := httptest.NewRequest("DELETE", "/tasks/"+taskID, nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response Domain.TaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Task deleted successfully", response.Message)
		
		mockTaskUsecase.AssertExpectations(t)
	})

	t.Run("Error - task not found", func(t *testing.T) {
		// Arrange
		controller, mockTaskUsecase, _ := setupTestController()
		router := setupGinContext()
		router.DELETE("/tasks/:id", controller.DeleteTask)

		taskID := primitive.NewObjectID().Hex()

		mockTaskUsecase.On("DeleteTask", taskID).Return(errors.New("task not found"))

		req := httptest.NewRequest("DELETE", "/tasks/"+taskID, nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Failed to delete task", response.Message)
		
		mockTaskUsecase.AssertExpectations(t)
	})

	t.Run("Error - invalid task ID format", func(t *testing.T) {
		// Arrange
		controller, mockTaskUsecase, _ := setupTestController()
		router := setupGinContext()
		router.DELETE("/tasks/:id", controller.DeleteTask)

		invalidID := "invalid-id"

		mockTaskUsecase.On("DeleteTask", invalidID).Return(errors.New("invalid task ID format"))

		req := httptest.NewRequest("DELETE", "/tasks/"+invalidID, nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Failed to delete task", response.Message)
		
		mockTaskUsecase.AssertExpectations(t)
	})
}

// Test constructor
func TestNewController(t *testing.T) {
	mockTaskUsecase := new(MockTaskUsecase)
	mockUserUsecase := new(MockUserUsecase)

	controller := NewController(mockTaskUsecase, mockUserUsecase)

	assert.NotNil(t, controller)
	assert.Equal(t, mockTaskUsecase, controller.taskUsecase)
	assert.Equal(t, mockUserUsecase, controller.userUsecase)
}