package repository

import (
	"context"
	"cvs/internal/models" // Importing domain models for user data
	"fmt"

	"github.com/jmoiron/sqlx" // Importing sqlx for database interactions
)

// UserRepository defines the interface for operations related to users.
// It includes methods for inserting, updating, retrieving, and deleting user records.
type UserRepository interface {
	InsertUser(ctx context.Context, user models.User) (int, error)         // Method to insert a new user
	UpdatePassword(ctx context.Context, user models.User) error            // Method to update a user's password
	UpdateRefreshToken(ctx context.Context, user models.User) error        // Method to update a user's refresh token
	GetUserById(ctx context.Context, userID int) (models.User, error)      // Method to retrieve a user by ID
	GetUserByEmail(ctx context.Context, email string) (models.User, error) // Method to retrieve a user by email
	GetAllIDs(ctx context.Context) ([]int, error)                          // Method to get all user IDs
	DeleteUser(ctx context.Context, clientID int) error                    // Method to delete a user by ID
}

// userRepository is a concrete implementation of the UserRepository interface.
// It holds a reference to the database connection.
type userRepository struct {
	db *sqlx.DB // Database connection
}

// NewUserRepository creates a new instance of userRepository.
// It initializes the repository with a database connection.
//
// Parameters:
//   - db: The database connection to be used by the repository.
//
// Returns:
//   - An instance of UserRepository.
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db} // Return a new instance of userRepository
}

// InsertUser inserts a new user into the database.
// It returns the newly created user's ID and an error if any occurs.
func (ur *userRepository) InsertUser(ctx context.Context, user models.User) (int, error) {
	const op = directoryPath + "user_repository.InsertUser" // Operation name for logging

	var clientID int // Variable to hold the newly created user's ID
	query := fmt.Sprintf(`
		INSERT INTO %s (
			email,
			password,
			refresh_token,
			session_id
		)
		values ($1, $2, $3, $4)
		RETURNING id;				
	`, userTable) // SQL query string for inserting data

	err := ur.db.GetContext(
		ctx,
		&clientID,
		query,
		user.Email,
		user.Password,
		user.RefreshToken,
		user.SessionID,
	) // Execute the SQL query and return the newly created user's ID
	if err != nil {
		return 0, repoError(op) // Return zero ID and wrapped error
	}

	return clientID, nil // Return the newly created user's ID and nil if no errors occurred
}

// UpdatePassword updates an existing user's password in the database.
// It returns an error if any occurs.
func (ur *userRepository) UpdatePassword(ctx context.Context, user models.User) error {
	const op = directoryPath + "user_repository.UpdatePassword" // Operation name for logging

	query := fmt.Sprintf(`
		UPDATE %s 
		SET password=$1,
			refresh_token=$2,
			session_id=$3,
			updated_at='now()'
		WHERE id=$4;`, userTable) // SQL query string for updating data

	rows, err := ur.db.ExecContext(
		ctx,
		query,
		user.Password,
		user.RefreshToken,
		user.SessionID,
		user.ID,
	) // Execute the SQL query with provided parameters
	rowsAffected, _ := rows.RowsAffected() // Get the number of rows affected by the update
	if err != nil || rowsAffected == 0 {   // Check for errors or if no rows were updated
		return repoError(op) // Return wrapped error
	}

	return nil // Return nil if no errors occurred
}

// UpdateRefreshToken updates an existing user's refresh token in the database.
// It returns an error if any occurs.
func (ur *userRepository) UpdateRefreshToken(ctx context.Context, user models.User) error {
	const op = directoryPath + "user_repository.UpdateRefreshToken" // Operation name for logging

	query := fmt.Sprintf(`
		UPDATE %s 
		SET refresh_token=$1,
			session_id=$2,
			updated_at='now()'
		WHERE id=$3;`, userTable) // SQL query string for updating data

	rows, err := ur.db.ExecContext(
		ctx,
		query,
		user.RefreshToken,
		user.SessionID,
		user.ID,
	) // Execute the SQL query with provided parameters
	rowsAffected, _ := rows.RowsAffected() // Get the number of rows affected by the update
	if err != nil || rowsAffected == 0 {   // Check for errors or if no rows were updated
		return repoError(op) // Return wrapped error
	}

	return nil // Return nil if no errors occurred
}

// GetUserById retrieves a user from the database by their ID.
// It returns the user and an error if any occurs.
func (ur *userRepository) GetUserById(ctx context.Context, userID int) (models.User, error) {
	const op = directoryPath + "user_repository.GetUserById" // Operation name for logging
	var user models.User                                     // Variable to hold retrieved user

	query := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d;`, userTable, userID) // SQL query string for selecting data

	err := ur.db.GetContext(ctx, &user, query) // Execute the SQL query and scan results into the user variable
	if err != nil {
		return user, repoError(op) // Return empty user and wrapped error
	}

	return user, nil // Return retrieved user and nil if no errors occurred
}

// GetUserByEmail retrieves a user from the database by their email address.
//
// This method constructs a SQL query to select a user based on their email and executes it.
// If successful, it returns the retrieved user. If an error occurs during execution,
// it logs the error and returns an empty user object along with the error.
//
// Parameters:
//   - ctx: A context.Context instance for managing request-scoped values, cancellation signals,
//     and deadlines across API boundaries.
//   - email: A string representing the email address of the user to be retrieved.
//
// Returns:
//   - models.User: The user associated with the provided email address.
//   - error: An error if any issues occur during the database query execution. If successful,
//     it returns nil for the error.
func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = directoryPath + "user_repository.GetUserByEmail" // Operation name for logging
	var user models.User                                        // Variable to hold retrieved user

	query := fmt.Sprintf(`SELECT * FROM %s WHERE email='%s';`, userTable, email) // SQL query string for selecting data

	err := ur.db.GetContext(ctx, &user, query) // Execute the SQL query and scan results into the user variable
	if err != nil {
		return user, repoError(op) // Return empty user and wrapped error
	}

	return user, nil // Return retrieved user and nil if no errors occurred
}

// GetAllIDs retrieves all unique user IDs from the database.
// It returns a slice of IDs and an error if any occurs.
func (ur *userRepository) GetAllIDs(ctx context.Context) ([]int, error) {
	const op = directoryPath + "user_repository.GetAllIDs" // Operation name for logging
	var allIDs []int                                       // Slice to hold all retrieved IDs

	query := fmt.Sprintf(`SELECT DISTINCT id FROM %s;`, userTable) // SQL query string for selecting distinct IDs

	err := ur.db.SelectContext(ctx, &allIDs, query) // Execute the SQL query and scan results into allIDs slice
	if err != nil {
		return allIDs, repoError(op) // Return empty slice and wrapped error
	}

	return allIDs, nil // Return retrieved IDs and nil if no errors occurred
}

// DeleteUser removes a specific user from the database by their ID.
// It returns an error if any occurs.
func (ur *userRepository) DeleteUser(ctx context.Context, clientID int) error {
	const op = directoryPath + "user_repository.DeleteUser" // Operation name for logging

	query := fmt.Sprintf(`
        DELETE FROM %s 
        WHERE id=$1`, userTable) // SQL query string for deleting data

	rows, err := ur.db.ExecContext(ctx, query, clientID) // Execute the SQL query with provided parameters
	rowsAffected, _ := rows.RowsAffected()               // Get number of rows affected by delete operation
	if err != nil || rowsAffected == 0 {                 // Check for errors or if no rows were deleted
		return repoError(op) // Return wrapped error
	}

	return nil // Return nil if no errors occurred
}
