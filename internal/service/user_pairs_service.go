package service

import (
	"context"
	"main/internal/models"
	"main/internal/repository"
	"time"
)

// UserPairsService defines the interface for working with user pair settings.
// This interface includes methods for adding, updating, retrieving, and deleting user pairs.
type UserPairsService interface {
	Add(ctx context.Context, pairData models.UserPairs) error
	UpdateExactValue(ctx context.Context, pairData models.UserPairs) error
	GetAllUserPairs(ctx context.Context, userID int) ([]models.UserPairs, error)
	GetPairsByExchange(ctx context.Context, exchange string) ([]string, error)
	DeletePair(ctx context.Context, pairData models.UserPairs) error
}

// userPairsService is a concrete implementation of UserPairsService.
// It holds a reference to the UserPairsRepository and a timeout duration.
type userPairsService struct {
	userPairsRepository repository.UserPairsRepository // Repository for accessing user pairs data
	contextTimeout      time.Duration                  // Timeout duration for context
}

// NewUserPairsService creates a new instance of userPairsService.
// It takes a UserPairsRepository and a timeout duration as parameters.
//
// Parameters:
//   - userPairsRepository: Repository for managing user pairs data.
//   - timeout: Duration to set context timeout for operations.
//
// Returns:
//   - An instance of UserPairsService.
func NewUserPairsService(userPairsRepository repository.UserPairsRepository, timeout time.Duration) UserPairsService {
	return &userPairsService{
		userPairsRepository: userPairsRepository,
		contextTimeout:      timeout,
	}
}

// Add inserts user pair data into the database.
// It validates the pair data before attempting to add it to the repository.
//
// Parameters:
//   - ctx: The context for managing request lifetime.
//   - pairData: The user pair data to be added.
//
// Returns:
//   - An error if the operation fails; otherwise, nil.
func (ups *userPairsService) Add(ctx context.Context, pairData models.UserPairs) error {
	// Validate the pair data using a separate validation function.
	if err := CheckPairData(pairData); err != nil {
		return err // Return validation error
	}

	ctx, cancel := context.WithTimeout(ctx, ups.contextTimeout) // Set up context with timeout
	defer cancel()                                              // Ensure cancellation of context when done

	// Attempt to add the pair data using the repository.
	if err := ups.userPairsRepository.Add(ctx, pairData); err != nil {
		return err // Return any errors from the repository
	}

	return nil // Return nil if successful
}

// UpdateExactValue updates existing pair settings in the database.
// It validates the pair data before attempting to update it in the repository.
//
// Parameters:
//   - ctx: The context for managing request lifetime.
//   - pairData: The user pair data with updated values.
//
// Returns:
//   - An error if the operation fails; otherwise, nil.
func (ups *userPairsService) UpdateExactValue(ctx context.Context, pairData models.UserPairs) error {
	// Validate the pair data before proceeding with the update.
	if err := CheckPairData(pairData); err != nil {
		return err // Return validation error
	}

	ctx, cancel := context.WithTimeout(ctx, ups.contextTimeout) // Set up context with timeout
	defer cancel()                                              // Ensure cancellation of context when done

	// Attempt to update the exact value using the repository.
	if err := ups.userPairsRepository.UpdateExactValue(ctx, pairData); err != nil {
		return err // Return any errors from the repository
	}

	return nil // Return nil if successful
}

// DeletePair removes a user pair from the database.
// It validates that the user ID and pair name are provided before attempting to delete.
//
// Parameters:
//   - ctx: The context for managing request lifetime.
//   - pairData: The user pair data to be deleted.
//
// Returns:
//   - An error if validation fails or if the operation fails; otherwise, nil.
func (ups *userPairsService) DeletePair(ctx context.Context, pairData models.UserPairs) error {
	// Validate that user ID is greater than zero.
	if pairData.UserID < 1 {
		err := errIdBelowOne // Custom error indicating invalid user ID
		return err           // Return validation error
	}

	// Validate that the pair name is not empty.
	if pairData.Pair == "" {
		err := errPairNameIsEmpty // Custom error indicating empty pair name
		return err                // Return validation error
	}

	ctx, cancel := context.WithTimeout(ctx, ups.contextTimeout) // Set up context with timeout
	defer cancel()                                              // Ensure cancellation of context when done

	// Attempt to delete the pair using the repository.
	if err := ups.userPairsRepository.DeletePair(ctx, pairData); err != nil {
		return err // Return any errors from the repository
	}

	return nil // Return nil if successful
}

// GetAllUserPairs retrieves all user pairs from the database for a given user ID.
//
// Parameters:
//   - ctx: The context for managing request lifetime.
//   - userID: The ID of the user whose pairs are to be retrieved.
//
// Returns:
//   - A slice of UserPairs and an error if any occurs during retrieval.
func (ups *userPairsService) GetAllUserPairs(ctx context.Context, userID int) ([]models.UserPairs, error) {
	userPairs, err := ups.userPairsRepository.GetAllUserPairs(ctx, userID)
	if err != nil {
		return userPairs, err // Return empty slice and error if retrieval fails
	}

	return userPairs, nil // Return retrieved pairs if successful
}

// GetPairsByExchange retrieves all user pairs associated with a given exchange name from the database.
//
// Parameters:
//   - ctx: The context for managing request lifetime.
//   - exchange: The name of the exchange whose pairs are to be retrieved.
//
// Returns:
//   - A slice of strings and an error if any occurs during retrieval.
func (ups *userPairsService) GetPairsByExchange(ctx context.Context, exchange string) ([]string, error) {
	exchangePairs, err := ups.userPairsRepository.GetPairsByExchange(ctx, exchange)
	if err != nil {
		return exchangePairs, err // Return empty slice and error if retrieval fails
	}

	return exchangePairs, nil // Return retrieved pairs if successful
}
