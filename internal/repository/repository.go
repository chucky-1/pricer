// Package repository receives current symbol prices and sends them to the channels
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
	channels   map[int32]map[string]chan *model.Symbol // map[symbol.ID]map[grpcID]chan
	muSymbols  sync.RWMutex
	symbols    map[int32]*model.Symbol // map[Symbol.ID]*symbol
}

// NewRepository is constructor
func NewRepository(rdb *redis.Client, ch chan *model.Symbol) *Repository {
	rep := Repository{
		rdb:      rdb,
		channels: make(map[int32]map[string]chan *model.Symbol),
		symbols:  make(map[int32]*model.Symbol),
	}
	go rep.listen(ch)
	return &rep
}

// Add is func initialization of stream
func (r *Repository) Add(list []int32, grpcID string, ch chan *model.Symbol) {
	r.muChannels.Lock()
	defer r.muChannels.Unlock()
	for _, symbolID := range list {
		mapWithChan, ok := r.channels[symbolID]
		if !ok {
			r.channels[symbolID] = make(map[string]chan *model.Symbol)
			r.channels[symbolID][grpcID] = ch
		} else {
			mapWithChan[grpcID] = ch
		}
	}
	go func() {
		r.sendPrimaryValues(list, ch)
	}()
}

// Del func unsubscribes from one or more symbols. It doesn't close the channel!
func (r *Repository) Del(list []int32, grpcID string) {
	r.muChannels.Lock()
	defer r.muChannels.Unlock()
	for _, symbolID := range list {
		mapWithChan, ok := r.channels[symbolID]
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
	var ch chan *model.Symbol
	for _, symbol := range r.channels {
		c, ok := symbol[grpcID]
		if ok {
			ch = c
			delete(symbol, grpcID)
			break
		}
	}
	for _, symbol := range r.channels {
		delete(symbol, grpcID)
	}
	if ch != nil {
		close(ch)
	}
	return nil
}

// listen func listens redis stream and sends shares to the channel if it is active
func (r *Repository) listen(ch chan *model.Symbol) {
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
		symbol := model.Symbol{
			ID:   int32(i),
			Bid:  float32(b),
			Ask:  float32(a),
			Time: nextID[:len(nextID)-2],
		}

		ch <- &symbol

		r.update(&symbol)

		// Init map for the symbols
		r.muChannels.Lock()
		_, ok = r.channels[symbol.ID]
		if !ok {
			r.channels[symbol.ID] = make(map[string]chan *model.Symbol)
		}
		r.muChannels.Unlock()

		// Send the price in the channels
		go func() {
			r.muChannels.RLock()
			defer r.muChannels.RUnlock()
			channels, ok := r.channels[symbol.ID]
			if ok {
				for _, c := range channels {
					c <- &symbol
				}
			}
		}()
	}
}

// update func updates the latest price
func (r *Repository)update(symbol *model.Symbol) {
	r.muSymbols.Lock()
	defer r.muSymbols.Unlock()
	r.symbols[symbol.ID] = symbol
}

// Send primary values sends symbols into chan
func (r *Repository)sendPrimaryValues(symbols []int32, ch chan *model.Symbol) {
	r.muSymbols.RLock()
	defer r.muSymbols.RUnlock()
	for _, id := range symbols {
		ch <- &model.Symbol{
			ID: id,
			Bid:  r.symbols[id].Bid,
			Ask:  r.symbols[id].Ask,
			Time: r.symbols[id].Time,
		}
	}
}
