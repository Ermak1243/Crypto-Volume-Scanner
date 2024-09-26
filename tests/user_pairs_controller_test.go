package tests

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"

	"main/api/server/controller"
	"main/internal/mocks"
	"main/internal/models"
	"main/internal/service/exchange"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddPairController(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	tests := []struct {
		name       string           // Name of the test case
		userID     int              // User ID for adding the pair
		pairData   models.UserPairs // Input data for adding the user pair
		mocksSetup func(
			userPairsMock *mocks.UserPairsService,
			userMock *mocks.UserService,
			allExchangesMock *mocks.AllExchanges,
			mockExchange *mocks.Exchange,
		) // Function to set up mock behavior
		expectedCode int // Expected HTTP status code after the request
	}{
		{
			name:   "Successful Addition",
			userID: 1,
			pairData: models.UserPairs{
				UserID:   1,
				Pair:     "BTC-ETH",
				Exchange: "Binance", // Assuming Exchange field is part of UserPairs
			},
			mocksSetup: func(
				userPairsMock *mocks.UserPairsService,
				userMock *mocks.UserService,
				allExchangesMock *mocks.AllExchanges,
				mockExchange *mocks.Exchange,
			) {
				userPairsMock.On("Add", mock.Anything, mock.Anything).Return(nil) // Mock successful addition
				userMock.On("SetUserIdIntoMemory", mock.Anything).Return(nil)     // Mock successful addition
				allExchangesMock.On("Get", "Binance").Return(mockExchange)        // Mock getting the exchange
				mockExchange.On("AddPairToSubscribedPairs", "BTC-ETH").Return()   // Mock adding pair to subscribed pairs
			},
			expectedCode: http.StatusOK, // Expecting 200 OK status
		},
		{
			name:   "Error Adding Pair - Service Error",
			userID: 1,
			pairData: models.UserPairs{
				UserID:   1,
				Pair:     "BTC-ETH",
				Exchange: "Binance", // Assuming Exchange field is part of UserPairs
			},
			mocksSetup: func(
				userPairsMock *mocks.UserPairsService,
				userMock *mocks.UserService,
				allExchangesMock *mocks.AllExchanges,
				mockExchange *mocks.Exchange,
			) {
				userPairsMock.On("Add", mock.Anything, mock.Anything).Return(errors.New("service error")) // Mock error during addition
			},
			expectedCode: http.StatusInternalServerError, // Expecting 500 Internal Server Error status due to service error
		},
	}

	for _, tt := range tests {
		tc := tt // Capture range variable for use in goroutine

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run each test case in parallel

			app := fiber.New() // Create a new Fiber application instance

			mockUserPairsService := mocks.NewUserPairsService(t) // Create a new mock UserPairs service
			mockUserService := mocks.NewUserService(t)           // Create a new mock User service
			mockAllExchangesStorage := mocks.NewAllExchanges(t)  // Create a new mock AllExchanges storage
			mockExchange := mocks.NewExchange(t)                 // Create a new mock Exchange instance

			tc.mocksSetup(mockUserPairsService, mockUserService, mockAllExchangesStorage, mockExchange) // Setup mocks for the current test case

			userPairsController := controller.NewUserPairsController(mockUserPairsService, mockUserService, nil, mockAllExchangesStorage)

			app.Post("/api/user/pairs", func(c *fiber.Ctx) error {
				c.Locals("user", models.User{ID: tc.userID}) // Add user to context locals
				return userPairsController.Add(c)            // Call Add method on UserPairsController
			})

			reqBody, _ := json.Marshal(tc.pairData)                                         // Marshal pairData into JSON format for request body
			req := httptest.NewRequest("POST", "/api/user/pairs", bytes.NewBuffer(reqBody)) // Create a new POST request with JSON body
			req.Header.Set("Content-Type", "application/json")                              // Set Content-Type header to application/json

			resp, err := app.Test(req, -1) // Execute the request against the Fiber app
			assert.NoError(t, err)         // Assert that there was no error during request execution

			assert.Equal(t, tc.expectedCode, resp.StatusCode) // Assert that the response status code matches expected
		})
	}
}

