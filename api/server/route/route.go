/*
Package route provides the setup for API routes in the Crypto Volume Finder application.

This package is responsible for defining and organizing the routes of the application using the Fiber web framework.
It groups related routes together and applies necessary middleware for authentication and other functionalities.

The main function in this package is `Setup`, which initializes the API routes and groups them under a common path.

Key functionalities provided by this package include:

1. **API Grouping**: The routes are organized under the `/api` path to separate them from other potential routes in the application.
2. **Documentation Routes**: A dedicated route group for API documentation, making it easier to access and view API specifications.
3. **User Routes**: Routes related to user operations, such as registration, login, and profile management.
4. **User Pairs Routes**: Routes specifically for managing user pairs, which require authentication to access.

The following functions are defined in this package:

- **Setup**: Configures the Fiber application with various route groups and applies middleware for authentication.

Example usage of this package can be seen in the main application file where these routes are applied to the Fiber app instance.
*/

package route

import (
	"cvs/api/server/middleware" // Importing middleware for route protection
	"cvs/internal/service"      // Importing services for business logic
	"cvs/internal/service/exchange"
	"cvs/internal/service/logger"

	"github.com/gofiber/fiber/v2" // Importing Fiber framework
)

// @title           Crypto Volume Finder API
// @version         1.0

// Setup initializes the API routes and middleware for the Fiber application.
//
// This function sets up the following route groups:
//
// 1. **Documentation Route Group**:
//   - Sets up a route group for API documentation under `/docs`.
//
// 2. **User Route Group**:
//   - Sets up a route group for user-related operations under `/user`.
//   - Applies JWT authentication middleware to protect user-related routes.
//
// 3. **User Pairs Route Group**:
//   - Sets up a nested route group under `/user/pairs` for managing user pairs,
//   - Requires authentication via JWT middleware.
//
// Parameters:
//   - fiber *fiber.App: The Fiber application instance to which the routes will be applied.
//   - userService service.UserService: The service responsible for user-related operations.
//   - userPairsService service.UserPairsService: The service responsible for managing user pairs.
//   - jwtService service.JwtService: The service responsible for handling JWT operations.
//   - foundVolumesService service.FoundVolumesService: The service responsible for managing found volumes.
//   - allExchangesStorage exchange.AllExchanges: The storage for all exchanges, allowing access to exchange-related operations.
//
// Example Usage:
//
//	func main() {
//	    app := fiber.New()
//	    route.Setup(app, userService, userPairsService, jwtService)
//	    app.Listen(":3000")
//	}
//
// @Title           Crypto Volume Finder API
// @Version         1.0
func Setup(
	fiber *fiber.App,
	userService service.UserService,
	userPairsService service.UserPairsService,
	jwtService service.JwtService,
	foundVolumesService service.FoundVolumesService,
	allExchangesStorage exchange.AllExchanges,
	logger logger.Logger,
) {
	api := fiber.Group("/api") // Create a new group for API routes

	// Group routes for documentation
	docsRoute := fiber.Group("/docs")
	NewDocsRouter(docsRoute) // Initialize documentation routes

	userRoute := api.Group("/user") // Create a group for user-related routes
	NewUserRouter(
		userRoute,
		userService,
		jwtService,
		allExchangesStorage,
		logger,
	) // Initialize user routes

	userPairsRoute := userRoute.Group("/pair").Use(middleware.IsAuthenticated(jwtService, userService)) // Create a protected group for user pairs
	NewUserPairsRouter(
		userPairsRoute,
		userPairsService,
		userService,
		foundVolumesService,
		allExchangesStorage,
		logger,
	) // Initialize user pairs routes
}
