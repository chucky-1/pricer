package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/chucky-1/pricer/internal/config"
	"github.com/chucky-1/pricer/internal/grpc/server"
	"github.com/chucky-1/pricer/internal/model"
	"github.com/chucky-1/pricer/internal/repository"
	"github.com/chucky-1/pricer/protocol"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"fmt"
	"net"
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

	// Initial dependencies
	channels := model.Channels{
		Chan:   map[int]map[int]chan *model.Stock{},
		ChanID: map[int]int{},
		UserID: map[int]int{},
	}
	ch := make(chan *model.Stock)
	rep := repository.NewRepository(rdb, &channels, ch)

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
			log.Infof("%s is update, new cost is %f", stock.Title, stock.Price)
		}
	}
}
