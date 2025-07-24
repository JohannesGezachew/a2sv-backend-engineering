package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Get MongoDB URI from environment variable
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	fmt.Printf("Testing MongoDB connection to: %s\n", uri)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("✅ Successfully connected to MongoDB!")

	// List databases
	databases, err := client.ListDatabaseNames(ctx, map[string]interface{}{})
	if err != nil {
		log.Printf("Failed to list databases: %v", err)
	} else {
		fmt.Println("Available databases:")
		for _, db := range databases {
			fmt.Printf("  - %s\n", db)
		}
	}

	// Test database access
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "taskmanager"
	}

	database := client.Database(dbName)
	collections, err := database.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		log.Printf("Failed to list collections in database '%s': %v", dbName, err)
	} else {
		fmt.Printf("Collections in database '%s':\n", dbName)
		if len(collections) == 0 {
			fmt.Println("  (no collections found)")
		} else {
			for _, collection := range collections {
				fmt.Printf("  - %s\n", collection)
			}
		}
	}

	fmt.Println("✅ MongoDB connection test completed successfully!")
}