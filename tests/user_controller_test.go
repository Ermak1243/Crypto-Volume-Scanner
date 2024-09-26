package tests

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"main/api/server/controller"
	"main/internal/mocks"
	"main/internal/models"
	"main/internal/service/exchange"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test for Signup method
func TestSignup(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	// Define a slice of test cases for the Signup functionality
	tests := []struct {
		name         string                                                       // Name of the test case
		newUserData  models.UserAuth                                              // New user data for signup
		mocksSetup   func(userMock *mocks.UserService, jwtMock *mocks.JwtService) // Function to set up mock behavior
		expectedCode int                                                          // Expected HTTP status code after the request
	}{
		{
			name: "Successful Signup",
			newUserData: models.UserAuth{
				Email:    "test@example.com",
				Password: "password123",
			},
			mocksSetup: func(userMock *mocks.UserService, jwtMock *mocks.JwtService) {
				userMock.On("InsertUser", mock.Anything, mock.Anything).Return(1, nil)                    // Mock successful user insertion
				userMock.On("UpdateRefreshToken", mock.Anything, mock.Anything).Return(nil)               // Mock successful refresh token update
				jwtMock.On("CreateAccessToken", 1, mock.Anything).Return("accessToken", int64(3600), nil) // Mock access token creation
				jwtMock.On("CreateRefreshToken", 1, mock.Anything).Return("refreshToken", nil)            // Mock refresh token creation
			},
			expectedCode: http.StatusOK, // Expecting 200 OK status
		},
		{
			name: "Invalid Input Data",
			newUserData: models.UserAuth{
				Email:    "",
				Password: "",
			},
			mocksSetup:   func(userMock *mocks.UserService, jwtMock *mocks.JwtService) {}, // No mocks needed for this case
			expectedCode: http.StatusBadRequest,                                           // Expecting 400 Bad Request status due to invalid input
		},
		{
			name: "Error Inserting User",
			newUserData: models.UserAuth{
				Email:    "test@example.com",
				Password: "password123",
			},
			mocksSetup: func(userMock *mocks.UserService, jwtMock *mocks.JwtService) {
				userMock.On("InsertUser", mock.Anything, mock.Anything).Return(0, errors.New("insert error")) // Mock error during user insertion
			},
			expectedCode: http.StatusInternalServerError, // Expecting 500 Internal Server Error status due to insertion failure
		},
	}

	for _, tt := range tests {
		tc := tt // Capture range variable for use in goroutine

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run each test case in parallel

			app := fiber.New() // Create a new Fiber application instance

			mockUserService := mocks.NewUserService(t)          // Create a new mock user service
			mockJwtService := mocks.NewJwtService(t)            // Create a new mock JWT service
			mockAllExchangesStorage := mocks.NewAllExchanges(t) // Create a new mock all exchanges storage

			if tc.mocksSetup != nil {
				tc.mocksSetup(mockUserService, mockJwtService) // Setup mocks for the current test case
			}

			uc := controller.NewUserController(mockUserService, mockJwtService, mockAllExchangesStorage) // Create a new UserController instance
			app.Post("/api/user/auth/signup", uc.Signup)                                                 // Define POST route for signup

			reqBody := `{"email":"` + tc.newUserData.Email + `","password":"` + tc.newUserData.Password + `"}`
			req := httptest.NewRequest("POST", "/api/user/auth/signup", strings.NewReader(reqBody)) // Create a new POST request with JSON body
			req.Header.Set("Content-Type", "application/json")                                      // Set Content-Type header to application/json

			resp, err := app.Test(req, -1) // Execute the request against the Fiber app

			assert.NoError(t, err)                            // Assert that there was no error during request execution
			assert.Equal(t, tc.expectedCode, resp.StatusCode) // Assert that the response status code matches expected
		})
	}
}

