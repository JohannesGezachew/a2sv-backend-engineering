package Infrastructure

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"task_manager/Domain"
)

// JWTServiceInterface defines the contract for JWT operations
type JWTServiceInterface interface {
	GenerateToken(user *Domain.User) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetJWTSecret() []byte
}

// JWTService implements JWT token operations
type JWTService struct {
	secret []byte
}

// NewJWTService creates a new instance of JWTService
func NewJWTService() JWTServiceInterface {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-super-secret-jwt-key-change-this-in-production"
	}
	
	return &JWTService{
		secret: []byte(secret),
	}
}

// GenerateToken generates a JWT token for a user
func (js *JWTService) GenerateToken(user *Domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.Hex(),
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(js.secret)
}

// ValidateToken validates a JWT token and returns the parsed token
func (js *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return js.secret, nil
	})
}

// GetJWTSecret returns the JWT secret key
func (js *JWTService) GetJWTSecret() []byte {
	return js.secret
}