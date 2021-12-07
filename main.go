package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/chucky-1/pricer/internal/config"
	"github.com/chucky-1/pricer/internal/consumer"
	"github.com/chucky-1/pricer/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"fmt"
)

func main() {
	// Configuration
	cfg := new(config.Config)
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("%v", err)
	}

	// Validator
	validate := validator.New()

	// Redis
	hostAndPort := fmt.Sprint(cfg.Host, ":", cfg.Port)
	rdb := redis.NewClient(&redis.Options{Addr: hostAndPort})

	ch := make(chan *model.Stock)
	consumer.NewConsumer(rdb, ch)

	// Initial stocks
	stocks := make(map[int]*model.Stock)

	// Business logic
	for {
		stock := <-ch
		err := validate.Struct(stock)
		if err != nil {
			log.Error("Struct isn't valid")
		} else {
			stocks[stock.ID] = stock
			log.Infof("%s is update, new cost is %f", stock.Title, stock.Price)
		}
	}
}