// Test for Tokens method
func TestTokens(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	// Define a slice of test cases for the Tokens functionality
	tests := []struct {
		name         string                                                       // Name of the test case
		userID       int                                                          // User ID for token generation
		refreshToken string                                                       // Refresh token for authentication
		mocksSetup   func(userMock *mocks.UserService, jwtMock *mocks.JwtService) // Function to set up mock behavior
		expectedCode int                                                          // Expected HTTP status code after the request
	}{
		{
			name:         "Successful Token Retrieval",
			userID:       1,
			refreshToken: "valid_refresh_token", // Valid refresh token for authentication
			mocksSetup: func(userMock *mocks.UserService, jwtMock *mocks.JwtService) {
				// Mock successful access and refresh token creation
				userMock.On("UpdateRefreshToken", mock.Anything, mock.Anything).Return(nil)
				jwtMock.On("CreateAccessToken", 1, mock.Anything).Return("newAccessToken", int64(3600), nil)
				jwtMock.On("CreateRefreshToken", 1, mock.Anything).Return("newRefreshToken", nil)
			},
			expectedCode: http.StatusOK, // Expecting 200 OK status
		},
		{
			name:         "Invalid Refresh Token",
			userID:       1,
			refreshToken: "",                      // Invalid refresh token (empty)
			expectedCode: http.StatusUnauthorized, // Expecting 401 Unauthorized status due to invalid refresh token
		},
		{
			name:         "Error Retrieving User",
			userID:       1,
			refreshToken: "valid_refresh_token",
			mocksSetup: func(userMock *mocks.UserService, jwtMock *mocks.JwtService) {
				// Mock error during access token creation
				jwtMock.On("CreateAccessToken", mock.Anything, mock.Anything).Return("", int64(0), errors.New("token creation error"))
			},
			expectedCode: http.StatusInternalServerError, // Expecting 500 Internal Server Error status due to retrieval failure
		},
	}

	for _, tt := range tests {
		tc := tt // Capture range variable for use in goroutine

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run each test case in parallel

			app := fiber.New() // Create a new Fiber application instance

			mockUserService := mocks.NewUserService(t)          // Create a new mock user service
			mockJwtService := mocks.NewJwtService(t)            // Create a new mock JWT service
			mockAllExchangesStorage := mocks.NewAllExchanges(t) // Create a new mock all exchanges storage

			if tc.mocksSetup != nil {
				tc.mocksSetup(mockUserService, mockJwtService) // Setup mocks for the current test case
			}

			userController := controller.NewUserController(mockUserService, mockJwtService, mockAllExchangesStorage) // Create a new UserController instance

			app.Get("/api/user/auth/tokens", func(c *fiber.Ctx) error {
				user := models.User{ID: tc.userID}    // Create a user model with the specified user ID
				user.SetRefreshToken(tc.refreshToken) // Set the refresh token for the user

				if tc.mocksSetup == nil {
					user.SetRefreshToken("1")
				}

				c.Locals("user", user)                                   // Store the user in context locals for retrieval in controller
				c.Request().Header.Set("Authorization", tc.refreshToken) // Set the Authorization header with the refresh token
				return userController.Tokens(c)                          // Call Tokens method on UserController
			})

			req := httptest.NewRequest("GET", "/api/user/auth/tokens", nil) // Create a new GET request

			resp, err := app.Test(req, -1) // Execute the request against the Fiber app
			assert.NoError(t, err)         // Assert that there was no error during request execution

			assert.Equal(t, tc.expectedCode, resp.StatusCode) // Assert that the response status code matches expected
		})
	}
}

func TestLogin(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	tests := []struct {
		name         string                                                       // Name of the test case
		userData     models.UserAuth                                              // User authentication data for login
		mocksSetup   func(userMock *mocks.UserService, jwtMock *mocks.JwtService) // Function to set up mock behavior
		expectedCode int                                                          // Expected HTTP status code after the request
	}{
		{
			name: "Successful Login",
			userData: models.UserAuth{
				Email:    "test@example.com",
				Password: "password123",
			},
			mocksSetup: func(userMock *mocks.UserService, jwtMock *mocks.JwtService) {
				user := models.User{ID: 1, Email: "test@example.com"}
				user.SetPassword("password123")                                                                 // Assume this sets a hashed password correctly
				userMock.On("UpdateRefreshToken", mock.Anything, mock.Anything).Return(nil)                     // Mock successful user retrieval
				userMock.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)              // Mock successful user retrieval
				jwtMock.On("CreateAccessToken", user.ID, mock.Anything).Return("accessToken", int64(3600), nil) // Mock access token creation
				jwtMock.On("CreateRefreshToken", user.ID, mock.Anything).Return("refreshToken", nil)            // Mock refresh token creation
			},
			expectedCode: http.StatusOK, // Expecting 200 OK status
		},
		{
			name: "User Not Found",
			userData: models.UserAuth{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			mocksSetup: func(userMock *mocks.UserService, jwtMock *mocks.JwtService) {
				userMock.On("GetUserByEmail", mock.Anything, "notfound@example.com").Return(models.User{}, errors.New("user not found")) // Mock user not found error
			},
			expectedCode: http.StatusBadRequest, // Expecting 400 Bad Request status due to user not found
		},
		{
			name: "Invalid Password",
			userData: models.UserAuth{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mocksSetup: func(userMock *mocks.UserService, jwtMock *mocks.JwtService) {
				user := models.User{ID: 1, Email: "test@example.com"}
				user.SetPassword("password123")
				userMock.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil) // Mock successful user retrieval
			},
			expectedCode: http.StatusBadRequest, // Expecting 400 Bad Request status due to invalid password
		},
		{
			name: "Error Generating Tokens",
			userData: models.UserAuth{
				Email:    "test@example.com",
				Password: "password123",
			},
			mocksSetup: func(userMock *mocks.UserService, jwtMock *mocks.JwtService) {
				user := models.User{ID: -1, Email: "test@example.com"}
				user.SetPassword("password123")
				userMock.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)
				jwtMock.On("CreateAccessToken", user.ID, mock.Anything).Return("", int64(0), errors.New("token error")) // Mock token generation error
			},
			expectedCode: http.StatusInternalServerError, // Expecting 500 Internal Server Error status due to token generation failure
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := fiber.New()

			mockUserService := mocks.NewUserService(t)
			mockJwtService := mocks.NewJwtService(t)
			mockAllExchangesStorage := mocks.NewAllExchanges(t) // Create a new mock all exchanges storage

			if tc.mocksSetup != nil {
				tc.mocksSetup(mockUserService, mockJwtService) // Setup mocks for the current test case
			}

			userController := controller.NewUserController(mockUserService, mockJwtService, mockAllExchangesStorage) // Create a new UserController instance
			app.Post("/api/user/auth/login", userController.Login)

			reqBody := `{"email":"` + tc.userData.Email + `","password":"` + tc.userData.Password + `"}`
			req := httptest.NewRequest("POST", "/api/user/auth/login", bytes.NewBufferString(reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode)

			var responseBody = make(map[string]interface{})
			json.NewDecoder(resp.Body).Decode(&responseBody)
		})
	}
}

