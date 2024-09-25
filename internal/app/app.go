package app

import (
	"context"
	"log"
	"main/api/server/route"           // Importing routing setup for the API
	"main/internal/config"            // Importing configuration management
	"main/internal/database/postgres" // Importing PostgreSQL database management
	"main/internal/repository"        // Importing repository interfaces and implementations
	"main/internal/service"           // Importing service layer for business logic
	"main/internal/service/exchange"  // Importing exchange service for trading functionality
	"os"
	"os/signal"
	"time"

	"github.com/goccy/go-json" // Importing JSON encoding/decoding library

	"github.com/gofiber/fiber/v2" // Importing Fiber web framework
)

var (
	ctx     = context.Background() // Background context for database operations
	timeout = 5 * time.Second      // Timeout duration for service operations
)

// Run initializes and starts the application.
func Run() {
	cfg := config.NewConfig("./configs/config.yaml") // Load configuration from YAML file

	// Initialize the PostgreSQL database connection
	postgresStorage := postgres.NewPostgresDB(cfg.Postgres)
	postgresStorage.Migration()     // Run database migrations to set up schema
	defer postgresStorage.CloseDB() // Ensure the database connection is closed when done

	db := postgresStorage.DB() // Get the underlying database connection

	// Initialize repositories for data access
	userPairsRepository := repository.NewUserPairsRepository(db) // User pairs repository for managing user pair data
	userRepository := repository.NewUserRepository(db)           // User repository for managing user data

	// Initialize services that contain business logic
	userPairsService := service.NewUserPairsService(userPairsRepository, timeout)                                                                    // Service for user pairs operations
	userService := service.NewUserService(userRepository, timeout)                                                                                   // Service for user operations
	httpRequestService := service.NewHttpRequestService(timeout)                                                                                     // Service for making HTTP requests
	jwtService := service.NewJwtService(cfg.JwtSecretKey, time.Duration(cfg.AccessTokenLifetimeHours), time.Duration(cfg.RefreshTokenLifetimeHours)) // Service for managing JWT tokens
	foundVolumeService := service.NewFoundVolumesService()                                                                                           // Service with found volumes storage                                                                                        // Service for storing found volumes
	userService.GetUsersIdFromDB(ctx)
	allExchangesStorage := exchange.NewAllExchangesService() // Initialize the AllExchanges service
	// Preload user IDs from the database

	// Initialize exchanges and their services
	exchange.InitAllExchanges(
		userService,
		userPairsService,
		httpRequestService,
		foundVolumeService,
		allExchangesStorage,
	)

	fiber := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,   // Set custom JSON encoder for responses
		JSONDecoder: json.Unmarshal, // Set custom JSON decoder for requests
		Immutable:   true,           // Enable immutable routes (for performance)
	})

	// Setup routes for the Fiber application with provided services
	route.Setup(
		fiber,
		userService,
		userPairsService,
		jwtService,
		foundVolumeService,
		allExchangesStorage,
	)

	// Channel for processing interrupt signals (e.g., Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt) // Listen for interrupt signals to gracefully shut down the server

	go func() {
		<-c // Wait for an interrupt signal
		log.Println("Gracefully shutting down...")
		fiber.Shutdown() // Shutdown the Fiber server gracefully
	}()

	if err := fiber.Listen(cfg.ServerPort); err != nil {
		log.Fatal(err) // Log fatal error if server fails to start
	}
}
