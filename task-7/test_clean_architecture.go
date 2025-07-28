package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"

	"task_manager/Delivery/routers"
	"task_manager/Domain"
)

// TestCleanArchitecture demonstrates the Clean Architecture implementation
func TestCleanArchitecture() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	// Set up test database configuration
	dbConfig := &routers.DatabaseConfig{
		URI:        getEnvOrDefault("MONGODB_URI", "mongodb://localhost:27017"),
		Database:   getEnvOrDefault("MONGODB_DATABASE", "taskmanager_test"),
		Collection: getEnvOrDefault("MONGODB_COLLECTION", "tasks_test"),
	}

	// Connect to MongoDB
	client, err := connectToMongoDB(dbConfig)
	if err != nil {
		log.Printf("Warning: Could not connect to MongoDB: %v", err)
		log.Println("Skipping database-dependent tests")
		return
	}
	defer disconnectFromMongoDB(client)

	// Initialize router with Clean Architecture
	router := routers.SetupRouter(client, dbConfig)

	// Test Clean Architecture layers
	fmt.Println("=== Testing Clean Architecture Implementation ===")
	
	// Test 1: User Registration (Domain -> Usecase -> Repository -> Infrastructure)
	fmt.Println("\n1. Testing User Registration Flow:")
	testUserRegistration(router)

	// Test 2: User Login (Authentication flow)
	fmt.Println("\n2. Testing User Login Flow:")
	token := testUserLogin(router)

	// Test 3: Task Creation (Admin-only operation)
	fmt.Println("\n3. Testing Task Creation Flow:")
	testTaskCreation(router, token)

	// Test 4: Task Retrieval (Read operation)
	fmt.Println("\n4. Testing Task Retrieval Flow:")
	testTaskRetrieval(router, token)

	fmt.Println("\n=== Clean Architecture Test Complete ===")
}

func testUserRegistration(router *gin.Engine) {
	userReq := Domain.UserRequest{
		Username: "testuser",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(userReq)
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Printf("   Status: %d\n", w.Code)
	if w.Code == http.StatusCreated || w.Code == http.StatusConflict {
		fmt.Println("   ✓ User registration flow working correctly")
	} else {
		fmt.Printf("   ✗ Unexpected status code: %d\n", w.Code)
		fmt.Printf("   Response: %s\n", w.Body.String())
	}
}

func testUserLogin(router *gin.Engine) string {
	loginReq := Domain.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Printf("   Status: %d\n", w.Code)
	if w.Code == http.StatusOK {
		var response Domain.LoginResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		fmt.Println("   ✓ User login flow working correctly")
		fmt.Printf("   Token received: %s...\n", response.Token[:20])
		return response.Token
	} else {
		fmt.Printf("   ✗ Login failed with status: %d\n", w.Code)
		fmt.Printf("   Response: %s\n", w.Body.String())
		return ""
	}
}

func testTaskCreation(router *gin.Engine, token string) {
	if token == "" {
		fmt.Println("   ✗ No token available, skipping task creation test")
		return
	}

	taskReq := Domain.TaskRequest{
		Title:       "Test Task",
		Description: "This is a test task created through Clean Architecture",
		DueDate:     "2024-12-31",
		Status:      "pending",
	}

	jsonData, _ := json.Marshal(taskReq)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Printf("   Status: %d\n", w.Code)
	if w.Code == http.StatusCreated {
		fmt.Println("   ✓ Task creation flow working correctly")
	} else {
		fmt.Printf("   ✗ Task creation failed with status: %d\n", w.Code)
		fmt.Printf("   Response: %s\n", w.Body.String())
	}
}

func testTaskRetrieval(router *gin.Engine, token string) {
	if token == "" {
		fmt.Println("   ✗ No token available, skipping task retrieval test")
		return
	}

	req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Printf("   Status: %d\n", w.Code)
	if w.Code == http.StatusOK {
		fmt.Println("   ✓ Task retrieval flow working correctly")
	} else {
		fmt.Printf("   ✗ Task retrieval failed with status: %d\n", w.Code)
		fmt.Printf("   Response: %s\n", w.Body.String())
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func connectToMongoDB(config *routers.DatabaseConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Printf("Successfully connected to MongoDB at %s", config.URI)
	return client, nil
}

func disconnectFromMongoDB(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
	}

	log.Println("Disconnected from MongoDB")
	return nil
}

func main() {
	TestCleanArchitecture()
}