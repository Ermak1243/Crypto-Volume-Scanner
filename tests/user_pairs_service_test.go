package tests

import (
	"context"
	"errors"
	"main/internal/mocks"
	"main/internal/models"
	"main/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserPairsService_Add(t *testing.T) {
	// Run tests in parallel to improve execution speed
	t.Parallel()

	// Define test cases for adding user pairs
	tests := []struct {
		name      string                           // Name of the test case
		pairData  models.UserPairs                 // Data for the user pair being tested
		mockRepo  func(*mocks.UserPairsRepository) // Mocking the repository behavior
		expectErr bool                             // Expectation of whether an error should occur
	}{
		{
			name: "Valid pair data",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "BTC/USD",
				Exchange:   "binance_spot",
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("Add", mock.Anything, mock.Anything).Return(nil) // Expect Add to be called with any arguments and return no error
			},
			expectErr: false, // No error expected for valid input
		},
		{
			name: "Empty pair name",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "", // Invalid data
				Exchange:   "binance_spot",
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("Add", mock.Anything, mock.Anything).Return(nil).Maybe() // Allow for Add to be called but expect it not to be in this case
			},
			expectErr: true, // Error expected due to empty pair name
		},
		{
			name: "Empty exchange name",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "BTC/USD",
				Exchange:   "", // Invalid data
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("Add", mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			expectErr: true, // Error expected due to empty exchange name
		},
		{
			name: "Exact value below one",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "BTC/USD",
				Exchange:   "binance_spot",
				ExactValue: 0, // Invalid data
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("Add", mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			expectErr: true, // Error expected due to exact value being below one
		},
		{
			name: "User ID below one",
			pairData: models.UserPairs{
				UserID:     0, // Invalid data
				Pair:       "BTC/USD",
				Exchange:   "binance_spot",
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("Add", mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			expectErr: true, // Error expected due to invalid user ID
		},
		{
			name: "Invalid pair format",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "", // Invalid format
				Exchange:   "binance_spot",
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("Add", mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			expectErr: true, // Error expected due to invalid pair format (empty)
		},
		{
			name: "Invalid exchange format",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "BTC/USD",
				Exchange:   "INVALID_EXCHANGE", // Invalid format
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("Add", mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			expectErr: true, // Error expected due to invalid exchange format
		},
	}

	// Iterate through each test case
	for _, tc := range tests {
		tc := tc // Capture the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			mockRepo := mocks.NewUserPairsRepository(t)                               // Create a new instance of the mocked repository
			userPairsService := service.NewUserPairsService(mockRepo, contextTimeout) // Create a new instance of the service with the mocked repository

			// Set up the mock expectations based on the test case
			tc.mockRepo(mockRepo)

			// Call the Add method on the service with the provided pair data
			err := userPairsService.Add(context.Background(), tc.pairData)

			if tc.expectErr {
				assert.Error(t, err) // Assert that an error occurred if one was expected
			} else {
				assert.NoError(t, err) // Assert that no error occurred for valid input
			}

			mockRepo.AssertExpectations(t) // Verify that all expectations were met on the mocked repository
		})
	}
}

func TestUserPairsService_UpdateExactValue(t *testing.T) {
	// Run tests in parallel to improve execution speed
	t.Parallel()

	// Define test cases for updating the exact value of user pairs
	tests := []struct {
		name      string                           // Name of the test case
		pairData  models.UserPairs                 // Data for the user pair being tested
		mockRepo  func(*mocks.UserPairsRepository) // Mocking the repository behavior
		expectErr bool                             // Expectation of whether an error should occur
	}{
		{
			name: "Valid pair data",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "BTC/USD",
				Exchange:   "binance_spot",
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("UpdateExactValue", mock.Anything, mock.Anything).Return(nil) // Expect UpdateExactValue to be called and return no error
			},
			expectErr: false, // No error expected for valid input
		},
		{
			name: "Empty pair name",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "", // Invalid data (empty pair name)
				Exchange:   "binance_spot",
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("UpdateExactValue", mock.Anything, mock.Anything).Return(nil).Maybe() // Allow for UpdateExactValue to be called but expect it not to be in this case
			},
			expectErr: true, // Error expected due to empty pair name
		},
		{
			name: "Empty exchange name",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "BTC/USD",
				Exchange:   "", // Invalid data (empty exchange name)
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("UpdateExactValue", mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			expectErr: true, // Error expected due to empty exchange name
		},
		{
			name: "Exact value below one",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "BTC/USD",
				Exchange:   "binance_spot",
				ExactValue: 0, // Invalid data (exact value must be greater than zero)
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("UpdateExactValue", mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			expectErr: true, // Error expected due to exact value being below one
		},
		{
			name: "User ID below one",
			pairData: models.UserPairs{
				UserID:     0, // Invalid data (user ID must be greater than zero)
				Pair:       "BTC/USD",
				Exchange:   "binance_spot",
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("UpdateExactValue", mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			expectErr: true, // Error expected due to invalid user ID
		},
		{
			name: "Invalid pair format",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "", // Invalid format (empty pair)
				Exchange:   "binance_spot",
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("UpdateExactValue", mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			expectErr: true, // Error expected due to invalid pair format (empty)
		},
		{
			name: "Invalid exchange format",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "BTC/USD",
				Exchange:   "INVALID_EXCHANGE", // Invalid format (not a recognized exchange)
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("UpdateExactValue", mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			expectErr: true, // Error expected due to invalid exchange format
		},
	}

	// Iterate through each test case
	for _, tc := range tests {
		tc := tc // Capture the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			mockRepo := mocks.NewUserPairsRepository(t)                               // Create a new instance of the mocked repository
			userPairsService := service.NewUserPairsService(mockRepo, contextTimeout) // Create a new instance of the service with the mocked repository

			// Set up the mock expectations based on the test case
			tc.mockRepo(mockRepo)

			// Call the UpdateExactValue method on the service with the provided pair data
			err := userPairsService.UpdateExactValue(context.Background(), tc.pairData)

			if tc.expectErr {
				assert.Error(t, err) // Assert that an error occurred if one was expected
			} else {
				assert.NoError(t, err) // Assert that no error occurred for valid input
			}

			mockRepo.AssertExpectations(t) // Verify that all expectations were met on the mocked repository
		})
	}
}

func TestUserPairsService_DeletePair(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Define test cases for deleting user pairs
	tests := []struct {
		name      string                           // Name of the test case
		pairData  models.UserPairs                 // Data for the user pair being tested
		mockRepo  func(*mocks.UserPairsRepository) // Mocking the repository behavior
		expectErr bool                             // Expectation of whether an error should occur
	}{
		{
			name: "Valid pair data",
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "BTC/USD",
				Exchange:   "binance_spot",
				ExactValue: 100,
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("DeletePair", mock.Anything, mock.Anything).Return(nil) // Expect DeletePair to be called and return no error
			},
			expectErr: false, // No error expected for valid input
		},
		{
			name: "Invalid pair data",
			pairData: models.UserPairs{
				UserID: 0, // Invalid data (user ID must be greater than zero)
				Pair:   "BTC/USD",
			},
			mockRepo: func(m *mocks.UserPairsRepository) {
				m.On("DeletePair", mock.Anything, mock.Anything).Return(nil).Maybe() // Allow for DeletePair to be called but expect it not to be in this case
			},
			expectErr: true, // Error expected due to invalid user ID
		},
	}

	// Iterate through each test case
	for _, tc := range tests {
		tc := tc // Capture the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			mockRepo := mocks.NewUserPairsRepository(t)                               // Create a new instance of the mocked repository
			userPairsService := service.NewUserPairsService(mockRepo, contextTimeout) // Create a new instance of the service with the mocked repository

			// Set up the mock expectations based on the test case
			tc.mockRepo(mockRepo)

			// Call the DeletePair method on the service with the provided pair data
			err := userPairsService.DeletePair(context.Background(), tc.pairData)

			if tc.expectErr {
				assert.Error(t, err) // Assert that an error occurred if one was expected
			} else {
				assert.NoError(t, err) // Assert that no error occurred for valid input
			}

			mockRepo.AssertExpectations(t) // Verify that all expectations were met on the mocked repository
		})
	}
}

func TestUserPairsService_GetAllUserPairs(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Define test cases for retrieving all user pairs
	tests := []struct {
		name       string             // Name of the test case
		mockReturn []models.UserPairs // Mocked return value for the repository method
		mockErr    error              // Mocked error to simulate repository behavior
		expectErr  bool               // Expectation of whether an error should occur
	}{
		{
			name: "Successful retrieval of user pairs",
			mockReturn: []models.UserPairs{
				{UserID: 1, Pair: "BTC/USD"},
				{UserID: 1, Pair: "ETH/USD"},
			},
			mockErr:   nil, // No error expected for successful retrieval
			expectErr: false,
		},
		{
			name:       "Error retrieving user pairs",
			mockReturn: nil,                          // No pairs returned
			mockErr:    errors.New("database error"), // Simulate a database error
			expectErr:  true,                         // Error expected in this case
		},
	}

	// Iterate through each test case
	for _, tc := range tests {
		tc := tc // Capture the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			mockRepo := mocks.NewUserPairsRepository(t)                               // Create a new instance of the mocked repository
			userPairsService := service.NewUserPairsService(mockRepo, contextTimeout) // Create a new instance of the service with the mocked repository

			// Set up the mock expectations based on the test case
			mockRepo.On("GetAllUserPairs", mock.Anything, mock.Anything).Return(tc.mockReturn, tc.mockErr)

			// Call the GetAllUserPairs method on the service with a valid user ID
			pairs, err := userPairsService.GetAllUserPairs(context.Background(), 1)

			if tc.expectErr {
				assert.Error(t, err) // Assert that an error occurred if one was expected
				assert.Nil(t, pairs) // Assert that no pairs were returned in case of an error
			} else {
				assert.NoError(t, err)                // Assert that no error occurred for valid input
				assert.Equal(t, tc.mockReturn, pairs) // Assert that the returned pairs match the mocked return value
			}

			mockRepo.AssertExpectations(t) // Verify that all expectations were met on the mocked repository
		})
	}
}

