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

// TaskRepositoryInterface defines the contract for task data access
type TaskRepositoryInterface interface {
	GetAll() ([]*Domain.Task, error)
	GetByID(id string) (*Domain.Task, error)
	Create(task *Domain.Task) error
	Update(id string, task *Domain.Task) error
	Delete(id string) error
}

// TaskRepository implements TaskRepositoryInterface with MongoDB
type TaskRepository struct {
	collection *mongo.Collection
}

// NewTaskRepository creates a new instance of TaskRepository
func NewTaskRepository(client *mongo.Client, dbName, collectionName string) TaskRepositoryInterface {
	collection := client.Database(dbName).Collection(collectionName)
	return &TaskRepository{
		collection: collection,
	}
}

// GetAll returns all tasks from MongoDB
func (tr *TaskRepository) GetAll() ([]*Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := tr.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*Domain.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetByID returns a task by its ObjectID from MongoDB
func (tr *TaskRepository) GetByID(id string) (*Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid task ID format")
	}

	var task Domain.Task
	err = tr.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return &task, nil
}

// Create creates a new task in MongoDB
func (tr *TaskRepository) Create(task *Domain.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	_, err := tr.collection.InsertOne(ctx, task)
	return err
}

// Update updates an existing task in MongoDB
func (tr *TaskRepository) Update(id string, task *Domain.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID format")
	}

	task.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"due_date":    task.DueDate,
			"status":      task.Status,
			"updated_at":  task.UpdatedAt,
		},
	}

	result, err := tr.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

// Delete deletes a task by its ObjectID from MongoDB
func (tr *TaskRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID format")
	}

	result, err := tr.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}