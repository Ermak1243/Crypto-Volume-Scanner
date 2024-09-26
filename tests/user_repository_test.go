package tests

import (
	"bytes"
	"testing"

	"main/internal/models"
	"main/internal/repository"

	"github.com/stretchr/testify/assert"
)

// TestInsertUser tests the InsertUser function of the UserRepository.
func TestInsertUser(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	// Define test cases
	tests := []struct {
		name    string      // Name of the test case
		user    models.User // User data to be inserted
		wantErr bool        // Expected outcome: true if an error is expected
	}{
		{
			name:    "Valid User",
			user:    models.User{Email: "newuser1@example.com", Password: []byte("newpassword123"), SessionID: 1}, // Include SessionID
			wantErr: false,                                                                                        // No error expected for valid user
		},
		{
			name:    "Invalid User - Empty Email",
			user:    models.User{Email: "", Password: []byte("password123"), SessionID: 1}, // Include SessionID
			wantErr: true,                                                                  // Error expected due to empty email
		},
		{
			name:    "Invalid User - Empty Password",
			user:    models.User{Email: "user2@example.com", Password: []byte{}, SessionID: 1}, // Include SessionID
			wantErr: true,                                                                      // Error expected due to empty password
		},
	}

	for _, tt := range tests {
		tc := tt // Create a copy of the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run this test case in parallel

			db := setupDB()  // Set up the database connection for testing
			defer db.Close() // Ensure the database connection is closed after the test

			userRepo := repository.NewUserRepository(db) // Initialize the user repository

			id, err := userRepo.InsertUser(ctx, tc.user)      // Attempt to insert the user into the database
			defer db.ExecContext(ctx, deleteUserQueryRow, id) // Clean up by deleting the user after the test

			if tc.wantErr {
				assert.Error(t, err) // Check that an error occurred if one was expected
			} else {
				assert.NoError(t, err) // Check that no error occurred

				var retrievedUser models.User
				query := `SELECT id, email, password, session_id FROM users WHERE id = $1` // Query to retrieve the inserted user
				db.GetContext(ctx, &retrievedUser, query, id)                              // Execute the query

				assert.NoError(t, err)                                                // Ensure no error occurred while retrieving the user
				assert.Equal(t, id, retrievedUser.ID)                                 // Check that the retrieved ID matches the inserted ID
				assert.Equal(t, tc.user.Email, retrievedUser.Email)                   // Verify that the email matches
				assert.True(t, bytes.Equal(tc.user.Password, retrievedUser.Password)) // Check that passwords match
				assert.Equal(t, tc.user.SessionID, retrievedUser.SessionID)           // Verify that SessionID matches
			}
		})
	}
}

// TestUpdatePassword tests the UpdatePassword function of the UserRepository.
func TestUpdatePassword(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	tests := []struct {
		name    string      // Name of the test case
		user    models.User // User data containing new password information
		wantErr bool        // Expected outcome: true if an error is expected
	}{
		{
			name:    "Valid User",
			user:    models.User{Email: "newuser3@example.com", Password: []byte("newpassword123"), SessionID: 1}, // Include SessionID
			wantErr: false,                                                                                        // No error expected for valid update
		},
		{
			name:    "Invalid User",
			user:    models.User{Email: "", Password: []byte(""), SessionID: 1}, // Include SessionID but invalid data
			wantErr: true,                                                       // Error expected due to invalid user data (empty fields)
		},
	}

	for _, tt := range tests {
		tc := tt // Create a copy of the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run this test case in parallel

			db := setupDB()  // Set up the database connection for testing
			defer db.Close() // Ensure the database connection is closed after the test

			if !tc.wantErr {
				id, err := insertUser(db, tc.user.Email, tc.user.Password) // Insert user if no error is expected
				defer db.ExecContext(ctx, deleteUserQueryRow, id)          // Clean up by deleting the user after the test
				assert.NoError(t, err)                                     // Ensure no error occurred during insertion

				tc.user.ID = id // Set the ID of the user for further operations
			}

			userRepo := repository.NewUserRepository(db) // Initialize the user repository

			err := userRepo.UpdatePassword(ctx, tc.user) // Attempt to update user's password in database
			if tc.wantErr {
				assert.Error(t, err) // Check that an error occurred if one was expected
			} else {
				assert.NoError(t, err) // Check that no error occurred

				var retrievedPassword []byte
				query := `SELECT password FROM users WHERE id = $1` // Query to retrieve updated password from database
				err = db.GetContext(ctx, &retrievedPassword, query, tc.user.ID)

				assert.NoError(t, err)                                           // Ensure no error occurred while retrieving updated password
				assert.True(t, bytes.Equal(retrievedPassword, tc.user.Password)) // Check that updated password matches expected value
			}
		})
	}
}

