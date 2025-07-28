package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"task_manager/Delivery/routers"
)

// GetDatabaseConfig returns database configuration from environment variables or defaults
func GetDatabaseConfig() *routers.DatabaseConfig {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	database := os.Getenv("MONGODB_DATABASE")
	if database == "" {
		database = "taskmanager"
	}

	collection := os.Getenv("MONGODB_COLLECTION")
	if collection == "" {
		collection = "tasks"
	}

	return &routers.DatabaseConfig{
		URI:        uri,
		Database:   database,
		Collection: collection,
	}
}

// ConnectToMongoDB establishes a connection to MongoDB
func ConnectToMongoDB(config *routers.DatabaseConfig) (*mongo.Client, error) {
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

// DisconnectFromMongoDB closes the MongoDB connection
func DisconnectFromMongoDB(client *mongo.Client) error {
	if client == nil {
		return nil // Nothing to disconnect
	}

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
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		// Try loading from parent directory as fallback
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env file found, using system environment variables")
		} else {
			log.Println("Loaded .env from parent directory")
		}
	} else {
		log.Println("Loaded .env from current directory")
	}

	// Get database configuration
	dbConfig := GetDatabaseConfig()
	log.Printf("Using MongoDB URI: %s", dbConfig.URI)
	log.Printf("Using Database: %s", dbConfig.Database)

	// Connect to MongoDB
	client, err := ConnectToMongoDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Initialize the router with Clean Architecture
	r := routers.SetupRouter(client, dbConfig)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting Task Management API server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Disconnect from MongoDB
	if err := DisconnectFromMongoDB(client); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}

	log.Println("Server exited")
}