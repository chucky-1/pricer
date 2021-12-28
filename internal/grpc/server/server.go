// Package server implements the server side of grpc
package server

import (
	"github.com/chucky-1/pricer/internal/model"
	"github.com/chucky-1/pricer/internal/repository"
	"github.com/chucky-1/pricer/protocol"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Server contains methods of application on service side
type Server struct {
	protocol.UnimplementedPricesServer
	rep *repository.Repository
}

// NewServer is constructor
func NewServer(rep *repository.Repository) *Server {
	return &Server{rep: rep}
}

// Subscribe listens on the channel and sends data to the client
func (s *Server) Subscribe(stream protocol.Prices_SubscribeServer) error {
	grpcID := uuid.New().String()
	ch := make(chan *model.Symbol)

	go func() {
		for {
			symbol, ok := <-ch
			if !ok {
				break
			}
			err := stream.Send(&protocol.SubscribeResponse{
				SymbolId: symbol.ID,
				Bid: symbol.Bid,
				Ask: symbol.Ask,
			})
			if err != nil {
				log.Error(err)
			}
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			err := s.rep.Close(grpcID)
			if err != nil {
				log.Error(err)
			}
			return stream.Context().Err()
		default:
			recv, err := stream.Recv()
			if err != nil {
				log.Error(err)
				continue
			}
			switch {
			case recv.Action.String() == "ADD":
				s.rep.Add(recv.SymbolId, grpcID, ch)
			case recv.Action.String() == "DEL":
				s.rep.Del(recv.SymbolId, grpcID)
			}
		}
	}
}
