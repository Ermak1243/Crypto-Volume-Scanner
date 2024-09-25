package controller

import (
	"context"
	"math/rand"
	"net/http"

	"main/internal/domain/models" // Importing the models package for user data structures
	"main/internal/service"       // Importing the service package for user and JWT services
	"main/internal/service/exchange"

	"github.com/gofiber/fiber/v2" // Importing the Fiber framework for building web applications
)

// userController handles user-related operations.
type userController struct {
	userService         service.UserService   // Service for managing user data
	allExchangesStorage exchange.AllExchanges // Storage for all exchanges
	jwtService          service.JwtService    // Service for managing JWT tokens
}

// NewUserController creates a new instance of userController.
// It takes a UserService and a JwtService as dependencies.
//
// Parameters:
//   - userService: A service for managing user data.
//   - jwtService: A service for managing JWT tokens.
//
// Returns:
//   - A pointer to a new userController instance.
func NewUserController(userService service.UserService, jwtService service.JwtService, allExchangesStorage exchange.AllExchanges) *userController {
	return &userController{
		userService:         userService,
		allExchangesStorage: allExchangesStorage,
		jwtService:          jwtService,
	}
}

// Signup handles the user registration process by parsing the incoming request,
// validating the user data, and inserting the new user into the database.
// It also generates access and refresh tokens for the newly registered user.
//
// The function performs the following steps:
// 1. Initializes a struct to hold new user data.
// 2. Parses the request body into the `newUserData` struct.
// 3. Creates a `User` object from the parsed email.
// 4. Sets the user's password and handles any errors that may occur.
// 5. Validates the user data (e.g., email format).
// 6. Attempts to insert the new user into the database and retrieves the user ID.
// 7. Generates access and refresh tokens for the newly created user.
// 8. Sets the refresh token for the user object and updates it in the database.
// 9. Returns a JSON response containing tokens data if successful, or an error message if any step fails.
//
// @Summary Sign up a new user
// @Description Create a new user account with email and password.
// @Description Returns the access token, the refresh token, and the time when the access token ceases to be valid. After the access token has ceased to be valid, you need to send a request along the path "/api/user/auth/token" to get a new pair of tokens.
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UserAuth true "User registration data"
// @Success 200 {object} models.Tokens "Successful response with tokens data"
// @Failure 400 {object} models.Response "Invalid input data"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/user/auth/signup [post]
func (uc *userController) Signup(c *fiber.Ctx) error {
	newUserData := models.UserAuth{} // Initialize a struct to hold new user data

	c.Status(http.StatusBadRequest) // Set response status to Bad Request initially

	// Parse the request body into the newUserData struct
	if err := c.BodyParser(&newUserData); err != nil {
		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format if parsing fails
		})
	}

	// Create a User object from the parsed email
	user := models.User{
		Email: newUserData.Email,
	}

	// Set the user's password using the provided password and handle any errors
	if err := user.SetPassword(newUserData.Password); err != nil {
		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format if setting password fails
		})
	}

	// Validate the user data (e.g., email format, etc.)
	if err := service.CheckUserData(user); err != nil {
		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format if validation fails
		})
	}

	c.Status(http.StatusInternalServerError) // Set response status to Internal Server Error for potential database issues

	user.SessionID = 1 // Set the user's session ID to an intermediate value

	// Insert the new user into the database and retrieve the user ID
	userId, err := uc.userService.InsertUser(c.Context(), user)
	if err != nil {
		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format if insertion fails
		})
	}

	user.ID = userId

	tokensData, err := uc.updateTokens(user)
	if err != nil {
		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format if updating refresh token fails
		})
	}

	return c.Status(http.StatusOK).JSON(tokensData) // Return tokens data in JSON format with a 200 OK status
}

// Tokens handles the refresh token operation for an authenticated user.
// It retrieves the refresh token from the request header and validates it.
// If valid, it generates new access and refresh tokens for the user.
//
// This method performs the following steps:
// 1. Retrieves the user object from the context, which was set during authentication.
// 2. Extracts the refresh token from the Authorization header of the request.
// 3. Validates the provided refresh token against the stored token for the user.
// 4. If validation is successful, generates new access and refresh tokens for the user.
// 5. Returns the newly generated tokens in JSON format upon successful operation.
//
// @Summary Get new tokens
// @Description Retrieve new access and refresh tokens for the authenticated user
// @Tags users
// @Param Authorization header string true "Refresh token"
// @Success 200 {object} models.Tokens "Successful response with new tokens"
// @Failure 401 {object} models.Response "Invalid refresh token"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/user/auth/tokens [get]
func (uc *userController) Tokens(c *fiber.Ctx) error {
	// Retrieve the user object from the context, which was set during authentication.
	user := c.Locals("user").(models.User)

	// Get the refresh token from the request header (Authorization header).
	refreshToken := c.Get("Authorization")

	// Compare the provided refresh token with the one stored for the user.
	err := user.CompareRefreshToken(refreshToken)
	if err != nil {
		c.Status(http.StatusUnauthorized) // Set response status to Unauthorized (401)

		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format
		})
	}

	newTokens, err := uc.updateTokens(user)
	if err != nil {
		c.Status(http.StatusInternalServerError)

		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format if updating refresh token fails
		})
	}

	return c.JSON(newTokens) // Return new tokens in JSON format with a 200 OK status
}

