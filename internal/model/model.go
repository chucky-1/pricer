// Package model has struct of essence
package model

import "github.com/google/uuid"

// Price is struct
type Price struct {
	ID   uuid.UUID `validate:"required"`
	Bid  float32   `validate:"required,gt=0"`
	Ask  float32   `validate:"required,gt=0"`
	Time string    `validate:"required"`
}
