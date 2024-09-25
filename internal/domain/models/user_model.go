package models

import (
	"errors"
	"time"

	"github.com/matthewhartstonge/argon2"
)

var argon = argon2.DefaultConfig()

type User struct {
	ID           int
	SessionID    int `db:"session_id"`
	Email        string
	RefreshToken []byte `db:"refresh_token"`
	Password     []byte
	CreatedAt    time.Time `json:"-" db:"created_at" default:"now()" `
	UpdatedAt    time.Time `json:"-" db:"updated_at" default:"now()"`
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := argon.HashEncoded([]byte(password))
	u.Password = hashedPassword

	return err
}

func (u *User) SetRefreshToken(refreshToken string) error {
	hashedToken, err := argon.HashEncoded([]byte(refreshToken))
	u.RefreshToken = hashedToken

	return err
}

func (u *User) ComparePassword(password string) error {
	ok, err := argon2.VerifyEncoded([]byte(password), u.Password)
	if !ok {
		return errors.New("comparison passwords failed")
	}

	return err
}

func (u *User) CompareRefreshToken(refreshToken string) error {
	ok, err := argon2.VerifyEncoded([]byte(refreshToken), u.RefreshToken)
	if !ok {
		return errors.New("comparison refresh tokens failed")
	}

	return err
}
