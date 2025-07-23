package main

import (
	"log"
	"task_manager/router"
)

func main() {
	// Initialize the router
	r := router.SetupRouter()

	// Start the server on port 8080
	log.Println("Starting Task Management API server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}