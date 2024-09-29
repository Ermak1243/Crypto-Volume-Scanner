package route

import (
	"cvs/api/server/controller" // Importing the controller package for handling requests
	"cvs/api/server/middleware" // Importing middleware for request authentication
	"cvs/internal/service"      // Importing service layer for business logic
	"cvs/internal/service/exchange"

	"github.com/gofiber/fiber/v2" // Importing Fiber framework for web server
)

// NewUserRouter sets up the user-related routes for the application.
//
// This function creates a new router group for user operations and defines the following routes:
//
// 1. **Authentication Routes**:
//   - POST /api/auth/signup: Endpoint for user registration.
//   - POST /api/auth/login: Endpoint for user login.
//   - GET /api/auth/tokens: Endpoint to retrieve tokens, requires authentication.
//
// 2. **User Management Routes**:
//   - PUT /api/user/update-password: Endpoint to update the user's password, requires authentication.
//   - DELETE /api/user/: Endpoint to delete the user's account, requires authentication.
//
// Parameters:
//   - group: A Fiber router group for organizing user-related routes.
//   - userService: A service responsible for user-related operations.
//   - jwtService: A service responsible for handling JWT operations.
func NewUserRouter(
	group fiber.Router,
	userService service.UserService,
	jwtService service.JwtService,
	allExchangesStorage exchange.AllExchanges,
) {
	uc := controller.NewUserController(userService, jwtService, allExchangesStorage) // Create a new instance of UserController

	authRoutes := group.Group("/auth")                                                        // Create a sub-group for authentication routes
	authRoutes.Post("/signup", uc.Signup)                                                     // Route for user signup
	authRoutes.Post("/login", uc.Login)                                                       // Route for user login
	authRoutes.Get("/tokens", middleware.IsAuthenticated(jwtService, userService), uc.Tokens) // Route to get tokens with authentication

	group.Put("/update-password", middleware.IsAuthenticated(jwtService, userService), uc.UpdatePassword) // Route to update password with authentication
	group.Delete("", middleware.IsAuthenticated(jwtService, userService), uc.DeleteUser)                  // Route to delete user account with authentication
}
