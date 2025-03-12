// backend/internal/api/middleware/auth.go

package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTMiddleware struct {
	secret     string
	expiration time.Duration
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func NewJWTMiddleware(secret string, expirationHours int) *JWTMiddleware {
	return &JWTMiddleware{
		secret:     secret,
		expiration: time.Duration(expirationHours) * time.Hour,
	}
}

// GenerateToken creates a new JWT token for a user
func (j *JWTMiddleware) GenerateToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.expiration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// AuthRequired is a middleware that checks for a valid JWT token
func (j *JWTMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := j.extractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims, err := j.validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Set user ID in context for later use
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// extractToken gets the token from the Authorization header
func (j *JWTMiddleware) extractToken(c *gin.Context) (string, error) {
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		return "", errors.New("no authorization header")
	}

	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}

// validateToken checks if the token is valid and returns the claims
func (j *JWTMiddleware) validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetUserID retrieves the user ID from the context
func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", errors.New("user ID not found in context")
	}

	userIDString, ok := userID.(string)
	if !ok {
		return "", errors.New("user ID is not a string")
	}

	return userIDString, nil
}

// RefreshToken generates a new token for a user
func (j *JWTMiddleware) RefreshToken(c *gin.Context) {
	oldToken, err := j.extractToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, err := j.validateToken(oldToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Generate new token
	newToken, err := j.GenerateToken(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}

// Optional middleware that will attempt to authenticate but not require it
func (j *JWTMiddleware) AuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := j.extractToken(c)
		if err == nil {
			claims, err := j.validateToken(token)
			if err == nil {
				c.Set("userID", claims.UserID)
			}
		}
		c.Next()
	}
}

// IsAuthenticated helper function to check if user is authenticated
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("userID")
	return exists
}
