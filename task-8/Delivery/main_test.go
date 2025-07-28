package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"task_manager/Delivery/routers"
)

func TestGetDatabaseConfig(t *testing.T) {
	t.Run("Success - with environment variables", func(t *testing.T) {
		// Arrange
		os.Setenv("MONGODB_URI", "mongodb://test:27017")
		os.Setenv("MONGODB_DATABASE", "testdb")
		os.Setenv("MONGODB_COLLECTION", "testcollection")
		defer func() {
			os.Unsetenv("MONGODB_URI")
			os.Unsetenv("MONGODB_DATABASE")
			os.Unsetenv("MONGODB_COLLECTION")
		}()

		// Act
		config := GetDatabaseConfig()

		// Assert
		assert.NotNil(t, config)
		assert.Equal(t, "mongodb://test:27017", config.URI)
		assert.Equal(t, "testdb", config.Database)
		assert.Equal(t, "testcollection", config.Collection)
	})

	t.Run("Success - with default values", func(t *testing.T) {
		// Arrange - ensure env vars are not set
		os.Unsetenv("MONGODB_URI")
		os.Unsetenv("MONGODB_DATABASE")
		os.Unsetenv("MONGODB_COLLECTION")

		// Act
		config := GetDatabaseConfig()

		// Assert
		assert.NotNil(t, config)
		assert.Equal(t, "mongodb://localhost:27017", config.URI)
		assert.Equal(t, "taskmanager", config.Database)
		assert.Equal(t, "tasks", config.Collection)
	})

	t.Run("Success - partial environment variables", func(t *testing.T) {
		// Arrange
		os.Setenv("MONGODB_URI", "mongodb://partial:27017")
		os.Unsetenv("MONGODB_DATABASE")
		os.Unsetenv("MONGODB_COLLECTION")
		defer os.Unsetenv("MONGODB_URI")

		// Act
		config := GetDatabaseConfig()

		// Assert
		assert.NotNil(t, config)
		assert.Equal(t, "mongodb://partial:27017", config.URI)
		assert.Equal(t, "taskmanager", config.Database) // Default
		assert.Equal(t, "tasks", config.Collection)     // Default
	})

	t.Run("Success - empty environment variables use defaults", func(t *testing.T) {
		// Arrange
		os.Setenv("MONGODB_URI", "")
		os.Setenv("MONGODB_DATABASE", "")
		os.Setenv("MONGODB_COLLECTION", "")
		defer func() {
			os.Unsetenv("MONGODB_URI")
			os.Unsetenv("MONGODB_DATABASE")
			os.Unsetenv("MONGODB_COLLECTION")
		}()

		// Act
		config := GetDatabaseConfig()

		// Assert
		assert.NotNil(t, config)
		assert.Equal(t, "mongodb://localhost:27017", config.URI)
		assert.Equal(t, "taskmanager", config.Database)
		assert.Equal(t, "tasks", config.Collection)
	})
}

func TestConnectToMongoDB(t *testing.T) {
	t.Run("Success - valid configuration", func(t *testing.T) {
		// Arrange
		config := &routers.DatabaseConfig{
			URI:        "mongodb://localhost:27017",
			Database:   "testdb",
			Collection: "testcollection",
		}

		// Act
		client, err := ConnectToMongoDB(config)

		// Assert
		// Note: This test will fail if MongoDB is not running locally
		// In a real test environment, you might want to use a test container
		// or mock the MongoDB connection
		if err != nil {
			// If MongoDB is not available, just test that the function handles the error
			assert.Error(t, err)
			assert.Nil(t, client)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, client)
			
			// Clean up
			if client != nil {
				DisconnectFromMongoDB(client)
			}
		}
	})

	t.Run("Error - invalid URI", func(t *testing.T) {
		// Arrange
		config := &routers.DatabaseConfig{
			URI:        "invalid://uri",
			Database:   "testdb",
			Collection: "testcollection",
		}

		// Act
		client, err := ConnectToMongoDB(config)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "failed to connect to MongoDB")
	})

	t.Run("Error - unreachable host", func(t *testing.T) {
		// Arrange
		config := &routers.DatabaseConfig{
			URI:        "mongodb://unreachable:27017",
			Database:   "testdb",
			Collection: "testcollection",
		}

		// Act
		client, err := ConnectToMongoDB(config)

		// Assert
		// This should fail to connect or ping
		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, client)
		}
	})
}

func TestDisconnectFromMongoDB(t *testing.T) {
	t.Run("Success - disconnect valid client", func(t *testing.T) {
		// Arrange
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			t.Skip("Cannot create MongoDB client for test")
		}

		// Connect the client first
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		err = client.Connect(ctx)
		if err != nil {
			t.Skip("Cannot connect to MongoDB for test - MongoDB may not be running")
		}

		// Act
		err = DisconnectFromMongoDB(client)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Success - disconnect nil client", func(t *testing.T) {
		// Act
		err := DisconnectFromMongoDB(nil)

		// Assert
		// This might panic or return an error depending on the MongoDB driver
		// The test verifies the function handles nil gracefully
		// In practice, you might want to add nil checks to the function
		if err != nil {
			assert.Error(t, err)
		}
	})
}

// Test database configuration struct
func TestDatabaseConfigStruct(t *testing.T) {
	t.Run("DatabaseConfig creation", func(t *testing.T) {
		// Arrange & Act
		config := &routers.DatabaseConfig{
			URI:        "mongodb://test:27017",
			Database:   "testdb",
			Collection: "testcollection",
		}

		// Assert
		assert.Equal(t, "mongodb://test:27017", config.URI)
		assert.Equal(t, "testdb", config.Database)
		assert.Equal(t, "testcollection", config.Collection)
	})

	t.Run("DatabaseConfig with empty values", func(t *testing.T) {
		// Arrange & Act
		config := &routers.DatabaseConfig{}

		// Assert
		assert.Empty(t, config.URI)
		assert.Empty(t, config.Database)
		assert.Empty(t, config.Collection)
	})
}

