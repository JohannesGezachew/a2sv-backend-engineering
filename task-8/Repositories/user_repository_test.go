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

// MockUserRepository for testing purposes
type MockUserRepositoryImpl struct {
	mock.Mock
}

func (m *MockUserRepositoryImpl) GetAll() ([]*Domain.User, error) {
	args := m.Called()
	return args.Get(0).([]*Domain.User), args.Error(1)
}

func (m *MockUserRepositoryImpl) GetByID(id string) (*Domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

func (m *MockUserRepositoryImpl) GetByUsername(username string) (*Domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

func (m *MockUserRepositoryImpl) Create(user *Domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepositoryImpl) Update(id string, user *Domain.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserRepositoryImpl) UpdateByUsername(username string, user *Domain.User) error {
	args := m.Called(username, user)
	return args.Error(0)
}

func (m *MockUserRepositoryImpl) CountUsers() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func TestUserRepository_GetAll(t *testing.T) {
	t.Run("Success - return all users", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		expectedUsers := []*Domain.User{
			{
				ID:        primitive.NewObjectID(),
				Username:  "user1",
				Password:  "hashedpass1",
				Role:      Domain.RoleUser,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        primitive.NewObjectID(),
				Username:  "admin1",
				Password:  "hashedpass2",
				Role:      Domain.RoleAdmin,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRepo.On("GetAll").Return(expectedUsers, nil)

		// Act
		users, err := mockRepo.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Len(t, users, 2)
		assert.Equal(t, Domain.RoleUser, users[0].Role)
		assert.Equal(t, Domain.RoleAdmin, users[1].Role)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - return empty list", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		expectedUsers := []*Domain.User{}
		mockRepo.On("GetAll").Return(expectedUsers, nil)

		// Act
		users, err := mockRepo.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Len(t, users, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - database connection error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		expectedError := errors.New("database connection failed")
		mockRepo.On("GetAll").Return([]*Domain.User(nil), expectedError)

		// Act
		users, err := mockRepo.GetAll()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, users)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	t.Run("Success - user found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		userID := primitive.NewObjectID().Hex()
		expectedUser := &Domain.User{
			ID:        primitive.NewObjectID(),
			Username:  "testuser",
			Password:  "hashedpassword",
			Role:      Domain.RoleUser,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.On("GetByID", userID).Return(expectedUser, nil)

		// Act
		user, err := mockRepo.GetByID(userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, Domain.RoleUser, user.Role)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - user not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		userID := primitive.NewObjectID().Hex()
		expectedError := errors.New("user not found")
		mockRepo.On("GetByID", userID).Return(nil, expectedError)

		// Act
		user, err := mockRepo.GetByID(userID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - invalid ID format", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		invalidID := "invalid-id-format"
		expectedError := errors.New("invalid user ID format")
		mockRepo.On("GetByID", invalidID).Return(nil, expectedError)

		// Act
		user, err := mockRepo.GetByID(invalidID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserRepository_GetByUsername(t *testing.T) {
	t.Run("Success - user found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		username := "testuser"
		expectedUser := &Domain.User{
			ID:        primitive.NewObjectID(),
			Username:  username,
			Password:  "hashedpassword",
			Role:      Domain.RoleAdmin,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.On("GetByUsername", username).Return(expectedUser, nil)

		// Act
		user, err := mockRepo.GetByUsername(username)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, Domain.RoleAdmin, user.Role)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - user not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		username := "nonexistentuser"
		expectedError := errors.New("user not found")
		mockRepo.On("GetByUsername", username).Return(nil, expectedError)

		// Act
		user, err := mockRepo.GetByUsername(username)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - database error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		username := "testuser"
		expectedError := errors.New("database query failed")
		mockRepo.On("GetByUsername", username).Return(nil, expectedError)

		// Act
		user, err := mockRepo.GetByUsername(username)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - user with special characters in username", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		username := "user@domain.com"
		expectedUser := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: username,
			Role:     Domain.RoleUser,
		}
		mockRepo.On("GetByUsername", username).Return(expectedUser, nil)

		// Act
		user, err := mockRepo.GetByUsername(username)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.Equal(t, username, user.Username)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserRepository_Create(t *testing.T) {
	t.Run("Success - create user", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		user := &Domain.User{
			Username: "newuser",
			Password: "hashedpassword",
			Role:     Domain.RoleUser,
		}
		mockRepo.On("Create", user).Return(nil)

		// Act
		err := mockRepo.Create(user)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - create admin user", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		user := &Domain.User{
			Username: "adminuser",
			Password: "hashedpassword",
			Role:     Domain.RoleAdmin,
		}
		mockRepo.On("Create", user).Return(nil)

		// Act
		err := mockRepo.Create(user)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - database insert failed", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		user := &Domain.User{
			Username: "newuser",
			Password: "hashedpassword",
			Role:     Domain.RoleUser,
		}
		expectedError := errors.New("database insert failed")
		mockRepo.On("Create", user).Return(expectedError)

		// Act
		err := mockRepo.Create(user)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - duplicate username", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		user := &Domain.User{
			Username: "existinguser",
			Password: "hashedpassword",
			Role:     Domain.RoleUser,
		}
		expectedError := errors.New("duplicate username")
		mockRepo.On("Create", user).Return(expectedError)

		// Act
		err := mockRepo.Create(user)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserRepository_Update(t *testing.T) {
	t.Run("Success - update existing user", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		userID := primitive.NewObjectID().Hex()
		user := &Domain.User{
			Username:  "updateduser",
			Password:  "newhashedpassword",
			Role:      Domain.RoleAdmin,
			UpdatedAt: time.Now(),
		}
		mockRepo.On("Update", userID, user).Return(nil)

		// Act
		err := mockRepo.Update(userID, user)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - user not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		userID := primitive.NewObjectID().Hex()
		user := &Domain.User{
			Username: "updateduser",
			Role:     Domain.RoleUser,
		}
		expectedError := errors.New("user not found")
		mockRepo.On("Update", userID, user).Return(expectedError)

		// Act
		err := mockRepo.Update(userID, user)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - invalid ID format", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		invalidID := "invalid-id"
		user := &Domain.User{
			Username: "updateduser",
			Role:     Domain.RoleUser,
		}
		expectedError := errors.New("invalid user ID format")
		mockRepo.On("Update", invalidID, user).Return(expectedError)

		// Act
		err := mockRepo.Update(invalidID, user)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - database update failed", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		userID := primitive.NewObjectID().Hex()
		user := &Domain.User{
			Username: "updateduser",
			Role:     Domain.RoleUser,
		}
		expectedError := errors.New("database update failed")
		mockRepo.On("Update", userID, user).Return(expectedError)

		// Act
		err := mockRepo.Update(userID, user)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserRepository_UpdateByUsername(t *testing.T) {
	t.Run("Success - update user by username", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		username := "testuser"
		user := &Domain.User{
			Username:  username,
			Role:      Domain.RoleAdmin,
			UpdatedAt: time.Now(),
		}
		mockRepo.On("UpdateByUsername", username, user).Return(nil)

		// Act
		err := mockRepo.UpdateByUsername(username, user)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - user not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		username := "nonexistentuser"
		user := &Domain.User{
			Username: username,
			Role:     Domain.RoleAdmin,
		}
		expectedError := errors.New("user not found")
		mockRepo.On("UpdateByUsername", username, user).Return(expectedError)

		// Act
		err := mockRepo.UpdateByUsername(username, user)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - database update failed", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		username := "testuser"
		user := &Domain.User{
			Username: username,
			Role:     Domain.RoleAdmin,
		}
		expectedError := errors.New("database update failed")
		mockRepo.On("UpdateByUsername", username, user).Return(expectedError)

		// Act
		err := mockRepo.UpdateByUsername(username, user)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - promote user to admin", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		username := "regularuser"
		user := &Domain.User{
			Username: username,
			Role:     Domain.RoleAdmin, // Promoting to admin
		}
		mockRepo.On("UpdateByUsername", username, user).Return(nil)

		// Act
		err := mockRepo.UpdateByUsername(username, user)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserRepository_CountUsers(t *testing.T) {
	t.Run("Success - return user count", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		expectedCount := int64(5)
		mockRepo.On("CountUsers").Return(expectedCount, nil)

		// Act
		count, err := mockRepo.CountUsers()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - return zero count", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		expectedCount := int64(0)
		mockRepo.On("CountUsers").Return(expectedCount, nil)

		// Act
		count, err := mockRepo.CountUsers()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - return large count", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		expectedCount := int64(1000000)
		mockRepo.On("CountUsers").Return(expectedCount, nil)

		// Act
		count, err := mockRepo.CountUsers()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - database count failed", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		expectedError := errors.New("database count failed")
		mockRepo.On("CountUsers").Return(int64(0), expectedError)

		// Act
		count, err := mockRepo.CountUsers()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, int64(0), count)
		mockRepo.AssertExpectations(t)
	})
}

// Test interface compliance
func TestUserRepositoryInterface(t *testing.T) {
	mockRepo := new(MockUserRepositoryImpl)
	var _ UserRepositoryInterface = mockRepo
	assert.NotNil(t, mockRepo)
}

// Test edge cases
func TestUserRepository_EdgeCases(t *testing.T) {
	t.Run("Create user with empty fields", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		user := &Domain.User{
			Username: "",
			Password: "",
			Role:     Domain.RoleUser,
		}
		mockRepo.On("Create", user).Return(nil)

		// Act
		err := mockRepo.Create(user)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get user with very long username", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		longUsername := "very-long-username-that-might-cause-issues-in-some-systems-but-should-be-handled-gracefully-by-the-repository"
		expectedUser := &Domain.User{
			Username: longUsername,
			Role:     Domain.RoleUser,
		}
		mockRepo.On("GetByUsername", longUsername).Return(expectedUser, nil)

		// Act
		user, err := mockRepo.GetByUsername(longUsername)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update user with zero time values", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		userID := primitive.NewObjectID().Hex()
		user := &Domain.User{
			Username:  "testuser",
			Role:      Domain.RoleUser,
			CreatedAt: time.Time{}, // Zero time
			UpdatedAt: time.Time{}, // Zero time
		}
		mockRepo.On("Update", userID, user).Return(nil)

		// Act
		err := mockRepo.Update(userID, user)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get user with empty username", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		emptyUsername := ""
		expectedError := errors.New("username cannot be empty")
		mockRepo.On("GetByUsername", emptyUsername).Return(nil, expectedError)

		// Act
		user, err := mockRepo.GetByUsername(emptyUsername)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Count users with concurrent access", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepositoryImpl)
		expectedCount := int64(42)
		mockRepo.On("CountUsers").Return(expectedCount, nil)

		// Act
		count, err := mockRepo.CountUsers()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)
		mockRepo.AssertExpectations(t)
	})
}