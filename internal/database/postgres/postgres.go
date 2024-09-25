package postgres

import (
	"context"
	"fmt"
	"log"
	"main/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres interface {
	Migration()
	DB() *sqlx.DB
	CloseDB()
}

type postgres struct {
	db *sqlx.DB
}

func NewPostgresDB(cfg config.PostgresConfig) Postgres {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.UserName, cfg.Password, cfg.DbName)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to PostgresConfig!")

	return &postgres{
		db: db,
	}
}

func (s *postgres) Migration() {
	_, err := s.db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id serial PRIMARY KEY,
			session_id integer NOT NULL CHECK (session_id > 0),  --the session ID is needed to link the access token and the refresh token
			email varchar(255) NOT NULL CHECK (email != ''),
			password bytea NOT NULL,
			refresh_token bytea NOT NULL,
			created_at timestamp DEFAULT now(),
			updated_at timestamp DEFAULT now(),
			UNIQUE (email),
			CONSTRAINT password_not_empty CHECK (octet_length(password) > 0)
		);

		CREATE TABLE IF NOT EXISTS user_pairs (
			user_id integer NOT NULL CHECK (user_id > 0) REFERENCES users(id) ON DELETE CASCADE,
			exchange varchar(255) NOT NULL CHECK (exchange != ''),
			pair varchar(255) NOT NULL CHECK (pair != ''),
			exact_value integer NOT NULL CHECK (exact_value > 0),
			UNIQUE (user_id, exchange, pair)  
		);
	`)
	if err != nil {
		fmt.Println("Migration error! ", err)
	}
}

func (s *postgres) DB() *sqlx.DB {
	return s.db
}

func (s *postgres) CloseDB() {
	err := s.db.Close()
	if err != nil {
		log.Println(err)
	}

	log.Println("Connection to Postgres closed.")
}
