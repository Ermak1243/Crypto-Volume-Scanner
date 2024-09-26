package models

type UserAuth struct {
	Email    string `json:"email" example:"example@example.com"`
	Password string `json:"password" example:"password"`
}
