package Usecases

import (
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"task_manager/Domain"
)

// MockUserRepository is a mock implementation of UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() ([]*Domain.User, error) {
	args := m.Called()
	return args.Get(0).([]*Domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id string) (*Domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*Domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *Domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(id string, user *Domain.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateByUsername(username string, user *Domain.User) error {
	args := m.Called(username, user)
	return args.Error(0)
}

func (m *MockUserRepository) CountUsers() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// MockPasswordService is a mock implementation of PasswordServiceInterface
type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) ComparePassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

// MockJWTService is a mock implementation of JWTServiceInterface
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(user *Domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func (m *MockJWTService) GetJWTSecret() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func TestUserUsecase_RegisterUser(t *testing.T) {
	t.Run("Success - register first user as admin", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		userReq := Domain.UserRequest{
			Username: "firstuser",
			Password: "password123",
		}
		hashedPassword := "hashed_password_123"

		mockUserRepo.On("GetByUsername", userReq.Username).Return(nil, errors.New("user not found"))
		mockPasswordService.On("HashPassword", userReq.Password).Return(hashedPassword, nil)
		mockUserRepo.On("CountUsers").Return(int64(0), nil)
		mockUserRepo.On("Create", mock.AnythingOfType("*Domain.User")).Return(nil)

		// Act
		user, err := userUsecase.RegisterUser(userReq)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userReq.Username, user.Username)
		assert.Equal(t, hashedPassword, user.Password)
		assert.Equal(t, Domain.RoleAdmin, user.Role) // First user should be admin

		mockUserRepo.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
	})

	t.Run("Success - register subsequent user as regular user", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		userReq := Domain.UserRequest{
			Username: "regularuser",
			Password: "password123",
		}
		hashedPassword := "hashed_password_123"

		mockUserRepo.On("GetByUsername", userReq.Username).Return(nil, errors.New("user not found"))
		mockPasswordService.On("HashPassword", userReq.Password).Return(hashedPassword, nil)
		mockUserRepo.On("CountUsers").Return(int64(1), nil) // Already has users
		mockUserRepo.On("Create", mock.AnythingOfType("*Domain.User")).Return(nil)

		// Act
		user, err := userUsecase.RegisterUser(userReq)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userReq.Username, user.Username)
		assert.Equal(t, hashedPassword, user.Password)
		assert.Equal(t, Domain.RoleUser, user.Role) // Subsequent users should be regular users

		mockUserRepo.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
	})

	t.Run("Error - username already exists", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		userReq := Domain.UserRequest{
			Username: "existinguser",
			Password: "password123",
		}
		existingUser := &Domain.User{
			Username: userReq.Username,
			Role:     Domain.RoleUser,
		}

		mockUserRepo.On("GetByUsername", userReq.Username).Return(existingUser, nil)

		// Act
		user, err := userUsecase.RegisterUser(userReq)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username already exists")
		assert.Nil(t, user)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error - password hashing fails", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		userReq := Domain.UserRequest{
			Username: "newuser",
			Password: "password123",
		}
		expectedError := errors.New("hashing error")

		mockUserRepo.On("GetByUsername", userReq.Username).Return(nil, errors.New("user not found"))
		mockPasswordService.On("HashPassword", userReq.Password).Return("", expectedError)

		// Act
		user, err := userUsecase.RegisterUser(userReq)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to hash password")
		assert.Nil(t, user)

		mockUserRepo.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
	})

	t.Run("Error - count users fails", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		userReq := Domain.UserRequest{
			Username: "newuser",
			Password: "password123",
		}
		hashedPassword := "hashed_password_123"
		expectedError := errors.New("database error")

		mockUserRepo.On("GetByUsername", userReq.Username).Return(nil, errors.New("user not found"))
		mockPasswordService.On("HashPassword", userReq.Password).Return(hashedPassword, nil)
		mockUserRepo.On("CountUsers").Return(int64(0), expectedError)

		// Act
		user, err := userUsecase.RegisterUser(userReq)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)

		mockUserRepo.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
	})

	t.Run("Error - create user fails", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		userReq := Domain.UserRequest{
			Username: "newuser",
			Password: "password123",
		}
		hashedPassword := "hashed_password_123"
		expectedError := errors.New("database create error")

		mockUserRepo.On("GetByUsername", userReq.Username).Return(nil, errors.New("user not found"))
		mockPasswordService.On("HashPassword", userReq.Password).Return(hashedPassword, nil)
		mockUserRepo.On("CountUsers").Return(int64(0), nil)
		mockUserRepo.On("Create", mock.AnythingOfType("*Domain.User")).Return(expectedError)

		// Act
		user, err := userUsecase.RegisterUser(userReq)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)

		mockUserRepo.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
	})
}

