package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/chucky-1/pricer/grpc/server"
	"github.com/chucky-1/pricer/internal/config"
	"github.com/chucky-1/pricer/internal/model"
	"github.com/chucky-1/pricer/internal/repository"
	"github.com/chucky-1/pricer/protocol"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"context"
	"fmt"
	"net"
	"strconv"
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

	rep := repository.NewRepository(model.Stocks{})
	ch := make(chan *model.Stock)

	// Listen to the redis stream
	go func() {
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
		}
	}()

	// Grpc
	go func() {
		hostAndPort = fmt.Sprint(cfg.HostGrpc, ":", cfg.PortGrpc)
		lis, err := net.Listen("tcp", hostAndPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		protocol.RegisterPricerServer(s, server.NewServer(rep))
		log.Infof("server listening at %v", lis.Addr())
		if err = s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Business logic
	for {
		stock := <-ch
		err := validate.Struct(stock)
		if err != nil {
			log.Error("Struct isn't valid")
		} else {
			err = rep.Add(stock)
			if err != nil {
				log.Error("Failed to add")
			} else {
				log.Infof("%s is update, new cost is %f", stock.Title, stock.Price)
			}
		}
	}
}