func TestUpdateExactValueController(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	tests := []struct {
		name         string                                 // Name of the test case
		userID       int                                    // User ID for updating the pair
		pairData     models.UserPairs                       // Input data for updating the user pair
		mocksSetup   func(userMock *mocks.UserPairsService) // Function to set up mock behavior
		expectedCode int                                    // Expected HTTP status code after the request
	}{
		{
			name:   "Successful Update",
			userID: 1,
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "BTC-ETH",
				ExactValue: 100, // Assuming there's a Value field to update
			},
			mocksSetup: func(userPairsMock *mocks.UserPairsService) {
				userPairsMock.On("UpdateExactValue", mock.Anything, mock.Anything).Return(nil) // Mock successful update
			},
			expectedCode: http.StatusOK, // Expecting 200 OK status
		},
		{
			name:   "Error Updating Pair",
			userID: 1,
			pairData: models.UserPairs{
				UserID:     1,
				Pair:       "BTC-ETH",
				ExactValue: 100,
			},
			mocksSetup: func(userPairsMock *mocks.UserPairsService) {
				userPairsMock.On("UpdateExactValue", mock.Anything, mock.Anything).Return(errors.New("update error")) // Mock error during update
			},
			expectedCode: http.StatusInternalServerError, // Expecting 500 Internal Server Error status due to update failure
		},
	}

	for _, tt := range tests {
		tc := tt // Capture range variable for use in goroutine

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run each test case in parallel

			app := fiber.New() // Create a new Fiber application instance

			mockUserPairsService := mocks.NewUserPairsService(t) // Create a new mock UserPairs service
			mockUserService := mocks.NewUserService(t)           // Create a new mock User service
			mockAllExchangesStorage := mocks.NewAllExchanges(t)  // Create a new mock AllExchanges storage

			tc.mocksSetup(mockUserPairsService) // Setup mocks for the current test case

			userPairsController := controller.NewUserPairsController(mockUserPairsService, mockUserService, nil, mockAllExchangesStorage)

			app.Put("/api/user/pairs", func(c *fiber.Ctx) error {
				c.Locals("user", models.User{ID: tc.userID})   // Add user to context locals
				return userPairsController.UpdateExactValue(c) // Call UpdateExactValue method on UserPairsController
			})

			reqBody, _ := json.Marshal(tc.pairData)                                        // Marshal pairData into JSON format for request body
			req := httptest.NewRequest("PUT", "/api/user/pairs", bytes.NewBuffer(reqBody)) // Create a new PUT request with JSON body
			req.Header.Set("Content-Type", "application/json")                             // Set Content-Type header to application/json

			resp, err := app.Test(req, -1) // Execute the request against the Fiber app
			assert.NoError(t, err)         // Assert that there was no error during request execution

			assert.Equal(t, tc.expectedCode, resp.StatusCode) // Assert that the response status code matches expected
		})
	}
}

func TestGetAllUserPairsController(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	tests := []struct {
		name         string                                 // Name of the test case
		userID       int                                    // User ID for which to retrieve pairs
		mocksSetup   func(userMock *mocks.UserPairsService) // Function to set up mock behavior
		expectedCode int                                    // Expected HTTP status code after the request
	}{
		{
			name:   "Successful Retrieval",
			userID: 1,
			mocksSetup: func(userPairsMock *mocks.UserPairsService) {
				userPairsMock.On("GetAllUserPairs", mock.Anything, 1).Return([]models.UserPairs{
					{UserID: 1, Pair: "BTC-ETH"},
					{UserID: 1, Pair: "ETH-LTC"},
				}, nil) // Mock successful retrieval of user pairs
			},
			expectedCode: http.StatusOK, // Expecting 200 OK status
		},
		{
			name:   "Error Retrieving User Pairs",
			userID: 1,
			mocksSetup: func(userPairsMock *mocks.UserPairsService) {
				userPairsMock.On("GetAllUserPairs", mock.Anything, 1).Return(nil, errors.New("retrieve error")) // Mock error during retrieval
			},
			expectedCode: http.StatusInternalServerError, // Expecting 500 Internal Server Error status due to retrieval failure
		},
	}

	for _, tt := range tests {
		tc := tt // Capture range variable for use in goroutine

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run each test case in parallel

			app := fiber.New() // Create a new Fiber application instance

			mockUserPairsService := mocks.NewUserPairsService(t) // Create a new mock UserPairs service
			mockUserService := mocks.NewUserService(t)           // Create a new mock User service
			mockAllExchangesStorage := mocks.NewAllExchanges(t)  // Create a new mock AllExchanges storage

			tc.mocksSetup(mockUserPairsService) // Setup mocks for the current test case

			userPairsController := controller.NewUserPairsController(mockUserPairsService, mockUserService, nil, mockAllExchangesStorage)

			app.Get("/api/user/pairs", func(c *fiber.Ctx) error {
				c.Locals("user", models.User{ID: tc.userID})  // Add user to context locals
				return userPairsController.GetAllUserPairs(c) // Call GetAllUserPairs method on UserPairsController
			})

			req := httptest.NewRequest("GET", "/api/user/pairs", nil) // Create a new GET request

			resp, err := app.Test(req, -1) // Execute the request against the Fiber app
			assert.NoError(t, err)         // Assert that there was no error during request execution

			assert.Equal(t, tc.expectedCode, resp.StatusCode) // Assert that the response status code matches expected
		})
	}
}

