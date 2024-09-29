package tests

import (
	"cvs/internal/models"
	"cvs/internal/service"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCheckUserDataService tests the CheckUserData function of the service package.
func TestCheckUserDataService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string      // Name of the test case
		inputUser   models.User // User data to be validated
		expectedErr error       // Expected error result from validation
	}{
		{
			name: "Ok", // Test case for valid user data
			inputUser: models.User{
				ID:       1,
				Email:    "test@test.test",
				Password: []byte("password"),
			},
			expectedErr: nil, // No error expected for valid input
		},
		{
			name: "Error. Email data is empty", // Test case for empty email
			inputUser: models.User{
				ID:       1,
				Email:    "",
				Password: []byte("password"),
			},
			expectedErr: errors.New("email data is empty"), // Expected error for empty email
		},
		{
			name: "Error. Invalid email format", // Test case for invalid email format
			inputUser: models.User{
				ID:       1,
				Email:    "test", // Invalid email format
				Password: []byte("password"),
			},
			expectedErr: errors.New("invalid email format"), // Expected error for invalid email format
		},
		{
			name: "Error. User password value is empty", // Test case for empty password
			inputUser: models.User{
				ID:       1,
				Email:    "test@test.test",
				Password: []byte(""), // Empty password
			},
			expectedErr: errors.New("user password value is empty"), // Expected error for empty password
		},
	}

	for _, test := range tests {
		tc := test // Create a copy of the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run this test case in parallel

			err := service.CheckUserData(tc.inputUser) // Call the function to validate user data

			if tc.name == "Ok" {
				assert.NoError(t, err) // Check that no error occurred for valid input
			} else {
				assert.EqualError(t, tc.expectedErr, err.Error()) // Check that the expected error matches the actual error
			}
		})
	}
}

// TestCheckPairDataService tests the CheckPairData function of the service package.
func TestCheckPairDataService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string           // Name of the test case
		inputPairData models.UserPairs // Pair data to be validated
		expectedErr   error            // Expected error result from validation
	}{
		{
			name: "Ok", // Test case for valid pair data
			inputPairData: models.UserPairs{
				UserID:     1,
				Exchange:   "binance_spot",
				Pair:       "BTC/USDT",
				ExactValue: 1,
			},
			expectedErr: nil, // No error expected for valid input
		},
		{
			name: "Error. Pair name is empty", // Test case for empty pair name
			inputPairData: models.UserPairs{
				UserID:     1,
				Exchange:   "binance_spot",
				Pair:       "", // Empty pair name
				ExactValue: 1,
			},
			expectedErr: errors.New("pair name is empty"), // Expected error for empty pair name
		},
		{
			name: "Error. Exchange name is empty", // Test case for empty exchange name
			inputPairData: models.UserPairs{
				UserID:     1,
				Exchange:   "", // Empty exchange name
				Pair:       "BTC/USDT",
				ExactValue: 1,
			},
			expectedErr: errors.New("exchange name is empty"), // Expected error for empty exchange name
		},
		{
			name: "Error. Exact value must be above zero", // Test case for exact value <= 0
			inputPairData: models.UserPairs{
				UserID:     1,
				Exchange:   "binance_spot",
				Pair:       "BTC/USDT",
				ExactValue: 0, // Invalid exact value (0)
			},
			expectedErr: errors.New("exact value must be above zero"), // Expected error for invalid exact value
		},
		{
			name: "Error. User id must be above zero", // Test case for invalid user ID (0)
			inputPairData: models.UserPairs{
				UserID:     0, // Invalid user ID (0)
				Exchange:   "binance_spot",
				Pair:       "BTC/USDT",
				ExactValue: 1,
			},
			expectedErr: errors.New("user id must be above zero"), // Expected error for invalid user ID
		},
		{
			name: "Error. Invalid pair name format", // Test case for invalid pair name format
			inputPairData: models.UserPairs{
				UserID:     1,
				Exchange:   "binance_spot",
				Pair:       "test.", // Invalid pair format (ends with a dot)
				ExactValue: 1,
			},
			expectedErr: errors.New("invalid pair name format"), // Expected error for invalid pair name format
		},
		{
			name: "Error. Invalid exchange name format", // Test case for invalid exchange name format
			inputPairData: models.UserPairs{
				UserID:     1,
				Exchange:   "test", // Invalid exchange format (not matching expected patterns)
				Pair:       "BTC/USDT",
				ExactValue: 1,
			},
			expectedErr: errors.New("invalid exchange name format"), // Expected error for invalid exchange name format
		},
	}

	for _, test := range tests {
		tc := test

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := service.CheckPairData(tc.inputPairData) // Call the function to validate pair data

			if tc.name == "Ok" {
				assert.NoError(t, err) // Check that no error occurred for valid input
			} else {
				assert.EqualError(t, tc.expectedErr, err.Error()) // Check that the expected error matches the actual error
			}
		})
	}
}
