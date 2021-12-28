// Package repository receives current prices and sends them to the channels
package repository

import (
	"github.com/chucky-1/pricer/internal/model"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"context"
	"strconv"
	"sync"
)

// Repository contains redis client and channels. All channels contain an activity flag
type Repository struct {
	rdb        *redis.Client
	muChannels sync.RWMutex
	channels   map[int32]map[string]chan *model.Price // map[price.ID]map[grpcID]chan
	muPrices   sync.RWMutex
	prices     map[int32]*model.Price // map[price.ID]*price
}

// NewRepository is constructor
func NewRepository(rdb *redis.Client, ch chan *model.Price) *Repository {
	rep := Repository{
		rdb:      rdb,
		channels: make(map[int32]map[string]chan *model.Price),
		prices:  make(map[int32]*model.Price),
	}
	go rep.listen(ch)
	return &rep
}

// Add is func initialization of stream
func (r *Repository) Add(priceID []int32, grpcID string, ch chan *model.Price) {
	r.muChannels.Lock()
	defer r.muChannels.Unlock()
	for _, id := range priceID {
		mapWithChan, ok := r.channels[id]
		if !ok {
			r.channels[id] = make(map[string]chan *model.Price)
			r.channels[id][grpcID] = ch
		} else {
			mapWithChan[grpcID] = ch
		}
	}
	go func() {
		r.sendPrimaryValues(priceID, ch)
	}()
}

// Del func unsubscribes from one or more prices. It doesn't close the channel!
func (r *Repository) Del(priceID []int32, grpcID string) {
	r.muChannels.Lock()
	defer r.muChannels.Unlock()
	for _, id := range priceID {
		mapWithChan, ok := r.channels[id]
		if !ok {
			continue
		}
		delete(mapWithChan, grpcID)
	}
}

// Close func closes the channel and delete it from map
func (r *Repository) Close(grpcID string) error {
	r.muChannels.Lock()
	defer r.muChannels.Unlock()
	var ch chan *model.Price
	for _, mapWithChan := range r.channels {
		c, ok := mapWithChan[grpcID]
		if ok {
			ch = c
			delete(mapWithChan, grpcID)
			break
		}
	}
	for _, mapWithChan := range r.channels {
		delete(mapWithChan, grpcID)
	}
	if ch != nil {
		close(ch)
	}
	return nil
}

// listen func listens redis stream and sends shares to the channel if it is active
func (r *Repository) listen(ch chan *model.Price) {
	var nextID = "$"
	for {
		entries, err := r.rdb.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{"stream", nextID},
			Count:   1,
			Block:   0,
		}).Result()
		if err != nil {
			log.Error(err)
			return
		}
		nextID = entries[0].Messages[0].ID
		m := entries[0].Messages[0].Values
		id, ok := m["ID"].(string)
		if !ok {
			log.Error("id is missing")
		}
		i, err := strconv.Atoi(id)
		if err != nil {
			log.Error(err)
		}
		bid, ok := m["Bid"].(string)
		if !ok {
			log.Error("price is missing")
		}
		b, err := strconv.ParseFloat(bid, 32)
		if err != nil {
			log.Error(err)
		}
		ask, ok := m["Ask"].(string)
		if !ok {
			log.Error("price is missing")
		}
		a, err := strconv.ParseFloat(ask, 32)
		if err != nil {
			log.Error(err)
		}
		price := model.Price{
			ID:   int32(i),
			Bid:  float32(b),
			Ask:  float32(a),
			Time: nextID[:len(nextID)-2],
		}

		ch <- &price

		r.update(&price)

		// Init map for the price
		r.muChannels.Lock()
		_, ok = r.channels[price.ID]
		if !ok {
			r.channels[price.ID] = make(map[string]chan *model.Price)
		}
		r.muChannels.Unlock()

		// Send the price in the channels
		go func() {
			r.muChannels.RLock()
			defer r.muChannels.RUnlock()
			channels, ok := r.channels[price.ID]
			if ok {
				for _, c := range channels {
					c <- &price
				}
			}
		}()
	}
}

// update func updates the latest price
func (r *Repository)update(price *model.Price) {
	r.muPrices.Lock()
	defer r.muPrices.Unlock()
	r.prices[price.ID] = price
}

// Send primary values sends prices into chan
func (r *Repository)sendPrimaryValues(priceID []int32, ch chan *model.Price) {
	r.muPrices.RLock()
	defer r.muPrices.RUnlock()
	for _, id := range priceID {
		ch <- &model.Price{
			ID: id,
			Bid:  r.prices[id].Bid,
			Ask:  r.prices[id].Ask,
			Time: r.prices[id].Time,
		}
	}
}