// TestUpdateRefreshToken tests the UpdateRefreshToken function of the UserRepository.
func TestUpdateRefreshToken(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	tests := []struct {
		name    string      // Name of the test case
		user    models.User // User data including RefreshToken
		wantErr bool        // Expected outcome: true if an error is expected
	}{
		{
			name:    "Update Refresh Token - Success",
			user:    models.User{ID: 1, Email: "newuser4453@example.com", Password: []byte("newpassword123"), RefreshToken: []byte("new_refresh_token"), SessionID: 1}, // Include SessionID
			wantErr: false,                                                                                                                                             // No error expected for successful update
		},
		{
			name:    "Update Refresh Token - Error",
			user:    models.User{ID: 1, Email: "newuser4453@example.com", Password: []byte("newpassword123"), RefreshToken: []byte("new_refresh_token"), SessionID: 1}, // Include SessionID
			wantErr: true,                                                                                                                                              // Error expected (e.g., user not found)
		},
	}

	for _, tt := range tests {
		tc := tt // Create a copy of the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run this test case in parallel

			db := setupDB()  // Set up the database connection for testing
			defer db.Close() // Ensure the database connection is closed after the test

			if !tc.wantErr {
				id, err := insertUser(db, tc.user.Email, tc.user.Password) // Insert user if no error is expected
				defer db.ExecContext(ctx, deleteUserQueryRow, id)          // Clean up by deleting the user after the test
				assert.NoError(t, err)                                     // Ensure no error occurred during insertion

				tc.user.ID = id // Set the ID of the user for further operations
			}

			userRepo := repository.NewUserRepository(db) // Initialize the user repository

			err := userRepo.UpdateRefreshToken(ctx, tc.user) // Attempt to update user's refresh token in database
			if tc.wantErr {
				assert.Error(t, err) // Check that an error occurred if one was expected
			} else {
				assert.NoError(t, err) // Check that no error occurred
			}
		})
	}
}

// TestGetUserByID tests the GetUserById function of the UserRepository.
func TestGetUserByID(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	tests := []struct {
		name    string      // Name of the test case
		wantErr bool        // Expected outcome: true if an error is expected
		user    models.User // User data to be retrieved
	}{
		{
			name:    "Get Existing User",
			wantErr: false, // No error expected when retrieving an existing user
			user: models.User{
				Email:    "newuser53790@example.comuser",
				Password: []byte("newpassword123"),
			},
		},
		{
			name:    "Error",
			wantErr: true, // Error expected when trying to retrieve a non-existing user
			user: models.User{
				Email:    "newuser52178@example.comuser",
				Password: []byte("newpassword123"),
			},
		},
	}

	for _, tt := range tests {
		tc := tt // Create a copy of the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run this test case in parallel

			db := setupDB()  // Set up the database connection for testing
			defer db.Close() // Ensure the database connection is closed after the test

			if !tc.wantErr {
				userID, err := insertUser(db, tc.user.Email, tc.user.Password) // Insert user if no error is expected
				defer db.ExecContext(ctx, deleteUserQueryRow, userID)          // Clean up by deleting the user after the test
				assert.NoError(t, err)                                         // Ensure no error occurred during insertion

				tc.user.ID = userID // Set the ID of the user for further operations
			}

			userRepo := repository.NewUserRepository(db) // Initialize the user repository

			user, err := userRepo.GetUserById(ctx, tc.user.ID) // Attempt to retrieve user by ID
			if !tc.wantErr {
				assert.NoError(t, err)                     // Check that no error occurred when retrieving existing user
				assert.Equal(t, tc.user.ID, user.ID)       // Verify that retrieved ID matches expected ID
				assert.Equal(t, tc.user.Email, user.Email) // Verify that retrieved email matches expected email
			} else {
				assert.Error(t, err) // Check that an error occurred when retrieving non-existing user
			}
		})
	}
}