func TestUpdatePasswordController(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	tests := []struct {
		name         string                                                          // Name of the test case
		userID       int                                                             // User ID for updating password
		oldPassword  []byte                                                          // Old password for validation
		newPassword  []byte                                                          // New password to be set
		mocksSetup   func(userMock *mocks.UserService, jwtService *mocks.JwtService) // Function to set up mock behavior
		expectedCode int                                                             // Expected HTTP status code after the request
	}{
		{
			name:        "Successful Password Update",
			userID:      1,
			oldPassword: []byte("oldpassword123"),
			newPassword: []byte("newpassword123"),
			mocksSetup: func(userMock *mocks.UserService, jwtMock *mocks.JwtService) {
				jwtMock.On("CreateAccessToken", mock.Anything, mock.Anything).Return("", int64(3600), nil) // Mock access token creation
				jwtMock.On("CreateRefreshToken", mock.Anything, mock.Anything).Return("", nil)             // Mock refresh token creation
				userMock.On("UpdatePassword", mock.Anything, mock.Anything).Return(nil)                    // Mock successful password update
			},
			expectedCode: http.StatusOK, // Expecting 200 OK status
		},
		{
			name:         "Invalid Old Password",
			userID:       1,
			oldPassword:  []byte("wrongpassword"),
			newPassword:  []byte("newpassword123"),
			expectedCode: http.StatusBadRequest, // Expecting 400 Bad Request status due to invalid old password
		},
		{
			name:        "Error Updating Password",
			userID:      1,
			oldPassword: []byte("oldpassword123"),
			newPassword: []byte("newpassword123"),
			mocksSetup: func(userMock *mocks.UserService, jwtMock *mocks.JwtService) {
				jwtMock.On("CreateAccessToken", mock.Anything, mock.Anything).Return("", int64(3600), nil)     // Mock access token creation
				jwtMock.On("CreateRefreshToken", mock.Anything, mock.Anything).Return("", nil)                 // Mock refresh token creation
				userMock.On("UpdatePassword", mock.Anything, mock.Anything).Return(errors.New("update error")) // Mock error during password update
			},
			expectedCode: http.StatusInternalServerError, // Expecting 500 Internal Server Error status due to update failure
		},
	}

	for _, tt := range tests {
		tc := tt // Capture range variable for use in goroutine

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run each test case in parallel

			app := fiber.New() // Create a new Fiber application instance

			mockUserService := mocks.NewUserService(t) // Create a new mock User service
			mockJwtService := mocks.NewJwtService(t)
			mockAllExchangesStorage := mocks.NewAllExchanges(t) // Create a new mock all exchanges storage

			if tc.mocksSetup != nil {
				tc.mocksSetup(mockUserService, mockJwtService) // Setup mocks for the current test case
			}

			userController := controller.NewUserController(mockUserService, mockJwtService, mockAllExchangesStorage) // Create a new UserController instance

			app.Put("/api/user/auth/update-password", func(c *fiber.Ctx) error {
				user := models.User{ID: tc.userID}
				user.SetPassword(string(tc.oldPassword))

				if tc.mocksSetup == nil {
					user.SetPassword("")
				}

				c.Locals("user", user) // Add user to context locals

				return userController.UpdatePassword(c) // Call UpdatePassword method on UserController
			})

			reqBody := models.PasswordUpdate{
				OldPassword:       string(tc.oldPassword),
				NewPassword:       string(tc.newPassword),
				NewPasswordRepeat: string(tc.newPassword), // Assuming you want to check if they match in your logic
			}
			body, _ := json.Marshal(reqBody) // Marshal request body into JSON format

			req := httptest.NewRequest("PUT", "/api/user/auth/update-password", bytes.NewBuffer(body)) // Create a new PUT request with JSON body
			req.Header.Set("Content-Type", "application/json")                                         // Set Content-Type header to application/json

			resp, err := app.Test(req, -1) // Execute the request against the Fiber app
			assert.NoError(t, err)         // Assert that there was no error during request execution

			assert.Equal(t, tc.expectedCode, resp.StatusCode) // Assert that the response status code matches expected
		})
	}
}

