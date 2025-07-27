package Usecases

import (
	"errors"

	"task_manager/Domain"
	"task_manager/Infrastructure"
	"task_manager/Repositories"
)

// UserUsecaseInterface defines the contract for user business logic
type UserUsecaseInterface interface {
	RegisterUser(userReq Domain.UserRequest) (*Domain.User, error)
	LoginUser(loginReq Domain.LoginRequest) (*Domain.User, string, error)
	GetUserProfile(userID string) (*Domain.User, error)
	GetAllUsers() ([]*Domain.User, error)
	PromoteUserToAdmin(username string) (*Domain.User, error)
}

// UserUsecase implements user business logic
type UserUsecase struct {
	userRepo        Repositories.UserRepositoryInterface
	passwordService Infrastructure.PasswordServiceInterface
	jwtService      Infrastructure.JWTServiceInterface
}

// NewUserUsecase creates a new instance of UserUsecase
func NewUserUsecase(
	userRepo Repositories.UserRepositoryInterface,
	passwordService Infrastructure.PasswordServiceInterface,
	jwtService Infrastructure.JWTServiceInterface,
) UserUsecaseInterface {
	return &UserUsecase{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

// RegisterUser creates a new user
func (uu *UserUsecase) RegisterUser(userReq Domain.UserRequest) (*Domain.User, error) {
	// Check if username already exists
	existingUser, _ := uu.userRepo.GetByUsername(userReq.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Hash the password
	hashedPassword, err := uu.passwordService.HashPassword(userReq.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Check if this is the first user (make them admin)
	userCount, err := uu.userRepo.CountUsers()
	if err != nil {
		return nil, err
	}

	role := Domain.RoleUser
	if userCount == 0 {
		role = Domain.RoleAdmin
	}

	user := &Domain.User{
		Username: userReq.Username,
		Password: hashedPassword,
		Role:     role,
	}

	err = uu.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser authenticates a user and returns user info with JWT token
func (uu *UserUsecase) LoginUser(loginReq Domain.LoginRequest) (*Domain.User, string, error) {
	user, err := uu.userRepo.GetByUsername(loginReq.Username)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Compare password with hash
	err = uu.passwordService.ComparePassword(user.Password, loginReq.Password)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := uu.jwtService.GenerateToken(user)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

// GetUserProfile returns user profile by ID
func (uu *UserUsecase) GetUserProfile(userID string) (*Domain.User, error) {
	return uu.userRepo.GetByID(userID)
}

// GetAllUsers returns all users (admin only)
func (uu *UserUsecase) GetAllUsers() ([]*Domain.User, error) {
	return uu.userRepo.GetAll()
}

// PromoteUserToAdmin promotes a user to admin role
func (uu *UserUsecase) PromoteUserToAdmin(username string) (*Domain.User, error) {
	user, err := uu.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if user.Role == Domain.RoleAdmin {
		return nil, errors.New("user is already an admin")
	}

	user.Role = Domain.RoleAdmin
	err = uu.userRepo.UpdateByUsername(username, user)
	if err != nil {
		return nil, err
	}

	// Return updated user
	return uu.userRepo.GetByUsername(username)
}