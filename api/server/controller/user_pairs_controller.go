package controller

import (
	"log"
	"net/http"

	"main/internal/models"
	"main/internal/service"
	"main/internal/service/exchange"

	"github.com/gofiber/fiber/v2"
)

// userPairsController handles operations related to user pairs.
type userPairsController struct {
	userPairsService    service.UserPairsService    // Service for managing user pairs
	userService         service.UserService         // Service for managing users data
	foundVolumesService service.FoundVolumesService // Service for managing found volumes
	allExchangesStorage exchange.AllExchanges       // Storage for all exchanges
}

// NewUserPairsController creates a new instance of userPairsController.
//
// This function initializes a userPairsController with the necessary services for managing user pairs
// and found volumes. It sets up the controller to handle requests related to user trading pairs
// and their associated volumes.
//
// Parameters:
//   - userPairsService: The service for managing user pairs data.
//   - foundVolumesService: The service for managing found volumes data.
//   - allExchangesStorage: The storage for all exchanges, allowing access to exchange-related operations.
//
// Returns:
//   - *userPairsController: A pointer to the initialized userPairsController instance.
func NewUserPairsController(
	userPairsService service.UserPairsService,
	userService service.UserService,
	foundVolumesService service.FoundVolumesService,
	allExchangesStorage exchange.AllExchanges,
) *userPairsController {
	return &userPairsController{
		userPairsService:    userPairsService,
		userService:         userService,
		foundVolumesService: foundVolumesService,
		allExchangesStorage: allExchangesStorage,
	}
}

// Add creates a new user pair in the database.
// It retrieves the authenticated user's ID from the context,
// parses the request body to obtain the new pair data,
// and calls the service to perform the addition.
//
// The function performs the following steps:
// 1. Initializes a `UserPairs` struct to hold the new pair data.
// 2. Retrieves the authenticated user's ID from context locals.
// 3. Parses the request body into the `pairData` struct.
// 4. Calls the service to add the new pair to the database.
// 5. Returns a JSON response indicating success or failure.
//
// @Summary Add a new user pair
// @Description Create a new pair for the authenticated user
// @Tags user-pairs
// @Accept json
// @Produce json
// @Param Authorization header string true "Access token"
// @Param pair body models.UserPairs true "User pair data"
// @Success 200 {object} models.Response "Successful response indicating the pair was added"
// @Failure 400 {object} models.Response "Invalid input data"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/user/pair/add [post]
func (uc *userPairsController) Add(c *fiber.Ctx) error {
	const op = directoryPath + "user_controller.Add"

	var pairData models.UserPairs                       // Initialize a UserPairs struct to hold the new pair data
	pairData.UserID = c.Locals("user").(models.User).ID // Retrieve authenticated user's ID from context locals

	// Parse the request body into pairData
	if err := c.BodyParser(&pairData); err != nil {
		log.Println(op, err)

		c.Status(http.StatusBadRequest)

		return c.JSON(models.Response{
			Result: "invalid input data", // Return error if parsing fails
		})
	}

	// Call the service to add the new pair to the database
	if err := uc.userPairsService.Add(c.Context(), pairData); err != nil {
		log.Println(op, err)

		c.Status(http.StatusInternalServerError)

		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format
		})
	}

	uc.userService.SetUserIdIntoMemory(pairData.UserID)

	exchange := uc.allExchangesStorage.Get(pairData.Exchange)
	exchange.AddPairToSubscribedPairs(pairData.Pair)

	return c.JSON(models.Response{
		Result: "pair added successfully",
	}) // Return success message in JSON format
}

// UpdateExactValue updates an existing user pair in the database.
// It retrieves the authenticated user's ID from the context,
// parses the request body to obtain the updated pair data,
// and calls the service to perform the update.
//
// The function performs the following steps:
// 1. Initializes a `UserPairs` struct to hold the updated pair data.
// 2. Retrieves the authenticated user's ID from context locals.
// 3. Parses the request body into the `pairData` struct.
// 4. Calls the service to update the existing pair in the database.
// 5. Returns a JSON response indicating success or failure.
//
// @Summary Update the exact value of a user pair
// @Description Update an existing pair for the authenticated user
// @Tags user-pairs
// @Accept json
// @Produce json
// @Param Authorization header string true "Access token"
// @Param pair body models.UserPairs true "User pair data"
// @Success 200 {object} models.Response "Successful response indicating the pair was updated"
// @Failure 400 {object} models.Response "Invalid input data"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/user/pair/update-exact-value [put]
func (uc *userPairsController) UpdateExactValue(c *fiber.Ctx) error {
	const op = directoryPath + "user_controller.UpdateExactValue"

	var pairData models.UserPairs                       // Initialize a UserPairs struct to hold the updated pair data
	pairData.UserID = c.Locals("user").(models.User).ID // Retrieve authenticated user's ID from context locals

	// Parse the request body into pairData
	if err := c.BodyParser(&pairData); err != nil {
		log.Println(op, err)

		c.Status(http.StatusBadRequest)

		return c.JSON(models.Response{
			Result: "invalid input data", // Return error if parsing fails
		})
	}

	// Call the service to update the existing pair in the database
	if err := uc.userPairsService.UpdateExactValue(c.Context(), pairData); err != nil {
		log.Println(op, err)

		c.Status(http.StatusInternalServerError)

		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format
		})
	}

	return c.JSON(models.Response{
		Result: "pair updated successfully",
	}) // Return success message in JSON format
}

