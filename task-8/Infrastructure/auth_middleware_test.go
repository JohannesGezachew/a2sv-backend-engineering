package Infrastructure

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"task_manager/Domain"
)

// MockJWTServiceForAuth is a mock implementation for auth middleware testing
type MockJWTServiceForAuth struct {
	mock.Mock
}

func (m *MockJWTServiceForAuth) GenerateToken(user *Domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTServiceForAuth) ValidateToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func (m *MockJWTServiceForAuth) GetJWTSecret() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func setupAuthTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthMiddleware_AuthenticateToken(t *testing.T) {
	t.Run("Success - valid token", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		// Create a valid token with claims
		claims := jwt.MapClaims{
			"user_id":  "507f1f77bcf86cd799439011",
			"username": "testuser",
			"role":     Domain.RoleUser,
		}
		token := &jwt.Token{
			Valid:  true,
			Claims: claims,
		}

		mockJWTService.On("ValidateToken", "valid.jwt.token").Return(token, nil)

		router.GET("/protected", authMiddleware.AuthenticateToken(), func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			username, _ := c.Get("username")
			role, _ := c.Get("role")
			
			c.JSON(http.StatusOK, gin.H{
				"user_id":  userID,
				"username": username,
				"role":     role,
			})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid.jwt.token")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "507f1f77bcf86cd799439011", response["user_id"])
		assert.Equal(t, "testuser", response["username"])
		assert.Equal(t, Domain.RoleUser, response["role"])
		
		mockJWTService.AssertExpectations(t)
	})

	t.Run("Error - missing authorization header", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		router.GET("/protected", authMiddleware.AuthenticateToken(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Authorization header required", response.Message)
		assert.Equal(t, "Missing Authorization header", response.Error)
	})

	t.Run("Error - invalid authorization header format", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		router.GET("/protected", authMiddleware.AuthenticateToken(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid authorization header format", response.Message)
		assert.Equal(t, "Authorization header must be in format: Bearer <token>", response.Error)
	})

	t.Run("Error - wrong bearer format", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		router.GET("/protected", authMiddleware.AuthenticateToken(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Basic token123")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid authorization header format", response.Message)
	})

	t.Run("Error - invalid token", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		mockJWTService.On("ValidateToken", "invalid.token").Return(nil, jwt.ErrSignatureInvalid)

		router.GET("/protected", authMiddleware.AuthenticateToken(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid or expired token", response.Message)
		
		mockJWTService.AssertExpectations(t)
	})

	t.Run("Error - token not valid", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		// Create an invalid token
		token := &jwt.Token{
			Valid: false,
		}

		mockJWTService.On("ValidateToken", "expired.token").Return(token, nil)

		router.GET("/protected", authMiddleware.AuthenticateToken(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer expired.token")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid or expired token", response.Message)
		
		mockJWTService.AssertExpectations(t)
	})

	t.Run("Error - invalid token claims", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		// Create a token with invalid claims type
		token := &jwt.Token{
			Valid:  true,
			Claims: jwt.RegisteredClaims{}, // Wrong claims type
		}

		mockJWTService.On("ValidateToken", "invalid.claims.token").Return(token, nil)

		router.GET("/protected", authMiddleware.AuthenticateToken(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.claims.token")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid token claims", response.Message)
		assert.Equal(t, "Could not parse token claims", response.Error)
		
		mockJWTService.AssertExpectations(t)
	})
}

func TestAuthMiddleware_RequireAdmin(t *testing.T) {
	t.Run("Success - admin user", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		router.Use(func(c *gin.Context) {
			c.Set("role", Domain.RoleAdmin)
			c.Next()
		})
		router.GET("/admin", authMiddleware.RequireAdmin(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		})

		req := httptest.NewRequest("GET", "/admin", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "admin access granted", response["message"])
	})

	t.Run("Error - role not found in context", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		router.GET("/admin", authMiddleware.RequireAdmin(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		})

		req := httptest.NewRequest("GET", "/admin", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "User role not found", response.Message)
		assert.Equal(t, "Authentication required", response.Error)
	})

	t.Run("Error - regular user trying to access admin endpoint", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		router.Use(func(c *gin.Context) {
			c.Set("role", Domain.RoleUser)
			c.Next()
		})
		router.GET("/admin", authMiddleware.RequireAdmin(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		})

		req := httptest.NewRequest("GET", "/admin", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusForbidden, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Access denied", response.Message)
		assert.Equal(t, "Admin privileges required", response.Error)
	})
}

func TestAuthMiddleware_RequireUser(t *testing.T) {
	t.Run("Success - admin user", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		router.Use(func(c *gin.Context) {
			c.Set("role", Domain.RoleAdmin)
			c.Next()
		})
		router.GET("/user", authMiddleware.RequireUser(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "user access granted"})
		})

		req := httptest.NewRequest("GET", "/user", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user access granted", response["message"])
	})

	t.Run("Success - regular user", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		router.Use(func(c *gin.Context) {
			c.Set("role", Domain.RoleUser)
			c.Next()
		})
		router.GET("/user", authMiddleware.RequireUser(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "user access granted"})
		})

		req := httptest.NewRequest("GET", "/user", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user access granted", response["message"])
	})

	t.Run("Error - role not found in context", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		router.GET("/user", authMiddleware.RequireUser(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "user access granted"})
		})

		req := httptest.NewRequest("GET", "/user", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "User role not found", response.Message)
		assert.Equal(t, "Authentication required", response.Error)
	})

	t.Run("Error - invalid role", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		router.Use(func(c *gin.Context) {
			c.Set("role", "invalid_role")
			c.Next()
		})
		router.GET("/user", authMiddleware.RequireUser(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "user access granted"})
		})

		req := httptest.NewRequest("GET", "/user", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusForbidden, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Access denied", response.Message)
		assert.Equal(t, "Valid user role required", response.Error)
	})
}