func TestDeleteUserController(t *testing.T) {
	// Define a slice of test cases for the DeleteUserController.
	tests := []struct {
		name         string                                                                                                // Name of the test case
		mocksSetup   func(userMock *mocks.UserService, allExchangesMock *mocks.AllExchanges, exchangeMock *mocks.Exchange) // Function to set up mock behavior
		expectedCode int                                                                                                   // Expected HTTP status code after the request
		expectedBody string                                                                                                // Expected response body in JSON format
	}{
		{
			name: "Successful User Deletion",
			mocksSetup: func(userMock *mocks.UserService, allExchangesMock *mocks.AllExchanges, exchangeMock *mocks.Exchange) {
				// Setup mock to return no error when DeleteUser is called.
				allExchangesMock.On("All").Return([]exchange.Exchange{exchangeMock})
				exchangeMock.On("ClearSubscribedPairsStorage").Return()
				userMock.On("DeleteUser", mock.Anything, 1).Return(nil)
				userMock.On("DeleteUserIdFromMemory", mock.Anything).Return(nil)
			},
			expectedCode: http.StatusOK,                            // Expecting 200 OK status
			expectedBody: `{"result":"user deleted successfully"}`, // Expected response body
		},
		{
			name: "Error Deleting User",
			mocksSetup: func(userMock *mocks.UserService, allExchangesMock *mocks.AllExchanges, exchangeMock *mocks.Exchange) {
				// Setup mock to return an error when DeleteUser is called.
				userMock.On("DeleteUser", mock.Anything, 1).Return(errors.New("deletion error"))
			},
			expectedCode: http.StatusInternalServerError,      // Expecting 500 Internal Server Error status
			expectedBody: `{"result":"user deletion failed"}`, // Expected response body
		},
	}

	// Iterate through each test case defined above.
	for _, tt := range tests {
		tc := tt // Capture range variable to avoid closure issues

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run each test case in parallel for efficiency

			app := fiber.New() // Create a new Fiber application instance

			mockUserService := mocks.NewUserService(t)          // Create a new mock user service
			mockAllExchangesStorage := mocks.NewAllExchanges(t) // Create a new mock all exchanges storage
			mockExchange := mocks.NewExchange(t)                // Create a new mock exchange
			if tc.mocksSetup != nil {
				tc.mocksSetup(mockUserService, mockAllExchangesStorage, mockExchange) // Setup mocks for the current test case
			}

			userController := controller.NewUserController(mockUserService, nil, mockAllExchangesStorage) // Create a new UserController instance
			app.Delete("/api/user", func(c *fiber.Ctx) error {
				user := models.User{ID: 1}         // Create a user model with ID 1
				user.SetPassword("oldpassword123") // Set a dummy password (not used in this test)

				c.Locals("user", user)              // Store the user in context locals for retrieval in controller
				return userController.DeleteUser(c) // Call the DeleteUser method on the controller
			})

			req := httptest.NewRequest("DELETE", "/api/user", nil) // Create a new DELETE request

			resp, err := app.Test(req, -1) // Execute the request against the Fiber app

			assert.NoError(t, err)                            // Assert that there was no error during request execution
			assert.Equal(t, tc.expectedCode, resp.StatusCode) // Assert that the response status code matches expected

			if tc.expectedBody != "" {
				bodyBytes, _ := io.ReadAll(resp.Body)                // Read the response body into bytes
				assert.JSONEq(t, tc.expectedBody, string(bodyBytes)) // Assert that the JSON response matches expected body
			}
		})
	}
}
