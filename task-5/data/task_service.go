package data

import (
"context"
"errors"
"time"

"go.mongodb.org/mongo-driver/bson"
"go.mongodb.org/mongo-driver/bson/primitive"
"go.mongodb.org/mongo-driver/mongo"

"task_manager/models"
)

// TaskService handles all task-related business logic with MongoDB
type TaskService struct {
collection *mongo.Collection
}

// NewTaskService creates a new instance of TaskService with MongoDB connection
func NewTaskService(client *mongo.Client, dbName, collectionName string) *TaskService {
collection := client.Database(dbName).Collection(collectionName)
return &TaskService{
collection: collection,
}
}

// GetAllTasks returns all tasks from MongoDB
func (ts *TaskService) GetAllTasks() ([]*models.Task, error) {
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

cursor, err := ts.collection.Find(ctx, bson.M{})
if err != nil {
return nil, err
}
defer cursor.Close(ctx)

var tasks []*models.Task
if err = cursor.All(ctx, &tasks); err != nil {
return nil, err
}

return tasks, nil
}

// GetTaskByID returns a task by its ObjectID from MongoDB
func (ts *TaskService) GetTaskByID(id string) (*models.Task, error) {
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

objectID, err := primitive.ObjectIDFromHex(id)
if err != nil {
return nil, errors.New("invalid task ID format")
}

var task models.Task
err = ts.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
if err != nil {
if err == mongo.ErrNoDocuments {
return nil, errors.New("task not found")
}
return nil, err
}

return &task, nil
}

// CreateTask creates a new task in MongoDB
func (ts *TaskService) CreateTask(taskReq models.TaskRequest) (*models.Task, error) {
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Parse due date if provided
var dueDate time.Time
var err error
if taskReq.DueDate != "" {
dueDate, err = time.Parse("2006-01-02", taskReq.DueDate)
if err != nil {
return nil, errors.New("invalid due date format, use YYYY-MM-DD")
}
}

// Validate status
if !isValidStatus(taskReq.Status) {
return nil, errors.New("invalid status, must be one of: pending, in_progress, completed")
}

task := &models.Task{
ID:          primitive.NewObjectID(),
Title:       taskReq.Title,
Description: taskReq.Description,
DueDate:     dueDate,
Status:      taskReq.Status,
CreatedAt:   time.Now(),
UpdatedAt:   time.Now(),
}

_, err = ts.collection.InsertOne(ctx, task)
if err != nil {
return nil, err
}

return task, nil
}

// UpdateTask updates an existing task in MongoDB
func (ts *TaskService) UpdateTask(id string, taskReq models.TaskRequest) (*models.Task, error) {
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

objectID, err := primitive.ObjectIDFromHex(id)
if err != nil {
return nil, errors.New("invalid task ID format")
}

// Parse due date if provided
var dueDate time.Time
if taskReq.DueDate != "" {
dueDate, err = time.Parse("2006-01-02", taskReq.DueDate)
if err != nil {
return nil, errors.New("invalid due date format, use YYYY-MM-DD")
}
}

// Validate status
if !isValidStatus(taskReq.Status) {
return nil, errors.New("invalid status, must be one of: pending, in_progress, completed")
}

update := bson.M{
"$set": bson.M{
"title":       taskReq.Title,
"description": taskReq.Description,
"due_date":    dueDate,
"status":      taskReq.Status,
"updated_at":  time.Now(),
},
}

result, err := ts.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
if err != nil {
return nil, err
}

if result.MatchedCount == 0 {
return nil, errors.New("task not found")
}

// Return the updated task
return ts.GetTaskByID(id)
}

// DeleteTask deletes a task by its ObjectID from MongoDB
func (ts *TaskService) DeleteTask(id string) error {
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

objectID, err := primitive.ObjectIDFromHex(id)
if err != nil {
return errors.New("invalid task ID format")
}

result, err := ts.collection.DeleteOne(ctx, bson.M{"_id": objectID})
if err != nil {
return err
}

if result.DeletedCount == 0 {
return errors.New("task not found")
}

return nil
}

// isValidStatus checks if the provided status is valid
func isValidStatus(status string) bool {
validStatuses := []string{"pending", "in_progress", "completed"}
for _, validStatus := range validStatuses {
if status == validStatus {
return true
}
}
return false
}
