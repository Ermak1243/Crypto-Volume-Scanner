package route

import (
	"cvs/api/server/controller" // Importing the controller package for handling user pair operations
	"cvs/internal/service"      // Importing service layer for business logic related to user pairs
	"cvs/internal/service/exchange"

	"github.com/gofiber/fiber/v2" // Importing Fiber framework for web server
)

// NewUserPairsRouter sets up the routes related to user pairs for the application.
//
// This function creates a new router group for user pair operations and defines the following routes:
//
// 1. **Add User Pair**:
//   - POST /api/user/pair/add: Endpoint to create a new user pair in the database.
//
// 2. **Update User Pair**:
//   - PUT /api/user/pair/update-exact-value: Endpoint to update an existing user pair in the database.
//
// 3. **Get All User Pairs**:
//   - GET /api/user/pair/all-pairs: Endpoint to retrieve all user pairs associated with the authenticated user.
//
// 4. **Delete User Pair**:
//   - DELETE /api/user/pair: Endpoint to delete a specific user pair from the database.
//
// 5. **Get All User Found Volumes**:
//   - GET /api/user/pair/found-volumes: Endpoint to retrieve all found volumes associated with the authenticated user.
//
// Parameters:
//   - group: A Fiber router group for organizing user pair-related routes.
//   - userPairsService: A service responsible for managing user pairs data.
//   - userService: A service responsible for managing user data.
//   - foundVolumesService: A service responsible for managing found volumes data.
//   - allExchangesStorage: A storage for all exchanges, allowing access to exchange-related operations.
func NewUserPairsRouter(
	group fiber.Router,
	userPairsService service.UserPairsService,
	userService service.UserService,
	foundVolumesService service.FoundVolumesService,
	allExchangesStorage exchange.AllExchanges,
) {
	upc := controller.NewUserPairsController(userPairsService, userService, foundVolumesService, allExchangesStorage) // Create a new instance of UserPairsController

	// Define routes for managing user pairs
	group.Post("/add", upc.Add)                            // Route for adding a new user pair
	group.Put("/update-exact-value", upc.UpdateExactValue) // Route for updating an existing user pair
	group.Get("/all-pairs", upc.GetAllUserPairs)           // Route for retrieving all user pairs
	group.Delete("/", upc.DeletePair)                      // Route for deleting a specific user pair
	group.Get("/found-volumes", upc.GetAllUserFoundVolumes)
}
