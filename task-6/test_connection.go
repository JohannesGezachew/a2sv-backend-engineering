package main

import (
	"log"

	"github.com/joho/godotenv"
	"task_manager/data"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get database configuration
	dbConfig := data.GetDatabaseConfig()
	
	log.Printf("Attempting to connect to MongoDB...")
	log.Printf("URI: %s", dbConfig.URI)
	log.Printf("Database: %s", dbConfig.Database)

	// Connect to MongoDB
	client, err := data.ConnectToMongoDB(dbConfig)
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}

	log.Println("✅ Successfully connected to MongoDB Atlas!")

	// Disconnect
	if err := data.DisconnectFromMongoDB(client); err != nil {
		log.Printf("Error disconnecting: %v", err)
	}

	log.Println("✅ Connection test completed successfully!")
}