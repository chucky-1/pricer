package main

import (
	"github.com/chucky-1/pricer/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"context"
)

// This is a test implementation of the client. This file will be removed.
func main() {
	conn, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)
	client := protocol.NewPricerClient(conn)

	id := protocol.Id{Id: int32(4)}

	stream, err := client.Send(context.Background(), &id)
	if err != nil {
		log.Error("Error for Send")
	}
	for {
		recv, err := stream.Recv()
		if err != nil {
			return
		}
		log.Infof("%s, %f", recv.Title, recv.Price)
	}
}
