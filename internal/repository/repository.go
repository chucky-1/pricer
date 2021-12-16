// Package repository receives current stock prices and sends them to active channels
package repository

import (
	"github.com/chucky-1/pricer/internal/model"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"context"
	"strconv"
	"sync"
	"time"
)

// Repository contains redis client and channels. All channels contain an activity flag
type Repository struct {
	rdb   *redis.Client
	mu    *sync.Mutex
	stock map[int32]map[string]chan *model.Stock // map[Stock ID]grpc
}

// NewRepository is constructor
func NewRepository(rdb *redis.Client, ch chan *model.Stock) *Repository {
	mu := new(sync.Mutex)
	stock := make(map[int32]map[string]chan *model.Stock)
	rep := Repository{rdb: rdb, mu: mu, stock: stock}
	go listen(rdb, stock, mu, ch)
	return &rep
}

// Send activates the stream, changing the flag to true
func (r *Repository) Send(list []int32, grpcID string) (chan *model.Stock, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch := make(chan *model.Stock)
	for _, stockID := range list {
		grpc, ok := r.stock[stockID]
		if !ok {
			r.stock[stockID] = make(map[string]chan *model.Stock)
			grpc = r.stock[stockID]
		}
		grpc[grpcID] = ch
	}
	go func() {
		err := sendPrimaryValues(r.rdb, list, ch)
		if err != nil {
			log.Error(err)
		}
	}()
	return ch, nil
}

// Send primary values from the database on client connection
func sendPrimaryValues(rdb *redis.Client, stocks []int32, ch chan *model.Stock) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, id := range stocks {
		lastPrice, err := rdb.Get(ctx, strconv.Itoa(int(id))).Result()
		if err != nil {
			return err
		}
		value, err := strconv.ParseFloat(lastPrice, 32)
		if err != nil {
			return err
		}
		v := float32(value)
		stock := model.Stock{
			ID:    id,
			Title: "unknown",
			Price: v,
		}
		ch <- &stock
	}
	return nil
}

// Close func closes the channel and delete it from map
func (r *Repository) Close(grpcID string) error {
	var ch chan *model.Stock
	for _, stock := range r.stock {
		c, ok := stock[grpcID]
		if ok {
			ch = c
			delete(stock, grpcID)
			break
		}
	}
	for _, stock := range r.stock {
		delete(stock, grpcID)
	}
	close(ch)
	return nil
}

// listen func listens redis stream and sends shares to the channel if it is active
func listen(rdb *redis.Client, mm map[int32]map[string]chan *model.Stock, mu *sync.Mutex, ch chan *model.Stock) {
	var nextID = "$"
	for {
		entries, err := rdb.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{"stream", nextID},
			Count:   1,
			Block:   0,
		}).Result()
		if err != nil {
			log.Error(err)
			return
		}
		nextID = entries[0].Messages[0].ID
		t, err := getTimeFromID(nextID)
		if err != nil {
			log.Error(err)
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
			ID:     int32(i),
			Title:  title,
			Price:  float32(p),
			Update: t,
		}

		ch <- &stock

		err = update(rdb, id, stock.Price)
		if err != nil {
			log.Error(err)
		}

		// Send the price in the channels
		go func() {
			grpc, ok := mm[stock.ID]
			if ok {
				for _, c := range grpc {
					c <- &stock
				}
			}
		}()
	}
}

// update func updates the latest price in the database
func update(rdb *redis.Client, id string, price float32) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := rdb.Set(ctx, id, price, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func getTimeFromID(id string) (time.Time, error) {
	mkr, err := strconv.Atoi(id[:len(id)-2])
	if err != nil {
		return time.Time{}, err
	}
	t := time.Unix(int64(mkr)/1000, 0)
	return t, nil
}
