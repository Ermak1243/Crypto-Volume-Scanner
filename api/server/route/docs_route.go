package route

import (
	_ "cvs/api/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// NewDocsRouter sets up the routes for API documentation.
//
// This function creates a new sub-group for documentation routes and defines the following routes:
//
// 1. **GoDoc Redirect**:
//   - GET /godoc: Endpoint to redirect to the GoDoc page for the application, which is hosted
//     on `http://localhost:6060/pkg/main`.
//
// 2. **Swagger UI**:
//   - GET /swagger/*: Endpoint to serve the Swagger UI for the API, which is hosted
//     on `http://localhost:8000/docs/swagger/*`.
//
// Parameters:
//   - group: A Fiber router group for organizing documentation routes.
func NewDocsRouter(group fiber.Router) {
	group.Get("/godoc", func(c *fiber.Ctx) error {
		return c.Redirect("http://localhost:6060/pkg/main", 302)
	})
	group.Get("/swagger/*", swagger.HandlerDefault)
}
