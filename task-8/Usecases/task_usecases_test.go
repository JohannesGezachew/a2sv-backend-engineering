package Usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"task_manager/Domain"
)

// MockTaskRepository is a mock implementation of TaskRepositoryInterface
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetAll() ([]*Domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]*Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByID(id string) (*Domain.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Create(task *Domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Update(id string, task *Domain.Task) error {
	args := m.Called(id, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestTaskUsecase_GetAllTasks(t *testing.T) {
	t.Run("Success - return all tasks", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		expectedTasks := []*Domain.Task{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Task 1",
				Description: "Description 1",
				Status:      Domain.StatusPending,
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Task 2",
				Description: "Description 2",
				Status:      Domain.StatusCompleted,
			},
		}
		mockRepo.On("GetAll").Return(expectedTasks, nil)

		// Act
		tasks, err := taskUsecase.GetAllTasks()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedTasks, tasks)
		assert.Len(t, tasks, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - return empty list", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		expectedTasks := []*Domain.Task{}
		mockRepo.On("GetAll").Return(expectedTasks, nil)

		// Act
		tasks, err := taskUsecase.GetAllTasks()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedTasks, tasks)
		assert.Len(t, tasks, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - repository error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		expectedError := errors.New("database connection error")
		mockRepo.On("GetAll").Return([]*Domain.Task(nil), expectedError)

		// Act
		tasks, err := taskUsecase.GetAllTasks()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, tasks)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskUsecase_GetTaskByID(t *testing.T) {
	t.Run("Success - task found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		taskID := primitive.NewObjectID().Hex()
		expectedTask := &Domain.Task{
			ID:          primitive.NewObjectID(),
			Title:       "Test Task",
			Description: "Test Description",
			Status:      Domain.StatusInProgress,
		}
		mockRepo.On("GetByID", taskID).Return(expectedTask, nil)

		// Act
		task, err := taskUsecase.GetTaskByID(taskID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedTask, task)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - task not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		taskID := primitive.NewObjectID().Hex()
		expectedError := errors.New("task not found")
		mockRepo.On("GetByID", taskID).Return(nil, expectedError)

		// Act
		task, err := taskUsecase.GetTaskByID(taskID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - invalid ID format", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		invalidID := "invalid-id"
		expectedError := errors.New("invalid task ID format")
		mockRepo.On("GetByID", invalidID).Return(nil, expectedError)

		// Act
		task, err := taskUsecase.GetTaskByID(invalidID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskUsecase_CreateTask(t *testing.T) {
	t.Run("Success - create task with all fields", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		taskReq := Domain.TaskRequest{
			Title:       "New Task",
			Description: "New Description",
			DueDate:     "2024-12-31",
			Status:      Domain.StatusPending,
		}
		mockRepo.On("Create", mock.AnythingOfType("*Domain.Task")).Return(nil)

		// Act
		task, err := taskUsecase.CreateTask(taskReq)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, taskReq.Title, task.Title)
		assert.Equal(t, taskReq.Description, task.Description)
		assert.Equal(t, taskReq.Status, task.Status)
		
		expectedDate, _ := time.Parse("2006-01-02", taskReq.DueDate)
		assert.Equal(t, expectedDate, task.DueDate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - create task without due date", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		taskReq := Domain.TaskRequest{
			Title:       "Task without due date",
			Description: "Description",
			Status:      Domain.StatusInProgress,
		}
		mockRepo.On("Create", mock.AnythingOfType("*Domain.Task")).Return(nil)

		// Act
		task, err := taskUsecase.CreateTask(taskReq)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, taskReq.Title, task.Title)
		assert.True(t, task.DueDate.IsZero())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - invalid status", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		taskReq := Domain.TaskRequest{
			Title:  "Task with invalid status",
			Status: "invalid_status",
		}

		// Act
		task, err := taskUsecase.CreateTask(taskReq)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
		assert.Nil(t, task)
		// No repository call expected for validation errors
	})

	t.Run("Error - invalid due date format", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		taskReq := Domain.TaskRequest{
			Title:   "Task with invalid date",
			DueDate: "invalid-date",
			Status:  Domain.StatusPending,
		}

		// Act
		task, err := taskUsecase.CreateTask(taskReq)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid due date format")
		assert.Nil(t, task)
		// No repository call expected for validation errors
	})

	t.Run("Error - repository error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		taskReq := Domain.TaskRequest{
			Title:  "Task",
			Status: Domain.StatusPending,
		}
		expectedError := errors.New("database error")
		mockRepo.On("Create", mock.AnythingOfType("*Domain.Task")).Return(expectedError)

		// Act
		task, err := taskUsecase.CreateTask(taskReq)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskUsecase_UpdateTask(t *testing.T) {
	t.Run("Success - update existing task", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		taskID := primitive.NewObjectID().Hex()
		existingTask := &Domain.Task{
			ID:          primitive.NewObjectID(),
			Title:       "Old Title",
			Description: "Old Description",
			Status:      Domain.StatusPending,
		}
		updatedTask := &Domain.Task{
			ID:          existingTask.ID,
			Title:       "Updated Title",
			Description: "Updated Description",
			Status:      Domain.StatusCompleted,
			DueDate:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		}
		taskReq := Domain.TaskRequest{
			Title:       "Updated Title",
			Description: "Updated Description",
			DueDate:     "2024-12-31",
			Status:      Domain.StatusCompleted,
		}

		mockRepo.On("GetByID", taskID).Return(existingTask, nil).Once()
		mockRepo.On("Update", taskID, mock.AnythingOfType("*Domain.Task")).Return(nil).Once()
		mockRepo.On("GetByID", taskID).Return(updatedTask, nil).Once()

		// Act
		task, err := taskUsecase.UpdateTask(taskID, taskReq)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, updatedTask.Title, task.Title)
		assert.Equal(t, updatedTask.Status, task.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - task not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		taskID := primitive.NewObjectID().Hex()
		taskReq := Domain.TaskRequest{
			Title:  "Updated Title",
			Status: Domain.StatusCompleted,
		}
		expectedError := errors.New("task not found")
		mockRepo.On("GetByID", taskID).Return(nil, expectedError)

		// Act
		task, err := taskUsecase.UpdateTask(taskID, taskReq)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - invalid status", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		taskID := primitive.NewObjectID().Hex()
		existingTask := &Domain.Task{
			ID:     primitive.NewObjectID(),
			Title:  "Existing Task",
			Status: Domain.StatusPending,
		}
		taskReq := Domain.TaskRequest{
			Title:  "Updated Title",
			Status: "invalid_status",
		}
		mockRepo.On("GetByID", taskID).Return(existingTask, nil)

		// Act
		task, err := taskUsecase.UpdateTask(taskID, taskReq)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskUsecase_DeleteTask(t *testing.T) {
	t.Run("Success - delete existing task", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		taskID := primitive.NewObjectID().Hex()
		mockRepo.On("Delete", taskID).Return(nil)

		// Act
		err := taskUsecase.DeleteTask(taskID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - task not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepository)
		taskUsecase := NewTaskUsecase(mockRepo)
		
		taskID := primitive.NewObjectID().Hex()
		expectedError := errors.New("task not found")
		mockRepo.On("Delete", taskID).Return(expectedError)

		// Act
		err := taskUsecase.DeleteTask(taskID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

// Additional standalone tests
func TestNewTaskUsecase(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	
	assert.NotNil(t, usecase)
	assert.Implements(t, (*TaskUsecaseInterface)(nil), usecase)
}

func TestTaskUsecaseInterface(t *testing.T) {
	// Test that our implementation satisfies the interface
	mockRepo := new(MockTaskRepository)
	var _ TaskUsecaseInterface = &TaskUsecase{taskRepo: mockRepo}
	var _ TaskUsecaseInterface = NewTaskUsecase(mockRepo)
}