package tests

import (
	"context"
	"errors"
	"main/internal/mocks"
	"main/internal/models"
	"main/internal/service"
	"strconv"
	"testing"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_InsertUser(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Create a mock user repository for testing
	mockUserRepository := mocks.NewUserRepository(t)
	userService := service.NewUserService(mockUserRepository, contextTimeout) // Create a new instance of the user service

	// Set up the expectation for the InsertUser method
	mockUserRepository.On("InsertUser", mock.Anything, mock.Anything).Return(1, nil)

	// Call the InsertUser method with a sample user
	userID, err := userService.InsertUser(context.Background(), models.User{
		Email:     "test@example.com",
		Password:  []byte("password123"),
		CreatedAt: time.Now(),
	})

	// Assert that no error occurred and the returned user ID is as expected
	assert.NoError(t, err)
	assert.Equal(t, 1, userID)

}

func TestUserService_UpdatePassword(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Create a mock user repository for testing
	mockUserRepository := mocks.NewUserRepository(t)
	userService := service.NewUserService(mockUserRepository, contextTimeout) // Create a new instance of the user service

	// Set up the expectation for the UpdatePassword method
	mockUserRepository.On("UpdatePassword", mock.Anything, mock.Anything).Return(nil)

	// Call the UpdatePassword method with a sample user ID and new password
	err := userService.UpdatePassword(context.Background(), models.User{
		ID:       1,
		Password: []byte("newPassword"),
	})

	// Assert that no error occurred during the password update
	assert.NoError(t, err)

}

func TestUserService_DeleteUser(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Create a mock user repository for testing
	mockUserRepository := mocks.NewUserRepository(t)
	userService := service.NewUserService(mockUserRepository, contextTimeout) // Create a new instance of the user service

	// Set up the expectation for the DeleteUser method
	mockUserRepository.On("DeleteUser", mock.Anything, 1).Return(nil)

	// Call the DeleteUser method with a sample user ID
	err := userService.DeleteUser(context.Background(), 1)

	// Assert that no error occurred during the deletion
	assert.NoError(t, err)

}

func TestUserService_GetUsersIDFromMemory(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Create a mock user repository for testing
	mockUserRepository := mocks.NewUserRepository(t)
	userService := service.NewUserService(mockUserRepository, contextTimeout) // Create a new instance of the user service

	// Call GetUsersIDFromMemory to retrieve users' IDs stored in memory
	usersIDs := userService.GetUsersIdFromMemory()

	// Assert that the returned value is of type ConcurrentMap[string, string]
	assert.IsType(t, cmap.ConcurrentMap[string, string]{}, usersIDs)
}
func TestUpdateRefreshTokenService(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Define test cases for updating refresh tokens
	tests := []struct {
		name string      // Name of the test case
		user models.User // User data for the test case
		err  error       // Expected error (nil for success)
	}{
		{
			name: "Successful Update",
			user: models.User{ID: 1, RefreshToken: []byte("new_refresh_token")},
			err:  nil, // No error expected for successful update
		},
		{
			name: "Error from UserRepository",
			user: models.User{ID: 1, RefreshToken: []byte("new_refresh_token")},
			err:  errors.New(""), // Simulate an error from the repository
		},
		{
			name: "Context Timeout",
			user: models.User{ID: 1, RefreshToken: []byte("new_refresh_token")},
			err:  errors.New(""), // Simulate a context timeout error
		},
	}

	// Iterate through each test case
	for _, tt := range tests {
		tc := tt // Capture the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			mockUserRepository := mocks.NewUserRepository(t)                                   // Create a new instance of the mocked repository
			mockUserRepository.On("UpdateRefreshToken", mock.Anything, tt.user).Return(tc.err) // Set up expectation

			userService := service.NewUserService(mockUserRepository, contextTimeout) // Create a new instance of the user service

			err := userService.UpdateRefreshToken(ctx, tt.user) // Call the method under test

			if tc.err == nil {
				assert.NoError(t, err) // Assert no error occurred for successful updates
			} else {
				assert.Error(t, err) // Assert that an error occurred as expected
			}
		})
	}
}

func TestGetUserById(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Define test cases for retrieving users by ID
	tests := []struct {
		name   string      // Name of the test case
		userID int         // User ID to retrieve
		err    error       // Expected error (nil for success)
		user   models.User // Expected user data (if any)
	}{
		{
			name:   "Get Existing User",
			userID: 1,
			err:    nil,
			user: models.User{
				ID:       1,
				Email:    "existinguser@example.com",
				Password: []byte("existingpassword"),
			},
		},
		{
			name:   "Get Non-Existent User",
			userID: 999,
			err:    errors.New(""), // Simulate an error when user does not exist
		},
	}

	// Iterate through each test case
	for _, tt := range tests {
		tc := tt // Capture the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			mockUserRepository := mocks.NewUserRepository(t)                                       // Create a new instance of the mocked repository
			mockUserRepository.On("GetUserById", mock.Anything, tc.userID).Return(tc.user, tc.err) // Set up expectation

			userService := service.NewUserService(mockUserRepository, time.Second*10) // Create a new instance of the user service

			_, err := userService.GetUserById(ctx, tc.userID) // Call the method under test

			if tc.err == nil {
				assert.NoError(t, err) // Assert no error occurred for existing users
			} else {
				assert.Error(t, err) // Assert that an error occurred as expected
			}
		})
	}
}