func TestUserUsecase_LoginUser(t *testing.T) {
	t.Run("Success - valid credentials", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		loginReq := Domain.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		user := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: loginReq.Username,
			Password: "hashed_password",
			Role:     Domain.RoleUser,
		}
		expectedToken := "jwt.token.here"

		mockUserRepo.On("GetByUsername", loginReq.Username).Return(user, nil)
		mockPasswordService.On("ComparePassword", user.Password, loginReq.Password).Return(nil)
		mockJWTService.On("GenerateToken", user).Return(expectedToken, nil)

		// Act
		resultUser, token, err := userUsecase.LoginUser(loginReq)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, user, resultUser)
		assert.Equal(t, expectedToken, token)

		mockUserRepo.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})

	t.Run("Error - user not found", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		loginReq := Domain.LoginRequest{
			Username: "nonexistentuser",
			Password: "password123",
		}
		expectedError := errors.New("user not found")

		mockUserRepo.On("GetByUsername", loginReq.Username).Return(nil, expectedError)

		// Act
		user, token, err := userUsecase.LoginUser(loginReq)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid credentials")
		assert.Nil(t, user)
		assert.Empty(t, token)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error - invalid password", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		loginReq := Domain.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}
		user := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: loginReq.Username,
			Password: "hashed_password",
			Role:     Domain.RoleUser,
		}
		expectedError := errors.New("password mismatch")

		mockUserRepo.On("GetByUsername", loginReq.Username).Return(user, nil)
		mockPasswordService.On("ComparePassword", user.Password, loginReq.Password).Return(expectedError)

		// Act
		resultUser, token, err := userUsecase.LoginUser(loginReq)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid credentials")
		assert.Nil(t, resultUser)
		assert.Empty(t, token)

		mockUserRepo.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
	})

	t.Run("Error - token generation fails", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		loginReq := Domain.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		user := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: loginReq.Username,
			Password: "hashed_password",
			Role:     Domain.RoleUser,
		}
		expectedError := errors.New("token generation error")

		mockUserRepo.On("GetByUsername", loginReq.Username).Return(user, nil)
		mockPasswordService.On("ComparePassword", user.Password, loginReq.Password).Return(nil)
		mockJWTService.On("GenerateToken", user).Return("", expectedError)

		// Act
		resultUser, token, err := userUsecase.LoginUser(loginReq)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to generate token")
		assert.Nil(t, resultUser)
		assert.Empty(t, token)

		mockUserRepo.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})
}

func TestUserUsecase_GetUserProfile(t *testing.T) {
	t.Run("Success - user found", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		userID := primitive.NewObjectID().Hex()
		expectedUser := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: "testuser",
			Role:     Domain.RoleUser,
		}

		mockUserRepo.On("GetByID", userID).Return(expectedUser, nil)

		// Act
		user, err := userUsecase.GetUserProfile(userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error - user not found", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		userID := primitive.NewObjectID().Hex()
		expectedError := errors.New("user not found")

		mockUserRepo.On("GetByID", userID).Return(nil, expectedError)

		// Act
		user, err := userUsecase.GetUserProfile(userID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_GetAllUsers(t *testing.T) {
	t.Run("Success - return all users", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		expectedUsers := []*Domain.User{
			{
				ID:       primitive.NewObjectID(),
				Username: "user1",
				Role:     Domain.RoleUser,
			},
			{
				ID:       primitive.NewObjectID(),
				Username: "admin1",
				Role:     Domain.RoleAdmin,
			},
		}

		mockUserRepo.On("GetAll").Return(expectedUsers, nil)

		// Act
		users, err := userUsecase.GetAllUsers()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Len(t, users, 2)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Success - return empty list", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		expectedUsers := []*Domain.User{}

		mockUserRepo.On("GetAll").Return(expectedUsers, nil)

		// Act
		users, err := userUsecase.GetAllUsers()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Len(t, users, 0)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error - repository error", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		expectedError := errors.New("database error")

		mockUserRepo.On("GetAll").Return([]*Domain.User(nil), expectedError)

		// Act
		users, err := userUsecase.GetAllUsers()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, users)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_PromoteUserToAdmin(t *testing.T) {
	t.Run("Success - promote user to admin", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		username := "usertoPromote"
		user := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: username,
			Role:     Domain.RoleUser,
		}
		promotedUser := &Domain.User{
			ID:       user.ID,
			Username: username,
			Role:     Domain.RoleAdmin,
		}

		mockUserRepo.On("GetByUsername", username).Return(user, nil).Once()
		mockUserRepo.On("UpdateByUsername", username, mock.AnythingOfType("*Domain.User")).Return(nil).Once()
		mockUserRepo.On("GetByUsername", username).Return(promotedUser, nil).Once()

		// Act
		resultUser, err := userUsecase.PromoteUserToAdmin(username)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, promotedUser, resultUser)
		assert.Equal(t, Domain.RoleAdmin, resultUser.Role)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error - user not found", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		username := "nonexistentuser"
		expectedError := errors.New("user not found")

		mockUserRepo.On("GetByUsername", username).Return(nil, expectedError)

		// Act
		user, err := userUsecase.PromoteUserToAdmin(username)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error - user is already admin", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		username := "adminuser"
		user := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: username,
			Role:     Domain.RoleAdmin,
		}

		mockUserRepo.On("GetByUsername", username).Return(user, nil)

		// Act
		resultUser, err := userUsecase.PromoteUserToAdmin(username)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user is already an admin")
		assert.Nil(t, resultUser)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error - update fails", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		username := "usertoPromote"
		user := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: username,
			Role:     Domain.RoleUser,
		}
		expectedError := errors.New("database update error")

		mockUserRepo.On("GetByUsername", username).Return(user, nil)
		mockUserRepo.On("UpdateByUsername", username, mock.AnythingOfType("*Domain.User")).Return(expectedError)

		// Act
		resultUser, err := userUsecase.PromoteUserToAdmin(username)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, resultUser)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error - get updated user fails", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		username := "usertoPromote"
		user := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: username,
			Role:     Domain.RoleUser,
		}
		expectedError := errors.New("user not found after update")

		mockUserRepo.On("GetByUsername", username).Return(user, nil).Once()
		mockUserRepo.On("UpdateByUsername", username, mock.AnythingOfType("*Domain.User")).Return(nil).Once()
		mockUserRepo.On("GetByUsername", username).Return(nil, expectedError).Once()

		// Act
		resultUser, err := userUsecase.PromoteUserToAdmin(username)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, resultUser)

		mockUserRepo.AssertExpectations(t)
	})
}

