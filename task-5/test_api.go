package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

type TaskResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func main() {
	fmt.Println("Testing Task Management API with MongoDB...")
	
	// Test 1: Health check
	fmt.Println("\n1. Testing health check...")
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Health check response: %s\n", string(body))
	
	// Test 2: Create a task
	fmt.Println("\n2. Creating a new task...")
	taskData := map[string]interface{}{
		"title":       "MongoDB Test Task",
		"description": "This is a test task with MongoDB integration",
		"due_date":    "2024-12-31",
		"status":      "pending",
	}
	
	jsonData, _ := json.Marshal(taskData)
	resp, err = http.Post(baseURL+"/tasks", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating task: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Create task response: %s\n", string(body))
	
	// Parse the response to get the task ID
	var createResponse TaskResponse
	json.Unmarshal(body, &createResponse)
	
	var taskID string
	if createResponse.Success {
		if taskData, ok := createResponse.Data.(map[string]interface{}); ok {
			if id, exists := taskData["id"]; exists {
				taskID = id.(string)
				fmt.Printf("Created task with ID: %s\n", taskID)
			}
		}
	}
	
	if taskID == "" {
		fmt.Println("Failed to get task ID from create response")
		return
	}
	
	// Test 3: Get all tasks
	fmt.Println("\n3. Getting all tasks...")
	resp, err = http.Get(baseURL + "/tasks")
	if err != nil {
		fmt.Printf("Error getting tasks: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Get all tasks response: %s\n", string(body))
	
	// Test 4: Get task by ID
	fmt.Printf("\n4. Getting task by ID (%s)...\n", taskID)
	resp, err = http.Get(baseURL + "/tasks/" + taskID)
	if err != nil {
		fmt.Printf("Error getting task by ID: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Get task by ID response: %s\n", string(body))
	
	// Test 5: Update task
	fmt.Printf("\n5. Updating task (%s)...\n", taskID)
	updateData := map[string]interface{}{
		"title":       "Updated MongoDB Test Task",
		"description": "This task has been updated with MongoDB",
		"due_date":    "2024-12-31",
		"status":      "in_progress",
	}
	
	jsonData, _ = json.Marshal(updateData)
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", baseURL+"/tasks/"+taskID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error updating task: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Update task response: %s\n", string(body))
	
	// Wait a moment before deletion
	time.Sleep(1 * time.Second)
	
	// Test 6: Delete task
	fmt.Printf("\n6. Deleting task (%s)...\n", taskID)
	req, _ = http.NewRequest("DELETE", baseURL+"/tasks/"+taskID, nil)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error deleting task: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Delete task response: %s\n", string(body))
	
	// Test 7: Try to get deleted task (should return 404)
	fmt.Printf("\n7. Trying to get deleted task (%s)...\n", taskID)
	resp, err = http.Get(baseURL + "/tasks/" + taskID)
	if err != nil {
		fmt.Printf("Error getting deleted task: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Get deleted task response: %s\n", string(body))
	
	// Test 8: Test invalid ObjectID format
	fmt.Println("\n8. Testing invalid ObjectID format...")
	resp, err = http.Get(baseURL + "/tasks/invalid-id")
	if err != nil {
		fmt.Printf("Error testing invalid ID: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Invalid ID response: %s\n", string(body))
	
	fmt.Println("\nMongoDB API testing completed!")
}