func TestGetUserByEmailService(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Define test cases for retrieving users by email
	tests := []struct {
		name  string      // Name of the test case
		email string      // Email to retrieve the user by
		err   error       // Expected error (nil for success)
		user  models.User // Expected user data (if any)
	}{
		{
			name:  "Get Existing User",
			email: "existinguser6870@example.com",
			err:   nil,
			user: models.User{
				ID:       1,
				Email:    "existinguser6870@example.com",
				Password: []byte("existingpassword"),
			},
		},
		{
			name:  "Get Non-Existent User",
			email: "nonexistent@example.com",
			err:   errors.New(""), // Simulate an error when email does not exist
		},
	}

	// Iterate through each test case
	for _, tt := range tests {
		tc := tt // Capture the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Allow this test case to run in parallel

			mockUserRepository := mocks.NewUserRepository(t)                                         // Create a new instance of the mocked repository
			mockUserRepository.On("GetUserByEmail", mock.Anything, tc.email).Return(tc.user, tc.err) // Set up expectation

			userService := service.NewUserService(mockUserRepository, time.Second*10) // Create a new instance of the user service

			_, err := userService.GetUserByEmail(ctx, tc.email) // Call the method under test

			if tc.err == nil {
				assert.NoError(t, err) // Assert no error occurred for existing users
			} else {
				assert.Error(t, err) // Assert that an error occurred as expected
			}
		})
	}
}

func TestUserService_GetUsersIdFromDB(t *testing.T) {
	t.Parallel() // Enable parallel execution for this test

	// Create a mock user repository for testing
	mockUserRepository := mocks.NewUserRepository(t)
	userService := service.NewUserService(mockUserRepository, contextTimeout) // Create a new instance of the user service

	// Set up the expectation for the GetAllIDs method to return a slice of user IDs
	mockUserRepository.On("GetAllIDs", mock.Anything).Return([]int{1, 2, 3}, nil)

	// Call the GetUsersIdFromDB method to retrieve user IDs from the database
	err := userService.GetUsersIdFromDB(context.Background())

	// Assert that no error occurred during the retrieval
	assert.NoError(t, err)

	// Retrieve user IDs stored in memory
	usersIDs := userService.GetUsersIdFromMemory()

	// Check that the expected user IDs are present in memory
	one, _ := usersIDs.Get("1")
	two, _ := usersIDs.Get("2")
	three, _ := usersIDs.Get("3")

	// Assert that the retrieved IDs match expected values
	assert.Equal(t, "1", one)
	assert.Equal(t, "2", two)
	assert.Equal(t, "3", three)

}

func TestSetUserIdIntoMemory(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	tests := []struct {
		name     string // Name of the test case
		userID   int    // User ID to be set into memory
		expected string // Expected value in memory after setting
	}{
		{
			name:     "Set Valid User ID",
			userID:   123,
			expected: "123", // Expecting the user ID to be stored as a string
		},
		{
			name:     "Set Another Valid User ID",
			userID:   456,
			expected: "456", // Expecting the user ID to be stored as a string
		},
	}

	for _, tt := range tests {
		tc := tt // Capture range variable for use in goroutine

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run each test case in parallel

			mockUserRepository := mocks.NewUserRepository(t)
			userService := service.NewUserService(mockUserRepository, contextTimeout) // Create a new instance of the user service

			userService.SetUserIdIntoMemory(tc.userID) // Call the method to set user ID into memory

			IDs := userService.GetUsersIdFromMemory() // Retrieve value from memory

			assert.True(t, IDs.Has(strconv.Itoa(tc.userID))) // Assert that the value exists in memory
		})
	}
}

func TestDeleteUserIdFromMemory(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	tests := []struct {
		name     string // Name of the test case
		userID   int    // User ID to be set into memory
		expected string // Expected value in memory after setting
	}{
		{
			name:     "Set Valid User ID",
			userID:   123,
			expected: "123", // Expecting the user ID to be stored as a string
		},
		{
			name:     "Set Another Valid User ID",
			userID:   456,
			expected: "456", // Expecting the user ID to be stored as a string
		},
	}

	for _, tt := range tests {
		tc := tt // Capture range variable for use in goroutine

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run each test case in parallel

			mockUserRepository := mocks.NewUserRepository(t)
			userService := service.NewUserService(mockUserRepository, contextTimeout) // Create a new instance of the user service

			userService.SetUserIdIntoMemory(tc.userID)    // Call the method to set user ID into memory
			userService.DeleteUserIdFromMemory(tc.userID) // Call the method to set user ID into memory

			IDs := userService.GetUsersIdFromMemory() // Retrieve value from memory

			assert.False(t, IDs.Has(strconv.Itoa(tc.userID))) // Assert that the value exists in memory
		})
	}
}
