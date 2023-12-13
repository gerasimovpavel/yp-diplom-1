package model

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	Number     string    `json:"number" db:"number"`
	UserID     uuid.UUID `json:"-" db:"user_id"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
	Status     string    `json:"status" db:"status"`
	Accrual    float64   `json:"accrual,omitempty" db:"accrual"`
}
