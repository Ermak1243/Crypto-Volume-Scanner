package config

import (
	"os" // Importing the os package for file operations

	"github.com/ilyakaznacheev/cleanenv" // Importing cleanenv for reading configuration from files
)

// PostgresConfig holds the configuration settings for connecting to a PostgreSQL database.
type PostgresConfig struct {
	UserName string `yaml:"db_user"`     // Database username
	Password string `yaml:"db_password"` // Database password
	DbName   string `yaml:"db_name"`     // Name of the database
	Host     string `yaml:"db_host"`     // Host where the database server is located
	Port     string `yaml:"db_port"`     // Port on which the database server is listening
	SslMode  string `yaml:"db_ssl_mode"` // SSL mode for database connection (e.g., "disable", "require")
}

// Config aggregates all configuration settings needed by the application.
type Config struct {
	Postgres                  PostgresConfig `yaml:"postgres"`                     // PostgreSQL configuration
	JwtSecretKey              string         `yaml:"jwt_secret_key"`               // Secret key used for signing JWTs
	LogLevel                  string         `yaml:"log_level"`                    // Logging level
	ServerPort                string         `yaml:"server_port"`                  // Port on which the server will run
	AccessTokenLifetimeHours  int            `yaml:"access_token_lifetime_hours"`  // Lifetime of access tokens in hours
	RefreshTokenLifetimeHours int            `yaml:"refresh_token_lifetime_hours"` // Lifetime of refresh tokens in hours
	ContextTimeout            int            `yaml:"context_timeout"`              // Timeout duration for context operations in seconds
}

// NewConfig creates a new configuration instance by loading settings from a specified path.
// Parameters:
//   - path: The file path to the configuration file (e.g., YAML file).
//
// Returns:
//   - A pointer to a Config instance containing the loaded settings.
func NewConfig(path string) *Config {
	return MustLoadPath(path) // Load configuration using MustLoadPath function
}

// MustLoadPath loads the configuration from the specified file path.
// It panics if there is an error loading the configuration or if the file does not exist.
//
// Parameters:
//   - configPath: The path to the configuration file.
//
// Returns:
//   - A pointer to a Config instance containing the loaded settings.
func MustLoadPath(configPath string) *Config {
	// Check if the configuration file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath) // Panic if the file does not exist
	}

	var cfg Config

	// Read and parse the configuration from the specified path into cfg
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error()) // Panic if there is an error reading the config
	}

	return &cfg // Return a pointer to the loaded Config instance
}
