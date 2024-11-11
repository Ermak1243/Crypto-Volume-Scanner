/*
Package middleware provides HTTP middleware functions for a Fiber web application.

This package contains functions to set up various middlewares that enhance the functionality,
security, and performance of the application. Middlewares are functions that intercept HTTP requests
and responses, allowing for processing such as logging, authentication, rate limiting, and CORS handling.

The following middlewares are configured in this package:

  - CORS Middleware: Manages Cross-Origin Resource Sharing settings to control which origins can access resources.
  - Logger Middleware: Logs incoming requests and responses to a specified log file for monitoring and debugging.
  - Rate Limiter Middleware: Limits the number of requests from a single IP address to prevent abuse and ensure fair usage.

The middleware functions included in this package are:

 1. **MiddlewaresSetup**: Configures and applies the necessary middlewares to the provided Fiber application instance.
 2. **IsAuthenticated**: A middleware that checks if the user is authenticated using JSON Web Tokens (JWT). It verifies the presence and validity of the JWT in the Authorization header.

Example usage of this package can be seen in the main application file where these middlewares are applied to the Fiber app instance.
*/
package middleware

import (
	"cvs/internal/models"  // Importing models for data structures
	"cvs/internal/service" // Importing service layer for business logic
	"net/http"

	"github.com/gofiber/fiber/v2"                    // Importing Fiber framework
	"github.com/gofiber/fiber/v2/middleware/cors"    // Importing CORS middleware
	"github.com/gofiber/fiber/v2/middleware/limiter" // Importing rate limiting middleware
	"github.com/gofiber/fiber/v2/middleware/logger"  // Importing logging middleware
)

// MiddlewaresSetup configures and applies various middlewares to the provided Fiber application.
//
// This function sets up the following middlewares:
//
// 1. CORS Middleware:
//   - Allows specific HTTP methods (POST, GET, DELETE, PUT) for cross-origin requests.
//   - Specifies allowed headers (Accept, Accept-Language, Content-Type) in requests.
//
// 2. Logger Middleware:
//   - Logs incoming requests and responses to a specified log file.
//   - The log output can be customized by modifying the logger configuration.
//
// 3. Rate Limiter Middleware:
//   - Limits the maximum number of requests per IP address to prevent abuse.
//   - Configured to allow a maximum of 1000 requests from a single IP address.
//
// Parameters:
//   - server *fiber.App: The Fiber application instance to which the middlewares will be applied.
//
// Example Usage:
//
//	func main() {
//	    app := fiber.New()
//	    middleware.MiddlewaresSetup(app)
//	    app.Listen(":3000")
//	}
func MiddlewaresSetup(server *fiber.App) {
	server.Use(
		cors.New(cors.Config{
			AllowMethods: "POST, GET, DELETE, PUT",                // Specify allowed HTTP methods
			AllowHeaders: "Accept, Accept-Language, Content-Type", // Specify allowed headers
		}),
		logger.New(),
		limiter.New(limiter.Config{
			Max: 1000, // Set maximum number of requests per IP address
		}),
	)
}

// IsAuthenticated is a middleware that checks if the user is authenticated using JWT.
//
// This middleware retrieves the JWT from the Authorization header and validates it by parsing
// the token to extract user ID and session ID. It then checks if the user exists in the database
// and whether the session ID matches. If authentication is successful, it stores the user information
// in context locals for later use; otherwise, it returns an error response.
//
// Parameters:
//   - jwtService service.JwtService: The service responsible for parsing JWT tokens.
//   - userService service.UserService: The service responsible for user-related operations.
//
// Returns:
//   - fiber.Handler: A Fiber handler function that performs authentication checks.
func IsAuthenticated(jwtService service.JwtService, userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jwt := c.Get("Authorization")     // Retrieve the JWT from the Authorization header
		c.Status(http.StatusUnauthorized) // Set the default response status to Unauthorized

		if jwt == "" {
			return c.JSON(fiber.Map{
				"result": "refresh token is required", // Return error if JWT is missing
			})
		}

		userID, sessionId, errParse := jwtService.Parse(jwt)              // Parse the JWT to extract user ID and session ID
		userFromDB, errDB := userService.GetUserById(c.Context(), userID) // Fetch user from database using user ID

		if errParse != nil || errDB != nil || userID < 1 || sessionId < 1 {
			return c.JSON(models.Response{
				Result: "user not found", // Return error message if user is not found or parsing fails
			})
		}

		if sessionId != userFromDB.SessionID {
			return c.JSON(models.Response{
				Result: "invalid token", // Return error if session ID does not match
			})
		}

		c.Locals("user", userFromDB) // Store the authenticated user in context locals for later use
		c.Status(http.StatusOK)

		return c.Next() // Proceed to the next middleware or handler
	}
}
