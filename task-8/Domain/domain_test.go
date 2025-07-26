package Domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "Valid status - pending",
			status:   StatusPending,
			expected: true,
		},
		{
			name:     "Valid status - in_progress",
			status:   StatusInProgress,
			expected: true,
		},
		{
			name:     "Valid status - completed",
			status:   StatusCompleted,
			expected: true,
		},
		{
			name:     "Invalid status - empty string",
			status:   "",
			expected: false,
		},
		{
			name:     "Invalid status - random string",
			status:   "invalid_status",
			expected: false,
		},
		{
			name:     "Invalid status - case sensitive",
			status:   "PENDING",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidStatus(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTaskStruct(t *testing.T) {
	t.Run("Task creation with all fields", func(t *testing.T) {
		id := primitive.NewObjectID()
		now := time.Now()
		dueDate := time.Now().Add(24 * time.Hour)

		task := Task{
			ID:          id,
			Title:       "Test Task",
			Description: "Test Description",
			DueDate:     dueDate,
			Status:      StatusPending,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		assert.Equal(t, id, task.ID)
		assert.Equal(t, "Test Task", task.Title)
		assert.Equal(t, "Test Description", task.Description)
		assert.Equal(t, dueDate, task.DueDate)
		assert.Equal(t, StatusPending, task.Status)
		assert.Equal(t, now, task.CreatedAt)
		assert.Equal(t, now, task.UpdatedAt)
	})

	t.Run("Task with empty fields", func(t *testing.T) {
		task := Task{}

		assert.True(t, task.ID.IsZero())
		assert.Empty(t, task.Title)
		assert.Empty(t, task.Description)
		assert.True(t, task.DueDate.IsZero())
		assert.Empty(t, task.Status)
		assert.True(t, task.CreatedAt.IsZero())
		assert.True(t, task.UpdatedAt.IsZero())
	})
}

func TestUserStruct(t *testing.T) {
	t.Run("User creation with all fields", func(t *testing.T) {
		id := primitive.NewObjectID()
		now := time.Now()

		user := User{
			ID:        id,
			Username:  "testuser",
			Password:  "hashedpassword",
			Role:      RoleUser,
			CreatedAt: now,
			UpdatedAt: now,
		}

		assert.Equal(t, id, user.ID)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "hashedpassword", user.Password)
		assert.Equal(t, RoleUser, user.Role)
		assert.Equal(t, now, user.CreatedAt)
		assert.Equal(t, now, user.UpdatedAt)
	})

	t.Run("User with admin role", func(t *testing.T) {
		user := User{
			Username: "admin",
			Role:     RoleAdmin,
		}

		assert.Equal(t, "admin", user.Username)
		assert.Equal(t, RoleAdmin, user.Role)
	})
}

func TestTaskRequest(t *testing.T) {
	t.Run("Valid task request", func(t *testing.T) {
		taskReq := TaskRequest{
			Title:       "New Task",
			Description: "Task description",
			DueDate:     "2024-12-31",
			Status:      StatusPending,
		}

		assert.Equal(t, "New Task", taskReq.Title)
		assert.Equal(t, "Task description", taskReq.Description)
		assert.Equal(t, "2024-12-31", taskReq.DueDate)
		assert.Equal(t, StatusPending, taskReq.Status)
	})

	t.Run("Task request with minimal fields", func(t *testing.T) {
		taskReq := TaskRequest{
			Title:  "Minimal Task",
			Status: StatusCompleted,
		}

		assert.Equal(t, "Minimal Task", taskReq.Title)
		assert.Empty(t, taskReq.Description)
		assert.Empty(t, taskReq.DueDate)
		assert.Equal(t, StatusCompleted, taskReq.Status)
	})
}

func TestUserRequest(t *testing.T) {
	t.Run("Valid user request", func(t *testing.T) {
		userReq := UserRequest{
			Username: "newuser",
			Password: "password123",
		}

		assert.Equal(t, "newuser", userReq.Username)
		assert.Equal(t, "password123", userReq.Password)
	})
}

func TestLoginRequest(t *testing.T) {
	t.Run("Valid login request", func(t *testing.T) {
		loginReq := LoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		assert.Equal(t, "testuser", loginReq.Username)
		assert.Equal(t, "password123", loginReq.Password)
	})
}

func TestPromoteRequest(t *testing.T) {
	t.Run("Valid promote request", func(t *testing.T) {
		promoteReq := PromoteRequest{
			Username: "usertoPromote",
		}

		assert.Equal(t, "usertoPromote", promoteReq.Username)
	})
}

func TestResponseStructs(t *testing.T) {
	t.Run("TaskResponse", func(t *testing.T) {
		response := TaskResponse{
			Success: true,
			Message: "Task created successfully",
			Data:    "test data",
		}

		assert.True(t, response.Success)
		assert.Equal(t, "Task created successfully", response.Message)
		assert.Equal(t, "test data", response.Data)
	})

	t.Run("UserResponse", func(t *testing.T) {
		response := UserResponse{
			Success: false,
			Message: "User creation failed",
		}

		assert.False(t, response.Success)
		assert.Equal(t, "User creation failed", response.Message)
		assert.Nil(t, response.Data)
	})

	t.Run("LoginResponse", func(t *testing.T) {
		user := &User{Username: "testuser", Role: RoleUser}
		response := LoginResponse{
			Success: true,
			Message: "Login successful",
			Token:   "jwt.token.here",
			User:    user,
		}

		assert.True(t, response.Success)
		assert.Equal(t, "Login successful", response.Message)
		assert.Equal(t, "jwt.token.here", response.Token)
		assert.Equal(t, user, response.User)
	})

	t.Run("ErrorResponse", func(t *testing.T) {
		response := ErrorResponse{
			Success: false,
			Message: "An error occurred",
			Error:   "detailed error message",
		}

		assert.False(t, response.Success)
		assert.Equal(t, "An error occurred", response.Message)
		assert.Equal(t, "detailed error message", response.Error)
	})
}

func TestJWTClaims(t *testing.T) {
	t.Run("JWT claims creation", func(t *testing.T) {
		claims := JWTClaims{
			UserID:   "507f1f77bcf86cd799439011",
			Username: "testuser",
			Role:     RoleAdmin,
		}

		assert.Equal(t, "507f1f77bcf86cd799439011", claims.UserID)
		assert.Equal(t, "testuser", claims.Username)
		assert.Equal(t, RoleAdmin, claims.Role)
	})
}

func TestConstants(t *testing.T) {
	t.Run("Role constants", func(t *testing.T) {
		assert.Equal(t, "admin", RoleAdmin)
		assert.Equal(t, "user", RoleUser)
	})

	t.Run("Status constants", func(t *testing.T) {
		assert.Equal(t, "pending", StatusPending)
		assert.Equal(t, "in_progress", StatusInProgress)
		assert.Equal(t, "completed", StatusCompleted)
	})
}