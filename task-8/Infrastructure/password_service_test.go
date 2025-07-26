package Infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PasswordServiceTestSuite struct {
	suite.Suite
	passwordService PasswordServiceInterface
}

func (suite *PasswordServiceTestSuite) SetupTest() {
	suite.passwordService = NewPasswordService()
}

func (suite *PasswordServiceTestSuite) TestHashPassword() {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false, // bcrypt can hash empty strings
		},
		{
			name:     "Long password (within bcrypt limit)",
			password: "this_is_a_long_password_that_is_within_bcrypt_72_byte_limit_test",
			wantErr:  false,
		},
		{
			name:     "Password with special characters",
			password: "p@ssw0rd!@#$%^&*()",
			wantErr:  false,
		},
		{
			name:     "Unicode password",
			password: "пароль123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			hashedPassword, err := suite.passwordService.HashPassword(tt.password)

			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Empty(suite.T(), hashedPassword)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotEmpty(suite.T(), hashedPassword)
				assert.NotEqual(suite.T(), tt.password, hashedPassword)
				assert.True(suite.T(), len(hashedPassword) > 0)
			}
		})
	}
}

func (suite *PasswordServiceTestSuite) TestComparePassword() {
	password := "testpassword123"
	hashedPassword, err := suite.passwordService.HashPassword(password)
	assert.NoError(suite.T(), err)

	tests := []struct {
		name           string
		hashedPassword string
		plainPassword  string
		wantErr        bool
	}{
		{
			name:           "Correct password",
			hashedPassword: hashedPassword,
			plainPassword:  password,
			wantErr:        false,
		},
		{
			name:           "Incorrect password",
			hashedPassword: hashedPassword,
			plainPassword:  "wrongpassword",
			wantErr:        true,
		},
		{
			name:           "Empty plain password",
			hashedPassword: hashedPassword,
			plainPassword:  "",
			wantErr:        true,
		},
		{
			name:           "Invalid hash format",
			hashedPassword: "invalid_hash",
			plainPassword:  password,
			wantErr:        true,
		},
		{
			name:           "Empty hash",
			hashedPassword: "",
			plainPassword:  password,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.passwordService.ComparePassword(tt.hashedPassword, tt.plainPassword)

			if tt.wantErr {
				assert.Error(suite.T(), err)
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

func (suite *PasswordServiceTestSuite) TestHashPasswordConsistency() {
	password := "consistencytest"

	// Hash the same password multiple times
	hash1, err1 := suite.passwordService.HashPassword(password)
	hash2, err2 := suite.passwordService.HashPassword(password)

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NotEqual(suite.T(), hash1, hash2) // Hashes should be different due to salt
	
	// But both should validate correctly
	assert.NoError(suite.T(), suite.passwordService.ComparePassword(hash1, password))
	assert.NoError(suite.T(), suite.passwordService.ComparePassword(hash2, password))
}

func (suite *PasswordServiceTestSuite) TestPasswordServiceInterface() {
	// Test that our implementation satisfies the interface
	var _ PasswordServiceInterface = &PasswordService{}
	var _ PasswordServiceInterface = suite.passwordService
}

func TestPasswordServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordServiceTestSuite))
}

// Additional standalone tests
func TestNewPasswordService(t *testing.T) {
	service := NewPasswordService()
	assert.NotNil(t, service)
	assert.Implements(t, (*PasswordServiceInterface)(nil), service)
}

func TestPasswordServiceEdgeCases(t *testing.T) {
	service := NewPasswordService()

	t.Run("Very long password (exceeds bcrypt limit)", func(t *testing.T) {
		// Test with a password longer than bcrypt's 72 byte limit
		longPassword := make([]byte, 100)
		for i := range longPassword {
			longPassword[i] = 'a'
		}
		
		hashedPassword, err := service.HashPassword(string(longPassword))
		// bcrypt will return an error for passwords > 72 bytes
		assert.Error(t, err)
		assert.Empty(t, hashedPassword)
	})

	t.Run("Password with null bytes", func(t *testing.T) {
		passwordWithNull := "password\x00with\x00null"
		
		hashedPassword, err := service.HashPassword(passwordWithNull)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		
		err = service.ComparePassword(hashedPassword, passwordWithNull)
		assert.NoError(t, err)
	})

	t.Run("Password with only special characters", func(t *testing.T) {
		specialPassword := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
		
		hashedPassword, err := service.HashPassword(specialPassword)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		
		err = service.ComparePassword(hashedPassword, specialPassword)
		assert.NoError(t, err)
	})

	t.Run("Password with mixed case and numbers", func(t *testing.T) {
		mixedPassword := "MyP@ssw0rd123"
		
		hashedPassword, err := service.HashPassword(mixedPassword)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		
		err = service.ComparePassword(hashedPassword, mixedPassword)
		assert.NoError(t, err)
	})

	t.Run("Compare with completely different password", func(t *testing.T) {
		originalPassword := "original123"
		differentPassword := "different456"
		
		hashedPassword, err := service.HashPassword(originalPassword)
		assert.NoError(t, err)
		
		err = service.ComparePassword(hashedPassword, differentPassword)
		assert.Error(t, err)
	})

	t.Run("Compare with truncated hash", func(t *testing.T) {
		password := "testpassword"
		hashedPassword, err := service.HashPassword(password)
		assert.NoError(t, err)
		
		// Truncate the hash to make it invalid
		truncatedHash := hashedPassword[:10]
		
		err = service.ComparePassword(truncatedHash, password)
		assert.Error(t, err)
	})

	t.Run("Compare with malformed hash", func(t *testing.T) {
		password := "testpassword"
		malformedHash := "not.a.valid.bcrypt.hash"
		
		err := service.ComparePassword(malformedHash, password)
		assert.Error(t, err)
	})
}