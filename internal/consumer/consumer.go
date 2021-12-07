// Package consumer listens stream
package consumer

import (
	"github.com/chucky-1/pricer/internal/model"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"context"
	"strconv"
)

// Consumer contains redis client
type Consumer struct {
	rdb *redis.Client
}

// NewConsumer is constructor
func NewConsumer(rdb *redis.Client, ch chan *model.Stock) {
	rep := &Consumer{rdb: rdb}
	go rep.listen(ch)
}

// listen func listens redis stream
func (r *Consumer) listen(ch chan *model.Stock) {
	for {
		entries, err := r.rdb.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{"stream", "$"},
			Count:   1,
			Block:   0,
		}).Result()
		if err != nil {
			log.Error(err)
			return
		}
		id, _ := entries[0].Messages[0].Values["ID"].(string)
		i, _ := strconv.Atoi(id)
		title, _ := entries[0].Messages[0].Values["Title"].(string)
		price, _ := entries[0].Messages[0].Values["Price"].(string)
		p, _ := strconv.ParseFloat(price, 32)
		stock := model.Stock{
			ID:    i,
			Title: title,
			Price: float32(p),
		}
		ch <- &stock
	}
}
