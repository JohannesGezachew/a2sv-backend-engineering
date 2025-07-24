package data

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"task_manager/models"
)

// UserService handles all user-related business logic with MongoDB
type UserService struct {
	collection *mongo.Collection
}

// NewUserService creates a new instance of UserService with MongoDB connection
func NewUserService(client *mongo.Client, dbName string) *UserService {
	collection := client.Database(dbName).Collection("users")
	return &UserService{
		collection: collection,
	}
}

// CreateUser creates a new user in MongoDB
func (us *UserService) CreateUser(userReq models.UserRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if username already exists
	existingUser, _ := us.GetUserByUsername(userReq.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Check if this is the first user (make them admin)
	userCount, err := us.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	role := models.RoleUser
	if userCount == 0 {
		role = models.RoleAdmin
	}

	user := &models.User{
		ID:        primitive.NewObjectID(),
		Username:  userReq.Username,
		Password:  string(hashedPassword),
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = us.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username from MongoDB
func (us *UserService) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := us.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID from MongoDB
func (us *UserService) GetUserByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var user models.User
	err = us.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// AuthenticateUser validates user credentials and returns the user if valid
func (us *UserService) AuthenticateUser(loginReq models.LoginRequest) (*models.User, error) {
	user, err := us.GetUserByUsername(loginReq.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// PromoteUserToAdmin promotes a user to admin role
func (us *UserService) PromoteUserToAdmin(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := us.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if user.Role == models.RoleAdmin {
		return nil, errors.New("user is already an admin")
	}

	update := bson.M{
		"$set": bson.M{
			"role":       models.RoleAdmin,
			"updated_at": time.Now(),
		},
	}

	_, err = us.collection.UpdateOne(ctx, bson.M{"username": username}, update)
	if err != nil {
		return nil, err
	}

	// Return updated user
	return us.GetUserByUsername(username)
}

// GetAllUsers returns all users from MongoDB (for admin purposes)
func (us *UserService) GetAllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := us.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}