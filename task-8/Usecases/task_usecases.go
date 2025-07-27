package Usecases

import (
	"errors"
	"time"

	"task_manager/Domain"
	"task_manager/Repositories"
)

// TaskUsecaseInterface defines the contract for task business logic
type TaskUsecaseInterface interface {
	GetAllTasks() ([]*Domain.Task, error)
	GetTaskByID(id string) (*Domain.Task, error)
	CreateTask(taskReq Domain.TaskRequest) (*Domain.Task, error)
	UpdateTask(id string, taskReq Domain.TaskRequest) (*Domain.Task, error)
	DeleteTask(id string) error
}

// TaskUsecase implements task business logic
type TaskUsecase struct {
	taskRepo Repositories.TaskRepositoryInterface
}

// NewTaskUsecase creates a new instance of TaskUsecase
func NewTaskUsecase(taskRepo Repositories.TaskRepositoryInterface) TaskUsecaseInterface {
	return &TaskUsecase{
		taskRepo: taskRepo,
	}
}

// GetAllTasks returns all tasks
func (tu *TaskUsecase) GetAllTasks() ([]*Domain.Task, error) {
	return tu.taskRepo.GetAll()
}

// GetTaskByID returns a task by its ID
func (tu *TaskUsecase) GetTaskByID(id string) (*Domain.Task, error) {
	return tu.taskRepo.GetByID(id)
}

// CreateTask creates a new task
func (tu *TaskUsecase) CreateTask(taskReq Domain.TaskRequest) (*Domain.Task, error) {
	// Validate status
	if !Domain.IsValidStatus(taskReq.Status) {
		return nil, errors.New("invalid status, must be one of: pending, in_progress, completed")
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

	task := &Domain.Task{
		Title:       taskReq.Title,
		Description: taskReq.Description,
		DueDate:     dueDate,
		Status:      taskReq.Status,
	}

	err = tu.taskRepo.Create(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// UpdateTask updates an existing task
func (tu *TaskUsecase) UpdateTask(id string, taskReq Domain.TaskRequest) (*Domain.Task, error) {
	// Check if task exists
	existingTask, err := tu.taskRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate status
	if !Domain.IsValidStatus(taskReq.Status) {
		return nil, errors.New("invalid status, must be one of: pending, in_progress, completed")
	}

	// Parse due date if provided
	var dueDate time.Time
	if taskReq.DueDate != "" {
		dueDate, err = time.Parse("2006-01-02", taskReq.DueDate)
		if err != nil {
			return nil, errors.New("invalid due date format, use YYYY-MM-DD")
		}
	}

	// Update task fields
	existingTask.Title = taskReq.Title
	existingTask.Description = taskReq.Description
	existingTask.DueDate = dueDate
	existingTask.Status = taskReq.Status

	err = tu.taskRepo.Update(id, existingTask)
	if err != nil {
		return nil, err
	}

	// Return updated task
	return tu.taskRepo.GetByID(id)
}

// DeleteTask deletes a task by its ID
func (tu *TaskUsecase) DeleteTask(id string) error {
	return tu.taskRepo.Delete(id)
}