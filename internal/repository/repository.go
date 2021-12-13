// Package repository receives current stock prices and sends them to active channels
package repository

import (
	"github.com/chucky-1/pricer/internal/model"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"context"
	"errors"
	"strconv"
	"sync"
)

// Repository contains redis client and channels. All channels contain an activity flag
type Repository struct {
	rdb    *redis.Client
	sub    *model.Subscribers
	memory *model.Memory
	mu     *sync.Mutex
}

// NewRepository is constructor
func NewRepository(rdb *redis.Client, sub *model.Subscribers, memory *model.Memory, mu *sync.Mutex, ch chan *model.Stock) *Repository {
	rep := Repository{rdb: rdb, sub: sub, memory: memory, mu: mu}
	go listen(rdb, memory, mu, ch)
	return &rep
}

// Send activates the stream, changing the flag to true
func (r *Repository) Send(list []int, userID string) (chan *model.Stock, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user := model.User{
		ID:     userID,
		Chan:   make(chan *model.Stock),
		Stocks: []int{},
	}
	for _, id := range list {
		user.Stocks = append(user.Stocks, id)

		s, ok := r.memory.Sub[id]
		if !ok {
			newSub := model.Subscribers{
				StockID: id,
				Users:   make(map[string]*model.User),
			}
			newSub.Users[user.ID] = &user
			r.memory.Sub[id] = &newSub
		} else {
			s.Users[user.ID] = &user
			r.memory.Sub[id] = s
		}
	}
	r.memory.User[userID] = &user
	return user.Chan, nil
}

// Close func closes the channel and delete it from map
func (r *Repository) Close(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.memory.User[userID]
	if !ok {
		return errors.New("chan didn't find")
	}

	for _, stockID := range user.Stocks {
		s, ok := r.memory.Sub[stockID]
		if !ok {
			continue
		}
		delete(s.Users, user.ID)
		r.memory.Sub[stockID] = s
	}
	close(user.Chan)
	return nil
}

// listen func listens redis stream and sends shares to the channel if it is active
func listen(rdb *redis.Client, memory *model.Memory, mu *sync.Mutex, ch chan *model.Stock) {
	for {
		entries, err := rdb.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{"stream", "$"},
			Count:   1,
			Block:   0,
		}).Result()
		if err != nil {
			log.Error(err)
			return
		}
		m := entries[0].Messages[0].Values
		id, ok := m["ID"].(string)
		if !ok {
			log.Error("id is missing")
		}
		i, err := strconv.Atoi(id)
		if err != nil {
			log.Error(err)
		}
		title, ok := m["Title"].(string)
		if !ok {
			log.Error("title is missing")
		}
		price, ok := m["Price"].(string)
		if !ok {
			log.Error("price is missing")
		}
		p, err := strconv.ParseFloat(price, 32)
		if err != nil {
			log.Error(err)
		}
		stock := model.Stock{
			ID:    i,
			Title: title,
			Price: float32(p),
		}
		ch <- &stock

		sub, ok := memory.Sub[stock.ID]
		if ok {
			for _, user := range sub.Users {
				user.Chan <- &stock
			}
		}
	}
}
