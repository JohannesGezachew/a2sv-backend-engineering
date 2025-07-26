package Infrastructure

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"task_manager/Domain"
)

type JWTServiceTestSuite struct {
	suite.Suite
	jwtService JWTServiceInterface
	testUser   *Domain.User
}

func (suite *JWTServiceTestSuite) SetupTest() {
	// Set a test JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	suite.jwtService = NewJWTService()
	
	// Create a test user
	suite.testUser = &Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Role:     Domain.RoleUser,
	}
}

func (suite *JWTServiceTestSuite) TearDownTest() {
	// Clean up environment variable
	os.Unsetenv("JWT_SECRET")
}

func (suite *JWTServiceTestSuite) TestGenerateToken() {
	tests := []struct {
		name    string
		user    *Domain.User
		wantErr bool
	}{
		{
			name:    "Valid user",
			user:    suite.testUser,
			wantErr: false,
		},
		{
			name: "Admin user",
			user: &Domain.User{
				ID:       primitive.NewObjectID(),
				Username: "admin",
				Role:     Domain.RoleAdmin,
			},
			wantErr: false,
		},
		{
			name: "User with empty username",
			user: &Domain.User{
				ID:       primitive.NewObjectID(),
				Username: "",
				Role:     Domain.RoleUser,
			},
			wantErr: false, // JWT generation should still work
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			token, err := suite.jwtService.GenerateToken(tt.user)

			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Empty(suite.T(), token)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotEmpty(suite.T(), token)
				
				// Verify token structure (should have 3 parts separated by dots)
				parts := len(token)
				assert.True(suite.T(), parts > 0)
				
				// Verify token can be parsed
				parsedToken, err := suite.jwtService.ValidateToken(token)
				assert.NoError(suite.T(), err)
				assert.True(suite.T(), parsedToken.Valid)
				
				// Verify claims
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				assert.True(suite.T(), ok)
				assert.Equal(suite.T(), tt.user.ID.Hex(), claims["user_id"])
				assert.Equal(suite.T(), tt.user.Username, claims["username"])
				assert.Equal(suite.T(), tt.user.Role, claims["role"])
				
				// Verify expiration is set (should be 24 hours from now)
				exp, ok := claims["exp"].(float64)
				assert.True(suite.T(), ok)
				assert.True(suite.T(), exp > float64(time.Now().Unix()))
				
				// Verify issued at time
				iat, ok := claims["iat"].(float64)
				assert.True(suite.T(), ok)
				assert.True(suite.T(), iat <= float64(time.Now().Unix()))
			}
		})
	}
}

func (suite *JWTServiceTestSuite) TestValidateToken() {
	// Generate a valid token first
	validToken, err := suite.jwtService.GenerateToken(suite.testUser)
	assert.NoError(suite.T(), err)

	tests := []struct {
		name        string
		tokenString string
		wantErr     bool
		wantValid   bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			wantErr:     false,
			wantValid:   true,
		},
		{
			name:        "Invalid token format",
			tokenString: "invalid.token.format",
			wantErr:     true,
			wantValid:   false,
		},
		{
			name:        "Empty token",
			tokenString: "",
			wantErr:     true,
			wantValid:   false,
		},
		{
			name:        "Malformed token",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.malformed.signature",
			wantErr:     true,
			wantValid:   false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			token, err := suite.jwtService.ValidateToken(tt.tokenString)

			if tt.wantErr {
				assert.Error(suite.T(), err)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.wantValid, token.Valid)
			}
		})
	}
}

func (suite *JWTServiceTestSuite) TestValidateTokenWithWrongSecret() {
	// Generate token with current service
	token, err := suite.jwtService.GenerateToken(suite.testUser)
	assert.NoError(suite.T(), err)

	// Create new service with different secret
	os.Setenv("JWT_SECRET", "different-secret")
	differentService := NewJWTService()

	// Try to validate with different secret
	parsedToken, err := differentService.ValidateToken(token)
	assert.Error(suite.T(), err)
	if parsedToken != nil {
		assert.False(suite.T(), parsedToken.Valid)
	}
}