// Login handles user authentication by processing the login request.
// It expects a JSON body containing the user's email and password.
//
// This method performs the following steps:
// 1. Parses the incoming request body to extract user credentials (email and password).
// 2. Retrieves the user from the database using the provided email address.
// 3. Validates the user's existence and checks if the provided password matches the stored password.
// 4. If authentication is successful, it generates new access and refresh tokens for the user.
// 5. Returns the newly generated tokens in JSON format.
//
// @Summary Log in a user
// @Description Authenticate a user and issue tokens if successful
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UserAuth true "User login data"
// @Success 200 {object} models.Tokens "New tokens data"
// @Failure 400 {object} models.Response "Invalid input data"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/user/auth/login [post]
func (uc *userController) Login(c *fiber.Ctx) error {
	userDataRequest := models.UserAuth{} // Initialize a struct to hold user credentials

	c.Status(http.StatusBadRequest) // Set response status to Bad Request initially

	// Parse the request body into the userDataRequest struct
	if err := c.BodyParser(&userDataRequest); err != nil {
		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format if parsing fails
		})
	}

	// Retrieve the user from the database using their email
	userFromDB, err := uc.userService.GetUserByEmail(c.Context(), userDataRequest.Email)
	if err != nil || userFromDB.Email != userDataRequest.Email {
		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format if user not found or email mismatch
		})
	}

	// Compare the provided password with the stored password for the user
	if err := userFromDB.ComparePassword(userDataRequest.Password); err != nil {
		return c.JSON(models.Response{
			Result: "invalid password", // Return error message in JSON format if password is invalid
		})
	}

	newTokens, err := uc.updateTokens(userFromDB)
	if err != nil {
		c.Status(http.StatusInternalServerError)

		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format if updating refresh token fails
		})
	}

	return c.Status(http.StatusOK).JSON(newTokens) // Return new tokens in JSON format with a 200 OK status
}

// UpdatePassword handles the request to update a user's password.
// It expects a JSON body containing the old and new passwords.
//
// This method performs the following steps:
// 1. Parses the incoming request body to extract the old and new passwords.
// 2. Retrieves the authenticated user object from the context.
// 3. Validates the provided old password against the stored password.
// 4. If validation is successful, it sets the new password and refresh token for the user.
// 5. Updates the user's password in the database and generates new access and refresh tokens.
// 6. Returns the newly generated tokens in JSON format upon successful update.
//
// @Summary Update user password
// @Description Update the password for the authenticated user.
// @Description Returns the access token, the refresh token, and the time when the access token ceases to be valid. After the access token has ceased to be valid, you need to send a request along the path "/api/user/auth/token" to get a new pair of tokens.
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Access token"
// @Param passwords body models.PasswordUpdate true "Passwords data"
// @Success 200 {object} models.Tokens "New tokens data"
// @Failure 400 {object} models.Response "Invalid password"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/user/update-password [put]
func (uc *userController) UpdatePassword(c *fiber.Ctx) error {
	passwordData := models.PasswordUpdate{} // Initialize a struct to hold password update data

	c.Status(http.StatusBadRequest) // Set response status to Bad Request initially

	// Parse the request body into the passwordData struct
	if err := c.BodyParser(&passwordData); err != nil {
		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format if parsing fails
		})
	}

	// Retrieve the user object from the context locals, which was set during authentication
	user := c.Locals("user").(models.User)

	// Compare the provided old password with the stored password for validation
	if err := user.ComparePassword(passwordData.OldPassword); err != nil {
		return c.JSON(models.Response{
			Result: "invalid old password", // Return error message in JSON format if old password is invalid
		})
	}

	c.Status(http.StatusInternalServerError) // Set response status to Internal Server Error (500)

	// Generate new access and refresh tokens for the user after updating their password
	newTokens, sessionId, err := uc.generateTokens(user.ID)
	if err != nil {
		return c.JSON(models.Response{
			Result: "user update failed", // Return error message in JSON format if token generation fails
		})
	}

	// Set the new refresh token in the user object
	user.SetRefreshToken(newTokens.Refresh)
	// Set the new password in the user object
	user.SetPassword(passwordData.NewPassword)
	user.SessionID = sessionId

	// Update the user's password in the database
	err = uc.userService.UpdatePassword(c.Context(), user)
	if err != nil {
		return c.JSON(models.Response{
			Result: "user update failed", // Return error message in JSON format if updating password fails
		})
	}

	return c.Status(http.StatusOK).JSON(newTokens) // Return new tokens in JSON format with a 200 OK status
}

