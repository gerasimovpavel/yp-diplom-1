package model

type Balance struct {
	Сurrent   float64 `json:"current" db:"current"`
	WithDrawn float64 `json:"with_drawn" db:"with_drawn"`
}
