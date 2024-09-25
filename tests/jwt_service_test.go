package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestJwtService_CreateAccessToken tests the CreateAccessToken function of the JwtService.
func TestJwtService_CreateAccessToken(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	// Define test cases with userId and sessionId
	tests := []struct {
		userId    int // User ID for token creation
		sessionId int // Session ID for token creation
	}{
		{1, 123}, // Test case 1
		{2, 456}, // Test case 2
	}

	for _, tt := range tests {
		t.Run("CreateAccessToken", func(t *testing.T) {
			// Create an access token using the userId and sessionId
			token, expiresAt, err := jwtService.CreateAccessToken(tt.userId, tt.sessionId)

			assert.NoError(t, err)                             // Ensure no error occurred during token creation
			assert.NotEmpty(t, token)                          // Ensure the token is not empty
			assert.True(t, time.Now().UnixMilli() < expiresAt) // Ensure the token is not expired
		})
	}
}

// TestJwtService_CreateRefreshToken tests the CreateRefreshToken function of the JwtService.
func TestJwtService_CreateRefreshToken(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	// Define test cases with userId and sessionId
	tests := []struct {
		userId    int // User ID for token creation
		sessionId int // Session ID for token creation
	}{
		{1, 123}, // Test case 1
		{2, 456}, // Test case 2
	}

	for _, tt := range tests {
		t.Run("CreateRefreshToken", func(t *testing.T) {
			// Create a refresh token using the userId and sessionId
			token, err := jwtService.CreateRefreshToken(tt.userId, tt.sessionId)

			assert.NoError(t, err)    // Ensure no error occurred during token creation
			assert.NotEmpty(t, token) // Ensure the refresh token is not empty
		})
	}
}

// TestJwtService_Parse_ValidToken tests the Parse function of the JwtService with a valid token.
func TestJwtService_Parse_ValidToken(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	userId := 1       // Define user ID for testing
	sessionId := 9879 // Define session ID for testing

	// Create a valid access token to be parsed later
	tokenString, _, _ := jwtService.CreateAccessToken(userId, sessionId)

	t.Run("Parse_ValidToken", func(t *testing.T) {
		// Parse the created token to retrieve user ID and session ID
		parsedUserId, parsedSessionId, err := jwtService.Parse(tokenString)

		assert.NoError(t, err)                      // Ensure no error occurred during parsing
		assert.Equal(t, userId, parsedUserId)       // Validate that parsed user ID matches expected user ID
		assert.Equal(t, sessionId, parsedSessionId) // Validate that parsed session ID matches expected session ID
	})
}

// TestJwtService_Parse_InvalidToken tests the Parse function of the JwtService with an invalid token.
func TestJwtService_Parse_InvalidToken(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	malformedToken := "invalid.token.string" // Define an invalid token string

	t.Run("Parse_InvalidToken", func(t *testing.T) {
		// Attempt to parse the malformed token and expect an error
		userId, sessionId, err := jwtService.Parse(malformedToken)

		assert.Error(t, err)          // Ensure an error occurred during parsing of invalid token
		assert.Equal(t, 0, userId)    // Validate that user ID is zero when parsing fails
		assert.Equal(t, 0, sessionId) // Validate that session ID is zero when parsing fails
	})
}
