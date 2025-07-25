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
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

type Task struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
	User    User   `json:"user"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func main() {
	fmt.Println("🚀 Starting comprehensive Task Management API testing...")
	fmt.Println(strings.Repeat("=", 60))

	// Test 1: Health Check
	fmt.Println("\n🏥 Test 1: Health Check")
	testHealthCheck()

	// Test 2: User Registration
	fmt.Println("\n📝 Test 2: User Registration")
	testUserRegistration()

	// Test 3: User Login
	fmt.Println("\n🔐 Test 3: User Login")
	adminToken := testUserLogin()

	// Test 4: Protected Route Access
	fmt.Println("\n🛡️ Test 4: Protected Route Access")
	testProtectedRoutes(adminToken)

	// Test 5: Task CRUD Operations (Admin)
	fmt.Println("\n📋 Test 5: Task CRUD Operations (Admin)")
	testTaskOperations(adminToken)

	// Test 6: User Role Testing
	fmt.Println("\n👥 Test 6: User Role Testing")
	testUserRoles(adminToken)

	// Test 7: Authorization Tests
	fmt.Println("\n🔒 Test 7: Authorization Tests")
	testAuthorization()

	// Test 8: Invalid Token Tests
	fmt.Println("\n❌ Test 8: Invalid Token Tests")
	testInvalidTokens()

	fmt.Println("\n✅ All tests completed!")
	fmt.Println(strings.Repeat("=", 60))
}

func testHealthCheck() {
	fmt.Println("  Testing health endpoint...")
	
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		fmt.Printf("    ❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("    Status: %d, Response: %s\n", resp.StatusCode, string(body))
	
	if resp.StatusCode == 200 {
		fmt.Printf("    ✅ Health check successful\n")
	} else {
		fmt.Printf("    ❌ Health check failed\n")
	}
}

func testUserRegistration() {
	users := []User{
		{Username: "admin_user", Password: "admin123"},
		{Username: "regular_user", Password: "user123"},
		{Username: "test_user", Password: "test123"},
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
		} else if resp.StatusCode == 409 {
			fmt.Printf("    ⚠️ User %s already exists\n", user.Username)
		} else {
			fmt.Printf("    ❌ Registration failed for %s\n", user.Username)
		}
	}
}

func testUserLogin() string {
	fmt.Println("  Testing login with admin user...")
	
	loginReq := LoginRequest{
		Username: "admin_user",
		Password: "admin123",
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
		fmt.Printf("  User role: %s\n", loginResp.User.Role)
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

	// Create a task (Admin only)
	fmt.Println("  Creating a new task (Admin operation)...")
	task := Task{
		Title:       "Test Task from API Test",
		Description: "This is a comprehensive test task",
		Status:      "pending",
		DueDate:     "2025-08-01",
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
		var apiResp APIResponse
		json.Unmarshal(body, &apiResp)
		taskData, _ := json.Marshal(apiResp.Data)
		json.Unmarshal(taskData, &createdTask)
		fmt.Printf("    ✅ Task created successfully with ID: %s\n", createdTask.ID)
	}

	// Get all tasks (Available to all authenticated users)
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

	// Update task if we have an ID (Admin only)
	if createdTask.ID != "" {
		fmt.Println("  Updating the task (Admin operation)...")
		createdTask.Title = "Updated Test Task"
		createdTask.Status = "in_progress"
		
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

		// Get task by ID
		fmt.Println("  Getting task by ID...")
		req, _ = http.NewRequest("GET", baseURL+"/api/v1/tasks/"+createdTask.ID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		resp, err = client.Do(req)
		if err != nil {
			fmt.Printf("    ❌ Get task by ID error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		body, _ = io.ReadAll(resp.Body)
		fmt.Printf("    Get by ID Status: %d, Response: %s\n", resp.StatusCode, string(body))

		if resp.StatusCode == 200 {
			fmt.Printf("    ✅ Task retrieved by ID successfully\n")
		}

		// Delete task (Admin only)
		fmt.Println("  Deleting the task (Admin operation)...")
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

func testUserRoles(adminToken string) {
	if adminToken == "" {
		fmt.Println("  ⚠️ No admin token available, skipping user role tests")
		return
	}

	client := &http.Client{}

	// Test getting all users (Admin only)
	fmt.Println("  Testing get all users (Admin only)...")
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("    ❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("    Status: %d, Response: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode == 200 {
		fmt.Printf("    ✅ Admin successfully retrieved all users\n")
	} else {
		fmt.Printf("    ❌ Failed to retrieve users\n")
	}

	// Test promoting a user (Admin only)
	fmt.Println("  Testing user promotion (Admin only)...")
	promoteReq := map[string]string{
		"username": "regular_user",
	}
	
	jsonData, _ := json.Marshal(promoteReq)
	req, _ = http.NewRequest("POST", baseURL+"/api/v1/users/promote", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("    ❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("    Status: %d, Response: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode == 200 {
		fmt.Printf("    ✅ User promoted successfully\n")
	} else {
		fmt.Printf("    ⚠️ User promotion failed (may already be admin)\n")
	}
}

func testAuthorization() {
	// Test accessing protected routes without token
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

	// Test regular user trying to access admin endpoints
	fmt.Println("  Testing regular user access to admin endpoints...")
	
	// First login as regular user
	regularLoginReq := LoginRequest{
		Username: "test_user",
		Password: "test123",
	}

	jsonData, _ = json.Marshal(regularLoginReq)
	resp, err = http.Post(baseURL+"/api/v1/login", "application/json", bytes.NewBuffer(jsonData))
	
	if err != nil {
		fmt.Printf("    ❌ Error logging in regular user: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	
	if resp.StatusCode == 200 {
		var loginResp LoginResponse
		json.Unmarshal(body, &loginResp)
		
		// Try to create a task (admin only operation)
		client := &http.Client{}
		task := Task{
			Title:       "Unauthorized Task",
			Description: "This should fail",
			Status:      "pending",
		}
		
		jsonData, _ = json.Marshal(task)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/tasks", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		req.Header.Set("Content-Type", "application/json")
		
		resp, err = client.Do(req)
		if err != nil {
			fmt.Printf("    ❌ Error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		body, _ = io.ReadAll(resp.Body)
		fmt.Printf("    Create task as regular user - Status: %d, Response: %s\n", resp.StatusCode, string(body))

		if resp.StatusCode == 403 {
			fmt.Printf("    ✅ Regular user properly blocked from admin operations\n")
		} else {
			fmt.Printf("    ❌ Regular user should not be able to create tasks\n")
		}
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

	// Test with no Authorization header
	fmt.Println("  Testing with no Authorization header...")
	req, _ = http.NewRequest("GET", baseURL+"/api/v1/users/profile", nil)
	
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("    ❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("    Status: %d, Response: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode == 401 {
		fmt.Printf("    ✅ Missing authorization header properly rejected\n")
	} else {
		fmt.Printf("    ❌ Should have rejected missing authorization header\n")
	}
}