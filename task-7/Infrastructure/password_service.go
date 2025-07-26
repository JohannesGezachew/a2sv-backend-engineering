package Infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordServiceInterface defines the contract for password operations
type PasswordServiceInterface interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

// PasswordService implements password hashing and comparison
type PasswordService struct{}

// NewPasswordService creates a new instance of PasswordService
func NewPasswordService() PasswordServiceInterface {
	return &PasswordService{}
}

// HashPassword hashes a plain text password
func (ps *PasswordService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a hashed password with a plain text password
func (ps *PasswordService) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}