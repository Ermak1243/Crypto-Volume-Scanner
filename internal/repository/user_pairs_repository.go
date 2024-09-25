package repository

import (
	"context"
	"fmt"
	"log"
	"main/internal/domain/models" // Importing domain models for user pairs

	"github.com/jmoiron/sqlx" // Importing sqlx for database interactions
)

// UserPairsRepository defines the interface for operations related to user pairs.
// It includes methods for adding, updating, retrieving, and deleting user pairs.
type UserPairsRepository interface {
	Add(ctx context.Context, pairData models.UserPairs) error                    // Method to add a new user pair
	UpdateExactValue(ctx context.Context, pairData models.UserPairs) error       // Method to update the exact value of a user pair
	GetAllUserPairs(ctx context.Context, userID int) ([]models.UserPairs, error) // Method to retrieve all user pairs for a given user ID
	GetPairsByExchange(ctx context.Context, exchange string) ([]string, error)   // Method to retrieve all pairs for a given exchange name
	DeletePair(ctx context.Context, pairData models.UserPairs) error             // Method to delete a specific user pair
}

// userPairsRepository is a concrete implementation of the UserPairsRepository interface.
// It holds a reference to the database connection.
type userPairsRepository struct {
	db *sqlx.DB // Database connection
}

// NewUserPairsRepository creates a new instance of userPairsRepository.
// It initializes the repository with a database connection.
//
// Parameters:
//   - db: The database connection to be used by the repository.
//
// Returns:
//   - An instance of UserPairsRepository.
func NewUserPairsRepository(db *sqlx.DB) UserPairsRepository {
	return &userPairsRepository{db} // Return a new instance of userPairsRepository
}

// Add inserts a new user pair into the database.
// It takes context and pair data as parameters and returns an error if any occurs.
func (upr *userPairsRepository) Add(ctx context.Context, pairData models.UserPairs) error {
	const op = directoryPath + "user_pairs_repository.Add" // Operation name for logging
	errFn := repoError("Add")                              // Error handling function

	queryString := fmt.Sprintf(`
		INSERT INTO %s (
			user_id,
			exchange, 
			pair,
			exact_value
		)
		values ($1, $2, $3, $4)
	`, userPairsTable) // SQL query string for inserting data

	_, err := upr.db.ExecContext(
		ctx,
		queryString,
		pairData.UserID,
		pairData.Exchange,
		pairData.Pair,
		pairData.ExactValue,
	) // Execute the SQL query with provided parameters
	if err != nil {
		log.Println(op, ": ", err) // Log any errors that occur during execution

		return errFn // Return wrapped error
	}

	return nil // Return nil if no errors occurred
}

// UpdateExactValue updates the exact value of an existing user pair in the database.
// It takes context and pair data as parameters and returns an error if any occurs.
func (upr *userPairsRepository) UpdateExactValue(ctx context.Context, pairData models.UserPairs) error {
	const op = directoryPath + "user_pairs_repository.UpdateExactValue" // Operation name for logging
	errFn := repoError("UpdateExactValue")                              // Error handling function

	queryString := fmt.Sprintf(`
		UPDATE %s 
		SET exact_value=$1
		WHERE user_id=$2 AND exchange=$3 AND pair=$4;
	`, userPairsTable) // SQL query string for updating data

	rows, err := upr.db.ExecContext(
		ctx,
		queryString,
		pairData.ExactValue,
		pairData.UserID,
		pairData.Exchange,
		pairData.Pair,
	) // Execute the SQL query with provided parameters
	rowsAffected, _ := rows.RowsAffected() // Get the number of rows affected by the update
	if err != nil || rowsAffected == 0 {   // Check for errors or if no rows were updated
		log.Println(op, ": ", err) // Log any errors that occur during execution

		return errFn // Return wrapped error
	}

	return nil // Return nil if no errors occurred
}

// GetAllUserPairs retrieves all user pairs associated with a given user ID from the database.
// It takes context and user ID as parameters and returns a slice of UserPairs and an error if any occurs.
func (upr *userPairsRepository) GetAllUserPairs(ctx context.Context, userID int) ([]models.UserPairs, error) {
	const op = directoryPath + "user_pairs_repository.GetAllUserPairs" // Operation name for logging
	errFn := repoError("GetAllUserPairs")                              // Error handling function
	var userPairs []models.UserPairs                                   // Slice to hold retrieved user pairs

	queryString := fmt.Sprintf(`
		SELECT * FROM %s WHERE user_id=%d;
	`, userPairsTable, userID) // SQL query string for selecting data

	err := upr.db.SelectContext(ctx, &userPairs, queryString) // Execute the SQL query and scan results into the slice
	if err != nil {
		log.Println(op, ": ", err) // Log any errors that occur during execution

		return userPairs, errFn // Return empty slice and wrapped error
	}

	return userPairs, nil // Return retrieved user pairs and nil if no errors occurred
}

// GetPairsByExchange retrieves all user pairs for a given exchange name from the database.
// It takes context and exchange name as parameters and returns a slice of strings and an error if any occurs.
func (upr *userPairsRepository) GetPairsByExchange(ctx context.Context, exchange string) ([]string, error) {
	const op = directoryPath + "user_pairs_repository.GetPairsByExchange" // Operation name for logging
	errFn := repoError("GetPairsByExchange")                              // Error handling function
	var exchangePairs []string                                            // Slice to hold retrieved user pairs

	queryString := fmt.Sprintf(`
		SELECT DISTINCT pair FROM %s WHERE exchange='%s';
	`, userPairsTable, exchange) // SQL query string for selecting data

	err := upr.db.SelectContext(ctx, &exchangePairs, queryString) // Execute the SQL query and scan results into the slice
	if err != nil {
		log.Println(op, ": ", err) // Log any errors that occur during execution

		return exchangePairs, errFn // Return empty slice and wrapped error
	}

	return exchangePairs, nil // Return retrieved user pairs and nil if no errors occurred
}

// DeletePair removes a specific user pair from the database.
// It takes context and pair data as parameters and returns an error if any occurs.
func (upr *userPairsRepository) DeletePair(ctx context.Context, pairData models.UserPairs) error {
	const op = directoryPath + "user_pairs_repository.DeletePair" // Operation name for logging
	errFn := repoError("DeletePair")                              // Error handling function

	queryString := fmt.Sprintf(`
		DELETE FROM %s 
		WHERE user_id=$1 AND pair=$2
	`, userPairsTable) // SQL query string for deleting data

	rows, err := upr.db.ExecContext(
		ctx,
		queryString,
		pairData.UserID,
		pairData.Pair,
	) // Execute the SQL query with provided parameters
	rowsAffected, _ := rows.RowsAffected() // Get the number of rows affected by the delete operation
	if err != nil || rowsAffected == 0 {   // Check for errors or if no rows were deleted
		log.Println(op, ": ", err) // Log any errors that occur during execution

		return errFn // Return wrapped error
	}

	return nil // Return nil if no errors occurred
}
