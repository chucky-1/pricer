// Package model has struct of essence
package model

import "time"

// Stock contains fields that describe the shares of companies
type Stock struct {
	ID     int       `validate:"required"`
	Title  string    `validate:"required"`
	Price  float32   `validate:"required,gt=0"`
	Update time.Time `validate:"required"`
}

// User struck
// Stocks is list of ID of Stock to which the user is subscribed
type User struct {
	ID     string
	Chan   chan *Stock
	Stocks []int
}

// Subscribers struct. It shows which users are subscribed to each stock
// Users is map[User.ID]
type Subscribers struct {
	StockID int
	Users   map[string]*User
}

// Memory struct
// User is map[User.ID]
// Sub is map[Subscribers.StockID]
type Memory struct {
	User map[string]*User
	Sub  map[int]*Subscribers
}