// TestGetUserByEmail tests the GetUserByEmail function of the UserRepository.
func TestGetUserByEmail(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	tests := []struct {
		name    string      // Name of the test case
		wantErr bool        // Expected outcome: true if an error is expected
		user    models.User // User data to be retrieved by email
	}{
		{
			name:    "Get Existing User",
			wantErr: false, // No error expected for retrieving an existing user
			user: models.User{
				Email:    "newuser5790@example.comuser",
				Password: []byte("newpassword123"),
			},
		},
		{
			name:    "Error",
			wantErr: true, // Error expected when trying to retrieve a non-existing user
			user: models.User{
				Email:    "newuser578@example.comuser",
				Password: []byte("newpassword123"),
			},
		},
	}

	for _, tt := range tests {
		tc := tt // Create a copy of the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run this test case in parallel

			db := setupDB()  // Set up the database connection for testing
			defer db.Close() // Ensure the database connection is closed after the test

			if !tc.wantErr {
				userID, err := insertUser(db, tc.user.Email, tc.user.Password) // Insert user if no error is expected
				defer db.ExecContext(ctx, deleteUserQueryRow, userID)          // Clean up by deleting the user after the test
				assert.NoError(t, err)                                         // Ensure no error occurred during insertion

				tc.user.ID = userID // Set the ID of the user for further operations
			}

			userRepo := repository.NewUserRepository(db) // Initialize the user repository

			user, err := userRepo.GetUserByEmail(ctx, tc.user.Email) // Attempt to retrieve user by email
			if !tc.wantErr {
				assert.NoError(t, err)                     // Check that no error occurred when retrieving existing user
				assert.Equal(t, tc.user.ID, user.ID)       // Verify that retrieved ID matches expected ID
				assert.Equal(t, tc.user.Email, user.Email) // Verify that retrieved email matches expected email
			} else {
				assert.Error(t, err) // Check that an error occurred when retrieving non-existing user
			}
		})
	}
}

// TestGetAllIDs tests the GetAllIDs function of the UserRepository.
func TestGetAllIDs(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	tests := []struct {
		name         string      // Name of the test case
		expected     int         // Expected number of IDs returned by GetAllIDs function
		user         models.User // User data to be inserted before testing (if needed)
		zeroUsersLen bool        // Flag indicating whether to expect zero users in DB or not
	}{
		{
			name:         "Get All IDs - Expect at least 1",
			expected:     1,
			user:         models.User{Email: "newuser4@example.comuser", Password: []byte("newpassword123")},
			zeroUsersLen: false,
		},
		{
			name:         "Get All IDs Empty",
			user:         models.User{Email: "", Password: []byte("")}, // No valid data provided
			expected:     0,
			zeroUsersLen: true,
		},
	}

	for _, tt := range tests {
		tc := tt // Create a copy of the current test case

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run this test case in parallel

			db := setupDB()  // Set up the database connection for testing
			defer db.Close() // Ensure the database connection is closed after the test

			if !tc.zeroUsersLen {
				id, err := insertUser(db, tc.user.Email, tc.user.Password) // Insert user if no error is expected
				defer db.ExecContext(ctx, deleteUserQueryRow, id)          // Clean up by deleting the user after the test
				assert.NoError(t, err)                                     // Ensure no error occurred during insertion

				tc.user.ID = id // Set ID of inserted user for further operations
			}

			userRepo := repository.NewUserRepository(db) // Initialize the user repository

			ids, err := userRepo.GetAllIDs(ctx) // Attempt to retrieve all user IDs

			assert.NoError(t, err)                          // Ensure no error occurred during retrieval
			assert.GreaterOrEqual(t, len(ids), tc.expected) // Check that returned IDs meet expectations
		})
	}
}

// TestDeleteUser tests the DeleteUser function of the UserRepository.
func TestDeleteUser(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	tests := []struct {
		name    string      // Name of the test case
		user    models.User // User data containing ID to be deleted
		wantErr bool        // Expected outcome: true if an error is expected during deletion
	}{
		{
			name:    "Valid Delete",
			user:    models.User{Email: "newuser6@example.comuser", Password: []byte("newpassword123")},
			wantErr: false, // No error expected for valid deletion
		},
		{
			name:    "Invalid Delete",
			user:    models.User{Email: "", Password: []byte("")}, // Invalid data (no email/password)
			wantErr: true,                                         // Error expected due to invalid input
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := setupDB()  // Set up the database connection for testing
			defer db.Close() // Ensure the database connection is closed after the test

			if !tc.wantErr {
				id, err := insertUser(db, tc.user.Email, tc.user.Password) // Insert user if no error is expected
				defer db.ExecContext(ctx, deleteUserQueryRow, id)          // Clean up by deleting the inserted user after test completion
				assert.NoError(t, err)                                     // Ensure no error occurred during insertion

				tc.user.ID = id // Set ID of inserted user for further operations
			}

			userRepo := repository.NewUserRepository(db) // Initialize the user repository

			err := userRepo.DeleteUser(ctx, tc.user.ID) // Attempt to delete user by ID

			if tc.wantErr {
				assert.Error(t, err) // Check that an error occurred if one was expected
			} else {
				assert.NoError(t, err) // Ensure no error occurred during deletion

				var deletedID int
				query := `SELECT id FROM users WHERE id=$1` // Query to check if ID still exists

				err = db.GetContext(ctx, &deletedID, query, tc.user.ID)

				assert.Error(t, err) // Check that an error occurs when trying to retrieve a deleted ID
			}
		})
	}
}
