package routers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	// Create a mock MongoDB client for testing
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	
	dbConfig := &DatabaseConfig{
		URI:        "mongodb://localhost:27017",
		Database:   "testdb",
		Collection: "tasks",
	}
	
	return SetupRouter(client, dbConfig)
}

func TestSetupRouter(t *testing.T) {
	t.Run("Success - router setup", func(t *testing.T) {
		// Arrange & Act
		router := setupTestRouter()

		// Assert
		assert.NotNil(t, router)
	})
}

func TestHealthEndpoint(t *testing.T) {
	t.Run("Success - health check", func(t *testing.T) {
		// Arrange
		router := setupTestRouter()
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "OK", response["status"])
		assert.Equal(t, "Task Management API is running", response["message"])
	})
}

func TestRouterEndpoints(t *testing.T) {
	router := setupTestRouter()

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		description    string
	}{
		{
			name:           "Register endpoint exists",
			method:         "POST",
			path:           "/api/v1/register",
			expectedStatus: http.StatusBadRequest, // No body provided, but endpoint exists
			description:    "Should reach the register handler",
		},
		{
			name:           "Login endpoint exists",
			method:         "POST",
			path:           "/api/v1/login",
			expectedStatus: http.StatusBadRequest, // No body provided, but endpoint exists
			description:    "Should reach the login handler",
		},
		{
			name:           "Profile endpoint requires auth",
			method:         "GET",
			path:           "/api/v1/users/profile",
			expectedStatus: http.StatusUnauthorized, // No auth header
			description:    "Should require authentication",
		},
		{
			name:           "Get all users requires auth",
			method:         "GET",
			path:           "/api/v1/users",
			expectedStatus: http.StatusUnauthorized, // No auth header
			description:    "Should require authentication",
		},
		{
			name:           "Promote user requires auth",
			method:         "POST",
			path:           "/api/v1/users/promote",
			expectedStatus: http.StatusUnauthorized, // No auth header
			description:    "Should require authentication",
		},
		{
			name:           "Get all tasks requires auth",
			method:         "GET",
			path:           "/api/v1/tasks",
			expectedStatus: http.StatusUnauthorized, // No auth header
			description:    "Should require authentication",
		},
		{
			name:           "Get task by ID requires auth",
			method:         "GET",
			path:           "/api/v1/tasks/507f1f77bcf86cd799439011",
			expectedStatus: http.StatusUnauthorized, // No auth header
			description:    "Should require authentication",
		},
		{
			name:           "Create task requires auth",
			method:         "POST",
			path:           "/api/v1/tasks",
			expectedStatus: http.StatusUnauthorized, // No auth header
			description:    "Should require authentication",
		},
		{
			name:           "Update task requires auth",
			method:         "PUT",
			path:           "/api/v1/tasks/507f1f77bcf86cd799439011",
			expectedStatus: http.StatusUnauthorized, // No auth header
			description:    "Should require authentication",
		},
		{
			name:           "Delete task requires auth",
			method:         "DELETE",
			path:           "/api/v1/tasks/507f1f77bcf86cd799439011",
			expectedStatus: http.StatusUnauthorized, // No auth header
			description:    "Should require authentication",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code, tc.description)
		})
	}
}

