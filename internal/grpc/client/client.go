package main

import (
	"context"
	"github.com/chucky-1/pricer/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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

	listID := protocol.ListID{
		Id: []int32{},
	}
	listID.Id = append(listID.Id, int32(1))
	listID.Id = append(listID.Id, int32(2))
	listID.Id = append(listID.Id, int32(3))
	stream, err := client.Send(context.Background(), &listID)
	if err != nil {
		log.Error("Error for Send")
	}
	for {
		recv, err := stream.Recv()
		if err != nil {
			return
		}
		log.Infof("%s is update, new cost is %f", recv.Title, recv.Price)
	}
}