// DeleteUser handles the request to delete a user's account.
// It retrieves the authenticated user from the context and deletes their account from the database.
//
// This method performs the following steps:
// 1. Retrieves the user object from the context locals, which was set during authentication.
// 2. Attempts to delete the user's account using their ID.
// 3. If successful, returns a success message; otherwise, returns an error message.
//
// @Summary Delete a user account
// @Description Delete the authenticated user's account
// @Tags users
// @Produce json
// @Param Authorization header string true "Access token"
// @Success 200 {object} models.Response "Successful response"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/user [delete]
func (uc *userController) DeleteUser(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User) // Retrieve user ID from context locals

	// Delete the user's account from the database using their ID.
	err := uc.userService.DeleteUser(c.Context(), user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError) // Set response status to Internal Server Error

		return c.JSON(models.Response{
			Result: "user deletion failed", // Return error message in JSON format
		})
	}

	uc.userService.DeleteUserIdFromMemory(user.ID)

	// Iterate over all exchanges and clear their subscribed pairs storage
	for _, exchange := range uc.allExchangesStorage.All() {
		exchange.ClearSubscribedPairsStorage()
	}

	return c.JSON(models.Response{
		Result: "user deleted successfully", // Return success message in JSON format
	})
}

// generateTokens generates new access and refresh tokens for a user.
//
// This method creates a random session ID for each token generation process,
// then generates an access token and a refresh token using the provided user ID.
//
// Parameters:
//   - userId: An integer representing the user's unique identifier.
//
// Returns:
//   - models.Tokens: A structure containing the newly generated access token,
//     refresh token, and expiration time of the access token.
//   - int: A randomly generated session ID associated with this token generation.
//   - error: An error if there was an issue creating either of the tokens;
//     if successful, it returns nil.
//
// Possible Errors:
//   - An error may occur during the creation of either the access or refresh tokens,
//     in which case it will be returned alongside an empty Tokens structure.
func (uc *userController) generateTokens(userId int) (models.Tokens, int, error) {
	// Generate a random session ID for this token generation process.
	sessionId := rand.Intn(9999)

	// Create an access token using the user ID and session ID.
	accessToken, expiresAt, err := uc.jwtService.CreateAccessToken(userId, sessionId)
	if err != nil {
		return models.Tokens{}, 0, err // Return an empty Tokens struct and error if token creation fails.
	}

	// Create a refresh token using the user ID and session ID.
	refreshToken, err := uc.jwtService.CreateRefreshToken(userId, sessionId)
	if err != nil {
		return models.Tokens{}, 0, err // Return an empty Tokens struct and error if token creation fails.
	}

	// Return a Tokens struct containing the generated access token,
	// refresh token, and expiration time.
	return models.Tokens{
		Access:    accessToken,
		Refresh:   refreshToken,
		ExpiresAt: expiresAt,
	}, sessionId, nil // Return session id and nil indicating no error occurred.
}

// updateRefreshToken updates the access and refresh tokens for the specified user.
//
// This function generates new access and refresh tokens for the authenticated user,
// sets the new refresh token in the user object, and updates it in the database.
//
// Parameters:
//   - user: A models.User structure representing the user for whom the tokens need to be updated.
//     The object must contain a valid user ID.
//
// Returns:
//   - models.Tokens: A structure containing the new access and refresh tokens.
//   - error: An error if there was an issue generating tokens or updating the token in the database.
//     If the function completes successfully, the returned error value will be nil.
//
// Possible Errors:
//   - An error may occur during token generation if the user does not exist or if there is
//     an issue creating the tokens.
//   - An error may occur when setting the refresh token in the user object.
//   - An error may occur when attempting to update the token in the database.
func (uc *userController) updateTokens(user models.User) (models.Tokens, error) {
	// Generate new access and refresh tokens for the authenticated user
	newTokens, sessionId, err := uc.generateTokens(user.ID)
	if err != nil {
		return models.Tokens{}, err // Return an empty Tokens struct and error if token generation fails
	}

	// Set the refresh token for the user object
	if err := user.SetRefreshToken(newTokens.Refresh); err != nil {
		return models.Tokens{}, err // Return an empty Tokens struct and error if setting the refresh token fails
	}
	user.SessionID = sessionId // Assign the new session ID to the user

	// Update the user's refresh token in the database
	err = uc.userService.UpdateRefreshToken(context.Background(), user)

	return newTokens, err // Return the new tokens and any error from updating the database
}