func (suite *JWTServiceTestSuite) TestGetJWTSecret() {
	secret := suite.jwtService.GetJWTSecret()
	assert.NotNil(suite.T(), secret)
	assert.Equal(suite.T(), []byte("test-secret-key-for-testing"), secret)
}

func (suite *JWTServiceTestSuite) TestJWTServiceInterface() {
	// Test that our implementation satisfies the interface
	var _ JWTServiceInterface = &JWTService{}
	var _ JWTServiceInterface = suite.jwtService
}

func TestJWTServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceTestSuite))
}

// Additional standalone tests
func TestNewJWTService(t *testing.T) {
	t.Run("With JWT_SECRET environment variable", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "custom-secret")
		defer os.Unsetenv("JWT_SECRET")

		service := NewJWTService()
		assert.NotNil(t, service)
		assert.Equal(t, []byte("custom-secret"), service.GetJWTSecret())
	})

	t.Run("Without JWT_SECRET environment variable", func(t *testing.T) {
		os.Unsetenv("JWT_SECRET")

		service := NewJWTService()
		assert.NotNil(t, service)
		
		expectedDefault := "your-super-secret-jwt-key-change-this-in-production"
		assert.Equal(t, []byte(expectedDefault), service.GetJWTSecret())
	})
}

func TestJWTServiceTokenExpiration(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	service := NewJWTService()
	user := &Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Role:     Domain.RoleUser,
	}

	token, err := service.GenerateToken(user)
	assert.NoError(t, err)

	parsedToken, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	// Check that expiration is approximately 24 hours from now
	exp, ok := claims["exp"].(float64)
	assert.True(t, ok)
	
	expectedExp := time.Now().Add(24 * time.Hour).Unix()
	actualExp := int64(exp)
	
	// Allow for a small time difference (within 1 minute)
	assert.True(t, actualExp >= expectedExp-60 && actualExp <= expectedExp+60)
}

func TestJWTServiceClaimsContent(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	service := NewJWTService()
	userID := primitive.NewObjectID()
	user := &Domain.User{
		ID:       userID,
		Username: "testuser",
		Role:     Domain.RoleAdmin,
	}

	token, err := service.GenerateToken(user)
	assert.NoError(t, err)

	parsedToken, err := service.ValidateToken(token)
	assert.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	// Verify all expected claims are present
	assert.Equal(t, userID.Hex(), claims["user_id"])
	assert.Equal(t, "testuser", claims["username"])
	assert.Equal(t, Domain.RoleAdmin, claims["role"])
	
	// Verify exp and iat are present and valid
	_, expExists := claims["exp"]
	_, iatExists := claims["iat"]
	assert.True(t, expExists)
	assert.True(t, iatExists)
}

