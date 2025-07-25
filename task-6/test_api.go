package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const baseURL = "http://localhost:8080"

type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

type Task struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func main() {
	fmt.Println("🚀 Starting comprehensive API testing...")
	fmt.Println(strings.Repeat("=", 50))

	// Test 1: User Registration
	fmt.Println("\n📝 Test 1: User Registration")
	testUserRegistration()

	// Test 2: User Login
	fmt.Println("\n🔐 Test 2: User Login")
	token := testUserLogin()

	// Test 3: Protected Route Access
	fmt.Println("\n🛡️ Test 3: Protected Route Access")
	testProtectedRoutes(token)

	// Test 4: Task CRUD Operations
	fmt.Println("\n📋 Test 4: Task CRUD Operations")
	testTaskOperations(token)

	// Test 5: Authorization Tests
	fmt.Println("\n🔒 Test 5: Authorization Tests")
	testAuthorization()

	// Test 6: Invalid Token Tests
	fmt.Println("\n❌ Test 6: Invalid Token Tests")
	testInvalidTokens()

	fmt.Println("\n✅ All tests completed!")
}

func testUserRegistration() {
	users := []User{
		{Username: "testuser1", Email: "test1@example.com", Password: "password123"},
		{Username: "testuser2", Email: "test2@example.com", Password: "password456"},
		{Username: "admin", Email: "admin@example.com", Password: "admin123"},
	}

	for i, user := range users {
		fmt.Printf("  Registering user %d: %s\n", i+1, user.Username)
		
		jsonData, _ := json.Marshal(user)
		resp, err := http.Post(baseURL+"/api/v1/register", "application/json", bytes.NewBuffer(jsonData))
		
		if err != nil {
			fmt.Printf("    ❌ Error: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("    Status: %d, Response: %s\n", resp.StatusCode, string(body))
		
		if resp.StatusCode == 201 {
			fmt.Printf("    ✅ User %s registered successfully\n", user.Username)
		} else {
			fmt.Printf("    ⚠️ Registration failed for %s\n", user.Username)
		}
	}
}

func testUserLogin() string {
	loginReq := LoginRequest{
		Username: "testuser1",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(loginReq)
	resp, err := http.Post(baseURL+"/api/v1/login", "application/json", bytes.NewBuffer(jsonData))
	
	if err != nil {
		fmt.Printf("  ❌ Login error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("  Status: %d, Response: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode == 200 {
		var loginResp LoginResponse
		json.Unmarshal(body, &loginResp)
		fmt.Printf("  ✅ Login successful, token received\n")
		return loginResp.Token
	}

	fmt.Printf("  ❌ Login failed\n")
	return ""
}

func testProtectedRoutes(token string) {
	if token == "" {
		fmt.Println("  ⚠️ No token available, skipping protected route tests")
		return
	}

	// Test accessing protected user profile
	fmt.Println("  Testing protected user profile access...")
	
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("    ❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("    Status: %d, Response: %s\n", resp.StatusCode, string(body))
	
	if resp.StatusCode == 200 {
		fmt.Printf("    ✅ Protected route access successful\n")
	} else {
		fmt.Printf("    ❌ Protected route access failed\n")
	}
}

func testTaskOperations(token string) {
	if token == "" {
		fmt.Println("  ⚠️ No token available, skipping task operations")
		return
	}

	client := &http.Client{}

	// Create a task
	fmt.Println("  Creating a new task...")
	task := Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
	}
	
	jsonData, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("    ❌ Create task error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("    Create Status: %d, Response: %s\n", resp.StatusCode, string(body))

	var createdTask Task
	if resp.StatusCode == 201 {
		json.Unmarshal(body, &createdTask)
		fmt.Printf("    ✅ Task created successfully with ID: %s\n", createdTask.ID)
	}

	// Get all tasks
	fmt.Println("  Fetching all tasks...")
	req, _ = http.NewRequest("GET", baseURL+"/api/v1/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("    ❌ Get tasks error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("    Get Status: %d, Response: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode == 200 {
		fmt.Printf("    ✅ Tasks retrieved successfully\n")
	}

	// Update task if we have an ID
	if createdTask.ID != "" {
		fmt.Println("  Updating the task...")
		createdTask.Title = "Updated Test Task"
		createdTask.Status = "completed"
		
		jsonData, _ = json.Marshal(createdTask)
		req, _ = http.NewRequest("PUT", baseURL+"/api/v1/tasks/"+createdTask.ID, bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		
		resp, err = client.Do(req)
		if err != nil {
			fmt.Printf("    ❌ Update task error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		body, _ = io.ReadAll(resp.Body)
		fmt.Printf("    Update Status: %d, Response: %s\n", resp.StatusCode, string(body))

		if resp.StatusCode == 200 {
			fmt.Printf("    ✅ Task updated successfully\n")
		}

		// Delete task
		fmt.Println("  Deleting the task...")
		req, _ = http.NewRequest("DELETE", baseURL+"/api/v1/tasks/"+createdTask.ID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		resp, err = client.Do(req)
		if err != nil {
			fmt.Printf("    ❌ Delete task error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		body, _ = io.ReadAll(resp.Body)
		fmt.Printf("    Delete Status: %d, Response: %s\n", resp.StatusCode, string(body))

		if resp.StatusCode == 200 {
			fmt.Printf("    ✅ Task deleted successfully\n")
		}
	}
}

func testAuthorization() {
	// Test accessing routes without token
	fmt.Println("  Testing access without authentication token...")
	
	resp, err := http.Get(baseURL + "/api/v1/tasks")
	if err != nil {
		fmt.Printf("    ❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("    Status: %d, Response: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode == 401 {
		fmt.Printf("    ✅ Unauthorized access properly blocked\n")
	} else {
		fmt.Printf("    ❌ Should have been unauthorized\n")
	}

	// Test with wrong credentials
	fmt.Println("  Testing login with wrong credentials...")
	loginReq := LoginRequest{
		Username: "wronguser",
		Password: "wrongpass",
	}

	jsonData, _ := json.Marshal(loginReq)
	resp, err = http.Post(baseURL+"/api/v1/login", "application/json", bytes.NewBuffer(jsonData))
	
	if err != nil {
		fmt.Printf("    ❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("    Status: %d, Response: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode == 401 {
		fmt.Printf("    ✅ Invalid credentials properly rejected\n")
	} else {
		fmt.Printf("    ❌ Should have rejected invalid credentials\n")
	}
}

func testInvalidTokens() {
	client := &http.Client{}
	
	// Test with invalid token
	fmt.Println("  Testing with invalid token...")
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/tasks", nil)
	req.Header.Set("Authorization", "Bearer invalid_token_here")
	
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("    ❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("    Status: %d, Response: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode == 401 {
		fmt.Printf("    ✅ Invalid token properly rejected\n")
	} else {
		fmt.Printf("    ❌ Should have rejected invalid token\n")
	}

	// Test with malformed Authorization header
	fmt.Println("  Testing with malformed Authorization header...")
	req, _ = http.NewRequest("GET", baseURL+"/api/v1/tasks", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("    ❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("    Status: %d, Response: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode == 401 {
		fmt.Printf("    ✅ Malformed header properly rejected\n")
	} else {
		fmt.Printf("    ❌ Should have rejected malformed header\n")
	}
}