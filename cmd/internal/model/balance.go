package model

type Balance struct {
	Current   float64 `json:"current" db:"current"`
	Accrual   float64 `json:"-" db:"accrual"`
	WithDrawn float64 `json:"withdrawn" db:"withdraw"`
}
