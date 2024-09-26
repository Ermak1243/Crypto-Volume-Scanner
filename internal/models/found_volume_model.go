package models

import "time"

type FoundVolume struct {
	Exchange        string    `json:"exchange"`
	Pair            string    `json:"pair"`
	Price           float64   `json:"price"`
	Index           int       `json:"index"`      // Number of rows between found volume index and best ask or best bid and found volume index
	Difference      float64   `json:"difference"` // Difference between found volume and best ask or best bid and found volume in percent
	Volume          float64   `json:"volume"`
	VolumeTimeFound time.Time `json:"volume_time_found"`
	Side            string    `json:"side"`
}