func TestUserPairsService_GetPairsByExchange(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Define test cases for retrieving pairs by exchange
	tests := []struct {
		name       string   // Name of the test case
		exchange   string   // Exchange name for which pairs are to be fetched
		mockReturn []string // Mocked return value for the repository method
		mockErr    error    // Mocked error to simulate repository behavior
		expectErr  bool     // Expectation of whether an error should occur
	}{
		{
			name:       "Successful retrieval from Binance",
			exchange:   "Binance",                        // Valid exchange with expected pairs
			mockReturn: []string{"BTC/USDT", "ETH/USDT"}, // Expecting these pairs for Binance
			mockErr:    nil,                              // No error expected for successful retrieval
			expectErr:  false,
		},
		{
			name:       "Successful retrieval from Coinbase",
			exchange:   "Coinbase",           // Another valid exchange with expected pairs
			mockReturn: []string{"ETH/USDT"}, // Expecting this pair for Coinbase
			mockErr:    nil,                  // No error expected for successful retrieval
			expectErr:  false,
		},
		{
			name:       "Error retrieving pairs from NonExistent exchange",
			exchange:   "NonExistent",                // Exchange that does not exist in the database
			mockReturn: nil,                          // No pairs returned
			mockErr:    errors.New("database error"), // Simulate a database error
			expectErr:  true,                         // Error expected in this case
		},
	}

	// Iterate through each test case
	for _, tc := range tests {
		tc := tc // Capture the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			mockRepo := mocks.NewUserPairsRepository(t)                               // Create a new instance of the mocked repository
			userPairsService := service.NewUserPairsService(mockRepo, contextTimeout) // Create a new instance of the service with the mocked repository

			// Set up the mock expectations based on the test case
			mockRepo.On("GetPairsByExchange", mock.Anything, tc.exchange).Return(tc.mockReturn, tc.mockErr)

			ctx, cancel := context.WithTimeout(context.Background(), contextTimeout) // Set up context with timeout
			defer cancel()

			// Call the GetPairsByExchange method on the service with the specified exchange
			pairs, err := userPairsService.GetPairsByExchange(ctx, tc.exchange)

			if tc.expectErr {
				assert.Error(t, err) // Assert that an error occurred if one was expected
				assert.Nil(t, pairs) // Assert that no pairs were returned in case of an error
			} else {
				assert.NoError(t, err)                // Assert that no error occurred for valid input
				assert.Equal(t, tc.mockReturn, pairs) // Assert that the returned pairs match the mocked return value
			}

			mockRepo.AssertExpectations(t) // Verify that all expectations were met on the mocked repository
		})
	}
}
