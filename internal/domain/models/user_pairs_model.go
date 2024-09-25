package models

type UserPairs struct {
	UserID     int     `json:"-" db:"user_id"`
	Exchange   string  `json:"exchange" example:"binance_spot"`
	Pair       string  `json:"pair" example:"BTC/USDT"`
	ExactValue float64 `json:"exact_value" db:"exact_value" example:"3"`
}