// GetAllUserPairs retrieves all user pairs associated with the authenticated user.
// It fetches the user's ID from the context and calls the service to get all pairs.
//
// The function performs the following steps:
// 1. Retrieves the authenticated user's ID from context locals.
// 2. Calls the service to get all pairs associated with the user's ID.
// 3. Returns a JSON response containing the list of user pairs or an error message.
//
// @Summary Retrieve all pairs for the authenticated user
// @Description Get all user pairs associated with the authenticated user's account
// @Tags user-pairs
// @Produce json
// @Param Authorization header string true "Access token"
// @Success 200 {array} models.UserPairs "List of user pairs"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/user/pair/all-pairs [get]
func (uc *userPairsController) GetAllUserPairs(c *fiber.Ctx) error {
	const op = directoryPath + "user_controller.GetAllUserPairs"

	userID := c.Locals("user").(models.User).ID // Retrieve authenticated user's ID from context locals

	// Call the service to get all pairs associated with the authenticated user's ID
	userPairs, err := uc.userPairsService.GetAllUserPairs(c.Context(), userID)
	if err != nil {
		log.Println(op, err)

		c.Status(http.StatusInternalServerError)

		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format
		})
	}

	return c.JSON(userPairs) // Return list of user pairs in JSON format
}

// GetAllUserFoundVolumes retrieves all found volumes associated with the authenticated user.
//
// This method extracts the user's ID from the context locals and calls the
// foundVolumesService to fetch all found volumes related to that user.
// If an error occurs during this process, it returns an appropriate error message.
//
// Parameters:
//   - c: A pointer to fiber.Ctx, which contains information about the HTTP request
//     and response, including context locals.
//
// Returns:
//   - error: If an error occurs while fetching found volumes, it returns an error
//     indicating that the operation failed. If successful, it returns nil.
//
// Possible Responses:
//   - On success, it returns a JSON response containing a list of found volumes
//     associated with the authenticated user.
//   - If an error occurs during retrieval, it sets the HTTP status to 500 (Internal Server Error)
//     and returns a JSON response containing the error message.
//
// @Summary Retrieve all found volumes for the authenticated user
// @Description This endpoint retrieves a list of all found volumes associated with the authenticated user.
// @Tags user-pairs
// @Accept json
// @Produce json
// @Param Authorization header string true "Access token"
// @Success 200 {array} models.FoundVolume "Success"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Router /api/user/pair/found-volumes [get]
func (uc *userPairsController) GetAllUserFoundVolumes(c *fiber.Ctx) error {
	const op = directoryPath + "user_controller.GetAllUserFoundVolumes"

	userID := c.Locals("user").(models.User).ID // Retrieve authenticated user's ID from context locals

	// Call the service to get all pairs associated with the authenticated user's ID
	foundVolumes, err := uc.foundVolumesService.GetAllFoundVolume(userID)
	if err != nil {
		log.Println(op, err)

		c.Status(http.StatusInternalServerError)

		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format
		})
	}

	return c.JSON(foundVolumes) // Return list of user pairs in JSON format
}

// DeletePair handles the HTTP request to delete a user pair from the database.
//
// This method retrieves the pair identifier from the query parameters and
// the authenticated user from the context. It then constructs a UserPairs
// object with the user ID and pair information. The function calls the
// userPairsService to perform the deletion.
//
// Query Parameters:
//   - pair: The identifier of the user pair to be deleted, extracted from the query string.
//
// Parameters:
//   - c: A pointer to fiber.Ctx, which contains information about the HTTP request
//     and response, including parameters and context locals.
//
// Returns:
//   - error: If an error occurs during the deletion process, it returns an error
//     indicating that the operation failed. If successful, it returns nil.
//
// Possible Responses:
//   - On success, it returns a JSON response with a message indicating that
//     the pair was deleted successfully.
//   - If an error occurs during deletion, it sets the HTTP status to 500 (Internal Server Error)
//     and returns a JSON response containing the error message.
//
// @Summary Delete a user pair
// @Description Remove an existing pair for the authenticated user
// @Tags user-pairs
// @Accept json
// @Produce json
// @Param Authorization header string true "Access token"
// @Param        pair   query      string  true  "The pair that should be deleted"
// @Success 200 {object} models.Response "Successful response indicating the pair was deleted"
// @Failure 400 {object} models.Response "Invalid input data"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/user/pair [delete]
func (uc *userPairsController) DeletePair(c *fiber.Ctx) error {
	const op = directoryPath + "user_controller.DeletePair"

	pair := c.Query("pair")                // Retrieve pair from query string
	user := c.Locals("user").(models.User) // Retrieve authenticated user from context locals

	userPairData := models.UserPairs{
		UserID: user.ID, // Set the UserID field to the authenticated user's ID
		Pair:   pair,    // Set the Pair field to the trading pair retrieved from the query
	}

	// Call the service to delete the specified pair from the database
	if err := uc.userPairsService.DeletePair(c.Context(), userPairData); err != nil {
		log.Println(op, err)

		c.Status(http.StatusInternalServerError) // Set HTTP status to 500 if an error occurs

		return c.JSON(models.Response{
			Result: err.Error(), // Return error message in JSON format
		})
	}

	uc.userService.DeleteUserIdFromMemory(user.ID) // Remove the user's ID from the in-memory storage

	// Iterate over all exchanges and remove the pair from their subscribed pairs
	for _, exchange := range uc.allExchangesStorage.All() {
		exchange.DeletePairFromSubscribedPairs(pair) // Remove the pair from each exchange's subscribed pairs

		userPairData.Exchange = exchange.ExchangeName() // Set the Exchange field to the exchange's name
		uc.foundVolumesService.DeleteFoundVolume(userPairData)
	}

	return c.JSON(models.Response{
		Result: "pair deleted successfully", // Return success message in JSON format
	})
}
