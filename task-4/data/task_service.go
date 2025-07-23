package data

import (
	"errors"
	"sync"
	"time"

	"task_manager/models"
)

// TaskService handles all task-related business logic
type TaskService struct {
	tasks  map[int]*models.Task
	nextID int
	mutex  sync.RWMutex
}

// NewTaskService creates a new instance of TaskService
func NewTaskService() *TaskService {
	return &TaskService{
		tasks:  make(map[int]*models.Task),
		nextID: 1,
	}
}

// GetAllTasks returns all tasks
func (ts *TaskService) GetAllTasks() []*models.Task {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	tasks := make([]*models.Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetTaskByID returns a task by its ID
func (ts *TaskService) GetTaskByID(id int) (*models.Task, error) {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	task, exists := ts.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

// CreateTask creates a new task
func (ts *TaskService) CreateTask(taskReq models.TaskRequest) (*models.Task, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

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
		ID:          ts.nextID,
		Title:       taskReq.Title,
		Description: taskReq.Description,
		DueDate:     dueDate,
		Status:      taskReq.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	ts.tasks[ts.nextID] = task
	ts.nextID++

	return task, nil
}

// UpdateTask updates an existing task
func (ts *TaskService) UpdateTask(id int, taskReq models.TaskRequest) (*models.Task, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	task, exists := ts.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}

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

	// Update task fields
	task.Title = taskReq.Title
	task.Description = taskReq.Description
	task.DueDate = dueDate
	task.Status = taskReq.Status
	task.UpdatedAt = time.Now()

	return task, nil
}

// DeleteTask deletes a task by its ID
func (ts *TaskService) DeleteTask(id int) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	_, exists := ts.tasks[id]
	if !exists {
		return errors.New("task not found")
	}

	delete(ts.tasks, id)
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