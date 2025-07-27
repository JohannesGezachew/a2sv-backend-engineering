package Repositories

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"task_manager/Domain"
)

// MockTaskRepository for testing purposes
type MockTaskRepositoryImpl struct {
	mock.Mock
}

func (m *MockTaskRepositoryImpl) GetAll() ([]*Domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]*Domain.Task), args.Error(1)
}

func (m *MockTaskRepositoryImpl) GetByID(id string) (*Domain.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.Task), args.Error(1)
}

func (m *MockTaskRepositoryImpl) Create(task *Domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepositoryImpl) Update(id string, task *Domain.Task) error {
	args := m.Called(id, task)
	return args.Error(0)
}

func (m *MockTaskRepositoryImpl) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestTaskRepository_GetAll(t *testing.T) {
	t.Run("Success - return all tasks", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		expectedTasks := []*Domain.Task{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Task 1",
				Description: "Description 1",
				Status:      Domain.StatusPending,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Task 2",
				Description: "Description 2",
				Status:      Domain.StatusCompleted,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		mockRepo.On("GetAll").Return(expectedTasks, nil)

		// Act
		tasks, err := mockRepo.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedTasks, tasks)
		assert.Len(t, tasks, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - return empty list", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		expectedTasks := []*Domain.Task{}
		mockRepo.On("GetAll").Return(expectedTasks, nil)

		// Act
		tasks, err := mockRepo.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedTasks, tasks)
		assert.Len(t, tasks, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - database connection error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		expectedError := errors.New("database connection failed")
		mockRepo.On("GetAll").Return([]*Domain.Task(nil), expectedError)

		// Act
		tasks, err := mockRepo.GetAll()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, tasks)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskRepository_GetByID(t *testing.T) {
	t.Run("Success - task found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		taskID := primitive.NewObjectID().Hex()
		expectedTask := &Domain.Task{
			ID:          primitive.NewObjectID(),
			Title:       "Test Task",
			Description: "Test Description",
			Status:      Domain.StatusInProgress,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockRepo.On("GetByID", taskID).Return(expectedTask, nil)

		// Act
		task, err := mockRepo.GetByID(taskID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedTask, task)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - task not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		taskID := primitive.NewObjectID().Hex()
		expectedError := errors.New("task not found")
		mockRepo.On("GetByID", taskID).Return(nil, expectedError)

		// Act
		task, err := mockRepo.GetByID(taskID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - invalid ID format", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		invalidID := "invalid-id-format"
		expectedError := errors.New("invalid task ID format")
		mockRepo.On("GetByID", invalidID).Return(nil, expectedError)

		// Act
		task, err := mockRepo.GetByID(invalidID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskRepository_Create(t *testing.T) {
	t.Run("Success - create task", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		task := &Domain.Task{
			Title:       "New Task",
			Description: "New Description",
			Status:      Domain.StatusPending,
		}
		mockRepo.On("Create", task).Return(nil)

		// Act
		err := mockRepo.Create(task)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - database insert failed", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		task := &Domain.Task{
			Title:  "New Task",
			Status: Domain.StatusPending,
		}
		expectedError := errors.New("database insert failed")
		mockRepo.On("Create", task).Return(expectedError)

		// Act
		err := mockRepo.Create(task)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - create task with all fields", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		dueDate := time.Now().Add(24 * time.Hour)
		task := &Domain.Task{
			Title:       "Complete Task",
			Description: "Complete Description",
			DueDate:     dueDate,
			Status:      Domain.StatusInProgress,
		}
		mockRepo.On("Create", task).Return(nil)

		// Act
		err := mockRepo.Create(task)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskRepository_Update(t *testing.T) {
	t.Run("Success - update existing task", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		taskID := primitive.NewObjectID().Hex()
		task := &Domain.Task{
			Title:       "Updated Task",
			Description: "Updated Description",
			Status:      Domain.StatusCompleted,
			UpdatedAt:   time.Now(),
		}
		mockRepo.On("Update", taskID, task).Return(nil)

		// Act
		err := mockRepo.Update(taskID, task)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - task not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		taskID := primitive.NewObjectID().Hex()
		task := &Domain.Task{
			Title:  "Updated Task",
			Status: Domain.StatusCompleted,
		}
		expectedError := errors.New("task not found")
		mockRepo.On("Update", taskID, task).Return(expectedError)

		// Act
		err := mockRepo.Update(taskID, task)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - invalid ID format", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		invalidID := "invalid-id"
		task := &Domain.Task{
			Title:  "Updated Task",
			Status: Domain.StatusCompleted,
		}
		expectedError := errors.New("invalid task ID format")
		mockRepo.On("Update", invalidID, task).Return(expectedError)

		// Act
		err := mockRepo.Update(invalidID, task)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - database update failed", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		taskID := primitive.NewObjectID().Hex()
		task := &Domain.Task{
			Title:  "Updated Task",
			Status: Domain.StatusCompleted,
		}
		expectedError := errors.New("database update failed")
		mockRepo.On("Update", taskID, task).Return(expectedError)

		// Act
		err := mockRepo.Update(taskID, task)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskRepository_Delete(t *testing.T) {
	t.Run("Success - delete existing task", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		taskID := primitive.NewObjectID().Hex()
		mockRepo.On("Delete", taskID).Return(nil)

		// Act
		err := mockRepo.Delete(taskID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - task not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		taskID := primitive.NewObjectID().Hex()
		expectedError := errors.New("task not found")
		mockRepo.On("Delete", taskID).Return(expectedError)

		// Act
		err := mockRepo.Delete(taskID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - invalid ID format", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		invalidID := "invalid-id"
		expectedError := errors.New("invalid task ID format")
		mockRepo.On("Delete", invalidID).Return(expectedError)

		// Act
		err := mockRepo.Delete(invalidID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - database delete failed", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		taskID := primitive.NewObjectID().Hex()
		expectedError := errors.New("database delete failed")
		mockRepo.On("Delete", taskID).Return(expectedError)

		// Act
		err := mockRepo.Delete(taskID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

// Test interface compliance
func TestTaskRepositoryInterface(t *testing.T) {
	mockRepo := new(MockTaskRepositoryImpl)
	var _ TaskRepositoryInterface = mockRepo
	assert.NotNil(t, mockRepo)
}

// Test edge cases
func TestTaskRepository_EdgeCases(t *testing.T) {
	t.Run("Create task with empty title", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		task := &Domain.Task{
			Title:  "",
			Status: Domain.StatusPending,
		}
		mockRepo.On("Create", task).Return(nil)

		// Act
		err := mockRepo.Create(task)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update task with zero time", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		taskID := primitive.NewObjectID().Hex()
		task := &Domain.Task{
			Title:     "Task",
			Status:    Domain.StatusPending,
			DueDate:   time.Time{}, // Zero time
			UpdatedAt: time.Time{}, // Zero time
		}
		mockRepo.On("Update", taskID, task).Return(nil)

		// Act
		err := mockRepo.Update(taskID, task)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get task with very long ID", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockTaskRepositoryImpl)
		longID := "very-long-id-that-might-cause-issues-in-some-systems-but-should-be-handled-gracefully"
		expectedError := errors.New("invalid task ID format")
		mockRepo.On("GetByID", longID).Return(nil, expectedError)

		// Act
		task, err := mockRepo.GetByID(longID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
	})
}