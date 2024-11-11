package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// JwtService defines the interface for JSON Web Token (JWT) operations.
// This interface includes methods for creating access and refresh tokens, as well as parsing tokens.
type JwtService interface {
	CreateAccessToken(userId, sessionId int) (string, int64, error) // Method to create an access token
	CreateRefreshToken(userId, sessionId int) (string, error)       // Method to create a refresh token
	Parse(token string) (userId int, sessionId int, err error)      // Method to parse a token
}

// jwtService is a concrete implementation of JwtService.
// It holds the secret key used for signing tokens and configuration for token lifetimes.
type jwtService struct {
	secretKey                 []byte        // Secret key for signing tokens
	accessTokenLifetimeHours  time.Duration // Duration in hours before the access token expires
	refreshTokenLifetimeHours time.Duration // Duration in hours before the refresh token expires
}

// NewJwtService creates a new instance of jwtService.
// It initializes the service with a secret key.
//
// Parameters:
//   - secretKey: The secret key used for signing tokens.
//
// Returns:
//   - An instance of JwtService.
func NewJwtService(
	secretKey string,
	accessTokenLifetimeHours,
	refreshTokenLifetimeHours time.Duration,
) JwtService {
	return &jwtService{
		secretKey:                 []byte(secretKey),         // Convert secret key to byte slice
		accessTokenLifetimeHours:  accessTokenLifetimeHours,  // Set access token lifetime in hours
		refreshTokenLifetimeHours: refreshTokenLifetimeHours, // Set refresh token lifetime in hours
	}
}

// CreateAccessToken generates a new access token for a given user ID.
// The token will expire in 20 hours.
//
// Parameters:
//   - userId: The ID of the user for whom the access token is created.
//
// Returns:
//   - The generated token as a string, its expiration time as an int64, and any error encountered.
func (js *jwtService) CreateAccessToken(userId, sessionId int) (string, int64, error) {
	expiresAt := time.Now().Add(time.Hour * js.accessTokenLifetimeHours).UnixMilli() // Set expiration time to 20 hours from now

	// Create a new JWT with standard claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":    userId,
			"session_id": sessionId,
			"exp":        expiresAt,
		},
	)

	tokenString, err := token.SignedString(js.secretKey) // Sign the token with the secret key
	if err != nil {
		return "", 0, err // Return empty string and zero expiration time if signing fails
	}

	return tokenString, expiresAt, nil // Return the signed token and its expiration time
}

// CreateRefreshToken generates a new refresh token for a given user ID.
// The refresh token will expire in 1200 hours (50 days).
//
// Parameters:
//   - userId: The ID of the user for whom the refresh token is created.
//
// Returns:
//   - The generated refresh token as a string and any error encountered.
func (js *jwtService) CreateRefreshToken(userId, sessionId int) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS384,
		jwt.MapClaims{
			"user_id":    userId,
			"session_id": sessionId,
			"exp":        time.Now().Add(time.Hour * js.refreshTokenLifetimeHours).UnixMilli(),
		},
	)

	tokenString, err := refreshToken.SignedString(js.secretKey) // Sign the refresh token with the secret key
	if err != nil {
		return "", err // Return empty string if signing fails
	}

	return tokenString, nil // Return the signed refresh token
}

// Parse validates and parses a given JWT token.
// It retrieves the user ID from the claims if valid.
//
// Parameters:
//   - token: The JWT token to be parsed.
//
// Returns:
//   - The user ID as a string and any error encountered.
func (js *jwtService) Parse(token string) (userId int, sessionId int, err error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // Validate signing method
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return js.secretKey, nil // Return the secret key for validation
	})
	if err != nil {
		return 0, 0, err // Return empty string if parsing fails
	}

	if !t.Valid { // Check if the token is valid
		return 0, 0, errors.New("invalid token") // Return error if invalid
	}

	claims, ok := t.Claims.(jwt.MapClaims) // Retrieve claims from the parsed token
	if !ok {
		return 0, 0, errors.New("invalid claims") // Return error if claims are not valid
	}

	return int(claims["user_id"].(float64)), int(claims["session_id"].(float64)), nil // Return the user ID if successful
}