func TestDeletePairController(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	tests := []struct {
		name       string // Name of the test case
		userID     int    // User ID for deleting the pair
		pairQuery  string // Query parameter for deleting the user pair
		mocksSetup func(
			userPairsMock *mocks.UserPairsService,
			userMock *mocks.UserService,
			allExchangesMock *mocks.AllExchanges,
			mockExchange *mocks.Exchange,
			mockFoundVolumes *mocks.FoundVolumesService,
		) // Function to set up mock behavior
		expectedCode int // Expected HTTP status code after the request
	}{
		{
			name:      "Successful Deletion",
			userID:    1,
			pairQuery: "BTC-ETH", // Pair to be deleted
			mocksSetup: func(
				userPairsMock *mocks.UserPairsService,
				userMock *mocks.UserService,
				allExchangesMock *mocks.AllExchanges,
				mockExchange *mocks.Exchange,
				mockFoundVolumes *mocks.FoundVolumesService,
			) {
				mockExchange.On("DeletePairFromSubscribedPairs", "BTC-ETH").Return()
				mockExchange.On("ExchangeName").Return("test-exchange")
				userPairsMock.On("DeletePair", mock.Anything, mock.Anything).Return(nil) // Mock successful deletion
				userMock.On("DeleteUserIdFromMemory", mock.Anything).Return(nil)         // Mock successful deletion
				allExchangesMock.On("All").Return([]exchange.Exchange{mockExchange})
				mockFoundVolumes.On("DeleteFoundVolume", mock.Anything).Return()
			},
			expectedCode: http.StatusOK, // Expecting 200 OK status
		},
		{
			name:      "Error Deleting Pair",
			userID:    1,
			pairQuery: "BTC-ETH", // Pair to be deleted
			mocksSetup: func(
				userPairsMock *mocks.UserPairsService,
				userMock *mocks.UserService,
				allExchangesMock *mocks.AllExchanges,
				mockExchange *mocks.Exchange,
				mockFoundVolumes *mocks.FoundVolumesService,
			) {
				userPairsMock.On("DeletePair", mock.Anything, mock.Anything).Return(errors.New("delete error")) // Mock error during deletion
			},
			expectedCode: http.StatusInternalServerError, // Expecting 500 Internal Server Error status due to deletion failure
		},
	}

	for _, tt := range tests {
		tc := tt // Capture range variable for use in goroutine

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run each test case in parallel

			app := fiber.New() // Create a new Fiber application instance

			mockUserPairsService := mocks.NewUserPairsService(t)       // Create a new mock UserPairs service
			mockUserService := mocks.NewUserService(t)                 // Create a new mock User service
			mockAllExchangesStorage := mocks.NewAllExchanges(t)        // Create a new mock AllExchanges storage
			mockFoundVolumesService := mocks.NewFoundVolumesService(t) // Create a new mock FoundVolumes service
			mockExchange := mocks.NewExchange(t)                       // Create a new mock Exchange instance

			tc.mocksSetup(
				mockUserPairsService,
				mockUserService,
				mockAllExchangesStorage,
				mockExchange,
				mockFoundVolumesService) // Setup mocks for the current test case

			userPairsController := controller.NewUserPairsController(mockUserPairsService, mockUserService, mockFoundVolumesService, mockAllExchangesStorage)

			app.Delete("/api/user/pairs", func(c *fiber.Ctx) error {
				c.Locals("user", models.User{ID: tc.userID}) // Add user to context locals
				return userPairsController.DeletePair(c)     // Call DeletePair method on UserPairsController
			})

			req := httptest.NewRequest("DELETE", "/api/user/pairs?pair="+tc.pairQuery, nil) // Create a new DELETE request with query parameter
			req.Header.Set("Content-Type", "application/json")                              // Set Content-Type header to application/json

			resp, err := app.Test(req, -1) // Execute the request against the Fiber app
			assert.NoError(t, err)         // Assert that there was no error during request execution

			assert.Equal(t, tc.expectedCode, resp.StatusCode) // Assert that the response status code matches expected
		})
	}
}
