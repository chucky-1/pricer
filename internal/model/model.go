// Package model has struct of essence
package model

// Stock contains fields that describe the shares of companies
// Update is time from redis ID
type Stock struct {
	ID     int32   `validate:"required"`
	Title  string  `validate:"required"`
	Price  float32 `validate:"required,gt=0"`
	Update string  `validate:"required"`
}