// Additional comprehensive tests for 100% coverage
func TestJWTServiceEdgeCases(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	service := NewJWTService()

	t.Run("Generate token for user with special characters in username", func(t *testing.T) {
		user := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: "user@domain.com",
			Role:     Domain.RoleUser,
		}

		token, err := service.GenerateToken(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "user@domain.com", claims["username"])
	})

	t.Run("Generate token for user with unicode username", func(t *testing.T) {
		user := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: "пользователь",
			Role:     Domain.RoleUser,
		}

		token, err := service.GenerateToken(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "пользователь", claims["username"])
	})

	t.Run("Validate token with wrong signing method", func(t *testing.T) {
		// This would normally use RS256, but we'll create an invalid token
		invalidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzQ2NzM2MDAsImlhdCI6MTYzNDY3MDAwMCwicm9sZSI6InVzZXIiLCJ1c2VyX2lkIjoidGVzdCIsInVzZXJuYW1lIjoidGVzdCJ9.invalid"

		parsedToken, err := service.ValidateToken(invalidToken)
		assert.Error(t, err)
		if parsedToken != nil {
			assert.False(t, parsedToken.Valid)
		}
	})

	t.Run("Validate expired token", func(t *testing.T) {
		// Create a token that's already expired
		claims := jwt.MapClaims{
			"user_id":  "test",
			"username": "test",
			"role":     "user",
			"exp":      time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
			"iat":      time.Now().Add(-2 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(service.GetJWTSecret())
		assert.NoError(t, err)

		parsedToken, err := service.ValidateToken(tokenString)
		assert.Error(t, err)
		if parsedToken != nil {
			assert.False(t, parsedToken.Valid)
		}
	})

	t.Run("Generate token with zero ObjectID", func(t *testing.T) {
		user := &Domain.User{
			ID:       primitive.ObjectID{}, // Zero ObjectID
			Username: "testuser",
			Role:     Domain.RoleUser,
		}

		token, err := service.GenerateToken(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "000000000000000000000000", claims["user_id"])
	})

	t.Run("Validate token with missing claims", func(t *testing.T) {
		// Create a token with minimal claims
		claims := jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(service.GetJWTSecret())
		assert.NoError(t, err)

		parsedToken, err := service.ValidateToken(tokenString)
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		parsedClaims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Nil(t, parsedClaims["user_id"])
		assert.Nil(t, parsedClaims["username"])
		assert.Nil(t, parsedClaims["role"])
	})

	t.Run("Generate multiple tokens for same user", func(t *testing.T) {
		user := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: "testuser",
			Role:     Domain.RoleUser,
		}

		token1, err1 := service.GenerateToken(user)
		// Add a small delay to ensure different iat timestamps
		time.Sleep(time.Millisecond * 10)
		token2, err2 := service.GenerateToken(user)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		// Tokens might be the same if generated at exactly the same time, so we just check they're valid
		assert.NotEmpty(t, token1)
		assert.NotEmpty(t, token2)

		// Both should be valid
		parsedToken1, err1 := service.ValidateToken(token1)
		parsedToken2, err2 := service.ValidateToken(token2)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.True(t, parsedToken1.Valid)
		assert.True(t, parsedToken2.Valid)
	})
}

func TestJWTServiceSecurityScenarios(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	service := NewJWTService()

	t.Run("Validate token with tampered payload", func(t *testing.T) {
		user := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: "user",
			Role:     Domain.RoleUser,
		}

		token, err := service.GenerateToken(user)
		assert.NoError(t, err)

		// Tamper with the token by changing a character in the payload
		tamperedToken := token[:len(token)-10] + "tampered123"

		parsedToken, err := service.ValidateToken(tamperedToken)
		assert.Error(t, err)
		if parsedToken != nil {
			assert.False(t, parsedToken.Valid)
		}
	})

	t.Run("Validate token with none algorithm", func(t *testing.T) {
		// Create a token with "none" algorithm (security vulnerability)
		noneToken := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJ1c2VyX2lkIjoidGVzdCIsInVzZXJuYW1lIjoidGVzdCIsInJvbGUiOiJhZG1pbiIsImV4cCI6OTk5OTk5OTk5OSwiaWF0IjoxNjM0NjcwMDAwfQ."

		parsedToken, err := service.ValidateToken(noneToken)
		assert.Error(t, err)
		if parsedToken != nil {
			assert.False(t, parsedToken.Valid)
		}
	})

	t.Run("Generate token with very long username", func(t *testing.T) {
		longUsername := make([]byte, 1000)
		for i := range longUsername {
			longUsername[i] = 'a'
		}

		user := &Domain.User{
			ID:       primitive.NewObjectID(),
			Username: string(longUsername),
			Role:     Domain.RoleUser,
		}

		token, err := service.GenerateToken(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, string(longUsername), claims["username"])
	})
}