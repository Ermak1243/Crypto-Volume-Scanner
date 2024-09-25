package tests

import (
	"context"
	"fmt"
	"main/internal/config"
	"main/internal/database/postgres"
	"main/internal/service"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable         = "users"
	pairsTable         = "user_pairs"
	contextTimeout     = 5 * time.Second
	contextTimeoutZero = 0 * time.Second
	confPath           = "../configs/tests_config.yaml"
)

var (
	ctx                = context.Background()
	deleteUserQueryRow = fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, usersTable)
	jwtService         = service.NewJwtService("secret_key", 20, 1200)
)

func setupDB() *sqlx.DB {
	cfg := config.NewConfig(confPath)

	postgresStorage := postgres.NewPostgresDB(cfg.Postgres)
	postgresStorage.Migration()

	return postgresStorage.DB()
}

// Helper function to create a user in the users table
func insertUser(db *sqlx.DB, email string, password []byte) (int, error) {
	var userID int
	refreshToken := []byte{}

	query := `INSERT INTO users (email, password, refresh_token, session_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.GetContext(context.Background(), &userID, query, email, password, refreshToken, 1)

	return userID, err
}

func insertUserPair(db *sqlx.DB, userID int, exchange string, pair string, exactValue int) error {
	query := `INSERT INTO user_pairs (user_id, exchange, pair, exact_value) VALUES ($1, $2, $3, $4)`
	_, err := db.ExecContext(context.Background(), query, userID, exchange, pair, exactValue)

	return err
}