func TestRouterNotFoundEndpoint(t *testing.T) {
	t.Run("404 for non-existent endpoint", func(t *testing.T) {
		// Arrange
		router := setupTestRouter()
		req := httptest.NewRequest("GET", "/non-existent", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDatabaseConfig(t *testing.T) {
	t.Run("DatabaseConfig struct", func(t *testing.T) {
		// Arrange & Act
		config := &DatabaseConfig{
			URI:        "mongodb://localhost:27017",
			Database:   "testdb",
			Collection: "tasks",
		}

		// Assert
		assert.Equal(t, "mongodb://localhost:27017", config.URI)
		assert.Equal(t, "testdb", config.Database)
		assert.Equal(t, "tasks", config.Collection)
	})
}

func TestRouterMiddlewareChain(t *testing.T) {
	t.Run("Public endpoints don't require auth", func(t *testing.T) {
		router := setupTestRouter()

		publicEndpoints := []struct {
			method string
			path   string
		}{
			{"POST", "/api/v1/register"},
			{"POST", "/api/v1/login"},
			{"GET", "/health"},
		}

		for _, endpoint := range publicEndpoints {
			req := httptest.NewRequest(endpoint.method, endpoint.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// These endpoints should not return 401 (unauthorized)
			// They might return 400 (bad request) due to missing body, but not 401
			assert.NotEqual(t, http.StatusUnauthorized, w.Code, 
				"Endpoint %s %s should not require authentication", endpoint.method, endpoint.path)
		}
	})

	t.Run("Protected endpoints require auth", func(t *testing.T) {
		router := setupTestRouter()

		protectedEndpoints := []struct {
			method string
			path   string
		}{
			{"GET", "/api/v1/users/profile"},
			{"GET", "/api/v1/users"},
			{"POST", "/api/v1/users/promote"},
			{"GET", "/api/v1/tasks"},
			{"GET", "/api/v1/tasks/507f1f77bcf86cd799439011"},
			{"POST", "/api/v1/tasks"},
			{"PUT", "/api/v1/tasks/507f1f77bcf86cd799439011"},
			{"DELETE", "/api/v1/tasks/507f1f77bcf86cd799439011"},
		}

		for _, endpoint := range protectedEndpoints {
			req := httptest.NewRequest(endpoint.method, endpoint.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// These endpoints should return 401 (unauthorized) without auth header
			assert.Equal(t, http.StatusUnauthorized, w.Code, 
				"Endpoint %s %s should require authentication", endpoint.method, endpoint.path)
		}
	})
}

func TestRouterAPIVersioning(t *testing.T) {
	t.Run("API v1 endpoints are properly versioned", func(t *testing.T) {
		router := setupTestRouter()

		// Test that v1 endpoints exist with correct methods
		v1Endpoints := []struct {
			method string
			path   string
		}{
			{"POST", "/api/v1/register"},
			{"POST", "/api/v1/login"},
			{"GET", "/api/v1/users/profile"},
			{"GET", "/api/v1/users"},
			{"POST", "/api/v1/users/promote"},
			{"GET", "/api/v1/tasks"},
		}

		for _, endpoint := range v1Endpoints {
			req := httptest.NewRequest(endpoint.method, endpoint.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Should not return 404 (not found) - endpoint exists
			// May return 401 (unauthorized) or 400 (bad request) for protected/invalid endpoints
			assert.NotEqual(t, http.StatusNotFound, w.Code, 
				"Endpoint %s %s should exist", endpoint.method, endpoint.path)
		}
	})

	t.Run("Non-versioned endpoints return 404", func(t *testing.T) {
		router := setupTestRouter()

		nonVersionedEndpoints := []string{
			"/register",
			"/login",
			"/users",
			"/tasks",
		}

		for _, endpoint := range nonVersionedEndpoints {
			req := httptest.NewRequest("GET", endpoint, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Should return 404 (not found) - endpoint doesn't exist without versioning
			assert.Equal(t, http.StatusNotFound, w.Code, 
				"Endpoint %s should not exist without versioning", endpoint)
		}
	})
}

func TestRouterHTTPMethods(t *testing.T) {
	t.Run("Endpoints respond to correct HTTP methods", func(t *testing.T) {
		router := setupTestRouter()

		methodTests := []struct {
			correctMethod string
			wrongMethod   string
			path          string
		}{
			{"POST", "GET", "/api/v1/register"},
			{"POST", "PUT", "/api/v1/login"},
			{"GET", "POST", "/api/v1/users/profile"},
			{"POST", "GET", "/api/v1/users/promote"},
			{"GET", "POST", "/api/v1/tasks"},
			{"POST", "GET", "/api/v1/tasks"},
			{"PUT", "POST", "/api/v1/tasks/507f1f77bcf86cd799439011"},
			{"DELETE", "GET", "/api/v1/tasks/507f1f77bcf86cd799439011"},
		}

		for _, test := range methodTests {
			// Test wrong method returns 405, 404, or 401 (for protected endpoints)
			req := httptest.NewRequest(test.wrongMethod, test.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should return method not allowed, not found, or unauthorized (for protected endpoints)
			// Protected endpoints may return 401 because auth middleware runs before method check
			assert.True(t, w.Code == http.StatusMethodNotAllowed || w.Code == http.StatusNotFound || w.Code == http.StatusUnauthorized,
				"Wrong method %s for %s should return 405, 404, or 401, got %d", 
				test.wrongMethod, test.path, w.Code)
		}
	})
}

func TestRouterContentType(t *testing.T) {
	t.Run("Router handles JSON content type", func(t *testing.T) {
		router := setupTestRouter()

		// Test with JSON content type
		req := httptest.NewRequest("POST", "/api/v1/register", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should not return 415 (unsupported media type)
		assert.NotEqual(t, http.StatusUnsupportedMediaType, w.Code)
	})
}

func TestRouterCORS(t *testing.T) {
	t.Run("Router handles OPTIONS requests", func(t *testing.T) {
		router := setupTestRouter()

		req := httptest.NewRequest("OPTIONS", "/api/v1/register", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// OPTIONS requests should be handled (not return 405)
		// Gin handles OPTIONS by default
		assert.NotEqual(t, http.StatusMethodNotAllowed, w.Code)
	})
}

// Integration test for router setup with different database configs
func TestSetupRouterWithDifferentConfigs(t *testing.T) {
	testConfigs := []struct {
		name   string
		config *DatabaseConfig
	}{
		{
			name: "Default config",
			config: &DatabaseConfig{
				URI:        "mongodb://localhost:27017",
				Database:   "taskmanager",
				Collection: "tasks",
			},
		},
		{
			name: "Custom config",
			config: &DatabaseConfig{
				URI:        "mongodb://custom:27017",
				Database:   "customdb",
				Collection: "customtasks",
			},
		},
		{
			name: "Test config",
			config: &DatabaseConfig{
				URI:        "mongodb://test:27017",
				Database:   "testdb",
				Collection: "testtasks",
			},
		},
	}

	for _, tc := range testConfigs {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			client, _ := mongo.NewClient(options.Client().ApplyURI(tc.config.URI))

			// Act
			router := SetupRouter(client, tc.config)

			// Assert
			assert.NotNil(t, router)

			// Test that health endpoint works
			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}