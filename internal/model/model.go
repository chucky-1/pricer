// Package model has struct of essence
package model

// Stock contains fields that describe the shares of companies
type Stock struct {
	ID    int     `validate:"required"`
	Title string  `validate:"required"`
	Price float32 `validate:"required,gt=0"`
}

// Stocks simulates a database
type Stocks map[int]*Stock
