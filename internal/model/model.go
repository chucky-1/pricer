// Package model has struct of essence
package model

// Price is struct
type Price struct {
	ID   int32   `validate:"required"`
	Bid  float32 `validate:"required,gt=0"`
	Ask  float32 `validate:"required,gt=0"`
	Time string  `validate:"required"`
}
