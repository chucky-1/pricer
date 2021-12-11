// Package model has struct of essence
package model

// Stock contains fields that describe the shares of companies
type Stock struct {
	ID    int     `validate:"required"`
	Title string  `validate:"required"`
	Price float32 `validate:"required,gt=0"`
}

// Channels contains collections of channels and condition of channels
// Chan is map[Stock.ID]map[ChanID]
// ChanID is map[Stock.ID]
// UserID key is userID, value is ChanID
type Channels struct {
	Chan   map[int]map[int]chan *Stock
	ChanID map[int]int
	UserID map[string]int
}
