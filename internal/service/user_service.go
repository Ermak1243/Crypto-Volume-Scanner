package service

import (
	"context"
	"main/internal/models"
	"main/internal/repository"
	"strconv"

	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

// UserService defines the interface for user-related operations.
// This interface includes methods for inserting, updating, retrieving, and deleting users.
type UserService interface {
	InsertUser(ctx context.Context, user models.User) (int, error)         // Insert a new user
	UpdatePassword(ctx context.Context, user models.User) error            // Update an existing user's password
	UpdateRefreshToken(c context.Context, user models.User) error          // Update an existing user's refresh token
	GetUsersIdFromDB(ctx context.Context) error                            // Get all user IDs from the database
	GetUserById(ctx context.Context, userID int) (models.User, error)      // Get a user by ID
	GetUserByEmail(ctx context.Context, email string) (models.User, error) // Get a user by email
	GetUsersIdFromMemory() cmap.ConcurrentMap[string, string]              // Get all user IDs from memory
	SetUserIdIntoMemory(userID int)                                        // Set a user ID into memory
	DeleteUserIdFromMemory(userID int)                                     // Delete a user ID from memory
	DeleteUser(ctx context.Context, userID int) error                      // Delete a user by ID
}

// userService is a concrete implementation of UserService.
// It holds a reference to the UserRepository and a concurrent map for storing user IDs.
type userService struct {
	userRepository repository.UserRepository          // Repository for accessing user data
	usersIDs       cmap.ConcurrentMap[string, string] // Concurrent map for storing user IDs in memory
	contextTimeout time.Duration                      // Timeout duration for context management
}

// NewUserService creates a new instance of userService.
// It initializes the usersIDs concurrent map and sets up the repository and timeout.
//
// Parameters:
//   - userRepository: Repository for managing user data.
//   - timeout: Duration to set context timeout for operations.
//
// Returns:
//   - An instance of UserService.
func NewUserService(userRepository repository.UserRepository, timeout time.Duration) UserService {
	usersIDs := cmap.New[string]() // Initialize a new concurrent map

	return &userService{
		userRepository: userRepository,
		usersIDs:       usersIDs,
		contextTimeout: timeout,
	}
}

// InsertUser adds a new user to the database.
// It returns the newly created user's ID and any error encountered during insertion.
//
// Parameters:
//   - c: The context for managing request lifetime.
//   - user: The user data to be inserted.
//
// Returns:
//   - The ID of the newly created user and an error if the operation fails.
func (us *userService) InsertUser(c context.Context, user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(c, us.contextTimeout) // Set up context with timeout
	defer cancel()                                           // Ensure cancellation of context when done

	userID, err := us.userRepository.InsertUser(ctx, user) // Call repository method to insert user

	return userID, err // Return the newly created user's ID and any errors
}

// UpdatePassword updates an existing user's password in the database.
//
// Parameters:
//   - c: The context for managing request lifetime.
//   - user: The user data containing the updated password.
//
// Returns:
//   - An error if the operation fails; otherwise, nil.
func (us *userService) UpdatePassword(c context.Context, user models.User) error {
	ctx, cancel := context.WithTimeout(c, us.contextTimeout) // Set up context with timeout
	defer cancel()                                           // Ensure cancellation of context when done

	err := us.userRepository.UpdatePassword(ctx, user) // Call repository method to update password

	return err // Return any errors from the repository
}

// UpdateRefreshToken updates an existing user's refresh token in the database.
//
// Parameters:
//   - c: The context for managing request lifetime.
//   - user: The user data containing the updated refresh token.
//
// Returns:
//   - An error if the operation fails; otherwise, nil.
func (us *userService) UpdateRefreshToken(c context.Context, user models.User) error {
	ctx, cancel := context.WithTimeout(c, us.contextTimeout) // Set up context with timeout
	defer cancel()                                           // Ensure cancellation of context when done

	err := us.userRepository.UpdateRefreshToken(ctx, user) // Call repository method to update refresh token

	return err // Return any errors from the repository
}

// DeleteUser removes a user's account from the database.
//
// Parameters:
//   - c: The context for managing request lifetime.
//   - userID: The ID of the user to be deleted.
//
// Returns:
//   - An error if the operation fails; otherwise, nil.
func (us *userService) DeleteUser(c context.Context, userID int) error {
	ctx, cancel := context.WithTimeout(c, us.contextTimeout) // Set up context with timeout
	defer cancel()                                           // Ensure cancellation of context when done

	err := us.userRepository.DeleteUser(ctx, userID) // Call repository method to delete user

	return err // Return any errors from the repository
}

// GetUserById retrieves a user's information by their ID from the database.
//
// Parameters:
//   - c: The context for managing request lifetime.
//   - userID: The ID of the user to retrieve.
//
// Returns:
//   - A User object and an error if any occurs during retrieval.
func (us *userService) GetUserById(c context.Context, userID int) (models.User, error) {
	ctx, cancel := context.WithTimeout(c, us.contextTimeout) // Set up context with timeout
	defer cancel()                                           // Ensure cancellation of context when done

	user, err := us.userRepository.GetUserById(ctx, userID) // Call repository method to get user by ID

	return user, err // Return retrieved User object and any errors
}

// GetUserByEmail retrieves a user from the repository by their email address.
// It takes a context and an email string as parameters and returns the User object and any error encountered.
func (us *userService) GetUserByEmail(c context.Context, email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(c, us.contextTimeout) // Create a context with a timeout duration defined in the service
	defer cancel()                                           // Ensure that the context is cancelled when this function exits to free up resources

	user, err := us.userRepository.GetUserByEmail(ctx, email) // Fetch user from repository using the provided email

	return user, err // Return the user and error (if any) to the caller
}

// GetUsersIDFromMemory returns a concurrent map containing all users' IDs stored in memory.
//
// Returns:
//   - A concurrent map containing users' IDs.
func (us *userService) GetUsersIdFromMemory() cmap.ConcurrentMap[string, string] {
	return us.usersIDs // Return the concurrent map of users' IDs
}

// SetUserIdIntoMemory adds a new user ID to the in-memory storage.
//
// Parameters:
//   - userID: The ID of the user to be added.
//
// It uses the concurrent map to store the ID.
func (us *userService) SetUserIdIntoMemory(userID int) {
	us.usersIDs.Set(strconv.Itoa(userID), strconv.Itoa(userID))
}

// DeleteUserIdFromMemory removes a user's ID from the in-memory storage.
//
// Parameters:
//   - userID: The ID of the user to be removed.
//
// It uses the concurrent map to remove the ID.
func (us *userService) DeleteUserIdFromMemory(userID int) {
	us.usersIDs.Remove(strconv.Itoa(userID))
}

// GetUsersIdFromDB retrieves all users' IDs from the database and stores them in memory.
//
// Parameters:
//   - c: The context for managing request lifetime.
//
// Returns:
//   - An error if any occurs during retrieval.
func (us *userService) GetUsersIdFromDB(c context.Context) error {
	ctx, cancel := context.WithTimeout(c, us.contextTimeout) // Set up context with timeout
	defer cancel()                                           // Ensure cancellation of context when done

	allIDs, err := us.userRepository.GetAllIDs(ctx) // Call repository method to get all IDs from DB

	// Fill all users' ID storage in memory
	allUsersIDs := allIDs
	for _, id := range allUsersIDs {
		idString := strconv.Itoa(id) // Convert integer ID to string

		us.usersIDs.Set(idString, idString) // Store ID in concurrent map
	}

	return err // Return any errors encountered during retrieval
}