// Additional standalone tests
func TestNewUserUsecase(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockPasswordService := new(MockPasswordService)
	mockJWTService := new(MockJWTService)

	usecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

	assert.NotNil(t, usecase)
	assert.Implements(t, (*UserUsecaseInterface)(nil), usecase)
}

func TestUserUsecaseInterface(t *testing.T) {
	// Test that our implementation satisfies the interface
	mockUserRepo := new(MockUserRepository)
	mockPasswordService := new(MockPasswordService)
	mockJWTService := new(MockJWTService)

	var _ UserUsecaseInterface = &UserUsecase{
		userRepo:        mockUserRepo,
		passwordService: mockPasswordService,
		jwtService:      mockJWTService,
	}
	var _ UserUsecaseInterface = NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)
}

// Edge case tests
func TestUserUsecase_EdgeCases(t *testing.T) {
	t.Run("Register user with minimum password length", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		userReq := Domain.UserRequest{
			Username: "testuser",
			Password: "123456", // Minimum 6 characters
		}
		hashedPassword := "hashed_123456"

		mockUserRepo.On("GetByUsername", userReq.Username).Return(nil, errors.New("user not found"))
		mockPasswordService.On("HashPassword", userReq.Password).Return(hashedPassword, nil)
		mockUserRepo.On("CountUsers").Return(int64(1), nil)
		mockUserRepo.On("Create", mock.AnythingOfType("*Domain.User")).Return(nil)

		// Act
		user, err := userUsecase.RegisterUser(userReq)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, Domain.RoleUser, user.Role)

		mockUserRepo.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
	})

	t.Run("Login with admin user", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		loginReq := Domain.LoginRequest{
			Username: "admin",
			Password: "adminpass",
		}
		adminUser := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: loginReq.Username,
			Password: "hashed_adminpass",
			Role:     Domain.RoleAdmin,
		}
		expectedToken := "admin.jwt.token"

		mockUserRepo.On("GetByUsername", loginReq.Username).Return(adminUser, nil)
		mockPasswordService.On("ComparePassword", adminUser.Password, loginReq.Password).Return(nil)
		mockJWTService.On("GenerateToken", adminUser).Return(expectedToken, nil)

		// Act
		resultUser, token, err := userUsecase.LoginUser(loginReq)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, adminUser, resultUser)
		assert.Equal(t, expectedToken, token)
		assert.Equal(t, Domain.RoleAdmin, resultUser.Role)

		mockUserRepo.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})

	t.Run("Get profile for admin user", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		mockJWTService := new(MockJWTService)
		userUsecase := NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

		userID := primitive.NewObjectID().Hex()
		adminUser := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: "admin",
			Role:     Domain.RoleAdmin,
		}

		mockUserRepo.On("GetByID", userID).Return(adminUser, nil)

		// Act
		user, err := userUsecase.GetUserProfile(userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, adminUser, user)
		assert.Equal(t, Domain.RoleAdmin, user.Role)

		mockUserRepo.AssertExpectations(t)
	})
}