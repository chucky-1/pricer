// Package repository works with database
package repository

import (
	"github.com/chucky-1/pricer/internal/model"

	"errors"
	"sync"
)

// Repository contains map which keeps model of stocks
type Repository struct {
	stocks model.Stocks
	mu     sync.Mutex
}

// NewRepository is constructor
func NewRepository(stocks model.Stocks) *Repository {
	return &Repository{stocks: stocks}
}

// Add func creates or updates a database
func (r *Repository) Add(stock *model.Stock) error {
	r.mu.Lock()
	r.stocks[stock.ID] = stock
	r.mu.Unlock()
	return nil
}

// Get func takes an element from database
func (r *Repository) Get(id int) (*model.Stock, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	stock, ok := r.stocks[id]
	if !ok {
		return nil, errors.New("stock didn't find")
	}
	return stock, nil
}
