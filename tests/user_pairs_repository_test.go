package tests

import (
	"context"
	"cvs/internal/models"
	"cvs/internal/repository"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	// Run tests in parallel to speed up execution
	t.Parallel()

	// Define test cases for adding user pairs
	tests := []struct {
		name     string           // Name of the test case
		pairData models.UserPairs // Data for the user pair being tested
		wantErr  bool             // Expectation of whether an error should occur
	}{
		{
			name: "Valid User Pair",
			pairData: models.UserPairs{
				Exchange:   "Binance",
				Pair:       "BTC/USDT",
				ExactValue: 45000,
			},
			wantErr: false, // No error expected for valid input
		},
		{
			name: "Invalid User Pair - Empty UserID",
			pairData: models.UserPairs{
				UserID:     0, // Invalid user ID (0)
				Exchange:   "Binance",
				Pair:       "BTC/USDT",
				ExactValue: 45000,
			},
			wantErr: true, // Error expected due to invalid UserID
		},
		{
			name: "Invalid User Pair - Non-existent User",
			pairData: models.UserPairs{
				UserID:     99999, // Assuming this user ID does not exist in the users table
				Exchange:   "Binance",
				Pair:       "BTC/USDT",
				ExactValue: 45000,
			},
			wantErr: true, // Error expected for non-existent user
		},
		{
			name: "Invalid User Pair - Empty Exchange",
			pairData: models.UserPairs{
				UserID:     1,
				Exchange:   "", // Empty exchange string
				Pair:       "BTC/USDT",
				ExactValue: 45000,
			},
			wantErr: true, // Error expected due to empty exchange
		},
		{
			name: "Invalid User Pair - Empty Pair",
			pairData: models.UserPairs{
				UserID:     1,
				Exchange:   "Binance",
				Pair:       "", // Empty pair string
				ExactValue: 45000,
			},
			wantErr: true, // Error expected due to empty pair
		},
		{
			name: "Invalid User Pair - Non-positive Exact Value",
			pairData: models.UserPairs{
				UserID:     1,
				Exchange:   "Binance",
				Pair:       "BTC/USDT",
				ExactValue: -100, // Invalid exact value (negative)
			},
			wantErr: true, // Error expected due to non-positive exact value
		},
	}

	// Iterate through each test case
	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			db := setupDB()  // Setup a new database connection for each test case
			defer db.Close() // Ensure the database connection is closed after the test

			// If the test case expects no error and is a valid user pair, create a user in the database
			if !tc.wantErr && tc.name == "Valid User Pair" {
				email := "validuser" + tc.name + "@example.com"                  // Create a unique email for each test case
				userID, err := insertUser(db, email, []byte("validpassword123")) // Insert a valid user into the database
				defer db.ExecContext(ctx, deleteUserQueryRow, userID)            // Clean up by deleting the user after the test

				assert.NoError(t, err) // Assert that there was no error inserting the user

				tc.pairData.UserID = userID // Set the valid user ID in the pair data for this test case
			}

			repo := repository.NewUserPairsRepository(db) // Create a new repository instance for user pairs
			err := repo.Add(ctx, tc.pairData)             // Attempt to add the user pair

			if tc.wantErr {
				assert.Error(t, err) // Assert that an error occurred if one was expected
			} else {
				assert.NoError(t, err) // Assert that no error occurred for valid input

				var retrievedPair models.UserPairs

				query := `SELECT user_id, exchange, pair, exact_value FROM user_pairs WHERE user_id = $1 AND pair = $2`
				err = db.GetContext(ctx, &retrievedPair, query, tc.pairData.UserID, tc.pairData.Pair) // Retrieve the added user pair from the database

				assert.NoError(t, err)                                            // Assert that there was no error retrieving the data
				assert.Equal(t, tc.pairData.Exchange, retrievedPair.Exchange)     // Check that the exchange matches what was added
				assert.Equal(t, tc.pairData.ExactValue, retrievedPair.ExactValue) // Check that the exact value matches what was added
			}
		})
	}
}

