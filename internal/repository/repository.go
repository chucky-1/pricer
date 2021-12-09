// Package repository receives current stock prices and sends them to active channels
package repository

import (
	"github.com/chucky-1/pricer/internal/model"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"context"
	"errors"
	"strconv"
)

// Repository contains redis client and channels. All channels contain an activity flag
type Repository struct {
	rdb      *redis.Client
	channels *model.Channels
}

// NewRepository is constructor
func NewRepository(rdb *redis.Client, channels *model.Channels, ch chan *model.Stock) *Repository {
	rep := Repository{rdb: rdb, channels: channels}
	go listen(rdb, ch, channels)
	return &rep
}

// Send activates the stream, changing the flag to true
func (r *Repository) Send(id int) (chan *model.Stock, error) {
	ch, ok := r.channels.Collect[id]
	if !ok {
		return nil, errors.New("stock didn't find")
	}
	r.channels.Active[id] = true
	return ch, nil
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

		c, ok := channels.Collect[stock.ID]
		if !ok {
			channels.Collect[stock.ID] = make(chan *model.Stock)
			channels.Active[stock.ID] = false
		} else {
			if channels.Active[stock.ID] {
				c <- &stock
			}
		}
	}
}
