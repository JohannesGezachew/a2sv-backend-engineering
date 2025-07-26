package Repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"task_manager/Domain"
)

// UserRepositoryInterface defines the contract for user data access
type UserRepositoryInterface interface {
	GetAll() ([]*Domain.User, error)
	GetByID(id string) (*Domain.User, error)
	GetByUsername(username string) (*Domain.User, error)
	Create(user *Domain.User) error
	Update(id string, user *Domain.User) error
	UpdateByUsername(username string, user *Domain.User) error
	CountUsers() (int64, error)
}

// UserRepository implements UserRepositoryInterface with MongoDB
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(client *mongo.Client, dbName string) UserRepositoryInterface {
	collection := client.Database(dbName).Collection("users")
	return &UserRepository{
		collection: collection,
	}
}

// GetAll returns all users from MongoDB
func (ur *UserRepository) GetAll() ([]*Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := ur.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*Domain.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetByID retrieves a user by ID from MongoDB
func (ur *UserRepository) GetByID(id string) (*Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var user Domain.User
	err = ur.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByUsername retrieves a user by username from MongoDB
func (ur *UserRepository) GetByUsername(username string) (*Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user Domain.User
	err := ur.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Create creates a new user in MongoDB
func (ur *UserRepository) Create(user *Domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := ur.collection.InsertOne(ctx, user)
	return err
}

// Update updates an existing user in MongoDB
func (ur *UserRepository) Update(id string, user *Domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	user.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"password":   user.Password,
			"role":       user.Role,
			"updated_at": user.UpdatedAt,
		},
	}

	result, err := ur.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// UpdateByUsername updates an existing user by username in MongoDB
func (ur *UserRepository) UpdateByUsername(username string, user *Domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"role":       user.Role,
			"updated_at": user.UpdatedAt,
		},
	}

	result, err := ur.collection.UpdateOne(ctx, bson.M{"username": username}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// CountUsers returns the total number of users in the database
func (ur *UserRepository) CountUsers() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := ur.collection.CountDocuments(ctx, bson.M{})
	return count, err
}