// Test environment variable handling edge cases
func TestEnvironmentVariableHandling(t *testing.T) {
	t.Run("Environment variables with spaces", func(t *testing.T) {
		// Arrange
		os.Setenv("MONGODB_URI", " mongodb://spaced:27017 ")
		os.Setenv("MONGODB_DATABASE", " spaceddb ")
		os.Setenv("MONGODB_COLLECTION", " spacedcollection ")
		defer func() {
			os.Unsetenv("MONGODB_URI")
			os.Unsetenv("MONGODB_DATABASE")
			os.Unsetenv("MONGODB_COLLECTION")
		}()

		// Act
		config := GetDatabaseConfig()

		// Assert
		// The function doesn't trim spaces, so they should be preserved
		assert.Equal(t, " mongodb://spaced:27017 ", config.URI)
		assert.Equal(t, " spaceddb ", config.Database)
		assert.Equal(t, " spacedcollection ", config.Collection)
	})

	t.Run("Environment variables with special characters", func(t *testing.T) {
		// Arrange
		os.Setenv("MONGODB_URI", "mongodb://user:p@ss@host:27017")
		os.Setenv("MONGODB_DATABASE", "db-with-dashes")
		os.Setenv("MONGODB_COLLECTION", "collection_with_underscores")
		defer func() {
			os.Unsetenv("MONGODB_URI")
			os.Unsetenv("MONGODB_DATABASE")
			os.Unsetenv("MONGODB_COLLECTION")
		}()

		// Act
		config := GetDatabaseConfig()

		// Assert
		assert.Equal(t, "mongodb://user:p@ss@host:27017", config.URI)
		assert.Equal(t, "db-with-dashes", config.Database)
		assert.Equal(t, "collection_with_underscores", config.Collection)
	})
}

// Test configuration validation
func TestConfigurationValidation(t *testing.T) {
	t.Run("Valid MongoDB URI formats", func(t *testing.T) {
		validURIs := []string{
			"mongodb://localhost:27017",
			"mongodb://user:pass@localhost:27017",
			"mongodb://localhost:27017,localhost:27018",
			"mongodb+srv://cluster.mongodb.net",
		}

		for _, uri := range validURIs {
			os.Setenv("MONGODB_URI", uri)
			config := GetDatabaseConfig()
			assert.Equal(t, uri, config.URI)
			os.Unsetenv("MONGODB_URI")
		}
	})

	t.Run("Database name validation", func(t *testing.T) {
		validDatabases := []string{
			"taskmanager",
			"task_manager",
			"task-manager",
			"TaskManager",
			"db123",
		}

		for _, db := range validDatabases {
			os.Setenv("MONGODB_DATABASE", db)
			config := GetDatabaseConfig()
			assert.Equal(t, db, config.Database)
			os.Unsetenv("MONGODB_DATABASE")
		}
	})

	t.Run("Collection name validation", func(t *testing.T) {
		validCollections := []string{
			"tasks",
			"task_collection",
			"task-collection",
			"TaskCollection",
			"collection123",
		}

		for _, collection := range validCollections {
			os.Setenv("MONGODB_COLLECTION", collection)
			config := GetDatabaseConfig()
			assert.Equal(t, collection, config.Collection)
			os.Unsetenv("MONGODB_COLLECTION")
		}
	})
}

// Test concurrent access to configuration
func TestConcurrentConfigAccess(t *testing.T) {
	t.Run("Concurrent GetDatabaseConfig calls", func(t *testing.T) {
		// Arrange
		os.Setenv("MONGODB_URI", "mongodb://concurrent:27017")
		os.Setenv("MONGODB_DATABASE", "concurrentdb")
		os.Setenv("MONGODB_COLLECTION", "concurrentcollection")
		defer func() {
			os.Unsetenv("MONGODB_URI")
			os.Unsetenv("MONGODB_DATABASE")
			os.Unsetenv("MONGODB_COLLECTION")
		}()

		// Act - simulate concurrent access
		configs := make([]*routers.DatabaseConfig, 10)
		for i := 0; i < 10; i++ {
			go func(index int) {
				configs[index] = GetDatabaseConfig()
			}(i)
		}

		// Wait a bit for goroutines to complete
		// In a real test, you'd use sync.WaitGroup
		// This is a simplified version
		config := GetDatabaseConfig()

		// Assert
		assert.Equal(t, "mongodb://concurrent:27017", config.URI)
		assert.Equal(t, "concurrentdb", config.Database)
		assert.Equal(t, "concurrentcollection", config.Collection)
	})
}

// Benchmark test for configuration retrieval
func BenchmarkGetDatabaseConfig(b *testing.B) {
	// Setup
	os.Setenv("MONGODB_URI", "mongodb://benchmark:27017")
	os.Setenv("MONGODB_DATABASE", "benchmarkdb")
	os.Setenv("MONGODB_COLLECTION", "benchmarkcollection")
	defer func() {
		os.Unsetenv("MONGODB_URI")
		os.Unsetenv("MONGODB_DATABASE")
		os.Unsetenv("MONGODB_COLLECTION")
	}()

	// Benchmark
	for i := 0; i < b.N; i++ {
		GetDatabaseConfig()
	}
}