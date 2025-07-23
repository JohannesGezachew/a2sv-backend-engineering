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

func main() {
	fmt.Println("Testing Task Management API...")
	
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
		"title":       "Test Task",
		"description": "This is a test task",
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
	fmt.Println("\n4. Getting task by ID (1)...")
	resp, err = http.Get(baseURL + "/tasks/1")
	if err != nil {
		fmt.Printf("Error getting task by ID: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Get task by ID response: %s\n", string(body))
	
	// Test 5: Update task
	fmt.Println("\n5. Updating task...")
	updateData := map[string]interface{}{
		"title":       "Updated Test Task",
		"description": "This task has been updated",
		"due_date":    "2024-12-31",
		"status":      "in_progress",
	}
	
	jsonData, _ = json.Marshal(updateData)
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", baseURL+"/tasks/1", bytes.NewBuffer(jsonData))
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
	fmt.Println("\n6. Deleting task...")
	req, _ = http.NewRequest("DELETE", baseURL+"/tasks/1", nil)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error deleting task: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Delete task response: %s\n", string(body))
	
	fmt.Println("\nAPI testing completed!")
}