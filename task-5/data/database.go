package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URI        string
	Database   string
	Collection string
}

// GetDatabaseConfig returns database configuration from environment variables or defaults
func GetDatabaseConfig() *DatabaseConfig {
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

	return &DatabaseConfig{
		URI:        uri,
		Database:   database,
		Collection: collection,
	}
}

// ConnectToMongoDB establishes a connection to MongoDB
func ConnectToMongoDB(config *DatabaseConfig) (*mongo.Client, error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
	}

	log.Println("Disconnected from MongoDB")
	return nil
}