func TestUpdateExactValue(t *testing.T) {
	// Run tests in parallel to improve execution speed
	t.Parallel()

	// Define test cases for updating the exact value of user pairs
	tests := []struct {
		name     string           // Name of the test case
		pairData models.UserPairs // Data for the user pair being tested
		wantErr  bool             // Expectation of whether an error should occur
	}{
		{
			name: "Valid Update",
			pairData: models.UserPairs{
				Exchange:   "Binance",
				Pair:       "BTC/USDT",
				ExactValue: 45000,
			},
			wantErr: false, // No error expected for valid input
		},
		{
			name: "Invalid Update - Non-existent User ID",
			pairData: models.UserPairs{
				UserID:     99999, // Assuming this user ID does not exist in the users table
				ExactValue: 50000,
			},
			wantErr: true, // Error expected due to non-existent UserID
		},
		{
			name: "Invalid Update - Non-existent Pair",
			pairData: models.UserPairs{
				UserID:     1, // Assuming this user ID exists in the users table
				Exchange:   "Binance",
				Pair:       "NON_EXISTENT_PAIR", // Assuming this pair does not exist
				ExactValue: 50000,
			},
			wantErr: true, // Error expected due to non-existent pair
		},
	}

	// Iterate through each test case
	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			db := setupDB()  // Setup a new database connection for each test case
			defer db.Close() // Ensure the database connection is closed after the test

			// Create a user and insert a valid pair for valid test cases
			if !tc.wantErr && tc.name == "Valid Update" {
				email := "validuser@example.com"                                 // Unique email for testing
				userID, err := insertUser(db, email, []byte("validpassword123")) // Insert a valid user into the database
				defer db.ExecContext(ctx, deleteUserQueryRow, userID)            // Clean up by deleting the user after the test
				assert.NoError(t, err)                                           // Assert that there was no error inserting the user

				// Insert a valid pair into the database for updating later
				insertQuery := `INSERT INTO user_pairs (user_id, exchange, pair, exact_value) VALUES ($1, $2, $3, $4)`
				_, err = db.ExecContext(ctx, insertQuery, userID, tc.pairData.Exchange, tc.pairData.Pair, tc.pairData.ExactValue)
				assert.NoError(t, err) // Assert that there was no error inserting the pair

				tc.pairData.UserID = userID // Set the valid user ID in the pair data for this test case
			}

			repo := repository.NewUserPairsRepository(db)  // Create a new repository instance for user pairs
			err := repo.UpdateExactValue(ctx, tc.pairData) // Attempt to update the exact value

			if tc.wantErr {
				assert.Error(t, err) // Assert that an error occurred if one was expected
			} else {
				assert.NoError(t, err) // Assert that no error occurred for valid input

				var retrievedPair models.UserPairs
				query := `SELECT exact_value FROM user_pairs WHERE user_id = $1 AND pair = $2`
				err = db.GetContext(ctx, &retrievedPair, query, tc.pairData.UserID, tc.pairData.Pair) // Retrieve the updated exact value from the database

				assert.NoError(t, err)                                            // Assert that there was no error retrieving the data
				assert.Equal(t, tc.pairData.ExactValue, retrievedPair.ExactValue) // Check that the exact value matches what was updated
			}
		})
	}
}

