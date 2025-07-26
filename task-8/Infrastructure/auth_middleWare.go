package Infrastructure

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"task_manager/Domain"
)

// AuthMiddleware provides authentication and authorization middleware
type AuthMiddleware struct {
	jwtService JWTServiceInterface
}

// NewAuthMiddleware creates a new instance of AuthMiddleware
func NewAuthMiddleware(jwtService JWTServiceInterface) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// AuthenticateToken validates JWT tokens
func (am *AuthMiddleware) AuthenticateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, Domain.ErrorResponse{
				Success: false,
				Message: "Authorization header required",
				Error:   "Missing Authorization header",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, Domain.ErrorResponse{
				Success: false,
				Message: "Invalid authorization header format",
				Error:   "Authorization header must be in format: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate token
		token, err := am.jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			errorMsg := "Token validation failed"
			if err != nil {
				errorMsg = err.Error()
			}
			c.JSON(http.StatusUnauthorized, Domain.ErrorResponse{
				Success: false,
				Message: "Invalid or expired token",
				Error:   errorMsg,
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, Domain.ErrorResponse{
				Success: false,
				Message: "Invalid token claims",
				Error:   "Could not parse token claims",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

// RequireAdmin ensures only admin users can access the endpoint
func (am *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, Domain.ErrorResponse{
				Success: false,
				Message: "User role not found",
				Error:   "Authentication required",
			})
			c.Abort()
			return
		}

		if role != Domain.RoleAdmin {
			c.JSON(http.StatusForbidden, Domain.ErrorResponse{
				Success: false,
				Message: "Access denied",
				Error:   "Admin privileges required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireUser ensures authenticated users (both admin and regular users) can access the endpoint
func (am *AuthMiddleware) RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, Domain.ErrorResponse{
				Success: false,
				Message: "User role not found",
				Error:   "Authentication required",
			})
			c.Abort()
			return
		}

		// Both admin and user roles are allowed
		if role != Domain.RoleAdmin && role != Domain.RoleUser {
			c.JSON(http.StatusForbidden, Domain.ErrorResponse{
				Success: false,
				Message: "Access denied",
				Error:   "Valid user role required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}