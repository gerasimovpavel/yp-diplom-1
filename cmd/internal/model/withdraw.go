package model

import "time"

type Withdraw struct {
	Order       string    `json:"order" db:"order"`
	Sum         float64   `json:"sum" db:"summa"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}
