// Package model has struct of essence
package model

// Symbol contains fields that describe the shares of companies
// Update is time from redis ID
type Symbol struct {
	ID   int32   `validate:"required"`
	Bid  float32 `validate:"required,gt=0"`
	Ask  float32 `validate:"required,gt=0"`
	Time string  `validate:"required"`
}
