package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("🧪 Task-6 API Test Runner")
	fmt.Println("=" * 30)

	// Check if server is running
	fmt.Println("📡 Checking if server is running...")
	if !isServerRunning() {
		fmt.Println("❌ Server is not running on localhost:8080")
		fmt.Println("Please start the server first with: go run main.go")
		os.Exit(1)
	}
	fmt.Println("✅ Server is running")

	// Wait a moment for server to be fully ready
	time.Sleep(2 * time.Second)

	// Run the API tests
	fmt.Println("\n🚀 Running API tests...")
	cmd := exec.Command("go", "run", "test_api.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ Test execution failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n🎉 Test execution completed!")
}

func isServerRunning() bool {
	// Simple check - try to make a request to the server
	cmd := exec.Command("curl", "-s", "-o", "nul", "-w", "%{http_code}", "http://localhost:8080/health")
	output, err := cmd.Output()
	
	if err != nil {
		return false
	}
	
	return string(output) == "200" || string(output) == "404" // 404 is fine, means server is running
}