func TestGetAllUserPairs(t *testing.T) {
	// Run tests in parallel to improve execution speed
	t.Parallel()

	// Define test cases for retrieving all user pairs
	tests := []struct {
		name          string // Name of the test case
		userID        int    // User ID for which pairs are to be fetched
		expectedCount int    // Expected number of pairs to be returned
		ctxTimeout    time.Duration
		wantErr       bool // Expectation of whether an error should occur
	}{
		{
			name:          "Get All User Pairs - Valid User",
			expectedCount: 2, // Expecting 2 pairs for a valid user
			ctxTimeout:    contextTimeout,
			wantErr:       false, // No error expected for valid input
		},
		{
			name:          "Get All User Pairs - No Pairs Found",
			userID:        -1, // Assuming this ID does not exist in the database
			expectedCount: 0,  // Expecting 0 pairs for a non-existent user
			ctxTimeout:    contextTimeoutZero,
			wantErr:       true, // Error expected due to non-existent user ID
		},
	}

	// Iterate through each test case
	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			db := setupDB()  // Setup a new database connection for each test case
			defer db.Close() // Ensure the database connection is closed after the test

			// Add pairs for the valid test case if no error is expected
			if !tc.wantErr {
				email := "validuser_" + tc.name + "@example.com"                 // Generate a unique email for testing
				userID, err := insertUser(db, email, []byte("validpassword123")) // Insert a valid user into the database
				defer db.ExecContext(ctx, deleteUserQueryRow, userID)            // Clean up by deleting the user after the test

				assert.NoError(t, err) // Assert that there was no error inserting the user

				// Insert valid pairs for the existing user
				err = insertUserPair(db, userID, "Binance", "BTC/USDT", 45000) // Insert first pair
				assert.NoError(t, err)                                         // Assert that there was no error inserting the pair

				err = insertUserPair(db, userID, "Coinbase", "ETH/USDT", 3000) // Insert second pair
				assert.NoError(t, err)                                         // Assert that there was no error inserting the pair

				tc.userID = userID // Set the valid user ID in the test case data
			}

			ctx, cancel := context.WithTimeout(ctx, tc.ctxTimeout) // Set up context with timeout
			defer cancel()

			repo := repository.NewUserPairsRepository(db)      // Create a new repository instance for user pairs
			pairs, err := repo.GetAllUserPairs(ctx, tc.userID) // Attempt to retrieve all pairs for the specified user ID

			if tc.wantErr {
				assert.Error(t, err) // Assert that an error occurred if one was expected
				return               // Exit early since we expect an error and no pairs should be returned
			}

			assert.NoError(t, err)                 // Assert that no error occurred for valid input
			assert.Len(t, pairs, tc.expectedCount) // Assert that the number of retrieved pairs matches the expected count

			for _, p := range pairs {
				assert.Equal(t, tc.userID, p.UserID) // Check that each retrieved pair belongs to the correct user ID
			}
		})
	}
}
func TestGetPairsByExchange(t *testing.T) {
	// Run tests in parallel to improve execution speed
	t.Parallel()

	// Define test cases for retrieving pairs by exchange
	tests := []struct {
		name          string   // Name of the test case
		exchange      string   // Exchange name for which pairs are to be fetched
		expectedPairs []string // Expected pairs to be returned
		wantErr       bool     // Expectation of whether an error should occur
	}{
		{
			name:          "Valid Exchange - Binance",
			exchange:      "Binance",                        // Valid exchange with expected pairs
			expectedPairs: []string{"BTC/USDT", "ETH/USDT"}, // Expecting these pairs for Binance
			wantErr:       false,                            // No error expected for valid input
		},
		{
			name:          "Valid Exchange - Coinbase",
			exchange:      "Coinbase",           // Another valid exchange with expected pairs
			expectedPairs: []string{"ETH/USDT"}, // Expecting this pair for Coinbase
			wantErr:       false,                // No error expected for valid input
		},
		{
			name:          "Invalid Exchange - NonExistent",
			exchange:      "NonExistent", // Exchange that does not exist in the database
			expectedPairs: []string{},    // Expecting an empty result for a non-existent exchange
			wantErr:       false,         // No error expected, just empty result
		},
	}

	// Iterate through each test case
	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			db := setupDB()  // Setup a new database connection for each test case
			defer db.Close() // Ensure the database connection is closed after the test

			var userID int

			if !tc.wantErr {
				email := fmt.Sprintf("testuser_%s@example.com", tc.name) // Generate a unique email for testing
				var err error

				userID, err = insertUser(db, email, []byte("validpassword123")) // Insert a valid user into the database
				defer db.ExecContext(ctx, deleteUserQueryRow, userID)           // Clean up by deleting the user after the test

				assert.NoError(t, err) // Assert that there was no error inserting the user

				// Insert user pairs based on the exchange being tested
				if tc.exchange == "Binance" {
					insertUserPair(db, userID, "Binance", "BTC/USDT", 45000) // Insert first pair for Binance
					insertUserPair(db, userID, "Binance", "ETH/USDT", 3000)  // Insert second pair for Binance
				} else if tc.exchange == "Coinbase" {
					insertUserPair(db, userID, "Coinbase", "ETH/USDT", 3000) // Insert pair for Coinbase
				}
			}

			ctx, cancel := context.WithTimeout(context.Background(), contextTimeout) // Set up context with timeout
			defer cancel()

			repo := repository.NewUserPairsRepository(db)           // Create a new repository instance for user pairs
			pairs, err := repo.GetPairsByExchange(ctx, tc.exchange) // Attempt to retrieve all pairs for the specified exchange

			if tc.wantErr {
				assert.Error(t, err)    // Assert that an error occurred if one was expected
				assert.Len(t, pairs, 0) // Assert that no pairs were returned if an error was expected
				return                  // Exit early since we expect an error and no pairs should be returned
			}

			assert.NoError(t, err)                      // Assert that no error occurred for valid input
			assert.Len(t, pairs, len(tc.expectedPairs)) // Assert that the number of retrieved pairs matches the expected count

			for i, p := range pairs {
				assert.Equal(t, tc.expectedPairs[i], p) // Check that each retrieved pair matches the expected pair
			}
		})
	}
}
func TestDeletePair(t *testing.T) {
	// Run tests in parallel to improve execution speed
	t.Parallel()

	// Define test cases for deleting user pairs
	tests := []struct {
		name     string           // Name of the test case
		pairData models.UserPairs // Data for the user pair being tested
		wantErr  bool             // Expectation of whether an error should occur
	}{
		{
			name: "Valid Delete Pair",
			pairData: models.UserPairs{
				Exchange: "Binance",
				Pair:     "BTC/USDT", // Valid pair to be deleted
			},
			wantErr: false, // No error expected for valid input
		},
		{
			name: "Invalid Delete Pair",
			pairData: models.UserPairs{
				Exchange: "Binance",
				Pair:     "NON_EXISTENT_PAIR", // Assuming this pair does not exist
			},
			wantErr: true, // Error expected due to non-existent pair
		},
	}

	// Iterate through each test case
	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			db := setupDB()  // Setup a new database connection for each test case
			defer db.Close() // Ensure the database connection is closed after the test

			repo := repository.NewUserPairsRepository(db) // Create a new repository instance for user pairs

			// For valid delete case, create a user and insert a valid pair to delete
			if !tc.wantErr && tc.name == "Valid Delete Pair" {
				email := "validuser_" + tc.name + "@example.com" // Generate a unique email for testing

				userID, err := insertUser(db, email, []byte("validpassword123")) // Insert a valid user into the database
				defer db.ExecContext(ctx, deleteUserQueryRow, userID)            // Clean up by deleting the user after the test

				assert.NoError(t, err) // Assert that there was no error inserting the user

				// Insert a valid pair into the database for deletion later
				err = insertUserPair(db, userID, tc.pairData.Exchange, tc.pairData.Pair, 45000)
				assert.NoError(t, err) // Assert that there was no error inserting the pair

				tc.pairData.UserID = userID // Set the valid user ID in the pair data for this test case
			}

			err := repo.DeletePair(ctx, tc.pairData) // Attempt to delete the specified pair

			if tc.wantErr {
				assert.Error(t, err) // Assert that an error occurred if one was expected
			} else {
				assert.NoError(t, err) // Assert that no error occurred for valid input

				var deletedPair models.UserPairs
				query := `SELECT * FROM user_pairs WHERE user_id = $1 AND pair = $2`
				err = db.GetContext(ctx, &deletedPair, query, tc.pairData.UserID, tc.pairData.Pair) // Attempt to retrieve the deleted pair

				assert.Error(t, err) // Assert that an error occurred when trying to retrieve a deleted pair (should not exist)
			}
		})
	}
}
