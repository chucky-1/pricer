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
	rdb      *redis.Client
	channels *model.Channels
	mu       sync.Mutex
}

// NewRepository is constructor
func NewRepository(rdb *redis.Client, channels *model.Channels, ch chan *model.Stock) *Repository {
	rep := Repository{rdb: rdb, channels: channels}
	go listen(rdb, ch, channels)
	return &rep
}

// Send activates the stream, changing the flag to true
func (r *Repository) Send(stockID int, userID string) (chan *model.Stock, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	mapChan, ok := r.channels.Chan[stockID]
	if !ok {
		return nil, errors.New("stock didn't find")
	}
	chanID := r.channels.ChanID[stockID]
	r.channels.ChanID[stockID]++
	ch := make(chan *model.Stock)
	mapChan[chanID] = ch
	r.channels.Chan[stockID] = mapChan
	r.channels.UserID[userID] = chanID
	return ch, nil
}

// Close func closes the channel and delete it from map
func (r *Repository) Close(stockID int, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	stock, ok := r.channels.Chan[stockID]
	if !ok {
		return errors.New("stock not found when closing")
	}
	chanID := r.channels.UserID[userID]
	ch := stock[chanID]

	close(ch)
	delete(r.channels.Chan[stockID], chanID)
	delete(r.channels.UserID, userID)
	return nil
}

// listen func listens redis stream and sends shares to the channel if it is active
func listen(rdb *redis.Client, ch chan *model.Stock, channels *model.Channels) {
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

		mapChan, ok := channels.Chan[stock.ID]
		if !ok {
			channels.Chan[stock.ID] = map[int]chan *model.Stock{}
			channels.ChanID[stock.ID] = 1
		} else {
			for _, chStock := range mapChan {
				chStock <- &stock
			}
		}
	}
}