// Test constructor
func TestNewAuthMiddleware(t *testing.T) {
	mockJWTService := new(MockJWTServiceForAuth)
	authMiddleware := NewAuthMiddleware(mockJWTService)

	assert.NotNil(t, authMiddleware)
	assert.Equal(t, mockJWTService, authMiddleware.jwtService)
}

// Integration test with multiple middleware layers
func TestAuthMiddleware_Integration(t *testing.T) {
	t.Run("Success - full authentication and authorization flow", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		// Create a valid admin token
		claims := jwt.MapClaims{
			"user_id":  "507f1f77bcf86cd799439011",
			"username": "admin",
			"role":     Domain.RoleAdmin,
		}
		token := &jwt.Token{
			Valid:  true,
			Claims: claims,
		}

		mockJWTService.On("ValidateToken", "admin.jwt.token").Return(token, nil)

		// Setup route with both authentication and admin authorization
		router.GET("/admin/users", 
			authMiddleware.AuthenticateToken(),
			authMiddleware.RequireAdmin(),
			func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				username, _ := c.Get("username")
				role, _ := c.Get("role")
				
				c.JSON(http.StatusOK, gin.H{
					"message":  "admin endpoint accessed",
					"user_id":  userID,
					"username": username,
					"role":     role,
				})
			})

		req := httptest.NewRequest("GET", "/admin/users", nil)
		req.Header.Set("Authorization", "Bearer admin.jwt.token")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "admin endpoint accessed", response["message"])
		assert.Equal(t, "507f1f77bcf86cd799439011", response["user_id"])
		assert.Equal(t, "admin", response["username"])
		assert.Equal(t, Domain.RoleAdmin, response["role"])
		
		mockJWTService.AssertExpectations(t)
	})

	t.Run("Error - regular user trying to access admin endpoint", func(t *testing.T) {
		// Arrange
		mockJWTService := new(MockJWTServiceForAuth)
		authMiddleware := NewAuthMiddleware(mockJWTService)
		router := setupAuthTestRouter()

		// Create a valid user token (not admin)
		claims := jwt.MapClaims{
			"user_id":  "507f1f77bcf86cd799439012",
			"username": "user",
			"role":     Domain.RoleUser,
		}
		token := &jwt.Token{
			Valid:  true,
			Claims: claims,
		}

		mockJWTService.On("ValidateToken", "user.jwt.token").Return(token, nil)

		// Setup route with both authentication and admin authorization
		router.GET("/admin/users", 
			authMiddleware.AuthenticateToken(),
			authMiddleware.RequireAdmin(),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "admin endpoint accessed"})
			})

		req := httptest.NewRequest("GET", "/admin/users", nil)
		req.Header.Set("Authorization", "Bearer user.jwt.token")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusForbidden, w.Code)
		
		var response Domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Access denied", response.Message)
		assert.Equal(t, "Admin privileges required", response.Error)
		
		mockJWTService.AssertExpectations(t)
	})
}